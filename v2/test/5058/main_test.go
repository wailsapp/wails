package test_5058

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func repoRoot() string {
	_, thisFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(thisFile), "..", "..", "..")
}

func TestSigningDocsUseUploadArtifactV4(t *testing.T) {
	websiteDir := filepath.Join(repoRoot(), "website")
	count := 0
	err := filepath.Walk(websiteDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Base(path) != "signing.mdx" {
			return nil
		}
		count++
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(data)
		if strings.Contains(content, "upload-artifact@v2") {
			rel, _ := filepath.Rel(websiteDir, path)
			t.Errorf("%s still references deprecated upload-artifact@v2", rel)
		}
		if strings.Contains(content, "upload-artifact@v3") {
			rel, _ := filepath.Rel(websiteDir, path)
			t.Errorf("%s references deprecated upload-artifact@v3", rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Walk error: %v", err)
	}
	if count == 0 {
		t.Fatal("no signing.mdx files found")
	}
	t.Logf("checked %d signing.mdx files", count)
}
