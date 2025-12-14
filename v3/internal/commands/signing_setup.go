package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/charmbracelet/huh"
	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/keychain"
)

// SigningSetup configures signing variables in platform Taskfiles
func SigningSetup(options *flags.SigningSetup) error {
	// Determine which platforms to configure
	platforms := options.Platforms
	if len(platforms) == 0 {
		// Auto-detect based on existing Taskfiles
		platforms = detectPlatforms()
		if len(platforms) == 0 {
			return fmt.Errorf("no platform Taskfiles found in build/ directory")
		}
	}

	for _, platform := range platforms {
		var err error
		switch platform {
		case "darwin":
			err = setupDarwinSigning()
		case "windows":
			err = setupWindowsSigning()
		case "linux":
			err = setupLinuxSigning()
		default:
			pterm.Warning.Printfln("Unknown platform: %s", platform)
			continue
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func detectPlatforms() []string {
	var platforms []string
	for _, p := range []string{"darwin", "windows", "linux"} {
		taskfile := filepath.Join("build", p, "Taskfile.yml")
		if _, err := os.Stat(taskfile); err == nil {
			platforms = append(platforms, p)
		}
	}
	return platforms
}

func setupDarwinSigning() error {
	pterm.DefaultHeader.Println("macOS Code Signing Setup")
	fmt.Println()

	// Determine signing method based on platform
	var signingMethod string
	onMacOS := runtime.GOOS == "darwin"

	if onMacOS {
		// On macOS, offer choice between native and cross-platform
		methodForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Signing method").
					Description("Choose how to sign macOS binaries").
					Options(
						huh.NewOption("Native (codesign) - requires macOS", "native"),
						huh.NewOption("Cross-platform (P12) - works on any OS", "crossplatform"),
					).
					Value(&signingMethod),
			),
		)
		if err := methodForm.Run(); err != nil {
			return err
		}
	} else {
		// Not on macOS, must use cross-platform
		signingMethod = "crossplatform"
		pterm.Info.Println("Not running on macOS - using cross-platform signing")
		fmt.Println()
	}

	if signingMethod == "native" {
		return setupDarwinSigningNative()
	}
	return setupDarwinSigningCrossPlatform()
}

// setupDarwinSigningNative configures native macOS signing using codesign
func setupDarwinSigningNative() error {
	// Get available signing identities
	identities, err := getMacOSSigningIdentities()
	if err != nil {
		pterm.Warning.Printfln("Could not list signing identities: %v", err)
		identities = []string{}
	}

	var signIdentity string
	var keychainProfile string
	var entitlements string
	var configureNotarization bool

	// Build identity options
	var identityOptions []huh.Option[string]
	for _, id := range identities {
		identityOptions = append(identityOptions, huh.NewOption(id, id))
	}
	identityOptions = append(identityOptions, huh.NewOption("Enter manually...", "manual"))

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select signing identity").
				Description("Choose your Developer ID certificate").
				Options(identityOptions...).
				Value(&signIdentity),
		).WithHideFunc(func() bool {
			return len(identities) == 0
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("Signing identity").
				Description("e.g., Developer ID Application: Your Company (TEAMID)").
				Placeholder("Developer ID Application: ...").
				Value(&signIdentity),
		).WithHideFunc(func() bool {
			return len(identities) > 0 && signIdentity != "manual"
		}),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Configure notarization?").
				Description("Required for distributing apps outside the App Store").
				Value(&configureNotarization),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("Keychain profile name").
				Description("The profile name used with 'wails3 signing credentials'").
				Placeholder("my-notarize-profile").
				Value(&keychainProfile),
		).WithHideFunc(func() bool {
			return !configureNotarization
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("Entitlements file (optional)").
				Description("Path to entitlements plist, leave empty to skip").
				Placeholder("build/darwin/entitlements.plist").
				Value(&entitlements),
		),
	)

	err = form.Run()
	if err != nil {
		return err
	}

	// Handle manual entry
	if signIdentity == "manual" {
		signIdentity = ""
	}

	// Update Taskfile
	taskfilePath := filepath.Join("build", "darwin", "Taskfile.yml")
	err = updateTaskfileVars(taskfilePath, map[string]string{
		"SIGN_IDENTITY":    signIdentity,
		"KEYCHAIN_PROFILE": keychainProfile,
		"ENTITLEMENTS":     entitlements,
	})
	if err != nil {
		return err
	}

	pterm.Success.Printfln("Updated %s", taskfilePath)

	if configureNotarization && keychainProfile != "" {
		fmt.Println()
		pterm.Info.Println("Next step: Store your notarization credentials:")
		fmt.Println()
		pterm.Println(pterm.LightBlue(fmt.Sprintf(`  wails3 signing credentials \
    --apple-id "your@email.com" \
    --team-id "TEAMID" \
    --password "app-specific-password" \
    --profile "%s"`, keychainProfile)))
		fmt.Println()
	}

	return nil
}

