//go:build linux
// +build linux

package packagemanager

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/wailsapp/wails/v2/internal/shell"
)

// Emerge represents the Emerge package manager
type Emerge struct {
	name string
	osid string
}

// NewEmerge creates a new Emerge instance
func NewEmerge(osid string) *Emerge {
	return &Emerge{
		name: "emerge",
		osid: osid,
	}
}

// Packages returns the libraries that we need for Wails to compile
// They will potentially differ on different distributions or versions
func (e *Emerge) Packages() packagemap {
	return packagemap{
		"libgtk-3": []*Package{
			{Name: "x11-libs/gtk+", SystemPackage: true, Library: true},
		},
		"libwebkit": []*Package{
			{Name: "net-libs/webkit-gtk", SystemPackage: true, Library: true},
		},
		"gcc": []*Package{
			{Name: "sys-devel/gcc", SystemPackage: true},
		},
		"pkg-config": []*Package{
			{Name: "dev-util/pkgconf", SystemPackage: true},
		},
		"npm": []*Package{
			{Name: "net-libs/nodejs", SystemPackage: true},
		},
		"docker": []*Package{
			{Name: "app-emulation/docker", SystemPackage: true, Optional: true},
		},
	}
}

// Name returns the name of the package manager
func (e *Emerge) Name() string {
	return e.name
}

// PackageInstalled tests if the given package name is installed
func (e *Emerge) PackageInstalled(pkg *Package) (bool, error) {
	if pkg.SystemPackage == false {
		return false, nil
	}
	stdout, _, err := shell.RunCommand(".", "emerge", "-s", pkg.Name+"$")
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if ok {
			return false, nil
		}
		return false, err
	}

	regex := `.*\*\s+` + regexp.QuoteMeta(pkg.Name) + `\n(?:\S|\s)+?Latest version installed: (.*)`
	installedRegex := regexp.MustCompile(regex)
	matches := installedRegex.FindStringSubmatch(stdout)
	pkg.Version = ""
	noOfMatches := len(matches)
	installed := false
	if noOfMatches > 1 && matches[1] != "[ Not Installed ]" {
		installed = true
		pkg.Version = strings.TrimSpace(matches[1])
	}
	return installed, err
}

// PackageAvailable tests if the given package is available for installation
func (e *Emerge) PackageAvailable(pkg *Package) (bool, error) {
	if pkg.SystemPackage == false {
		return false, nil
	}
	stdout, _, err := shell.RunCommand(".", "emerge", "-s", pkg.Name+"$")
	// We add a space to ensure we get a full match, not partial match
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if ok {
			return false, nil
		}
		return false, err
	}

	installedRegex := regexp.MustCompile(`.*\*\s+` + regexp.QuoteMeta(pkg.Name) + `\n(?:\S|\s)+?Latest version available: (.*)`)
	matches := installedRegex.FindStringSubmatch(stdout)
	pkg.Version = ""
	noOfMatches := len(matches)
	available := false
	if noOfMatches > 1 {
		available = true
		pkg.Version = strings.TrimSpace(matches[1])
	}
	return available, nil
}

// InstallCommand returns the package manager specific command to install a package
func (e *Emerge) InstallCommand(pkg *Package) string {
	if pkg.SystemPackage == false {
		return pkg.InstallCommand[e.osid]
	}
	return "sudo emerge " + pkg.Name
}
