package main

import (
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"log"
)

type Echo struct {
}

func (e *Echo) Echo(message string) string {
	return message
}

func main() {

	// Create application with options
	app, err := wails.CreateAppWithOptions(&options.App{
		Title:         "Runtime Tester!",
		Width:         850,
		Height:        620,
		DisableResize: false,
		Fullscreen:    false,
		RGBA:          0xFFFFFF00,
		Mac: &mac.Options{
			// TitleBar: mac.TitleBarHidden(),
			// TitleBar: mac.TitleBarHiddenInset(),
			TitleBar: mac.TitleBarDefault(),
			// Appearance:                    mac.NSAppearanceNameVibrantLight,
			WebviewIsTransparent:          true,
			WindowBackgroundIsTranslucent: true,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

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

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
