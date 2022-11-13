// Package main warthog main package
package main

import (
	"flag"
	"runtime"
	"sync"

	"github.com/asticode/go-astilectron"
	"golang.org/x/net/context"

	db "github.com/forest33/warthog/adapter/database"
	"github.com/forest33/warthog/adapter/grpc"
	"github.com/forest33/warthog/business/entity"
	"github.com/forest33/warthog/business/usecase"
	"github.com/forest33/warthog/pkg/logger"
	"github.com/forest33/warthog/pkg/resources"

	"github.com/forest33/warthog/pkg/database"
)

var (
	zlog    *logger.Zerolog
	homeDir string
)

var (
	// AppName application name
	AppName string
	// application version
	AppVersion string
	// AppURL application homepage
	AppURL = "https://github.com/forest33/warthog"
	// BuiltAt build date
	BuiltAt string
	// VersionAstilectron Astilectron version
	VersionAstilectron string
	// VersionElectron Electron version
	VersionElectron string
	// UseBootstrap true if using bootstrap
	UseBootstrap = "false"
)

var (
	cfg = &entity.Config{}
	dbi *database.Database

	guiConfigRepo *db.GUIConfigRepository
	workspaceRepo *db.WorkspaceRepository
	grpcClient    *grpc.Client

	guiConfigUseCase *usecase.GUIConfigUseCase
	workspaceUseCase *usecase.WorkspaceUseCase
	grpcUseCase      *usecase.GrpcUseCase

	guiCfg *entity.GUIConfig
	ast    *astilectron.Astilectron
	window *astilectron.Window
	tray   *astilectron.Tray
	menu   *astilectron.Menu

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
	initUseCases()
	initGUIConfig()

	if UseBootstrap == "true" {
		withBootstrap()
	} else {
		withoutBootstrap()
	}
}

func initAdapters() {
	guiConfigRepo = db.NewGUIConfigRepository(ctx, dbi)
	workspaceRepo = db.NewWorkspaceRepository(ctx, dbi)
	grpcClient = grpc.New(ctx, cfg.Grpc)
}

func initUseCases() {
	guiConfigUseCase = usecase.NewGUIConfigUseCase(ctx, &wg, zlog, guiConfigRepo)
	workspaceUseCase = usecase.NewWorkspaceUseCase(ctx, zlog, workspaceRepo, workspaceID)
	grpcUseCase = usecase.NewGrpcUseCase(ctx, zlog, grpcClient, workspaceRepo)
	usecase.SetWorkspaceUseCase(workspaceUseCase)
}

func shutdown() {
	guiConfigUseCase.Stop()
	wg.Wait()
	cancel()
	dbi.Close()
}
