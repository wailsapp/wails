//go:build linux && !android && !server

package application

import (
	"testing"

	"github.com/godbus/dbus/v5"
)

// trayHandlerCounts records which of the SystemTray callbacks a dbusmenu
// event reached, so each assertion can pin the whole set rather than the one
// callback it is about — an event running a callback it should not is the
// failure mode here.
type trayHandlerCounts struct {
	clicks, opens, closes int
}

// newCountingTray builds a tray with all three callbacks attached and an App
// with no logger installed: the dbusmenu callbacks log through
// globalApplication, which is nil in a test binary.
func newCountingTray(t *testing.T, counts *trayHandlerCounts) *linuxSystemTray {
	t.Helper()

	prev := globalApplication
	globalApplication = &App{}
	t.Cleanup(func() { globalApplication = prev })

	return &linuxSystemTray{parent: &SystemTray{
		clickHandler: func() { counts.clicks++ },
		onMenuOpen:   func() { counts.opens++ },
		onMenuClose:  func() { counts.closes++ },
	}}
}

func (c trayHandlerCounts) check(t *testing.T, after string, clicks, opens, closes int) {
	t.Helper()
	if c.clicks != clicks || c.opens != opens || c.closes != closes {
		t.Errorf("after %s: click=%d onMenuOpen=%d onMenuClose=%d, want %d/%d/%d",
			after, c.clicks, c.opens, c.closes, clicks, opens, closes)
	}
}

// The dbusmenu "opened" event is the host announcing that it is about to show
// the context menu, which is the secondary button. Running the click handler
// here as well made every right-click do whatever a left-click does — on an
// app whose click handler toggles a window, right-clicking toggled it.
func TestLinuxSystemTrayMenuEventsDoNotClick(t *testing.T) {
	var counts trayHandlerCounts
	tray := newCountingTray(t, &counts)

	if err := tray.Event(0, "opened", dbus.Variant{}, 0); err != nil {
		t.Fatalf(`Event("opened") returned %v`, err)
	}
	counts.check(t, `Event("opened")`, 0, 1, 0)

	if err := tray.Event(0, "closed", dbus.Variant{}, 0); err != nil {
		t.Fatalf(`Event("closed") returned %v`, err)
	}
	counts.check(t, `Event("closed")`, 0, 1, 1)
}

// The counterpart: ItemIsMenu is published as false, so the primary button
// arrives as Activate, which is now the only path to the click handler.
func TestLinuxSystemTrayActivateClicks(t *testing.T) {
	var counts trayHandlerCounts
	tray := newCountingTray(t, &counts)

	if err := tray.Activate(12, 34); err != nil {
		t.Fatalf("Activate returned %v", err)
	}
	counts.check(t, "Activate", 1, 0, 0)

	if tray.lastClickX != 12 || tray.lastClickY != 34 {
		t.Errorf("Activate recorded (%d,%d), want (12,34)", tray.lastClickX, tray.lastClickY)
	}
}

// A tray with no handlers attached must not panic on the same events.
func TestLinuxSystemTrayEventsWithoutHandlers(t *testing.T) {
	prev := globalApplication
	globalApplication = &App{}
	t.Cleanup(func() { globalApplication = prev })

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
