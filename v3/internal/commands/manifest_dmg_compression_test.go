package commands

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestDMGSignsPrivateBundleBeforePortableWriter(t *testing.T) {
	root := t.TempDir()
	spec := packageTestSpec("darwin", "arm64", "dmg")
	spec.Binary, spec.Output = "input.app", "dist/example.dmg"
	spec.BundleSigning = &manifest.SigningPlatform{Enabled: true, Identity: "Developer ID", Notarize: true}
	if err := os.Mkdir(filepath.Join(root, spec.Binary), 0755); err != nil {
		t.Fatal(err)
	}
	previousSign, previousPackage, previousPrepare := manifestSign, manifestToolPackage, manifestPrepareDMGForSigning
	t.Cleanup(func() {
		manifestSign, manifestToolPackage, manifestPrepareDMGForSigning = previousSign, previousPackage, previousPrepare
	})
	var calls []string
	manifestSign = func(options *flags.Sign) error {
		if filepath.Ext(options.Input) == ".app" {
			calls = append(calls, "bundle")
			if options.Input == filepath.Join(root, spec.Binary) || !options.HardenedRuntime || options.Notarize {
				t.Fatal("bundle signing must be private, hardened, and defer notarization")
			}
			return os.WriteFile(filepath.Join(options.Input, "signed"), []byte("signature"), 0600)
		}
		calls = append(calls, "image")
		return nil
	}
	manifestToolPackage = func(options *flags.ToolPackage) error {
		if _, err := os.Stat(filepath.Join(options.Out, spec.Project.BinaryName+".app", "signed")); err != nil {
			t.Fatal("image writer received unsigned bundle", err)
		}
		calls = append(calls, "compress")
		return os.WriteFile(filepath.Join(options.Out, spec.Project.BinaryName+".dmg"), []byte("image"), 0600)
	}
	manifestPrepareDMGForSigning = func(context.Context, string, string, string, bool) error {
		t.Fatal("pre-signed bundle must not mount or recompress image")
		return nil
	}
	handler := &manifestHandler{root: root}
	if _, err := handler.packageArtifact(t.Context(), spec); err != nil {
		t.Fatal(err)
	}
	if _, err := handler.sign(t.Context(), pipeline.SignSpec{BundleSigned: true, TargetOS: "darwin", Format: "dmg", Input: spec.Output, Config: *spec.BundleSigning}); err != nil {
		t.Fatal(err)
	}
	if len(calls) != 3 || calls[0] != "bundle" || calls[1] != "compress" || calls[2] != "image" {
		t.Fatalf("signing/packaging sequence: %v", calls)
	}
	if _, err := os.Stat(filepath.Join(root, spec.Binary, "signed")); !os.IsNotExist(err) {
		t.Fatal("source bundle was changed", err)
	}
}
