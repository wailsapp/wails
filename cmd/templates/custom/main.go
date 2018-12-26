package main

import (
	wails "github.com/wailsapp/wails"
)

var html = `<h1> Custom Project </h1>`

func main() {

	// Initialise the app
	app := wails.CreateApp(&wails.AppConfig{
		Width:  800,
		Height: 600,
		Title:  "My Project",
		HTML:   html,
	})
	app.Run()
}
