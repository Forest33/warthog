// Package entity provides entities for business logic.
package entity

import (
	"errors"
	"os"
	"strconv"

	"github.com/forest33/warthog/pkg/structs"
)

// protobuf data types.
const (
	TypeString   = "string"
	TypeBytes    = "bytes"
	TypeInt32    = "int32"
	TypeInt64    = "int64"
	TypeUInt32   = "uint32"
	TypeUInt64   = "uint64"
	TypeSInt32   = "sint32"
	TypeSInt64   = "sint64"
	TypeFixed32  = "fixed32"
	TypeFixed64  = "fixed64"
	TypeSFixed32 = "sfixed32"
	TypeSFixed64 = "sfixed64"
	TypeDouble   = "double"
	TypeFloat    = "float"
	TypeBool     = "bool"
	TypeEnum     = "enum"
	TypeMessage  = "message"
)

// Query gRPC request.
type Query struct {
	ServerID int64
	Service  string
	Method   string
	Data     map[string]interface{}
	Metadata []string
}

// QueryResponse gRPC response.
type QueryResponse struct {
	Time       string              `json:"time"`
	JsonString string              `json:"json_string"`
	SpentTime  string              `json:"spent_time"`
	Header     map[string][]string `json:"header"`
	Trailer    map[string][]string `json:"trailer"`
	Error      *Error              `json:"error"`
	Sent       uint                `json:"sent"`
	Received   uint                `json:"received"`
}

// Model creates Query from UI request.
func (r *Query) Model(server map[string]interface{}) error {
	if server == nil {
		return errors.New("empty data")
	}

	if v, ok := server["server_id"]; ok && v != nil {
		if id, ok := v.(float64); !ok {
			return errors.New("server id not a float")
		} else {
			r.ServerID = int64(id)
		}
	}
	if v, ok := server["service"]; ok {
		if r.Service, ok = v.(string); !ok {
			return errors.New("service not a string")
		}
	}
	if v, ok := server["method"]; ok {
		if r.Method, ok = v.(string); !ok {
			return errors.New("method not a string")
		}
	}
	if v, ok := server["data"]; ok {
		if r.Data, ok = v.(map[string]interface{}); !ok {
			return errors.New("data not a map[string]interface{}")
		}
	}
	if v, ok := server["metadata"]; ok {
		if m, ok := v.(map[string]interface{}); !ok {
			return errors.New("metadata not a map[string]interface{}")
		} else {
			r.Metadata = make([]string, 0, len(m)*2)
			for k, v := range m {
				if mv, ok := v.(string); !ok {
					return errors.New("metadata value not a string")
				} else {
					r.Metadata = append(r.Metadata, k, mv)
				}
			}
		}
	}

	return nil
}

// GetBool transforms to bool.
func GetBool(f *Field, val interface{}) (interface{}, error) {
	if f.Repeated {
		v, ok := val.([]interface{})
		if !ok {
			return nil, errors.New("field not a []interface{}")
		}
		return structs.MapWithError(v, func(i interface{}) (bool, error) {
			b, ok := i.(bool)
			if !ok {
				return false, errors.New("slice value not a bool")
			}
			return b, nil
		})
	}

	b, ok := val.(bool)
	if !ok {
		return nil, errors.New("field not a bool")
	}

	return b, nil
}

// GetString transforms to string.
func GetString(f *Field, val interface{}) (interface{}, error) {
	if f.Repeated {
		v, ok := val.([]interface{})
		if !ok {
			return nil, errors.New("field not a []interface{}")
		}
		return structs.MapWithError(v, func(i interface{}) (string, error) {
			s, ok := i.(string)
			if !ok {
				return "", errors.New("slice value not a string")
			}
			return s, nil
		})
	}

	v, ok := val.(string)
	if !ok {
		return "", errors.New("field not a string")
	}

	return v, nil
}

func getBytesValue(val map[string]interface{}) ([]byte, error) {
	if v, ok := val["file"]; ok {
		if fName, ok := v.(string); ok && len(fName) > 0 {
			data, err := os.ReadFile(fName)
			if err != nil {
				return nil, err
			}
			return data, nil
		}
	}

	v, ok := val["value"].(string)
	if !ok {
		return nil, errors.New("value not a string")
	}

	return []byte(v), nil
}

