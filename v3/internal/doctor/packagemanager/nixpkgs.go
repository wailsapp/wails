//go:build linux

package packagemanager

import (
	"encoding/json"
)

// Nixpkgs represents the Nixpkgs manager
type Nixpkgs struct {
	name string
	osid string
}

type NixPackageDetail struct {
	Name    string
	Pname   string
	Version string
}

var available map[string]NixPackageDetail

// NewNixpkgs creates a new Nixpkgs instance
func NewNixpkgs(osid string) *Nixpkgs {
	available = map[string]NixPackageDetail{}

	return &Nixpkgs{
		name: "nixpkgs",
		osid: osid,
	}
}

// Packages returns the libraries that we need for Wails to compile
// They will potentially differ on different distributions or versions
func (n *Nixpkgs) Packages() Packagemap {
	// Currently, only support checking the default channel.
	channel := "nixpkgs"
	if n.osid == "nixos" {
		channel = "nixos"
	}

	return Packagemap{
		"gtk3": []*Package{
			{Name: channel + ".gtk3", SystemPackage: true, Library: true},
		},
		"webkit2gtk": []*Package{
			{Name: channel + ".webkitgtk", SystemPackage: true, Library: true},
		},
		"gcc": []*Package{
			{Name: channel + ".gcc", SystemPackage: true},
		},
		"pkg-config": []*Package{
			{Name: channel + ".pkg-config", SystemPackage: true},
		},
		"npm": []*Package{
			{Name: channel + ".nodejs", SystemPackage: true},
		},
	}
}

// Name returns the name of the package manager
func (n *Nixpkgs) Name() string {
	return n.name
}

// PackageInstalled tests if the given package name is installed
func (n *Nixpkgs) PackageInstalled(pkg *Package) (bool, error) {
	if !pkg.SystemPackage {
		if pkg.InstallCheck != nil {
			return pkg.InstallCheck(), nil
		}
		return false, nil
	}

	stdout, err := execCmd("nix-env", "--json", "-qA", pkg.Name)
	if err != nil {
		return false, nil
	}

	var attributes map[string]NixPackageDetail
	err = json.Unmarshal([]byte(stdout), &attributes)
	if err != nil {
		return false, err
	}

	// Did we get one?
	installed := false
	for attribute, detail := range attributes {
		if attribute == pkg.Name {
			installed = true
			pkg.Version = detail.Version
		}
		break
	}

	// If on NixOS, package may be installed via system config, so check the nix store.
	detail, ok := available[pkg.Name]
	if !installed && n.osid == "nixos" && ok {
		cmd := "nix-store --query --requisites /run/current-system | cut -d- -f2- | sort | uniq | grep '^" + detail.Pname + "'"

		if pkg.Library {
			cmd += " | grep 'dev$'"
		}

		stdout, err = execCmd("sh", "-c", cmd)
		if err != nil {
			return false, nil
		}

		if len(stdout) > 0 {
			installed = true
		}
	}

	return installed, nil
}

// PackageAvailable tests if the given package is available for installation
func (n *Nixpkgs) PackageAvailable(pkg *Package) (bool, error) {
	if pkg.SystemPackage == false {
		return false, nil
	}

	stdout, err := execCmd("nix-env", "--json", "-qaA", pkg.Name)
	if err != nil {
		return false, nil
	}

	var attributes map[string]NixPackageDetail
	err = json.Unmarshal([]byte(stdout), &attributes)
	if err != nil {
		return false, err
	}

	// Grab first version.
	for attribute, detail := range attributes {
		pkg.Version = detail.Version
		available[attribute] = detail
		break
	}

	return len(pkg.Version) > 0, nil
}

// InstallCommand returns the package manager specific command to install a package
func (n *Nixpkgs) InstallCommand(pkg *Package) string {
	if pkg.SystemPackage == false {
		return pkg.InstallCommand
	}
	return "nix-env -iA " + pkg.Name
}
