package buildmeta

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

// Source reads a named top-level string declaration. Literal strings, local
// string aliases, parentheses and concatenation are supported; calls are never
// evaluated. It does not execute init functions or import the application.
func Source(path, symbol string) (Info, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return Info{}, err
	}
	declarations := map[string]ast.Expr{}
	for _, declaration := range file.Decls {
		group, ok := declaration.(*ast.GenDecl)
		if !ok || (group.Tok != token.CONST && group.Tok != token.VAR) {
			continue
		}
		for _, spec := range group.Specs {
			value := spec.(*ast.ValueSpec)
			if len(value.Names) != len(value.Values) {
				continue
			}
			for i, name := range value.Names {
				if name.Name == "_" {
					continue
				}
				if _, exists := declarations[name.Name]; exists {
					return Info{}, fmt.Errorf("duplicate declaration %q in version source", name.Name)
				}
				declarations[name.Name] = value.Values[i]
			}
		}
	}
	value, err := sourceString(ast.NewIdent(symbol), declarations, map[string]bool{})
	if err != nil {
		return Info{}, fmt.Errorf("application version %s in %s: %w", symbol, path, err)
	}
	return Parse(value)
}

func sourceString(expression ast.Expr, declarations map[string]ast.Expr, visiting map[string]bool) (string, error) {
	switch expression := expression.(type) {
	case *ast.BasicLit:
		if expression.Kind == token.STRING {
			return strconv.Unquote(expression.Value)
		}
	case *ast.ParenExpr:
		return sourceString(expression.X, declarations, visiting)
	case *ast.Ident:
		value, exists := declarations[expression.Name]
		if !exists || visiting[expression.Name] {
			return "", fmt.Errorf("missing or cyclic string declaration %q", expression.Name)
		}
		visiting[expression.Name] = true
		defer delete(visiting, expression.Name)
		return sourceString(value, declarations, visiting)
	case *ast.BinaryExpr:
		if expression.Op != token.ADD {
			break
		}
		left, err := sourceString(expression.X, declarations, visiting)
		if err != nil {
			return "", err
		}
		right, err := sourceString(expression.Y, declarations, visiting)
		return left + right, err
	}
	return "", fmt.Errorf("version must be a static string literal, alias or concatenation")
}
