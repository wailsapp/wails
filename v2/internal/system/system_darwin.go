//go:build darwin
// +build darwin

package system

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"
	"github.com/wailsapp/wails/v2/internal/system/packagemanager"
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

	i.Dependencies = append(i.Dependencies, checkXCodeSelect())
	i.Dependencies = append(i.Dependencies, checkNodejs())
	i.Dependencies = append(i.Dependencies, checkNPM())
	i.Dependencies = append(i.Dependencies, checkXCodeBuild())
	i.Dependencies = append(i.Dependencies, checkUPX())
	i.Dependencies = append(i.Dependencies, checkNSIS())
	return nil
}

func checkXCodeSelect() *packagemanager.Dependency {
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
	return &packagemanager.Dependency{
		Name:           "Xcode command line tools ",
		PackageName:    "N/A",
		Installed:      installed,
		InstallCommand: "xcode-select --install",
		Version:        version,
		Optional:       false,
		External:       false,
	}
}

func checkXCodeBuild() *packagemanager.Dependency {
	// Check for xcode
	output, err := exec.Command("xcodebuild", "-version").Output()
	installed := true
	version := ""
	if err != nil {
		installed = false
	} else if l := strings.Split(string(output), "\n"); len(l) >= 2 {
		version = fmt.Sprintf("%s (%s)",
			strings.TrimPrefix(l[0], "Xcode "),
			strings.TrimPrefix(l[1], "Build version "))
	} else {
		version = "N/A"
	}

	return &packagemanager.Dependency{
		Name:           "Xcode",
		PackageName:    "N/A",
		Installed:      installed,
		InstallCommand: "Available at https://apps.apple.com/us/app/xcode/id497799835",
		Version:        version,
		Optional:       true,
		External:       false,
	}
}
