package entity

const (
	CmdDevTools            GUICommand = "dev.tools.show"
	CmdGetWorkspace        GUICommand = "workspace.get"
	CmdSortingWorkspace    GUICommand = "workspace.sorting"
	CmdDeleteWorkspace     GUICommand = "workspace.delete"
	CmdExpandWorkspace     GUICommand = "workspace.expand"
	CmdStateWorkspace      GUICommand = "workspace.state"
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
)

type GUICommand string

func (c GUICommand) String() string {
	return string(c)
}
