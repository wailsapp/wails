package main

import (
	_ "embed"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/icons"
	"log"
	"runtime"
)

var windowShowing bool

func createWindow(app *application.App) {
	if windowShowing {
		return
	}
	// Log the time taken to create the window
	window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Width:            500,
		Height:           500,
		Name:             "Systray Demo Window",
		AlwaysOnTop:      true,
		Hidden:           true,
		BackgroundColour: application.NewRGB(33, 37, 41),
		DisableResize:    true,
		Windows: application.WindowsWindow{
			HiddenOnTaskbar: true,
		},
	})
	windowShowing = true

	window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
		windowShowing = false
	})

	window.Show()
}

func main() {
	app := application.New(application.Options{
		Name:        "Systray Demo",
		Description: "A demo of the Systray API",
		Assets:      application.AlphaAssets,
		Windows: application.WindowsOptions{
			DisableQuitOnLastWindowClosed: true,
		},
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
	})

	systemTray := app.NewSystemTray()
	menu := app.NewMenu()
	menu.Add("Quit").OnClick(func(data *application.Context) {
		app.Quit()
	})
	systemTray.SetMenu(menu)

	if runtime.GOOS == "darwin" {
		systemTray.SetTemplateIcon(icons.SystrayMacTemplate)
	}

	systemTray.OnClick(func() {
		createWindow(app)
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
