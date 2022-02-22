//go:build linux
// +build linux

package system

import (
	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
	"github.com/wailsapp/wails/v2/internal/system/packagemanager"
)

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
			if dep.Name == "npm" {
				locallyInstalled := checkNPM()
				if locallyInstalled.Installed {
					dep.Installed = true
					dep.Version = locallyInstalled.Version
				}
			}
			if dep.Name == "docker" {
				locallyInstalled := checkDocker()
				if locallyInstalled.Installed {
					dep.Installed = true
					dep.Version = locallyInstalled.Version
				}
			}
			if dep.Name == "upx" {
				locallyInstalled := checkUPX()
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
