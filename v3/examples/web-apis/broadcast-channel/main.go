package main

import (
	"embed"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/*
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "Broadcast Channel Demo",
		Description: "Cross-window communication via BroadcastChannel API",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Broadcast Channel Demo",
		Width:  800,
		Height: 550,
		URL:    "/",
	})

	app.Run()
}
