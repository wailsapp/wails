package parser

import (
	"go/ast"
	"os"
)

func findNewCalls(imports []ImportSpecInfo) *ImportSpecInfo {

	var result *ast.Expr
	var inImport *ImportSpecInfo

	for _, imp := range imports {
		for _, decl := range imp.File.Decls {
			ast.Inspect(decl, func(n ast.Node) bool {
				// check if the current node is a function call
				callExpr, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				// check if the function being called is .New()
				selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
				if !ok || selExpr.Sel.Name != "New" {
					return true
				}

				// Get the name of the thing that New is being called on.
				var receiverName string
				switch rcv := selExpr.X.(type) {
				case *ast.Ident:
					receiverName = rcv.Name
				case *ast.SelectorExpr:
					receiverName = rcv.Sel.Name
				default:
					receiverName = "unknown"
				}

				// Check if the receiver is a package name
				for _, i := range imports {
					if i.Identifier() == receiverName {
						if i.ImportSpec.Path.Value == `"github.com/wailsapp/wails/exp/pkg/application"` {
							// We have a call to application.New()
							// Parse out the first argument and check it is an option.Application struct
							// check callExpr.Args[0] is a struct literal
							if len(callExpr.Args) != 1 {
								return true
							}
							if result != nil {
								println("Found more than one call to application.New()")
								os.Exit(1)
							}
							result = &(callExpr.Args[0])
							inImport = &imp

							return false
						}
					}
				}
				return true
			})
		}
	}
	inImport.Options = result
	return inImport
}
