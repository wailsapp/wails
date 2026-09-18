//go:build linux && cgo && !android && !server

package application

import (
	"github.com/godbus/dbus/v5"
	"github.com/wailsapp/wails/v3/pkg/events"
)

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
			dbus.WithMatchObjectPath("/org/freedesktop/portal/desktop"),
		); err != nil {
			a.parent.warning(
				"[WARNING] Failed to subscribe to portal SettingChanged; theme changes will not fire: %v",
				err,
			)
			return
		}

		c := make(chan *dbus.Signal, 10)
		conn.Signal(c)

		getTheme := func(body []interface{}) (string, bool) {
			if len(body) < 3 {
				return "", false
			}
			if entry, ok := body[0].(string); !ok || entry != "org.gnome.desktop.interface" {
				return "", false
			}
			if entry, ok := body[1].(string); !ok || entry != "color-scheme" {
				return "", false
			}
			variant, ok := body[2].(dbus.Variant)
			if !ok {
				return "", false
			}
			value, ok := variant.Value().(string)
			if !ok {
				return "", false
			}
			return value, true
		}

		for v := range c {
			theme, ok := getTheme(v.Body)
			if !ok {
				continue
			}

			if theme != a.theme {
				a.theme = theme
				event := newApplicationEvent(events.Linux.SystemThemeChanged)
				event.Context().setIsDarkMode(a.isDarkMode())
				applicationEvents <- event
			}

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
