// Package entity provides entities for business logic.
package entity

import (
	"fmt"
	"strconv"

	"github.com/forest33/warthog/pkg/structs"
)

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

// Query gRPC request
type Query struct {
	Service  string
	Method   string
	Data     map[string]interface{}
	Metadata []string
}

// QueryResponse gRPC response
type QueryResponse struct {
	JsonString string              `json:"json_string"`
	SpentTime  string              `json:"spent_time"`
	Header     map[string][]string `json:"header"`
	Trailer    map[string][]string `json:"trailer"`
}

// Model creates Query from UI request
func (r *Query) Model(server map[string]interface{}) error {
	if server == nil {
		return fmt.Errorf("empty data")
	}

	if v, ok := server["service"]; ok {
		r.Service = v.(string)
	}
	if v, ok := server["method"]; ok {
		r.Method = v.(string)
	}
	if v, ok := server["data"]; ok {
		r.Data = v.(map[string]interface{})
	}
	if v, ok := server["metadata"]; ok {
		r.Metadata = make([]string, 0, len(v.(map[string]interface{}))*2)
		for k, v := range v.(map[string]interface{}) {
			r.Metadata = append(r.Metadata, k, v.(string))
		}
	}

	return nil
}

// GetBool transforms to bool
func GetBool(f *Field, val interface{}) interface{} {
	if f.Repeated {
		return structs.Map(val.([]interface{}), func(i interface{}) bool { return i.(bool) })
	}
	return val.(bool)
}

// GetString transforms to string
func GetString(f *Field, val interface{}) interface{} {
	if f.Repeated {
		return structs.Map(val.([]interface{}), func(i interface{}) string { return i.(string) })
	}
	return val.(string)
}

// GetBytes transforms to bytes
func GetBytes(f *Field, val interface{}) interface{} {
	if f.Repeated {
		return structs.Map(val.([]interface{}), func(i interface{}) []byte { return []byte(i.(string)) })
	}
	return []byte(val.(string))
}

// GetInt32 transforms to int32
func GetInt32(f *Field, val interface{}) (interface{}, error) {
	if f.Repeated {
		return structs.MapWithError(val.([]interface{}), func(i interface{}) (int32, error) {
			v, err := strconv.ParseInt(i.(string), 10, 32)
			if err != nil {
				return 0, err
			}
			return int32(v), nil
		})
	}
	v, err := strconv.ParseInt(val.(string), 10, 32)
	if err != nil {
		return nil, err
	}
	return int32(v), nil
}

// GetInt64 transforms to int64
func GetInt64(f *Field, val interface{}) (interface{}, error) {
	if f.Repeated {
		return structs.MapWithError(val.([]interface{}), func(i interface{}) (int64, error) {
			return strconv.ParseInt(i.(string), 10, 64)
		})
	}
	return strconv.ParseInt(val.(string), 10, 32)
}

// GetUInt32 transforms to uint32
func GetUInt32(f *Field, val interface{}) (interface{}, error) {
	if f.Repeated {
		return structs.MapWithError(val.([]interface{}), func(i interface{}) (uint32, error) {
			v, err := strconv.ParseUint(i.(string), 10, 32)
			if err != nil {
				return 0, err
			}
			return uint32(v), nil
		})
	}
	v, err := strconv.ParseUint(val.(string), 10, 32)
	if err != nil {
		return nil, err
	}
	return uint32(v), nil
}

// GetUInt64 transforms to uint64
func GetUInt64(f *Field, val interface{}) (interface{}, error) {
	if f.Repeated {
		return structs.MapWithError(val.([]interface{}), func(i interface{}) (uint64, error) {
			return strconv.ParseUint(i.(string), 10, 64)
		})
	}
	return strconv.ParseUint(val.(string), 10, 32)
}

// GetFloat32 transforms to float32
func GetFloat32(f *Field, val interface{}) (interface{}, error) {
	if f.Repeated {
		return structs.MapWithError(val.([]interface{}), func(i interface{}) (float32, error) {
			v, err := strconv.ParseFloat(i.(string), 32)
			if err != nil {
				return 0, err
			}
			return float32(v), nil
		})
	}
	v, err := strconv.ParseFloat(val.(string), 32)
	if err != nil {
		return nil, err
	}
	return float32(v), nil
}

// GetFloat64 transforms to float64
func GetFloat64(f *Field, val interface{}) (interface{}, error) {
	if f.Repeated {
		return structs.MapWithError(val.([]interface{}), func(i interface{}) (float64, error) {
			return strconv.ParseFloat(i.(string), 64)
		})
	}
	return strconv.ParseFloat(val.(string), 64)
}
