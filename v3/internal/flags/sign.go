package flags

// Sign contains flags for the sign command
type Sign struct {
	Input   string `name:"input" description:"Path to the file to sign"`
	Output  string `name:"output" description:"Output path (optional, defaults to in-place signing)"`
	Verbose bool   `name:"verbose" description:"Enable verbose output"`

	// Windows certificate signing
	Certificate string `name:"certificate" description:"Path to PKCS#12 (.pfx/.p12) certificate file (Windows)"`
	Password    string `name:"password" description:"Certificate password (reads from keychain if not provided)"`
	Thumbprint  string `name:"thumbprint" description:"Certificate thumbprint in Windows certificate store"`
	Timestamp   string `name:"timestamp" description:"Timestamp server URL"`

	// macOS native signing (requires macOS)
	Identity        string `name:"identity" description:"Signing identity (e.g., 'Developer ID Application: ...')"`
	Entitlements    string `name:"entitlements" description:"Path to entitlements plist file"`
	HardenedRuntime bool   `name:"hardened-runtime" description:"Enable hardened runtime (default: true for notarization)"`
	KeychainProfile string `name:"keychain-profile" description:"Keychain profile for notarization credentials (native macOS)"`

	// macOS cross-platform signing via Quill (works on any OS)
	P12Certificate string `name:"p12" description:"Path to P12 certificate file for cross-platform macOS signing"`
	NotaryKey      string `name:"notary-key" description:"Path to Apple API key (.p8) for notarization"`
	NotaryKeyID    string `name:"notary-key-id" description:"Apple API Key ID for notarization"`
	NotaryIssuer   string `name:"notary-issuer" description:"Apple Team ID (Issuer) for notarization"`

	// Notarization (works with both native and Quill)
	Notarize bool `name:"notarize" description:"Submit for Apple notarization after signing"`

	// Linux PGP signing
	PGPKey      string `name:"pgp-key" description:"Path to PGP private key file"`
	PGPPassword string `name:"pgp-password" description:"PGP key password (reads from keychain if not provided)"`
	Role        string `name:"role" description:"DEB signing role (origin, maint, archive, builder)"`
}
