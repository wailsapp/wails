package types

// White-box tests for branches unreachable through the participle parser:
// direct calls into unexported helpers with synthetic Param/Library graphs,
// cyclic inheritance chains, and a template-fault seam used by the external
// (types_test) package to exercise the generation error paths.

import (
	"bytes"
	"io"
	"io/fs"
	"testing"
	"time"
)

// ---- template-fault seam (exported for the external types_test package) ----

// SetTemplateSourceForTest swaps the template filesystem and returns a restore
// function. Test-only; lets tests drive renderTemplate's error arms.
func SetTemplateSourceForTest(fsys fs.FS) func() {
	old := templateSource
	templateSource = fsys
	return func() { templateSource = old }
}

// memFile is an in-memory fs.File for the fault filesystem.
type memFile struct {
	name string
	r    *bytes.Reader
}

func (m *memFile) Stat() (fs.FileInfo, error) { return memInfo{m.name, int64(m.r.Len())}, nil }
func (m *memFile) Read(p []byte) (int, error) { return m.r.Read(p) }
func (m *memFile) Close() error               { return nil }

type memInfo struct {
	name string
	size int64
}

func (i memInfo) Name() string       { return i.name }
func (i memInfo) Size() int64         { return i.size }
func (i memInfo) Mode() fs.FileMode   { return 0o444 }
func (i memInfo) ModTime() time.Time  { return time.Time{} }
func (i memInfo) IsDir() bool         { return false }
func (i memInfo) Sys() any            { return nil }

// FaultFS returns an fs.FS backed by the real embedded templates except that
// "templates/<target>" misbehaves per mode: "missing" (read error), "parse"
// (content that fails to parse), or "exec" (content that fails to execute).
func FaultFS(target, mode string) fs.FS { return faultFS{target: target, mode: mode} }

type faultFS struct{ target, mode string }

func (f faultFS) Open(name string) (fs.File, error) {
	if name == "templates/"+f.target {
		switch f.mode {
		case "missing":
			return nil, fs.ErrNotExist
		case "parse":
			return &memFile{name, bytes.NewReader([]byte("{{ end }}"))}, nil
		case "exec":
			return &memFile{name, bytes.NewReader([]byte("{{.FieldThatDoesNotExist}}"))}, nil
		}
	}
	return templates.Open(name)
}

// ---- renderTemplate arms exercised directly ----

func TestRenderTemplateErrors(t *testing.T) {
	cases := []struct {
		mode string
		want string
	}{
		{"missing", "read template"},
		{"parse", "parse template"},
		{"exec", "execute template"},
	}
	for _, tc := range cases {
		t.Run(tc.mode, func(t *testing.T) {
			restore := SetTemplateSourceForTest(FaultFS("com.tmpl", tc.mode))
			defer restore()
			err := renderTemplate("COM", "com.tmpl", struct{ PackageName string }{"x"}, io.Discard)
			if err == nil || !contains(err.Error(), tc.want) {
				t.Errorf("renderTemplate %s = %v, want %q", tc.mode, err, tc.want)
			}
		})
	}
}

func contains(s, sub string) bool { return bytes.Contains([]byte(s), []byte(sub)) }

// ---- simple Capture / helper methods ----

func TestBooleanCapture(t *testing.T) {
	var b Boolean
	if err := b.Capture([]string{"true"}); err != nil || !bool(b) {
		t.Errorf("Capture(true) = (%v, %v), want (true, nil)", bool(b), err)
	}
	if err := b.Capture([]string{"false"}); err != nil || bool(b) {
		t.Errorf("Capture(false) = (%v, %v), want (false, nil)", bool(b), err)
	}
}

func TestInterfaceMethodNameCaptureEmpty(t *testing.T) {
	var n InterfaceMethodName
	if err := n.Capture(nil); err != nil {
		t.Errorf("Capture(nil) = %v, want nil", err)
	}
	if n != "" {
		t.Errorf("name mutated on empty capture: %q", n)
	}
	// remove_ prefix branch (add_ is covered by real IDL).
	if err := n.Capture([]string{"remove_Thing"}); err != nil || n != "RemoveThing" {
		t.Errorf("Capture(remove_Thing) = (%q, %v), want RemoveThing", n, err)
	}
}

func TestPropCaptureEmpty(t *testing.T) {
	var p Prop
	if err := p.Capture(nil); err != nil {
		t.Errorf("Capture(nil) = %v, want nil", err)
	}
	if p != "" {
		t.Errorf("prop mutated on empty capture: %q", p)
	}
}

