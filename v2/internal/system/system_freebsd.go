//go:build freebsd
// +build freebsd

package system

import (
	"os/exec"
	"strings"

	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
	"github.com/wailsapp/wails/v2/internal/system/packagemanager"
)

func checkPkgConfig() *packagemanager.Dependency {
	// Check for pkg-config
	output, err := exec.Command("pkg-config", "-version").Output()
	installed := true
	version := ""
	if err != nil {
		installed = false
	} else {
		version = strings.TrimSpace(strings.Split(string(output), "\n")[0])
	}
	return &packagemanager.Dependency{
		Name:           "pkg-config ",
		PackageName:    "pkgconf",
		Installed:      installed,
		InstallCommand: "pkg install pkgconf",
		Version:        version,
		Optional:       true,
		External:       false,
	}
}

func (i *Info) discover() error {
	var err error
	osinfo, err := operatingsystem.Info()
	if err != nil {
		return err
	}
	i.OS = osinfo

	i.Dependencies = append(i.Dependencies, checkPkgConfig())
	i.Dependencies = append(i.Dependencies, checkNodejs())
	i.Dependencies = append(i.Dependencies, checkNPM())
	i.Dependencies = append(i.Dependencies, checkUPX())
	i.Dependencies = append(i.Dependencies, checkNSIS())
	return nil
}
