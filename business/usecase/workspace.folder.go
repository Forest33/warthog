// Package usecase provides business logic.
package usecase

import (
	"github.com/forest33/warthog/business/entity"
)

// CreateFolder creates folder on workspace.
func (uc *WorkspaceUseCase) CreateFolder(payload map[string]interface{}) *entity.GUIResponse {
	req := &entity.FolderRequest{}
	if err := req.Model(payload); err != nil {
		return entity.ErrorGUIResponse(err)
	}

	folder, err := uc.workspaceRepo.Create(&entity.Workspace{
		ParentID: req.ParentID,
		Type:     entity.WorkspaceTypeFolder,
		Title:    req.Title,
	})
	if err != nil {
		uc.log.Error().Msgf("failed to create folder: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	return uc.successFolderResponse(folder, req.TypeFilter)
}

// UpdateFolder updates folder on workspace.
func (uc *WorkspaceUseCase) UpdateFolder(payload map[string]interface{}) *entity.GUIResponse {
	req := &entity.FolderRequest{}
	if err := req.Model(payload); err != nil {
		return entity.ErrorGUIResponse(err)
	}

	folder, err := uc.workspaceRepo.Update(&entity.Workspace{
		ID:       req.ID,
		ParentID: req.ParentID,
		Title:    req.Title,
	})
	if err != nil {
		uc.log.Error().Msgf("failed to update folder: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	return uc.successFolderResponse(folder, req.TypeFilter)
}

// DeleteFolder deletes folder on workspace.
func (uc *WorkspaceUseCase) DeleteFolder(payload map[string]interface{}) *entity.GUIResponse {
	req := &entity.FolderRequest{}
	if err := req.Model(payload); err != nil {
		return entity.ErrorGUIResponse(err)
	}

	if err := uc.workspaceRepo.Delete(req.ID); err != nil {
		uc.log.Error().Msgf("failed to delete folder: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	return &entity.GUIResponse{
		Status: entity.GUIResponseStatusOK,
	}
}

func (uc *WorkspaceUseCase) successFolderResponse(folder *entity.Workspace, typeFilter []entity.WorkspaceType) *entity.GUIResponse {
	w, err := uc.workspaceRepo.Get()
	if err != nil {
		uc.log.Error().Msgf("failed to get workspace: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	return &entity.GUIResponse{
		Status: entity.GUIResponseStatusOK,
		Payload: &entity.FolderResponse{
			Folder: folder,
			Tree:   entity.MakeWorkspaceTree(w, &entity.WorkspaceTreeFilter{Type: typeFilter}, folder.ID),
		},
	}
}
