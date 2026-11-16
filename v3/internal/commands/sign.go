package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/bundlearchive"
	"github.com/wailsapp/wails/v3/internal/defaults"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/keychain"
	"github.com/wailsapp/wails/v3/internal/packagesign"
)

// resolveSigningDefaults fills in any signing options not supplied via flags from
// the global config (~/.config/wails/defaults.yaml). This makes the signing
// configuration written by `wails3 setup` take effect at sign time for every
// platform — flags (and therefore project Taskfile vars) still win when present.
func resolveSigningDefaults(options *flags.Sign) {
	explicitAppleIdentity := options.Identity != ""
	cfg, err := defaults.Load()
	if err != nil {
		return
	}

	// macOS
	if options.Identity == "" && cfg.Signing.Darwin.Identity != "" {
		options.Identity = cfg.Signing.Darwin.Identity
	}
	if options.KeychainProfile == "" && !usesNotaryAPI(options) && cfg.Signing.Darwin.KeychainProfile != "" {
		options.KeychainProfile = cfg.Signing.Darwin.KeychainProfile
	}
	if options.Entitlements == "" && cfg.Signing.Darwin.Entitlements != "" {
		options.Entitlements = cfg.Signing.Darwin.Entitlements
	}

	// Windows defaults apply only to non-Apple inputs.
	appleInput := appleSigningInput(options.Input)
	if appleInput {
		if options.Certificate == "" && !explicitAppleIdentity {
			options.Certificate = cfg.Signing.Darwin.P12Path
		}
		if options.KeychainProfile == "" {
			if options.NotaryKey == "" {
				options.NotaryKey = cfg.Signing.Darwin.APIKeyPath
			}
			if options.NotaryKeyID == "" {
				options.NotaryKeyID = cfg.Signing.Darwin.APIKeyID
			}
			if options.NotaryIssuerID == "" {
				options.NotaryIssuerID = cfg.Signing.Darwin.APIIssuerID
			}
		}
	}
	if !appleInput && options.Certificate == "" && options.Thumbprint == "" && cfg.Signing.Windows.CertificatePath != "" {
		options.Certificate = cfg.Signing.Windows.CertificatePath
	}
	if !appleInput && options.Thumbprint == "" && cfg.Signing.Windows.Thumbprint != "" {
		options.Thumbprint = cfg.Signing.Windows.Thumbprint
	}
	if !appleInput && options.Timestamp == "" && cfg.Signing.Windows.TimestampServer != "" {
		options.Timestamp = cfg.Signing.Windows.TimestampServer
	}

	// Linux (PGP). The signer needs an exported key *file*, so use the stored
	// key path; the bare key ID is not usable on its own.
	if options.PGPKey == "" && cfg.Signing.Linux.GPGKeyPath != "" {
		options.PGPKey = cfg.Signing.Linux.GPGKeyPath
	}
	if options.Role == "" && cfg.Signing.Linux.SignRole != "" {
		options.Role = cfg.Signing.Linux.SignRole
	}
}

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

	// Backfill any unset options from the global signing config so the values
	// configured via `wails3 setup` are honored across all platforms.
	resolveSigningDefaults(options)

	// Determine what type of signing to do based on file extension and flags
	ext := strings.ToLower(filepath.Ext(options.Input))

	// macOS app bundle (directory)
	if info.IsDir() && appleBundleInput(options.Input) {
		return signMacOSApp(options)
	}
	if ext == ".dmg" {
		return signMacOSDiskImage(options)
	}

	// macOS binary or Windows executable
	if ext == ".exe" || ext == ".dll" || ext == ".msi" || ext == ".msix" || ext == ".appx" {
		return signWindows(options)
	}

	// Linux packages
	if ext == ".deb" {
		return signDEB(options)
	}
	if ext == ".rpm" {
		return signRPM(options)
	}
	if strings.HasSuffix(strings.ToLower(options.Input), ".pkg.tar.zst") {
		return signArchPackage(context.Background(), options.Input, options.PGPKey, options.PGPPassword)
	}

	// macOS binary (no extension typically)
	if appleSigningInput(options.Input) || runtime.GOOS == "darwin" && options.Identity != "" {
		return signMacOSBinary(options)
	}

	return fmt.Errorf("unsupported file type: %s", ext)
}

func signMacOSDiskImage(options *flags.Sign) error {
	resolveSigningDefaults(options)
	if portableAppleSigning(options) {
		return signPortableApple(context.Background(), options)
	}
	if options.Identity == "" {
		return fmt.Errorf("--identity is required for macOS disk image signing (set via `wails3 setup` or --identity flag)")
	}
	args := []string{"--force", "--sign", options.Identity, options.Input}
	cmd := exec.Command("codesign", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("codesign disk image failed: %w", err)
	}
	pterm.Success.Printfln("Signed: %s", options.Input)
	if !options.Notarize {
		return nil
	}
	return notarizeMacOSDiskImage(context.Background(), options)
}

