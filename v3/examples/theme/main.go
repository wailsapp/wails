package main

import (
	"embed"
	// "fmt"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {
	windowService := &WindowService{}
	app := application.New(application.Options{
		Name:        "customEventProcessor Demo",
		Description: "A demo of the customEventProcessor API",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Windows: application.WindowsOptions{},
		Services: []application.Service{
			application.NewService(windowService),
		},
	})

	// Listen for the theme‑change event and log the payload
	// app.Event.On("applicationThemeChanged", func(ev *application.CustomEvent) {
	// 	fmt.Printf("[Go] applicationThemeChanged received, data = %v\n", ev.Data)
	// })

	windowService.app = app
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Window 1",
		Name:  "Window 1",
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Window 2",
		Name:  "Window 2",
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Window 3",
		Name:  "Window 3",
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
