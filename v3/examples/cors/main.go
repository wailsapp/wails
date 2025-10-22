package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend
var assets embed.FS

// GreetService is our backend service
type GreetService struct{}

// Greet returns a greeting message
func (s *GreetService) Greet(name string) string {
	return fmt.Sprintf("Hello %s! This message came from Wails backend via CORS.", name)
}

// GetTime returns the current server time
func (s *GreetService) GetTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// TestCORS tests if CORS is working correctly
func (s *GreetService) TestCORS(ctx context.Context) map[string]interface{} {
	// Get the window from context if available
	window := ctx.Value(application.WindowKey)

	return map[string]interface{}{
		"cors_enabled": true,
		"backend_time": time.Now().Unix(),
		"message":      "CORS is working! This response came from the Wails backend.",
		"window_info":  fmt.Sprintf("%v", window != nil),
	}
}

func main() {
	// Create the application
	app := application.New(application.Options{
		Name:        "Wails CORS Example",
		Description: "An example demonstrating CORS support with external URLs",

		// Configure CORS to allow our external URL
		CORS: application.CORSConfig{
			Enabled: true,
			AllowedOrigins: []string{
				"https://app-local.wails-awesome.io:3000",  // Our external URL
				"https://localhost:3000",                    // Alternative localhost
				"http://localhost:5173",                     // Vite dev server (if using)
			},
			AllowedMethods: []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders: []string{
				"Content-Type",
				"X-Wails-Window-ID",
				"X-Wails-Window-Name",
				"X-Wails-Client-ID",
				"Authorization",
			},
			MaxAge: 5 * time.Minute,
		},

		Services: []application.Service{
			application.NewService(&GreetService{}),
		},

		Assets: application.AssetOptions{
			// Use embedded assets as fallback
			Handler: application.AssetFileServerFS(assets),
		},
	})

	// Create window with external URL
	window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "CORS Example - External URL",
		Width:  1024,
		Height: 768,
		URL:    "https://app-local.wails-awesome.io:3000", // External URL

		DevToolsEnabled:       true,
		DefaultContextMenuEnabled: true,
	})

	// Show window
	window.Show()

	// Run the application
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}