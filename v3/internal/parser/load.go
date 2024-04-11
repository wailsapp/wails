package parser

import (
	"cmp"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"slices"

	"github.com/pterm/pterm"
	"golang.org/x/tools/go/packages"
)

// WailsAppPkgPath is the official import path of Wails v3's application package
const WailsAppPkgPath = "github.com/wailsapp/wails/v3/pkg/application"

// LoadPackages loads the packages specified by the given patterns.
//
// The resulting package instances include syntax trees and types.
// If full is true, they also include detailed type information.
func LoadPackages(buildFlags []string, full bool, patterns ...string) ([]*packages.Package, error) {
	rewrittenPatterns := make([]string, len(patterns))
	for i, pattern := range patterns {
		rewrittenPatterns[i] = "pattern=" + pattern
	}

	loadMode := packages.NeedName | packages.NeedSyntax | packages.NeedTypes | packages.NeedCompiledGoFiles
	if full {
		loadMode |= packages.NeedTypesInfo
	}

	pkgs, err := packages.Load(&packages.Config{
		Mode: loadMode,
		Logf: func(format string, args ...interface{}) {
			pterm.Debug.Printf(format+"\n", args...)
		},
		BuildFlags: buildFlags,
		ParseFile: func(fset *token.FileSet, filename string, src []byte) (file *ast.File, err error) {
			file, err = parser.ParseFile(fset, filename, src, parser.AllErrors|parser.ParseComments|parser.SkipObjectResolution)
			return
		},
	}, rewrittenPatterns...)

	// Sort (*ast.File)s by start position, ascending
	for _, pkg := range pkgs {
		slices.SortFunc(pkg.Syntax, func(f1 *ast.File, f2 *ast.File) int {
			return cmp.Compare(f1.FileStart, f2.FileStart)
		})
	}

	return pkgs, err
}

func LoadPackage(buildFlags []string, full bool, pattern string) (*packages.Package, error) {
	pkgs, err := LoadPackages(buildFlags, full, pattern)
	if err != nil {
		return nil, err
	}
	if len(pkgs) <= 0 {
		return nil, errors.New("package not found: " + pattern)
	}
	return pkgs[0], nil
}
