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
		URL:    "/",
		Mac: application.MacWindow{
			TitleBar: application.MacTitleBar{
				ToolbarStyle: application.MacToolbarStyleExpanded,
			},
		},
	})

	toolbar := newDaymarkToolbar(app)
	window.SetToolbar(toolbar.NativeToolbar())
	installDaymarkMenu(app, toolbar)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
