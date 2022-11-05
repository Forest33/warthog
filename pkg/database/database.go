package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/Forest33/warthog/business/entity"
	"github.com/Forest33/warthog/pkg/logger"
	"github.com/Forest33/warthog/pkg/resources"
)

type Database struct {
	Connector     *sqlx.DB
	cfg           *entity.DatabaseConfig
	binDataConfig *BinDataConfig
	log           *logger.Zerolog
}

func NewConnector(cfg *entity.DatabaseConfig, binDataConfig *BinDataConfig, log *logger.Zerolog) (*Database, error) {
	connector, err := sqlx.Connect(cfg.DriverName, resources.GetDatabase())
	if err != nil {
		return nil, err
	}

	log.Debug().
		Str("datasource", cfg.DatasourceName).
		Str("driver", cfg.DriverName).
		Msg("initialize database")

	db := &Database{
		cfg:           cfg,
		binDataConfig: binDataConfig,
		log:           log,
		Connector:     connector,
	}

	if _, err := db.migrate(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *Database) Close() {
	if err := db.Connector.Close(); err != nil {
		db.log.Error().Msgf("failed to close database: %v", err)
	}
}
