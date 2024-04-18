package parser

import (
	"errors"
	"go/ast"
	"go/token"
	"go/types"
	"slices"

	"github.com/pterm/pterm"
	"github.com/samber/lo"
	"golang.org/x/exp/maps"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
)

// WailsAppOptions is the import path of Wails v3's application.Options
const WailsAppOptions = WailsAppPkgPath + ".Options"

var ErrBadApplicationOptions = errors.New("could not find field Bind of struct application.Options: is the wails v3 module properly installed?")

// FindServices finds all named types that are listed in the Bind field
// of the Wails application.Options struct in the given packages.
//
// each pkg must import the Wails application package (github.com/wailsapp/wails/v3/pkg/application).
//
// As a precondition, the Syntax field of each package object in pkgs
// must be sorted by file start position, ascending. The condition
// is satisfied by default by package objects returned by [LoadPackages].
//
// The resulting slice contains no duplicate elements.
//
// FindServices supports only some kind of expressions,
// but emits warnings for unsupported forms.
func FindServices(pkgs []*packages.Package) ([]*types.TypeName, error) {
	if len(pkgs) == 0 {
		return []*types.TypeName{}, nil
	}

	// Analyse uses of the field application.Options.Bind

	analyser, err := newBindFieldAnalyser(pkgs[0])
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		analyser.Analyse(pkg)
	}

	// If the analysis found new global variables that contribute to bindings,
	// process them. Keep checking until they are exhausted. Note that this
	// loop must eventually terminate, because there are only finitely many
	// globals in a program and we schedule them at most once.
	for {
		globals := analyser.ScheduledGlobals()
		if len(globals) == 0 {
			break
		}

		for _, global := range globals {
			inPkgs := false
			globalAnalyser := newVarAnalyser(analyser, global)

			for _, pkg := range pkgs {
				if global.Pkg() != nil && global.Pkg().Path() == pkg.Types.Path() {
					inPkgs = true
				}
				globalAnalyser.Analyse(pkg)
			}

			if !inPkgs {
				// The current global does not belong to a package under analysis
				pterm.Warning.Printfln(
					"global %s contributes to bindings, but its declaring package is not under analysis: assignments from that package will not be tracked",
					types.ObjectString(global, nil),
				)
			}
		}
	}

	return analyser.Result(), nil
}

// bindingAnalyser tracks the information required during analysis
// of a struct field or variable that contributes to bindings,
// as well as the results of the analysis.
type bindingAnalyser struct {
	fieldOrVar  *types.Var
	fieldStruct types.Type
	fieldIndex  int

	// found tracks all service types found so far
	found map[*types.TypeName]bool

	// scheduled tracks global variables that need to be analysed in the next round
	scheduled map[*types.Var]bool

	// visited tracks local and global variables that do be inspected further
	visited map[*types.Var]bool
}

// newBindFieldAnalyser constructs an analyser
// for the struct field application.Options.Bind.
//
// pkg must import the Wails application package.
func newBindFieldAnalyser(pkg *packages.Package) (*bindingAnalyser, error) {

	// Retrieve type of application.Options
	var options *types.Named
	for _, tv := range pkg.TypesInfo.Types {
		named, ok := tv.Type.(*types.Named)
		if ok && tv.Type.String() == WailsAppOptions {
			options = named
			break
		}
	}
	if options == nil {
		return nil, ErrBadApplicationOptions
	}

	// Retrieve underlying struct
	fieldStruct, ok := options.Underlying().(*types.Struct)
	if !ok {
		return nil, ErrBadApplicationOptions
	}

	// Retrieve field Options.Bind
	var field *types.Var
	fieldIndex := 0
	for ; fieldIndex < fieldStruct.NumFields(); fieldIndex++ {
		field = fieldStruct.Field(fieldIndex)
		if field.Name() == "Bind" {
			break
		}
	}
	if field.Name() != "Bind" {
		return nil, ErrBadApplicationOptions
	}

	return &bindingAnalyser{
		field,
		fieldStruct,
		fieldIndex,
		make(map[*types.TypeName]bool),
		make(map[*types.Var]bool),
		make(map[*types.Var]bool),
	}, nil
}

