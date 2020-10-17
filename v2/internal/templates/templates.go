package templates

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/leaanthony/gosod"
	"github.com/leaanthony/slicer"
	"github.com/olekukonko/tablewriter"
	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
)

// Cahce for the templates
// We use this because we need different views of the same data
var templateCache []Template = nil

// Data contains the data we wish to embed during template installation
type Data struct {
	ProjectName    string
	BinaryName     string
	WailsVersion   string
	NPMProjectName string
	Author         string
	WailsDirectory string
}

// Options for installing a template
type Options struct {
	ProjectName  string
	TemplateName string
	BinaryName   string
	TargetDir    string
	Logger       *clilogger.CLILogger
}

// Template holds data relating to a template
// including the metadata stored in template.json
type Template struct {

	// Template details
	Name        string `json:"name"`
	ShortName   string `json:"shortname"`
	Author      string `json:"author"`
	Description string `json:"description"`
	HelpURL     string `json:"helpurl"`

	// Other data
	Directory string `json:"-"`
}

func parseTemplate(directory string) (Template, error) {
	templateJSON := filepath.Join(directory, "template.json")
	var result Template
	data, err := ioutil.ReadFile(templateJSON)
	if err != nil {
		return result, err
	}

	result.Directory = directory
	err = json.Unmarshal(data, &result)
	return result, err
}

// TemplateShortNames returns a slicer of short template names
func TemplateShortNames() (*slicer.StringSlicer, error) {

	var result slicer.StringSlicer

	// If the cache isn't loaded, load it
	if templateCache == nil {
		err := loadTemplateCache()
		if err != nil {
			return nil, err
		}
	}

	for _, template := range templateCache {
		result.Add(template.ShortName)
	}

	return &result, nil
}

// List returns the list of available templates
func List() ([]Template, error) {

	// If the cache isn't loaded, load it
	if templateCache == nil {
		err := loadTemplateCache()
		if err != nil {
			return nil, err
		}
	}

	return templateCache, nil
}

// getTemplateByShortname returns the template with the given short name
func getTemplateByShortname(shortname string) (Template, error) {

	var result Template

	// If the cache isn't loaded, load it
	if templateCache == nil {
		err := loadTemplateCache()
		if err != nil {
			return result, err
		}
	}

	for _, template := range templateCache {
		if template.ShortName == shortname {
			return template, nil
		}
	}

	return result, fmt.Errorf("shortname '%s' is not a valid template shortname", shortname)
}

// Loads the template cache
func loadTemplateCache() error {

	// Get local template directory
	templateDir := fs.RelativePath("templates")

	// Get directories
	files, err := ioutil.ReadDir(templateDir)
	if err != nil {
		return err
	}

	// Reset cache
	templateCache = []Template{}

	for _, file := range files {
		if file.IsDir() {
			templateDir := filepath.Join(templateDir, file.Name())
			template, err := parseTemplate(templateDir)
			if err != nil {
				// Cannot parse this template, continue
				continue
			}
			templateCache = append(templateCache, template)
		}
	}

	return nil
}

// Install the given template
func Install(options *Options) error {

	// Get cwd
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Did the user want to install in current directory?
	if options.TargetDir == "." {
		// Yes - use cwd
		options.TargetDir = cwd
	} else {
		// Get the absolute path of the given directory
		targetDir, err := filepath.Abs(filepath.Join(cwd, options.TargetDir))
		if err != nil {
			return err
		}
		options.TargetDir = targetDir
		if !fs.DirExists(options.TargetDir) {
			err := fs.Mkdir(options.TargetDir)
			if err != nil {
				return err
			}
		}
	}

	// Get template
	template, err := getTemplateByShortname(options.TemplateName)
	if err != nil {
		return err
	}

	// Use Gosod to install the template
	installer, err := gosod.TemplateDir(template.Directory)
	if err != nil {
		return err
	}

	// Ignore template.json files
	installer.IgnoreFilename("template.json")

	// Setup the data.
	// We use the directory name for the binary name, like Go
	BinaryName := filepath.Base(options.TargetDir)
	NPMProjectName := strings.ToLower(strings.ReplaceAll(BinaryName, " ", ""))
	localWailsDirectory := fs.RelativePath("../..")
	templateData := &Data{
		ProjectName:    options.ProjectName,
		BinaryName:     filepath.Base(options.TargetDir),
		NPMProjectName: NPMProjectName,
		WailsDirectory: localWailsDirectory,
	}

	// Extract the template
	err = installer.Extract(options.TargetDir, templateData)
	if err != nil {
		return err
	}

	// Calculate the directory name
	return nil
}

// OutputList prints the list of available tempaltes to the given logger
func OutputList(logger *clilogger.CLILogger) error {
	templates, err := List()
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(logger.Writer)
	table.SetHeader([]string{"Template", "Short Name", "Description"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	for _, template := range templates {
		table.Append([]string{template.Name, template.ShortName, template.Description})
	}
	table.Render()
	return nil
}
