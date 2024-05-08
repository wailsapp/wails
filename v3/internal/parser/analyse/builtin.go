package analyse

import (
	"go/ast"
	"go/types"

	"github.com/pterm/pterm"
)

// processBuiltinSink handles a call to a builtin function
// where one of the parameters references a target variable.
//
// If the callee is copy and the target variable is referenced
// as the destination (first parameter), the source expression
// (second parameter) is fed to [Analyser.processExpression]
// as if the call was an assignment of the source to the target.
//
// If the callee is copy, the target variable is referenced
// as the source (second parameter) and the path has reference semantics,
// the destination (first parameter) is fed to [Analyser.processExpression]
// to be scheduled as an additional target,
// as if the call was an assignment of the reference to the destination.
//
// If the callee is append, the target variable is referenced
// as an element (second or later parameter) and the path
// has reference semantics, processBuiltinSink prepends an indexing step
// to the path and instructs processReference to keep analysing
// the context.
//
// Every other builtin is ignored silently.
//
// processBuiltinSink operates under the assumption
// that the given call expression is valid.
func (analyser *Analyser) processBuiltinSink(pkgi int, callee types.Object, call *ast.CallExpr, paramIndex int, pathIn Path) (path Path, stop bool) {
	pkg := analyser.pkgs[pkgi]
	path = pathIn
	stop = true // Unless set to false, processReference will stop processing the context.

	switch callee {
	case types.Universe.Lookup("copy"):
		if !path.At(0).IsIndexing() {
			// We are not interested elements of this slice.
			return
		}

		if paramIndex == 0 && !path.Consume(1).HasWeakRef() {
			// Call to copy where current expression is dst
			// and we are tracking assignments to its elements: process src.
			analyser.processExpression(pkgi, nil, call.Args[1], path.Set(0, IndexingStep))
		} else if paramIndex == 1 && path.Consume(1).HasRef() {
			// Call to copy where current expression is src
			// and its elements have reference semantics:
			// schedule dst as new weak target.
			analyser.processExpression(pkgi, nil, call.Args[0], path.Consume(1).Weaken().Prepend(WeakIndexingStep))
		}

	case types.Universe.Lookup("append"):
		if paramIndex == 0 && path.At(0).IsIndexing() {
			// Call to append where current expression is the input slice.

			if !path.Consume(1).HasWeakRef() {
				// Process appended elements as assignments
				// to elements of the current expression.
				if call.Ellipsis.IsValid() {
					analyser.processExpression(pkgi, nil, call.Args[1], path.Clone().Set(0, IndexingStep))
				} else {
					for _, el := range call.Args[1:] {
						analyser.processExpression(pkgi, nil, el, path.Clone().Consume(1))
					}
				}
			}

			// Process call expression as weak reference.
			path = path.Consume(1).Weaken().Prepend(WeakIndexingStep)
		} else if paramIndex > 0 && path.HasRef() {
			// Call to append where current expression is among elements
			// and has reference semantics: update path and continue.
			if call.Ellipsis.IsValid() {
				// Elements of the current expression are being appended
				// to another slice: ensure that they have reference semantics.
				if !path.At(0).IsIndexing() || !path.Consume(1).HasRef() {
					return
				}

				// Consume indexing step, it will be added back below.
				path = path.Consume(1)
			}

			// Process call expression as a slice of weak references.
			path = path.Weaken().Prepend(WeakIndexingStep)
		}

		stop = false // Let processReference process the surrounding expression.

	case types.Universe.Lookup("panic"):
		pterm.Warning.Printfln(
			"%s: use of panic and recover to provide bindings is not supported",
			pkg.Fset.Position(call.Pos()),
		)
	}

	return
}

// processBuiltinSource handles a call to a builtin function
// whose return value has been marked as a binding source.
//
// If the callee is append, processBuiltinSource feeds each
// element argument individually to [Analyser.processExpression];
// the slice argument is then returned for further processing.
//
// If the callee is recover, processBuiltinSource emits a warning.
//
// Every other builtin is ignored silently.
func (analyser *Analyser) processBuiltinSource(pkgi int, callee types.Object, call *ast.CallExpr, pathIn Path) (expr ast.Expr, path Path, stop bool) {
	pkg := analyser.pkgs[pkgi]
	path = pathIn
	stop = true // Unless set to false, processExpression will stop processing subexpressions.

	switch callee {
	case types.Universe.Lookup("append"):
		if len(call.Args) == 1 {
			if tuple, ok := pkg.TypesInfo.TypeOf(call.Args[0]).(*types.Tuple); ok {
				// Unique argument is a multi-value expression:
				// use processTuple to analyse each element.
				for i := 1; i < tuple.Len(); i++ {
					analyser.processTuple(pkgi, call.Args, nil, i, tuple.Len(), path.Clone().Consume(1))
				}

				// Analyse initial slice.
				analyser.processTuple(pkgi, call.Args, nil, 0, tuple.Len(), path)
				return
			}
		}

		// Process appended elements.
		if len(call.Args) > 1 {
			if call.Ellipsis.IsValid() {
				// Process variadic argument as is.
				analyser.processExpression(pkgi, nil, call.Args[1], path.Clone())
			} else {
				// Process each element separately.

				if !path.At(0).IsIndexing() {
					// This should not happen...
					panic("unexpected append expression")
				}

				for _, el := range call.Args[1:] {
					analyser.processExpression(pkgi, nil, el, path.Clone().Consume(1))
				}
			}
		}

		// Let processExpression continue with the first argument.
		expr = call.Args[0]
		stop = false

	case types.Universe.Lookup("recover"):
		pterm.Warning.Printfln(
			"%s: use of panic and recover to provide bindings is not supported",
			pkg.Fset.Position(call.Pos()),
		)
	}

	return
}
