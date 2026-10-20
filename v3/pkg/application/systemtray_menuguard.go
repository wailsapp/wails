package application

import (
	"sync"
	"time"
)

// trayMenuOpenTimeout bounds how long a tray's context menu is assumed open
// after the host last announced it. dbusmenu hosts are not guaranteed to send
// a "closed" event (the menu can be dismissed by clicking away or by the
// host going away), so an "open" state with no deadline could suppress
// hide-on-focus-loss forever.
const trayMenuOpenTimeout = 5 * time.Second

// trayMenuGuard records whether a system tray's context menu is open, so that
// the focus loss caused by the menu grabbing focus is not mistaken for the
// user clicking away from the attached window (#4494).
//
// Only the Linux dbusmenu implementation marks it open; on every other
// platform it stays closed and has no effect. The zero value is ready to use.
type trayMenuGuard struct {
	mu       sync.Mutex
	open     bool
	deadline time.Time
}

// opened marks the menu open until now+trayMenuOpenTimeout. Calling it again
// (the host announces both AboutToShow and "opened", and AboutToShow fires per
// submenu) extends the deadline.
func (g *trayMenuGuard) opened(now time.Time) {
	g.mu.Lock()
	g.open = true
	g.deadline = now.Add(trayMenuOpenTimeout)
	g.mu.Unlock()
}

// closed marks the menu closed.
func (g *trayMenuGuard) closed() {
	g.mu.Lock()
	g.open = false
	g.mu.Unlock()
}

// active reports whether the menu is considered open at now. An open state
// past its deadline is expired and cleared.
func (g *trayMenuGuard) active(now time.Time) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.open && !now.Before(g.deadline) {
		g.open = false
	}
	return g.open
}

// trayMenuActiveFor reports whether a system tray that has window attached
// currently has its context menu open.
func (a *App) trayMenuActiveFor(window Window) bool {
	if a == nil || window == nil {
		return false
	}
	now := time.Now()
	a.systemTraysLock.Lock()
	defer a.systemTraysLock.Unlock()
	for _, tray := range a.systemTrays {
		if tray.attachedWindow.Window == window && tray.menuGuard.active(now) {
			return true
		}
	}
	return false
}
