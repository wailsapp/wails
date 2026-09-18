package main

import (
	"embed"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/*
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "Canvas 2D Demo",
		Description: "Demonstrates the Canvas 2D API for graphics rendering",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Canvas 2D Demo",
		Width:  800,
		Height: 600,
		URL:    "/",
	})

	app.Run()
}
