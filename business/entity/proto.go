// Package entity provides entities for business logic.
package entity

import (
	"github.com/jhump/protoreflect/desc/protoparse/ast"
)

// protobuf types.
const (
	ProtoTypeEnum    = "TYPE_ENUM"
	ProtoTypeMessage = "TYPE_MESSAGE"
)

// ProtobufError protobuf parsing warning or error.
type ProtobufError struct {
	Pos             ast.SourcePos `json:"pos"`
	Code            uint32        `json:"code"`
	CodeDescription string        `json:"code_description"`
	Warning         string        `json:"warning"`
	Err             error         `json:"err"`
}

// Error returns error string.
func (pe ProtobufError) Error() string {
	return pe.Err.Error()
}
