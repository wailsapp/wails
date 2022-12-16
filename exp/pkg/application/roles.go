package application

import "runtime"

// Heavily inspired by Electron (c) 2013-2020 Github Inc.
// Electron License: https://github.com/electron/electron/blob/master/LICENSE

// Role is a type to identify menu roles
type Role uint

// These constants need to be kept in sync with `v2/internal/frontend/desktop/darwin/Role.h`
const (
	NoRole             Role = iota
	AppMenu            Role = iota
	EditMenu           Role = iota
	ViewMenu           Role = iota
	ServicesMenu       Role = iota
	Hide               Role = iota
	HideOthers         Role = iota
	UnHide             Role = iota
	About              Role = iota
	Undo               Role = iota
	Redo               Role = iota
	Cut                Role = iota
	Copy               Role = iota
	Paste              Role = iota
	PasteAndMatchStyle Role = iota
	SelectAll          Role = iota
	Delete             Role = iota
	SpeechMenu         Role = iota
	Quit               Role = iota
	FileMenu           Role = iota
	Close              Role = iota
	Reload             Role = iota
	ForceReload        Role = iota
	ToggleDevTools     Role = iota
	ResetZoom          Role = iota
	ZoomIn             Role = iota
	ZoomOut            Role = iota
	ToggleFullscreen   Role = iota

	//MinimizeRole           Role =
	//QuitRole               Role =
	//TogglefullscreenRole   Role = "togglefullscreen"
	//ViewMenuRole           Role = "viewMenu"
	//WindowMenuRole         Role = "windowMenu"

	//FrontRole              Role = "front"
	//ZoomRole               Role = "zoom"
	//WindowSubMenuRole      Role = "windowSubMenu"
	//HelpSubMenuRole        Role = "helpSubMenu"
	//SeparatorItemRole      Role = "separatorItem"
)

func newFileMenu() *MenuItem {
	fileMenu := NewMenu()
	if runtime.GOOS == "darwin" {
		fileMenu.AddRole(Close)
	} else {
		fileMenu.AddRole(Quit)
	}
	subMenu := newSubMenuItem("File")
	subMenu.submenu = fileMenu
	return subMenu
}

/*
	{
	  label: 'View',
	  submenu: [
	    { role: 'reload' },
	    { role: 'forceReload' },
	    { role: 'toggleDevTools' },
	    { type: 'separator' },
	    { role: 'resetZoom' },
	    { role: 'zoomIn' },
	    { role: 'zoomOut' },
	    { type: 'separator' },
	    { role: 'togglefullscreen' }
	  ]
	},
*/
func newViewMenu() *MenuItem {
	viewMenu := NewMenu()
	viewMenu.AddRole(Reload)
	viewMenu.AddRole(ForceReload)
	viewMenu.AddRole(ToggleDevTools)
	viewMenu.AddSeparator()
	//viewMenu.AddRole(ResetZoom)
	//viewMenu.AddRole(ZoomIn)
	//viewMenu.AddRole(ZoomOut)
	viewMenu.AddSeparator()
	viewMenu.AddRole(ToggleFullscreen)
	subMenu := newSubMenuItem("View")
	subMenu.submenu = viewMenu
	return subMenu
}
