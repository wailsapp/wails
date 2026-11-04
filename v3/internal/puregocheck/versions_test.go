package main

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"
)

const girSource = `<?xml version="1.0"?>
<repository xmlns="http://www.gtk.org/introspection/core/1.0" xmlns:c="http://www.gtk.org/introspection/c/1.0">
  <namespace name="Gtk">
    <function name="new" c:identifier="gtk_new" version="4.6"/>
    <method name="m" c:identifier="gtk_method" version="4.14"/>
    <constructor name="c" c:identifier="gtk_ctor"/>
    <function name="x" c:identifier="gtk_nover"/>
  </namespace>
</repository>`

func TestReadGIRIndexesFunctionsByCName(t *testing.T) {
	dir := t.TempDir()
	writeTree(t, dir, map[string]string{"Gtk-4.0.gir": girSource})
	gir, err := readGIR(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]since{
		"gtk_new":    {"Gtk-4.0", "4.6"},
		"gtk_method": {"Gtk-4.0", "4.14"},
		"gtk_ctor":   {"Gtk-4.0", ""},
		"gtk_nover":  {"Gtk-4.0", ""},
	}
	for name, w := range want {
		if got, ok := gir[name]; !ok || got != w {
			t.Errorf("gir[%s] = %+v (%v), want %+v", name, got, ok, w)
		}
	}
}

func TestReadGIRNeedsFiles(t *testing.T) {
	if _, err := readGIR(t.TempDir()); err == nil {
		t.Error("readGIR accepted a directory with no .gir files")
	}
}

func TestLoadFloors(t *testing.T) {
	dir := t.TempDir()
	writeTree(t, dir, map[string]string{"ok.json": `{"Gtk-4.0":"4.8"}`, "bad.json": `{`})
	fl, err := loadFloors(filepath.Join(dir, "ok.json"))
	if err != nil || fl["Gtk-4.0"] != "4.8" {
		t.Errorf("loadFloors = %v, %v", fl, err)
	}
	if _, err := loadFloors(filepath.Join(dir, "bad.json")); err == nil {
		t.Error("loadFloors accepted invalid JSON")
	}
	if _, err := loadFloors(filepath.Join(dir, "missing.json")); err == nil {
		t.Error("loadFloors accepted a missing file")
	}
}

func TestCheckVersionsGatesSymbolsNewerThanTheFloor(t *testing.T) {
	r := &report{Symbols: []*entry{
		{Name: "gtk_new", Where: []string{"a.go:1"}},
		{Name: "gtk_old", Where: []string{"a.go:2"}},
		{Name: "gtk_opt", Optional: true, Where: []string{"a.go:3"}},
		{Name: "gtk_nover", Where: []string{"a.go:4"}},
		{Name: "other", Where: []string{"a.go:5"}},
		{Name: "unlisted", Where: []string{"a.go:6"}},
	}}
	gir := map[string]since{
		"gtk_new":   {"Gtk-4.0", "4.14"},
		"gtk_old":   {"Gtk-4.0", "4.6"},
		"gtk_opt":   {"Gtk-4.0", "4.12"},
		"gtk_nover": {"Gtk-4.0", ""},
		"other":     {"Foo-1.0", "9.0"},
	}
	checkVersions(r, gir, floors{"Gtk-4.0": "4.8"})

	byName := map[string]*entry{}
	for _, e := range r.Symbols {
		byName[e.Name] = e
	}
	if e := byName["gtk_new"]; !e.Gated || e.Since != "4.14" || e.Namespace != "Gtk-4.0" {
		t.Errorf("gtk_new = %+v, want gated at 4.14", e)
	}
	if e := byName["gtk_old"]; e.Gated {
		t.Errorf("gtk_old (4.6, below floor 4.8) gated: %+v", e)
	}
	if e := byName["gtk_opt"]; !e.Gated {
		t.Errorf("optional gtk_opt should still be gated (counted), got %+v", e)
	}
	if e := byName["other"]; e.Gated || e.Namespace != "Foo-1.0" {
		t.Errorf("other (no floor for its namespace) = %+v", e)
	}
	if e := byName["gtk_nover"]; e.Gated || e.Namespace != "Gtk-4.0" {
		t.Errorf("gtk_nover (no recorded version) = %+v", e)
	}
	if e := byName["unlisted"]; e.Namespace != "" {
		t.Errorf("symbol missing from the GIR got a namespace: %+v", e)
	}
	if len(r.Versions) != 1 || !strings.Contains(r.Versions[0], "gtk_new") {
		t.Errorf("Versions = %v, want one breach for the required gtk_new", r.Versions)
	}
}

func TestCheckGuardsFlagsUnguardedOptionalCalls(t *testing.T) {
	fset := token.NewFileSet()
	src := `package x
func unguarded() { foo(1) }
func guarded() { if _native_foo != nil { foo(1) } }
func required() { bar() }
func foo() { foo() }
func bodyless()
`
	f, err := parser.ParseFile(fset, "x.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	c := &collector{
		fset: fset, rel: "x.go", r: &report{},
		symbols: map[string]*entry{
			"foo": {Name: "foo", Optional: true},
			"bar": {Name: "bar"},
		},
	}
	c.checkGuards(f)
	if len(c.r.Unguarded) != 1 || !strings.Contains(c.r.Unguarded[0], "unguarded calls optional foo") {
		t.Errorf("Unguarded = %v, want only unguarded calling foo", c.r.Unguarded)
	}
}

func TestGuardedRecognisesVersionChecks(t *testing.T) {
	ok := body(t, `func f() { if gtkAtLeast(4, 10) { gtk_file_dialog_open(d) } }`)
	if !guarded(ok, "gtk_file_dialog_open") {
		t.Error("the _get_minor_version check was not recognised")
	}
	ok = body(t, `func f() { if haveSymbol(x) { gtk_file_dialog_open(d) } }`)
	if !guarded(ok, "gtk_file_dialog_open") {
		t.Error("haveSymbol was not recognised")
	}
}
