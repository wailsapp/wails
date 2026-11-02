package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func body(t *testing.T, src string) *ast.BlockStmt {
	t.Helper()
	f, err := parser.ParseFile(token.NewFileSet(), "x.go", "package x\n"+src, 0)
	if err != nil {
		t.Fatal(err)
	}
	return f.Decls[0].(*ast.FuncDecl).Body
}

func TestTrivial(t *testing.T) {
	for src, want := range map[string]bool{
		"func f() {}":                        true,
		"func f() bool { return false }":     true,
		"func f() error { return nil }":      true,
		"func f() { panic(\"todo\") }":       true,
		"func f() T { return T{} }":          true,
		"func f() { C.flash(w) }":            false,
		"func f() int { return g() }":        false,
		"func f() { if x { return } }":       false,
		"func f() []int { return []int{1} }": false,
	} {
		if got := trivial(body(t, src)); got != want {
			t.Errorf("trivial(%q) = %v, want %v", src, got, want)
		}
	}
}

func TestNewer(t *testing.T) {
	for _, c := range []struct {
		a, b string
		want bool
	}{{"4.14", "4.8", true}, {"4.8", "4.14", false}, {"2.40", "2.40", false}, {"3.24.1", "3.24", true}} {
		if got := newer(c.a, c.b); got != c.want {
			t.Errorf("newer(%s, %s) = %v", c.a, c.b, got)
		}
	}
}

// collectSrc runs the collector over one source file as if it were a package
// file built for goos.
func collectSrc(t *testing.T, goos, src string) *collector {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "x.go", "package x\n"+src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	c := &collector{
		fset: fset, rel: "x.go", optional: map[string]bool{}, r: &report{GOOS: goos},
		symbols: map[string]*entry{}, classes: map[string]*entry{}, selectors: map[string]*entry{},
	}
	c.file(f)
	return c
}

// A send guarded by respondsTo in the same function is optional, as the
// cgo-era selectors newer than the floor must be. A bare send stays required,
// and a guard only covers its own selector and its own function.
func TestRespondsToGuardedSendIsOptional(t *testing.T) {
	for _, c := range []struct {
		name, sel, src string
		optional       bool
	}{
		{"guarded in the same function", "setInspectable:", `
func f(w id) {
	if respondsTo(w, "setInspectable:") {
		w.send("setInspectable:", true)
	}
}`, true},
		{"bare send", "setInspectable:", `
func f(w id) {
	w.send("setInspectable:", true)
}`, false},
		{"guard in another function does not count", "setInspectable:", `
func g(w id) bool { return respondsTo(w, "setInspectable:") }
func f(w id) { w.send("setInspectable:", true) }`, false},
		{"guard names another selector", "setInspectable:", `
func f(w id) {
	if respondsTo(w, "setTintColor:") {
		w.send("setInspectable:", true)
	}
}`, false},
		{"one bare send elsewhere makes it required", "a:", `
func f(w id) { if respondsTo(w, "a:") { w.send("a:") } }
func g(w id) { w.send("a:") }`, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			e := collectSrc(t, "darwin", c.src).selectors[c.sel]
			if e == nil {
				t.Fatalf("selector %s not recorded", c.sel)
			}
			if e.Optional != c.optional {
				t.Errorf("%s: Optional = %v, want %v", c.sel, e.Optional, c.optional)
			}
		})
	}
}

func TestGuarded(t *testing.T) {
	ok := body(t, `func f() float64 {
	if _native_gdk_monitor_get_scale != nil {
		return gdk_monitor_get_scale(m)
	}
	return 1
}`)
	if !guarded(ok, "gdk_monitor_get_scale") {
		t.Error("nil-checked call reported unguarded")
	}
	bad := body(t, `func f() float64 { return gdk_monitor_get_scale(m) }`)
	if guarded(bad, "gdk_monitor_get_scale") {
		t.Error("bare call reported guarded")
	}
	version := body(t, `func f() { if gtkAtLeast(4, 10) { gtk_file_dialog_open(d) } }`)
	if !guarded(version, "gtk_file_dialog_open") {
		t.Error("version-checked call reported unguarded")
	}
}

// Synthetic darwin purego sources for the stacknarrow rule. The prelude
// defines the backend's helper shapes so the argument expressions type-check
// without the real package.
const stackPrelude = `
type id uintptr
type NSPoint struct{ X, Y float64 }

func class(string) id { return 0 }
func (o id) send(sel string, args ...any) id { return o }
func sendSuperFrom(definer id, self id, sel string, args ...any) id { return 0 }
func get[T any](o id, sel string, args ...any) T { var z T; return z }
func register(fn any, lib int, name string) {}
`

