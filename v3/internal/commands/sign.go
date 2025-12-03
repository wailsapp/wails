package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/signing"
	"github.com/wailsapp/wails/v3/internal/term"
	"gopkg.in/yaml.v3"
)

// SignOptions defines the options for the sign command
type SignOptions struct {
	// Input is the path to the binary or app bundle to sign
	Input string `name:"input" description:"Path to the binary or .app bundle to sign"`

	// AppPath is an alias for Input (for backwards compatibility with macOS)
	AppPath string `name:"app" description:"Path to the .app bundle to sign (alias for --input)"`

	// Output is the path for the signed binary (optional, defaults to signing in place)
	Output string `name:"output" description:"Path for the signed binary (defaults to in-place)"`

	// Platform is the target platform (auto-detected from file if not specified)
	Platform string `name:"platform" description:"Target platform: windows, darwin (auto-detected from file)"`

	// Certificate is the path to the .pfx/.p12 certificate file
	Certificate string `name:"certificate" description:"Path to .pfx/.p12 certificate file (for cross-platform signing)"`

	// CertificatePassword is the password for the certificate file
	CertificatePassword string `name:"certificate-password" description:"Password for the certificate file (or use $SIGN_CERT_PASSWORD)"`

	// Identity is the signing identity (macOS native only)
	Identity string `name:"identity" description:"Code signing identity (macOS: use '-' for ad-hoc, or leave empty to auto-detect)"`

	// Thumbprint is the certificate thumbprint (Windows native only)
	Thumbprint string `name:"thumbprint" description:"Certificate thumbprint in Windows cert store (Windows native only)"`

	// Entitlements is the path to the entitlements file (macOS)
	Entitlements string `name:"entitlements" description:"Path to the entitlements plist file (macOS)"`

	// HardenedRuntime enables hardened runtime (macOS, required for notarization)
	HardenedRuntime bool `name:"hardened-runtime" description:"Enable hardened runtime (macOS, required for notarization)" default:"true"`

	// TimestampServer is the URL of the timestamp server
	TimestampServer string `name:"timestamp" description:"Timestamp server URL (default: http://timestamp.digicert.com)"`

	// Description is the application description (Windows Authenticode)
	Description string `name:"description" description:"Application description (Windows Authenticode)"`

	// URL is the application URL (Windows Authenticode)
	URL string `name:"url" description:"Application URL (Windows Authenticode)"`

	// Notarize submits the app for notarization after signing (macOS)
	Notarize bool `name:"notarize" description:"Submit for notarization after signing (macOS)"`

	// KeychainProfile is the name of the keychain profile for notarization (macOS)
	KeychainProfile string `name:"keychain-profile" description:"Keychain profile name for notarization credentials"`

	// AppleID is the Apple ID for notarization (macOS)
	AppleID string `name:"apple-id" description:"Apple ID for notarization (or use $APPLE_ID env var)"`

	// TeamID is the Apple Developer Team ID (macOS)
	TeamID string `name:"team-id" description:"Apple Developer Team ID (or use $APPLE_TEAM_ID env var)"`

	// AppPassword is the app-specific password for notarization (macOS)
	AppPassword string `name:"app-password" description:"App-specific password for notarization (or use $APPLE_APP_PASSWORD env var)"`

	// Config is the path to a signing configuration file
	Config string `name:"config" description:"Path to signing configuration file (YAML)"`

	// PGPKey is the path to PGP private key file (for Linux package signing)
	PGPKey string `name:"pgp-key" description:"Path to PGP private key file (for Linux .deb/.rpm signing)"`

	// PGPKeyPassword is the password for the PGP key
	PGPKeyPassword string `name:"pgp-key-password" description:"Password for PGP key (or use $PGP_KEY_PASSWORD)"`

	// Role is the signing role for DEB packages (builder, origin, maint, archive)
	Role string `name:"role" description:"Signing role for DEB packages (builder, origin, maint, archive)" default:"builder"`

	// Verify runs signature verification after signing
	Verify bool `name:"verify" description:"Verify signature after signing" default:"true"`

	// Verbose enables verbose output
	Verbose bool `name:"verbose" description:"Enable verbose output"`

	// NoColour disables colour output
	NoColour bool `name:"n" description:"Disable colour output"`
}

