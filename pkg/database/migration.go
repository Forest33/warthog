package database

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	bin "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

const (
	DefaultMigrationsDir = "migrations"
	MigrationsTable      = "schema_migrations"
	SQLiteDriver         = "sqlite3"
)

type BinDataConfig struct {
	AssetDirFunc AssetDirFunc
	Dir          string
	AssetFunc    bin.AssetFunc
}

type AssetDirFunc func(name string) ([]string, error)

func (db *Database) migrate() (uint, error) {
	if db.binDataConfig == nil {
		return 0, errors.New("nil bin data config")
	}

	var (
		driver database.Driver
		err    error
	)

	switch db.cfg.DriverName {
	case SQLiteDriver:
		driver, err = sqlite3.WithInstance(db.Connector.DB, &sqlite3.Config{
			MigrationsTable: MigrationsTable,
			DatabaseName:    db.cfg.DatasourceName,
		})
	default:
		return 0, errors.New("no available driver for database migrating")
	}

	if err != nil {
		return 0, errors.Wrap(err, "can't create migrations driver")
	}

	defer func() {
		err = multierr.Combine(err)
	}()

	names, err := db.binDataConfig.AssetDirFunc(db.binDataConfig.Dir)
	if err != nil {
		return 0, errors.Wrap(err, "can't get names migrations")
	}

	assetFunc := dirAssetFunc(db.binDataConfig.Dir, db.binDataConfig.AssetFunc)
	sourceInstance, err := bin.WithInstance(bin.Resource(names, assetFunc))
	if err != nil {
		return 0, errors.Wrap(err, "parsing migrations is failed")
	}

	m, err := migrate.NewWithInstance("go-bindata", sourceInstance, db.cfg.DriverName, driver)
	if err != nil {
		return 0, errors.Wrap(err, "failed create migrate instance")
	}

	migratingErr := m.Up()
	version, dirty, err := m.Version()
	if dirty {
		db.log.Error().Msg("dirty migration")
	}
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return 0, errors.Wrap(err, "can't get migrations version")
	}
	if errors.Is(migratingErr, migrate.ErrNoChange) {
		return 0, err
	}

	switch handlingErr := migratingErr.(type) {
	case database.Error:
		err = handlingErr.OrigErr
	case migrate.ErrDirty:
		prevVersion, err := sourceInstance.Prev(uint(handlingErr.Version))
		if err != nil && !os.IsExist(err) {
			return 0, errors.New("can't rollback to previous database version")
		} else if err != nil {
			break
		}
		if prevVersion == 0 {
			return 0, errors.New("can't rollback to previous database version")
		} else {
			err = m.Force(int(prevVersion))
		}
		if err != nil {
			break
		}
		return db.migrate()
	default:
		err = errors.Wrap(err, "database migration failed")
	}

	return version, err
}

func dirAssetFunc(dir string, assetFunc bin.AssetFunc) bin.AssetFunc {
	return func(name string) ([]byte, error) {
		data, err := assetFunc(fmt.Sprintf("%s/%s", dir, name))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return data, nil
	}
}
