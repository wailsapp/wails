package analyse

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

// RefMap maps fields and variables to the list of all ast nodes
// that refer to them or define them within some package.
type RefMap = map[*types.Var][]ast.Node

// BuildRefMap builds the ref map for the given package.
func BuildRefMap(pkg *packages.Package) RefMap {
	result := make(RefMap)

	for ident, obj := range pkg.TypesInfo.Defs {
		if v, ok := obj.(*types.Var); ok {
			refs, ok := result[v]
			if !ok {
				refs = make([]ast.Node, 0, 4)
			}
			result[v] = append(refs, ident)
		}
	}

	for ident, obj := range pkg.TypesInfo.Uses {
		if v, ok := obj.(*types.Var); ok {
			refs, ok := result[v]
			if !ok {
				refs = make([]ast.Node, 0, 4)
			}
			result[v] = append(refs, ident)
		}
	}

	for node, obj := range pkg.TypesInfo.Implicits {
		if v, ok := obj.(*types.Var); ok {
			refs, ok := result[v]
			if !ok {
				refs = make([]ast.Node, 0, 4)
			}
			result[v] = append(refs, node)
		}
	}

	return result
}
