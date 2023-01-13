package parser

import (
	"go/ast"
)

func isOptionsApplication(expr *ast.Expr) bool {
	cl, ok := (*expr).(*ast.CompositeLit)
	if ok {
		se, ok := cl.Type.(*ast.SelectorExpr)
		if ok {
			_, ok := se.X.(*ast.Ident)
			if ok {
				if se.Sel.Name == "Application" && se.X.(*ast.Ident).Name == "options" {
					return true
				}
			}
		}
	}

	return false
}
