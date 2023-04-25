package system

import (
	"os/exec"
	"strings"

	"github.com/wailsapp/wails/v2/internal/shell"
	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
	"github.com/wailsapp/wails/v2/internal/system/packagemanager"
)

var (
	IsAppleSilicon bool
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

func checkNodejs() *packagemanager.Dependency {

	// Check for Nodejs
	output, err := exec.Command("node", "-v").Output()
	installed := true
	version := ""
	if err != nil {
		installed = false
	} else {
		if len(output) > 0 {
			version = strings.TrimSpace(strings.Split(string(output), "\n")[0])[1:]
		}
	}
	return &packagemanager.Dependency{
		Name:           "Nodejs",
		PackageName:    "N/A",
		Installed:      installed,
		InstallCommand: "Available at https://nodejs.org/en/download/",
		Version:        version,
		Optional:       false,
		External:       false,
	}
}

func checkNPM() *packagemanager.Dependency {

	// Check for npm
	output, err := exec.Command("npm", "-version").Output()
	installed := true
	version := ""
	if err != nil {
		installed = false
	} else {
		version = strings.TrimSpace(strings.Split(string(output), "\n")[0])
	}
	return &packagemanager.Dependency{
		Name:           "npm ",
		PackageName:    "N/A",
		Installed:      installed,
		InstallCommand: "Available at https://nodejs.org/en/download/",
		Version:        version,
		Optional:       false,
		External:       false,
	}
}

func checkUPX() *packagemanager.Dependency {

	// Check for npm
	output, err := exec.Command("upx", "-V").Output()
	installed := true
	version := ""
	if err != nil {
		installed = false
	} else {
		version = strings.TrimSpace(strings.Split(string(output), "\n")[0])
	}
	return &packagemanager.Dependency{
		Name:           "upx ",
		PackageName:    "N/A",
		Installed:      installed,
		InstallCommand: "Available at https://upx.github.io/",
		Version:        version,
		Optional:       true,
		External:       false,
	}
}

func checkNSIS() *packagemanager.Dependency {

	// Check for nsis installer
	output, err := exec.Command("makensis", "-VERSION").Output()
	installed := true
	version := ""
	if err != nil {
		installed = false
	} else {
		version = strings.TrimSpace(strings.Split(string(output), "\n")[0])
	}
	return &packagemanager.Dependency{
		Name:           "nsis ",
		PackageName:    "N/A",
		Installed:      installed,
		InstallCommand: "More info at https://wails.io/docs/guides/windows-installer/",
		Version:        version,
		Optional:       true,
		External:       false,
	}
}

func checkLibrary(name string) func() *packagemanager.Dependency {
	return func() *packagemanager.Dependency {
		output, _, _ := shell.RunCommand(".", "pkg-config", "--cflags", name)
		installed := len(strings.TrimSpace(output)) > 0

		return &packagemanager.Dependency{
			Name:           "lib" + name + " ",
			PackageName:    "N/A",
			Installed:      installed,
			InstallCommand: "Install via your package manager",
			Version:        "N/A",
			Optional:       false,
			External:       false,
		}
	}
}

func checkDocker() *packagemanager.Dependency {

	// Check for npm
	output, err := exec.Command("docker", "version").Output()
	installed := true
	version := ""

	// Docker errors if it is not running so check for that
	if len(output) == 0 && err != nil {
		installed = false
	} else {
		// Version is in a line like: " Version:           20.10.5"
		versionOutput := strings.Split(string(output), "\n")
		for _, line := range versionOutput[1:] {
			splitLine := strings.Split(line, ":")
			if len(splitLine) > 1 {
				key := strings.TrimSpace(splitLine[0])
				if key == "Version" {
					version = strings.TrimSpace(splitLine[1])
					break
				}
			}
		}
	}
	return &packagemanager.Dependency{
		Name:           "docker ",
		PackageName:    "N/A",
		Installed:      installed,
		InstallCommand: "Available at https://www.docker.com/products/docker-desktop",
		Version:        version,
		Optional:       true,
		External:       false,
	}
}