// stackScan runs the darwin scan over one synthetic purego source file.
func stackScan(t *testing.T, body string) *report {
	t.Helper()
	root := t.TempDir()
	writeTree(t, root, map[string]string{
		"pkg/application/x_purego.go": "//go:build darwin && purego\n\npackage application\n" + stackPrelude + body,
	})
	r, err := scan(root, "darwin", []string{"purego"})
	if err != nil {
		t.Fatal(err)
	}
	return r
}

const keyEventSelector = "keyEventWithType:location:modifierFlags:timestamp:windowNumber:context:characters:charactersIgnoringModifiers:isARepeat:keyCode:"

func TestStackNarrowKeyEventSendIsFlagged(t *testing.T) {
	// self, _cmd, then eventType..charactersIgnoringModifiers fill x2..x7 and
	// the NSPoint takes v0. isARepeat (C parameter 11, BOOL) and keyCode
	// (C parameter 12, unsigned short) spill to the stack.
	r := stackScan(t, `
func keyEvent(o id) {
	o.send("`+keyEventSelector+`",
		uint(10), NSPoint{}, uint(0), float64(1), 0, id(0), id(0), id(0), true, uint16(13))
}`)
	if len(r.StackNarrow) != 2 {
		t.Fatalf("findings = %d, want 2: %+v", len(r.StackNarrow), r.StackNarrow)
	}
	for i, want := range []string{"param 11 (bool)", "param 12 (uint16)"} {
		if !strings.Contains(r.StackNarrow[i].Text, want) || r.StackNarrow[i].Name != keyEventSelector {
			t.Errorf("finding %d = %+v, want %s on the key event selector", i, r.StackNarrow[i], want)
		}
	}
	if r.StackNarrow[0].Where == "" || !strings.Contains(r.StackNarrow[0].Where, "x_purego.go:") {
		t.Errorf("Where = %q, want the file and line of the call", r.StackNarrow[0].Where)
	}
}

func TestStackNarrowKeyEventSendWithWideArgsIsNotFlagged(t *testing.T) {
	r := stackScan(t, `
func keyEvent(o id) {
	o.send("`+keyEventSelector+`",
		uint(10), NSPoint{}, uint(0), float64(1), 0, id(0), id(0), id(0), uint64(1), uint64(13))
}`)
	if len(r.StackNarrow) != 0 {
		t.Errorf("findings = %+v, want none", r.StackNarrow)
	}
	if r.StackChecked != 1 {
		t.Errorf("StackChecked = %d, want 1", r.StackChecked)
	}
}

func TestStackNarrowSendForms(t *testing.T) {
	// Every form passes self and _cmd in x0/x1: six 8-byte arguments fill
	// x2..x7, the next two go to the stack without narrowing, and the bool
	// after them spills at C parameter 11.
	eight := "uint64(1), uint64(2), uint64(3), uint64(4), uint64(5), uint64(6), uint64(7), uint64(8), true"
	for _, c := range []struct {
		name, call string
		want       int
	}{
		{"method send", `o.send("x:", ` + eight + `)`, 1},
		{"sendSuperFrom", `sendSuperFrom(o, o, "x:", ` + eight + `)`, 1},
		{"generic get", `get[bool](o, "x:", ` + eight + `)`, 1},
		{"narrow argument in a register is not flagged", `o.send("x:", uint16(1))`, 0},
		{"spread arguments are not classified", `o.send("x:", args...)`, 0},
		{"interface argument is not classified", `o.send("x:", v)`, 0},
	} {
		t.Run(c.name, func(t *testing.T) {
			r := stackScan(t, `
func f(o id, args []any, v any) {
	`+c.call+`
}`)
			if len(r.StackNarrow) != c.want {
				t.Errorf("findings = %+v, want %d", r.StackNarrow, c.want)
			}
			if c.want == 1 && r.StackNarrow[0].Name != "x:" {
				t.Errorf("name = %q, want the selector x:", r.StackNarrow[0].Name)
			}
		})
	}
}

func TestStackNarrowSendSelectorNotLiteral(t *testing.T) {
	r := stackScan(t, `
func f(o id, sel string) { o.send(sel, true, true, true, true, true, true, true, true, true) }`)
	if len(r.StackNarrow) == 0 || r.StackNarrow[0].Name != "<dynamic>" {
		t.Errorf("findings = %+v, want <dynamic> name", r.StackNarrow)
	}
}

