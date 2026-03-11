//go:build darwin

package application

import "fmt"

// getOppositeMacAppearance returns the macOS appearance that represents
// the opposite light/dark variant.
func getOppositeMacAppearance(name string) (MacAppearanceType, error) {
	if name == "NSAppearanceNameDarkAqua" {
		return "NSAppearanceNameAqua", nil
	}

	// If opposite appearance doesnt match then send the default Dark Appearance
	err := fmt.Errorf("unknown appearance name: %s", name)
	return "NSAppearanceNameDarkAqua", err
}

// isMacAppearanceDark reports whether the current window appearance
// corresponds to a dark macOS appearance.
func isMacAppearanceDark(appr string) bool {
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
	currentDark := isMacAppearanceDark(currentAppearance)

	switch globalApplication.theme {
	case AppSystemDefault:
		w.clearAppearance()
		return
	case AppDark:
		if !currentDark {
			appr, _ := getOppositeMacAppearance(currentAppearance)
			w.setAppearanceByName(appr)
		}
	case AppLight:
		if currentDark {
			appr, _ := getOppositeMacAppearance(currentAppearance)
			w.setAppearanceByName(appr)
		}
	}
}

// setTheme sets the theme for the window. If WinThemeApplication is provided,
// the window will follow global application theme settings.
func (w *macosWebviewWindow) setTheme(theme WinTheme) {
	switch theme {
	case WinSystemDefault:
		w.parent.followApplicationTheme = false
		w.clearAppearance()
		return
	case WinAppDefault:
		w.parent.followApplicationTheme = true
		w.syncTheme()
		return
	}

	currentAppearance := w.getEffectiveAppearanceName()
	isDark := isMacAppearanceDark(currentAppearance)
	w.parent.followApplicationTheme = false

	switch theme {
	case WinDark:
		if !isDark {
			appr, _ := getOppositeMacAppearance(currentAppearance)
			w.setAppearanceByName(appr)
		}
	case WinLight:
		if isDark {
			appr, _ := getOppositeMacAppearance(currentAppearance)
			w.setAppearanceByName(appr)
		}
	}
}

// getTheme returns the current theme configuration for the window.
func (w *macosWebviewWindow) getTheme() WinTheme {
	if w.parent.followApplicationTheme {
		return WinAppDefault
	}

	explicitAppearance := w.getExplicitAppearanceName()

	if explicitAppearance == "" {
		return WinSystemDefault
	}

	if isMacAppearanceDark(w.getEffectiveAppearanceName()) {
		return WinDark
	}

	return WinLight
}
