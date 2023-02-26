// Package entity provides entities for business logic.
package entity

import (
	"fmt"

	"github.com/jhump/protoreflect/desc/protoparse/ast"
	"google.golang.org/grpc/status"
)

// UI response statuses
const (
	GUIResponseStatusOK    GUIResponseStatus = "ok"
	GUIResponseStatusError GUIResponseStatus = "error"
)

// GUIResponseStatus response status
type GUIResponseStatus string

// GUIRequest UI request
type GUIRequest struct {
	Cmd     GUICommand  `json:"name"`
	Payload interface{} `json:"payload"`
}

// GUIResponse UI response
type GUIResponse struct {
	Status  GUIResponseStatus `json:"status"`
	Error   Error             `json:"error,omitempty"`
	Payload interface{}       `json:"payload,omitempty"`
}

// Error UI error response
type Error struct {
	Pos             ast.SourcePos `json:"pos"`
	Code            uint32        `json:"code"`
	CodeDescription string        `json:"code_description"`
	Message         string        `json:"message"`
}

// Info UI info message
type Info struct {
	Message string `json:"message"`
}

// Error returns error string
func (e Error) Error() string {
	if e.Code > 0 {
		return fmt.Sprintf("error code: %d message: %s", e.Code, e.Message)
	}
	return e.Message
}

// ErrorGUIResponse returns UI error response
func ErrorGUIResponse(err error, payload ...interface{}) *GUIResponse {
	resp := &GUIResponse{
		Status: GUIResponseStatusError,
		Error: Error{
			Code:            uint32(status.Code(err)),
			CodeDescription: status.Code(err).String(),
			Message:         err.Error(),
		},
	}

	if len(payload) > 1 && len(payload)%2 == 0 {
		resp.Payload = make(map[string]interface{}, len(payload)/2)
	}

	for i := 0; i < len(payload); i += 2 {
		resp.Payload.(map[string]interface{})[payload[i].(string)] = payload[i+1]
	}

	return resp
}
