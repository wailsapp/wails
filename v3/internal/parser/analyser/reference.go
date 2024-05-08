package analyser

import (
	"go/ast"
	"go/token"
	"go/types"
	"slices"
)

// processReference analyses the context of an AST node
// that references a target variable or field.
// It locates either assignments to the target variable
// (or one of its subfields), which are then processed as binding sources,
// or variables that may contain a reference to the current
// target (or one of its subfields), which are scheduled as additional targets.
//
// path represents a sequence of access operations (see [Path]) that,
// if performed upon the target variable or field,
// should yield an assignable reference to a binding value.
// processReference may alter the contents of the path slice,
// hence the caller should pass in a fresh copy or cede ownership.
//
// processReference keeps walking _up_ the AST
// as long as the context expression
// matches the expected sequence of operations,
// consuming them as they are performed by the syntax.
//
// If processReference encounters a non-matching operation,
// it returns immediately.
//
// If the expression is used by an address operator or to initialise
// a member of a composite literal, processReference prepends an indirection,
// field selection or indexing operation to the path and keeps walking.
//
// If processReference encounters a statement or function call
// that uses the expression under analysis, it forwards it
// to the appropriate method for further analysis and stops.
//
// When node is an [ast.Field] that defines the target variable
// as a function parameter or result, processReference feeds it
// to [Analyser.processParamOrResultDefinition],
// that may warn about some unsupported cases
// and if necessary analyse all return statements in the function body.
//
// When node is an [ast.CaseClause] that defines the target variable
// as a type switch variable, processReference prepends a type assertion
// to the path as specified by the case clause, then retrieves
// the type switch assertion, feeds the expression to [Analyser.processExpression]
// and returns.
func (analyser *Analyser) processReference(pkgi int, variable *types.Var, node ast.Node, path Path) {
	pkg := analyser.pkgs[pkgi]
	context := FindAstPath(pkg, node.Pos(), node.End())

	// Handle implicit definitions
	switch def := context[0].(type) {
	case *ast.Field:
		// Target is an anonymous parameter or result.
		// Rewind context and let the main loop below handle this.
		context = slices.Insert(context, 0, nil)

	case *ast.CaseClause:
		// Each case clause in a type switch defines
		// a different instance of the switch variable.
		// Retrieve the switch guard and analyse it
		// as an assignment to the target variable.
		//
		// We treat this as a special case because the main loop below
		// must ignore case clauses from selects and regular switches.

		if path.HasWeakRef() {
			// Path has weak steps, don't track assignments to the target.
			return
		}

		// NOTE: the switch guard may be an ast.ExprStmt,
		// but in that case no type switch variable would be defined,
		// hence it must be a single-value ast.AssignStmt.
		//
		// AST shape is
		//   context[2]      context[1]    context[0]
		// TypeSwitchStmt -> BlockStmt  -> CaseClause.
		//   |
		//   +-> Assign: AssignStmt -> Rhs[0]: TypeAssertExpr -> X: Expr

		// Retrieve source expression from switch guard.
		src := context[2].(*ast.TypeSwitchStmt).Assign.(*ast.AssignStmt).Rhs[0].(*ast.TypeAssertExpr).X
		styp := pkg.TypesInfo.TypeOf(src)
		if !IsValidType(styp) || !types.IsInterface(styp) {
			// Invalid source expression: stop here.
			return
		}

		if len(def.List) != 1 {
			// Default or multi-type clause, forward path as is.
			analyser.processExpression(pkgi, nil, src, path)
		} else {
			// Single type clause: if concrete, prepend type assertion to path.
			if ctyp := pkg.TypesInfo.TypeOf(def.List[0]); !types.IsInterface(ctyp) {
				path = path.Prepend(TypeAssertionStep, analyser.paths.TypeStep(ctyp))
			}
			analyser.processExpression(pkgi, nil, src, path)
		}

		return
	}

	// Walk up the AST.
	for context = Reparen(context); len(context) > 1; context = Reparen(context[1:]) {
		switch parent := context[1].(type) {
		case *ast.Field:
			// Target is either a function parameter, a function result field,
			// or a struct field (the root target application.Options.Bind).

			// AST shape is
			//          context[3]             context[2]    context[1]    context[0]
			// StructType/FuncType/FuncDecl -> FieldList  ->   Field    -> Ident/nil
			//
			// Struct field definitions are ignored silently.

			if _, ok := context[3].(*ast.StructType); !ok {
				analyser.processParamOrResultDefinition(pkgi, variable, context, path)
			}

			return

		case *ast.SelectorExpr:
			var stop bool
			path, stop = analyser.processSelectorSink(pkgi, parent, context[0].(ast.Expr), path)
			if stop {
				return
			}

		case *ast.IndexExpr:
			if context[0] != parent.X {
				// Current expr is not the subject expression: stop here.
				return
			}

			if path.At(0).IsIndirection() && path.At(1) == ArrayIndexingStep {
				// Allow one automatic pointer indirection if indexing a pointer to array.
				path = path.Consume(1)
			}

			if !path.At(0).IsIndexing() {
				// Context does not match path: stop here.
				return
			}

			// Consume indexing step and continue.
			path = path.Consume(1)

		case *ast.SliceExpr:
			if context[0] != parent.X {
				// Current expr is not the subject expression: stop here.
				return
			}

			if path.At(0).IsIndirection() && path.At(1) == ArrayIndexingStep {
				// Allow one automatic pointer indirection if slicing a pointer to array.
				path = path.Consume(1)
			}

			if !path.At(0).IsIndexing() {
				// Context does not match path: stop here.
				return
			}

			// After slicing we still need to perform an indexing step:
			// continue without removing it from the path,
			// but change array/strong indexing to weak indexing.
			path.Set(0, WeakIndexingStep)

		case *ast.StarExpr:
			if !path.At(0).IsIndirection() {
				// Context does not match path: stop here.
				return
			}

			// Consume indirection step and continue.
			path = path.Consume(1)

		case *ast.UnaryExpr:
			switch parent.Op {
			case token.AND:
				// Target address taken: prepend indirection step and continue.
				path = path.Prepend(WeakIndirectionStep)

			case token.ARROW:
				if !path.At(0).IsChanRecv() {
					// Context does not match path: stop here.
					return
				}

				if ch, ok := pkg.TypesInfo.TypeOf(parent.X).(*types.Chan); !ok || ch.Dir() == types.SendOnly {
					// Invalid receive expression: stop here.
					return
				}

				// Consume receive step and continue.
				path = path.Consume(1)

			default:
				// Invalid or irrelevant operation: stop here.
				return
			}

		case *ast.TypeAssertExpr:
			if context[0] != parent.X {
				// Current expr is the type: stop here.
				return
			}

			if parent.Type == nil {
				// Handle type switch.
				if _, ok := context[2].(*ast.AssignStmt); !ok {
					// Switch without assignment: stop here.
					return
				}

				// Range over case clauses and extract type assertions.
				sw := context[3].(*ast.TypeSwitchStmt)
				for _, stmt := range sw.Body.List {
					clause := stmt.(*ast.CaseClause)
					clauseVar, ok := pkg.TypesInfo.Implicits[clause].(*types.Var)
					if !ok {
						// Skip invalid clause.
						continue
					}

					clausePath := path

					if len(clause.List) == 1 {
						typ := pkg.TypesInfo.TypeOf(clause.List[0])
						if !IsValidType(typ) {
							// Skip invalid clause.
							continue
						}

						if !types.IsInterface(typ) {
							// Single concrete type clause: take type assertion off path.
							if path.At(0) != TypeAssertionStep || path.At(1) != analyser.paths.TypeStep(typ) {
								// Context does not match path: skip clause.
								continue
							}
							clausePath = clausePath.Consume(2)
						}
					}

					analyser.schedule(target{
						pkgi:     pkgi,
						variable: clauseVar,
						path:     analyser.paths.Get(clausePath),
					})
				}

				return
			}

			typ := pkg.TypesInfo.TypeOf(parent.Type)
			if !IsValidType(typ) {
				// Invalid type: stop here.
				return
			}

			if path.At(0) != TypeAssertionStep || path.At(1) != analyser.paths.TypeStep(typ) {
				// Context does not match path: stop here.
				return
			}

			// Consume type assertion and continue.
			path = path.Consume(2)

		case *ast.KeyValueExpr:
			if context[0] != parent.Value {
				// Current expr is the key: stop here.
				return
			}

			// Continue and analyse composite literal.

		case *ast.CompositeLit:
			var stop bool
			path, stop = analyser.processLiteralSink(pkgi, parent, context[0].(ast.Expr), path)
			if stop {
				return
			}

		case *ast.CallExpr:
			var stop bool
			path, stop = analyser.processCallSink(pkgi, parent, context[0].(ast.Expr), path)
			if stop {
				return
			}

		case *ast.SendStmt:
			analyser.processSendStmt(pkgi, parent, context[0].(ast.Expr), path)
			return

		case *ast.RangeStmt:
			analyser.processRangeClause(pkgi, parent, context[0].(ast.Expr), path)
			return

		case *ast.AssignStmt:
			analyser.processAssignment(pkgi, parent, context[0].(ast.Expr), path)
			return

		case *ast.ValueSpec:
			analyser.processVarDecl(pkgi, parent, context[0].(ast.Expr), path)
			return

		case *ast.ReturnStmt:
			analyser.processReturnSink(pkgi, parent, context, path)
			return

		default:
			// Stop as soon as we encounter an irrelevant context node.
			return
		}
	}
}
