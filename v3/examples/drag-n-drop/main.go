package main

import (
	"embed"
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed assets
var assets embed.FS

func main() {

	app := application.New(application.Options{
		Name:        "Drag-n-drop Demo",
		Description: "A demo of the Drag-n-drop API",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	window1 := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "Window Level Drag and Drop",
		URL:   "/window1/index.html",
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
		DragAndDrop: application.DragAndDropTypeWindow,
	})

	window2 := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "Webview Level Drag and Drop",
		URL:   "/window2/index.html",
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
		DragAndDrop: application.DragAndDropTypeWebview,
	})

	handleDraggedFiles := func(event *application.WindowEvent) {
		ctx := event.Context()
		files := ctx.DroppedFiles()
		x, y := ctx.Location()
		app.Events.Emit(&application.WailsEvent{
			Name: "files",
			Data: files,
		})
		app.Logger.Info("Files Dropped!", "x", x, "y", y, "files", files)
	}

	window1.On(events.Common.WindowFilesDropped, handleDraggedFiles)
	window2.On(events.Common.WindowFilesDropped, handleDraggedFiles)

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
