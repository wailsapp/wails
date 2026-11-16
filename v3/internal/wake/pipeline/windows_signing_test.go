package pipeline

import (
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"testing"
)

func TestPortableWindowsSigningHostMatrix(t *testing.T) {
	for _, hostOS := range []string{"linux", "darwin", "windows"} {
		for _, format := range []string{"binary", "exe", "dll", "nsis", "msix", "appx"} {
			spec := SignSpec{TargetOS: "windows", Format: format, Config: manifest.SigningPlatform{Enabled: true, Certificate: "app.pfx"}}
			host := NewHostCapabilities(hostOS, "amd64", nil, nil)
			if err := validateSigningHost(spec, host); err != nil {
				t.Errorf("%s/%s: %v", hostOS, format, err)
			}
			spec.Config.Credential = "PFX_PASSWORD"
			if err := validateSigningHost(spec, host); err == nil {
				t.Errorf("%s/%s accepted missing password", hostOS, format)
			}
			host = NewHostCapabilities(hostOS, "amd64", nil, []string{"PFX_PASSWORD"})
			if err := validateSigningHost(spec, host); err != nil {
				t.Errorf("%s/%s credential: %v", hostOS, format, err)
			}
		}
	}
	for _, format := range []string{"binary", "msix", "msi"} {
		spec := SignSpec{TargetOS: "windows", Format: format, Config: manifest.SigningPlatform{Enabled: true, Thumbprint: "123"}}
		for _, hostOS := range []string{"linux", "darwin", "windows"} {
			if err := validateSigningHost(spec, NewHostCapabilities(hostOS, "amd64", nil, nil)); err == nil {
				t.Errorf("%s/%s accepted missing native signer", hostOS, format)
			}
		}
		if err := validateSigningHost(spec, NewHostCapabilities("windows", "amd64", []string{"signtool.exe"}, nil)); err != nil {
			t.Fatal(err)
		}
	}
	spec := SignSpec{TargetOS: "windows", Format: "msi", Config: manifest.SigningPlatform{Enabled: true, Certificate: "app.pfx"}}
	if err := validateSigningHost(spec, NewHostCapabilities("linux", "amd64", nil, nil)); err == nil {
		t.Fatal("MSI accepted off Windows")
	}
}
