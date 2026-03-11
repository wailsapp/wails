package application

import "github.com/wailsapp/wails/v3/pkg/events"

// WinTheme represents the theme preference for a window.
type WinTheme string

const (
	// WinAppDefault indicates the window should follow the application theme.
	WinAppDefault WinTheme = "application"
	// WinDark forces the window to use a dark theme.
	WinDark WinTheme = "dark"
	// WinLight forces the window to use a light theme.
	WinLight WinTheme = "light"
	// WinSystemDefault indicates the window should follow the system theme.
	WinSystemDefault WinTheme = "system"
)

// String returns the string representation of the window theme.
func (t WinTheme) String() string {
	return string(t)
}

// Valid returns true if the theme is a recognized WinTheme value.
func (t WinTheme) Valid() bool {
	switch t {
	case WinAppDefault, WinDark, WinLight, WinSystemDefault:
		return true
	}
	return false
}

// GetTheme returns the current theme of the window.
func (w *WebviewWindow) GetTheme() WinTheme {
	if w.impl == nil {
		return WinAppDefault
	}
	return w.impl.getTheme()
}

// SetTheme sets the theme for the window.
func (w *WebviewWindow) SetTheme(theme WinTheme) {
	if !theme.Valid() {
		return
	}
	if w.impl != nil {
		w.impl.setTheme(theme)
		// Notify listeners of the theme change
		w.emit(events.WindowEventType(events.Common.ThemeChanged))
	}
}
