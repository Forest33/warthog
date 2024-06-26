package main

import (
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"

	"github.com/forest33/warthog/business/entity"
	"github.com/forest33/warthog/pkg/resources"
)

const (
	defaultMinWindowWidth     = 800
	defaultMinWindowHeight    = 600
	defaultVersionAstilectron = "0.56.0"
	defaultVersionElectron    = "13.6.9"
)

func getAstilectronOptions() astilectron.Options {
	if VersionAstilectron == "" {
		VersionAstilectron = defaultVersionAstilectron
	}
	if VersionElectron == "" {
		VersionElectron = defaultVersionElectron
	}

	iconPath := resources.GetApplicationIcon()

	zlog.Debug().
		Str("path", iconPath).
		Msgf("application icon")

	return astilectron.Options{
		AppName:            applicationName,
		AppIconDarwinPath:  iconPath,
		AppIconDefaultPath: iconPath,
		BaseDirectoryPath:  homeDir,
		SingleInstance:     *settings.SingleInstance,
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
		Width:     astikit.IntPtr(settings.WindowWidth),
		Height:    astikit.IntPtr(settings.WindowHeight),
		MinWidth:  astikit.IntPtr(defaultMinWindowWidth),
		MinHeight: astikit.IntPtr(defaultMinWindowHeight),
		X:         settings.WindowX,
		Y:         settings.WindowY,
		Title:     astikit.StrPtr(applicationName),
		Custom: &astilectron.WindowCustomOptions{
			HideOnClose: astikit.BoolPtr(!entity.IsDebug()),
		},
		WebPreferences: &astilectron.WebPreferences{
			NodeIntegrationInWorker: astikit.BoolPtr(true),
			EnableRemoteModule:      astikit.BoolPtr(true),
		},
	}
}

func getMenuOptions() []*astilectron.MenuItemOptions {
	return []*astilectron.MenuItemOptions{
		{
			Label: astikit.StrPtr("File"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label:   astikit.StrPtr("DevTools"),
					Role:    astilectron.MenuItemRoleToggleDevTools,
					Visible: astikit.BoolPtr(entity.IsDebug()),
				},
				{
					Label:   astikit.StrPtr("Refresh"),
					Role:    astilectron.MenuItemRoleForceReload,
					Visible: astikit.BoolPtr(entity.IsDebug()),
				},
				{
					Label:       astikit.StrPtr("Settings"),
					Accelerator: astilectron.NewAccelerator("CommandOrControl+,"),
					OnClick:     menuSettings,
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
					Label: astikit.StrPtr("Check for Updates..."),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						menuCheckUpdates()
						return
					},
				},
				{
					Label: astikit.StrPtr("About"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						menuAbout()
						return
					},
				},
			},
		},
	}
}

func getTrayOptions() *astilectron.TrayOptions {
	if !*settings.SingleInstance {
		return nil
	}

	iconPath := resources.GetTrayIcon()

	zlog.Debug().Str("path", iconPath).Msgf("tray icon")

	return &astilectron.TrayOptions{
		Image:   astikit.StrPtr(iconPath),
		Tooltip: astikit.StrPtr(applicationName),
	}
}

func getTrayMenuOptions() []*astilectron.MenuItemOptions {
	if !*settings.SingleInstance {
		return nil
	}

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

		settingsUseCase.Set(&entity.Settings{
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

		settingsUseCase.Set(&entity.Settings{
			WindowWidth:  *e.Bounds.Width,
			WindowHeight: *e.Bounds.Height,
		})

		return false
	})
}

func initGrpcResponse() {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case resp := <-grpcClient.GetResponseChannel():
				req := &entity.GUIRequest{
					Cmd:     entity.CmdQueryResponse,
					Payload: resp,
				}
				err := window.SendMessage(req, func(_ *astilectron.EventMessage) {})
				if err != nil {
					zlog.Error().Msgf("failed to send gRPC response: %v", err)
				}
			}
		}
	}()
}

func initAsyncMessages() {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case resp := <-grpcUseCase.GetInfoChannel():
				req := &entity.GUIRequest{
					Cmd:     entity.CmdMessageInfo,
					Payload: resp,
				}
				err := window.SendMessage(req, func(_ *astilectron.EventMessage) {})
				if err != nil {
					zlog.Error().Msgf("failed to send info message: %v", err)
				}
			case resp := <-grpcUseCase.GetErrorChannel():
				req := &entity.GUIRequest{
					Cmd:     entity.CmdMessageError,
					Payload: resp,
				}
				err := window.SendMessage(req, func(_ *astilectron.EventMessage) {})
				if err != nil {
					zlog.Error().Msgf("failed to send error message: %v", err)
				}
			}
		}
	}()
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

func menuAbout() {
	if b := strings.Split(BuiltAt, " m="); len(b) > 0 {
		if bt, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", b[0]); err == nil {
			BuiltAt = bt.Format("2006-01-02 15:04:05")
		}
	}

	req := &entity.GUIRequest{
		Cmd: entity.CmdMenuAbout,
		Payload: map[string]interface{}{
			"app_name":            AppName,
			"app_version":         AppVersion,
			"app_url":             AppURL,
			"astilectron_version": VersionAstilectron,
			"electron_version":    VersionElectron,
			"built_at":            BuiltAt,
			"go_version":          strings.ReplaceAll(runtime.Version(), "go", ""),
		},
	}

	err := window.SendMessage(req, func(_ *astilectron.EventMessage) {})
	if err != nil {
		zlog.Error().Msgf("failed to send message: %v", err)
	}
}

func menuCheckUpdates() {
	settingsUseCase.CheckUpdates(getCheckUpdatesCallback(false))
}

func menuSettings(e astilectron.Event) (deleteListener bool) {
	err := window.SendMessage(&entity.GUIRequest{Cmd: entity.CmdMenuSettings}, func(_ *astilectron.EventMessage) {})
	if err != nil {
		zlog.Error().Msgf("failed to send message: %v", err)
	}
	return false
}

type applicationState struct {
	State    *entity.WorkspaceState `json:"state"`
	Settings *entity.Settings       `json:"settings"`
}

func getApplicationState() *entity.GUIResponse {
	state, err := workspaceUseCase.GetState()
	if err != nil {
		return entity.ErrorGUIResponse(err)
	}

	if settings.IsCheckUpdates() {
		settingsUseCase.CheckUpdates(getCheckUpdatesCallback(true))
	}

	return &entity.GUIResponse{
		Status: entity.GUIResponseStatusOK,
		Payload: &applicationState{
			State:    state,
			Settings: settings,
		},
	}
}

func getCheckUpdatesCallback(background bool) func(r *entity.GithubRelease) {
	return func(r *entity.GithubRelease) {
		req := &entity.GUIRequest{
			Cmd: entity.CmdCheckUpdates,
		}

		var version string
		if r == nil {
			version = AppVersion
		} else {
			version = strings.ReplaceAll(r.Name, "v", "")
		}

		if version != AppVersion {
			req.Payload = map[string]interface{}{
				"version": version,
				"url":     r.HtmlUrl,
			}
		} else if background {
			return
		}

		err := window.SendMessage(req, func(_ *astilectron.EventMessage) {})
		if err != nil {
			zlog.Error().Msgf("failed to send message: %v", err)
		}
	}
}