func TestGetBaseClass(t *testing.T) {
	if got := (&InterfaceDeclaration{BaseClass: "IUnknown"}).GetBaseClass(); got != "" {
		t.Errorf("GetBaseClass(IUnknown) = %q, want empty", got)
	}
	if got := (&InterfaceDeclaration{BaseClass: "ICoreWebView2"}).GetBaseClass(); got != "ICoreWebView2" {
		t.Errorf("GetBaseClass = %q, want ICoreWebView2", got)
	}
}

// ---- ResolveGoType: every mapping arm ----

func TestResolveGoType(t *testing.T) {
	cases := []struct{ idl, ptr, dir, want string }{
		{"LPCWSTR", "*", "in", "[]string"}, // [in] array of strings
		{"LPWSTR", "*", "out", "string"},   // [out] single string
		{"LPWSTR", "", "in", "string"},
		{"HRESULT", "", "in", "uintptr"},
		{"UINT64", "", "in", "uint64"},
		{"UINT32", "", "in", "uint32"},
		{"DWORD", "", "in", "uint32"},
		{"UINT", "", "in", "uint"},
		{"INT64", "", "in", "int64"},
		{"INT32", "", "in", "int32"},
		{"INT", "", "in", "int"},
		{"BOOL", "", "in", "bool"},
		{"BYTE", "", "in", "uint8"},
		{"double", "", "in", "float64"},
		{"IUnknown", "", "in", "IUnknown"},
		{"EventRegistrationToken", "", "in", "EventRegistrationToken"},
		{"ICoreWebView2Settings", "*", "out", "ICoreWebView2Settings"}, // pass-through
	}
	for _, c := range cases {
		if got := ResolveGoType(c.idl, c.ptr, c.dir); got != c.want {
			t.Errorf("ResolveGoType(%q,%q,%q) = %q, want %q", c.idl, c.ptr, c.dir, got, c.want)
		}
	}
}

// ---- synthetic graph helpers ----

func newLibrary() *Library {
	l := &Library{Name: "Lib"}
	l.packageName = "lib"
	l.structSizes = map[string]int{}
	l.interfaceBases = map[string]string{}
	return l
}

// wire attaches a param to a fresh method/interface/declaration chain rooted at l.
func wire(l *Library, p *Param) *Param {
	p.decl = &InterfaceMethod{decl: &InterfaceDeclaration{decl: &Declaration{library: l}}}
	return p
}

// ---- struct sizeOf: unknown field type ----

func TestSizeOfUnknownField(t *testing.T) {
	d := &StructDeclaration{Name: "Bad", Fields: []*StructField{{Type: "FLOAT", Name: "x"}}}
	if _, err := d.sizeOf(); err == nil {
		t.Error("sizeOf with unknown field type should error")
	}
}

// ---- Declaration.Process: unrecognised declaration ----

func TestDeclarationProcessUnknown(t *testing.T) {
	l := newLibrary()
	d := &Declaration{} // every variant nil/empty
	if err := d.Process(l); err == nil {
		t.Error("Process of an empty declaration should error")
	}
}

// ---- CountVariableName: no paired count output ----

func TestCountVariableNameNoPair(t *testing.T) {
	l := newLibrary()
	arr := wire(l, &Param{Name: "items", GoType: "[]*IThing"})
	// Its method has no uint32 sibling.
	arr.decl.outputParams = []*Param{arr}
	if _, err := arr.CountVariableName(); err == nil {
		t.Error("CountVariableName without a uint32 pair should error")
	}

	// And the happy path: a uint32 sibling resolves.
	count := wire(l, &Param{Name: "count", GoType: "uint32"})
	arr.decl.outputParams = []*Param{arr, count}
	count.decl = arr.decl
	if name, err := arr.CountVariableName(); err != nil || name != "count" {
		t.Errorf("CountVariableName = (%q, %v), want count", name, err)
	}
}

// ---- processVtableCallInput: the error and size arms ----

func TestProcessVtableCallInputErrors(t *testing.T) {
	l := newLibrary()

	// [in] triple pointer is unsupported.
	tri := wire(l, &Param{Name: "x", Type: "IThing", Pointer: "***", GoType: "IThing"})
	if err := tri.processVtableCallInput(); err == nil {
		t.Error("[in] triple pointer should be unsupported")
	}

	// Unknown by-value type (not builtin, not a declared struct).
	unk := wire(l, &Param{Name: "x", Type: "Mystery", GoType: "Mystery"})
	if err := unk.processVtableCallInput(); err == nil {
		t.Error("unknown by-value type should error")
	}

	// A declared struct of an unsupported size (12 bytes) hits the size default.
	l.structSizes["Odd"] = 12
	odd := wire(l, &Param{Name: "x", Type: "Odd", GoType: "Odd"})
	if err := odd.processVtableCallInput(); err == nil {
		t.Error("by-value struct of unsupported size should error")
	}
}

