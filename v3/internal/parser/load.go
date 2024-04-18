package parser

import (
	"bytes"
	"cmp"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os/exec"
	"path/filepath"
	"slices"
	"sync"

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

type listPackage struct {
	Name    string
	Dir     string
	GoFiles []string
}

type loaderPackage struct {
	listPkg *listPackage
	astPkg  *ast.Package
	err     error
}

// CREDIT: https://cs.opensource.google/go/x/tools/+/refs/tags/v0.20.0:go/packages/golist.go;l=359
func LoadAstPackages(patterns ...string) (map[string]*ast.Package, error) {
	result := make(map[string]*ast.Package)
	if len(patterns) == 0 {
		return result, nil
	}

	// find go files
	listargs := append([]string{"list", "-find", "-json=Name,Dir,GoFiles"}, patterns...)
	cmd := exec.Command("go", listargs...)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lPkgs := []*loaderPackage{}
	buf := bytes.NewBufferString(string(output))
	for dec := json.NewDecoder(buf); dec.More(); {
		p := new(listPackage)
		err := dec.Decode(p)
		if err != nil {
			return nil, err
		}
		lPkgs = append(lPkgs, &loaderPackage{
			listPkg: p,
		})
	}

	// load packages concurrently
	var wg sync.WaitGroup
	for _, lPkg := range lPkgs {
		wg.Add(1)
		go func(lPkg *loaderPackage) {
			lPkg.astPkg, lPkg.err = LoadAstPackage(lPkg.listPkg)
			wg.Done()
		}(lPkg)
	}
	wg.Wait()

	// load packages
	for i, lPkg := range lPkgs {
		if lPkg.err != nil {
			return result, lPkg.err
		}
		result[patterns[i]] = lPkg.astPkg
	}

	return result, nil
}

func LoadAstPackage(pkg *listPackage) (*ast.Package, error) {
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
