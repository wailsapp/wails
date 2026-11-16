package appimage

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestNativeImageCompanionsAndModes(t *testing.T) {
	options := appImageFixture(t)
	if err := os.WriteFile(filepath.Join(options.Root, "runtime"), minimalELF(runtime.GOARCH, true), 0755); err != nil {
		t.Fatal(err)
	}
	options.RuntimePath = "runtime"
	if err := CreateNative(t.Context(), options); err != nil {
		t.Fatal(err)
	}
	if err := verify(filepath.Join(options.Root, options.Output), runtime.GOARCH); err != nil {
		t.Fatal(err)
	}
	reader, err := exec.LookPath("unsquashfs")
	if err != nil {
		t.Skip("official SquashFS verifier unavailable")
	}
	destination := filepath.Join(t.TempDir(), "image")
	command := exec.Command(reader, "-no-progress", "-o", "64", "-d", destination, filepath.Join(options.Root, options.Output))
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("official reader: %s, %v", output, err)
	}
	for _, name := range []string{"app", "helper"} {
		source, err := os.ReadFile(filepath.Join(options.Root, name))
		if err != nil {
			t.Fatal(err)
		}
		filename := filepath.Join(destination, "usr/bin", name)
		data, err := os.ReadFile(filename)
		if err != nil || string(data) != string(source) {
			t.Fatalf("companion %s changed", name)
		}
		if info, err := os.Stat(filename); err != nil || info.Mode().Perm() != 0755 {
			t.Fatalf("executable mode: %v, %v", info, err)
		}
	}
}

func TestNativeImageFailurePreservesPublishedImage(t *testing.T) {
	options := appImageFixture(t)
	options.RuntimePath = "runtime"
	options.UpdateInfo = "zsync|https://example.invalid/image.zsync"
	if err := os.WriteFile(filepath.Join(options.Root, "runtime"), minimalELF(runtime.GOARCH, true), 0755); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(options.Root, options.Output)
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("previous image"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := CreateNative(t.Context(), options); err == nil || !strings.Contains(err.Error(), "update information") {
		t.Fatalf("expected missing update section: %v", err)
	}
	data, err := os.ReadFile(target)
	if err != nil || string(data) != "previous image" {
		t.Fatalf("previous image lost: %q, %v", data, err)
	}
}

func TestNativeOptionsRejectUnsafeNamesAndInputOverwrite(t *testing.T) {
	options := appImageFixture(t)
	for _, name := range []string{"../app", "/app", "a\\b", ""} {
		candidate := options
		candidate.Binaries = append([]Binary(nil), options.Binaries...)
		candidate.Binaries[0].Name = name
		if _, err := normalizeNativeOptions(candidate); err == nil {
			t.Fatalf("accepted binary name %q", name)
		}
	}
	options.Output = "app"
	if _, err := normalizeNativeOptions(options); err == nil {
		t.Fatal("allowed overwriting source binary")
	}
}