// newVarAnalyser constructs an analyser for the variable described by v.
// Results are stored in the parent analyser.
func newVarAnalyser(parent *bindingAnalyser, v *types.Var) *bindingAnalyser {
	return &bindingAnalyser{
		v,
		nil,
		-1,
		parent.found,
		parent.scheduled,
		parent.visited,
	}
}

// Result retrieves all service types found so far.
// The resulting slice contains no duplicate elements.
func (analyser *bindingAnalyser) Result() []*types.TypeName {
	return lo.Keys(analyser.found)
}

// ScheduledGlobals collects into a slice the globals
// that are currently scheduled for analysis,
// marks them as visited and returns them.
// The resulting slice contains no duplicate elements.
func (analyser *bindingAnalyser) ScheduledGlobals() []*types.Var {
	scheduled := lo.Keys(analyser.scheduled)
	maps.Copy(analyser.visited, analyser.scheduled)
	clear(analyser.scheduled)
	return scheduled
}

// scheduleGlobal schedules a global variable for analysis
// unless it has been marked as visited.
func (analyser *bindingAnalyser) scheduleGlobal(global *types.Var) {
	if !analyser.visited[global] {
		analyser.scheduled[global] = true
	}
}

// Analyse analyses all assignments performed in pkg to the variable or field
// for which the analyser has been configured.
// The field or variable is marked as visited to break cycles (typically
// introduced by appends).
func (analyser *bindingAnalyser) Analyse(pkg *packages.Package) {
	analyser.visited[analyser.fieldOrVar] = true

	// Visit references to the field or variable under analysis
	for ident, obj := range pkg.TypesInfo.Uses {
		if obj != analyser.fieldOrVar {
			continue
		}

		path := FindAstPath(pkg, ident.Pos(), ident.End())
		if len(path) == 0 {
			pterm.Warning.Printfln(
				"%s: no source file found for reference to field application.Options.Bind in package \"%s\"",
				pkg.Fset.PositionFor(ident.Pos(), true),
				pkg.PkgPath,
			)
			continue
		}

		analyser.analyseReference(pkg, path)
	}

	if analyser.fieldStruct == nil {
		// We are analysing a variable
		if analyser.fieldOrVar.Pos().IsValid() && analyser.fieldOrVar.Pkg() != nil && analyser.fieldOrVar.Pkg().Path() == pkg.Types.Path() {
			// The variable belongs to the current package: analyse its declaration

			path := FindAstPath(pkg, analyser.fieldOrVar.Pos(), analyser.fieldOrVar.Pos())

			// No need to reparen, as identifiers in variable declarations
			// may not be enclosed in parentheses
			if len(path) > 1 {
				var value ast.Expr

				switch decl := path[1].(type) {
				case *ast.ValueSpec:
					index := slices.Index(decl.Names, path[0].(*ast.Ident))
					if index >= 0 && index < len(decl.Values) {
						value = decl.Values[index]
					}

				case *ast.AssignStmt:
					index := slices.Index(decl.Lhs, path[0].(ast.Expr))
					if index >= 0 && index < len(decl.Rhs) {
						value = decl.Rhs[index]
					}
				}

				if value != nil {
					analyser.analyseFieldExpr(pkg, value)
				}
			}
		}
	} else {
		// Analyse expressions whose type is the struct containing the field
		// under analysis. We need this step to catch unkeyed struct literals.
		for expr, tv := range pkg.TypesInfo.Types {
			if !tv.IsValue() || !types.Identical(tv.Type, analyser.fieldStruct) {
				continue
			}

			analyser.analyseStructExpr(pkg, expr)
		}
	}
}

