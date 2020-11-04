package backendjs

import (
	"go/token"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/internal/fs"
	"golang.org/x/tools/go/packages"
)

// GenerateBackendJSPackage will generate a Javascript/Typescript
// package in `<project>/frontend/backend` that defines which methods
// and structs are bound to your frontend
func GenerateBackendJSPackage() error {

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	p := NewParser()

	err = p.parseProject(dir)
	if err != nil {
		return err
	}

	err = p.generateModule()

	return err
}

func (p *Parser) parseProject(projectPath string) error {
	mode := packages.NeedName |
		packages.NeedFiles |
		packages.NeedSyntax |
		packages.NeedTypes |
		packages.NeedImports |
		packages.NeedTypesInfo

	var fset = token.NewFileSet()
	cfg := &packages.Config{Fset: fset, Mode: mode, Dir: projectPath}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return errors.Wrap(err, "Problem loading packages")
	}
	if packages.PrintErrors(pkgs) > 0 {
		return errors.Wrap(err, "Errors during parsing")
	}
	for _, pkg := range pkgs {
		parsedPackage, err := p.parsePackage(pkg, fset)
		if err != nil {
			return err
		}
		p.Packages[parsedPackage.Name] = parsedPackage
	}

	// Resolve all the loose ends from parsing
	err = p.resolve()
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) generateModule() error {

	moduleDir, err := createBackendJSDirectory()
	if err != nil {
		return err
	}

	for _, pkg := range p.Packages {

		// Calculate directory
		dir := filepath.Join(moduleDir, pkg.Name)

		// Create the directory if it doesn't exist
		fs.Mkdir(dir)

		err := generatePackage(pkg, dir)
		if err != nil {
			return err
		}
	}
	return nil
}

func createBackendJSDirectory() (string, error) {

	// Calculate the package directory
	// Note this is *always* called from the project directory
	// so using paths relative to CWD is fine
	dir, err := fs.RelativeToCwd("./frontend/backend")
	if err != nil {
		return "", errors.Wrap(err, "Error creating backend js directory")
	}

	// Remove directory if it exists - REGENERATION!
	err = os.RemoveAll(dir)
	if err != nil {
		return "", errors.Wrap(err, "Error removing module directory")
	}

	// Make the directory
	err = fs.Mkdir(dir)

	return dir, err
}
