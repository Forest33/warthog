// Package entity provides entities for business logic.
package entity

import (
	"fmt"

	"github.com/forest33/warthog/pkg/structs"
)

// ServerRequest read/create/delete server request
type ServerRequest struct {
	ID       int64  `json:"id"`
	FolderID int64  `json:"folder_id"`
	Title    string `json:"title"`
	WorkspaceItemServer
}

// ServerResponse read/create/update server response
type ServerResponse struct {
	Server *Workspace           `json:"server"`
	Query  *Workspace           `json:"query"`
	Tree   []*WorkspaceTreeNode `json:"tree"`
}

// ServerUpdateRequest update server request
type ServerUpdateRequest struct {
	ID      int64                  `json:"id"`
	Service string                 `json:"service"`
	Method  string                 `json:"method"`
	Request map[string]interface{} `json:"request"`
}

// WorkspaceItemServer stored server data
type WorkspaceItemServer struct {
	Addr              string                            `json:"addr,omitempty"`
	UseReflection     bool                              `json:"use_reflection,omitempty"`
	ProtoFiles        []string                          `json:"proto_files,omitempty"`
	ImportPath        []string                          `json:"import_path,omitempty"`
	NoTLS             bool                              `json:"no_tls,omitempty"`
	Insecure          bool                              `json:"insecure,omitempty"`
	RootCertificate   string                            `json:"root_certificate,omitempty"`
	ClientCertificate string                            `json:"client_certificate,omitempty"`
	ClientKey         string                            `json:"client_key,omitempty"`
	Request           map[string]map[string]interface{} `json:"request"`
}

// Model creates ServerRequest from UI request
func (r *ServerRequest) Model(req map[string]interface{}) error {
	if req == nil {
		return fmt.Errorf("no data")
	}

	if v, ok := req["id"]; ok && v != nil {
		r.ID = int64(v.(float64))
	}
	if v, ok := req["folder_id"]; ok && v != nil {
		r.FolderID = int64(v.(float64))
	}
	if v, ok := req["title"]; ok && v != nil {
		r.Title = v.(string)
	}

	r.WorkspaceItemServer = WorkspaceItemServer{}

	return r.WorkspaceItemServer.Model(req)
}

// Model creates ServerUpdateRequest from UI request
func (r *ServerUpdateRequest) Model(req map[string]interface{}) error {
	if req == nil {
		return fmt.Errorf("no data")
	}

	if v, ok := req["id"]; ok && v != nil {
		r.ID = int64(v.(float64))
	}
	if v, ok := req["service"]; ok && v != nil {
		r.Service = v.(string)
	}
	if v, ok := req["method"]; ok && v != nil {
		r.Method = v.(string)
	}
	if v, ok := req["request"]; ok && v != nil {
		r.Request = v.(map[string]interface{})
	}

	return nil
}

// Model creates WorkspaceItemServer from UI request
func (s *WorkspaceItemServer) Model(server map[string]interface{}) error {
	if server == nil {
		return fmt.Errorf("no data")
	}

	if v, ok := server["addr"]; ok && v != nil {
		s.Addr = v.(string)
	}
	if v, ok := server["use_reflection"]; ok && v != nil {
		s.UseReflection = v.(bool)
	}
	if v, ok := server["proto_files"]; ok && v != nil {
		s.ProtoFiles = structs.Map(v.([]interface{}), func(p interface{}) string { return p.(string) })
	}
	if v, ok := server["import_path"]; ok && v != nil {
		s.ImportPath = structs.Map(v.([]interface{}), func(p interface{}) string { return p.(string) })
	}
	if v, ok := server["no_tls"]; ok && v != nil {
		s.NoTLS = v.(bool)
	}
	if v, ok := server["insecure"]; ok && v != nil {
		s.Insecure = v.(bool)
	}
	if v, ok := server["root_certificate"]; ok && v != nil {
		s.RootCertificate = v.(string)
	}
	if v, ok := server["client_certificate"]; ok && v != nil {
		s.ClientCertificate = v.(string)
	}
	if v, ok := server["client_key"]; ok && v != nil {
		s.ClientKey = v.(string)
	}

	return nil
}
