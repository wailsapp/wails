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
		Name:        "Screen Demo",
		Description: "A demo of the Screen API",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.NewWebviewWindowWithOptions(&application.WebviewWindowOptions{
		Title:  "Screen Demo",
		Width:  800,
		Height: 600,
		Assets: application.AssetOptions{
			FS: assets,
		},
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 50,
		},
	})

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
