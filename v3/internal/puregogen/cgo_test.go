package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRewriteCgoSource(t *testing.T) {
	dir := t.TempDir()
	source, output := filepath.Join(dir, "reference.go"), filepath.Join(dir, "generated.go")
	input := `package native
/* void native_init(void); */
import "C"
func init() { C.native_init() }
func callback(value C.int, data C.gpointer) { if value == C.FLAG { C.native_init() } }
func atom() C.int { return C.ATOM }
`
	if err := os.WriteFile(source, []byte(input), 0600); err != nil {
		t.Fatal(err)
	}
	names := map[string]string{"native_init": "_Cfunc_native_init", "int": "_Ctype_int", "gpointer": "_Ctype_gpointer", "FLAG": "_Ciconst_FLAG", "ATOM": "_Cmacro_ATOM"}
	macros := map[string]bool{"ATOM": true}
	generate := func(check bool) error {
		return rewriteCgoSource(source, output, "linux && purego", names, macros, []string{"callback"}, map[string]string{"atom": "requireNative"}, check)
	}
	if err := generate(false); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"ensureGTK3Libraries()", "_Cfunc_native_init()", "callback(value _Ctype_int, data uintptr)", "value == _Ciconst_FLAG", "return _Cmacro_ATOM()", "func atom() _Ctype_int { requireNative(); return _Cmacro_ATOM() }"} {
		if !strings.Contains(string(data), want) {
			t.Errorf("missing %q in generated source:\n%s", want, data)
		}
	}
	if strings.Contains(string(data), `import "C"`) || strings.Contains(string(data), "void native_init") {
		t.Fatal("generated source retains C dependency")
	}
	if err := generate(true); err != nil {
		t.Fatalf("fresh output fails drift check: %v", err)
	}
	if err := os.WriteFile(output, append(data, '\n'), 0600); err != nil {
		t.Fatal(err)
	}
	if err := generate(true); err == nil {
		t.Fatal("stale output passes drift check")
	}
	delete(names, "native_init")
	if err := generate(false); err == nil {
		t.Fatal("unknown native call silently accepted")
	}
}
