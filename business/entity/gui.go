// Package entity provides entities for business logic.
package entity

import (
	"fmt"
	"time"

	"google.golang.org/grpc/status"
)

// UI response statuses
const (
	GUIResponseStatusOK    GUIResponseStatus = "ok"
	GUIResponseStatusError GUIResponseStatus = "error"
)

// GUIConfig UI settings
type GUIConfig struct {
	WindowWidth  int
	WindowHeight int
	WindowX      *int
	WindowY      *int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// GUIResponseStatus response status
type GUIResponseStatus string

// GUIRequest UI request
type GUIRequest struct {
	Cmd     GUICommand             `json:"name"`
	Payload map[string]interface{} `json:"payload"`
}

// GUIResponse UI response
type GUIResponse struct {
	Status  GUIResponseStatus `json:"status"`
	Error   Error             `json:"error,omitempty"`
	Payload interface{}       `json:"payload,omitempty"`
}

// Error UI error response
type Error struct {
	Code            uint32 `json:"code"`
	CodeDescription string `json:"code_description"`
	Message         string `json:"message"`
}

// Error returns error string
func (e Error) Error() string {
	if e.Code > 0 {
		return fmt.Sprintf("error code: %d message: %s", e.Code, e.Message)
	}
	return e.Message
}

// ErrorGUIResponse returns UI error response
func ErrorGUIResponse(err error) *GUIResponse {
	return &GUIResponse{
		Status: GUIResponseStatusError,
		Error: Error{
			Code:            uint32(status.Code(err)),
			CodeDescription: status.Code(err).String(),
			Message:         err.Error(),
		},
	}
}
