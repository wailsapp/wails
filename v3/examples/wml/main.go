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
			FS: assets,
		},
	})

	app.NewWebviewWindowWithOptions(&application.WebviewWindowOptions{
		Title:  "Wails ML Demo",
		Width:  800,
		Height: 600,
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
	})

	app.Events.On("button-pressed", func(_ *application.WailsEvent) {
		println("Button Pressed!")
	})
	app.Events.On("hover", func(_ *application.WailsEvent) {
		println("Hover time!")
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
