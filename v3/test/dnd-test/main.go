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
		Name:        "DnD Test",
		Description: "Test drag and drop behavior - internal vs external",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	win := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:                  "DnD Test",
		Width:                  900,
		Height:                 700,
		EnableFileDrop:         true,
		OpenInspectorOnStartup: true,
	})

	// Listen for file drop events
	win.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
		files := event.Context().DroppedFiles()
		details := event.Context().DropTargetDetails()

		log.Printf("[Go] Files dropped: %v", files)
		if details != nil {
			log.Printf("[Go] Drop target: id=%s, classes=%v, x=%d, y=%d",
				details.ElementID, details.ClassList, details.X, details.Y)
		}
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
