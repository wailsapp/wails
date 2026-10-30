package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const stubTree = "//go:build darwin && purego\n\npackage application\n\nfunc stubbed() bool { return true }\n"
const stubCgo = "//go:build darwin && !purego\n\npackage application\n\nimport \"C\"\n\nfunc stubbed() bool { return C.foo() != 0 }\n"

func TestRunRejectsBadUsage(t *testing.T) {
	for _, args := range [][]string{
		nil,
		{"bogus"},
		{"-nope", "static"},
		{"-goos", "plan9", "resolve"},
	} {
		if got := run(args); got != 2 {
			t.Errorf("run(%q) = %d, want 2", args, got)
		}
	}
}

func TestRunStaticStatus(t *testing.T) {
	clean := t.TempDir()
	writeTree(t, clean, map[string]string{"pkg/application/a_purego.go": "//go:build darwin && purego\n\npackage application\n\nfunc f() { w.send(\"init\") }\n"})
	jsonPath := filepath.Join(t.TempDir(), "report.json")
	if got := run([]string{"-root", clean, "-goos", "darwin", "-tags", "purego", "-json", jsonPath, "static"}); got != 0 {
		t.Errorf("clean tree: run = %d, want 0", got)
	}
	if data, err := os.ReadFile(jsonPath); err != nil || !strings.Contains(string(data), `"goos": "darwin"`) {
		t.Errorf("-json report not written: %v", err)
	}

	stubbed := t.TempDir()
	writeTree(t, stubbed, map[string]string{
		"pkg/application/stub_purego.go": stubTree,
		"pkg/application/stub_cgo.go":    stubCgo,
	})
	if got := run([]string{"-root", stubbed, "-goos", "darwin", "-tags", "purego", "static"}); got != 1 {
		t.Errorf("stubbed tree: run = %d, want 1", got)
	}
	if got := run([]string{"-root", stubbed, "-goos", "darwin", "-tags", "purego", "-allow", writeAllow(t, "stubbed"), "static"}); got != 0 {
		t.Errorf("allowed stub: run = %d, want 0", got)
	}
	if got := run([]string{"-root", t.TempDir(), "-json", filepath.Join(t.TempDir(), "no", "such", "dir.json"), "static"}); got != 2 {
		t.Errorf("unwritable -json: run = %d, want 2", got)
	}
}

func writeAllow(t *testing.T, names ...string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "allow.txt")
	if err := os.WriteFile(p, []byte("# intentional\n\n"+strings.Join(names, "\n")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRunGIRVersionChecks(t *testing.T) {
	root := t.TempDir()
	writeTree(t, root, map[string]string{
		"pkg/application/gtk_linux.go": "//go:build linux\n\npackage application\n\nvar _ = RegisterLibFunc(&fn, lib, \"gtk_method\")\n",
	})
	gir := t.TempDir()
	writeTree(t, gir, map[string]string{"Gtk-4.0.gir": girSource})
	floors := filepath.Join(t.TempDir(), "floors.json")
	if err := os.WriteFile(floors, []byte(`{"Gtk-4.0":"4.8"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	args := []string{"-root", root, "-goos", "linux", "-gir", gir, "-floors", floors, "-allow", writeAllow(t), "static"}
	if got := run(args); got != 1 {
		t.Errorf("required gtk_new above the floor: run = %d, want 1", got)
	}
	if got := run([]string{"-root", root, "-goos", "linux", "-gir", t.TempDir(), "static"}); got != 2 {
		t.Errorf("GIR directory without .gir files: run = %d, want 2", got)
	}
	bad := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(bad, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := run([]string{"-root", root, "-goos", "linux", "-gir", gir, "-floors", bad, "static"}); got != 2 {
		t.Errorf("unreadable floors: run = %d, want 2", got)
	}
}

func TestRunResolveOnThisOS(t *testing.T) {
	// An empty tree names nothing, so resolve only loads the platform
	// libraries and must succeed.
	if got := run([]string{"-root", t.TempDir(), "-goos", runtime.GOOS, "resolve"}); got != 0 {
		t.Errorf("resolve on an empty tree = %d, want 0", got)
	}
}

func TestApplyAllowMovesNamedStubs(t *testing.T) {
	r := &report{Stubs: []stub{{Name: "A.keep"}, {Name: "b"}}}
	applyAllow(r, writeAllow(t, "b"))
	if len(r.Allowed) != 1 || r.Allowed[0].Name != "b" || len(r.Stubs) != 1 || r.Stubs[0].Name != "A.keep" {
		t.Errorf("applyAllow: stubs=%+v allowed=%+v", r.Stubs, r.Allowed)
	}
	applyAllow(r, filepath.Join(t.TempDir(), "missing.txt"))
	if len(r.Stubs) != 1 {
		t.Error("a missing allow file must change nothing")
	}
}

func TestSummaryExitStatus(t *testing.T) {
	cases := []struct {
		name string
		r    report
		want int
	}{
		{"clean", report{GOOS: "linux"}, 0},
		{"absent optional only", report{GOOS: "linux", Absent: []string{"symbol x"}}, 0},
		{"stub", report{GOOS: "linux", Stubs: []stub{{Name: "f"}}}, 1},
		{"missing", report{GOOS: "linux", Missing: []string{"symbol y"}}, 1},
		{"version", report{GOOS: "linux", Versions: []string{"v"}}, 1},
		{"unguarded", report{GOOS: "linux", Unguarded: []string{"u"}}, 1},
	}
	for _, c := range cases {
		r := c.r
		if got := summary(&r); got != c.want {
			t.Errorf("%s: summary = %d, want %d", c.name, got, c.want)
		}
	}
}

func TestSummaryCountsDarwinAndResolvedSections(t *testing.T) {
	r := &report{
		GOOS:      "darwin",
		Symbols:   []*entry{{Name: "a"}, {Name: "b", Optional: true, Gated: true}},
		Classes:   []*entry{{Name: "C", Ours: true}, {Name: "D"}},
		Selectors: []*entry{{Name: "s:", Optional: true}},
		Dynamic:   []string{"x.go:1 RegisterLibFunc"},
		Resolved:  true,
		Loaded:    []string{"macOS 27.0"},
		Missing:   []string{"selector q:"},
		Absent:    []string{"class Z"},
		Allowed:   []stub{{Name: "f"}},
	}
	if got := summary(r); got != 1 {
		t.Errorf("summary with a missing selector = %d, want 1", got)
	}
	if n := countGated(r.Symbols); n != 1 {
		t.Errorf("countGated = %d, want 1", n)
	}
}

func TestSortedEntriesOrdersByNameAndWhere(t *testing.T) {
	m := map[string]*entry{
		"b": {Name: "b", Where: []string{"z.go:1", "a.go:9"}},
		"a": {Name: "a"},
	}
	got := sortedEntries(m)
	if len(got) != 2 || got[0].Name != "a" || got[1].Name != "b" {
		t.Fatalf("sortedEntries order = %v", got)
	}
	if got[1].Where[0] != "a.go:9" {
		t.Errorf("Where not sorted: %v", got[1].Where)
	}
}
