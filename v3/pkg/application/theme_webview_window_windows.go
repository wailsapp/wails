//go:build windows

package application

import "github.com/wailsapp/wails/v3/pkg/w32"

// resolveWindowsEffectiveTheme determines the realized Theme for the window by resolving
// application-level and window-level theme settings.
func resolveWindowsEffectiveTheme(winTheme WinTheme, appTheme AppTheme) Theme {
	switch winTheme {
	case WinThemeDark:
		return Dark
	case WinThemeLight:
		return Light
	case WinThemeSystem:
		return SystemDefault
	default:
		// For WinThemeApplication and/or Unset values we default to following
		switch appTheme {
		case AppDark:
			return Dark
		case AppLight:
			return Light
		case AppSystemDefault:
			return SystemDefault
		default:
			return SystemDefault
		}
	}
}

// syncTheme synchronizes the window's appearance with the application-wide theme,
// assuming the window is configured to follow the application theme.
func (w *windowsWebviewWindow) syncTheme() {
	if !w.parent.followApplicationTheme {
		return
	}

	switch globalApplication.theme {
	case AppSystemDefault:
		w.theme = SystemDefault
		w.updateTheme(w32.IsCurrentlyDarkMode())
	case AppDark:
		if w.theme != Dark {
			w.theme = Dark
			w32.AllowDarkModeForWindow(w.hwnd, true)
			w.updateTheme(true)
		}
	case AppLight:
		if w.theme != Light {
			w.theme = Light
			w.updateTheme(false)
		}
	}
}

// setTheme sets the theme for the window. If WinThemeApplication is provided,
// the window will follow the application-wide theme settings.
func (w *windowsWebviewWindow) setTheme(theme WinTheme) {
	if theme == WinThemeApplication {
		w.parent.followApplicationTheme = true
		w.syncTheme()
		return
	}

	w.parent.followApplicationTheme = false
	switch theme {
	case WinThemeDark:
		w.theme = Dark
		w.updateTheme(true)
	case WinThemeLight:
		w.theme = Light
		w.updateTheme(false)
	case WinThemeSystem:
		w.theme = SystemDefault
		w.updateTheme(w32.IsCurrentlyDarkMode())
	default:
		w.theme = SystemDefault
		w.updateTheme(w32.IsCurrentlyDarkMode())
	}
}

// getTheme returns the current theme configuration for the window.
func (w *windowsWebviewWindow) getTheme() WinTheme {
	if w.parent.followApplicationTheme {
		return WinThemeApplication
	}

	if w.theme == SystemDefault {
		return WinThemeSystem
	}

	if w.theme == Dark {
		return WinThemeDark
	}

	return WinThemeLight
}
