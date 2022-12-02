// Package grpc provides basic gRPC functions.
package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
	"runtime/debug"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/forest33/warthog/business/entity"
)

func (c *Client) createMessage(method *entity.Method, data map[string]interface{}, metadata []string) (ms *dynamic.Message, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
			err = fmt.Errorf("%v<br><pre>%s</pre>", err, debug.Stack())
		}
	}()

	ms = dynamic.NewMessage(method.Descriptor.GetInputType())

	for _, f := range method.Input {
		if _, ok := data[getProtoFQN(f)]; !ok {
			continue
		}

		arg, err := c.getArgument(f, data[getProtoFQN(f)])
		if err != nil {
			return nil, err
		}

		if arg == nil {
			continue
		}

		if err := ms.TrySetField(f.Descriptor, arg); err != nil {
			return nil, err
		}
	}

	if *c.cfg.RequestTimeout > 0 && method.Type == entity.MethodTypeUnary {
		c.queryCtx, c.queryCancel = context.WithTimeout(c.ctx, time.Second*time.Duration(*c.cfg.RequestTimeout))
	} else {
		c.queryCtx, c.queryCancel = context.WithCancel(c.ctx)
	}

	c.queryCtx = addMetadata(c.queryCtx, metadata)

	return
}

// Query executes a gRPC request
func (c *Client) Query(method *entity.Method, data map[string]interface{}, metadata []string) {
	ms, err := c.createMessage(method, data, metadata)
	if err != nil {
		c.responseError(err, "")
		return
	}

	var isNew bool

	switch method.Type {
	case entity.MethodTypeUnary:
		c.unary(method, ms)
	case entity.MethodTypeClientStream:
		isNew, err = c.clientStream(method)
	case entity.MethodTypeServerStream:
		c.serverStream(method, ms)
	case entity.MethodTypeBidiStream:
		c.bidiStream(method)
	}

	if err != nil {
		return
	}

	if isNew || method.Type == entity.MethodTypeClientStream || method.Type == entity.MethodTypeBidiStream {
		c.request(ms)
	}
}

// CancelQuery aborting a running gRPC request
func (c *Client) CancelQuery() {
	if c.queryCancel == nil {
		return
	}
	c.queryCancel()
}

// CloseStream stops a running gRPC stream
func (c *Client) CloseStream() {
	c.closeStreamCh <- struct{}{}
}

// GetResponseChannel returns response channel
func (c *Client) GetResponseChannel() chan *entity.QueryResponse {
	return c.responseCh
}

// GetSentCounter returns sent messages counter
func (c *Client) GetSentCounter() uint {
	return c.sentMessages
}

func (c *Client) unary(method *entity.Method, ms *dynamic.Message) {
	var (
		header  metadata.MD
		trailer metadata.MD
	)

	stub := grpcdynamic.NewStub(c.conn)
	c.queryStartTime = time.Now()
	resp, err := stub.InvokeRpc(c.queryCtx, method.Descriptor, ms, grpc.Header(&header), grpc.Trailer(&trailer))
	c.response(resp, header, trailer, err)

	c.sentMessages = 0
	c.receivedMessaged = 0
}

func (c *Client) clientStream(method *entity.Method) (bool, error) {
	var (
		header  metadata.MD
		trailer metadata.MD
		stream  *grpcdynamic.ClientStream
		err     error
	)

	isNew := c.startRequest()
	if isNew {
		stub := grpcdynamic.NewStub(c.conn)
		stream, err = stub.InvokeRpcClientStream(c.queryCtx, method.Descriptor, grpc.Header(&header), grpc.Trailer(&trailer))
		if err != nil {
			c.responseError(err, "")
			return isNew, err
		}

		go func() {
			defer func() {
				c.stopRequest()
			}()

			c.queryStartTime = time.Now()

			for {
				select {
				case <-c.queryCtx.Done():
					c.response(nil, header, trailer, status.FromContextError(context.Canceled).Err())
					_, _ = stream.CloseAndReceive()
					c.log.Debug().Msg("stream canceled")
					return
				case <-c.closeStreamCh:
					data, err := stream.CloseAndReceive()
					c.response(data, header, trailer, err)
					c.log.Debug().Msg("close & receive stream")
					return
				case ms := <-c.requestCh:
					if err := stream.SendMsg(ms); err != nil {
						c.response(nil, header, trailer, err)
						return
					}
				}
			}
		}()
	}

	return isNew, nil
}

