//go:build integration

package browser

import (
	"os"
	"path/filepath"
	"testing"
)

// TestOpenURL_Success launches the default browser with a real URL.
// Run with: go test -tags integration ./internal/browser/...
func TestOpenURL_Success(t *testing.T) {
	if err := OpenURL("https://example.com"); err != nil {
		t.Errorf("OpenURL returned unexpected error: %v", err)
	}
}

// TestOpenFile_Success opens a real temp file via the OS file handler.
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
