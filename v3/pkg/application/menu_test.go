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
		name          string
		menu          *application.Menu
		index         int
		expectedLabel string
		shouldBeNil   bool
	}{
		{
			name: "Get first item",
			menu: application.NewMenuFromItems(
				application.NewMenuItem("Item 1"),
				application.NewMenuItem("Item 2"),
			),
			index:         0,
			expectedLabel: "Item 1",
			shouldBeNil:   false,
		},
		{
			name: "Get last item",
			menu: application.NewMenuFromItems(
				application.NewMenuItem("Item 1"),
				application.NewMenuItem("Item 2"),
			),
			index:         1,
			expectedLabel: "Item 2",
			shouldBeNil:   false,
		},
		{
			name: "Index out of bounds (negative)",
			menu: application.NewMenuFromItems(
				application.NewMenuItem("Item 1"),
				application.NewMenuItem("Item 2"),
			),
			index:       -1,
			shouldBeNil: true,
		},
		{
			name: "Index out of bounds (too large)",
			menu: application.NewMenuFromItems(
				application.NewMenuItem("Item 1"),
				application.NewMenuItem("Item 2"),
			),
			index:       2,
			shouldBeNil: true,
		},
		{
			name:        "Empty menu",
			menu:        application.NewMenu(),
			index:       0,
			shouldBeNil: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			item := test.menu.ItemAt(test.index)

			if test.shouldBeNil {
				if item != nil {
					t.Errorf("Expected nil item, but got item with label: %s", item.Label())
				}
			} else {
				if item == nil {
					t.Errorf("Expected item with label %s, but got nil", test.expectedLabel)
				} else if item.Label() != test.expectedLabel {
					t.Errorf("Expected item with label %s, but got %s", test.expectedLabel, item.Label())
				}
			}
		})
	}
}

func TestMenu_InsertAt(t *testing.T) {
	tests := []struct {
		name           string
		setupMenu      func() *application.Menu
		index          int
		label          string
		expectedLabels []string
	}{
		{
			name: "Insert at beginning",
			setupMenu: func() *application.Menu {
				return application.NewMenuFromItems(
					application.NewMenuItem("Item 1"),
					application.NewMenuItem("Item 2"),
				)
			},
			index:          0,
			label:          "New Item",
			expectedLabels: []string{"New Item", "Item 1", "Item 2"},
		},
		{
			name: "Insert in middle",
			setupMenu: func() *application.Menu {
				return application.NewMenuFromItems(
					application.NewMenuItem("Item 1"),
					application.NewMenuItem("Item 2"),
				)
			},
			index:          1,
			label:          "New Item",
			expectedLabels: []string{"Item 1", "New Item", "Item 2"},
		},
		{
			name: "Insert at end",
			setupMenu: func() *application.Menu {
				return application.NewMenuFromItems(
					application.NewMenuItem("Item 1"),
					application.NewMenuItem("Item 2"),
				)
			},
			index:          2,
			label:          "New Item",
			expectedLabels: []string{"Item 1", "Item 2", "New Item"},
		},
		{
			name: "Insert with negative index (should insert at beginning)",
			setupMenu: func() *application.Menu {
				return application.NewMenuFromItems(
					application.NewMenuItem("Item 1"),
					application.NewMenuItem("Item 2"),
				)
			},
			index:          -1,
			label:          "New Item",
			expectedLabels: []string{"New Item", "Item 1", "Item 2"},
		},
		{
			name: "Insert with index too large (should append)",
			setupMenu: func() *application.Menu {
				return application.NewMenuFromItems(
					application.NewMenuItem("Item 1"),
					application.NewMenuItem("Item 2"),
				)
			},
			index:          10,
			label:          "New Item",
			expectedLabels: []string{"Item 1", "Item 2", "New Item"},
		},
		{
			name: "Insert into empty menu",
			setupMenu: func() *application.Menu {
				return application.NewMenu()
			},
			index:          0,
			label:          "New Item",
			expectedLabels: []string{"New Item"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			menu := test.setupMenu()
			menu.InsertAt(test.index, test.label)

			// Verify the menu has the correct number of items
			if menu.Count() != len(test.expectedLabels) {
				t.Errorf("Expected menu to have %d items, but got %d", len(test.expectedLabels), menu.Count())
			}

			// Verify each item has the expected label
			for i, expectedLabel := range test.expectedLabels {
				item := menu.ItemAt(i)
				if item == nil {
					t.Errorf("Expected item at index %d, but got nil", i)
				} else if item.Label() != expectedLabel {
					t.Errorf("Expected item at index %d to have label %s, but got %s", i, expectedLabel, item.Label())
				}
			}
		})
	}
}

