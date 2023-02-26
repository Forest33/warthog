// Package entity provides entities for business logic.
package entity

import (
	"errors"

	"github.com/forest33/warthog/pkg/structs"
)

// Settings application settings
type Settings struct {
	WindowWidth           int   `json:"window_width"`
	WindowHeight          int   `json:"window_height"`
	WindowX               *int  `json:"window_x"`
	WindowY               *int  `json:"window_y"`
	SingleInstance        *bool `json:"single_instance"`
	ConnectTimeout        *int  `json:"connect_timeout"`
	RequestTimeout        *int  `json:"request_timeout"`
	K8SRequestTimeout     *int  `json:"k8s_request_timeout"`
	NonBlockingConnection *bool `json:"non_blocking_connection"`
	SortMethodsByName     *bool `json:"sort_methods_by_name"`
	MaxLoopDepth          *int  `json:"max_loop_depth"`
}

// DefaultSettings settings by default
var DefaultSettings = &Settings{
	WindowWidth:           1024,
	WindowHeight:          768,
	WindowX:               structs.Ref(50),
	WindowY:               structs.Ref(50),
	SingleInstance:        structs.Ref(true),
	ConnectTimeout:        structs.Ref(10),
	RequestTimeout:        structs.Ref(30),
	K8SRequestTimeout:     structs.Ref(30),
	NonBlockingConnection: structs.Ref(true),
	SortMethodsByName:     structs.Ref(true),
	MaxLoopDepth:          structs.Ref(10),
}

// Model creates Settings from UI request
func (s *Settings) Model(payload map[string]interface{}) error {
	if payload == nil {
		return errors.New("no data")
	}

	if v, ok := payload["single_instance"]; ok && v != nil {
		s.SingleInstance = structs.Ref(v.(bool))
	}
	if v, ok := payload["connect_timeout"]; ok && v != nil {
		s.ConnectTimeout = structs.Ref(int(v.(float64)))
	}
	if v, ok := payload["request_timeout"]; ok && v != nil {
		s.RequestTimeout = structs.Ref(int(v.(float64)))
	}
	if v, ok := payload["k8s_request_timeout"]; ok && v != nil {
		s.K8SRequestTimeout = structs.Ref(int(v.(float64)))
	}
	if v, ok := payload["non_blocking_connection"]; ok && v != nil {
		s.NonBlockingConnection = structs.Ref(v.(bool))
	}
	if v, ok := payload["sort_methods_by_name"]; ok && v != nil {
		s.SortMethodsByName = structs.Ref(v.(bool))
	}
	if v, ok := payload["max_loop_depth"]; ok && v != nil {
		s.MaxLoopDepth = structs.Ref(int(v.(float64)))
	}

	return nil
}
