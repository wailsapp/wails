//go:build linux

package doctor

import (
	"github.com/wailsapp/wails/v3/internal/doctor/packagemanager"
	"github.com/wailsapp/wails/v3/internal/operatingsystem"
)

func getInfo() (map[string]string, bool) {
	result := make(map[string]string)
	ok := true

	info, _ := operatingsystem.Info()

	pm := packagemanager.Find(info.ID)
	deps, _ := packagemanager.Dependencies(pm)
	for _, dep := range deps {
		var status string

		switch true {
		case !dep.Installed:
			if dep.Optional {
				status = "[Optional] "
			} else {
				ok = false
			}
			status += "not installed."
			if dep.InstallCommand != "" {
				status += " Install with: " + dep.InstallCommand
			}
		case dep.Version != "":
			status = dep.Version
		}

		result[dep.Name] = status
	}

	return result, ok
}
