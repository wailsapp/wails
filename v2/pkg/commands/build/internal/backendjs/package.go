package backendjs

import (
	"bytes"
	"go/ast"
	"go/token"
	"html/template"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/internal/fs"
	"golang.org/x/tools/go/packages"
)

// Package defines a parsed package
type Package struct {
	Name    string
	Structs map[string]*Struct
}

func parsePackage(pkg *packages.Package, fset *token.FileSet) (*Package, error) {
	result := &Package{
		Name:    pkg.Name,
		Structs: make(map[string]*Struct),
	}
	for _, fileAst := range pkg.Syntax {
		var parseError error
		ast.Inspect(fileAst, func(n ast.Node) bool {
			if typeDecl, ok := n.(*ast.TypeSpec); ok {
				if structType, ok := typeDecl.Type.(*ast.StructType); ok {

					// spew.Dump(structType)
					structName := typeDecl.Name.Name
					// findInFields(structTy.Fields, n, pkg.TypesInfo, fset)
					structDef, err := parseStruct(structType, structName)
					if err != nil {
						parseError = err
						return false
					}

					// Parse comments
					structDef.Comments = parseComments(typeDecl.Doc)

					result.Structs[structName] = structDef
				}
			}
			return true
		})
		if parseError != nil {
			return nil, parseError
		}
	}
	return result, nil
}

func generatePackage(pkg *Package, moduledir string) error {

	// Get path to local file
	typescriptTemplateFile := fs.RelativePath("./package.d.template")

	// Load typescript template
	typescriptTemplateData := fs.MustLoadString(typescriptTemplateFile)
	typescriptTemplate, err := template.New("typescript").Parse(typescriptTemplateData)
	if err != nil {
		return errors.Wrap(err, "Error creating template")
	}

	// Execute javascript template
	var buffer bytes.Buffer
	err = typescriptTemplate.Execute(&buffer, pkg)
	if err != nil {
		return errors.Wrap(err, "Error generating code")
	}

	// Save typescript file
	err = ioutil.WriteFile(filepath.Join(moduledir, "index.d.ts"), buffer.Bytes(), 0755)
	if err != nil {
		return errors.Wrap(err, "Error writing backend package file")
	}
	return nil
}
