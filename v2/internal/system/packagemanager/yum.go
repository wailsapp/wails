// +build linux

package packagemanager

import (
	"os/exec"
	"strings"

	"github.com/wailsapp/wails/v2/internal/shell"
)

// Yum represents the Yum manager
type Yum struct {
	name string
	osid string
}

// NewYum creates a new Yum instance
func NewYum(osid string) *Yum {
	return &Yum{
		name: "yum",
		osid: osid,
	}
}

// Packages returns the libraries that we need for Wails to compile
// They will potentially differ on different distributions or versions
func (y *Yum) Packages() packagemap {
	return packagemap{
		"libgtk-3": []*Package{
			{Name: "gtk3-devel", SystemPackage: true, Library: true},
		},
		"libwebkit": []*Package{
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
		},
		"docker": []*Package{
			{
				SystemPackage: false,
				Optional:      true,
				InstallCommand: map[string]string{
					"centos": "Follow the guide: https://docs.docker.com/engine/install/centos/",
					"fedora": "Follow the guide: https://docs.docker.com/engine/install/fedora/",
				},
			},
		},
	}
}

// Name returns the name of the package manager
func (y *Yum) Name() string {
	return y.name
}

// PackageInstalled tests if the given package name is installed
func (y *Yum) PackageInstalled(pkg *Package) (bool, error) {
	if pkg.SystemPackage == false {
		return false, nil
	}
	stdout, _, err := shell.RunCommand(".", "yum", "info", "installed", pkg.Name)
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
func (y *Yum) PackageAvailable(pkg *Package) (bool, error) {
	if pkg.SystemPackage == false {
		return false, nil
	}
	stdout, _, err := shell.RunCommand(".", "yum", "info", pkg.Name)
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
func (y *Yum) InstallCommand(pkg *Package) string {
	if pkg.SystemPackage == false {
		return pkg.InstallCommand[y.osid]
	}
	return "sudo yum install " + pkg.Name
}

func (y *Yum) getPackageVersion(pkg *Package, output string) {
	splitOutput := strings.Split(output, " ")
	if len(splitOutput) > 0 {
		pkg.Version = splitOutput[1]
	}
}
