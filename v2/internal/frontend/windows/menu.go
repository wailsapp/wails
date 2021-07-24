package windows

import (
	"fmt"
	"github.com/tadvi/winc"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

var checkboxMap = map[*menu.MenuItem][]*winc.MenuItem{}

func toggleCheckBox(menuItem *menu.MenuItem) {
	menuItem.Checked = !menuItem.Checked
	for _, wincMenu := range checkboxMap[menuItem] {
		wincMenu.SetChecked(menuItem.Checked)
	}
}

func addCheckBoxToMap(menuItem *menu.MenuItem, wincMenuItem *winc.MenuItem) {
	if checkboxMap[menuItem] == nil {
		checkboxMap[menuItem] = []*winc.MenuItem{}
	}
	checkboxMap[menuItem] = append(checkboxMap[menuItem], wincMenuItem)
}

func processApplicationMenu(window *Window, menuToProcess *menu.Menu) {
	mainMenu := window.NewMenu()
	for _, menuItem := range menuToProcess.Items {
		println("Adding menu:", menuItem.Label)
		submenu := mainMenu.AddSubMenu(menuItem.Label)
		for _, menuItem := range menuItem.SubMenu.Items {
			fmt.Printf("Processing: %#v\n", menuItem)
			processMenuItem(submenu, menuItem)
		}
	}
	mainMenu.Show()
}

func processMenuItem(parent *winc.MenuItem, menuItem *menu.MenuItem) {
	if menuItem.Hidden {
		return
	}
	switch menuItem.Type {
	case menu.SeparatorType:
		parent.SetSeparator()
	case menu.TextType:
		newItem := parent.AddItem(menuItem.Label, winc.NoShortcut)
		if menuItem.Tooltip != "" {
			newItem.SetToolTip(menuItem.Tooltip)
		}
		if menuItem.Click != nil {
			newItem.OnClick().Bind(func(e *winc.Event) {
				menuItem.Click(&menu.CallbackData{
					MenuItem: menuItem,
				})
			})
		}
		newItem.SetEnabled(!menuItem.Disabled)

	case menu.CheckboxType:
		newItem := parent.AddItem(menuItem.Label, winc.NoShortcut)
		newItem.SetCheckable(true)
		newItem.SetChecked(menuItem.Checked)
		if menuItem.Tooltip != "" {
			newItem.SetToolTip(menuItem.Tooltip)
		}
		if menuItem.Click != nil {
			newItem.OnClick().Bind(func(e *winc.Event) {
				toggleCheckBox(menuItem)
				menuItem.Click(&menu.CallbackData{
					MenuItem: menuItem,
				})
			})
		}
		newItem.SetEnabled(!menuItem.Disabled)
		addCheckBoxToMap(menuItem, newItem)
	case menu.RadioType:
	case menu.SubmenuType:
		submenu := parent.AddSubMenu(menuItem.Label)
		for _, menuItem := range menuItem.SubMenu.Items {
			processMenuItem(submenu, menuItem)
		}
	}
}
