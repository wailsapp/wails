package parser

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/internal/fs"
)

// GenerateWailsFrontendPackage will generate a Javascript/Typescript
// package in `<project>/frontend/wails` that defines which methods
// and structs are bound to your frontend
func GenerateWailsFrontendPackage() error {

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	p := NewParser()

	err = p.ParseProject(dir)
	if err != nil {
		return err
	}

	err = p.generateModule()

	return err
}

func (p *Parser) generateModule() error {

	moduleDir, err := createBackendJSDirectory()
	if err != nil {
		return err
	}

	packagesToGenerate := p.packagesToGenerate()

	for _, pkg := range packagesToGenerate {

		err := generatePackage(pkg, moduleDir)
		if err != nil {
			return err
		}
	}

	// Copy the standard files
	srcFile := fs.RelativePath("./package.json")
	tgtFile := filepath.Join(moduleDir, "package.json")
	err = fs.CopyFile(srcFile, tgtFile)
	if err != nil {
		return err
	}

	// Generate the globals.d.ts file
	err = generateGlobalsTS(moduleDir, packagesToGenerate)
	if err != nil {
		return err
	}

	// Generate the index.js file
	err = generateIndexJS(moduleDir, packagesToGenerate)
	if err != nil {
		return err
	}
	// Generate the index.d.ts file
	err = generateIndexTS(moduleDir, packagesToGenerate)
	if err != nil {
		return err
	}

	return nil
}

func createBackendJSDirectory() (string, error) {

	// Calculate the package directory
	// Note this is *always* called from the project directory
	// so using paths relative to CWD is fine
	dir, err := fs.RelativeToCwd("./frontend/backend")
	if err != nil {
		return "", errors.Wrap(err, "Error creating backend module directory")
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
	err = ioutil.WriteFile(filepath.Join(moduledir, "_"+pkg.Name+".d.ts"), buffer.Bytes(), 0755)
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
	err = ioutil.WriteFile(filepath.Join(moduledir, "_"+pkg.Name+".js"), buffer.Bytes(), 0755)
	if err != nil {
		return errors.Wrap(err, "Error writing backend package file")
	}

	return nil
}

func generateIndexJS(dir string, packages []*Package) error {

	// Get path to local file
	templateFile := fs.RelativePath("./index.template")

	// Load template
	templateData := fs.MustLoadString(templateFile)
	packagesTemplate, err := template.New("index").Parse(templateData)
	if err != nil {
		return errors.Wrap(err, "Error creating template")
	}

	// Execute template
	var buffer bytes.Buffer
	err = packagesTemplate.Execute(&buffer, packages)
	if err != nil {
		return errors.Wrap(err, "Error generating code")
	}

	// Calculate target filename
	indexJS := filepath.Join(dir, "index.js")

	err = ioutil.WriteFile(indexJS, buffer.Bytes(), 0755)
	if err != nil {
		return errors.Wrap(err, "Error writing backend package index.js file")
	}

	return nil
}
func generateIndexTS(dir string, packages []*Package) error {

	// Get path to local file
	templateFile := fs.RelativePath("./index.d.template")

	// Load template
	templateData := fs.MustLoadString(templateFile)
	indexTSTemplate, err := template.New("index.d").Parse(templateData)
	if err != nil {
		return errors.Wrap(err, "Error creating template")
	}

	// Execute template
	var buffer bytes.Buffer
	err = indexTSTemplate.Execute(&buffer, packages)
	if err != nil {
		return errors.Wrap(err, "Error generating code")
	}

	// Calculate target filename
	indexJS := filepath.Join(dir, "index.d.ts")

	err = ioutil.WriteFile(indexJS, buffer.Bytes(), 0755)
	if err != nil {
		return errors.Wrap(err, "Error writing backend package index.d.ts file")
	}

	return nil
}

func generateGlobalsTS(dir string, packages []*Package) error {

	// Get path to local file
	templateFile := fs.RelativePath("./globals.d.template")

	// Load template
	templateData := fs.MustLoadString(templateFile)
	packagesTemplate, err := template.New("globals").Parse(templateData)
	if err != nil {
		return errors.Wrap(err, "Error creating template")
	}

	// Execute template
	var buffer bytes.Buffer
	err = packagesTemplate.Execute(&buffer, packages)
	if err != nil {
		return errors.Wrap(err, "Error generating code")
	}

	// Calculate target filename
	indexJS := filepath.Join(dir, "globals.d.ts")

	err = ioutil.WriteFile(indexJS, buffer.Bytes(), 0755)
	if err != nil {
		return errors.Wrap(err, "Error writing backend package globals.d.ts file")
	}

	return nil
}
