package backendjs

import (
	"fmt"
	"go/ast"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// Struct defines a parsed struct
type Struct struct {
	Name     string
	Comments []string
	Fields   []*Field
	Methods  []*Method
	IsBound  bool

	// This indicates that the struct is passed as data
	// between the frontend and backend
	IsUsedAsData bool
}

// StructName is used to define a fully qualified struct name
// EG: mypackage.Person
type StructName struct {
	Name    string
	Package string
}

// ToString returns a text representation of the struct name
func (s *StructName) ToString() string {
	result := ""
	if s.Package != "" {
		result = s.Package + "."
	}
	return result + s.Name
}

// Field defines a parsed struct field
type Field struct {
	Name     string
	Type     string
	Struct   *StructName
	Comments []string
}

// JSType returns the Javascript type for this field
func (f *Field) JSType() string {
	return goTypeToJS(f)
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

func (p *Parser) ParseStruct(structType *ast.StructType, name string, pkg *Package) (*Struct, error) {

	// Check if we've seen this struct before
	result := pkg.Structs[name]

	// If we haven't create a new one
	if result == nil {
		result = &Struct{Name: name}
	}

	for _, field := range structType.Fields.List {
		result.Fields = append(result.Fields, p.ParseField(field, pkg)...)
	}
	return result, nil
}

func (p *Parser) parseStructNameFromStarExpr(starExpr *ast.StarExpr) *StructName {
	pkg := ""
	name := ""
	// Determine the FQN
	switch x := starExpr.X.(type) {
	case *ast.SelectorExpr:
		switch i := x.X.(type) {
		case *ast.Ident:
			pkg = i.Name
		default:
			println("one")
			fieldNotSupported(x)
		}

		name = x.Sel.Name

	case *ast.StarExpr:
		switch s := x.X.(type) {
		case *ast.Ident:
			name = s.Name
		default:
			println("two")
			fieldNotSupported(x)
		}
	case *ast.Ident:
		name = x.Name
	default:
		println("three")

		fieldNotSupported(x)
	}
	return &StructName{
		Name:    name,
		Package: pkg,
	}
}

func (p *Parser) ParseField(field *ast.Field, pkg *Package) []*Field {
	var result []*Field

	var fieldType string
	var structName *StructName
	// Determine type
	switch t := field.Type.(type) {
	case *ast.Ident:
		fieldType = t.Name
	case *ast.StarExpr:
		fieldType = "struct"
		structName = p.parseStructNameFromStarExpr(t)
		// Save external reference if we have it
		if structName.Package == "" {
			pkg.structsUsedAsData.AddUnique(structName.Name)
		} else {
			// Save this reference to the external struct
			referencedPackage := p.Packages[structName.Package]
			if referencedPackage == nil {
				// Should we be ignoring this?
			} else {
				pkg.packageReferences.AddUnique(structName.Package)
				referencedPackage.structsUsedAsData.AddUnique(structName.Name)
			}

		}
	default:
		fieldNotSupported(t)
	}

	// Loop over names
	for _, name := range field.Names {

		// Create a field per name
		thisField := &Field{
			Comments: p.parseComments(field.Doc),
		}
		thisField.Name = name.Name
		thisField.Type = fieldType
		thisField.Struct = structName

		result = append(result, thisField)

	}

	return result
}

func fieldNotSupported(t interface{}) {
	println("Field type not supported:")
	spew.Dump(t)
	os.Exit(1)
}

// Method defines a struct method
type Method struct {
	Name     string
	Comments []string
	Inputs   []*Field
	Returns  []*Field
}

// InputsAsTSText generates a string with the method inputs
// formatted in a way acceptable to Typescript
func (m *Method) InputsAsTSText() string {
	var inputs []string

	for _, input := range m.Inputs {
		inputText := fmt.Sprintf("%s: %s", input.Name, goTypeToTS(input))
		inputs = append(inputs, inputText)
	}

	return strings.Join(inputs, ", ")
}

// OutputsAsTSText generates a string with the method inputs
// formatted in a way acceptable to Javascript
func (m *Method) OutputsAsTSText() string {

	if len(m.Returns) == 0 {
		return "void"
	}

	var result []string

	for _, output := range m.Returns {
		result = append(result, goTypeToTS(output))
	}
	return strings.Join(result, ", ")
}

// InputsAsJSText generates a string with the method inputs
// formatted in a way acceptable to Javascript
func (m *Method) InputsAsJSText() string {
	var inputs []string

	for _, input := range m.Inputs {
		inputs = append(inputs, input.Name)
	}

	return strings.Join(inputs, ", ")
}
