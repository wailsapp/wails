package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name:        "Show macOS Toolbar",
		Description: "A demo of the ShowToolbarWhenFullscreen option",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create window
	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "Toolbar hidden (default behaviour)",
		HTML:  "<html><body><h1>Switch this window to fullscreen: the toolbar will be hidden</h1></body></html>",
		CSS:   `body { background-color: blue; color: white; height: 100vh; display: flex; justify-content: center; align-items: center; }`,
		Mac: application.MacWindow{
			TitleBar: application.MacTitleBar{
				UseToolbar:           true,
				HideToolbarSeparator: true,
			},
		},
	})

	// Create window
	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "Toolbar visible",
		HTML:  "<html><body><h1>Switch this window to fullscreen: the toolbar will stay visible</h1></body></html>",
		CSS:   `body { background-color: red; color: white; height: 100vh; display: flex; justify-content: center; align-items: center; }`,
		Mac: application.MacWindow{
			TitleBar: application.MacTitleBar{
				UseToolbar:                true,
				HideToolbarSeparator:      true,
				ShowToolbarWhenFullscreen: true,
			},
		},
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