// A C function with a fixed signature is placed by purego's Apple bundling
// (struct_arm64.go), so its narrow spill is counted as bundled, not reported.
func TestStackNarrowFixedSignatureIsBundledNotFlagged(t *testing.T) {
	r := stackScan(t, `
var cfn func(a1, a2, a3, a4, a5, a6, a7, a8 uintptr, a9 uint16, a10 uintptr) int32

func init() { register(&cfn, 0, "cfn") }

func call() { cfn(1, 2, 3, 4, 5, 6, 7, 8, 9, 10) }`)
	if len(r.StackNarrow) != 0 {
		t.Errorf("findings = %+v, want none", r.StackNarrow)
	}
	if r.StackBundled != 1 {
		t.Errorf("StackBundled = %d, want 1 (the uint16 at parameter 9)", r.StackBundled)
	}
}

func TestStackNarrowFixedSignaturePointersNotCounted(t *testing.T) {
	r := stackScan(t, `
var cfn func(a1, a2, a3, a4, a5, a6, a7, a8, a9, a10 uintptr) int32

func init() { register(&cfn, 0, "cfn") }

func call() { cfn(1, 2, 3, 4, 5, 6, 7, 8, 9, 10) }`)
	if len(r.StackNarrow) != 0 || r.StackBundled != 0 {
		t.Errorf("findings = %+v, bundled = %d, want none", r.StackNarrow, r.StackBundled)
	}
}

func TestStackNarrowVariadicFixedCallIsNotClassified(t *testing.T) {
	r := stackScan(t, `
var vfn func(a ...any) int32

func init() { register(&vfn, 0, "vfn") }

func call() { vfn(true, uint16(1)) }`)
	if r.StackUnknown == 0 {
		t.Errorf("StackUnknown = 0, want the variadic call counted as not classified")
	}
}

