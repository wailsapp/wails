package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

// main is the application's entry point. It creates the Wails application,
// registers the LayoutService (save / open / export through native dialogs)
// and opens the builder window.
func main() {
	app := application.New(application.Options{
		Name:        "ui-builder",
		Description: "A drag & drop UI builder built with Wails",
		Services: []application.Service{
			application.NewService(&LayoutService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// The builder is a three-pane tool (palette / canvas / inspector), so it
	// opens at a comfortable desktop size and refuses to shrink below the point
	// where the panes stop being usable.
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "ui-builder",
		Width:     1320,
		Height:    840,
		MinWidth:  960,
		MinHeight: 600,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 48,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(6, 7, 15),
		URL:              "/",
	})

	// Run the application. This blocks until the application has been exited.
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
