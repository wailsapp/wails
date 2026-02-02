package main

import (
	"fmt"
	"path/filepath"

	"github.com/wailsapp/wails/v2/internal/fs"
)

// DirectorySafetyResult represents the result of a directory safety check
type DirectorySafetyResult int

const (
	// DirectorySafe indicates the directory is safe to use (empty or doesn't exist)
	DirectorySafe DirectorySafetyResult = iota
	// DirectoryNeedsConfirmation indicates the directory is non-empty and needs user confirmation
	DirectoryNeedsConfirmation
)

// DirectorySafetyError is returned when a directory is non-empty in CI mode
type DirectorySafetyError struct {
	TargetDir string
}

func (e *DirectorySafetyError) Error() string {
	return fmt.Sprintf("target directory '%s' is not empty. Aborting to prevent data loss. Use an empty directory or remove existing files first", e.TargetDir)
}

// CheckDirectorySafety checks if the target directory is safe to initialize a project in.
// It returns:
//   - DirectorySafe if the directory is empty or doesn't exist
//   - DirectoryNeedsConfirmation if the directory is non-empty (and not in CI mode)
//   - DirectorySafetyError if in CI mode and directory is non-empty
//   - Other errors for filesystem issues
func CheckDirectorySafety(targetDir string, ciMode bool, force bool) (DirectorySafetyResult, error) {
	// If no target directory is specified, the default behavior creates a new directory
	// with the project name, so we don't need to check safety
	if targetDir == "" {
		return DirectorySafe, nil
	}

	// Get absolute path
	absTargetDir, err := filepath.Abs(targetDir)
	if err != nil {
		return DirectorySafe, err
	}

	// If directory doesn't exist, it's safe
	if !fs.DirExists(absTargetDir) {
		return DirectorySafe, nil
	}

	// Check if directory is empty
	isEmpty, err := fs.DirIsEmpty(absTargetDir)
	if err != nil {
		return DirectorySafe, fmt.Errorf("failed to check target directory: %w", err)
	}

	// If directory is empty, it's safe
	if isEmpty {
		return DirectorySafe, nil
	}

	// Directory is non-empty
	// If force flag is set, skip confirmation
	if force {
		return DirectorySafe, nil
	}

	// In CI mode, we can't prompt for confirmation, so return an error
	if ciMode {
		return DirectorySafe, &DirectorySafetyError{TargetDir: absTargetDir}
	}

	// Directory is non-empty and needs user confirmation
	return DirectoryNeedsConfirmation, nil
}

// GetAbsoluteTargetDir returns the absolute path of the target directory
func GetAbsoluteTargetDir(targetDir string) (string, error) {
	if targetDir == "" {
		return "", nil
	}
	return filepath.Abs(targetDir)
}
