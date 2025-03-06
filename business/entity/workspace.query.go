// Package entity provides entities for business logic.
package entity

import (
	"errors"
)

// SavedQuery saved query.
type SavedQuery struct {
	Input    interface{} `json:"input"`
	Metadata interface{} `json:"metadata"`
}

// WorkspaceItemQuery stored query data.
type WorkspaceItemQuery struct {
	Service     string      `json:"service"`
	Method      string      `json:"method"`
	Description string      `json:"description"`
	Request     *SavedQuery `json:"request"`
}

// QueryRequest read/create/update/delete query.
type QueryRequest struct {
	ID       int64  `json:"id"`
	ServerID int64  `json:"server_id"`
	Title    string `json:"title"`
	WorkspaceItemQuery
}

// Model creates QueryRequest from UI request.
func (r *QueryRequest) Model(req map[string]interface{}) error {
	if req == nil {
		return errors.New("no data")
	}

	if v, ok := req["id"]; ok && v != nil {
		if id, ok := v.(float64); !ok {
			return errors.New("id not a float")
		} else {
			r.ID = int64(id)
		}
	}
	if v, ok := req["server_id"]; ok && v != nil {
		if id, ok := v.(float64); !ok {
			return errors.New("server id not a float")
		} else {
			r.ID = int64(id)
		}
	}
	if v, ok := req["title"]; ok && v != nil {
		if r.Title, ok = v.(string); !ok {
			return errors.New("title not a string")
		}
	}

	r.WorkspaceItemQuery = WorkspaceItemQuery{}

	return r.WorkspaceItemQuery.Model(req)
}

// Model creates WorkspaceItemQuery from UI request.
func (s *WorkspaceItemQuery) Model(req map[string]interface{}) error {
	if req == nil {
		return errors.New("no data")
	}

	if v, ok := req["service"]; ok && v != nil {
		if s.Service, ok = v.(string); !ok {
			return errors.New("service not a string")
		}
	}
	if v, ok := req["method"]; ok && v != nil {
		if s.Method, ok = v.(string); !ok {
			return errors.New("method not a string")
		}
	}
	if v, ok := req["request"]; ok && v != nil {
		r, ok := req["request"].(map[string]interface{})
		if !ok {
			return errors.New("request not a map[string]interface{}")
		}
		sq := &SavedQuery{}
		sq.Model(r)
		s.Request = sq
	}
	if v, ok := req["description"]; ok && v != nil {
		if s.Description, ok = v.(string); !ok {
			return errors.New("description not a string")
		}
	}

	return nil
}

// Model creates SavedQuery from UI request.
func (s *SavedQuery) Model(req map[string]interface{}) {
	if req == nil {
		return
	}

	if v, ok := req["input"]; ok && v != nil {
		s.Input = v
	}
	if v, ok := req["metadata"]; ok && v != nil {
		s.Metadata = v
	}
}
