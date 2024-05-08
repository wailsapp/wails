package analyser

import (
	"go/ast"
	"go/token"
	"go/types"
	"slices"

	"github.com/pterm/pterm"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/types/typeutil"
)

// processCallSink handles a call expression where one of the parameters
// references a target variable or field.
//
// If the callee is a builtin, it delegates to [Analyser.processBuiltinSink].
//
// If the expression is a type conversion to or from unsafe.Pointer,
// processCallSink emits a warning. If it is a valid type conversion
// from concrete to interface type, processCallSink prepends
// a type assertion to the path and instructs processReference
// to keep analysing the context.
//
// Otherwise, processCallSink prepends a ParamStep to the path
// and feeds the function expression to [Analyser.processExpression] to have the
// parameter field added as a new target variable.
func (analyser *Analyser) processCallSink(pkgi int, call *ast.CallExpr, param ast.Expr, pathIn Path) (path Path, stop bool) {
	pkg := analyser.pkgs[pkgi]
	path = pathIn
	stop = true // Unless set to false, processReference will stop processing the context.

	paramIndex := slices.Index(call.Args, param)
	if paramIndex < 0 {
		// Current expr is not an argument to the function: stop here.
		return
	}

	if !IsValidType(pkg.TypesInfo.TypeOf(call)) {
		// Invalid call expr: stop here.
		return
	}

	// Retrieve callee object and type info.
	callee := typeutil.Callee(pkg.TypesInfo, call)
	tv := pkg.TypesInfo.Types[call.Fun]

	// Handle builtin and slice function calls.
	if tv.IsBuiltin() {
		return analyser.processBuiltinSink(pkgi, callee, call, paramIndex, path)
	} else if callee != nil && callee.Pkg().Path() == "slices" {
		return analyser.processSliceFuncSink(pkgi, callee, call, paramIndex, path)
	}

	// From now on we only care about paths with reference semantics.
	if !path.HasRef() {
		return
	}

	path = path.Weaken()

	if tv.IsType() {
		// Call expression is a conversion: retrieve source type.
		typ := pkg.TypesInfo.TypeOf(param)

		if types.Identical(tv.Type, types.Typ[types.UnsafePointer]) {
			// Conversion to unsafe pointer: emit a warning and stop here.
			pterm.Warning.Printfln(
				"%s: use of unsafe features to provide bindings is not supported",
				pkg.Fset.Position(call.Pos()),
			)
			return
		} else if types.Identical(typ, types.Typ[types.UnsafePointer]) {
			// Conversion from unsafe pointer: emit a warning and stop here.
			pterm.Warning.Printfln(
				"%s: use of unsafe features to provide bindings is not supported",
				pkg.Fset.Position(param.Pos()),
			)
			return
		}

		// If a concrete type is being converted to an interface type,
		// prepend a type assertion to the path.
		if types.IsInterface(tv.Type) && !types.IsInterface(typ) {
			path = path.Prepend(TypeAssertionStep, analyser.paths.TypeStep(typ))
		}

		stop = false // Let processReference process the surrounding expression.
		return
	}

	signature := tv.Type.(*types.Signature)
	srcType := pkg.TypesInfo.TypeOf(param)

	if signature.Variadic() && paramIndex >= signature.Params().Len()-1 {
		// Handle variadic parameters.
		paramIndex = signature.Params().Len() - 1
		if !call.Ellipsis.IsValid() {
			// Current expression becomes an element
			// of the newly allocated variadic parameter slice:
			// prepend an indexing step to the path.

			// If a concrete type is being converted to an interface type,
			// prepend a type assertion to the path.
			targetType := signature.Params().At(paramIndex).Type().(*types.Slice).Elem()
			if types.IsInterface(targetType) && !types.IsInterface(srcType) {
				path = path.Prepend(TypeAssertionStep, analyser.paths.TypeStep(srcType))
			}

			path = path.Prepend(WeakIndexingStep)
		}
	} else {
		// If a concrete type is being converted to an interface type,
		// prepend a type assertion to the path.
		targetType := signature.Params().At(paramIndex).Type()
		if types.IsInterface(targetType) && !types.IsInterface(srcType) {
			path = path.Prepend(TypeAssertionStep, analyser.paths.TypeStep(srcType))
		}
	}

	// Prepend parameter access step to path and process function expression.
	analyser.processExpression(pkgi, nil, call.Fun, path.Prepend(ParamStep, Step(paramIndex)))
	return
}

