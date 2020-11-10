package parser

import (
	"fmt"
	"go/ast"

	"github.com/davecgh/go-spew/spew"
)

// Field defines a parsed struct field
type Field struct {

	// Name of the field
	Name string

	// The type of the field.
	// "struct" if it's a struct
	Type string

	// A pointer to the struct if the Type is "struct"
	Struct *Struct

	// User comments on the field
	Comments []string

	// Indicates if the Field is an array of type "Type"
	IsArray bool
}

// JSType returns the Javascript type for this field
func (f *Field) JSType() string {
	return string(goTypeToJS(f))
}

// TypeAsTSType converts the Field type to something TS wants
func (f *Field) TypeAsTSType(pkgName string) string {
	var result = ""
	switch f.Type {
	case "string":
		result = "string"
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		result = "number"
	case "float32", "float64":
		result = "number"
	case "bool":
		result = "boolean"
	case "struct":
		if f.Struct.Package != nil {
			if f.Struct.Package.Name != pkgName {
				result = f.Struct.Package.Name + "."
			}
		}
		result = result + f.Struct.Name
	default:
		result = "any"
	}

	return result
}

func (p *Parser) parseField(file *ast.File, field *ast.Field, pkg *Package) ([]*Field, error) {
	var result []*Field

	var fieldType string
	var strct *Struct
	var isArray bool

	// Determine type
	switch t := field.Type.(type) {
	case *ast.Ident:
		fieldType = t.Name
	case *ast.StarExpr:
		fieldType = "struct"
		packageName, structName, err := parseStructNameFromStarExpr(t)
		if err != nil {
			return nil, err
		}

		// If this is an external package, find it
		if packageName != "" {
			referencedGoPackage := pkg.getImportByName(packageName, file)
			referencedPackage := p.getPackageByID(referencedGoPackage.ID)

			// If we found the struct, save it as an external package reference
			if referencedPackage != nil {
				pkg.addExternalReference(referencedPackage)
			}

			// We save this to pkg anyway, because we want to know if this package
			// was NOT found
			pkg = referencedPackage
		}

		// If this is a package in our project, parse the struct!
		if pkg != nil {

			// Parse the struct
			strct, err = p.parseStruct(pkg, structName)
			if err != nil {
				return nil, err
			}

		}

	case *ast.ArrayType:
		isArray = true
		// Parse the Elt (There must be a better way!)
		switch t := t.Elt.(type) {
		case *ast.Ident:
			fieldType = t.Name
		case *ast.StarExpr:
			fieldType = "struct"
			packageName, structName, err := parseStructNameFromStarExpr(t)
			if err != nil {
				return nil, err
			}

			// If this is an external package, find it
			if packageName != "" {
				referencedGoPackage := pkg.getImportByName(packageName, file)
				referencedPackage := p.getPackageByID(referencedGoPackage.ID)

				// If we found the struct, save it as an external package reference
				if referencedPackage != nil {
					pkg.addExternalReference(referencedPackage)
				}

				// We save this to pkg anyway, because we want to know if this package
				// was NOT found
				pkg = referencedPackage
			}

			// If this is a package in our project, parse the struct!
			if pkg != nil {

				// Parse the struct
				strct, err = p.parseStruct(pkg, structName)
				if err != nil {
					return nil, err
				}

			}
		default:
			// We will default to "Array<any>" for eg nested arrays
			fieldType = "any"
		}

	default:
		spew.Dump(t)
		return nil, fmt.Errorf("unsupported field found in struct: %+v", t)
	}

	// Loop over names if we have
	if len(field.Names) > 0 {
		for _, name := range field.Names {

			// Create a field per name
			thisField := &Field{
				Comments: parseComments(field.Doc),
			}
			thisField.Name = name.Name
			thisField.Type = fieldType
			thisField.Struct = strct
			thisField.IsArray = isArray

			result = append(result, thisField)
		}
		return result, nil
	}

	// When we have no name
	thisField := &Field{
		Comments: parseComments(field.Doc),
	}
	thisField.Type = fieldType
	thisField.Struct = strct
	thisField.IsArray = isArray
	result = append(result, thisField)

	return result, nil
}
