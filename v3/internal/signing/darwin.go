//go:build darwin

package signing

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// SigningIdentity represents a code signing identity found on the system
type SigningIdentity struct {
	// Hash is the SHA-1 hash of the certificate
	Hash string
	// Name is the common name of the certificate
	Name string
	// IsValid indicates if the certificate is valid (not expired, not revoked)
	IsValid bool
}

// ListSigningIdentities returns all available code signing identities on the system
func ListSigningIdentities() ([]SigningIdentity, error) {
	cmd := exec.Command("security", "find-identity", "-v", "-p", "codesigning")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list signing identities: %w", err)
	}

	return parseSigningIdentities(string(output))
}

// parseSigningIdentities parses the output of `security find-identity`
func parseSigningIdentities(output string) ([]SigningIdentity, error) {
	var identities []SigningIdentity

	// Pattern matches lines like:
	// 1) ABC123... "Developer ID Application: Company Name (TEAMID)"
	pattern := regexp.MustCompile(`^\s*\d+\)\s+([A-Fa-f0-9]+)\s+"([^"]+)"(?:\s+\(CSSMERR_TP_CERT_EXPIRED\))?`)

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		matches := pattern.FindStringSubmatch(line)
		if len(matches) >= 3 {
			isValid := !strings.Contains(line, "CSSMERR_TP_CERT_EXPIRED") &&
				!strings.Contains(line, "CSSMERR_TP_CERT_REVOKED")

			identities = append(identities, SigningIdentity{
				Hash:    matches[1],
				Name:    matches[2],
				IsValid: isValid,
			})
		}
	}

	return identities, scanner.Err()
}

// FindDeveloperIDIdentity finds a "Developer ID Application" identity for distribution
func FindDeveloperIDIdentity() (*SigningIdentity, error) {
	identities, err := ListSigningIdentities()
	if err != nil {
		return nil, err
	}

	for _, id := range identities {
		if id.IsValid && strings.HasPrefix(id.Name, "Developer ID Application:") {
			return &id, nil
		}
	}

	return nil, fmt.Errorf("no valid 'Developer ID Application' identity found")
}

// SignOptions defines options for code signing
type SignOptions struct {
	// AppPath is the path to the .app bundle or binary to sign
	AppPath string
	// Identity is the signing identity (use "-" for ad-hoc)
	Identity string
	// Entitlements is the path to the entitlements plist file
	Entitlements string
	// HardenedRuntime enables the hardened runtime
	HardenedRuntime bool
	// Deep performs deep signing (signs all nested code)
	Deep bool
	// Force replaces any existing signature
	Force bool
	// Verbose enables verbose output
	Verbose bool
}

// Sign signs a macOS application or binary
func Sign(options SignOptions) error {
	if options.AppPath == "" {
		return fmt.Errorf("app path is required")
	}

	// Check if path exists
	if _, err := os.Stat(options.AppPath); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", options.AppPath)
	}

	// Default to ad-hoc signing if no identity specified
	if options.Identity == "" {
		options.Identity = "-"
	}

	args := []string{}

	if options.Force {
		args = append(args, "--force")
	}

	if options.Deep {
		args = append(args, "--deep")
	}

	if options.Verbose {
		args = append(args, "--verbose")
	}

	args = append(args, "--sign", options.Identity)

	if options.HardenedRuntime {
		args = append(args, "--options", "runtime")
	}

	if options.Entitlements != "" {
		if _, err := os.Stat(options.Entitlements); os.IsNotExist(err) {
			return fmt.Errorf("entitlements file does not exist: %s", options.Entitlements)
		}
		args = append(args, "--entitlements", options.Entitlements)
	}

	args = append(args, options.AppPath)

	cmd := exec.Command("codesign", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("codesign failed: %w", err)
	}

	return nil
}

// SignAppBundle signs an entire .app bundle with proper handling of nested components
// This follows Apple's recommended signing order: frameworks first, then the main app
func SignAppBundle(bundlePath string, identity string, entitlements string, hardenedRuntime bool) error {
	if !strings.HasSuffix(bundlePath, ".app") {
		return fmt.Errorf("path must be an .app bundle: %s", bundlePath)
	}

	contentsPath := filepath.Join(bundlePath, "Contents")
	macOSPath := filepath.Join(contentsPath, "MacOS")
	frameworksPath := filepath.Join(contentsPath, "Frameworks")

	// Sign frameworks first if they exist
	if info, err := os.Stat(frameworksPath); err == nil && info.IsDir() {
		entries, err := os.ReadDir(frameworksPath)
		if err != nil {
			return fmt.Errorf("failed to read Frameworks directory: %w", err)
		}

		for _, entry := range entries {
			frameworkPath := filepath.Join(frameworksPath, entry.Name())
			if err := Sign(SignOptions{
				AppPath:         frameworkPath,
				Identity:        identity,
				HardenedRuntime: hardenedRuntime,
				Force:           true,
			}); err != nil {
				return fmt.Errorf("failed to sign framework %s: %w", entry.Name(), err)
			}
		}
	}

	// Sign any dylibs in MacOS directory
	if info, err := os.Stat(macOSPath); err == nil && info.IsDir() {
		entries, err := os.ReadDir(macOSPath)
		if err != nil {
			return fmt.Errorf("failed to read MacOS directory: %w", err)
		}

		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".dylib") {
				dylibPath := filepath.Join(macOSPath, entry.Name())
				if err := Sign(SignOptions{
					AppPath:         dylibPath,
					Identity:        identity,
					HardenedRuntime: hardenedRuntime,
					Force:           true,
				}); err != nil {
					return fmt.Errorf("failed to sign dylib %s: %w", entry.Name(), err)
				}
			}
		}
	}

	// Sign the main app bundle
	return Sign(SignOptions{
		AppPath:         bundlePath,
		Identity:        identity,
		Entitlements:    entitlements,
		HardenedRuntime: hardenedRuntime,
		Force:           true,
		Deep:            false, // We already signed nested components
	})
}

// VerifySignature verifies the code signature of an application
func VerifySignature(appPath string) error {
	cmd := exec.Command("codesign", "--verify", "--verbose=2", appPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("signature verification failed: %s", stderr.String())
	}

	return nil
}

// VerifyNotarization checks if an application has been notarized
func VerifyNotarization(appPath string) error {
	cmd := exec.Command("spctl", "--assess", "--verbose=2", "--type", "execute", appPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("notarization verification failed: %s", stderr.String())
	}

	return nil
}