// Sign performs code signing
func Sign(options *SignOptions) error {
	if options.NoColour {
		term.DisableColor()
	}

	DisableFooter = true
	term.Header("Code Signing")

	// Load config file if specified
	if options.Config != "" {
		if err := loadSigningConfig(options); err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
	}

	// Handle backwards compatibility: --app is alias for --input
	if options.Input == "" && options.AppPath != "" {
		options.Input = options.AppPath
	}

	// Auto-detect input path if not specified
	if options.Input == "" {
		options.Input = findSignableFile()
		if options.Input == "" {
			return fmt.Errorf("no signable file found. Specify path with --input or --app flag")
		}
		pterm.Printf("Found: %s\n", options.Input)
	}

	// Validate input path
	info, err := os.Stat(options.Input)
	if err != nil {
		return fmt.Errorf("cannot access input: %w", err)
	}

	// Detect platform from file
	platform := detectPlatformFromOptions(options, info)
	pterm.Printf("Target platform: %s\n", platform)

	// Show which backend will be used
	signer, err := signing.DefaultRegistry.GetSigner(platform)
	if err != nil {
		return fmt.Errorf("no signer available: %w", err)
	}
	pterm.Printf("Signing backend: %s\n", signer.Backend())

	// Build the sign request
	req := signing.SignRequest{
		InputPath:       options.Input,
		OutputPath:      options.Output,
		Platform:        platform,
		TimestampServer: options.TimestampServer,
		Description:     options.Description,
		URL:             options.URL,
		Entitlements:    options.Entitlements,
		HardenedRuntime: options.HardenedRuntime,
		BundleSign:      info.IsDir() && strings.HasSuffix(options.Input, ".app"),
		Verbose:         options.Verbose,
		Certificate: signing.CertificateConfig{
			PKCS12Path:     options.Certificate,
			PKCS12Password: getEnvOrValue(options.CertificatePassword, "SIGN_CERT_PASSWORD"),
			Identity:       options.Identity,
			Thumbprint:     options.Thumbprint,
			PGPKeyPath:     options.PGPKey,
			PGPKeyPassword: getEnvOrValue(options.PGPKeyPassword, "PGP_KEY_PASSWORD"),
		},
	}

	// Platform-specific setup
	if platform == signing.PlatformMacOS {
		if err := setupMacOSSigning(options, &req); err != nil {
			return err
		}
	}

	// Linux-specific setup
	if platform == signing.PlatformLinux {
		if err := setupLinuxSigning(options, &req); err != nil {
			return err
		}
	}

	pterm.Println()
	pterm.Printf("Signing %s...\n", filepath.Base(options.Input))

	// Perform signing
	ctx := context.Background()
	result, err := signer.Sign(ctx, req)
	if err != nil {
		return fmt.Errorf("signing failed: %w", err)
	}

	pterm.Success.Printf("Signing complete! (backend: %s)\n", result.Backend)

	// Verify if requested
	if options.Verify {
		pterm.Println()
		pterm.Println("Verifying signature...")
		if err := signer.Verify(ctx, result.OutputPath); err != nil {
			pterm.Warning.Printf("Signature verification: %v\n", err)
		} else {
			pterm.Success.Println("Signature verified!")
		}
	}

	// Notarize if requested (macOS only)
	if options.Notarize {
		if platform != signing.PlatformMacOS {
			return fmt.Errorf("notarization is only supported for macOS")
		}
		if runtime.GOOS != "darwin" {
			return fmt.Errorf("notarization requires macOS (use rcodesign for cross-platform notarization)")
		}

		pterm.Println()

		if options.Identity == "-" {
			return fmt.Errorf("cannot notarize ad-hoc signed apps. Use a Developer ID certificate")
		}

		notarizeOpts := signing.NotarizeOptions{
			AppPath:             result.OutputPath,
			KeychainProfile:     options.KeychainProfile,
			AppleID:             getEnvOrValue(options.AppleID, "APPLE_ID"),
			TeamID:              getEnvOrValue(options.TeamID, "APPLE_TEAM_ID"),
			AppSpecificPassword: getEnvOrValue(options.AppPassword, "APPLE_APP_PASSWORD"),
			Verbose:             options.Verbose,
		}

		if err := signing.NotarizeAndStaple(notarizeOpts); err != nil {
			return fmt.Errorf("notarization failed: %w", err)
		}

		pterm.Success.Println("Notarization complete!")
	}

	return nil
}