func TestMenu_InsertItemAt(t *testing.T) {
	tests := []struct {
		name           string
		setupMenu      func() *application.Menu
		index          int
		itemLabel      string
		expectedLabels []string
	}{
		{
			name: "Insert existing item at beginning",
			setupMenu: func() *application.Menu {
				return application.NewMenuFromItems(
					application.NewMenuItem("Item 1"),
					application.NewMenuItem("Item 2"),
				)
			},
			index:          0,
			itemLabel:      "Existing Item",
			expectedLabels: []string{"Existing Item", "Item 1", "Item 2"},
		},
		{
			name: "Insert existing item in middle",
			setupMenu: func() *application.Menu {
				return application.NewMenuFromItems(
					application.NewMenuItem("Item 1"),
					application.NewMenuItem("Item 2"),
				)
			},
			index:          1,
			itemLabel:      "Existing Item",
			expectedLabels: []string{"Item 1", "Existing Item", "Item 2"},
		},
		{
			name: "Insert existing item at end",
			setupMenu: func() *application.Menu {
				return application.NewMenuFromItems(
					application.NewMenuItem("Item 1"),
					application.NewMenuItem("Item 2"),
				)
			},
			index:          2,
			itemLabel:      "Existing Item",
			expectedLabels: []string{"Item 1", "Item 2", "Existing Item"},
		},
		{
			name: "Insert with negative index (should insert at beginning)",
			setupMenu: func() *application.Menu {
				return application.NewMenuFromItems(
					application.NewMenuItem("Item 1"),
					application.NewMenuItem("Item 2"),
				)
			},
			index:          -1,
			itemLabel:      "Existing Item",
			expectedLabels: []string{"Existing Item", "Item 1", "Item 2"},
		},
		{
			name: "Insert with index too large (should append)",
			setupMenu: func() *application.Menu {
				return application.NewMenuFromItems(
					application.NewMenuItem("Item 1"),
					application.NewMenuItem("Item 2"),
				)
			},
			index:          10,
			itemLabel:      "Existing Item",
			expectedLabels: []string{"Item 1", "Item 2", "Existing Item"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			menu := test.setupMenu()
			existingItem := application.NewMenuItem(test.itemLabel)
			menu.InsertItemAt(test.index, existingItem)

			// Verify the menu has the correct number of items
			if menu.Count() != len(test.expectedLabels) {
				t.Errorf("Expected menu to have %d items, but got %d", len(test.expectedLabels), menu.Count())
			}

			// Verify each item has the expected label
			for i, expectedLabel := range test.expectedLabels {
				item := menu.ItemAt(i)
				if item == nil {
					t.Errorf("Expected item at index %d, but got nil", i)
				} else if item.Label() != expectedLabel {
					t.Errorf("Expected item at index %d to have label %s, but got %s", i, expectedLabel, item.Label())
				}
			}
		})
	}
}

