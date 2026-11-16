package voidpackage

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"howett.net/plist"
)

func fixture(t testing.TB) Options {
	t.Helper()
	root := t.TempDir()
	write := func(name, data string) string {
		p := filepath.Join(root, name)
		if err := os.WriteFile(p, []byte(data), 0755); err != nil {
			t.Fatal(err)
		}
		return p
	}
	return Options{Metadata: Metadata{Name: "wake-fixture", Version: "1.2.3", Arch: "amd64", Description: "Fixture package\nLong description", License: "MIT", Depends: []string{"dependency>=1.0_1"}}, MTime: time.Unix(1700000000, 0).UTC(), Contents: []Entry{
		{Source: write("app", "main payload"), Destination: "/usr/bin/app", Mode: 0755},
		{Source: write("settings", "configuration"), Destination: "/etc/app.conf", Type: "config", Mode: 0640},
		{Source: "../bin/app", Destination: "/usr/lib/app-link", Type: "symlink", Mode: 0777},
		{Destination: "/var/lib/app", Type: "dir", Mode: 0700},
	}}
}

type archiveEntry struct {
	header tar.Header
	data   []byte
}

func decode(t testing.TB, data []byte) map[string]archiveEntry {
	t.Helper()
	compressed, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	defer compressed.Close()
	archive := tar.NewReader(compressed)
	entries := map[string]archiveEntry{}
	for {
		header, err := archive.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(archive)
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := entries[header.Name]; ok {
			t.Fatalf("duplicate archive name %q", header.Name)
		}
		entries[header.Name] = archiveEntry{*header, body}
	}
	return entries
}

func TestArchiveMetadataHashesAndEntries(t *testing.T) {
	options := fixture(t)
	var output bytes.Buffer
	if err := Write(context.Background(), &output, options); err != nil {
		t.Fatal(err)
	}
	first := append([]byte(nil), output.Bytes()...)
	output.Reset()
	if err := Write(context.Background(), &output, options); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, output.Bytes()) {
		t.Fatal("identical inputs are not reproducible")
	}
	entries := decode(t, output.Bytes())
	if string(entries["./usr/bin/app"].data) != "main payload" || entries["./usr/bin/app"].header.Mode != 0755 || entries["./etc/app.conf"].header.Mode != 0640 || entries["./var/lib/app/"].header.Mode != 0700 {
		t.Fatal("payload/modes changed")
	}
	if entries["./usr/lib/app-link"].header.Linkname != "../bin/app" {
		t.Fatal("link changed")
	}
	var props map[string]any
	if _, err := plist.Unmarshal(entries["./props.plist"].data, &props); err != nil {
		t.Fatal(err)
	}
	if props["pkgver"] != "wake-fixture-1.2.3_1" || props["architecture"] != "x86_64" || props["installed_size"] != uint64(25) {
		t.Fatalf("bad properties: %#v", props)
	}
	var files map[string][]map[string]any
	if _, err := plist.Unmarshal(entries["./files.plist"].data, &files); err != nil {
		t.Fatal(err)
	}
	for _, kind := range []string{"files", "conf_files"} {
		for _, item := range files[kind] {
			sum := sha256.Sum256(entries["."+item["file"].(string)].data)
			if item["sha256"] != hex.EncodeToString(sum[:]) {
				t.Fatalf("bad manifest hash: %#v", item)
			}
		}
	}
	if len(files["files"]) != 1 || len(files["conf_files"]) != 1 || files["links"][0]["target"] != "/usr/bin/app" {
		t.Fatalf("bad manifest: %#v", files)
	}
}

func TestFileNameAndArchitecture(t *testing.T) {
	for _, spec := range []struct{ arch, want string }{{"amd64", "x86_64"}, {"arm64", "aarch64"}, {"386", "i686"}, {"aarch64-musl", "aarch64-musl"}, {"noarch", "noarch"}} {
		name, err := FileName(Metadata{Name: "example-app", Version: "1.2.3beta.4", Revision: "2", Arch: spec.arch})
		if err != nil || name != "example-app-1.2.3beta.4_2."+spec.want+".xbps" {
			t.Fatalf("%s: %s %v", spec.arch, name, err)
		}
	}
	for _, metadata := range []Metadata{{Name: "../app", Version: "1", Arch: "amd64"}, {Name: "app", Version: "1_2", Arch: "amd64"}, {Name: "app", Version: "1", Arch: "arm"}, {Name: "app", Version: "1", Arch: "amd64", Revision: "0"}} {
		if _, err := FileName(metadata); err == nil {
			t.Fatalf("accepted invalid identity: %#v", metadata)
		}
	}
}

func TestInvalidContentIsRejectedBeforeWriting(t *testing.T) {
	for _, entries := range [][]Entry{
		{{Destination: "/../escape", Type: "dir"}},
		{{Destination: "relative", Type: "dir"}},
		{{Destination: "/props.plist", Type: "dir"}},
		{{Destination: "/app", Type: "dir"}, {Destination: "/app", Type: "dir"}},
		{{Destination: "/link", Type: "symlink", Source: "/elsewhere"}, {Destination: "/link/nested", Type: "dir"}},
		{{Destination: "/unsupported", Type: "ghost"}},
	} {
		options := fixture(t)
		options.Contents = entries
		var output bytes.Buffer
		if err := Write(context.Background(), &output, options); err == nil || output.Len() != 0 {
			t.Fatalf("invalid entries wrote archive: %v", err)
		}
	}
}

type interceptWriter struct {
	bytes.Buffer
	intercept func()
}

func (w *interceptWriter) Write(data []byte) (int, error) {
	if w.intercept != nil {
		f := w.intercept
		w.intercept = nil
		f()
	}
	return w.Buffer.Write(data)
}

func TestChangingPayloadCannotPublishFalseHash(t *testing.T) {
	options := fixture(t)
	output := &interceptWriter{intercept: func() {
		if err := os.WriteFile(options.Contents[0].Source, []byte("evil payload"), 0755); err != nil {
			t.Fatal(err)
		}
	}}
	if err := Write(context.Background(), output, options); err == nil {
		t.Fatal("changed payload accepted")
	}
}

func TestCancellationDuringArchiveWrite(t *testing.T) {
	options := fixture(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	output := &interceptWriter{intercept: cancel}
	if err := Write(ctx, output, options); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}

func BenchmarkWrite(b *testing.B) {
	options := fixture(b)
	if err := os.WriteFile(options.Contents[0].Source, bytes.Repeat([]byte("payload-"), 256*1024), 0755); err != nil {
		b.Fatal(err)
	}
	b.SetBytes(2 * 1024 * 1024)
	b.ReportAllocs()
	for b.Loop() {
		if err := Write(context.Background(), io.Discard, options); err != nil {
			b.Fatal(err)
		}
	}
}
