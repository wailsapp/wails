//go:build linux && !android && !server

package application

import (
	"testing"

	"github.com/godbus/dbus/v5"
)

// Regression test for #5344 — the filter must reject every signal that
// isn't NameOwnerChanged with a 3-element body before any Body[2] indexing.
func TestShouldReregisterOnNameOwnerChanged(t *testing.T) {
	cases := []struct {
		name string
		sig  *dbus.Signal
		want bool
	}{
		{
			name: "nil signal",
			sig:  nil,
			want: false,
		},
		{
			name: "Notifications.ActionInvoked (2 args, observed in #5344)",
			sig: &dbus.Signal{
				Name: "org.freedesktop.Notifications.ActionInvoked",
				Body: []interface{}{uint32(0x3d), "default"},
			},
			want: false,
		},
		{
			name: "Notifications.NotificationClosed (2 args, observed in #5344)",
			sig: &dbus.Signal{
				Name: "org.freedesktop.Notifications.NotificationClosed",
				Body: []interface{}{uint32(0x3d), uint32(2)},
			},
			want: false,
		},
		{
			name: "Local.Disconnected (0 args)",
			sig: &dbus.Signal{
				Name: "org.freedesktop.DBus.Local.Disconnected",
				Body: []interface{}{},
			},
			want: false,
		},
		{
			name: "NameAcquired (1 arg, similar matching path)",
			sig: &dbus.Signal{
				Name: "org.freedesktop.DBus.NameAcquired",
				Body: []interface{}{":1.42"},
			},
			want: false,
		},
		{
			name: "NameOwnerChanged with malformed short body (1 arg)",
			sig: &dbus.Signal{
				Name: dbusNameOwnerChanged,
				Body: []interface{}{dbusStatusNotifierWatcher},
			},
			want: false,
		},
		{
			// Pins the length-guard boundary: weakening `< 3` to `< 2` panics here.
			name: "NameOwnerChanged with malformed short body (2 args, boundary)",
			sig: &dbus.Signal{
				Name: dbusNameOwnerChanged,
				Body: []interface{}{dbusStatusNotifierWatcher, ":1.42"},
			},
			want: false,
		},
		{
			// Pins the type assertion: non-string Body[2] must not register.
			name: "NameOwnerChanged with non-string new owner (spec violation)",
			sig: &dbus.Signal{
				Name: dbusNameOwnerChanged,
				Body: []interface{}{dbusStatusNotifierWatcher, "", uint32(42)},
			},
			want: false,
		},
		{
			name: "NameOwnerChanged with empty new owner (watcher disappeared)",
			sig: &dbus.Signal{
				Name: dbusNameOwnerChanged,
				Body: []interface{}{dbusStatusNotifierWatcher, ":1.42", ""},
			},
			want: false,
		},
		{
			name: "NameOwnerChanged with populated new owner (watcher restarted)",
			sig: &dbus.Signal{
				Name: dbusNameOwnerChanged,
				Body: []interface{}{dbusStatusNotifierWatcher, "", ":1.99"},
			},
			want: true,
		},
		{
			// Pins the Body[0] guard: a NameOwnerChanged for a different
			// bus name (e.g. some unrelated process taking a name) must
			// not trigger a re-register on our watcher.
			name: "NameOwnerChanged for a different bus name (must NOT trigger register)",
			sig: &dbus.Signal{
				Name: dbusNameOwnerChanged,
				Body: []interface{}{"org.freedesktop.Notifications", "", ":1.50"},
			},
			want: false,
		},
		{
			// Pins the Body[0] type assertion: a non-string watched-name
			// (spec violation) must not trigger.
			name: "NameOwnerChanged with non-string Body[0] (spec violation)",
			sig: &dbus.Signal{
				Name: dbusNameOwnerChanged,
				Body: []interface{}{uint32(1), "", ":1.50"},
			},
			want: false,
		},
		{
			// Pins the name guard: PropertiesChanged has the same 3-arg
			// shape, so removing the name filter flips this case to true.
			name: "PropertiesChanged (3 args, non-empty Body[2]) must NOT trigger register",
			sig: &dbus.Signal{
				Name: "org.freedesktop.DBus.Properties.PropertiesChanged",
				Body: []interface{}{
					"com.canonical.dbusmenu",
					map[string]dbus.Variant{},
					[]string{"Version"},
				},
			},
			want: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := shouldReregisterOnNameOwnerChanged(tc.sig); got != tc.want {
				t.Errorf("shouldReregisterOnNameOwnerChanged(%+v) = %v, want %v",
					tc.sig, got, tc.want)
			}
		})
	}
}
