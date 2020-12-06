package main

import (
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

func main() {

	// Create application with options
	app := wails.CreateAppWithOptions(&options.App{
		Title:     "minmax",
		Width:     800,
		Height:    600,
		MinWidth:  400,
		MinHeight: 300,
		MaxWidth:  1024,
		MaxHeight: 768,
	})

	app.Bind(newBasic())

	app.Run()
}
