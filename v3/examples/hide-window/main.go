package main

import (
	_ "embed"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/icons"
	"log"
	"runtime"
)

func main() {
	app := application.New(application.Options{
		Name:        "Hide Window Demo",
		Description: "A test of Hidden window and display it",
		Assets:      application.AlphaAssets,
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	systemTray := app.NewSystemTray()

	window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Width:         500,
		Height:        800,
		Frameless:     false,
		AlwaysOnTop:   false,
		Hidden:        false,
		DisableResize: false,
		Windows: application.WindowsWindow{
			HiddenOnTaskbar: true,
		},
	})

	window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		window.Hide()
		e.Cancel()
	})

	if runtime.GOOS == "darwin" {
		systemTray.SetTemplateIcon(icons.SystrayMacTemplate)
	}

	// Click Dock icon tigger application show
	app.OnApplicationEvent(events.Mac.ApplicationShouldHandleReopen, func(event *application.ApplicationEvent) {
		println("reopen")
		window.Show()
	})

	myMenu := app.NewMenu()
	myMenu.Add("Show").OnClick(func(ctx *application.Context) {
		window.Show()
	})

	myMenu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	systemTray.SetMenu(myMenu)
	systemTray.OnClick(func() {
		window.Show()
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
