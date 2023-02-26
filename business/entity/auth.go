package entity

import (
	"encoding/json"
	"strings"
)

const (
	AuthTypeNone   = "none"
	AuthTypeBasic  = "basic"
	AuthTypeBearer = "bearer"
	AuthTypeJWT    = "jwt"
	AuthTypeGCE    = "google"
)

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
