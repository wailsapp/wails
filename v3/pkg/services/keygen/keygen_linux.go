//go:build linux

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

// linuxKeygen implements platform-specific functionality for Linux
type linuxKeygen struct{}

// NewPlatformKeygen creates a new platform-specific keygen implementation
func NewPlatformKeygen() platformKeygen {
	return &linuxKeygen{}
}

// GetMachineFingerprint generates a unique machine fingerprint for Linux
func (l *linuxKeygen) GetMachineFingerprint() (string, error) {
	// Try multiple sources for machine identification
	var identifiers []string

	// 1. Try /etc/machine-id (systemd)
	if machineID, err := os.ReadFile("/etc/machine-id"); err == nil {
		id := strings.TrimSpace(string(machineID))
		if id != "" {
			identifiers = append(identifiers, id)
		}
	}

	// 2. Try /var/lib/dbus/machine-id (D-Bus)
	if machineID, err := os.ReadFile("/var/lib/dbus/machine-id"); err == nil {
		id := strings.TrimSpace(string(machineID))
		if id != "" && (len(identifiers) == 0 || identifiers[0] != id) {
			identifiers = append(identifiers, id)
		}
	}

	// 3. Try DMI product UUID
	if dmiUUID, err := os.ReadFile("/sys/class/dmi/id/product_uuid"); err == nil {
		uuid := strings.TrimSpace(string(dmiUUID))
		if uuid != "" {
			identifiers = append(identifiers, uuid)
		}
	}

	// 4. Try motherboard serial
	if boardSerial, err := os.ReadFile("/sys/class/dmi/id/board_serial"); err == nil {
		serial := strings.TrimSpace(string(boardSerial))
		if serial != "" && serial != "None" && serial != "To be filled by O.E.M." {
			identifiers = append(identifiers, serial)
		}
	}

	// 5. Fallback to MAC addresses
	if len(identifiers) == 0 {
		macs, err := getMACAddresses()
		if err == nil && len(macs) > 0 {
			identifiers = append(identifiers, macs...)
		}
	}

	if len(identifiers) == 0 {
		return "", fmt.Errorf("unable to generate machine fingerprint")
	}

	// Combine all identifiers
	fingerprint := strings.Join(identifiers, "-")
	return hashFingerprint(fingerprint), nil
}

// InstallUpdatePlatform installs an update on Linux
func (l *linuxKeygen) InstallUpdatePlatform(updatePath string) error {
	// Verify the update file exists
	if _, err := os.Stat(updatePath); err != nil {
		return fmt.Errorf("update file not found: %w", err)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(updatePath))
	switch ext {
	case ".deb":
		return l.installDEB(updatePath)
	case ".rpm":
		return l.installRPM(updatePath)
	case ".appimage":
		return l.installAppImage(updatePath)
	case ".tar.gz", ".tgz":
		return l.installTarGz(updatePath)
	case ".sh":
		return l.installScript(updatePath)
	default:
		// Try to make it executable and run
		return l.installBinary(updatePath)
	}
}

// GetCacheDir returns the cache directory for Linux
func (l *linuxKeygen) GetCacheDir() string {
	// Follow XDG Base Directory specification
	xdgCache := os.Getenv("XDG_CACHE_HOME")
	if xdgCache != "" {
		return filepath.Join(xdgCache, getAppName(), "keygen")
	}

	// Fallback to ~/.cache
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to temp directory
		return filepath.Join(os.TempDir(), "keygen-cache")
	}

	return filepath.Join(home, ".cache", getAppName(), "keygen")
}

// Helper methods

