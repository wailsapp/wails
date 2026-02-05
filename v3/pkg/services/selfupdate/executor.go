package selfupdate

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Executor handles the platform-specific update installation process.
type Executor struct {
	config      *ExecutorConfig
	prepareFunc PrepareFunc
}

// ExecutorConfig holds configuration for the update executor.
type ExecutorConfig struct {
	// TargetPath is the path to the executable to update.
	// If empty, defaults to the current executable.
	TargetPath string

	// BackupPath is where to store the backup of the current executable.
	// If empty, defaults to TargetPath + ".backup"
	BackupPath string

	// TempDir is the directory for temporary files during update.
	// If empty, defaults to os.TempDir().
	TempDir string

	// KeepBackup controls whether to keep the backup after successful update.
	// If false, the backup is deleted after restart.
	KeepBackup bool

	// PrepareFunc is called after download to prepare the update.
	// This handles archive extraction, bundle navigation, etc.
	PrepareFunc PrepareFunc
}

// NewExecutor creates a new Executor with the given configuration.
func NewExecutor(config *ExecutorConfig) (*Executor, error) {
	if config == nil {
		config = &ExecutorConfig{}
	}

	// Default to current executable
	if config.TargetPath == "" {
		exe, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("failed to get executable path: %w", err)
		}
		// Resolve symlinks
		exe, err = filepath.EvalSymlinks(exe)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve executable path: %w", err)
		}
		config.TargetPath = exe
	}

	// Default backup path
	if config.BackupPath == "" {
		config.BackupPath = config.TargetPath + ".backup"
	}

	// Default temp dir
	if config.TempDir == "" {
		config.TempDir = os.TempDir()
	}

	return &Executor{
		config:      config,
		prepareFunc: config.PrepareFunc,
	}, nil
}

