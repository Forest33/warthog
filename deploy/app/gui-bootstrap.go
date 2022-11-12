package main

import (
	"encoding/json"

	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"

	"github.com/forest33/warthog/business/entity"
	"github.com/forest33/warthog/pkg/resources"
)

func getWindow() []*bootstrap.Window {
	homepage := resources.GetHomepage()
	zlog.Debug().Str("path", homepage).Msg("homepage")

	return []*bootstrap.Window{{
		Homepage:       homepage,
		MessageHandler: bootstrapMessageHandler,
		Options:        getWindowOptions(),
	}}
}

type bootstrapResponse struct {
	Data   interface{}              `json:"data"`
	Status entity.GUIResponseStatus `json:"status"`
	Error  entity.Error             `json:"error"`
}

func bootstrapMessageHandler(_ *astilectron.Window, m bootstrap.MessageIn) (interface{}, error) {
	req := &entity.GUIRequest{
		Cmd: entity.GUICommand(m.Name),
	}

	if len(m.Payload) > 0 {
		var payload map[string]interface{}
		if err := json.Unmarshal(m.Payload, &payload); err != nil {
			zlog.Error().Msgf("failed to unmarshal GUI event: %v", err)
			return err, err
		}
		req.Payload = payload
	}

	resp := eventsHandler(req)
	if resp != nil {
		return &bootstrapResponse{
			Data:   resp.Payload,
			Status: resp.Status,
			Error:  resp.Error,
		}, nil
	}

	return nil, nil
}

func withBootstrap() {
	zlog.Info().
		Bool("debug", entity.IsDebug()).
		Msg("Application using bootstrap")

	options := bootstrap.Options{
		Asset:              Asset,
		AssetDir:           AssetDir,
		RestoreAssets:      RestoreAssets,
		AstilectronOptions: getAstilectronOptions(),
		Windows:            getWindow(),
		Debug:              false,
		Logger:             zlog,
		TrayOptions:        getTrayOptions(),
		TrayMenuOptions:    getTrayMenuOptions(),
		MenuOptions:        getMenuOptions(),
		OnWait: func(a *astilectron.Astilectron, w []*astilectron.Window, _ *astilectron.Menu, _t *astilectron.Tray, _ *astilectron.Menu) error {
			ast = a
			window = w[0]
			tray = _t

			initGUIEvents()
			loadWorkspace()

			return nil
		},
	}

	err := bootstrap.Run(options)
	if err != nil {
		zlog.Fatal(err.Error())
	}
}
