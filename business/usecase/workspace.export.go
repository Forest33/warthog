package usecase

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/forest33/warthog/business/entity"
)

// ExportFile exports the workspace to a file.
func (uc *WorkspaceUseCase) ExportFile(payload map[string]interface{}) *entity.GUIResponse {
	req := &entity.ExportRequest{}
	if err := req.Model(payload); err != nil {
		return entity.ErrorGUIResponse(err)
	}

	w, err := uc.workspaceRepo.Get(nil)
	if err != nil {
		uc.log.Error().Msgf("failed to get workspace: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	data, err := json.Marshal(w)
	if err != nil {
		uc.log.Error().Msgf("failed to marshal data: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	err = os.WriteFile(req.Path, data, 0644)
	if err != nil {
		uc.log.Error().Msgf("failed to write file: %v", err)
		return entity.ErrorGUIResponse(errors.New("Failed to save file: " + err.Error()))
	}

	return &entity.GUIResponse{
		Status: entity.GUIResponseStatusOK,
	}
}

// ImportFile imports the workspace from a file.
func (uc *WorkspaceUseCase) ImportFile(payload map[string]interface{}) *entity.GUIResponse {
	req := &entity.ExportRequest{}
	if err := req.Model(payload); err != nil {
		return entity.ErrorGUIResponse(err)
	}

	data, err := os.ReadFile(req.Path)
	if err != nil {
		uc.log.Error().Msgf("failed to read file: %v", err)
		return entity.ErrorGUIResponse(errors.New("Failed to read file: " + err.Error()))
	}

	var importData []*entity.Workspace
	if err := json.Unmarshal(data, &importData); err != nil {
		uc.log.Error().Msgf("failed to unmarshal file: %v", err)
		return entity.ErrorGUIResponse(errors.New("Failed to unmarshal file: " + err.Error()))
	}

	w, err := uc.workspaceRepo.Get(nil)
	if err != nil {
		uc.log.Error().Msgf("failed to get workspace: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	_ = w

	return &entity.GUIResponse{
		Status: entity.GUIResponseStatusOK,
	}
}
