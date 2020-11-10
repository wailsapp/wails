package parser

import (
	"fmt"
	"go/ast"

	"github.com/pkg/errors"
)

// Struct represents a struct that is used by the frontend
// in a Wails project
type Struct struct {

	// The name of the struct
	Name string

	// The package this was declared in
	Package *Package

	// Comments for the struct
	Comments []string

	// The fields used in this struct
	Fields []*Field

	// The methods available to the front end
	Methods []*Method

	// Indicates if this struct is bound to the app
	IsBound bool

	// Indicates if this struct is used as data
	IsUsedAsData bool
}

func parseStructNameFromStarExpr(starExpr *ast.StarExpr) (string, string, error) {
	pkg := ""
	name := ""
	// Determine the FQN
	switch x := starExpr.X.(type) {
	case *ast.SelectorExpr:
		switch i := x.X.(type) {
		case *ast.Ident:
			pkg = i.Name
		default:
			// TODO: Store warnings?
			return "", "", errors.WithStack(fmt.Errorf("unknown type in selector for *ast.SelectorExpr: %+v", i))
		}

		name = x.Sel.Name

	// TODO: IS this used?
	case *ast.StarExpr:
		switch s := x.X.(type) {
		case *ast.Ident:
			name = s.Name
		default:
			// TODO: Store warnings?
			return "", "", errors.WithStack(fmt.Errorf("unknown type in selector for *ast.StarExpr: %+v", s))
		}
	case *ast.Ident:
		name = x.Name
	default:
		// TODO: Store warnings?
		return "", "", errors.WithStack(fmt.Errorf("unknown type in selector for *ast.StarExpr: %+v", starExpr))
	}
	return pkg, name, nil
}
