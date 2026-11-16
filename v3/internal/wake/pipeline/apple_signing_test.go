package pipeline

import (
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestPortableAppleSigningHostMatrix(t *testing.T) {
	for _, host := range []string{"linux", "darwin", "windows"} {
		for _, signing := range []manifest.SigningPlatform{{Enabled: true, Identity: "-"}, {Enabled: true, Certificate: "developer.p12"}} {
			for _, format := range []string{"app", "binary", "dmg"} {
				spec := SignSpec{TargetOS: "darwin", Format: format, BundleSigned: true, Config: signing}
				if err := validateSigningHost(spec, testHost(host, "amd64")); err != nil {
					t.Fatalf("%s/%s: %v", host, format, err)
				}
				if format == "dmg" {
					spec.BundleSigned = false
					if err := validateSigningHost(spec, testHost(host, "amd64")); err == nil {
						t.Fatal("image-editing hooks bypassed mount requirement", host)
					}
				}
			}
		}
	}
}

func TestPortableAppleSigningRejectsUnavailableCredentialsAndAdHocNotarization(t *testing.T) {
	spec := SignSpec{TargetOS: "darwin", Format: "app", Config: manifest.SigningPlatform{Enabled: true, Certificate: "identity.pfx", Credential: "APPLE_PASSWORD"}}
	if err := validateSigningHost(spec, testHost("linux", "amd64")); err == nil || !strings.Contains(err.Error(), "APPLE_PASSWORD") {
		t.Fatal("missing PFX password accepted", err)
	}
	spec.Config = manifest.SigningPlatform{Enabled: true, Identity: "-", Notarize: true, NotarizationCredential: "profile"}
	if err := validateSigningHost(spec, testHost("darwin", "arm64", "xcrun")); err == nil || !strings.Contains(err.Error(), "Developer ID") {
		t.Fatal("ad-hoc notarization accepted", err)
	}
}

func TestAppleAPINotarizationNeedsNoHostTools(t *testing.T) {
	for _, hostOS := range []string{"linux", "windows", "darwin"} {
		for _, format := range []string{"app", "dmg"} {
			t.Run(hostOS+"/"+format, func(t *testing.T) {
				host := NewHostCapabilities(hostOS, "arm64", nil, []string{"APPLE_PFX_PASSWORD"})
				spec := SignSpec{TargetOS: "darwin", Format: format, BundleSigned: true, Config: manifest.SigningPlatform{Enabled: true, Certificate: "developer.p12", Credential: "APPLE_PFX_PASSWORD", TimestampServer: "https://timestamp.example.test", Notarize: true, NotarizationKey: "AuthKey.p8", NotarizationKeyID: "KEYID", NotarizationIssuerID: "ISSUER"}}
				if err := validateSigningHost(spec, host); err != nil {
					t.Fatalf("portable API path required host tooling: %v", err)
				}
				spec.Config.TimestampServer = ""
				if err := validateSigningHost(spec, host); err == nil || !strings.Contains(err.Error(), "timestamp") {
					t.Fatalf("missing secure timestamp accepted: %v", err)
				}
			})
		}
	}
}

func TestAppleKeychainNotarizationStillRequiresMacAndXcrun(t *testing.T) {
	spec := SignSpec{TargetOS: "darwin", Format: "app", Config: manifest.SigningPlatform{Enabled: true, Certificate: "developer.p12", TimestampServer: "https://timestamp.example.test", Notarize: true, NotarizationCredential: "profile"}}
	for _, hostOS := range []string{"linux", "windows"} {
		if err := validateSigningHost(spec, NewHostCapabilities(hostOS, "arm64", []string{"xcrun"}, nil)); err == nil || !strings.Contains(err.Error(), "darwin") {
			t.Fatalf("%s accepted keychain profile: %v", hostOS, err)
		}
	}
	if err := validateSigningHost(spec, NewHostCapabilities("darwin", "arm64", nil, nil)); err == nil || !strings.Contains(err.Error(), "xcrun") {
		t.Fatalf("missing xcrun accepted: %v", err)
	}
	if err := validateSigningHost(spec, NewHostCapabilities("darwin", "arm64", []string{"xcrun"}, nil)); err != nil {
		t.Fatal(err)
	}
}
