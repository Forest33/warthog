package main

import (
	"github.com/forest33/warthog/deploy/app/migrations"
	"github.com/forest33/warthog/pkg/database"
)

//go:generate go-bindata -o ./migrations/migrations.bindata.go -pkg migrations -ignore=\\*.go ./migrations/...

func initDatabase() {
	binDataCfg := &database.BinDataConfig{
		Dir:          database.DefaultMigrationsDir,
		AssetDirFunc: migrations.AssetDir,
		AssetFunc:    migrations.Asset,
	}

	var err error
	dbi, err = database.NewConnector(cfg.Database, binDataCfg, zlog)
	if err != nil {
		zlog.Fatal(err.Error())
	}
}
