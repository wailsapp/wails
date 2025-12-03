package signing

// Config defines the signing configuration for all platforms
type Config struct {
	MacOS   *MacOSConfig   `yaml:"macos,omitempty" json:"macos,omitempty"`
	Windows *WindowsConfig `yaml:"windows,omitempty" json:"windows,omitempty"`
}

// MacOSConfig defines macOS-specific signing configuration
type MacOSConfig struct {
	// Identity is the signing identity (e.g., "Developer ID Application: Company Name (TEAMID)")
	// Use "-" for ad-hoc signing, or leave empty to auto-detect
	Identity string `yaml:"identity,omitempty" json:"identity,omitempty"`

	// Entitlements is the path to the entitlements plist file
	Entitlements string `yaml:"entitlements,omitempty" json:"entitlements,omitempty"`

	// HardenedRuntime enables the hardened runtime (required for notarization)
	HardenedRuntime bool `yaml:"hardened_runtime,omitempty" json:"hardened_runtime,omitempty"`

	// Notarization contains notarization-specific settings
	Notarization *NotarizationConfig `yaml:"notarization,omitempty" json:"notarization,omitempty"`
}

// NotarizationConfig defines macOS notarization settings
type NotarizationConfig struct {
	// Enabled determines whether to notarize the app
	Enabled bool `yaml:"enabled,omitempty" json:"enabled,omitempty"`

	// AppleID is the Apple ID used for notarization (can use env var reference like "${APPLE_ID}")
	AppleID string `yaml:"apple_id,omitempty" json:"apple_id,omitempty"`

	// TeamID is the Apple Developer Team ID
	TeamID string `yaml:"team_id,omitempty" json:"team_id,omitempty"`

	// KeychainProfile is the name of the keychain profile storing credentials
	// Created via: xcrun notarytool store-credentials "profile-name"
	KeychainProfile string `yaml:"keychain_profile,omitempty" json:"keychain_profile,omitempty"`

	// AppSpecificPassword is an app-specific password for notarization
	// Use this instead of keychain_profile if preferred
	AppSpecificPassword string `yaml:"app_specific_password,omitempty" json:"app_specific_password,omitempty"`
}

// WindowsConfig defines Windows-specific signing configuration
type WindowsConfig struct {
	// CertificatePath is the path to the .pfx certificate file
	CertificatePath string `yaml:"certificate_path,omitempty" json:"certificate_path,omitempty"`

	// CertificatePassword is the password for the certificate
	CertificatePassword string `yaml:"certificate_password,omitempty" json:"certificate_password,omitempty"`

	// CertificateThumbprint is the thumbprint of a certificate in the Windows Certificate Store
	// Use this instead of certificate_path for certificates stored in the system
	CertificateThumbprint string `yaml:"certificate_thumbprint,omitempty" json:"certificate_thumbprint,omitempty"`

	// TimestampServer is the URL of the timestamp server (default: http://timestamp.digicert.com)
	TimestampServer string `yaml:"timestamp_server,omitempty" json:"timestamp_server,omitempty"`

	// SignAlgorithm is the signing algorithm (default: SHA256)
	SignAlgorithm string `yaml:"sign_algorithm,omitempty" json:"sign_algorithm,omitempty"`
}

// DefaultWindowsConfig returns the default Windows signing configuration
func DefaultWindowsConfig() *WindowsConfig {
	return &WindowsConfig{
		TimestampServer: "http://timestamp.digicert.com",
		SignAlgorithm:   "SHA256",
	}
}

// DefaultMacOSConfig returns the default macOS signing configuration
func DefaultMacOSConfig() *MacOSConfig {
	return &MacOSConfig{
		Identity:        "-", // Ad-hoc signing by default
		HardenedRuntime: true,
	}
}
