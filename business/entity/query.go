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

type Query struct {
	Service  string
	Method   string
	Data     map[string]interface{}
	Metadata []string
}

type QueryResponse struct {
	JsonString string              `json:"json_string"`
	SpentTime  string              `json:"spent_time"`
	Header     map[string][]string `json:"header"`
	Trailer    map[string][]string `json:"trailer"`
}

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

func GetBool(f *Field, val interface{}) interface{} {
	if f.Repeated {
		return structs.Map(val.([]interface{}), func(i interface{}) bool { return i.(bool) })
	}
	return val.(bool)
}

func GetString(f *Field, val interface{}) interface{} {
	if f.Repeated {
		return structs.Map(val.([]interface{}), func(i interface{}) string { return i.(string) })
	}
	return val.(string)
}

func GetBytes(f *Field, val interface{}) interface{} {
	if f.Repeated {
		return structs.Map(val.([]interface{}), func(i interface{}) []byte { return []byte(i.(string)) })
	}
	return []byte(val.(string))
}

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

func GetInt64(f *Field, val interface{}) (interface{}, error) {
	if f.Repeated {
		return structs.MapWithError(val.([]interface{}), func(i interface{}) (int64, error) {
			return strconv.ParseInt(i.(string), 10, 64)
		})
	}
	return strconv.ParseInt(val.(string), 10, 32)
}

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

func GetUInt64(f *Field, val interface{}) (interface{}, error) {
	if f.Repeated {
		return structs.MapWithError(val.([]interface{}), func(i interface{}) (uint64, error) {
			return strconv.ParseUint(i.(string), 10, 64)
		})
	}
	return strconv.ParseUint(val.(string), 10, 32)
}

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

func GetFloat64(f *Field, val interface{}) (interface{}, error) {
	if f.Repeated {
		return structs.MapWithError(val.([]interface{}), func(i interface{}) (float64, error) {
			return strconv.ParseFloat(i.(string), 64)
		})
	}
	return strconv.ParseFloat(val.(string), 64)
}
