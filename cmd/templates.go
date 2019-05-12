package cmd

import (
	"bytes"
	"encoding/json"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	mewn "github.com/leaanthony/mewn"
	mewnlib "github.com/leaanthony/mewn/lib"
	"github.com/leaanthony/slicer"
)

// TemplateMetadata holds all the metadata for a Wails template
type TemplateMetadata struct {
	Name             string `json:"name"`
	ShortDescription string `json:"shortdescription"`
	Description      string `json:"description"`
	Install          string `json:"install"`
	Build            string `json:"build"`
	Author           string `json:"author"`
	Created          string `json:"created"`
	FrontendDir      string `json:"frontenddir"`
	Serve            string `json:"serve"`
	Bridge           string `json:"bridge"`
}

// TemplateDetails holds information about a specific template
type TemplateDetails struct {
	BasePath string
	Path     string
	Metadata *TemplateMetadata
}

// TemplateList is a list of available templates
type TemplateList struct {
	details map[string]*TemplateDetails
}

// NewTemplateList creates a new TemplateList object
func NewTemplateList(filenames *mewnlib.FileGroup) *TemplateList {
	// Iterate each template and store information

	result := &TemplateList{details: make(map[string]*TemplateDetails)}

	entries := slicer.String()
	entries.AddSlice(filenames.Entries())

	// Find all template.json files
	metadataFiles := entries.Filter(func(filename string) bool {
		match, _ := regexp.MatchString("(.)+template.json$", filename)
		return match
	})

	// Load each metadata file
	metadataFiles.Each(func(filename string) {
		fileData := filenames.Bytes(filename)
		var metadata TemplateMetadata
		err := json.Unmarshal(fileData, &metadata)
		if err != nil {
			log.Fatalf("corrupt metadata for template: %s", filename)
		}
		path := strings.Split(filename, "/")[0]
		thisTemplate := &TemplateDetails{Path: path, Metadata: &metadata}
		result.details[filename] = thisTemplate
	})

	return result
}

// Template holds details about a Wails template
type Template struct {
	Name        string
	Path        string
	Description string
}

// TemplateHelper is a utility object to help with processing templates
type TemplateHelper struct {
	TemplateList *TemplateList
	Files        *mewnlib.FileGroup
}

// NewTemplateHelper creates a new template helper
func NewTemplateHelper() *TemplateHelper {
	files := mewn.Group("./templates")

	return &TemplateHelper{
		TemplateList: NewTemplateList(files),
		Files:        files,
	}
}

// InstallTemplate installs the template given in the project options to the
// project path given
func (t *TemplateHelper) InstallTemplate(projectPath string, projectOptions *ProjectOptions) error {

	// Get template files
	templatePath := projectOptions.selectedTemplate.Path

	templateFilenames := slicer.String()
	templateFilenames.AddSlice(projectOptions.templates.Files.Entries())

	templateJSONFilename := filepath.Join(templatePath, "template.json")

	templateFiles := templateFilenames.Filter(func(filename string) bool {
		filename = filepath.FromSlash(filename)
		return strings.HasPrefix(filename, templatePath) && filename != templateJSONFilename
	})

	var err error
	templateFiles.Each(func(templateFile string) {

		// Setup filenames
		relativeFilename := strings.TrimPrefix(templateFile, templatePath)[1:]
		targetFilename, err := filepath.Abs(filepath.Join(projectOptions.OutputDirectory, relativeFilename))
		if err != nil {
			return
		}
		filedata := projectOptions.templates.Files.Bytes(templateFile)

		// If file is a template, process it
		if strings.HasSuffix(templateFile, ".template") {
			templateData := projectOptions.templates.Files.String(templateFile)
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
