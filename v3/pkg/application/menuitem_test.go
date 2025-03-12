package application_test

import (
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestMenuItem_GetAccelerator(t *testing.T) {
	tests := []struct {
		name        string
		menuItem    *application.MenuItem
		expectedAcc string
	}{
		{
			name:        "Get existing accelerator",
			menuItem:    application.NewMenuItem("Item 1").SetAccelerator("ctrl+a"),
			expectedAcc: "Ctrl+A",
		},
		{
			name:        "Get non-existing accelerator",
			menuItem:    application.NewMenuItem("Item 2"),
			expectedAcc: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			acc := test.menuItem.GetAccelerator()
			if acc != test.expectedAcc {
				t.Errorf("Expected accelerator to be %v, but got %v", test.expectedAcc, acc)
			}
		})
	}
}

func TestMenuItem_RemoveAccelerator(t *testing.T) {
	tests := []struct {
		name     string
		menuItem *application.MenuItem
	}{
		{
			name:     "Remove existing accelerator",
			menuItem: application.NewMenuItem("Item 1").SetAccelerator("Ctrl+A"),
		},
		{
			name:     "Remove non-existing accelerator",
			menuItem: application.NewMenuItem("Item 2"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.menuItem.RemoveAccelerator()
			acc := test.menuItem.GetAccelerator()
			if acc != "" {
				t.Errorf("Expected accelerator to be removed, but got %v", acc)
			}
		})
	}
}

func TestMenuItem_SetHidden(t *testing.T) {
	tests := []struct {
		name          string
		setupMenuItem func() *application.MenuItem
		setHidden     bool
		expectedState bool
	}{
		{
			name: "Hide regular menu item",
			setupMenuItem: func() *application.MenuItem {
				return application.NewMenuItem("Regular Item")
			},
			setHidden:     true,
			expectedState: true,
		},
		{
			name: "Show regular menu item",
			setupMenuItem: func() *application.MenuItem {
				return application.NewMenuItem("Regular Item").SetHidden(true)
			},
			setHidden:     false,
			expectedState: false,
		},
		{
			name: "Hide checkbox menu item",
			setupMenuItem: func() *application.MenuItem {
				return application.NewMenuItemCheckbox("Checkbox Item", true)
			},
			setHidden:     true,
			expectedState: true,
		},
		{
			name: "Hide radio menu item",
			setupMenuItem: func() *application.MenuItem {
				return application.NewMenuItemRadio("Radio Item", true)
			},
			setHidden:     true,
			expectedState: true,
		},
		{
			name: "Hide separator",
			setupMenuItem: func() *application.MenuItem {
				return application.NewMenuItemSeparator()
			},
			setHidden:     true,
			expectedState: true,
		},
		{
			name: "Hide submenu",
			setupMenuItem: func() *application.MenuItem {
				submenu := application.NewSubmenu("Submenu", application.NewMenu())
				return submenu
			},
			setHidden:     true,
			expectedState: true,
		},
		{
			name: "Show submenu",
			setupMenuItem: func() *application.MenuItem {
				submenu := application.NewSubmenu("Submenu", application.NewMenu())
				submenu.SetHidden(true)
				return submenu
			},
			setHidden:     false,
			expectedState: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			menuItem := test.setupMenuItem()

			// Set the hidden state
			menuItem.SetHidden(test.setHidden)

			// Verify the hidden state
			if menuItem.Hidden() != test.expectedState {
				t.Errorf("Expected hidden state to be %v, but got %v", test.expectedState, menuItem.Hidden())
			}
		})
	}
}

func TestMenuItem_SubmenuVisibility(t *testing.T) {
	t.Run("Submenu visibility through SetHidden", func(t *testing.T) {
		// Create a menu with a submenu
		menu := application.NewMenu()

		// Create a submenu using NewSubMenuItem directly to get the MenuItem
		submenuItem := application.NewSubMenuItem("Submenu")
		menu.InsertItemAt(0, submenuItem)

		// Get the submenu from the menu item
		submenu := submenuItem.GetSubmenu()

		// Add some items to the submenu
		submenu.Add("Submenu Item 1")
		submenu.Add("Submenu Item 2")

		// Initially, the submenu should be visible
		if submenuItem.Hidden() {
			t.Errorf("Expected submenu to be visible initially, but it was hidden")
		}

		// Hide the submenu
		submenuItem.SetHidden(true)
		if !submenuItem.Hidden() {
			t.Errorf("Expected submenu to be hidden after SetHidden(true), but it was visible")
		}

		// Show the submenu again
		submenuItem.SetHidden(false)
		if submenuItem.Hidden() {
			t.Errorf("Expected submenu to be visible after SetHidden(false), but it was hidden")
		}
	})
}
