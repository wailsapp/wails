package main

import (
	"log"
	"time"

	"github.com/wailsapp/wails/exp/pkg/application"
	"github.com/wailsapp/wails/exp/pkg/events"
	"github.com/wailsapp/wails/exp/pkg/options"
)

func main() {
	app := application.New(&options.Application{
		Mac: &options.Mac{
			//ActivationPolicy: options.ActivationPolicyAccessory,
		},
	})
	app.On(events.Mac.ApplicationDidFinishLaunching, func() {
		println("WOOOOO!")
	})
	app.On(events.Mac.ApplicationWillTerminate, func() {
		println("TERMINATION!!")
	})
	myWindow := app.NewWindow(&options.Window{
		Title:         "Basic",
		Width:         600,
		Height:        400,
		AlwaysOnTop:   true,
		DisableResize: false,
		//MinWidth:       100,
		//MinHeight:      100,
		//MaxWidth:       1000,
		//MaxHeight:      1000,
		EnableDevTools: true,
		BackgroundColour: &options.RGBA{
			Red:   255,
			Green: 255,
			Blue:  255,
			Alpha: 30,
		},
		StartState: options.WindowStateMaximised,
		Mac: &options.MacWindow{
			Backdrop:   options.MacBackdropTranslucent,
			TitleBar:   options.TitleBarHiddenInset(),
			Appearance: options.NSAppearanceNameDarkAqua,
		},
	})

	myWindow2 := app.NewWindow(&options.Window{
		Title:         "#2",
		Width:         1024,
		Height:        768,
		AlwaysOnTop:   false,
		URL:           "https://google.com",
		DisableResize: true,
		Mac: &options.MacWindow{
			Backdrop: options.MacBackdropTranslucent,
			TitleBar: options.TitleBarHidden(),
		},
	})

	go func() {
		time.Sleep(5 * time.Second)
		myWindow.SetTitle("Wooooo")
		myWindow.SetAlwaysOnTop(true)
		myWindow2.EnableDevTools()
		myWindow2.SetTitle("OMFG")
		myWindow2.NavigateToURL("https://wails.io")
		myWindow.SetMinSize(600, 600)
		myWindow.SetMaxSize(650, 650)

	}()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
