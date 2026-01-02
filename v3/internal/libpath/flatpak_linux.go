//go:build linux

package libpath

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// getFlatpakLibPaths returns cached library paths from installed Flatpak runtimes.
func getFlatpakLibPaths() []string {
	return cache.getFlatpak()
}

// discoverFlatpakLibPaths scans for Flatpak runtime library directories.
// Uses `flatpak --installations` and scans for runtime lib directories.
func discoverFlatpakLibPaths() []string {
	var paths []string

	// Get system and user installation directories
	installDirs := []string{
		"/var/lib/flatpak",                         // System default
		os.ExpandEnv("$HOME/.local/share/flatpak"), // User default
	}

	// Try to get actual installation path from flatpak
	if out, err := exec.Command("flatpak", "--installations").Output(); err == nil {
		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			if line != "" {
				installDirs = append(installDirs, line)
			}
		}
	}

	// Scan for runtime lib directories
	for _, installDir := range installDirs {
		runtimeDir := filepath.Join(installDir, "runtime")
		if _, err := os.Stat(runtimeDir); err != nil {
			continue
		}

		// Look for lib directories in runtimes
		// Structure: runtime/<name>/<arch>/<version>/<hash>/files/lib
		matches, err := filepath.Glob(filepath.Join(runtimeDir, "*", "*", "*", "*", "files", "lib"))
		if err == nil {
			paths = append(paths, matches...)
		}
	}

	return paths
}
