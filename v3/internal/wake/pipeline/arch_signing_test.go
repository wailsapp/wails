package pipeline

import (
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestArchSigningPublishesPackageAndDetachedSignatureOnEveryHost(t *testing.T) {
	for _, hostOS := range []string{"linux", "darwin", "windows"} {
		t.Run(hostOS, func(t *testing.T) {
			config := testConfig(t)
			config.Project.Kind = "cli"
			config.Signing.Linux = manifest.SigningPlatform{Enabled: true, Certificate: "private-key.asc"}
			plan, err := PlanBuildForHost(config, Request{Verb: "sign", TargetOS: "linux", TargetArch: "arm64", Formats: []string{"archlinux"}, CGOEnabled: "0", SkipFrontend: true, SkipBindings: true}, NewHostCapabilities(hostOS, "amd64", []string{"go"}, nil))
			if err != nil {
				t.Fatal(err)
			}
			signature := plan.Nodes["package:linux/arm64:archlinux:sign"]
			if signature.Cache != CacheNever || signature.Artifact.Format != "sig" || !strings.HasSuffix(signature.Output, ".pkg.tar.zst.sig") || len(signature.Inputs) != 0 {
				t.Fatalf("incorrect signing node: %#v", signature)
			}
			var publication Node
			for _, key := range plan.Artifacts {
				if node := plan.Nodes[key]; node.Kind == PublishArtifact && node.Artifact.Format == "archlinux" {
					publication = node
				}
			}
			spec, ok := publication.Spec.(PublishSpec)
			if !ok || publication.Cache != CacheNever || len(spec.Associated) != 1 || spec.Associated[0].Destination != spec.Destination+".sig" || strings.Contains(spec.Destination, ".signed") {
				t.Fatalf("incorrect signed publication: %#v", publication)
			}
			collect := plan.Nodes["collect:artifacts"]
			references := collect.Spec.(CollectSpec).Artifacts
			if collect.Cache != CacheNever || len(references) != 2 || references[0].Path == references[1].Path {
				t.Fatalf("package/signature receipt: %#v", collect)
			}
		})
	}
}
