//go:build linux

package packagemanager

import (
	"regexp"
	"strings"
)

type Eopkg struct {
	name string
	osid string
}

// NewEopkg creates a new Eopkg instance
func NewEopkg(osid string) *Eopkg {
	result := &Eopkg{
		name: "eopkg",
		osid: osid,
	}
	result.intialiseName()
	return result
}

// Packages returns the packages that we need for Wails to compile
// They will potentially differ on different distributions or versions
func (e *Eopkg) Packages() Packagemap {
	return Packagemap{
		"gtk3": []*Package{
			{Name: "libgtk-3-devel", SystemPackage: true, Library: true},
		},
		"webkit2gtk": []*Package{
			{Name: "libwebkit-gtk-devel", SystemPackage: true, Library: true},
		},
		"gcc": []*Package{
			{Name: "gcc", SystemPackage: true},
		},
		"pkg-config": []*Package{
			{Name: "pkg-config", SystemPackage: true},
		},
		"npm": []*Package{
			{Name: "nodejs", SystemPackage: true},
		},
	}
}

// Name returns the name of the package manager
func (e *Eopkg) Name() string {
	return e.name
}

// PackageInstalled tests if the given package is installed
func (e *Eopkg) PackageInstalled(pkg *Package) (bool, error) {
	if !pkg.SystemPackage {
		if pkg.InstallCheck != nil {
			return pkg.InstallCheck(), nil
		}
		return false, nil
	}
	stdout, err := execCmd("eopkg", "info", pkg.Name)
	return strings.HasPrefix(stdout, "Installed"), err
}

// PackageAvailable tests if the given package is available for installation
func (e *Eopkg) PackageAvailable(pkg *Package) (bool, error) {
	if pkg.SystemPackage == false {
		return false, nil
	}
	stdout, err := execCmd("eopkg", "info", pkg.Name)
	// We add a space to ensure we get a full match, not partial match
	output := e.removeEscapeSequences(stdout)
	installed := strings.Contains(output, "Package found in Solus repository")
	e.getPackageVersion(pkg, output)
	return installed, err
}

// InstallCommand returns the package manager specific command to install a package
func (e *Eopkg) InstallCommand(pkg *Package) string {
	if pkg.SystemPackage == false {
		return pkg.InstallCommand
	}
	return "sudo eopkg it " + pkg.Name
}

func (e *Eopkg) removeEscapeSequences(in string) string {
	escapechars, _ := regexp.Compile(`\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])`)
	return escapechars.ReplaceAllString(in, "")
}

func (e *Eopkg) intialiseName() {
	result := "eopkg"
	stdout, err := execCmd("eopkg", "--version")
	if err == nil {
		result = strings.TrimSpace(stdout)
	}
	e.name = result
}

func (e *Eopkg) getPackageVersion(pkg *Package, output string) {

	versionRegex := regexp.MustCompile(`.*Name.*version:\s+(.*)+, release: (.*)`)
	matches := versionRegex.FindStringSubmatch(output)
	pkg.Version = ""
	noOfMatches := len(matches)
	if noOfMatches > 1 {
		pkg.Version = matches[1]
		if noOfMatches > 2 {
			pkg.Version += " (r" + matches[2] + ")"
		}
	}
}
