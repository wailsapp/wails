//go:build darwin && purego && !ios && !server

package application

import "testing"

// newPuregoTestTray registers a tray with no status item, so setMenu does not
// touch AppKit and the displayed callbacks can run without a status bar.
func newPuregoTestTray(t *testing.T, id uint, menu *Menu) *macosSystemTray {
	t.Helper()
	tray := &macosSystemTray{id: id, menu: menu}
	systemTrayMap[id] = tray
	t.Cleanup(func() { delete(systemTrayMap, id) })
	return tray
}

// TestPuregoSystrayDisplayedCallback_Flag checks that the displayed callback
// sets and clears Menu.displayed for the tray's current menu.
func TestPuregoSystrayDisplayedCallback_Flag(t *testing.T) {
	menu := NewMenu()
	newPuregoTestTray(t, 4721, menu)

	systrayMenuDisplayedCallback(4721, 1)
	if !menu.displayed.Load() {
		t.Fatal("open: menu.displayed = false, want true")
	}

	systrayMenuDisplayedCallback(4721, 0)
	if menu.displayed.Load() {
		t.Fatal("close: menu.displayed = true, want false")
	}
}

// TestPuregoSystrayDisplayedCallback_UnknownTray checks that a callback for an
// unregistered tray id is a no-op rather than a nil dereference.
func TestPuregoSystrayDisplayedCallback_UnknownTray(t *testing.T) {
	menu := NewMenu()
	newPuregoTestTray(t, 4722, menu)

	systrayMenuDisplayedCallback(4799, 1)
	if menu.displayed.Load() {
		t.Fatal("callback for an unknown tray marked another tray's menu displayed")
	}
}

// TestPuregoMenuClear_WhileDisplayedKeepsIDs checks the reason for tracking the
// open menu: while it is on screen, Clear keeps the retired item IDs
// resolvable so a click on a stale row still finds its handler (#4719). Once
// the menu closes, the next Clear releases the previous generation.
func TestPuregoMenuClear_WhileDisplayedKeepsIDs(t *testing.T) {
	menu := NewMenu()
	// A non-nil impl (stubMenuImpl, menu_internal_test.go) selects the
	// retained-ID path that a systray menu takes once a native menu exists.
	menu.impl = stubMenuImpl{}
	newPuregoTestTray(t, 4723, menu)

	first := menu.Add("first")
	t.Cleanup(func() { removeMenuItemByID(first.id) })
	firstID := first.id

	systrayMenuDisplayedCallback(4723, 1)
	menu.Clear()
	if getMenuItemByID(firstID) == nil {
		t.Fatal("Clear while displayed released a retired ID that the open menu still shows")
	}

	systrayMenuDisplayedCallback(4723, 0)
	second := menu.Add("second")
	t.Cleanup(func() { removeMenuItemByID(second.id) })
	menu.Clear()
	if getMenuItemByID(firstID) != nil {
		t.Fatal("Clear after close kept a two-generations-old ID, want it released")
	}
	if getMenuItemByID(second.id) == nil {
		t.Fatal("Clear after close released the just-retired generation, want it kept")
	}
}

// TestPuregoSystrayDisplayed_SetMenuWhileOpen checks that a menu replaced by
// SetMenu while open is the one that gets released on close, and that the new
// menu is marked only once it is actually shown.
func TestPuregoSystrayDisplayed_SetMenuWhileOpen(t *testing.T) {
	old, replacement := NewMenu(), NewMenu()
	tray := newPuregoTestTray(t, 4724, old)

	systrayMenuDisplayedCallback(4724, 1)
	if !old.displayed.Load() {
		t.Fatal("open: old menu not displayed")
	}

	// SetMenu replaces the menu while the old one is on screen.
	tray.setMenu(replacement)
	if replacement.displayed.Load() {
		t.Fatal("SetMenu marked the replacement displayed before it was shown")
	}

	systrayMenuDisplayedCallback(4724, 0)
	if old.displayed.Load() {
		t.Fatal("close: replaced menu still marked displayed (leaked)")
	}
	if replacement.displayed.Load() {
		t.Fatal("close: replacement menu marked displayed although it was never shown")
	}
	if tray.displayedMenu != nil {
		t.Fatal("close: displayedMenu not cleared")
	}

	// Reopening shows the replacement and releases nothing else.
	systrayMenuDisplayedCallback(4724, 1)
	if !replacement.displayed.Load() {
		t.Fatal("reopen: replacement menu not displayed")
	}
	systrayMenuDisplayedCallback(4724, 0)
	if replacement.displayed.Load() {
		t.Fatal("second close: replacement menu still displayed")
	}
}

