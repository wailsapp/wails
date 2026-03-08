//go:build darwin

package application

import "fmt"

// getOppositeAppearance returns the macOS appearance that represents
// the opposite light/dark variant while preserving vibrancy and accessibility.
func (w *macosWebviewWindow) getOppositeAppearance(name string) (MacAppearanceType, error) {
	complementary := map[string]MacAppearanceType{
		"NSAppearanceNameAqua":     "NSAppearanceNameDarkAqua",
		"NSAppearanceNameDarkAqua": "NSAppearanceNameAqua",

		"NSAppearanceNameVibrantLight": "NSAppearanceNameVibrantDark",
		"NSAppearanceNameVibrantDark":  "NSAppearanceNameVibrantLight",

		"NSAppearanceNameAccessibilityHighContrastAqua":     "NSAppearanceNameAccessibilityHighContrastDarkAqua",
		"NSAppearanceNameAccessibilityHighContrastDarkAqua": "NSAppearanceNameAccessibilityHighContrastAqua",

		"NSAppearanceNameAccessibilityHighContrastVibrantLight": "NSAppearanceNameAccessibilityHighContrastVibrantDark",
		"NSAppearanceNameAccessibilityHighContrastVibrantDark":  "NSAppearanceNameAccessibilityHighContrastVibrantLight",
	}

	if result, ok := complementary[name]; ok {
		return result, nil
	}

	return "", fmt.Errorf("unknown appearance name: %s", name)
}

// isAppearanceDark reports whether the current window appearance
// corresponds to a dark macOS appearance.
func (w *macosWebviewWindow) isAppearanceDark() bool {
	appr := w.getAppearanceName()
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

// Sync Window Apperance with Application Theme if the window is set to follow application theme
func (w *macosWebviewWindow) syncTheme() {
	if !w.parent.followApplicationTheme {
		return
	}

	currentAppearance := w.getAppearanceName()
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

	currentAppearance := w.getAppearanceName()
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

func (w *macosWebviewWindow) getTheme() WinTheme {
	if w.parent.followApplicationTheme {
		return WinThemeApplication
	}

	if !w.hasExplicitAppearance() {
		return WinThemeSystem
	}

	if w.isAppearanceDark() {
		return WinThemeDark
	}

	return WinThemeLight
}
