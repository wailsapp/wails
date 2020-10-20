package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/kennygrant/sanitize"
	"github.com/leaanthony/slicer"
)

// TemplateMetadata holds all the metadata for a Wails template
type TemplateMetadata struct {
	Name                 string                `json:"name"`
	Version              string                `json:"version"`
	ShortDescription     string                `json:"shortdescription"`
	Description          string                `json:"description"`
	Install              string                `json:"install"`
	Build                string                `json:"build"`
	Author               string                `json:"author"`
	Created              string                `json:"created"`
	FrontendDir          string                `json:"frontenddir"`
	Serve                string                `json:"serve"`
	Bridge               string                `json:"bridge"`
	WailsDir             string                `json:"wailsdir"`
	TemplateDependencies []*TemplateDependency `json:"dependencies,omitempty"`

	// List of platforms that this template is supported on.
	// No value means all platforms. A platform name is the same string
	// as `runtime.GOOS` will return, eg: "darwin". NOTE: This is
	// case sensitive.
	Platforms []string `json:"platforms,omitempty"`
}

// PlatformSupported returns true if this template supports the
// currently running platform
func (m *TemplateMetadata) PlatformSupported() bool {

	// Default is all platforms supported
	if len(m.Platforms) == 0 {
		return true
	}

	// Check that the platform is in the list
	platformsSupported := slicer.String(m.Platforms)
	return platformsSupported.Contains(runtime.GOOS)
}

// TemplateDependency defines a binary dependency for the template
// EG: ng for angular
type TemplateDependency struct {
	Bin  string `json:"bin"`
	Help string `json:"help"`
}

// TemplateDetails holds information about a specific template
type TemplateDetails struct {
	Name     string
	Path     string
	Metadata *TemplateMetadata
	fs       *FSHelper
}

// TemplateHelper is a utility object to help with processing templates
type TemplateHelper struct {
	templateDir      *Dir
	fs               *FSHelper
	metadataFilename string
}

// NewTemplateHelper creates a new template helper
func NewTemplateHelper() *TemplateHelper {

	templateDir, err := fs.LocalDir("./templates")
	if err != nil {
		log.Fatal("Unable to find the template directory. Please reinstall Wails.")
	}

	return &TemplateHelper{
		templateDir:      templateDir,
		metadataFilename: "template.json",
	}
}

// IsValidTemplate returns true if the given template name resides on disk
func (t *TemplateHelper) IsValidTemplate(templateName string) bool {
	pathToTemplate := filepath.Join(t.templateDir.fullPath, templateName)
	return t.fs.DirExists(pathToTemplate)
}

// SanitizeFilename sanitizes the given string to make a valid filename
func (t *TemplateHelper) SanitizeFilename(name string) string {
	return sanitize.Name(name)
}

// CreateNewTemplate creates a new template based on the given directory name and string
func (t *TemplateHelper) CreateNewTemplate(dirname string, details *TemplateMetadata) (string, error) {

	// Check if this template has already been created
	if t.IsValidTemplate(dirname) {
		return "", fmt.Errorf("cannot create template in directory '%s' - already exists", dirname)
	}

	targetDir := filepath.Join(t.templateDir.fullPath, dirname)
	err := t.fs.MkDir(targetDir)
	if err != nil {
		return "", err
	}
	targetMetadata := filepath.Join(targetDir, t.metadataFilename)
	err = t.fs.SaveAsJSON(details, targetMetadata)

	return targetDir, err
}

// LoadMetadata loads the template's 'metadata.json' file
func (t *TemplateHelper) LoadMetadata(dir string) (*TemplateMetadata, error) {
	templateFile := filepath.Join(dir, t.metadataFilename)
	result := &TemplateMetadata{}
	if !t.fs.FileExists(templateFile) {
		return nil, nil
	}
	rawJSON, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawJSON, &result)
	return result, err
}

// GetTemplateDetails returns a map of Template structs containing details
// of the found templates
func (t *TemplateHelper) GetTemplateDetails() (map[string]*TemplateDetails, error) {

	// Get the subdirectory details
	templateDirs, err := t.templateDir.GetSubdirs()
	if err != nil {
		return nil, err
	}

	result := make(map[string]*TemplateDetails)

	for name, dir := range templateDirs {
		result[name] = &TemplateDetails{
			Path: dir,
		}
		metadata, err := t.LoadMetadata(dir)
		if err != nil {
			return nil, err
		}

		result[name].Metadata = metadata
		if metadata.Name != "" {
			result[name].Name = metadata.Name
		} else {
			// Ignore bad templates?
			result[name] = nil
		}
	}

	return result, nil
}

// GetTemplateFilenames returns all the filenames of the given template
func (t *TemplateHelper) GetTemplateFilenames(template *TemplateDetails) (*slicer.StringSlicer, error) {

	// Get the subdirectory details
	templateDir, err := t.fs.Directory(template.Path)
	if err != nil {
		return nil, err
	}
	return templateDir.GetAllFilenames()
}

// InstallTemplate installs the template given in the project options to the
// project path given
func (t *TemplateHelper) InstallTemplate(projectPath string, projectOptions *ProjectOptions) error {

	// Check dependencies before installing
	dependencies := projectOptions.selectedTemplate.Metadata.TemplateDependencies
	if dependencies != nil {
		programHelper := NewProgramHelper()
		logger := NewLogger()
		errors := []string{}
		for _, dep := range dependencies {
			program := programHelper.FindProgram(dep.Bin)
			if program == nil {
				errors = append(errors, dep.Help)
			}
		}
		if len(errors) > 0 {
			mainError := "template dependencies not installed"
			if len(errors) == 1 {
				mainError = errors[0]
			} else {
				for _, error := range errors {
					logger.Red(error)
				}
			}
			return fmt.Errorf(mainError)
		}
	}

	// Get template files
	templateFilenames, err := t.GetTemplateFilenames(projectOptions.selectedTemplate)
	if err != nil {
		return err
	}

	templatePath := projectOptions.selectedTemplate.Path

	// Save the version
	projectOptions.WailsVersion = Version

	templateJSONFilename := filepath.Join(templatePath, t.metadataFilename)

	templateFiles := templateFilenames.Filter(func(filename string) bool {
		filename = filepath.FromSlash(filename)
		return strings.HasPrefix(filename, templatePath) && filename != templateJSONFilename
	})

	templateFiles.Each(func(templateFile string) {

		// Setup filenames
		relativeFilename := strings.TrimPrefix(templateFile, templatePath)[1:]
		targetFilename, err := filepath.Abs(filepath.Join(projectOptions.OutputDirectory, relativeFilename))
		if err != nil {
			return
		}
		filedata, err := t.fs.LoadAsBytes(templateFile)
		if err != nil {
			return
		}

		// If file is a template, process it
		if strings.HasSuffix(templateFile, ".template") {
			templateData := string(filedata)
			tmpl := template.New(templateFile)
			tmpl.Parse(templateData)
			var tpl bytes.Buffer
			err = tmpl.Execute(&tpl, projectOptions)
			if err != nil {
				return
			}

			// Remove template suffix
			targetFilename = strings.TrimSuffix(targetFilename, ".template")

			// Set the filedata to the template result
			filedata = tpl.Bytes()
		}

		// Normal file, just copy it
		err = fs.CreateFile(targetFilename, filedata)
		if err != nil {
			return
		}
	})

	if err != nil {
		return err
	}

	return nil
}