// ---- methodInputArg: type strategies not expressible in the bundled IDL ----

func TestMethodInputArgStrategies(t *testing.T) {
	l := newLibrary()
	l.enums.Add("MyEnum")
	imports := map[string]bool{}

	ok := []*Param{
		wire(l, &Param{Name: "s", GoType: "string"}),
		wire(l, &Param{Name: "ss", GoType: "[]string"}),
		wire(l, &Param{Name: "is", GoType: "[]*IThing"}),
		wire(l, &Param{Name: "b", GoType: "bool"}),
		wire(l, &Param{Name: "f64", GoType: "float64"}),
		wire(l, &Param{Name: "f32", GoType: "float32"}),
		wire(l, &Param{Name: "i64", GoType: "int64"}),
		wire(l, &Param{Name: "u64", GoType: "uint64"}),
		wire(l, &Param{Name: "tok", GoType: "EventRegistrationToken"}),
		wire(l, &Param{Name: "col", GoType: "COREWEBVIEW2_COLOR"}),
		wire(l, &Param{Name: "pt", GoType: "POINT"}),
		wire(l, &Param{Name: "rc", GoType: "RECT"}),
		wire(l, &Param{Name: "en", Type: "MyEnum", GoType: "MyEnum"}),
		wire(l, &Param{Name: "h", GoType: "HWND"}),
		wire(l, &Param{Name: "n", GoType: "uint32"}),
		wire(l, &Param{Name: "ptr", GoType: "IThing", Pointer: "*"}),
	}
	for i, p := range ok {
		if _, err := methodInputArg(p, i, imports); err != nil {
			t.Errorf("methodInputArg(%s) unexpected error: %v", p.Name, err)
		}
	}

	// No strategy → error.
	bad := wire(l, &Param{Name: "weird", Type: "Weird", GoType: "Weird"})
	if _, err := methodInputArg(bad, 0, imports); err == nil {
		t.Error("methodInputArg should error for an unsupported type")
	}
}

// ---- methodOutputArg: type strategies not expressible in the bundled IDL ----

func TestMethodOutputArgStrategies(t *testing.T) {
	l := newLibrary()
	l.enums.Add("MyEnum")
	imports := map[string]bool{}
	none := map[*Param]int{}

	simple := []*Param{
		wire(l, &Param{Name: "b", GoType: "bool"}),
		wire(l, &Param{Name: "s", GoType: "string"}),
		wire(l, &Param{Name: "ss", GoType: "[]string"}),
		wire(l, &Param{Name: "is", GoType: "[]*IThing"}),
		wire(l, &Param{Name: "f64", GoType: "float64"}),
		wire(l, &Param{Name: "i64", GoType: "int64"}),
		wire(l, &Param{Name: "u64", GoType: "uint64"}),
		wire(l, &Param{Name: "tok", GoType: "EventRegistrationToken"}),
		wire(l, &Param{Name: "u16", GoType: "uint16"}),
		wire(l, &Param{Name: "i16", GoType: "int16"}),
		wire(l, &Param{Name: "u8", GoType: "uint8"}),
		wire(l, &Param{Name: "i8", GoType: "int8"}),
		wire(l, &Param{Name: "n", GoType: "uint32"}),
		wire(l, &Param{Name: "h", GoType: "HWND"}),
		wire(l, &Param{Name: "ptr", GoType: "*IThing", OutputGoType: "*IThing"}),
		wire(l, &Param{Name: "pod", GoType: "POINT", OutputGoType: "POINT"}), // struct default
	}
	for i, p := range simple {
		if got := methodOutputArg(p, i, none, imports); got == nil {
			t.Errorf("methodOutputArg(%s) returned nil", p.Name)
		}
	}

	// The array-count branch: a uint32 that is the count for some array param.
	count := wire(l, &Param{Name: "count", GoType: "uint32"})
	if got := methodOutputArg(count, 0, map[*Param]int{count: 2}, imports); got == nil {
		t.Error("methodOutputArg(count) returned nil")
	}
}

// TestErrorValuesHRESULTNonHRESULT covers the non-HRESULT branch, which the
// method template only reaches for HRESULT-returning methods.
func TestErrorValuesHRESULTNonHRESULT(t *testing.T) {
	m := &InterfaceMethod{ReturnType: "void"}
	if got := m.ErrorValuesHRESULT(); got != "err" {
		t.Errorf("ErrorValuesHRESULT(void) = %q, want err", got)
	}
}

