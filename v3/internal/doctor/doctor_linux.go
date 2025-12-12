//go:build linux

package doctor

import (
	"os"
	"strings"

	"github.com/wailsapp/wails/v3/internal/doctor/packagemanager"
	"github.com/wailsapp/wails/v3/internal/operatingsystem"
)

func getInfo() (map[string]string, bool) {
	result := make(map[string]string)

	// Check for NVIDIA driver
	result["NVIDIA Driver"] = getNvidiaDriverInfo()

	return result, true
}

func getNvidiaDriverInfo() string {
	version, err := os.ReadFile("/sys/module/nvidia/version")
	if err != nil {
		return "N/A"
	}

	versionStr := strings.TrimSpace(string(version))

	srcVersion, err := os.ReadFile("/sys/module/nvidia/srcversion")
	if err != nil {
		return versionStr
	}

	return versionStr + " (" + strings.TrimSpace(string(srcVersion)) + ")"
}

func checkPlatformDependencies(result map[string]string, ok *bool) {
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
				*ok = false
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

	checkCommonDependencies(result, ok)
}
