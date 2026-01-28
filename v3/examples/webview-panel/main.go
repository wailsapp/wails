package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed assets/*
var assets embed.FS

// App struct holds application state
type App struct {
	panel *application.WebviewPanel
}

// SwitchPanel switches the embedded panel to a different URL
func (a *App) SwitchPanel(url string) {
	if a.panel != nil {
		log.Printf("🔄 Switching panel to: %s", url)
		a.panel.SetURL(url)
	}
}

func main() {
	appState := &App{}

	app := application.New(application.Options{
		Name:        "WebviewPanel Demo",
		Description: "Demonstrates embedding multiple webviews with switching capability",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Services: []application.Service{
			application.NewService(appState),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create the main window with our custom UI
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "WebviewPanel Demo - Embedded Webviews",
		Width:            1200,
		Height:           800,
		BackgroundType:   application.BackgroundTypeSolid,
		BackgroundColour: application.NewRGB(26, 26, 46),
		URL:              "/index.html",
	})

	// Create the embedded panel after the window is shown
	window.OnWindowEvent(events.Common.WindowShow, func(*application.WindowEvent) {
		// Based on actual HTML measurements:
		// - Header: 41px height (with 15px/25px padding already included)
		// - Content area: padding 20px
		// - Panel container: 1142×591, border 1px
		// - Panel should fit inside container: 1140×589
		const (
			headerHeight    = 41 + 15*2 // 41px header + vertical padding
			contentPadding  = 20
			containerBorder = 1
			// Panel container inner size (container size minus borders)
			panelContainerWidth  = 1140 // 1142 - 2
			panelContainerHeight = 589  // 591 - 2
		)

		// Panel position: content padding + container border
		panelX := contentPadding + containerBorder
		panelY := headerHeight + contentPadding + containerBorder

		// Create a panel positioned inside the content area
		// Using AnchorFill means the panel will maintain these margins when the window is resized
		appState.panel = window.NewPanel(application.WebviewPanelOptions{
			Name:            "embedded-content",
			URL:             "https://wails.io",
			X:               panelX,
			Y:               panelY,
			Width:           panelContainerWidth,
			Height:          panelContainerHeight,
			Anchor:          application.AnchorFill, // Maintain margins on all sides when resizing
			Visible:         boolPtr(true),
			DevToolsEnabled: boolPtr(true),
		})
	})

	// Run the application
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func boolPtr(b bool) *bool {
	return &b
}
