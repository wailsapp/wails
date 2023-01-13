package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type ImportSpecInfo struct {
	FileName         string
	ImportSpec       *ast.ImportSpec
	File             *ast.File
	Alias            string
	Options          *ast.Expr
	BoundStructNames []*structInfo
}

func (i ImportSpecInfo) Identifier() string {
	if i.Alias != "" {
		return i.Alias
	}
	return strings.Trim(filepath.Base(i.ImportSpec.Path.Value), `"`)
}

func findFilesImportingPackage(dir, pkg string) ([]ImportSpecInfo, error) {
	var importSpecInfo []ImportSpecInfo

	// Recursively search for .go files in the given directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
		if err != nil {
			return err
		}

		// Search for import statement in the AST
		for _, imp := range f.Imports {
			if imp.Path.Value == "\""+pkg+"\"" {
				var alias string
				if imp.Name != nil {
					alias = imp.Name.Name
				}
				importSpecInfo = append(importSpecInfo, ImportSpecInfo{
					FileName:   path,
					ImportSpec: imp,
					Alias:      alias,
					File:       f,
				})
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return importSpecInfo, nil
}
