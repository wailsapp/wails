package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed assets
var assets embed.FS

func main() {

	app := application.New(application.Options{
		Name:        "Mouse Hover Demo",
		Description: "A demo of WindowMouseEnter and WindowMouseLeave events",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a main window that demonstrates mouse enter/leave events
	mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Mouse Hover Demo - Main Window",
		Width:  600,
		Height: 400,
	})

	// Log mouse enter/leave events for the main window
	mainWindow.OnWindowEvent(events.Common.WindowMouseEnter, func(e *application.WindowEvent) {
		app.Logger.Info("Main Window: Mouse entered!")
	})

	mainWindow.OnWindowEvent(events.Common.WindowMouseLeave, func(e *application.WindowEvent) {
		app.Logger.Info("Main Window: Mouse left!")
	})

	// Create a secondary "popup-like" window that auto-focuses on mouse enter
	// This is useful for tray popup windows where you want immediate interaction
	popupWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:             "Auto-Focus Window (hover to focus)",
		Width:             300,
		Height:            200,
		X:                 650,
		Y:                 100,
		FocusOnMouseEnter: true, // Automatically focus when mouse enters
	})

	// Log mouse enter/leave events for the popup window
	popupWindow.OnWindowEvent(events.Common.WindowMouseEnter, func(e *application.WindowEvent) {
		app.Logger.Info("Popup Window: Mouse entered! (window auto-focused)")
	})

	popupWindow.OnWindowEvent(events.Common.WindowMouseLeave, func(e *application.WindowEvent) {
		app.Logger.Info("Popup Window: Mouse left!")
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
