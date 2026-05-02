//go:build linux

package packagemanager

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

type PackageManager interface {
	Name() string
	Packages() Packagemap
	PackageInstalled(*Package) (bool, error)
	PackageAvailable(*Package) (bool, error)
	InstallCommand(*Package) string
}

type Dependency struct {
	Name           string
	PackageName    string
	Installed      bool
	InstallCommand string
	Version        string
	Optional       bool
	External       bool
}

type DependencyList []*Dependency
