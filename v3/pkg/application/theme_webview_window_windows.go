//go:build windows

package application

import "github.com/wailsapp/wails/v3/pkg/w32"

// resolveWindowsEffectiveTheme determines the realized Theme for the window by resolving
// application-level and window-level theme settings. It also returns whether the window follows the application theme.
func resolveWindowsEffectiveTheme(winTheme WinTheme, appTheme AppTheme) (theme, bool) {
	switch winTheme {
	case WinDark:
		return dark, false
	case WinLight:
		return light, false
	case WinSystemDefault:
		return systemDefault, false
	default:
		// For WinThemeApplication and/or Unset values we default to following
		switch appTheme {
		case AppDark:
			return dark, true
		case AppLight:
			return light, true
		case AppSystemDefault:
			return systemDefault, true
		default:
			return systemDefault, true
		}
	}
}

// syncTheme synchronizes the window's appearance with the application-wide theme,
// assuming the window is configured to follow the application theme.
// Theme updates are expected to run on the UI thread.
// SystemThemeChanged events dispatch via InvokeAsync, ensuring
// that window theme state is mutated from a single thread.
// But if required, Mutex can be added to make sure w.theme does not
// cause any Race condition.
func (w *windowsWebviewWindow) syncTheme() {
	if !w.parent.followApplicationTheme {
		return
	}

	switch globalApplication.theme {
	case AppSystemDefault:
		w.theme = systemDefault
		w.updateTheme(w32.IsCurrentlyDarkMode())
	case AppDark:
		if w.theme != dark {
			w.theme = dark
			w32.AllowDarkModeForWindow(w.hwnd, true)
			w.updateTheme(true)
		}
	case AppLight:
		if w.theme != light {
			w.theme = light
			w.updateTheme(false)
		}
	}
}

// setTheme sets the theme for the window. If WinThemeApplication is provided,
// the window will follow the application-wide theme settings.
func (w *windowsWebviewWindow) setTheme(theme WinTheme) {
	if theme == WinAppDefault {
		w.parent.followApplicationTheme = true
		w.syncTheme()
		return
	}

	w.parent.followApplicationTheme = false
	switch theme {
	case WinDark:
		w.theme = dark
		w.updateTheme(true)
	case WinLight:
		w.theme = light
		w.updateTheme(false)
	case WinSystemDefault:
		w.theme = systemDefault
		w.updateTheme(w32.IsCurrentlyDarkMode())
	default:
		w.theme = systemDefault
		w.updateTheme(w32.IsCurrentlyDarkMode())
	}
}

// getTheme returns the current theme configuration for the window.
func (w *windowsWebviewWindow) getTheme() WinTheme {
	if w.parent.followApplicationTheme {
		return WinAppDefault
	}

	if w.theme == systemDefault {
		return WinSystemDefault
	}

	if w.theme == dark {
		return WinDark
	}

	return WinLight
}
