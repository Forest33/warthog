package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/forest33/warthog/business/entity"
	"github.com/forest33/warthog/pkg/database"
	"github.com/forest33/warthog/pkg/database/types"
)

const (
	guiConfigTable       = "gui_config"
	guiConfigTableFields = "window_width, window_height, window_x, window_y, created_at, updated_at"
)

// GUIConfigRepository object capable of interacting with GUIConfigRepository
type GUIConfigRepository struct {
	db  *database.Database
	ctx context.Context
}

// NewGUIConfigRepository creates a new GUIConfigRepository
func NewGUIConfigRepository(ctx context.Context, db *database.Database) *GUIConfigRepository {
	return &GUIConfigRepository{
		db:  db,
		ctx: ctx,
	}
}

type guiConfigDTO struct {
	WindowWidth  int    `db:"window_width"`
	WindowHeight int    `db:"window_height"`
	WindowX      int    `db:"window_x"`
	WindowY      int    `db:"window_y"`
	CreatedAt    string `db:"created_at"`
	UpdatedAt    string `db:"updated_at"`
}

func (dto *guiConfigDTO) entity() (*entity.GUIConfig, error) {
	out := &entity.GUIConfig{
		WindowWidth:  dto.WindowWidth,
		WindowHeight: dto.WindowHeight,
		WindowX:      &dto.WindowX,
		WindowY:      &dto.WindowY,
	}

	var err error
	out.CreatedAt, err = types.StrToDateTime(dto.CreatedAt)
	if err != nil {
		return nil, err
	}
	out.UpdatedAt, err = types.StrToDateTime(dto.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// Get returns GUIConfig
func (repo *GUIConfigRepository) Get() (*entity.GUIConfig, error) {
	dto := &guiConfigDTO{}

	err := repo.db.Connector.GetContext(repo.ctx, dto, fmt.Sprintf("SELECT %s FROM %s", guiConfigTableFields, guiConfigTable))
	if err != nil {
		return nil, err
	}

	return dto.entity()
}

// Update updates GUIConfig
func (repo *GUIConfigRepository) Update(in *entity.GUIConfig) (*entity.GUIConfig, error) {
	dto := &guiConfigDTO{}
	attrs := make([]string, 0, 4)
	mapper := make(map[string]interface{}, 4)

	if in.WindowWidth > 0 {
		attrs = append(attrs, "window_width = :window_width")
		mapper["window_width"] = in.WindowWidth
	}
	if in.WindowHeight > 0 {
		attrs = append(attrs, "window_height = :window_height")
		mapper["window_height"] = in.WindowHeight
	}
	if in.WindowX != nil {
		attrs = append(attrs, "window_x = :window_x")
		mapper["window_x"] = in.WindowX
	}
	if in.WindowY != nil {
		attrs = append(attrs, "window_y = :window_y")
		mapper["window_y"] = in.WindowY
	}
	if len(attrs) == 0 {
		return repo.Get()
	}

	query, args, err := repo.db.Connector.BindNamed(fmt.Sprintf(`
			UPDATE %s SET %s, updated_at = datetime('now','localtime')
			RETURNING %s;`, guiConfigTable, strings.Join(attrs, ","), guiConfigTableFields), mapper)
	if err != nil {
		return nil, err
	}
	if err := repo.db.Connector.GetContext(repo.ctx, dto, query, args...); err != nil {
		return nil, err
	}

	return dto.entity()
}
