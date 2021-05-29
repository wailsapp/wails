// +build windows

package system

import (
	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
	"github.com/wailsapp/wails/v2/internal/system/packagemanager"
	"os/exec"
	"strings"
)

func (i *Info) discover() error {

	var err error
	osinfo, err := operatingsystem.Info()
	if err != nil {
		return err
	}
	i.OS = osinfo

	// Check for gcc
	output, err := exec.Command("gcc", "--version").Output()
	installed := true
	version := ""
	if err != nil {
		installed = false
	} else {
		version = strings.TrimSpace(strings.Split(string(output), "\n")[0])
	}
	gccDependency := &packagemanager.Dependancy{
		Name:           "gcc ",
		PackageName:    "N/A",
		Installed:      installed,
		InstallCommand: "",
		Version:        version,
		Optional:       false,
		External:       false,
	}
	i.Dependencies = append(i.Dependencies, gccDependency)
	i.Dependencies = append(i.Dependencies, checkNPM())
	i.Dependencies = append(i.Dependencies, checkUPX())

	return nil
}
