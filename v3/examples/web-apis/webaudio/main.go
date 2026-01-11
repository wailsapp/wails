package main

import (
	"embed"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/*
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "Web Audio Demo",
		Description: "Demonstrates the Web Audio API for audio processing",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Web Audio Demo",
		Width:  800,
		Height: 600,
		URL:    "/",
	})

	app.Run()
}
