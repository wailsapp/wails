package parser

import (
	"go/ast"
)

// getCachedStruct attempts to get an already parsed struct from the
// struct cache
func (p *Parser) getCachedStruct(packageName string, structName string) *Struct {
	fqn := packageName + "." + structName
	return p.parsedStructs[fqn]
}

// ParseStruct will attempt to parse the given struct using
// the package it references
func (p *Parser) ParseStruct(packageName string, structName string) (*Struct, error) {

	// Check the cache
	result := p.getCachedStruct(packageName, structName)
	if result != nil {
		return result, nil
	}

	// Find the package
	pkg := p.getPackageByName(packageName)
	if pkg == nil {
		// TODO: Find package via imports?
		println("Cannot find package", packageName)
		return nil, nil
	}

	// Iterate through the whole package looking for the bound structs
	for _, fileAst := range pkg.Syntax {

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
								result = p.newStruct(pkg, structDefinitionName)

								// Save comments
								result.Comments = p.parseComments(genDecl.Doc)

								parseError = p.parseStructMethods(result)
								if parseError != nil {
									return false
								}

								// Parse the struct fields
								parseError = p.parseStructFields(structType, result)

								// Cache this struct
								key := result.FullyQualifiedName()
								p.parsedStructs[key] = result

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

func (p *Parser) parseStructFields(structType *ast.StructType, boundStruct *Struct) error {

	// Parse the fields
	for _, field := range structType.Fields.List {
		fields, err := p.parseField(field, boundStruct.Package.Name)
		if err != nil {
			return err
		}
		boundStruct.Fields = append(boundStruct.Fields, fields...)
	}

	return nil
}
