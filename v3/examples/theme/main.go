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
		// We Start With Dark Theme
		Theme: application.AppDark,
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

	windowService.app = app
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Window 1",
		Name:  "Window 1",
		// Both Mac and Windows will follow light theme
		Mac: application.MacWindow{
			Appearance: "NSAppearanceNameAqua",
		},
		Windows: application.WindowsWindow{
			Theme: "light",
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Window 2",
		Name:  "Window 2",
		// Both Mac and Widnows will follow Application Theme
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Window 3",
		Name:  "Window 3",
		// Both Mac and Widnows will follow Dark Theme
		Mac: application.MacWindow{
			Appearance: "NSAppearanceNameDarkAqua",
		},
		Windows: application.WindowsWindow{
			Theme: "dark",
		},
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
