package grpc

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Forest33/warthog/business/entity"
)

func (c *Client) Query(method *entity.Method, data map[string]interface{}, requestMetadata []string) (qResp *entity.QueryResponse, err error) {
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

	ms := dynamic.NewMessage(method.Descriptor.GetInputType())

	for _, f := range method.Input {
		if _, ok := data[f.Fqn]; !ok {
			continue
		}

		arg, err := c.getArgument(f, data[f.Fqn])
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

	var ctx context.Context
	ctx, c.cancelQuery = context.WithTimeout(c.ctx, time.Second*time.Duration(c.cfg.QueryTimeout))
	defer func() {
		c.cancelQueryMux.Lock()
		defer c.cancelQueryMux.Unlock()

		if c.cancelQuery != nil {
			c.cancelQuery()
			c.cancelQuery = nil
		}
	}()

	ctx = addMetadata(ctx, requestMetadata)

	var header, trailer metadata.MD

	var (
		resp     proto.Message
		jsonResp string
	)

	stub := grpcdynamic.NewStub(c.conn)
	startTime := time.Now()

	switch {
	case method.Descriptor.IsClientStreaming() && method.Descriptor.IsServerStreaming():
		stream, err := stub.InvokeRpcBidiStream(ctx, method.Descriptor, grpc.Header(&header), grpc.Trailer(&trailer))
		if err != nil {
			return nil, err
		}
		if err = stream.SendMsg(ms); err != nil {
			return nil, err
		}
		if header, err = stream.Header(); err != nil {
			return nil, err
		}
		if resp, err = stream.RecvMsg(); err != nil {
			return nil, err
		}
		trailer = stream.Trailer()
	case method.Descriptor.IsClientStreaming():
		stream, err := stub.InvokeRpcClientStream(ctx, method.Descriptor, grpc.Header(&header), grpc.Trailer(&trailer))
		if err != nil {
			return nil, err
		}
		if err = stream.SendMsg(ms); err != nil {
			return nil, err
		}
		resp, err = stream.CloseAndReceive()
	case method.Descriptor.IsServerStreaming():
		stream, err := stub.InvokeRpcServerStream(ctx, method.Descriptor, ms, grpc.Header(&header), grpc.Trailer(&trailer))
		if err != nil {
			return nil, err
		}
		if header, err = stream.Header(); err != nil {
			return nil, err
		}
		if resp, err = stream.RecvMsg(); err != nil {
			return nil, err
		}
		trailer = stream.Trailer()
	default:
		resp, err = stub.InvokeRpc(ctx, method.Descriptor, ms, grpc.Header(&header), grpc.Trailer(&trailer))
	}

	spentTime := time.Since(startTime)
	if err != nil {
		return nil, err
	}

	if resp != nil {
		jsonResp, err = c.getResponse(resp)
		if err != nil {
			return nil, err
		}
	}

	return &entity.QueryResponse{
		JsonString: jsonResp,
		SpentTime:  spentTime.String(),
		Header:     header,
		Trailer:    trailer,
	}, nil
}

func (c *Client) CancelQuery() {
	c.cancelQueryMux.Lock()
	defer c.cancelQueryMux.Unlock()

	if c.cancelQuery != nil {
		c.cancelQuery()
		c.cancelQuery = nil
	}
}

func (c *Client) getResponse(m proto.Message) (string, error) {
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
		return "", fmt.Errorf("unknown response type: %s", reflect.TypeOf(m).Elem().String())
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
		resp = entity.GetBytes(field, data)
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
						Fqn:  field.Fqn,
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
			//md := field.Descriptor.GetFile().FindMessage(field.Message.Type)
			md := field.Message.Descriptor
			if md == nil {
				return nil, fmt.Errorf("failed to find message type %s", field.Message.Type)
			}
			if !field.Repeated {
				ms := dynamic.NewMessage(md)
				for _, mf := range field.Message.Fields {
					v, err := c.getArgument(mf, data.(map[string]interface{})[mf.Fqn])
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
						v, err := c.getArgument(mf, d.(map[string]interface{})[mf.Fqn])
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
