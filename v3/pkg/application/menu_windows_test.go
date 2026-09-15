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
