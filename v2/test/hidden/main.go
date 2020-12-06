package main

import (
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

func main() {

	// Create application with options
	app := wails.CreateAppWithOptions(&options.App{
		Title:       "Hidden Demo",
		Width:       1024,
		Height:      768,
		StartHidden: true,
	})

	app.Bind(newBasic())

	app.Run()
}
