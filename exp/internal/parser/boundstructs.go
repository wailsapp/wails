package parser

import "go/ast"

type structInfo struct {
	packageName string
	structName  string
}

func extractBindExprs(imp *ImportSpecInfo) {
	// Check if the given expression is a composite literal
	if comp, ok := (*imp.Options).(*ast.CompositeLit); ok {
		// Iterate through the composite literal fields
		for _, field := range comp.Elts {
			// Check if the field key is "Bind"
			if kv, ok := field.(*ast.KeyValueExpr); ok {
				if ident, ok := kv.Key.(*ast.Ident); ok {
					if ident.Name == "Bind" {
						// Extract the expressions in the Bind field
						if arr, ok := kv.Value.(*ast.CompositeLit); ok {
							for _, elt := range arr.Elts {
								if unaryExpr, ok := elt.(*ast.UnaryExpr); ok {
									if addr, ok := unaryExpr.X.(*ast.CompositeLit); ok {
										// Extract the type of the struct
										if selExpr, ok := addr.Type.(*ast.SelectorExpr); ok {
											if ident, ok := selExpr.X.(*ast.Ident); ok {
												imp.BoundStructNames = append(imp.BoundStructNames, &structInfo{
													packageName: ident.Name,
													structName:  selExpr.Sel.Name,
												})
											}
										} else if ident, ok := addr.Type.(*ast.Ident); ok {
											imp.BoundStructNames = append(imp.BoundStructNames, &structInfo{
												packageName: "",
												structName:  ident.Name,
											})
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
