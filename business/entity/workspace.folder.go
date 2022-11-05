package entity

import (
	"fmt"

	"github.com/Forest33/warthog/pkg/structs"
)

type WorkspaceItemFolder struct {
}

type FolderRequest struct {
	ID         int64           `json:"id"`
	ParentID   *int64          `json:"parent_id"`
	Title      string          `json:"title"`
	TypeFilter []WorkspaceType `json:"type_filter"`
}

type FolderResponse struct {
	Folder *Workspace           `json:"folder"`
	Tree   []*WorkspaceTreeNode `json:"tree"`
}

func (r *FolderRequest) Model(folder map[string]interface{}) error {
	if folder == nil {
		return fmt.Errorf("empty data")
	}

	if v, ok := folder["id"]; ok && v != nil {
		r.ID = int64(v.(float64))
	}
	if v, ok := folder["parent_id"]; ok && v != nil {
		r.ParentID = structs.Ref(int64(v.(float64)))
	}
	if v, ok := folder["title"]; ok && v != nil {
		r.Title = v.(string)
	}
	if v, ok := folder["type_filter"]; ok && v != nil {
		r.TypeFilter = structs.Map(v.([]interface{}), func(t interface{}) WorkspaceType { return WorkspaceType(t.(string)) })
	}

	return nil
}
