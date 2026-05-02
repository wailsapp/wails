//go:build linux

package packagemanager

import (
	"bytes"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func execCmd(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var stdo, stde bytes.Buffer
	cmd.Stdout = &stdo
	cmd.Stderr = &stde
	cmd.Env = append(os.Environ(), "LANGUAGE=en_US.utf-8")
	err := cmd.Run()
	return stdo.String(), err
}

// A list of package manager commands
var pmcommands = []string{
	"eopkg",
	"apt",
	"dnf",
	"pacman",
	"emerge",
	"zypper",
	"nix-env",
}

// commandExists returns true if the given command can be found on the shell
func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	if err != nil {
		return false
	}
	return true
}

// Find will attempt to find the system package manager
func Find(osid string) PackageManager {

	// Loop over pmcommands
	for _, pmname := range pmcommands {
		if commandExists(pmname) {
			return newPackageManager(pmname, osid)
		}
	}
	return nil
}

func newPackageManager(pmname string, osid string) PackageManager {
	switch pmname {
	case "eopkg":
		return NewEopkg(osid)
	case "apt":
		return NewApt(osid)
	case "dnf":
		return NewDnf(osid)
	case "pacman":
		return NewPacman(osid)
	case "emerge":
		return NewEmerge(osid)
	case "zypper":
		return NewZypper(osid)
	case "nix-env":
		return NewNixpkgs(osid)
	}
	return nil
}

// Dependencies scans the system for required dependencies
// Returns a list of dependencies search for, whether they were found
// and whether they were installed
func Dependencies(p PackageManager) (DependencyList, error) {

	var dependencies DependencyList

	for name, packages := range p.Packages() {
		dependency := &Dependency{Name: name}
		for _, pkg := range packages {
			dependency.Optional = pkg.Optional
			dependency.External = !pkg.SystemPackage
			dependency.InstallCommand = p.InstallCommand(pkg)
			packageavailable, err := p.PackageAvailable(pkg)
			if err != nil {
				return nil, err
			}
			if packageavailable {
				dependency.Version = pkg.Version
				dependency.PackageName = pkg.Name
				installed, err := p.PackageInstalled(pkg)
				if err != nil {
					return nil, err
				}
				if installed {
					dependency.Installed = true
					dependency.Version = pkg.Version
					if !pkg.SystemPackage {
						dependency.Version = AppVersion(name)
					}
				} else {
					dependency.InstallCommand = p.InstallCommand(pkg)
				}
				break
			}
		}
		dependencies = append(dependencies, dependency)
	}

	// Sort dependencies
	sort.Slice(dependencies, func(i, j int) bool {
		return dependencies[i].Name < dependencies[j].Name
	})

	return dependencies, nil
}

// AppVersion returns the version for application related to the given package
func AppVersion(name string) string {

	if name == "gcc" {
		return gccVersion()
	}

	if name == "pkg-config" {
		return pkgConfigVersion()
	}

	if name == "npm" {
		return npmVersion()
	}

	return ""

}

func gccVersion() string {

	var version string
	var err error

	// Try "-dumpfullversion"
	version, err = execCmd("gcc", "-dumpfullversion")
	if err != nil {

		// Try -dumpversion
		// We ignore the error as this function is not for testing whether the
		// application exists, only that we can get the version number
		dumpversion, err := execCmd("gcc", "-dumpversion")
		if err == nil {
			version = dumpversion
		}
	}
	return strings.TrimSpace(version)
}

func pkgConfigVersion() string {
	version, _ := execCmd("pkg-config", "--version")
	return strings.TrimSpace(version)
}

func npmVersion() string {
	version, _ := execCmd("npm", "--version")
	return strings.TrimSpace(version)
}
