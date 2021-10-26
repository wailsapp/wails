package build

import (
	"fmt"
	"github.com/leaanthony/winicon"
	"github.com/tc-hib/winres"
	"github.com/wailsapp/wails/v2/internal/fs"
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

func generateIcoFile(options *Options) error {
	// Check ico file exists already
	icoFile := filepath.Join(options.ProjectData.Path, "build", "windows", "icon.ico")
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

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer func() {
		os.Chdir(currentDir)
	}()
	windowsDir := filepath.Join(options.ProjectData.Path, "build", "windows")
	err = os.Chdir(windowsDir)
	if err != nil {
		return err
	}
	rs := winres.ResourceSet{}
	icon := filepath.Join(windowsDir, "icon.ico")
	iconFile, err := os.Open(icon)
	if err != nil {
		return err
	}
	defer iconFile.Close()
	ico, err := winres.LoadICO(iconFile)
	if err != nil {
		return err
	}
	err = rs.SetIcon(winres.RT_ICON, ico)
	if err != nil {
		return err
	}

	ManifestFilename := options.ProjectData.Name + ".exe.manifest"
	manifestData, err := os.ReadFile(ManifestFilename)
	xmlData, err := winres.AppManifestFromXML(manifestData)
	if err != nil {
		return err
	}
	rs.SetManifest(xmlData)

	targetFile := filepath.Join(options.ProjectData.Path, options.ProjectData.Name+"-res.syso")
	fout, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer fout.Close()

	archs := map[string]winres.Arch{
		"amd64": winres.ArchAMD64,
	}
	targetArch, supported := archs[options.Arch]
	if !supported {
		return fmt.Errorf("arch '%s' not supported", options.Arch)
	}

	err = rs.WriteObject(fout, targetArch)
	if err != nil {
		return err
	}
	return nil
}
