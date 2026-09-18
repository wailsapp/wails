package main

import (
	"embed"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/*
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "Permissions API Demo",
		Description: "Permission queries",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Permissions API Demo",
		Width:  900,
		Height: 700,
		URL:    "/",
	})
	app.Run()
}
