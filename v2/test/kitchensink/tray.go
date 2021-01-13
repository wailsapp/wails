package main

import (
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"strconv"
	"sync"
)

// Tray struct
type Tray struct {
	runtime *wails.Runtime

	dynamicMenuCounter int
	lock               sync.Mutex
	dynamicMenuItems   map[string]*menu.MenuItem

	trayMenu       *menu.TrayMenu
	secondTrayMenu *menu.TrayMenu

	done bool
}

// WailsInit is called at application startup
func (t *Tray) WailsInit(runtime *wails.Runtime) error {
	// Perform your setup here
	t.runtime = runtime

	//// Auto switch between light / dark tray icons
	//t.runtime.Events.OnThemeChange(func(darkMode bool) {
	//	if darkMode {
	//		t.runtime.Tray.SetIcon("light")
	//		return
	//	}
	//
	//	t.runtime.Tray.SetIcon("dark")
	//})

	return nil
}

func (t *Tray) showWindow(_ *menu.CallbackData) {
	t.runtime.Window.Show()
}

func (t *Tray) hideWindow(_ *menu.CallbackData) {
	t.runtime.Window.Hide()
}

func (t *Tray) unminimiseWindow(_ *menu.CallbackData) {
	t.runtime.Window.Unminimise()
}

func (t *Tray) minimiseWindow(_ *menu.CallbackData) {
	t.runtime.Window.Minimise()
}

func (t *Tray) WailsShutdown() {
	t.done = true
}

func (t *Tray) incrementcounter() int {
	t.dynamicMenuCounter++
	return t.dynamicMenuCounter
}

func (t *Tray) decrementcounter() int {
	t.dynamicMenuCounter--
	return t.dynamicMenuCounter
}

func (t *Tray) SvelteIcon(_ *menu.CallbackData) {
	t.secondTrayMenu.Icon = "svelte"
	t.runtime.Menu.UpdateTrayMenu(t.secondTrayMenu)
}
func (t *Tray) NoIcon(_ *menu.CallbackData) {
	t.secondTrayMenu.Icon = ""
	t.runtime.Menu.UpdateTrayMenu(t.secondTrayMenu)
}
func (t *Tray) LightIcon(_ *menu.CallbackData) {
	t.secondTrayMenu.Icon = "light"
	t.runtime.Menu.UpdateTrayMenu(t.secondTrayMenu)
}
func (t *Tray) DarkIcon(_ *menu.CallbackData) {
	t.secondTrayMenu.Icon = "dark"
	t.runtime.Menu.UpdateTrayMenu(t.secondTrayMenu)
}

//func (t *Tray) removeMenu(_ *menu.MenuItem) {
//
//	// Lock because this method will be called in a goroutine
//	t.lock.Lock()
//	defer t.lock.Unlock()
//
//	// Get the id of the last dynamic menu
//	menuID := "Dynamic Menu Item " + strconv.Itoa(t.dynamicMenuCounter)
//
//	// Remove the last menu item by ID
//	t.runtime.Tray.RemoveByID(menuID)
//
//	// Update the counter
//	counter := t.decrementcounter()
//
//	// If we deleted the last dynamic menu, remove the "Remove Last Item" menu
//	if counter == 0 {
//		t.runtime.Tray.RemoveByID("Remove Last Item")
//	} else {
//		// Update label
//		menuText := "Dynamic Menu Item " + strconv.Itoa(counter)
//		removeMenu := t.runtime.Tray.GetByID("Remove Last Item")
//		// Test if the remove menu hasn't already been removed in another thread
//		if removeMenu == nil {
//			return
//		}
//		removeMenu.Label = "Remove " + menuText
//	}
//
//	// 	parent.Append(menu.Text(menuText, menuText, menu.Key("[")))
//	t.runtime.Tray.Update()
//}

//func (t *Tray) SetIcon(trayIconID string) {
//	t.runtime.Tray.SetIcon(trayIconID)
//}

func (t *Tray) createTrayMenus() []*menu.TrayMenu {
	trayMenu := &menu.TrayMenu{}
	trayMenu.Label = "Test Tray Label"
	trayMenu.Menu = menu.NewMenuFromItems(
		menu.Text("Show Window", nil, t.showWindow),
		menu.Text("Hide Window", nil, t.hideWindow),
		menu.Text("Minimise Window", nil, t.minimiseWindow),
		menu.Text("Unminimise Window", nil, t.unminimiseWindow),
	)
	t.trayMenu = trayMenu

	secondTrayMenu := &menu.TrayMenu{}
	secondTrayMenu.Label = "Another tray label"
	secondTrayMenu.Icon = "svelte"
	secondTrayMenu.Menu = menu.NewMenuFromItems(
		menu.Text("Update Label", nil, func(_ *menu.CallbackData) {
			// Lock because this method will be called in a goroutine
			t.lock.Lock()
			defer t.lock.Unlock()

			counter := t.incrementcounter()
			trayLabel := "Updated Label " + strconv.Itoa(counter)
			secondTrayMenu.Label = trayLabel
			t.runtime.Menu.UpdateTrayMenu(t.secondTrayMenu)
		}),
		menu.SubMenu("Select Icon", menu.NewMenuFromItems(
			menu.Text("Svelte", nil, t.SvelteIcon),
			menu.Text("Light", nil, t.LightIcon),
			menu.Text("Dark", nil, t.DarkIcon),
			menu.Text("None", nil, t.NoIcon),
		)),
	)
	t.secondTrayMenu = secondTrayMenu
	return []*menu.TrayMenu{
		trayMenu,
		secondTrayMenu,
	}
}
