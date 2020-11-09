package parser

import (
	"go/token"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

// Parser is the Wails project parser
type Parser struct {

	// Placeholders for Go's parser
	fileSet *token.FileSet

	// The packages we parse
	// The map key is the package ID
	packages map[string]*Package
}

// NewParser creates a new Wails project parser
func NewParser() *Parser {
	return &Parser{
		fileSet:  token.NewFileSet(),
		packages: make(map[string]*Package),
	}
}

// ParseProject will parse the Wails project in the given directory
func (p *Parser) ParseProject(dir string) error {

	var err error

	err = p.loadPackages(dir)
	if err != nil {
		return err
	}

	// Find all the bound structs
	for _, pkg := range p.packages {
		err = p.findBoundStructs(pkg)
		if err != nil {
			return err
		}
	}

	// Parse the structs
	for _, pkg := range p.packages {
		err = p.parseBoundStructs(pkg)
		if err != nil {
			return err
		}
	}

	// Resolve package names
	// We do this because some packages may have the same name
	p.resolvePackageNames()

	spew.Dump(p.packages)

	return nil
}

func (p *Parser) loadPackages(projectPath string) error {
	mode := packages.NeedName |
		packages.NeedFiles |
		packages.NeedSyntax |
		packages.NeedTypes |
		packages.NeedImports |
		packages.NeedTypesInfo |
		packages.NeedModule

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

	// Create a map of packages
	for _, pkg := range pkgs {
		p.packages[pkg.ID] = newPackage(pkg)
	}

	return nil
}

func (p *Parser) getPackageByID(id string) *Package {
	return p.packages[id]
}
