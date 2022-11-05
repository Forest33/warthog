package usecase

import (
	"context"
	"errors"

	"github.com/forest33/warthog/business/entity"
	"github.com/forest33/warthog/pkg/logger"
)

type WorkspaceUseCase struct {
	ctx                context.Context
	log                *logger.Zerolog
	workspaceRepo      WorkspaceRepo
	startupWorkspaceID *int64
}

func NewWorkspaceUseCase(ctx context.Context, log *logger.Zerolog, workspaceRepo WorkspaceRepo, startupWorkspaceID *int64) *WorkspaceUseCase {
	uc := &WorkspaceUseCase{
		ctx:                ctx,
		log:                log,
		workspaceRepo:      workspaceRepo,
		startupWorkspaceID: startupWorkspaceID,
	}

	return uc
}

func (uc *WorkspaceUseCase) Get(payload interface{}) *entity.GUIResponse {
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

	tree := entity.MakeWorkspaceTree(w, filter)

	return &entity.GUIResponse{
		Status:  entity.GUIResponseStatusOK,
		Payload: tree,
	}
}

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

func (uc *WorkspaceUseCase) GetState() *entity.GUIResponse {
	workspaces, err := uc.workspaceRepo.Get()
	if err != nil {
		uc.log.Error().Msgf("failed to get workspace: %v", err)
		return entity.ErrorGUIResponse(err)
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

	return &entity.GUIResponse{
		Status:  entity.GUIResponseStatusOK,
		Payload: state,
	}
}

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

func (uc *WorkspaceUseCase) GetBreadcrumb(id int64) ([]string, error) {
	w, err := uc.workspaceRepo.Get()
	if err != nil {
		return nil, err
	}
	return entity.GetBreadcrumb(w, id), nil
}
