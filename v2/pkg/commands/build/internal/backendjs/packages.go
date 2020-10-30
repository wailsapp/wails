package backendjs

import (
	"bytes"
	"io/ioutil"
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

	err = generatePackageFiles(packages)
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
			{
				Name: "TwoInputsAndOutput",
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
				Outputs: []*Parameter{
					{
						Name: "result",
						Type: reflect.Bool,
					},
				},
			},
		},
	})

	return result, nil
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
		packageFile, err := fs.RelativeToCwd("./frontend/backend/" + thisPackage.Name)
		if err != nil {
			return errors.Wrap(err, "Error calculating package path")
		}

		// Execute javascript template
		var buffer bytes.Buffer
		err = javascriptTemplate.Execute(&buffer, thisPackage)
		if err != nil {
			return errors.Wrap(err, "Error generating code")
		}

		// Save javascript file
		err = ioutil.WriteFile(packageFile+".js", buffer.Bytes(), 0755)
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
		err = ioutil.WriteFile(packageFile+".d.ts", buffer.Bytes(), 0755)
		if err != nil {
			return errors.Wrap(err, "Error writing backend package file")
		}
	}

	return nil
}
