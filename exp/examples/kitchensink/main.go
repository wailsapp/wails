package main

import (
	"log"
	"sync"
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
		println("ApplicationDidFinishLaunching")
	})
	app.On(events.Mac.ApplicationWillTerminate, func() {
		println("ApplicationWillTerminate")
	})
	app.On(events.Mac.ApplicationDidBecomeActive, func() {
		println("ApplicationDidBecomeActive")
	})
	app.On(events.Mac.ApplicationDidChangeBackingProperties, func() {
		println("ApplicationDidChangeBackingProperties")
	})

	app.On(events.Mac.ApplicationDidChangeEffectiveAppearance, func() {
		println("ApplicationDidChangeEffectiveAppearance")
	})
	app.On(events.Mac.ApplicationDidHide, func() {
		println("ApplicationDidHide")
	})

	myWindow := app.NewWindow(&options.Window{
		Title:         "Kitchen Sink",
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
			Appearance: options.NSAppearanceNameDarkAqua,
		},
	})

	myWindow.On(events.Mac.WindowWillClose, func() {
		println(myWindow.ID(), "WindowWillClose")
	})
	myWindow.On(events.Mac.WindowDidResize, func() {
		w, h := myWindow.Size()
		println(myWindow.ID(), "WindowDidResize", w, h)
	})
	myWindow.On(events.Mac.WindowDidMove, func() {
		x, y := myWindow.Position()
		println(myWindow.ID(), "WindowDidMove", x, y)
	})
	myWindow.On(events.Mac.WindowDidMiniaturize, func() {
		println(myWindow.ID(), "WindowDidMiniaturize")
	})
	myWindow.On(events.Mac.WindowDidDeminiaturize, func() {
		println(myWindow.ID(), "WindowDidDeminiaturize")
	})
	myWindow.On(events.Mac.WindowDidBecomeKey, func() {
		println(myWindow.ID(), "WindowDidBecomeKey")
	})
	myWindow.On(events.Mac.WindowDidResignKey, func() {
		println(myWindow.ID(), "WindowDidResignKey")
	})
	myWindow.On(events.Mac.WindowDidBecomeMain, func() {
		println(myWindow.ID(), "WindowDidBecomeMain")
	})
	myWindow.On(events.Mac.WindowDidResignMain, func() {
		println(myWindow.ID(), "WindowDidResignMain")
	})
	myWindow.On(events.Mac.WindowWillEnterFullScreen, func() {
		println(myWindow.ID(), "WindowWillEnterFullScreen")
	})
	myWindow.On(events.Mac.WindowDidEnterFullScreen, func() {
		println(myWindow.ID(), "WindowDidEnterFullScreen")
	})
	myWindow.On(events.Mac.WindowWillExitFullScreen, func() {
		println(myWindow.ID(), "WindowWillExitFullScreen")
	})
	myWindow.On(events.Mac.WindowDidExitFullScreen, func() {
		println(myWindow.ID(), "WindowDidExitFullScreen")
	})
	myWindow.On(events.Mac.WindowWillEnterVersionBrowser, func() {
		println(myWindow.ID(), "WindowWillEnterVersionBrowser")
	})
	myWindow.On(events.Mac.WindowDidEnterVersionBrowser, func() {
		println(myWindow.ID(), "WindowDidEnterVersionBrowser")
	})
	myWindow.On(events.Mac.WindowWillExitVersionBrowser, func() {
		println(myWindow.ID(), "WindowWillExitVersionBrowser")
	})
	myWindow.On(events.Mac.WindowDidExitVersionBrowser, func() {
		println(myWindow.ID(), "WindowDidExitVersionBrowser")
	})

	var myWindow2 *application.Window
	var myWindow2Lock sync.RWMutex
	myWindow2 = app.NewWindow(&options.Window{
		Title:       "#2",
		Width:       1024,
		Height:      768,
		AlwaysOnTop: false,
		URL:         "https://google.com",
		Mac: &options.MacWindow{
			Backdrop: options.MacBackdropTranslucent,
		},
	})
	//myWindow2.On(events.Mac.WindowDidMove, func() {
	//	myWindow2Lock.RLock()
	//	x, y := myWindow2.Position()
	//	println(myWindow2.ID(), "WindowDidMove: ", x, y)
	//	myWindow2Lock.RUnlock()
	//})
	//

	go func() {
		time.Sleep(5 * time.Second)
		myWindow2Lock.RLock()
		myWindow.SetTitle("Wooooo")
		myWindow.SetAlwaysOnTop(true)
		myWindow2.EnableDevTools()
		myWindow2.SetTitle("OMFG")
		myWindow2.NavigateToURL("https://wails.io")
		myWindow.SetMinSize(600, 600)
		myWindow.SetMaxSize(650, 650)
		myWindow.Center()
		myWindow2Lock.RUnlock()

	}()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
