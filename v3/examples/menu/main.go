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
		UseGlobalMenuByDefault: true,
		Name:                   "Menu Demo",
		Description:            "A demo of the menu system",
		Assets:                 application.AlphaAssets,
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

	disabledMenuItem := myMenu.Add("Show window").SetEnabled(false).OnClick(func(ctx *application.Context) {
		// Create the window with HTML content
		app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
			Title:  "Success",
			Width:  300,
			Height: 300,
			HTML:   "It worked",
		}).Show()
	})
	myMenu.Add("Enable show window").OnClick(func(ctx *application.Context) {
		disabledMenuItem.SetEnabled(true)
	})

	// Add an example of URL-based window creation
	myMenu.Add("Show window with URL").OnClick(func(ctx *application.Context) {
		// Create a window with a URL
		app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
			Title:  "URL Example",
			Width:  800,
			Height: 600,
			URL:    "https://wails.io",
		}).Show()
	})

	myMenu.Add("Show window with raw HTML").OnClick(func(ctx *application.Context) {
		// Create a window with a URL
		app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
			Title:  "Raw HTML Example",
			Width:  800,
			Height: 600,
			HTML:   "<html><body><h1>Hello World</h1></body></html>",
		}).Show()
	})

	// Not Enabled menu item (left for backwards compatibility)
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

	// ---- New index-based menu operations demo ----
	indexMenu := menu.AddSubmenu("Index Operations")

	// Add some initial items
	indexMenu.Add("Item 0")
	indexMenu.Add("Item 2")
	indexMenu.Add("Item 4")

	// Demonstrate inserting items at specific indices
	indexMenu.InsertAt(1, "Item 1").OnClick(func(*application.Context) {
		println("Item 1 clicked")
	})

	indexMenu.InsertAt(3, "Item 3").OnClick(func(*application.Context) {
		println("Item 3 clicked")
	})

	// Demonstrate inserting different types of items at specific indices
	indexMenu.AddSeparator()
	indexMenu.InsertCheckboxAt(6, "Checkbox at index 6", true).OnClick(func(ctx *application.Context) {
		println("Checkbox at index 6 clicked, checked:", ctx.ClickedMenuItem().Checked())
	})

	indexMenu.InsertRadioAt(7, "Radio at index 7", true).OnClick(func(ctx *application.Context) {
		println("Radio at index 7 clicked")
	})

	indexMenu.InsertSeparatorAt(8)

	// Create a submenu and insert it at a specific index
	submenuAtIndex := indexMenu.InsertSubmenuAt(9, "Inserted Submenu")
	submenuAtIndex.Add("Submenu Item 1")
	submenuAtIndex.Add("Submenu Item 2")

	// Demonstrate ItemAt to access items by index
	indexMenu.AddSeparator()
	indexMenu.Add("Get Item at Index").OnClick(func(*application.Context) {
		// Get the item at index 2 and change its label
		if item := indexMenu.ItemAt(2); item != nil {
			println("Item at index 2:", item.Label())
			item.SetLabel("Item 2 (Modified)")
		}
	})

	// Demonstrate Count method
	indexMenu.Add("Count Items").OnClick(func(*application.Context) {
		println("Menu has", indexMenu.Count(), "items")
	})

	// Demonstrate visibility control for different item types
	visibilityMenu := menu.AddSubmenu("Visibility Control")

	// Regular menu item
	regularItem := visibilityMenu.Add("Regular Item")

	// Checkbox menu item
	checkboxItem := visibilityMenu.AddCheckbox("Checkbox Item", true)

	// Radio menu item
	radioItem := visibilityMenu.AddRadio("Radio Item", true)

	// Separator
	visibilityMenu.AddSeparator()
	separatorIndex := visibilityMenu.Count() - 1
	separatorItem := visibilityMenu.ItemAt(separatorIndex)

	// Submenu - get the MenuItem for the submenu to control visibility
	submenuMenuItem := application.NewSubMenuItem("Submenu")
	visibilityMenu.InsertItemAt(visibilityMenu.Count(), submenuMenuItem)
	submenuContent := submenuMenuItem.GetSubmenu()
	submenuContent.Add("Submenu Content")

	// Controls for toggling visibility
	visibilityMenu.AddSeparator()
	visibilityMenu.Add("Toggle Regular Item").OnClick(func(*application.Context) {
		regularItem.SetHidden(!regularItem.Hidden())
		println("Regular item hidden:", regularItem.Hidden())
	})

	visibilityMenu.Add("Toggle Checkbox Item").OnClick(func(*application.Context) {
		checkboxItem.SetHidden(!checkboxItem.Hidden())
		println("Checkbox item hidden:", checkboxItem.Hidden())
	})

	visibilityMenu.Add("Toggle Radio Item").OnClick(func(*application.Context) {
		radioItem.SetHidden(!radioItem.Hidden())
		println("Radio item hidden:", radioItem.Hidden())
	})

	visibilityMenu.Add("Toggle Separator").OnClick(func(*application.Context) {
		separatorItem.SetHidden(!separatorItem.Hidden())
		println("Separator hidden:", separatorItem.Hidden())
	})

	// For submenu visibility, we need to toggle the visibility of the MenuItem that contains the submenu
	visibilityMenu.Add("Toggle Submenu").OnClick(func(ctx *application.Context) {
		// Log the current state before toggling
		println("Submenu hidden before toggle:", submenuMenuItem.Hidden())

		// Toggle the visibility
		submenuMenuItem.SetHidden(!submenuMenuItem.Hidden())

		// Log the new state after toggling
		println("Submenu hidden after toggle:", submenuMenuItem.Hidden())

		// Update the menu item label to reflect the current state
		if submenuMenuItem.Hidden() {
			ctx.ClickedMenuItem().SetLabel("Show Submenu")
		} else {
			ctx.ClickedMenuItem().SetLabel("Hide Submenu")
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


	// ---- New index-based menu operations demo ----
	indexMenu := menu.AddSubmenu("Index Operations")

	// Add some initial items
	indexMenu.Add("Item 0")
	indexMenu.Add("Item 2")
	indexMenu.Add("Item 4")

	// Demonstrate inserting items at specific indices
	indexMenu.InsertAt(1, "Item 1").OnClick(func(*application.Context) {
		println("Item 1 clicked")
	})

	indexMenu.InsertAt(3, "Item 3").OnClick(func(*application.Context) {
		println("Item 3 clicked")
	})

	// Demonstrate inserting different types of items at specific indices
	indexMenu.AddSeparator()
	indexMenu.InsertCheckboxAt(6, "Checkbox at index 6", true).OnClick(func(ctx *application.Context) {
		println("Checkbox at index 6 clicked, checked:", ctx.ClickedMenuItem().Checked())
	})

	indexMenu.InsertRadioAt(7, "Radio at index 7", true).OnClick(func(ctx *application.Context) {
		println("Radio at index 7 clicked")
	})

	indexMenu.InsertSeparatorAt(8)

	// Create a submenu and insert it at a specific index
	submenuAtIndex := indexMenu.InsertSubmenuAt(9, "Inserted Submenu")
	submenuAtIndex.Add("Submenu Item 1")
	submenuAtIndex.Add("Submenu Item 2")

	// Demonstrate ItemAt to access items by index
	indexMenu.AddSeparator()
	indexMenu.Add("Get Item at Index").OnClick(func(*application.Context) {
		// Get the item at index 2 and change its label
		if item := indexMenu.ItemAt(2); item != nil {
			println("Item at index 2:", item.Label())
			item.SetLabel("Item 2 (Modified)")
		}
	})

	// Demonstrate Count method
	indexMenu.Add("Count Items").OnClick(func(*application.Context) {
		println("Menu has", indexMenu.Count(), "items")
	})

	// Demonstrate visibility control for different item types
	visibilityMenu := menu.AddSubmenu("Visibility Control")

	// Regular menu item
	regularItem := visibilityMenu.Add("Regular Item")

	// Checkbox menu item
	checkboxItem := visibilityMenu.AddCheckbox("Checkbox Item", true)

	// Radio menu item
	radioItem := visibilityMenu.AddRadio("Radio Item", true)

	// Separator
	visibilityMenu.AddSeparator()
	separatorIndex := visibilityMenu.Count() - 1
	separatorItem := visibilityMenu.ItemAt(separatorIndex)

	// Submenu - get the MenuItem for the submenu to control visibility
	submenuMenuItem := application.NewSubMenuItem("Submenu")
	visibilityMenu.InsertItemAt(visibilityMenu.Count(), submenuMenuItem)
	submenuContent := submenuMenuItem.GetSubmenu()
	submenuContent.Add("Submenu Content")

	// Controls for toggling visibility
	visibilityMenu.AddSeparator()
	visibilityMenu.Add("Toggle Regular Item").OnClick(func(*application.Context) {
		regularItem.SetHidden(!regularItem.Hidden())
		println("Regular item hidden:", regularItem.Hidden())
	})

	visibilityMenu.Add("Toggle Checkbox Item").OnClick(func(*application.Context) {
		checkboxItem.SetHidden(!checkboxItem.Hidden())
		println("Checkbox item hidden:", checkboxItem.Hidden())
	})

	visibilityMenu.Add("Toggle Radio Item").OnClick(func(*application.Context) {
		radioItem.SetHidden(!radioItem.Hidden())
		println("Radio item hidden:", radioItem.Hidden())
	})

	visibilityMenu.Add("Toggle Separator").OnClick(func(*application.Context) {
		separatorItem.SetHidden(!separatorItem.Hidden())
		println("Separator hidden:", separatorItem.Hidden())
	})

	// For submenu visibility, we need to toggle the visibility of the MenuItem that contains the submenu
	visibilityMenu.Add("Toggle Submenu").OnClick(func(ctx *application.Context) {
		// Log the current state before toggling
		println("Submenu hidden before toggle:", submenuMenuItem.Hidden())

		// Toggle the visibility
		submenuMenuItem.SetHidden(!submenuMenuItem.Hidden())

		// Log the new state after toggling
		println("Submenu hidden after toggle:", submenuMenuItem.Hidden())

		// Update the menu item label to reflect the current state
		if submenuMenuItem.Hidden() {
			ctx.ClickedMenuItem().SetLabel("Show Submenu")
		} else {
			ctx.ClickedMenuItem().SetLabel("Hide Submenu")
		}
	})

	app.SetMenu(menu)

	window := app.NewWebviewWindow().SetBackgroundColour(application.NewRGB(33, 37, 41))
	window.SetMenu(menu)

	err := app.Run()

	if err != nil {
		log.Fatal(err.Error())
	}
}
