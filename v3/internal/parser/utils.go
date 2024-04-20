package parser

import (
	"cmp"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"slices"

	"github.com/pterm/pterm"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

// FindAstPath returns the node that encloses the source interval [start, end),
// and all its ancestors up to the AST root, in the syntax tree for pkg.
//
// If no source file can be found for the specified interval,
// FindAstPath returns an empty slice.
//
// See [astutil.PathEnclosingInterval].
func FindAstPath(pkg *packages.Package, start token.Pos, end token.Pos) []ast.Node {
	// Perform a binary search to find the file enclosing the node
	fileIndex, exact := slices.BinarySearchFunc(pkg.Syntax, start, func(f *ast.File, p token.Pos) int {
		return cmp.Compare(f.FileStart, p)
	})

	// If exact is true, pkg.Syntax[fileIndex] is the file we are looking for;
	// otherwise, it is the first file whose start position is _after_ ident.Pos()
	if !exact {
		fileIndex--
	}

	// When exact is false, the search could theoretically fail (this is bad and should never happen)
	if fileIndex < 0 || start < pkg.Syntax[fileIndex].FileStart || pkg.Syntax[fileIndex].FileEnd < end {
		return nil
	}

	path, _ := astutil.PathEnclosingInterval(pkg.Syntax[fileIndex], start, end)
	return path
}

// Reparen is the opposite of [ast.Unparen]: it travels up the given path
// until the immediate context path[1] is an unparenthesised expression.
func Reparen(path []ast.Node) []ast.Node {
	for ; len(path) > 1; path = path[1:] {
		if _, ok := path[1].(*ast.ParenExpr); !ok {
			break
		}
	}

	return path
}

func aliasToNamed(alias *types.Alias) *types.Named {
	return types.NewNamed(alias.Obj(), alias.Underlying(), nil)
}

// filteredPrefixPrinter is used to remove duplicate output of (pterm.PrefixPrinter)s
type filteredPrefixPrinter struct {
	printer  *pterm.PrefixPrinter
	messages map[string]bool
}

func newFilteredPrefixPrinter(printer *pterm.PrefixPrinter) *filteredPrefixPrinter {
	return &filteredPrefixPrinter{
		printer:  printer,
		messages: make(map[string]bool),
	}
}

func (p *filteredPrefixPrinter) Printfln(format string, a ...interface{}) *pterm.TextPrinter {
	message := fmt.Sprintf(format, a...)
	if _, ok := p.messages[message]; ok {
		return nil
	}
	p.messages[message] = true
	return p.printer.Println(message)
}

var filteredWarning = newFilteredPrefixPrinter(&pterm.Warning)
