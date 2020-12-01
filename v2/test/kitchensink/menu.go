package main

import (
	"fmt"

	wails "github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

// Menu struct
type Menu struct {
	runtime *wails.Runtime
}

// WailsInit is called at application startup
func (m *Menu) WailsInit(runtime *wails.Runtime) error {
	// Perform your setup here
	m.runtime = runtime

	// Setup Menu Listeners
	m.runtime.Menu.On("hello", func(m *menu.MenuItem) {
		fmt.Printf("The '%s' menu was clicked\n", m.Label)
	})
	m.runtime.Menu.On("checkbox-menu", func(m *menu.MenuItem) {
		fmt.Printf("The '%s' menu was clicked\n", m.Label)
		fmt.Printf("It is now %v\n", m.Checked)
		// m.Checked = false
		// runtime.Menu.Update()
	})
	m.runtime.Menu.On("😀option-1", func(m *menu.MenuItem) {
		fmt.Printf("We can use UTF-8 IDs: %s\n", m.Label)
	})
	return nil
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
			menu.SubMenu("Accelerators", []*menu.MenuItem{
				menu.SubMenu("Modifiers", []*menu.MenuItem{
					menu.TextWithAccelerator("Shift accelerator", "Shift", menu.ShiftAccel("o")),
					menu.TextWithAccelerator("Control accelerator", "Control", menu.ControlAccel("o")),
					menu.TextWithAccelerator("Command accelerator", "Command", menu.CmdOrCtrlAccel("o")),
					menu.TextWithAccelerator("Option accelerator", "Option", menu.OptionOrAltAccel("o")),
				}),
				menu.SubMenu("System Keys", []*menu.MenuItem{
					menu.TextWithAccelerator("Backspace", "Backspace", menu.Accel("Backspace")),
					menu.TextWithAccelerator("Tab", "Tab", menu.Accel("Tab")),
					menu.TextWithAccelerator("Return", "Return", menu.Accel("Return")),
					menu.TextWithAccelerator("Escape", "Escape", menu.Accel("Escape")),
					menu.TextWithAccelerator("Left", "Left", menu.Accel("Left")),
					menu.TextWithAccelerator("Right", "Right", menu.Accel("Right")),
					menu.TextWithAccelerator("Up", "Up", menu.Accel("Up")),
					menu.TextWithAccelerator("Down", "Down", menu.Accel("Down")),
					menu.TextWithAccelerator("Space", "Space", menu.Accel("Space")),
					menu.TextWithAccelerator("Delete", "Delete", menu.Accel("Delete")),
					menu.TextWithAccelerator("Home", "Home", menu.Accel("Home")),
					menu.TextWithAccelerator("End", "End", menu.Accel("End")),
					menu.TextWithAccelerator("Page Up", "Page Up", menu.Accel("Page Up")),
					menu.TextWithAccelerator("Page Down", "Page Down", menu.Accel("Page Down")),
					menu.TextWithAccelerator("Insert", "Insert", menu.Accel("Insert")),
					menu.TextWithAccelerator("PrintScreen", "PrintScreen", menu.Accel("PrintScreen")),
					menu.TextWithAccelerator("ScrollLock", "ScrollLock", menu.Accel("ScrollLock")),
					menu.TextWithAccelerator("NumLock", "NumLock", menu.Accel("NumLock")),
				}),
				menu.SubMenu("Function Keys", []*menu.MenuItem{
					menu.TextWithAccelerator("F1", "F1", menu.Accel("F1")),
					menu.TextWithAccelerator("F2", "F2", menu.Accel("F2")),
					menu.TextWithAccelerator("F3", "F3", menu.Accel("F3")),
					menu.TextWithAccelerator("F4", "F4", menu.Accel("F4")),
					menu.TextWithAccelerator("F5", "F5", menu.Accel("F5")),
					menu.TextWithAccelerator("F6", "F6", menu.Accel("F6")),
					menu.TextWithAccelerator("F7", "F7", menu.Accel("F7")),
					menu.TextWithAccelerator("F8", "F8", menu.Accel("F8")),
					menu.TextWithAccelerator("F9", "F9", menu.Accel("F9")),
					menu.TextWithAccelerator("F10", "F10", menu.Accel("F10")),
					menu.TextWithAccelerator("F11", "F11", menu.Accel("F11")),
					menu.TextWithAccelerator("F12", "F12", menu.Accel("F12")),
					menu.TextWithAccelerator("F13", "F13", menu.Accel("F13")),
					menu.TextWithAccelerator("F14", "F14", menu.Accel("F14")),
					menu.TextWithAccelerator("F15", "F15", menu.Accel("F15")),
					menu.TextWithAccelerator("F16", "F16", menu.Accel("F16")),
					menu.TextWithAccelerator("F17", "F17", menu.Accel("F17")),
					menu.TextWithAccelerator("F18", "F18", menu.Accel("F18")),
					menu.TextWithAccelerator("F19", "F19", menu.Accel("F19")),
					menu.TextWithAccelerator("F20", "F20", menu.Accel("F20")),
					menu.TextWithAccelerator("F21", "F21", menu.Accel("F21")),
					menu.TextWithAccelerator("F22", "F22", menu.Accel("F22")),
					menu.TextWithAccelerator("F23", "F23", menu.Accel("F23")),
					menu.TextWithAccelerator("F24", "F24", menu.Accel("F24")),
					menu.TextWithAccelerator("F25", "F25", menu.Accel("F25")),
					menu.TextWithAccelerator("F26", "F26", menu.Accel("F26")),
					menu.TextWithAccelerator("F27", "F27", menu.Accel("F27")),
					menu.TextWithAccelerator("F28", "F28", menu.Accel("F28")),
					menu.TextWithAccelerator("F29", "F29", menu.Accel("F29")),
					menu.TextWithAccelerator("F30", "F30", menu.Accel("F30")),
					menu.TextWithAccelerator("F31", "F31", menu.Accel("F31")),
					menu.TextWithAccelerator("F32", "F32", menu.Accel("F32")),
					menu.TextWithAccelerator("F33", "F33", menu.Accel("F33")),
					menu.TextWithAccelerator("F34", "F34", menu.Accel("F34")),
					menu.TextWithAccelerator("F35", "F35", menu.Accel("F35")),
				}),
				menu.SubMenu("Standard Keys", []*menu.MenuItem{
					menu.TextWithAccelerator("Backtick", "Backtick", menu.Accel("`")),
					menu.TextWithAccelerator("Plus", "Plus", menu.Accel("+")),
				}),
			}),
			&menu.MenuItem{
				Label:       "Disabled Menu",
				Type:        menu.TextType,
				Accelerator: menu.ComboAccel("p", menu.CmdOrCtrl, menu.Shift),
				Disabled:    true,
			},
			&menu.MenuItem{
				Label:  "Hidden Menu",
				Type:   menu.TextType,
				Hidden: true,
			},
			&menu.MenuItem{
				ID:          "checkbox-menu",
				Label:       "Checkbox Menu",
				Type:        menu.CheckboxType,
				Accelerator: menu.CmdOrCtrlAccel("l"),
				Checked:     true,
			},
			menu.Separator(),
			menu.Radio("😀 Option 1", "😀option-1", true),
			menu.Radio("😺 Option 2", "option-2", false),
			menu.Radio("❤️ Option 3", "option-3", false),
		}),
	})

	myMenu.Append(windowMenu)
	return myMenu
}
