package cmd

import (
	"path/filepath"
)

// FrontendHelper handles the frontend wails assets
// like bridge, runtime, etc
type FrontendHelper struct {
	fs *FSHelper
}

// NewFrontendHelper creates a new frontend helper
func NewFrontendHelper() *FrontendHelper {
	return &FrontendHelper{
		fs: NewFSHelper(),
	}
}

// InstallRuntime installs the wails runtime and associated libraries in the project
func (f *FrontendHelper) InstallRuntime(projectOptions *ProjectOptions) error {

	// Determine target location
	wailsLocalDir := filepath.FromSlash(projectOptions.selectedTemplate.Metadata.WailsDir)
	runtimeDir, err := filepath.Abs(filepath.Join(projectOptions.OutputDirectory, projectOptions.FrontEnd.Dir, wailsLocalDir))
	if err != nil {
		return err
	}

	// Create directory if needed
	err = f.fs.MkDirs(runtimeDir)
	if err != nil {
		return err
	}

	// Copy runtime assets
	runtimeAssetsDir, err := f.fs.LocalDir("../runtimeassets/frontend")
	if err != nil {
		return err
	}

	return runtimeAssetsDir.CopyFiles(runtimeDir)

}
