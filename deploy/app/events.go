package main

import (
	"fmt"

	"github.com/forest33/warthog/business/entity"
)

func eventsHandler(r *entity.GUIRequest) *entity.GUIResponse {
	resp := &entity.GUIResponse{}

	zlog.Debug().
		Str("cmd", r.Cmd.String()).
		Interface("payload", r.Payload).
		Msg("event")

	var payload map[string]interface{}
	if r.Payload != nil {
		payload = r.Payload.(map[string]interface{})
	}

	switch r.Cmd {
	case entity.CmdApplicationState:
		resp = getApplicationState()
	case entity.CmdSettingsUpdate:
		resp = settingsUseCase.Update(payload)
	case entity.CmdGetWorkspace:
		resp = workspaceUseCase.Get(payload)
	case entity.CmdSortingWorkspace:
		resp = workspaceUseCase.Sorting(payload)
	case entity.CmdDeleteWorkspace:
		resp = workspaceUseCase.Delete(payload)
	case entity.CmdExpandWorkspace:
		resp = workspaceUseCase.Expand(payload)
	case entity.CmdCreateFolder:
		resp = workspaceUseCase.CreateFolder(payload)
	case entity.CmdUpdateFolder:
		resp = workspaceUseCase.UpdateFolder(payload)
	case entity.CmdDeleteFolder:
		resp = workspaceUseCase.DeleteFolder(payload)
	case entity.CmdCreateServer:
		resp = workspaceUseCase.CreateServer(payload)
	case entity.CmdUpdateServer:
		resp = workspaceUseCase.UpdateServer(payload)
	case entity.CmdUpdateServerRequest:
		resp = workspaceUseCase.UpdateServerRequest(payload)
	case entity.CmdUpdateQuery:
		resp = workspaceUseCase.UpdateQuery(payload)
	case entity.CmdLoadServer:
		resp = grpcUseCase.LoadServer(payload)
	case entity.CmdRunQuery:
		resp = grpcUseCase.Query(payload)
	case entity.CmdCancelQuery:
		grpcUseCase.CancelQuery()
	case entity.CmdCloseStream:
		grpcUseCase.CloseStream()
	case entity.CmdDevTools:
		_ = window.OpenDevTools()
	default:
		resp = entity.ErrorGUIResponse(fmt.Errorf("unknown command %s", r.Cmd))
	}

	return resp
}
