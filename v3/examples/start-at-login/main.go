package main

import (
	"fmt"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// GreetService is a simple service that provides greeting functionality
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
	return fmt.Sprintf("Hello %s! Welcome to the Start At Login demo.", name)
}

// ToggleStartAtLogin allows the frontend to toggle the start at login setting
func (g *GreetService) ToggleStartAtLogin() (bool, error) {
	app := application.Get()
	
	// Check current status
	isEnabled, err := app.StartsAtLogin()
	if err != nil {
		return false, fmt.Errorf("failed to check start at login status: %w", err)
	}
	
	// Toggle the setting
	newStatus := !isEnabled
	if err := app.SetStartAtLogin(newStatus); err != nil {
		return false, fmt.Errorf("failed to set start at login: %w", err)
	}
	
	return newStatus, nil
}

// GetStartAtLoginStatus returns the current start at login status
func (g *GreetService) GetStartAtLoginStatus() (bool, error) {
	app := application.Get()
	return app.StartsAtLogin()
}

func main() {
	app := application.New(application.Options{
		Name:        "Start At Login Demo",
		Description: "A demo application showing how to use the Start At Login feature",
		Services: []application.Service{
			application.NewService(&GreetService{}),
		},
		Assets: application.AlphaAssets,
		// Uncomment the line below to enable start at login when the app first runs
		// StartAtLogin: true,
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}