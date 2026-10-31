//go:build linux

package main

import (
	"os"
	"strings"
	"testing"
)

// These tests talk to a real X server, so they need DISPLAY (run them under
// xvfb-run). Without one they are skipped; the pure mapping tests in
// keys_test.go always run.
func requireDisplay(t *testing.T) *xDriver {
	t.Helper()
	if os.Getenv("DISPLAY") == "" {
		t.Skip("DISPLAY is empty: run under xvfb-run to test the X keymap")
	}
	d, err := newDriver()
	if err != nil {
		t.Fatalf("newDriver: %v", err)
	}
	return d.(*xDriver)
}

// The keymap decides Shift from the server's own table. On the default US
// keymap of Xvfb, capitals and the shifted punctuation need Shift, and the
// unshifted keys do not.
func TestKeymapShiftDecisions(t *testing.T) {
	d := requireDisplay(t)
	tests := []struct {
		keysym    string
		wantShift bool
	}{
		{"a", false}, {"A", true}, {"k", false}, {"K", true},
		{"slash", false}, {"question", true},
		{"period", false}, {"minus", false}, {"underscore", true},
		{"space", false}, {"Return", false},
	}
	for _, tc := range tests {
		_, shift, err := d.x.keycodeFor(tc.keysym)
		if err != nil {
			t.Errorf("keycodeFor(%q): %v", tc.keysym, err)
			continue
		}
		if shift != tc.wantShift {
			t.Errorf("keycodeFor(%q) shift = %v, want %v", tc.keysym, shift, tc.wantShift)
		}
	}
}

func TestKeycodeForUnknownKeysym(t *testing.T) {
	d := requireDisplay(t)
	if _, _, err := d.x.keycodeFor("NoSuchKeysym"); err == nil {
		t.Error("keycodeFor of an unknown keysym returned no error")
	}
}

func TestFindReportsMissingWindow(t *testing.T) {
	d := requireDisplay(t)
	_, err := d.x.find("no such window title for e2e")
	if err == nil || !strings.Contains(err.Error(), "no such window title for e2e") {
		t.Errorf("find of a missing title: error = %v", err)
	}
}

func TestMenuBarOffsetOverride(t *testing.T) {
	t.Setenv("E2E_MENUBAR_HEIGHT", "41")
	if got := menuBarOffset(); got != 41 {
		t.Errorf("menuBarOffset with E2E_MENUBAR_HEIGHT=41 = %d", got)
	}
	t.Setenv("E2E_MENUBAR_HEIGHT", "not a number")
	if got := menuBarOffset(); got != menuBarHeight {
		t.Errorf("menuBarOffset with a bad value = %d, want %d", got, menuBarHeight)
	}
}
