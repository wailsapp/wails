package main

import (
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// installDaymarkMenu gives the native toolbar actions keyboard and menu
// equivalents, as a production macOS application should.
func installDaymarkMenu(app *application.App, toolbar *daymarkToolbar) {
	menu := app.NewMenu()
	if runtime.GOOS == "darwin" {
		menu.AddRole(application.AppMenu)
	}

	fileMenu := menu.AddSubmenu("File")
	fileMenu.Add("New Note").SetAccelerator("CmdOrCtrl+n").OnClick(func(*application.Context) {
		app.Event.Emit("toolbar:new")
	})
	fileMenu.Add("Save Note").SetAccelerator("CmdOrCtrl+s").OnClick(func(*application.Context) {
		toolbar.saveNote()
	})
	fileMenu.AddSeparator()
	fileMenu.AddRole(application.CloseWindow)

	menu.AddRole(application.EditMenu)

	viewMenu := menu.AddSubmenu("View")
	viewMenu.Add("Toggle Focus").SetAccelerator("CmdOrCtrl+Shift+f").OnClick(func(*application.Context) {
		toolbar.toggleFocus()
	})
	viewMenu.Add("Toggle Inspector").SetAccelerator("CmdOrCtrl+Option+i").OnClick(func(*application.Context) {
		app.Event.Emit("toolbar:details")
	})

	menu.AddRole(application.WindowMenu)
	menu.AddRole(application.HelpMenu)
	app.Menu.Set(menu)
}
