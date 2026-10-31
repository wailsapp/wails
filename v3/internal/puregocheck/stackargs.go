package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"
)

// The stacknarrow rule (see abi.go for the ABI reasoning).
//
// Every call that purego places on the darwin/arm64 stack is classified with
// go/types: its arguments get the AAPCS64 register allocation, and a narrow
// argument (bool, int8..32, uint8..32, float32, or a composite with such a
// member or a size that is not a multiple of 8) that spills to the stack is a
// finding. Two kinds of call are checked:
//
//   - Objective-C sends. The receiver and the selector take x0 and x1; the
//     Go arguments follow. These go through purego's variadic path, which
//     does not pack the stack, so every spill is a finding.
//   - Fixed-signature C functions registered with RegisterLibFunc (any call
//     whose first argument is &fn of a func type). purego v0.11.1 bundles
//     their stacked arguments with the Apple layout, so a narrow spill there
//     is counted (bundled) but not reported.
//
// Calls whose argument types cannot be resolved are counted (unknown) and not
// guessed at. The rule runs for darwin only.

// stackFinding is one narrow argument purego places on the stack.
type stackFinding struct {
	Name  string `json:"name"`  // C symbol or selector
	Where string `json:"where"` // file:line of the call
	Text  string `json:"text"`  // the report line
}

// stackResult is what checkStackPackage found in one package.
type stackResult struct {
	findings []stackFinding
	checked  int      // calls classified
	bundled  int      // fixed-signature narrow spills purego packs correctly
	unknown  int      // calls skipped: an argument or parameter type is not mapped
	skipped  []string // file:line of each unknown call
	// forwarded counts the calls inside the send helpers that pass the
	// helper's own args... on. Their arguments are classified at each call of
	// the helper instead.
	forwarded int
}

// forwards reports whether call is a send helper (one sendArgs recognises)
// passing its own variadic parameter on unchanged.
func forwards(info *types.Info, fd *ast.FuncDecl, call *ast.CallExpr) bool {
	if fd == nil || !call.Ellipsis.IsValid() || len(call.Args) == 0 {
		return false
	}
	name := fd.Name.Name
	helper := name == "get" || name == "sendSuperFrom" || name == "sendSuperFromBool" || (fd.Recv != nil && isSendName(name))
	params := fd.Type.Params.List
	if !helper || len(params) == 0 {
		return false
	}
	last := params[len(params)-1]
	if _, variadic := last.Type.(*ast.Ellipsis); !variadic || len(last.Names) == 0 {
		return false
	}
	arg, ok := call.Args[len(call.Args)-1].(*ast.Ident)
	return ok && info.Uses[arg] != nil && info.Uses[arg] == info.Defs[last.Names[len(last.Names)-1]]
}

// skip counts a call that cannot be classified and records where it is.
func (res *stackResult) skip(where string) {
	res.unknown++
	res.skipped = append(res.skipped, where)
}

// stackImporter resolves the imports of the darwin purego build (objc.ID is
// uintptr, which the argument sizes depend on). One importer serves a whole
// scan so that each dependency is type-checked once. The source importer keeps
// a pointer to build.Default, so the darwin/arm64 context is left in place for
// the importer's lifetime; buildFiles copies build.Default and sets its own
// fields, so it is not affected.
func stackImporter(fset *token.FileSet, tags []string) types.Importer {
	build.Default.GOOS = "darwin"
	build.Default.GOARCH = "arm64"
	build.Default.CgoEnabled = false
	build.Default.BuildTags = append([]string{}, tags...)
	return importer.ForCompiler(fset, "source", nil)
}

// checkStackPackage type-checks one package (its purego-selected files) and
// classifies every call that places arguments on the stack.
// checkStackPackage classifies the calls of one package. root is the v3 module
// root, so that reported positions are relative to it.
func checkStackPackage(fset *token.FileSet, imp types.Importer, root string, files []*ast.File) stackResult {
	var res stackResult
	info := &types.Info{Types: map[ast.Expr]types.TypeAndValue{}, Uses: map[*ast.Ident]types.Object{}, Defs: map[*ast.Ident]types.Object{}}
	conf := types.Config{Importer: imp, Sizes: types.SizesFor("gc", "arm64"), Error: func(error) {}}
	conf.Check("purego", fset, files, info)

	// The C functions: every object whose address is passed to a
	// registration call (RegisterLibFunc, RegisterFunc, the register helpers).
	registered := map[types.Object]bool{}
	for _, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || len(call.Args) < 2 {
				return true
			}
			u, ok := call.Args[0].(*ast.UnaryExpr)
			if !ok || u.Op != token.AND {
				return true
			}
			obj := objectOf(info, u.X)
			if _, isFn := types.Unalias(info.TypeOf(u.X)).Underlying().(*types.Signature); !isFn || obj == nil {
				return true
			}
			registered[obj] = true
			return true
		})
	}

	for _, f := range files {
		var enclosing *ast.FuncDecl
		ast.Inspect(f, func(n ast.Node) bool {
			switch d := n.(type) {
			case *ast.FuncDecl:
				enclosing = d
			case *ast.GenDecl:
				enclosing = nil
			}
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			if forwards(info, enclosing, call) {
				res.forwarded++
				return true
			}
			if obj := objectOf(info, call.Fun); obj != nil {
				if registered[obj] {
					res.fixedCall(info, call, where(fset, root, call.Pos()))
					return true
				}
			}
			res.objcSend(info, call, where(fset, root, call.Pos()))
			return true
		})
	}
	return res
}