func (c *Client) serverStream(method *entity.Method, ms *dynamic.Message) {
	var (
		header  metadata.MD
		trailer metadata.MD
		isBreak bool
	)

	if c.startRequest() {
		stub := grpcdynamic.NewStub(c.conn)
		c.queryStartTime = time.Now()
		stream, err := stub.InvokeRpcServerStream(c.queryCtx, method.Descriptor, ms, grpc.Header(&header), grpc.Trailer(&trailer))
		if err != nil {
			c.responseError(err, "")
			return
		}

		go func() {
			defer func() {
				c.stopRequest()
			}()

			for !isBreak {
				data, err := stream.RecvMsg()
				if err == io.EOF {
					isBreak = true
				} else if status.Code(err) == codes.Canceled {
					c.responseError(err, time.Since(c.queryStartTime).String())
					return
				} else if err != nil {
					c.responseError(err, "")
					c.log.Error().Msgf("failed to receive message: %v", err)
					return
				}
				header, hErr := stream.Header()
				if hErr != nil {
					c.log.Error().Msgf("failed to get message header: %v", err)
				}
				trailer = stream.Trailer()
				c.response(data, header, trailer, nil)
			}
		}()
	}
}

func (c *Client) bidiStream(method *entity.Method) {
	var (
		header  metadata.MD
		trailer metadata.MD
		isBreak bool
	)

	if c.startRequest() {
		stub := grpcdynamic.NewStub(c.conn)
		c.queryStartTime = time.Now()
		stream, err := stub.InvokeRpcBidiStream(c.queryCtx, method.Descriptor, grpc.Header(&header), grpc.Trailer(&trailer))
		if err != nil {
			c.responseError(err, "")
			return
		}

		go func() {
			defer func() {
				c.stopRequest()
			}()

			for !isBreak {
				data, err := stream.RecvMsg()
				if err == io.EOF || isStreamEOF(err) {
					isBreak = true
				} else if status.Code(err) == codes.Canceled {
					c.responseError(err, time.Since(c.queryStartTime).String())
					return
				} else if err != nil {
					c.responseError(err, time.Since(c.queryStartTime).String())
					c.log.Error().Msgf("failed to receive message: %v", err)
					return
				}
				header, hErr := stream.Header()
				if hErr != nil {
					c.log.Error().Msgf("failed to get message header: %v", err)
				}
				trailer = stream.Trailer()
				c.response(data, header, trailer, err)
			}
		}()

		go func() {
			for {
				select {
				case <-c.queryCtx.Done():
					c.response(nil, header, trailer, status.FromContextError(context.Canceled).Err())
					_ = stream.CloseSend()
					return
				case <-c.closeStreamCh:
					if err := stream.CloseSend(); err != nil {
						c.log.Error().Msgf("failed to close stream: %v", err)
					}
					c.log.Debug().Msg("close & send stream")
					return
				case ms := <-c.requestCh:
					if ms == nil {
						continue
					}
					if err := stream.SendMsg(ms); err != nil {
						c.response(nil, header, trailer, err)
						return
					}
				}
			}
		}()
	}
}

func (c *Client) getResponse(m proto.Message) (string, error) {
	if m == nil {
		return "", nil
	}

	switch t := m.(type) {
	case *dynamic.Message:
		buf, err := t.MarshalJSONPB(&jsonpb.Marshaler{Indent: "  ", OrigName: true})
		if err != nil {
			return "", err
		}
		return string(buf), nil
	case *emptypb.Empty:
		return "google.protobuf.Empty", nil
	default:
		err := fmt.Errorf("unknown response type: %s", reflect.TypeOf(m).Elem().String())
		return err.Error(), err
	}
}

