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
	fileMenu := menu.AddRole(application.FileMenu)
	_ = fileMenu
	//fileMenu.FindByRole(application.Open).OnClick(func(context *application.Context) {
	//	selection, err := application.OpenFileDialog().PromptForSingleSelection()
	//	if err != nil {
	//		println("Error: " + err.Error())
	//		return
	//	}
	//	println("You selected: " + selection)
	//})
	menu.AddRole(application.EditMenu)
	menu.AddRole(application.WindowMenu)
	menu.AddRole(application.HelpMenu)

	// Let's make a "Demo" menu
	myMenu := menu.AddSubmenu("Demo")

	// Disabled menu item
	myMenu.Add("Not Enabled").SetEnabled(false)

	// Click callbacks
	myMenu.Add("Click Me!").SetAccelerator("CmdOrCtrl+l").SetBitmap(clickBitmap).OnClick(func(ctx *application.Context) {
		switch ctx.ClickedMenuItem().Label() {
		case "Click Me!":
			ctx.ClickedMenuItem().SetLabel("Thanks mate!")
		case "Thanks mate!":
			ctx.ClickedMenuItem().SetLabel("Click Me!")
		}
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
	myMenu.Add("Hide the beatles").OnClick(func(ctx *application.Context) {
		if beatles.Hidden() {
			ctx.ClickedMenuItem().SetLabel("Hide the beatles!")
			beatles.SetHidden(false)
		} else {
			beatles.SetHidden(true)
			ctx.ClickedMenuItem().SetLabel("Unhide the beatles!")
		}
	})
	app.SetMenu(menu)

	app.NewWebviewWindow().SetBackgroundColour(application.NewRGB(33, 37, 41))

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
