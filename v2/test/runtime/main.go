package main

import (
	wails "github.com/wailsapp/wails/v2"
)

type Echo struct {
}

func (e *Echo) Echo(message string) string {
	return message
}

func main() {

	// Create application with options
	app := wails.CreateAppWithOptions(&wails.Options{
		Title:         "Runtime Tester!",
		Width:         850,
		Height:        620,
		DisableResize: false,
		Fullscreen:    false,
		Mac: wails.MacOptions{
			HideTitle:                  true,
			HideTitleBar:               false,
			TitlebarAppearsTransparent: true,
			FullSizeContent:            true,
		},
	})

	// You can also use the simplified call:
	// app := wails.CreateApp("Tester!", 1024, 768)

	// ------------- Assets ------------
	// assets := wails.Assets()

	// textFile, err := assets.Read("hello.txt")
	// if err != nil {
	// 	println("Unable to load asset: hello.txt")
	// }
	// println(textFile)
	// ---------------------------------

	app.Bind(newCalc("te"))
	app.Bind(&Echo{})
	app.Bind(&RuntimeTest{})

	app.Run()
}
