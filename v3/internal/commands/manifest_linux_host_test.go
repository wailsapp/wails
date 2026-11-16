package commands

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func linuxPackageFixture(t testing.TB) (*manifestHandler, pipeline.PackageSpec) {
	t.Helper()
	root := t.TempDir()
	if err := GenerateBuildAssets(&BuildAssetsOptions{Dir: filepath.Join(root, "assets"), Name: "example", BinaryName: "example", ProductName: "Example", ProductVersion: "1.0.0", Silent: true}); err != nil {
		t.Fatal(err)
	}
	for _, directory := range []string{"bin", "dist"} {
		if err := os.MkdirAll(filepath.Join(root, directory), 0700); err != nil {
			t.Fatal(err)
		}
	}
	for name, data := range map[string]string{"example": "primary payload", "helper": "companion payload"} {
		if err := os.WriteFile(filepath.Join(root, "bin", name), []byte(data), 0644); err != nil {
			t.Fatal(err)
		}
	}
	spec := packageTestSpec("linux", "amd64", "deb")
	spec.Assets = "assets"
	spec.Project.Kind = "cli"
	spec.Companions = []pipeline.NamedBinary{{Name: "helper", Path: "bin/helper", Kind: "cli"}}
	return &manifestHandler{root: root}, spec
}

func TestLinuxPackageBinariesExecutableWithWindowsSourceModes(t *testing.T) {
	handler, spec := linuxPackageFixture(t)
	if _, err := handler.packageLinux(spec); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(handler.root, spec.Output))
	if err != nil {
		t.Fatal(err)
	}
	payload := debTarMember(t, data, "data.tar.gz")
	for name, want := range map[string]string{"usr/local/bin/example": "primary payload", "usr/bin/helper": "companion payload"} {
		archive := tar.NewReader(bytes.NewReader(payload))
		found := false
		for {
			entry, err := archive.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatal(err)
			}
			if strings.TrimPrefix(entry.Name, "./") != name {
				continue
			}
			found = true
			if entry.Mode&0777 != 0755 {
				t.Errorf("%s mode=%#o, want 0755", name, entry.Mode)
			}
			contents, err := io.ReadAll(archive)
			if err != nil || string(contents) != want {
				t.Errorf("%s payload=%q err=%v", name, contents, err)
			}
			break
		}
		if !found {
			t.Errorf("missing %s", name)
		}
	}
}

// DEB's short-named ar members allow an independent stdlib-only payload check.
func debTarMember(t *testing.T, data []byte, member string) []byte {
	t.Helper()
	if !bytes.HasPrefix(data, []byte("!<arch>\n")) {
		t.Fatal("not a DEB ar archive")
	}
	for offset := 8; offset+60 <= len(data); {
		header := data[offset : offset+60]
		size, err := strconv.Atoi(strings.TrimSpace(string(header[48:58])))
		if err != nil || size < 0 || size > len(data)-offset-60 {
			t.Fatal("invalid ar member size")
		}
		name := strings.TrimRight(strings.TrimSpace(string(header[:16])), "/")
		if name == member {
			reader, err := gzip.NewReader(bytes.NewReader(data[offset+60 : offset+60+size]))
			if err != nil {
				t.Fatal(err)
			}
			defer reader.Close()
			contents, err := io.ReadAll(reader)
			if err != nil {
				t.Fatal(err)
			}
			return contents
		}
		offset += 60 + size + size%2
	}
	t.Fatalf("missing %s", member)
	return nil
}
