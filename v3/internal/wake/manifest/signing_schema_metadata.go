package manifest

import "strings"

type signingFieldMetadata struct{ description, example string }

var signingFieldDescriptions = map[string]map[string]signingFieldMetadata{
	"certificate": {
		"windows": {"Project-relative PFX signing certificate file; alternatively configure thumbprint.", `"signing/windows.pfx"`},
		"linux":   {"Arch: project-relative exported armored PGP private-key file. DEB/RPM: PGP key identifier passed to the external signing tool.", `"signing/private-key.asc"`},
		"android": {"Project-relative Android signing keystore file.", `"signing/upload.jks"`},
	},
	"identity": {
		"darwin":  {"macOS code-signing identity from the keychain; use - for ad-hoc signing.", `"Developer ID Application: Example"`},
		"ios":     {"Apple code-signing identity from the macOS keychain.", `"Apple Development: Example"`},
		"linux":   {"DEB signing role; defaults to builder. Not used by the Arch detached signer.", `"builder"`},
		"android": {"Android signing key alias when key_alias is omitted.", `"upload"`},
	},
	"credential": {
		"windows": {"Name of the environment variable containing the PFX certificate password; never the password itself.", `"WINDOWS_SIGNING_PASSWORD"`},
		"linux":   {"Arch: optional environment variable containing the private-key passphrase; omit for an unencrypted key.", `"ARCH_SIGNING_PASSPHRASE"`},
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
		"windows": {"Thumbprint of a certificate in the Windows certificate store.", `"0123456789ABCDEF"`},
	},
	"timestamp_server": {
		"windows": {"URL of the Windows signing timestamp service.", `"https://timestamp.example.com"`},
	},
	"key_alias": {
		"android": {"Alias of the Android signing key.", `"upload"`},
	},
}

func platformSigningMetadata(path, name string) (signingFieldMetadata, bool) {
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