// Apply applies the update from the given reader.
// This is the main entry point for installing an update.
//
// On Windows, this will rename the current exe and copy the new one.
// On other platforms, the binary is replaced atomically when possible.
func (e *Executor) Apply(ctx context.Context, updateData io.Reader, providerConfig *ProviderConfig) error {
	// Create temp file for the download
	tempFile, err := os.CreateTemp(e.config.TempDir, "wails-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tempPath := tempFile.Name()

	// Copy update data to temp file
	if _, err := io.Copy(tempFile, updateData); err != nil {
		tempFile.Close()
		os.Remove(tempPath)
		return fmt.Errorf("failed to write update to temp file: %w", err)
	}
	tempFile.Close()

	// Prepare the update (extract archives, find binary, etc.)
	binaryPath := tempPath
	if e.prepareFunc != nil {
		prepared, err := e.prepareFunc(ctx, tempPath, providerConfig)
		if err != nil {
			os.Remove(tempPath)
			return fmt.Errorf("failed to prepare update: %w", err)
		}
		binaryPath = prepared
	}

	// Apply based on platform
	switch runtime.GOOS {
	case "windows":
		return e.applyWindows(ctx, binaryPath)
	default:
		return e.applyUnix(ctx, binaryPath)
	}
}

// applyUnix applies the update on Unix-like systems (Linux, macOS).
// This performs an atomic replacement of the binary when possible.
func (e *Executor) applyUnix(_ context.Context, newBinaryPath string) error {
	targetPath := e.config.TargetPath

	// For macOS .app bundles, we need to find the actual target
	if runtime.GOOS == "darwin" {
		targetPath = e.resolveMacOSTarget(targetPath, newBinaryPath)
	}

	// Get permissions from the current binary
	info, err := os.Stat(targetPath)
	if err != nil {
		return fmt.Errorf("failed to stat target: %w", err)
	}

	// Create backup of current binary
	if err := e.createBackup(targetPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Set permissions on new binary
	if err := os.Chmod(newBinaryPath, info.Mode()); err != nil {
		if restoreErr := e.restoreBackup(targetPath); restoreErr != nil {
			return fmt.Errorf("failed to set permissions: %w (restore also failed: %v)", err, restoreErr)
		}
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// Rename new binary to target (atomic on same filesystem)
	if err := os.Rename(newBinaryPath, targetPath); err != nil {
		// If rename fails (cross-filesystem), use atomic copy via temp file (H2)
		if err := e.atomicCopyFile(newBinaryPath, targetPath, info.Mode()); err != nil {
			if restoreErr := e.restoreBackup(targetPath); restoreErr != nil {
				return fmt.Errorf("failed to install update: %w (restore also failed: %v)", err, restoreErr)
			}
			return fmt.Errorf("failed to install update: %w", err)
		}
		os.Remove(newBinaryPath)
	}

	return nil
}

// applyWindows applies the update on Windows.
// Windows locks running executables, so we rename the current one first.
func (e *Executor) applyWindows(_ context.Context, newBinaryPath string) error {
	targetPath := e.config.TargetPath
	oldPath := targetPath + ".old"

	// Remove any previous .old file
	os.Remove(oldPath)

	// Rename current exe to .old
	if err := os.Rename(targetPath, oldPath); err != nil {
		return fmt.Errorf("failed to rename current executable: %w", err)
	}

	// Copy new binary to target
	if err := copyFile(newBinaryPath, targetPath); err != nil {
		// Try to restore (H4: check error)
		if restoreErr := os.Rename(oldPath, targetPath); restoreErr != nil {
			return fmt.Errorf("failed to install update: %w (CRITICAL: restore also failed: %v; original binary is at %s)",
				err, restoreErr, oldPath)
		}
		return fmt.Errorf("failed to install update: %w", err)
	}

	// Clean up temp file
	os.Remove(newBinaryPath)

	return nil
}

// resolveMacOSTarget resolves the actual binary path for macOS.
// If we're updating an .app bundle, we need to find the correct location.
func (e *Executor) resolveMacOSTarget(currentTarget, newBinaryPath string) string {
	// Check if the new binary is an .app bundle
	if filepath.Ext(newBinaryPath) == ".app" {
		// Find the current app bundle root
		bundleRoot := findMacOSBundleRoot(currentTarget)
		if bundleRoot != "" {
			return bundleRoot
		}
	}
	return currentTarget
}

// findMacOSBundleRoot finds the .app bundle root from an executable path.
// macOS apps are structured as: MyApp.app/Contents/MacOS/MyApp
func findMacOSBundleRoot(execPath string) string {
	// Walk up the directory tree looking for .app
	current := execPath
	for range 4 { // Max 4 levels up
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		if filepath.Ext(parent) == ".app" {
			return parent
		}
		current = parent
	}
	return ""
}

// GetExecutablePath returns the path to the current executable,
// resolving symlinks and handling macOS .app bundles.
func GetExecutablePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return "", fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	return exe, nil
}

// GetWorkingDirectory returns the working directory for the application.
// On macOS, this handles .app bundle structure by navigating up from
// the Contents/MacOS directory to the bundle root's parent.
func GetWorkingDirectory() (string, error) {
	exe, err := GetExecutablePath()
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(exe)

	// On macOS, navigate out of the .app bundle
	if runtime.GOOS == "darwin" {
		bundleRoot := findMacOSBundleRoot(exe)
		if bundleRoot != "" {
			dir = filepath.Dir(bundleRoot)
		}
	}

	return dir, nil
}

// createBackup creates a backup of the file at the given path.
func (e *Executor) createBackup(path string) error {
	return copyFile(path, e.config.BackupPath)
}

// restoreBackup restores the backup to the given path (H3: uses copy for cross-fs safety).
func (e *Executor) restoreBackup(path string) error {
	// Try rename first (same filesystem, atomic)
	if err := os.Rename(e.config.BackupPath, path); err == nil {
		return nil
	}
	// Fall back to copy (cross-filesystem)
	if err := copyFile(e.config.BackupPath, path); err != nil {
		return fmt.Errorf("failed to restore backup from %s: %w", e.config.BackupPath, err)
	}
	os.Remove(e.config.BackupPath)
	return nil
}

// atomicCopyFile copies src to dst atomically by writing to a temp file in the
// same directory as dst, then renaming (H2).
func (e *Executor) atomicCopyFile(src, dst string, mode os.FileMode) error {
	dir := filepath.Dir(dst)
	tmpFile, err := os.CreateTemp(dir, ".wails-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file in target dir: %w", err)
	}
	tmpPath := tmpFile.Name()

	srcFile, err := os.Open(src)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return err
	}

	if _, err := io.Copy(tmpFile, srcFile); err != nil {
		srcFile.Close()
		tmpFile.Close()
		os.Remove(tmpPath)
		return err
	}
	srcFile.Close()

	if err := tmpFile.Chmod(mode); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return err
	}

	if err := tmpFile.Sync(); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return err
	}
	tmpFile.Close()

	// Rename is atomic on the same filesystem (guaranteed since temp is in same dir)
	if err := os.Rename(tmpPath, dst); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return nil
}

// CleanupOldVersions removes any leftover files from previous updates.
// This should be called on application startup.
func CleanupOldVersions() {
	exe, err := GetExecutablePath()
	if err != nil {
		return
	}

	// Clean up .old and .backup files
	for _, suffix := range []string{".old", ".backup"} {
		oldPath := exe + suffix
		if _, err := os.Stat(oldPath); err == nil {
			// File exists, try to remove it
			// On Windows, this might fail if we just started,
			// so we'll retry a few times
			for range 3 {
				if err := os.Remove(oldPath); err == nil {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

// copyFile copies a file from src to dst, preserving permissions.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return dstFile.Sync()
}

// CanUpdate checks if the current process has permission to update itself.
// It verifies write access to the directory containing the executable,
// since updates use rename which requires directory write permission (H1).
func CanUpdate() bool {
	exe, err := GetExecutablePath()
	if err != nil {
		return false
	}

	// Check if we can write to the directory containing the executable
	dir := filepath.Dir(exe)
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}

	// Check if directory is writable by attempting to create a temp file
	tmpFile, err := os.CreateTemp(dir, ".wails-update-check-*")
	if err != nil {
		return false
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	os.Remove(tmpPath)

	_ = info // used for stat check above
	return true
}
