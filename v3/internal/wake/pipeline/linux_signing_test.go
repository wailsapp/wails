package pipeline

import (
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestExportedLinuxSigningKeysDoNotRequireHostTools(t *testing.T) {
	for _, host := range []string{"linux", "darwin", "windows"} {
		for _, format := range []string{"deb", "rpm"} {
			config := testConfig(t)
			config.Project.Kind = "cli"
			config.Signing.Linux = manifest.SigningPlatform{Enabled: true, Certificate: "keys/release.asc", Credential: "KEY_PASSWORD"}
			plan, err := PlanBuildForHost(config, Request{Verb: "sign", TargetOS: "linux", TargetArch: "arm64", Formats: []string{format}, CGOEnabled: "0", SkipFrontend: true, SkipBindings: true}, NewHostCapabilities(host, "amd64", []string{"go"}, []string{"KEY_PASSWORD"}))
			if err != nil {
				t.Fatalf("%s -> %s: %v", host, format, err)
			}
			for _, node := range plan.Nodes {
				if node.Kind == SignArtifact && (node.Cache != CacheNever || len(node.Inputs) != 0) {
					t.Fatal("private key or signed output entered reusable cache")
				}
			}
		}
	}
}
