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
	ShowAll            Role = iota
	BringAllToFront    Role = iota
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
	CloseWindow        Role = iota
	Reload             Role = iota
	ForceReload        Role = iota
	OpenDevTools       Role = iota
	ResetZoom          Role = iota
	ZoomIn             Role = iota
	ZoomOut            Role = iota
	ToggleFullscreen   Role = iota

	Minimize   Role = iota
	Zoom       Role = iota
	FullScreen Role = iota

	NewFile        Role = iota
	Open           Role = iota
	Save           Role = iota
	SaveAs         Role = iota
	StartSpeaking  Role = iota
	StopSpeaking   Role = iota
	Revert         Role = iota
	Print          Role = iota
	PageLayout     Role = iota
	Find           Role = iota
	FindAndReplace Role = iota
	FindNext       Role = iota
	FindPrevious   Role = iota
	Front          Role = iota
	Help           Role = iota
)

func NewFileMenu() *MenuItem {
	fileMenu := NewMenu()
	if runtime.GOOS == "darwin" {
		fileMenu.AddRole(CloseWindow)
	} else {
		fileMenu.AddRole(Quit)
	}
	subMenu := NewSubMenuItem("File")
	subMenu.submenu = fileMenu
	return subMenu
}

func NewViewMenu() *MenuItem {
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
	subMenu := NewSubMenuItem("View")
	subMenu.submenu = viewMenu
	return subMenu
}

func NewAppMenu() *MenuItem {
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
	subMenu := NewSubMenuItem(globalApplication.options.Name)
	subMenu.submenu = appMenu
	return subMenu
}

func NewEditMenu() *MenuItem {
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
	subMenu := NewSubMenuItem("Edit")
	subMenu.submenu = editMenu
	return subMenu
}

func NewWindowMenu() *MenuItem {
	menu := NewMenu()
	menu.AddRole(Minimize)
	menu.AddRole(Zoom)
	if runtime.GOOS == "darwin" {
		menu.AddSeparator()
		menu.AddRole(Front)
		//menu.AddSeparator()
		//menu.AddRole(Window)
	} else {
		menu.AddRole(CloseWindow)
	}
	subMenu := NewSubMenuItem("Window")
	subMenu.submenu = menu
	return subMenu
}

func NewHelpMenu() *MenuItem {
	menu := NewMenu()
	menu.Add("Learn More").OnClick(func(ctx *Context) {
		globalApplication.CurrentWindow().SetURL("https://wails.io")
	})
	subMenu := NewSubMenuItem("Help")
	subMenu.submenu = menu
	return subMenu
}
