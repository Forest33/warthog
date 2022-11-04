package main

import (
	"flag"
	"runtime"

	"github.com/asticode/go-astilectron"
	"golang.org/x/net/context"

	dbadapter "warthog/adapter/database"
	"warthog/adapter/grpc"
	"warthog/business/entity"
	"warthog/business/usecase"
	"warthog/pkg/logger"
	"warthog/pkg/resources"

	"warthog/pkg/database"
)

var (
	zlog    *logger.Zerolog
	homeDir string
)

var (
	AppName            string
	AppVersion         string = "debug"
	BuiltAt            string = "debug"
	VersionAstilectron string = "0.56.0"
	VersionElectron    string = "13.6.9"
	UseBootstrap       string
)

var (
	cfg = &entity.Config{}
	db  *database.Database

	guiConfigRepo *dbadapter.GUIConfigRepository
	workspaceRepo *dbadapter.WorkspaceRepository
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
	guiConfigRepo = dbadapter.NewGUIConfigRepository(ctx, db)
	workspaceRepo = dbadapter.NewWorkspaceRepository(ctx, db)
	grpcClient = grpc.New(ctx, cfg.Grpc)
}

func initUseCases() {
	guiConfigUseCase = usecase.NewGUIConfigUseCase(ctx, zlog, guiConfigRepo)
	workspaceUseCase = usecase.NewWorkspaceUseCase(ctx, zlog, workspaceRepo, workspaceID)
	grpcUseCase = usecase.NewGrpcUseCase(ctx, zlog, grpcClient, workspaceRepo)
	usecase.SetWorkspaceUseCase(workspaceUseCase)
}

func shutdown() {
	cancel()
	db.Close()
}
