//go:build darwin
// +build darwin

package system

import (
	"github.com/wailsapp/wails/v2/internal/system/packagemanager"
	"os/exec"
	"strings"
	"syscall"

	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
)

// Determine if the app is running on Apple Silicon
// Credit: https://www.yellowduck.be/posts/detecting-apple-silicon-via-go/
func init() {
	r, err := syscall.Sysctl("sysctl.proc_translated")
	if err != nil {
		return
	}

	IsAppleSilicon = r == "\x00\x00\x00" || r == "\x01\x00\x00"
}

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
	xcodeDep := &packagemanager.Dependency{
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
	i.Dependencies = append(i.Dependencies, checkUPX())
	i.Dependencies = append(i.Dependencies, checkNSIS())
	return nil
}
