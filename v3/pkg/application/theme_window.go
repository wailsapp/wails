package application

import "github.com/wailsapp/wails/v3/pkg/events"

type WinTheme string

const (
	WinThemeApplication WinTheme = "application"
	WinThemeDark        WinTheme = "dark"
	WinThemeLight       WinTheme = "light"
	WinThemeSystem      WinTheme = "system"
)

func (t WinTheme) String() string {
	return string(t)
}

func (t WinTheme) Valid() bool {
	switch t {
	case WinThemeApplication, WinThemeDark, WinThemeLight, WinThemeSystem:
		return true
	}
	return false
}

// GetTheme returns the current theme for current Windows - Windows OS
func (w *WebviewWindow) GetTheme() string {
	if w.impl == nil {
		return WinThemeApplication.String()
	}
	return w.impl.getTheme().String()
}

// SetTheme sets the theme for the current Window - Windows OS
func (w *WebviewWindow) SetTheme(theme WinTheme) {
	if !theme.Valid() {
		return
	}
	if w.impl != nil {
		w.impl.setTheme(theme)
	}
	// actual := w.GetTheme()
	// Notify listeners of the theme change
	// w.EmitEvent("windowThemeChanged", map[string]any{
	// 	"windowID":   w.ID(),
	// 	"windowName": w.Name(),
	// 	"theme":      actual,
	// })
	w.emit(events.WindowEventType(events.Common.ThemeChanged))
}
