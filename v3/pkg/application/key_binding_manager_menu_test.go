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
	manager.removeMenuBinding(item)
	if !manager.Process(key, visible) {
		t.Fatal("removing menu binding removed application key binding")
	}
}
