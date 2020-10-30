package backendjs

import (
	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/internal/fs"
)

// GenerateBackendJSPackage will attempt to create the backend javascript package
// used by the frontend to access methods and structs
func GenerateBackendJSPackage() error {

	// Create directory
	err := createBackendJSDirectory()
	if err != nil {
		return errors.Wrap(err, "Error creating backend directory:")
	}

	// Generate Packages
	err = generatePackages()
	if err != nil {
		return errors.Wrap(err, "Error generating method wrappers:")
	}

	return nil
}

func createBackendJSDirectory() error {

	// Calculate the package directory
	// Note this is *always* called from the project directory
	// so using paths relative to CWD is fine
	dir, err := fs.RelativeToCwd("./frontend/backend")
	if err != nil {
		return errors.Wrap(err, "Error creating backend js directory")
	}

	// Only create the directory if it doesn't exit
	if !fs.DirExists(dir) {
		return fs.Mkdir(dir)
	}

	return nil
}