// setupMacOSSigning sets up macOS-specific signing options
func setupMacOSSigning(options *SignOptions, req *signing.SignRequest) error {
	// Only do auto-detection on macOS
	if runtime.GOOS == "darwin" {
		// Auto-detect signing identity if not specified and no certificate provided
		if req.Certificate.Identity == "" && req.Certificate.PKCS12Path == "" {
			pterm.Println("No signing identity specified, looking for Developer ID certificate...")
			identity, err := signing.FindDeveloperIDIdentity()
			if err != nil {
				pterm.Warning.Println("No Developer ID certificate found, using ad-hoc signing")
				req.Certificate.Identity = "-"
			} else {
				req.Certificate.Identity = identity.Name
				pterm.Printf("Using identity: %s\n", identity.Name)
			}
		}

		// Check for entitlements file
		if options.Entitlements == "" {
			defaultEntitlements := filepath.Join(filepath.Dir(options.Input), "..", "entitlements.plist")
			if _, err := os.Stat(defaultEntitlements); err == nil {
				req.Entitlements = defaultEntitlements
				pterm.Printf("Using entitlements: %s\n", req.Entitlements)
			}
		}
	} else {
		// Cross-platform macOS signing requires a certificate
		if req.Certificate.PKCS12Path == "" {
			return fmt.Errorf("cross-platform macOS signing requires --certificate flag with a .p12 file")
		}
	}

	return nil
}

// detectPlatformFromOptions detects the target platform from options and file
func detectPlatformFromOptions(options *SignOptions, info os.FileInfo) signing.Platform {
	// Explicit platform takes precedence
	if options.Platform != "" {
		switch strings.ToLower(options.Platform) {
		case "windows", "win":
			return signing.PlatformWindows
		case "darwin", "macos", "mac":
			return signing.PlatformMacOS
		case "linux":
			return signing.PlatformLinux
		}
	}

	// Detect from file extension/type
	input := strings.ToLower(options.Input)
	switch {
	case strings.HasSuffix(input, ".exe"),
		strings.HasSuffix(input, ".dll"),
		strings.HasSuffix(input, ".msi"),
		strings.HasSuffix(input, ".msix"),
		strings.HasSuffix(input, ".appx"),
		strings.HasSuffix(input, ".ps1"):
		return signing.PlatformWindows

	case strings.HasSuffix(input, ".app"),
		strings.HasSuffix(input, ".dmg"),
		strings.HasSuffix(input, ".pkg"):
		return signing.PlatformMacOS

	case strings.HasSuffix(input, ".deb"),
		strings.HasSuffix(input, ".rpm"):
		return signing.PlatformLinux

	default:
		// Fall back to current platform
		return signing.Platform(runtime.GOOS)
	}
}

// setupLinuxSigning sets up Linux-specific signing options
func setupLinuxSigning(options *SignOptions, req *signing.SignRequest) error {
	// Check for PGP key
	if req.Certificate.PGPKeyPath == "" {
		return fmt.Errorf("Linux package signing requires --pgp-key flag with a PGP private key file\n\nTo generate a new key:\n  wails3 signing generate-key --name \"Your Name\" --email \"you@example.com\"")
	}

	// Validate PGP key exists
	if _, err := os.Stat(req.Certificate.PGPKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("PGP key file not found: %s", req.Certificate.PGPKeyPath)
	}

	// Use role from options if provided (for DEB packages)
	if options.Role != "" {
		req.Description = options.Role
	}

	return nil
}

// ListIdentitiesOptions defines the options for listing signing identities
type ListIdentitiesOptions struct {
	NoColour bool `name:"n" description:"Disable colour output"`
}

