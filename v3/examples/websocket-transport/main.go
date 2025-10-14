package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/*
var assets embed.FS

func main() {
	// Create WebSocket transport on port 9099 with default Base64/JSON codec
	// To use a different codec (e.g., raw JSON), use:
	// wsTransport := application.NewWebSocketTransport(":9099", application.WithCodec(application.NewRawJSONCodec()))
	wsTransport := application.NewWebSocketTransport(":9099")

	app := application.New(application.Options{
		Name:        "WebSocket Transport Example",
		Description: "Example demonstrating custom WebSocket-based IPC transport",
		Services: []application.Service{
			application.NewService(&GreetService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		// Use WebSocket transport instead of default HTTP
		Transport: wsTransport,
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:           "WebSocket Transport Example",
		URL:             "/",
		Width:           800,
		Height:          600,
		DevToolsEnabled: true,
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
