package manifest

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverManifestSymlinks(t *testing.T) {
	for _, dangling := range []bool{false, true} {
		name := "valid"
		if dangling {
			name = "dangling"
		}
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			target := filepath.Join(root, "native.yaml")
			if !dangling {
				if err := os.WriteFile(target, []byte(minimalYAML), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			path := filepath.Join(root, Filename)
			if err := os.Symlink(target, path); err != nil {
				t.Skipf("symlinks unavailable: %v", err)
			}
			discoveredRoot, discoveredPath, err := Discover(root)
			if dangling {
				if err == nil || errors.Is(err, ErrNotFound) {
					t.Fatalf("existing dangling manifest must fail discovery, got %v", err)
				}
				if !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("read error must retain its cause, got %v", err)
				}
				return
			}
			if err != nil || discoveredRoot != root || discoveredPath != path {
				t.Fatalf("Discover() = %q, %q, %v", discoveredRoot, discoveredPath, err)
			}
			if _, err := Load(root, ""); err != nil {
				t.Fatalf("load valid symlink: %v", err)
			}
		})
	}
}