// ListSigningIdentities lists available code signing identities
func ListSigningIdentities(options *ListIdentitiesOptions) error {
	if options.NoColour {
		term.DisableColor()
	}

	DisableFooter = true
	term.Header("Signing Identities")

	if runtime.GOOS != "darwin" {
		pterm.Println("Native signing identities are only available on macOS.")
		pterm.Println()
		pterm.Println("For cross-platform signing, use a .pfx/.p12 certificate file:")
		pterm.Println("  wails3 sign --input myapp.exe --certificate cert.pfx --certificate-password PASSWORD")
		return nil
	}

	identities, err := signing.ListSigningIdentities()
	if err != nil {
		return err
	}

	if len(identities) == 0 {
		pterm.Warning.Println("No signing identities found")
		pterm.Println("\nTo create a signing identity:")
		pterm.Println("  1. Open Keychain Access")
		pterm.Println("  2. Go to Certificate Assistant > Create a Certificate...")
		pterm.Println("  Or enroll in the Apple Developer Program for distribution certificates")
		return nil
	}

	pterm.Printf("Found %d signing %s:\n\n", len(identities), pluralize("identity", len(identities)))

	for _, id := range identities {
		status := pterm.Green("valid")
		if !id.IsValid {
			status = pterm.Red("invalid")
		}
		pterm.Printf("  %s [%s]\n", id.Name, status)
		pterm.Printf("    Hash: %s\n\n", id.Hash)
	}

	return nil
}

// SigningInfoOptions defines options for the signing info command
type SigningInfoOptions struct {
	NoColour bool `name:"n" description:"Disable colour output"`
}

// SigningInfo shows information about signing capabilities on this system
func SigningInfo(options *SigningInfoOptions) error {
	if options.NoColour {
		term.DisableColor()
	}

	DisableFooter = true
	term.Header("Signing Capabilities")

	pterm.Printf("Current platform: %s/%s\n\n", runtime.GOOS, runtime.GOARCH)

	// Check Windows signing
	pterm.Println("Windows signing:")
	winSigners := signing.DefaultRegistry.ListAvailableSigners(signing.PlatformWindows)
	if len(winSigners) == 0 {
		pterm.Printf("  %s No signers available\n", pterm.Red("✗"))
	} else {
		for _, s := range winSigners {
			pterm.Printf("  %s %s\n", pterm.Green("✓"), s.Backend())
		}
	}

	pterm.Println()

	// Check macOS signing
	pterm.Println("macOS signing:")
	macSigners := signing.DefaultRegistry.ListAvailableSigners(signing.PlatformMacOS)
	if len(macSigners) == 0 {
		pterm.Printf("  %s No signers available\n", pterm.Red("✗"))
	} else {
		for _, s := range macSigners {
			note := ""
			if s.Backend() == signing.BackendRelic {
				note = " (limited - use rcodesign for full support)"
			}
			pterm.Printf("  %s %s%s\n", pterm.Green("✓"), s.Backend(), note)
		}
	}

	pterm.Println()

	// Check Linux signing
	pterm.Println("Linux package signing:")
	linuxSigners := signing.DefaultRegistry.ListAvailableSigners(signing.PlatformLinux)
	if len(linuxSigners) == 0 {
		pterm.Printf("  %s No signers available\n", pterm.Red("✗"))
	} else {
		for _, s := range linuxSigners {
			pterm.Printf("  %s %s (DEB/RPM packages)\n", pterm.Green("✓"), s.Backend())
		}
	}

	pterm.Println()
	pterm.Println("Cross-platform signing:")
	pterm.Println("  Windows binaries can be signed from any platform using the relic library.")
	pterm.Println("  Linux packages (DEB/RPM) can be signed from any platform using the relic library.")
	pterm.Println("  macOS binaries should be signed on macOS or using rcodesign.")
	pterm.Println()
	pterm.Println("For more information: https://wails.io/docs/guides/signing")

	return nil
}

// StoreCredentialsOptions defines the options for storing notarization credentials
type StoreCredentialsOptions struct {
	ProfileName string `name:"profile" description:"Name for the credential profile" required:"true"`
	AppleID     string `name:"apple-id" description:"Apple ID" required:"true"`
	TeamID      string `name:"team-id" description:"Apple Developer Team ID" required:"true"`
	NoColour    bool   `name:"n" description:"Disable colour output"`
}

// StoreNotarizationCredentials stores notarization credentials in the keychain
func StoreNotarizationCredentials(options *StoreCredentialsOptions) error {
	if options.NoColour {
		term.DisableColor()
	}

	DisableFooter = true
	term.Header("Store Notarization Credentials")

	if runtime.GOOS != "darwin" {
		return fmt.Errorf("this command is only available on macOS")
	}

	pterm.Println("This will store your notarization credentials in the macOS keychain.")
	pterm.Println("You will be prompted for your app-specific password.")
	pterm.Println()

	// Prompt for password (this is interactive)
	pterm.Print("Enter your app-specific password: ")

	// Read password securely (this requires interactive input)
	return fmt.Errorf("interactive password input not yet implemented. Please run:\n  xcrun notarytool store-credentials %s --apple-id %s --team-id %s",
		options.ProfileName, options.AppleID, options.TeamID)
}

