//go:build windows
// +build windows

package windows

import (
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

var checkboxMap = map[*menu.MenuItem][]*winc.MenuItem{}
var radioGroupMap = map[*menu.MenuItem][]*winc.MenuItem{}

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

func toggleRadioItem(menuItem *menu.MenuItem) {
	menuItem.Checked = !menuItem.Checked
	for _, wincMenu := range radioGroupMap[menuItem] {
		wincMenu.SetChecked(menuItem.Checked)
	}
}

func addRadioItemToMap(menuItem *menu.MenuItem, wincMenuItem *winc.MenuItem) {
	if radioGroupMap[menuItem] == nil {
		radioGroupMap[menuItem] = []*winc.MenuItem{}
	}
	radioGroupMap[menuItem] = append(radioGroupMap[menuItem], wincMenuItem)
}

func (w *Window) SetApplicationMenu(menu *menu.Menu) {
	w.applicationMenu = menu
	processMenu(w, menu)
}

func processMenu(window *Window, menu *menu.Menu) {
	mainMenu := window.NewMenu()
	for _, menuItem := range menu.Items {
		submenu := mainMenu.AddSubMenu(menuItem.Label)
		if menuItem.SubMenu != nil {
			for _, menuItem := range menuItem.SubMenu.Items {
				processMenuItem(submenu, menuItem)
			}
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
		parent.AddSeparator()
	case menu.TextType:
		shortcut := acceleratorToWincShortcut(menuItem.Accelerator)
		newItem := parent.AddItem(menuItem.Label, shortcut)
		//if menuItem.Tooltip != "" {
		//	newItem.SetToolTip(menuItem.Tooltip)
		//}
		if menuItem.Click != nil {
			newItem.OnClick().Bind(func(e *winc.Event) {
				menuItem.Click(&menu.CallbackData{
					MenuItem: menuItem,
				})
			})
		}
		newItem.SetEnabled(!menuItem.Disabled)

	case menu.CheckboxType:
		shortcut := acceleratorToWincShortcut(menuItem.Accelerator)
		newItem := parent.AddItem(menuItem.Label, shortcut)
		newItem.SetCheckable(true)
		newItem.SetChecked(menuItem.Checked)
		//if menuItem.Tooltip != "" {
		//	newItem.SetToolTip(menuItem.Tooltip)
		//}
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
		shortcut := acceleratorToWincShortcut(menuItem.Accelerator)
		newItem := parent.AddItemRadio(menuItem.Label, shortcut)
		newItem.SetCheckable(true)
		newItem.SetChecked(menuItem.Checked)
		//if menuItem.Tooltip != "" {
		//	newItem.SetToolTip(menuItem.Tooltip)
		//}
		if menuItem.Click != nil {
			newItem.OnClick().Bind(func(e *winc.Event) {
				toggleRadioItem(menuItem)
				menuItem.Click(&menu.CallbackData{
					MenuItem: menuItem,
				})
			})
		}
		newItem.SetEnabled(!menuItem.Disabled)
		addRadioItemToMap(menuItem, newItem)
	case menu.SubmenuType:
		submenu := parent.AddSubMenu(menuItem.Label)
		for _, menuItem := range menuItem.SubMenu.Items {
			processMenuItem(submenu, menuItem)
		}
	}
}

func (f *Frontend) MenuSetApplicationMenu(menu *menu.Menu) {
	f.mainWindow.SetApplicationMenu(menu)
}

func (f *Frontend) MenuUpdateApplicationMenu() {
	processMenu(f.mainWindow, f.mainWindow.applicationMenu)
}
