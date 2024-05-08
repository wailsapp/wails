package analyse

import (
	"go/ast"
	"go/types"
	"slices"
)

// processSliceFuncSink handles a call to a function
// from the standard slices package
// where one of the parameters references a target variable or field.
//
// processSliceFuncSink operates under the assumption
// that the given call expression is valid.
func (analyser *Analyser) processSliceFuncSink(pkgi int, callee types.Object, call *ast.CallExpr, paramIndex int, pathIn Path) (path Path, stop bool) {
	path = pathIn
	stop = true // Unless set to false, processReference will stop processing the context.

	targetSecondParam := false

	switch name := callee.Name(); name {
	case "Clip", "Clone", "Compact", "Delete", "Grow":
		// Call to Clip, Clone, Compact, Delete, or Grow
		// where current expression is the input slice.
		if paramIndex != 0 || !path.At(0).IsIndexing() {
			return
		}

		// Process call expression as a weak reference.
		path = path.Weaken()
		stop = false // Let processReference process the surrounding expression.

	case "Concat":
		if call.Ellipsis.IsValid() {
			// Call to Concat where current expression is a slice of slices
			// to concatenate.
			if !path.At(0).IsIndexing() || !path.Consume(1).HasRef() {
				return
			}

			// Consume indexing step and process
			// call expression as a weak reference.
			path = path.Consume(1).Weaken()
		} else {
			// Call to Concat where current expression is one of the slices
			// to concatenate.
			if !path.HasRef() {
				return
			}

			// Process call expression as a weak reference.
			path = path.Weaken()
		}

		stop = false // Let processReference process the surrounding expression.

	case "Insert", "Replace": // Handled just like the append builtin (see processBuiltinSink).
		variadicArgIndex := 2
		if name == "Replace" {
			variadicArgIndex = 3
		}

		if paramIndex == 0 && path.At(0).IsIndexing() {
			// Call to Insert or Replace where current expression
			// is the input slice.
			if !path.Consume(1).HasWeakRef() {
				// Process inserted or replacement elements as assignments
				// to elements of the current expression.
				if call.Ellipsis.IsValid() {
					analyser.processExpression(pkgi, nil, call.Args[variadicArgIndex], path.Clone().Set(0, IndexingStep))
				} else {
					for _, el := range call.Args[variadicArgIndex:] {
						analyser.processExpression(pkgi, nil, el, path.Clone().Consume(1))
					}
				}
			}

			// Process call expression as weak reference.
			path = path.Weaken()
		} else if paramIndex >= variadicArgIndex && path.HasRef() {
			// Call to Insert or Replace where current expression is
			// among elements and has reference semantics:
			// update path and continue.
			if call.Ellipsis.IsValid() {
				// Elements of the current expression are being copied
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

	case "MaxFunc", "MinFunc":
		// Call to MaxFunc or MinFunc where current expression
		// is the input slice: if elements of the slice have
		// reference semantics, we must target predicate function params
		// and process the call expression.
		if paramIndex != 0 || !path.At(0).IsIndexing() || !path.Consume(1).HasRef() {
			return
		}

		// Consume indexing step and weaken element reference.
		path = path.Consume(1).Weaken()

		paramPath := path.CloneAndGrow(2)
		analyser.processExpression(pkgi, nil, call.Args[1], paramPath.Prepend(ParamStep, 0).Clone())
		analyser.processExpression(pkgi, nil, call.Args[1], paramPath.Prepend(ParamStep, 1))

		stop = false // Let processReference process the surrounding expression.

	case "CompareFunc", "EqualFunc":
		// We must target first or second parameter of predicate
		// depending on the position of the current expression.
		if paramIndex == 1 {
			targetSecondParam = true
			paramIndex = 0 // Fool next case into accepting param index.
		}
		fallthrough

	case "BinarySearchFunc":
		// Call to CompareFunc, EqualFunc or BinarySearch
		// where current expression is one of the input slices:
		// if elements of the slice have reference semantics
		// we must target predicate function params.
		if paramIndex != 0 || !path.At(0).IsIndexing() || !path.Consume(1).HasRef() {
			return
		}

		var targetParam Step = 0
		if targetSecondParam {
			targetParam = 1
		}

		analyser.processExpression(pkgi, nil, call.Args[2], path.Consume(1).Weaken().Prepend(ParamStep, targetParam))

		// Discard return value, it is never a target reference.

	case "CompactFunc", "IsSortedFunc", "SortFunc", "SortStableFunc":
		// These functions use binary predicates.
		targetSecondParam = true
		fallthrough

	case "ContainsFunc", "DeleteFunc", "IndexFunc":
		// Call to CompactFunc, IsSortedFunc, SortFunc, SortStableFunc,
		// ContainsFunc, DeleteFunc or IndexFunc
		// where current expression is the input slice:
		// if elements of the slice have reference semantics
		// we must target predicate function params.
		if paramIndex != 0 || !path.At(0).IsIndexing() || !path.Consume(1).HasRef() {
			return
		}

		// Consume indexing step, weaken element reference and add capacity.
		path = slices.Grow(path.Consume(1).Weaken(), 2)

		if targetSecondParam {
			// Target second parameter when the predicate is binary.
			analyser.processExpression(pkgi, nil, call.Args[1], path.Prepend(ParamStep, 1).Clone())
		}
		analyser.processExpression(pkgi, nil, call.Args[1], path.Prepend(ParamStep, 0))

		// Discard return value, it is never a target reference.
	}

	return
}

// processSliceFuncSource handles a call to a function
// from the standard slices package
// whose return value has been marked as a binding source.
func (analyser *Analyser) processSliceFuncSource(pkgi int, callee types.Object, call *ast.CallExpr, pathIn Path) (expr ast.Expr, path Path, stop bool) {
	pkg := analyser.pkgs[pkgi]
	path = pathIn
	stop = true // Unless set to false, processReference will stop processing the context.

	switch name := callee.Name(); name {
	case "Clip", "Clone", "Compact":
		expr = call.Args[0]
		stop = false

	case "Delete", "Grow":
		if len(call.Args) == 1 {
			// Unique argument is a multi-value expression:
			// use processTuple to analyse each element.
			argCount := 2
			if name == "Delete" {
				argCount = 3
			}

			analyser.processTuple(pkgi, call.Args, nil, 0, argCount, path)
		}

		expr = call.Args[0]
		stop = false

	case "CompactFunc", "DeleteFunc":
		if len(call.Args) == 1 {
			// Unique argument is a multi-value expression:
			// use processTuple to analyse each element.
			if path.Consume(1).HasRef() {
				if !path.At(0).IsIndexing() {
					// This should not happen...
					panic("unexpected call to slices." + name)
				}

				// Path has reference semantics: target comparison function parameters.
				paramPath := path.Consume(1).CloneAndGrow(2).Weaken()
				if name == "CompactFunc" {
					analyser.processTuple(pkgi, call.Args, nil, 1, 2, paramPath.Prepend(ParamStep, 1).Clone())
				}
				analyser.processTuple(pkgi, call.Args, nil, 1, 2, paramPath.Prepend(ParamStep, 0))
			}

			analyser.processTuple(pkgi, call.Args, nil, 0, 2, path)
			return
		}

		if path.Consume(1).HasRef() {
			if !path.At(0).IsIndexing() {
				// This should not happen...
				panic("unexpected call to slices." + name)
			}

			// Path has reference semantics: target comparison/predicate function parameters.
			paramPath := path.Consume(1).CloneAndGrow(2).Weaken()
			if name == "CompactFunc" {
				analyser.processExpression(pkgi, nil, call.Args[1], paramPath.Prepend(ParamStep, 1).Clone())
			}
			analyser.processExpression(pkgi, nil, call.Args[1], paramPath.Prepend(ParamStep, 0))
		}

		expr = call.Args[0]
		stop = false

	case "Concat":
		if len(call.Args) == 0 {
			return
		}

		if len(call.Args) == 1 {
			if tuple, ok := pkg.TypesInfo.TypeOf(call.Args[0]).(*types.Tuple); ok {
				// Unique argument is a multi-value expression:
				// use processTuple to analyse each element.
				for i := 0; i < tuple.Len()-1; i++ {
					analyser.processTuple(pkgi, call.Args, nil, i, tuple.Len(), path.Clone())
				}
				// Save one clone operation
				analyser.processTuple(pkgi, call.Args, nil, tuple.Len()-1, tuple.Len(), path)
				return
			}
		}

		if call.Ellipsis.IsValid() {
			path = path.Prepend(IndexingStep)
		} else {
			for _, slice := range call.Args[1:] {
				analyser.processExpression(pkgi, nil, slice, path.Clone())
			}
		}

		expr = call.Args[0]
		stop = false

	case "Insert", "Replace":
		variadicArgIndex := 2
		if name == "Replace" {
			variadicArgIndex = 3
		}

		if len(call.Args) == 1 {
			// Unique argument is a multi-value expression:
			// use processTuple to analyse each element.
			tuple := pkg.TypesInfo.TypeOf(call.Args[0]).(*types.Tuple)

			if variadicArgIndex < tuple.Len() {
				if !path.At(0).IsIndexing() {
					// This should not happen...
					panic("unexpected call to slices." + name)
				}

				for i := variadicArgIndex; i < tuple.Len(); i++ {
					analyser.processTuple(pkgi, call.Args, nil, i, tuple.Len(), path.Clone().Consume(1))
				}
			}

			analyser.processTuple(pkgi, call.Args, nil, 0, tuple.Len(), path)
			return
		}

		if variadicArgIndex < len(call.Args) {
			if call.Ellipsis.IsValid() {
				analyser.processExpression(pkgi, nil, call.Args[variadicArgIndex], path.Clone())
			} else {
				if !path.At(0).IsIndexing() {
					// This should not happen...
					panic("unexpected call to slices." + name)
				}

				for _, el := range call.Args[variadicArgIndex:] {
					analyser.processExpression(pkgi, nil, el, path.Clone().Consume(1))
				}
			}
		}

		expr = call.Args[0]
		stop = false

	case "MaxFunc", "MinFunc":
		if len(call.Args) == 1 {
			// Unique argument is a multi-value expression:
			// use processTuple to analyse each element.
			if path.HasRef() {
				// Path has reference semantics: target comparison function parameters.
				paramPath := path.CloneAndGrow(2).Weaken()
				analyser.processTuple(pkgi, call.Args, nil, 1, 2, paramPath.Prepend(ParamStep, 0).Clone())
				analyser.processTuple(pkgi, call.Args, nil, 1, 2, paramPath.Prepend(ParamStep, 1))
			}

			analyser.processTuple(pkgi, call.Args, nil, 0, 2, path.Prepend(IndexingStep))
			return
		}

		if path.HasRef() {
			// Path has reference semantics: target comparison function parameters.
			paramPath := path.CloneAndGrow(2).Weaken()
			analyser.processExpression(pkgi, nil, call.Args[1], paramPath.Prepend(ParamStep, 0).Clone())
			analyser.processExpression(pkgi, nil, call.Args[1], paramPath.Prepend(ParamStep, 1))
		}

		expr = call.Args[0]
		path = path.Prepend(IndexingStep)
		return
	}

	return
}
