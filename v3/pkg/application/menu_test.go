package application_test

import (
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestMenu_FindByLabel(t *testing.T) {
	tests := []struct {
		name        string
		menu        *application.Menu
		label       string
		shouldError bool
	}{
		{
			name: "Find top-level item",
			menu: application.NewMenuFromItems(
				application.NewMenuItem("Target"),
			),
			label:       "Target",
			shouldError: false,
		},
		{
			name: "Find item in submenu",
			menu: application.NewMenuFromItems(
				application.NewMenuItem("Item 1"),
				application.NewSubmenu("Submenu", application.NewMenuFromItems(
					application.NewMenuItem("Subitem 1"),
					application.NewMenuItem("Target"),
				)),
			),
			label:       "Target",
			shouldError: false,
		},
		{
			name: "Not find item",
			menu: application.NewMenuFromItems(
				application.NewMenuItem("Item 1"),
				application.NewSubmenu("Submenu", application.NewMenuFromItems(
					application.NewMenuItem("Subitem 1"),
					application.NewMenuItem("Target"),
				)),
			),
			label:       "Random",
			shouldError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			found := test.menu.FindByLabel(test.label)
			if test.shouldError && found != nil {
				t.Errorf("Expected error, but found %v", found)
			}
			if !test.shouldError && found == nil {
				t.Errorf("Expected item, but found none")
			}
		})
	}
}

func TestMenu_ItemAt(t *testing.T) {
	tests := []struct {
		name        string
		menu        *application.Menu
		index       int
		shouldError bool
	}{
		{
			name: "Valid index",
			menu: application.NewMenuFromItems(
				application.NewMenuItem("Item 1"),
				application.NewMenuItem("Item 2"),
				application.NewMenuItem("Target"),
			),
			index:       2,
			shouldError: false,
		},
		{
			name: "Index out of bounds (negative)",
			menu: application.NewMenuFromItems(
				application.NewMenuItem("Item 1"),
				application.NewMenuItem("Item 2"),
			),
			index:       -1,
			shouldError: true,
		},
		{
			name: "Index out of bounds (too large)",
			menu: application.NewMenuFromItems(
				application.NewMenuItem("Item 1"),
				application.NewMenuItem("Item 2"),
			),
			index:       2,
			shouldError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			item := test.menu.ItemAt(test.index)
			if test.shouldError && item != nil {
				t.Errorf("Expected error, but found %v", item)
			}
			if !test.shouldError && item == nil {
				t.Errorf("Expected item, but found none")
			}
		})
	}
}

func TestMenu_RemoveMenuItem(t *testing.T) {
	itemToRemove := application.NewMenuItem("Target")
	itemToKeep := application.NewMenuItem("Item 1")

	tests := []struct {
		name       string
		menu       *application.Menu
		item       *application.MenuItem
		shouldFind bool
	}{
		{
			name:       "Remove existing item",
			menu:       application.NewMenuFromItems(itemToKeep, itemToRemove),
			item:       itemToRemove,
			shouldFind: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.menu.RemoveMenuItem(test.item)
			found := test.menu.FindByLabel(test.item.Label())
			if !test.shouldFind && found != nil {
				t.Errorf("Expected item to be removed, but found %v", found)
			}
		})
	}
}

