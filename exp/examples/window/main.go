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

	sizeMenu := menu.AddSubmenu("Size")
	sizeMenu.Add("Set Size (800,600)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.SetSize(800, 600)
		})
	})

	sizeMenu.Add("Set Size (Random)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.SetSize(rand.Intn(800)+200, rand.Intn(600)+200)
		})
	})
	sizeMenu.Add("Set Min Size (200,200)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.SetMinSize(200, 200)
		})
	})
	sizeMenu.Add("Set Max Size (600,600)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.SetFullscreenButtonEnabled(false)
			w.SetMaxSize(600, 600)
		})
	})
	sizeMenu.Add("Reset Min Size").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.SetMinSize(0, 0)
		})
	})

	sizeMenu.Add("Reset Max Size").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.SetMaxSize(0, 0)
			w.SetFullscreenButtonEnabled(true)
		})
	})
	positionMenu := menu.AddSubmenu("Position")
	positionMenu.Add("Set Position (0,0)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.SetPosition(0, 0)
		})
	})
	positionMenu.Add("Set Position (Random)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.SetPosition(rand.Intn(1000), rand.Intn(800))
		})
	})
	positionMenu.Add("Center").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.Center()
		})
	})
	stateMenu := menu.AddSubmenu("State")
	stateMenu.Add("Minimise (for 2 secs)").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.Minimise()
			time.Sleep(2 * time.Second)
			w.Restore()
		})
	})
	stateMenu.Add("Maximise").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.Maximise()
		})
	})
	stateMenu.Add("Fullscreen").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.Fullscreen()
		})
	})
	stateMenu.Add("UnFullscreen").OnClick(func(ctx *application.Context) {
		currentWindow(func(w *application.Window) {
			w.UnFullscreen()
		})
	})
	stateMenu.Add("Restore").OnClick(func(ctx *application.Context) {
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
