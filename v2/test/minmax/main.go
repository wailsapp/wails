package main

import (
	wails "github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

func main() {

	// Create application with options
	app := wails.CreateAppWithOptions(&options.App{
		Title:     "minmax",
		Width:     1024,
		Height:    768,
		MinWidth:  800,
		MinHeight: 600,
		MaxWidth:  1280,
		MaxHeight: 1024,
	})

	app.Bind(newBasic())

	app.Run()
}