// setupDarwinSigningCrossPlatform configures cross-platform macOS signing using P12 certificates
func setupDarwinSigningCrossPlatform() error {
	var p12Path string
	var p12Password string
	var configureNotarization bool
	var notaryKeyPath string
	var notaryKeyID string
	var notaryIssuer string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("P12 certificate path").
				Description("Path to your Developer ID certificate exported as .p12").
				Placeholder("certs/developer-id.p12").
				Value(&p12Path).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("P12 certificate path is required")
					}
					return nil
				}),

			huh.NewInput().
				Title("P12 password").
				Description("Stored securely in system keychain").
				EchoMode(huh.EchoModePassword).
				Value(&p12Password).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("P12 password is required")
					}
					return nil
				}),
		),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Configure notarization?").
				Description("Required for distributing apps outside the App Store").
				Value(&configureNotarization),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("Apple API key path (.p8)").
				Description("Download from App Store Connect > Users and Access > Keys").
				Placeholder("certs/AuthKey_XXXXXXXX.p8").
				Value(&notaryKeyPath).
				Validate(func(s string) error {
					if configureNotarization && s == "" {
						return fmt.Errorf("API key path is required for notarization")
					}
					return nil
				}),

			huh.NewInput().
				Title("API Key ID").
				Description("The Key ID from App Store Connect (e.g., ABC123DEFG)").
				Placeholder("ABC123DEFG").
				Value(&notaryKeyID).
				Validate(func(s string) error {
					if configureNotarization && s == "" {
						return fmt.Errorf("API Key ID is required for notarization")
					}
					return nil
				}),

			huh.NewInput().
				Title("Team ID (Issuer ID)").
				Description("Your Apple Developer Team ID").
				Placeholder("TEAMID123").
				Value(&notaryIssuer).
				Validate(func(s string) error {
					if configureNotarization && s == "" {
						return fmt.Errorf("Team ID is required for notarization")
					}
					return nil
				}),
		).WithHideFunc(func() bool {
			return !configureNotarization
		}),
	)

	err := form.Run()
	if err != nil {
		return err
	}

	// Store P12 password in keychain
	err = keychain.Set(keychain.KeyMacOSP12Password, p12Password)
	if err != nil {
		return fmt.Errorf("failed to store P12 password in keychain: %w", err)
	}
	pterm.Success.Println("P12 password stored in system keychain")

	// Store notarization credentials in keychain if configured
	if configureNotarization {
		err = keychain.Set(keychain.KeyNotaryKeyID, notaryKeyID)
		if err != nil {
			return fmt.Errorf("failed to store notary key ID in keychain: %w", err)
		}

		err = keychain.Set(keychain.KeyNotaryIssuer, notaryIssuer)
		if err != nil {
			return fmt.Errorf("failed to store notary issuer in keychain: %w", err)
		}
		pterm.Success.Println("Notarization credentials stored in system keychain")
	}

	// Update Taskfile with non-sensitive values
	taskfilePath := filepath.Join("build", "darwin", "Taskfile.yml")
	vars := map[string]string{
		"P12_CERTIFICATE": p12Path,
	}
	if configureNotarization {
		vars["NOTARY_KEY"] = notaryKeyPath
	}

	err = updateTaskfileVars(taskfilePath, vars)
	if err != nil {
		return err
	}

	pterm.Success.Printfln("Updated %s", taskfilePath)

	fmt.Println()
	pterm.Info.Println("Cross-platform signing configured!")
	fmt.Println()
	pterm.Println("To sign a macOS binary from any platform:")
	pterm.Println(pterm.LightBlue(fmt.Sprintf("  wails3 tool sign --input myapp --p12 %s", p12Path)))
	fmt.Println()

	if configureNotarization {
		pterm.Println("To sign and notarize:")
		pterm.Println(pterm.LightBlue(fmt.Sprintf("  wails3 tool sign --input myapp --p12 %s --notarize --notary-key %s",
			p12Path, notaryKeyPath)))
		fmt.Println()
	}

	return nil
}

