package application

import (
	"testing"
)

func TestNewMenu(t *testing.T) {
	menu := NewMenu()
	if menu == nil {
		t.Fatal("NewMenu returned nil")
	}
	if menu.items != nil {
		t.Error("items should be nil initially")
	}
}

func TestMenu_Add(t *testing.T) {
	menu := NewMenu()

	item := menu.Add("Test Item")
	if item == nil {
		t.Fatal("Add returned nil")
	}
	if item.label != "Test Item" {
		t.Errorf("label = %q, want %q", item.label, "Test Item")
	}
	if len(menu.items) != 1 {
		t.Errorf("items count = %d, want 1", len(menu.items))
	}

	// Clean up
	menu.Destroy()
}

func TestMenu_AddSeparator(t *testing.T) {
	menu := NewMenu()

	menu.AddSeparator()
	if len(menu.items) != 1 {
		t.Errorf("items count = %d, want 1", len(menu.items))
	}
	if menu.items[0].itemType != separator {
		t.Error("item should be a separator")
	}

	// Clean up
	menu.Destroy()
}

func TestMenu_AddCheckbox(t *testing.T) {
	menu := NewMenu()

	item := menu.AddCheckbox("Check Me", true)
	if item == nil {
		t.Fatal("AddCheckbox returned nil")
	}
	if item.label != "Check Me" {
		t.Errorf("label = %q, want %q", item.label, "Check Me")
	}
	if item.itemType != checkbox {
		t.Error("item should be a checkbox")
	}
	if !item.checked {
		t.Error("checkbox should be checked")
	}

	// Clean up
	menu.Destroy()
}

func TestMenu_AddRadio(t *testing.T) {
	menu := NewMenu()

	item := menu.AddRadio("Radio Option", false)
	if item == nil {
		t.Fatal("AddRadio returned nil")
	}
	if item.label != "Radio Option" {
		t.Errorf("label = %q, want %q", item.label, "Radio Option")
	}
	if item.itemType != radio {
		t.Error("item should be a radio")
	}
	if item.checked {
		t.Error("radio should not be checked")
	}

	// Clean up
	menu.Destroy()
}

func TestMenu_AddSubmenu(t *testing.T) {
	menu := NewMenu()

	subMenu := menu.AddSubmenu("Submenu")
	if subMenu == nil {
		t.Fatal("AddSubmenu returned nil")
	}
	if len(menu.items) != 1 {
		t.Errorf("items count = %d, want 1", len(menu.items))
	}
	if menu.items[0].itemType != submenu {
		t.Error("item should be a submenu")
	}

	// Clean up
	menu.Destroy()
}

func TestMenu_SetLabel(t *testing.T) {
	menu := NewMenu()

	menu.SetLabel("My Menu")
	if menu.label != "My Menu" {
		t.Errorf("label = %q, want %q", menu.label, "My Menu")
	}
}

func TestMenu_ItemAt(t *testing.T) {
	menu := NewMenu()
	menu.Add("Item 0")
	menu.Add("Item 1")
	menu.Add("Item 2")

	item := menu.ItemAt(1)
	if item == nil {
		t.Fatal("ItemAt returned nil")
	}
	if item.label != "Item 1" {
		t.Errorf("label = %q, want %q", item.label, "Item 1")
	}

	// Clean up
	menu.Destroy()
}

func TestMenu_ItemAt_OutOfBounds(t *testing.T) {
	menu := NewMenu()
	menu.Add("Item")

	if menu.ItemAt(-1) != nil {
		t.Error("ItemAt(-1) should return nil")
	}
	if menu.ItemAt(1) != nil {
		t.Error("ItemAt(1) should return nil for single item menu")
	}
	if menu.ItemAt(100) != nil {
		t.Error("ItemAt(100) should return nil")
	}

	// Clean up
	menu.Destroy()
}

func TestMenu_FindByLabel(t *testing.T) {
	menu := NewMenu()
	menu.Add("First")
	menu.Add("Second")
	menu.Add("Third")

	found := menu.FindByLabel("Second")
	if found == nil {
		t.Fatal("FindByLabel returned nil")
	}
	if found.label != "Second" {
		t.Errorf("label = %q, want %q", found.label, "Second")
	}

	// Clean up
	menu.Destroy()
}

func TestMenu_FindByLabel_NotFound(t *testing.T) {
	menu := NewMenu()
	menu.Add("First")

	found := menu.FindByLabel("NonExistent")
	if found != nil {
		t.Error("FindByLabel should return nil for non-existent label")
	}

	// Clean up
	menu.Destroy()
}

func TestMenu_FindByLabel_InSubmenu(t *testing.T) {
	menu := NewMenu()
	submenu := menu.AddSubmenu("Submenu")
	submenu.Add("Nested Item")

	found := menu.FindByLabel("Nested Item")
	if found == nil {
		t.Fatal("FindByLabel should find item in submenu")
	}
	if found.label != "Nested Item" {
		t.Errorf("label = %q, want %q", found.label, "Nested Item")
	}

	// Clean up
	menu.Destroy()
}

func TestMenu_RemoveMenuItem(t *testing.T) {
	menu := NewMenu()
	item1 := menu.Add("First")
	item2 := menu.Add("Second")
	menu.Add("Third")

	menu.RemoveMenuItem(item2)
	if len(menu.items) != 2 {
		t.Errorf("items count = %d, want 2", len(menu.items))
	}
	if menu.FindByLabel("Second") != nil {
		t.Error("Second item should be removed")
	}

	// Clean up
	removeMenuItemByID(item1.id)
}

