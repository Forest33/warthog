package entity

import (
	"fmt"
	"time"

	"google.golang.org/grpc/status"
)

const (
	GUIResponseStatusOK    GUIResponseStatus = "ok"
	GUIResponseStatusError GUIResponseStatus = "error"
)

type GUIConfig struct {
	WindowWidth  int
	WindowHeight int
	WindowX      *int
	WindowY      *int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type GUIResponseStatus string

type GUIRequest struct {
	Cmd     GUICommand             `json:"name"`
	Payload map[string]interface{} `json:"payload"`
}

type GUIResponse struct {
	Status  GUIResponseStatus `json:"status"`
	Error   Error             `json:"error,omitempty"`
	Payload interface{}       `json:"payload,omitempty"`
}

type Error struct {
	Code            uint32 `json:"code"`
	CodeDescription string `json:"code_description"`
	Message         string `json:"message"`
}

func (e Error) Error() string {
	if e.Code > 0 {
		return fmt.Sprintf("error code: %s message: %s", e.Code, e.Message)
	}
	return e.Message
}

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
