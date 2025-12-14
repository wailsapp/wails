package commands

import (
	"fmt"
	"io"
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
	// 1. If --use-docker is set, use Docker with Quill
	// 2. If --p12 is set (cross-platform signing), use Quill (native or Docker)
	// 3. If on macOS with --identity, use native codesign
	// 4. If on macOS without --identity but with --p12, use Quill
	// 5. If not on macOS and no --p12, error

	if options.UseDocker {
		return signMacOSWithDocker(options)
	}

	if options.P12Certificate != "" {
		// Cross-platform signing with Quill
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
		return fmt.Errorf("macOS signing on non-macOS requires --p12 certificate for cross-platform signing, or --use-docker")
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
// Cross-platform macOS signing via Quill
// =============================================================================

// signMacOSWithQuill signs a macOS binary using Quill (works on any platform)
func signMacOSWithQuill(options *flags.Sign) error {
	if options.P12Certificate == "" {
		return fmt.Errorf("--p12 is required for cross-platform macOS signing")
	}

	// Check if P12 file exists
	if _, err := os.Stat(options.P12Certificate); err != nil {
		return fmt.Errorf("P12 certificate not found: %w", err)
	}

	// Check if quill is available
	quillPath, err := exec.LookPath("quill")
	if err != nil {
		// Quill not found - try Docker
		if options.Verbose {
			pterm.Info.Println("Quill not found locally, trying Docker...")
		}
		return signMacOSWithDocker(options)
	}

	// Get P12 password from keychain
	password, err := keychain.Get(keychain.KeyMacOSP12Password)
	if err != nil {
		return fmt.Errorf("P12 password not found: %w\nStore it with: wails3 setup signing", err)
	}

	if options.Verbose {
		pterm.Info.Printfln("Signing macOS binary with Quill: %s", options.Input)
	}

	// Write password to temp file for Quill's --password-file flag
	passFile, err := os.CreateTemp("", "wails-sign-pass-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(passFile.Name())

	if err := os.Chmod(passFile.Name(), 0600); err != nil {
		return fmt.Errorf("failed to set temp file permissions: %w", err)
	}
	if _, err := passFile.WriteString(password); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	passFile.Close()

	// Build quill sign command
	args := []string{
		"sign",
		"--p12", options.P12Certificate,
		"--password-file", passFile.Name(),
		options.Input,
	}

	cmd := exec.Command(quillPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("quill sign failed: %w", err)
	}

	pterm.Success.Printfln("Signed: %s", options.Input)

	// Notarize if requested
	if options.Notarize {
		return notarizeMacOSWithQuill(options, passFile.Name())
	}

	return nil
}

// notarizeMacOSWithQuill notarizes a macOS binary using Quill
func notarizeMacOSWithQuill(options *flags.Sign, passFilePath string) error {
	// For notarization, we need the Apple API key credentials
	if options.NotaryKey == "" {
		return fmt.Errorf("--notary-key is required for notarization with Quill")
	}

	// Get notary credentials - first check flags, then keychain
	notaryKeyID := options.NotaryKeyID
	if notaryKeyID == "" {
		var err error
		notaryKeyID, err = keychain.Get(keychain.KeyNotaryKeyID)
		if err != nil {
			return fmt.Errorf("notary key ID not found: %w\nProvide --notary-key-id or run: wails3 setup signing", err)
		}
	}

	notaryIssuer := options.NotaryIssuer
	if notaryIssuer == "" {
		var err error
		notaryIssuer, err = keychain.Get(keychain.KeyNotaryIssuer)
		if err != nil {
			return fmt.Errorf("notary issuer not found: %w\nProvide --notary-issuer or run: wails3 setup signing", err)
		}
	}

	if options.Verbose {
		pterm.Info.Println("Submitting for notarization with Quill...")
	}

	quillPath, err := exec.LookPath("quill")
	if err != nil {
		return fmt.Errorf("quill not found: %w", err)
	}

	args := []string{
		"notarize",
		"--p12", options.P12Certificate,
		"--password-file", passFilePath,
		"--notary-key", options.NotaryKey,
		"--notary-key-id", notaryKeyID,
		"--notary-issuer", notaryIssuer,
		options.Input,
	}

	cmd := exec.Command(quillPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("quill notarize failed: %w", err)
	}

	pterm.Success.Println("Notarization complete")
	return nil
}

// =============================================================================
// Docker-based macOS signing (uses Quill in container)
// =============================================================================

const dockerCrossImage = "wails-cross"

// signMacOSWithDocker signs a macOS binary using Quill in a Docker container
func signMacOSWithDocker(options *flags.Sign) error {
	if options.P12Certificate == "" {
		return fmt.Errorf("--p12 is required for Docker-based macOS signing")
	}

	// Check Docker is available
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker not found: %w\nInstall Docker or Quill for cross-platform macOS signing", err)
	}

	// Check Docker image exists
	checkCmd := exec.Command("docker", "image", "inspect", dockerCrossImage)
	if err := checkCmd.Run(); err != nil {
		return fmt.Errorf("docker image '%s' not found\nBuild it with: wails3 task setup:docker", dockerCrossImage)
	}

	// Check if P12 file exists
	p12Abs, err := filepath.Abs(options.P12Certificate)
	if err != nil {
		return fmt.Errorf("failed to resolve P12 path: %w", err)
	}
	if _, err := os.Stat(p12Abs); err != nil {
		return fmt.Errorf("P12 certificate not found: %w", err)
	}

	// Get input file absolute path
	inputAbs, err := filepath.Abs(options.Input)
	if err != nil {
		return fmt.Errorf("failed to resolve input path: %w", err)
	}

	// Get P12 password from keychain
	password, err := keychain.Get(keychain.KeyMacOSP12Password)
	if err != nil {
		return fmt.Errorf("P12 password not found: %w\nStore it with: wails3 setup signing", err)
	}

	if options.Verbose {
		pterm.Info.Printfln("Signing macOS binary with Docker: %s", options.Input)
	}

	// Build docker command - password passed via stdin
	dockerArgs := []string{
		"run", "--rm", "-i",
		"-v", p12Abs + ":/cert.p12:ro",
		"-v", inputAbs + ":/input",
		"--entrypoint", "/usr/local/bin/sign.sh",
		dockerCrossImage,
		"/input",
		"--p12", "/cert.p12",
	}

	// Add notarization args if requested
	if options.Notarize {
		if options.NotaryKey == "" {
			return fmt.Errorf("--notary-key is required for notarization")
		}

		notaryKeyAbs, err := filepath.Abs(options.NotaryKey)
		if err != nil {
			return fmt.Errorf("failed to resolve notary key path: %w", err)
		}

		// Get notary credentials
		notaryKeyID := options.NotaryKeyID
		if notaryKeyID == "" {
			notaryKeyID, err = keychain.Get(keychain.KeyNotaryKeyID)
			if err != nil {
				return fmt.Errorf("notary key ID not found: %w", err)
			}
		}

		notaryIssuer := options.NotaryIssuer
		if notaryIssuer == "" {
			notaryIssuer, err = keychain.Get(keychain.KeyNotaryIssuer)
			if err != nil {
				return fmt.Errorf("notary issuer not found: %w", err)
			}
		}

		// Update docker args to include notary key mount
		dockerArgs = []string{
			"run", "--rm", "-i",
			"-v", p12Abs + ":/cert.p12:ro",
			"-v", notaryKeyAbs + ":/notary-key.p8:ro",
			"-v", inputAbs + ":/input",
			"--entrypoint", "/usr/local/bin/sign.sh",
			dockerCrossImage,
			"/input",
			"--p12", "/cert.p12",
			"--notarize",
			"--notary-key", "/notary-key.p8",
			"--notary-key-id", notaryKeyID,
			"--notary-issuer", notaryIssuer,
		}
	}

	cmd := exec.Command("docker", dockerArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Pass password via stdin (secure - not in process list or env)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start docker: %w", err)
	}

	// Write password to stdin
	if _, err := io.WriteString(stdin, password); err != nil {
		return fmt.Errorf("failed to write password: %w", err)
	}
	stdin.Close()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("docker signing failed: %w", err)
	}

	pterm.Success.Printfln("Signed: %s", options.Input)
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