func (c *Client) getArgument(field *entity.Field, data interface{}) (interface{}, error) {
	if isEmpty(data) {
		return nil, nil
	}

	var (
		resp interface{}
		err  error
	)

	switch field.Type {
	case entity.TypeString:
		resp = entity.GetString(field, data)
	case entity.TypeBytes:
		resp, err = entity.GetBytes(field, data)
	case entity.TypeInt32, entity.TypeSInt32, entity.TypeSFixed32:
		resp, err = entity.GetInt32(field, data)
	case entity.TypeInt64, entity.TypeSInt64, entity.TypeSFixed64:
		resp, err = entity.GetInt64(field, data)
	case entity.TypeUInt32, entity.TypeFixed32:
		resp, err = entity.GetUInt32(field, data)
	case entity.TypeUInt64, entity.TypeFixed64:
		resp, err = entity.GetUInt64(field, data)
	case entity.TypeDouble:
		resp, err = entity.GetFloat64(field, data)
	case entity.TypeFloat:
		resp, err = entity.GetFloat32(field, data)
	case entity.TypeBool:
		resp = entity.GetBool(field, data)
	case entity.TypeEnum:
		resp, err = entity.GetInt32(field, data)
	case entity.TypeMessage:
		if field.Map != nil {
			if field.Map.ProtoValueType == entity.TypeMessage {
				obj := make(map[interface{}]interface{}, len(data.(map[string]interface{})))
				for k, v := range data.(map[string]interface{}) {
					key, err := c.getArgument(&entity.Field{Type: field.Map.KeyType}, k)
					if err != nil {
						return nil, err
					}
					if key == nil {
						continue
					}
					val, err := c.getArgument(&entity.Field{
						FQN:  field.FQN,
						Type: field.Map.ProtoValueType,
						Message: &entity.Message{
							Type:       field.Map.ValueTypeFqn,
							Fields:     field.Map.Fields,
							Descriptor: field.Map.ValueDescriptor,
						},
						Descriptor: field.Descriptor,
					}, v)
					if err != nil {
						return nil, err
					}
					if val != nil {
						obj[key] = val
					}
				}
				resp = obj
			} else {
				obj := make(map[interface{}]interface{}, len(data.(map[string]interface{})))
				for k, v := range data.(map[string]interface{}) {
					key, err := c.getArgument(&entity.Field{Type: field.Map.KeyType}, k)
					if err != nil {
						return nil, err
					}
					if key == nil {
						continue
					}
					val, err := c.getArgument(&entity.Field{Type: field.Map.ValueType}, v)
					if err != nil {
						return nil, err
					}
					if val != nil {
						obj[key] = val
					}
				}
				resp = obj
			}
		} else {
			md := field.Message.Descriptor
			if md == nil {
				return nil, fmt.Errorf("failed to find message type %s", field.Message.Type)
			}
			if !field.Repeated {
				ms := dynamic.NewMessage(md)
				for _, mf := range field.Message.Fields {
					v, err := c.getArgument(mf, data.(map[string]interface{})[getProtoFQN(mf)])
					if err != nil {
						return nil, err
					}
					if v != nil {
						if err := ms.TrySetField(mf.Descriptor, v); err != nil {
							return nil, err
						}
					}
				}
				resp = ms
			} else {
				messages := make([]*dynamic.Message, 0, len(data.([]interface{})))
				for _, d := range data.([]interface{}) {
					ms := dynamic.NewMessage(md)
					var hasFields bool
					for _, mf := range field.Message.Fields {
						v, err := c.getArgument(mf, d.(map[string]interface{})[getProtoFQN(mf)])
						if err != nil {
							return nil, err
						}
						if v != nil {
							if err := ms.TrySetField(mf.Descriptor, v); err != nil {
								return nil, err
							}
							hasFields = true
						}
					}
					if hasFields {
						messages = append(messages, ms)
					}
				}
				resp = messages
			}
		}
	default:
		return nil, fmt.Errorf("unknown data type %s", field.Type)
	}

	return resp, err
}

func isEmpty(v interface{}) bool {
	if v == nil {
		return true
	}

	switch t := v.(type) {
	case string:
		return len(t) == 0
	case []string:
		return len(t) == 0
	case []interface{}:
		if len(t) == 1 {
			return t[0] == nil || (reflect.TypeOf(t[0]).Kind() == reflect.String && len(t[0].(string)) == 0)
		}
		return len(t) == 0
	}

	return false
}

func addMetadata(ctx context.Context, data []string) context.Context {
	if len(data) == 0 {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, data...)
}

func getProtoFQN(f *entity.Field) string {
	if f.ProtoFQN != "" {
		return f.ProtoFQN
	}
	return f.FQN
}

func (c *Client) request(ms *dynamic.Message) {
	c.sentMessages++
	c.requestCh <- ms
}

func (c *Client) response(data proto.Message, header metadata.MD, trailer metadata.MD, err error) {
	spent := time.Since(c.queryStartTime).String()

	if err != nil {
		c.responseError(err, spent)
		return
	}

	if data != nil {
		c.receivedMessaged++
	}

	resp := &entity.QueryResponse{
		Time:      time.Now().Format("15:04:05.99"),
		SpentTime: spent,
		Header:    header,
		Trailer:   trailer,
		Error:     toError(err),
		Sent:      c.sentMessages,
		Received:  c.receivedMessaged,
	}

	jsonResp, err := c.getResponse(data)
	resp.JsonString = jsonResp

	c.responseCh <- resp
}

func (c *Client) responseError(err error, spent string) {
	c.responseCh <- &entity.QueryResponse{
		Error:     toError(err),
		SpentTime: spent,
	}
}

func (c *Client) startRequest() bool {
	if c.requestCh != nil {
		return false
	}

	c.requestCh = make(chan *dynamic.Message, requestChanCapacity)
	c.sentMessages = 0
	c.receivedMessaged = 0

	return true
}

func (c *Client) stopRequest() {
	close(c.requestCh)
	c.requestCh = nil
	c.sentMessages = 0
	c.receivedMessaged = 0
}

func toError(err error) *entity.Error {
	if err == nil {
		return nil
	}
	return &entity.Error{
		Code:            uint32(status.Code(err)),
		CodeDescription: status.Code(err).String(),
		Message:         status.Convert(err).Message(),
	}
}

func isStreamEOF(err error) bool {
	if err == nil {
		return false
	}
	if s, ok := status.FromError(err); ok {
		return s.Message() == "EOF"
	}
	return false
}
