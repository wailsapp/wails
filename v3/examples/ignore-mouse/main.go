package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name:        "WebviewWindow Demo",
		Description: "A demo of the WebviewWindow API",
		Assets:      application.AlphaAssets,
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Width:             800,
		Height:            600,
		Title:             "Ignore Mouse Example",
		URL:               "https://wails.io",
		IgnoreMouseEvents: false,
	})

	window.SetIgnoreMouseEvents(true)
	log.Println("IgnoreMouseEvents set", window.IsIgnoreMouseEvents())

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
