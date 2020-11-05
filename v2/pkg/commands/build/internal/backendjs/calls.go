package backendjs

import (
	"go/ast"

	"github.com/davecgh/go-spew/spew"
)

func (p *Parser) parseCallExpressions(x *ast.CallExpr, pkg *Package) {
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
							pkg.structMethodsThatWereBound.Add(fn.Name)
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
										pkg.structPointerLiteralsThatWereBound.Add(t.Name)
									} else {
										e, ok := cl.Type.(*ast.SelectorExpr)
										if ok {
											var thisType = ""
											var thisPackage = ""
											switch x := e.X.(type) {
											case *ast.Ident:
												thisPackage = x.Name
											default:
												println("Identifier in binding not supported:")
												spew.Dump(ue.X)
												return
											}
											thisType = e.Sel.Name
											// Save this reference in the package
											pkg := p.getOrCreatePackage(thisPackage)
											pkg.structPointerLiteralsThatWereBound.Add(thisType)
										} else {
											println("Binding not supported:")
											spew.Dump(ue.X)
										}
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
									pkg.structLiteralsThatWereBound.Add(t.Name)

								}
							} else {
								// Also check for when we bind a variable
								// myVariable := &MyStruct{}
								// app.Bind( myVariable )
								i, ok := x.Args[0].(*ast.Ident)
								if ok {
									pkg.variablesThatWereBound.Add(i.Name)
								} else {

								}
							}
						}
					}
				}
			}
		}
	}
}
