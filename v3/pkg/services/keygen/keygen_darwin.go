//go:build darwin

package keygen

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// darwinKeygen implements platform-specific functionality for macOS
type darwinKeygen struct{}

// NewPlatformKeygen creates a new platform-specific keygen implementation
func NewPlatformKeygen() platformKeygen {
	return &darwinKeygen{}
}

// GetMachineFingerprint generates a unique machine fingerprint for macOS
func (d *darwinKeygen) GetMachineFingerprint() (string, error) {
	// Get hardware UUID
	cmd := exec.Command("system_profiler", "SPHardwareDataType")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get hardware info: %w", err)
	}

	// Parse hardware UUID
	lines := strings.Split(string(output), "\n")
	var hardwareUUID string
	for _, line := range lines {
		if strings.Contains(line, "Hardware UUID:") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				hardwareUUID = strings.TrimSpace(parts[1])
				break
			}
		}
	}

	if hardwareUUID == "" {
		return "", fmt.Errorf("hardware UUID not found")
	}

	// Get serial number as additional entropy
	cmd = exec.Command("ioreg", "-c", "IOPlatformExpertDevice", "-d", "2")
	output, err = cmd.Output()
	if err != nil {
		// Use hardware UUID only if we can't get serial
		return hashFingerprint(hardwareUUID), nil
	}

	// Parse serial number
	var serialNumber string
	lines = strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "IOPlatformSerialNumber") {
			parts := strings.Split(line, "=")
			if len(parts) >= 2 {
				serialNumber = strings.Trim(strings.TrimSpace(parts[1]), "\"")
				break
			}
		}
	}

	// Combine hardware UUID and serial number
	fingerprint := hardwareUUID
	if serialNumber != "" {
		fingerprint += "-" + serialNumber
	}

	return hashFingerprint(fingerprint), nil
}

// InstallUpdatePlatform installs an update on macOS
func (d *darwinKeygen) InstallUpdatePlatform(updatePath string) error {
	// Verify the update file exists
	if _, err := os.Stat(updatePath); err != nil {
		return fmt.Errorf("update file not found: %w", err)
	}

	// Check file extension
	ext := filepath.Ext(updatePath)
	switch ext {
	case ".dmg":
		return d.installDMG(updatePath)
	case ".pkg":
		return d.installPKG(updatePath)
	case ".app":
		return d.installApp(updatePath)
	case ".zip":
		return d.installZip(updatePath)
	default:
		// Try to make it executable and run
		if err := os.Chmod(updatePath, 0755); err != nil {
			return fmt.Errorf("failed to make update executable: %w", err)
		}

		// Launch the updater
		cmd := exec.Command(updatePath)
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start updater: %w", err)
		}

		return nil
	}
}

// GetCacheDir returns the cache directory for macOS
func (d *darwinKeygen) GetCacheDir() string {
	// Use ~/Library/Caches/[app-name]/keygen
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to temp directory
		return filepath.Join(os.TempDir(), "keygen-cache")
	}

	// Try to get app name from bundle
	appName := getAppName()

	return filepath.Join(home, "Library", "Caches", appName, "keygen")
}

// Helper methods

// installDMG handles DMG installation
func (d *darwinKeygen) installDMG(dmgPath string) error {
	// Mount the DMG
	mountCmd := exec.Command("hdiutil", "attach", dmgPath, "-nobrowse", "-noautoopen")
	output, err := mountCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to mount DMG: %w", err)
	}

	// Parse mount point
	lines := strings.Split(string(output), "\n")
	var mountPoint string
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 3 && strings.HasPrefix(parts[2], "/Volumes/") {
			mountPoint = parts[2]
			break
		}
	}

	if mountPoint == "" {
		return fmt.Errorf("failed to find mount point")
	}

	// Defer unmount
	defer func() {
		exec.Command("hdiutil", "detach", mountPoint).Run()
	}()

	// Find .app in the mounted volume
	entries, err := os.ReadDir(mountPoint)
	if err != nil {
		return fmt.Errorf("failed to read mounted volume: %w", err)
	}

	var appPath string
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".app") {
			appPath = filepath.Join(mountPoint, entry.Name())
			break
		}
	}

	if appPath == "" {
		return fmt.Errorf("no .app found in DMG")
	}

	// Copy to Applications
	destPath := filepath.Join("/Applications", filepath.Base(appPath))

	// Remove existing app if present
	os.RemoveAll(destPath)

	// Copy new app
	copyCmd := exec.Command("cp", "-R", appPath, destPath)
	if err := copyCmd.Run(); err != nil {
		return fmt.Errorf("failed to copy app: %w", err)
	}

	// Open the new app
	openCmd := exec.Command("open", destPath)
	return openCmd.Start()
}

// installPKG handles PKG installation
func (d *darwinKeygen) installPKG(pkgPath string) error {
	// Use installer command (requires admin privileges)
	cmd := exec.Command("open", pkgPath)
	return cmd.Start()
}

// installApp handles direct .app installation
func (d *darwinKeygen) installApp(appPath string) error {
	// Get app name
	appName := filepath.Base(appPath)
	destPath := filepath.Join("/Applications", appName)

	// Remove existing app if present
	os.RemoveAll(destPath)

	// Copy new app
	copyCmd := exec.Command("cp", "-R", appPath, destPath)
	if err := copyCmd.Run(); err != nil {
		return fmt.Errorf("failed to copy app: %w", err)
	}

	// Open the new app
	openCmd := exec.Command("open", destPath)
	return openCmd.Start()
}

// installZip handles ZIP file installation
func (d *darwinKeygen) installZip(zipPath string) error {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "keygen-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Unzip
	unzipCmd := exec.Command("unzip", "-q", zipPath, "-d", tempDir)
	if err := unzipCmd.Run(); err != nil {
		return fmt.Errorf("failed to unzip: %w", err)
	}

	// Find .app in extracted files
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read extracted files: %w", err)
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".app") {
			appPath := filepath.Join(tempDir, entry.Name())
			return d.installApp(appPath)
		}
	}

	// No .app found, look for executable
	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := filepath.Join(tempDir, entry.Name())
			info, err := os.Stat(filePath)
			if err == nil && info.Mode()&0111 != 0 {
				// File is executable
				if err := os.Chmod(filePath, 0755); err != nil {
					return fmt.Errorf("failed to make file executable: %w", err)
				}

				cmd := exec.Command(filePath)
				return cmd.Start()
			}
		}
	}

	return fmt.Errorf("no installable content found in ZIP")
}

// getAppName tries to get the application name
func getAppName() string {
	// Try to get from bundle
	cmd := exec.Command("defaults", "read", "/Applications/"+os.Args[0]+"/Contents/Info.plist", "CFBundleName")
	if output, err := cmd.Output(); err == nil {
		name := strings.TrimSpace(string(output))
		if name != "" {
			return name
		}
	}

	// Try from executable name
	if len(os.Args) > 0 {
		return filepath.Base(os.Args[0])
	}

	return "WailsApp"
}

// hashFingerprint creates a consistent hash from the fingerprint components
func hashFingerprint(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// init sets the platform implementation when the service is created
func init() {
	// This will be called when creating a new service
}
