package main

import (
	"fmt"
	"path/filepath"

	"github.com/wailsapp/wails/v2/internal/fs"
)

// GetAbsoluteTargetDir returns the absolute path of the target directory.
// If targetDir is empty, it returns an empty string (no target directory specified).
func GetAbsoluteTargetDir(targetDir string) (string, error) {
	if targetDir == "" {
		return "", nil
	}
	return filepath.Abs(targetDir)
}

// CheckDirectorySafety checks if the target directory is safe to initialize a project in.
// It returns an error if the directory exists and is non-empty, unless force is true.
// The absTargetDir should be an absolute path obtained from GetAbsoluteTargetDir.
func CheckDirectorySafety(absTargetDir string, force bool) error {
	// If no target directory is specified, the default behavior creates a new directory
	// with the project name, so we don't need to check safety
	if absTargetDir == "" {
		return nil
	}

	// If directory doesn't exist, it's safe
	if !fs.DirExists(absTargetDir) {
		return nil
	}

	// Check if directory is empty
	isEmpty, err := fs.DirIsEmpty(absTargetDir)
	if err != nil {
		return fmt.Errorf("failed to check target directory: %w", err)
	}

	// If directory is empty, it's safe
	if isEmpty {
		return nil
	}

	// Directory is non-empty - fail unless force flag is set
	if force {
		return nil
	}

	return fmt.Errorf("target directory '%s' is not empty. Aborting to prevent data loss. Use -f to force init in a non-empty directory", absTargetDir)
}
