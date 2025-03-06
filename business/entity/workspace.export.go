package entity

import "errors"

// ExportRequest export request.
type ExportRequest struct {
	Path string `json:"path"`
}

// Model creates ExportRequest from UI request.
func (r *ExportRequest) Model(req map[string]interface{}) error {
	if req == nil {
		return errors.New("no data")
	}

	if v, ok := req["path"]; ok && v != nil {
		if r.Path, ok = v.(string); !ok {
			return errors.New("path not a string")
		}
	}

	return nil
}
