package pipeline

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestDMGSignsBeforeCompressionUnlessHooksNeedUnsignedImage(t *testing.T) {
	for _, phase := range []manifest.HookPhase{"", manifest.AfterPackage, manifest.BeforeSign} {
		t.Run(string(phase), func(t *testing.T) {
			config := testConfig(t)
			config.Signing.Darwin = manifest.SigningPlatform{Enabled: true, Identity: "Developer ID"}
			config.Selected = manifest.Profile{Name: "release", Targets: []manifest.ProfileTarget{{Target: "darwin/arm64", Formats: []string{"dmg"}, Sign: true}}}
			if phase != "" {
				if err := os.WriteFile(filepath.Join(config.Root, "hook.sh"), []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
					t.Fatal(err)
				}
				config.Hooks = map[manifest.HookPhase]manifest.Hook{phase: {Script: "hook.sh"}}
			}
			plan, err := PlanBuild(config, Request{Verb: "build"})
			if err != nil {
				t.Fatal(err)
			}
			packaging := plan.Nodes["package:darwin/arm64:dmg"]
			spec := packaging.Spec.(PackageSpec)
			wantEarly := phase == ""
			if (spec.BundleSigning != nil) != wantEarly {
				t.Fatalf("bundle signing with hook %q: %v", phase, spec.BundleSigning)
			}
			if wantEarly && (packaging.Cache != CacheNever || packaging.Claims.Exclusive != "sign") {
				t.Fatal("signed DMG needs uncached, exclusive signing")
			}
			signer := plan.Nodes["package:darwin/arm64:dmg:sign"].Spec.(SignSpec)
			if signer.BundleSigned != wantEarly {
				t.Fatal("outer signer would repeat or omit inner bundle signing")
			}
		})
	}
}
