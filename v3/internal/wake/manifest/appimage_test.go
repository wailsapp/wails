package manifest

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestAppImageRuntimeOptionsRoundTrip(t *testing.T) {
	root := t.TempDir()
	source := `version: 3
project: {name: image, product_name: Image, identifier: com.example.image, version: 1.0.0}
packages:
  appimage:
    runtime_root: packaging/linux-root
    runtime_path: packaging/runtime
    gtk_version: 3
`
	if err := os.WriteFile(filepath.Join(root, Filename), []byte(source), 0644); err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(root, "")
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := EncodeConfig(loaded.Config)
	if err != nil {
		t.Fatal(err)
	}
	restored, err := decodeYAML(root, "round-trip.yaml", encoded, "")
	if err != nil || !reflect.DeepEqual(restored.Config.Package.Linux.AppImage, loaded.Config.Package.Linux.AppImage) {
		t.Fatal("runtime configuration did not round trip", err)
	}
	for _, invalid := range []string{
		strings.Replace(source, "gtk_version: 3", "gtk_version: 2", 1),
		strings.Replace(source, "packaging/linux-root", "../outside", 1),
		strings.Replace(source, "packaging/runtime", "/runtime", 1),
	} {
		if _, err := decodeYAML(root, "invalid.yaml", []byte(invalid), ""); err == nil {
			t.Fatal("invalid runtime configuration accepted")
		}
	}
}
