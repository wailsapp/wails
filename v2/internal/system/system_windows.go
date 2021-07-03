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
		InstallCommand: "Available at https://jmeubank.github.io/tdm-gcc/",
		Version:        version,
		Optional:       false,
		External:       false,
	}
	i.Dependencies = append(i.Dependencies, gccDependency)
	i.Dependencies = append(i.Dependencies, checkNPM())
	i.Dependencies = append(i.Dependencies, checkUPX())
	i.Dependencies = append(i.Dependencies, checkDocker())

	return nil
}

// IsAppleSilicon returns true if the app is running on Apple Silicon
// Credit: https://www.yellowduck.be/posts/detecting-apple-silicon-via-go/
// NOTE: Not applicable to windows
func IsAppleSilicon() bool {
	return false
}
