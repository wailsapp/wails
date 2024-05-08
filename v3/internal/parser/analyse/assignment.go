package analyse

import (
	"go/ast"
	"go/token"
	"go/types"
	"slices"
)

// processSendStmt handles a target reference
// that occurs in the context of a channel send statement.
//
// If the occurrence is in channel position and the path begins
// with a send step, the value expression is fed to [Analyser.processExpression].
//
// If the occurrence is in value position and the path
// has reference semantics, the channel expression is fed
// to [Analyser.processExpression] to be scheduled as a target.
func (analyser *Analyser) processSendStmt(pkgi int, send *ast.SendStmt, ref ast.Expr, path Path) {
	pkg := analyser.pkgs[pkgi]

	ch, ok := pkg.TypesInfo.TypeOf(send.Chan).(*types.Chan)
	v := pkg.TypesInfo.TypeOf(send.Value)

	if ok && ch.Dir() != types.RecvOnly && IsAssignableTo(v, ch.Elem()) {
		if ref == send.Chan && path.At(0).IsChanSend() {
			// Curent target is the channel and we are tracking
			// sends into it: process value.
			analyser.processExpression(pkgi, nil, send.Value, path.Consume(1))
		} else if ref == send.Value && path.HasRef() {
			// We are sending a reference to the current target:
			// track receives from the channel.

			path = path.Weaken()

			// If a concrete type is being converted to an interface type,
			// prepend a type assertion to the path.
			if types.IsInterface(ch.Elem()) && !types.IsInterface(v) {
				path = path.Prepend(TypeAssertionStep, analyser.paths.TypeStep(v))
			}

			analyser.processExpression(pkgi, nil, send.Chan, path.Prepend(ChanRecvStep))
		}
	}
}

// processRangeClause handles a target reference
// that occurs in the context of a range clause.
//
// If the occurrence is in value position and the path is not weak,
// an indexing step is prepended to the path
// and the range expression is fed to [Analyser.processExpression].
//
// If the occurrence is in range position, it is indexable
// and its elements have reference semantics,
// the value expression is fed to [Analyser.processExpression]
// to be scheduled as a target.
func (analyser *Analyser) processRangeClause(pkgi int, rng *ast.RangeStmt, ref ast.Expr, path Path) {
	pkg := analyser.pkgs[pkgi]

	// Retrieve element type.
	rtype := pkg.TypesInfo.TypeOf(rng.X)
	if rtype == nil {
		// Invalid range expression: stop here.
		return
	}

	isPtr, isArray := false, false

	switch typ := rtype.Underlying().(type) {
	case *types.Pointer:
		array, ok := typ.Elem().Underlying().(*types.Array)
		if !ok {
			// Not a pointer to array: stop here.
			return
		}

		rtype = array.Elem()
		isPtr, isArray = true, true

	case *types.Array:
		rtype = typ.Elem()
		isArray = true

	case *types.Slice:
		rtype = typ.Elem()

	case *types.Map:
		rtype = typ.Elem()

	default:
		// Invalid range expression: stop here.
		return
	}

	if ref == rng.Value && !path.HasWeakRef() {
		// Current expr is the value and we are tracking assignments to it.
		// Inject an indexing step and analyse range expression.

		// Validate value expression and element assignment.
		ltype := validateLhs(pkg.TypesInfo, rng.Value)
		if ltype == nil || !IsAssignableTo(rtype, ltype) {
			return
		}

		if isPtr {
			path = path.Prepend(IndirectionStep, ArrayIndexingStep)
		} else if isArray {
			path = path.Prepend(ArrayIndexingStep)
		} else {
			path = path.Prepend(IndexingStep)
		}

		analyser.processExpression(pkgi, nil, rng.X, path)

	} else if ref == rng.X {
		// Current expr is range expression.
		// This is equivalent to an implicit indexing step
		// followed by assignment.

		if rng.Value == nil {
			// Loop ranges only over keys: ignore silently.
			return
		}

		if isPtr {
			if path.At(0).IsIndirection() && path.At(1) == ArrayIndexingStep {
				// Allow one automatic pointer indirection if slicing a pointer to array.
				path = path.Consume(1)
			} else {
				// Context does not match path: stop here.
				return
			}
		}

		if !path.At(0).IsIndexing() || !path.Consume(1).HasRef() {
			// Context does not match path
			// or elements do not have reference semantics: stop here.
			return
		}

		// Consume indexing step and weaken element path.
		path = path.Consume(1).Weaken()

		// Validate value expression and element assignment.
		ltype := validateLhs(pkg.TypesInfo, rng.Value)
		if ltype == nil || !IsAssignableTo(rtype, ltype) {
			return
		}

		analyser.processExpression(pkgi, nil, rng.Value, path)
	}
}

