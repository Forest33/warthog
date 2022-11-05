package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/forest33/warthog/business/entity"
	"github.com/forest33/warthog/pkg/database"
	"github.com/forest33/warthog/pkg/database/types"
	"github.com/forest33/warthog/pkg/structs"
)

const (
	workspaceTable       = "workspace"
	workspaceTableFields = "id, parent_id, has_child, type, title, data, sort, expanded, created_at, updated_at"
)

type WorkspaceRepository struct {
	db  *database.Database
	ctx context.Context
}

func NewWorkspaceRepository(ctx context.Context, db *database.Database) *WorkspaceRepository {
	return &WorkspaceRepository{
		db:  db,
		ctx: ctx,
	}
}

type workspaceDTO struct {
	ID        int64          `db:"id"`
	ParentID  sql.NullInt64  `db:"parent_id"`
	HasChild  bool           `db:"has_child"`
	Type      string         `db:"type"`
	Title     string         `db:"title"`
	Data      sql.NullString `db:"data"`
	Sort      int64          `db:"sort"`
	Expanded  bool           `db:"expanded"`
	CreatedAt string         `db:"created_at"`
	UpdatedAt string         `db:"updated_at"`
}

func newWorkspaceDTO(in *entity.Workspace) (dto *workspaceDTO, err error) {
	dto = &workspaceDTO{
		ParentID: types.RefInt64ToSQL(in.ParentID),
		Type:     in.Type.String(),
		Title:    in.Title,
	}

	if in.Data != nil {
		data, err := json.Marshal(in.Data)
		if err != nil {
			return nil, err
		}
		dto.Data = types.StringToSQL(string(data))
	}

	return
}

func (dto *workspaceDTO) entity() (*entity.Workspace, error) {
	out := &entity.Workspace{
		ID:       dto.ID,
		ParentID: types.SQLToRefInt64(dto.ParentID),
		HasChild: &dto.HasChild,
		Type:     entity.WorkspaceType(dto.Type),
		Title:    dto.Title,
		Expanded: &dto.Expanded,
	}

	var err error
	if dto.Data.Valid {
		switch out.Type {
		case entity.WorkspaceTypeFolder:
			out.Data = &entity.WorkspaceItemFolder{}
		case entity.WorkspaceTypeServer:
			out.Data = &entity.WorkspaceItemServer{}
		case entity.WorkspaceTypeQuery:
			out.Data = &entity.WorkspaceItemQuery{}
		default:
			return nil, fmt.Errorf("unknown message type: %v", out.Type)
		}
		if err := json.Unmarshal([]byte(dto.Data.String), &out.Data); err != nil {
			return nil, err
		}
	}
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

func (repo *WorkspaceRepository) GetByID(id int64) (*entity.Workspace, error) {
	dto := &workspaceDTO{ID: id}

	rows, err := repo.db.Connector.NamedQueryContext(repo.ctx, fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE id = :id`, workspaceTableFields, workspaceTable), dto)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	if rows.Next() {
		err := rows.StructScan(&dto)
		if err != nil {
			return nil, err
		}
	}

	return dto.entity()
}

func (repo *WorkspaceRepository) GetByParentID(parentID int64, tx *sqlx.Tx) ([]*entity.Workspace, error) {
	dto := &workspaceDTO{ParentID: types.Int64ToSQL(parentID)}
	res := make([]*entity.Workspace, 0, 10)

	type idb interface {
		NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	}

	var db idb = repo.db.Connector
	if tx != nil {
		db = tx
	}

	rows, err := db.NamedQuery(fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE parent_id = :parent_id`, workspaceTableFields, workspaceTable), dto)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		err := rows.StructScan(&dto)
		if err != nil {
			return nil, err
		}
		w, err := dto.entity()
		if err != nil {
			return nil, err
		}
		res = append(res, w)
	}

	return res, nil
}

func (repo *WorkspaceRepository) Get() ([]*entity.Workspace, error) {
	var dto []*workspaceDTO

	err := repo.db.Connector.SelectContext(repo.ctx, &dto, fmt.Sprintf(`
		SELECT %s 
		FROM %s
		ORDER BY type, sort;`, workspaceTableFields, workspaceTable))
	if err != nil {
		return nil, err
	}

	return structs.MapWithError(dto, func(w *workspaceDTO) (*entity.Workspace, error) { return w.entity() })
}

