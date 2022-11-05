package grpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/grpcreflect"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	"github.com/forest33/warthog/business/entity"
)

func (c *Client) AddProtobuf(path ...string) {
	c.protoPath = path
}

func (c *Client) AddImport(path ...string) {
	c.importPath = path
}

func (c *Client) LoadFromProtobuf() ([]*entity.Service, error) {
	if len(c.protoPath) == 0 {
		return nil, fmt.Errorf("empty path to protobuf's")
	}

	parser := protoparse.Parser{
		ImportPaths:     c.importPath,
		ErrorReporter:   nil, // TODO implement
		WarningReporter: nil, // TODO implement
	}

	services := make([]*entity.Service, 0, len(c.protoPath))

	for _, p := range c.protoPath {
		fd, err := parser.ParseFiles(p)
		if err != nil {
			return nil, err
		}
		if len(fd) != 1 {
			return nil, fmt.Errorf("wrong parse result")
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

	return services, nil
}

func (c *Client) LoadFromReflection() ([]*entity.Service, error) {
	ctx, cancel := context.WithTimeout(c.ctx, time.Second*time.Duration(c.cfg.ConnectTimeout))
	defer cancel()

	stub := rpb.NewServerReflectionClient(c.conn)
	client := grpcreflect.NewClient(ctx, stub)
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
			Input:      c.getFields(md.GetInputType().GetFields(), nil),
		})
	}

	if c.cfg.SortMethodsByName {
		c.sortMethodsByName(methods)
	}

	return methods
}

func getMethodType(md *desc.MethodDescriptor) string {
	switch {
	case md.IsClientStreaming() && md.IsServerStreaming():
		return entity.MethodTypeClientServerStream
	case md.IsClientStreaming():
		return entity.MethodTypeClientStream
	case md.IsServerStreaming():
		return entity.MethodTypeServerStream
	}
	return entity.MethodTypeUnary
}

func (c *Client) getFields(fd []*desc.FieldDescriptor, parent *desc.FieldDescriptor) []*entity.Field {
	fields := make([]*entity.Field, 0, len(fd))

	for _, f := range fd {
		var (
			fieldType = f.GetType().String()
		)

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
				Fqn:        f.GetFullyQualifiedName(),
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
				mapField := &entity.Field{
					Descriptor: f,
					Fqn:        f.GetFullyQualifiedName(),
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
					mapField.Map.ValueDescriptor = f.GetMapValueType().GetMessageType()
					mapField.Map.Fields = c.getFields(f.GetMapValueType().GetMessageType().GetFields(), f)
				} else if valueType == entity.ProtoTypeEnum {
					mapField.Map.Fields = c.getFields([]*desc.FieldDescriptor{f.GetMapValueType()}, f)
				}
				fields = append(fields, mapField)
			} else {
				msgField := &entity.Field{
					Descriptor: f,
					Fqn:        f.GetFullyQualifiedName(),
					Name:       f.GetName(),
					Type:       getTypeName(fieldType),
					ParentType: getParentTypeName(parent),
					Repeated:   f.IsRepeated(),
					OneOf:      getOneOf(f),
					Message: &entity.Message{
						Name:       f.GetMessageType().GetName(),
						Type:       f.GetMessageType().GetFullyQualifiedName(),
						Fields:     c.getFields(f.GetMessageType().GetFields(), f),
						Descriptor: f.GetMessageType(),
					},
				}
				fields = append(fields, msgField)
			}
		default:
			fields = append(fields, &entity.Field{
				Descriptor: f,
				Fqn:        f.GetFullyQualifiedName(),
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
