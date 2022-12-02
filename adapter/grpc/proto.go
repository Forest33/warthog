// Package grpc provides basic gRPC functions.
package grpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc/status"

	"github.com/forest33/warthog/business/entity"
)

// AddProtobuf adds protobuf files
func (c *Client) AddProtobuf(path ...string) {
	c.protoPath = path
}

// AddImport adds import paths
func (c *Client) AddImport(path ...string) {
	c.importPath = path
}

// LoadFromProtobuf loads services from protobuf
func (c *Client) LoadFromProtobuf() ([]*entity.Service, []*entity.ProtobufError, *entity.ProtobufError) {
	if len(c.protoPath) == 0 {
		return nil, nil, &entity.ProtobufError{Err: fmt.Errorf("empty path to protobuf's")}
	}

	var (
		protoWarn []*entity.ProtobufError
		protoErr  *entity.ProtobufError
	)

	warningReporter := func(err protoparse.ErrorWithPos) {
		if protoWarn == nil {
			protoWarn = make([]*entity.ProtobufError, 0, 1)
		}
		protoWarn = append(protoWarn, &entity.ProtobufError{
			Pos:     err.GetPosition(),
			Warning: err.Unwrap().Error(),
		})
	}

	errorReporter := func(err protoparse.ErrorWithPos) error {
		protoErr = &entity.ProtobufError{
			Code:            uint32(status.Code(err)),
			CodeDescription: status.Code(err).String(),
			Pos:             err.GetPosition(),
			Err:             err.Unwrap(),
		}
		return err
	}

	parser := protoparse.Parser{
		ImportPaths:     c.importPath,
		WarningReporter: warningReporter,
		ErrorReporter:   errorReporter,
	}

	services := make([]*entity.Service, 0, len(c.protoPath))

	for _, p := range c.protoPath {
		fd, err := parser.ParseFiles(p)
		if err != nil {
			return nil, nil, protoErr
		}
		if len(fd) != 1 {
			return nil, nil, &entity.ProtobufError{Err: fmt.Errorf("wrong parse result")}
		}

		for _, sd := range fd[0].GetServices() {
			methods := c.getMethods(sd)
			services = append(services, &entity.Service{
				Name:    sd.GetFullyQualifiedName(),
				Methods: methods,
			})
		}
	}

	c.sortServicesByName(services)

	return services, protoWarn, nil
}

// LoadFromReflection loads services using reflection
func (c *Client) LoadFromReflection() ([]*entity.Service, error) {
	ctx, cancel := context.WithTimeout(c.ctx, time.Second*time.Duration(*c.cfg.ConnectTimeout))
	defer cancel()

	client := grpcreflect.NewClientAuto(ctx, c.conn)
	list, err := client.ListServices()
	if err != nil {
		return nil, err
	}

	services := make([]*entity.Service, 0, len(list))
	for _, name := range list {
		sd, err := client.ResolveService(name)
		if err != nil {
			return nil, err
		}
		methods := c.getMethods(sd)
		services = append(services,
			&entity.Service{
				Name:    sd.GetFullyQualifiedName(),
				Methods: methods,
			},
		)
	}

	c.sortServicesByName(services)

	return services, nil
}

func (c *Client) getMethods(sd *desc.ServiceDescriptor) []*entity.Method {
	methodDesc := sd.GetMethods()
	methods := make([]*entity.Method, 0, len(methodDesc))

	for _, md := range methodDesc {
		methods = append(methods, &entity.Method{
			Name:       md.GetName(),
			Type:       getMethodType(md),
			Descriptor: md,
			Input:      c.getFields(md.GetInputType().GetFields(), nil, nil),
		})
	}

	if *c.cfg.SortMethodsByName {
		c.sortMethodsByName(methods)
	}

	return methods
}

func getMethodType(md *desc.MethodDescriptor) string {
	switch {
	case md.IsClientStreaming() && md.IsServerStreaming():
		return entity.MethodTypeBidiStream
	case md.IsClientStreaming():
		return entity.MethodTypeClientStream
	case md.IsServerStreaming():
		return entity.MethodTypeServerStream
	}
	return entity.MethodTypeUnary
}

