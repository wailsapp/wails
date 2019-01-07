package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alecthomas/template"
)

const templateSuffix = ".template"

// TemplateHelper helps with creating projects
type TemplateHelper struct {
	system      *SystemHelper
	fs          *FSHelper
	templateDir string
	// templates   map[string]string
	templateSuffix   string
	metadataFilename string
}

// Template defines a single template
type Template struct {
	Name     string
	Dir      string
	Metadata map[string]interface{}
}

// NewTemplateHelper creates a new template helper
func NewTemplateHelper() *TemplateHelper {
	result := TemplateHelper{
		system:           NewSystemHelper(),
		fs:               NewFSHelper(),
		templateSuffix:   ".template",
		metadataFilename: "template.json",
	}
	// Calculate template base dir
	_, filename, _, _ := runtime.Caller(1)
	result.templateDir = filepath.Join(path.Dir(filename), "templates")
	// result.templateDir = filepath.Join(result.system.homeDir, "go", "src", "github.com", "wailsapp", "wails", "cmd", "templates")
	return &result
}

// GetTemplateNames returns a map of all available templates
func (t *TemplateHelper) GetTemplateNames() (map[string]string, error) {
	templateDirs, err := t.fs.GetSubdirs(t.templateDir)
	if err != nil {
		return nil, err
	}
	return templateDirs, nil
}

// GetTemplateDetails returns a map of Template structs containing details
// of the found templates
func (t *TemplateHelper) GetTemplateDetails() (map[string]*Template, error) {
	templateDirs, err := t.fs.GetSubdirs(t.templateDir)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*Template)

	for name, dir := range templateDirs {
		result[name] = &Template{
			Dir: dir,
		}
		metadata, err := t.LoadMetadata(dir)
		if err != nil {
			return nil, err
		}
		result[name].Metadata = metadata
		if metadata["name"] != nil {
			result[name].Name = metadata["name"].(string)
		} else {
			// Ignore bad templates?
			result[name] = nil
		}
	}

	return result, nil
}

// LoadMetadata loads the template's 'metadata.json' file
func (t *TemplateHelper) LoadMetadata(dir string) (map[string]interface{}, error) {
	templateFile := filepath.Join(dir, t.metadataFilename)
	result := make(map[string]interface{})
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

// TemplateExists returns true if the given template name exists
func (t *TemplateHelper) TemplateExists(templateName string) (bool, error) {
	templates, err := t.GetTemplateNames()
	if err != nil {
		return false, err
	}
	_, exists := templates[templateName]
	return exists, nil
}

// InstallTemplate installs the template given in the project options to the
// project path given
func (t *TemplateHelper) InstallTemplate(projectPath string, projectOptions *ProjectOptions) error {

	// Get template files
	template, err := t.getTemplateFiles(projectOptions.Template)
	if err != nil {
		return err
	}

	// Copy files to target
	err = template.Install(projectPath, projectOptions)
	if err != nil {
		return err
	}

	return nil
}

// templateFiles categorises files found in a template
type templateFiles struct {
	BaseDir       string
	StandardFiles []string
	Templates     []string
	Dirs          []string
}

// newTemplateFiles returns a new TemplateFiles struct
func (t *TemplateHelper) newTemplateFiles(dir string) *templateFiles {
	pathsep := string(os.PathSeparator)
	// Ensure base directory has trailing slash
	if !strings.HasSuffix(dir, pathsep) {
		dir = dir + pathsep
	}
	return &templateFiles{
		BaseDir: dir,
	}
}

// AddStandardFile adds the given file to the list of standard files
func (t *templateFiles) AddStandardFile(filename string) {
	localPath := strings.TrimPrefix(filename, t.BaseDir)
	t.StandardFiles = append(t.StandardFiles, localPath)
}

// AddTemplate adds the given file to the list of template files
func (t *templateFiles) AddTemplate(filename string) {
	localPath := strings.TrimPrefix(filename, t.BaseDir)
	t.Templates = append(t.Templates, localPath)
}

// AddDir adds the given directory to the list of template dirs
func (t *templateFiles) AddDir(dir string) {
	localPath := strings.TrimPrefix(dir, t.BaseDir)
	t.Dirs = append(t.Dirs, localPath)
}

// getTemplateFiles returns a struct categorising files in
// the template directory
func (t *TemplateHelper) getTemplateFiles(templateName string) (*templateFiles, error) {

	templates, err := t.GetTemplateNames()
	if err != nil {
		return nil, err
	}
	templateDir := templates[templateName]
	result := t.newTemplateFiles(templateDir)
	var localPath string
	err = filepath.Walk(templateDir, func(dir string, info os.FileInfo, err error) error {
		if dir == templateDir {
			return nil
		}
		if err != nil {
			return err
		}

		// Don't copy template metadata
		localPath = strings.TrimPrefix(dir, templateDir+string(filepath.Separator))
		if localPath == t.metadataFilename {
			return nil
		}

		// Categorise the file
		switch {
		case info.IsDir():
			result.AddDir(dir)
		case strings.HasSuffix(info.Name(), templateSuffix):
			result.AddTemplate(dir)
		default:
			result.AddStandardFile(dir)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error processing template '%s' in path '%q': %v", templateName, templateDir, err)
	}
	return result, err
}

// Install the template files into the given project path
func (t *templateFiles) Install(projectPath string, projectOptions *ProjectOptions) error {

	fs := NewFSHelper()

	// Create directories
	var targetDir string
	for _, dirname := range t.Dirs {
		targetDir = filepath.Join(projectPath, dirname)
		fs.MkDir(targetDir)
	}

	// Copy standard files
	var targetFile, sourceFile string
	var err error
	for _, filename := range t.StandardFiles {
		sourceFile = filepath.Join(t.BaseDir, filename)
		targetFile = filepath.Join(projectPath, filename)

		err = fs.CopyFile(sourceFile, targetFile)
		if err != nil {
			return err
		}
	}

	// Do we have template files?
	if len(t.Templates) > 0 {

		// Iterate over the templates
		var templateFile string
		var tmpl *template.Template
		for _, filename := range t.Templates {

			// Load template text
			templateFile = filepath.Join(t.BaseDir, filename)
			templateText, err := fs.LoadAsString(templateFile)
			if err != nil {
				return err
			}

			// Apply template
			tmpl = template.New(templateFile)
			tmpl.Parse(templateText)

			// Write the template to a buffer
			var tpl bytes.Buffer
			err = tmpl.Execute(&tpl, projectOptions)
			if err != nil {
				fmt.Println("ERROR!!! " + err.Error())
				return err
			}

			// Save buffer to disk
			targetFilename := strings.TrimSuffix(filename, templateSuffix)
			targetFile = filepath.Join(projectPath, targetFilename)
			err = ioutil.WriteFile(targetFile, tpl.Bytes(), 0644)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
