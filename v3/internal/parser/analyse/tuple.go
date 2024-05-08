package analyse

import (
	"go/ast"
	"go/types"
)

// processTuple selects an expression to process
// corresponding to the given index in a tuple of expressions.
//
// The length parameter specifies the expected number of values
// that tuple should provide. This requirement may be satisfied
// either by the presence of many expressions, with len(tuple) == length,
// or by providing a unique expression (comma-ok or function call)
// whose type is a tuple of appropriate length.
//
// The purpose of this step is to find the correct target for analysis
// on the rhs of a multi-value assignment or in the result tuple
// of a return statement. When the tuple contains multiple entries,
// we always pass the index-th entry to [Analyser.processExpression].
// However, when the tuple has just one entry which is a function call,
// the target corresponding to the given index must be selected among
// the result fields of the function.
func (analyser *Analyser) processTuple(pkgi int, tuple []ast.Expr, targetType types.Type, index int, length int, path Path) {
	pkg := analyser.pkgs[pkgi]

	// Assignments and returns in Go may take three forms:
	//   - single-valued: then length == 1 and we can analyse
	//     the unique element in the tuple without further processing;
	//   - multi-valued, with multiple values provided explicitly:
	//     then len(tuple) > 1 and we just have to select the right element;
	//   - multi-valued, with multiple values provided implicitly
	//     by a single expression: in this case we need to inspect
	//     the type of the unique element in the tuple.

	if len(tuple) != 1 && len(tuple) != length {
		// Mismatch between actual and expected tuple length
		// is a type-checking error and we can ignore it silently.
		return
	}

	// Handle the first two cases.
	if length == 1 || len(tuple) > 1 {
		analyser.processExpression(pkgi, targetType, tuple[index], path)
		return
	}

	// Handle the third case (multiple values provided by a single expression).
	tv, ok := pkg.TypesInfo.Types[tuple[0]]
	if !ok {
		// Invalid expression, ignore it silently
		return
	}

	typeTuple, ok := tv.Type.(*types.Tuple)
	if !ok {
		// We know already that length != 1
		// hence we must reject single expressions
		// whose type is not a tuple.
		return
	}

	// tuple[0] must be a either a call expression or a comma-ok form
	// (i.e. a map indexing expression, channel receive or type assertion).

	if typeTuple.Len() != length {
		// Mismatch between actual and expected tuple length
		// is a type-checking error and we can ignore it silently.
		return
	}

	typ := typeTuple.At(index)

	if len(path) == 0 && !types.IsInterface(typ.Type()) {
		// Path is empty and we reached a concrete type:
		// report it without further processing.
		analyser.reportResult(pkgi, tuple[0].Pos(), typ.Type())
		return
	}

	// Handle comma-ok forms
	if tv.HasOk() {
		if index == 0 {
			// For the arbitrarily typed value,
			// pass on the comma-ok expression unchanged.
			// processExpression knows how to handle it.
			analyser.processExpression(pkgi, targetType, tuple[0], path)
		} else if index == 1 && len(path) == 0 {
			// Boolean values only matter when path is empty,
			// as booleans have no subfields and would be ignored anyways.
			analyser.reportResult(pkgi, tuple[0].Pos(), types.Typ[types.Bool])
		}

		return
	}

	// Analyse function expression.
	// Note that we don't have to check for builtins or slice functions here
	// as the ones we handle specially never return multiple values.
	analyser.processExpression(pkgi, nil, tuple[0].(*ast.CallExpr).Fun, path.Prepend(ResultStep, Step(index)))
}
