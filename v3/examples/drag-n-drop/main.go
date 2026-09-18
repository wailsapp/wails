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
		Name:        "File Drop Demo",
		Description: "A demo of file drag and drop",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	win := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:          "File Drop Demo",
		Width:          800,
		Height:         600,
		EnableFileDrop: true,
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
	})

	// Listen for file drop events
	win.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
		files := event.Context().DroppedFiles()
		details := event.Context().DropTargetDetails()

		log.Printf("Files dropped: %v", files)
		if details != nil {
			log.Printf("Drop target: id=%s, classes=%v, x=%d, y=%d",
				details.ElementID, details.ClassList, details.X, details.Y)
		}

		// Emit event to frontend
		application.Get().Event.Emit("files-dropped", map[string]any{
			"files":   files,
			"details": details,
		})
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
