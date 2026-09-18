//go:build linux

package packagemanager

import (
	"strings"
)

// Xbps represents the Xbps manager
type Xbps struct {
	name string
	osid string
}

// NewXbps creates a new Xbps instance
func NewXbps(osid string) *Xbps {
	return &Xbps{
		name: "xbps",
		osid: osid,
	}
}

// Packages returns the libraries that we need for Wails to compile
func (x *Xbps) Packages() Packagemap {
	return Packagemap{
		"libgtk-3": []*Package{
			{Name: "gtk+3-devel", SystemPackage: true, Library: true},
		},
		"libwebkit": []*Package{
			{Name: "libwebkit2gtk41-devel", SystemPackage: true, Library: true},
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
		"docker": []*Package{
			{Name: "docker", SystemPackage: true, Optional: true},
		},
		"upx": []*Package{
			{Name: "upx", SystemPackage: true, Optional: true},
		},
	}
}

// Name returns the name of the package manager
func (x *Xbps) Name() string {
	return x.name
}

// PackageInstalled tests if the given package name is installed
func (x *Xbps) PackageInstalled(pkg *Package) (bool, error) {
	if !pkg.SystemPackage {
		return false, nil
	}

	output, err := execCmd("xbps-query", pkg.Name)
	if err != nil {
		return false, nil
	}

	installed := strings.Contains(output, "state: installed")
	if installed {
		x.getPackageVersion(pkg, output)
	}

	return installed, nil
}

// PackageAvailable tests if the given package is available for installation
func (x *Xbps) PackageAvailable(pkg *Package) (bool, error) {
	if !pkg.SystemPackage {
		return false, nil
	}

	output, err := execCmd("xbps-query", "-Rs", pkg.Name)
	if err != nil {
		// xbps-query -Rs exits non-zero when no packages match, writing nothing
		// to stdout (the error message goes to stderr). Treat empty-stdout errors
		// as "not available", the same way PackageInstalled silences all errors.
		if output == "" {
			return false, nil
		}
		return false, err
	}

	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		pkgField := fields[1]

		lastDash := strings.LastIndex(pkgField, "-")
		if lastDash <= 0 {
			continue
		}

		name := pkgField[:lastDash]

		if strings.EqualFold(name, pkg.Name) {
			x.getPackageVersion(pkg, pkgField)
			return true, nil
		}
	}

	return false, nil
}

// InstallCommand returns the package manager specific command to install a package
func (x *Xbps) InstallCommand(pkg *Package) string {
	if !pkg.SystemPackage {
		return pkg.InstallCommand
	}

	return "sudo xbps-install -Sy " + pkg.Name
}

func (x *Xbps) getPackageVersion(pkg *Package, output string) {
	lowerName := strings.ToLower(pkg.Name)
	prefix := lowerName + "-"
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(strings.ToLower(line), lowerName) {
			fields := strings.Fields(line)
			for _, f := range fields {
				if strings.HasPrefix(strings.ToLower(f), prefix) {
					pkg.Version = f[len(prefix):]
					return
				}
			}
		}
	}
}
