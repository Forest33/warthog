package usecase

import (
	"github.com/Forest33/warthog/business/entity"
)

func (uc *WorkspaceUseCase) UpdateQuery(payload map[string]interface{}) *entity.GUIResponse {
	req := &entity.QueryRequest{}
	if err := req.Model(payload); err != nil {
		return entity.ErrorGUIResponse(err)
	}

	var (
		query *entity.Workspace
		err   error
	)

	if req.ID != 0 {
		query, err = uc.workspaceRepo.Update(&entity.Workspace{
			ID:    req.ID,
			Title: req.Title,
			Data:  req.WorkspaceItemQuery,
		})
	} else {
		query, err = uc.workspaceRepo.Create(&entity.Workspace{
			ParentID: &req.ServerID,
			Type:     entity.WorkspaceTypeQuery,
			Title:    req.Title,
			Data:     req.WorkspaceItemQuery,
		})
	}
	if err != nil {
		uc.log.Error().Msgf("failed to update query: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	w, err := uc.workspaceRepo.Get()
	if err != nil {
		uc.log.Error().Msgf("failed to get workspace: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	server, err := uc.workspaceRepo.GetByID(*query.ParentID)
	if err != nil {
		uc.log.Error().
			Int64("id", req.ID).
			Msgf("failed to get server: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	server.Breadcrumb, err = uc.GetBreadcrumb(query.ID)
	if err != nil {
		uc.log.Error().Msgf("failed to make breadcrumb: %v", err)
		return entity.ErrorGUIResponse(err)
	}

	return &entity.GUIResponse{
		Status: entity.GUIResponseStatusOK,
		Payload: &entity.ServerResponse{
			Server: server,
			Query:  query,
			Tree:   entity.MakeWorkspaceTree(w, nil),
		},
	}
}
