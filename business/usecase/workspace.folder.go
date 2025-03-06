// Package usecase provides business logic.
package usecase

import (
	"errors"

	"github.com/forest33/warthog/business/entity"
	"github.com/forest33/warthog/pkg/structs"
)

// CreateFolder creates folder on workspace.
func (uc *WorkspaceUseCase) CreateFolder(payload map[string]interface{}) *entity.GUIResponse {
	req := &entity.FolderRequest{}
	if err := req.Model(payload); err != nil {
		return entity.ErrorGUIResponse(err)
	}

	check, err := uc.workspaceRepo.Get(&entity.WorkspaceFilter{
		Title:    req.Title,
		ParentID: req.ParentID,
	})
	if err != nil {
		uc.log.Error().Msgf("failed to get folder: %v", err)
		return entity.ErrorGUIResponse(err)
	} else if check != nil {
		return uc.folderResponse(nil, req.TypeFilter, entity.ErrFolderAlreadyExists)
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

	return uc.folderResponse(folder, req.TypeFilter, nil)
}

// UpdateFolder updates folder on workspace.
func (uc *WorkspaceUseCase) UpdateFolder(payload map[string]interface{}) *entity.GUIResponse {
	req := &entity.FolderRequest{}
	if err := req.Model(payload); err != nil {
		return entity.ErrorGUIResponse(err)
	}

	existsFolder, err := uc.workspaceRepo.GetByID(req.ID)
	if err != nil {
		uc.log.Error().Msgf("failed to get folder: %v", err)
		return entity.ErrorGUIResponse(err)
	} else if existsFolder == nil {
		uc.log.Error().Msgf("failed to get folder: %v", err)
		return entity.ErrorGUIResponse(errors.New("failed to get folder"))
	}

	check, err := uc.workspaceRepo.Get(&entity.WorkspaceFilter{
		Title:    req.Title,
		ParentID: existsFolder.ParentID,
	})
	if err != nil {
		uc.log.Error().Msgf("failed to get folder: %v", err)
		return entity.ErrorGUIResponse(err)
	} else if check != nil {
		return uc.folderResponse(nil, req.TypeFilter, entity.ErrFolderAlreadyExists)
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

	return uc.folderResponse(folder, req.TypeFilter, nil)
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

func (uc *WorkspaceUseCase) folderResponse(folder *entity.Workspace, typeFilter []entity.WorkspaceType, actionErr error) *entity.GUIResponse {
	w, err := uc.workspaceRepo.Get(nil)
	if err != nil {
		uc.log.Error().Msgf("failed to get workspace: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	if actionErr == nil {
		return &entity.GUIResponse{
			Status: entity.GUIResponseStatusOK,
			Payload: &entity.FolderResponse{
				Folder: folder,
				Tree:   entity.MakeWorkspaceTree(w, &entity.WorkspaceTreeFilter{Type: typeFilter}, folder.ID),
			},
		}
	}

	return &entity.GUIResponse{
		Status: structs.If(actionErr == nil, entity.GUIResponseStatusOK, entity.GUIResponseStatusError),
		Payload: &entity.FolderResponse{
			Tree: entity.MakeWorkspaceTree(w, &entity.WorkspaceTreeFilter{Type: typeFilter}, 0),
		},
		Error: entity.ErrorGUIResponse(actionErr).Error,
	}
}
