package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func GenerateConstants(goData []byte) (string, error) {

	// Create a new token file set and parser
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "", goData, parser.AllErrors)
	if err != nil {
		return "", err
	}

	// Extract constant declarations and generate JavaScript constants
	var jsConstants []string
	for _, decl := range f.Decls {
		if gd, ok := decl.(*ast.GenDecl); ok && gd.Tok == token.CONST {
			for _, spec := range gd.Specs {
				if vs, ok := spec.(*ast.ValueSpec); ok {
					for i, name := range vs.Names {
						value := vs.Values[i]
						if value != nil {
							jsConstants = append(jsConstants, fmt.Sprintf("export const %s = %s;", name.Name, jsValue(value)))
						}
					}
				}
			}
		}
	}

	// Join the JavaScript constants into a single string
	jsCode := strings.Join(jsConstants, "\n")

	return jsCode, nil
}

func jsValue(expr ast.Expr) string {
	// Implement conversion from Go constant value to JavaScript value here.
	// You can add more cases for different types if needed.
	switch e := expr.(type) {
	case *ast.BasicLit:
		return e.Value
	case *ast.Ident:
		return e.Name
	default:
		return ""
	}
}
