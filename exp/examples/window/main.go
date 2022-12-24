package main

import (
	_ "embed"
	"log"
	"strconv"

	"github.com/wailsapp/wails/exp/pkg/events"

	"github.com/wailsapp/wails/exp/pkg/application"
	"github.com/wailsapp/wails/exp/pkg/options"
)

func main() {
	app := application.New()
	app.SetName("Window Demo")
	app.SetDescription("A demo of the windowing capabilities")
	app.On(events.Mac.ApplicationDidFinishLaunching, func() {
		log.Println("ApplicationDidFinishLaunching")
	})

	// Create a custom menu
	menu := app.NewMenu()
	menu.AddRole(application.AppMenu)

	windowCounter := 1

	// Let's make a "Demo" menu
	myMenu := menu.AddSubmenu("New")
	myMenu.Add("New Blank Window").OnClick(func(ctx *application.Context) {
		app.NewWindow().SetTitle("Window " + strconv.Itoa(windowCounter)).Run()
		windowCounter++
	})
	myMenu.Add("New Window").OnClick(func(ctx *application.Context) {
		app.NewWindow().
			SetTitle("Window " + strconv.Itoa(windowCounter)).
			SetBackgroundColour(&options.RGBA{
				Red:   255,
				Green: 0,
				Blue:  0,
				Alpha: 0,
			}).
			Run()
		windowCounter++
	})
	myMenu.Add("New Webview").OnClick(func(ctx *application.Context) {
		app.NewWindow().
			SetTitle("Webview " + strconv.Itoa(windowCounter)).
			SetURL("https://wails.app").
			Run()
		windowCounter++
	})

	// Disabled menu item
	adjustMenu := menu.AddSubmenu("Adjust")
	adjustMenu.Add("Set Position (0,0)").OnClick(func(ctx *application.Context) {
		app.CurrentWindow().SetPosition(0, 0)
		windowCounter++
	})
	app.SetMenu(menu)
	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
