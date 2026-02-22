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
		Name:        "Systray Window Only",
		Description: "Tests systray with attached window only (no menu)",
		Assets:      application.AlphaAssets,
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
	})

	systray := app.SystemTray.New()

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Width:         400,
		Height:        300,
		Name:          "Window Only Test",
		Title:         "Window Only - Left-click systray to toggle",
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

	if runtime.GOOS == "darwin" {
		systray.SetTemplateIcon(icons.SystrayMacTemplate)
	}

	systray.AttachWindow(window).WindowOffset(5)

	log.Println("Window-only test started")
	log.Println("Expected behavior:")
	log.Println("  - Left-click: Toggle window visibility")
	log.Println("  - Right-click: Nothing (no menu)")

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
