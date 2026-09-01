//go:build linux && cgo && !android && !server

package application

import (
	"github.com/godbus/dbus/v5"
	"github.com/wailsapp/wails/v3/pkg/events"
)

const (
	portalBusName       = "org.freedesktop.portal.Desktop"
	portalObjectPath    = "/org/freedesktop/portal/desktop"
	portalSettingsIface = "org.freedesktop.portal.Settings"

	// appearanceNamespace is the standardised namespace every portal
	// implementation publishes; it is the only one read or watched. GNOME also
	// mirrors the preference under org.gnome.desktop.interface as a string, but
	// honouring that in the signal filter alone would be inert: the handler
	// resolves through portalColorScheme, which speaks only this namespace.
	appearanceNamespace = "org.freedesktop.appearance"
	colorSchemeKey      = "color-scheme"

	colorSchemePreferDark = 1
)

// isDarkMode reports the desktop colour-scheme preference, read from the
// freedesktop Settings portal on every call.
//
// Read on demand rather than served from state maintained by
// monitorThemeChanges: that monitor is started per backend, so any cache it
// owns is only as correct as its startup wiring, and a cache fed from
// SettingChanged payloads reports light on every desktop until the first signal
// arrives. An on-demand read is right whether or not the monitor is running.
// Callers that need this on a hot path should cache it themselves.
func (a *linuxApp) isDarkMode() bool {
	scheme, ok := portalColorScheme()
	return ok && scheme == colorSchemePreferDark
}

// portalColorScheme reads org.freedesktop.appearance color-scheme: 0 is no
// preference, 1 prefers dark, 2 prefers light. ok is false when the portal is
// unreachable or the value is not the documented type.
func portalColorScheme() (uint32, bool) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return 0, false
	}

	obj := conn.Object(portalBusName, portalObjectPath)
	call := obj.Call(portalSettingsIface+".Read", 0, appearanceNamespace, colorSchemeKey)
	if call.Err != nil {
		return 0, false
	}

	var outer dbus.Variant
	if err := call.Store(&outer); err != nil {
		return 0, false
	}
	// Portal v1 Read double-wraps the value; other implementations, and ReadOne,
	// return it singly wrapped. Accept both, because rejecting one shape here
	// silently reports light -- the failure this whole path exists to remove.
	if inner, ok := outer.Value().(dbus.Variant); ok {
		scheme, ok := inner.Value().(uint32)
		return scheme, ok
	}
	scheme, ok := outer.Value().(uint32)
	return scheme, ok
}

// isColorSchemeChange reports whether a signal is a colour-scheme
// SettingChanged in the standardised appearance namespace.
func isColorSchemeChange(sig *dbus.Signal) bool {
	if sig.Name != portalSettingsIface+".SettingChanged" {
		return false
	}
	if len(sig.Body) < 2 {
		return false
	}
	namespace, _ := sig.Body[0].(string)
	key, _ := sig.Body[1].(string)
	return namespace == appearanceNamespace && key == colorSchemeKey
}

// monitorThemeChanges emits Linux.SystemThemeChanged when the desktop colour
// scheme changes. The portal is re-read rather than the signal payload trusted,
// so the emitted value always agrees with isDarkMode regardless of which
// namespace fired and what type it carried.
func (a *linuxApp) monitorThemeChanges() {
	go func() {
		defer handlePanic()
		conn, err := dbus.ConnectSessionBus()
		if err != nil {
			a.parent.warning(
				"[WARNING] Failed to connect to session bus; monitoring for theme changes will not function: %v",
				err,
			)
			return
		}
		defer conn.Close()

		if err = conn.AddMatchSignal(
			dbus.WithMatchSender(portalBusName),
			dbus.WithMatchObjectPath(portalObjectPath),
			dbus.WithMatchInterface(portalSettingsIface),
			dbus.WithMatchMember("SettingChanged"),
		); err != nil {
			a.parent.warning(
				"[WARNING] Failed to subscribe to portal SettingChanged; theme changes will not fire: %v",
				err,
			)
			return
		}

		c := make(chan *dbus.Signal, 10)
		conn.Signal(c)

		last := a.isDarkMode()
		for v := range c {
			if !isColorSchemeChange(v) {
				continue
			}
			dark := a.isDarkMode()
			if dark == last {
				continue
			}
			last = dark

			event := newApplicationEvent(events.Linux.SystemThemeChanged)
			event.Context().setIsDarkMode(dark)
			applicationEvents <- event
		}
	}()
}

// monitorPowerEvents subscribes to systemd-logind's PrepareForSleep signal on
// the system bus and translates it into Linux.SystemWillSleep (arg=true, just
// before suspend) and Linux.SystemDidWake (arg=false, immediately on resume).
// Mirrors NSWorkspace willSleep/didWake on macOS and WM_POWERBROADCAST on
// Windows.
//
// On systems without systemd or logind/elogind reachable on the system bus
// (Alpine, Void, some Devuan setups), we log a warning and exit cleanly so
// the rest of the app keeps working.
func (a *linuxApp) monitorPowerEvents() {
	go func() {
		defer handlePanic()
		conn, err := dbus.ConnectSystemBus()
		if err != nil {
			a.parent.warning(
				"[WARNING] Failed to connect to system bus; sleep/wake events will not fire: %v",
				err,
			)
			return
		}
		defer conn.Close()

		// Probe for logind/elogind ownership of org.freedesktop.login1 on the
		// system bus. Without this check, AddMatchSignal would succeed on any
		// systemd-less distro and the goroutine would block forever on a
		// channel that never receives — silently masking the missing service.
		var hasOwner bool
		if err := conn.BusObject().Call(
			"org.freedesktop.DBus.NameHasOwner", 0, "org.freedesktop.login1",
		).Store(&hasOwner); err != nil {
			a.parent.warning(
				"[WARNING] Failed to probe org.freedesktop.login1; sleep/wake events will not fire: %v",
				err,
			)
			return
		}
		if !hasOwner {
			a.parent.warning(
				"[WARNING] systemd-logind/elogind not reachable on the system bus; sleep/wake events will not fire",
			)
			return
		}

		// Constrain the sender to logind's well-known name so a hostile
		// connection on the system bus can't spoof PrepareForSleep signals.
		if err = conn.AddMatchSignal(
			dbus.WithMatchSender("org.freedesktop.login1"),
			dbus.WithMatchInterface("org.freedesktop.login1.Manager"),
			dbus.WithMatchMember("PrepareForSleep"),
			dbus.WithMatchObjectPath("/org/freedesktop/login1"),
		); err != nil {
			a.parent.warning(
				"[WARNING] Failed to subscribe to logind PrepareForSleep; sleep/wake events will not fire: %v",
				err,
			)
			return
		}

		c := make(chan *dbus.Signal, 4)
		conn.Signal(c)

		for v := range c {
			if v.Name != "org.freedesktop.login1.Manager.PrepareForSleep" {
				continue
			}
			if len(v.Body) < 1 {
				continue
			}
			willSleep, ok := v.Body[0].(bool)
			if !ok {
				continue
			}
			if willSleep {
				applicationEvents <- newApplicationEvent(events.Linux.SystemWillSleep)
			} else {
				applicationEvents <- newApplicationEvent(events.Linux.SystemDidWake)
			}
		}
	}()
}
