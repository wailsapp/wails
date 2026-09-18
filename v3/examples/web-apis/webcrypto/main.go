package main

import (
	"embed"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/*
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "Web Crypto Demo",
		Description: "Demonstrates the Web Cryptography API",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Web Crypto Demo",
		Width:  800,
		Height: 600,
		URL:    "/",
	})

	app.Run()
}
