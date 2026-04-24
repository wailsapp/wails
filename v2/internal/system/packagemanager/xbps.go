//go:build linux
// +build linux

package packagemanager

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/wailsapp/wails/v2/internal/shell"
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
func (x *Xbps) Packages() packagemap {
	return packagemap{
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
	if pkg.SystemPackage == false {
		return false, nil
	}
	cmd := exec.Command("xbps-query", pkg.Name)
	var stdo, stde bytes.Buffer
	cmd.Stdout = &stdo
	cmd.Stderr = &stde
	cmd.Env = append(os.Environ(), "LANGUAGE=en")
	err := cmd.Run()
	// xbps-query exits 0 only if the package is found and installed
	installed := err == nil && strings.Contains(stdo.String(), "state: installed")
	if installed {
		x.getPackageVersion(pkg, stdo.String())
	}
	return installed, nil
}

// PackageAvailable tests if the given package is available for installation
func (x *Xbps) PackageAvailable(pkg *Package) (bool, error) {
	if pkg.SystemPackage == false {
		return false, nil
	}
	// -Rs searches the remote repositories
	stdout, _, err := shell.RunCommand(".", "xbps-query", "-Rs", pkg.Name)
	available := false
	for _, line := range strings.Split(stdout, "\n") {
		// Each line is like: "[-] pkgname-version ..."  or  "[*] pkgname-version ..."
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		// fields[1] is "pkgname-version", strip the version to compare
		pkgfield := fields[1]
		dashIdx := strings.LastIndex(pkgfield, "-")
		if dashIdx > 0 && strings.EqualFold(pkgfield[:dashIdx], pkg.Name) {
			available = true
			x.getPackageVersion(pkg, pkgfield)
			break
		}
	}
	return available, err
}

// InstallCommand returns the package manager specific command to install a package
func (x *Xbps) InstallCommand(pkg *Package) string {
	if pkg.SystemPackage == false {
		return pkg.InstallCommand[x.osid]
	}
	return "sudo xbps-install -S " + pkg.Name
}

func (x *Xbps) getPackageVersion(pkg *Package, output string) {
	// xbps-query output: "pkgname-version_revision" or inline in -Rs output
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, pkg.Name) {
			fields := strings.Fields(line)
			for _, f := range fields {
				if strings.HasPrefix(f, pkg.Name+"-") {
					// strip pkgname- prefix, keep version_revision
					pkg.Version = strings.TrimPrefix(f, pkg.Name+"-")
					return
				}
			}
		}
	}
}
