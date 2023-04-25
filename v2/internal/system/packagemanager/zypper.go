//go:build linux
// +build linux

package packagemanager

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/wailsapp/wails/v2/internal/shell"
)

// Zypper represents the Zypper package manager
type Zypper struct {
	name string
	osid string
}

// NewZypper creates a new Zypper instance
func NewZypper(osid string) *Zypper {
	return &Zypper{
		name: "zypper",
		osid: osid,
	}
}

// Packages returns the libraries that we need for Wails to compile
// They will potentially differ on different distributions or versions
func (z *Zypper) Packages() packagemap {
	return packagemap{
		"libgtk-3": []*Package{
			{Name: "gtk3-devel", SystemPackage: true, Library: true},
		},
		"libwebkit": []*Package{
			{Name: "webkit2gtk3-soup2-devel", SystemPackage: true, Library: true},
			{Name: "webkit2gtk3-devel", SystemPackage: true, Library: true},
		},
		"gcc": []*Package{
			{Name: "gcc-c++", SystemPackage: true},
		},
		"pkg-config": []*Package{
			{Name: "pkg-config", SystemPackage: true},
			{Name: "pkgconf-pkg-config", SystemPackage: true},
		},
		"npm": []*Package{
			{Name: "npm10", SystemPackage: true},
		},
		"docker": []*Package{
			{Name: "docker", SystemPackage: true, Optional: true},
		},
	}
}

// Name returns the name of the package manager
func (z *Zypper) Name() string {
	return z.name
}

// PackageInstalled tests if the given package name is installed
func (z *Zypper) PackageInstalled(pkg *Package) (bool, error) {
	if pkg.SystemPackage == false {
		return false, nil
	}
	var env []string
	env = shell.SetEnv(env, "LANGUAGE", "en_US.utf-8")
	stdout, _, err := shell.RunCommandWithEnv(env, ".", "zypper", "info", pkg.Name)
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if ok {
			return false, nil
		}
		return false, err
	}
	reg := regexp.MustCompile(`.*Installed\s*:\s*(Yes)\s*`)
	matches := reg.FindStringSubmatch(stdout)
	pkg.Version = ""
	noOfMatches := len(matches)
	if noOfMatches > 1 {
		z.getPackageVersion(pkg, stdout)
	}
	return noOfMatches > 1, err
}

// PackageAvailable tests if the given package is available for installation
func (z *Zypper) PackageAvailable(pkg *Package) (bool, error) {
	if pkg.SystemPackage == false {
		return false, nil
	}
	var env []string
	env = shell.SetEnv(env, "LANGUAGE", "en_US.utf-8")
	stdout, _, err := shell.RunCommandWithEnv(env, ".", "zypper", "info", pkg.Name)
	// We add a space to ensure we get a full match, not partial match
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if ok {
			return false, nil
		}
		return false, err
	}

	available := strings.Contains(stdout, "Information for package")
	if available {
		z.getPackageVersion(pkg, stdout)
	}

	return available, nil
}

// InstallCommand returns the package manager specific command to install a package
func (z *Zypper) InstallCommand(pkg *Package) string {
	if pkg.SystemPackage == false {
		return pkg.InstallCommand[z.osid]
	}
	return "sudo zypper in " + pkg.Name
}

func (z *Zypper) getPackageVersion(pkg *Package, output string) {

	reg := regexp.MustCompile(`.*Version.*:(.*)`)
	matches := reg.FindStringSubmatch(output)
	pkg.Version = ""
	noOfMatches := len(matches)
	if noOfMatches > 1 {
		pkg.Version = strings.TrimSpace(matches[1])
	}
}