func TestMenu_InsertAt(t *testing.T) {
	tests := []struct {
		name          string
		menu          *application.Menu
		index         int
		label         string
		expectedIndex int
	}{
		{
			name:          "Insert at beginning",
			menu:          application.NewMenuFromItems(application.NewMenuItem("Item 1"), application.NewMenuItem("Item 2")),
			index:         0,
			label:         "New Item",
			expectedIndex: 0,
		},
		{
			name:          "Insert in middle",
			menu:          application.NewMenuFromItems(application.NewMenuItem("Item 1"), application.NewMenuItem("Item 2")),
			index:         1,
			label:         "New Item",
			expectedIndex: 1,
		},
		{
			name:          "Insert at end",
			menu:          application.NewMenuFromItems(application.NewMenuItem("Item 1"), application.NewMenuItem("Item 2")),
			index:         2,
			label:         "New Item",
			expectedIndex: 2,
		},
		{
			name:          "Insert at negative index (should insert at beginning)",
			menu:          application.NewMenuFromItems(application.NewMenuItem("Item 1"), application.NewMenuItem("Item 2")),
			index:         -1,
			label:         "New Item",
			expectedIndex: 0,
		},
		{
			name:          "Insert beyond end (should insert at end)",
			menu:          application.NewMenuFromItems(application.NewMenuItem("Item 1"), application.NewMenuItem("Item 2")),
			index:         10,
			label:         "New Item",
			expectedIndex: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			item := test.menu.InsertAt(test.index, test.label)
			if item == nil {
				t.Errorf("Expected item to be created, but got nil")
				return
			}
			if item.Label() != test.label {
				t.Errorf("Expected label %s, but got %s", test.label, item.Label())
			}

			// Check if the item is at the expected index
			itemAtIndex := test.menu.ItemAt(test.expectedIndex)
			if itemAtIndex != item {
				t.Errorf("Expected item to be at index %d, but it's not", test.expectedIndex)
			}
		})
	}
}

func TestMenu_InsertItemAt(t *testing.T) {
	tests := []struct {
		name          string
		menu          *application.Menu
		index         int
		item          *application.MenuItem
		expectedIndex int
	}{
		{
			name:          "Insert at beginning",
			menu:          application.NewMenuFromItems(application.NewMenuItem("Item 1"), application.NewMenuItem("Item 2")),
			index:         0,
			item:          application.NewMenuItem("New Item"),
			expectedIndex: 0,
		},
		{
			name:          "Insert in middle",
			menu:          application.NewMenuFromItems(application.NewMenuItem("Item 1"), application.NewMenuItem("Item 2")),
			index:         1,
			item:          application.NewMenuItem("New Item"),
			expectedIndex: 1,
		},
		{
			name:          "Insert at end",
			menu:          application.NewMenuFromItems(application.NewMenuItem("Item 1"), application.NewMenuItem("Item 2")),
			index:         2,
			item:          application.NewMenuItem("New Item"),
			expectedIndex: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			menu := test.menu.InsertItemAt(test.index, test.item)
			if menu == nil {
				t.Errorf("Expected menu to be returned, but got nil")
				return
			}

			// Check if the item is at the expected index
			itemAtIndex := test.menu.ItemAt(test.expectedIndex)
			if itemAtIndex != test.item {
				t.Errorf("Expected item to be at index %d, but it's not", test.expectedIndex)
			}
		})
	}
}

func TestMenu_InsertSeparatorAt(t *testing.T) {
	tests := []struct {
		name          string
		menu          *application.Menu
		index         int
		expectedIndex int
	}{
		{
			name:          "Insert at beginning",
			menu:          application.NewMenuFromItems(application.NewMenuItem("Item 1"), application.NewMenuItem("Item 2")),
			index:         0,
			expectedIndex: 0,
		},
		{
			name:          "Insert in middle",
			menu:          application.NewMenuFromItems(application.NewMenuItem("Item 1"), application.NewMenuItem("Item 2")),
			index:         1,
			expectedIndex: 1,
		},
		{
			name:          "Insert at end",
			menu:          application.NewMenuFromItems(application.NewMenuItem("Item 1"), application.NewMenuItem("Item 2")),
			index:         2,
			expectedIndex: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			menu := test.menu.InsertSeparatorAt(test.index)
			if menu == nil {
				t.Errorf("Expected menu to be returned, but got nil")
				return
			}

			// Check if a separator is at the expected index
			itemAtIndex := test.menu.ItemAt(test.expectedIndex)
			if itemAtIndex == nil || !itemAtIndex.IsSeparator() {
				t.Errorf("Expected separator to be at index %d, but it's not", test.expectedIndex)
			}
		})
	}
}

