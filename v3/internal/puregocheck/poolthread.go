package main

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strconv"
)

// findPoolThreads is the poolthread rule (darwin): every purego function that
// pushes an NSAutoreleasePool must pin its goroutine with runtime.LockOSThread.
// A pool belongs to the thread that pushed it, and a goroutine can move to
// another thread between two purego calls; draining the pool there corrupts
// both threads' pool stacks (a SIGSEGV in the drain, found on the macOS VM).
func findPoolThreads(fset *token.FileSet, root string, pure []parsed) []string {
	var out []string
	for _, p := range pure {
		for _, d := range p.file.Decls {
			fd, ok := d.(*ast.FuncDecl)
			if !ok || fd.Body == nil {
				continue
			}
			var pool token.Pos
			locked := false
			ast.Inspect(fd.Body, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.BasicLit:
					if s, err := strconv.Unquote(x.Value); err == nil && s == "NSAutoreleasePool" && !pool.IsValid() {
						pool = x.Pos()
					}
				case *ast.SelectorExpr:
					if id, ok := x.X.(*ast.Ident); ok && id.Name == "runtime" && x.Sel.Name == "LockOSThread" {
						locked = true
					}
				}
				return true
			})
			if pool.IsValid() && !locked {
				pos := fset.Position(pool)
				rel, _ := filepath.Rel(root, pos.Filename)
				out = append(out, rel+":"+strconv.Itoa(pos.Line)+" "+fd.Name.Name)
			}
		}
	}
	return out
}