func TestStackNarrowAllowFileSuppresses(t *testing.T) {
	r := stackScan(t, `
func keyEvent(o id) {
	o.send("`+keyEventSelector+`", uint(10), NSPoint{}, uint(0), float64(1), 0, id(0), id(0), id(0), true, uint16(13))
}`)
	allow := filepath.Join(t.TempDir(), "allow.txt")
	if err := os.WriteFile(allow, []byte("# intentional\n"+keyEventSelector+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	applyAllow(r, allow)
	if len(r.StackNarrow) != 0 || len(r.AllowedStack) != 2 {
		t.Errorf("after allow: findings = %d, allowed = %d, want 0 and 2", len(r.StackNarrow), len(r.AllowedStack))
	}
}

func TestStackNarrowRunExitStatus(t *testing.T) {
	root := t.TempDir()
	writeTree(t, root, map[string]string{
		"pkg/application/x_purego.go": "//go:build darwin && purego\n\npackage application\n" + stackPrelude +
			"\nfunc keyEvent(o id) {\n\to.send(\"" + keyEventSelector + "\", uint(10), NSPoint{}, uint(0), float64(1), 0, id(0), id(0), id(0), true, uint16(13))\n}\n",
	})
	if got := run([]string{"-root", root, "-goos", "darwin", "-tags", "purego", "static"}); got != 1 {
		t.Errorf("run = %d, want 1 for a stacknarrow finding", got)
	}
}

// Calls that pass a narrow value to a fixed-signature, non-variadic C function
// on linux are not touched: the rule is darwin only.
func TestStackNarrowLinuxIsNotChecked(t *testing.T) {
	root := t.TempDir()
	writeTree(t, root, map[string]string{
		"pkg/application/x_purego.go": "//go:build linux && purego\n\npackage application\n" + stackPrelude +
			"\nfunc keyEvent(o id) { o.send(\"" + keyEventSelector + "\", uint(10), NSPoint{}, uint(0), float64(1), 0, id(0), id(0), id(0), true, uint16(13)) }\n",
	})
	r, err := scan(root, "linux", []string{"purego"})
	if err != nil {
		t.Fatal(err)
	}
	if len(r.StackNarrow) != 0 || r.StackChecked != 0 {
		t.Errorf("linux: findings = %d, checked = %d, want 0 and 0", len(r.StackNarrow), r.StackChecked)
	}
}

func TestAllocateAAPCS64(t *testing.T) {
	in := func(n int) cArg { return cArg{kind: kindInt, size: 8} }
	fl := cArg{kind: kindFloat, size: 8}
	for _, c := range []struct {
		name    string
		args    []cArg
		preInt  int
		spilled []int
	}{
		{"eight integers fit", repeat(in(0), 8), 0, nil},
		{"ninth integer spills", repeat(in(0), 9), 0, []int{8}},
		{"self and _cmd leave six", repeat(in(0), 7), 2, []int{6}},
		{"composite that exactly fits the last integer registers", append(repeat(in(0), 6), cArg{kind: kindSmall, n: 2}), 0, nil},
		{"HFA that exactly fits the last float registers", append(repeat(fl, 6), cArg{kind: kindHFA, n: 2}), 0, nil},
		{"composite that does not fit exhausts integers", append(repeat(in(0), 7), cArg{kind: kindSmall, n: 2}, in(0)), 0, []int{7, 8}},
		{"HFA that does not fit exhausts floats", append(repeat(fl, 7), cArg{kind: kindHFA, n: 2}, fl), 0, []int{7, 8}},
		{"HFA of four fits eight floats", repeat(fl, 4), 0, nil},
		{"floats and integers use separate registers", append(repeat(fl, 8), repeat(in(0), 8)...), 0, nil},
		{"zero-sized takes nothing", []cArg{{kind: kindZero}, in(0)}, 0, nil},
	} {
		t.Run(c.name, func(t *testing.T) {
			got := allocate(c.args, c.preInt)
			if !equalInts(got, c.spilled) {
				t.Errorf("allocate = %v, want %v", got, c.spilled)
			}
		})
	}
}

func repeat(a cArg, n int) []cArg {
	out := make([]cArg, n)
	for i := range out {
		out[i] = a
	}
	return out
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestClassifyGoTypes(t *testing.T) {
	sizes := types.SizesFor("gc", "arm64")
	f64 := types.Typ[types.Float64]
	point := types.NewStruct([]*types.Var{types.NewField(0, nil, "X", f64, false), types.NewField(0, nil, "Y", f64, false)}, nil)
	u16 := types.NewStruct([]*types.Var{types.NewField(0, nil, "A", types.Typ[types.Uint16], false), types.NewField(0, nil, "B", types.Typ[types.Uint16], false)}, nil)
	big := types.NewArray(f64, 5)
	for _, c := range []struct {
		name string
		t    types.Type
		want cArg
		ok   bool
	}{
		{"bool", types.Typ[types.Bool], cArg{kind: kindInt, size: 1, narrow: true}, true},
		{"uint16", types.Typ[types.Uint16], cArg{kind: kindInt, size: 2, narrow: true}, true},
		{"uintptr", types.Typ[types.Uintptr], cArg{kind: kindInt, size: 8}, true},
		{"float32", types.Typ[types.Float32], cArg{kind: kindFloat, size: 4, narrow: true}, true},
		{"float64", f64, cArg{kind: kindFloat, size: 8}, true},
		{"string is a pointer", types.Typ[types.String], cArg{kind: kindInt, size: 8}, true},
		{"pointer", types.NewPointer(types.Typ[types.Int8]), cArg{kind: kindInt, size: 8}, true},
		{"NSPoint is an HFA", point, cArg{kind: kindHFA, n: 2, size: 16}, true},
		{"uint16 pair is a narrow small composite", u16, cArg{kind: kindSmall, n: 1, size: 4, narrow: true}, true},
		{"large composite is by reference", big, cArg{kind: kindInt, size: 8}, true},
		{"empty struct", types.NewStruct(nil, nil), cArg{kind: kindZero}, true},
		{"interface is not mapped", types.NewInterfaceType(nil, nil), cArg{}, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			got, ok := classify(c.t, sizes)
			if ok != c.ok || got != c.want {
				t.Errorf("classify = %+v, %v; want %+v, %v", got, ok, c.want, c.ok)
			}
		})
	}
}

// A send helper passing its own args... on is a forward: its arguments are
// classified at each call of the helper, so it is neither skipped nor checked.
// Spreading any other slice is still not classified.
func TestStackNarrowHelperForwardIsCounted(t *testing.T) {
	r := stackScan(t, `
func (o id) sendAll(sel string, args ...any) id { return o.send(sel, args...) }

func get(o id, sel string, args ...any) id { return o.send(sel, args...) }

func other(o id, sel string, args ...any) id { return o.send(sel, args...) }

func spread(o id, sel string, rest []any) id { return o.send(sel, rest...) }`)
	if r.StackForwarded != 1 {
		t.Errorf("StackForwarded = %d, want 1 (get; sendAll and other are not send helpers)", r.StackForwarded)
	}
	if r.StackUnknown != 3 || len(r.StackSkipped) != 3 {
		t.Errorf("StackUnknown = %d, skipped %v, want 3", r.StackUnknown, r.StackSkipped)
	}
}

// An autorelease pool pushed without pinning the goroutine to its thread is
// reported; the same pool after runtime.LockOSThread is not.
func TestPoolThreadRule(t *testing.T) {
	r := stackScan(t, `
func bad(c func(string) id) { c("NSAutoreleasePool") }

func good(c func(string) id) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	c("NSAutoreleasePool")
}`)
	if len(r.PoolThread) != 1 || !strings.HasSuffix(r.PoolThread[0], " bad") {
		t.Errorf("PoolThread = %v, want only bad", r.PoolThread)
	}
}
