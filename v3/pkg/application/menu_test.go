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
