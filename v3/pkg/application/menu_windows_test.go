//go:build windows && !server && !cef

package application

import (
	"runtime"
	"syscall"
	"testing"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/w32"
)

// These cover three defects in the shipping Windows backend's setMenu, found
// while porting it to the CEF frontend and filed as #6102, #6103 and #6104.
//
// Replacement tests use a hidden native window so SetMenu can succeed. The
// redraw (#6103) is a repaint side effect and is not asserted here.

// newTestMenuWindow creates a hidden window on the test's OS thread and
// releases its menu and window before unlocking that thread.
func newTestMenuWindow(t *testing.T) *windowsWebviewWindow {
	t.Helper()
	runtime.LockOSThread()
	t.Cleanup(runtime.UnlockOSThread)
	hwnd := w32.CreateWindowEx(0, w32.MustStringToUTF16Ptr("STATIC"), nil,
		w32.WS_OVERLAPPEDWINDOW, 0, 0, 100, 100, 0, 0, 0, nil)
	if hwnd == 0 {
		t.Fatal("failed to create test window")
	}
	w := &windowsWebviewWindow{parent: &WebviewWindow{}, hwnd: hwnd}
	t.Cleanup(func() {
		if w.menu != nil {
			w32.SetMenu(hwnd, 0)
			w.menu.Destroy()
		}
		w32.DestroyWindow(hwnd)
	})
	return w
}

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

// SetMenu(nil) on a window with no menu does nothing. Menu.Update
// dereferences its receiver once the application is running, so this used to
// panic (#6104).
func TestWindowsSetMenuNilWithoutAMenu(t *testing.T) {
	withTestApplication(t)

	w := &windowsWebviewWindow{parent: &WebviewWindow{}}
	w.setMenu(nil)

	if w.menu != nil {
		t.Fatalf("a nil menu left one on the window: %#v", w.menu)
	}
}

// SetMenu(nil) removes the window's menu bar: the window has no HMENU, the
// removed one is freed, and its accelerators no longer fire in the window
// (pg-ifz).
func TestWindowsSetMenuNilRemovesTheMenu(t *testing.T) {
	withTestApplication(t)

	w := newTestMenuWindow(t)
	w.parent.menuBindings = map[string]*MenuItem{}

	menu := NewMenu()
	menu.Add("Action").SetAccelerator("Ctrl+M").OnClick(func(*Context) {})
	w.setMenu(menu)
	if w.menu == nil || w32.GetMenu(w.hwnd) == 0 {
		t.Fatal("setMenu put no menu on the window")
	}
	if len(w.parent.menuBindings) != 1 {
		t.Fatalf("menu bindings after setMenu = %v, want the accelerator", w.parent.menuBindings)
	}
	removed := w.menu.menu

	w.setMenu(nil)

	if w.menu != nil {
		t.Error("SetMenu(nil) kept the window's menu")
	}
	if got := w32.GetMenu(w.hwnd); got != 0 {
		t.Errorf("window still shows HMENU %v after SetMenu(nil)", got)
	}
	if w32.IsMenu(removed) {
		t.Error("the removed menu's HMENU is still live, so it leaked")
	}
	if len(w.parent.menuBindings) != 0 {
		t.Errorf("removed menu kept accelerators %v", w.parent.menuBindings)
	}

	// The bar comes back with the next SetMenu.
	w.setMenu(menu)
	if w.menu == nil || w32.GetMenu(w.hwnd) == 0 {
		t.Error("setMenu after SetMenu(nil) put no menu on the window")
	}
}

// If the window cannot let go of its menu, SetMenu(nil) keeps it rather than
// destroying a menu the window still shows.
func TestWindowsSetMenuNilKeepsTheMenuOnFailure(t *testing.T) {
	withTestApplication(t)

	w := newTestMenuWindow(t)
	w.parent.menuBindings = map[string]*MenuItem{}
	menu := NewMenu()
	menu.Add("Action").SetAccelerator("Ctrl+M").OnClick(func(*Context) {})
	w.setMenu(menu)
	kept := w.menu

	hwnd := w.hwnd
	w.hwnd = 0
	w.setMenu(nil)
	w.hwnd = hwnd

	if w.menu != kept || !w32.IsMenu(kept.menu) {
		t.Fatal("a failed SetMenu(nil) dropped or destroyed the window's menu")
	}
	if len(w.parent.menuBindings) != 1 {
		t.Errorf("a failed SetMenu(nil) dropped the accelerators: %v", w.parent.menuBindings)
	}
}

