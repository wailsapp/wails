package parser

import (
	"fmt"
	"go/ast"
)

// Field defines a parsed struct field
type Field struct {
	Name     string
	Type     string
	Struct   *Struct
	Comments []string

	// This struct reference is to temporarily hold the name
	// of the struct during parsing
	structReference *StructReference
}

func (p *Parser) parseField(field *ast.Field, thisPackageName string) ([]*Field, error) {
	var result []*Field

	var fieldType string
	var structReference *StructReference

	// Determine type
	switch t := field.Type.(type) {
	case *ast.Ident:
		fieldType = t.Name
	case *ast.StarExpr:
		fieldType = "struct"
		packageName, structName, err := p.parseStructNameFromStarExpr(t)
		if err != nil {
			return nil, err
		}

		// If we don't ahve a package name, it means it's in this package
		if packageName == "" {
			packageName = thisPackageName
		}

		// Temporarily store the struct reference
		structReference = newStructReference(packageName, structName)

	default:
		return nil, fmt.Errorf("Unsupported field found in struct: %+v", t)
	}

	// Loop over names if we have
	if len(field.Names) > 0 {
		for _, name := range field.Names {

			// Create a field per name
			thisField := &Field{
				Comments: p.parseComments(field.Doc),
			}
			thisField.Name = name.Name
			thisField.Type = fieldType
			thisField.structReference = structReference

			result = append(result, thisField)
		}
		return result, nil
	}

	// When we have no name
	thisField := &Field{
		Comments: p.parseComments(field.Doc),
	}
	thisField.Type = fieldType
	thisField.structReference = structReference

	result = append(result, thisField)

	return result, nil
}

func (p *Parser) resolveFieldReferences(fields []*Field) error {

	// Loop over fields
	for _, field := range fields {

		// If we have a struct reference but no actual struct,
		// we need to resolve it
		if field.structReference != nil && field.Struct == nil {
			fqn := field.structReference.FullyQualifiedName()
			println("Need to resolve struct reference: ", fqn)
			// Check the cache for the struct
			structPointer, err := p.ParseStruct(field.structReference.Package, field.structReference.Name)
			if err != nil {
				return err
			}
			field.Struct = structPointer
			if field.Struct != nil {
				// Save the fact that the struct is used as data
				field.Struct.UsedAsData = true
				println("Resolved struct reference:", fqn)

				// Resolve *its* references
				err = p.resolveStructReferences(field.Struct)
				if err != nil {
					return err
				}
			} else {
				println("Unable to resolve struct reference:", fqn)
			}

		}
	}

	return nil
}
