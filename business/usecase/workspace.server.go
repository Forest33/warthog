package usecase

import (
	"github.com/forest33/warthog/business/entity"
)

func (uc *WorkspaceUseCase) CreateServer(payload map[string]interface{}) *entity.GUIResponse {
	req := &entity.ServerRequest{}
	if err := req.Model(payload); err != nil {
		return entity.ErrorGUIResponse(err)
	}

	server, err := uc.workspaceRepo.Create(&entity.Workspace{
		ParentID: &req.FolderID,
		Type:     entity.WorkspaceTypeServer,
		Title:    req.Title,
		Data:     req.WorkspaceItemServer,
	})
	if err != nil {
		uc.log.Error().Msgf("failed to create server: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	w, err := uc.workspaceRepo.Get()
	if err != nil {
		uc.log.Error().Msgf("failed to get workspace: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	return &entity.GUIResponse{
		Status: entity.GUIResponseStatusOK,
		Payload: &entity.ServerResponse{
			Server: server,
			Tree:   entity.MakeWorkspaceTree(w, nil),
		},
	}
}

func (uc *WorkspaceUseCase) UpdateServer(payload map[string]interface{}) *entity.GUIResponse {
	req := &entity.ServerRequest{}
	if err := req.Model(payload); err != nil {
		return entity.ErrorGUIResponse(err)
	}

	server, err := uc.workspaceRepo.GetByID(req.ID)
	if err != nil {
		return entity.ErrorGUIResponse(err)
	}

	data := server.Data.(*entity.WorkspaceItemServer)
	if data.Request != nil {
		req.WorkspaceItemServer.Request = data.Request
	}

	server, err = uc.workspaceRepo.Update(&entity.Workspace{
		ID:       req.ID,
		ParentID: &req.FolderID,
		Title:    req.Title,
		Data:     req.WorkspaceItemServer,
	})
	if err != nil {
		uc.log.Error().Msgf("failed to update server: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	w, err := uc.workspaceRepo.Get()
	if err != nil {
		uc.log.Error().Msgf("failed to get workspace: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	return &entity.GUIResponse{
		Status: entity.GUIResponseStatusOK,
		Payload: &entity.ServerResponse{
			Server: server,
			Tree:   entity.MakeWorkspaceTree(w, nil),
		},
	}
}

func (uc *WorkspaceUseCase) UpdateServerRequest(payload map[string]interface{}) *entity.GUIResponse {
	req := &entity.ServerUpdateRequest{}
	if err := req.Model(payload); err != nil {
		return entity.ErrorGUIResponse(err)
	}

	server, err := uc.workspaceRepo.GetByID(req.ID)
	if err != nil {
		return entity.ErrorGUIResponse(err)
	}

	data := server.Data.(*entity.WorkspaceItemServer)
	if data.Request != nil {
		if _, ok := data.Request[req.Service]; !ok {
			data.Request[req.Service] = make(map[string]interface{}, 1)
		}
		data.Request[req.Service][req.Method] = req.Request
	} else {
		data.Request = make(map[string]map[string]interface{}, 1)
		data.Request[req.Service] = make(map[string]interface{}, 1)
		data.Request[req.Service][req.Method] = req.Request
	}

	_, err = uc.workspaceRepo.Update(&entity.Workspace{
		ID:   req.ID,
		Data: server.Data,
	})
	if err != nil {
		uc.log.Error().Msgf("failed to update server: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	w, err := uc.workspaceRepo.Get()
	if err != nil {
		uc.log.Error().Msgf("failed to get workspace: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	return &entity.GUIResponse{
		Status: entity.GUIResponseStatusOK,
		Payload: &entity.ServerResponse{
			Server: server,
			Tree:   entity.MakeWorkspaceTree(w, nil),
		},
	}
}
