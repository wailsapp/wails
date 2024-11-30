//go:build linux

package packagemanager

import (
	"os/exec"
	"strings"
)

// Dnf represents the Dnf manager
type Dnf struct {
	name string
	osid string
}

// NewDnf creates a new Dnf instance
func NewDnf(osid string) *Dnf {
	return &Dnf{
		name: "dnf",
		osid: osid,
	}
}

// Packages returns the libraries that we need for Wails to compile
// They will potentially differ on different distributions or versions
func (y *Dnf) Packages() Packagemap {
	return Packagemap{
		"gtk3": []*Package{
			{Name: "gtk3-devel", SystemPackage: true, Library: true},
		},
		"webkit2gtk": []*Package{
			{Name: "webkit2gtk4.1-devel", SystemPackage: true, Library: true},
			{Name: "webkit2gtk3-devel", SystemPackage: true, Library: true},
			// {Name: "webkitgtk3-devel", SystemPackage: true, Library: true},
		},
		"gcc": []*Package{
			{Name: "gcc-c++", SystemPackage: true},
		},
		"pkg-config": []*Package{
			{Name: "pkgconf-pkg-config", SystemPackage: true},
		},
		"npm": []*Package{
			{Name: "npm", SystemPackage: true},
			{Name: "nodejs-npm", SystemPackage: true},
		},
	}
}

// Name returns the name of the package manager
func (y *Dnf) Name() string {
	return y.name
}

// PackageInstalled tests if the given package name is installed
func (y *Dnf) PackageInstalled(pkg *Package) (bool, error) {
	if !pkg.SystemPackage {
		if pkg.InstallCheck != nil {
			return pkg.InstallCheck(), nil
		}
		return false, nil
	}
	stdout, err := execCmd("dnf", "info", "installed", pkg.Name)
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if ok {
			return false, nil
		}
		return false, err
	}

	splitoutput := strings.Split(stdout, "\n")
	for _, line := range splitoutput {
		if strings.HasPrefix(line, "Version") {
			splitline := strings.Split(line, ":")
			pkg.Version = strings.TrimSpace(splitline[1])
		}
	}

	return true, err
}

// PackageAvailable tests if the given package is available for installation
func (y *Dnf) PackageAvailable(pkg *Package) (bool, error) {
	if pkg.SystemPackage == false {
		return false, nil
	}
	stdout, err := execCmd("dnf", "info", pkg.Name)
	// We add a space to ensure we get a full match, not partial match
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if ok {
			return false, nil
		}
		return false, err
	}
	splitoutput := strings.Split(stdout, "\n")
	for _, line := range splitoutput {
		if strings.HasPrefix(line, "Version") {
			splitline := strings.Split(line, ":")
			pkg.Version = strings.TrimSpace(splitline[1])
		}
	}
	return true, nil
}

// InstallCommand returns the package manager specific command to install a package
func (y *Dnf) InstallCommand(pkg *Package) string {
	if pkg.SystemPackage == false {
		return pkg.InstallCommand
	}
	return "sudo dnf install " + pkg.Name
}

func (y *Dnf) getPackageVersion(pkg *Package, output string) {
	splitOutput := strings.Split(output, " ")
	if len(splitOutput) > 0 {
		pkg.Version = splitOutput[1]
	}
}
