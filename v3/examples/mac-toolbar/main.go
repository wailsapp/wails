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
		Name:        "Daymark",
		Description: "A native macOS toolbar editor",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Daymark",
		Width:  1180,
		Height: 760,
		URL:    "/editor.html",
		Mac: application.MacWindow{
			// Finder's window structure is an ordinary opaque AppKit window
			// with edge-to-edge content. AppKit supplies the toolbar material;
			// the document is not made into a custom glass surface.
			Backdrop:      application.MacBackdropNormal,
			ContentLayout: application.MacContentLayoutEdgeToEdge,
			TitleBar: application.MacTitleBar{
				AppearsTransparent:   false,
				FullSizeContent:      true,
				HideToolbarSeparator: true,
				ToolbarStyle:         application.MacToolbarStyleUnified,
			},
		},
	})

	// The split view must be configured before the window is shown, and its
	// installation precedes toolbar attachment so the toolbar's tracking
	// separator can align with the sidebar divider.
	split := newDaymarkSplit(app)
	window.SetSplitView(split.NativeSplitView())

	toolbar := newDaymarkToolbar(app, split)
	window.SetToolbar(toolbar.NativeToolbar())
	installDaymarkMenu(app, toolbar, split)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