// processAssignment handles a target reference
// that occurs in the context of an assignment.
//
// If the occurrence is on the lhs and the path is not weak,
// the corresponding expression on the rhs
// is fed to [Analyser.processExpression].
//
// If the occurrence is on the rhs and the path has reference semantics,
// the corresponding variable or expression on the lhs is fed to
// to [Analyser.processExpression] to be scheduled as a target.
func (analyser *Analyser) processAssignment(pkgi int, assignment *ast.AssignStmt, ref ast.Expr, path Path) {
	pkg := analyser.pkgs[pkgi]

	if assignment.Tok != token.ASSIGN && assignment.Tok != token.DEFINE {
		// Ignore arithmetic assignments.
		return
	}

	if index := slices.Index(assignment.Lhs, ref); index >= 0 {
		if path.HasWeakRef() {
			// Path has weak steps: do not track assignments to target.
			return
		}

		// Reference occurs on the lhs and path is not weak:
		// validate lhs and analyse rhs.
		ltype := validateLhs(pkg.TypesInfo, ref)
		if ltype == nil {
			return
		}

		analyser.processTuple(pkgi, assignment.Rhs, ltype, index, len(assignment.Lhs), path)

	} else if index := slices.Index(assignment.Rhs, ref); index >= 0 {
		if !path.HasRef() || len(assignment.Lhs) < len(assignment.Rhs) {
			// Path does not have reference semantics
			// or assignment is malformed: stop here.
			return
		}

		// Reference occurs on the rhs and path has reference semantics:
		// validate assignment and analyse lhs.
		lhs := assignment.Lhs[index]
		if !validateAssignment(pkg.TypesInfo, lhs, ref) {
			return
		}

		path = path.Weaken()

		// If a concrete type is being converted to an interface type,
		// prepend a type assertion to the path.
		rtype := pkg.TypesInfo.TypeOf(ref)
		if types.IsInterface(pkg.TypesInfo.TypeOf(lhs)) && !types.IsInterface(rtype) {
			path = path.Prepend(TypeAssertionStep, analyser.paths.TypeStep(rtype))
		}

		analyser.processExpression(pkgi, nil, lhs, path)
	}
}

// processVarDecl handles a target reference
// that occurs in the context of a variable declaration.
//
// If the occurrence is on the lhs and the path is not weak,
// the corresponding expression on the rhs is fed to [Analyser.processExpression].
//
// If the occurrence is on the rhs and the path has reference semantics,
// the corresponding variable on the lhs is fed to [Analyser.processExpression]
// to be scheduled as a target.
func (analyser *Analyser) processVarDecl(pkgi int, decl *ast.ValueSpec, ref ast.Expr, path Path) {
	pkg := analyser.pkgs[pkgi]

	// Try to cast ref as identifier, in case it occurs on the lhs.
	ident, _ := ref.(*ast.Ident)

	if index := slices.Index(decl.Names, ident); ident != nil && index >= 0 {
		if path.HasWeakRef() {
			// Path has weak steps: do not track assignments to target.
			return
		}

		// Reference occurs on the lhs and path is not weak:
		// validate lhs and analyse rhs.
		ltype := validateLhs(pkg.TypesInfo, ident)
		if ltype == nil {
			return
		}

		analyser.processTuple(pkgi, decl.Values, ltype, index, len(decl.Names), path)

	} else if index := slices.Index(decl.Values, ref); index >= 0 {
		if !path.HasRef() || len(decl.Names) != len(decl.Values) {
			// Path does not have reference semantics
			// or declaration is malformed: stop here.
			return
		}

		// Reference occurs on the rhs and path has reference semantics:
		// validate assignment and analyse lhs.
		lhs := decl.Names[index]
		if !validateAssignment(pkg.TypesInfo, lhs, ref) {
			return
		}

		path = path.Weaken()

		// If a concrete type is being converted to an interface type,
		// prepend a type assertion to the path.
		rtype := pkg.TypesInfo.TypeOf(ref)
		if types.IsInterface(pkg.TypesInfo.TypeOf(lhs)) && !types.IsInterface(rtype) {
			path = path.Prepend(TypeAssertionStep, analyser.paths.TypeStep(rtype))
		}

		analyser.processExpression(pkgi, nil, lhs, path)
	}
}

// validateLhs checks whether the expression in lhs is valid and assignable to.
// If yes, it returns the type of the expression, otherwise it returns nil.
func validateLhs(info *types.Info, lhs ast.Expr) types.Type {
	if ident, ok := lhs.(*ast.Ident); ok {
		variable, ok := info.ObjectOf(ident).(*types.Var)
		if ok && IsValidType(variable.Type()) {
			return variable.Type()
		}
	} else {
		tv, ok := info.Types[lhs]
		if ok && tv.Assignable() && IsValidType(tv.Type) {
			return tv.Type
		}
	}

	return nil
}

// validateAssignment returns true if and only if expressions lhs and rhs
// are well-formed, well-typed, and rhs can be assigned to lhs.
func validateAssignment(info *types.Info, lhs ast.Expr, rhs ast.Expr) bool {
	var ltype types.Type

	if ident, ok := lhs.(*ast.Ident); ok {
		variable, ok := info.ObjectOf(ident).(*types.Var)
		if !ok {
			return false
		}
		ltype = variable.Type()
	} else {
		tv, ok := info.Types[lhs]
		if !ok || !tv.Assignable() {
			return false
		}
		ltype = tv.Type
	}

	rtype := info.TypeOf(rhs)
	if tv, ok := info.Types[rhs]; ok && tv.HasOk() {
		rtype = tv.Type.(*types.Tuple).At(0).Type()
	}

	return IsAssignableTo(rtype, ltype)
}
