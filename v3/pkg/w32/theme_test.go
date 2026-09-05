//go:build windows

package w32

import (
	"go/ast"
	"go/parser"
	gotoken "go/token"
	"path/filepath"
	"runtime"
	"strconv"
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

// TestDarkModeDetectionReadsExpectedRegistryValues guards the Personalize key
// names: application windows follow AppsUseLightTheme, while system surfaces
// such as the taskbar and its tray icons follow SystemUsesLightTheme.
func TestDarkModeDetectionReadsExpectedRegistryValues(t *testing.T) {
	const personalizeKey = `SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`

	file := parseThemeFile(t)

	for _, testCase := range []struct {
		function string
		value    string
	}{
		{function: "IsCurrentlyDarkMode", value: "AppsUseLightTheme"},
		{function: "IsSystemCurrentlyDarkMode", value: "SystemUsesLightTheme"},
	} {
		t.Run(testCase.function, func(t *testing.T) {
			declaration := findFunctionDeclaration(t, file, testCase.function)
			assertOpensRegistryKey(t, declaration, personalizeKey)
			identifier := assertReadsIntegerValue(t, declaration, testCase.value)
			assertReportsDarkModeWhenZero(t, declaration, identifier)
		})
	}
}

func parseThemeFile(t *testing.T) *ast.File {
	t.Helper()

	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate theme_test.go")
	}

	themeFile := filepath.Join(filepath.Dir(testFile), "theme.go")
	file, err := parser.ParseFile(gotoken.NewFileSet(), themeFile, nil, 0)
	if err != nil {
		t.Fatalf("parse theme.go: %v", err)
	}
	return file
}

func findFunctionDeclaration(t *testing.T, file *ast.File, name string) *ast.FuncDecl {
	t.Helper()

	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if ok && function.Recv == nil && function.Name.Name == name {
			return function
		}
	}

	t.Fatalf("%s not declared in theme.go", name)
	return nil
}

func assertOpensRegistryKey(t *testing.T, declaration *ast.FuncDecl, path string) {
	t.Helper()

	var opened bool
	ast.Inspect(declaration, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok || len(call.Args) != 3 {
			return true
		}

		selector, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || selector.Sel.Name != "OpenKey" {
			return true
		}

		if stringLiteralValue(call.Args[1]) == path {
			opened = true
			return false
		}
		return true
	})

	if !opened {
		t.Fatalf("%s does not open the %s registry key", declaration.Name.Name, path)
	}
}

// assertReadsIntegerValue checks that the function reads valueName and returns
// the identifier the value is assigned to.
func assertReadsIntegerValue(t *testing.T, declaration *ast.FuncDecl, valueName string) string {
	t.Helper()

	var identifier string
	ast.Inspect(declaration, func(node ast.Node) bool {
		assignment, ok := node.(*ast.AssignStmt)
		if !ok || len(assignment.Lhs) == 0 || len(assignment.Rhs) != 1 {
			return true
		}

		call, ok := assignment.Rhs[0].(*ast.CallExpr)
		if !ok || len(call.Args) != 1 {
			return true
		}

		selector, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || selector.Sel.Name != "GetIntegerValue" {
			return true
		}

		if stringLiteralValue(call.Args[0]) != valueName {
			return true
		}

		target, ok := assignment.Lhs[0].(*ast.Ident)
		if !ok {
			return true
		}

		identifier = target.Name
		return false
	})

	if identifier == "" {
		t.Fatalf("%s does not read the %q registry value", declaration.Name.Name, valueName)
	}
	return identifier
}

func assertReportsDarkModeWhenZero(t *testing.T, declaration *ast.FuncDecl, identifier string) {
	t.Helper()

	var inverted bool
	ast.Inspect(declaration, func(node ast.Node) bool {
		statement, ok := node.(*ast.ReturnStmt)
		if !ok || len(statement.Results) != 1 {
			return true
		}

		comparison, ok := statement.Results[0].(*ast.BinaryExpr)
		if !ok || comparison.Op != gotoken.EQL {
			return true
		}

		left, ok := comparison.X.(*ast.Ident)
		if !ok || left.Name != identifier {
			return true
		}

		right, ok := comparison.Y.(*ast.BasicLit)
		if ok && right.Kind == gotoken.INT && right.Value == "0" {
			inverted = true
			return false
		}
		return true
	})

	if !inverted {
		t.Fatalf("%s does not report dark mode as %s == 0", declaration.Name.Name, identifier)
	}
}

func stringLiteralValue(expression ast.Expr) string {
	literal, ok := expression.(*ast.BasicLit)
	if !ok || literal.Kind != gotoken.STRING {
		return ""
	}

	value, err := strconv.Unquote(literal.Value)
	if err != nil {
		return ""
	}
	return value
}
