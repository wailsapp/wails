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
// Must be called on the UI thread.
func (w *windowsWebviewWindow) syncTheme() {
	if !w.parent.followApplicationTheme {
		return
	}

	switch globalApplication.theme {
	case AppSystemDefault:
		w.themeMu.Lock()
		w.theme = systemDefault
		w.themeMu.Unlock()
		w.updateTheme(w32.IsCurrentlyDarkMode())
	case AppDark:
		w.themeMu.Lock()
		changed := w.theme != dark
		if changed {
			w.theme = dark
		}
		w.themeMu.Unlock()
		if changed {
			w32.AllowDarkModeForWindow(w.hwnd, true)
			w.updateTheme(true)
		}
	case AppLight:
		w.themeMu.Lock()
		changed := w.theme != light
		if changed {
			w.theme = light
		}
		w.themeMu.Unlock()
		if changed {
			w32.AllowDarkModeForWindow(w.hwnd, false)
			w.updateTheme(false)
		}
	}
}

// setTheme sets the theme for the window. If WinAppDefault is provided,
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
		w.themeMu.Lock()
		w.theme = dark
		w.themeMu.Unlock()
		w32.AllowDarkModeForWindow(w.hwnd, true)
		w.updateTheme(true)
	case WinLight:
		w.themeMu.Lock()
		w.theme = light
		w.themeMu.Unlock()
		w32.AllowDarkModeForWindow(w.hwnd, false)
		w.updateTheme(false)
	case WinSystemDefault:
		w.themeMu.Lock()
		w.theme = systemDefault
		w.themeMu.Unlock()
		w.updateTheme(w32.IsCurrentlyDarkMode())
	default:
		w.themeMu.Lock()
		w.theme = systemDefault
		w.themeMu.Unlock()
		w.updateTheme(w32.IsCurrentlyDarkMode())
	}
}

// getTheme returns the current theme configuration for the window.
func (w *windowsWebviewWindow) getTheme() WinTheme {
	if w.parent.followApplicationTheme {
		return WinAppDefault
	}

	w.themeMu.RLock()
	t := w.theme
	w.themeMu.RUnlock()

	switch t {
	case dark:
		return WinDark
	case light:
		return WinLight
	default:
		return WinSystemDefault
	}
}
