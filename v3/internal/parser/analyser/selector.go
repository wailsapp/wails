package analyser

import (
	"go/ast"
	"go/types"

	"github.com/pterm/pterm"
)

// processSelectorSink handles a selector expression
// whose receiver references a target variable or field.
func (analyser *Analyser) processSelectorSink(pkgi int, selExpr *ast.SelectorExpr, recv ast.Expr, pathIn Path) (path Path, stop bool) {
	pkg := analyser.pkgs[pkgi]
	path = pathIn
	stop = true // Unless set to false, processReference will stop processing the context.

	if recv != selExpr.X {
		// Current expr is the field identifier: continue.
		stop = false
		return
	}

	sel, ok := pkg.TypesInfo.Selections[selExpr]
	if !ok {
		// Selection is not valid: ignore silently.
		return
	}

	index := sel.Index()

	// Consume matching initial segment from path and continue.
	for ; len(index) > 1; index = index[1:] {
		if path.At(0).IsIndirection() {
			// Allow one automatic pointer indirection
			// for each field selection step.
			path = path.Consume(1)
		}

		if path.At(0) != Step(index[0]) {
			// Selection index does not match path: stop here.
			return
		}

		path = path.Consume(1)
	}

	if sel.Kind() == types.FieldVal {
		// Selector selects a struct field, consume last field index.
		if path.At(0).IsIndirection() {
			// Allow one automatic pointer indirection
			// for each field selection step.
			path = path.Consume(1)
		}

		if path.At(0) != Step(index[0]) {
			// Selection index does not match path: stop here.
			return
		}

		path = path.Consume(1)
	} else {
		// Selector selects a concrete or abstract method.
		if !path.HasRef() {
			// Receiver path does not have reference semantics: stop here.
			return
		}

		// Check receiver type.
		recv := sel.Obj().(*types.Func).Type().(*types.Signature).Recv()
		if recv != nil && types.IsInterface(recv.Type()) && path.At(0) != TypeAssertionStep {
			// Selector selects an abstract method.
			// Interface targets may occur in four situations:
			//   - when resolving the type of a binding expression;
			//   - when resolving an abstract method;
			//   - when a concrete target has been assigned
			//     to an interface variable or field;
			//   - when a binding expression contains a type assertion.
			// We only care about the last two cases,
			// which are characterised by the presence
			// of a type assertion on the path.
			return
		}

		// Reprocess selector as a method expression.
		// processExpression will resolve the method and schedule
		// the receiver as an additional target.
		analyser.processExpression(pkgi, nil, selExpr, path.Weaken().Prepend(ParamStep, -1))
		return
	}

	stop = false
	return
}

// processSelectorSource handles a selector expression
// that has been marked as a binding source.
//
// If the expression is a qualified identifier,
// processExpression is instructed to process the identifier part.
//
// If the expression resolves to a concrete method,
// it is feeded to [Analyser.processFunc].
//
// Otherwise, the type tree is visited to detect
// implicit pointer indirections and the resulting
// sequence of access operations is recorded on path
// and processExpression is instructed to process the receiver expression.
func (analyser *Analyser) processSelectorSource(pkgi int, selExpr *ast.SelectorExpr, pathIn Path) (expr ast.Expr, path Path, stop bool) {
	pkg := analyser.pkgs[pkgi]
	path = pathIn
	stop = true // Unless set to false, processExpression will stop processing subexpressions.

	sel, ok := pkg.TypesInfo.Selections[selExpr]
	if !ok {
		// Selector is a qualified identifier: process the identifier.
		expr = selExpr.Sel
		stop = false
		return
	}

	if sel.Kind() == types.FieldVal {
		// Selector selects a struct field.

		if sel.Obj().(*types.Var) == analyser.root.Field && len(path) == 1 && path.At(0).IsIndexing() {
			// Ignore Bind field of struct application.Options,
			// it is being tracked already.
			return
		}
	} else {
		// Selector selects a concrete or abstract method.

		if !path.At(0).IsFunc() {
			// This should not happen...
			panic("unexpected method selector")
		}

		fn := sel.Obj().(*types.Func)
		recv := fn.Type().(*types.Signature).Recv()

		// Check receiver type.
		if recv != nil && types.IsInterface(recv.Type()) {
			// Selector selects an abstract method.

			if sel.Kind() == types.MethodExpr {
				// Interface method expr not supported: emit a warning and stop here.
				pterm.Warning.Printfln(
					"%s: method expressions with interface receiver are not supported",
					pkg.Fset.Position(selExpr.Pos()),
				)
				return
			}

			// Prepend method lookup step to path.
			path = path.Prepend(MethodLookupStep, analyser.paths.MethodStep(fn.Pkg(), fn.Name()))
		} else {
			// Selector selects a concrete method: handle it right away.

			if sel.Kind() == types.MethodExpr && path.At(0) == ParamStep {
				// Selector is a method expression.
				// The method receiver is treated as a regular argument here,
				// hence parameter index 0 actually points to the receiver,
				// index 1 to the first parameter, and so on.
				// Decrement parameter index to bring it in line
				// with instanced method calls.
				path.Set(1, path.At(1)-1)

				if path.At(1) == -1 && len(sel.Index()) > 1 {
					// Embedded method expr with receiver target is not supported:
					// emit a warning and stop here.
					pterm.Warning.Printfln(
						"%s: receiver of embedded method expression has been marked as a binding source: this is not supported",
						pkg.Fset.Position(selExpr.Pos()),
					)
					return
				}
			}

			analyser.processFunc(pkgi, selExpr.Pos(), fn, path)
			return
		}
	}

	// Visit type tree to detect implicit pointer dereferences
	// and record them on path.

	index := sel.Index()
	prefix := make([]Step, 0, 2*len(sel.Index()))

	// Preserve path strength.
	indirection := IndirectionStep
	if path.HasWeakRef() {
		indirection = WeakIndirectionStep
	}

	recv := sel.Recv()
	if ptr, ok := recv.Underlying().(*types.Pointer); ok {
		// Record one implicit pointer dereference before next selection.
		recv = ptr.Elem()
		prefix = append(prefix, indirection)
	}

	for ; len(index) > 1; index = index[1:] {
		recv = recv.Underlying().(*types.Struct).Field(index[0]).Type()
		prefix = append(prefix, Step(index[0]))

		if ptr, ok := recv.(*types.Pointer); ok {
			// Record one implicit pointer dereference before next selection.
			recv = ptr.Elem()
			prefix = append(prefix, indirection)
		}
	}

	if !types.IsInterface(recv) {
		// Selector selects a struct field, append last field index to prefix.
		prefix = append(prefix, Step(index[0]))
	}

	// Record selection on path and analyse receiver.
	path = path.Prepend(prefix...)
	expr = selExpr.X
	stop = false
	return
}
