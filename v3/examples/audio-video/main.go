package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist/*
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "Audio/Video Example",
		Description: "A demo of HTML5 Audio/Video with the Wails runtime npm module",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Audio/Video Example",
		Width:  900,
		Height: 700,
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
