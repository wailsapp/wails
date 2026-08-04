package main

import (
	"embed"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "mac-toolbar",
		Description: "A native macOS toolbar editor",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Daymark",
		Width:  1180,
		Height: 760,
		URL:    "/",
		Mac: application.MacWindow{
			TitleBar: application.MacTitleBar{
				ToolbarStyle: application.MacToolbarStyleExpanded,
			},
		},
	})

	toolbar := application.NewMacToolbar()
	newNote := toolbar.AddButton("New", "square.and.pencil").SetBordered(true)
	mode := toolbar.AddGroup("Mode", application.ToolbarGroupSelectOne)
	mode.AddButton("Write", "pencil").OnClick(func(*application.Context) {
		mode.SetSelectedIndex(0)
		app.Event.Emit("toolbar:mode", "write")
	})
	mode.AddButton("Preview", "doc.richtext").OnClick(func(*application.Context) {
		mode.SetSelectedIndex(1)
		app.Event.Emit("toolbar:mode", "preview")
	})
	toolbar.AddFlexibleSpace()
	search := toolbar.AddSearch("Search notes").SetTooltip("Search your notes")
	share := toolbar.AddButton("Share", "square.and.arrow.up").SetBordered(true)
	details := toolbar.AddButton("Details", "info.circle").SetBordered(true)
	save := toolbar.AddButton("Save", "checkmark.circle").SetProminent(true).SetBadgeCount(0)
	focus := toolbar.AddButton("Focus", "arrow.up.left.and.arrow.down.right").SetBordered(true)

	focused := false
	dirty := false
	newNote.OnClick(func(*application.Context) { app.Event.Emit("toolbar:new") })
	search.OnSearch(func(_ *application.Context, query string) { app.Event.Emit("toolbar:search", query) })
	share.OnClick(func(*application.Context) { app.Event.Emit("toolbar:share") })
	details.OnClick(func(*application.Context) { app.Event.Emit("toolbar:details") })
	save.OnClick(func(*application.Context) {
		dirty = false
		save.SetBadgeCount(0).SetLabel("Saved")
		app.Event.Emit("toolbar:save")
	})
	focus.OnClick(func(*application.Context) {
		focused = !focused
		details.SetHidden(focused)
		app.Event.Emit("toolbar:focus", focused)
	})

	menu := app.NewMenu()
	if runtime.GOOS == "darwin" {
		menu.AddRole(application.AppMenu)
	}
	menu.AddRole(application.EditMenu)
	menu.AddRole(application.WindowMenu)
	menu.AddRole(application.HelpMenu)
	daymarkMenu := menu.AddSubmenu("Daymark")
	daymarkMenu.Add("New Note").SetAccelerator("CmdOrCtrl+n").OnClick(func(*application.Context) {
		app.Event.Emit("toolbar:new")
	})
	daymarkMenu.Add("Save Note").SetAccelerator("CmdOrCtrl+s").OnClick(func(*application.Context) {
		dirty = false
		save.SetBadgeCount(0).SetLabel("Saved")
		app.Event.Emit("toolbar:save")
	})
	daymarkMenu.Add("Toggle Focus").SetAccelerator("CmdOrCtrl+Shift+f").OnClick(func(*application.Context) {
		focused = !focused
		details.SetHidden(focused)
		app.Event.Emit("toolbar:focus", focused)
	})
	daymarkMenu.Add("Toggle Inspector").SetAccelerator("CmdOrCtrl+Option+i").OnClick(func(*application.Context) {
		app.Event.Emit("toolbar:details")
	})
	app.Menu.Set(menu)

	app.Event.On("editor:dirty", func(e *application.CustomEvent) {
		isDirty, ok := e.Data.(bool)
		if !ok || isDirty == dirty {
			return
		}
		dirty = isDirty
		if dirty {
			save.SetBadgeCount(1).SetLabel("Save")
		} else {
			save.SetBadgeCount(0).SetLabel("Saved")
		}
	})

	window.SetToolbar(toolbar)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
