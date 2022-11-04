package entity

import (
	"fmt"
)

type WorkspaceItemQuery struct {
	Service     string      `json:"service"`
	Method      string      `json:"method"`
	Request     interface{} `json:"request"`
	Description string      `json:"description"`
}

type QueryRequest struct {
	ID       int64  `json:"id"`
	ServerID int64  `json:"server_id"`
	Title    string `json:"title"`
	WorkspaceItemQuery
}

func (r *QueryRequest) Model(req map[string]interface{}) error {
	if req == nil {
		return fmt.Errorf("empty data")
	}

	if v, ok := req["id"]; ok && v != nil {
		r.ID = int64(v.(float64))
	}
	if v, ok := req["server_id"]; ok && v != nil {
		r.ServerID = int64(v.(float64))
	}
	if v, ok := req["title"]; ok && v != nil {
		r.Title = v.(string)
	}

	r.WorkspaceItemQuery = WorkspaceItemQuery{}

	return r.WorkspaceItemQuery.Model(req)
}

func (s *WorkspaceItemQuery) Model(req map[string]interface{}) error {
	if req == nil {
		return fmt.Errorf("empty data")
	}

	if v, ok := req["service"]; ok && v != nil {
		s.Service = v.(string)
	}
	if v, ok := req["method"]; ok && v != nil {
		s.Method = v.(string)
	}
	if v, ok := req["request"]; ok && v != nil {
		s.Request = v
	}
	if v, ok := req["description"]; ok && v != nil {
		s.Description = v.(string)
	}

	return nil
}
