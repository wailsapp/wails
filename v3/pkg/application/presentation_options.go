package application

import (
	"errors"
	"fmt"
	"strings"
)

// This file is the cross-platform surface of the macOS presentation options
// (NSApplication.presentationOptions): the bitmask type, its validation and
// the App methods that apply it at runtime. The platform functions it calls
// live in presentation_options_darwin.go, with documented no-op stubs in
// presentation_options_other.go. MacOptions.PresentationOptions applies a
// value when the application finishes launching.

// MacPresentationOptions is a bitmask mirroring NSApplicationPresentationOptions.
// It controls the Dock, the menu bar and the system controls (process
// switching, force quit, log out, hiding) while the application is active.
// The zero value, MacPresentationDefault, keeps the standard behaviour.
//
// AppKit rejects some combinations; Validate reports them so they can be
// caught before they reach the native setter, which would otherwise raise:
//   - MacPresentationAutoHideDock and MacPresentationHideDock are exclusive;
//   - MacPresentationAutoHideMenuBar and MacPresentationHideMenuBar are
//     exclusive, and either one requires AutoHideDock or HideDock;
//   - MacPresentationAutoHideToolbar requires both MacPresentationFullScreen
//     and MacPresentationAutoHideMenuBar.
type MacPresentationOptions uint

const (
	// MacPresentationDefault keeps the standard Dock, menu bar and system
	// controls.
	MacPresentationDefault MacPresentationOptions = 0
	// MacPresentationAutoHideDock hides the Dock until the pointer reaches
	// its edge of the screen.
	MacPresentationAutoHideDock MacPresentationOptions = 1 << 0
	// MacPresentationHideDock hides the Dock entirely.
	MacPresentationHideDock MacPresentationOptions = 1 << 1
	// MacPresentationAutoHideMenuBar hides the menu bar until the pointer
	// reaches the top of the screen. Requires AutoHideDock or HideDock.
	MacPresentationAutoHideMenuBar MacPresentationOptions = 1 << 2
	// MacPresentationHideMenuBar hides the menu bar entirely. Requires
	// AutoHideDock or HideDock.
	MacPresentationHideMenuBar MacPresentationOptions = 1 << 3
	// MacPresentationDisableAppleMenu disables every item in the Apple menu.
	MacPresentationDisableAppleMenu MacPresentationOptions = 1 << 4
	// MacPresentationDisableProcessSwitching disables Command-Tab and the
	// Dock as ways to switch to another application.
	MacPresentationDisableProcessSwitching MacPresentationOptions = 1 << 5
	// MacPresentationDisableForceQuit disables the Force Quit window
	// (Command-Option-Escape).
	MacPresentationDisableForceQuit MacPresentationOptions = 1 << 6
	// MacPresentationDisableSessionTermination disables log out, restart and
	// shut down while the application is active.
	MacPresentationDisableSessionTermination MacPresentationOptions = 1 << 7
	// MacPresentationDisableHideApplication disables the Hide command.
	MacPresentationDisableHideApplication MacPresentationOptions = 1 << 8
	// MacPresentationDisableMenuBarTransparency draws an opaque menu bar.
	MacPresentationDisableMenuBarTransparency MacPresentationOptions = 1 << 9
	// MacPresentationFullScreen marks the application as being in
	// application-level full screen mode (the kiosk style used before
	// per-window full screen). It is distinct from a window's own full
	// screen state.
	MacPresentationFullScreen MacPresentationOptions = 1 << 10
	// MacPresentationAutoHideToolbar detaches the full-screen window toolbar
	// so it hides and shows with the menu bar. Requires FullScreen and
	// AutoHideMenuBar.
	MacPresentationAutoHideToolbar MacPresentationOptions = 1 << 11
	// MacPresentationDisableCursorLocationAssistance disables the shake to
	// locate cursor feature (macOS 10.11.2+).
	MacPresentationDisableCursorLocationAssistance MacPresentationOptions = 1 << 12
)

// macPresentationKnown is every bit AppKit defines.
const macPresentationKnown = MacPresentationAutoHideDock | MacPresentationHideDock |
	MacPresentationAutoHideMenuBar | MacPresentationHideMenuBar |
	MacPresentationDisableAppleMenu | MacPresentationDisableProcessSwitching |
	MacPresentationDisableForceQuit | MacPresentationDisableSessionTermination |
	MacPresentationDisableHideApplication | MacPresentationDisableMenuBarTransparency |
	MacPresentationFullScreen | MacPresentationAutoHideToolbar |
	MacPresentationDisableCursorLocationAssistance

