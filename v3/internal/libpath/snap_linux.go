//go:build linux

package libpath

import (
	"os"
	"path/filepath"
)

// getSnapLibPaths returns cached library paths from installed Snap packages.
func getSnapLibPaths() []string {
	return cache.getSnap()
}

// discoverSnapLibPaths scans for Snap package library directories.
// Scans /snap/*/current/usr/lib* directories.
func discoverSnapLibPaths() []string {
	var paths []string

	snapDir := "/snap"
	if _, err := os.Stat(snapDir); err != nil {
		return paths
	}

	// Find all snap packages with lib directories
	patterns := []string{
		filepath.Join(snapDir, "*", "current", "usr", "lib"),
		filepath.Join(snapDir, "*", "current", "usr", "lib64"),
		filepath.Join(snapDir, "*", "current", "usr", "lib", "*-linux-gnu"),
		filepath.Join(snapDir, "*", "current", "lib"),
		filepath.Join(snapDir, "*", "current", "lib", "*-linux-gnu"),
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err == nil {
			paths = append(paths, matches...)
		}
	}

	return paths
}
