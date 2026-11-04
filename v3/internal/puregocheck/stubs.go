package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"slices"
)

type decl struct {
	fn    *ast.FuncDecl
	where string
}

// findStubs pairs the functions and methods that only the purego build
// defines with the same-named ones that only the cgo build defines, and
// reports each pair where the purego body does nothing but the cgo body does.
func findStubs(fset *token.FileSet, root string, pure, ref []parsed) []stub {
	paths := func(ps []parsed) []string {
		var out []string
		for _, p := range ps {
			out = append(out, p.path)
		}
		return out
	}
	purePaths, refPaths := paths(pure), paths(ref)
	collect := func(ps []parsed, exclude []string) map[string]decl {
		out := map[string]decl{}
		for _, p := range ps {
			if slices.Contains(exclude, p.path) {
				continue // shared by both builds
			}
			rel, _ := filepath.Rel(root, p.path)
			for _, d := range p.file.Decls {
				fd, ok := d.(*ast.FuncDecl)
				if !ok {
					continue
				}
				out[funcKey(fd)] = decl{fd, fmt.Sprintf("%s:%d", rel, fset.Position(fd.Pos()).Line)}
			}
		}
		return out
	}
	pureDecls := collect(pure, refPaths)
	refDecls := collect(ref, purePaths)
	var out []stub
	for key, p := range pureDecls {
		r, ok := refDecls[key]
		if ok && trivial(p.fn.Body) && !trivial(r.fn.Body) {
			out = append(out, stub{Name: key, Purego: p.where, Cgo: r.where})
		}
	}
	slices.SortFunc(out, func(a, b stub) int {
		if a.Name < b.Name {
			return -1
		}
		if a.Name > b.Name {
			return 1
		}
		return 0
	})
	return out
}

func funcKey(fd *ast.FuncDecl) string {
	if fd.Recv == nil || len(fd.Recv.List) == 0 {
		return fd.Name.Name
	}
	t := fd.Recv.List[0].Type
	if s, ok := t.(*ast.StarExpr); ok {
		t = s.X
	}
	if i, ok := t.(*ast.IndexExpr); ok {
		t = i.X
	}
	if id, ok := t.(*ast.Ident); ok {
		return id.Name + "." + fd.Name.Name
	}
	return fd.Name.Name
}

// trivial reports whether a body does nothing observable: empty, a bare or
// constant return, or only a call to panic/notYet.
func trivial(body *ast.BlockStmt) bool {
	if body == nil {
		return true
	}
	for _, s := range body.List {
		switch st := s.(type) {
		case *ast.ReturnStmt:
			for _, r := range st.Results {
				if !constant(r) {
					return false
				}
			}
		case *ast.ExprStmt:
			call, ok := st.X.(*ast.CallExpr)
			if !ok {
				return false
			}
			if n := funName(call.Fun); n != "panic" && n != "notYet" && n != "notImplemented" {
				return false
			}
		case *ast.EmptyStmt:
		default:
			return false
		}
	}
	return true
}

func constant(e ast.Expr) bool {
	switch x := e.(type) {
	case *ast.BasicLit:
		return true
	case *ast.Ident:
		return x.Name == "nil" || x.Name == "true" || x.Name == "false"
	case *ast.CompositeLit:
		return len(x.Elts) == 0
	case *ast.UnaryExpr:
		return constant(x.X)
	}
	return false
}