func (c *Client) getFields(fd []*desc.FieldDescriptor, parent *desc.FieldDescriptor, fqn map[string]int) []*entity.Field {
	fields := make([]*entity.Field, 0, len(fd))

	for _, f := range fd {
		fieldType := f.GetType().String()

		switch fieldType {
		case entity.ProtoTypeEnum:
			values := f.GetEnumType().GetValues()
			enumValues := make([]*entity.EnumValue, 0, len(values))
			for _, v := range values {
				enumValues = append(enumValues, &entity.EnumValue{
					Name:   v.GetName(),
					Number: v.GetNumber(),
				})
			}
			fields = append(fields, &entity.Field{
				Descriptor: f,
				FQN:        f.GetFullyQualifiedName(),
				Name:       f.GetName(),
				Type:       getTypeName(fieldType),
				ParentType: getParentTypeName(parent),
				Repeated:   f.IsRepeated(),
				OneOf:      getOneOf(f),
				Enum: &entity.Enum{
					ValueType: getTypeName(f.GetEnumType().GetName()),
					Values:    enumValues,
				},
			})
		case entity.ProtoTypeMessage:
			if f.IsMap() {
				keyType, valueType, valueTypeName := getMapFieldTypeName(f)
				loopFQN, protoFQN := getFQN(f, fqn)
				mapField := &entity.Field{
					Descriptor: f,
					FQN:        loopFQN,
					ProtoFQN:   protoFQN,
					Name:       f.GetName(),
					Type:       getTypeName(fieldType),
					ParentType: getParentTypeName(parent),
					Repeated:   f.IsRepeated(),
					OneOf:      getOneOf(f),
					Map: &entity.Map{
						KeyType:        getTypeName(keyType),
						ValueType:      getTypeName(valueTypeName),
						ProtoValueType: getTypeName(f.GetMapValueType().GetType().String()),
					},
				}
				if f.GetMapValueType().GetMessageType() != nil {
					mapField.Map.ValueTypeFqn = f.GetMapValueType().GetMessageType().GetFullyQualifiedName()
				}
				if valueType == entity.ProtoTypeMessage {
					messageFields, fqn := c.getMessageFields(f.GetMapValueType().GetMessageType().GetFields(), fqn)
					if len(messageFields) == 0 {
						break
					}

					mapField.Map.ValueDescriptor = f.GetMapValueType().GetMessageType()
					mapField.Map.Fields = c.getFields(messageFields, f, fqn)
				} else if valueType == entity.ProtoTypeEnum {
					mapField.Map.Fields = c.getFields([]*desc.FieldDescriptor{f.GetMapValueType()}, f, fqn)
				}
				fields = append(fields, mapField)
			} else {
				messageFields, fqn := c.getMessageFields(f.GetMessageType().GetFields(), fqn)
				if len(messageFields) == 0 {
					break
				}

				loopFQN, protoFQN := getFQN(f, fqn)
				msgField := &entity.Field{
					Descriptor: f,
					FQN:        loopFQN,
					ProtoFQN:   protoFQN,
					Name:       f.GetName(),
					Type:       getTypeName(fieldType),
					ParentType: getParentTypeName(parent),
					Repeated:   f.IsRepeated(),
					OneOf:      getOneOf(f),
					Message: &entity.Message{
						Name:       f.GetMessageType().GetName(),
						Type:       f.GetMessageType().GetFullyQualifiedName(),
						Fields:     c.getFields(messageFields, f, fqn),
						Descriptor: f.GetMessageType(),
					},
				}
				fields = append(fields, msgField)
			}
		default:
			fields = append(fields, &entity.Field{
				Descriptor: f,
				FQN:        f.GetFullyQualifiedName(),
				Name:       f.GetName(),
				Type:       getTypeName(fieldType),
				ParentType: getParentTypeName(parent),
				Repeated:   f.IsRepeated(),
				OneOf:      getOneOf(f),
			})
		}
	}

	return fields
}

func getMapFieldTypeName(f *desc.FieldDescriptor) (string, string, string) {
	var valueTypeName string
	keyType := f.GetMapKeyType().GetType().String()
	valueType := f.GetMapValueType().GetType().String()
	if valueType == entity.ProtoTypeMessage {
		if mt := f.GetMapValueType().GetMessageType(); mt != nil {
			valueTypeName = mt.GetName()
		}
	} else if valueType == entity.ProtoTypeEnum {
		valueTypeName = f.GetMapValueType().GetEnumType().GetName()
	} else {
		valueTypeName = valueType
	}
	return keyType, valueType, valueTypeName
}

func getTypeName(t string) string {
	if strings.HasPrefix(t, "TYPE_") {
		return strings.ToLower(strings.ReplaceAll(t, "TYPE_", ""))
	}
	return t
}

func getParentTypeName(f *desc.FieldDescriptor) string {
	if f == nil {
		return ""
	}
	return getTypeName(f.GetType().String())
}

func getOneOf(f *desc.FieldDescriptor) *entity.OneOf {
	oneOf := f.GetOneOf()
	if oneOf == nil {
		return nil
	}

	return &entity.OneOf{
		Fqn:  oneOf.GetFullyQualifiedName(),
		Name: oneOf.GetName(),
	}
}

func (c *Client) getMessageFields(fields []*desc.FieldDescriptor, fqn map[string]int) ([]*desc.FieldDescriptor, map[string]int) {
	messageFields := make([]*desc.FieldDescriptor, 0, len(fields))
	if fqn == nil {
		fqn = make(map[string]int, len(fields))
	}

	for _, mf := range fields {
		name := mf.GetFullyQualifiedName()
		if count, ok := fqn[name]; ok && count >= *c.cfg.MaxLoopDepth {
			continue
		}
		fqn[name]++
		messageFields = append(messageFields, mf)
	}

	return messageFields, fqn
}

func getFQN(f *desc.FieldDescriptor, fqn map[string]int) (string, string) {
	name := f.GetFullyQualifiedName()
	if fqn == nil {
		return name, name
	}
	if _, ok := fqn[name]; !ok {
		return name, name
	}
	return fmt.Sprintf("%s[%d]", name, fqn[name]-1), name
}
