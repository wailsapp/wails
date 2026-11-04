package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCleanupScopesArtifacts(t *testing.T) {
	root := filepath.Join(t.TempDir(), "project.git-data")
	files := map[string]bool{
		"internal/commands/myapp.DEB":          false,
		"internal/commands/myapp.RPM":          false,
		"internal/commands/myapp.ARCHLINUX":    false,
		"v3/internal/commands/myapp.DEB":       false,
		"bin/tool.test":                        false,
		"bin/testbuild-demo":                   false,
		"bin/appimage_testfiles":               false,
		"test-template-old/a.txt":              false,
		"test-template-old/nested/b.txt":       false,
		"downloads/myapp.DEB":                  true,
		"packages/myapp.RPM":                   true,
		"examples/internal/commands/myapp.DEB": true,
		".git/objects/keep.test":               true,
		"normal/file.txt":                      true,
		"test-template-file":                   true,
		"keep/testbuild-directory/nested.txt":  true,
		"keep/cache.test/nested.txt":           true,
	}
	for name := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("fixture"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	count, err := cleanup(root, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if count != 8 {
		t.Fatalf("removed %d items, want 8", count)
	}
	for name, keep := range files {
		_, err := os.Stat(filepath.Join(root, name))
		if keep && err != nil {
			t.Errorf("protected %s: %v", name, err)
		}
		if !keep && !os.IsNotExist(err) {
			t.Errorf("artifact %s was not removed: %v", name, err)
		}
	}
}

func TestCleanupMissingRoot(t *testing.T) {
	count, err := cleanup(filepath.Join(t.TempDir(), "missing"), io.Discard)
	if count != 0 || !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("count=%d error=%v", count, err)
	}
}

func TestCleanupDeleteFailure(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("requires POSIX directory permissions without root override")
	}
	root := t.TempDir()
	file := filepath.Join(root, "locked.test")
	if err := os.WriteFile(file, []byte("fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(root, 0500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chmod(root, 0700); err != nil {
			t.Error(err)
		}
	})
	var output bytes.Buffer
	count, err := cleanup(root, &output)
	if count != 0 || !errors.Is(err, os.ErrPermission) {
		t.Fatalf("count=%d error=%v", count, err)
	}
	if output.Len() != 0 {
		t.Fatalf("failed deletion reported success: %s", output.String())
	}
	if _, err := os.Stat(file); err != nil {
		t.Fatalf("failed deletion changed file: %v", err)
	}
}

func TestCleanupSkipsUnreadableGit(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("requires POSIX directory permissions without root override")
	}
	root := t.TempDir()
	git := filepath.Join(root, ".git")
	if err := os.Mkdir(git, 0700); err != nil {
		t.Fatal(err)
	}
	protected := filepath.Join(git, "keep.test")
	if err := os.WriteFile(protected, []byte("git fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "tool.test"), []byte("artifact"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(git, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chmod(git, 0700); err != nil {
			t.Error(err)
		}
	})
	count, err := cleanup(root, io.Discard)
	if err != nil || count != 1 {
		t.Fatalf("count=%d error=%v", count, err)
	}
	if err := os.Chmod(git, 0700); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(protected); err != nil || string(data) != "git fixture" {
		t.Fatalf("protected file=%q error=%v", data, err)
	}
}

func BenchmarkCleanupGit(b *testing.B) {
	root := b.TempDir()
	for directory := 0; directory < 50; directory++ {
		path := filepath.Join(root, ".git", "objects", fmt.Sprint(directory))
		if err := os.MkdirAll(path, 0700); err != nil {
			b.Fatal(err)
		}
		for file := 0; file < 40; file++ {
			if err := os.WriteFile(filepath.Join(path, fmt.Sprint(file)), []byte("fixture"), 0600); err != nil {
				b.Fatal(err)
			}
		}
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if count, err := cleanup(root, io.Discard); err != nil || count != 0 {
			b.Fatalf("count=%d error=%v", count, err)
		}
	}
}
