package application

// AppTheme represents the theme preference for the application.
type AppTheme string

const (
	// AppSystemDefault follows the system theme (light or dark).
	AppSystemDefault AppTheme = "system"
	// AppDark forces the application to use a dark theme.
	AppDark          AppTheme = "dark"
	// AppLight forces the application to use a light theme.
	AppLight         AppTheme = "light"
)

// String returns the string representation of the application theme.
func (t AppTheme) String() string {
	return string(t)
}

// Valid returns true if the theme is a recognized AppTheme value.
func (t AppTheme) Valid() bool {
	switch t {
	case AppSystemDefault, AppDark, AppLight:
		return true
	}
	return false
}

// GetTheme returns the current application-level theme setting.
func (a *App) GetTheme() string {
	return a.theme.String()
}

// SetTheme sets the application-level theme preference.
// This will apply the theme to the application and any windows configured to follow it.
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
	a.Event.Emit("applicationThemeChanged", map[string]any{
		"theme": a.theme.String(),
	})
}
