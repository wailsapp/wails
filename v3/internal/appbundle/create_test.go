package appbundle

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func fixture(t testing.TB) Options {
	t.Helper()
	root := t.TempDir()
	paths := map[string]string{}
	for _, name := range []string{"Info.plist", "icons.icns", "Assets.car", "app", "helper"} {
		path := filepath.Join(root, name)
		if err := os.WriteFile(path, []byte(name), 0755); err != nil {
			t.Fatal(err)
		}
		paths[name] = path
	}
	return Options{Output: filepath.Join(root, "out", "My App.app"), InfoPlist: paths["Info.plist"], Icon: paths["icons.icns"], Assets: paths["Assets.car"], Binaries: map[string]string{"app": paths["app"], "helper": paths["helper"]}}
}

func TestCreateReplacesFreshBundleAndPreservesModes(t *testing.T) {
	options := fixture(t)
	stale := filepath.Join(options.Output, "Contents", "stale")
	if err := os.MkdirAll(filepath.Dir(stale), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(stale, []byte("old"), 0600); err != nil {
		t.Fatal(err)
	}
	prepared := false
	options.Prepare = func(ctx context.Context, path string) error {
		if path == options.Output {
			t.Fatal("prepare received live output")
		}
		if _, err := os.Stat(stale); err != nil {
			t.Fatal("previous bundle removed before preparation")
		}
		prepared = true
		return os.WriteFile(filepath.Join(path, "signature"), []byte("signed"), 0600)
	}
	if err := Create(context.Background(), options); err != nil {
		t.Fatal(err)
	}
	if !prepared {
		t.Fatal("prepare was skipped")
	}
	if _, err := os.Stat(stale); !os.IsNotExist(err) {
		t.Fatalf("stale bundle entry remains: %v", err)
	}
	for name, want := range map[string]string{"Contents/Info.plist": "Info.plist", "Contents/Resources/icons.icns": "icons.icns", "Contents/Resources/Assets.car": "Assets.car", "Contents/MacOS/app": "app", "Contents/MacOS/helper": "helper", "signature": "signed"} {
		data, err := os.ReadFile(filepath.Join(options.Output, name))
		if err != nil || string(data) != want {
			t.Fatalf("%s: %q, %v", name, data, err)
		}
	}
	info, err := os.Stat(filepath.Join(options.Output, "Contents", "MacOS", "app"))
	if err != nil {
		t.Fatal(err)
	}
	sourceInfo, err := os.Stat(options.Binaries["app"])
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != sourceInfo.Mode().Perm() {
		t.Fatalf("executable mode = %o", info.Mode().Perm())
	}
}

func TestFailurePreservesExistingBundle(t *testing.T) {
	for _, failure := range []string{"missing source", "signing", "cancellation", "invalid name"} {
		t.Run(failure, func(t *testing.T) {
			options := fixture(t)
			if err := Create(context.Background(), options); err != nil {
				t.Fatal(err)
			}
			previous := filepath.Join(options.Output, "previous")
			if err := os.WriteFile(previous, []byte("keep"), 0600); err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			switch failure {
			case "missing source":
				options.InfoPlist += ".missing"
			case "signing":
				options.Prepare = func(context.Context, string) error { return errors.New("signing failed") }
			case "cancellation":
				options.Prepare = func(context.Context, string) error { cancel(); return nil }
			case "invalid name":
				options.Binaries["../escape"] = "invalid"
			}
			if err := Create(ctx, options); err == nil {
				t.Fatal("expected failure")
			}
			data, err := os.ReadFile(previous)
			if err != nil || string(data) != "keep" {
				t.Fatalf("old bundle lost: %q, %v", data, err)
			}
			matches, err := filepath.Glob(filepath.Join(filepath.Dir(options.Output), ".wails-app-*"))
			if err != nil || len(matches) != 0 {
				t.Fatalf("staging leak: %v, %v", matches, err)
			}
		})
	}
}

func BenchmarkCreateAppBundle(b *testing.B) {
	options := fixture(b)
	executable := options.Binaries["app"]
	if err := os.WriteFile(executable, bytes.Repeat([]byte("x"), 4<<20), 0755); err != nil {
		b.Fatal(err)
	}
	b.SetBytes(4 << 20)
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if err := Create(context.Background(), options); err != nil {
			b.Fatal(err)
		}
	}
}

func TestReplaceRollsBackPublicationFailure(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "app.app")
	if err := os.Mkdir(output, 0755); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(output, "keep")
	if err := os.WriteFile(marker, []byte("old"), 0600); err != nil {
		t.Fatal(err)
	}
	backup := filepath.Join(root, "previous.app")
	if err := replace(filepath.Join(root, "missing.app"), output, backup); err == nil {
		t.Fatal("expected missing staged bundle failure")
	}
	if data, err := os.ReadFile(marker); err != nil || string(data) != "old" {
		t.Fatalf("rollback lost previous bundle: %q, %v", data, err)
	}
	if _, err := os.Stat(backup); !os.IsNotExist(err) {
		t.Fatalf("rollback left a backup: %v", err)
	}
}