// TestPuregoSystrayDisplayed_ReopenWithReplacement checks the case where the
// tray reports a second open (menuWillOpen:) while a SetMenu replacement is in
// place without an intervening close: the old menu is released, not leaked.
func TestPuregoSystrayDisplayed_ReopenWithReplacement(t *testing.T) {
	old, replacement := NewMenu(), NewMenu()
	tray := newPuregoTestTray(t, 4725, old)

	systrayMenuDisplayedCallback(4725, 1)
	tray.setMenu(replacement)
	systrayMenuDisplayedCallback(4725, 1)

	if old.displayed.Load() {
		t.Fatal("reopen: previously displayed menu left marked displayed")
	}
	if !replacement.displayed.Load() {
		t.Fatal("reopen: replacement menu not marked displayed")
	}
	systrayMenuDisplayedCallback(4725, 0)
	if replacement.displayed.Load() {
		t.Fatal("close: replacement menu still displayed")
	}
}

// TestPuregoStatusItemController_MenuDelegateSelectors is the AppKit-level
// check: the registered controller class answers the NSMenuDelegate selectors
// that drive the displayed flag, both on the class and on a live instance.
func TestPuregoStatusItemController_MenuDelegateSelectors(t *testing.T) {
	cls := statusItemControllerClassRef()

	for _, sel := range []string{"menuWillOpen:", "menuDidClose:", "statusItemClicked:"} {
		if !get[bool](cls, "instancesRespondToSelector:", sel_(sel)) {
			t.Errorf("controller class does not respond to %s", sel)
		}
	}
	if get[bool](cls, "instancesRespondToSelector:", sel_("menuDidOpen:")) {
		t.Error("controller class responds to a selector it does not implement")
	}

	inst := cls.send("alloc").send("init")
	t.Cleanup(func() { inst.send("release") })
	if !respondsTo(inst, "menuWillOpen:") {
		t.Error("controller instance does not respond to menuWillOpen:")
	}
	if !respondsTo(inst, "menuDidClose:") {
		t.Error("controller instance does not respond to menuDidClose:")
	}
}

// TestPuregoStatusItemController_DelegateDrivesDisplayed sends the delegate
// messages through the real controller instance, so the selector-to-IMP wiring
// (not just the Go helper) is exercised: willOpen sets the flag, didClose
// clears it.
func TestPuregoStatusItemController_DelegateDrivesDisplayed(t *testing.T) {
	menu := NewMenu()
	tray := newPuregoTestTray(t, 4726, menu)

	inst := statusItemControllerClassRef().send("alloc").send("init")
	t.Cleanup(func() { inst.send("release") })
	systrayControllerMu.Lock()
	systrayControllerMap[uintptr(inst.ptr())] = tray
	systrayControllerMu.Unlock()
	t.Cleanup(func() {
		systrayControllerMu.Lock()
		delete(systrayControllerMap, uintptr(inst.ptr()))
		systrayControllerMu.Unlock()
	})

	nsMenu := class("NSMenu").send("alloc").send("init")
	t.Cleanup(func() { nsMenu.send("release") })
	// showMenu/pre-click attach the controller as the menu delegate before the
	// menu opens; menuDidClose: must detach it again.
	nsMenu.send("setDelegate:", inst)

	inst.send("menuWillOpen:", nsMenu)
	if !menu.displayed.Load() {
		t.Fatal("menuWillOpen: did not mark the tray menu displayed")
	}

	inst.send("menuDidClose:", nsMenu)
	if menu.displayed.Load() {
		t.Fatal("menuDidClose: did not clear the displayed flag")
	}
	if !nsMenu.send("delegate").isNil() {
		t.Fatal("menuDidClose: left the menu delegate attached")
	}
}
