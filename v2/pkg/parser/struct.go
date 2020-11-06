package parser

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/packages"
)

type Struct struct {
	Package  *packages.Package
	Name     string
	Comments []string
	Fields   []*Field
	Methods  []*Method

	// This is true when this struct is used as a datatype
	UsedAsData bool
}

// newStruct creates a new struct and stores in the cache
func (p *Parser) newStruct(pkg *packages.Package, name string) *Struct {

	result := &Struct{
		Package: pkg,
		Name:    name,
	}

	return result
}

// FullyQualifiedName returns the fully qualified name of this struct
func (s *Struct) FullyQualifiedName() string {
	return s.Package.Name + "." + s.Name
}

func (p *Parser) parseStructNameFromStarExpr(starExpr *ast.StarExpr) (string, string, error) {
	pkg := ""
	name := ""
	// Determine the FQN
	switch x := starExpr.X.(type) {
	case *ast.SelectorExpr:
		switch i := x.X.(type) {
		case *ast.Ident:
			pkg = i.Name
		default:
			return "", "", fmt.Errorf("Unsupported Selector expression: %+v", i)
		}

		name = x.Sel.Name

	case *ast.StarExpr:
		switch s := x.X.(type) {
		case *ast.Ident:
			name = s.Name
		default:
			return "", "", fmt.Errorf("Unsupported Star expression: %+v", s)
		}
	case *ast.Ident:
		name = x.Name
	default:
		return "", "", fmt.Errorf("Unsupported Star.X expression: %+v", x)
	}
	return pkg, name, nil
}

// StructReference defines a reference to a fully qualified struct
type StructReference struct {
	Package string
	Name    string
}

func newStructReference(packageName string, structName string) *StructReference {
	return &StructReference{Package: packageName, Name: structName}
}

// FullyQualifiedName returns a string representing the struct reference
func (s *StructReference) FullyQualifiedName() string {
	return s.Package + "." + s.Name
}

func (p *Parser) resolveStructReferences(boundStruct *Struct) error {

	var err error

	// Resolve field references
	err = p.resolveFieldReferences(boundStruct.Fields)
	if err != nil {
		return nil
	}

	// Check if method fields need resolving
	for _, method := range boundStruct.Methods {

		// Resolve method inputs
		err = p.resolveFieldReferences(method.Inputs)
		if err != nil {
			return nil
		}

		// Resolve method outputs
		err = p.resolveFieldReferences(method.Returns)
		if err != nil {
			return nil
		}
	}

	return nil
}