// analyseReference explores the context of references to the field
// or variable under analysis, looking for assignments, copies,
// and element assignments.
func (analyser *bindingAnalyser) analyseReference(pkg *packages.Package, path []ast.Node) {
	// Remember whether we traversed a slice expression
	sliceExpr := false

	// Enter a loop to avoid recursion when possible
	for {
		path = Reparen(path)
		if len(path) < 2 {
			return
		}

		switch ctx := path[1].(type) {
		case *ast.KeyValueExpr:
			// The context is a KeyValue expr: if our field is the key,
			// analyse the assigned value; otherwise ignore it silently.
			// If the current target is not a field, this is a map entry
			// and must be ignored.

			if ctx.Key != path[0] || !analyser.fieldOrVar.IsField() {
				return
			}

			analyser.analyseFieldExpr(pkg, ctx.Value)
			return

		case *ast.SelectorExpr:
			// The context is a selector expr: if it selects our field or var,
			// travel one step up the path and repeat the analysis
			if ctx.Sel != path[0] {
				return
			}

			path = path[1:]
			continue // Recurse

		case *ast.AssignStmt:
			// If the reference is assigned to, analyse its value;
			// otherwise it is an assignee and we ignore it silently

			if sliceExpr {
				// Ignore invalid assignment to a slice expression
				return
			}

			index := slices.Index(ctx.Lhs, path[0].(ast.Expr))
			if index < 0 || index >= len(ctx.Rhs) {
				return
			}

			analyser.analyseFieldExpr(pkg, ctx.Rhs[index])
			return

		case *ast.CallExpr:
			// If the reference is the destination of a call to the copy builtin,
			// analyse the source

			if len(ctx.Args) < 1 || ctx.Args[0] != path[0] {
				return
			}

			callee := typeutil.Callee(pkg.TypesInfo, ctx)
			if callee == nil || callee != types.Universe.Lookup("copy") {
				return
			}

			analyser.analyseFieldExpr(pkg, ctx.Args[1])
			return

		case *ast.IndexExpr:
			// If the reference is the subject of an index expression,
			// check whether the element is being assigned to:
			// if yes, analyse the value

			if ctx.X != path[0] || len(path) < 3 {
				return
			}

			stmt, ok := path[2].(*ast.AssignStmt)
			if !ok {
				return
			}

			index := slices.Index(stmt.Lhs, ast.Expr(ctx))
			if index < 0 || index >= len(stmt.Rhs) {
				return
			}

			analyser.analyseServiceExpr(pkg, stmt.Rhs[index])
			return

		case *ast.SliceExpr:
			// A slice expression might be the destination of a call to copy,
			// so we keep traversing the path

			sliceExpr = true
			path = path[1:]
			continue // Recurse

		case *ast.UnaryExpr:
			// We cannot track indirect assignments through pointers:
			// if the address of the field or variable is taken, emit a warning

			if ctx.Op == token.AND {
				pterm.Warning.Printfln(
					"%s: address of field or variable under analysis taken here: indirect assignments will not be tracked",
					pkg.Fset.PositionFor(ctx.Pos(), true),
				)
			}

			return
		}

		return
	}
}

// analyseStructExpr analyses expressions
// whose type is the struct application.Options,
// catching assignments to the Bind field through unkeyed literals.
//
// If it detects an explicit conversion to application.Options
// from another type, it emits a warning.
func (analyser *bindingAnalyser) analyseStructExpr(pkg *packages.Package, expr ast.Expr) {
	switch x := expr.(type) {
	case *ast.CompositeLit:
		// Ignore literals with keyed fields

		for _, field := range x.Elts {
			if _, ok := field.(*ast.KeyValueExpr); ok {
				return
			}
		}

		if analyser.fieldIndex < len(x.Elts) {
			analyser.analyseFieldExpr(pkg, x.Elts[analyser.fieldIndex])
		}

	case *ast.CallExpr:
		// Detect explicit conversions from other types and emit a warning

		tv, ok := pkg.TypesInfo.Types[x.Fun]
		if !ok || !tv.IsType() {
			return
		}

		if len(x.Args) != 1 {
			return
		}

		tv, ok = pkg.TypesInfo.Types[x.Args[0]]
		if !ok || !tv.IsValue() || tv.Type == nil || types.Identical(tv.Type, types.Typ[types.Invalid]) || types.Identical(tv.Type, analyser.fieldStruct) {
			return
		}

		pterm.Warning.Printfln(
			"%s: unsupported conversion from type %s to %s",
			pkg.Fset.PositionFor(expr.Pos(), true),
			tv.Type,
			analyser.fieldStruct,
		)
	}
}

