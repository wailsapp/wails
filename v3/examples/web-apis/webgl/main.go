package main

import (
	"embed"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/*
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "WebGL Demo",
		Description: "Demonstrates the WebGL API for 3D graphics",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "WebGL Demo",
		Width:  800,
		Height: 600,
		URL:    "/",
	})

	app.Run()
}
