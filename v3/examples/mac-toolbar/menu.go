package main

import (
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// installDaymarkMenu gives the native toolbar actions keyboard and menu
// equivalents, as a production macOS application should.
func installDaymarkMenu(app *application.App, toolbar *daymarkToolbar, split *daymarkSplit) {
	menu := app.NewMenu()
	if runtime.GOOS == "darwin" {
		menu.AddRole(application.AppMenu)
	}

	fileMenu := menu.AddSubmenu("File")
	fileMenu.Add("New Note").SetAccelerator("CmdOrCtrl+n").OnClick(func(*application.Context) {
		split.NewNote()
	})
	fileMenu.Add("Save Note").SetAccelerator("CmdOrCtrl+s").OnClick(func(*application.Context) {
		toolbar.saveNote()
	})
	fileMenu.AddSeparator()
	fileMenu.AddRole(application.CloseWindow)

	menu.AddRole(application.EditMenu)

	viewMenu := menu.AddSubmenu("View")
	viewMenu.Add("Toggle Sidebar").SetAccelerator("Ctrl+Cmd+s").OnClick(func(*application.Context) {
		split.ToggleSidebar()
	})
	viewMenu.Add("Toggle Focus").SetAccelerator("CmdOrCtrl+Shift+f").OnClick(func(*application.Context) {
		toolbar.toggleFocus()
	})
	viewMenu.Add("Toggle Inspector").SetAccelerator("CmdOrCtrl+Option+i").OnClick(func(*application.Context) {
		split.ToggleInspector()
	})
	viewMenu.AddSeparator()
	layoutMenu := viewMenu.AddSubmenu("Content Layout")
	edgeToEdge := layoutMenu.AddRadio("Edge to Edge", true)
	belowToolbar := layoutMenu.AddRadio("Below Toolbar", false)
	edgeToEdge.OnClick(func(*application.Context) {
		edgeToEdge.SetChecked(true)
		belowToolbar.SetChecked(false)
		split.SetContentLayout(application.MacContentLayoutEdgeToEdge)
	})
	belowToolbar.OnClick(func(*application.Context) {
		edgeToEdge.SetChecked(false)
		belowToolbar.SetChecked(true)
		split.SetContentLayout(application.MacContentLayoutBelowToolbar)
	})

	menu.AddRole(application.WindowMenu)
	menu.AddRole(application.HelpMenu)
	app.Menu.Set(menu)
}