// analyseFieldExpr visits expressions assigned directly or indirectly
// to the fields and variables under analysis.
func (analyser *bindingAnalyser) analyseFieldExpr(pkg *packages.Package, expr ast.Expr) {
	// Enter a loop to avoid recursion where possible
	for {
		// Unwrap expression and check its type
		expr = ast.Unparen(expr)

		tv, ok := pkg.TypesInfo.Types[expr]
		if !ok || !tv.IsValue() || tv.Type == nil || types.Identical(tv.Type, types.Typ[types.Invalid]) || !types.AssignableTo(tv.Type, analyser.fieldOrVar.Type()) {
			// Ignore invalid expressions or implicit conversions
			// Do not emit a warning as this is a typing error anyways
			return
		}

		// If the expression is a selector, extract its terminal identifier
		if sel, ok := expr.(*ast.SelectorExpr); ok {
			expr = sel.Sel
		}

		switch x := expr.(type) {
		case *ast.Ident:
			// If the expression is an identifier, find the local, global
			// or field it refers to and schedule it for analysis,
			// unless it is already marked as visited

			switch obj := pkg.TypesInfo.ObjectOf(x).(type) {
			case nil, *types.Nil:
				// Silently ignore nil values and undefined identifiers
				return

			case *types.Var:
				if analyser.visited[obj] {
					return
				}

				if obj.IsField() {
					// Ignore and report extraneous struct fields,
					// as tracking them is too much work.
					break
				}

				if !obj.Exported() || (obj.Parent() != nil && obj.Parent().Parent() != nil && obj.Parent().Parent() != types.Universe) {
					// obj is a local or unexported variable, analyse it right away
					newVarAnalyser(analyser, obj).Analyse(pkg)
				} else {
					// obj is a global variable, schedule it for later as it needs
					analyser.scheduleGlobal(obj)
				}

				return
			}

		case *ast.CompositeLit:
			// If the expression is a composite literal, analyse the elements

			for _, elt := range x.Elts {
				if kv, ok := elt.(*ast.KeyValueExpr); ok {
					analyser.analyseServiceExpr(pkg, kv.Value)
				} else {
					analyser.analyseServiceExpr(pkg, elt)
				}
			}

			return

		case *ast.SliceExpr:
			// If the expression is a slice expr, analyse its subject

			expr = x.X
			continue // Recurse

		case *ast.CallExpr:
			// If the expression is a call expr, handle functions
			// whose behaviour is known and standardised

			callee := typeutil.Callee(pkg.TypesInfo, x)
			if callee == nil {
				break
			}

			// Handle make and append builtins
			switch callee {
			case types.Universe.Lookup("make"):
				return
			case types.Universe.Lookup("append"):
				if len(x.Args) > 1 {
					if x.Ellipsis.IsValid() {
						if len(x.Args) > 2 {
							return
						}

						analyser.analyseFieldExpr(pkg, x.Args[1])
					} else {
						for _, arg := range x.Args[1:] {
							analyser.analyseServiceExpr(pkg, arg)
						}
					}
				}

				if len(x.Args) > 0 {
					expr = x.Args[0]
					continue // Recurse
				}

				return
			}

			if callee.Pkg() == nil || callee.Pkg().Path() != "slices" || callee.Parent() != callee.Pkg().Scope() {
				break
			}

			// Handle functions from the slices standard package
			switch callee.Name() {
			case "Clip", "Clone", "Compact", "CompactFunc", "Delete", "DeleteFunc", "Grow":
				if len(x.Args) < 1 {
					return
				}

				expr = x.Args[0]
				continue // Recurse

			case "Concat":
				if x.Ellipsis.IsValid() && len(x.Args) == 1 {
					break
				}

				if len(x.Args) < 1 {
					return
				}

				for _, arg := range x.Args[1:] {
					analyser.analyseFieldExpr(pkg, arg)
				}

				expr = x.Args[0]
				continue // Recurse

			case "Insert":
				if len(x.Args) > 2 {
					if x.Ellipsis.IsValid() {
						if len(x.Args) > 3 {
							return
						}

						analyser.analyseFieldExpr(pkg, x.Args[2])
					} else {
						for _, arg := range x.Args[2:] {
							analyser.analyseServiceExpr(pkg, arg)
						}
					}
				}

				if len(x.Args) > 0 {
					expr = x.Args[0]
					continue // Recurse
				}

				return

			case "Replace":
				if len(x.Args) > 3 {
					if x.Ellipsis.IsValid() {
						if len(x.Args) > 4 {
							return
						}

						analyser.analyseFieldExpr(pkg, x.Args[3])
					} else {
						for _, arg := range x.Args[3:] {
							analyser.analyseServiceExpr(pkg, arg)
						}
					}
				}

				if len(x.Args) > 0 {
					expr = x.Args[0]
					continue // Recurse
				}

				return
			}
		}

		// Report failure
		pterm.Warning.Printfln(
			"%s: ignoring unsupported expression assigned (directly or indirectly) to field application.Options.Bind",
			pkg.Fset.PositionFor(expr.Pos(), true),
		)

		return
	}
}

