package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/*
var assets embed.FS

// App owns navigation policy for the embedded browser.
type App struct{ panel *application.WebviewPanel }

// SwitchPanel selects one of the destinations allowed by this application.
// The frontend cannot supply an arbitrary URL or execute code in the panel.
func (a *App) SwitchPanel(site string) error {
	destinations := map[string]string{
		"local":  "/panel.html",
		"wails":  "https://wails.io",
		"google": "https://www.google.com",
	}
	url, ok := destinations[site]
	if !ok {
		return fmt.Errorf("unknown site: %s", site)
	}
	if a.panel == nil {
		return fmt.Errorf("panel is unavailable")
	}
	a.panel.SetURL(url)
	return nil
}

func main() {
	state := &App{}
	app := application.New(application.Options{
		Name:        "WebviewPanel Demo",
		Description: "Embedded native webviews with responsive layout",
		Assets:      application.AssetOptions{Handler: application.BundledAssetFileServer(assets)},
		Services:    []application.Service{application.NewService(state)},
		Mac:         application.MacOptions{ApplicationShouldTerminateAfterLastWindowClosed: true},
	})
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "WebviewPanel Demo", Width: 1200, Height: 800,
		BackgroundColour: application.NewRGB(26, 26, 46), URL: "/index.html",
	})
	// Configure before Run: a repeated show event must not create another panel.
	// The frontend measures the placeholder, sets bounds, and then shows the view.
	visible := false
	state.panel = window.NewPanel(application.WebviewPanelOptions{
		Name: "embedded-content", URL: "/panel.html", Visible: &visible,
		BackgroundColour: application.NewRGB(248, 250, 252),
	})
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
