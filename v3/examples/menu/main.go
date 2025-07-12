package main

import (
	_ "embed"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed icon.png
var clickBitmap []byte

func main() {

	app := application.New(application.Options{
		Name:        "Menu Demo",
		Description: "A demo of the menu system",
		Assets:      application.AlphaAssets,
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a custom menu
	menu := app.NewMenu()
	if runtime.GOOS == "darwin" {
		menu.AddRole(application.AppMenu)
	}
	menu.AddRole(application.FileMenu)
	menu.AddRole(application.EditMenu)
	menu.AddRole(application.WindowMenu)
	menu.AddRole(application.HelpMenu)

	// Let's make a "Demo" menu
	myMenu := menu.AddSubmenu("Demo")

	// Hidden menu item that can be unhidden
	hidden := myMenu.Add("I was hidden").SetHidden(true)
	myMenu.Add("Toggle the hidden menu").OnClick(func(ctx *application.Context) {
		hidden.SetHidden(!hidden.Hidden())
	})

	// Disabled menu item
	myMenu.Add("Not Enabled").SetEnabled(false)

	// Click callbacks
	myMenu.Add("Click Me!").SetAccelerator("CmdOrCtrl+l").OnClick(func(ctx *application.Context) {
		switch ctx.ClickedMenuItem().Label() {
		case "Click Me!":
			ctx.ClickedMenuItem().SetLabel("Thanks mate!")
		case "Thanks mate!":
			ctx.ClickedMenuItem().SetLabel("Click Me!")
		}
	})

	// You can control the current window from the menu
	myMenu.Add("Lock WebviewWindow Resize").OnClick(func(ctx *application.Context) {
		if app.Window.Current().Resizable() {
			app.Window.Current().SetResizable(false)
			ctx.ClickedMenuItem().SetLabel("Unlock WebviewWindow Resize")
		} else {
			app.Window.Current().SetResizable(true)
			ctx.ClickedMenuItem().SetLabel("Lock WebviewWindow Resize")
		}
	})

	myMenu.AddSeparator()

	// Checkboxes will tell you their new state so you don't need to track it
	myMenu.AddCheckbox("My checkbox", true).OnClick(func(context *application.Context) {
		println("Clicked checkbox. Checked:", context.ClickedMenuItem().Checked())
	})
	myMenu.AddSeparator()

	// Callbacks can be shared. This is useful for radio groups
	radioCallback := func(ctx *application.Context) {
		menuItem := ctx.ClickedMenuItem()
		menuItem.SetLabel(menuItem.Label() + "!")
	}

	// Radio groups are created implicitly by placing radio items next to each other in a menu
	myMenu.AddRadio("Radio 1", true).OnClick(radioCallback)
	myMenu.AddRadio("Radio 2", false).OnClick(radioCallback)
	myMenu.AddRadio("Radio 3", false).OnClick(radioCallback)

	// Submenus are also supported
	submenu := myMenu.AddSubmenu("Submenu")
	submenu.Add("Submenu item 1")
	submenu.Add("Submenu item 2")
	submenu.Add("Submenu item 3")

	myMenu.AddSeparator()

	beatles := myMenu.Add("Hello").OnClick(func(*application.Context) {
		println("The beatles would be proud")
	})
	myMenu.Add("Toggle the menuitem above").OnClick(func(*application.Context) {
		if beatles.Enabled() {
			beatles.SetEnabled(false)
			beatles.SetLabel("Goodbye")
		} else {
			beatles.SetEnabled(true)
			beatles.SetLabel("Hello")
		}
	})
	myMenu.Add("Hide the beatles").OnClick(func(ctx *application.Context) {
		if beatles.Hidden() {
			ctx.ClickedMenuItem().SetLabel("Hide the beatles!")
			beatles.SetHidden(false)
		} else {
			beatles.SetHidden(true)
			ctx.ClickedMenuItem().SetLabel("Unhide the beatles!")
		}
	})

	myMenu.AddSeparator()

	coffee := myMenu.Add("Request Coffee").OnClick(func(*application.Context) {
		println("Coffee dispatched. Productivity +10!")
	})

	myMenu.Add("Toggle coffee availability").OnClick(func(*application.Context) {
		if coffee.Enabled() {
			coffee.SetEnabled(false)
			coffee.SetLabel("Coffee Machine Broken")
			println("Alert: Developer morale critically low.")
		} else {
			coffee.SetEnabled(true)
			coffee.SetLabel("Request Coffee")
			println("All systems nominal. Coffee restored.")
		}
	})

	myMenu.Add("Hide the coffee option").OnClick(func(ctx *application.Context) {
		if coffee.Hidden() {
			ctx.ClickedMenuItem().SetLabel("Hide the coffee option")
			coffee.SetHidden(false)
			println("Coffee menu item has been resurrected!")
		} else {
			coffee.SetHidden(true)
			ctx.ClickedMenuItem().SetLabel("Unhide the coffee option")
			println("The coffee option has vanished into the void.")
		}
	})

	app.Menu.Set(menu)

	window := app.Window.New().SetBackgroundColour(application.NewRGB(33, 37, 41))
	window.SetMenu(menu)

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
