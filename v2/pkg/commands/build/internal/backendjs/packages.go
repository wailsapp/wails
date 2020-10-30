package backendjs

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/internal/fs"
)

// Package defines a single package that contains bound structs
type Package struct {
	Name     string
	Comments []string
	Methods  []*Method
}

func generatePackages() error {

	packages, err := parsePackages()
	if err != nil {
		return errors.Wrap(err, "Error parsing struct packages:")
	}

	err = generateJSFiles(packages)
	if err != nil {
		return errors.Wrap(err, "Error generating struct js file:")
	}

	return nil
}

func parsePackages() ([]*Package, error) {

	// STUB!
	var result []*Package

	result = append(result, &Package{
		Name:     "mypackage",
		Comments: []string{"// mypackage is awesome"},
		Methods: []*Method{
			{
				Name: "Naked",
			},
		},
	})

	return result, nil
}

func generateJSFiles(packages []*Package) error {

	err := generateIndexJS(packages)
	if err != nil {
		return errors.Wrap(err, "Error generating index.js file")
	}
	return nil
}

func generateIndexJS(packages []*Package) error {

	// Get path to local file
	templateFile := fs.RelativePath("./package.template")

	// Load template
	templateData := fs.MustLoadString(templateFile)
	packagesTemplate, err := template.New("packages").Parse(templateData)
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
	indexJS, err := fs.RelativeToCwd("./frontend/backend/index.js")
	if err != nil {
		return errors.Wrap(err, "Error creating backend js directory")
	}

	err = ioutil.WriteFile(indexJS, buffer.Bytes(), 0755)
	if err != nil {
		return errors.Wrap(err, "Error writing backend package index.js file")
	}

	return nil
}
