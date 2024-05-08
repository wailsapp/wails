package parser

import (
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/packages"
)

// LoadPackages loads the packages specified by the given patterns.
//
// The resulting package instances include syntax trees and types.
// If full is true, they also include detailed type information.
func LoadPackages(buildFlags []string, full bool, patterns ...string) ([]*packages.Package, error) {
	rewrittenPatterns := make([]string, len(patterns))
	for i, pattern := range patterns {
		rewrittenPatterns[i] = "pattern=" + pattern
	}

	// Global file set for syntax-only mode.
	var fset *token.FileSet

	loadMode := packages.NeedName | packages.NeedCompiledGoFiles | packages.NeedSyntax
	if full {
		loadMode |= packages.NeedTypes | packages.NeedTypesInfo
	} else {
		fset = token.NewFileSet()
	}

	pkgs, err := packages.Load(&packages.Config{
		Mode:       loadMode,
		Logf:       nil,
		BuildFlags: buildFlags,
		Fset:       fset,
		ParseFile: func(fset *token.FileSet, filename string, src []byte) (file *ast.File, err error) {
			file, err = parser.ParseFile(fset, filename, src, parser.ParseComments|parser.SkipObjectResolution)
			return
		},
	}, rewrittenPatterns...)

	// If in syntax only mode, add global file set to each package.
	if !full {
		for _, pkg := range pkgs {
			pkg.Fset = fset
		}
	}

	return pkgs, err
}
