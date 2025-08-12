package generator

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/wailsapp/wails/v3/internal/generator/config"
	"golang.org/x/tools/go/packages"
)

// ResolveSystemPaths resolves paths for the Wails system
func ResolveSystemPaths(buildFlags []string) (paths *config.SystemPaths, err error) {
	// Resolve context pkg path.
	contextPkgPaths, err := ResolvePatterns(buildFlags, "context")
	if err != nil {
		return
	} else if len(contextPkgPaths) < 1 {
		err = ErrNoContextPackage
		return
	} else if len(contextPkgPaths) > 1 {
		// This should never happen...
		panic("context package path matched multiple packages")
	}

	// Resolve wails app pkg path.
	wailsAppPkgPaths, err := ResolvePatterns(buildFlags, config.WailsAppPkgPath)
	if err != nil {
		return
	} else if len(wailsAppPkgPaths) < 1 {
		err = ErrNoApplicationPackage
		return
	} else if len(wailsAppPkgPaths) > 1 {
		// This should never happen...
		panic("wails application package path matched multiple packages")
	}

	paths = &config.SystemPaths{
		ContextPackage:     contextPkgPaths[0],
		ApplicationPackage: wailsAppPkgPaths[0],
	}
	return
}

// ResolvePatterns returns a slice containing all package paths
// that match the given patterns, according to the underlying build tool
// and within the context of the current working directory.
func ResolvePatterns(buildFlags []string, patterns ...string) (paths []string, err error) {
	rewrittenPatterns := make([]string, len(patterns))
	for i, pattern := range patterns {
		rewrittenPatterns[i] = "pattern=" + pattern
	}

	pkgs, err := packages.Load(&packages.Config{
		Mode:       packages.NeedName,
		BuildFlags: buildFlags,
	}, rewrittenPatterns...)

	for _, pkg := range pkgs {
		paths = append(paths, pkg.PkgPath)
	}

	return
}

// LoadPackages loads the packages specified by the given patterns
// and their whole dependency tree. It returns a slice containing
// all packages that match the given patterns and all of their direct
// and indirect dependencies.
//
// The returned slice is in post-order w.r.t. the dependency relation,
// i.e. if package A depends on package B, then package B precedes package A.
//
// All returned package instances include syntax trees and full type information.
//
// Syntax is loaded in the context of a global [token.FileSet],
// which is available through the field [packages.Package.Fset]
// on each returned package. Therefore, source positions
// are canonical across all loaded packages.
func LoadPackages(buildFlags []string, patterns ...string) (pkgs []*packages.Package, err error) {
	rewrittenPatterns := make([]string, len(patterns))
	for i, pattern := range patterns {
		rewrittenPatterns[i] = "pattern=" + pattern
	}

	// Global file set.
	fset := token.NewFileSet()

	roots, err := packages.Load(&packages.Config{
		// NOTE: some Go maintainers now believe deprecation was an error and recommend using Load* modes
		// (see e.g. https://github.com/golang/go/issues/48226#issuecomment-1948792315).
		Mode:       packages.LoadAllSyntax,
		BuildFlags: buildFlags,
		Fset:       fset,
		ParseFile: func(fset *token.FileSet, filename string, src []byte) (file *ast.File, err error) {
			file, err = parser.ParseFile(fset, filename, src, parser.ParseComments|parser.SkipObjectResolution)
			return
		},
	}, rewrittenPatterns...)

	// Flatten dependency tree.
	packages.Visit(roots, nil, func(pkg *packages.Package) {
		if pkg.Fset != fset {
			panic("fileset missing or not the global one")
		}
		pkgs = append(pkgs, pkg)
	})

	return
}