func (repo *WorkspaceRepository) Create(in *entity.Workspace) (*entity.Workspace, error) {
	dto, err := newWorkspaceDTO(in)
	if err != nil {
		return nil, err
	}

	var (
		query string
		args  []interface{}
	)

	tx := repo.db.Connector.MustBegin()
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	query, args, err = tx.BindNamed(fmt.Sprintf(`
			INSERT INTO %s (parent_id, has_child, type, title, data)
			VALUES (:parent_id, :has_child, :type, :title, :data)
			RETURNING %s;`, workspaceTable, workspaceTableFields), dto)
	if err != nil {
		return nil, err
	}
	if err := tx.GetContext(repo.ctx, dto, query, args...); err != nil {
		return nil, err
	}

	if in.ParentID != nil {
		_, err = tx.NamedExecContext(repo.ctx, fmt.Sprintf(`
			UPDATE %s SET has_child = TRUE, updated_at = datetime('now','localtime')
			WHERE id = :id`, workspaceTable), &workspaceDTO{ID: *in.ParentID})
		if err != nil {
			return nil, err
		}
	}

	return dto.entity()
}

func (repo *WorkspaceRepository) Update(in *entity.Workspace) (*entity.Workspace, error) {
	dto := &workspaceDTO{}
	attrs := make([]string, 0, 6)
	mapper := make(map[string]interface{}, 7)

	if in.ParentID != nil {
		attrs = append(attrs, "parent_id = :parent_id")
		mapper["parent_id"] = *in.ParentID
	}
	if in.HasChild != nil {
		attrs = append(attrs, "has_child = :has_child")
		mapper["has_child"] = *in.HasChild
	}
	if len(in.Title) > 0 {
		attrs = append(attrs, "title = :title")
		mapper["title"] = in.Title
	}
	if in.Data != nil {
		data, err := json.Marshal(in.Data)
		if err != nil {
			return nil, err
		}
		attrs = append(attrs, "data = :data")
		mapper["data"] = types.StringToSQL(string(data))
	}
	if in.Sort != nil {
		attrs = append(attrs, "sort = :sort")
		mapper["sort"] = *in.Sort
	}
	if in.Expanded != nil {
		attrs = append(attrs, "expanded = :expanded")
		mapper["expanded"] = *in.Expanded
	}
	if len(attrs) == 0 {
		return repo.GetByID(in.ID)
	}

	attrs = append(attrs, "id = :id")
	mapper["id"] = in.ID

	var (
		query string
		args  []interface{}
		err   error
	)

	tx := repo.db.Connector.MustBegin()
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	query, args, err = tx.BindNamed(fmt.Sprintf(`
			UPDATE %s SET %s, updated_at = datetime('now','localtime')
			WHERE id = :id
			RETURNING %s;`, workspaceTable, strings.Join(attrs, ","), workspaceTableFields), mapper)
	if err != nil {
		return nil, err
	}
	if err := tx.GetContext(repo.ctx, dto, query, args...); err != nil {
		return nil, err
	}

	if in.ParentID != nil {
		_, err = tx.NamedExecContext(repo.ctx, fmt.Sprintf(`
			UPDATE %s SET has_child = TRUE, updated_at = datetime('now','localtime')
			WHERE id = :id`, workspaceTable), &workspaceDTO{ID: *in.ParentID})
		if err != nil {
			return nil, err
		}
	}

	return dto.entity()
}

func (repo *WorkspaceRepository) Delete(id int64) error {
	workspace, err := repo.GetByID(id)
	if err != nil {
		return err
	}

	if *workspace.HasChild {
		return errors.New("workspace has a child")
	}

	tx := repo.db.Connector.MustBegin()
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	_, err = tx.NamedExecContext(repo.ctx, fmt.Sprintf(`
			DELETE FROM %s 
			WHERE id = :id;`, workspaceTable), &workspaceDTO{ID: id})
	if err != nil {
		return err
	}

	if workspace.ParentID != nil {
		var child []*entity.Workspace
		child, err = repo.GetByParentID(*workspace.ParentID, tx)
		if err != nil {
			return err
		}
		if len(child) == 0 {
			_, err = tx.NamedExecContext(repo.ctx, fmt.Sprintf(`
					UPDATE %s SET has_child = FALSE, updated_at = datetime('now','localtime')
					WHERE id = :id`, workspaceTable), &workspaceDTO{ID: *workspace.ParentID})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
