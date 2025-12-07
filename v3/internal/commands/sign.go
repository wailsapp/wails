package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/keychain"
)

// Sign signs a binary or package
func Sign(options *flags.Sign) error {
	if options.Input == "" {
		return fmt.Errorf("--input is required")
	}

	// Check input file exists
	info, err := os.Stat(options.Input)
	if err != nil {
		return fmt.Errorf("input file not found: %w", err)
	}

	// Determine what type of signing to do based on file extension and flags
	ext := strings.ToLower(filepath.Ext(options.Input))

	// macOS app bundle (directory)
	if info.IsDir() && strings.HasSuffix(options.Input, ".app") {
		return signMacOSApp(options)
	}

	// macOS binary or Windows executable
	if ext == ".exe" || ext == ".msi" || ext == ".msix" || ext == ".appx" {
		return signWindows(options)
	}

	// Linux packages
	if ext == ".deb" {
		return signDEB(options)
	}
	if ext == ".rpm" {
		return signRPM(options)
	}

	// macOS binary (no extension typically)
	if runtime.GOOS == "darwin" && options.Identity != "" {
		return signMacOSBinary(options)
	}

	return fmt.Errorf("unsupported file type: %s", ext)
}

func signMacOSApp(options *flags.Sign) error {
	if options.Identity == "" {
		return fmt.Errorf("--identity is required for macOS signing")
	}

	if options.Verbose {
		pterm.Info.Printfln("Signing macOS app bundle: %s", options.Input)
	}

	// Build codesign command
	args := []string{
		"--force",
		"--deep",
		"--sign", options.Identity,
	}

	if options.Entitlements != "" {
		args = append(args, "--entitlements", options.Entitlements)
	}

	if options.HardenedRuntime || options.Notarize {
		args = append(args, "--options", "runtime")
	}

	args = append(args, options.Input)

	cmd := exec.Command("codesign", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("codesign failed: %w", err)
	}

	pterm.Success.Printfln("Signed: %s", options.Input)

	// Notarize if requested
	if options.Notarize {
		return notarizeMacOSApp(options)
	}

	return nil
}

func signMacOSBinary(options *flags.Sign) error {
	if options.Identity == "" {
		return fmt.Errorf("--identity is required for macOS signing")
	}

	if options.Verbose {
		pterm.Info.Printfln("Signing macOS binary: %s", options.Input)
	}

	args := []string{
		"--force",
		"--sign", options.Identity,
	}

	if options.Entitlements != "" {
		args = append(args, "--entitlements", options.Entitlements)
	}

	if options.HardenedRuntime {
		args = append(args, "--options", "runtime")
	}

	args = append(args, options.Input)

	cmd := exec.Command("codesign", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("codesign failed: %w", err)
	}

	pterm.Success.Printfln("Signed: %s", options.Input)
	return nil
}

func notarizeMacOSApp(options *flags.Sign) error {
	if options.KeychainProfile == "" {
		return fmt.Errorf("--keychain-profile is required for notarization")
	}

	if options.Verbose {
		pterm.Info.Println("Submitting for notarization...")
	}

	// Create a zip for notarization
	zipPath := options.Input + ".zip"
	zipCmd := exec.Command("ditto", "-c", "-k", "--keepParent", options.Input, zipPath)
	if err := zipCmd.Run(); err != nil {
		return fmt.Errorf("failed to create zip for notarization: %w", err)
	}
	defer os.Remove(zipPath)

	// Submit for notarization
	args := []string{
		"notarytool", "submit",
		zipPath,
		"--keychain-profile", options.KeychainProfile,
		"--wait",
	}

	cmd := exec.Command("xcrun", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("notarization failed: %w", err)
	}

	// Staple the ticket
	stapleCmd := exec.Command("xcrun", "stapler", "staple", options.Input)
	stapleCmd.Stdout = os.Stdout
	stapleCmd.Stderr = os.Stderr

	if err := stapleCmd.Run(); err != nil {
		return fmt.Errorf("stapling failed: %w", err)
	}

	pterm.Success.Println("Notarization complete and ticket stapled")
	return nil
}

func signWindows(options *flags.Sign) error {
	// Get password from keychain if not provided
	password := options.Password
	if password == "" && options.Certificate != "" {
		var err error
		password, err = keychain.Get(keychain.KeyWindowsCertPassword)
		if err != nil {
			pterm.Warning.Printfln("Could not get password from keychain: %v", err)
			// Continue without password - might work for some certificates
		}
	}

	if options.Verbose {
		pterm.Info.Printfln("Signing Windows executable: %s", options.Input)
	}

	// Try native signtool first on Windows
	if runtime.GOOS == "windows" {
		err := signWindowsNative(options, password)
		if err == nil {
			return nil
		}
		if options.Verbose {
			pterm.Warning.Printfln("Native signing failed, trying built-in: %v", err)
		}
	}

	// Use built-in signing (works cross-platform)
	return signWindowsBuiltin(options, password)
}

