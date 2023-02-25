package main

import (
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {

	app := application.New(application.Options{
		Name:        "Menu Demo",
		Description: "A demo of the menu system",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a custom menu
	menu := app.NewMenu()
	menu.AddRole(application.AppMenu)

	// Let's make a "Demo" menu
	myMenu := menu.AddSubmenu("Demo")

	// Disabled menu item
	myMenu.Add("Not Enabled").SetEnabled(false)

	// Click callbacks
	myMenu.Add("Click Me!").OnClick(func(ctx *application.Context) {
		ctx.ClickedMenuItem().SetLabel("Thanks mate!")
	})

	// You can control the current window from the menu
	myMenu.Add("Lock WebviewWindow Resize").OnClick(func(ctx *application.Context) {
		if app.CurrentWindow().Resizable() {
			app.CurrentWindow().SetResizable(false)
			ctx.ClickedMenuItem().SetLabel("Unlock WebviewWindow Resize")
		} else {
			app.CurrentWindow().SetResizable(true)
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

	app.SetMenu(menu)

	app.NewWebviewWindow()

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
