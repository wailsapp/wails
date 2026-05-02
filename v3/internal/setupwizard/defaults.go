package setupwizard

import (
	"github.com/wailsapp/wails/v3/internal/defaults"
)

// Re-export types for convenience
type GlobalDefaults = defaults.GlobalDefaults
type AuthorDefaults = defaults.AuthorDefaults
type ProjectDefaults = defaults.ProjectDefaults

// DefaultGlobalDefaults returns sensible defaults for first-time users
func DefaultGlobalDefaults() GlobalDefaults {
	return defaults.Default()
}

// GetDefaultsPath returns the path to the defaults.yaml file
func GetDefaultsPath() (string, error) {
	return defaults.GetDefaultsPath()
}

// LoadGlobalDefaults loads the global defaults from the config file
func LoadGlobalDefaults() (GlobalDefaults, error) {
	return defaults.Load()
}

// SaveGlobalDefaults saves the global defaults to the config file
func SaveGlobalDefaults(d GlobalDefaults) error {
	return defaults.Save(d)
}
