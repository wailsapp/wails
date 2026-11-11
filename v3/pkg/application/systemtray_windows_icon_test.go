//go:build windows && !server

package application

import (
	"testing"

	"github.com/wailsapp/wails/v3/pkg/icons"
)

// A SetIcon made on the main thread is applied on a goroutine, so an older
// call can reach iconLock after a newer one. The older one must not win.
func TestWindowsSystemTrayDropsAnOlderIconUpdate(t *testing.T) {
	s := &windowsSystemTray{parent: &SystemTray{}}

	s.applyIcon(icons.SystrayLight, 2)
	newest := s.lightModeIcon
	if newest == 0 {
		t.Fatal("the newest icon update was not applied")
	}

	s.applyIcon(icons.SystrayDark, 1)
	if s.lightModeIcon != newest {
		t.Fatal("an older icon update replaced a newer one")
	}
	if s.lightApplied != 2 {
		t.Fatalf("lightApplied = %d, want 2", s.lightApplied)
	}
}

// The same rule for the dark mode icon, which has its own sequence.
func TestWindowsSystemTrayDropsAnOlderDarkModeIconUpdate(t *testing.T) {
	s := &windowsSystemTray{parent: &SystemTray{}}

	s.applyDarkModeIcon(icons.SystrayDark, 2)
	newest := s.darkModeIcon
	s.applyDarkModeIcon(icons.SystrayLight, 1)
	if s.darkModeIcon != newest || s.darkApplied != 2 {
		t.Fatalf("an older dark mode icon update replaced a newer one (applied=%d)", s.darkApplied)
	}
}
