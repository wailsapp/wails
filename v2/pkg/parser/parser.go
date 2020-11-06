package parser

import (
	"fmt"
	"go/token"

	"github.com/davecgh/go-spew/spew"
	"github.com/leaanthony/slicer"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

type Parser struct {

	// Placeholders for Go's parser
	goPackages      []*packages.Package
	fileSet         *token.FileSet
	internalMethods *slicer.StringSlicer

	// This is a map of structs that have been parsed
	// The key is <package>.<structname>
	parsedStructs map[string]*Struct

	// The list of struct names that are bound
	BoundStructReferences []*StructReference

	// The list of structs that are bound
	BoundStructs []*Struct
}

func NewParser() *Parser {
	return &Parser{
		fileSet:         token.NewFileSet(),
		internalMethods: slicer.String([]string{"WailsInit", "WailsShutdown"}),
		parsedStructs:   make(map[string]*Struct),
	}
}

// ParseProject will parse the Wails project in the given directory
func (p *Parser) ParseProject(dir string) error {

	var err error

	err = p.loadPackages(dir)
	if err != nil {
		return err
	}

	err = p.findBoundStructs()
	if err != nil {
		return err
	}

	err = p.parseBoundStructs()
	if err != nil {
		return err
	}

	spew.Dump(p.BoundStructs)
	println("******* Parsed Structs *******")
	fmt.Printf("%+v\n", p.parsedStructs)

	return err
}

func (p *Parser) loadPackages(projectPath string) error {
	mode := packages.NeedName |
		packages.NeedFiles |
		packages.NeedSyntax |
		packages.NeedTypes |
		packages.NeedImports |
		packages.NeedTypesInfo

	cfg := &packages.Config{Fset: p.fileSet, Mode: mode, Dir: projectPath}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return errors.Wrap(err, "Problem loading packages")
	}
	// Check for errors
	var parseError error
	for _, pkg := range pkgs {
		for _, err := range pkg.Errors {
			if parseError == nil {
				parseError = errors.New(err.Error())
			} else {
				parseError = errors.Wrap(parseError, err.Error())
			}
		}
	}

	if parseError != nil {
		return parseError
	}

	p.goPackages = pkgs

	return nil
}

func (p *Parser) getPackageByName(packageName string) *packages.Package {
	for _, pkg := range p.goPackages {
		if pkg.Name == packageName {
			return pkg
		}
	}
	return nil
}

func (p *Parser) getWailsImportName(pkg *packages.Package) (string, bool) {
	// Scan the imports for the wails v2 import
	for key, details := range pkg.Imports {
		if key == "github.com/wailsapp/wails/v2" {
			return details.Name, true
		}
	}
	return "", false
}

// findBoundStructs will search through the Wails project looking
// for which structs have been bound using the `Bind()` method
func (p *Parser) findBoundStructs() error {

	// Try each of the packages to find the Bind() calls
	for _, pkg := range p.goPackages {

		// Does this package import Wails?
		wailsImportName, imported := p.getWailsImportName(pkg)
		if !imported {
			continue
		}

		// Do we create an app using CreateApp?
		appVariableName, created := p.getApplicationVariableName(pkg, wailsImportName)
		if !created {
			continue
		}

		boundStructReferences := p.findBoundStructsInPackage(pkg, appVariableName)
		p.BoundStructReferences = append(p.BoundStructReferences, boundStructReferences...)
	}

	return nil
}

func (p *Parser) parseBoundStructs() error {

	// Iterate the structs
	for _, boundStructReference := range p.BoundStructReferences {
		// Parse the struct
		boundStruct, err := p.ParseStruct(boundStructReference.Package, boundStructReference.Name)
		if err != nil {
			return err
		}

		p.BoundStructs = append(p.BoundStructs, boundStruct)
	}

	// Resolve the references between the structs
	// This is when a field of one struct is a struct type
	for _, boundStruct := range p.BoundStructs {
		err := p.resolveStructReferences(boundStruct)
		if err != nil {
			return err
		}
	}

	return nil
}
