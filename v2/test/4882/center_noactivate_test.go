package test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"
)

// formSource is the winc Form implementation that defines Center().
const formSource = "../../internal/frontend/desktop/windows/winc/form.go"

// TestCenterSetsNoActivate is the regression guard for #4882: launching an app
// with StartHidden:true must not steal focus from the foreground window.
// Form.Center() prevents that by passing SWP_NOACTIVATE to SetWindowPos.
//
// winc is Windows-only (build-tagged), so instead of linking it we assert
// against the source of Center() directly: this fails if the flag is ever
// dropped from that call, which is the exact regression that reintroduces the
// bug. (The previous version of this test only OR'd together local copies of
// the constants and would pass even if form.go reverted — see PR review.)
func TestCenterSetsNoActivate(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filepath.FromSlash(formSource), nil, 0)
	if err != nil {
		t.Fatalf("parsing %s: %v", formSource, err)
	}

	var center *ast.FuncDecl
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Recv != nil && fn.Name.Name == "Center" {
			center = fn
			break
		}
	}
	if center == nil {
		t.Fatalf("could not find Center() in %s", formSource)
	}

	var foundSetWindowPos, foundNoActivate bool
	ast.Inspect(center, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "SetWindowPos" {
			return true
		}
		foundSetWindowPos = true
		for _, arg := range call.Args {
			if exprReferences(arg, "SWP_NOACTIVATE") {
				foundNoActivate = true
			}
		}
		return true
	})

	if !foundSetWindowPos {
		t.Fatal("Center() no longer calls SetWindowPos; update this regression guard")
	}
	if !foundNoActivate {
		t.Error("Center()'s SetWindowPos call is missing SWP_NOACTIVATE — this reintroduces the StartHidden focus-steal bug (#4882)")
	}
}

// exprReferences reports whether expr references an identifier with the given
// name, e.g. the w32.SWP_NOACTIVATE term inside a bitwise-or of flags.
func exprReferences(expr ast.Expr, name string) bool {
	found := false
	ast.Inspect(expr, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.SelectorExpr:
			if v.Sel.Name == name {
				found = true
			}
		case *ast.Ident:
			if v.Name == name {
				found = true
			}
		}
		return !found
	})
	return found
}
