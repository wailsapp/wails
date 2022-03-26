package windows

import (
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/win32"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

func (w *Window) updateTheme() {

	if win32.IsCurrentlyHighContrastMode() {
		return
	}

	if !win32.SupportsThemes() {
		return
	}

	// Only process if there's a theme change
	isDarkMode := win32.IsCurrentlyDarkMode()
	if w.isDarkMode == isDarkMode {
		return
	}
	w.isDarkMode = isDarkMode

	// Default use system theme
	winOptions := w.frontendOptions.Windows
	var customTheme *windows.ThemeSettings
	if winOptions != nil {
		customTheme = winOptions.CustomTheme
		if winOptions.Theme == windows.Dark {
			isDarkMode = true
		}
		if winOptions.Theme == windows.Light {
			isDarkMode = false
		}
	}

	win32.SetTheme(w.Handle(), isDarkMode)

	// Custom theme
	if win32.SupportsCustomThemes() && customTheme != nil {
		if isDarkMode {
			win32.SetTitleBarColour(w.Handle(), customTheme.DarkModeTitleBar)
			win32.SetTitleTextColour(w.Handle(), customTheme.DarkModeTitleText)
			win32.SetBorderColour(w.Handle(), customTheme.DarkModeBorder)
		} else {
			win32.SetTitleBarColour(w.Handle(), customTheme.LightModeTitleBar)
			win32.SetTitleTextColour(w.Handle(), customTheme.LightModeTitleText)
			win32.SetBorderColour(w.Handle(), customTheme.LightModeBorder)
		}
	}
}
