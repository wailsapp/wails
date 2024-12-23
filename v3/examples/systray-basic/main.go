package main

import (
	_ "embed"
	"github.com/wailsapp/wails/v3/pkg/events"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

func main() {
	app := application.New(application.Options{
		Name:        "Systray Demo",
		Description: "A demo of the Systray API",
		Assets:      application.AlphaAssets,
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
	})

	systemTray := app.NewSystemTray()

	window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Width:         500,
		Height:        500,
		Name:          "Systray Demo Window",
		Frameless:     true,
		AlwaysOnTop:   true,
		Hidden:        true,
		DisableResize: true,
		Windows: application.WindowsWindow{
			HiddenOnTaskbar: true,
		},
		KeyBindings: map[string]func(window *application.WebviewWindow){
			"F12": func(window *application.WebviewWindow) {
				systemTray.OpenMenu()
			},
		},
	})

	// Register a hook to hide the window when the window is closing
	window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		// Hide the window
		window.Hide()
		// Cancel the event so it doesn't get destroyed
		e.Cancel()
	})

	if runtime.GOOS == "darwin" {
		systemTray.SetTemplateIcon(icons.SystrayMacTemplate)
	}

	systemTray.AttachWindow(window).WindowOffset(5)

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
