package main

import (
	"embed"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/*
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "Clipboard API Demo",
		Description: "Demonstrates the Clipboard API for copy/paste",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Clipboard API Demo",
		Width:  800,
		Height: 600,
		URL:    "/",
	})

	app.Run()
}
