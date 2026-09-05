package main

import (
	"embed"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/*
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "Streams API Demo",
		Description: "Readable and writable streams",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Streams API Demo",
		Width:  900,
		Height: 700,
		URL:    "/",
	})
	app.Run()
}
