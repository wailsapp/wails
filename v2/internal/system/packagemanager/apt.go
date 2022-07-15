//go:build linux
// +build linux

package packagemanager

import (
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/wailsapp/wails/v2/internal/shell"
)

// Apt represents the Apt manager
type Apt struct {
	name string
	osid string
}

// NewApt creates a new Apt instance
func NewApt(osid string) *Apt {
	return &Apt{
		name: "apt",
		osid: osid,
	}
}

// Packages returns the libraries that we need for Wails to compile
// They will potentially differ on different distributions or versions
func (a *Apt) Packages() packagemap {
	return packagemap{
		"libgtk-3": []*Package{
			{Name: "libgtk-3-dev", SystemPackage: true, Library: true},
		},
		"libwebkit": []*Package{
			{Name: "libwebkit2gtk-4.0-dev", SystemPackage: true, Library: true},
		},
		"gcc": []*Package{
			{Name: "build-essential", SystemPackage: true},
		},
		"pkg-config": []*Package{
			{Name: "pkg-config", SystemPackage: true},
		},
		"npm": []*Package{
			{Name: "npm", SystemPackage: true},
		},
		"docker": []*Package{
			{Name: "docker.io", SystemPackage: true, Optional: true},
		},
		"nsis": []*Package{
			{Name: "nsis", SystemPackage: true, Optional: true},
		},
	}
}

// Name returns the name of the package manager
func (a *Apt) Name() string {
	return a.name
}

// PackageInstalled tests if the given package name is installed
func (a *Apt) PackageInstalled(pkg *Package) (bool, error) {
	if pkg.SystemPackage == false {
		return false, nil
	}
	cmd := exec.Command("apt", "list", "-qq", pkg.Name)
	var stdo, stde bytes.Buffer
	cmd.Stdout = &stdo
	cmd.Stderr = &stde
	cmd.Env = append(os.Environ(), "LANGUAGE=en")
	err := cmd.Run()
	return strings.Contains(stdo.String(), "[installed]"), err
}

// PackageAvailable tests if the given package is available for installation
func (a *Apt) PackageAvailable(pkg *Package) (bool, error) {
	if pkg.SystemPackage == false {
		return false, nil
	}
	stdout, _, err := shell.RunCommand(".", "apt", "list", "-qq", pkg.Name)
	// We add a space to ensure we get a full match, not partial match
	output := a.removeEscapeSequences(stdout)
	installed := strings.HasPrefix(output, pkg.Name)
	a.getPackageVersion(pkg, output)
	return installed, err
}

// InstallCommand returns the package manager specific command to install a package
func (a *Apt) InstallCommand(pkg *Package) string {
	if pkg.SystemPackage == false {
		return pkg.InstallCommand[a.osid]
	}
	return "sudo apt install " + pkg.Name
}

func (a *Apt) removeEscapeSequences(in string) string {
	escapechars, _ := regexp.Compile(`\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])`)
	return escapechars.ReplaceAllString(in, "")
}

func (a *Apt) getPackageVersion(pkg *Package, output string) {

	splitOutput := strings.Split(output, " ")
	if len(splitOutput) > 1 {
		pkg.Version = splitOutput[1]
	}
}
