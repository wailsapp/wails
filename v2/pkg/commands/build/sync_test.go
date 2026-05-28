package build

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyDirContents(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	files := []struct {
		path    string
		content string
	}{
		{"index.html", "<html></html>"},
		{"assets/main.js", "console.log('hi')"},
		{"assets/style.css", "body {}"},
		{"assets/sub/deep.txt", "deep"},
	}

	for _, f := range files {
		fullPath := filepath.Join(srcDir, f.path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(f.content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if err := copyDirContents(srcDir, dstDir); err != nil {
		t.Fatalf("copyDirContents() error = %v", err)
	}

	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(dstDir, f.path))
		if err != nil {
			t.Errorf("file %s not found in destination: %v", f.path, err)
			continue
		}
		if string(data) != f.content {
			t.Errorf("file %s: got %q, want %q", f.path, string(data), f.content)
		}
	}
}
