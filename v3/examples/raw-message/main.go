package main

import (
	"embed"
	_ "embed"
	"fmt"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

// main creates and runs the Wails application "Raw Message Demo" with embedded assets,
// a single webview window named "Window 1", and a RawMessageHandler that logs incoming
// raw messages along with their origin information. It logs and exits if the app fails to run.
func main() {

	app := application.New(application.Options{
		Name:        "Raw Message Demo",
		Description: "A demo of sending raw messages from the frontend",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
			println(fmt.Sprintf("Raw message received from Window %s with message: %s, origin %s, topOrigin %s, isMainFrame %t", window.Name(), message, originInfo.Origin, originInfo.TopOrigin, originInfo.IsMainFrame))
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Window 1",
		Name:  "Window 1",
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}