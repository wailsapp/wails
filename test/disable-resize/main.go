package main

import (
	wails "github.com/wailsapp/wails/v2"
)

func main() {

	// Create application with options
	app := wails.CreateAppWithOptions(&wails.Options{
		Title:         "disable resize",
		Width:         1024,
		Height:        768,
		DisableResize: true,
	})

	app.Bind(newBasic())

	app.Run()
}
