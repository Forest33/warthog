package entity

import (
	"github.com/jhump/protoreflect/desc"
)

const (
	ReflectionServiceFQN         = "grpc.reflection.v1alpha.ServerReflection"
	MethodTypeUnary              = "u"
	MethodTypeClientStream       = "cs"
	MethodTypeServerStream       = "ss"
	MethodTypeClientServerStream = "css"
)

// Service gRPC service
type Service struct {
	Name    string    `json:"name"`
	Methods []*Method `json:"methods,omitempty"`
}

// Method gRPC method
type Method struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Input      []*Field               `json:"input,omitempty"`
	Descriptor *desc.MethodDescriptor `json:"-"`
}

// LoadServerResponse server data, methods, and saved queries
type LoadServerResponse struct {
	Server   *Workspace `json:"server"`
	Services []*Service `json:"services"`
	Query    *Workspace `json:"query"`
}

// Field protobuf field
type Field struct {
	FQN        string                `json:"fqn"`
	ProtoFQN   string                `json:"proto_fqn"`
	Name       string                `json:"name"`
	Type       string                `json:"type"`
	ParentType string                `json:"parent_type"`
	Repeated   bool                  `json:"repeated"`
	Enum       *Enum                 `json:"enum,omitempty"`
	Map        *Map                  `json:"map,omitempty"`
	Message    *Message              `json:"message,omitempty"`
	OneOf      *OneOf                `json:"oneof,omitempty"`
	Descriptor *desc.FieldDescriptor `json:"-"`
}

// Map protobuf map
type Map struct {
	KeyType         string                  `json:"key_type"`
	ValueType       string                  `json:"value_type"`
	ValueTypeFqn    string                  `json:"-"`
	ProtoValueType  string                  `json:"-"`
	ValueDescriptor *desc.MessageDescriptor `json:"-"`
	Fields          []*Field                `json:"fields,omitempty"`
}

// Message protobuf message
type Message struct {
	Name       string                  `json:"name"`
	Type       string                  `json:"type"`
	Fields     []*Field                `json:"fields,omitempty"`
	Descriptor *desc.MessageDescriptor `json:"-"`
}

// Enum protobuf enum
type Enum struct {
	ValueType string       `json:"value_type"`
	Values    []*EnumValue `json:"values"`
}

// EnumValue protobuf enum value
type EnumValue struct {
	Name   string `json:"name"`
	Number int32  `json:"number"`
}

// OneOf protobuf oneof
type OneOf struct {
	Fqn  string `json:"fqn"`
	Name string `json:"name"`
}
