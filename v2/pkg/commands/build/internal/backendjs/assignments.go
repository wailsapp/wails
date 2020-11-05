package backendjs

import (
	"go/ast"

	"github.com/davecgh/go-spew/spew"
)

func (p *Parser) parseAssignment(assignStmt *ast.AssignStmt, pkg *Package) {
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
							pkg.variablesThatWereAssignedByFunctions[i.Name] = fe.Name
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
								pkg.variablesThatWereAssignedByStructLiterals[i.Name] = t.Name
							}
						}
					} else {
						e, ok := cl.Type.(*ast.SelectorExpr)
						if ok {
							var thisType = ""
							var thisPackage = ""
							switch x := e.X.(type) {
							case *ast.Ident:
								thisPackage = x.Name
							}
							thisType = e.Sel.Name
							if len(assignStmt.Lhs) == 1 {
								i, ok := assignStmt.Lhs[0].(*ast.Ident)
								if ok {
									sn := &StructName{
										Name:    thisType,
										Package: thisPackage,
									}
									pkg.variablesThatWereAssignedByExternalStructLiterals[i.Name] = sn
								}
							}
						}
					}
				}
			} else {
				cl, ok := rhs.(*ast.CompositeLit)
				if ok {
					t, ok := cl.Type.(*ast.Ident)
					if ok {
						if len(assignStmt.Lhs) == 1 {
							i, ok := assignStmt.Lhs[0].(*ast.Ident)
							if ok {
								pkg.variablesThatWereAssignedByStructLiterals[i.Name] = t.Name
							} else {
								println("herer")
							}
						}
					} else {
						println("herer")
					}
				} else {
					println("herer")
					spew.Dump(rhs)
				}
			}
		}
	}
}
