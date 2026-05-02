package flags

// Sign contains flags for the sign command
type Sign struct {
	Input       string `name:"input" description:"Path to the file to sign"`
	Output      string `name:"output" description:"Output path (optional, defaults to in-place signing)"`
	Verbose     bool   `name:"verbose" description:"Enable verbose output"`

	// Windows/macOS certificate signing
	Certificate string `name:"certificate" description:"Path to PKCS#12 (.pfx/.p12) certificate file"`
	Password    string `name:"password" description:"Certificate password (reads from keychain if not provided)"`
	Thumbprint  string `name:"thumbprint" description:"Certificate thumbprint in Windows certificate store"`
	Timestamp   string `name:"timestamp" description:"Timestamp server URL"`

	// macOS specific
	Identity         string `name:"identity" description:"Signing identity (e.g., 'Developer ID Application: ...')"`
	Entitlements     string `name:"entitlements" description:"Path to entitlements plist file"`
	HardenedRuntime  bool   `name:"hardened-runtime" description:"Enable hardened runtime (default: true for notarization)"`
	Notarize         bool   `name:"notarize" description:"Submit for Apple notarization after signing"`
	KeychainProfile  string `name:"keychain-profile" description:"Keychain profile for notarization credentials"`

	// Linux PGP signing
	PGPKey      string `name:"pgp-key" description:"Path to PGP private key file"`
	PGPPassword string `name:"pgp-password" description:"PGP key password (reads from keychain if not provided)"`
	Role        string `name:"role" description:"DEB signing role (origin, maint, archive, builder)"`
}
