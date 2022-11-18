// Package entity provides entities for business logic.
package entity

import (
	"os"
)

const (
	debugEnv = "WARTHOG_DEBUG"
)

// Config application configuration
type Config struct {
	Application *ApplicationConfig `json:"application"`
	Database    *DatabaseConfig    `json:"database"`
	Grpc        *GrpcConfig        `json:"grpc"`
	Logger      *LoggerConfig      `json:"logger"`
	Runtime     *RuntimeConfig     `json:"runtime"`
}

// ApplicationConfig base application params
type ApplicationConfig struct {
	Homepage        string `json:"homepage" default:"resources/index.html"`
	HomepageWin     string `json:"homepage_win" default:"../index.html"`
	IconsPath       string `json:"icons_path" default:"resources/icons"`
	AppIconLinux    string `json:"app_icon_linux" default:"app.png"`
	AppIconDarwin   string `json:"app_icon_darwin" default:"warthog.icns"`
	AppIconWindows  string `json:"app_icon_windows" default:"app.ico"`
	TrayIconLinux   string `json:"tray_icon_linux" default:"tray.png"`
	TrayIconDarwin  string `json:"tray_icon_darwin" default:"tray24.png"`
	TrayIconWindows string `json:"tray_icon_windows" default:"tray.ico"`
	SingleInstance  bool   `json:"single_instance" default:"true"`
}

// DatabaseConfig database settings
type DatabaseConfig struct {
	DatasourceName string `json:"datasource_name" default:"warthog.db"`
	DriverName     string `json:"driver_name" default:"sqlite3"`
}

// GrpcConfig gRPC settings
type GrpcConfig struct {
	ConnectTimeout    int  `json:"connect_timeout" default:"10"`
	QueryTimeout      int  `json:"query_timeout" default:"30"`
	NonBlocking       bool `json:"non_blocking" default:"true"`
	SortMethodsByName bool `json:"sort_methods_by_name" default:"true"`
	MaxLoopDepth      int  `json:"max_loop_depth" default:"100"`
}

// LoggerConfig logger settings
type LoggerConfig struct {
	Level             string `json:"level" default:"debug"`
	TimeFieldFormat   string `json:"time_field_format" default:"2006-01-02T15:04:05Z07:00"`
	PrettyPrint       bool   `json:"pretty_print" default:"true"`
	DisableSampling   bool   `json:"disable_sampling" default:"false"`
	RedirectStdLogger bool   `json:"redirect_std_logger" default:"false"`
	ErrorStack        bool   `json:"error_stack" default:"false"`
	ShowCaller        bool   `json:"show_caller" default:"false"`
}

// RuntimeConfig runtime settings
type RuntimeConfig struct {
	GoMaxProcs int `json:"go_max_procs" default:"0"`
}

// IsDebug returns true if application runs on debug mode
func IsDebug() bool {
	return os.Getenv(debugEnv) != ""
}
