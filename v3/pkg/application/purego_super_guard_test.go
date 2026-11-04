package application

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"
)

// TestPuregoSuperCallsUseSendSuperFrom guards the super-call rule of the purego
// backend. purego's objc.ID.SendSuper and objc.SendSuper[T] start the method
// lookup from the receiver's dynamic class, which recurses forever when the
// receiver is an instance of a subclass of the definer or was isa-swizzled by
// AppKit. Super calls must go through sendSuperFrom or sendSuperFromBool, which
// start from the definer (see darwin_purego_cocoa.go).
//
// The test reads the source of every non-test purego file (*purego*.go, which
// includes darwin_purego_cocoa.go, where the helpers live), so it runs on every
// OS and needs no build tag.
func TestPuregoSuperCallsUseSendSuperFrom(t *testing.T) {
	all, err := filepath.Glob("*purego*.go")
	if err != nil {
		t.Fatal(err)
	}
	var files []string
	for _, name := range all {
		if !strings.HasSuffix(name, "_test.go") {
			files = append(files, name)
		}
	}
	if len(files) == 0 {
		t.Fatal("no purego source files found: the guard must run in the package directory")
	}
	fset := token.NewFileSet()
	var bad []string
	for _, name := range files {
		f, err := parser.ParseFile(fset, name, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		ast.Inspect(f, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			// objc.ID(x).SendSuper, objc.SendSuper[T] (an IndexExpr over the
			// selector), and id.sendSuper all reach here as a SelectorExpr.
			if sel.Sel.Name == "SendSuper" || sel.Sel.Name == "sendSuper" {
				bad = append(bad, fset.Position(sel.Pos()).String()+": "+sel.Sel.Name)
			}
			return true
		})
	}
	if len(bad) > 0 {
		t.Errorf("purego super calls must use sendSuperFrom or sendSuperFromBool (purego's SendSuper looks up from the dynamic class):\n  %s",
			strings.Join(bad, "\n  "))
	}
}
