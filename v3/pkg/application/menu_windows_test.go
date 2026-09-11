//go:build windows && !server && !cef

package application

import (
	"testing"

	"github.com/wailsapp/wails/v3/pkg/w32"
)

// These cover three defects in the shipping Windows backend's setMenu, found
// while porting it to the CEF frontend and filed as #6102, #6103 and #6104.
//
// Two of the three are testable here because setMenu does its own bookkeeping
// before it touches the window: a window handle of zero makes SetMenu and
// DrawMenuBar no-ops without changing anything this asserts on. The redraw
// (#6103) is a repaint side effect with nothing to assert on and is not
// covered here.

// withTestApplication gives the package a global application to log through,
// since setMenu debug-logs on a path these tests take, and puts back whatever
// was there.
func withTestApplication(t *testing.T) {
	t.Helper()
	previous := globalApplication
	globalApplication = nil
	app := New(Options{Name: "menu test"})
	// Marked as running, because that is the only state in which any of this
	// misbehaves: Menu.Update returns before touching its receiver otherwise,
	// so a nil menu would not panic and a real one would not be rebuilt.
	app.runLock.Lock()
	app.running = true
	app.runLock.Unlock()
	t.Cleanup(func() { globalApplication = previous })
}

// A nil menu means "no menu", which is what macOS makes of it. Menu.Update
// dereferences its receiver once the application is running, so this used to
// panic rather than do nothing (#6104).
func TestWindowsSetMenuIgnoresANilMenu(t *testing.T) {
	withTestApplication(t)

	w := &windowsWebviewWindow{parent: &WebviewWindow{}}
	w.setMenu(nil)

	if w.menu != nil {
		t.Fatalf("a nil menu left one on the window: %#v", w.menu)
	}
}

// setMenu builds a new Win32Menu, so the one it replaces is dropped. Nothing
// else ever frees it - Win32Menu.Update only frees the previous handle of the
// same object - so its HMENU leaked (#6102).
func TestWindowsSetMenuFreesTheMenuItReplaces(t *testing.T) {
	withTestApplication(t)

	w := &windowsWebviewWindow{parent: &WebviewWindow{}}

	first := NewMenu()
	first.Add("First")
	w.setMenu(first)

	if w.menu == nil {
		t.Fatal("setMenu built no menu")
	}
	replaced := w.menu.menu
	if !w32.IsMenu(replaced) {
		t.Fatalf("the first menu was not a live HMENU: %v", replaced)
	}

	second := NewMenu()
	second.Add("Second")
	w.setMenu(second)

	if w.menu.menu == replaced {
		t.Fatal("setMenu reused the previous HMENU rather than building one")
	}
	if w32.IsMenu(replaced) {
		t.Error("the replaced menu's HMENU is still live, so it leaked (#6102)")
	}
	if !w32.IsMenu(w.menu.menu) {
		t.Error("setMenu destroyed the menu it had just put on the window")
	}
}

// The same Menu twice is the case where freeing the old Win32Menu can reach
// state the new one now owns: both mappings hold the same *MenuItem values,
// whose impl the second build has reassigned. The old menu must still go, and
// the new one must survive.
func TestWindowsSetMenuWithTheSameMenuTwice(t *testing.T) {
	withTestApplication(t)

	w := &windowsWebviewWindow{parent: &WebviewWindow{}}

	menu := NewMenu()
	menu.Add("Only")

	w.setMenu(menu)
	replaced := w.menu.menu

	w.setMenu(menu)

	if w32.IsMenu(replaced) {
		t.Error("the replaced menu's HMENU is still live, so it leaked (#6102)")
	}
	if !w32.IsMenu(w.menu.menu) {
		t.Error("setMenu destroyed the menu it had just put on the window")
	}
}