func TestMenu_Clear(t *testing.T) {
	menu := NewMenu()
	menu.Add("First")
	menu.Add("Second")
	menu.Add("Third")

	menu.Clear()
	if menu.items != nil {
		t.Error("items should be nil after Clear")
	}
}

func TestMenu_Append(t *testing.T) {
	menu1 := NewMenu()
	menu1.Add("Item 1")

	menu2 := NewMenu()
	menu2.Add("Item 2")
	menu2.Add("Item 3")

	menu1.Append(menu2)
	if len(menu1.items) != 3 {
		t.Errorf("items count = %d, want 3", len(menu1.items))
	}

	// Clean up
	menu1.Destroy()
}

func TestMenu_Append_Nil(t *testing.T) {
	menu := NewMenu()
	menu.Add("Item")

	menu.Append(nil)
	if len(menu.items) != 1 {
		t.Error("Append(nil) should not change menu")
	}

	// Clean up
	menu.Destroy()
}

func TestMenu_Prepend(t *testing.T) {
	menu1 := NewMenu()
	menu1.Add("Item 3")

	menu2 := NewMenu()
	menu2.Add("Item 1")
	menu2.Add("Item 2")

	menu1.Prepend(menu2)
	if len(menu1.items) != 3 {
		t.Errorf("items count = %d, want 3", len(menu1.items))
	}
	if menu1.items[0].label != "Item 1" {
		t.Errorf("First item should be 'Item 1', got %q", menu1.items[0].label)
	}

	// Clean up
	menu1.Destroy()
}

func TestMenu_Clone(t *testing.T) {
	menu := NewMenu()
	menu.SetLabel("Original Menu")
	menu.Add("Item 1")
	menu.Add("Item 2")

	clone := menu.Clone()
	if clone == menu {
		t.Error("Clone should return different pointer")
	}
	if clone.label != menu.label {
		t.Error("Clone should have same label")
	}
	if len(clone.items) != len(menu.items) {
		t.Error("Clone should have same number of items")
	}

	// Clean up
	menu.Destroy()
	clone.Destroy()
}

func TestNewMenuFromItems(t *testing.T) {
	item1 := NewMenuItem("Item 1")
	item2 := NewMenuItem("Item 2")
	item3 := NewMenuItem("Item 3")

	menu := NewMenuFromItems(item1, item2, item3)
	if menu == nil {
		t.Fatal("NewMenuFromItems returned nil")
	}
	if len(menu.items) != 3 {
		t.Errorf("items count = %d, want 3", len(menu.items))
	}

	// Clean up
	menu.Destroy()
}

func TestNewSubmenu(t *testing.T) {
	items := NewMenu()
	items.Add("Sub Item 1")
	items.Add("Sub Item 2")

	item := NewSubmenu("My Submenu", items)
	if item == nil {
		t.Fatal("NewSubmenu returned nil")
	}
	if item.label != "My Submenu" {
		t.Errorf("label = %q, want %q", item.label, "My Submenu")
	}
	if item.submenu != items {
		t.Error("submenu should be the provided menu")
	}

	// Clean up
	removeMenuItemByID(item.id)
}

func TestMenu_ProcessRadioGroups(t *testing.T) {
	menu := NewMenu()

	// Add some non-radio items
	menu.Add("Regular Item")

	// Add a group of radio items
	radio1 := menu.AddRadio("Radio 1", true)
	radio2 := menu.AddRadio("Radio 2", false)
	radio3 := menu.AddRadio("Radio 3", false)

	// Add separator to end the group
	menu.AddSeparator()

	// Add another group
	radio4 := menu.AddRadio("Radio A", true)
	radio5 := menu.AddRadio("Radio B", false)

	// Process radio groups
	menu.processRadioGroups()

	// First group should be linked
	if len(radio1.radioGroupMembers) != 3 {
		t.Errorf("First group should have 3 members, got %d", len(radio1.radioGroupMembers))
	}
	if len(radio2.radioGroupMembers) != 3 {
		t.Errorf("radio2 should have 3 members, got %d", len(radio2.radioGroupMembers))
	}
	if len(radio3.radioGroupMembers) != 3 {
		t.Errorf("radio3 should have 3 members, got %d", len(radio3.radioGroupMembers))
	}

	// Second group should be linked
	if len(radio4.radioGroupMembers) != 2 {
		t.Errorf("Second group should have 2 members, got %d", len(radio4.radioGroupMembers))
	}
	if len(radio5.radioGroupMembers) != 2 {
		t.Errorf("radio5 should have 2 members, got %d", len(radio5.radioGroupMembers))
	}

	// Clean up
	menu.Destroy()
}

func TestMenu_SetContextData(t *testing.T) {
	menu := NewMenu()
	menu.Add("Item 1")
	menu.Add("Item 2")

	data := &ContextMenuData{Data: "test-data"}
	menu.setContextData(data)

	// Verify data was set on items
	for _, item := range menu.items {
		if item.contextMenuData != data {
			t.Error("Context data should be set on all items")
		}
	}

	// Clean up
	menu.Destroy()
}
