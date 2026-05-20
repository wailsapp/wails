package application

import (
	"testing"
)

func TestMenuItemType_Constants(t *testing.T) {
	// Verify menu item type constants are distinct
	if text != 0 {
		t.Error("text should be 0")
	}
	if separator != 1 {
		t.Error("separator should be 1")
	}
	if checkbox != 2 {
		t.Error("checkbox should be 2")
	}
	if radio != 3 {
		t.Error("radio should be 3")
	}
	if submenu != 4 {
		t.Error("submenu should be 4")
	}
}

func TestNewMenuItem(t *testing.T) {
	item := NewMenuItem("Test Label")

	if item == nil {
		t.Fatal("NewMenuItem returned nil")
	}
	if item.label != "Test Label" {
		t.Errorf("label = %q, want %q", item.label, "Test Label")
	}
	if item.itemType != text {
		t.Error("itemType should be text")
	}
	if item.id == 0 {
		t.Error("id should be non-zero")
	}

	// Clean up
	removeMenuItemByID(item.id)
}

func TestNewMenuItemSeparator(t *testing.T) {
	item := NewMenuItemSeparator()

	if item == nil {
		t.Fatal("NewMenuItemSeparator returned nil")
	}
	if item.itemType != separator {
		t.Error("itemType should be separator")
	}
}

func TestNewMenuItemCheckbox(t *testing.T) {
	item := NewMenuItemCheckbox("Checkbox", true)

	if item == nil {
		t.Fatal("NewMenuItemCheckbox returned nil")
	}
	if item.label != "Checkbox" {
		t.Errorf("label = %q, want %q", item.label, "Checkbox")
	}
	if item.itemType != checkbox {
		t.Error("itemType should be checkbox")
	}
	if !item.checked {
		t.Error("checked should be true")
	}

	// Clean up
	removeMenuItemByID(item.id)
}

func TestNewMenuItemCheckbox_Unchecked(t *testing.T) {
	item := NewMenuItemCheckbox("Unchecked", false)

	if item.checked {
		t.Error("checked should be false")
	}

	// Clean up
	removeMenuItemByID(item.id)
}

func TestNewMenuItemRadio(t *testing.T) {
	item := NewMenuItemRadio("Radio", true)

	if item == nil {
		t.Fatal("NewMenuItemRadio returned nil")
	}
	if item.label != "Radio" {
		t.Errorf("label = %q, want %q", item.label, "Radio")
	}
	if item.itemType != radio {
		t.Error("itemType should be radio")
	}
	if !item.checked {
		t.Error("checked should be true")
	}

	// Clean up
	removeMenuItemByID(item.id)
}

func TestNewSubMenuItem(t *testing.T) {
	item := NewSubMenuItem("Submenu")

	if item == nil {
		t.Fatal("NewSubMenuItem returned nil")
	}
	if item.label != "Submenu" {
		t.Errorf("label = %q, want %q", item.label, "Submenu")
	}
	if item.itemType != submenu {
		t.Error("itemType should be submenu")
	}
	if item.submenu == nil {
		t.Error("submenu should not be nil")
	}
	if item.submenu.label != "Submenu" {
		t.Errorf("submenu.label = %q, want %q", item.submenu.label, "Submenu")
	}

	// Clean up
	removeMenuItemByID(item.id)
}

func TestMenuItemMap_AddGetRemove(t *testing.T) {
	item := NewMenuItem("Test")
	id := item.id

	// Item should be in map
	retrieved := getMenuItemByID(id)
	if retrieved != item {
		t.Error("getMenuItemByID should return the same item")
	}

	// Remove item
	removeMenuItemByID(id)

	// Item should be gone
	retrieved = getMenuItemByID(id)
	if retrieved != nil {
		t.Error("getMenuItemByID should return nil after removal")
	}
}

func TestGetMenuItemByID_NotFound(t *testing.T) {
	result := getMenuItemByID(999999)
	if result != nil {
		t.Error("getMenuItemByID should return nil for non-existent ID")
	}
}

func TestMenuItem_UniqueIDs(t *testing.T) {
	item1 := NewMenuItem("Item 1")
	item2 := NewMenuItem("Item 2")
	item3 := NewMenuItem("Item 3")

	if item1.id == item2.id || item2.id == item3.id || item1.id == item3.id {
		t.Error("Menu items should have unique IDs")
	}

	// Clean up
	removeMenuItemByID(item1.id)
	removeMenuItemByID(item2.id)
	removeMenuItemByID(item3.id)
}

func TestMenuItem_Label(t *testing.T) {
	item := NewMenuItem("Original")

	if item.Label() != "Original" {
		t.Errorf("Label() = %q, want %q", item.Label(), "Original")
	}

	item.SetLabel("Updated")
	if item.label != "Updated" {
		t.Errorf("label = %q, want %q", item.label, "Updated")
	}

	// Clean up
	removeMenuItemByID(item.id)
}

func TestMenuItem_Enabled(t *testing.T) {
	item := NewMenuItem("Test")

	if item.disabled {
		t.Error("disabled should default to false")
	}

	item.SetEnabled(false)
	if !item.disabled {
		t.Error("disabled should be true after SetEnabled(false)")
	}

	item.SetEnabled(true)
	if item.disabled {
		t.Error("disabled should be false after SetEnabled(true)")
	}

	// Clean up
	removeMenuItemByID(item.id)
}

