package application

type AppTheme string

const (
	AppSystemDefault AppTheme = "system"
	AppDark          AppTheme = "dark"
	AppLight         AppTheme = "light"
)

func (t AppTheme) String() string {
	return string(t)
}

func (t AppTheme) Valid() bool {
	switch t {
	case AppSystemDefault, AppDark, AppLight:
		return true
	}
	return false
}

// GetTheme returns the app-level theme setting as a string.
func (a *App) GetTheme() string {
	return a.theme.String()
}

// SetTheme sets the app-level theme preference.
func (a *App) SetTheme(theme AppTheme) {
	if !theme.Valid() {
		return
	}

	if theme == a.theme {
		return
	}
	a.theme = theme

	if a.impl != nil {
		a.impl.setTheme(theme)
	}

	// Notify listeners of the theme change
	// a.Event.Emit("applicationThemeChanged", theme.String())
	a.Event.Emit("applicationThemeChanged", map[string]any{
		"theme": a.theme.String(),
	})
}
