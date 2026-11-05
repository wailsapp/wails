package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestMSIXCopyFile(t *testing.T) {
	dir := t.TempDir()
	source, target := filepath.Join(dir, "source"), filepath.Join(dir, "target")
	content := bytes.Repeat([]byte("application binary\x00"), 4096)
	if err := os.WriteFile(source, content, 0600); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"new", "existing", "same file", "symlink", "hardlink"} {
		t.Run(name, func(t *testing.T) {
			destination := target
			if name == "existing" {
				if err := os.WriteFile(target, bytes.Repeat([]byte("old"), len(content)), 0600); err != nil {
					t.Fatal(err)
				}
			}
			if name == "same file" {
				destination = source
			}
			if name == "symlink" || name == "hardlink" {
				destination = filepath.Join(dir, name)
				link := os.Link
				if name == "symlink" {
					link = os.Symlink
				}
				if err := link(source, destination); err != nil {
					t.Skipf("cannot create %s fixture: %v", name, err)
				}
			}
			if err := copyFile(source, destination, 0644); err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(destination)
			if err != nil || !bytes.Equal(got, content) {
				t.Fatalf("copy length=%d, error=%v", len(got), err)
			}
		})
	}
}

func TestMSIXCopyRejectsDirectorySource(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(t.TempDir(), "target")
	if err := os.WriteFile(target, []byte("existing"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := copyFile(dir, target, 0644); err == nil {
		t.Fatal("accepted directory source")
	}
	got, err := os.ReadFile(target)
	if err != nil || string(got) != "existing" {
		t.Fatalf("destination=%q, error=%v", got, err)
	}
}

func BenchmarkMSIXCopyFile(b *testing.B) {
	dir := b.TempDir()
	source, target := filepath.Join(dir, "source"), filepath.Join(dir, "target")
	const size = 16 << 20
	file, err := os.Create(source)
	if err != nil {
		b.Fatal(err)
	}
	if err := file.Truncate(size); err != nil {
		b.Fatal(err)
	}
	if err := file.Close(); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.SetBytes(size)
	b.ResetTimer()
	for b.Loop() {
		if err := copyFile(source, target, 0644); err != nil {
			b.Fatal(err)
		}
	}
}
