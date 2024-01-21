package application

import "runtime"

// Heavily inspired by Electron (c) 2013-2020 Github Inc.
// Electron License: https://github.com/electron/electron/blob/master/LICENSE

// Role is a type to identify menu roles
type Role uint

// These constants need to be kept in sync with `v2/internal/frontend/desktop/darwin/Role.h`
const (
	NoRole       Role = iota
	AppMenu      Role = iota
	EditMenu     Role = iota
	ViewMenu     Role = iota
	WindowMenu   Role = iota
	ServicesMenu Role = iota
	HelpMenu     Role = iota

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
	ShowDevTools       Role = iota
	ResetZoom          Role = iota
	ZoomIn             Role = iota
	ZoomOut            Role = iota
	ToggleFullscreen   Role = iota

	Minimize   Role = iota
	Zoom       Role = iota
	FullScreen Role = iota
	//Front      Role = iota
	//WindowRole Role = iota

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

func newViewMenu() *MenuItem {
	viewMenu := NewMenu()
	viewMenu.AddRole(Reload)
	viewMenu.AddRole(ForceReload)
	addDevToolMenuItem(viewMenu)
	viewMenu.AddSeparator()
	viewMenu.AddRole(ResetZoom)
	viewMenu.AddRole(ZoomIn)
	viewMenu.AddRole(ZoomOut)
	viewMenu.AddSeparator()
	viewMenu.AddRole(ToggleFullscreen)
	subMenu := newSubMenuItem("View")
	subMenu.submenu = viewMenu
	return subMenu
}

func newAppMenu() *MenuItem {
	if runtime.GOOS != "darwin" {
		return nil
	}
	appMenu := NewMenu()
	appMenu.AddRole(About)
	appMenu.AddSeparator()
	appMenu.AddRole(ServicesMenu)
	appMenu.AddSeparator()
	appMenu.AddRole(Hide)
	appMenu.AddRole(HideOthers)
	appMenu.AddRole(UnHide)
	appMenu.AddSeparator()
	appMenu.AddRole(Quit)
	subMenu := newSubMenuItem(globalApplication.options.Name)
	subMenu.submenu = appMenu
	return subMenu
}

func newEditMenu() *MenuItem {
	editMenu := NewMenu()
	editMenu.AddRole(Undo)
	editMenu.AddRole(Redo)
	editMenu.AddSeparator()
	editMenu.AddRole(Cut)
	editMenu.AddRole(Copy)
	editMenu.AddRole(Paste)
	if runtime.GOOS == "darwin" {
		editMenu.AddRole(PasteAndMatchStyle)
		editMenu.AddRole(PasteAndMatchStyle)
		editMenu.AddRole(Delete)
		editMenu.AddRole(SelectAll)
		editMenu.AddSeparator()
		editMenu.AddRole(SpeechMenu)
	} else {
		editMenu.AddRole(Delete)
		editMenu.AddSeparator()
		editMenu.AddRole(SelectAll)
	}
	subMenu := newSubMenuItem("Edit")
	subMenu.submenu = editMenu
	return subMenu
}

func newWindowMenu() *MenuItem {
	menu := NewMenu()
	menu.AddRole(Minimize)
	menu.AddRole(Zoom)
	if runtime.GOOS == "darwin" {
		menu.AddSeparator()
		menu.AddRole(FullScreen)
	} else {
		menu.AddRole(Close)
	}
	subMenu := newSubMenuItem("Window")
	subMenu.submenu = menu
	return subMenu
}

func newHelpMenu() *MenuItem {
	menu := NewMenu()
	menu.Add("Learn More").OnClick(func(ctx *Context) {
		globalApplication.CurrentWindow().SetURL("https://wails.io")
	})
	subMenu := newSubMenuItem("Help")
	subMenu.submenu = menu
	return subMenu
}
