// Package entity provides entities for business logic.
package entity

// UI events
const (
	CmdApplicationState    GUICommand = "application.state"
	CmdSettingsUpdate      GUICommand = "settings.update"
	CmdGetWorkspace        GUICommand = "workspace.get"
	CmdSortingWorkspace    GUICommand = "workspace.sorting"
	CmdDeleteWorkspace     GUICommand = "workspace.delete"
	CmdDuplicateWorkspace  GUICommand = "workspace.duplicate"
	CmdExpandWorkspace     GUICommand = "workspace.expand"
	CmdCreateServer        GUICommand = "server.create"
	CmdUpdateServer        GUICommand = "server.update"
	CmdUpdateServerRequest GUICommand = "server.update.request"
	CmdLoadServer          GUICommand = "server.load"
	CmdCreateFolder        GUICommand = "folder.create"
	CmdUpdateFolder        GUICommand = "folder.update"
	CmdDeleteFolder        GUICommand = "folder.delete"
	CmdUpdateQuery         GUICommand = "query.update"
	CmdRunQuery            GUICommand = "query.run"
	CmdCancelQuery         GUICommand = "query.cancel"
	CmdCloseStream         GUICommand = "query.close.stream"
	CmdQueryResponse       GUICommand = "query.response"
	CmdDevTools            GUICommand = "dev.tools.show"
	CmdMenuSettings        GUICommand = "menu.settings"
	CmdMenuAbout           GUICommand = "menu.about"
	CmdMessageInfo         GUICommand = "message.info"
	CmdMessageError        GUICommand = "message.error"
)

// GUICommand UI command
type GUICommand string

// String returns UI command string
func (c GUICommand) String() string {
	return string(c)
}