// Replacing a menu drops the accelerators only the old menu held, and keeps
// those the replacement still uses.
func TestWindowsSetMenuDropsReplacedAccelerators(t *testing.T) {
	withTestApplication(t)

	w := newTestMenuWindow(t)
	w.parent.menuBindings = map[string]*MenuItem{}

	shared := NewMenu()
	kept := shared.Add("Kept").SetAccelerator("Ctrl+K").OnClick(func(*Context) {})
	first := NewMenu()
	gone := first.Add("Gone").SetAccelerator("Ctrl+G").OnClick(func(*Context) {})
	first.Append(shared)
	w.setMenu(first)

	w.setMenu(shared)

	bound := map[*MenuItem]bool{}
	for _, item := range w.parent.menuBindings {
		bound[item] = true
	}
	if bound[gone] {
		t.Error("the replaced menu's accelerator still fires")
	}
	if !bound[kept] {
		t.Errorf("menu bindings = %v, want the shared item's accelerator kept", w.parent.menuBindings)
	}
}

// setMenu builds a new Win32Menu, so the one it replaces is dropped. Nothing
// else ever frees it - Win32Menu.Update only frees the previous handle of the
// same object - so its HMENU leaked (#6102).
func TestWindowsSetMenuFreesTheMenuItReplaces(t *testing.T) {
	withTestApplication(t)

	w := newTestMenuWindow(t)

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

	w := newTestMenuWindow(t)

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

// A failed native assignment must keep the previous menu alive and release
// the replacement, including when no menu was previously installed.
func TestWindowsSetMenuRetainsPreviousOnFailure(t *testing.T) {
	getGuiResources := syscall.NewLazyDLL("user32.dll").NewProc("GetGuiResources")
	getCurrentProcess := syscall.NewLazyDLL("kernel32.dll").NewProc("GetCurrentProcess")
	process, _, _ := getCurrentProcess.Call()
	for _, existing := range []bool{false, true} {
		name := "without previous menu"
		if existing {
			name = "with previous menu"
		}
		t.Run(name, func(t *testing.T) {
			withTestApplication(t)
			w := newTestMenuWindow(t)
			if existing {
				first := NewMenu()
				first.Add("First")
				w.setMenu(first)
				if w.menu == nil {
					t.Fatal("setMenu built no initial menu")
				}
			}
			previous := w.menu
			// Keep the real window alive, but pass an invalid handle to SetMenu.
			w.hwnd = 0
			second := NewMenu()
			item := second.Add("Second")
			// Initialise the model's separate menu before measuring replacement
			// handles. Rebuilding it must not change the USER object count.
			second.Update()
			before, _, _ := getGuiResources.Call(process, 1) // GR_USEROBJECTS
			if before == 0 {
				t.Fatal("failed to read USER object count")
			}
			w.setMenu(second)
			after, _, _ := getGuiResources.Call(process, 1)

			if w.menu != previous {
				t.Error("failed SetMenu replaced the window's menu")
			}
			if previous != nil && !w32.IsMenu(previous.menu) {
				t.Error("failed SetMenu destroyed the previous menu")
			}
			if item.impl != nil {
				t.Error("failed SetMenu retained a binding to the rejected menu")
			}
			if after != before {
				t.Error("failed SetMenu leaked the rejected replacement")
			}
		})
	}
}

func TestWindowsSetMenuRetainsItemBindingsOnFailure(t *testing.T) {
	withTestApplication(t)
	w := newTestMenuWindow(t)
	menu := NewMenu()
	item := menu.Add("Only")
	w.setMenu(menu)
	if w.menu == nil {
		t.Fatal("setMenu built no initial menu")
	}
	previous := w.menu
	impl := item.impl.(*windowsMenuItem)
	added := menu.Add("New")
	w.hwnd = 0
	w.setMenu(menu)

	if w.menu != previous || item.impl != impl {
		t.Fatal("failed SetMenu did not restore the previous menu and item binding")
	}
	if added.impl != nil {
		t.Error("failed SetMenu retained a native binding for a new item")
	}
	item.SetEnabled(false)
	info := w32.MENUITEMINFO{FMask: w32.MIIM_STATE}
	info.CbSize = uint32(unsafe.Sizeof(info))
	if !w32.GetMenuItemInfo(previous.menu, uint32(impl.id), false, &info) {
		t.Fatal("failed to read the retained native menu item")
	}
	if info.FState&w32.MFS_DISABLED == 0 {
		t.Error("item update did not reach the retained native menu")
	}
}
