package build

import (
	"fmt"
	"github.com/leaanthony/winicon"
	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/internal/shell"
	"github.com/wailsapp/wails/v2/pkg/buildassets"
	"os"
	"path/filepath"
)

func packageApplication(options *Options) error {
	// Generate icon
	var err error
	err = generateIcoFile(options)
	if err != nil {
		return err
	}

	// Ensure RC file is present
	err = generateRCFile(options)
	if err != nil {
		return err
	}

	// Ensure Manifest is present
	err = generateManifest(options)
	if err != nil {
		return err
	}

	// Create syso file
	err = compileResources(options)
	if err != nil {
		return err
	}

	return nil
}

func generateManifest(options *Options) error {
	filename := options.ProjectData.Name + ".exe.manifest"
	manifestFile := filepath.Join(options.ProjectData.Path, "build", "windows", filename)
	if !fs.FileExists(manifestFile) {
		return buildassets.RegenerateManifest(manifestFile)
	}
	return nil
}

func generateRCFile(options *Options) error {
	filename := options.ProjectData.Name + ".rc"
	rcFile := filepath.Join(options.ProjectData.Path, "build", "windows", filename)
	if !fs.FileExists(rcFile) {
		return buildassets.RegenerateRCFile(options.ProjectData.Path, options.ProjectData.Name)
	}
	return nil
}

func generateIcoFile(options *Options) error {
	// Check ico file exists already
	icoFile := filepath.Join(options.ProjectData.Path, "build", "windows", "app.ico")
	if !fs.FileExists(icoFile) {
		// Check icon exists
		appicon := filepath.Join(options.ProjectData.Path, "build", "appicon.png")
		if !fs.FileExists(appicon) {
			return fmt.Errorf("application icon missing: %s", appicon)
		}
		// Load icon
		input, err := os.Open(appicon)
		if err != nil {
			return err
		}
		output, err := os.OpenFile(icoFile, os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		err = winicon.GenerateIcon(input, output, []int{256, 128, 64, 48, 32, 16})
		if err != nil {
			return err
		}
	}
	return nil
}

func compileResources(options *Options) error {
	windowsBuildDir := filepath.Join(options.ProjectData.Path, "build", "windows")
	sourcefile := filepath.Join(options.ProjectData.BuildDir, "windows", options.ProjectData.Name+".rc")
	targetFile := filepath.Join(options.ProjectData.Path, options.ProjectData.Name+"-res.syso")
	_, _, err := shell.RunCommand(windowsBuildDir, "windres", "-o", targetFile, sourcefile)
	return err
}
