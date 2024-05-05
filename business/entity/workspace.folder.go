// Package entity provides entities for business logic.
package entity

import (
	"errors"

	"github.com/forest33/warthog/pkg/structs"
)

// WorkspaceItemFolder stored folder data.
type WorkspaceItemFolder struct {
}

// FolderRequest read/create/update/delete folder request.
type FolderRequest struct {
	ID         int64           `json:"id"`
	ParentID   *int64          `json:"parent_id"`
	Title      string          `json:"title"`
	TypeFilter []WorkspaceType `json:"type_filter"`
}

// FolderResponse read/create/update folder response.
type FolderResponse struct {
	Folder *Workspace           `json:"folder"`
	Tree   []*WorkspaceTreeNode `json:"tree"`
}

// Model creates FolderRequest from UI request.
func (r *FolderRequest) Model(folder map[string]interface{}) error {
	if folder == nil {
		return errors.New("no data")
	}

	if v, ok := folder["id"]; ok && v != nil {
		if id, ok := v.(float64); !ok {
			return errors.New("id not a float")
		} else {
			r.ID = int64(id)
		}
	}
	if v, ok := folder["parent_id"]; ok && v != nil {
		if id, ok := v.(float64); !ok {
			return errors.New("parent id not a float")
		} else {
			r.ParentID = structs.Ref(int64(id))
		}
	}
	if v, ok := folder["title"]; ok && v != nil {
		if r.Title, ok = v.(string); !ok {
			return errors.New("title not a string")
		}
	}
	if v, ok := folder["type_filter"]; ok && v != nil {
		tf, ok := v.([]interface{})
		if !ok {
			return errors.New("type filter not a []interface{}")
		}
		var err error
		r.TypeFilter, err = structs.MapWithError(tf, func(t interface{}) (WorkspaceType, error) {
			s, ok := t.(string)
			if !ok {
				return "", errors.New("type filter value not a string")
			}
			return WorkspaceType(s), nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}