// loadSigningConfig loads signing options from a config file
func loadSigningConfig(options *SignOptions) error {
	data, err := os.ReadFile(options.Config)
	if err != nil {
		return err
	}

	var config struct {
		Signing signing.Config `yaml:"signing"`
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}

	if config.Signing.MacOS != nil {
		if options.Identity == "" && config.Signing.MacOS.Identity != "" {
			options.Identity = config.Signing.MacOS.Identity
		}
		if options.Entitlements == "" && config.Signing.MacOS.Entitlements != "" {
			options.Entitlements = config.Signing.MacOS.Entitlements
		}
		if config.Signing.MacOS.Notarization != nil {
			if options.KeychainProfile == "" {
				options.KeychainProfile = config.Signing.MacOS.Notarization.KeychainProfile
			}
			if options.AppleID == "" {
				options.AppleID = config.Signing.MacOS.Notarization.AppleID
			}
			if options.TeamID == "" {
				options.TeamID = config.Signing.MacOS.Notarization.TeamID
			}
			if options.AppPassword == "" {
				options.AppPassword = config.Signing.MacOS.Notarization.AppSpecificPassword
			}
		}
	}

	if config.Signing.Windows != nil {
		if options.Certificate == "" && config.Signing.Windows.CertificatePath != "" {
			options.Certificate = config.Signing.Windows.CertificatePath
		}
		if options.CertificatePassword == "" && config.Signing.Windows.CertificatePassword != "" {
			options.CertificatePassword = config.Signing.Windows.CertificatePassword
		}
		if options.Thumbprint == "" && config.Signing.Windows.CertificateThumbprint != "" {
			options.Thumbprint = config.Signing.Windows.CertificateThumbprint
		}
		if options.TimestampServer == "" && config.Signing.Windows.TimestampServer != "" {
			options.TimestampServer = config.Signing.Windows.TimestampServer
		}
	}

	return nil
}

// findSignableFile looks for a signable file in common locations
func findSignableFile() string {
	// Look for .app bundles (macOS)
	if runtime.GOOS == "darwin" {
		matches, _ := filepath.Glob("bin/*.app")
		if len(matches) > 0 {
			return matches[0]
		}
	}

	// Look for .exe files (Windows)
	matches, _ := filepath.Glob("bin/*.exe")
	if len(matches) > 0 {
		return matches[0]
	}

	// Look in current directory
	if runtime.GOOS == "darwin" {
		matches, _ = filepath.Glob("*.app")
		if len(matches) > 0 {
			return matches[0]
		}
	}

	matches, _ = filepath.Glob("*.exe")
	if len(matches) > 0 {
		return matches[0]
	}

	// Look in build/bin
	if runtime.GOOS == "darwin" {
		matches, _ = filepath.Glob("build/bin/*.app")
		if len(matches) > 0 {
			return matches[0]
		}
	}

	matches, _ = filepath.Glob("build/bin/*.exe")
	if len(matches) > 0 {
		return matches[0]
	}

	return ""
}

// getEnvOrValue returns the environment variable value if it exists, otherwise the provided value
func getEnvOrValue(value, envVar string) string {
	if value != "" {
		return value
	}
	return os.Getenv(envVar)
}

func pluralize(word string, count int) string {
	if count == 1 {
		return word
	}
	return word + "s"
}

// GenerateKeyOptions defines the options for generating a PGP key
type GenerateKeyOptions struct {
	// Name is the name associated with the key
	Name string `name:"name" description:"Name for the key (e.g., 'John Doe')" required:"true"`
	// Email is the email associated with the key
	Email string `name:"email" description:"Email for the key" required:"true"`
	// Comment is an optional comment
	Comment string `name:"comment" description:"Optional comment (e.g., 'Package Signing Key')"`
	// Output is the output directory for the keys
	Output string `name:"output" description:"Output directory for key files" default:"."`
	// Password is the password to encrypt the private key
	Password string `name:"password" description:"Password to encrypt the private key (recommended)"`
	// KeyBits is the RSA key size
	KeyBits int `name:"bits" description:"RSA key size" default:"4096"`
	// NoColour disables colour output
	NoColour bool `name:"n" description:"Disable colour output"`
}

