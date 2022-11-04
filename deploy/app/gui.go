package main

import (
	"os"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"

	"warthog/business/entity"
	"warthog/pkg/resources"
)

const (
	defaultWindowWidth     = 1024
	defaultWindowHeight    = 768
	defaultWindowX         = 50
	defaultWindowY         = 50
	defaultMinWindowWidth  = 800
	defaultMinWindowHeight = 600
)

func getAstilectronOptions() astilectron.Options {
	iconPath := resources.GetApplicationIcon()

	zlog.Debug().
		Str("path", iconPath).
		Msgf("application icon")

	return astilectron.Options{
		AppName:            applicationName,
		AppIconDarwinPath:  iconPath,
		AppIconDefaultPath: iconPath,
		BaseDirectoryPath:  homeDir,
		SingleInstance:     true,
		VersionAstilectron: VersionAstilectron,
		VersionElectron:    VersionElectron,
		ElectronSwitches:   os.Args[1:],
	}
}

func getWindowOptions() *astilectron.WindowOptions {
	return &astilectron.WindowOptions{
		Center:    astikit.BoolPtr(true),
		Frame:     astikit.BoolPtr(true),
		Show:      astikit.BoolPtr(true),
		Width:     astikit.IntPtr(guiCfg.WindowWidth),
		Height:    astikit.IntPtr(guiCfg.WindowHeight),
		MinWidth:  astikit.IntPtr(defaultMinWindowWidth),
		MinHeight: astikit.IntPtr(defaultMinWindowHeight),
		X:         guiCfg.WindowX,
		Y:         guiCfg.WindowY,
		Title:     astikit.StrPtr(applicationName),
		Custom: &astilectron.WindowCustomOptions{
			HideOnClose: astikit.BoolPtr(!entity.IsDebug()),
		},
		WebPreferences: &astilectron.WebPreferences{
			NodeIntegrationInWorker: astikit.BoolPtr(true),
			EnableRemoteModule:      astikit.BoolPtr(true),
			//DevTools:                astikit.BoolPtr(entity.IsDebug()),
		},
	}
}

func getMenuOptions() []*astilectron.MenuItemOptions {
	return []*astilectron.MenuItemOptions{
		{
			Label: astikit.StrPtr("File"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astikit.StrPtr("DevTools"),
					Role:  astilectron.MenuItemRoleToggleDevTools,
					//Visible: astikit.BoolPtr(entity.IsDebug()),
				},
				{
					Label: astikit.StrPtr("Refresh"),
					Role:  astilectron.MenuItemRoleForceReload,
					//Visible: astikit.BoolPtr(entity.IsDebug()),
				},
				{
					Label: astikit.StrPtr("Exit"),
					Role:  astilectron.MenuItemRoleQuit,
				},
			},
		},
		{
			Label: astikit.StrPtr("Edit"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astikit.StrPtr("Undo"),
					Role:  astilectron.MenuItemRoleUndo,
				},
				{
					Label: astikit.StrPtr("Redo"),
					Role:  astilectron.MenuItemRoleRedo,
				},

				{Type: astilectron.MenuItemTypeSeparator},

				{
					Label: astikit.StrPtr("Cut"),
					Role:  astilectron.MenuItemRoleCut,
				},
				{
					Label: astikit.StrPtr("Copy"),
					Role:  astilectron.MenuItemRoleCopy,
				},
				{
					Label: astikit.StrPtr("Paste"),
					Role:  astilectron.MenuItemRolePaste,
				},
				{
					Label: astikit.StrPtr("Select all"),
					Role:  astilectron.MenuItemRoleSelectAll,
				},
			},
		},
		{
			Label: astikit.StrPtr("View"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astikit.StrPtr("Reset zoom"),
					Role:  astilectron.MenuItemRoleResetZoom,
				},
				{
					Label: astikit.StrPtr("Zoom in"),
					Role:  astilectron.MenuItemRoleZoomIn,
				},
				{
					Label: astikit.StrPtr("Zoom out"),
					Role:  astilectron.MenuItemRoleZoomOut,
				},
			},
		},
		{
			Label: astikit.StrPtr("Help"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astikit.StrPtr("About"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						return
					},
				},
			},
		},
	}
}

func getTrayOptions() *astilectron.TrayOptions {
	iconPath := resources.GetTrayIcon()

	zlog.Debug().Str("path", iconPath).Msgf("tray icon")

	return &astilectron.TrayOptions{
		Image:   astikit.StrPtr(iconPath),
		Tooltip: astikit.StrPtr(applicationName),
	}
}

func getTrayMenuOptions() []*astilectron.MenuItemOptions {
	return []*astilectron.MenuItemOptions{
		{
			Label: astikit.StrPtr("Show"),
			OnClick: func(e astilectron.Event) (deleteListener bool) {
				_ = window.Show()
				return
			},
		},
		{
			Label: astikit.StrPtr("Exit"),
			Role:  astilectron.MenuItemRoleQuit,
		},
	}
}

func initGUIEvents() {
	window.On(astilectron.EventNameWindowEventMove, func(e astilectron.Event) bool {
		x := *e.Bounds.X
		y := *e.Bounds.Y

		zlog.Debug().
			Int("x", x).
			Int("y", y).
			Msg("window.event.move")

		guiConfigUseCase.Set(&entity.GUIConfig{
			WindowX: &x,
			WindowY: &y,
		})

		return false
	})

	window.On(astilectron.EventNameWindowEventResize, func(e astilectron.Event) bool {
		zlog.Debug().
			Int("width", *e.Bounds.Width).
			Int("height", *e.Bounds.Height).
			Msg("window.event.resize")

		guiConfigUseCase.Set(&entity.GUIConfig{
			WindowWidth:  *e.Bounds.Width,
			WindowHeight: *e.Bounds.Height,
		})

		return false
	})
}

func initGUIConfig() {
	guiCfg = &entity.GUIConfig{
		WindowWidth:  defaultWindowWidth,
		WindowHeight: defaultWindowHeight,
		WindowX:      astikit.IntPtr(defaultWindowX),
		WindowY:      astikit.IntPtr(defaultWindowY),
	}

	cfg, err := guiConfigUseCase.Get()
	if err == nil {
		guiCfg = cfg
	}
}

func loadWorkspace() {
	if workspaceID == nil || *workspaceID == 0 {
		return
	}

	req := &entity.GUIRequest{
		Cmd:     entity.CmdLoadServer,
		Payload: map[string]interface{}{"id": *workspaceID},
	}

	err := window.SendMessage(req, func(_ *astilectron.EventMessage) {})
	if err != nil {
		zlog.Error().Msgf("failed to send message: %v", err)
	}
}
