package backendjs

import "go/ast"

func (p *Parser) parseCallExpressions(x *ast.CallExpr) {
	f, ok := x.Fun.(*ast.SelectorExpr)
	if ok {
		n, ok := f.X.(*ast.Ident)
		if ok {
			//Check this is the Bind() call associated with the app variable
			if n.Name == p.applicationVariable && f.Sel.Name == "Bind" {
				if len(x.Args) == 1 {
					ce, ok := x.Args[0].(*ast.CallExpr)
					if ok {
						fn, ok := ce.Fun.(*ast.Ident)
						if ok {
							// We found a bind method using a function call
							// EG: app.Bind( newMyStruct() )
							p.structMethodsThatWereBound.Add(fn.Name)
						}
					} else {
						// We also want to check for Bind( &MyStruct{} )
						ue, ok := x.Args[0].(*ast.UnaryExpr)
						if ok {
							if ue.Op.String() == "&" {
								cl, ok := ue.X.(*ast.CompositeLit)
								if ok {
									t, ok := cl.Type.(*ast.Ident)
									if ok {
										// We have found Bind( &MyStruct{} )
										p.structPointerLiteralsThatWereBound.Add(t.Name)
									}
								}
							}
						} else {
							// Let's check when the user binds a struct,
							// rather than a struct pointer: Bind( MyStruct{} )
							// We do this to provide better hints to the user
							cl, ok := x.Args[0].(*ast.CompositeLit)
							if ok {
								t, ok := cl.Type.(*ast.Ident)
								if ok {
									p.structLiteralsThatWereBound.Add(t.Name)
								}
							} else {
								// Also check for when we bind a variable
								// myVariable := &MyStruct{}
								// app.Bind( myVariable )
								i, ok := x.Args[0].(*ast.Ident)
								if ok {
									p.variablesThatWereBound.Add(i.Name)
								}
							}
						}
					}
				}
			}
		}
	}
}