// TestProcessVtableCallInputFloat32 covers the float32 marshalling arm; no IDL
// type maps to float32, so it is only reachable by a direct call.
func TestProcessVtableCallInputFloat32(t *testing.T) {
	l := newLibrary()
	p := wire(l, &Param{Name: "f", GoType: "float32"})
	if err := p.processVtableCallInput(); err != nil {
		t.Fatalf("processVtableCallInput(float32): %v", err)
	}
	if len(p.callWordsAMD64) != 1 {
		t.Errorf("float32 should marshal to one word, got %v", p.callWordsAMD64)
	}
}

// ---- cyclic inheritance: chainRoot / generateVtbl / emitQIHelperTest / ABI ----

func cyclicLibrary() *Library {
	l := newLibrary()
	l.interfaceBases = map[string]string{"A": "B", "B": "A"} // neither is IUnknown
	return l
}

func TestChainRootCycle(t *testing.T) {
	if _, err := cyclicLibrary().chainRoot("A"); err == nil {
		t.Error("chainRoot on a cyclic chain should error")
	}
}

func TestGenerateVtblChainRootError(t *testing.T) {
	l := cyclicLibrary()
	d := &InterfaceDeclaration{
		Name:      "A",
		BaseClass: "B", // present in interfaceBases, so it passes the declared-base check
		decl:      &Declaration{library: l},
	}
	if err := d.generateVtbl("lib", io.Discard); err == nil {
		t.Error("generateVtbl should surface the chainRoot cycle error")
	}
}

func TestEmitQIHelperTestChainRootError(t *testing.T) {
	l := cyclicLibrary()
	d := &InterfaceDeclaration{Name: "A", BaseClass: "B"}
	var buf bytes.Buffer
	if err := emitQIHelperTest(&buf, d, l); err == nil {
		t.Error("emitQIHelperTest should surface the chainRoot cycle error")
	}
}

func TestGenerateABITestFileTooDeep(t *testing.T) {
	l := cyclicLibrary()
	l.Declarations = []*Declaration{{Interface: &InterfaceDeclaration{Name: "A", BaseClass: "B"}}}
	if _, err := l.generateABITestFile(); err == nil {
		t.Error("generateABITestFile should error on a cyclic inheritance chain")
	}
}

// TestLibraryProcessSizeOfError covers Library.Process surfacing a struct
// sizing error (the field type is rejected by sizeOf, not the parser).
func TestLibraryProcessSizeOfError(t *testing.T) {
	l := newLibrary()
	l.Declarations = []*Declaration{
		{Struct: &StructDeclaration{Name: "Bad", Fields: []*StructField{{Type: "FLOAT", Name: "x"}}}},
	}
	if err := l.Process(); err == nil {
		t.Error("Library.Process should surface a struct sizing error")
	}
}

// TestGenerateTestFileQIError covers generateTestFile returning the QI-helper
// error when the inheritance chain is cyclic.
func TestGenerateTestFileQIError(t *testing.T) {
	l := cyclicLibrary()
	d := &InterfaceDeclaration{Name: "A", BaseClass: "B"} // derived → emits a QI helper
	if _, err := d.generateTestFile(l); err == nil {
		t.Error("generateTestFile should surface the QI-helper chainRoot error")
	}
}

// TestGenerateTestsABIError covers generateTests returning the ABI-file error.
// The interface's struct field BaseClass is IUnknown (so per-interface test
// generation is skipped) while interfaceBases encodes a cycle the ABI walk hits.
func TestGenerateTestsABIError(t *testing.T) {
	l := cyclicLibrary()
	l.Declarations = []*Declaration{{Interface: &InterfaceDeclaration{Name: "A", BaseClass: "IUnknown"}}}
	if _, err := l.generateTests(); err == nil {
		t.Error("generateTests should surface the ABI-file cycle error")
	}
}

// ---- IDL.Generate with no libraries returns (nil, nil) ----

func TestIDLGenerateNoLibraries(t *testing.T) {
	files, err := (&IDL{}).Generate()
	if err != nil || files != nil {
		t.Errorf("empty IDL.Generate = (%v, %v), want (nil, nil)", files, err)
	}
	tests, err := (&IDL{}).GenerateTests()
	if err != nil || tests != nil {
		t.Errorf("empty IDL.GenerateTests = (%v, %v), want (nil, nil)", tests, err)
	}
}
