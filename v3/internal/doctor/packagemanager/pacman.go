//go:build linux

package packagemanager

import (
	"os/exec"
	"regexp"
	"strings"
)

// Pacman represents the Pacman package manager
type Pacman struct {
	name string
	osid string
}

// NewPacman creates a new Pacman instance
func NewPacman(osid string) *Pacman {
	return &Pacman{
		name: "pacman",
		osid: osid,
	}
}

// Packages returns the libraries that we need for Wails to compile
// They will potentially differ on different distributions or versions
func (p *Pacman) Packages() Packagemap {
	return Packagemap{
		"gtk3": []*Package{
			{Name: "gtk3", SystemPackage: true, Library: true},
		},
		"webkit2gtk": []*Package{
			{Name: "webkit2gtk-4.1", SystemPackage: true, Library: true},
		},
		"gcc": []*Package{
			{Name: "gcc", SystemPackage: true},
		},
		"pkg-config": []*Package{
			{Name: "pkgconf", SystemPackage: true},
		},
		"npm": []*Package{
			{Name: "npm", SystemPackage: true},
		},
	}
}

// Name returns the name of the package manager
func (p *Pacman) Name() string {
	return p.name
}

// PackageInstalled tests if the given package name is installed
func (p *Pacman) PackageInstalled(pkg *Package) (bool, error) {
	if !pkg.SystemPackage {
		if pkg.InstallCheck != nil {
			return pkg.InstallCheck(), nil
		}
		return false, nil
	}
	stdout, err := execCmd("pacman", "-Q", pkg.Name)
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if ok {
			return false, nil
		}
		return false, err
	}

	splitoutput := strings.Split(stdout, "\n")
	for _, line := range splitoutput {
		if strings.HasPrefix(line, pkg.Name) {
			splitline := strings.Split(line, " ")
			pkg.Version = strings.TrimSpace(splitline[1])
		}
	}

	return true, err
}

// PackageAvailable tests if the given package is available for installation
func (p *Pacman) PackageAvailable(pkg *Package) (bool, error) {
	if pkg.SystemPackage == false {
		return false, nil
	}
	output, err := execCmd("pacman", "-Si", pkg.Name)
	// We add a space to ensure we get a full match, not partial match
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if ok {
			return false, nil
		}
		return false, err
	}

	reg := regexp.MustCompile(`.*Version.*?:\s+(.*)`)
	matches := reg.FindStringSubmatch(output)
	pkg.Version = ""
	noOfMatches := len(matches)
	if noOfMatches > 1 {
		pkg.Version = strings.TrimSpace(matches[1])
	}

	return true, nil
}

// InstallCommand returns the package manager specific command to install a package
func (p *Pacman) InstallCommand(pkg *Package) string {
	if pkg.SystemPackage == false {
		return pkg.InstallCommand
	}
	return "sudo pacman -S " + pkg.Name
}
