package parser

import "go/ast"

func (p *Package) getApplicationVariableName(file *ast.File, wailsImportName string) string {

	// Iterate through the whole file looking for the application name
	applicationVariableName := ""

	// Inspect the file
	ast.Inspect(file, func(n ast.Node) bool {
		// Parse Assignments looking for application name
		if assignStmt, ok := n.(*ast.AssignStmt); ok {

			// Check the RHS is of the form:
			//   `app := wails.CreateApp()` or
			//   `app := wails.CreateAppWithOptions`
			for _, rhs := range assignStmt.Rhs {
				ce, ok := rhs.(*ast.CallExpr)
				if !ok {
					continue
				}
				se, ok := ce.Fun.(*ast.SelectorExpr)
				if !ok {
					continue
				}
				i, ok := se.X.(*ast.Ident)
				if !ok {
					continue
				}
				// Have we found the wails import name?
				if i.Name == wailsImportName {
					// Check we are calling a function to create the app
					if se.Sel.Name == "CreateApp" || se.Sel.Name == "CreateAppWithOptions" {
						if len(assignStmt.Lhs) == 1 {
							i, ok := assignStmt.Lhs[0].(*ast.Ident)
							if ok {
								// Found the app variable name
								applicationVariableName = i.Name
								return false
							}
						}
					}
				}
			}
		}
		return true
	})
	return applicationVariableName
}
