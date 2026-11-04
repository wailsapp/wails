package dirscan

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
)

func TestWalkBoundsConcurrencySkipsDirectoriesAndJoinsCallbacks(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"a/one", "a/two", "b/three", "skip/hidden"} {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(name), 0644); err != nil {
			t.Fatal(err)
		}
	}
	for _, workers := range []int{1, 2, 8} {
		var mu sync.Mutex
		found := map[string]bool{}
		var active, peak atomic.Int32
		err := Walk(root, workers, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			count := active.Add(1)
			defer active.Add(-1)
			for old := peak.Load(); count > old && !peak.CompareAndSwap(old, count); old = peak.Load() {
			}
			if entry.IsDir() && entry.Name() == "skip" {
				return fs.SkipDir
			}
			if !entry.IsDir() {
				relative, _ := filepath.Rel(root, path)
				mu.Lock()
				found[filepath.ToSlash(relative)] = true
				mu.Unlock()
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(found) != 3 || found["skip/hidden"] {
			t.Fatalf("wrong files: %v", found)
		}
		if peak.Load() > int32(workers) || active.Load() != 0 {
			t.Fatalf("callbacks peak=%d active=%d workers=%d", peak.Load(), active.Load(), workers)
		}
	}
}

func TestWalkPreservesErrorsAndStopsWithoutFollowingSymlinks(t *testing.T) {
	root := t.TempDir()
	sentinel := errors.New("visit failed")
	if err := os.WriteFile(filepath.Join(root, "file"), nil, 0644); err != nil {
		t.Fatal(err)
	}
	if err := Walk(root, 4, func(string, fs.DirEntry, error) error { return sentinel }); !errors.Is(err, sentinel) {
		t.Fatal(err)
	}
	if err := Walk(filepath.Join(root, "missing"), 4, func(_ string, _ fs.DirEntry, err error) error { return err }); !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
	if err := Walk(root, 4, func(string, fs.DirEntry, error) error { return fs.SkipAll }); err != nil {
		t.Fatal(err)
	}
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "outside-file"), nil, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "link")); err != nil {
		t.Skipf("symlink fixture unavailable: %v", err)
	}
	var outsideVisited atomic.Bool
	if err := Walk(root, 4, func(_ string, entry fs.DirEntry, err error) error {
		if entry != nil && entry.Name() == "outside-file" {
			outsideVisited.Store(true)
		}
		return err
	}); err != nil {
		t.Fatal(err)
	}
	if outsideVisited.Load() {
		t.Fatal("followed a symlink outside the project")
	}
}
