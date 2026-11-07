package application

import (
	"strings"
	"testing"
)

// TestRemoveDifferentCase check that removing a key binding which differs only in case works as expected
func TestRemoveDifferentCase(t *testing.T) {
	const key = "Ctrl+K"
	app := &App{keyBindings: make(map[string]func(Window))}
	manager := newKeyBindingManager(app)
	manager.Add(strings.ToLower(key), func(Window) {})
	manager.Remove(key)
	if len(manager.GetAll()) != 0 {
		t.Fatal("accelerator with different case not removed")
	}
}

// TestLowercaseFires checks that a binding fires even if added with a different case (lowercase)
func TestLowercaseFires(t *testing.T) {
	const key = "Ctrl+K"
	app := &App{keyBindings: make(map[string]func(Window))}
	manager := newKeyBindingManager(app)

	called := false
	manager.Add(strings.ToLower(key), func(Window) { called = true })
	if !manager.Process(key, &WebviewWindow{}) {
		t.Fatal("key binding was not processed")
	}
	if !called {
		t.Fatal("key binding callback was not called")
	}
}

// TestRemovedDoesNotFire checks that adding then removing the same key binding using
// case causes the binding to not fire
func TestRemovedDoesNotFire(t *testing.T) {
	const key = "Ctrl+K"
	app := &App{keyBindings: make(map[string]func(Window))}
	manager := newKeyBindingManager(app)

	called := false
	manager.Add(key, func(Window) { called = true })
	manager.Remove(strings.ToLower(key))
	if manager.Process(key, &WebviewWindow{}) {
		t.Fatal("removed binding was processed")
	}
	if called {
		t.Fatal("removed binding callback was called")
	}
}

// TestSameKeyDiffCase checks that overriding a key binding using different case
// results only in the last callback added being called
func TestSameKeyDiffCase(t *testing.T) {
	const key = "Ctrl+K"
	app := &App{keyBindings: make(map[string]func(Window))}
	manager := newKeyBindingManager(app)

	called := ""
	manager.Add(key, func(Window) { called = "first" })
	// Override the previous binding (should ignore case) with a new callback
	manager.Add(strings.ToLower(key), func(Window) { called = "override" })
	if !manager.Process(key, &WebviewWindow{}) {
		t.Fatal("removed binding was processed")
	}
	if called != "last" {
		t.Fatal("correct callback was not called")
	}
}
