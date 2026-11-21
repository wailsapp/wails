package sourceignore

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestPatternsMatchGit(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Git oracle unavailable")
	}
	tests := []struct {
		pattern string
		paths   []string
	}{
		{"broken[\ntriple***star\nliteral\\[", []string{"broken[", "tripleXstar", "literal["}},
		{"?.txt\né.go\n[é].tmp", []string{"é.txt", "é.go", "é.tmp", "e.tmp"}},
		{"*.go\n!keep.go", []string{"a.go", "keep.go", "dir/a.go", "dir/keep.go", "a.gox"}},
		{"/root\nfolder/", []string{"root", "dir/root", "folder/", "folder/file", "dir/folder/", "folder"}},
		{"a/**/z", []string{"a/z", "a/b/z", "a/b/c/z", "x/a/z", "a/z/f"}},
		{"**/a/b", []string{"a/b", "x/a/b", "x/y/a/b", "a/z/b"}},
		{"abc/**", []string{"abc/", "abc/a", "abc/a/b", "x/abc/a"}},
		{"a**b\nfoo?bar\n[!a-c].txt\n[[:digit:]].log\n[]a].tmp", []string{"aXXb", "a/b", "fooXbar", "foobar", "d.txt", "b.txt", "7.log", "a.log", "].tmp", "a.tmp"}},
		{"\\#name\n\\!name\nspace\\ \ntrim   \ntrailing\\", []string{"#name", "!name", "space ", "space", "trim", "trim ", "trailing", "trailing\\"}},
		{"hidden/\n!hidden/keep", []string{"hidden/", "hidden/keep", "hidden/other"}},
		{"/*\n!/foo\n/foo/*\n!/foo/bar", []string{"foo/", "foo/bar/", "foo/bar/file", "foo/else", "else"}},
		{"*.tmp\n!dir/keep.tmp", []string{"dir/keep.tmp", "else/keep.tmp", "dir/file.tmp"}},
	}
	for _, test := range tests {
		t.Run(test.pattern, func(t *testing.T) {
			root := t.TempDir()
			if data, err := exec.Command("git", "init", "-q", root).CombinedOutput(); err != nil {
				t.Fatalf("git init: %s %v", data, err)
			}
			if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte(test.pattern+"\n"), 0600); err != nil {
				t.Fatal(err)
			}
			m := New(strings.Split(test.pattern, "\n")...)
			for _, path := range test.paths {
				clean := strings.TrimSuffix(path, "/")
				if strings.HasSuffix(path, "/") {
					if err := os.MkdirAll(filepath.Join(root, filepath.FromSlash(clean)), 0700); err != nil {
						t.Fatal(err)
					}
				}
				cmd := exec.Command("git", "-c", "core.excludesFile="+os.DevNull, "-C", root, "check-ignore", "--no-index", "-q", "--", clean)
				cmd.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL="+os.DevNull)
				err := cmd.Run()
				want := err == nil
				if err != nil {
					if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 1 {
						t.Fatal(err)
					}
				}
				got := m.Match(strings.Split(strings.TrimSuffix(path, "/"), "/"), strings.HasSuffix(path, "/"))
				if strings.HasSuffix(path, "/") {
					if err := os.RemoveAll(filepath.Join(root, filepath.FromSlash(clean))); err != nil {
						t.Fatal(err)
					}
				}
				if got != want {
					t.Errorf("path %q: excluded=%v, git=%v", path, got, want)
				}
			}
		})
	}
}

func TestReadNestedRulesAndIgnoredParents(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		".gitignore":             "*.tmp\nclosed/\n",
		"nested/.gitignore":      "!keep.tmp\n/local\n",
		"nested/deep/.gitignore": "keep.tmp\n",
		"closed/.gitignore":      "!keep.tmp\n",
		"é/.gitignore":           "/x\n",
	}
	for name, data := range files {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
	}
	m, err := Read(root)
	if err != nil {
		t.Fatal(err)
	}
	for path, want := range map[string]bool{"x.tmp": true, "nested/keep.tmp": false, "nested/x.tmp": true, "nested/local": true, "nested/deep/local": false, "nested/deep/keep.tmp": true, "closed/keep.tmp": true, "é/x": true, "é/deep/x": false} {
		if got := m.Match(strings.Split(path, "/"), false); got != want {
			t.Errorf("%s=%v want %v", path, got, want)
		}
	}
	if _, err := Read(filepath.Join(root, "missing")); err == nil {
		t.Fatal("missing root accepted")
	}
}

func TestReadDoesNotFollowIgnoreSymlinks(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "rules")
	if err := os.WriteFile(outside, []byte("*\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, ".gitignore")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	m, err := Read(root)
	if err != nil {
		t.Fatal(err)
	}
	if m.Match([]string{"file.go"}, false) {
		t.Fatal("followed .gitignore symlink")
	}
}

func BenchmarkMatch(b *testing.B) {
	m := New("node_modules/", "build/", "*.tmp", "*.generated.go", "!keep.generated.go", "assets/**/cache/*", "nested/[[:digit:]].tmp")
	paths := [][]string{{"src", "backend", "service.go"}, {"node_modules", "module", "file.js"}, {"assets", "icons", "cache", "large.png"}, {"nested", "keep.generated.go"}}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m.Match(paths[i%len(paths)], false)
	}
}

func TestReadToleratesDirectoriesDeletedMidWalk(t *testing.T) {
	if runtime.GOOS == "windows" {
		// Windows refuses to open a directory another goroutine is deleting (a sharing violation, not "not exist"),
		// and Read rightly reports that; the race this guards is the POSIX unlink-while-walking one.
		t.Skip("concurrent deletion fails differently on Windows")
	}
	root := t.TempDir()
	staging := filepath.Join(root, "staging")
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; ; i++ {
			select {
			case <-stop:
				return
			default:
			}
			dir := filepath.Join(staging, "stage-"+strings.Repeat("x", i%7), "nested")
			_ = os.MkdirAll(dir, 0o755)
			_ = os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("*.tmp\n"), 0o644)
			_ = os.RemoveAll(staging)
		}
	}()
	defer func() { close(stop); <-done }()
	for range 2000 {
		if _, err := Read(root); err != nil {
			t.Fatalf("Read failed on a directory deleted mid-walk: %v", err)
		}
	}
}