// installDEB handles DEB package installation
func (l *linuxKeygen) installDEB(debPath string) error {
	// Check if we have dpkg
	if _, err := exec.LookPath("dpkg"); err != nil {
		return fmt.Errorf("dpkg not found, cannot install .deb package")
	}

	// Try to install with apt if available (handles dependencies better)
	if _, err := exec.LookPath("apt"); err == nil {
		cmd := exec.Command("pkexec", "apt", "install", "-y", debPath)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	// Fallback to dpkg
	cmd := exec.Command("pkexec", "dpkg", "-i", debPath)
	return cmd.Run()
}

// installRPM handles RPM package installation
func (l *linuxKeygen) installRPM(rpmPath string) error {
	// Check for dnf first (newer)
	if _, err := exec.LookPath("dnf"); err == nil {
		cmd := exec.Command("pkexec", "dnf", "install", "-y", rpmPath)
		return cmd.Run()
	}

	// Try yum
	if _, err := exec.LookPath("yum"); err == nil {
		cmd := exec.Command("pkexec", "yum", "install", "-y", rpmPath)
		return cmd.Run()
	}

	// Try zypper (openSUSE)
	if _, err := exec.LookPath("zypper"); err == nil {
		cmd := exec.Command("pkexec", "zypper", "install", "-y", rpmPath)
		return cmd.Run()
	}

	// Fallback to rpm
	if _, err := exec.LookPath("rpm"); err == nil {
		cmd := exec.Command("pkexec", "rpm", "-U", rpmPath)
		return cmd.Run()
	}

	return fmt.Errorf("no RPM package manager found")
}

// installAppImage handles AppImage installation
func (l *linuxKeygen) installAppImage(appImagePath string) error {
	// Make it executable
	if err := os.Chmod(appImagePath, 0755); err != nil {
		return fmt.Errorf("failed to make AppImage executable: %w", err)
	}

	// Get current executable location
	currentExe, err := os.Executable()
	if err != nil {
		// Just run the AppImage
		cmd := exec.Command(appImagePath)
		return cmd.Start()
	}

	// If current app is in a standard location, try to replace it
	if strings.Contains(currentExe, "/usr/") || strings.Contains(currentExe, "/opt/") {
		// Need elevated privileges
		destPath := currentExe + ".new"
		copyCmd := exec.Command("pkexec", "cp", appImagePath, destPath)
		if err := copyCmd.Run(); err != nil {
			return fmt.Errorf("failed to copy AppImage: %w", err)
		}

		// Replace old with new
		moveCmd := exec.Command("pkexec", "mv", destPath, currentExe)
		if err := moveCmd.Run(); err != nil {
			exec.Command("pkexec", "rm", destPath).Run()
			return fmt.Errorf("failed to replace executable: %w", err)
		}

		// Restart application
		cmd := exec.Command(currentExe)
		return cmd.Start()
	}

	// For user-installed apps, just replace directly
	destPath := currentExe + ".new"
	if err := copyFile(appImagePath, destPath); err != nil {
		return fmt.Errorf("failed to copy AppImage: %w", err)
	}

	// Create update script
	scriptPath := filepath.Join(os.TempDir(), "wails-update.sh")
	scriptContent := fmt.Sprintf(`#!/bin/sh
sleep 2
mv "%s" "%s.old" 2>/dev/null
mv "%s" "%s"
chmod +x "%s"
exec "%s" "$@"
rm -f "$0"
`, currentExe, currentExe, destPath, currentExe, currentExe, currentExe)

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		os.Remove(destPath)
		return fmt.Errorf("failed to create update script: %w", err)
	}

	// Run update script
	cmd := exec.Command("/bin/sh", scriptPath)
	if err := cmd.Start(); err != nil {
		os.Remove(scriptPath)
		os.Remove(destPath)
		return fmt.Errorf("failed to start update process: %w", err)
	}

	// Exit current process
	os.Exit(0)
	return nil
}

