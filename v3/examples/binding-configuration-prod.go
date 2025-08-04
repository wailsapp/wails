// Production Configuration Example
// This example shows how to configure HTTP-only bindings for production use
package main

import (
	"context"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// UserService provides user management functionality
type UserService struct{}

func (u *UserService) GetUserInfo(ctx context.Context, userID string) map[string]interface{} {
	// In production, this would fetch from a database
	return map[string]interface{}{
		"id":       userID,
		"name":     "Production User",
		"role":     "admin",
		"lastSeen": time.Now(),
	}
}

func (u *UserService) ProcessLargeDataset(ctx context.Context, data []map[string]interface{}) []map[string]interface{} {
	// Simulate processing large datasets
	processed := make([]map[string]interface{}, len(data))
	for i, item := range data {
		processed[i] = map[string]interface{}{
			"original": item,
			"processed": true,
			"timestamp": time.Now(),
		}
	}
	return processed
}

func main() {
	app := application.New(application.Options{
		Name:        "Production HTTP Bindings Example",
		Description: "Example showing production configuration for HTTP-only bindings",
		
		// Production binding configuration
		Bindings: application.BindingConfig{
			// Longer timeout for production workloads
			Timeout: 10 * time.Minute,
			
			// Strict CORS configuration for production
			CORS: application.CORSConfig{
				// Enable CORS with strict controls
				Enabled: true,
				
				// Only allow specific production origins
				AllowedOrigins: []string{
					"https://myapp.com",
					"https://app.myapp.com",
					"https://*.myapp.com", // Allow subdomains
				},
				
				// Limited HTTP methods for security
				AllowedMethods: []string{"GET", "POST", "OPTIONS"},
				
				// Minimal required headers only
				AllowedHeaders: []string{
					"Content-Type",
					"x-wails-client-id",
					"x-wails-window-name",
					"Authorization", // For authenticated requests
				},
				
				// Longer cache time for production efficiency
				MaxAge: 24 * time.Hour,
			},
			
			// Enable streaming for large data processing
			EnableStreaming: true,
		},
		
		Services: []application.Service{
			application.NewService(&UserService{}),
		},
	})

	err := app.Run()
	if err != nil {
		app.Logger.Error("Failed to run application", "error", err)
	}
}