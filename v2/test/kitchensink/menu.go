package main

import (
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"math/rand"
	"strconv"
	"sync"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

// Menu struct
type Menu struct {
	runtime *wails.Runtime

	dynamicMenuCounter        int
	lock                      sync.Mutex
	dynamicMenuItems          map[string]*menu.MenuItem
	anotherDynamicMenuCounter int
}

// WailsInit is called at application startup
func (m *Menu) WailsInit(runtime *wails.Runtime) error {
	// Perform your setup here
	m.runtime = runtime

	// Setup Menu Listeners
	m.runtime.Menu.On("hello", func(mi *menu.MenuItem) {
		fmt.Printf("The '%s' menu was clicked\n", mi.Label)
	})
	m.runtime.Menu.On("checkbox-menu", func(mi *menu.MenuItem) {
		fmt.Printf("The '%s' menu was clicked\n", mi.Label)
		fmt.Printf("It is now %v\n", mi.Checked)
	})
	m.runtime.Menu.On("😀option-1", func(mi *menu.MenuItem) {
		fmt.Printf("We can use UTF-8 IDs: %s\n", mi.Label)
	})

	m.runtime.Menu.On("show-dynamic-menus-2", func(mi *menu.MenuItem) {
		mi.Hidden = true
		// Create dynamic menu items 2 submenu
		m.createDynamicMenuTwo()
	})

	// Setup dynamic menus
	m.runtime.Menu.On("Add Menu Item", m.addMenu)
	return nil
}

func (m *Menu) incrementcounter() int {
	m.dynamicMenuCounter++
	return m.dynamicMenuCounter
}

func (m *Menu) decrementcounter() int {
	m.dynamicMenuCounter--
	return m.dynamicMenuCounter
}

func (m *Menu) addMenu(mi *menu.MenuItem) {

	// Lock because this method will be called in a gorouting
	m.lock.Lock()
	defer m.lock.Unlock()

	// Get this menu's parent
	parent := mi.Parent()
	counter := m.incrementcounter()
	menuText := "Dynamic Menu Item " + strconv.Itoa(counter)
	parent.Append(menu.Text(menuText, menuText, nil))
	// 	parent.Append(menu.Text(menuText, menuText, menu.Key("[")))

	// If this is the first dynamic menu added, let's add a remove menu item
	if counter == 1 {
		removeMenu := menu.Text("Remove "+menuText,
			"Remove Last Item", keys.CmdOrCtrl("-"))
		parent.Prepend(removeMenu)
		m.runtime.Menu.On("Remove Last Item", m.removeMenu)
	} else {
		removeMenu := m.runtime.Menu.GetByID("Remove Last Item")
		// Test if the remove menu hasn't already been removed in another thread
		if removeMenu != nil {
			removeMenu.Label = "Remove " + menuText
		}
	}
	m.runtime.Menu.Update()
}

func (m *Menu) removeMenu(_ *menu.MenuItem) {

	// Lock because this method will be called in a goroutine
	m.lock.Lock()
	defer m.lock.Unlock()

	// Get the id of the last dynamic menu
	menuID := "Dynamic Menu Item " + strconv.Itoa(m.dynamicMenuCounter)

	// Remove the last menu item by ID
	m.runtime.Menu.RemoveByID(menuID)

	// Update the counter
	counter := m.decrementcounter()

	// If we deleted the last dynamic menu, remove the "Remove Last Item" menu
	if counter == 0 {
		m.runtime.Menu.RemoveByID("Remove Last Item")
	} else {
		// Update label
		menuText := "Dynamic Menu Item " + strconv.Itoa(counter)
		removeMenu := m.runtime.Menu.GetByID("Remove Last Item")
		// Test if the remove menu hasn't already been removed in another thread
		if removeMenu == nil {
			return
		}
		removeMenu.Label = "Remove " + menuText
	}

	// 	parent.Append(menu.Text(menuText, menuText, menu.Key("[")))
	m.runtime.Menu.Update()
}

func (m *Menu) createDynamicMenuTwo() {

	// Create our submenu
	dm2 := menu.SubMenu("Dynamic Menus 2", menu.NewMenuFromItems(
		menu.Text("Insert Before Random Menu Item",
			"Insert Before Random", keys.CmdOrCtrl("]")),
		menu.Text("Insert After Random Menu Item",
			"Insert After Random", keys.CmdOrCtrl("[")),
		menu.Separator(),
	))

	m.runtime.Menu.On("Insert Before Random", m.insertBeforeRandom)
	m.runtime.Menu.On("Insert After Random", m.insertAfterRandom)

	// Initialise out map
	m.dynamicMenuItems = make(map[string]*menu.MenuItem)

	// Create some random menu items
	m.anotherDynamicMenuCounter = 5
	for index := 0; index < m.anotherDynamicMenuCounter; index++ {
		text := "Other Dynamic Menu Item " + strconv.Itoa(index+1)
		item := menu.Text(text, text, nil)
		m.dynamicMenuItems[text] = item
		dm2.Append(item)
	}

	// Insert this menu after Dynamic Menu Item 1
	dm1 := m.runtime.Menu.GetByID("Dynamic Menus 1")
	if dm1 == nil {
		return
	}

	dm1.InsertAfter(dm2)
	m.runtime.Menu.Update()
}

func (m *Menu) insertBeforeRandom(_ *menu.MenuItem) {

	// Lock because this method will be called in a goroutine
	m.lock.Lock()
	defer m.lock.Unlock()

	// Pick a random menu
	var randomItemID string
	var count int
	var random = rand.Intn(len(m.dynamicMenuItems))
	for randomItemID = range m.dynamicMenuItems {
		if count == random {
			break
		}
		count++
	}
	m.anotherDynamicMenuCounter++
	text := "Other Dynamic Menu Item " + strconv.Itoa(
		m.anotherDynamicMenuCounter+1)
	newItem := menu.Text(text, text, nil)
	m.dynamicMenuItems[text] = newItem

	item := m.runtime.Menu.GetByID(randomItemID)
	if item == nil {
		return
	}

	m.runtime.Log.Info(fmt.Sprintf(
		"Inserting menu item '%s' before menu item '%s'", newItem.Label,
		item.Label))

	item.InsertBefore(newItem)
	m.runtime.Menu.Update()
}

func (m *Menu) insertAfterRandom(_ *menu.MenuItem) {

	// Lock because this method will be called in a goroutine
	m.lock.Lock()
	defer m.lock.Unlock()

	// Pick a random menu
	var randomItemID string
	var count int
	var random = rand.Intn(len(m.dynamicMenuItems))
	for randomItemID = range m.dynamicMenuItems {
		if count == random {
			break
		}
		count++
	}
	m.anotherDynamicMenuCounter++
	text := "Other Dynamic Menu Item " + strconv.Itoa(
		m.anotherDynamicMenuCounter+1)
	newItem := menu.Text(text, text, nil)

	item := m.runtime.Menu.GetByID(randomItemID)
	m.dynamicMenuItems[text] = newItem

	m.runtime.Log.Info(fmt.Sprintf(
		"Inserting menu item '%s' after menu item '%s'", newItem.Label,
		item.Label))

	item.InsertAfter(newItem)
	m.runtime.Menu.Update()
}

func createApplicationMenu() *menu.Menu {

	// Create menu
	myMenu := menu.DefaultMacMenu()

	windowMenu := menu.SubMenu("Test", menu.NewMenuFromItems(
		menu.Togglefullscreen(),
		menu.Minimize(),
		menu.Zoom(),

		menu.Separator(),

		menu.Copy(),
		menu.Cut(),
		menu.Delete(),

		menu.Separator(),

		menu.Front(),

		menu.SubMenu("Test Submenu", menu.NewMenuFromItems(
			menu.Text("Plain text", "plain text", nil),
			menu.Text("Show Dynamic Menus 2 Submenu", "show-dynamic-menus-2", nil),
			menu.SubMenu("Accelerators", menu.NewMenuFromItems(
				menu.SubMenu("Modifiers", menu.NewMenuFromItems(
					menu.Text("Shift accelerator", "Shift", keys.Shift("o")),
					menu.Text("Control accelerator", "Control", keys.Control("o")),
					menu.Text("Command accelerator", "Command", keys.CmdOrCtrl("o")),
					menu.Text("Option accelerator", "Option", keys.OptionOrAlt("o")),
					menu.Text("Combo accelerator", "Combo", keys.Combo("o", keys.CmdOrCtrlKey, keys.ShiftKey)),
				)),
				menu.SubMenu("System Keys", menu.NewMenuFromItems(
					menu.Text("Backspace", "Backspace", keys.Key("Backspace")),
					menu.Text("Tab", "Tab", keys.Key("Tab")),
					menu.Text("Return", "Return", keys.Key("Return")),
					menu.Text("Escape", "Escape", keys.Key("Escape")),
					menu.Text("Left", "Left", keys.Key("Left")),
					menu.Text("Right", "Right", keys.Key("Right")),
					menu.Text("Up", "Up", keys.Key("Up")),
					menu.Text("Down", "Down", keys.Key("Down")),
					menu.Text("Space", "Space", keys.Key("Space")),
					menu.Text("Delete", "Delete", keys.Key("Delete")),
					menu.Text("Home", "Home", keys.Key("Home")),
					menu.Text("End", "End", keys.Key("End")),
					menu.Text("Page Up", "Page Up", keys.Key("Page Up")),
					menu.Text("Page Down", "Page Down", keys.Key("Page Down")),
					menu.Text("NumLock", "NumLock", keys.Key("NumLock")),
				)),
				menu.SubMenu("Function Keys", menu.NewMenuFromItems(
					menu.Text("F1", "F1", keys.Key("F1")),
					menu.Text("F2", "F2", keys.Key("F2")),
					menu.Text("F3", "F3", keys.Key("F3")),
					menu.Text("F4", "F4", keys.Key("F4")),
					menu.Text("F5", "F5", keys.Key("F5")),
					menu.Text("F6", "F6", keys.Key("F6")),
					menu.Text("F7", "F7", keys.Key("F7")),
					menu.Text("F8", "F8", keys.Key("F8")),
					menu.Text("F9", "F9", keys.Key("F9")),
					menu.Text("F10", "F10", keys.Key("F10")),
					menu.Text("F11", "F11", keys.Key("F11")),
					menu.Text("F12", "F12", keys.Key("F12")),
					menu.Text("F13", "F13", keys.Key("F13")),
					menu.Text("F14", "F14", keys.Key("F14")),
					menu.Text("F15", "F15", keys.Key("F15")),
					menu.Text("F16", "F16", keys.Key("F16")),
					menu.Text("F17", "F17", keys.Key("F17")),
					menu.Text("F18", "F18", keys.Key("F18")),
					menu.Text("F19", "F19", keys.Key("F19")),
					menu.Text("F20", "F20", keys.Key("F20")),
				)),
				menu.SubMenu("Standard Keys", menu.NewMenuFromItems(
					menu.Text("Backtick", "Backtick", keys.Key("`")),
					menu.Text("Plus", "Plus", keys.Key("+")),
				)),
			)),
			menu.SubMenuWithID("Dynamic Menus 1", "Dynamic Menus 1", menu.NewMenuFromItems(
				menu.Text("Add Menu Item", "Add Menu Item", keys.CmdOrCtrl("+")),
				menu.Separator(),
			)),
			&menu.MenuItem{
				Label:       "Disabled Menu",
				Type:        menu.TextType,
				Accelerator: keys.Combo("p", keys.CmdOrCtrlKey, keys.ShiftKey),
				Disabled:    true,
			},
			&menu.MenuItem{
				Label:  "Hidden Menu",
				Type:   menu.TextType,
				Hidden: true,
			},
			&menu.MenuItem{
				ID:          "checkbox-menu 1",
				Label:       "Checkbox Menu 1",
				Type:        menu.CheckboxType,
				Accelerator: keys.CmdOrCtrl("l"),
				Checked:     true,
			},
			menu.Checkbox("Checkbox Menu 2", "checkbox-menu 2", false, nil),
			menu.Separator(),
			menu.Radio("😀 Option 1", "😀option-1", true, nil),
			menu.Radio("😺 Option 2", "option-2", false, nil),
			menu.Radio("❤️ Option 3", "option-3", false, nil),
		)),
	))

	myMenu.Append(windowMenu)
	return myMenu
}
