package main

import (
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"math/rand"
	"strconv"
	"sync"
	"time"

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

	done bool
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

	t.runtime.Tray.On("Minimise Window", func(mi *menu.MenuItem) {
		t.runtime.Window.Minimise()
	})

	t.runtime.Tray.On("Unminimise Window", func(mi *menu.MenuItem) {
		t.runtime.Window.Unminimise()
	})

	// Auto switch between light / dark tray icons
	t.runtime.Events.OnThemeChange(func(darkMode bool) {
		if darkMode {
			t.runtime.Tray.SetIcon("light")
			return
		}

		t.runtime.Tray.SetIcon("dark")
	})

	// Start ticker
	go t.startTicker()

	return nil
}

func (t *Tray) WailsShutdown() {
	t.done = true
}

func (t *Tray) startTicker() {
	time.Sleep(1 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	for t.done == false {
		select {
		case <-ticker.C:
			r := rand.Intn(100)
			t.runtime.Tray.SetLabel(fmt.Sprintf("CPU: %d", r))
		}
	}
	ticker.Stop()
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
	parent.Append(menu.Text(menuText, menuText, nil))
	// 	parent.Append(menu.Text(menuText, menuText, menu.Key("[")))

	// If this is the first dynamic menu added, let's add a remove menu item
	if counter == 1 {
		removeMenu := menu.Text("Remove "+menuText,
			"Remove Last Item", keys.CmdOrCtrl("-"))
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

	// 	parent.Append(menu.Text(menuText, menuText, menu.Key("[")))
	t.runtime.Tray.Update()
}

func createApplicationTray() *menu.Menu {
	trayMenu := &menu.Menu{}
	trayMenu.Append(menu.Text("Show Window", "Show Window", nil))
	trayMenu.Append(menu.Text("Hide Window", "Hide Window", nil))
	trayMenu.Append(menu.Text("Minimise Window", "Minimise Window", nil))
	trayMenu.Append(menu.Text("Unminimise Window", "Unminimise Window", nil))
	return trayMenu
}
