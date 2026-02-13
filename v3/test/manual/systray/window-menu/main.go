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
		Name:        "Systray Window + Menu",
		Description: "Tests systray with both attached window and menu",
		Assets:      application.AlphaAssets,
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
	})

	systray := app.SystemTray.New()

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Width:         400,
		Height:        300,
		Name:          "Window Menu Test",
		Title:         "Window + Menu - Left-click toggles, Right-click shows menu",
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
	menu.Add("Show Window").OnClick(func(ctx *application.Context) {
		systray.ShowWindow()
	})
	menu.Add("Hide Window").OnClick(func(ctx *application.Context) {
		systray.HideWindow()
	})
	menu.AddSeparator()
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	if runtime.GOOS == "darwin" {
		systray.SetTemplateIcon(icons.SystrayMacTemplate)
	}

	systray.AttachWindow(window).WindowOffset(5).SetMenu(menu)

	log.Println("Window + Menu test started")
	log.Println("Expected behavior:")
	log.Println("  - Left-click: Toggle window visibility")
	log.Println("  - Right-click: Show menu")

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
