// Package usecase provides business logic.
package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/forest33/warthog/business/entity"
	"github.com/forest33/warthog/pkg/logger"
)

// WorkspaceUseCase object capable of interacting with WorkspaceUseCase
type WorkspaceUseCase struct {
	ctx                context.Context
	log                *logger.Zerolog
	workspaceRepo      WorkspaceRepo
	startupWorkspaceID *int64
	subscribers        []func(e entity.WorkspaceEvent, payload interface{})
}

// NewWorkspaceUseCase creates a new WorkspaceUseCase
func NewWorkspaceUseCase(ctx context.Context, log *logger.Zerolog, workspaceRepo WorkspaceRepo, startupWorkspaceID *int64) *WorkspaceUseCase {
	uc := &WorkspaceUseCase{
		ctx:                ctx,
		log:                log,
		workspaceRepo:      workspaceRepo,
		startupWorkspaceID: startupWorkspaceID,
	}

	return uc
}

// Get returns workspace tree
func (uc *WorkspaceUseCase) Get(payload map[string]interface{}) *entity.GUIResponse {
	req := &entity.WorkspaceRequest{}
	if err := req.Model(payload); err != nil {
		return entity.ErrorGUIResponse(err)
	}

	filter := &entity.WorkspaceTreeFilter{}
	if req.Type != nil {
		filter.Type = req.Type
	}

	w, err := uc.workspaceRepo.Get()
	if err != nil {
		uc.log.Error().Msgf("failed to get workspace: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	tree := entity.MakeWorkspaceTree(w, filter, req.SelectedID)

	return &entity.GUIResponse{
		Status:  entity.GUIResponseStatusOK,
		Payload: tree,
	}
}

// Sorting gets the sorted workspace tree and stores it in the database
func (uc *WorkspaceUseCase) Sorting(payload map[string]interface{}) *entity.GUIResponse {
	req := &entity.WorkspaceSortingRequest{}
	if err := req.Model(payload); err != nil {
		return entity.ErrorGUIResponse(err)
	}

	for i, w := range req.Nodes {
		w.Sort = &i
		_, err := uc.workspaceRepo.Update(w)
		if err != nil {
			uc.log.Error().Msgf("failed to update workspace: %v", err)
			return entity.ErrorGUIResponse(err)
		}
	}

	return uc.Get(nil)
}

// Expand stores expand/collapse status on database
func (uc *WorkspaceUseCase) Expand(payload map[string]interface{}) *entity.GUIResponse {
	req := &entity.WorkspaceExpandRequest{}
	if err := req.Model(payload); err != nil {
		return entity.ErrorGUIResponse(err)
	}

	_, err := uc.workspaceRepo.Update(&entity.Workspace{
		ID:       req.ID,
		Expanded: &req.Expand,
	})
	if err != nil {
		uc.log.Error().Msgf("failed to update workspace: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	return &entity.GUIResponse{
		Status: entity.GUIResponseStatusOK,
	}
}

// GetState returns count of folders/servers/queries
func (uc *WorkspaceUseCase) GetState() (*entity.WorkspaceState, error) {
	workspaces, err := uc.workspaceRepo.Get()
	if err != nil {
		uc.log.Error().Msgf("failed to get workspace: %v", err)
		return nil, err
	}

	state := &entity.WorkspaceState{StartupWorkspaceID: uc.startupWorkspaceID}
	for _, w := range workspaces {
		switch w.Type {
		case entity.WorkspaceTypeFolder:
			state.Folders++
		case entity.WorkspaceTypeServer:
			state.Servers++
		case entity.WorkspaceTypeQuery:
			state.Queries++
		}
	}

	return state, nil
}

// Delete deletes workspace item
func (uc *WorkspaceUseCase) Delete(payload map[string]interface{}) *entity.GUIResponse {
	if payload == nil {
		return entity.ErrorGUIResponse(errors.New("nil payload"))
	}
	if _, ok := payload["id"]; !ok {
		return entity.ErrorGUIResponse(errors.New("no workspace id"))
	}

	if err := uc.workspaceRepo.Delete(int64(payload["id"].(float64))); err != nil {
		return entity.ErrorGUIResponse(err)
	}

	return uc.Get(nil)
}

// Duplicate duplicates workspace item
func (uc *WorkspaceUseCase) Duplicate(payload map[string]interface{}) *entity.GUIResponse {
	if payload == nil {
		return entity.ErrorGUIResponse(errors.New("nil payload"))
	}
	if _, ok := payload["id"]; !ok {
		return entity.ErrorGUIResponse(errors.New("no workspace id"))
	}

	item, err := uc.workspaceRepo.GetByID(int64(payload["id"].(float64)))
	if err != nil {
		return entity.ErrorGUIResponse(errors.New("failed to get workspace item"))
	}

	if item.Type != entity.WorkspaceTypeServer && item.Type != entity.WorkspaceTypeQuery {
		return entity.ErrorGUIResponse(errors.New("wrong workspace item type"))
	}

	item.Title = fmt.Sprintf("%s %s", item.Title, entity.WorkspaceDuplicatePostfix)

	if _, err := uc.workspaceRepo.Create(item); err != nil {
		return entity.ErrorGUIResponse(err)
	}

	return uc.Get(nil)
}

// GetBreadcrumb returns the breadcrumbs
func (uc *WorkspaceUseCase) GetBreadcrumb(id int64) ([]string, error) {
	w, err := uc.workspaceRepo.Get()
	if err != nil {
		return nil, err
	}
	return entity.GetBreadcrumb(w, id), nil
}

// Subscribe event subscription
func (uc *WorkspaceUseCase) Subscribe(handler func(e entity.WorkspaceEvent, payload interface{})) {
	if uc.subscribers == nil {
		uc.subscribers = make([]func(e entity.WorkspaceEvent, payload interface{}), 0, 1)
	}
	uc.subscribers = append(uc.subscribers, handler)
}

// Publish sending an event
func (uc *WorkspaceUseCase) Publish(e entity.WorkspaceEvent, payload interface{}) {
	if uc.subscribers == nil {
		return
	}
	for _, s := range uc.subscribers {
		go s(e, payload)
	}
}
