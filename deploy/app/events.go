package main

import (
	"fmt"

	"warthog/business/entity"
)

func eventsHandler(r *entity.GUIRequest) *entity.GUIResponse {
	resp := &entity.GUIResponse{}

	zlog.Debug().
		Str("cmd", r.Cmd.String()).
		Interface("payload", r.Payload).
		Msg("event")

	switch r.Cmd {
	case entity.CmdGetWorkspace:
		resp = workspaceUseCase.Get(r.Payload)
	case entity.CmdSortingWorkspace:
		resp = workspaceUseCase.Sorting(r.Payload)
	case entity.CmdDeleteWorkspace:
		resp = workspaceUseCase.Delete(r.Payload)
	case entity.CmdExpandWorkspace:
		resp = workspaceUseCase.Expand(r.Payload)
	case entity.CmdStateWorkspace:
		resp = workspaceUseCase.GetState()
	case entity.CmdCreateFolder:
		resp = workspaceUseCase.CreateFolder(r.Payload)
	case entity.CmdUpdateFolder:
		resp = workspaceUseCase.UpdateFolder(r.Payload)
	case entity.CmdDeleteFolder:
		resp = workspaceUseCase.DeleteFolder(r.Payload)
	case entity.CmdCreateServer:
		resp = workspaceUseCase.CreateServer(r.Payload)
	case entity.CmdUpdateServer:
		resp = workspaceUseCase.UpdateServer(r.Payload)
	case entity.CmdUpdateServerRequest:
		resp = workspaceUseCase.UpdateServerRequest(r.Payload)
	case entity.CmdUpdateQuery:
		resp = workspaceUseCase.UpdateQuery(r.Payload)
	case entity.CmdLoadServer:
		resp = grpcUseCase.LoadServer(r.Payload)
	case entity.CmdRunQuery:
		resp = grpcUseCase.Query(r.Payload)
	case entity.CmdCancelQuery:
		grpcUseCase.CancelQuery()
	case entity.CmdDevTools:
		_ = window.OpenDevTools()
	default:
		resp = entity.ErrorGUIResponse(fmt.Errorf("unknown command %s", r.Cmd))
	}

	return resp
}
