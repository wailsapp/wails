//go:build unix

package gosod

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

// TestExtract_MkdirAllFails relies on Unix permission semantics (mode 0555
// prevents subdirectory creation). Not portable to Windows.
func TestExtract_MkdirAllFails(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("root can always mkdir; skipping")
	}
	fsys := fstest.MapFS{
		"f.txt": {Data: []byte("x")},
	}
	td := New(fsys)
	tmp := t.TempDir()
	roDir := filepath.Join(tmp, "readonly")
	if err := os.Mkdir(roDir, 0555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(roDir, 0755) })
	targetDir := filepath.Join(roDir, "newsubdir")
	err := td.Extract(targetDir, nil)
	if err == nil {
		t.Fatal("expected error when target directory cannot be created")
	}
}
