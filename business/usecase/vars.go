package usecase

import (
	"warthog/business/entity"
)

var (
	workspaceUseCase *WorkspaceUseCase
)

type WorkspaceRepo interface {
	Get() ([]*entity.Workspace, error)
	GetByID(id int64) (*entity.Workspace, error)
	Create(in *entity.Workspace) (*entity.Workspace, error)
	Update(in *entity.Workspace) (*entity.Workspace, error)
	Delete(id int64) error
}

func SetWorkspaceUseCase(uc *WorkspaceUseCase) {
	workspaceUseCase = uc
}
