//go:build linux

package updater

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// applyUpdate applies the downloaded update on Linux
func applyUpdate(ctx context.Context, downloadPath string, info *UpdateInfo) error {
	// Determine if this is a patch or full update
	isPatch := info.PatchURL != "" && strings.HasSuffix(downloadPath, ".patch")

	if isPatch {
		return applyPatch(ctx, downloadPath, info)
	}

	return applyFullUpdate(ctx, downloadPath, info)
}

// applyFullUpdate extracts and replaces the application
func applyFullUpdate(ctx context.Context, archivePath string, info *UpdateInfo) error {
	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Resolve symlinks
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	// Check if this is an AppImage
	if isAppImage(execPath) {
		return applyAppImageUpdate(ctx, archivePath, execPath, info)
	}

	// Regular binary update
	return applyBinaryUpdate(ctx, archivePath, execPath, info)
}

// isAppImage checks if the executable is an AppImage
func isAppImage(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".appimage") ||
		os.Getenv("APPIMAGE") != ""
}

// applyAppImageUpdate replaces an AppImage with a new version
func applyAppImageUpdate(ctx context.Context, archivePath, execPath string, info *UpdateInfo) error {
	// Create temp directory for extraction
	tempDir, err := os.MkdirTemp("", "wails-update-extract-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Extract the archive
	if err := extractTarGz(archivePath, tempDir); err != nil {
		return fmt.Errorf("failed to extract update: %w", err)
	}

	// Find the new AppImage
	newAppImage, err := findAppImage(tempDir)
	if err != nil {
		return fmt.Errorf("failed to find AppImage in update: %w", err)
	}

	// Make it executable
	if err := os.Chmod(newAppImage, 0755); err != nil {
		return fmt.Errorf("failed to set executable permission: %w", err)
	}

	// Create backup
	backupPath := execPath + ".backup"
	if err := os.Rename(execPath, backupPath); err != nil {
		return fmt.Errorf("failed to backup current AppImage: %w", err)
	}

	// Move new AppImage into place
	if err := os.Rename(newAppImage, execPath); err != nil {
		os.Rename(backupPath, execPath)
		return fmt.Errorf("failed to install update: %w", err)
	}

	// Remove backup
	os.Remove(backupPath)

	// Restart
	return restartApp(execPath)
}

// applyBinaryUpdate replaces a regular binary
func applyBinaryUpdate(ctx context.Context, archivePath, execPath string, info *UpdateInfo) error {
	// Create temp directory for extraction
	tempDir, err := os.MkdirTemp("", "wails-update-extract-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Extract the archive
	if err := extractTarGz(archivePath, tempDir); err != nil {
		return fmt.Errorf("failed to extract update: %w", err)
	}

	// Find the new binary
	appName := filepath.Base(execPath)
	newBinary := filepath.Join(tempDir, appName)
	if _, err := os.Stat(newBinary); os.IsNotExist(err) {
		// Try to find any executable
		newBinary, err = findExecutable(tempDir)
		if err != nil {
			return fmt.Errorf("failed to find executable in update: %w", err)
		}
	}

	// Make it executable
	if err := os.Chmod(newBinary, 0755); err != nil {
		return fmt.Errorf("failed to set executable permission: %w", err)
	}

	// Create backup
	backupPath := execPath + ".backup"
	if err := os.Rename(execPath, backupPath); err != nil {
		return fmt.Errorf("failed to backup current binary: %w", err)
	}

	// Move new binary into place
	if err := os.Rename(newBinary, execPath); err != nil {
		os.Rename(backupPath, execPath)
		return fmt.Errorf("failed to install update: %w", err)
	}

	// Remove backup
	os.Remove(backupPath)

	// Restart
	return restartApp(execPath)
}

// applyPatch applies a bsdiff patch
func applyPatch(ctx context.Context, patchPath string, info *UpdateInfo) error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	// Create temp file for patched binary
	tempFile, err := os.CreateTemp("", "wails-patched-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempPath)

	// Apply the patch using bspatch
	cmd := exec.CommandContext(ctx, "bspatch", execPath, tempPath, patchPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to apply patch: %w", err)
	}

	// Make it executable
	if err := os.Chmod(tempPath, 0755); err != nil {
		return fmt.Errorf("failed to set executable permission: %w", err)
	}

	// Create backup
	backupPath := execPath + ".backup"
	if err := os.Rename(execPath, backupPath); err != nil {
		return fmt.Errorf("failed to backup current binary: %w", err)
	}

	// Move patched binary into place
	if err := os.Rename(tempPath, execPath); err != nil {
		os.Rename(backupPath, execPath)
		return fmt.Errorf("failed to install patched binary: %w", err)
	}

	// Remove backup
	os.Remove(backupPath)

	// Restart
	return restartApp(execPath)
}

// findAppImage finds an AppImage file in a directory
func findAppImage(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if strings.HasSuffix(strings.ToLower(entry.Name()), ".appimage") {
			return filepath.Join(dir, entry.Name()), nil
		}
	}

	return "", fmt.Errorf("no AppImage found")
}

// findExecutable finds an executable file in a directory
func findExecutable(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		if info.Mode()&0111 != 0 {
			return path, nil
		}
	}

	return "", fmt.Errorf("no executable found")
}

// extractTarGz extracts a .tar.gz archive
func extractTarGz(archivePath, destDir string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		targetPath := filepath.Join(destDir, header.Name)

		// Security: prevent path traversal
		if !strings.HasPrefix(targetPath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return err
			}

		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return err
			}

			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()

		case tar.TypeSymlink:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return err
			}
			if err := os.Symlink(header.Linkname, targetPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// restartApp launches the app and exits the current process
func restartApp(appPath string) error {
	cmd := exec.Command(appPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to restart app: %w", err)
	}

	os.Exit(0)
	return nil
}
