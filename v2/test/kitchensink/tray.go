package main

import (
	"strconv"
	"sync"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

// Tray struct
type Tray struct {
	runtime *wails.Runtime

	dynamicMenuCounter        int
	lock                      sync.Mutex
	dynamicMenuItems          map[string]*menu.MenuItem
	anotherDynamicMenuCounter int
}

// WailsInit is called at application startup
func (t *Tray) WailsInit(runtime *wails.Runtime) error {
	// Perform your setup here
	t.runtime = runtime

	// Setup Menu Listeners
	t.runtime.Tray.On("Show Window", func(mi *menu.MenuItem) {
		t.runtime.Window.Show()
	})
	t.runtime.Tray.On("Hide Window", func(mi *menu.MenuItem) {
		t.runtime.Window.Hide()
	})

	return nil
}

func (t *Tray) incrementcounter() int {
	t.dynamicMenuCounter++
	return t.dynamicMenuCounter
}

func (t *Tray) decrementcounter() int {
	t.dynamicMenuCounter--
	return t.dynamicMenuCounter
}

func (t *Tray) addMenu(mi *menu.MenuItem) {

	// Lock because this method will be called in a gorouting
	t.lock.Lock()
	defer t.lock.Unlock()

	// Get this menu's parent
	parent := mi.Parent()
	counter := t.incrementcounter()
	menuText := "Dynamic Menu Item " + strconv.Itoa(counter)
	parent.Append(menu.Text(menuText, menuText))
	// 	parent.Append(menu.TextWithAccelerator(menuText, menuText, menu.Accel("[")))

	// If this is the first dynamic menu added, let's add a remove menu item
	if counter == 1 {
		removeMenu := menu.TextWithAccelerator("Remove "+menuText,
			"Remove Last Item", menu.CmdOrCtrlAccel("-"))
		parent.Prepend(removeMenu)
		t.runtime.Tray.On("Remove Last Item", t.removeMenu)
	} else {
		removeMenu := t.runtime.Tray.GetByID("Remove Last Item")
		// Test if the remove menu hasn't already been removed in another thread
		if removeMenu != nil {
			removeMenu.Label = "Remove " + menuText
		}
	}
	t.runtime.Tray.Update()
}

func (t *Tray) removeMenu(_ *menu.MenuItem) {

	// Lock because this method will be called in a goroutine
	t.lock.Lock()
	defer t.lock.Unlock()

	// Get the id of the last dynamic menu
	menuID := "Dynamic Menu Item " + strconv.Itoa(t.dynamicMenuCounter)

	// Remove the last menu item by ID
	t.runtime.Tray.RemoveByID(menuID)

	// Update the counter
	counter := t.decrementcounter()

	// If we deleted the last dynamic menu, remove the "Remove Last Item" menu
	if counter == 0 {
		t.runtime.Tray.RemoveByID("Remove Last Item")
	} else {
		// Update label
		menuText := "Dynamic Menu Item " + strconv.Itoa(counter)
		removeMenu := t.runtime.Tray.GetByID("Remove Last Item")
		// Test if the remove menu hasn't already been removed in another thread
		if removeMenu == nil {
			return
		}
		removeMenu.Label = "Remove " + menuText
	}

	// 	parent.Append(menu.TextWithAccelerator(menuText, menuText, menu.Accel("[")))
	t.runtime.Tray.Update()
}

func createApplicationTray() *menu.Menu {
	trayMenu := &menu.Menu{}
	trayMenu.Append(menu.Text("Show Window", "Show Window"))
	trayMenu.Append(menu.Text("Hide Window", "Hide Window"))
	return trayMenu
}
