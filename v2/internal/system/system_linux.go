//go:build linux
// +build linux

package system

import (
	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
	"github.com/wailsapp/wails/v2/internal/system/packagemanager"
)

func checkGCC() *packagemanager.Dependency {

	version := packagemanager.AppVersion("gcc")

	return &packagemanager.Dependency{
		Name:           "gcc ",
		PackageName:    "N/A",
		Installed:      version != "",
		InstallCommand: "Install via your package manager",
		Version:        version,
		Optional:       false,
		External:       false,
	}
}

func checkPkgConfig() *packagemanager.Dependency {

	version := packagemanager.AppVersion("pkg-config")

	return &packagemanager.Dependency{
		Name:           "pkg-config ",
		PackageName:    "N/A",
		Installed:      version != "",
		InstallCommand: "Install via your package manager",
		Version:        version,
		Optional:       false,
		External:       false,
	}
}

func checkLocallyInstalled(checker func() *packagemanager.Dependency, dependency *packagemanager.Dependency) {
	if !dependency.Installed {
		applyLocalCheck(dependency, checker())
	}
}

func applyLocalCheck(dependency, locallyInstalled *packagemanager.Dependency) {
	if !locallyInstalled.Installed {
		return
	}
	dependency.Installed = true
	dependency.Version = locallyInstalled.Version
	// doctor reports a dependency without a package name as "Not Found".
	if dependency.PackageName == "" {
		dependency.PackageName = locallyInstalled.PackageName
	}
}

var checkerFunctions = map[string]func() *packagemanager.Dependency{
	"Nodejs":     checkNodejs,
	"npm":        checkNPM,
	"docker":     checkDocker,
	"upx":        checkUPX,
	"gcc":        checkGCC,
	"pkg-config": checkPkgConfig,
	"libgtk-3":   checkLibrary("libgtk-3"),
	"libwebkit":  checkLibrary("libwebkit"),
}

func (i *Info) discover() error {

	var err error
	osinfo, err := operatingsystem.Info()
	if err != nil {
		return err
	}
	i.OS = osinfo

	i.PM = packagemanager.Find(osinfo.ID)
	if i.PM != nil {
		dependencies, err := packagemanager.Dependencies(i.PM)
		if err != nil {
			return err
		}
		for _, dep := range dependencies {
			checker := checkerFunctions[dep.Name]
			if checker != nil {
				checkLocallyInstalled(checker, dep)
			}
			if dep.Name == "nsis" {
				applyLocalCheck(dep, checkNSIS())
			}
		}
		i.Dependencies = dependencies
	}

	return nil
}
