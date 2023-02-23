// Package entity provides entities for business logic.
package entity

import (
	"encoding/json"
	"errors"
	"strings"

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
	ID      int64       `json:"id"`
	Service string      `json:"service"`
	Method  string      `json:"method"`
	Request *SavedQuery `json:"request"`
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
	Request           map[string]map[string]*SavedQuery `json:"request"`
	Auth              *Auth                             `json:"auth"`
}

// Auth authentication data
type Auth struct {
	Type         string                 `json:"type,omitempty"`
	Login        string                 `json:"login,omitempty"`
	Password     string                 `json:"password,omitempty"`
	Token        string                 `json:"token,omitempty"`
	Algorithm    string                 `json:"algorithm,omitempty"`
	Secret       string                 `json:"secret,omitempty"`
	PrivateKey   string                 `json:"private_key,omitempty"`
	SecretBase64 bool                   `json:"secret_base64,omitempty"`
	HeaderPrefix string                 `json:"header_prefix,omitempty"`
	Payload      map[string]interface{} `json:"payload,omitempty"`
	GoogleScopes []string               `json:"google_scopes,omitempty"`
	GoogleToken  string                 `json:"google_token,omitempty"`
}

// Model creates ServerRequest from UI request
func (r *ServerRequest) Model(req map[string]interface{}) error {
	if req == nil {
		return errors.New("no data")
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
		return errors.New("no data")
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
		sq := &SavedQuery{}
		sq.Model(req["request"].(map[string]interface{}))
		r.Request = sq
	}

	return nil
}

// Model creates WorkspaceItemServer from UI request
func (s *WorkspaceItemServer) Model(server map[string]interface{}) error {
	if server == nil {
		return errors.New("no data")
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
	if v, ok := server["auth"]; ok && v != nil {
		s.Auth = &Auth{}
		if err := s.Auth.Model(v.(map[string]interface{})); err != nil {
			return err
		}
	}

	return nil
}

// Model creates Auth from UI request
func (s *Auth) Model(auth map[string]interface{}) error {
	authType, ok := auth["type"]
	if !ok || authType == AuthTypeNone {
		s.Type = AuthTypeNone
		return nil
	}

	s.Type = authType.(string)

	if v, ok := auth["login"]; ok {
		s.Login = v.(string)
	}
	if v, ok := auth["password"]; ok {
		s.Password = v.(string)
	}
	if v, ok := auth["token"]; ok {
		s.Token = v.(string)
	}
	if v, ok := auth["algorithm"]; ok {
		s.Algorithm = v.(string)
	}
	if v, ok := auth["secret"]; ok {
		s.Secret = v.(string)
	}
	if v, ok := auth["private_key"]; ok {
		s.PrivateKey = v.(string)
	}
	if v, ok := auth["secret_base64"]; ok {
		s.SecretBase64 = v.(bool)
	}
	if v, ok := auth["header_prefix"]; ok {
		s.HeaderPrefix = strings.TrimSpace(v.(string))
	}
	if v, ok := auth["payload"]; ok {
		s.Payload = map[string]interface{}{}
		if err := json.Unmarshal([]byte(v.(string)), &s.Payload); err != nil {
			return err
		}
	}
	if v, ok := auth["google_token"]; ok {
		s.GoogleToken = v.(string)
	}
	if v, ok := auth["google_scopes"]; ok {
		s.GoogleScopes = strings.Split(v.(string), ",")
	}

	return nil
}
