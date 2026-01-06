package main

import (
	_ "embed"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

func main() {
	app := application.New(application.Options{
		Name:        "Systray Custom Handlers",
		Description: "Tests systray with custom click handlers overriding defaults",
		Assets:      application.AlphaAssets,
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
	})

	systray := app.SystemTray.New()

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Width:         400,
		Height:        300,
		Name:          "Custom Handlers Test",
		Title:         "Custom Handlers - Check console for click events",
		Frameless:     true,
		AlwaysOnTop:   true,
		Hidden:        true,
		DisableResize: true,
		Windows: application.WindowsWindow{
			HiddenOnTaskbar: true,
		},
	})

	window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		window.Hide()
		e.Cancel()
	})

	menu := app.Menu.New()
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	if runtime.GOOS == "darwin" {
		systray.SetTemplateIcon(icons.SystrayMacTemplate)
	}

	systray.AttachWindow(window).WindowOffset(5).SetMenu(menu)

	systray.OnClick(func() {
		log.Println("Custom left-click handler called!")
		log.Println("  -> Toggling window (custom behavior with logging)")
		systray.ToggleWindow()
	})

	systray.OnRightClick(func() {
		log.Println("Custom right-click handler called!")
		log.Println("  -> Opening menu (custom behavior)")
		systray.OpenMenu()
	})

	systray.OnDoubleClick(func() {
		log.Println("Double-click detected!")
	})

	log.Println("Custom handlers test started")
	log.Println("Expected behavior:")
	log.Println("  - Left-click: Custom handler logs + toggles window")
	log.Println("  - Right-click: Custom handler logs + opens menu")
	log.Println("  - Double-click: Logs to console")

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
