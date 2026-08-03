package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "mac-split-sidebar",
		Description: "A demo of MacSplitWindow: sidebar, content and inspector webview panes",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	splitWindow, err := application.NewSplitWindow(application.MacSplitWindowOptions{
		Title:        "Reader",
		Width:        980,
		Height:       640,
		Orientation:  application.SplitHorizontal,
		AutosaveName: "mac-split-sidebar-example",
		Panes: []application.MacSplitPaneOptions{
			{
				Name:                       "sidebar",
				Content:                    application.MacWebviewContent{URL: "/sidebar.html"},
				Behavior:                   application.PaneBehaviorSidebar,
				MinThickness:               180,
				MaxThickness:               280,
				PreferredThicknessFraction: 0.22,
				Collapsible:                true,
			},
			{
				Name:     "content",
				Content:  application.MacWebviewContent{URL: "/content.html"},
				Behavior: application.PaneBehaviorDefault,
			},
			{
				Name:                       "inspector",
				Content:                    application.MacWebviewContent{URL: "/inspector.html"},
				Behavior:                   application.PaneBehaviorInspector,
				MinThickness:               220,
				MaxThickness:               320,
				PreferredThicknessFraction: 0.22,
				Collapsible:                true,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// The sidebar-toggle toolbar item needs no pane-tracking wiring: it maps
	// straight to AppKit's own toggleSidebar: action, which
	// NSSplitViewController already implements and the responder chain
	// finds without any Go-side involvement.
	splitWindow.Window().SetToolbar(&application.MacToolbar{
		Items: []application.MacToolbarItem{
			{ID: "sidebar-toggle", Kind: application.ToolbarSidebarToggle},
			{ID: "flex", Kind: application.ToolbarFlexibleSpace},
		},
	})

	splitWindow.Show()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
