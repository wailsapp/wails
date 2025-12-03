//go:build darwin

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

// applyUpdate applies the downloaded update on macOS
func applyUpdate(ctx context.Context, downloadPath string, info *UpdateInfo) error {
	// Determine if this is a patch or full update
	isPatch := info.PatchURL != "" && strings.HasSuffix(downloadPath, ".patch")

	if isPatch {
		return applyPatch(ctx, downloadPath, info)
	}

	return applyFullUpdate(ctx, downloadPath, info)
}

// applyFullUpdate extracts and replaces the application bundle
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

	// Find the .app bundle (go up from Contents/MacOS/binary)
	appBundlePath := findAppBundle(execPath)
	if appBundlePath == "" {
		return fmt.Errorf("could not find .app bundle")
	}

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

	// Find the extracted .app bundle
	extractedApp, err := findExtractedApp(tempDir)
	if err != nil {
		return fmt.Errorf("failed to find extracted app: %w", err)
	}

	// Create backup of current app
	backupPath := appBundlePath + ".backup"
	if err := os.Rename(appBundlePath, backupPath); err != nil {
		return fmt.Errorf("failed to backup current app: %w", err)
	}

	// Move new app into place
	if err := os.Rename(extractedApp, appBundlePath); err != nil {
		// Restore backup on failure
		os.Rename(backupPath, appBundlePath)
		return fmt.Errorf("failed to install update: %w", err)
	}

	// Remove backup
	os.RemoveAll(backupPath)

	// Clear quarantine attribute
	exec.Command("xattr", "-cr", appBundlePath).Run()

	// Launch the new version and exit
	return restartApp(appBundlePath)
}

// applyPatch applies a bsdiff patch to the current binary
func applyPatch(ctx context.Context, patchPath string, info *UpdateInfo) error {
	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	// Find the .app bundle
	appBundlePath := findAppBundle(execPath)
	if appBundlePath == "" {
		return fmt.Errorf("could not find .app bundle")
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

	// Make executable
	if err := os.Chmod(execPath, 0755); err != nil {
		os.Rename(backupPath, execPath)
		return fmt.Errorf("failed to set executable permissions: %w", err)
	}

	// Remove backup
	os.Remove(backupPath)

	// Re-sign the app (ad-hoc)
	exec.Command("codesign", "--force", "--deep", "--sign", "-", appBundlePath).Run()

	// Restart
	return restartApp(appBundlePath)
}

// findAppBundle finds the .app bundle path from an executable path
func findAppBundle(execPath string) string {
	// Walk up the path looking for .app
	current := execPath
	for current != "/" {
		if strings.HasSuffix(current, ".app") {
			return current
		}
		current = filepath.Dir(current)
	}
	return ""
}

// findExtractedApp finds the .app bundle in the extraction directory
func findExtractedApp(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".app") {
			return filepath.Join(dir, entry.Name()), nil
		}
	}

	return "", fmt.Errorf("no .app bundle found in extracted files")
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
	// Use 'open' to launch the app
	cmd := exec.Command("open", appPath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to restart app: %w", err)
	}

	// Exit the current process
	os.Exit(0)
	return nil
}
