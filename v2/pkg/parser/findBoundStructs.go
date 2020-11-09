package parser

import "go/ast"

// findBoundStructs will search through the Wails project looking
// for which structs have been bound using the `Bind()` method
func (p *Parser) findBoundStructs(pkg *Package) error {

	// Iterate through the files in the package looking for the bound structs
	for _, fileAst := range pkg.gopackage.Syntax {

		// Find the wails import name
		wailsImportName := pkg.getWailsImportName(fileAst)

		// If this file doesn't import wails, continue
		if wailsImportName == "" {
			continue
		}

		applicationVariableName := pkg.getApplicationVariableName(fileAst, wailsImportName)
		if applicationVariableName == "" {
			continue
		}

		var parseError error

		ast.Inspect(fileAst, func(n ast.Node) bool {
			// Parse Call expressions looking for bind calls
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			// Check this is the right kind of expression (something.something())
			f, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			ident, ok := f.X.(*ast.Ident)
			if !ok {
				return true
			}

			if ident.Name != applicationVariableName {
				return true
			}

			if f.Sel.Name != "Bind" {
				return true
			}

			if len(callExpr.Args) != 1 {
				return true
			}

			// Work out what was bound
			switch boundItem := callExpr.Args[0].(type) {

			// // app.Bind( someFunction() )
			// case *ast.CallExpr:
			// 	switch fn := boundItem.Fun.(type) {
			// 	case *ast.Ident:
			// 		boundStructs = append(boundStructs, newStruct(pkg.Name, fn.Name))
			// 		println("Found bound function:", fn.Name)
			// 	case *ast.SelectorExpr:
			// 		ident, ok := fn.X.(*ast.Ident)
			// 		if !ok {
			// 			return true
			// 		}
			// 		packageName := ident.Name
			// 		functionName := fn.Sel.Name
			// 		println("Found bound function:", packageName+"."+functionName)

			// 		strct := p.getFunctionReturnType(packageName, functionName)
			// 		if strct == nil {
			// 			// Unable to resolve function
			// 			return true
			// 		}
			// 		boundStructs = append(boundStructs, strct)
			// 	}

			// Binding struct pointer literals
			case *ast.UnaryExpr:

				if boundItem.Op.String() != "&" {
					return true
				}

				cl, ok := boundItem.X.(*ast.CompositeLit)
				if !ok {
					return true
				}

				switch boundStructExp := cl.Type.(type) {

				// app.Bind( &myStruct{} )
				case *ast.Ident:
					pkg.boundStructs.Add(boundStructExp.Name)

				// app.Bind( &mypackage.myStruct{} )
				case *ast.SelectorExpr:
					var structName = ""
					var packageName = ""
					switch x := boundStructExp.X.(type) {
					case *ast.Ident:
						packageName = x.Name
					default:
						// TODO: Save these warnings
						// println("Identifier in binding not supported:")
						return true
					}
					structName = boundStructExp.Sel.Name
					referencedPackage := pkg.getImportByName(packageName, fileAst)
					packageWrapper := p.getPackageByID(referencedPackage.ID)
					packageWrapper.boundStructs.Add(structName)
				}

			// // Binding struct literals
			// case *ast.CompositeLit:
			// 	switch literal := boundItem.Type.(type) {

			// 	// app.Bind( myStruct{} )
			// 	case *ast.Ident:
			// 		structName := literal.Name
			// 		boundStructReference := newStructReference(p.GoPackage, structName)
			// 		p.addBoundStructReference(boundStructReference)

			// 	// app.Bind( mypackage.myStruct{} )
			// 	case *ast.SelectorExpr:
			// 		var structName = ""
			// 		var packageName = ""
			// 		switch x := literal.X.(type) {
			// 		case *ast.Ident:
			// 			packageName = x.Name
			// 		default:
			// 			// TODO: Save these warnings
			// 			// println("Identifier in binding not supported:")
			// 			return true
			// 		}
			// 		structName = literal.Sel.Name

			// 		referencedPackage := p.getImportByName(pkg, packageName)
			// 		boundStructReference := newStructReference(referencedPackage, structName)
			// 		p.addBoundStructReference(boundStructReference)
			// 	}

			default:
				// TODO: Save these warnings
				// println("Unsupported bind expression:")
				// spew.Dump(boundItem)
			}

			return true
		})

		if parseError != nil {
			return parseError
		}
	}

	return nil
}
