package analyse

import (
	"go/ast"
	"go/types"
	"slices"
)

// processLiteralSink handles a composite literal
// where one fields references a target variable or field.
//
// It validates the expression, prepends the appropriate
// indexing or selection steps to the path
// and instructs processReference to keep analysing the context.
func (analyser *Analyser) processLiteralSink(pkgi int, lit *ast.CompositeLit, el ast.Expr, pathIn Path) (path Path, stop bool) {
	pkg := analyser.pkgs[pkgi]
	path = pathIn
	stop = true // Unless set to false, processReference will stop processing the surrounding expression.

	// Retrieve composite literal type
	typ := pkg.TypesInfo.TypeOf(lit)

	if !IsValidType(typ) || !path.HasRef() {
		// Literal is invalid or path does not have
		// reference semantics: stop here.
		return
	}

	index := slices.Index(lit.Elts, el)
	if index < 0 {
		// Current expr is not an element of the literal: stop here.
		return
	}

	// Compute appropriate access steps.
	step := WeakIndexingStep

	switch u := typ.Underlying().(type) {
	case *types.Struct:
		if kv, ok := el.(*ast.KeyValueExpr); ok {
			// Keyed literal: find field index in struct.
			field := pkg.TypesInfo.Uses[kv.Key.(*ast.Ident)]
			for index = 0; index < u.NumFields(); index++ {
				if u.Field(index) == field {
					break
				}
			}

			if index == u.NumFields() {
				// Something's wrong...
				panic("could not find keyed field index in struct")
			}
		}
		// For unkeyed literals, field index is equal to index in literal.

		if types.Identical(typ, analyser.root.Type) && index == analyser.root.Index && len(path) == 1 && path.At(0).IsIndexing() {
			// Reference is being assigned to field application.Options.Bind.
			// We already track this by other means.
			return
		}

		step = Step(index)

	case *types.Array:
		step = ArrayIndexingStep
	}

	// Current expr is an element of the literal:
	// prepend indexing or selection step
	// to weakened path and continue.
	path = path.Weaken().Prepend(step)

	stop = false // Let processReference continue.
	return
}

// processLiteralSource handles a composite literal
// where one or more fields have been marked as binding sources.
//
// It extracts element expressions selected by the first step
// on the given path and feeds them to [Analyser.processExpression].
func (analyser *Analyser) processLiteralSource(pkgi int, exprType types.Type, lit *ast.CompositeLit, pathIn Path) (expr ast.Expr, path Path, stop bool) {
	pkg := analyser.pkgs[pkgi]
	path = pathIn
	stop = true // Unless set to false, processExpression will stop processing subexpressions.

	if len(lit.Elts) == 0 {
		// Empty literal: stop here.
		// Zero values yield no result anyways.
		return
	}

	switch typ := exprType.Underlying().(type) {
	case *types.Array, *types.Slice, *types.Map:
		// We are indexing an array, slice or map:
		// process all elements.
		if !path.At(0).IsIndexing() {
			// This should not happen...
			panic("unexpected array, slice or map literal")
		}

		// Consume indexing step.
		path = path.Consume(1)

		for i := 0; i < len(lit.Elts)-1; i++ {
			switch el := lit.Elts[i].(type) {
			case *ast.KeyValueExpr:
				analyser.processExpression(pkgi, nil, el.Value, path.Clone())
			default:
				analyser.processExpression(pkgi, nil, el, path.Clone())
			}
		}

		// Let processExpression continue with the last element.
		switch el := lit.Elts[len(lit.Elts)-1].(type) {
		case *ast.KeyValueExpr:
			expr = el.Value
		default:
			expr = el
		}

	case *types.Struct:
		if !path.At(0).IsSelection() || int(path.At(0)) >= typ.NumFields() {
			// This should not happen...
			panic("unexpected struct literal")
		}

		// Save and consume selection step.
		index := int(path.At(0))
		path = path.Consume(1)

		// For unkeyed literals we can index the element slice directly.
		if _, keyed := lit.Elts[0].(*ast.KeyValueExpr); !keyed {
			expr = lit.Elts[index]
		} else {
			// Retrieve selected field.
			field := typ.Field(index)

			for _, el := range lit.Elts {
				kv := el.(*ast.KeyValueExpr)
				if pkg.TypesInfo.ObjectOf(kv.Key.(*ast.Ident)) == field {
					expr = kv.Value
					break
				}
			}

			if expr == nil {
				// This should not happen...
				panic("required struct field not found in keyed literal")
			}
		}
	}

	stop = false // Let processExpression continue.
	return
}
