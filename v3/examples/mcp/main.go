package main

import (
	"embed"
	"fmt"
	"log"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/mcp"
)

//go:embed assets
var assets embed.FS

// GreetService is a demo service so MCP clients can exercise the
// call_bound_method tool.
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
	return fmt.Sprintf("Hello %s!", name)
}

func (g *GreetService) Add(a, b int) int {
	return a + b
}

func (g *GreetService) Shout(text string) string {
	return strings.ToUpper(text) + "!!!"
}

func main() {
	app := application.New(application.Options{
		Name:        "MCP Demo",
		Description: "A playground app controlled by LLMs over the Model Context Protocol",
		Services: []application.Service{
			application.NewService(&GreetService{}),
			// The Route gives the MCP service a same-origin callback channel
			// on the asset server for JavaScript results.
			application.NewServiceWithOptions(mcp.New(), application.ServiceOptions{
				Route: "/wails-mcp",
			}),
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "MCP Demo",
		Name:   "main",
		Width:  900,
		Height: 700,
		URL:    "/",
	})

	// Log events emitted from the frontend or via the MCP emit_event tool.
	app.Event.On("playground:event", func(event *application.CustomEvent) {
		app.Logger.Info("playground:event received", "data", event.Data)
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
