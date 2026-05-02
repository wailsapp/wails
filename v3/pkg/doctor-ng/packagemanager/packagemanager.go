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

var pmCommands = []string{
	"eopkg",
	"apt",
	"dnf",
	"pacman",
	"emerge",
	"zypper",
	"nix-env",
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func Find(osid string) PackageManager {
	for _, pmname := range pmCommands {
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

	sort.Slice(dependencies, func(i, j int) bool {
		return dependencies[i].Name < dependencies[j].Name
	})

	return dependencies, nil
}

func AppVersion(name string) string {
	switch name {
	case "gcc":
		return gccVersion()
	case "pkg-config":
		return pkgConfigVersion()
	case "npm":
		return npmVersion()
	}
	return ""
}

func gccVersion() string {
	version, err := execCmd("gcc", "-dumpfullversion")
	if err != nil {
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
