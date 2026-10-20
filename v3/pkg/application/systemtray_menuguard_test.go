package application

import (
	"testing"
	"time"
)

func TestTrayMenuGuard(t *testing.T) {
	base := time.Unix(1000, 0)
	var g trayMenuGuard

	if g.active(base) {
		t.Fatal("zero value must be inactive")
	}

	g.opened(base)
	if !g.active(base.Add(trayMenuOpenTimeout - time.Nanosecond)) {
		t.Error("must be active before the timeout")
	}

	g.closed()
	if g.active(base) {
		t.Error("closed must clear the flag")
	}

	g.opened(base)
	if g.active(base.Add(trayMenuOpenTimeout)) {
		t.Error("must expire at the timeout when no closed event arrives")
	}
	if g.active(base) {
		t.Error("expiry must be sticky")
	}

	// A repeated open (AboutToShow then "opened") extends the deadline.
	g.opened(base)
	g.opened(base.Add(3 * time.Second))
	if !g.active(base.Add(trayMenuOpenTimeout + time.Second)) {
		t.Error("second open must extend the deadline")
	}
}

func TestTrayMenuActiveFor(t *testing.T) {
	win := &WebviewWindow{}
	other := &WebviewWindow{}
	attached := &SystemTray{}
	attached.attachedWindow.Window = win
	detached := &SystemTray{}
	app := &App{systemTrays: map[uint]*SystemTray{1: attached, 2: detached}}

	if app.trayMenuActiveFor(win) {
		t.Fatal("no menu open yet")
	}
	detached.menuGuard.opened(time.Now())
	if app.trayMenuActiveFor(win) {
		t.Error("menu of a tray without this window must not suppress")
	}
	attached.menuGuard.opened(time.Now())
	if !app.trayMenuActiveFor(win) {
		t.Error("open menu of the attached tray must suppress")
	}
	if app.trayMenuActiveFor(other) {
		t.Error("other windows must not be affected")
	}
	attached.menuGuard.closed()
	if app.trayMenuActiveFor(win) {
		t.Error("closed menu must stop suppressing")
	}
	var nilApp *App
	if nilApp.trayMenuActiveFor(win) || app.trayMenuActiveFor(nil) {
		t.Error("nil app/window must be inactive")
	}
}
