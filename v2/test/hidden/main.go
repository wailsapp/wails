package main

import (
	wails "github.com/wailsapp/wails/v2"
)

func main() {

	// Create application with options
	app := wails.CreateAppWithOptions(&wails.Options{
		Title:       "Hidden Demo",
		Width:       1024,
		Height:      768,
		StartHidden: true,
		Frameless:   true,
	})

	app.Bind(newBasic())

	app.Run()
}
