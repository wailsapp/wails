package main

import (
	_ "embed"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/options"
)

func main() {
	app := application.New(options.Application{
		Name:        "Menu Demo",
		Description: "A demo of the menu system",
		Mac: options.Mac{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})
	/*
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

	*/

	menuCallback := func(ctx *application.Context) {
		menuItem := ctx.ClickedMenuItem()
		menuItem.SetLabel("Clicked!")
	}

	radioCallback := func(ctx *application.Context) {
		menuItem := ctx.ClickedMenuItem()
		menuItem.SetLabel(menuItem.Label() + "!")
	}

	myMenu := app.NewMenu()
	file1 := myMenu.Add("File")
	file1.SetTooltip("Create New Tray Menu")
	file1.OnClick(menuCallback)
	myMenu.Add("Create New Tray Menu").
		SetAccelerator("CmdOrCtrl+N").
		SetTooltip("ROFLCOPTER!!!!").
		OnClick(func(ctx *application.Context) {
			mySystray := app.NewSystemTray()
			mySystray.SetLabel("Wails")
			if runtime.GOOS == "darwin" {
				mySystray.SetTemplateIcon(application.DefaultMacTemplateIcon)
			} else {
				mySystray.SetIcon(application.DefaultApplicationIcon)
			}
			myMenu := app.NewMenu()
			myMenu.Add("Item 1")
			myMenu.AddSeparator()
			myMenu.Add("Kill this menu").OnClick(func(ctx *application.Context) {
				mySystray.Destroy()
			})
			mySystray.SetMenu(myMenu)

		})
	myMenu.Add("Not Enabled").SetEnabled(false)
	myMenu.AddSeparator()
	myMenu.AddCheckbox("My checkbox", true).OnClick(menuCallback)
	myMenu.AddSeparator()
	myMenu.AddRadio("Radio 1", true).OnClick(radioCallback)
	myMenu.AddRadio("Radio 2", false).OnClick(radioCallback)
	myMenu.AddRadio("Radio 3", false).OnClick(radioCallback)

	submenu := myMenu.AddSubmenu("Submenu")
	submenu.Add("Submenu item 1").OnClick(menuCallback)
	submenu.Add("Submenu item 2").OnClick(menuCallback)
	submenu.Add("Submenu item 3").OnClick(menuCallback)
	myMenu.AddSeparator()
	file4 := myMenu.Add("File 4").OnClick(func(*application.Context) {
		println("File 4 clicked")
	})

	myMenu.Add("Click to toggle").OnClick(func(*application.Context) {
		enabled := file4.Enabled()
		println("Enabled: ", enabled)
		file4.SetEnabled(!enabled)
	})
	myMenu.Add("File 5").OnClick(menuCallback)

	mySystray := app.NewSystemTray()
	mySystray.SetLabel("Wails is awesome")
	if runtime.GOOS == "darwin" {
		mySystray.SetTemplateIcon(application.DefaultMacTemplateIcon)
	} else {
		mySystray.SetIcon(application.DefaultApplicationIcon)
	}
	mySystray.SetMenu(myMenu)
	mySystray.SetIconPosition(application.NSImageLeading)

	myWindow := app.NewWebviewWindowWithOptions(&options.WebviewWindow{
		Title:         "Kitchen Sink",
		Width:         600,
		Height:        400,
		AlwaysOnTop:   true,
		DisableResize: false,
		BackgroundColour: &options.RGBA{
			Red:   255,
			Green: 255,
			Blue:  255,
			Alpha: 30,
		},
		StartState: options.WindowStateMaximised,
		Mac: options.MacWindow{
			Backdrop:   options.MacBackdropTranslucent,
			Appearance: options.NSAppearanceNameDarkAqua,
		},
	})
	/*
		myWindow.On(events.Mac.WindowWillClose, func() {
			println(myWindow.ID(), "WindowWillClose")
		})
		myWindow.On(events.Mac.WindowDidResize, func() {
			//w, h := myWindow.Size()
			//println(myWindow.ID(), "WindowDidResize", w, h)
		})
		myWindow.On(events.Mac.WindowDidMove, func() {
			//x, y := myWindow.Position()
			//println(myWindow.ID(), "WindowDidMove", x, y)
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
	*/
	var myWindow2 *application.WebviewWindow
	var myWindow2Lock sync.RWMutex
	myWindow2 = app.NewWebviewWindowWithOptions(&options.WebviewWindow{
		Title:       "#2",
		Width:       1024,
		Height:      768,
		AlwaysOnTop: false,
		URL:         "https://google.com",
		Mac: options.MacWindow{
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
		myWindow2.SetTitle("OMG")
		myWindow2.SetURL("https://wails.io")
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
