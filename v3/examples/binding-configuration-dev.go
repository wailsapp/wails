// Development Configuration Example
// This example shows how to configure HTTP-only bindings for development use
package main

import (
	"context"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// GreetService provides greeting functionality
type GreetService struct{}

func (g *GreetService) Greet(ctx context.Context, name string) string {
	return "Hello " + name + ", it's development time!"
}

func main() {
	app := application.New(application.Options{
		Name:        "Development HTTP Bindings Example",
		Description: "Example showing development configuration for HTTP-only bindings",
		
		// Development binding configuration
		Bindings: application.BindingConfig{
			// Shorter timeout for development
			Timeout: 5 * time.Minute,
			
			// CORS configuration for development
			CORS: application.CORSConfig{
				// Enable CORS for development
				Enabled: true,
				
				// Allow all origins in development (empty list = allow all)
				AllowedOrigins: []string{},
				
				// Standard HTTP methods
				AllowedMethods: []string{"GET", "POST", "OPTIONS"},
				
				// Standard headers plus development-specific ones
				AllowedHeaders: []string{
					"Content-Type",
					"x-wails-client-id",
					"x-wails-window-name",
					"x-dev-token", // Development-specific header
				},
				
				// Shorter cache time for development
				MaxAge: 1 * time.Hour,
			},
			
			// Enable streaming for testing large payloads
			EnableStreaming: true,
		},
		
		Services: []application.Service{
			application.NewService(&GreetService{}),
		},
	})

	err := app.Run()
	if err != nil {
		app.Logger.Error("Failed to run application", "error", err)
	}
}