//go:build !darwin

package signing

import "fmt"

// SigningIdentity represents a code signing identity found on the system
type SigningIdentity struct {
	Hash    string
	Name    string
	IsValid bool
}

// ListSigningIdentities returns all available code signing identities on the system
func ListSigningIdentities() ([]SigningIdentity, error) {
	return nil, fmt.Errorf("signing identities are only available on macOS")
}

// FindDeveloperIDIdentity finds a "Developer ID Application" identity for distribution
func FindDeveloperIDIdentity() (*SigningIdentity, error) {
	return nil, fmt.Errorf("signing identities are only available on macOS")
}

// SignOptions defines options for code signing
type SignOptions struct {
	AppPath         string
	Identity        string
	Entitlements    string
	HardenedRuntime bool
	Deep            bool
	Force           bool
	Verbose         bool
}

// Sign signs a macOS application or binary
func Sign(options SignOptions) error {
	return fmt.Errorf("code signing is only available on macOS")
}

// SignAppBundle signs an entire .app bundle
func SignAppBundle(bundlePath string, identity string, entitlements string, hardenedRuntime bool) error {
	return fmt.Errorf("code signing is only available on macOS")
}

// VerifySignature verifies the code signature of an application
func VerifySignature(appPath string) error {
	return fmt.Errorf("signature verification is only available on macOS")
}

// VerifyNotarization checks if an application has been notarized
func VerifyNotarization(appPath string) error {
	return fmt.Errorf("notarization verification is only available on macOS")
}
