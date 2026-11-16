package packagesign

import (
	"os"
	"path/filepath"
	"testing"
)

func TestKeyFileReference(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "export"), []byte("key"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		reference string
		file      bool
	}{
		{"", false}, {"release@example.com", false}, {"0123456789ABCDEF", false},
		{"export", true}, {"missing.asc", true}, {"./missing", true}, {filepath.Join(root, "missing"), true},
	} {
		path, file := KeyFile(root, test.reference)
		if file != test.file || (file && !filepath.IsAbs(path)) {
			t.Errorf("reference %q: %q %v", test.reference, path, file)
		}
	}
}
