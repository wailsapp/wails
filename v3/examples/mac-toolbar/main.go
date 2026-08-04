package main

import (
	"embed"
	"log"
	"runtime"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "Daymark",
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
	newNote := toolbar.AddButton("New").SetSymbol("square.and.pencil").SetBordered(true)
	mode := toolbar.AddGroup("Mode", application.ToolbarGroupSelectOne)
	mode.AddButton("Write").SetSymbol("pencil").OnClick(func(*application.Context) {
		mode.SetSelectedIndex(0)
		app.Event.Emit("toolbar:mode", "write")
	})
	mode.AddButton("Preview").SetSymbol("doc.richtext").OnClick(func(*application.Context) {
		mode.SetSelectedIndex(1)
		app.Event.Emit("toolbar:mode", "preview")
	})
	toolbar.AddFlexibleSpace()
	search := toolbar.AddSearch("Search notes").SetTooltip("Search your notes")
	share := toolbar.AddButton("Share").SetSymbol("square.and.arrow.up").SetBordered(true)
	details := toolbar.AddButton("Details").SetSymbol("info.circle").SetBordered(true)
	save := toolbar.AddButton("Save").SetSymbol("checkmark.circle").SetProminent(true).SetBadgeCount(0)
	focus := toolbar.AddButton("Focus").SetSymbol("arrow.up.left.and.arrow.down.right").SetBordered(true)

	var stateLock sync.Mutex
	focused := false
	dirty := false
	setSaved := func() {
		stateLock.Lock()
		dirty = false
		stateLock.Unlock()
		save.SetBadgeCount(0).SetLabel("Saved")
		app.Event.Emit("toolbar:save")
	}
	toggleFocus := func() {
		stateLock.Lock()
		focused = !focused
		isFocused := focused
		stateLock.Unlock()
		details.SetHidden(isFocused)
		if isFocused {
			focus.SetLabel("Exit Focus").SetSymbol("arrow.down.right.and.arrow.up.left")
		} else {
			focus.SetLabel("Focus").SetSymbol("arrow.up.left.and.arrow.down.right")
		}
		app.Event.Emit("toolbar:focus", isFocused)
	}

	newNote.OnClick(func(*application.Context) { app.Event.Emit("toolbar:new") })
	search.OnSearch(func(_ *application.Context, query string) { app.Event.Emit("toolbar:search", query) })
	share.OnClick(func(*application.Context) { app.Event.Emit("toolbar:share") })
	details.OnClick(func(*application.Context) { app.Event.Emit("toolbar:details") })
	save.OnClick(func(*application.Context) { setSaved() })
	focus.OnClick(func(*application.Context) { toggleFocus() })

	menu := app.NewMenu()
	if runtime.GOOS == "darwin" {
		menu.AddRole(application.AppMenu)
	}
	fileMenu := menu.AddSubmenu("File")
	fileMenu.Add("New Note").SetAccelerator("CmdOrCtrl+n").OnClick(func(*application.Context) {
		app.Event.Emit("toolbar:new")
	})
	fileMenu.Add("Save Note").SetAccelerator("CmdOrCtrl+s").OnClick(func(*application.Context) {
		setSaved()
	})
	fileMenu.AddSeparator()
	fileMenu.AddRole(application.CloseWindow)
	menu.AddRole(application.EditMenu)
	viewMenu := menu.AddSubmenu("View")
	viewMenu.Add("Toggle Focus").SetAccelerator("CmdOrCtrl+Shift+f").OnClick(func(*application.Context) {
		toggleFocus()
	})
	viewMenu.Add("Toggle Inspector").SetAccelerator("CmdOrCtrl+Option+i").OnClick(func(*application.Context) {
		app.Event.Emit("toolbar:details")
	})
	menu.AddRole(application.WindowMenu)
	menu.AddRole(application.HelpMenu)
	app.Menu.Set(menu)

	app.Event.On("editor:dirty", func(e *application.CustomEvent) {
		isDirty, ok := e.Data.(bool)
		if !ok {
			return
		}
		stateLock.Lock()
		if isDirty == dirty {
			stateLock.Unlock()
			return
		}
		dirty = isDirty
		stateLock.Unlock()
		if isDirty {
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