func setupWindowsSigning() error {
	pterm.DefaultHeader.Println("Windows Code Signing Setup")
	fmt.Println()

	var certSource string
	var certPath string
	var certPassword string
	var thumbprint string
	var timestampServer string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Certificate source").
				Options(
					huh.NewOption("Certificate file (.pfx/.p12)", "file"),
					huh.NewOption("Windows certificate store (thumbprint)", "store"),
				).
				Value(&certSource),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("Certificate path").
				Description("Path to your .pfx or .p12 file").
				Placeholder("certs/signing.pfx").
				Value(&certPath),

			huh.NewInput().
				Title("Certificate password").
				Description("Stored securely in system keychain").
				EchoMode(huh.EchoModePassword).
				Value(&certPassword),
		).WithHideFunc(func() bool {
			return certSource != "file"
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("Certificate thumbprint").
				Description("SHA-1 thumbprint of the certificate in Windows store").
				Placeholder("ABC123DEF456...").
				Value(&thumbprint),
		).WithHideFunc(func() bool {
			return certSource != "store"
		}),

		huh.NewGroup(
			huh.NewInput().
				Title("Timestamp server (optional)").
				Description("Leave empty for default: http://timestamp.digicert.com").
				Placeholder("http://timestamp.digicert.com").
				Value(&timestampServer),
		),
	)

	err := form.Run()
	if err != nil {
		return err
	}

	// Store password in keychain if provided
	if certPassword != "" {
		err = keychain.Set(keychain.KeyWindowsCertPassword, certPassword)
		if err != nil {
			return fmt.Errorf("failed to store password in keychain: %w", err)
		}
		pterm.Success.Println("Certificate password stored in system keychain")
	}

	// Update Taskfile (no passwords stored here)
	taskfilePath := filepath.Join("build", "windows", "Taskfile.yml")
	vars := map[string]string{
		"TIMESTAMP_SERVER": timestampServer,
	}

	if certSource == "file" {
		vars["SIGN_CERTIFICATE"] = certPath
	} else {
		vars["SIGN_THUMBPRINT"] = thumbprint
	}

	err = updateTaskfileVars(taskfilePath, vars)
	if err != nil {
		return err
	}

	pterm.Success.Printfln("Updated %s", taskfilePath)
	return nil
}

