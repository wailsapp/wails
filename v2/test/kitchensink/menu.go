package main

import (
	"fmt"
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
	// 	parent.Append(menu.Text(menuText, menuText, menu.Accel("[")))

	// If this is the first dynamic menu added, let's add a remove menu item
	if counter == 1 {
		removeMenu := menu.Text("Remove "+menuText,
			"Remove Last Item", menu.CmdOrCtrlAccel("-"))
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

	// 	parent.Append(menu.Text(menuText, menuText, menu.Accel("[")))
	m.runtime.Menu.Update()
}

func (m *Menu) createDynamicMenuTwo() {

	// Create our submenu
	dm2 := menu.SubMenu("Dynamic Menus 2", []*menu.MenuItem{
		menu.Text("Insert Before Random Menu Item",
			"Insert Before Random", menu.CmdOrCtrlAccel("]")),
		menu.Text("Insert After Random Menu Item",
			"Insert After Random", menu.CmdOrCtrlAccel("[")),
		menu.Separator(),
	})

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

	windowMenu := menu.SubMenu("Test", []*menu.MenuItem{
		menu.Togglefullscreen(),
		menu.Minimize(),
		menu.Zoom(),

		menu.Separator(),

		menu.Copy(),
		menu.Cut(),
		menu.Delete(),

		menu.Separator(),

		menu.Front(),

		menu.SubMenu("Test Submenu", []*menu.MenuItem{
			menu.Text("Plain text", "plain text", nil),
			menu.Text("Show Dynamic Menus 2 Submenu", "show-dynamic-menus-2", nil),
			menu.SubMenu("Accelerators", []*menu.MenuItem{
				menu.SubMenu("Modifiers", []*menu.MenuItem{
					menu.Text("Shift accelerator", "Shift", menu.ShiftAccel("o")),
					menu.Text("Control accelerator", "Control", menu.ControlAccel("o")),
					menu.Text("Command accelerator", "Command", menu.CmdOrCtrlAccel("o")),
					menu.Text("Option accelerator", "Option", menu.OptionOrAltAccel("o")),
				}),
				menu.SubMenu("System Keys", []*menu.MenuItem{
					menu.Text("Backspace", "Backspace", menu.Accel("Backspace")),
					menu.Text("Tab", "Tab", menu.Accel("Tab")),
					menu.Text("Return", "Return", menu.Accel("Return")),
					menu.Text("Escape", "Escape", menu.Accel("Escape")),
					menu.Text("Left", "Left", menu.Accel("Left")),
					menu.Text("Right", "Right", menu.Accel("Right")),
					menu.Text("Up", "Up", menu.Accel("Up")),
					menu.Text("Down", "Down", menu.Accel("Down")),
					menu.Text("Space", "Space", menu.Accel("Space")),
					menu.Text("Delete", "Delete", menu.Accel("Delete")),
					menu.Text("Home", "Home", menu.Accel("Home")),
					menu.Text("End", "End", menu.Accel("End")),
					menu.Text("Page Up", "Page Up", menu.Accel("Page Up")),
					menu.Text("Page Down", "Page Down", menu.Accel("Page Down")),
					menu.Text("NumLock", "NumLock", menu.Accel("NumLock")),
				}),
				menu.SubMenu("Function Keys", []*menu.MenuItem{
					menu.Text("F1", "F1", menu.Accel("F1")),
					menu.Text("F2", "F2", menu.Accel("F2")),
					menu.Text("F3", "F3", menu.Accel("F3")),
					menu.Text("F4", "F4", menu.Accel("F4")),
					menu.Text("F5", "F5", menu.Accel("F5")),
					menu.Text("F6", "F6", menu.Accel("F6")),
					menu.Text("F7", "F7", menu.Accel("F7")),
					menu.Text("F8", "F8", menu.Accel("F8")),
					menu.Text("F9", "F9", menu.Accel("F9")),
					menu.Text("F10", "F10", menu.Accel("F10")),
					menu.Text("F11", "F11", menu.Accel("F11")),
					menu.Text("F12", "F12", menu.Accel("F12")),
					menu.Text("F13", "F13", menu.Accel("F13")),
					menu.Text("F14", "F14", menu.Accel("F14")),
					menu.Text("F15", "F15", menu.Accel("F15")),
					menu.Text("F16", "F16", menu.Accel("F16")),
					menu.Text("F17", "F17", menu.Accel("F17")),
					menu.Text("F18", "F18", menu.Accel("F18")),
					menu.Text("F19", "F19", menu.Accel("F19")),
					menu.Text("F20", "F20", menu.Accel("F20")),
				}),
				menu.SubMenu("Standard Keys", []*menu.MenuItem{
					menu.Text("Backtick", "Backtick", menu.Accel("`")),
					menu.Text("Plus", "Plus", menu.Accel("+")),
				}),
			}),
			menu.SubMenuWithID("Dynamic Menus 1", "Dynamic Menus 1", []*menu.MenuItem{
				menu.Text("Add Menu Item", "Add Menu Item", menu.CmdOrCtrlAccel("+")),
				menu.Separator(),
			}),
			{
				Label:       "Disabled Menu",
				Type:        menu.TextType,
				Accelerator: menu.ComboAccel("p", menu.CmdOrCtrl, menu.Shift),
				Disabled:    true,
			},
			{
				Label:  "Hidden Menu",
				Type:   menu.TextType,
				Hidden: true,
			},
			{
				ID:          "checkbox-menu 1",
				Label:       "Checkbox Menu 1",
				Type:        menu.CheckboxType,
				Accelerator: menu.CmdOrCtrlAccel("l"),
				Checked:     true,
			},
			menu.Checkbox("Checkbox Menu 2", "checkbox-menu 2", false, nil),
			menu.Separator(),
			menu.Radio("😀 Option 1", "😀option-1", true, nil),
			menu.Radio("😺 Option 2", "option-2", false, nil),
			menu.Radio("❤️ Option 3", "option-3", false, nil),
		}),
	})

	myMenu.Append(windowMenu)
	return myMenu
}
