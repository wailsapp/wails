package main

import (
	"embed"
	"fmt"
	"log"
	"strconv"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed assets
var assets embed.FS

// The counter is the state that survives a relaunch: every increment is
// stored with the window's restorable state, and OnRestore reads it back.
var counter atomic.Int64

// kiosk is the presentation the Kiosk toggle applies: Dock and menu bar
// hidden until the pointer reaches them, no Command-Tab, no Hide. Force
// Quit stays available on purpose.
const kiosk = application.MacPresentationAutoHideDock |
	application.MacPresentationAutoHideMenuBar |
	application.MacPresentationDisableProcessSwitching |
	application.MacPresentationDisableHideApplication

func main() {
	app := application.New(application.Options{
		Name:        "Windows Extra",
		Description: "Sheets, popovers, presentation options and state restoration",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
			// Opt in to secure coding for the restorable state. The
			// delegate method exists either way, so macOS 14 never logs
			// its "Secure coding is automatically enabled" warning.
			SupportsSecureRestorableState: true,
		},
	})

	// Never leave the Dock or menu bar hidden behind us.
	app.OnShutdown(func() {
		_ = app.SetPresentationOptions(application.MacPresentationDefault)
	})

	// The popover shown from the toolbar and from the tray. It hosts a
	// native control strip; a popover is created lazily on first show and
	// can be shown from several anchors in turn.
	popover := newCounterPopover(app)
	trayPopover := newTrayPopover(app)

	tray := app.SystemTray.New()
	tray.SetSymbol("rectangle.stack").SetTooltip("Windows Extra")
	tray.OnClick(func() {
		if trayPopover.IsShown() {
			trayPopover.Close()
			return
		}
		if err := tray.ShowPopover(trayPopover); err != nil {
			app.Logger.Error("tray popover", "error", err)
		}
	})

	installMenu(app, popover)

	// State restoration: macOS asks for every window that was visible when
	// the application last terminated (see the README for when that
	// happens). Recreate it from the saved counter.
	app.Window.OnRestore(func(id string, state application.RestorationState) application.Window {
		if id != "main" {
			return nil
		}
		value, _ := strconv.ParseInt(state.Get("counter"), 10, 64)
		counter.Store(value)
		app.Logger.Info("restoring the main window", "counter", value)
		return newMainWindow(app, popover)
	})

	// Create the default window only when nothing is being restored, or
	// a relaunch would show two main windows.
	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
		if app.Window.WillRestoreWindows() {
			return
		}
		newMainWindow(app, popover)
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

// newMainWindow creates the restorable main window with its toolbar.
func newMainWindow(app *application.App, popover *application.MacPopover) *application.WebviewWindow {
	toolbar := application.NewMacToolbar()
	sheetItem := toolbar.AddButton("Sheet").SetSymbol("rectangle.bottomthird.inset.filled").SetTooltip("Present a sheet")
	popoverItem := toolbar.AddButton("Popover").SetSymbol("bubble.left").SetTooltip("Show a popover")
	toolbar.AddFlexibleSpace()
	countItem := toolbar.AddButton("Count").SetSymbol("plus.circle").SetTooltip("Increment the restored counter")
	kioskItem := toolbar.AddButton("Kiosk").SetSymbol("rectangle.inset.filled").SetTooltip("Toggle kiosk presentation")

	// A toolbar passed in the options is validated when the window is
	// created, so every button needs its OnClick before NewWithOptions. The
	// closures capture the window and sheet variables assigned just below.
	var window, sheet *application.WebviewWindow
	sheetItem.OnClick(func(*application.Context) {
		if err := window.PresentSheet(sheet); err != nil {
			app.Logger.Error("present sheet", "error", err)
		}
	})
	popoverItem.OnClick(func(*application.Context) {
		if popover.IsShown() {
			popover.Close()
			return
		}
		if err := popoverItem.ShowPopover(popover); err != nil {
			app.Logger.Error("toolbar popover", "error", err)
		}
	})
	countItem.OnClick(func(*application.Context) {
		increment(app, window)
	})
	kioskItem.OnClick(func(*application.Context) {
		toggleKiosk(app)
	})

	window = app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:   "main",
		Title:  "Windows Extra",
		Width:  760,
		Height: 520,
		URL:    "/",
		Mac: application.MacWindow{
			RestorationID: "main",
			TitleBar:      application.MacTitleBar{ToolbarStyle: application.MacToolbarStyleUnified},
			Toolbar:       toolbar,
		},
	})
	window.SetRestorationData(map[string]string{"counter": strconv.FormatInt(counter.Load(), 10)})
	sheet = newSheetWindow(app, window)

	// The page asks for the counter when its runtime is ready and increments
	// it from its own button; both paths go through increment/publish.
	window.OnWindowEvent(events.Common.WindowRuntimeReady, func(*application.WindowEvent) {
		publish(app, window)
	})
	app.Event.On("counter:increment", func(*application.CustomEvent) {
		increment(app, window)
	})
	app.Event.On("popover:show", func(event *application.CustomEvent) {
		// Anchor to the rectangle the page reported, in window points with
		// the origin at the top left.
		rect := application.Rect{}
		if values, ok := event.Data.(map[string]any); ok {
			rect.X = int(number(values["x"]))
			rect.Y = int(number(values["y"]))
			rect.Width = int(number(values["width"]))
			rect.Height = int(number(values["height"]))
		}
		if err := popover.ShowRelativeTo(rect, window, application.MacRectEdgeMaxY); err != nil {
			app.Logger.Error("rect popover", "error", err)
		}
	})
	return window
}

