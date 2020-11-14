package parser

import (
	"go/ast"
	"strings"

	"github.com/leaanthony/slicer"
	"golang.org/x/tools/go/packages"
)

// Package is a wrapper around the go parsed package
type Package struct {

	// A unique Name for this package.
	// This is calculated and may not be the same as the one
	// defined in Go - but that's ok!
	Name string

	// the package we are wrapping
	Gopackage *packages.Package

	// a list of struct names that are bound in this package
	boundStructs slicer.StringSlicer

	// Structs used in this package
	parsedStructs map[string]*Struct

	// A list of external packages we reference from this package
	externalReferences slicer.InterfaceSlicer
}

func newPackage(pkg *packages.Package) *Package {
	return &Package{
		Gopackage:     pkg,
		parsedStructs: make(map[string]*Struct),
	}
}

func (p *Package) getWailsImportName(file *ast.File) string {
	// Scan the imports for the wails v2 import
	for _, details := range file.Imports {
		if details.Path.Value == `"github.com/wailsapp/wails/v2"` {
			if details.Name != nil {
				return details.Name.Name
			}

			// Get the import name from the package
			imp := p.getImportByPath("github.com/wailsapp/wails/v2")
			if imp != nil {
				return imp.Name
			}
		}
	}
	return ""
}

func (p *Package) getImportByName(importName string, file *ast.File) *packages.Package {

	// Check if the file has aliased the import
	for _, imp := range file.Imports {
		if imp.Name != nil {
			if imp.Name.Name == importName {
				// Yes it has. Get the import by path
				return p.getImportByPath(imp.Path.Value)
			}
		}
	}

	// We need to find which package import has this name
	for _, imp := range p.Gopackage.Imports {
		if imp.Name == importName {
			return imp
		}
	}

	// Looks like this package is outside the project...
	return nil
}

func (p *Package) getImportByPath(packagePath string) *packages.Package {
	packagePath = strings.Trim(packagePath, "\"")
	return p.Gopackage.Imports[packagePath]
}

func (p *Package) getStruct(structName string) *Struct {
	return p.parsedStructs[structName]
}

func (p *Package) addStruct(strct *Struct) {
	p.parsedStructs[strct.Name] = strct
}

// HasBoundStructs returns true if any of its structs
// are bound
func (p *Package) HasBoundStructs() bool {

	for _, strct := range p.parsedStructs {
		if strct.IsBound {
			return true
		}
	}

	return false
}

// HasDataStructs returns true if any of its structs
// are used as data
func (p *Package) HasDataStructs() bool {
	for _, strct := range p.parsedStructs {
		if strct.IsUsedAsData {
			return true
		}
	}

	return false
}

// ShouldGenerate returns true when this package should be generated
func (p *Package) ShouldGenerate() bool {
	return p.HasBoundStructs() || p.HasDataStructs()
}

// DeclarationReferences returns a list of external packages
// we reference from this package
func (p *Package) DeclarationReferences() []string {

	var referenceNames slicer.StringSlicer

	// Generics can't come soon enough!
	p.externalReferences.Each(func(p interface{}) {
		referenceNames.Add(p.(*Package).Name)
	})

	return referenceNames.AsSlice()
}

// addExternalReference saves the given package as an external reference
func (p *Package) addExternalReference(pkg *Package) {
	p.externalReferences.AddUnique(pkg)
}

// Structs returns the structs that we want to generate
func (p *Package) Structs() []*Struct {

	var result []*Struct

	for _, elem := range p.parsedStructs {
		result = append(result, elem)
	}

	return result
}
