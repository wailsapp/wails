//go:build linux

package packagemanager

// Package contains information about a system package
type Package struct {
	Name           string
	Version        string
	InstallCommand string
	InstallCheck   func() bool
	SystemPackage  bool
	Library        bool
	Optional       bool
}

type Packagemap = map[string][]*Package

// PackageManager is a common interface across all package managers
type PackageManager interface {
	Name() string
	Packages() Packagemap
	PackageInstalled(*Package) (bool, error)
	PackageAvailable(*Package) (bool, error)
	InstallCommand(*Package) string
}

// Dependency represents a system package that we require
type Dependency struct {
	Name           string
	PackageName    string
	Installed      bool
	InstallCommand string
	Version        string
	Optional       bool
	External       bool
}

// DependencyList is a list of Dependency instances
type DependencyList []*Dependency

// InstallAllRequiredCommand returns the command you need to use to install all required dependencies
func (d DependencyList) InstallAllRequiredCommand() string {

	result := ""
	for _, dependency := range d {
		if !dependency.Installed && !dependency.Optional {
			result += "  - " + dependency.Name + ": " + dependency.InstallCommand + "\n"
		}
	}

	return result
}

// InstallAllOptionalCommand returns the command you need to use to install all optional dependencies
func (d DependencyList) InstallAllOptionalCommand() string {

	result := ""
	for _, dependency := range d {
		if !dependency.Installed && dependency.Optional {
			result += "  - " + dependency.Name + ": " + dependency.InstallCommand + "\n"
		}
	}

	return result
}
