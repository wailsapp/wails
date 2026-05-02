package main

import (
	"embed"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/*
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "Server-Sent Events Demo",
		Description: "EventSource/SSE API demonstration",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Server-Sent Events Demo",
		Width:  900,
		Height: 700,
		URL:    "/",
	})
	app.Run()
}
