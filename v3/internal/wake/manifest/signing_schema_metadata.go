package manifest

import "strings"

type signingFieldMetadata struct{ description, example string }

var signingFieldDescriptions = map[string]map[string]signingFieldMetadata{
	"certificate": {
		"darwin":  {"Project-relative exported PKCS#12 identity for portable Go signing; omit to use a macOS keychain identity or ad-hoc signing.", `"signing/developer-id.p12"`},
		"windows": {"Project-relative exported RSA PFX identity for portable PE/MSIX signing on any desktop host. AES/PBKDF2 and 3DES supported; native Windows SignTool handles RC2 exports. Alternatively configure thumbprint.", `"signing/windows.pfx"`},
		"linux":   {"Exported armored or binary OpenPGP private-key file for native DEB, RPM or Arch signing. Bare GnuPG key identifiers retain the external DEB/RPM signer.", `"signing/private-key.asc"`},
		"android": {"Project-relative Android signing keystore file.", `"signing/upload.jks"`},
	},
	"identity": {
		"darwin":  {"macOS code-signing identity from the keychain; use - for portable ad-hoc signing, or configure an exported certificate.", `"Developer ID Application: Example"`},
		"ios":     {"Apple code-signing identity from the macOS keychain.", `"Apple Development: Example"`},
		"linux":   {"DEB signing role; defaults to builder. Not used by the Arch detached signer.", `"builder"`},
		"android": {"Android signing key alias when key_alias is omitted.", `"upload"`},
	},
	"credential": {
		"darwin":  {"Name of the environment variable containing the exported PKCS#12 identity password; never the password itself.", `"APPLE_SIGNING_PASSWORD"`},
		"windows": {"Name of the environment variable containing the PFX certificate password; never the password itself.", `"WINDOWS_SIGNING_PASSWORD"`},
		"linux":   {"Optional environment variable containing the exported private-key passphrase; omit for an unencrypted key.", `"LINUX_SIGNING_PASSPHRASE"`},
		"android": {"Name of the environment variable containing the keystore password; never the password itself.", `"ANDROID_SIGNING_PASSWORD"`},
	},
	"entitlements": {
		"darwin": {"Project-relative Apple entitlements file.", `"signing/entitlements.plist"`},
		"ios":    {"Project-relative Apple entitlements file.", `"signing/entitlements.plist"`},
	},
	"provisioning_profile": {
		"ios": {"Project-relative Apple provisioning profile.", `"signing/app.mobileprovision"`},
	},
	"thumbprint": {
		"windows": {"Thumbprint of a certificate in the Windows certificate store; requires Windows with SignTool.", `"0123456789ABCDEF"`},
	},
	"timestamp_server": {
		"darwin":  {"RFC 3161 timestamp URL for portable certificate signing; required for notarized certificate signatures.", `"https://timestamp.example.com"`},
		"windows": {"RFC 3161 timestamp service URL; defaults to http://timestamp.digicert.com. Timestamp signatures and trust chains are verified.", `"https://timestamp.example.com"`},
	},
	"key_alias": {
		"android": {"Alias of the Android signing key.", `"upload"`},
	},
}

func platformSigningMetadata(path, name string) (signingFieldMetadata, bool) {
	switch path {
	case "darwin.notarization.api_key":
		return signingFieldMetadata{"Project-relative App Store Connect team API private-key .p8 file. Uses the portable Notary API instead of notarytool; requires key_id and issuer_id.", `"signing/AuthKey.p8"`}, true
	case "darwin.notarization.key_id":
		return signingFieldMetadata{"App Store Connect team API key ID belonging to api_key.", `"ABCDEFGHIJ"`}, true
	case "darwin.notarization.issuer_id":
		return signingFieldMetadata{"App Store Connect team issuer ID for API notarization.", `"00000000-0000-0000-0000-000000000000"`}, true
	}
	if path == "darwin.notarization.credential" {
		return signingFieldMetadata{"Name of the Apple notarytool keychain profile used for notarisation.", `"release-notary"`}, true
	}
	if !strings.Contains(path, ".signing.") {
		return signingFieldMetadata{}, false
	}
	platforms, ok := signingFieldDescriptions[name]
	if !ok {
		return signingFieldMetadata{}, false
	}
	platform, _, _ := strings.Cut(path, ".")
	if metadata, ok := platforms[platform]; ok {
		return metadata, true
	}
	return signingFieldMetadata{"Not used by this platform's signing operation.", `""`}, true
}
