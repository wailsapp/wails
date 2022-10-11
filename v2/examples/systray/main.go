package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed iconLightMode.png
var lightModeIcon []byte

//go:embed iconDarkMode.png
var darkModeIcon []byte

func main() {

	var runtimeContext context.Context

	// Create a new Wails application using the current options
	mainApp := application.NewWithOptions(&options.App{
		Assets:            assets,
		StartHidden:       true,
		HideWindowOnClose: true,
		OnStartup: func(ctx context.Context) {
			runtimeContext = ctx
		},
	})

	// ------------------------------------
	// Create a systray for the application
	// Currently we only support PNG for icons

	systray := mainApp.NewSystemTray(&options.SystemTray{
		// This is the icon used when the system in using light mode
		LightModeIcon: &options.SystemTrayIcon{
			Data: lightModeIcon,
		},
		// This is the icon used when the system in using dark mode
		DarkModeIcon: &options.SystemTrayIcon{
			Data: darkModeIcon,
		},
		Tooltip: "Systray Example",
		OnLeftClick: func() {
			runtime.Show(runtimeContext)
		},
	})

	// ------------------------------------
	// Create some menus items that will be reused

	// This is a menuitem we will show/hide at runtime
	visibleNotVisible := menu.Label("visible?").Show()

	var counter int
	icons := [][]byte{lightModeIcon, darkModeIcon}
	iconCounter := 0

	disabledEnabledMenu := menu.Label("disabled").Disable().OnClick(func(c *menu.CallbackData) {
		println("Disabled item clicked!")
	})

	// This checkbox menuitem will print the current checked state to the console when clicked.
	// When a checkbox item is clicked, the state of the `Checked` variable is toggled.
	// The UI automatically reflects the current state, even if this item is used multiple times.
	mycheckbox := menu.Label("checked").SetChecked(true).OnClick(func(c *menu.CallbackData) {
		println("My checked state is: ", c.MenuItem.Checked)
	})

	// This radio callback will be used by all the radio items.
	// The CallbackData has a pointer back to the menuitem so we can determine
	// which item was selected
	radioCallback := func(data *menu.CallbackData) {
		println("Radio item clicked:", data.MenuItem.Label)
	}

	// We create 3 radio items , with the first being selected. They all share a callback.
	radio1 := menu.Radio("Radio 1", true, nil, radioCallback)
	radio2 := menu.Radio("Radio 2", false, nil, radioCallback)
	radio3 := menu.Radio("Radio 3", false, nil, radioCallback)

	// Now we set the menu of the systray
	systray.SetMenu(menu.NewMenuFromItems(
		visibleNotVisible,
		// This menu item changes its label when clicked.
		menu.Label("Click Me!").OnClick(func(c *menu.CallbackData) {
			c.MenuItem.SetLabel(fmt.Sprintf("Clicked %d times", counter))
			counter++
			systray.Update()
		}),

		// We add a checkbox
		menu.Separator(),
		mycheckbox,

		// Next we create 2 radio groups containing the same menu items.
		// It is perfectly fine to reuse radio item groups - the state and UI will
		// stay in sync. Warning: Using the same radio item in different groups will
		// lead to unspecified behaviour!
		menu.Separator(),
		radio1,
		radio2,
		radio3,

		menu.Separator(),
		mycheckbox,

		menu.Label("Toggle items!").OnClick(func(c *menu.CallbackData) {

			iconCounter++

			// Swap light and dark mode icons
			systray.SetIcons(&options.SystemTrayIcon{
				Data: icons[iconCounter%2],
			}, &options.SystemTrayIcon{
				Data: icons[(iconCounter+1)%2],
			})

			// Do some toggling
			if iconCounter%2 == 0 {
				visibleNotVisible.Show()
				disabledEnabledMenu.Disable()
			} else {
				visibleNotVisible.Hide()
				disabledEnabledMenu.Enable()
			}

			// Update the menu
			err := systray.Update()
			if err != nil {
				panic(err)
			}
		}),

		// We create a checkbox item that is initially unchecked.
		menu.Label("unchecked").SetChecked(false).OnClick(func(c *menu.CallbackData) {
			println("My checked state is: ", c.MenuItem.Checked)
			systray.SetTooltip("My updated tooltip!")
		}),

		// This menu item will toggle between enabled and disabled each time the "Toggle items!" menu
		// option is selected
		disabledEnabledMenu,

		// We now add a submenu, reusing the checkbox item and submenu we created earlier
		menu.SubMenu("submenu", menu.NewMenuFromItems(
			mycheckbox,
			menu.Label("submenu item").OnClick(func(data *menu.CallbackData) {
				println("submenu item clicked")
			}),
			menu.Separator(),
			radio1,
			radio2,
			radio3,
		)),
		menu.Separator(),
		menu.Label("quit").OnClick(func(_ *menu.CallbackData) {
			println("Quitting application")
			mainApp.Quit()
		}),
	))

	println("Check out the system tray!")

	// Now we run the application
	err := mainApp.Run()

	if err != nil {
		println("Error:", err.Error())
	}
}
