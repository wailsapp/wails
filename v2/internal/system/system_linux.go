//go:build linux
// +build linux

package system

import (
	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
	"github.com/wailsapp/wails/v2/internal/system/packagemanager"
)

func checkLocallyInstalled(checker func() *packagemanager.Dependancy, dependency *packagemanager.Dependancy) {
	if !dependency.Installed {
		locallyInstalled := checker()
		if locallyInstalled.Installed {
			dependency.Installed = true
			dependency.Version = locallyInstalled.Version
		}
	}
}

var checkerFunctions = map[string]func() *packagemanager.Dependancy{
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
		dependencies, err := packagemanager.Dependancies(i.PM)
		if err != nil {
			return err
		}
		for _, dep := range dependencies {
			checker := checkerFunctions[dep.Name]
			if checker() != nil {
				checkLocallyInstalled(checker, dep)
				continue
			}
			if dep.Name == "nsis" {
				locallyInstalled := checkNSIS()
				if locallyInstalled.Installed {
					dep.Installed = true
					dep.Version = locallyInstalled.Version
				}
			}
		}
		i.Dependencies = dependencies
	}

	return nil
}
