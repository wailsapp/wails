// Package menu provides all the functions and structs related to menus in a Wails application.
// Heavily inspired by Electron (c) 2013-2020 Github Inc.
// Electron License: https://github.com/electron/electron/blob/master/LICENSE
package menu

// Role is a type to identify menu roles
type Role int

// These constants need to be kept in sync with `v2/internal/frontend/desktop/darwin/Role.h`
const (
	AppMenuRole    Role = 1
	EditMenuRole        = 2
	WindowMenuRole      = 3
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
	//HideRole               Role = "hide"
	//HideOthersRole         Role = "hideOthers"
	//UnhideRole             Role = "unhide"
	//FrontRole              Role = "front"
	//ZoomRole               Role = "zoom"
	//WindowSubMenuRole      Role = "windowSubMenu"
	//HelpSubMenuRole        Role = "helpSubMenu"
	//SeparatorItemRole      Role = "separatorItem"
)

/*
// About provides a MenuItem with the About role
func About() *MenuItem {
	return &MenuItem{
		Role: AboutRole,
	}
}

// Undo provides a MenuItem with the Undo role
func Undo() *MenuItem {
	return &MenuItem{
		Role: UndoRole,
	}
}

// Redo provides a MenuItem with the Redo role
func Redo() *MenuItem {
	return &MenuItem{
		Role: RedoRole,
	}
}

// Cut provides a MenuItem with the Cut role
func Cut() *MenuItem {
	return &MenuItem{
		Role: CutRole,
	}
}

// Copy provides a MenuItem with the Copy role
func Copy() *MenuItem {
	return &MenuItem{
		Role: CopyRole,
	}
}

// Paste provides a MenuItem with the Paste role
func Paste() *MenuItem {
	return &MenuItem{
		Role: PasteRole,
	}
}

// PasteAndMatchStyle provides a MenuItem with the PasteAndMatchStyle role
func PasteAndMatchStyle() *MenuItem {
	return &MenuItem{
		Role: PasteAndMatchStyleRole,
	}
}

// SelectAll provides a MenuItem with the SelectAll role
func SelectAll() *MenuItem {
	return &MenuItem{
		Role: SelectAllRole,
	}
}

// Delete provides a MenuItem with the Delete role
func Delete() *MenuItem {
	return &MenuItem{
		Role: DeleteRole,
	}
}

// Minimize provides a MenuItem with the Minimize role
func Minimize() *MenuItem {
	return &MenuItem{
		Role: MinimizeRole,
	}
}

// Quit provides a MenuItem with the Quit role
func Quit() *MenuItem {
	return &MenuItem{
		Role: QuitRole,
	}
}

// ToggleFullscreen provides a MenuItem with the ToggleFullscreen role
func ToggleFullscreen() *MenuItem {
	return &MenuItem{
		Role: TogglefullscreenRole,
	}
}

// FileMenu provides a MenuItem with the whole default "File" menu (Close / Quit)
func FileMenu() *MenuItem {
	return &MenuItem{
		Role: FileMenuRole,
	}
}
*/

// EditMenu provides a MenuItem with the whole default "Edit" menu (Undo, Copy, etc.).
func EditMenu() *MenuItem {
	return &MenuItem{
		Role: EditMenuRole,
	}
}

/*
// ViewMenu provides a MenuItem with the whole default "View" menu (Reload, Toggle Developer Tools, etc.)
func ViewMenu() *MenuItem {
	return &MenuItem{
		Role: ViewMenuRole,
	}
}
*/

// WindowMenu provides a MenuItem with the whole default "Window" menu (Minimize, Zoom, etc.).
// On MacOS currently all options in there won't work if the window is frameless.
func WindowMenu() *MenuItem {
	return &MenuItem{
		Role: WindowMenuRole,
	}
}

// These roles are Mac only

// AppMenu provides a MenuItem with the whole default "App" menu (About, Services, etc.)
func AppMenu() *MenuItem {
	return &MenuItem{
		Role: AppMenuRole,
	}
}

/*
// Hide provides a MenuItem that maps to the hide action.
func Hide() *MenuItem {
	return &MenuItem{
		Role: HideRole,
	}
}

// HideOthers provides a MenuItem that maps to the hideOtherApplications action.
func HideOthers() *MenuItem {
	return &MenuItem{
		Role: HideOthersRole,
	}
}

// UnHide provides a MenuItem that maps to the unHideAllApplications action.
func UnHide() *MenuItem {
	return &MenuItem{
		Role: UnhideRole,
	}
}

// Front provides a MenuItem that maps to the arrangeInFront action.
func Front() *MenuItem {
	return &MenuItem{
		Role: FrontRole,
	}
}

// Zoom provides a MenuItem that maps to the performZoom action.
func Zoom() *MenuItem {
	return &MenuItem{
		Role: ZoomRole,
	}
}

// WindowSubMenu provides a MenuItem with the "Window" submenu.
func WindowSubMenu() *MenuItem {
	return &MenuItem{
		Role: WindowSubMenuRole,
	}
}

// HelpSubMenu provides a MenuItem with the "Help" submenu.
func HelpSubMenu() *MenuItem {
	return &MenuItem{
		Role: HelpSubMenuRole,
	}
}
*/