func setupLinuxSigning() error {
	pterm.DefaultHeader.Println("Linux Package Signing Setup")
	fmt.Println()

	var keySource string
	var keyPath string
	var keyPassword string
	var signRole string

	// For key generation
	var genName string
	var genEmail string
	var genPassword string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("PGP key source").
				Options(
					huh.NewOption("Use existing key", "existing"),
					huh.NewOption("Generate new key", "generate"),
				).
				Value(&keySource),
		),

		// Existing key options
		huh.NewGroup(
			huh.NewInput().
				Title("PGP private key path").
				Description("Path to your ASCII-armored private key file").
				Placeholder("signing-key.asc").
				Value(&keyPath),

			huh.NewInput().
				Title("Key password (if encrypted)").
				Description("Stored securely in system keychain").
				EchoMode(huh.EchoModePassword).
				Value(&keyPassword),
		).WithHideFunc(func() bool {
			return keySource != "existing"
		}),

		// Key generation options
		huh.NewGroup(
			huh.NewInput().
				Title("Name").
				Description("Name for the PGP key").
				Placeholder("Your Name").
				Value(&genName).
				Validate(func(s string) error {
					if keySource == "generate" && s == "" {
						return fmt.Errorf("name is required")
					}
					return nil
				}),

			huh.NewInput().
				Title("Email").
				Description("Email for the PGP key").
				Placeholder("you@example.com").
				Value(&genEmail).
				Validate(func(s string) error {
					if keySource == "generate" && s == "" {
						return fmt.Errorf("email is required")
					}
					return nil
				}),

			huh.NewInput().
				Title("Key password (optional but recommended)").
				Description("Stored securely in system keychain").
				EchoMode(huh.EchoModePassword).
				Value(&genPassword),
		).WithHideFunc(func() bool {
			return keySource != "generate"
		}),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("DEB signing role").
				Description("Role for signing Debian packages").
				Options(
					huh.NewOption("builder (default)", "builder"),
					huh.NewOption("origin", "origin"),
					huh.NewOption("maint", "maint"),
					huh.NewOption("archive", "archive"),
				).
				Value(&signRole),
		),
	)

	err := form.Run()
	if err != nil {
		return err
	}

	// Generate key if requested
	if keySource == "generate" {
		keyPath = "signing-key.asc"
		pubKeyPath := "signing-key.pub.asc"

		pterm.Info.Println("Generating PGP key pair...")

		// Call the key generation
		err = generatePGPKeyForSetup(genName, genEmail, genPassword, keyPath, pubKeyPath)
		if err != nil {
			return fmt.Errorf("failed to generate key: %w", err)
		}

		keyPassword = genPassword

		pterm.Success.Printfln("Generated %s and %s", keyPath, pubKeyPath)
		fmt.Println()
		pterm.Info.Println("Distribute the public key to users so they can verify your packages:")
		pterm.Println(pterm.LightBlue(fmt.Sprintf("  # For apt: sudo cp %s /etc/apt/trusted.gpg.d/", pubKeyPath)))
		pterm.Println(pterm.LightBlue(fmt.Sprintf("  # For rpm: sudo rpm --import %s", pubKeyPath)))
		fmt.Println()
	}

	// Store password in keychain if provided
	if keyPassword != "" {
		err = keychain.Set(keychain.KeyPGPPassword, keyPassword)
		if err != nil {
			return fmt.Errorf("failed to store password in keychain: %w", err)
		}
		pterm.Success.Println("PGP key password stored in system keychain")
	}

	// Update Taskfile (no passwords stored here)
	taskfilePath := filepath.Join("build", "linux", "Taskfile.yml")
	vars := map[string]string{
		"PGP_KEY": keyPath,
	}
	if signRole != "" && signRole != "builder" {
		vars["SIGN_ROLE"] = signRole
	}

	err = updateTaskfileVars(taskfilePath, vars)
	if err != nil {
		return err
	}

	pterm.Success.Printfln("Updated %s", taskfilePath)
	return nil
}

// getMacOSSigningIdentities returns available signing identities on macOS
func getMacOSSigningIdentities() ([]string, error) {
	if runtime.GOOS != "darwin" {
		return nil, fmt.Errorf("not running on macOS")
	}

	// Run security find-identity to get available codesigning identities
	cmd := exec.Command("security", "find-identity", "-v", "-p", "codesigning")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run security find-identity: %w", err)
	}

	var identities []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// Lines look like: 1) ABC123... "Developer ID Application: Company Name (TEAMID)"
		// We want to extract the quoted part
		if strings.Contains(line, "\"") {
			start := strings.Index(line, "\"")
			end := strings.LastIndex(line, "\"")
			if start != -1 && end > start {
				identity := line[start+1 : end]
				// Filter for Developer ID certificates (most useful for distribution)
				if strings.Contains(identity, "Developer ID") {
					identities = append(identities, identity)
				}
			}
		}
	}

	return identities, nil
}