// installTarGz handles tar.gz archive installation
func (l *linuxKeygen) installTarGz(tarPath string) error {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "keygen-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Extract tar.gz
	cmd := exec.Command("tar", "-xzf", tarPath, "-C", tempDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract tar.gz: %w", err)
	}

	// Look for executable or install script
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read extracted files: %w", err)
	}

	// First, look for install script
	for _, entry := range entries {
		name := strings.ToLower(entry.Name())
		if name == "install.sh" || name == "setup.sh" || name == "install" {
			scriptPath := filepath.Join(tempDir, entry.Name())
			return l.installScript(scriptPath)
		}
	}

	// Look for executable
	appName := strings.ToLower(getAppName())
	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := filepath.Join(tempDir, entry.Name())
			info, err := os.Stat(filePath)
			if err == nil {
				// Check if it's executable or matches app name
				if info.Mode()&0111 != 0 || strings.Contains(strings.ToLower(entry.Name()), appName) {
					return l.installBinary(filePath)
				}
			}
		}
	}

	return fmt.Errorf("no installable content found in archive")
}

// installScript handles shell script installation
func (l *linuxKeygen) installScript(scriptPath string) error {
	// Make it executable
	if err := os.Chmod(scriptPath, 0755); err != nil {
		return fmt.Errorf("failed to make script executable: %w", err)
	}

	// Run the script
	cmd := exec.Command("/bin/sh", scriptPath)
	return cmd.Start()
}

// installBinary handles direct binary installation
func (l *linuxKeygen) installBinary(binaryPath string) error {
	// Make it executable
	if err := os.Chmod(binaryPath, 0755); err != nil {
		return fmt.Errorf("failed to make binary executable: %w", err)
	}

	// Get current executable
	currentExe, err := os.Executable()
	if err != nil {
		// Just run the new binary
		cmd := exec.Command(binaryPath)
		return cmd.Start()
	}

	// Copy to temporary location
	destPath := currentExe + ".new"
	if err := copyFile(binaryPath, destPath); err != nil {
		return fmt.Errorf("failed to copy binary: %w", err)
	}

	// Make new binary executable
	if err := os.Chmod(destPath, 0755); err != nil {
		os.Remove(destPath)
		return fmt.Errorf("failed to make new binary executable: %w", err)
	}

	// Create update script
	scriptPath := filepath.Join(os.TempDir(), "wails-update.sh")
	scriptContent := fmt.Sprintf(`#!/bin/sh
sleep 2
mv "%s" "%s.old" 2>/dev/null
mv "%s" "%s"
exec "%s" "$@"
rm -f "$0"
`, currentExe, currentExe, destPath, currentExe, currentExe)

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		os.Remove(destPath)
		return fmt.Errorf("failed to create update script: %w", err)
	}

	// Run update script
	cmd := exec.Command("/bin/sh", scriptPath)
	if err := cmd.Start(); err != nil {
		os.Remove(scriptPath)
		os.Remove(destPath)
		return fmt.Errorf("failed to start update process: %w", err)
	}

	// Exit current process
	os.Exit(0)
	return nil
}

// getMACAddresses returns MAC addresses of network interfaces
func getMACAddresses() ([]string, error) {
	var macs []string

	// Read network interfaces
	interfaces, err := os.ReadDir("/sys/class/net")
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		// Skip loopback
		if iface.Name() == "lo" {
			continue
		}

		// Read MAC address
		macPath := filepath.Join("/sys/class/net", iface.Name(), "address")
		macBytes, err := os.ReadFile(macPath)
		if err != nil {
			continue
		}

		mac := strings.TrimSpace(string(macBytes))
		if mac != "" && mac != "00:00:00:00:00:00" {
			macs = append(macs, mac)
		}
	}

	return macs, nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0644)
}

// getAppName tries to get the application name
func getAppName() string {
	// Try to get from executable name
	if exe, err := os.Executable(); err == nil {
		name := filepath.Base(exe)
		if name != "" {
			return name
		}
	}

	// Try from process name
	if len(os.Args) > 0 {
		return filepath.Base(os.Args[0])
	}

	return "wailsapp"
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
