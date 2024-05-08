package analyse

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/pterm/pterm"
)

// processExpression analyses an AST expression
// that has been assigned to a target variable or field.
// It finds either a new target to schedule, or the concrete type
// of a binding expression, which is then reported as a result.
//
// If targetType is not nil, processExpression checks that
// the expression be assignable to a variable of that type.
//
// path represents a sequence of access operations (see [Path]) that,
// if applied to the target expression, should yield a binding expression.
// processExpression may alter the contents of the path slice,
// hence the caller should pass in a fresh copy or cede ownership.
//
// processExpression keeps walking _down_ the AST and extending
// the given path appropriately, until one of the following conditions
// is met:
//
//   - the current expr is a reference to a variable: then the variable
//     is scheduled as an additional target with the current path;
//   - the current expr is a static function call: then the relevant
//     result field is scheduled as a target with the current path;
//   - the current expr is a dynamic function call: then the function variable
//     or interface method is analysed to resolve it
//     to one or more possible static calls;
//   - the current expr can be destructured by the head operation
//     on the current path: then the operation is consumed from the path
//     and performed on the current expr, thus obtaining
//     a simpler expression to analyse.
//
// processExpression is a sort of dual to [Analyser.processReference]:
// it processes increasingly simpler expressions,
// accumulating path steps and scheduling target variables
// which are then consumed by processReference
// as it processes increasingly complex expressions,
// consuming path steps in symmetrical fashion.
func (analyser *Analyser) processExpression(pkgi int, targetType types.Type, expr ast.Expr, path Path) {
	pkg := analyser.pkgs[pkgi]

	// Retrieve expression type.
	exprType := pkg.TypesInfo.TypeOf(expr)
	if tv, ok := pkg.TypesInfo.Types[expr]; ok && tv.HasOk() {
		// Comma-ok expression, get type of first result.
		// We only do this for the top-level expression:
		// nested comma-ok expressions always have the type
		// of their first result.
		exprType = tv.Type.(*types.Tuple).At(0).Type()
	}

	// Validate expression type.
	// We only do this for the top-level expression:
	// if it is valid, all sub-expressions are valid as well.
	//
	// Because of this property, processExpression can get by
	// with way less error checks than processReference.
	if targetType == nil {
		targetType = exprType
		if !IsValidType(exprType) {
			return
		}
	} else if !IsAssignableTo(exprType, targetType) {
		return
	}

	var prevExpr ast.Expr

mainloop:
	for expr = ast.Unparen(expr); ; expr, exprType = ast.Unparen(expr), pkg.TypesInfo.TypeOf(expr) {
		// Safeguard against bugs.
		if prevExpr == expr {
			panic("expression analyser loop stuck")
		}
		prevExpr = expr

		if !types.IsInterface(exprType) {
			// The expression has concrete type.

			if types.Identical(exprType, types.Typ[types.UntypedNil]) {
				// Ignore untyped nil. Zero values yield no result anyways.
				return
			}

			if len(path) == 0 {
				// We reached a concrete result type: report and stop.
				analyser.reportResult(pkgi, expr.Pos(), exprType)
				return
			}

			// If possible, take a type assertion off the path.
			if path.At(0) == TypeAssertionStep {
				if path.At(1) != analyser.paths.TypeStep(exprType) {
					// This is not the type we are looking for:
					// ignore silently.
					return
				}

				path = path.Consume(2)
			}

			// If required, perform a method lookup.
			if path.At(0) == MethodLookupStep {
				pkg, name := analyser.paths.MethodStepParams(path.At(1))

				obj, _, indirect := types.LookupFieldOrMethod(exprType, false, pkg, name)
				if fn, ok := obj.(*types.Func); ok && !indirect {
					analyser.processFunc(pkgi, expr.Pos(), fn, path.Consume(2))
				}

				return
			}
		}

		switch x := expr.(type) {
		case *ast.Ident:
			switch obj := pkg.TypesInfo.ObjectOf(x).(type) {
			case *types.Nil:
				// Ignore nil. Zero values yield no result anyways.

			case *types.Var:
				// Identifier is defined and refers to a variable or field.
				// Ignore struct fields.
				if !obj.IsField() {
					analyser.schedule(target{
						pkgi:     varPkgi(analyser, pkgi, obj),
						variable: obj,
						path:     analyser.paths.Get(path),
					})
				}

			case *types.Func:
				// Identifier is defined and refers to a function.
				analyser.processFunc(pkgi, x.Pos(), obj, path)

			default:
				break mainloop
			}

			return

		case *ast.SelectorExpr:
			var stop bool
			expr, path, stop = analyser.processSelectorSource(pkgi, x, path)
			if stop {
				return
			}

		case *ast.IndexListExpr:
			// Current expression is a function instantiation.
			if !path.At(0).IsFunc() {
				// This should not happen...
				panic("unexpected function instantiation")
			}

			// Discard type parameters and analyse function.
			expr = x.X

		case *ast.IndexExpr:
			if pkg.TypesInfo.Types[x.Index].IsType() {
				// Current expression is a function instantiation.
				if !path.At(0).IsFunc() {
					// This should not happen...
					panic("unexpected function instantiation")
				}

				// Discard type parameters.
			} else {
				// Current expression is an indexing operation:
				// prepend appropriate indexing steps.

				switch pkg.TypesInfo.TypeOf(x.X).Underlying().(type) {
				case *types.Pointer:
					// Array indexing step with implicit indirection.
					// Preserve path strength.
					if path.HasWeakRef() {
						path = path.Prepend(WeakIndirectionStep, ArrayIndexingStep)
					} else {
						path = path.Prepend(IndirectionStep, ArrayIndexingStep)
					}

				case *types.Array:
					// Array indexing step.
					path = path.Prepend(ArrayIndexingStep)

				default:
					// Preserve path strength.
					if path.HasWeakRef() {
						path = path.Prepend(WeakIndexingStep)
					} else {
						path = path.Prepend(IndexingStep)
					}
				}
			}

			// Analyse subject function, array or slice.
			expr = x.X

		case *ast.SliceExpr:
			// Current expression is a slicing operation:
			// no need to add indexing steps.
			// If the subject is an array, change a (Weak)IndexingStep
			// to ArrayIndexingStep and prepend an optional indirection.

			switch pkg.TypesInfo.TypeOf(x.X).Underlying().(type) {
			case *types.Pointer:
				// Array slicing step with implicit indirection.
				// Preserve path strength.
				if path.HasWeakRef() {
					path = path.Prepend(WeakIndirectionStep)
				} else {
					path = path.Prepend(IndirectionStep, ArrayIndexingStep)
				}

				if path.At(1).IsIndexing() {
					path.Set(1, ArrayIndexingStep)
				}

			case *types.Array:
				// Array slicing step.
				if path.At(0).IsIndexing() {
					path.Set(0, ArrayIndexingStep)
				}
			}

			// Analyse subject array or slice.
			expr = x.X

		case *ast.StarExpr:
			// Current expression is a pointer indirection operation.
			// Preserve path strength.
			if path.HasWeakRef() {
				path = path.Prepend(WeakIndirectionStep)
			} else {
				path = path.Prepend(IndirectionStep)
			}

			// Analyse subject pointer.
			expr = x.X

		case *ast.UnaryExpr:
			switch x.Op {
			case token.AND:
				// Current expression is an address operation:
				// take an indirection step off the path.

				if !path.At(0).IsIndirection() {
					// This should not happen...
					panic("unexpected address operator")
				}

				// Consume indirection step and analyse subject expr.
				path = path.Consume(1)
				expr = x.X

			case token.ARROW:
				// Current expression is a channel receive operation:
				// analyse sends into the channel.
				path = path.Prepend(ChanSendStep)
				expr = x.X

			default:
				break mainloop
			}

		case *ast.TypeAssertExpr:
			// Current expression is a type assertion.

			if !types.IsInterface(exprType) {
				// Asserted type is concrete:
				// prepend type assertion step to the current path.
				path = path.Prepend(TypeAssertionStep, analyser.paths.TypeStep(exprType))
			}

			// Analyse subject interface value.
			expr = x.X

		case *ast.CompositeLit:
			// Current expression is a composite literal:
			// use first step on path to select one or more elements.
			var stop bool
			expr, path, stop = analyser.processLiteralSource(pkgi, exprType, x, path)
			if stop {
				return
			}

		case *ast.FuncLit:
			// Current expression is a function literal:
			// schedule parameter or result field as additional target.
			if !path.At(0).IsFunc() {
				// This should not happen...
				panic("unexpected function literal")
			}

			analyser.processSignature(pkgi, x.Pos(), exprType.(*types.Signature), path)
			return

		case *ast.CallExpr:
			// Current expression is a function call in single-value context.
			// Delegate to [Analyser.processCallSource].
			var stop bool
			expr, path, stop = analyser.processCallSource(pkgi, exprType, x, path)
			if stop {
				return
			}

		default:
			break mainloop
		}
	}

	pterm.Warning.Printfln(
		"%s: unsupported expression has been marked as a binding source",
		pkg.Fset.Position(expr.Pos()),
	)
}
