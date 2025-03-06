// Package entity provides entities for business logic.
package entity

import (
	"errors"
	"time"

	"github.com/forest33/warthog/pkg/structs"
)

// workspace types.
const (
	WorkspaceTypeFolder WorkspaceType = "f"
	WorkspaceTypeServer WorkspaceType = "s"
	WorkspaceTypeQuery  WorkspaceType = "r"

	WorkspaceDuplicatePostfix = "Copy"

	WorkspaceEventServerUpdated = "server.updated"
)

var (
	// ErrWorkspaceNotExists error workspace not exists.
	ErrWorkspaceNotExists = errors.New("workspace not exists")
)

// Workspace workspace item.
type Workspace struct {
	ID         int64         `json:"id"`
	ParentID   *int64        `json:"parent_id"`
	HasChild   *bool         `json:"has_child"`
	Type       WorkspaceType `json:"type"`
	Title      string        `json:"text"`
	Breadcrumb []string      `json:"breadcrumb"`
	Data       interface{}   `json:"data"`
	Sort       *int          `json:"sort"`
	Expanded   *bool         `json:"expanded"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

// WorkspaceType workspace type.
type WorkspaceType string

// String returns workspace type string.
func (t WorkspaceType) String() string {
	return string(t)
}

// WorkspaceEvent workspace event.
type WorkspaceEvent string

// String returns workspace event string.
func (e WorkspaceEvent) String() string {
	return string(e)
}

// WorkspaceRequest workspace request by type.
type WorkspaceRequest struct {
	Type       []WorkspaceType `json:"type"`
	SelectedID int64           `json:"selected_id"`
}

// WorkspaceFilter filter for workspace search.
type WorkspaceFilter struct {
	Title    string
	ParentID *int64
}

// Model creates WorkspaceRequest from UI request.
func (r *WorkspaceRequest) Model(payload map[string]interface{}) error {
	if payload == nil {
		return nil
	}

	if v, ok := payload["type"]; ok && v != nil {
		t, ok := v.([]interface{})
		if !ok {
			return errors.New("type not a []interface{}")
		}
		var err error
		r.Type, err = structs.MapWithError(t, func(t interface{}) (WorkspaceType, error) {
			s, ok := t.(string)
			if !ok {
				return "", errors.New("type value not a string")
			}
			return WorkspaceType(s), nil
		})
		if err != nil {
			return err
		}
	}
	if v, ok := payload["selected_id"]; ok && v != nil {
		f, ok := v.(float64)
		if !ok {
			return errors.New("selected id not a float")
		}
		r.SelectedID = int64(f)
	}

	return nil
}

// WorkspaceSortingRequest workspace sorting request.
type WorkspaceSortingRequest struct {
	Nodes []*Workspace `json:"nodes"`
}

// Model creates WorkspaceSortingRequest from UI request.
func (r *WorkspaceSortingRequest) Model(payload map[string]interface{}) error {
	if payload == nil {
		return errors.New("no nodes")
	}
	if _, ok := payload["nodes"]; !ok {
		return errors.New("no nodes")
	}

	nodes, ok := payload["nodes"].([]interface{})
	if !ok {
		return errors.New("nodes not a []interface{}")
	}

	r.Nodes = make([]*Workspace, len(nodes))
	for i, n := range nodes {
		var (
			id       int64
			parentID *int64
		)
		node, ok := n.(map[string]interface{})
		if !ok {
			return errors.New("node not a map[string]interface{}")
		}
		if v, ok := node["id"]; ok && v != nil {
			f, ok := v.(float64)
			if !ok {
				return errors.New("node id not a float")
			}
			id = int64(f)
		} else {
			return errors.New("id not exists")
		}
		if v, ok := node["parent_id"]; ok && v != nil {
			f, ok := v.(float64)
			if !ok {
				return errors.New("node parent id not a float")
			}
			parentID = structs.Ref(int64(f))
		}
		r.Nodes[i] = &Workspace{ID: id, ParentID: parentID}
	}

	return nil
}

// WorkspaceExpandRequest workspace expand/collapse request.
type WorkspaceExpandRequest struct {
	ID     int64 `json:"id"`
	Expand bool  `json:"expand"`
}

// Model creates WorkspaceExpandRequest from UI request.
func (r *WorkspaceExpandRequest) Model(payload map[string]interface{}) error {
	if payload == nil {
		return errors.New("no data")
	}

	if v, ok := payload["id"]; ok && v != nil {
		f, ok := v.(float64)
		if !ok {
			return errors.New("id not a float")
		}
		r.ID = int64(f)
	}
	if v, ok := payload["expand"]; ok && v != nil {
		if r.Expand, ok = v.(bool); !ok {
			return errors.New("expand not a bool")
		}
	}

	return nil
}

// WorkspaceTreeFilter filtering workspace by type.
type WorkspaceTreeFilter struct {
	Type []WorkspaceType
}

// WorkspaceTreeNode workspace tree node.
type WorkspaceTreeNode struct {
	Data  *Workspace           `json:"data"`
	Text  string               `json:"text"`
	Nodes []*WorkspaceTreeNode `json:"nodes"`
}

// GetBreadcrumb returns breadcrumb.
func GetBreadcrumb(w []*Workspace, id int64) []string {
	nodeMap := structs.SliceToMap(w, func(w *Workspace) int64 { return w.ID })
	return makeBreadcrumb(nodeMap, id, []string{})
}

func makeBreadcrumb(nodeMap map[int64]*Workspace, id int64, breadcrumb []string) []string {
	if _, ok := nodeMap[id]; !ok {
		return breadcrumb
	}
	breadcrumb = append([]string{nodeMap[id].Title}, breadcrumb...)
	if nodeMap[id].ParentID == nil {
		return breadcrumb
	}
	return makeBreadcrumb(nodeMap, *nodeMap[id].ParentID, breadcrumb)
}

func getExpandedNodes(w []*Workspace, nodeMap map[int64]int, id int64) map[int64]struct{} {
	if _, ok := nodeMap[id]; !ok {
		return map[int64]struct{}{}
	}

	parentNodes := make(map[int64]struct{}, 1)
	parent := w[nodeMap[id]].ParentID
	for parent != nil {
		parentNodes[*parent] = struct{}{}
		parent = w[nodeMap[*parent]].ParentID
	}

	return parentNodes
}

// MakeWorkspaceTree creates workspace tree for UI.
func MakeWorkspaceTree(w []*Workspace, filter *WorkspaceTreeFilter, selectedID int64) []*WorkspaceTreeNode {
	nodeMap := make(map[int64]int, len(w))
	expandedNodes := make(map[int64]struct{}, len(w))
	list := make([]*WorkspaceTreeNode, len(w))
	tree := make([]*WorkspaceTreeNode, 0, len(w))
	onlyTypes := map[WorkspaceType]struct{}{}
	if filter != nil {
		onlyTypes = structs.SliceToMapOfStruct(filter.Type, func(t WorkspaceType) WorkspaceType { return t })
	}

	for i, item := range w {
		if len(onlyTypes) > 0 {
			if _, ok := onlyTypes[item.Type]; !ok {
				continue
			}
		}

		nodeMap[item.ID] = i
		item.Breadcrumb = []string{item.Title}
		list[i] = &WorkspaceTreeNode{
			Data: item,
			Text: item.Title,
		}
	}

	if selectedID != 0 {
		expandedNodes = getExpandedNodes(w, nodeMap, selectedID)
	}

	for _, item := range list {
		if item == nil {
			continue
		}
		if _, ok := expandedNodes[item.Data.ID]; ok {
			item.Data.Expanded = structs.Ref(true)
		}
		if item.Data.ParentID != nil && *item.Data.ParentID != 0 {
			item.Data.Breadcrumb = append([]string{list[nodeMap[*item.Data.ParentID]].Text}, item.Data.Breadcrumb...)
			if list[nodeMap[*item.Data.ParentID]].Nodes == nil {
				list[nodeMap[*item.Data.ParentID]].Nodes = make([]*WorkspaceTreeNode, 0, 10)
			}
			list[nodeMap[*item.Data.ParentID]].Nodes = append(list[nodeMap[*item.Data.ParentID]].Nodes, item)
		} else {
			tree = append(tree, item)
		}
	}

	return tree
}

// WorkspaceState count of folders/servers/queries.
type WorkspaceState struct {
	Folders            int    `json:"folders"`
	Servers            int    `json:"servers"`
	Queries            int    `json:"queries"`
	StartupWorkspaceID *int64 `json:"startup_workspace_id"`
}
