package parser

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os/exec"
	"path/filepath"
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

	loadMode := packages.NeedName | packages.NeedSyntax | packages.NeedCompiledGoFiles
	if full {
		loadMode |= packages.NeedTypesInfo | packages.NeedTypes
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

func LoadPackagesParallel(buildFlags []string, full bool, patterns ...string) ([]*packages.Package, error) {

	type Result struct {
		pkg *packages.Package
		err error
	}

	jobs := make(chan string, len(patterns))
	results := make(chan Result, len(patterns))

	worker := func() {
		for pattern := range jobs {
			pkg, err := LoadPackage(buildFlags, full, pattern)
			results <- Result{
				pkg: pkg,
				err: err,
			}
		}
	}

	for i := 0; i < 5; i++ {
		go worker()
	}

	for _, pattern := range patterns {
		jobs <- pattern
	}
	close(jobs)

	pkgs := []*packages.Package{}
	for i := 0; i < len(patterns); i++ {
		res := <-results
		if res.err != nil {
			return nil, res.err
		}
		pkgs = append(pkgs, res.pkg)
	}
	close(results)

	return pkgs, nil
}

type ListPackage struct {
	Name    string
	Dir     string
	GoFiles []string
}

// CREDIT: https://cs.opensource.google/go/x/tools/+/refs/tags/v0.20.0:go/packages/golist.go;l=359
func LoadAstPackages(patterns ...string) (map[string]*ast.Package, error) {
	result := make(map[string]*ast.Package)
	if len(patterns) == 0 {
		return result, nil
	}

	// find go files
	listargs := append([]string{"list", "-json=Name,Dir,GoFiles"}, patterns...)
	cmd := exec.Command("go", listargs...)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lPkgs := []*ListPackage{}
	buf := bytes.NewBufferString(string(output))
	for dec := json.NewDecoder(buf); dec.More(); {
		p := new(ListPackage)
		err := dec.Decode(p)
		if err != nil {
			return nil, err
		}
		lPkgs = append(lPkgs, p)
	}

	// load packages
	for i, lPkg := range lPkgs {
		astPkg, err := LoadAstPackage(lPkg)
		if err != nil {
			return result, err
		}
		result[patterns[i]] = astPkg
	}

	return result, nil
}

func LoadAstPackage(pkg *ListPackage) (*ast.Package, error) {
	fset := token.NewFileSet()
	files := make(map[string]*ast.File)
	for _, filename := range pkg.GoFiles {
		goFilePath := filepath.Join(pkg.Dir, filename)
		file, err := parser.ParseFile(fset, goFilePath, nil, parser.AllErrors|parser.ParseComments|parser.SkipObjectResolution)
		if err != nil {
			return nil, err
		}
		files[goFilePath] = file
	}
	return &ast.Package{
		Name:  pkg.Name,
		Files: files,
	}, nil
}
