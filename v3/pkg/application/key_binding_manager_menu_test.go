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
