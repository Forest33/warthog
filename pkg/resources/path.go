// Package resources provides basic operations with application resources
package resources

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/forest33/warthog/business/entity"
	"github.com/forest33/warthog/pkg/logger"
)

const (
	applicationDir        = ".warthog"
	applicationWindowsDir = "Warthog"
)

var (
	cfg     *entity.Config
	log     *logger.Zerolog
	homeDir string
)

// Init initialize package
func Init(c *entity.Config, l *logger.Zerolog) {
	cfg = c
	log = l
}

// CreateApplicationDir creates application folder
func CreateApplicationDir() string {
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to get home directory: %v", err)
	}

	homeDir = filepath.Join(userHome, applicationDir)
	if runtime.GOOS == "windows" {
		homeDir = filepath.Join(userHome, applicationWindowsDir)
	}

	if _, err := os.Stat(homeDir); os.IsNotExist(err) {
		_ = os.MkdirAll(homeDir, 0755)
	}

	return homeDir
}

// GetApplicationIcon returns path to application icon
func GetApplicationIcon() string {
	if runtime.GOOS == "darwin" && entity.IsDebug() {
		return ""
	} else if runtime.GOOS == "windows" {
		if dir, err := os.UserConfigDir(); err == nil {
			return filepath.Join(dir, applicationWindowsDir, cfg.Application.IconsPath, cfg.Application.AppIconWindows)
		}
	}
	return getPath(cfg.Application.IconsPath, cfg.Application.AppIconLinux, cfg.Application.AppIconDarwin, cfg.Application.AppIconWindows)
}

// GetTrayIcon returns path to tray icon
func GetTrayIcon() string {
	if runtime.GOOS != "windows" {
		return getPath(cfg.Application.IconsPath, cfg.Application.TrayIconLinux, cfg.Application.TrayIconDarwin, cfg.Application.TrayIconWindows)
	}
	if dir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(dir, applicationWindowsDir, cfg.Application.IconsPath, cfg.Application.TrayIconWindows)
	}
	return cfg.Application.TrayIconWindows
}

// GetHomepage returns path to application homepage
func GetHomepage() string {
	if !entity.IsDebug() {
		if runtime.GOOS != "windows" {
			return filepath.Join(homeDir, cfg.Application.Homepage)
		}
		return cfg.Application.HomepageWin
	}
	return cfg.Application.Homepage
}

// GetDatabase returns path to SQLite database
func GetDatabase() string {
	return filepath.Join(homeDir, cfg.Database.DatasourceName)
}

func getPath(root, linux, darwin, windows string) string {
	if !entity.IsDebug() {
		switch runtime.GOOS {
		case "linux":
			return filepath.Join(homeDir, root, linux)
		case "darwin":
			ex, _ := os.Executable()
			return filepath.Join(filepath.Dir(ex), "../Resources", darwin)
		default:
			return filepath.Join(homeDir, root, windows)
		}
	}

	switch runtime.GOOS {
	case "linux":
		return filepath.Join(root, linux)
	case "darwin":
		return filepath.Join(root, darwin)
	default:
		return filepath.Join(root, windows)
	}
}