func number(value any) float64 {
	if f, ok := value.(float64); ok {
		return f
	}
	return 0
}

// increment bumps the counter, stores it with the restorable state and
// tells the page.
func increment(app *application.App, window *application.WebviewWindow) {
	value := counter.Add(1)
	window.SetRestorationData(map[string]string{"counter": strconv.FormatInt(value, 10)})
	publish(app, window)
}

func publish(app *application.App, window *application.WebviewWindow) {
	value := counter.Load()
	window.SetSubtitle(fmt.Sprintf("Counter %d (survives relaunch)", value))
	app.Event.Emit("counter", value)
	app.Event.Emit("presentation", app.PresentationOptions().String())
}

// newSheetWindow creates the window presented as a sheet. It is a titled
// window (AppKit removes the title bar buttons while it is a sheet) created
// hidden so it never flashes on screen before it is attached. The page's
// OK and Cancel buttons end the sheet with a response code.
func newSheetWindow(app *application.App, parent *application.WebviewWindow) *application.WebviewWindow {
	sheet := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:      "sheet",
		Title:     "Sheet",
		Width:     420,
		Height:    240,
		URL:       "/sheet.html",
		Hidden:    true,
		Frameless: false,
	})
	sheet.OnSheetEnd(func(code int) {
		app.Logger.Info("sheet ended", "code", code, "parent", parent.Name())
		app.Event.Emit("sheet:ended", code)
	})
	app.Event.On("sheet:end", func(event *application.CustomEvent) {
		sheet.EndSheet(int(number(event.Data)))
	})
	return sheet
}

// newCounterPopover builds the popover shown from the toolbar item and from
// a rectangle in the page: a label and two buttons in a native strip.
func newCounterPopover(app *application.App) *application.MacPopover {
	strip := application.NewMacAccessory(application.MacAccessoryLayoutBottom)
	label := strip.AddLabel("Counter actions")
	strip.AddFlexibleSpace()
	strip.AddSymbolButton("plus").OnClick(func(*application.Context) {
		if window, ok := app.Window.Get("main"); ok {
			if webview, ok := window.(*application.WebviewWindow); ok {
				increment(app, webview)
				label.SetText(fmt.Sprintf("Counter is %d", counter.Load()))
			}
		}
	})
	popover := application.NewMacPopover(application.MacPopoverOptions{
		Width:    280,
		Behavior: application.MacPopoverBehaviorTransient,
		Content:  strip,
	})
	strip.AddButton("Close").OnClick(func(*application.Context) { popover.Close() })
	popover.OnClose(func() {
		app.Logger.Info("counter popover closed")
	})
	return popover
}

// newTrayPopover builds the popover shown from the status item: the native
// replacement for a window positioned under the tray icon.
func newTrayPopover(app *application.App) *application.MacPopover {
	strip := application.NewMacAccessory(application.MacAccessoryLayoutBottom)
	strip.AddLabel("Windows Extra")
	strip.AddFlexibleSpace()
	strip.AddButton("Show Window").OnClick(func(*application.Context) {
		if window, ok := app.Window.Get("main"); ok {
			window.Show().Focus()
		}
	})
	strip.AddButton("Quit").OnClick(func(*application.Context) { app.Quit() })
	return application.NewMacPopover(application.MacPopoverOptions{
		Width:    320,
		Behavior: application.MacPopoverBehaviorTransient,
		Content:  strip,
	})
}

// toggleKiosk switches between the default presentation and the kiosk set.
func toggleKiosk(app *application.App) {
	target := kiosk
	if app.PresentationOptions() != application.MacPresentationDefault {
		target = application.MacPresentationDefault
	}
	if err := app.SetPresentationOptions(target); err != nil {
		app.Logger.Error("presentation options", "error", err)
		return
	}
	app.Event.Emit("presentation", app.PresentationOptions().String())
}

func installMenu(app *application.App, popover *application.MacPopover) {
	menu := app.NewMenu()
	menu.AddRole(application.AppMenu)
	menu.AddRole(application.FileMenu)
	menu.AddRole(application.EditMenu)
	view := menu.AddSubmenu("View")
	view.Add("Toggle Kiosk Presentation").SetAccelerator("CmdOrCtrl+K").OnClick(func(*application.Context) {
		toggleKiosk(app)
	})
	view.Add("Hide Dock Only").OnClick(func(*application.Context) {
		if err := app.SetPresentationOptions(application.MacPresentationAutoHideDock); err != nil {
			app.Logger.Error("presentation options", "error", err)
		}
		app.Event.Emit("presentation", app.PresentationOptions().String())
	})
	view.Add("Invalid Combination (logs an error)").OnClick(func(*application.Context) {
		// Hiding the menu bar without the Dock is rejected in Go before
		// AppKit sees it.
		err := app.SetPresentationOptions(application.MacPresentationHideMenuBar)
		app.Logger.Error("presentation options", "error", err)
	})
	view.Add("Reset Presentation").OnClick(func(*application.Context) {
		_ = app.SetPresentationOptions(application.MacPresentationDefault)
		app.Event.Emit("presentation", app.PresentationOptions().String())
	})
	view.AddSeparator()
	view.Add("Close Popover").OnClick(func(*application.Context) { popover.Close() })
	menu.AddRole(application.WindowMenu)
	app.Menu.Set(menu)
}
