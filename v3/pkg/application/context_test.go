package application

import (
	"testing"
)

func TestNewContext(t *testing.T) {
	ctx := newContext()
	if ctx == nil {
		t.Fatal("newContext() returned nil")
	}
	if ctx.data == nil {
		t.Error("newContext() should initialize data map")
	}
}

func TestContext_ClickedMenuItem_NotExists(t *testing.T) {
	ctx := newContext()
	result := ctx.ClickedMenuItem()
	if result != nil {
		t.Error("ClickedMenuItem() should return nil when not set")
	}
}

func TestContext_ClickedMenuItem_Exists(t *testing.T) {
	ctx := newContext()
	menuItem := &MenuItem{label: "Test"}
	ctx.withClickedMenuItem(menuItem)

	result := ctx.ClickedMenuItem()
	if result == nil {
		t.Fatal("ClickedMenuItem() should return the menu item")
	}
	if result != menuItem {
		t.Error("ClickedMenuItem() should return the same menu item")
	}
}

func TestContext_IsChecked_NotSet(t *testing.T) {
	ctx := newContext()
	if ctx.IsChecked() {
		t.Error("IsChecked() should return false when not set")
	}
}

func TestContext_IsChecked_True(t *testing.T) {
	ctx := newContext()
	ctx.withChecked(true)
	if !ctx.IsChecked() {
		t.Error("IsChecked() should return true when set to true")
	}
}

func TestContext_IsChecked_False(t *testing.T) {
	ctx := newContext()
	ctx.withChecked(false)
	if ctx.IsChecked() {
		t.Error("IsChecked() should return false when set to false")
	}
}

func TestContext_ContextMenuData_Empty(t *testing.T) {
	ctx := newContext()
	result := ctx.ContextMenuData()
	if result != "" {
		t.Errorf("ContextMenuData() should return empty string when not set, got %q", result)
	}
}

func TestContext_ContextMenuData_Exists(t *testing.T) {
	ctx := newContext()
	data := &ContextMenuData{Data: "test-data"}
	ctx.withContextMenuData(data)

	result := ctx.ContextMenuData()
	if result != "test-data" {
		t.Errorf("ContextMenuData() = %q, want %q", result, "test-data")
	}
}

func TestContext_ContextMenuData_NilData(t *testing.T) {
	ctx := newContext()
	ctx.withContextMenuData(nil)

	result := ctx.ContextMenuData()
	if result != "" {
		t.Errorf("ContextMenuData() should return empty string for nil data, got %q", result)
	}
}

func TestContext_ContextMenuData_WrongType(t *testing.T) {
	ctx := newContext()
	// Manually set wrong type to test type assertion
	ctx.data[contextMenuData] = 123

	result := ctx.ContextMenuData()
	if result != "" {
		t.Errorf("ContextMenuData() should return empty string for non-string type, got %q", result)
	}
}

func TestContext_WithClickedMenuItem_Chaining(t *testing.T) {
	ctx := newContext()
	menuItem := &MenuItem{label: "Test"}

	returnedCtx := ctx.withClickedMenuItem(menuItem)
	if returnedCtx != ctx {
		t.Error("withClickedMenuItem should return the same context for chaining")
	}
}

func TestContext_WithContextMenuData_Chaining(t *testing.T) {
	ctx := newContext()
	data := &ContextMenuData{Data: "test"}

	returnedCtx := ctx.withContextMenuData(data)
	if returnedCtx != ctx {
		t.Error("withContextMenuData should return the same context for chaining")
	}
}
