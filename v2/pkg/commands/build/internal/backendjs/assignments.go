package backendjs

import "go/ast"

func (p *Parser) parseAssignment(assignStmt *ast.AssignStmt) {
	for _, rhs := range assignStmt.Rhs {
		ce, ok := rhs.(*ast.CallExpr)
		if ok {
			se, ok := ce.Fun.(*ast.SelectorExpr)
			if ok {
				i, ok := se.X.(*ast.Ident)
				if ok {
					// Have we found the wails package name?
					if i.Name == p.wailsPackageVariable {
						// Check we are calling a function to create the app
						if se.Sel.Name == "CreateApp" || se.Sel.Name == "CreateAppWithOptions" {
							if len(assignStmt.Lhs) == 1 {
								i, ok := assignStmt.Lhs[0].(*ast.Ident)
								if ok {
									// Found the app variable name
									p.applicationVariable = i.Name
								}
							}
						}
					}
				}
			} else {
				// Check for function assignment
				// a := newMyStruct()
				fe, ok := ce.Fun.(*ast.Ident)
				if ok {
					if len(assignStmt.Lhs) == 1 {
						i, ok := assignStmt.Lhs[0].(*ast.Ident)
						if ok {
							// Store the variable -> Function mapping
							// so we can later resolve the type
							p.variablesThatWereAssignedByFunctions[i.Name] = fe.Name
						}
					}
				}
			}
		} else {
			// Check for literal assignment of struct
			// EG: myvar := MyStruct{}
			ue, ok := rhs.(*ast.UnaryExpr)
			if ok {
				cl, ok := ue.X.(*ast.CompositeLit)
				if ok {
					t, ok := cl.Type.(*ast.Ident)
					if ok {
						if len(assignStmt.Lhs) == 1 {
							i, ok := assignStmt.Lhs[0].(*ast.Ident)
							if ok {
								p.variablesThatWereAssignedByStructLiterals[i.Name] = t.Name
							}
						}
					}
				}
			}
		}
	}
}
