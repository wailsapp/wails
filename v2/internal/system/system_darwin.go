// +build darwin

package system

import (
	"os/exec"
	"strings"

	"github.com/wailsapp/wails/v2/internal/system/packagemanager"

	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
)

func (i *Info) discover() error {
	var err error
	osinfo, err := operatingsystem.Info()
	if err != nil {
		return err
	}
	i.OS = osinfo

	// Check for xcode command line tools
	output, err := exec.Command("xcode-select", "-v").Output()
	installed := true
	version := ""
	if err != nil {
		installed = false
	} else {
		version = strings.TrimPrefix(string(output), "xcode-select version ")
		version = strings.TrimSpace(version)
		version = strings.TrimSuffix(version, ".")
	}
	xcodeDep := &packagemanager.Dependancy{
		Name:           "xcode command line tools ",
		PackageName:    "N/A",
		Installed:      installed,
		InstallCommand: "xcode-select --install",
		Version:        version,
		Optional:       false,
		External:       false,
	}
	i.Dependencies = append(i.Dependencies, xcodeDep)
	i.Dependencies = append(i.Dependencies, checkNPM())
	return nil
}