func signWindowsNative(options *flags.Sign, password string) error {
	// Find signtool.exe
	signtool, err := findSigntool()
	if err != nil {
		return err
	}

	args := []string{"sign"}

	if options.Certificate != "" {
		args = append(args, "/f", options.Certificate)
		if password != "" {
			args = append(args, "/p", password)
		}
	} else if options.Thumbprint != "" {
		args = append(args, "/sha1", options.Thumbprint)
	} else {
		return fmt.Errorf("either --certificate or --thumbprint is required")
	}

	// Add timestamp server
	timestamp := options.Timestamp
	if timestamp == "" {
		timestamp = "http://timestamp.digicert.com"
	}
	args = append(args, "/tr", timestamp, "/td", "SHA256", "/fd", "SHA256")

	args = append(args, options.Input)

	cmd := exec.Command(signtool, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("signtool failed: %w", err)
	}

	pterm.Success.Printfln("Signed: %s", options.Input)
	return nil
}

func findSigntool() (string, error) {
	// Check if signtool is in PATH
	path, err := exec.LookPath("signtool.exe")
	if err == nil {
		return path, nil
	}

	// Common Windows SDK locations
	sdkPaths := []string{
		`C:\Program Files (x86)\Windows Kits\10\bin\10.0.22621.0\x64\signtool.exe`,
		`C:\Program Files (x86)\Windows Kits\10\bin\10.0.19041.0\x64\signtool.exe`,
		`C:\Program Files (x86)\Windows Kits\10\bin\x64\signtool.exe`,
	}

	for _, p := range sdkPaths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("signtool.exe not found")
}

func signWindowsBuiltin(options *flags.Sign, password string) error {
	// This would use a Go library for Authenticode signing
	// For now, we'll return an error indicating it needs implementation
	// In a full implementation, you'd use something like:
	// - github.com/AkarinLiu/osslsigncode-go
	// - or implement PE signing directly

	if options.Certificate == "" {
		return fmt.Errorf("--certificate is required for cross-platform signing")
	}

	return fmt.Errorf("built-in Windows signing not yet implemented - please use signtool.exe on Windows, or install osslsigncode")
}

func signDEB(options *flags.Sign) error {
	if options.PGPKey == "" {
		return fmt.Errorf("--pgp-key is required for DEB signing")
	}

	// Get password from keychain if not provided
	password := options.PGPPassword
	if password == "" {
		var err error
		password, err = keychain.Get(keychain.KeyPGPPassword)
		if err != nil {
			// Password might not be required if key is unencrypted
			if options.Verbose {
				pterm.Warning.Printfln("Could not get PGP password from keychain: %v", err)
			}
		}
	}

	if options.Verbose {
		pterm.Info.Printfln("Signing DEB package: %s", options.Input)
	}

	role := options.Role
	if role == "" {
		role = "builder"
	}

	// Use dpkg-sig for signing
	args := []string{
		"-k", options.PGPKey,
		"--sign", role,
	}

	if password != "" {
		// dpkg-sig reads from GPG_TTY or gpg-agent
		// For scripted use, we need to use gpg with passphrase
		args = append(args, "--gpg-options", fmt.Sprintf("--batch --passphrase %s", password))
	}

	args = append(args, options.Input)

	cmd := exec.Command("dpkg-sig", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Fallback: try using gpg directly to sign
		return signDEBWithGPG(options, password, role)
	}

	pterm.Success.Printfln("Signed: %s", options.Input)
	return nil
}

func signDEBWithGPG(options *flags.Sign, password, role string) error {
	// Alternative approach using ar and gpg directly
	// This is more portable but more complex
	return fmt.Errorf("dpkg-sig not found - please install dpkg-sig or use a Linux system")
}

func signRPM(options *flags.Sign) error {
	if options.PGPKey == "" {
		return fmt.Errorf("--pgp-key is required for RPM signing")
	}

	// Get password from keychain if not provided
	password := options.PGPPassword
	if password == "" {
		var err error
		password, err = keychain.Get(keychain.KeyPGPPassword)
		if err != nil {
			if options.Verbose {
				pterm.Warning.Printfln("Could not get PGP password from keychain: %v", err)
			}
		}
	}

	if options.Verbose {
		pterm.Info.Printfln("Signing RPM package: %s", options.Input)
	}

	// RPM signing requires the key to be imported to GPG keyring
	// and uses rpmsign command
	args := []string{
		"--addsign",
		options.Input,
	}

	cmd := exec.Command("rpmsign", args...)

	// Set up passphrase via environment if needed
	if password != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("GPG_PASSPHRASE=%s", password))
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("rpmsign failed: %w", err)
	}

	pterm.Success.Printfln("Signed: %s", options.Input)
	return nil
}
