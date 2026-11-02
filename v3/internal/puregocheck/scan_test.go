package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTree(t *testing.T, dir string, files map[string]string) {
	t.Helper()
	for name, src := range files {
		p := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func parseDecls(t *testing.T, src string) (*token.FileSet, *ast.File) {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "x.go", "package x\n"+src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	return fset, f
}

func TestCallRecordsSymbolsClassesAndSelectors(t *testing.T) {
	c := collectSrc(t, "darwin", `
func f() {
	Dlsym(h, "g_free")
	RegisterLibFunc(&fn, lib, "gtk_thing")
	RegisterLibFunc(&fn2, lib, name)
	GetClass("NSApp")
	class("NSWindow")
	classExists("NSFoo")
	objc.RegisterClass("Ours", "NSSuper")
	registerDelegateClass("Del")
	sel_("description")
	w.send("init")
	w.get(x, "setTitle:")
	w.respondsTo(x, "optionalSel:")
}`)
	checks := []struct {
		m        map[string]*entry
		name     string
		optional bool
		ours     bool
	}{
		{c.symbols, "g_free", true, false},
		{c.symbols, "gtk_thing", false, false},
		{c.classes, "NSApp", false, false},
		{c.classes, "NSWindow", false, false},
		{c.classes, "NSFoo", true, false},
		{c.classes, "Ours", false, true},
		{c.classes, "NSSuper", false, false},
		{c.classes, "Del", false, true},
		{c.selectors, "description", false, false},
		{c.selectors, "init", false, false},
		{c.selectors, "setTitle:", false, false},
		{c.selectors, "optionalSel:", true, false},
	}
	for _, want := range checks {
		e := want.m[want.name]
		if e == nil {
			t.Errorf("%s not recorded", want.name)
			continue
		}
		if e.Optional != want.optional || e.Ours != want.ours {
			t.Errorf("%s: optional=%v ours=%v, want optional=%v ours=%v", want.name, e.Optional, e.Ours, want.optional, want.ours)
		}
	}
	if len(c.r.Dynamic) != 1 || !strings.Contains(c.r.Dynamic[0], "RegisterLibFunc") {
		t.Errorf("Dynamic = %v, want one non-literal RegisterLibFunc", c.r.Dynamic)
	}
}

func TestCallOnlyRecordsObjCOnDarwin(t *testing.T) {
	c := collectSrc(t, "linux", `
func f() {
	GetClass("NSApp")
	w.send("init")
	Dlsym(h, "g_free")
}`)
	if len(c.classes) != 0 || len(c.selectors) != 0 {
		t.Errorf("linux recorded ObjC names: classes=%d selectors=%d", len(c.classes), len(c.selectors))
	}
	if c.symbols["g_free"] == nil {
		t.Error("C symbols are recorded on every GOOS")
	}
}

func TestCompositeRecordsMethodDefsAndSymbolTables(t *testing.T) {
	src := `
var table = map[string]any{"g_strfreev": &fn}
var defs = []objc.MethodDef{{Cmd: sel_("windowWillClose:")}}
func setup() {
	_ = []objc.MethodDef{{Cmd: sel_("windowDidResize:")}}
	_ = []string{"resignKey:"}
}`
	c := collectSrc(t, "darwin", src)
	if e := c.symbols["g_strfreev"]; e == nil || e.Optional {
		t.Errorf("symbol table entry g_strfreev = %+v, want required", e)
	}
	for _, sel := range []string{"windowWillClose:", "windowDidResize:", "resignKey:"} {
		if e := c.selectors[sel]; e == nil || !e.Ours {
			t.Errorf("%s = %+v, want recorded as implemented by the backend", sel, e)
		}
	}
	if e := collectSrc(t, "linux", src).selectors["resignKey:"]; e != nil {
		t.Errorf("linux recorded the MethodDef-table string resignKey: as %+v", e)
	}
}

func TestNodeRecordsLibraryNames(t *testing.T) {
	c := collectSrc(t, "linux", `var libs = []string{"libwebkitgtk-6.0.so.4", "not a library", "libfoo.dylib"}`)
	got := strings.Join(c.r.libCandidates, " ")
	if got != "libwebkitgtk-6.0.so.4 libfoo.dylib" {
		t.Errorf("libCandidates = %q", got)
	}
}

func TestAddKeepsRequiredAndOurs(t *testing.T) {
	m := map[string]*entry{}
	add(m, "x", "a.go:1", true, false)
	if !m["x"].Optional || m["x"].Ours {
		t.Fatalf("first optional use = %+v", m["x"])
	}
	add(m, "x", "a.go:2", false, false)
	add(m, "x", "a.go:3", true, false)
	if m["x"].Optional {
		t.Error("a later required use must make the entry required")
	}
	add(m, "x", "a.go:2", false, true)
	if !m["x"].Ours {
		t.Error("ours flag not set")
	}
	if len(m["x"].Where) != 3 {
		t.Errorf("Where = %v, want three distinct sites", m["x"].Where)
	}
	for i := 0; i < 20; i++ {
		add(m, "x", "site"+string(rune('a'+i)), false, false)
	}
	if len(m["x"].Where) > 8 {
		t.Errorf("Where grows past 8 sites: %d", len(m["x"].Where))
	}
}

func TestLitAndFunName(t *testing.T) {
	if s, ok := lit(&ast.BasicLit{Kind: token.STRING, Value: `"abc"`}); !ok || s != "abc" {
		t.Errorf("lit string = %q, %v", s, ok)
	}
	if _, ok := lit(&ast.BasicLit{Kind: token.INT, Value: "3"}); ok {
		t.Error("lit accepted an integer")
	}
	if _, ok := lit(&ast.BasicLit{Kind: token.STRING, Value: "not-quoted"}); ok {
		t.Error("lit accepted an unquoted string")
	}
	_, f := parseDecls(t, `func x() { a(); b.c(); d[0](); e[int, int](); (g())() }`)
	var names []string
	ast.Inspect(f, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			names = append(names, funName(call.Fun))
		}
		return true
	})
	if got := strings.Join(names, ","); got != "a,c,d,e,,g" {
		t.Errorf("funName = %q", got)
	}
}

