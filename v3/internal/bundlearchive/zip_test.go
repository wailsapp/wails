package bundlearchive

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestBundleArchive(t *testing.T) {
	root := filepath.Join(t.TempDir(), "Example.app")
	binary := filepath.Join(root, "Contents", "MacOS", "example")
	if err := os.MkdirAll(filepath.Dir(binary), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(binary, []byte("application payload"), 0755); err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" {
		if err := os.Symlink("MacOS/example", filepath.Join(root, "Contents", "current")); err != nil {
			t.Fatal(err)
		}
	}
	var output bytes.Buffer
	if err := Write(t.Context(), &output, root); err != nil {
		t.Fatal(err)
	}
	reader, err := zip.NewReader(bytes.NewReader(output.Bytes()), int64(output.Len()))
	if err != nil {
		t.Fatal(err)
	}
	wanted := map[string]string{"Example.app/Contents/MacOS/example": "application payload"}
	if runtime.GOOS != "windows" {
		wanted["Example.app/Contents/current"] = "MacOS/example"
	}
	for _, entry := range reader.File {
		if entry.FileInfo().IsDir() {
			continue
		}
		stream, err := entry.Open()
		if err != nil {
			t.Fatal(err)
		}
		content, err := io.ReadAll(stream)
		stream.Close()
		if err != nil || string(content) != wanted[entry.Name] {
			t.Fatalf("entry %s: %q, %v", entry.Name, content, err)
		}
		delete(wanted, entry.Name)
		if runtime.GOOS != "windows" {
			if entry.Name == "Example.app/Contents/current" && entry.Mode()&os.ModeSymlink == 0 {
				t.Fatal("symlink was flattened")
			}
			if entry.Name == "Example.app/Contents/MacOS/example" && entry.Mode().Perm() != 0755 {
				t.Fatalf("executable mode = %v", entry.Mode())
			}
		}
	}
	if len(wanted) != 0 {
		t.Fatalf("missing entries: %v", wanted)
	}
}

func TestArchiveRejectsEscapingSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation requires Windows privileges")
	}
	root := t.TempDir()
	if err := os.Symlink("../outside", filepath.Join(root, "link")); err != nil {
		t.Fatal(err)
	}
	if err := Write(t.Context(), io.Discard, root); err == nil {
		t.Fatal("escaping symlink accepted")
	}
}

func TestArchiveCancellationAndOutputFailure(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if err := Write(ctx, io.Discard, t.TempDir()); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled archive: %v", err)
	}
	broken := errors.New("output unavailable")
	if err := Write(t.Context(), failingWriter{broken}, t.TempDir()); !errors.Is(err, broken) {
		t.Fatalf("output failure: %v", err)
	}
}

type failingWriter struct{ err error }

func (w failingWriter) Write([]byte) (int, error) { return 0, w.err }

func BenchmarkBundleArchive(b *testing.B) {
	root := filepath.Join(b.TempDir(), "Example.app")
	if err := os.Mkdir(root, 0755); err != nil {
		b.Fatal(err)
	}
	data := make([]byte, 4<<20)
	var value uint32 = 1
	for index := 0; index < len(data)/2; index++ {
		value = value*1664525 + 1013904223
		data[index] = byte(value >> 24)
	}
	if err := os.WriteFile(filepath.Join(root, "binary"), data, 0755); err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(len(data)))
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if err := Write(b.Context(), io.Discard, root); err != nil {
			b.Fatal(err)
		}
	}
}
