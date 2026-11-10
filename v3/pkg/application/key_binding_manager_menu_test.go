package application

import (
	"testing"
	"time"
)

func TestApplicationMenuBindingRespectsWindowRemoval(t *testing.T) {
	app := &App{keyBindings: make(map[string]func(Window))}
	manager := newKeyBindingManager(app)
	menu := NewMenu()
	clicks := make(chan struct{}, 2)
	item := menu.AddSubmenu("File").Add("Action").SetAccelerator("Ctrl+M").OnClick(func(*Context) { clicks <- struct{}{} })
	app.applicationMenu = menu
	manager.addMenuBinding(item)
	if got := manager.GetAll(); len(got) != 1 {
		t.Fatalf("menu binding is absent from GetAll: %v", got)
	}
	key := item.accelerator.String()
	removed := &WebviewWindow{menuRemoved: true}
	visible := &WebviewWindow{}
	if manager.Process(key, removed) {
		t.Fatal("removed application menu shortcut was handled")
	}
	if !manager.Process(key, visible) {
		t.Fatal("another window lost the application menu shortcut")
	}
	select {
	case <-clicks:
	case <-time.After(time.Second):
		t.Fatal("menu callback did not fire")
	}
	called := false
	manager.Add(key, func(Window) { called = true })
	if !manager.Process(key, removed) || !called {
		t.Fatal("removal suppressed an unrelated application key binding")
	}
	called = false
	if !manager.Process(key, visible) || !called {
		t.Fatal("menu shortcut masked an application key binding")
	}
	if got := manager.GetAll(); len(got) != 1 {
		t.Fatalf("duplicate shortcut in GetAll: %v", got)
	}
	manager.removeMenuBinding(item)
	if !manager.Process(key, visible) {
		t.Fatal("removing menu binding removed application key binding")
	}
	manager.addMenuBinding(item)
	manager.Remove(key)
	if manager.Process(key, visible) || len(manager.GetAll()) != 0 {
		t.Fatal("Remove left a binding behind")
	}
}

// A disabled item's shortcut does nothing, as its greyed-out row does; enabled
// again, it works. The window bindings (Windows) and the application's behave
// the same.
func TestMenuBindingsSkipDisabledItems(t *testing.T) {
	app := &App{keyBindings: make(map[string]func(Window))}
	manager := newKeyBindingManager(app)
	clicks := make(chan struct{}, 4)
	item := NewMenu().Add("Action").SetAccelerator("Ctrl+D").OnClick(func(*Context) { clicks <- struct{}{} })
	key := item.accelerator.String()
	manager.addMenuBinding(item)
	window := &WebviewWindow{menuBindings: map[string]*MenuItem{}}
	window.addMenuBinding(item.accelerator, item)
	previous := globalApplication
	globalApplication = app
	app.KeyBinding = manager
	t.Cleanup(func() { globalApplication = previous })

	expectClick := func(want bool) {
		t.Helper()
		select {
		case <-clicks:
			if !want {
				t.Fatal("a disabled item's shortcut fired")
			}
		case <-time.After(200 * time.Millisecond):
			if want {
				t.Fatal("an enabled item's shortcut did not fire")
			}
		}
	}

	item.disabled = true
	if manager.Process(key, window) {
		t.Error("the application handled a disabled item's shortcut")
	}
	if window.processKeyBinding(key) {
		t.Error("the window handled a disabled item's shortcut")
	}
	expectClick(false)

	item.disabled = false
	if !manager.Process(key, window) {
		t.Error("the application did not handle an enabled item's shortcut")
	}
	expectClick(true)
	if !window.processKeyBinding(key) {
		t.Error("the window did not handle an enabled item's shortcut")
	}
	expectClick(true)
}

// After SetAccelerator a rebuild must not leave the old shortcut bound:
// removeMenuBindingsTo drops an item's bindings whatever their key, and
// leaves other items' alone.
func TestRemoveMenuBindingsToDropsChangedAccelerators(t *testing.T) {
	app := &App{keyBindings: make(map[string]func(Window))}
	manager := newKeyBindingManager(app)
	window := &WebviewWindow{menuBindings: map[string]*MenuItem{}}
	menu := NewMenu()
	rebound := menu.Add("Rebound").SetAccelerator("Ctrl+E")
	other := menu.Add("Other").SetAccelerator("Ctrl+K")
	oldKey := rebound.accelerator.String()
	for _, item := range []*MenuItem{rebound, other} {
		manager.addMenuBinding(item)
		window.addMenuBinding(item.accelerator, item)
	}
	rebound.SetAccelerator("Ctrl+O")

	built := map[int]*MenuItem{1: rebound}
	manager.removeMenuBindingsTo(built)
	window.removeMenuBindingsTo(built)
	manager.removeMenuBindingsTo(nil) // nothing to drop: a no-op

	if manager.menuBindings[oldKey] != nil || window.menuBindings[oldKey] != nil {
		t.Error("the old accelerator is still bound after the rebuild")
	}
	if manager.menuBindings[other.accelerator.String()] != other || window.menuBindings[other.accelerator.String()] != other {
		t.Error("another item's binding was dropped")
	}
}