func TestMentions(t *testing.T) {
	_, f := parseDecls(t, `func f() { MethodDef{} }
func g() { x := 1; _ = x }`)
	if !mentions(f.Decls[0], "MethodDef") || mentions(f.Decls[1], "MethodDef") {
		t.Error("mentions did not find the identifier only where it occurs")
	}
}

func TestRespondedSelectors(t *testing.T) {
	_, f := parseDecls(t, `func f(w id) {
	respondsTo(w, "a:")
	respondsTo(w, sel)
	respondsTo(w)
	other(w, "b:")
}`)
	got := respondedSelectors(f.Decls[0].(*ast.FuncDecl).Body)
	if len(got) != 1 || !got["a:"] {
		t.Errorf("respondedSelectors = %v, want only a:", got)
	}
}

func TestBuildFilesSelectsByGOOSTagsAndCgo(t *testing.T) {
	dir := t.TempDir()
	writeTree(t, dir, map[string]string{
		"a.go":        "package x\n",
		"b_linux.go":  "package x\n",
		"c.go":        "//go:build purego\n\npackage x\n",
		"d.go":        "package x\n\nimport \"C\"\n",
		"e_test.go":   "package x\n",
		"sub/skip.go": "package y\n",
	})
	names := func(ps []parsed) string {
		var out []string
		for _, p := range ps {
			out = append(out, filepath.Base(p.path))
		}
		return strings.Join(out, ",")
	}
	fset := token.NewFileSet()
	got, err := buildFiles(fset, dir, "darwin", []string{"purego"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if names(got) != "a.go,c.go" {
		t.Errorf("pure darwin files = %s", names(got))
	}
	got, err = buildFiles(fset, dir, "darwin", nil, true)
	if err != nil {
		t.Fatal(err)
	}
	if names(got) != "a.go,d.go" {
		t.Errorf("cgo darwin files = %s", names(got))
	}
	got, err = buildFiles(fset, dir, "linux", nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if names(got) != "a.go,b_linux.go" {
		t.Errorf("linux files = %s", names(got))
	}
}

func TestBuildFilesReportsParseErrors(t *testing.T) {
	dir := t.TempDir()
	writeTree(t, dir, map[string]string{"bad.go": "package x\nfunc (\n"})
	if _, err := buildFiles(token.NewFileSet(), dir, "darwin", nil, false); err == nil {
		t.Error("a file that does not parse was accepted")
	}
}

func TestCgoTags(t *testing.T) {
	if got := cgoTags("darwin", []string{"purego", "devtools"}); strings.Join(got, ",") != "devtools" {
		t.Errorf("darwin cgo tags = %v", got)
	}
	if got := cgoTags("linux", []string{"purego", "gtk3"}); strings.Join(got, ",") != "gtk3" {
		t.Errorf("linux cgo tags = %v", got)
	}
}

func TestManifestOptional(t *testing.T) {
	root := t.TempDir()
	writeTree(t, root, map[string]string{
		"internal/puregogen/a.json":   `{"functions":[{"name":"f1","optional":true},{"name":"f2"}]}`,
		"internal/puregogen/bad.json": `{not json`,
	})
	got := manifestOptional(root)
	if !got["f1"] || got["f2"] || len(got) != 1 {
		t.Errorf("manifestOptional = %v, want only f1", got)
	}
	if len(manifestOptional(t.TempDir())) != 0 {
		t.Error("a root without manifests should give no optional names")
	}
}

func TestFuncKey(t *testing.T) {
	_, f := parseDecls(t, `func plain() {}
func (w *T) ptr() {}
func (w T) val() {}
func (w *G[K]) gen() {}
func (w pkg.T) sel() {}`)
	want := []string{"plain", "T.ptr", "T.val", "G.gen", "sel"}
	for i, d := range f.Decls {
		if got := funcKey(d.(*ast.FuncDecl)); got != want[i] {
			t.Errorf("funcKey[%d] = %q, want %q", i, got, want[i])
		}
	}
}

func TestTrivialAndConstantEdges(t *testing.T) {
	for src, want := range map[string]bool{
		"func f() { panic(x) }":         true,
		"func f() { notYet() }":         true,
		"func f() { notImplemented() }": true,
		"func f() { return }":           true,
		"func f() { x := 1; _ = x }":    false,
		"func f() { x }":                false,
		"func f() { g() }":              false,
	} {
		if got := trivial(body(t, src)); got != want {
			t.Errorf("trivial(%q) = %v, want %v", src, got, want)
		}
	}
	if !trivial(nil) {
		t.Error("a nil body is trivial")
	}
	_, f := parseDecls(t, `var a = -1
var b = nil
var c = true
var d = x
var e = []int{}
var g = []int{1}
var h = f()`)
	want := []bool{true, true, true, false, true, false, false}
	for i, spec := range f.Decls {
		val := spec.(*ast.GenDecl).Specs[0].(*ast.ValueSpec).Values[0]
		if got := constant(val); got != want[i] {
			t.Errorf("constant(decl %d) = %v, want %v", i, got, want[i])
		}
	}
}

func TestFindStubsPairsPureAndCgoTwins(t *testing.T) {
	root := t.TempDir()
	writeTree(t, root, map[string]string{
		"pkg/application/stub_purego.go": "//go:build darwin && purego\n\npackage application\n\nfunc stubbed() bool { return true }\n\nfunc real() bool { return true }\n",
		"pkg/application/stub_cgo.go":    "//go:build darwin && !purego\n\npackage application\n\nimport \"C\"\n\nfunc stubbed() bool { return C.foo() != 0 }\n\nfunc real() bool { return false }\n",
	})
	r, err := scan(root, "darwin", []string{"purego"})
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Stubs) != 1 || r.Stubs[0].Name != "stubbed" {
		t.Errorf("Stubs = %+v, want only stubbed", r.Stubs)
	}
}

func TestScanReportsParseErrors(t *testing.T) {
	root := t.TempDir()
	writeTree(t, root, map[string]string{"pkg/application/bad.go": "package x\nfunc (\n"})
	if _, err := scan(root, "darwin", []string{"purego"}); err == nil {
		t.Error("scan accepted a package file that does not parse")
	}
}

func TestScanRecordsGuardedSendAndUnguardedCall(t *testing.T) {
	root := t.TempDir()
	writeTree(t, root, map[string]string{
		"pkg/application/inspector_purego.go": `//go:build darwin && purego

package application

func enable(w id) {
	if respondsTo(w, "setInspectable:") {
		w.send("setInspectable:", true)
	}
}

func bare(w id) {
	w.send("setTitle:", true)
}
`,
	})
	r, err := scan(root, "darwin", []string{"purego"})
	if err != nil {
		t.Fatal(err)
	}
	opt := map[string]bool{}
	for _, e := range r.Selectors {
		opt[e.Name] = e.Optional
	}
	if !opt["setInspectable:"] {
		t.Error("guarded setInspectable: should be optional")
	}
	if opt["setTitle:"] {
		t.Error("unguarded setTitle: should be required")
	}
}
