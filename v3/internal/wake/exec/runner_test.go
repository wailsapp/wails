package exec

import (
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

// Precondition `sh:` strings are Go templates and must be expanded against the
// task's resolved vars before being run as shell commands. The build Taskfiles
// guard the garble check with `{{if eq .OBFUSCATED "true"}}...{{else}}true{{end}}`;
// without expansion the raw template reaches the shell, fails to parse, and the
// precondition's `msg:` surfaces as a spurious "garble is required" error on
// every build. See the obfuscation false-positive regression.
func TestCheckPreconditionsExpandsTemplates(t *testing.T) {
	const guard = `{{if eq .OBFUSCATED "true"}}command -v definitely-not-a-real-binary >/dev/null 2>&1{{else}}true{{end}}`

	newTask := func() *ast.Task {
		return &ast.Task{
			Precondition: []*ast.Precondition{{Sh: guard, Msg: "garble is required for obfuscated builds"}},
		}
	}

	t.Run("unset var resolves to true and passes", func(t *testing.T) {
		if err := checkPreconditions(newTask(), map[string]*ast.Var{}); err != nil {
			t.Fatalf("expected precondition to pass when OBFUSCATED is unset, got: %v", err)
		}
	})

	t.Run("non-true var resolves to true and passes", func(t *testing.T) {
		vars := map[string]*ast.Var{"OBFUSCATED": {Value: "false"}}
		if err := checkPreconditions(newTask(), vars); err != nil {
			t.Fatalf("expected precondition to pass when OBFUSCATED is false, got: %v", err)
		}
	})

	t.Run("true var runs the guard and reports msg on failure", func(t *testing.T) {
		vars := map[string]*ast.Var{"OBFUSCATED": {Value: "true"}}
		err := checkPreconditions(newTask(), vars)
		if err == nil {
			t.Fatal("expected precondition to fail when OBFUSCATED is true and the tool is missing")
		}
		if got := err.Error(); got == "" || !contains(got, "garble is required") {
			t.Fatalf("expected garble msg, got: %v", err)
		}
	})
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
