package pipeline

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/cache"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestVoidPackageTracksLifecycleScriptContents(t *testing.T) {
	for _, script := range []string{"pre-install.sh", "post-install.sh", "pre-remove.sh", "post-remove.sh"} {
		t.Run(script, func(t *testing.T) {
			root := t.TempDir()
			config := manifest.PackageFormat{PreInstall: "pre-install.sh", PostInstall: "post-install.sh", PreRemove: "pre-remove.sh", PostRemove: "post-remove.sh"}
			for _, name := range []string{config.PreInstall, config.PostInstall, config.PreRemove, config.PostRemove} {
				if err := os.WriteFile(filepath.Join(root, name), []byte("echo before\n"), 0644); err != nil {
					t.Fatal(err)
				}
			}
			node := Node{Kind: PackageArtifact, Spec: PackageSpec{TargetOS: "linux", Format: "void", Config: config}, Inputs: packageInputs(root, "linux", "void", config)}
			store, err := cache.OpenCache(root)
			if err != nil {
				t.Fatal(err)
			}
			action := func() string {
				t.Helper()
				inputs, err := snapshotNodeInputs(store, node)
				if err != nil {
					t.Fatal(err)
				}
				key, err := cache.ActionKey(string(node.Kind), nodeActionSpec(node, "void-writer"), inputs, nil)
				if err != nil {
					t.Fatal(err)
				}
				return key
			}
			before := action()
			if err := os.WriteFile(filepath.Join(root, script), []byte("echo after\n"), 0644); err != nil {
				t.Fatal(err)
			}
			store.InvalidateObservations()
			if before == action() {
				t.Fatalf("editing %s did not invalidate the Void package action", script)
			}
		})
	}
}
