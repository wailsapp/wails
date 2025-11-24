package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts up
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Simple echo method for testing
func (a *App) Echo(message string) map[string]interface{} {
	return map[string]interface{}{
		"echo":      message,
		"timestamp": time.Now().Format(time.RFC3339),
		"from":      "HTTP-only bindings",
	}
}

// Add two numbers
func (a *App) Add(x, y float64) map[string]interface{} {
	return map[string]interface{}{
		"result": x + y,
		"operation": fmt.Sprintf("%.2f + %.2f = %.2f", x, y, x+y),
	}
}

// Simulate a long-running operation
func (a *App) SlowOperation(ctx context.Context, seconds int) map[string]interface{} {
	// Simulate work with cancellation support
	for i := 0; i < seconds; i++ {
		select {
		case <-ctx.Done():
			return map[string]interface{}{
				"status": "cancelled",
				"progress": fmt.Sprintf("%d/%d seconds", i, seconds),
			}
		default:
			time.Sleep(1 * time.Second)
		}
	}
	
	return map[string]interface{}{
		"status": "completed",
		"duration": fmt.Sprintf("%d seconds", seconds),
	}
}

// Return large data to test performance
func (a *App) GetLargeData() map[string]interface{} {
	data := make(map[string]interface{})
	data["metadata"] = map[string]interface{}{
		"size": 1000,
		"generated": time.Now().Format(time.RFC3339),
		"type": "performance_test",
	}
	
	items := make([]map[string]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		items[i] = map[string]interface{}{
			"id":   i,
			"name": fmt.Sprintf("Item %d", i),
			"data": fmt.Sprintf("Some data for item %d with extra content to increase size", i),
		}
	}
	data["items"] = items
	
	return data
}

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	wailsApp := application.New(application.Options{
		Name:        "HTTP-Only Bindings Demo",
		Description: "A demo app showcasing HTTP-only bindings with CORS support",
		
		Services: []application.Service{},
		
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(embedFiles),
			
			// Configure HTTP-only bindings
			Bindings: application.BindingConfig{
				// 10 minute timeout for long operations
				Timeout: 10 * time.Minute,
				
				// Enable CORS for external URLs
				CORS: application.CORSConfig{
					Enabled: true,
					AllowedOrigins: []string{
						"https://example.com",
						"https://*.example.com",
						"http://localhost:3000",
						"http://localhost:8080",
						"http://127.0.0.1:3000",
						"http://127.0.0.1:8080",
					},
					AllowedMethods: []string{"GET", "POST", "OPTIONS"},
					AllowedHeaders: []string{
						"Content-Type",
						"x-wails-client-id",
						"x-wails-window-name",
						"x-wails-window-id",
					},
					MaxAge: 24 * time.Hour,
				},
			},
		},

		OnStartup: app.Startup,
	})

	// Bind methods
	err := wailsApp.NewWebviewWindow(application.WebviewWindowOptions{
		Title: "HTTP-Only Bindings Demo",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})
	
	if err != nil {
		log.Fatal(err)
	}

	err = wailsApp.Run()

	if err != nil {
		println("Error:", err.Error())
	}
}

//go:embed all:frontend/dist
var embedFiles embed.FS