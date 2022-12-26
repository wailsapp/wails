package main

import (
	_ "embed"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/wailsapp/wails/exp/pkg/events"

	"github.com/wailsapp/wails/exp/pkg/application"
)

func main() {
	app := application.New()
	app.SetName("Window Demo")
	app.SetDescription("A demo of the windowing capabilities")
	app.On(events.Mac.ApplicationDidFinishLaunching, func() {
		log.Println("ApplicationDidFinishLaunching")
	})

	currentWindow := func(fn func(window *application.Window)) {
		if app.CurrentWindow() != nil {
			fn(app.CurrentWindow())
		} else {
			println("Current Window is nil")
		}
	}

	// Create a custom menu
	menu := app.NewMenu()
	menu.AddRole(application.AppMenu)

	windowCounter := 1

	// Let's make a "Demo" menu
	myMenu := menu.AddSubmenu("New")

	myMenu.Add("New Window").
		SetAccelerator("CmdOrCtrl+N").
		OnClick(func(ctx *application.Context) {
			app.NewWindow().
				SetTitle("Window "+strconv.Itoa(windowCounter)).
				SetPosition(rand.Intn(1000), rand.Intn(800)).
				SetURL("https://wails.io").
				Run()
			windowCounter++
		})

	// Disabled menu item
	adjustMenu := menu.AddSubmenu("Adjust")
	adjustMenu.Add("Set Position (0,0)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.SetPosition(0, 0)
		})
	})
	adjustMenu.Add("Set Position (Random)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.SetPosition(rand.Intn(1000), rand.Intn(800))
		})
	})
	adjustMenu.Add("Set Size (800,600)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.SetSize(800, 600)
		})
	})

	adjustMenu.Add("Set Size (Random)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.SetSize(rand.Intn(800)+200, rand.Intn(600)+200)
		})
	})
	adjustMenu.Add("Set Min Size (200,200)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.SetMinSize(200, 200)
		})
	})
	adjustMenu.Add("Set Max Size (600,600)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.SetMaxSize(600, 600)
		})
	})
	adjustMenu.Add("Reset Min Size").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.SetMinSize(0, 0)
		})
	})

	adjustMenu.Add("Reset Max Size").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.SetMaxSize(0, 0)
		})
	})
	adjustMenu.Add("Center").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.Center()
		})
	})
	adjustMenu.Add("Minimise (for 2 secs)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.Minimise()
			time.Sleep(2 * time.Second)
			w.Restore()
		})
	})
	adjustMenu.Add("Maximise").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.Maximise()
		})
	})
	adjustMenu.Add("Fullscreen").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.Fullscreen()
		})
	})
	adjustMenu.Add("Restore").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.Restore()
		})
	})

	app.NewWindow()

	app.SetMenu(menu)
	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}

}
