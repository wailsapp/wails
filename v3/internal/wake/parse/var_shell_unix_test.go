//go:build !windows

package parse

import (
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

func TestVariableShellFailureIncludesDiagnostic(t *testing.T) {
	const diagnostic = "toolchain helper is unavailable"
	t.Setenv("WAKE_VARIABLE_DIAGNOSTIC", diagnostic)
	variable := &ast.Var{Shell: `printf '%s\n' "$WAKE_VARIABLE_DIAGNOSTIC" >&2; exit 7`}
	err := ResolveVarShell(variable)
	if err == nil || !strings.Contains(err.Error(), diagnostic) {
		t.Errorf("error = %v, want captured stderr diagnostic", err)
	}
	var exitError *exec.ExitError
	if !errors.As(err, &exitError) || exitError.ExitCode() != 7 {
		t.Errorf("error = %v, want wrapped exit status 7", err)
	}
	if variable.Value != "" {
		t.Errorf("failed command assigned variable value %q", variable.Value)
	}
}

func TestVariableShellSuccessUsesOnlyStdout(t *testing.T) {
	t.Setenv("WAKE_VARIABLE_VALUE", "expected")
	t.Setenv("WAKE_VARIABLE_DIAGNOSTIC", "successful command warning")
	variable := &ast.Var{Shell: `printf '%s\n' "$WAKE_VARIABLE_VALUE"; printf '%s\n' "$WAKE_VARIABLE_DIAGNOSTIC" >&2`}
	if err := ResolveAllVarShells(map[string]*ast.Var{"VALUE": variable}); err != nil {
		t.Fatal(err)
	}
	if variable.Value != "expected" {
		t.Errorf("variable value = %q, want stdout without diagnostic", variable.Value)
	}
}
