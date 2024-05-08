package analyser

import (
	"errors"
	"go/ast"
	"go/token"
	"go/types"
)

// WailsAppPkgPath is the official import path of Wails v3's application package
const WailsAppPkgPath = "github.com/wailsapp/wails/v3/pkg/application"

// ErrNoApplicationPackage indicates that the Wails application package
// does not appear in the dependency graph of the packages under analysis.
var ErrNoApplicationPackage = errors.New(WailsAppPkgPath + ": package not found")

// ErrBadApplicationPackage indicates that the Wails application package
// was found but did not have the expected shape.
var ErrBadApplicationPackage = errors.New(WailsAppPkgPath + ": field Options.Bind is missing or mistyped: is the Wails v3 module properly installed?")

// RootTarget stores the root target for the static analyser,
// that is the Option struct from the Wails application package,
// its field named Bind, and the index of the field within the struct.
type RootTarget struct {
	Type  *types.Named
	Field *types.Var
	Index int
}

// FindRootTarget performs a BFS over the dependency graph
// of the given packages, looking for the Wails application package.
// If the search succeeds, it tries to look up the root target
// within the resulting package (see [RootTarget]).
//
// If the search fails, the returned error may be either
// [ErrNoApplicationPackage] or [ErrBadApplicationPackage].
func FindRootTarget(pkgs []*types.Package) (RootTarget, error) {
	current := make([]*types.Package, 0, len(pkgs))
	next := make([]*types.Package, 1, len(pkgs))
	visited := make(map[*types.Package]bool, 2*len(pkgs))

	// Create a fake root package
	next[0] = types.NewPackage("", "")
	next[0].SetImports(pkgs)

	for len(next) > 0 {
		// Advance to next level
		current, next = next, current[:0]

		for _, pkg := range current {
			for _, dep := range pkg.Imports() {
				if dep.Path() == WailsAppPkgPath {
					return findRootTargetInApplicationPackage(dep)
				}

				if !IsStdImportPath(dep.Path()) && !visited[dep] {
					visited[dep] = true
					next = append(next, dep)
				}
			}
		}
	}

	return RootTarget{}, ErrNoApplicationPackage
}

// typeSliceAny caches the expected underlying type
// of the field application.Options.Bind, that is []interface{}
var typeSliceAny = types.NewSlice(types.Universe.Lookup("any").Type().Underlying())

// findRootTargetInApplicationPackage looks up a struct named `Option`
// with a field named `Bind` within the given package.
// If the lookup succeeds, it checks that the field has type `[]any`
// and does not come from an embedded struct.
func findRootTargetInApplicationPackage(app *types.Package) (tgt RootTarget, err error) {
	err = ErrBadApplicationPackage

	// Type-check the expression "Options{}.Bind"
	// within the scope of the Wails application package.
	expr := &ast.SelectorExpr{
		X:   &ast.CompositeLit{Type: &ast.Ident{Name: "Options"}},
		Sel: &ast.Ident{Name: "Bind"},
	}

	info := types.Info{
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}

	if types.CheckExpr(token.NewFileSet(), app, token.NoPos, expr, &info) != nil {
		return
	}

	// Retrieve selection result
	sel := info.Selections[expr]
	if sel == nil || sel.Kind() != types.FieldVal {
		return
	}

	// Retrieve and check receiver type
	recv, ok := sel.Recv().(*types.Named)
	if !ok {
		return
	}

	// Check underlying type of receiver
	if _, ok := recv.Underlying().(*types.Struct); !ok {
		return
	}

	// Retrieve selected field and check its type
	field, ok := sel.Obj().(*types.Var)
	if !ok || !types.Identical(field.Type().Underlying(), typeSliceAny) {
		return
	}

	// Check field depth
	if len(sel.Index()) != 1 {
		err = ErrBadApplicationPackage
		return
	}

	// Save Options struct, Bind field and canonical path
	tgt = RootTarget{
		Type:  recv,
		Field: field,
		Index: sel.Index()[0],
	}
	err = nil

	return
}

// processRootTarget performs static analysis on the root target.
func (analyser *Analyser) processRootTarget() {
	for i := range analyser.pkgs {
		analyser.processRootTargetInPackage(i)
	}
}

// processRootTargetInPackage scans the package indexed by `pkgi`
// in `analyser.pkgs` for
//   - all composite literals whose type is identical to `analyser.root.Type`;
//   - all assignable selector expressions that refer to `analyser.root.Field`.
//
// The former are fed to [Analyser.processExpression]
// with path `[analyser.root.Index, IndexingStep]`,
// the latter to [Analyser.processReference]
// with path `[IndexingStep]`.
//
// The purpose of this step is to find all assignments to the root field
// (i.e., `application.Options.Bind`) in the given package.
// There are two ways to perform such an assignment:
// either some expression is assigned to the field
// in a composite literal of type application.Options,
// or the field is selected through a selector expression
// (whose receiver might also be a struct embedding application.Options)
// which in turn occurs on the left side of an assignment.
func (analyser *Analyser) processRootTargetInPackage(pkgi int) {
	pkg := analyser.pkgs[pkgi]
	path := NewPath(Step(analyser.root.Index), IndexingStep)

	for expr, tv := range pkg.TypesInfo.Types {
		switch x := expr.(type) {
		case *ast.CompositeLit:
			if types.Identical(tv.Type, analyser.root.Type) {
				analyser.processExpression(pkgi, nil, expr, path)
			}
		case *ast.SelectorExpr:
			if pkg.TypesInfo.Uses[x.Sel] == analyser.root.Field && tv.Assignable() {
				analyser.processReference(pkgi, analyser.root.Field, expr, path.Consume(1))
			}
		}
	}
}
