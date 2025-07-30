package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "HTTP CORS Workaround Example",
		Description: "Demonstrates using Wails HTTP API to bypass CORS",
		Assets: application.AssetOptions{
			FS: assets,
		},
	})

	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "HTTP CORS Workaround",
		Width:  800,
		Height: 600,
		URL:    "/",
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}