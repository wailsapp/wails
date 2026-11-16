package commands

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestVoidAdapterPackagesExecutablesAndDesktopAssets(t *testing.T) {
	root := t.TempDir()
	spec := packageTestSpec("linux", "arm64", "void")
	spec.Project.Description = "Void package adapter fixture"
	spec.Config.Dependencies = []string{"gtk4"}
	spec.Companions = []pipeline.NamedBinary{{Name: "helper", Path: "bin/helper"}}
	inputs := map[string]string{"bin/example": "primary", "bin/helper": "companion", ".wails/assets/appicon.png": "icon", ".wails/assets/linux/example.desktop": "desktop"}
	for name, value := range inputs {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(value), 0755); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := (&manifestHandler{root: root}).packageVoid(t.Context(), spec); err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(filepath.Join(root, spec.Output))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	compressed, err := gzip.NewReader(file)
	if err != nil {
		t.Fatal(err)
	}
	defer compressed.Close()
	archive := tar.NewReader(compressed)
	want := map[string]string{"./usr/bin/example": "primary", "./usr/bin/helper": "companion", "./usr/share/icons/hicolor/256x256/apps/example.png": "icon", "./usr/share/applications/example.desktop": "desktop"}
	for {
		header, err := archive.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if expected, ok := want[header.Name]; ok {
			data, err := io.ReadAll(archive)
			if err != nil || string(data) != expected {
				t.Fatalf("%s: %q, %v", header.Name, data, err)
			}
			delete(want, header.Name)
		}
	}
	if len(want) != 0 {
		t.Fatalf("missing package contents: %v", want)
	}
}