// GetBytes transforms to bytes.
func GetBytes(f *Field, val interface{}) (interface{}, error) {
	if f.Repeated {
		v, ok := val.([]interface{})
		if !ok {
			return nil, errors.New("field not a []interface{}")
		}
		bytes := make([][]byte, len(v))
		for i, item := range v {
			iv, ok := item.(map[string]interface{})
			if !ok {
				return "", errors.New("slice value not a map[string]interface{}")
			}
			b, err := getBytesValue(iv)
			if err != nil {
				return nil, err
			}
			bytes[i] = b
		}
		return bytes, nil
	}

	v, ok := val.(map[string]interface{})
	if !ok {
		return "", errors.New("value not a map[string]interface{}")
	}

	return getBytesValue(v)
}

// GetInt32 transforms to int32.
func GetInt32(f *Field, val interface{}) (interface{}, error) {
	if f.Repeated {
		v, ok := val.([]interface{})
		if !ok {
			return nil, errors.New("field not a []interface{}")
		}
		return structs.MapWithError(v, func(i interface{}) (int32, error) {
			v, err := strconv.ParseInt(i.(string), 10, 32)
			if err != nil {
				return 0, err
			}
			return int32(v), nil
		})
	}
	v, err := strconv.ParseInt(val.(string), 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}

// GetInt64 transforms to int64.
func GetInt64(f *Field, val interface{}) (interface{}, error) {
	if f.Repeated {
		v, ok := val.([]interface{})
		if !ok {
			return nil, errors.New("field not a []interface{}")
		}
		return structs.MapWithError(v, func(i interface{}) (int64, error) {
			s, ok := i.(string)
			if !ok {
				return 0, errors.New("slice value not a string")
			}
			return strconv.ParseInt(s, 10, 64)
		})
	}

	s, ok := val.(string)
	if !ok {
		return 0, errors.New("value not a string")
	}

	return strconv.ParseInt(s, 10, 32)
}

// GetUInt32 transforms to uint32.
func GetUInt32(f *Field, val interface{}) (interface{}, error) {
	if f.Repeated {
		v, ok := val.([]interface{})
		if !ok {
			return nil, errors.New("field not a []interface{}")
		}
		return structs.MapWithError(v, func(i interface{}) (uint32, error) {
			s, ok := i.(string)
			if !ok {
				return 0, errors.New("slice value not a string")
			}
			v, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return 0, err
			}
			return uint32(v), nil
		})
	}

	s, ok := val.(string)
	if !ok {
		return 0, errors.New("value not a string")
	}

	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return nil, err
	}

	return uint32(v), nil
}

// GetUInt64 transforms to uint64.
func GetUInt64(f *Field, val interface{}) (interface{}, error) {
	if f.Repeated {
		v, ok := val.([]interface{})
		if !ok {
			return nil, errors.New("field not a []interface{}")
		}
		return structs.MapWithError(v, func(i interface{}) (uint64, error) {
			s, ok := i.(string)
			if !ok {
				return 0, errors.New("slice value not a string")
			}
			return strconv.ParseUint(s, 10, 64)
		})
	}

	s, ok := val.(string)
	if !ok {
		return 0, errors.New("value not a string")
	}

	return strconv.ParseUint(s, 10, 32)
}

// GetFloat32 transforms to float32.
func GetFloat32(f *Field, val interface{}) (interface{}, error) {
	if f.Repeated {
		v, ok := val.([]interface{})
		if !ok {
			return nil, errors.New("field not a []interface{}")
		}
		return structs.MapWithError(v, func(i interface{}) (float32, error) {
			s, ok := i.(string)
			if !ok {
				return 0, errors.New("slice value not a string")
			}
			v, err := strconv.ParseFloat(s, 32)
			if err != nil {
				return 0, err
			}
			return float32(v), nil
		})
	}

	s, ok := val.(string)
	if !ok {
		return 0, errors.New("value not a string")
	}

	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return nil, err
	}

	return float32(v), nil
}

// GetFloat64 transforms to float64.
func GetFloat64(f *Field, val interface{}) (interface{}, error) {
	if f.Repeated {
		v, ok := val.([]interface{})
		if !ok {
			return nil, errors.New("field not a []interface{}")
		}
		return structs.MapWithError(v, func(i interface{}) (float64, error) {
			s, ok := i.(string)
			if !ok {
				return 0, errors.New("slice value not a string")
			}
			return strconv.ParseFloat(s, 64)
		})
	}

	s, ok := val.(string)
	if !ok {
		return 0, errors.New("value not a string")
	}

	return strconv.ParseFloat(s, 64)
}
