package main

import (
	"embed"
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/*
var assets embed.FS

func main() {

	app := application.New(application.Options{
		Name:        "Wails ML Demo",
		Description: "A demo of the Wails ML API",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:  "Wails ML Demo",
		Width:  1280,
		Height: 1024,
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
	})

	app.OnEvent("button-pressed", func(_ *application.CustomEvent) {
		println("Button Pressed!")
	})
	app.OnEvent("hover", func(_ *application.CustomEvent) {
		println("Hover time!")
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
