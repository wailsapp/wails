package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/*
var assets embed.FS

func main() {
	// Create WebSocket transport on port 9099
	wsTransport := NewWebSocketTransport(":9099")

	app := application.New(application.Options{
		Name:        "WebSocket Transport Example",
		Description: "Example demonstrating custom WebSocket-based IPC transport with event support",
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
		// Events are automatically forwarded to the transport since it implements EventTransport
		Transport: wsTransport,
	})

	// ✨ NO MANUAL EVENT WIRING NEEDED! ✨
	// Events are automatically forwarded to the WebSocket transport because it implements
	// the EventTransport interface. The following code is no longer necessary:
	//
	// app.Events.On("greet:count", func(event *application.WailsEvent) {
	//     wsTransport.SendEvent(event)
	// })
	//
	// All events emitted via app.Events.Emit() are automatically broadcast to connected clients!

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
