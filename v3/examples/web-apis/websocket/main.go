package main

import (
	"embed"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/*
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "WebSocket Demo",
		Description: "Demonstrates the WebSocket API for real-time communication",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "WebSocket Demo",
		Width:  800,
		Height: 600,
		URL:    "/",
	})

	app.Run()
}
