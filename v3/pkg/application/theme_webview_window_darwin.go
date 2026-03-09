//go:build darwin

package application

import "fmt"

// getOppositeAppearance returns the macOS appearance that represents
// the opposite light/dark variant.
func (w *macosWebviewWindow) getOppositeAppearance(name string) (MacAppearanceType, error) {
	if name == "NSAppearanceNameDarkAqua" {
		return "NSAppearanceNameAqua", nil
	}

	// If opposite appearance doesnt match then send the default Dark Appearance
	err := fmt.Errorf("unknown appearance name: %s", name)
	return "NSAppearanceNameDarkAqua", err
}

// isAppearanceDark reports whether the current window appearance
// corresponds to a dark macOS appearance.
func (w *macosWebviewWindow) isAppearanceDark() bool {
	appr := w.getEffectiveAppearanceName()
	// Check if the appearance name contains "Dark"
	switch appr {
	case "NSAppearanceNameDarkAqua",
		"NSAppearanceNameVibrantDark",
		"NSAppearanceNameAccessibilityHighContrastDarkAqua",
		"NSAppearanceNameAccessibilityHighContrastVibrantDark":
		return true
	default:
		return false
	}
}

// syncTheme synchronizes the window's appearance with the application-wide theme
// when the window is configured to follow global application theme settings.
func (w *macosWebviewWindow) syncTheme() {
	if !w.parent.followApplicationTheme {
		return
	}

	currentAppearance := w.getEffectiveAppearanceName()
	currentDark := w.isAppearanceDark()

	switch globalApplication.theme {
	case AppSystemDefault:
		w.clearAppearance()
		return
	case AppDark:
		if !currentDark {
			appr, _ := w.getOppositeAppearance(currentAppearance)
			w.setAppearanceByName(appr)
		}
	case AppLight:
		if currentDark {
			appr, _ := w.getOppositeAppearance(currentAppearance)
			w.setAppearanceByName(appr)
		}
	}
}

// setTheme sets the theme for the window. If WinThemeApplication is provided,
// the window will follow global application theme settings.
func (w *macosWebviewWindow) setTheme(theme WinTheme) {
	switch theme {
	case WinThemeSystem:
		w.parent.followApplicationTheme = false
		w.clearAppearance()
		return
	case WinThemeApplication:
		w.parent.followApplicationTheme = true
		w.syncTheme()
		return
	}

	currentAppearance := w.getEffectiveAppearanceName()
	isDark := w.isAppearanceDark()
	w.parent.followApplicationTheme = false

	switch theme {
	case WinThemeDark:
		if !isDark {
			appr, _ := w.getOppositeAppearance(currentAppearance)
			w.setAppearanceByName(appr)
		}
	case WinThemeLight:
		if isDark {
			appr, _ := w.getOppositeAppearance(currentAppearance)
			w.setAppearanceByName(appr)
		}
	}
}

// getTheme returns the current theme configuration for the window.
func (w *macosWebviewWindow) getTheme() WinTheme {
	if w.parent.followApplicationTheme {
		return WinThemeApplication
	}

	explicitAppearance := w.getExplicitAppearanceName()

	if !explicitAppearance {
		return WinThemeSystem
	}

	if w.isAppearanceDark() {
		return WinThemeDark
	}

	return WinThemeLight
}
