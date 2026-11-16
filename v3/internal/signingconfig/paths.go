package signingconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

// IsGlobalPath requires the field's current value to match its recorded origin.
// Changing a resolved field cannot transfer trust to another credential path.
func IsGlobalPath(config manifest.SigningPlatform, field, value string) bool {
	var trusted string
	switch field {
	case "certificate":
		trusted = config.GlobalPaths.Certificate
	case "entitlements":
		trusted = config.GlobalPaths.Entitlements
	case "notarization_key":
		trusted = config.GlobalPaths.NotarizationKey
	}
	return value != "" && filepath.IsAbs(value) && value == trusted
}

// ResolvePath allows external files only when supplied by global defaults.
// Derived absolute paths within the project retain ordinary containment checks.
func ResolvePath(root string, config manifest.SigningPlatform, field, value string) (string, error) {
	if value == "" {
		return "", nil
	}
	var path string
	var err error
	if IsGlobalPath(config, field, value) {
		path, err = filepath.EvalSymlinks(value)
	} else {
		relative := value
		if filepath.IsAbs(value) {
			relative, err = filepath.Rel(root, value)
			if err != nil {
				return "", err
			}
		}
		path, err = manifest.ResolveProjectPath(root, field, relative, true)
	}
	if err != nil {
		return "", err
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("signing %s must be a regular file", field)
	}
	return path, nil
}
