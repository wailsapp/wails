package exec

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestRecursiveMatchRespectsPrefix(t *testing.T) {
	cases := []struct {
		path    string
		pattern string
		want    bool
	}{
		{"frontend/dist/foo.js", "frontend/dist/**/*", true},
		{"frontend/dist/assets/app.css", "frontend/dist/**/*", true},
		{".wake/cache.json", "frontend/dist/**/*", false}, // regression: prefix was ignored
		{"bin/badge", "frontend/dist/**/*", false},
		{"src/main.go", "**/*.go", true},
		{"src/main.ts", "**/*.go", false},
		{"frontend/bindings/index.ts", "frontend/bindings/**/*", true},
	}
	for _, c := range cases {
		if got := recursiveMatch(c.path, c.pattern); got != c.want {
			t.Errorf("recursiveMatch(%q, %q) = %v, want %v", c.path, c.pattern, got, c.want)
		}
	}
}

func TestRecursiveGlobConfinesToPrefix(t *testing.T) {
	root := t.TempDir()
	mustWrite := func(rel string) {
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	mustWrite("sub/a.js")
	mustWrite("sub/nested/b.js")
	mustWrite("other/c.js")
	mustWrite(".wake/cache.json")

	got := recursiveGlob(root, "sub/**/*")
	sort.Strings(got)
	want := []string{
		filepath.Join(root, "sub/a.js"),
		filepath.Join(root, "sub/nested/b.js"),
	}
	if len(got) != len(want) {
		t.Fatalf("recursiveGlob returned %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("recursiveGlob[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestClassifyGoCmd(t *testing.T) {
	cases := []struct {
		cmd  string
		want goCmdKind
	}{
		{"go build -o bin/app", goCmdBuild},
		{"go mod tidy", goCmdModTidy},
		{"go test ./...", goCmdNone},
		{"npm install", goCmdNone},
		{"GOOS=darwin go build", goCmdNone}, // env prefix not stripped here; classifier expects clean cmd
	}
	for _, c := range cases {
		if got := classifyGoCmd(c.cmd); got != c.want {
			t.Errorf("classifyGoCmd(%q) = %v, want %v", c.cmd, got, c.want)
		}
	}
}

func TestParseOutputFlag(t *testing.T) {
	if got := parseOutputFlag(`go build -buildvcs=false -o bin/badge`); got != "bin/badge" {
		t.Errorf("parseOutputFlag = %q, want bin/badge", got)
	}
	if got := parseOutputFlag(`go build -buildvcs=false`); got != "" {
		t.Errorf("parseOutputFlag (no -o) = %q, want empty", got)
	}
}
