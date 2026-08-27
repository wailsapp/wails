//go:build linux && !android && !server

package application

import (
	"testing"

	"github.com/godbus/dbus/v5"
)

// withStubApplication installs an App with no logger for the duration of the
// test: the dbusmenu callbacks log through globalApplication, which is nil in a
// test binary.
func withStubApplication(t *testing.T) {
	t.Helper()
	prev := globalApplication
	globalApplication = &App{}
	t.Cleanup(func() { globalApplication = prev })
}

// The dbusmenu "opened" event is the host announcing that it is about to show
// the context menu, which is the secondary button. Running the click handler
// here as well made every right-click do whatever a left-click does — on an
// app whose click handler toggles a window, right-clicking toggled it.
func TestLinuxSystemTrayMenuEventsDoNotClick(t *testing.T) {
	withStubApplication(t)

	var clicks, opens, closes int
	tray := &linuxSystemTray{parent: &SystemTray{
		clickHandler: func() { clicks++ },
		onMenuOpen:   func() { opens++ },
		onMenuClose:  func() { closes++ },
	}}

	if err := tray.Event(0, "opened", dbus.Variant{}, 0); err != nil {
		t.Fatalf(`Event("opened") returned %v`, err)
	}
	if opens != 1 {
		t.Errorf(`Event("opened") ran onMenuOpen %d times, want 1`, opens)
	}
	if clicks != 0 {
		t.Errorf(`Event("opened") ran the click handler %d times, want 0`, clicks)
	}

	if err := tray.Event(0, "closed", dbus.Variant{}, 0); err != nil {
		t.Fatalf(`Event("closed") returned %v`, err)
	}
	if closes != 1 {
		t.Errorf(`Event("closed") ran onMenuClose %d times, want 1`, closes)
	}
	if clicks != 0 {
		t.Errorf(`Event("closed") ran the click handler %d times, want 0`, clicks)
	}
}

// The counterpart: ItemIsMenu is published as false, so the primary button
// arrives as Activate, which is now the only path to the click handler.
func TestLinuxSystemTrayActivateClicks(t *testing.T) {
	withStubApplication(t)

	var clicks int
	tray := &linuxSystemTray{parent: &SystemTray{
		clickHandler: func() { clicks++ },
	}}

	if err := tray.Activate(12, 34); err != nil {
		t.Fatalf("Activate returned %v", err)
	}
	if clicks != 1 {
		t.Errorf("Activate ran the click handler %d times, want 1", clicks)
	}
	if tray.lastClickX != 12 || tray.lastClickY != 34 {
		t.Errorf("Activate recorded (%d,%d), want (12,34)", tray.lastClickX, tray.lastClickY)
	}
}

// A tray with no handlers attached must not panic on the same events.
func TestLinuxSystemTrayEventsWithoutHandlers(t *testing.T) {
	withStubApplication(t)

	tray := &linuxSystemTray{parent: &SystemTray{}}
	for _, eventID := range []string{"opened", "closed"} {
		if err := tray.Event(0, eventID, dbus.Variant{}, 0); err != nil {
			t.Fatalf("Event(%q) returned %v", eventID, err)
		}
	}
	if err := tray.Activate(0, 0); err != nil {
		t.Fatalf("Activate returned %v", err)
	}
}
