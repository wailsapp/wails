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
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Assets: application.AssetOptions{
			FS: assets,
		},
	})

	window := app.NewWebviewWindowWithOptions(&application.WebviewWindowOptions{
		Title: "Drag-n-drop Demo",
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
		EnableDragAndDrop: true,
	})

	window.On(events.FilesDropped, func(ctx *application.WindowEventContext) {
		files := ctx.DroppedFiles()
		app.Events.Emit(&application.WailsEvent{
			Name: "files",
			Data: files,
		})
		log.Printf("[Go] FilesDropped received: %+v\n", files)
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