// GeneratePGPKey generates a new PGP key pair for Linux package signing
func GeneratePGPKey(options *GenerateKeyOptions) error {
	if options.NoColour {
		term.DisableColor()
	}

	DisableFooter = true
	term.Header("Generate PGP Signing Key")

	// Determine output paths
	baseName := strings.ReplaceAll(strings.ToLower(options.Name), " ", "_")
	privateKeyPath := filepath.Join(options.Output, baseName+"_private.asc")
	publicKeyPath := filepath.Join(options.Output, baseName+"_public.asc")

	// Check if files already exist
	if _, err := os.Stat(privateKeyPath); err == nil {
		return fmt.Errorf("private key already exists: %s", privateKeyPath)
	}
	if _, err := os.Stat(publicKeyPath); err == nil {
		return fmt.Errorf("public key already exists: %s", publicKeyPath)
	}

	pterm.Printf("Generating %d-bit RSA key pair...\n", options.KeyBits)
	pterm.Printf("  Name: %s\n", options.Name)
	pterm.Printf("  Email: %s\n", options.Email)
	if options.Comment != "" {
		pterm.Printf("  Comment: %s\n", options.Comment)
	}
	pterm.Println()

	config := signing.PGPKeyConfig{
		Name:     options.Name,
		Email:    options.Email,
		Comment:  options.Comment,
		KeyBits:  options.KeyBits,
		Password: options.Password,
	}

	result, err := signing.GeneratePGPKey(config, privateKeyPath, publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	pterm.Success.Println("Key pair generated successfully!")
	pterm.Println()
	pterm.Printf("Private key: %s\n", result.PrivateKeyPath)
	pterm.Printf("Public key:  %s\n", result.PublicKeyPath)
	pterm.Printf("Key ID:      %s\n", result.KeyID)
	pterm.Printf("Fingerprint: %s\n", result.Fingerprint)
	pterm.Println()

	if options.Password == "" {
		pterm.Warning.Println("Private key is NOT encrypted. Consider using --password for production keys.")
	} else {
		pterm.Info.Println("Private key is encrypted with your password.")
	}

	pterm.Println()
	pterm.Println("To sign a package:")
	pterm.Printf("  wails3 sign --input myapp.deb --pgp-key %s\n", result.PrivateKeyPath)
	pterm.Println()
	pterm.Println("Distribute the public key to users so they can verify your packages:")
	pterm.Printf("  %s\n", result.PublicKeyPath)

	return nil
}

// KeyInfoOptions defines options for displaying PGP key information
type KeyInfoOptions struct {
	// KeyPath is the path to the PGP key file
	KeyPath string `name:"key" description:"Path to PGP key file" required:"true"`
	// NoColour disables colour output
	NoColour bool `name:"n" description:"Disable colour output"`
}

// PGPKeyInfo displays information about a PGP key
func PGPKeyInfo(options *KeyInfoOptions) error {
	if options.NoColour {
		term.DisableColor()
	}

	DisableFooter = true
	term.Header("PGP Key Information")

	info, err := signing.GetPGPKeyInfo(options.KeyPath)
	if err != nil {
		return fmt.Errorf("failed to read key: %w", err)
	}

	pterm.Printf("Key ID:      %s\n", info.KeyID)
	pterm.Printf("Fingerprint: %s\n", info.Fingerprint)
	pterm.Printf("Created:     %s\n", info.CreatedAt.Format("2006-01-02 15:04:05"))

	if info.ExpiresAt != nil {
		status := pterm.Green("valid")
		if info.IsExpired() {
			status = pterm.Red("expired")
		}
		pterm.Printf("Expires:     %s [%s]\n", info.ExpiresAt.Format("2006-01-02"), status)
	} else {
		pterm.Println("Expires:     never")
	}

	pterm.Println()
	pterm.Println("User IDs:")
	for _, uid := range info.UserIDs {
		pterm.Printf("  %s\n", uid)
	}

	pterm.Println()
	if info.HasPrivate {
		if info.IsEncrypted {
			pterm.Info.Println("This is a private key (encrypted)")
		} else {
			pterm.Warning.Println("This is a private key (NOT encrypted)")
		}
	} else {
		pterm.Info.Println("This is a public key")
	}

	return nil
}
