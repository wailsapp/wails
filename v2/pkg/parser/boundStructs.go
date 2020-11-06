package parser

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

func (p *Parser) getImportByName(pkg *packages.Package, importName string) *packages.Package {
	// Find package path
	for _, imp := range pkg.Imports {
		if imp.Name == importName {
			return imp
		}
	}
	return nil
}

func (p *Parser) findBoundStructsInPackage(pkg *packages.Package, applicationVariableName string) []*StructReference {

	var boundStructs []*StructReference

	// Iterate through the whole package looking for the bound structs
	for _, fileAst := range pkg.Syntax {
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
					boundStruct := newStructReference(pkg.Name, boundStructExp.Name)
					boundStructs = append(boundStructs, boundStruct)

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
					referencedPackage := p.getImportByName(pkg, packageName)
					boundStruct := newStructReference(referencedPackage.Name, structName)
					boundStructs = append(boundStructs, boundStruct)
				}

			// Binding struct literals
			case *ast.CompositeLit:
				switch literal := boundItem.Type.(type) {

				// app.Bind( myStruct{} )
				case *ast.Ident:
					structName := literal.Name
					boundStruct := newStructReference(pkg.Name, structName)
					boundStructs = append(boundStructs, boundStruct)

				// app.Bind( mypackage.myStruct{} )
				case *ast.SelectorExpr:
					var structName = ""
					var packageName = ""
					switch x := literal.X.(type) {
					case *ast.Ident:
						packageName = x.Name
					default:
						// TODO: Save these warnings
						// println("Identifier in binding not supported:")
						return true
					}
					structName = literal.Sel.Name

					referencedPackage := p.getImportByName(pkg, packageName)
					boundStruct := newStructReference(referencedPackage.Name, structName)
					boundStructs = append(boundStructs, boundStruct)
				}

			default:
				// TODO: Save these warnings
				// println("Unsupported bind expression:")
				// spew.Dump(boundItem)
			}

			return true
		})
	}
	return boundStructs
}
