// Package main warthog main package.
package main

import (
	"context"
	"flag"
	"runtime"
	"sync"

	"github.com/asticode/go-astilectron"

	db "github.com/forest33/warthog/adapter/database"
	"github.com/forest33/warthog/adapter/grpc"
	"github.com/forest33/warthog/adapter/k8s"
	"github.com/forest33/warthog/business/entity"
	"github.com/forest33/warthog/business/usecase"
	"github.com/forest33/warthog/pkg/database"
	"github.com/forest33/warthog/pkg/logger"
	"github.com/forest33/warthog/pkg/resources"
)

//go:generate golangci-lint run -v ../.././...

var (
	zlog    *logger.Zerolog
	homeDir string
)

var (
	// AppName application name.
	AppName string
	// AppVersion application version.
	AppVersion string
	// AppURL application homepage.
	AppURL = "https://github.com/forest33/warthog"
	// BuiltAt build date.
	BuiltAt string
	// VersionAstilectron Astilectron version.
	VersionAstilectron string
	// VersionElectron Electron version.
	VersionElectron string
	// UseBootstrap true if using bootstrap.
	UseBootstrap = "false"
)

var (
	cfg = &entity.Config{}
	dbi *database.Database

	settingsRepo  *db.SettingsRepository
	workspaceRepo *db.WorkspaceRepository
	grpcClient    *grpc.Client
	k8sClient     *k8s.Client

	settingsUseCase  *usecase.SettingsUseCase
	workspaceUseCase *usecase.WorkspaceUseCase
	grpcUseCase      *usecase.GrpcUseCase

	settings *entity.Settings
	ast      *astilectron.Astilectron
	window   *astilectron.Window
	tray     *astilectron.Tray
	menu     *astilectron.Menu

	ctx    context.Context
	cancel context.CancelFunc
	wg     = sync.WaitGroup{}

	workspaceID *int64
)

const (
	applicationName = "Warthog"
)

func init() {
	if cfg.Runtime.GoMaxProcs == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(cfg.Runtime.GoMaxProcs)
	}

	ctx, cancel = context.WithCancel(context.Background())

	workspaceID = flag.Int64("workspace-id", 0, "load workspace with id")
	flag.Parse()
}

func main() {
	defer shutdown()

	zlog = logger.NewZerolog(logger.ZeroConfig{
		Level: func() string {
			if entity.IsDebug() {
				return cfg.Logger.Level
			}
			return "info"
		}(),
		TimeFieldFormat:   cfg.Logger.TimeFieldFormat,
		PrettyPrint:       cfg.Logger.PrettyPrint,
		DisableSampling:   cfg.Logger.DisableSampling,
		RedirectStdLogger: cfg.Logger.RedirectStdLogger,
		ErrorStack:        cfg.Logger.ErrorStack,
		ShowCaller:        cfg.Logger.ShowCaller,
	})

	if workspaceID != nil {
		zlog.Debug().Int64("id", *workspaceID).Msg("use workspace")
	}

	resources.Init(cfg, zlog)
	homeDir = resources.CreateApplicationDir()

	initDatabase()
	initAdapters()
	initClients()
	initUseCases()

	if UseBootstrap == "true" {
		withBootstrap()
	} else {
		withoutBootstrap()
	}
}

func initAdapters() {
	settingsRepo = db.NewSettingsRepository(ctx, dbi)
	workspaceRepo = db.NewWorkspaceRepository(ctx, dbi, zlog)
}

func initClients() {
	grpcClient = grpc.New(ctx, zlog)
	k8sClient = k8s.New(ctx, zlog)
}

func initUseCases() {
	settingsUseCase = usecase.NewSettingsUseCase(ctx, &wg, zlog, settingsRepo, grpcClient)

	settings := initSettings()
	grpcClient.SetSettings(settings)
	k8sClient.SetSettings(settings)

	workspaceUseCase = usecase.NewWorkspaceUseCase(ctx, zlog, workspaceRepo, workspaceID)
	usecase.SetWorkspaceUseCase(workspaceUseCase)

	grpcUseCase = usecase.NewGrpcUseCase(ctx, zlog, grpcClient, k8sClient, workspaceRepo)
}

func initSettings() *entity.Settings {
	settings = entity.DefaultSettings
	if s, err := settingsUseCase.Get(); err == nil {
		settings = s
	}
	return settings
}

func shutdown() {
	settingsUseCase.Stop()
	wg.Wait()
	cancel()
	dbi.Close()
}
