//go:build ios

package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/*
var assets embed.FS

// App struct for binding methods
type App struct{}

// Greet returns a greeting message
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s from Wails v3 on iOS!", name)
}

func main() {
	// Create application with options
	app := application.New(application.Options{
		Name:        "Wails iOS PoC",
		Description: "Proof of concept for Wails v3 on iOS",
		Assets: application.AssetOptions{
			FS: assets,
		},
		Services: []application.Service{
			application.NewService(&App{}),
		},
		LogLevel: application.LogLevelDebug,
	})

	// Run the application
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}