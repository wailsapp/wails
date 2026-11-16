// Package signingconfig resolves user-level signing defaults for runtime plans.
package signingconfig

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/internal/defaults"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

// ApplyDefaults returns a runtime copy. Only signed selected-profile targets or
// explicitly requested signing platforms read global defaults. Never persist
// the returned configuration as project YAML.
func ApplyDefaults(config manifest.Config, requestedPlatforms ...string) (manifest.Config, error) {
	requested := make(map[string]bool)
	for _, platform := range requestedPlatforms {
		requested[platform] = true
	}
	for _, target := range config.Selected.Targets {
		if target.Sign {
			platform, _, _ := strings.Cut(target.Target, "/")
			requested[platform] = true
		}
	}
	if !requested["darwin"] && !requested["windows"] && !requested["linux"] {
		return config, nil
	}
	global, err := defaults.Load()
	if err != nil {
		return config, fmt.Errorf("load global signing defaults: %w", err)
	}
	return apply(config, requested, global.Signing)
}

func apply(config manifest.Config, requested map[string]bool, global defaults.SigningDefaults) (manifest.Config, error) {
	for _, platform := range []string{"darwin", "windows", "linux"} {
		if !requested[platform] {
			continue
		}
		var signing *manifest.SigningPlatform
		switch platform {
		case "darwin":
			signing = &config.Signing.Darwin
		case "windows":
			signing = &config.Signing.Windows
		case "linux":
			signing = &config.Signing.Linux
		}
		signing.Enabled = true
		fillPath := func(value *string, trusted *string, fallback string) error {
			if *value != "" || fallback == "" {
				return nil
			}
			path, err := filepath.Abs(fallback)
			if err != nil {
				return err
			}
			*value, *trusted = path, path
			return nil
		}
		switch platform {
		case "darwin":
			explicitIdentity := signing.Identity != "" || signing.Certificate != ""
			if !explicitIdentity {
				signing.Identity = global.Darwin.Identity
			}
			if !explicitIdentity {
				if err := fillPath(&signing.Certificate, &signing.GlobalPaths.Certificate, global.Darwin.P12Path); err != nil {
					return config, err
				}
			}
			if err := fillPath(&signing.Entitlements, &signing.GlobalPaths.Entitlements, global.Darwin.Entitlements); err != nil {
				return config, err
			}
			notarize := signing.Notarize
			for _, target := range config.Selected.Targets {
				notarize = notarize || (strings.HasPrefix(target.Target, "darwin/") && target.Notarize)
			}
			if notarize && signing.NotarizationCredential == "" && signing.NotarizationKey == "" && signing.NotarizationKeyID == "" && signing.NotarizationIssuerID == "" {
				signing.NotarizationCredential = global.Darwin.KeychainProfile
				if signing.NotarizationCredential == "" {
					if err := fillPath(&signing.NotarizationKey, &signing.GlobalPaths.NotarizationKey, global.Darwin.APIKeyPath); err != nil {
						return config, err
					}
					signing.NotarizationKeyID, signing.NotarizationIssuerID = global.Darwin.APIKeyID, global.Darwin.APIIssuerID
				}
			}
		case "windows":
			if signing.Certificate == "" && signing.Thumbprint == "" {
				if err := fillPath(&signing.Certificate, &signing.GlobalPaths.Certificate, global.Windows.CertificatePath); err != nil {
					return config, err
				}
				if signing.Certificate == "" {
					signing.Thumbprint = global.Windows.Thumbprint
				}
			}
			if signing.TimestampServer == "" {
				signing.TimestampServer = global.Windows.TimestampServer
			}
		case "linux":
			if err := fillPath(&signing.Certificate, &signing.GlobalPaths.Certificate, global.Linux.GPGKeyPath); err != nil {
				return config, err
			}
			if signing.Identity == "" {
				signing.Identity = global.Linux.SignRole
			}
		}
	}
	return config, nil
}
