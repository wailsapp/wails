//go:build linux

package packagemanager

import (
	"regexp"
	"strings"
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
func (a *Apt) Packages() Packagemap {
	return Packagemap{
		"gtk3": []*Package{
			{Name: "libgtk-3-dev", SystemPackage: true, Library: true},
		},
		"webkit2gtk": []*Package{
			{Name: "libwebkit2gtk-4.1-dev", SystemPackage: true, Library: true},
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
	}
}

// Name returns the name of the package manager
func (a *Apt) Name() string {
	return a.name
}

func (a *Apt) listPackage(name string) (string, error) {
	return execCmd("apt", "list", "-qq", name)
}

// PackageInstalled tests if the given package name is installed
func (a *Apt) PackageInstalled(pkg *Package) (bool, error) {
	if !pkg.SystemPackage {
		if pkg.InstallCheck != nil {
			return pkg.InstallCheck(), nil
		}
		return false, nil
	}
	output, err := a.listPackage(pkg.Name)
	// apt list -qq returns "all" if you have packages installed globally and locally
	return strings.Contains(output, "installed") || strings.Contains(output, " all"), err
}

// PackageAvailable tests if the given package is available for installation
func (a *Apt) PackageAvailable(pkg *Package) (bool, error) {
	if !pkg.SystemPackage {
		return true, nil
	}
	output, err := a.listPackage(pkg.Name)
	// We add a space to ensure we get a full match, not partial match
	escapechars, _ := regexp.Compile(`\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])`)
	escapechars.ReplaceAllString(output, "")
	installed := strings.HasPrefix(output, pkg.Name)
	a.getPackageVersion(pkg, output)
	return installed, err
}

// InstallCommand returns the package manager specific command to install a package
func (a *Apt) InstallCommand(pkg *Package) string {
	if !pkg.SystemPackage {
		return pkg.InstallCommand
	}
	return "sudo apt install " + pkg.Name
}

func (a *Apt) getPackageVersion(pkg *Package, output string) {
	splitOutput := strings.Split(output, " ")
	if len(splitOutput) > 1 {
		pkg.Version = splitOutput[1]
	}
}