// analyseServiceExpr deduces the type of expressions used as elements
// of the Bind slice, verifies them and records them
func (analyser *bindingAnalyser) analyseServiceExpr(pkg *packages.Package, expr ast.Expr) {
	// Unwrap conversions and assertions to interface types:
	// we might be able to reach an expression with concrete type
unwrap:
	for {
		expr = ast.Unparen(expr)

		switch x := expr.(type) {
		case *ast.CallExpr:
			tv, ok := pkg.TypesInfo.Types[x.Fun]
			if !ok || !tv.IsType() {
				break unwrap
			}

			if _, ok := tv.Type.Underlying().(*types.Interface); !ok {
				return
			}

			if len(x.Args) != 1 {
				break unwrap
			}

			expr = x.Args[0]

		case *ast.TypeAssertExpr:
			tv, ok := pkg.TypesInfo.Types[x.Type]
			if !ok || !tv.IsType() {
				break unwrap
			}

			if _, ok := tv.Type.Underlying().(*types.Interface); !ok {
				return
			}

			expr = x.X

		default:
			break unwrap
		}
	}

	// Retrieve the type of the expression
	tv, ok := pkg.TypesInfo.Types[expr]
	if !ok {
		pterm.Warning.Printfln(
			"%s: ignoring invalid service expression",
			pkg.Fset.PositionFor(expr.Pos(), true),
		)
		return
	}

	// Unalias and validate type
	typ := types.Unalias(tv.Type)
	if typ == nil || types.Identical(typ, types.Typ[types.Invalid]) {
		pterm.Warning.Printfln(
			"%s: ignoring service expression with invalid type",
			pkg.Fset.PositionFor(expr.Pos(), true),
		)
		return
	}

	// Accept nil silently
	if types.Identical(typ, types.Typ[types.UntypedNil]) {
		return
	}

	// Ensure the concrete type is a pointer
	ptr, ok := typ.(*types.Pointer)
	if !ok {
		pterm.Warning.Printfln(
			"%s: ignoring service expression with non-pointer type %s",
			pkg.Fset.PositionFor(expr.Pos(), true),
			typ,
		)
		return
	}

	// Retrieve, unalias and validate the element type
	elem := types.Unalias(ptr.Elem())
	if elem == nil || types.Identical(elem, types.Typ[types.Invalid]) {
		pterm.Warning.Printfln(
			"%s: ignoring service expression with invalid element type",
			pkg.Fset.PositionFor(expr.Pos(), true),
		)
		return
	}

	// Ensure the element type is not anonymous
	srv, ok := elem.(*types.Named)
	if !ok {
		pterm.Warning.Printfln(
			"%s: ignoring service expression with anonymous or basic element type %s",
			pkg.Fset.PositionFor(expr.Pos(), true),
			elem,
		)
		return
	}

	// Validate named type
	srvu := srv.Underlying()
	if srvu == nil || types.Identical(srvu, types.Typ[types.Invalid]) {
		pterm.Warning.Printfln(
			"%s: ignoring service expression with invalid element type",
			pkg.Fset.PositionFor(expr.Pos(), true),
		)
		return
	}

	// Ensure it is a struct type
	if _, ok := srvu.(*types.Struct); !ok {
		pterm.Warning.Printfln(
			"%s: ignoring service expression with non-struct element type %s",
			pkg.Fset.PositionFor(expr.Pos(), true),
			elem,
		)
		return
	}

	// Ensure it is not generic
	if srv.TypeParams() != nil {
		pterm.Warning.Printfln(
			"%s: ignoring service expression with generic element type %s",
			pkg.Fset.PositionFor(expr.Pos(), true),
			elem,
		)
		return
	}

	// Record service type
	analyser.found[srv.Obj()] = true
}
