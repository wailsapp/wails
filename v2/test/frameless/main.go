package main

import (
	wails "github.com/wailsapp/wails/v2"
)

func main() {

	// Create application with options
	app := wails.CreateAppWithOptions(&wails.Options{
		Title:     "Frameless Demo",
		Width:     1024,
		Height:    768,
		MinWidth:  800,
		MinHeight: 600,
		MaxWidth:  1280,
		MaxHeight: 1024,
		Frameless: true,
	})

	app.Bind(newBasic())

	app.Run()
}