// ErrMacPresentationOptionsInvalid is wrapped by the errors Validate returns
// for combinations AppKit rejects.
var ErrMacPresentationOptionsInvalid = errors.New("invalid macOS presentation options")

// Has reports whether every bit in flag is set.
func (o MacPresentationOptions) Has(flag MacPresentationOptions) bool {
	return o&flag == flag
}

// Validate reports the first combination AppKit would reject, or nil when
// the options can be applied. Unknown bits are rejected too.
func (o MacPresentationOptions) Validate() error {
	if unknown := o &^ macPresentationKnown; unknown != 0 {
		return fmt.Errorf("%w: unknown bits %#x", ErrMacPresentationOptionsInvalid, uint(unknown))
	}
	if o.Has(MacPresentationAutoHideDock) && o.Has(MacPresentationHideDock) {
		return fmt.Errorf("%w: AutoHideDock and HideDock are mutually exclusive", ErrMacPresentationOptionsInvalid)
	}
	if o.Has(MacPresentationAutoHideMenuBar) && o.Has(MacPresentationHideMenuBar) {
		return fmt.Errorf("%w: AutoHideMenuBar and HideMenuBar are mutually exclusive", ErrMacPresentationOptionsInvalid)
	}
	hidesMenuBar := o.Has(MacPresentationAutoHideMenuBar) || o.Has(MacPresentationHideMenuBar)
	hidesDock := o.Has(MacPresentationAutoHideDock) || o.Has(MacPresentationHideDock)
	if hidesMenuBar && !hidesDock {
		return fmt.Errorf("%w: hiding the menu bar requires AutoHideDock or HideDock", ErrMacPresentationOptionsInvalid)
	}
	if o.Has(MacPresentationAutoHideToolbar) && !(o.Has(MacPresentationFullScreen) && o.Has(MacPresentationAutoHideMenuBar)) {
		return fmt.Errorf("%w: AutoHideToolbar requires FullScreen and AutoHideMenuBar", ErrMacPresentationOptionsInvalid)
	}
	return nil
}

// String lists the set flags by name, or "Default" for the zero value.
func (o MacPresentationOptions) String() string {
	if o == MacPresentationDefault {
		return "Default"
	}
	names := []struct {
		flag MacPresentationOptions
		name string
	}{
		{MacPresentationAutoHideDock, "AutoHideDock"},
		{MacPresentationHideDock, "HideDock"},
		{MacPresentationAutoHideMenuBar, "AutoHideMenuBar"},
		{MacPresentationHideMenuBar, "HideMenuBar"},
		{MacPresentationDisableAppleMenu, "DisableAppleMenu"},
		{MacPresentationDisableProcessSwitching, "DisableProcessSwitching"},
		{MacPresentationDisableForceQuit, "DisableForceQuit"},
		{MacPresentationDisableSessionTermination, "DisableSessionTermination"},
		{MacPresentationDisableHideApplication, "DisableHideApplication"},
		{MacPresentationDisableMenuBarTransparency, "DisableMenuBarTransparency"},
		{MacPresentationFullScreen, "FullScreen"},
		{MacPresentationAutoHideToolbar, "AutoHideToolbar"},
		{MacPresentationDisableCursorLocationAssistance, "DisableCursorLocationAssistance"},
	}
	var parts []string
	for _, entry := range names {
		if o.Has(entry.flag) {
			parts = append(parts, entry.name)
		}
	}
	if unknown := o &^ macPresentationKnown; unknown != 0 {
		parts = append(parts, fmt.Sprintf("unknown(%#x)", uint(unknown)))
	}
	return strings.Join(parts, "|")
}

// SetPresentationOptions replaces the application's presentation options
// (NSApplication.presentationOptions) while it is running. Invalid
// combinations are rejected with an error wrapping
// ErrMacPresentationOptionsInvalid before anything reaches AppKit. Pass
// MacPresentationDefault to restore the standard Dock and menu bar. Off
// macOS it returns ErrMacOnly.
func (a *App) SetPresentationOptions(options MacPresentationOptions) error {
	if err := options.Validate(); err != nil {
		return err
	}
	if !macPresentationOptionsSupported() {
		return ErrMacOnly
	}
	if a == nil || a.impl == nil {
		return errors.New("the application has not started")
	}
	return macSetPresentationOptions(options)
}

// PresentationOptions returns the application's current presentation
// options. It is MacPresentationDefault off macOS and before the
// application starts.
func (a *App) PresentationOptions() MacPresentationOptions {
	if a == nil || a.impl == nil || !macPresentationOptionsSupported() {
		return MacPresentationDefault
	}
	return macPresentationOptions()
}
