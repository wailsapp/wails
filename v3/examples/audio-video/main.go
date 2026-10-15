package main

import (
	"embed"
	"io/fs"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/media"
)

//go:embed assets
var assets embed.FS

func main() {
	mediaFiles, err := fs.Sub(assets, "assets")
	if err != nil {
		log.Fatal(err)
	}
	mediaHandler, err := media.NewHandler(mediaFiles, 32<<20)
	if err != nil {
		log.Fatal(err)
	}
	app := application.New(application.Options{
		Name: "Local media",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})
	app.HandleStream("media", mediaHandler)
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Local audio and video", Width: 800, Height: 650,
	})
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
