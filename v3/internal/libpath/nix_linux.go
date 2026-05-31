//go:build linux

package libpath

import "os"

// getNixLibPaths returns cached library paths for Nix/NixOS installations.
func getNixLibPaths() []string {
	return cache.getNix()
}

// discoverNixLibPaths scans for Nix library paths.
func discoverNixLibPaths() []string {
	var paths []string

	nixProfileLib := os.ExpandEnv("$HOME/.nix-profile/lib")
	if _, err := os.Stat(nixProfileLib); err == nil {
		paths = append(paths, nixProfileLib)
	}

	// System Nix store - packages expose libs through profiles
	nixStoreLib := "/run/current-system/sw/lib"
	if _, err := os.Stat(nixStoreLib); err == nil {
		paths = append(paths, nixStoreLib)
	}

	return paths
}
