//go:build windows

package updater

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// applyUpdate applies the downloaded update on Windows
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

	// Get application directory
	appDir := filepath.Dir(execPath)
	appName := filepath.Base(execPath)

	// Create temp directory for extraction
	tempDir, err := os.MkdirTemp("", "wails-update-extract-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}

	// Extract the archive
	if err := extractZip(archivePath, tempDir); err != nil {
		os.RemoveAll(tempDir)
		return fmt.Errorf("failed to extract update: %w", err)
	}

	// Find the new executable
	newExec := filepath.Join(tempDir, appName)
	if _, err := os.Stat(newExec); os.IsNotExist(err) {
		// Try without .exe if not found
		newExec = filepath.Join(tempDir, strings.TrimSuffix(appName, ".exe"))
		if _, err := os.Stat(newExec); os.IsNotExist(err) {
			os.RemoveAll(tempDir)
			return fmt.Errorf("could not find executable in update archive")
		}
	}

	// Create update script that will run after we exit
	scriptPath := filepath.Join(os.TempDir(), "wails-update.bat")
	script := fmt.Sprintf(`@echo off
timeout /t 2 /nobreak > nul
move /y "%s" "%s.old" > nul 2>&1
move /y "%s" "%s" > nul 2>&1
del "%s.old" > nul 2>&1
rmdir /s /q "%s" > nul 2>&1
start "" "%s"
del "%%~f0"
`, execPath, execPath, newExec, execPath, execPath, tempDir, execPath)

	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		os.RemoveAll(tempDir)
		return fmt.Errorf("failed to create update script: %w", err)
	}

	// Run the script hidden
	cmd := exec.Command("cmd", "/C", scriptPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Start(); err != nil {
		os.RemoveAll(tempDir)
		os.Remove(scriptPath)
		return fmt.Errorf("failed to start update script: %w", err)
	}

	// Exit the current process
	os.Exit(0)
	return nil
}

// applyPatch applies a bsdiff patch to the current binary
func applyPatch(ctx context.Context, patchPath string, info *UpdateInfo) error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Create temp file for patched binary
	tempFile, err := os.CreateTemp("", "wails-patched-*.exe")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()

	// Look for bspatch in the app directory or PATH
	bspatchPath := findBspatch()
	if bspatchPath == "" {
		os.Remove(tempPath)
		return fmt.Errorf("bspatch not found")
	}

	// Apply the patch
	cmd := exec.CommandContext(ctx, bspatchPath, execPath, tempPath, patchPath)
	if err := cmd.Run(); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to apply patch: %w", err)
	}

	// Create update script
	scriptPath := filepath.Join(os.TempDir(), "wails-update.bat")
	script := fmt.Sprintf(`@echo off
timeout /t 2 /nobreak > nul
move /y "%s" "%s.old" > nul 2>&1
move /y "%s" "%s" > nul 2>&1
del "%s.old" > nul 2>&1
start "" "%s"
del "%%~f0"
`, execPath, execPath, tempPath, execPath, execPath, execPath)

	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to create update script: %w", err)
	}

	// Run the script hidden
	cmd = exec.Command("cmd", "/C", scriptPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Start(); err != nil {
		os.Remove(tempPath)
		os.Remove(scriptPath)
		return fmt.Errorf("failed to start update script: %w", err)
	}

	os.Exit(0)
	return nil
}

// findBspatch looks for the bspatch executable
func findBspatch() string {
	// Check in app directory first
	execPath, _ := os.Executable()
	appDir := filepath.Dir(execPath)
	bspatchPath := filepath.Join(appDir, "bspatch.exe")
	if _, err := os.Stat(bspatchPath); err == nil {
		return bspatchPath
	}

	// Check in PATH
	path, err := exec.LookPath("bspatch.exe")
	if err == nil {
		return path
	}

	return ""
}

// extractZip extracts a .zip archive
func extractZip(archivePath, destDir string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		targetPath := filepath.Join(destDir, file.Name)

		// Security: prevent path traversal
		if !strings.HasPrefix(targetPath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, file.Mode()); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		inFile, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, inFile)
		inFile.Close()
		outFile.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
