// MCP example — a playground Wails application controllable by LLM agents.
//
// Build with the mcp tag to start the built-in MCP server automatically:
//
//	WAILS_MCP=1 wails3 dev
//	WAILS_MCP=1 wails3 build
//	go run -tags mcp .
//
// Then connect any MCP client to the logged endpoint, e.g. Claude Code:
//
//	claude mcp add --transport http my-app http://127.0.0.1:9099/mcp
//
// No code changes are required: the MCP server starts as part of the
// application when the build tag is present.
package main

import (
	"embed"
	"fmt"
	"log"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

// GreetService is a demo service so MCP clients can exercise call_bound_method.
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