// updateTaskfileVars updates the vars section of a Taskfile
func updateTaskfileVars(path string, vars map[string]string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	lines := strings.Split(string(content), "\n")
	var result []string
	inVars := false
	varsInserted := false
	remainingVars := make(map[string]string)
	for k, v := range vars {
		remainingVars[k] = v
	}

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Detect vars section
		if trimmed == "vars:" {
			inVars = true
			result = append(result, line)
			continue
		}

		// Detect end of vars section (next top-level key or tasks:)
		if inVars && len(line) > 0 && line[0] != ' ' && line[0] != '\t' && !strings.HasPrefix(trimmed, "#") {
			// Insert any remaining vars before leaving vars section
			for k, v := range remainingVars {
				if v != "" {
					result = append(result, fmt.Sprintf("  %s: %q", k, v))
				}
			}
			remainingVars = make(map[string]string)
			varsInserted = true
			inVars = false
		}

		if inVars {
			// Check if this line is a var we want to update
			updated := false
			for k, v := range remainingVars {
				commentedKey := "# " + k + ":"
				uncommentedKey := k + ":"

				if strings.Contains(trimmed, commentedKey) || strings.HasPrefix(trimmed, uncommentedKey) {
					if v != "" {
						// Uncomment and set value
						result = append(result, fmt.Sprintf("  %s: %q", k, v))
					} else {
						// Keep as comment
						result = append(result, line)
					}
					delete(remainingVars, k)
					updated = true
					break
				}
			}
			if !updated {
				result = append(result, line)
			}
		} else {
			result = append(result, line)
		}

		// If we're at the end and haven't inserted vars yet, we need to add vars section
		if i == len(lines)-1 && !varsInserted && len(remainingVars) > 0 {
			// Find where to insert (after includes, before tasks)
			// For simplicity, just append warning
			pterm.Warning.Println("Could not find vars section in Taskfile, please add manually")
		}
	}

	return os.WriteFile(path, []byte(strings.Join(result, "\n")), 0644)
}

// generatePGPKeyForSetup generates a PGP key pair for signing packages
func generatePGPKeyForSetup(name, email, password, privatePath, publicPath string) error {
	// Create a new entity (key pair)
	config := &packet.Config{
		DefaultHash:            0, // Use default
		DefaultCipher:          0, // Use default
		DefaultCompressionAlgo: 0, // Use default
	}

	entity, err := openpgp.NewEntity(name, "", email, config)
	if err != nil {
		return fmt.Errorf("failed to create PGP entity: %w", err)
	}

	// Encrypt the private key if password is provided
	if password != "" {
		err = entity.PrivateKey.Encrypt([]byte(password))
		if err != nil {
			return fmt.Errorf("failed to encrypt private key: %w", err)
		}
		// Also encrypt subkeys
		for _, subkey := range entity.Subkeys {
			if subkey.PrivateKey != nil {
				err = subkey.PrivateKey.Encrypt([]byte(password))
				if err != nil {
					return fmt.Errorf("failed to encrypt subkey: %w", err)
				}
			}
		}
	}

	// Write private key
	privateFile, err := os.Create(privatePath)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %w", err)
	}
	defer privateFile.Close()

	privateArmor, err := armor.Encode(privateFile, openpgp.PrivateKeyType, nil)
	if err != nil {
		return fmt.Errorf("failed to create armor encoder: %w", err)
	}

	err = entity.SerializePrivate(privateArmor, config)
	if err != nil {
		return fmt.Errorf("failed to serialize private key: %w", err)
	}
	privateArmor.Close()

	// Write public key
	publicFile, err := os.Create(publicPath)
	if err != nil {
		return fmt.Errorf("failed to create public key file: %w", err)
	}
	defer publicFile.Close()

	publicArmor, err := armor.Encode(publicFile, openpgp.PublicKeyType, nil)
	if err != nil {
		return fmt.Errorf("failed to create armor encoder: %w", err)
	}

	err = entity.Serialize(publicArmor)
	if err != nil {
		return fmt.Errorf("failed to serialize public key: %w", err)
	}
	publicArmor.Close()

	return nil
}
