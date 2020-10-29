package build

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/internal/fs"
)

// GenerateBackendJSPackage will attempt to create the backend javascript package
// used by the frontend to access methods and structs
func GenerateBackendJSPackage(options *Options) error {

	// Create directory
	err := createBackendJSDirectory()
	if err != nil {
		return errors.Wrap(err, "Error Generating Backend JS Package")
	}

	// Generate index.js

	// Generate Method wrappers
	// err = generateMethodWrappers()

	// Generate Structs

	//

	return nil
}

func createBackendJSDirectory() error {

	// Path to package dir
	packageDir, err := filepath.Abs("./frontend/backend")
	if err != nil {
		return err
	}
	return fs.Mkdir(packageDir)
}
