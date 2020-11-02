package backendjs

import (
	"go/ast"
	"os"

	"github.com/davecgh/go-spew/spew"
)

// Struct defines a parsed struct
type Struct struct {
	Name     string
	Comments []string
	Fields   []*Field
}

type StructName struct {
	Name    string
	Package string
}

// Field defines a parsed struct field
type Field struct {
	Name     string
	Type     string
	Struct   *StructName
	Comments []string
}

// TypeAsTSType converts the Field type to something TS wants
func (f *Field) TypeAsTSType() string {
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
		if f.Struct.Package != "" {
			result = f.Struct.Package + "."
		}
		result = result + f.Struct.Name
	default:
		result = "any"
	}

	return result
}

func parseStruct(structType *ast.StructType, name string) (*Struct, error) {
	result := &Struct{Name: name}

	for _, field := range structType.Fields.List {
		result.Fields = append(result.Fields, parseField(field)...)
	}
	return result, nil
}

func parseField(field *ast.Field) []*Field {
	var result []*Field

	var fieldType string
	var structName *StructName
	// Determine type
	switch t := field.Type.(type) {
	case *ast.Ident:
		fieldType = t.Name
	case *ast.StarExpr:
		pkg := ""
		name := ""
		// Determine the FQN
		switch x := t.X.(type) {
		case *ast.SelectorExpr:
			switch i := x.X.(type) {
			case *ast.Ident:
				pkg = i.Name
			default:
				println("one")
				FieldNotSupported(x)
			}

			name = x.Sel.Name

		case *ast.StarExpr:
			switch s := x.X.(type) {
			case *ast.Ident:
				name = s.Name
			default:
				println("two")
				FieldNotSupported(x)
			}
		case *ast.Ident:
			name = x.Name
		default:
			println("three")

			FieldNotSupported(x)
		}
		fieldType = "struct"
		structName = &StructName{
			Name:    name,
			Package: pkg,
		}

	default:
		FieldNotSupported(t)
	}

	// Loop over names
	for _, name := range field.Names {

		// Create a field per name
		thisField := &Field{
			Comments: parseComments(field.Doc),
		}
		thisField.Name = name.Name
		thisField.Type = fieldType
		thisField.Struct = structName

		result = append(result, thisField)

	}

	return result
}

func FieldNotSupported(t interface{}) {
	println("Field type not supported:")
	spew.Dump(t)
	os.Exit(1)
}