// processReturnSink processes a return statement
// where one of the results references a target variable or field.
//
// It schedules the corresponding result field from the function signature
// as an additional target.
//
// This is necessary to let [Analyser.processParamOrResultDefinition]
// discover and warn about cases where a function
// returns a reference to a binding source,
// but we don't know anything about the call site.
func (analyser *Analyser) processReturnSink(pkgi int, ret *ast.ReturnStmt, context []ast.Node, path Path) {
	pkg := analyser.pkgs[pkgi]

	if !path.HasRef() {
		// We are not returning a reference: stop here.
		return
	}

	// Find index of result field
	fieldIndex := slices.Index(ret.Results, context[0].(ast.Expr))

	// Find index of surrounding function in AST
	fnIndex := slices.IndexFunc(context[2:], func(n ast.Node) bool {
		_, isDecl := n.(*ast.FuncDecl)
		_, isLit := n.(*ast.FuncLit)
		return isDecl || isLit
	})

	if fieldIndex < 0 || fnIndex < 0 {
		// This should never happen...
		panic("return field or surrounding function not found")
	}

	var signature *types.Signature

	switch fn := context[2+fnIndex].(type) {
	case *ast.FuncDecl:
		signature, _ = pkg.TypesInfo.TypeOf(fn.Type).(*types.Signature)
	case *ast.FuncLit:
		signature, _ = pkg.TypesInfo.TypeOf(fn.Type).(*types.Signature)
	}

	if signature == nil || fieldIndex >= signature.Results().Len() {
		// Bad function type? Stop here.
		return
	}

	// Enqueue result field.
	analyser.schedule(target{
		pkgi:     pkgi,
		variable: signature.Results().At(fieldIndex),
		path:     analyser.paths.Get(path.Weaken()),
	})
}

