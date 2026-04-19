package build

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v2/internal/fs"
)

func TestCopyDir(t *testing.T) {
	tmpDir := t.TempDir()

	srcDir := filepath.Join(tmpDir, "src")
	dstDir := filepath.Join(tmpDir, "dst")

	if err := fs.MkDirs(filepath.Join(srcDir, "subdir"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "subdir", "file2.txt"), []byte("world"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := copyDir(srcDir, dstDir); err != nil {
		t.Fatalf("copyDir() error = %v", err)
	}

	got, err := os.ReadFile(filepath.Join(dstDir, "file1.txt"))
	if err != nil {
		t.Fatalf("failed to read copied file: %v", err)
	}
	if string(got) != "hello" {
		t.Errorf("file1.txt = %q, want %q", string(got), "hello")
	}

	got2, err := os.ReadFile(filepath.Join(dstDir, "subdir", "file2.txt"))
	if err != nil {
		t.Fatalf("failed to read copied nested file: %v", err)
	}
	if string(got2) != "world" {
		t.Errorf("file2.txt = %q, want %q", string(got2), "world")
	}
}

func TestCopyDir_SourceNotExist(t *testing.T) {
	tmpDir := t.TempDir()

	err := copyDir(filepath.Join(tmpDir, "nonexistent"), filepath.Join(tmpDir, "dst"))
	if err == nil {
		t.Error("copyDir() expected error for non-existent source, got nil")
	}
}
