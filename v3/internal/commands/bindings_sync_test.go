package commands

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for name, content := range files {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o777); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o666); err != nil {
			t.Fatal(err)
		}
	}
}

func readTree(t *testing.T, root string) map[string]string {
	t.Helper()
	tree := make(map[string]string)
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		tree[filepath.ToSlash(rel)] = string(data)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return tree
}

func checkTree(t *testing.T, root string, want map[string]string) {
	t.Helper()
	got := readTree(t, root)
	if len(got) != len(want) {
		t.Errorf("got %d files, want %d: %v", len(got), len(want), got)
	}
	for name, content := range want {
		if got[name] != content {
			t.Errorf("file %q: got %q, want %q", name, got[name], content)
		}
	}
}

func TestSyncDirsNoDestination(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, ".bindings-tmp-1")
	dst := filepath.Join(root, "bindings")
	files := map[string]string{"svc/a.ts": "a", "index.ts": "idx"}
	writeFiles(t, src, files)

	if err := syncDirs(src, dst); err != nil {
		t.Fatal(err)
	}
	checkTree(t, dst, files)
	if _, err := os.Lstat(src); !os.IsNotExist(err) {
		t.Errorf("source directory should be gone, got err=%v", err)
	}
}

func TestSyncDirsUpdatesAndDeletes(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, ".bindings-tmp-1")
	dst := filepath.Join(root, "bindings")
	writeFiles(t, src, map[string]string{
		"svc/a.ts": "new a",
		"svc/b.ts": "same b",
		"new.ts":   "brand new",
	})
	writeFiles(t, dst, map[string]string{
		"svc/a.ts":       "old a",
		"svc/b.ts":       "same b",
		"stale.ts":       "stale",
		"gone/nested.ts": "stale dir",
	})

	// Backdate the unchanged file so we can verify it is left untouched.
	past := time.Now().Add(-time.Hour)
	unchanged := filepath.Join(dst, "svc", "b.ts")
	if err := os.Chtimes(unchanged, past, past); err != nil {
		t.Fatal(err)
	}

	if err := syncDirs(src, dst); err != nil {
		t.Fatal(err)
	}
	checkTree(t, dst, map[string]string{
		"svc/a.ts": "new a",
		"svc/b.ts": "same b",
		"new.ts":   "brand new",
	})
	if _, err := os.Lstat(filepath.Join(dst, "gone")); !os.IsNotExist(err) {
		t.Errorf("stale directory should be gone, got err=%v", err)
	}
	info, err := os.Stat(unchanged)
	if err != nil {
		t.Fatal(err)
	}
	if !info.ModTime().Equal(past) {
		t.Errorf("unchanged file was rewritten: mtime %v, want %v", info.ModTime(), past)
	}
	if _, err := os.Lstat(src); !os.IsNotExist(err) {
		t.Errorf("source directory should be gone, got err=%v", err)
	}
}

func TestSyncDirsTypeMismatch(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, ".bindings-tmp-1")
	dst := filepath.Join(root, "bindings")
	// In src, "thing" is a file and "other" is a directory; in dst the
	// types are reversed.
	writeFiles(t, src, map[string]string{
		"thing":    "now a file",
		"other/x":  "in dir",
		"plain.ts": "p",
	})
	writeFiles(t, dst, map[string]string{
		"thing/nested": "was a dir",
		"other":        "was a file",
	})

	if err := syncDirs(src, dst); err != nil {
		t.Fatal(err)
	}
	checkTree(t, dst, map[string]string{
		"thing":    "now a file",
		"other/x":  "in dir",
		"plain.ts": "p",
	})
}

func TestSyncDirsDestinationIsFile(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, ".bindings-tmp-1")
	dst := filepath.Join(root, "bindings")
	files := map[string]string{"a.ts": "a"}
	writeFiles(t, src, files)
	if err := os.WriteFile(dst, []byte("not a directory"), 0o666); err != nil {
		t.Fatal(err)
	}

	if err := syncDirs(src, dst); err != nil {
		t.Fatal(err)
	}
	checkTree(t, dst, files)
	if _, err := os.Lstat(src); !os.IsNotExist(err) {
		t.Errorf("source directory should be gone, got err=%v", err)
	}
}

func TestSyncDirsIdenticalContent(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, ".bindings-tmp-1")
	dst := filepath.Join(root, "bindings")
	files := map[string]string{"a.ts": "a", "sub/b.ts": "b"}
	writeFiles(t, src, files)
	writeFiles(t, dst, files)

	if err := syncDirs(src, dst); err != nil {
		t.Fatal(err)
	}
	checkTree(t, dst, files)
}
