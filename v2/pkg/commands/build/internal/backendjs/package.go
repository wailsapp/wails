package backendjs

import (
	"bytes"
	"go/ast"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"text/template"

	"github.com/davecgh/go-spew/spew"
	"github.com/leaanthony/slicer"
	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/internal/fs"
	"golang.org/x/tools/go/packages"
)

// Package defines a parsed package
type Package struct {
	Name    string
	Structs map[string]*Struct

	// These are the structs declared in this package
	// that are used as data by either this or other packages
	structsUsedAsData slicer.StringSlicer

	// A list of functions that return struct pointers
	functionsThatReturnStructPointers map[string]string

	// A list of functions that return structs
	functionsThatReturnStructs map[string]string

	// A list of struct literals that were bound to the application
	// EG: app.Bind( &mystruct{} )
	structLiteralsThatWereBound slicer.StringSlicer

	// A list of struct pointer literals that were bound to the application
	// EG: app.Bind( &mystruct{} )
	structPointerLiteralsThatWereBound slicer.StringSlicer

	// A list of methods that returns structs to the Bind method
	// EG: app.Bind( newMyStruct() )
	structMethodsThatWereBound slicer.StringSlicer

	// A list of variables that were used for binding
	// Eg: myVar := &mystruct{}; app.Bind( myVar )
	variablesThatWereBound slicer.StringSlicer

	// A list of variables that were assigned using a function call
	// EG: myVar := newStruct()
	variablesThatWereAssignedByFunctions map[string]string

	// A map of variables that were assigned using a struct literal
	// EG: myVar := MyStruct{}
	variablesThatWereAssignedByStructLiterals map[string]string

	// A map of variables that were assigned using a struct literal
	// in a different package
	// EG: myVar := mypackage.MyStruct{}
	variablesThatWereAssignedByExternalStructLiterals map[string]*StructName
}

func newPackage(name string) *Package {
	return &Package{
		Name:                                 name,
		Structs:                              make(map[string]*Struct),
		functionsThatReturnStructPointers:    make(map[string]string),
		functionsThatReturnStructs:           make(map[string]string),
		variablesThatWereAssignedByFunctions: make(map[string]string),
		variablesThatWereAssignedByStructLiterals:         make(map[string]string),
		variablesThatWereAssignedByExternalStructLiterals: make(map[string]*StructName),
	}
}

func (p *Parser) parsePackage(pkg *packages.Package, fset *token.FileSet) (*Package, error) {
	result := p.Packages[pkg.Name]
	if result == nil {
		result = newPackage(pkg.Name)
	}

	// Get the absolute path to the project's main.go file
	pathToMain, err := fs.RelativeToCwd("main.go")
	if err != nil {
		return nil, err
	}

	// Work out if this is the main package
	goFiles := slicer.String(pkg.GoFiles)
	if goFiles.Contains(pathToMain) {
		// This is the program entrypoint file
		// Scan the imports for the wails v2 import
		for key, details := range pkg.Imports {
			if key == "github.com/wailsapp/wails/v2" {
				p.wailsPackageVariable = details.Name
			}
		}
	}

	for _, fileAst := range pkg.Syntax {
		var parseError error
		ast.Inspect(fileAst, func(n ast.Node) bool {
			// if typeDecl, ok := n.(*ast.TypeSpec); ok {
			// 	// Parse struct definitions
			// 	if structType, ok := typeDecl.Type.(*ast.StructType); ok {
			// 		structName := typeDecl.Name.Name
			// 		// findInFields(structTy.Fields, n, pkg.TypesInfo, fset)
			// 		structDef, err := p.ParseStruct(structType, structName, result)
			// 		if err != nil {
			// 			parseError = err
			// 			return false
			// 		}

			// 		// Parse comments
			// 		structDef.Comments = p.parseComments(typeDecl.Doc)

			// 		result.Structs[structName] = structDef
			// 	}
			// }

			if genDecl, ok := n.(*ast.GenDecl); ok {
				println("GenDecl:")
				for _, spec := range genDecl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if structType, ok := typeSpec.Type.(*ast.StructType); ok {
							structName := typeSpec.Name.Name
							structDef, err := p.ParseStruct(structType, structName, result)
							if err != nil {
								parseError = err
								return false
							}

							// Parse comments
							structDef.Comments = p.parseComments(genDecl.Doc)

							result.Structs[structName] = structDef
						}
					}
				}
				spew.Dump(genDecl)
			}

			// Capture call expressions
			if callExpr, ok := n.(*ast.CallExpr); ok {
				p.parseCallExpressions(callExpr, result)
			}

			// Parse Assignments
			if assignStmt, ok := n.(*ast.AssignStmt); ok {
				p.parseAssignment(assignStmt, result)
			}

			// Parse Function declarations
			if funcDecl, ok := n.(*ast.FuncDecl); ok {
				p.parseFunctionDeclaration(funcDecl, result)
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

	// Get path to local file
	javascriptTemplateFile := fs.RelativePath("./package.template")

	// Load javascript template
	javascriptTemplateData := fs.MustLoadString(javascriptTemplateFile)
	javascriptTemplate, err := template.New("javascript").Parse(javascriptTemplateData)
	if err != nil {
		return errors.Wrap(err, "Error creating template")
	}

	// Reset the buffer
	buffer.Reset()

	err = javascriptTemplate.Execute(&buffer, pkg)
	if err != nil {
		return errors.Wrap(err, "Error generating code")
	}

	// Save javascript file
	err = ioutil.WriteFile(filepath.Join(moduledir, "index.js"), buffer.Bytes(), 0755)
	if err != nil {
		return errors.Wrap(err, "Error writing backend package file")
	}

	return nil
}

// DeclarationReferences returns the typescript declaration references for the package
func (p *Package) DeclarationReferences() []string {
	var result []string
	for _, strct := range p.Structs {
		if strct.IsBound {
			refs := strct.packageReferences.AsSlice()
			result = append(result, refs...)
		}
	}
	return result
}

// StructIsUsedAsData returns true if the given struct name has
// been used in structs, inputs or outputs by other packages
func (p *Package) StructIsUsedAsData(structName string) bool {
	return p.structsUsedAsData.Contains(structName)
}

func (p *Package) resolveBoundStructLiterals() {
	p.structLiteralsThatWereBound.Each(func(structName string) {
		strct := p.Structs[structName]
		if strct == nil {
			println("Warning: Cannot find bound struct", structName, "in package", p.Name)
			return
		}
		println("Bound struct", strct.Name, "in package", p.Name)
		strct.IsBound = true
	})
}

func (p *Package) resolveBoundStructPointerLiterals() {
	p.structPointerLiteralsThatWereBound.Each(func(structName string) {
		strct := p.Structs[structName]
		if strct == nil {
			println("Warning: Cannot find bound struct", structName, "in package", p.Name)
			return
		}
		println("Bound struct pointer", strct.Name, "in package", p.Name)
		strct.IsBound = true
	})
}

// ShouldBeGenerated indicates if the package should be generated
// The package should be generated only if we have structs that are
// bound or structs that are used as data
func (p *Package) ShouldBeGenerated() bool {
	for _, strct := range p.Structs {
		if strct.IsBound || strct.IsUsedAsData {
			return true
		}
	}
	return false
}