// where formats a position as root-relative file:line.
func where(fset *token.FileSet, root string, pos token.Pos) string {
	p := fset.Position(pos)
	rel, err := filepath.Rel(root, p.Filename)
	if err != nil {
		rel = p.Filename
	}
	return fmt.Sprintf("%s:%d", rel, p.Line)
}

// objectOf returns the object an identifier or qualified identifier names.
func objectOf(info *types.Info, e ast.Expr) types.Object {
	switch x := e.(type) {
	case *ast.Ident:
		return info.Uses[x]
	case *ast.SelectorExpr:
		return info.Uses[x.Sel]
	case *ast.ParenExpr:
		return objectOf(info, x.X)
	}
	return nil
}

// fixedCall classifies a call of a registered fixed-signature C function.
func (res *stackResult) fixedCall(info *types.Info, call *ast.CallExpr, where string) {
	sig, ok := types.Unalias(info.TypeOf(call.Fun)).Underlying().(*types.Signature)
	if !ok || sig.Variadic() || call.Ellipsis.IsValid() || len(call.Args) != sig.Params().Len() {
		res.skip(where)
		return
	}
	sizes := types.SizesFor("gc", "arm64")
	args := make([]cArg, 0, sig.Params().Len())
	for i := 0; i < sig.Params().Len(); i++ {
		a, ok := classify(sig.Params().At(i).Type(), sizes)
		if !ok {
			res.skip(where)
			return
		}
		args = append(args, a)
	}
	res.checked++
	for _, i := range allocate(args, 0) {
		if args[i].narrow {
			res.bundled++
		}
	}
}

// objcSend classifies a call that sends an Objective-C message through one of
// the backend's helpers. It is a no-op for every other call.
func (res *stackResult) objcSend(info *types.Info, call *ast.CallExpr, where string) {
	args, sel, ok := sendArgs(info, call)
	if !ok {
		return
	}
	if call.Ellipsis.IsValid() {
		res.skip(where)
		return
	}
	name := "<dynamic>"
	if sel != nil {
		if s, ok := lit(sel); ok {
			name = s
		}
	}
	sizes := types.SizesFor("gc", "arm64")
	cargs := make([]cArg, 0, len(args))
	for _, a := range args {
		t := info.TypeOf(a)
		if t == nil {
			res.skip(where)
			return
		}
		c, ok := classify(types.Default(t), sizes)
		if !ok {
			res.skip(where)
			return
		}
		cargs = append(cargs, c)
	}
	res.checked++
	for _, i := range allocate(cargs, 2) {
		if !cargs[i].narrow {
			continue
		}
		t := types.Default(info.TypeOf(args[i]))
		// Parameter numbers count the C signature: self is 1, _cmd is 2.
		text := fmt.Sprintf("STACKNARROW %s param %d (%s) is passed on the stack; purego misplaces it on darwin/arm64",
			name, i+3, types.TypeString(t, nil))
		res.findings = append(res.findings, stackFinding{
			Name:  name,
			Where: where,
			Text:  text,
		})
	}
}

// sendArgs recognises the backend's message-send helpers and returns the
// Go arguments that become C arguments after self and _cmd, and the selector
// expression. ok is false for calls that are not sends.
//
//	recv.send(sel, args...)           method on id (and Send/SendSuper)
//	objc.Send[T](recv, sel, args...)  purego package function
//	get[T](recv, sel, args...)        the backend's typed send
//	sendSuperFrom(definer, self, sel, args...)  and sendSuperFromBool
func sendArgs(info *types.Info, call *ast.CallExpr) (args []ast.Expr, sel ast.Expr, ok bool) {
	fun := call.Fun
	index := false
	if ix, isIx := fun.(*ast.IndexExpr); isIx {
		fun = ix.X
		index = true
	} else if ix, isIx := fun.(*ast.IndexListExpr); isIx {
		fun = ix.X
		index = true
	}
	switch f := fun.(type) {
	case *ast.Ident:
		switch f.Name {
		case "get":
			if index && len(call.Args) >= 2 {
				return call.Args[2:], call.Args[1], true
			}
		case "sendSuperFrom", "sendSuperFromBool":
			if len(call.Args) >= 3 {
				return call.Args[3:], call.Args[2], true
			}
		}
	case *ast.SelectorExpr:
		if !isSendName(f.Sel.Name) {
			return nil, nil, false
		}
		if pkg, isIdent := f.X.(*ast.Ident); isIdent {
			if _, isPkg := info.Uses[pkg].(*types.PkgName); isPkg {
				// objc.Send[T](recv, sel, args...)
				if index && len(call.Args) >= 2 {
					return call.Args[2:], call.Args[1], true
				}
				return nil, nil, false
			}
		}
		// recv.send(sel, args...): the receiver is self.
		if len(call.Args) >= 1 {
			return call.Args[1:], call.Args[0], true
		}
	}
	return nil, nil, false
}

func isSendName(name string) bool {
	switch name {
	case "send", "Send", "SendSuper", "sendSuper":
		return true
	}
	return false
}

// stackNarrowSummary is the stacknarrow part of the report summary.
func stackNarrowSummary(r *report) string {
	var b strings.Builder
	fmt.Fprintf(&b, "  stack-argument calls checked: %d, fixed-signature narrow spills packed by purego: %d, send-helper forwards: %d, not classified: %d\n",
		r.StackChecked, r.StackBundled, r.StackForwarded, r.StackUnknown)
	return b.String()
}
