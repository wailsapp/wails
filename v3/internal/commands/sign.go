package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/anchore/quill/quill"
	"github.com/anchore/quill/quill/pki/load"
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
		return signMacOS(options)
	}

	// Windows executable
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

	// macOS binary (no extension typically) - check for macOS signing options
	if options.Identity != "" || options.P12Certificate != "" {
		return signMacOS(options)
	}

	return fmt.Errorf("unsupported file type or missing signing options: %s", ext)
}

// signMacOS routes to the appropriate macOS signing method
func signMacOS(options *flags.Sign) error {
	// Determine signing method:
	// 1. If --p12 is set, use Quill library (cross-platform)
	// 2. If on macOS with --identity, use native codesign
	// 3. If not on macOS and no --p12, error

	if options.P12Certificate != "" {
		// Cross-platform signing with Quill library
		return signMacOSWithQuill(options)
	}

	if runtime.GOOS == "darwin" && options.Identity != "" {
		// Native macOS signing
		info, err := os.Stat(options.Input)
		if err != nil {
			return err
		}
		if info.IsDir() && strings.HasSuffix(options.Input, ".app") {
			return signMacOSAppNative(options)
		}
		return signMacOSBinaryNative(options)
	}

	if runtime.GOOS != "darwin" {
		return fmt.Errorf("macOS signing on non-macOS requires --p12 certificate for cross-platform signing")
	}

	return fmt.Errorf("--identity or --p12 is required for macOS signing")
}

// =============================================================================
// Native macOS signing (requires macOS with codesign)
// =============================================================================

func signMacOSAppNative(options *flags.Sign) error {
	if options.Identity == "" {
		return fmt.Errorf("--identity is required for native macOS signing")
	}

	if options.Verbose {
		pterm.Info.Printfln("Signing macOS app bundle (native): %s", options.Input)
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
		return notarizeMacOSNative(options)
	}

	return nil
}

func signMacOSBinaryNative(options *flags.Sign) error {
	if options.Identity == "" {
		return fmt.Errorf("--identity is required for native macOS signing")
	}

	if options.Verbose {
		pterm.Info.Printfln("Signing macOS binary (native): %s", options.Input)
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

	if options.Notarize {
		return notarizeMacOSNative(options)
	}

	return nil
}

func notarizeMacOSNative(options *flags.Sign) error {
	if options.KeychainProfile == "" {
		return fmt.Errorf("--keychain-profile is required for native notarization")
	}

	if options.Verbose {
		pterm.Info.Println("Submitting for notarization (native)...")
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

// =============================================================================
// Cross-platform macOS signing via Quill library
// =============================================================================

// signMacOSWithQuill signs a macOS binary using the Quill library (works on any platform)
func signMacOSWithQuill(options *flags.Sign) error {
	if options.P12Certificate == "" {
		return fmt.Errorf("--p12 is required for cross-platform macOS signing")
	}

	// Check if P12 file exists
	if _, err := os.Stat(options.P12Certificate); err != nil {
		return fmt.Errorf("P12 certificate not found: %w", err)
	}

	// Get P12 password from keychain
	password, err := keychain.Get(keychain.KeyMacOSP12Password)
	if err != nil {
		return fmt.Errorf("P12 password not found: %w\nStore it with: wails3 setup signing", err)
	}

	if options.Verbose {
		pterm.Info.Printfln("Signing macOS binary with Quill: %s", options.Input)
	}

	// Load P12 certificate contents - password stays in memory, never touches disk
	p12Contents, err := load.P12(options.P12Certificate, password)
	if err != nil {
		return fmt.Errorf("failed to load P12 certificate: %w", err)
	}

	// Create signing config from P12 contents
	cfg, err := quill.NewSigningConfigFromP12(options.Input, *p12Contents, false)
	if err != nil {
		return fmt.Errorf("failed to create signing config: %w", err)
	}

	// Add entitlements if specified
	if options.Entitlements != "" {
		cfg = cfg.WithEntitlements(options.Entitlements)
	}

	// Sign the binary
	if err := quill.Sign(*cfg); err != nil {
		return fmt.Errorf("signing failed: %w", err)
	}

	pterm.Success.Printfln("Signed: %s", options.Input)

	// Notarize if requested
	if options.Notarize {
		return notarizeMacOSWithQuill(options)
	}

	return nil
}

// notarizeMacOSWithQuill notarizes a macOS binary using the Quill library
func notarizeMacOSWithQuill(options *flags.Sign) error {
	// For notarization, we need the Apple API key credentials
	if options.NotaryKey == "" {
		return fmt.Errorf("--notary-key is required for notarization")
	}

	// Read the API key file
	privateKey, err := os.ReadFile(options.NotaryKey)
	if err != nil {
		return fmt.Errorf("failed to read notary key: %w", err)
	}

	// Get notary credentials - first check flags, then keychain
	notaryKeyID := options.NotaryKeyID
	if notaryKeyID == "" {
		notaryKeyID, err = keychain.Get(keychain.KeyNotaryKeyID)
		if err != nil {
			return fmt.Errorf("notary key ID not found: %w\nProvide --notary-key-id or run: wails3 setup signing", err)
		}
	}

	notaryIssuer := options.NotaryIssuer
	if notaryIssuer == "" {
		notaryIssuer, err = keychain.Get(keychain.KeyNotaryIssuer)
		if err != nil {
			return fmt.Errorf("notary issuer not found: %w\nProvide --notary-issuer or run: wails3 setup signing", err)
		}
	}

	if options.Verbose {
		pterm.Info.Println("Submitting for notarization...")
	}

	// Create notarization config
	cfg := quill.NewNotarizeConfig(notaryIssuer, notaryKeyID, string(privateKey))

	// Submit for notarization and wait for result
	status, err := quill.Notarize(options.Input, *cfg)
	if err != nil {
		return fmt.Errorf("notarization failed: %w", err)
	}

	if options.Verbose {
		pterm.Info.Printfln("Notarization status: %s", status)
	}

	pterm.Success.Println("Notarization complete")
	return nil
}

// =============================================================================
// Windows signing
// =============================================================================

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

// =============================================================================
// Linux signing
// =============================================================================

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