// processParamOrResultDefinition handles the definition for a target
// that is either a function parameter or a function result.
// It determines:
//   - whether to emit warnings about the analyser's
//     inability to follow the call graph;
//   - whether to analyse the function's return statements
//     as assignments to the target.
func (analyser *Analyser) processParamOrResultDefinition(pkgi int, variable *types.Var, context []ast.Node, path Path) {
	pkg := analyser.pkgs[pkgi]

	// AST shape is
	// context[4]    context[3]    context[2]    context[1]    context[0]
	//  FuncLit   ->  FuncType  -> FieldList  ->   Field    -> Ident/nil
	//                FuncDecl  -> FieldList  ->   Field    -> Ident/nil

	ident, _ := context[0].(*ast.Ident)
	field := context[1].(*ast.Field)
	fn := context[3]

	if _, ok := fn.(*ast.FuncDecl); !ok {
		fn = context[4]
	}

	// Retrieve function signature.
	var signature *types.Signature
	if decl, ok := fn.(*ast.FuncDecl); ok {
		obj, _ := pkg.TypesInfo.ObjectOf(decl.Name).(*types.Func)
		if obj != nil {
			signature, _ = obj.Type().(*types.Signature)
		}
	} else {
		signature, _ = pkg.TypesInfo.TypeOf(fn.(ast.Expr)).(*types.Signature)
	}
	if signature == nil {
		// Invalid function type expression, ignore silently.
		return
	}

	// Find parameter field index.
	index := -2

	if signature.Recv() != nil && signature.Recv().Origin() == variable {
		index = -1
	} else {
		for i := 0; i < signature.Params().Len(); i++ {
			if signature.Params().At(i).Origin() == variable {
				index = i
				break
			}
		}
	}

	if index >= -1 {
		// Target is a parameter.

		if !path.HasWeakRef() {
			// We are tracking assignments to a function parameter: emit a warning.
			if ident == nil {
				pterm.Warning.Printfln(
					"%s: function parameter #%d is a binding source: values passed at call sites may not be detected",
					pkg.Fset.Position(field.Pos()),
					index,
				)
			} else {
				pterm.Warning.Printfln(
					"%s: function parameter '%s' is a binding source: values passed at call sites may not be detected",
					pkg.Fset.Position(ident.Pos()),
					ident.Name,
				)
			}
		}

		return
	}

	// Target is not a parameter, hence it must be a result field.

	// Find result field index.
	index = -1
	for i := 0; i < signature.Results().Len(); i++ {
		if signature.Results().At(i) == variable {
			index = i
			break
		}
	}

	if index < 0 {
		// This should never happen...
		panic("result field not found in signature")
	}

	if path.HasWeakRef() {
		// Function returns a reference to a binding source: emit a warning.
		if ident == nil {
			pterm.Warning.Printfln(
				"%s: function result #%d references a binding source: uses at call sites may not be detected",
				pkg.Fset.Position(field.Pos()),
				index,
			)
		} else {
			pterm.Warning.Printfln(
				"%s: function result '%s' references a binding source: uses at call sites may not be detected",
				pkg.Fset.Position(ident.Pos()),
				ident.Name,
			)
		}
	} else {
		// We are tracking assignments to the result field: track returns too.
		var body *ast.BlockStmt

		// Retrieve function body.
		switch fnDef := fn.(type) {
		case *ast.FuncLit:
			// signature belongs to a function literal.
			body = fnDef.Body
		case *ast.FuncDecl:
			// signature belongs to a function declaration.
			body = fnDef.Body
		default:
			// This should never happen...
			panic("function body not available for static analysis")
		}

		analyser.processReturnSources(pkgi, body, variable.Type(), index, signature.Results().Len(), path)
	}
}

// processCallSource handles a call expression in single-value context
// where the unique result field has been marked as a binding source.
//
// If the expression is a type conversion to or from unsafe.Pointer,
// processCallSink emits a warning. For every other type conversion,
// it instructs processExpression to analyse the subject expression.
//
// If the callee is a builtin, processCallSource
// delegates to [Analyser.processBuiltinSource].
// If it is a function from the slices package,
// it delegates to [Analyser.processSliceFuncSource].
//
// Otherwise, processCallSink prepends a ResultStep to the path
// and feeds the call expression to [Analyser.processExpression] to have the
// result field added as a new target variable.
func (analyser *Analyser) processCallSource(pkgi int, exprType types.Type, call *ast.CallExpr, pathIn Path) (expr ast.Expr, path Path, stop bool) {
	pkg := analyser.pkgs[pkgi]
	path = pathIn
	stop = true // Unless set to false, processExpression will stop processing subexpressions.

	tv := pkg.TypesInfo.Types[call.Fun]

	if tv.IsType() {
		expr = call.Args[0]

		if types.Identical(exprType, types.Typ[types.UnsafePointer]) {
			// Conversion to unsafe pointer: emit a warning and stop here.
			pterm.Warning.Printfln(
				"%s: use of unsafe features to provide bindings is not supported",
				pkg.Fset.Position(call.Pos()),
			)
		} else if types.Identical(pkg.TypesInfo.TypeOf(expr), types.Typ[types.UnsafePointer]) {
			// Conversion from unsafe pointer: emit a warning and stop here.
			pterm.Warning.Printfln(
				"%s: use of unsafe features to provide bindings is not supported",
				pkg.Fset.Position(expr.Pos()),
			)
		} else {
			// Safe conversion: continue.
			stop = false
		}

		return
	}

	callee := typeutil.Callee(pkg.TypesInfo, call)

	// Handle builtin and slice function calls.
	if tv.IsBuiltin() {
		return analyser.processBuiltinSource(pkgi, callee, call, path)
	} else if callee != nil && callee.Pkg().Path() == "slices" {
		return analyser.processSliceFuncSource(pkgi, callee, call, path)
	}

	// Prepend result lookup to path
	// and analyse function expression.
	path = path.Prepend(ResultStep, 0)
	expr = call.Fun
	stop = false
	return
}

