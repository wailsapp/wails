package application

import "github.com/wailsapp/wails/v3/pkg/events"

// WinTheme represents the theme preference for a window.
type WinTheme string

const (
	// WinThemeApplication indicates the window should follow the application theme.
	WinThemeApplication WinTheme = "application"
	// WinThemeDark forces the window to use a dark theme.
	WinThemeDark WinTheme = "dark"
	// WinThemeLight forces the window to use a light theme.
	WinThemeLight WinTheme = "light"
	// WinThemeSystem indicates the window should follow the system theme.
	WinThemeSystem WinTheme = "system"
)

// String returns the string representation of the window theme.
func (t WinTheme) String() string {
	return string(t)
}

// Valid returns true if the theme is a recognized WinTheme value.
func (t WinTheme) Valid() bool {
	switch t {
	case WinThemeApplication, WinThemeDark, WinThemeLight, WinThemeSystem:
		return true
	}
	return false
}

// GetTheme returns the current theme of the window.
func (w *WebviewWindow) GetTheme() string {
	if w.impl == nil {
		return WinThemeApplication.String()
	}
	return w.impl.getTheme().String()
}

// SetTheme sets the theme for the window.
func (w *WebviewWindow) SetTheme(theme WinTheme) {
	if !theme.Valid() {
		return
	}
	if w.impl != nil {
		w.impl.setTheme(theme)
	}
	// Notify listeners of the theme change
	w.emit(events.WindowEventType(events.Common.ThemeChanged))
}
