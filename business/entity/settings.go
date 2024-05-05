// Package entity provides entities for business logic.
package entity

import (
	"errors"

	"github.com/forest33/warthog/pkg/structs"
)

// Settings application settings.
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
	EmitDefaults          *bool `json:"emit_defaults"`
	CheckUpdates          *bool `json:"check_updates"`
}

// DefaultSettings settings by default.
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
	EmitDefaults:          structs.Ref(false),
	CheckUpdates:          structs.Ref(true),
}

// Model creates Settings from UI request.
func (s *Settings) Model(payload map[string]interface{}) error {
	if payload == nil {
		return errors.New("no data")
	}

	if v, ok := payload["single_instance"]; ok && v != nil {
		b, ok := v.(bool)
		if !ok {
			return errors.New("single instance not a bool")
		}
		s.SingleInstance = &b
	}
	if v, ok := payload["connect_timeout"]; ok && v != nil {
		f, ok := v.(float64)
		if !ok {
			return errors.New("connection timeout not a float")
		}
		s.ConnectTimeout = structs.Ref(int(f))
	}
	if v, ok := payload["request_timeout"]; ok && v != nil {
		f, ok := v.(float64)
		if !ok {
			return errors.New("request timeout not a float")
		}
		s.RequestTimeout = structs.Ref(int(f))
	}
	if v, ok := payload["k8s_request_timeout"]; ok && v != nil {
		f, ok := v.(float64)
		if !ok {
			return errors.New("k8s request timeout not a float")
		}
		s.K8SRequestTimeout = structs.Ref(int(f))
	}
	if v, ok := payload["non_blocking_connection"]; ok && v != nil {
		b, ok := v.(bool)
		if !ok {
			return errors.New("non blocking connection not a bool")
		}
		s.NonBlockingConnection = &b
	}
	if v, ok := payload["sort_methods_by_name"]; ok && v != nil {
		b, ok := v.(bool)
		if !ok {
			return errors.New("sort methods by name not a bool")
		}
		s.SortMethodsByName = &b
	}
	if v, ok := payload["max_loop_depth"]; ok && v != nil {
		f, ok := v.(float64)
		if !ok {
			return errors.New("max loop depth not a float")
		}
		s.MaxLoopDepth = structs.Ref(int(f))
	}
	if v, ok := payload["emit_defaults"]; ok && v != nil {
		b, ok := v.(bool)
		if !ok {
			return errors.New("emit defaults not a bool")
		}
		s.EmitDefaults = &b
	}
	if v, ok := payload["check_updates"]; ok && v != nil {
		b, ok := v.(bool)
		if !ok {
			return errors.New("check updates not a bool")
		}
		s.CheckUpdates = &b
	}

	return nil
}

func (s *Settings) IsEmitDefaults() bool {
	return s.EmitDefaults != nil && *s.EmitDefaults
}

func (s *Settings) IsCheckUpdates() bool {
	return s.CheckUpdates != nil && *s.CheckUpdates
}