// processReturnSources visits all return statements in a function body
// and feeds the index-th value to the expression analyser.
func (analyser *Analyser) processReturnSources(pkgi int, body *ast.BlockStmt, targetType types.Type, index int, length int, path Path) {
	astutil.Apply(body, func(c *astutil.Cursor) bool {
		switch n := c.Node().(type) {
		case *ast.ReturnStmt:
			if n.Results != nil {
				analyser.processTuple(pkgi, n.Results, targetType, index, length, path)
			}
			return false
		case ast.Stmt:
			// Visit nested block statements.
			// This is the lazy approach.
			// The more accurate one would be
			// to list here all statement types
			// that may have nested blocks.
			return true
		case *ast.FuncLit:
			// Do not visit nested func literals:
			// their return statements must be ignored.
			// This should be handled by the default case as well,
			// but we make it explicit for increased safety.
			return false
		default:
			return false
		}
	}, nil)
}

// processFunc retrieves a function signature and feeds it to [Analyser.processSignature].
func (analyser *Analyser) processFunc(pkgi int, pos token.Pos, fn *types.Func, path Path) {
	if !IsValidType(fn.Type()) {
		return
	}

	if fn.Pkg().Path() == "reflect" {
		// Warn about unsupported use of reflection and stop.
		pterm.Warning.Printfln(
			"%s: use of reflection to provide bindings is not supported",
			analyser.pkgs[pkgi].Fset.Position(pos),
		)
		return
	} else if fn.Pkg().Path() == "unsafe" {
		// Warn about unsupported use of unsafe features and stop.
		pterm.Warning.Printfln(
			"%s: use of unsafe features to provide bindings is not supported",
			analyser.pkgs[pkgi].Fset.Position(pos),
		)
		return
	}

	analyser.processSignature(pkgi, pos, fn.Type().(*types.Signature), path)
}

// processSignature retrieves a parameter or result field
// from a function signature and enqueues it as an additional target.
func (analyser *Analyser) processSignature(pkgi int, pos token.Pos, signature *types.Signature, path Path) {
	pkg := analyser.pkgs[pkgi]

	if !path.At(0).IsFunc() {
		// We are not interested in functions: stop here.
		return
	}

	if signature.TypeParams() != nil || signature.RecvTypeParams() != nil {
		// We cannot always handle generic functions correctly: emit a warning.
		pterm.Warning.Printfln(
			"%s: parameters or results of generic function have been marked as binding sources: this case is not fully supported",
			pkg.Fset.Position(pos),
		)
	}

	step, index := path.At(0), int(path.At(1))
	path = path.Consume(2)

	var field *types.Var

	switch step {
	case ParamStep:
		if index < 0 {
			field = signature.Recv().Origin()
		} else if index < signature.Params().Len() {
			field = signature.Params().At(index).Origin()
		}
	case ResultStep:
		if index < signature.Results().Len() {
			field = signature.Results().At(signature.Results().Len() - 1).Origin()
		}
	}

	if field == nil {
		// This should never happen...
		panic("could not find required function parameter or result")
	}

	tgt := target{
		pkgi:     varPkgi(analyser, pkgi, field),
		variable: field,
		path:     analyser.paths.Get(path),
	}

	if tgt.pkgi < len(analyser.pkgs) {
		analyser.schedule(tgt)
	} else {
		// This should only happen when targeting a function parameter.
		pterm.Warning.Printfln(
			"%s: parameters or results of function have been marked as binding sources, but its declaring package is not under analysis",
			pkg.Fset.Position(pos),
		)
	}
}
