package pipeline

import (
	"fmt"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func validateDarwinSigningHost(spec SignSpec, host HostCapabilities) error {
	signing := spec.Config
	portable := signing.Certificate != "" || signing.Identity == "-"
	if !portable {
		if host.os != "darwin" {
			return hostConfigurationField("sign", fmt.Errorf("keychain signing requires a darwin host; use an exported certificate or identity '-' for portable signing"))
		}
		if signing.Identity == "" {
			return hostConfigurationField("identity", fmt.Errorf("darwin signing requires an identity or exported certificate"))
		}
		if err := requireHostTools(host, "darwin keychain signing", "codesign"); err != nil {
			return hostConfigurationField("sign", err)
		}
	}
	if signing.Certificate != "" && signing.Credential != "" && !host.hasCredential(signing.Credential) {
		return hostConfigurationField("credential", fmt.Errorf("Apple signing credential %q is unavailable", signing.Credential))
	}
	if spec.Format == "dmg" && !spec.BundleSigned {
		if host.os != "darwin" {
			return hostConfigurationField("sign", fmt.Errorf("signing an application after DMG-editing hooks requires macOS and hdiutil"))
		}
		if err := requireHostTools(host, "signing applications after DMG-editing hooks", "hdiutil"); err != nil {
			return hostConfigurationField("sign", err)
		}
	}
	if !signing.Notarize {
		return nil
	}
	if signing.Identity == "-" && signing.Certificate == "" {
		return hostConfigurationField("identity", fmt.Errorf("Apple notarization requires a Developer ID certificate; ad-hoc signing is insufficient"))
	}
	if err := manifest.ValidateNotarization(signing); err != nil {
		return err
	}
	if signing.NotarizationCredential == "" && signing.NotarizationKey == "" {
		return hostConfigurationField("notarization.credential", fmt.Errorf("notarization requires a keychain profile or team API key"))
	}
	if portable && signing.TimestampServer == "" {
		return hostConfigurationField("timestamp_server", fmt.Errorf("notarized portable signatures require an RFC 3161 timestamp server"))
	}
	if signing.NotarizationKey != "" {
		return nil
	}
	if host.os != "darwin" {
		return hostConfigurationField("sign", fmt.Errorf("notarytool keychain profiles require a darwin host; configure team API credentials for portable notarization"))
	}
	return hostConfigurationField("sign", requireHostTools(host, "darwin notarytool keychain profile", "xcrun"))
}
