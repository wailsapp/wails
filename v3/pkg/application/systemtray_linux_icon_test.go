//go:build linux && !android && !server

package application

import "testing"

// SetIconName and SetIconThemePath called before the tray runs are stored on
// the SystemTray, copied into the Linux impl, and published as the initial
// StatusNotifierItem properties when run exports them.
func TestLinuxSystemTrayIconNameBeforeRun(t *testing.T) {
	tray := newSystemTray(1)
	tray.SetIconName("app-connected-symbolic").SetIconThemePath("/opt/app/icons")

	impl := newSystemTrayImpl(tray).(*linuxSystemTray)
	props := impl.createPropSpec()["org.kde.StatusNotifierItem"]

	if got := props["IconName"].Value; got != "app-connected-symbolic" {
		t.Errorf("IconName = %q, want %q", got, "app-connected-symbolic")
	}
	if got := props["IconThemePath"].Value; got != "/opt/app/icons" {
		t.Errorf("IconThemePath = %q, want %q", got, "/opt/app/icons")
	}
	if _, ok := props["IconPixmap"]; !ok {
		t.Error("IconPixmap missing: hosts that ignore IconName still draw it")
	}
}

// Without SetIconName the properties stay empty, so hosts keep drawing
// IconPixmap exactly as before.
func TestLinuxSystemTrayIconNameDefaultsEmpty(t *testing.T) {
	impl := newSystemTrayImpl(newSystemTray(1)).(*linuxSystemTray)
	props := impl.createPropSpec()["org.kde.StatusNotifierItem"]

	if got := props["IconName"].Value; got != "" {
		t.Errorf("IconName = %q, want empty", got)
	}
	if got := props["IconThemePath"].Value; got != "" {
		t.Errorf("IconThemePath = %q, want empty", got)
	}
}

// Before run there are no exported properties or bus connection; the setters
// must record the value for run to publish rather than dereference either.
func TestLinuxSystemTrayIconNameSettersBeforeExport(t *testing.T) {
	impl := &linuxSystemTray{parent: &SystemTray{}}

	impl.setIconName("app-warning-symbolic")
	impl.setIconThemePath("/opt/app/icons")

	if impl.iconName != "app-warning-symbolic" || impl.iconThemePath != "/opt/app/icons" {
		t.Errorf("got (%q, %q), want (%q, %q)",
			impl.iconName, impl.iconThemePath, "app-warning-symbolic", "/opt/app/icons")
	}
}
