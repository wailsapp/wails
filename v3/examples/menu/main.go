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

	// Create a top-level menu for the submenu label update functionality
	newMenu := menu.AddSubmenu("New")
	newMenu.Add("New Window (shared menu)").OnClick(func(*application.Context) {
		window := app.NewWebviewWindow().SetBackgroundColour(application.NewRGB(33, 37, 41))
		window.SetMenu(menu)
	})
	newMenu.Add("New Window (cloned menu)").OnClick(func(*application.Context) {
		window := app.NewWebviewWindow().SetBackgroundColour(application.NewRGB(33, 37, 41))
		window.SetMenu(menu.Clone())
	})

	// Create a menu for demonstrating the index-based API
	indexMenu := menu.AddSubmenu("Index API")

	// Add a hidden item at the start of the menu (to demonstrate the bug in issue #4088)
	hiddenFirstItem := indexMenu.Add("Hidden First Item").SetHidden(true)
	indexMenu.Add("Toggle Hidden First Item").OnClick(func(*application.Context) {
		hiddenFirstItem.SetHidden(!hiddenFirstItem.Hidden())
		menu.Update()
		println("Toggled visibility of first item - Now hidden:", hiddenFirstItem.Hidden())
	})
	indexMenu.AddSeparator()

	// Add some initial items to the menu
	indexMenu.Add("Item 1")
	indexMenu.Add("Item 2")
	indexMenu.Add("Item 3")
	indexMenu.AddSeparator()

	// Demonstrate InsertAt - insert a new item at index 1 (between Item 1 and Item 2)
	indexMenu.Add("Insert at index 1").OnClick(func(*application.Context) {
		item := indexMenu.InsertAt(1, "Inserted Item")
		item.OnClick(func(*application.Context) {
			println("Clicked on inserted item!")
		})
		menu.Update()
		println("Inserted item at index 1")
	})

	// Demonstrate InsertCheckboxAt - insert a checkbox at index 2
	indexMenu.Add("Insert checkbox at index 2").OnClick(func(*application.Context) {
		checkbox := indexMenu.InsertCheckboxAt(2, "Inserted Checkbox", true)
		checkbox.OnClick(func(ctx *application.Context) {
			println("Checkbox clicked, checked:", ctx.ClickedMenuItem().Checked())
		})
		menu.Update()
		println("Inserted checkbox at index 2")
	})

	indexMenu.AddSeparator()

	// Create a section for checkboxes
	indexMenu.Add("Checkbox Items Demo").SetEnabled(false)

	// Add initial checkboxes
	checkbox1 := indexMenu.AddCheckbox("Checkbox A", true).OnClick(func(ctx *application.Context) {
		println("Checkbox A clicked, checked:", ctx.ClickedMenuItem().Checked())
	})

	// Demonstrate adding a checkbox next to an existing checkbox
	indexMenu.Add("Insert checkbox next to Checkbox A").OnClick(func(*application.Context) {
		// Find the index of Checkbox A
		var checkboxIndex int
		for i := 0; ; i++ {
			item := indexMenu.ItemAt(i)
			if item == nil {
				break
			}
			if item == checkbox1 {
				checkboxIndex = i
				break
			}
		}

		// Insert a new checkbox right after Checkbox A
		newCheckbox := indexMenu.InsertCheckboxAt(checkboxIndex+1, "Checkbox B", false)
		newCheckbox.OnClick(func(ctx *application.Context) {
			println("Checkbox B clicked, checked:", ctx.ClickedMenuItem().Checked())
		})
		menu.Update()
		println("Inserted Checkbox B after Checkbox A")
	})

	indexMenu.AddSeparator()

	// Create a section for radio items
	indexMenu.Add("Radio Items Demo").SetEnabled(false)

	// Add initial radio items
	indexMenu.AddRadio("Radio A", true).OnClick(func(ctx *application.Context) {
		println("Radio A clicked, checked:", ctx.ClickedMenuItem().Checked())
	})
	radio2 := indexMenu.AddRadio("Radio B", false).OnClick(func(ctx *application.Context) {
		println("Radio B clicked, checked:", ctx.ClickedMenuItem().Checked())
	})

	// Demonstrate adding a radio item next to existing radio items
	indexMenu.Add("Insert radio next to Radio B").OnClick(func(*application.Context) {
		// Find the index of Radio B
		var radioIndex int
		for i := 0; ; i++ {
			item := indexMenu.ItemAt(i)
			if item == nil {
				break
			}
			if item == radio2 {
				radioIndex = i
				break
			}
		}

		// Insert a new radio item right after Radio B
		newRadio := indexMenu.InsertRadioAt(radioIndex+1, "Radio C", false)
		newRadio.OnClick(func(ctx *application.Context) {
			println("Radio C clicked, checked:", ctx.ClickedMenuItem().Checked())
		})
		menu.Update()
		println("Inserted Radio C after Radio B")
	})

	indexMenu.AddSeparator()

	// Demonstrate InsertSeparatorAt
	indexMenu.Add("Insert separator at index 3").OnClick(func(*application.Context) {
		indexMenu.InsertSeparatorAt(3)
		menu.Update()
		println("Inserted separator at index 3")
	})

	// Demonstrate InsertSubmenuAt
	indexMenu.Add("Insert submenu at end").OnClick(func(*application.Context) {
		// Get the count of items to determine the end index
		var count int
		for i := 0; ; i++ {
			if indexMenu.ItemAt(i) == nil {
				count = i
				break
			}
		}

		// Insert a submenu at the end
		submenu := indexMenu.InsertSubmenuAt(count, "Inserted Submenu")
		submenu.Add("Submenu Item 1")
		submenu.Add("Submenu Item 2")
		menu.Update()
		println("Inserted submenu at end")
	})

	indexMenu.AddSeparator()

	// Create a section for enabling/disabling menu items
	indexMenu.Add("Enable/Disable Items Demo").SetEnabled(false)

	// Add an item that we'll enable/disable
	targetItem := indexMenu.Add("Target Item")
	targetItem.OnClick(func(*application.Context) {
		println("Target item clicked!")
	})

	// Demonstrate enabling/disabling a menu item at a specific index
	indexMenu.Add("Toggle enable/disable for Target Item").OnClick(func(*application.Context) {
		// Find the index of the target item
		var targetIndex int
		for i := 0; ; i++ {
			item := indexMenu.ItemAt(i)
			if item == nil {
				break
			}
			if item == targetItem {
				targetIndex = i
				break
			}
		}

		// Toggle the enabled state of the item at the found index
		item := indexMenu.ItemAt(targetIndex)
		if item != nil {
			item.SetEnabled(!item.Enabled())
			menu.Update()
			println("Toggled enabled state for item at index", targetIndex, "- Now enabled:", item.Enabled())
		}
	})

	indexMenu.AddSeparator()

	// Create a section for hiding/unhiding menu items
	indexMenu.Add("Hide/Unhide Items Demo").SetEnabled(false)

	// Add an item that we'll hide/unhide
	hidableItem := indexMenu.Add("Hidable Item")
	hidableItem.OnClick(func(*application.Context) {
		println("Hidable item clicked!")
	})

	// Demonstrate hiding/unhiding a menu item at a specific index
	indexMenu.Add("Toggle hide/unhide for Hidable Item").OnClick(func(*application.Context) {
		// Find the index of the hidable item
		var hidableIndex int
		for i := 0; ; i++ {
			item := indexMenu.ItemAt(i)
			if item == nil {
				break
			}
			if item == hidableItem {
				hidableIndex = i
				break
			}
		}

		// Toggle the hidden state of the item at the found index
		item := indexMenu.ItemAt(hidableIndex)
		if item != nil {
			item.SetHidden(!item.Hidden())
			menu.Update()
			println("Toggled hidden state for item at index", hidableIndex, "- Now hidden:", item.Hidden())
		}
	})

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

	// Create a hidable submenu for demonstration
	hidableSubmenu := myMenu.AddSubmenu("Hidable Submenu")
	hidableSubmenu.Add("Hidable Submenu item 1")
	hidableSubmenu.Add("Hidable Submenu item 2")
	hidableSubmenu.Add("Hidable Submenu item 3")

	// Get a reference to the menu item that contains the hidable submenu
	hidableSubmenuItem := myMenu.FindByLabel("Hidable Submenu")

	// Add a button to toggle the visibility of the hidable submenu
	myMenu.Add("Toggle Hidable Submenu").OnClick(func(ctx *application.Context) {
		if hidableSubmenuItem.Hidden() {
			hidableSubmenuItem.SetHidden(false)
			ctx.ClickedMenuItem().SetLabel("Hide Submenu")
			println("Submenu is now visible")
		} else {
			hidableSubmenuItem.SetHidden(true)
			ctx.ClickedMenuItem().SetLabel("Show Submenu")
			println("Submenu is now hidden")
		}
		// Update the menu to reflect the changes
		menu.Update()
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

	// Demonstrate the new menu manipulation API
	myMenu.AddSeparator()
	myMenu.Add("Menu Manipulation API Demo").SetEnabled(false)

	// Create a new menu item to append
	newItem := application.NewMenuItem("Appended Item")
	newItem.OnClick(func(*application.Context) {
		println("Clicked on appended item!")
	})

	// Demonstrate AppendItem with Update(true)
	myMenu.Add("Append a new item").OnClick(func(*application.Context) {
		myMenu.AppendItem(newItem)
		// Use Update(true) to update everything in one call
		menu.Update()
		println("Appended a new item to the menu")
	})

	// Demonstrate Remove with Update(true)
	myMenu.Add("Remove the last item").OnClick(func(*application.Context) {
		// Get the index of the last item
		var lastIndex int
		for i := 0; ; i++ {
			if myMenu.ItemAt(i) == nil {
				lastIndex = i - 1
				break
			}
		}
		if lastIndex >= 0 {
			myMenu.Remove(lastIndex)
			// Use Update(true) to update everything in one call
			menu.Update()
			println("Removed the last item from the menu")
		} else {
			println("No items to remove")
		}
	})

	// Create a new top-level menu for demonstrating disabled items with callbacks
	disabledMenu := menu.AddSubmenu("Disabled Callback Demo")

	// Add a disabled menu item with an onClick callback
	disabledItem := disabledMenu.Add("Disabled Item with Callback").SetEnabled(false)
	disabledItem.OnClick(func(*application.Context) {
		println("Disabled item was clicked! This callback is still valid.")
	})

	// Add a menu item to toggle the enabled state of the disabled item
	disabledMenu.Add("Toggle Enabled State").OnClick(func(*application.Context) {
		disabledItem.SetEnabled(!disabledItem.Enabled())
		menu.Update()
		println("Toggled enabled state - Now enabled:", disabledItem.Enabled())
	})

	// Create a menu update demo menu
	menuUpdateDemo := menu.AddSubmenu("Menu Update Demo")

	// Add some initial items to the menu
	menuUpdateDemo.Add("Original Item 1")
	menuUpdateDemo.Add("Original Item 2")
	menuUpdateDemo.Add("Original Item 3")

	// Add a button to reset the menu
	menuUpdateDemo.Add("Reset menu").OnClick(func(ctx *application.Context) {
		menuUpdateDemo.Clear()
		menuUpdateDemo.Add("New option")
		menu.Update()
		// Explicitly set the menu on the application and window to ensure it's updated
		app.SetMenu(menu)
		if app.CurrentWindow() != nil {
			app.CurrentWindow().SetMenu(menu)
		}
		println("Menu has been reset")
	})

	// Create a top-level menu for the submenu label update functionality
	labelMenu := menu.AddSubmenu("Label Update")

	// Create a submenu for label operations
	labelOperations := labelMenu.AddSubmenu("Label Operations")

	// Add a button to the parent menu to update the submenu's label from outside
	labelMenu.Add("Update Submenu Label From Parent").OnClick(func(*application.Context) {
		// Find the submenu item in the Label Update menu
		submenuItem := labelMenu.FindByLabel("Label Operations")
		if submenuItem != nil {
			// Change the label
			submenuItem.SetLabel("Parent Updated Label")
			// Update the menu to reflect the changes
			menu.Update()
			println("Submenu label changed from parent to: Parent Updated Label")

			// Verify the change
			updatedItem := labelMenu.FindByLabel("Parent Updated Label")
			if updatedItem != nil {
				println("Verification: Found submenu with updated label:", updatedItem.Label())
			} else {
				println("Verification failed: Could not find submenu with updated label")
			}
		} else {
			println("Submenu not found!")
		}
	})

	// Add a button to update the submenu label in the Demo menu
	labelOperations.Add("Update Demo Submenu Label").OnClick(func(*application.Context) {
		// Find the submenu item in the Demo menu
		submenuItem := myMenu.FindByLabel("Submenu")
		if submenuItem != nil {
			// Change the label
			submenuItem.SetLabel("Updated Submenu Label")
			// Update the menu to reflect the changes
			menu.Update()
			println("Submenu label changed to: Updated Submenu Label")

			// Verify the change
			updatedItem := myMenu.FindByLabel("Updated Submenu Label")
			if updatedItem != nil {
				println("Verification: Found submenu with updated label:", updatedItem.Label())
			} else {
				println("Verification failed: Could not find submenu with updated label")
			}
		} else {
			println("Submenu not found!")
		}
	})

	// Add a button to update the Label Operations submenu's own label
	labelOperations.Add("Update This Submenu's Label").OnClick(func(*application.Context) {
		// Find the submenu item in the Label Update menu
		submenuItem := labelMenu.FindByLabel("Label Operations")
		if submenuItem != nil {
			// Change the label
			submenuItem.SetLabel("Updated Operations")
			// Update the menu to reflect the changes
			menu.Update()
			println("Submenu label changed to: Updated Operations")

			// Verify the change
			updatedItem := labelMenu.FindByLabel("Updated Operations")
			if updatedItem != nil {
				println("Verification: Found submenu with updated label:", updatedItem.Label())
			} else {
				println("Verification failed: Could not find submenu with updated label")
			}
		} else {
			println("Submenu not found!")
		}
	})

	// Create a nested submenu to demonstrate deeper nesting
	nestedSubmenu := labelOperations.AddSubmenu("Nested Submenu")
	nestedSubmenu.Add("Nested Item 1")
	nestedSubmenu.Add("Nested Item 2")

	// Add a button to the nested submenu to change its own label
	nestedSubmenu.Add("Change My Own Label").OnClick(func(*application.Context) {
		// Find the nested submenu item
		submenuItem := labelOperations.FindByLabel("Nested Submenu")
		if submenuItem != nil {
			// Change the label
			submenuItem.SetLabel("Self-Updated Submenu")
			// Update the menu to reflect the changes
			menu.Update()
			println("Nested submenu changed its own label to: Self-Updated Submenu")

			// Verify the change
			updatedItem := labelOperations.FindByLabel("Self-Updated Submenu")
			if updatedItem != nil {
				println("Verification: Found nested submenu with self-updated label:", updatedItem.Label())
			} else {
				println("Verification failed: Could not find nested submenu with self-updated label")
			}
		} else {
			println("Nested submenu not found!")
		}
	})

	// Add a button to update the nested submenu's label
	labelOperations.Add("Update Nested Submenu Label").OnClick(func(*application.Context) {
		// Find the nested submenu item
		submenuItem := labelOperations.FindByLabel("Nested Submenu")
		if submenuItem != nil {
			// Change the label
			submenuItem.SetLabel("Updated Nested Submenu")
			// Update the menu to reflect the changes
			menu.Update()
			println("Nested submenu label changed to: Updated Nested Submenu")

			// Verify the change
			updatedItem := labelOperations.FindByLabel("Updated Nested Submenu")
			if updatedItem != nil {
				println("Verification: Found nested submenu with updated label:", updatedItem.Label())
			} else {
				println("Verification failed: Could not find nested submenu with updated label")
			}
		} else {
			println("Nested submenu not found!")
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