func TestMenuItem_Checked(t *testing.T) {
	item := NewMenuItemCheckbox("Test", false)

	if item.checked {
		t.Error("checked should be false")
	}

	item.SetChecked(true)
	if !item.checked {
		t.Error("checked should be true after SetChecked(true)")
	}

	item.SetChecked(false)
	if item.checked {
		t.Error("checked should be false after SetChecked(false)")
	}

	// Clean up
	removeMenuItemByID(item.id)
}

func TestMenuItem_Hidden(t *testing.T) {
	item := NewMenuItem("Test")

	if item.hidden {
		t.Error("hidden should default to false")
	}

	item.SetHidden(true)
	if !item.hidden {
		t.Error("hidden should be true after SetHidden(true)")
	}

	item.SetHidden(false)
	if item.hidden {
		t.Error("hidden should be false after SetHidden(false)")
	}

	// Clean up
	removeMenuItemByID(item.id)
}

func TestMenuItem_Tooltip(t *testing.T) {
	item := NewMenuItem("Test")

	if item.tooltip != "" {
		t.Error("tooltip should default to empty string")
	}

	item.SetTooltip("Tooltip text")
	if item.tooltip != "Tooltip text" {
		t.Errorf("tooltip = %q, want %q", item.tooltip, "Tooltip text")
	}

	// Clean up
	removeMenuItemByID(item.id)
}

func TestMenuItem_Bitmap(t *testing.T) {
	item := NewMenuItem("Test")

	if item.bitmap != nil {
		t.Error("bitmap should default to nil")
	}

	bitmap := []byte{0x89, 0x50, 0x4E, 0x47}
	item.SetBitmap(bitmap)
	if len(item.bitmap) != 4 {
		t.Error("bitmap should be set")
	}

	// Clean up
	removeMenuItemByID(item.id)
}

func TestMenuItem_OnClick(t *testing.T) {
	item := NewMenuItem("Test")

	if item.callback != nil {
		t.Error("callback should default to nil")
	}

	called := false
	item.OnClick(func(ctx *Context) {
		called = true
	})

	if item.callback == nil {
		t.Error("callback should be set after OnClick")
	}

	// Call the callback
	item.callback(nil)
	if !called {
		t.Error("callback should have been called")
	}

	// Clean up
	removeMenuItemByID(item.id)
}

func TestMenuItem_SetAccelerator(t *testing.T) {
	item := NewMenuItem("Test")

	if item.accelerator != nil {
		t.Error("accelerator should default to nil")
	}

	result := item.SetAccelerator("Ctrl+A")
	if result != item {
		t.Error("SetAccelerator should return the same item for chaining")
	}

	if item.accelerator == nil {
		t.Error("accelerator should be set after SetAccelerator")
	}

	// Clean up
	removeMenuItemByID(item.id)
}

func TestMenuItem_GetAccelerator(t *testing.T) {
	item := NewMenuItem("Test")

	if item.GetAccelerator() != "" {
		t.Error("GetAccelerator should return empty string when not set")
	}

	item.SetAccelerator("Ctrl+B")
	acc := item.GetAccelerator()
	if acc == "" {
		t.Error("GetAccelerator should return non-empty string when set")
	}

	// Clean up
	removeMenuItemByID(item.id)
}

func TestMenuItem_RemoveAccelerator(t *testing.T) {
	item := NewMenuItem("Test")
	item.SetAccelerator("Ctrl+C")

	item.RemoveAccelerator()
	if item.accelerator != nil {
		t.Error("accelerator should be nil after RemoveAccelerator")
	}

	// Clean up
	removeMenuItemByID(item.id)
}

func TestMenuItem_Clone(t *testing.T) {
	original := NewMenuItem("Original")
	original.SetTooltip("Tooltip")
	original.SetChecked(true)
	original.SetEnabled(false)

	clone := original.Clone()

	if clone == original {
		t.Error("Clone should return a different pointer")
	}
	if clone.label != original.label {
		t.Error("Clone should have same label")
	}
	if clone.tooltip != original.tooltip {
		t.Error("Clone should have same tooltip")
	}
	if clone.checked != original.checked {
		t.Error("Clone should have same checked state")
	}
	if clone.disabled != original.disabled {
		t.Error("Clone should have same disabled state")
	}
	// Note: Clone preserves the ID (shallow clone behavior)

	// Clean up
	removeMenuItemByID(original.id)
	removeMenuItemByID(clone.id)
}

func TestMenuItem_Chaining(t *testing.T) {
	item := NewMenuItem("Test").
		SetTooltip("Tooltip").
		SetEnabled(false).
		SetHidden(true).
		SetAccelerator("Ctrl+X")

	if item.tooltip != "Tooltip" {
		t.Error("Chaining SetTooltip failed")
	}
	if !item.disabled {
		t.Error("Chaining SetEnabled failed")
	}
	if !item.hidden {
		t.Error("Chaining SetHidden failed")
	}
	if item.accelerator == nil {
		t.Error("Chaining SetAccelerator failed")
	}

	// Clean up
	removeMenuItemByID(item.id)
}
