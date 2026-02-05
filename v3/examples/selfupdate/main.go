package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/selfupdate"
)

//go:embed assets/*
var assets embed.FS

// Version is set at build time using:
// go build -ldflags "-X main.Version=1.0.0"
var Version = "0.0.1-dev"

func main() {
	// Create the selfupdate service
	updateService := selfupdate.New(&selfupdate.Config{
		CurrentVersion: Version,
		Provider:       "github",
		GitHub: &selfupdate.GitHubConfig{
			Owner: "wailsapp",  // Replace with your GitHub org/user
			Repo:  "wails",     // Replace with your repository name
			// Token: os.Getenv("GITHUB_TOKEN"), // Optional: for private repos
		},
		// AutoCheck: true, // Uncomment to check for updates on startup
	})

	app := application.New(application.Options{
		Name:        "Self-Update Example",
		Description: "Demonstrates the selfupdate service",
		Services: []application.Service{
			application.NewService(updateService),
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:           "Self-Update Example",
		Width:           600,
		Height:          450,
		URL:             "/",
		DevToolsEnabled: true,
	})

	// Listen for update available event (emitted when AutoCheck is true)
	app.Event.On("selfupdate:available", func(event *application.CustomEvent) {
		log.Printf("Update available: %v", event.Data)
	})

	// Listen for download progress
	app.Event.On("selfupdate:progress", func(event *application.CustomEvent) {
		log.Printf("Download progress: %v", event.Data)
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
