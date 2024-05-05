// Package database provides CRUD operations with database.
package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/forest33/warthog/business/entity"
	"github.com/forest33/warthog/pkg/database"
)

const (
	settingsTable       = "settings"
	settingsTableFields = `window_width, window_height, window_x, window_y, single_instance, connect_timeout,
							request_timeout, k8s_request_timeout, non_blocking_connection, sort_methods_by_name, max_loop_depth, 
							emit_defaults, check_updates`
)

// SettingsRepository object capable of interacting with SettingsRepository.
type SettingsRepository struct {
	db  *database.Database
	ctx context.Context
}

// NewSettingsRepository creates a new SettingsRepository.
func NewSettingsRepository(ctx context.Context, db *database.Database) *SettingsRepository {
	return &SettingsRepository{
		db:  db,
		ctx: ctx,
	}
}

type settingsDTO struct {
	WindowWidth           int  `db:"window_width"`
	WindowHeight          int  `db:"window_height"`
	WindowX               int  `db:"window_x"`
	WindowY               int  `db:"window_y"`
	SingleInstance        bool `db:"single_instance"`
	ConnectTimeout        int  `db:"connect_timeout"`
	RequestTimeout        int  `db:"request_timeout"`
	K8SRequestTimeout     int  `db:"k8s_request_timeout"`
	NonBlockingConnection bool `db:"non_blocking_connection"`
	SortMethodsByName     bool `db:"sort_methods_by_name"`
	MaxLoopDepth          int  `db:"max_loop_depth"`
	EmitDefaults          bool `db:"emit_defaults"`
	CheckUpdates          bool `db:"check_updates"`
}

func (dto *settingsDTO) entity() *entity.Settings {
	return &entity.Settings{
		WindowWidth:           dto.WindowWidth,
		WindowHeight:          dto.WindowHeight,
		WindowX:               &dto.WindowX,
		WindowY:               &dto.WindowY,
		SingleInstance:        &dto.SingleInstance,
		ConnectTimeout:        &dto.ConnectTimeout,
		RequestTimeout:        &dto.RequestTimeout,
		K8SRequestTimeout:     &dto.K8SRequestTimeout,
		NonBlockingConnection: &dto.NonBlockingConnection,
		SortMethodsByName:     &dto.SortMethodsByName,
		MaxLoopDepth:          &dto.MaxLoopDepth,
		EmitDefaults:          &dto.EmitDefaults,
		CheckUpdates:          &dto.CheckUpdates,
	}
}

// Get returns Settings.
func (repo *SettingsRepository) Get() (*entity.Settings, error) {
	dto := &settingsDTO{}

	err := repo.db.Connector.GetContext(repo.ctx, dto, fmt.Sprintf("SELECT %s FROM %s", settingsTableFields, settingsTable))
	if err != nil {
		return nil, err
	}

	return dto.entity(), nil
}

// Update updates Settings.
func (repo *SettingsRepository) Update(in *entity.Settings) (*entity.Settings, error) {
	dto := &settingsDTO{}
	attrs := make([]string, 0, 13)
	mapper := make(map[string]interface{}, 13)

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
	if in.SingleInstance != nil {
		attrs = append(attrs, "single_instance = :single_instance")
		mapper["single_instance"] = in.SingleInstance
	}
	if in.ConnectTimeout != nil {
		attrs = append(attrs, "connect_timeout = :connect_timeout")
		mapper["connect_timeout"] = in.ConnectTimeout
	}
	if in.RequestTimeout != nil {
		attrs = append(attrs, "request_timeout = :request_timeout")
		mapper["request_timeout"] = in.RequestTimeout
	}
	if in.K8SRequestTimeout != nil {
		attrs = append(attrs, "k8s_request_timeout = :k8s_request_timeout")
		mapper["k8s_request_timeout"] = in.K8SRequestTimeout
	}
	if in.NonBlockingConnection != nil {
		attrs = append(attrs, "non_blocking_connection = :non_blocking_connection")
		mapper["non_blocking_connection"] = in.NonBlockingConnection
	}
	if in.SortMethodsByName != nil {
		attrs = append(attrs, "sort_methods_by_name = :sort_methods_by_name")
		mapper["sort_methods_by_name"] = in.SortMethodsByName
	}
	if in.MaxLoopDepth != nil {
		attrs = append(attrs, "max_loop_depth = :max_loop_depth")
		mapper["max_loop_depth"] = in.MaxLoopDepth
	}
	if in.EmitDefaults != nil {
		attrs = append(attrs, "emit_defaults = :emit_defaults")
		mapper["emit_defaults"] = in.EmitDefaults
	}
	if in.CheckUpdates != nil {
		attrs = append(attrs, "check_updates = :check_updates")
		mapper["check_updates"] = in.CheckUpdates
	}
	if len(attrs) == 0 {
		return repo.Get()
	}

	query, args, err := repo.db.Connector.BindNamed(fmt.Sprintf(`
			UPDATE %s SET %s, updated_at = datetime('now','localtime')
			RETURNING %s;`, settingsTable, strings.Join(attrs, ","), settingsTableFields), mapper)
	if err != nil {
		return nil, err
	}
	if err := repo.db.Connector.GetContext(repo.ctx, dto, query, args...); err != nil {
		return nil, err
	}

	return dto.entity(), nil
}
