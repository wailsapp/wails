package application

import (
	"errors"
	"testing"
)

func TestMacPresentationOptionValuesMatchAppKit(t *testing.T) {
	// NSApplicationPresentationOptions bit positions, from NSApplication.h.
	want := map[MacPresentationOptions]uint{
		MacPresentationAutoHideDock:                    1 << 0,
		MacPresentationHideDock:                        1 << 1,
		MacPresentationAutoHideMenuBar:                 1 << 2,
		MacPresentationHideMenuBar:                     1 << 3,
		MacPresentationDisableAppleMenu:                1 << 4,
		MacPresentationDisableProcessSwitching:         1 << 5,
		MacPresentationDisableForceQuit:                1 << 6,
		MacPresentationDisableSessionTermination:       1 << 7,
		MacPresentationDisableHideApplication:          1 << 8,
		MacPresentationDisableMenuBarTransparency:      1 << 9,
		MacPresentationFullScreen:                      1 << 10,
		MacPresentationAutoHideToolbar:                 1 << 11,
		MacPresentationDisableCursorLocationAssistance: 1 << 12,
	}
	for flag, value := range want {
		if uint(flag) != value {
			t.Fatalf("%s = %#x, want %#x", flag, uint(flag), value)
		}
	}
	if MacPresentationDefault != 0 {
		t.Fatal("MacPresentationDefault must be the zero value")
	}
}

func TestMacPresentationOptionsValidateAcceptsDocumentedCombinations(t *testing.T) {
	valid := []MacPresentationOptions{
		MacPresentationDefault,
		MacPresentationAutoHideDock,
		MacPresentationHideDock,
		MacPresentationAutoHideDock | MacPresentationAutoHideMenuBar,
		MacPresentationHideDock | MacPresentationHideMenuBar,
		MacPresentationHideDock | MacPresentationAutoHideMenuBar | MacPresentationDisableProcessSwitching |
			MacPresentationDisableForceQuit | MacPresentationDisableSessionTermination |
			MacPresentationDisableHideApplication | MacPresentationDisableAppleMenu,
		MacPresentationFullScreen | MacPresentationAutoHideDock | MacPresentationAutoHideMenuBar | MacPresentationAutoHideToolbar,
		MacPresentationDisableMenuBarTransparency | MacPresentationDisableCursorLocationAssistance,
	}
	for _, options := range valid {
		if err := options.Validate(); err != nil {
			t.Fatalf("Validate(%s) = %v, want nil", options, err)
		}
	}
}

func TestMacPresentationOptionsValidateRejectsInvalidPairs(t *testing.T) {
	invalid := map[string]MacPresentationOptions{
		"both dock modes":                 MacPresentationAutoHideDock | MacPresentationHideDock,
		"both menu bar modes":             MacPresentationHideDock | MacPresentationAutoHideMenuBar | MacPresentationHideMenuBar,
		"auto hide menu bar without dock": MacPresentationAutoHideMenuBar,
		"hide menu bar without dock":      MacPresentationHideMenuBar | MacPresentationDisableForceQuit,
		"auto hide toolbar alone":         MacPresentationAutoHideToolbar,
		"auto hide toolbar without menu":  MacPresentationAutoHideToolbar | MacPresentationFullScreen,
		"auto hide toolbar without full":  MacPresentationAutoHideToolbar | MacPresentationAutoHideDock | MacPresentationAutoHideMenuBar,
		"unknown bit":                     MacPresentationOptions(1 << 20),
	}
	for name, options := range invalid {
		err := options.Validate()
		if !errors.Is(err, ErrMacPresentationOptionsInvalid) {
			t.Fatalf("%s: Validate(%s) = %v, want ErrMacPresentationOptionsInvalid", name, options, err)
		}
	}
}

func TestMacPresentationOptionsString(t *testing.T) {
	if got := MacPresentationDefault.String(); got != "Default" {
		t.Fatalf("Default.String() = %q", got)
	}
	options := MacPresentationHideDock | MacPresentationHideMenuBar | MacPresentationDisableForceQuit
	if got := options.String(); got != "HideDock|HideMenuBar|DisableForceQuit" {
		t.Fatalf("String() = %q", got)
	}
	if !options.Has(MacPresentationHideDock) || options.Has(MacPresentationAutoHideDock) {
		t.Fatal("Has reported the wrong flags")
	}
}

func TestSetPresentationOptionsValidatesBeforeThePlatform(t *testing.T) {
	var nilApp *App
	err := nilApp.SetPresentationOptions(MacPresentationAutoHideDock | MacPresentationHideDock)
	if !errors.Is(err, ErrMacPresentationOptionsInvalid) {
		t.Fatalf("invalid options on a nil app = %v, want ErrMacPresentationOptionsInvalid", err)
	}
	if got := nilApp.PresentationOptions(); got != MacPresentationDefault {
		t.Fatalf("PresentationOptions on a nil app = %s", got)
	}
	notStarted := &App{}
	err = notStarted.SetPresentationOptions(MacPresentationAutoHideDock)
	if macPresentationOptionsSupported() {
		if err == nil || errors.Is(err, ErrMacOnly) {
			t.Fatalf("SetPresentationOptions before start = %v, want a not-started error", err)
		}
	} else if !errors.Is(err, ErrMacOnly) {
		t.Fatalf("SetPresentationOptions off macOS = %v, want ErrMacOnly", err)
	}
	if got := notStarted.PresentationOptions(); got != MacPresentationDefault {
		t.Fatalf("PresentationOptions before start = %s", got)
	}
}
