package entity

import (
	"encoding/json"
	"errors"
	"strings"
)

const (
	// AuthTypeNone do not use authentication.
	AuthTypeNone = "none"
	// AuthTypeBasic basic authentication.
	AuthTypeBasic = "basic"
	// AuthTypeBearer bearer token authentication.
	AuthTypeBearer = "bearer"
	// AuthTypeJWT jwt token authentication.
	AuthTypeJWT = "jwt"
	// AuthTypeGCE Google Compute Engine authentication.
	AuthTypeGCE = "google"
)

// Auth authentication data.
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

// Model creates Auth from UI request.
func (s *Auth) Model(auth map[string]interface{}) error {
	authType, ok := auth["type"]
	if !ok || authType == AuthTypeNone {
		s.Type = AuthTypeNone
		return nil
	}

	if s.Type, ok = authType.(string); !ok {
		return errors.New("authType not a string")
	}

	if v, ok := auth["login"]; ok {
		if s.Login, ok = v.(string); !ok {
			return errors.New("login not a string")
		}
	}
	if v, ok := auth["password"]; ok {
		if s.Password, ok = v.(string); !ok {
			return errors.New("password not a string")
		}
	}
	if v, ok := auth["token"]; ok {
		if s.Token, ok = v.(string); !ok {
			return errors.New("token not a string")
		}
	}
	if v, ok := auth["algorithm"]; ok {
		if s.Algorithm, ok = v.(string); !ok {
			return errors.New("algorithm not a string")
		}
	}
	if v, ok := auth["secret"]; ok {
		if s.Secret, ok = v.(string); !ok {
			return errors.New("secret not a string")
		}
	}
	if v, ok := auth["private_key"]; ok {
		if s.PrivateKey, ok = v.(string); !ok {
			return errors.New("private key not a string")
		}
	}
	if v, ok := auth["secret_base64"]; ok {
		if s.SecretBase64, ok = v.(bool); !ok {
			return errors.New("secret base64 not a boolean")
		}
	}
	if v, ok := auth["header_prefix"]; ok {
		if hp, ok := v.(string); !ok {
			return errors.New("header prefix not a string")
		} else {
			s.HeaderPrefix = strings.TrimSpace(hp)
		}
	}
	if v, ok := auth["payload"]; ok {
		if p, ok := v.(string); !ok {
			return errors.New("payload not a string")
		} else {
			s.Payload = map[string]interface{}{}
			if err := json.Unmarshal([]byte(p), &s.Payload); err != nil {
				return err
			}
		}
	}
	if v, ok := auth["google_token"]; ok {
		if s.GoogleToken, ok = v.(string); !ok {
			return errors.New("google token not a string")
		}
	}
	if v, ok := auth["google_scopes"]; ok {
		if gs, ok := v.(string); !ok {
			return errors.New("google scopes not a string")
		} else {
			s.GoogleScopes = strings.Split(gs, ",")
		}
	}

	return nil
}
