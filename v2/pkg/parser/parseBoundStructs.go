package parser

import "go/ast"

func (p *Parser) parseBoundStructs(pkg *Package) error {

	// Loop over the bound structs
	for _, structName := range pkg.boundStructs.AsSlice() {
		strct, err := p.parseStruct(pkg, structName)
		if err != nil {
			return err
		}
		strct.IsBound = true
	}

	return nil
}

// ParseStruct will attempt to parse the given struct using
// the package it references
func (p *Parser) parseStruct(pkg *Package, structName string) (*Struct, error) {

	// Check the parser cache for this struct
	result := pkg.getStruct(structName)
	if result != nil {
		return result, nil
	}

	// Iterate through the whole package looking for the bound structs
	for _, fileAst := range pkg.Gopackage.Syntax {

		// Track errors
		var parseError error

		ast.Inspect(fileAst, func(n ast.Node) bool {
			if genDecl, ok := n.(*ast.GenDecl); ok {
				for _, spec := range genDecl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if structType, ok := typeSpec.Type.(*ast.StructType); ok {
							structDefinitionName := typeSpec.Name.Name
							if structDefinitionName == structName {

								// Create the new struct
								result = &Struct{Name: structName, Package: pkg}

								// Save comments
								result.Comments = parseComments(genDecl.Doc)

								parseError = p.parseStructMethods(result)
								if parseError != nil {
									return false
								}

								// Parse the struct fields
								parseError = p.parseStructFields(fileAst, structType, result)

								// Save this struct
								pkg.addStruct(result)

								return false
							}
						}
					}
				}
			}
			return true
		})

		// If we got an error, return it
		if parseError != nil {
			return nil, parseError
		}
	}
	return result, nil
}
