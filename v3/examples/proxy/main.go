package main

import (
	"embed"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend
var assets embed.FS

// GreetService is our backend service
type GreetService struct{}

// Greet returns a greeting message
func (s *GreetService) Greet(name string) string {
	return fmt.Sprintf("Hello %s! This message came from Wails backend via proxy.", name)
}

// GetTime returns the current server time
func (s *GreetService) GetTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func main() {
	// Create the application
	app := application.New(application.Options{
		Name:        "Wails Proxy Example",
		Description: "An example demonstrating proxy support for external URLs",

		// Enable debug logging
		LogLevel: slog.LevelDebug,

		Services: []application.Service{
			application.NewService(&GreetService{}),
		},

		Assets: application.AssetOptions{
			// Set the external URL to proxy through the asset server
			ProxyTo: "https://app-local.wails-awesome.io:3000",
		},
	})

	// Create window - it will load from local asset server which proxies to external URL
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Proxy Example - External URL",
		Width:  1024,
		Height: 768,
		// No explicit URL needed - will use local asset server

		DevToolsEnabled:            true,
		DefaultContextMenuDisabled: false, // Enable default context menu
	})

	// Show window
	window.Show()

	// Run the application
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
