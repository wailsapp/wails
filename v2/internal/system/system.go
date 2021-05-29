package system

import (
	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
	"github.com/wailsapp/wails/v2/internal/system/packagemanager"
	"os/exec"
	"strings"
)

// Info holds information about the current operating system,
// package manager and required dependancies
type Info struct {
	OS           *operatingsystem.OS
	PM           packagemanager.PackageManager
	Dependencies packagemanager.DependencyList
}

// GetInfo scans the system for operating system details,
// the system package manager and the status of required
// dependancies.
func GetInfo() (*Info, error) {
	var result Info
	err := result.discover()
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func checkNPM() *packagemanager.Dependancy {

	// Check for npm
	output, err := exec.Command("npm", "-version").Output()
	installed := true
	version := ""
	if err != nil {
		installed = false
	} else {
		version = strings.TrimSpace(strings.Split(string(output), "\n")[0])
	}
	return &packagemanager.Dependancy{
		Name:           "npm ",
		PackageName:    "N/A",
		Installed:      installed,
		InstallCommand: "Install from https://nodejs.org/en/download/",
		Version:        version,
		Optional:       false,
		External:       false,
	}
}

func checkUPX() *packagemanager.Dependancy {

	// Check for npm
	output, err := exec.Command("upx", "-V").Output()
	installed := true
	version := ""
	if err != nil {
		installed = false
	} else {
		version = strings.TrimSpace(strings.Split(string(output), "\n")[0])
	}
	return &packagemanager.Dependancy{
		Name:           "upx ",
		PackageName:    "N/A",
		Installed:      installed,
		InstallCommand: "Install from https://upx.github.io/",
		Version:        version,
		Optional:       true,
		External:       false,
	}
}
