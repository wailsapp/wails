package main

import (
	wails "github.com/wailsapp/wails"
)

var html = `<h1> Basic Template </h1>`

func main() {

	// Initialise the app
	app := wails.CreateApp(&wails.AppConfig{
		Width:  1024,
		Height: 768,
		Title:  "My Project",
		HTML:   html,
	})
	app.Run()
}
