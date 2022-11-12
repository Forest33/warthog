package usecase

import (
	"github.com/forest33/warthog/business/entity"
)

var (
	workspaceUseCase *WorkspaceUseCase
)

// WorkspaceRepo is the common interface implemented WorkspaceRepository methods
type WorkspaceRepo interface {
	Get() ([]*entity.Workspace, error)
	GetByID(id int64) (*entity.Workspace, error)
	Create(in *entity.Workspace) (*entity.Workspace, error)
	Update(in *entity.Workspace) (*entity.Workspace, error)
	Delete(id int64) error
}

// SetWorkspaceUseCase sets WorkspaceUseCase instance
func SetWorkspaceUseCase(uc *WorkspaceUseCase) {
	workspaceUseCase = uc
}
