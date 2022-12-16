package application

// Heavily inspired by Electron (c) 2013-2020 Github Inc.
// Electron License: https://github.com/electron/electron/blob/master/LICENSE

// Role is a type to identify menu roles
type Role uint

// These constants need to be kept in sync with `v2/internal/frontend/desktop/darwin/Role.h`
const (
	NoRole       Role = 0
	AppMenu      Role = 1
	EditMenu     Role = 2
	ServicesMenu Role = 3
	Hide         Role = 4
	HideOthers   Role = 5
	UnHide       Role = 6
	//AboutRole              Role = "about"
	//UndoRole               Role = "undo"
	//RedoRole               Role = "redo"
	//CutRole                Role = "cut"
	//CopyRole               Role = "copy"
	//PasteRole              Role = "paste"
	//PasteAndMatchStyleRole Role = "pasteAndMatchStyle"
	//SelectAllRole          Role = "selectAll"
	//DeleteRole             Role = "delete"
	//MinimizeRole           Role = "minimize"
	//QuitRole               Role = "quit"
	//TogglefullscreenRole   Role = "togglefullscreen"
	//FileMenuRole           Role = "fileMenu"
	//ViewMenuRole           Role = "viewMenu"
	//WindowMenuRole         Role = "windowMenu"

	//FrontRole              Role = "front"
	//ZoomRole               Role = "zoom"
	//WindowSubMenuRole      Role = "windowSubMenu"
	//HelpSubMenuRole        Role = "helpSubMenu"
	//SeparatorItemRole      Role = "separatorItem"
)
