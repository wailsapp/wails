package commands

import (
	"runtime"
	"testing"
)

func TestArchToMSIX(t *testing.T) {
	tests := []struct {
		goarch string
		msix   string
	}{
		{"amd64", "x64"},
		{"386", "x86"},
		{"arm64", "arm64"},
		{"arm", "arm"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.goarch, func(t *testing.T) {
			got := archToMSIX(tt.goarch)
			if got != tt.msix {
				t.Errorf("archToMSIX(%q) = %q, want %q", tt.goarch, got, tt.msix)
			}
		})
	}
}

func TestBuildAssetsDefaultArchNotX64(t *testing.T) {
	if runtime.GOARCH == "amd64" {
		t.Skip("skipping on amd64 where x64 happens to be correct")
	}

	options := &BuildAssetsOptions{
		Name:               "TestApp",
		ProductName:        "Test App",
		ProductDescription: "Test",
		ProductVersion:     "1.0.0",
		ProductCompany:     "Test",
		Silent:             true,
	}

	GenerateBuildAssets(options)

	expected := archToMSIX(runtime.GOARCH)
	if options.ProcessorArchitecture != expected {
		t.Errorf("ProcessorArchitecture = %q, want %q (from runtime.GOARCH=%q)", options.ProcessorArchitecture, expected, runtime.GOARCH)
	}
}
