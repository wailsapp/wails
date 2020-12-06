package main

import (
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

func main() {

	// Create application with options
	app := wails.CreateAppWithOptions(&options.App{
		Title:         "disable resize",
		Width:         1024,
		Height:        768,
		DisableResize: true,
	})

	app.Bind(newBasic())

	app.Run()
}
