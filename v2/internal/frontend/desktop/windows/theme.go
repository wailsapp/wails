//go:build windows

package windows

import (
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/win32"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

func (w *Window) UpdateTheme() {

	// Don't redraw theme if nothing has changed
	if !w.themeChanged {
		return
	}
	w.themeChanged = false

	if win32.IsCurrentlyHighContrastMode() {
		return
	}

	if !win32.SupportsThemes() {
		return
	}

	var isDarkMode bool
	switch w.theme {
	case windows.SystemDefault:
		isDarkMode = win32.IsCurrentlyDarkMode()
	case windows.Dark:
		isDarkMode = true
	case windows.Light:
		isDarkMode = false
	}
	win32.SetTheme(w.Handle(), isDarkMode)

	// Custom theme processing
	winOptions := w.frontendOptions.Windows
	var customTheme *windows.ThemeSettings
	if winOptions != nil {
		customTheme = winOptions.CustomTheme
	}
	// Custom theme
	if win32.SupportsCustomThemes() && customTheme != nil {
		if w.isActive {
			if isDarkMode {
				win32.SetTitleBarColour(w.Handle(), customTheme.DarkModeTitleBar)
				win32.SetTitleTextColour(w.Handle(), customTheme.DarkModeTitleText)
				win32.SetBorderColour(w.Handle(), customTheme.DarkModeBorder)
			} else {
				win32.SetTitleBarColour(w.Handle(), customTheme.LightModeTitleBar)
				win32.SetTitleTextColour(w.Handle(), customTheme.LightModeTitleText)
				win32.SetBorderColour(w.Handle(), customTheme.LightModeBorder)
			}
		} else {
			if isDarkMode {
				win32.SetTitleBarColour(w.Handle(), customTheme.DarkModeTitleBarInactive)
				win32.SetTitleTextColour(w.Handle(), customTheme.DarkModeTitleTextInactive)
				win32.SetBorderColour(w.Handle(), customTheme.DarkModeBorderInactive)
			} else {
				win32.SetTitleBarColour(w.Handle(), customTheme.LightModeTitleBarInactive)
				win32.SetTitleTextColour(w.Handle(), customTheme.LightModeTitleTextInactive)
				win32.SetBorderColour(w.Handle(), customTheme.LightModeBorderInactive)
			}
		}
	}
}
