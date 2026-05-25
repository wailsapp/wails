package browser

import (
	"os"
	"path/filepath"
	"testing"
)

// TestOpenURL_Success calls OpenURL with a real URL.
// On macOS this invokes `open <url>` which starts the default browser.
// The test only asserts that Start() returns nil (process launched).
func TestOpenURL_Success(t *testing.T) {
	if err := OpenURL("https://example.com"); err != nil {
		t.Errorf("OpenURL returned unexpected error: %v", err)
	}
}

// TestOpenFile_Success calls OpenFile with a path to a temp file.
func TestOpenFile_Success(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(path, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if err := OpenFile(path); err != nil {
		t.Errorf("OpenFile returned unexpected error: %v", err)
	}
}

// TestOpen_StartError exercises the cmd.Start() error path by making the
// platform command unresolvable (empty PATH).
func TestOpen_StartError(t *testing.T) {
	t.Setenv("PATH", "/nonexistent_path_that_does_not_exist")
	err := OpenURL("https://example.com")
	if err == nil {
		t.Error("open() should have returned an error with an empty PATH, but got nil")
	}
}