func TestMenu_SpecializedInsertFunctions(t *testing.T) {
	t.Run("InsertSeparatorAt", func(t *testing.T) {
		menu := application.NewMenuFromItems(
			application.NewMenuItem("Item 1"),
			application.NewMenuItem("Item 2"),
		)

		menu.InsertSeparatorAt(1)

		// Verify the separator was inserted
		if menu.Count() != 3 {
			t.Errorf("Expected menu to have 3 items, but got %d", menu.Count())
		}

		separator := menu.ItemAt(1)
		if separator == nil {
			t.Errorf("Expected separator at index 1, but got nil")
		} else if !separator.IsSeparator() {
			t.Errorf("Expected item at index 1 to be a separator, but it wasn't")
		}
	})

	t.Run("InsertCheckboxAt", func(t *testing.T) {
		menu := application.NewMenuFromItems(
			application.NewMenuItem("Item 1"),
			application.NewMenuItem("Item 2"),
		)

		menu.InsertCheckboxAt(1, "Checkbox", true)

		// Verify the checkbox was inserted
		if menu.Count() != 3 {
			t.Errorf("Expected menu to have 3 items, but got %d", menu.Count())
		}

		checkbox := menu.ItemAt(1)
		if checkbox == nil {
			t.Errorf("Expected checkbox at index 1, but got nil")
		} else if !checkbox.IsCheckbox() {
			t.Errorf("Expected item at index 1 to be a checkbox, but it wasn't")
		} else if checkbox.Label() != "Checkbox" {
			t.Errorf("Expected checkbox to have label 'Checkbox', but got '%s'", checkbox.Label())
		} else if !checkbox.Checked() {
			t.Errorf("Expected checkbox to be checked, but it wasn't")
		}
	})

	t.Run("InsertRadioAt", func(t *testing.T) {
		menu := application.NewMenuFromItems(
			application.NewMenuItem("Item 1"),
			application.NewMenuItem("Item 2"),
		)

		menu.InsertRadioAt(1, "Radio", true)

		// Verify the radio button was inserted
		if menu.Count() != 3 {
			t.Errorf("Expected menu to have 3 items, but got %d", menu.Count())
		}

		radio := menu.ItemAt(1)
		if radio == nil {
			t.Errorf("Expected radio button at index 1, but got nil")
		} else if !radio.IsRadio() {
			t.Errorf("Expected item at index 1 to be a radio button, but it wasn't")
		} else if radio.Label() != "Radio" {
			t.Errorf("Expected radio button to have label 'Radio', but got '%s'", radio.Label())
		} else if !radio.Checked() {
			t.Errorf("Expected radio button to be checked, but it wasn't")
		}
	})

	t.Run("InsertSubmenuAt", func(t *testing.T) {
		menu := application.NewMenuFromItems(
			application.NewMenuItem("Item 1"),
			application.NewMenuItem("Item 2"),
		)

		submenu := menu.InsertSubmenuAt(1, "Submenu")
		submenu.Add("Submenu Item")

		// Verify the submenu was inserted
		if menu.Count() != 3 {
			t.Errorf("Expected menu to have 3 items, but got %d", menu.Count())
		}

		submenuItem := menu.ItemAt(1)
		if submenuItem == nil {
			t.Errorf("Expected submenu at index 1, but got nil")
		} else if !submenuItem.IsSubmenu() {
			t.Errorf("Expected item at index 1 to be a submenu, but it wasn't")
		} else if submenuItem.Label() != "Submenu" {
			t.Errorf("Expected submenu to have label 'Submenu', but got '%s'", submenuItem.Label())
		}

		// Verify the submenu has the expected item
		submenuFromItem := submenuItem.GetSubmenu()
		if submenuFromItem == nil {
			t.Errorf("Expected to get submenu from item, but got nil")
		} else if submenuFromItem.Count() != 1 {
			t.Errorf("Expected submenu to have 1 item, but got %d", submenuFromItem.Count())
		} else {
			submenuItemFromSubmenu := submenuFromItem.ItemAt(0)
			if submenuItemFromSubmenu == nil {
				t.Errorf("Expected item in submenu, but got nil")
			} else if submenuItemFromSubmenu.Label() != "Submenu Item" {
				t.Errorf("Expected submenu item to have label 'Submenu Item', but got '%s'", submenuItemFromSubmenu.Label())
			}
		}
	})
}

func TestMenu_Count(t *testing.T) {
	tests := []struct {
		name          string
		menu          *application.Menu
		expectedCount int
	}{
		{
			name: "Menu with items",
			menu: application.NewMenuFromItems(
				application.NewMenuItem("Item 1"),
				application.NewMenuItem("Item 2"),
				application.NewMenuItem("Item 3"),
			),
			expectedCount: 3,
		},
		{
			name:          "Empty menu",
			menu:          application.NewMenu(),
			expectedCount: 0,
		},
		{
			name: "Menu with separator",
			menu: application.NewMenuFromItems(
				application.NewMenuItem("Item 1"),
				application.NewMenuItemSeparator(),
				application.NewMenuItem("Item 2"),
			),
			expectedCount: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			count := test.menu.Count()
			if count != test.expectedCount {
				t.Errorf("Expected count to be %d, but got %d", test.expectedCount, count)
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