func TestMenu_InsertCheckboxAt(t *testing.T) {
	tests := []struct {
		name          string
		menu          *application.Menu
		index         int
		label         string
		checked       bool
		expectedIndex int
	}{
		{
			name:          "Insert checked checkbox",
			menu:          application.NewMenuFromItems(application.NewMenuItem("Item 1"), application.NewMenuItem("Item 2")),
			index:         1,
			label:         "Checkbox",
			checked:       true,
			expectedIndex: 1,
		},
		{
			name:          "Insert unchecked checkbox",
			menu:          application.NewMenuFromItems(application.NewMenuItem("Item 1"), application.NewMenuItem("Item 2")),
			index:         1,
			label:         "Checkbox",
			checked:       false,
			expectedIndex: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			item := test.menu.InsertCheckboxAt(test.index, test.label, test.checked)
			if item == nil {
				t.Errorf("Expected item to be created, but got nil")
				return
			}
			if item.Label() != test.label {
				t.Errorf("Expected label %s, but got %s", test.label, item.Label())
			}
			if item.Checked() != test.checked {
				t.Errorf("Expected checked to be %v, but got %v", test.checked, item.Checked())
			}
			if !item.IsCheckbox() {
				t.Errorf("Expected item to be a checkbox, but it's not")
			}

			// Check if the item is at the expected index
			itemAtIndex := test.menu.ItemAt(test.expectedIndex)
			if itemAtIndex != item {
				t.Errorf("Expected item to be at index %d, but it's not", test.expectedIndex)
			}
		})
	}
}

func TestMenu_InsertRadioAt(t *testing.T) {
	tests := []struct {
		name          string
		menu          *application.Menu
		index         int
		label         string
		checked       bool
		expectedIndex int
	}{
		{
			name:          "Insert checked radio",
			menu:          application.NewMenuFromItems(application.NewMenuItem("Item 1"), application.NewMenuItem("Item 2")),
			index:         1,
			label:         "Radio",
			checked:       true,
			expectedIndex: 1,
		},
		{
			name:          "Insert unchecked radio",
			menu:          application.NewMenuFromItems(application.NewMenuItem("Item 1"), application.NewMenuItem("Item 2")),
			index:         1,
			label:         "Radio",
			checked:       false,
			expectedIndex: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			item := test.menu.InsertRadioAt(test.index, test.label, test.checked)
			if item == nil {
				t.Errorf("Expected item to be created, but got nil")
				return
			}
			if item.Label() != test.label {
				t.Errorf("Expected label %s, but got %s", test.label, item.Label())
			}
			if item.Checked() != test.checked {
				t.Errorf("Expected checked to be %v, but got %v", test.checked, item.Checked())
			}
			if !item.IsRadio() {
				t.Errorf("Expected item to be a radio, but it's not")
			}

			// Check if the item is at the expected index
			itemAtIndex := test.menu.ItemAt(test.expectedIndex)
			if itemAtIndex != item {
				t.Errorf("Expected item to be at index %d, but it's not", test.expectedIndex)
			}
		})
	}
}

func TestMenu_InsertSubmenuAt(t *testing.T) {
	tests := []struct {
		name          string
		menu          *application.Menu
		index         int
		label         string
		expectedIndex int
	}{
		{
			name:          "Insert submenu",
			menu:          application.NewMenuFromItems(application.NewMenuItem("Item 1"), application.NewMenuItem("Item 2")),
			index:         1,
			label:         "Submenu",
			expectedIndex: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			submenu := test.menu.InsertSubmenuAt(test.index, test.label)
			if submenu == nil {
				t.Errorf("Expected submenu to be created, but got nil")
				return
			}

			// Check if a submenu item is at the expected index
			itemAtIndex := test.menu.ItemAt(test.expectedIndex)
			if itemAtIndex == nil || !itemAtIndex.IsSubmenu() {
				t.Errorf("Expected submenu to be at index %d, but it's not", test.expectedIndex)
			}

			// Check if the submenu's label is correct
			if itemAtIndex.Label() != test.label {
				t.Errorf("Expected label %s, but got %s", test.label, itemAtIndex.Label())
			}

			// Check if the submenu is accessible through GetSubmenu
			if itemAtIndex.GetSubmenu() != submenu {
				t.Errorf("Expected submenu to be accessible through GetSubmenu, but it's not")
			}
		})
	}
}