func notarizeMacOSDiskImage(ctx context.Context, options *flags.Sign) error {
	if usesNotaryAPI(options) {
		return notarizeMacOSWithAPI(ctx, options)
	}
	if options.KeychainProfile == "" {
		return fmt.Errorf("--keychain-profile is required for notarization (set via `wails3 setup` or --keychain-profile flag)")
	}
	if err := submitNotarizationProfile(ctx, options.Input, options.KeychainProfile); err != nil {
		return err
	}
	if err := stapleNotarization(ctx, options.Input); err != nil {
		return fmt.Errorf("disk image stapling failed: %w", err)
	}
	pterm.Success.Println("Notarization complete and ticket stapled")
	return nil
}

func signMacOSApp(options *flags.Sign) error {
	resolveSigningDefaults(options)
	if portableAppleSigning(options) {
		return signPortableApple(context.Background(), options)
	}
	if options.Identity == "" {
		return fmt.Errorf("--identity is required for macOS signing (set via `wails3 setup` or --identity flag)")
	}

	if options.Verbose {
		pterm.Info.Printfln("Signing macOS app bundle: %s", options.Input)
	}

	// Finder and resource-fork metadata can appear on bundles in synced folders
	// after a previous package run. codesign rejects those attributes.
	if output, err := exec.Command("xattr", "-cr", options.Input).CombinedOutput(); err != nil {
		return fmt.Errorf("could not clear bundle metadata before signing: %w: %s", err, strings.TrimSpace(string(output)))
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
		return fmt.Errorf("codesign failed (if Finder metadata is reattached, move the bundle to an unsynced directory): %w", err)
	}
	if !macOSOutputIsLocal(filepath.Dir(options.Input)) {
		// An arbitrary input may be in a synced folder. Allow asynchronous Finder
		// metadata to settle before accepting the signature. Generated local
		// output carries a path-specific marker and avoids this delay.
		time.Sleep(2500 * time.Millisecond)
	}
	if output, err := exec.Command("codesign", "--verify", "--deep", "--strict", options.Input).CombinedOutput(); err != nil {
		return fmt.Errorf("signed bundle failed strict verification (synced folders may reattach Finder metadata; try an unsynced build directory): %w: %s", err, strings.TrimSpace(string(output)))
	}

	pterm.Success.Printfln("Signed: %s", options.Input)

	// Notarize if requested
	if options.Notarize {
		return notarizeMacOSApp(context.Background(), options)
	}

	return nil
}

func signMacOSBinary(options *flags.Sign) error {
	resolveSigningDefaults(options)
	if portableAppleSigning(options) {
		return signPortableApple(context.Background(), options)
	}
	if options.Identity == "" {
		return fmt.Errorf("--identity is required for macOS signing (set via `wails3 setup` or --identity flag)")
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

func notarizeMacOSApp(ctx context.Context, options *flags.Sign) error {
	if usesNotaryAPI(options) {
		return notarizeMacOSWithAPI(ctx, options)
	}
	if options.KeychainProfile == "" {
		return fmt.Errorf("--keychain-profile is required for notarization (set via `wails3 setup` or --keychain-profile flag)")
	}

	if options.Verbose {
		pterm.Info.Println("Submitting for notarization...")
	}

	archiveDir, err := os.MkdirTemp("", "wails-notarization-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(archiveDir)
	zipPath := filepath.Join(archiveDir, "application.zip")
	if err := bundlearchive.Create(ctx, zipPath, options.Input); err != nil {
		return fmt.Errorf("failed to create zip for notarization: %w", err)
	}

	if err := submitNotarizationProfile(ctx, zipPath, options.KeychainProfile); err != nil {
		return err
	}

	if err := stapleNotarization(ctx, options.Input); err != nil {
		return fmt.Errorf("stapling failed: %w", err)
	}

	pterm.Success.Println("Notarization complete and ticket stapled")
	return nil
}

func signWindows(options *flags.Sign) error {
	password := options.Password
	if password == "" && options.Certificate != "" {
		password, _ = keychain.Get(keychain.KeyWindowsCertPassword)
	}
	return signWindowsWithPassword(context.Background(), options, password)
}

func signDEB(options *flags.Sign) error {
	if key, isFile := packagesign.KeyFile("", options.PGPKey); isFile {
		return signExportedPackage(options, "deb", key)
	}
	if options.PGPKey == "" {
		return fmt.Errorf("no PGP signing key found for DEB signing — configure linux.signing in wails.yaml, pass --pgp-key, or run `wails3 setup` / `wails3 setup signing`")
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
	if key, isFile := packagesign.KeyFile("", options.PGPKey); isFile {
		return signExportedPackage(options, "rpm", key)
	}
	if options.PGPKey == "" {
		return fmt.Errorf("no PGP signing key found for RPM signing — configure linux.signing in wails.yaml, pass --pgp-key, or run `wails3 setup` / `wails3 setup signing`")
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
