//go:build windows

package w32

import (
	"go/ast"
	"go/parser"
	gotoken "go/token"
	"path/filepath"
	"runtime"
	"testing"
)

func TestAllowDarkModeForWindowPassesHWNDAndAllowFlag(t *testing.T) {
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate theme_test.go")
	}

	themeFile := filepath.Join(filepath.Dir(testFile), "theme.go")
	file, err := parser.ParseFile(gotoken.NewFileSet(), themeFile, nil, 0)
	if err != nil {
		t.Fatalf("parse theme.go: %v", err)
	}

	var syscallCall *ast.CallExpr
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok || len(call.Args) == 0 {
			return true
		}

		selector, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || selector.Sel.Name != "SyscallN" {
			return true
		}

		proc, ok := call.Args[0].(*ast.Ident)
		if ok && proc.Name == "procAllowDarkModeForWindow" {
			syscallCall = call
			return false
		}
		return true
	})

	if syscallCall == nil {
		t.Fatal("AllowDarkModeForWindow syscall not found")
	}
	if len(syscallCall.Args) != 3 {
		t.Fatalf("SyscallN argument count = %d; want 3 (procedure, HWND, allow flag)", len(syscallCall.Args))
	}

	assertUintptrIdent(t, syscallCall.Args[1], "hwnd")
	assertUintptrIdent(t, syscallCall.Args[2], "allowInt")
}

func assertUintptrIdent(t *testing.T, expression ast.Expr, name string) {
	t.Helper()

	conversion, ok := expression.(*ast.CallExpr)
	if !ok || len(conversion.Args) != 1 {
		t.Fatalf("argument is not a single-value uintptr conversion: %#v", expression)
	}

	typeName, ok := conversion.Fun.(*ast.Ident)
	if !ok || typeName.Name != "uintptr" {
		t.Fatalf("argument conversion = %#v; want uintptr", conversion.Fun)
	}

	identifier, ok := conversion.Args[0].(*ast.Ident)
	if !ok || identifier.Name != name {
		t.Fatalf("converted argument = %#v; want %s", conversion.Args[0], name)
	}
}
