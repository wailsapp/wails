package backendjs

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"reflect"
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
		Comments: []string{"mypackage is awesome"},
		Methods: []*Method{
			{
				Name:     "Naked",
				Comments: []string{"Naked is a method that does nothing"},
			},
		},
	})
	result = append(result, &Package{
		Name:     "otherpackage",
		Comments: []string{"otherpackage is awesome"},
		Methods: []*Method{
			{
				Name:     "OneInput",
				Comments: []string{"OneInput does stuff"},
				Inputs: []*Parameter{
					{
						Name: "name",
						Type: reflect.String,
					},
				},
			},
			{
				Name: "TwoInputs",
				Inputs: []*Parameter{
					{
						Name: "name",
						Type: reflect.String,
					},
					{
						Name: "age",
						Type: reflect.Uint8,
					},
				},
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

	err = generatePackageFiles(packages)
	if err != nil {
		return errors.Wrap(err, "Error generating packages")
	}

	return nil
}

func generateIndexJS(packages []*Package) error {

	// Get path to local file
	templateFile := fs.RelativePath("./index.template")

	// Load template
	javascriptTemplateData := fs.MustLoadString(templateFile)
	packagesTemplate, err := template.New("index").Parse(javascriptTemplateData)
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
		return errors.Wrap(err, "Error calculating index js path")
	}

	err = ioutil.WriteFile(indexJS, buffer.Bytes(), 0755)
	if err != nil {
		return errors.Wrap(err, "Error writing backend package index.js file")
	}

	return nil
}

func generatePackageFiles(packages []*Package) error {

	// Get path to local file
	javascriptTemplateFile := fs.RelativePath("./package.template")

	// Load javascript template
	javascriptTemplateData := fs.MustLoadString(javascriptTemplateFile)
	javascriptTemplate, err := template.New("javascript").Parse(javascriptTemplateData)
	if err != nil {
		return errors.Wrap(err, "Error creating template")
	}

	// Get path to local file
	typescriptTemplateFile := fs.RelativePath("./package.d.template")

	// Load typescript template
	typescriptTemplateData := fs.MustLoadString(typescriptTemplateFile)
	typescriptTemplate, err := template.New("typescript").Parse(typescriptTemplateData)
	if err != nil {
		return errors.Wrap(err, "Error creating template")
	}

	// Iterate over each package
	for _, thisPackage := range packages {

		// Calculate target directory
		packageDir, err := fs.RelativeToCwd("./frontend/backend/" + thisPackage.Name)
		if err != nil {
			return errors.Wrap(err, "Error calculating package path")
		}

		// Make the dir but ignore if it already exists
		fs.Mkdir(packageDir)

		// Execute javascript template
		var buffer bytes.Buffer
		err = javascriptTemplate.Execute(&buffer, thisPackage)
		if err != nil {
			return errors.Wrap(err, "Error generating code")
		}

		// Save javascript file
		err = ioutil.WriteFile(filepath.Join(packageDir, "index.js"), buffer.Bytes(), 0755)
		if err != nil {
			return errors.Wrap(err, "Error writing backend package file")
		}

		// Clear buffer
		buffer.Reset()

		// Execute typescript template
		err = typescriptTemplate.Execute(&buffer, thisPackage)
		if err != nil {
			return errors.Wrap(err, "Error generating code")
		}

		// Save typescript file
		err = ioutil.WriteFile(filepath.Join(packageDir, "index.d.ts"), buffer.Bytes(), 0755)
		if err != nil {
			return errors.Wrap(err, "Error writing backend package file")
		}
	}

	return nil
}
