//go:build linux && cgo && !android && !server

package application

import (
	"os"
	"strings"
)

// newGlobalShortcutImpl selects the appropriate Linux backend.
//
// On X11 sessions the XGrabKey-based backend is used: it is self-contained,
// requires no portal support and grabs the exact accelerator requested.
//
// On Wayland sessions there is, by design, no way for a client to grab keys
// directly. The only sanctioned mechanism is the XDG Desktop Portal's
// org.freedesktop.portal.GlobalShortcuts interface, so the portal backend is
// used there. Note that under the portal the compositor (and ultimately the
// user) decides the final key binding; see portalGlobalShortcuts.
func newGlobalShortcutImpl(manager *GlobalShortcutManager) globalShortcutImpl {
	if isWaylandSession() {
		return newPortalGlobalShortcuts(manager)
	}
	return newX11GlobalShortcuts(manager)
}

// isWaylandSession reports whether the process is running under a Wayland
// session. XDG_SESSION_TYPE is authoritative when set; otherwise the presence
// of WAYLAND_DISPLAY is used as a fallback.
func isWaylandSession() bool {
	switch strings.ToLower(os.Getenv("XDG_SESSION_TYPE")) {
	case "wayland":
		return true
	case "x11":
		return false
	}
	return os.Getenv("WAYLAND_DISPLAY") != ""
}
