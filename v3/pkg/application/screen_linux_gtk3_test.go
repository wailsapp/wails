//go:build linux && cgo && gtk3 && !android && !server

package application

import (
	"errors"
	"testing"
)

func TestGetScreensWithoutActiveWindow(t *testing.T) {
	// A service-only application has no active GTK window. Screen discovery must
	// either use the default display or report that no display is available.
	screens, err := getScreens(nil)
	if err != nil && !errors.Is(err, errScreenDisplayUnavailable) {
		t.Fatalf("getScreens(nil) returned an unexpected error: %v", err)
	}
	if err == nil && len(screens) == 0 {
		t.Fatal("getScreens(nil) returned neither screens nor an error")
	}
}

func TestGetScreensForDisplayRejectsNil(t *testing.T) {
	screens, err := getScreensForDisplay(nil)
	if !errors.Is(err, errScreenDisplayUnavailable) {
		t.Fatalf("getScreensForDisplay(nil) error = %v, want %v", err, errScreenDisplayUnavailable)
	}
	if screens != nil {
		t.Fatalf("getScreensForDisplay(nil) screens = %v, want nil", screens)
	}
}
