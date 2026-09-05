package main

import (
	"embed"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/*
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "Speech Synthesis Demo",
		Description: "Text-to-speech API",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Speech Synthesis Demo",
		Width:  900,
		Height: 700,
		URL:    "/",
	})
	app.Run()
}
