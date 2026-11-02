package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReplaceFiles(t *testing.T) {
	root := t.TempDir()
	files := map[string]struct{ input, want string }{
		"a.txt":       {"before before β\n", "after after β\n"},
		"nested/b.md": {"before after\n", "after after\n"},
		"skip.txt":    {"before\n", "before\n"},
		"keep.go":     {"before\n", "before\n"},
		"noext":       {"before\n", "before\n"},
	}
	modes := make(map[string]os.FileMode)
	for name, tc := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(tc.input), 0600); err != nil {
			t.Fatal(err)
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		modes[name] = info.Mode().Perm()
	}
	if err := run([]string{"replace", "-help=false", "--dir", root,
		"-old=before", "--new", "after", "--ext", ".txt,.md", "-ignore", "skip.txt"}, io.Discard); err != nil {
		t.Fatal(err)
	}
	for name, tc := range files {
		checkFile(t, filepath.Join(root, name), tc.want, modes[name])
	}
}

func checkFile(t *testing.T, path, want string, mode os.FileMode) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != want {
		t.Errorf("%s = %q, want %q", path, data, want)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != mode {
		t.Errorf("%s changed mode", path)
	}

}

func TestMissingDirectory(t *testing.T) {
	err := ReplaceInFiles(&ReplaceInFilesOptions{Dir: filepath.Join(t.TempDir(), "missing")})
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("error = %v, want missing-directory error", err)
	}
}

func TestHelp(t *testing.T) {
	for _, args := range [][]string{nil, {"--help"}, {"-h"}, {"replace", "--help"}, {"replace", "-h"}} {
		var output bytes.Buffer
		if err := run(args, &output); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(output.String(), "replace") {
			t.Errorf("missing command help for %v: %q", args, output.String())
		}
	}
}

func TestInvalidArguments(t *testing.T) {
	for _, args := range [][]string{{"unknown"}, {"--version"}, {"replace", "-unknown"},
		{"replace", "-dir"}, {"replace", "extra"}} {
		if err := run(args, io.Discard); err == nil {
			t.Errorf("accepted invalid arguments: %v", args)
		}
	}
}
