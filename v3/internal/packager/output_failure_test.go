package packager

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPackageFailurePreservesOutput(t *testing.T) {
	for _, format := range []PackageType{DEB, RPM, APK, IPK, ARCH} {
		t.Run(string(format), func(t *testing.T) {
			root := t.TempDir()
			config := filepath.Join(root, "nfpm.yaml")
			data := "name: example\nversion: 1.0.0\narch: amd64\nmaintainer: Example <example@example.com>\ncontents:\n- src: " + filepath.Join(root, "missing-file") + "\n  dst: /usr/bin/example\n"
			if err := os.WriteFile(config, []byte(data), 0600); err != nil {
				t.Fatal(err)
			}
			output := filepath.Join(root, "existing.package")
			if err := os.WriteFile(output, []byte("previous package"), 0600); err != nil {
				t.Fatal(err)
			}
			if err := CreatePackageFromConfig(format, config, output); err == nil {
				t.Fatal("expected missing-content failure")
			}
			got, err := os.ReadFile(output)
			if err != nil || string(got) != "previous package" {
				t.Errorf("previous output=%q, error=%v", got, err)
			}
			entries, err := os.ReadDir(root)
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 2 {
				t.Errorf("temporary output retained: %v", entries)
			}
			previous, err := os.Stat(output)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(root, "missing-file"), []byte("application binary"), 0600); err != nil {
				t.Fatal(err)
			}
			if err := CreatePackageFromConfig(format, config, output); err != nil {
				t.Fatal(err)
			}
			got, err = os.ReadFile(output)
			if err != nil || len(got) == 0 || string(got) == "previous package" {
				t.Fatalf("successful build did not replace previous output: %d bytes, error=%v", len(got), err)
			}
			current, err := os.Stat(output)
			if err != nil {
				t.Fatal(err)
			}
			if current.Mode().Perm() != previous.Mode().Perm() {
				t.Errorf("permissions=%v, want %v", current.Mode().Perm(), previous.Mode().Perm())
			}
			entries, err = os.ReadDir(root)
			if err != nil || len(entries) != 3 {
				t.Errorf("unexpected output after success: %v, error=%v", entries, err)
			}
		})
	}
}
