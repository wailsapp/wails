package commands

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/templates"

	"github.com/pterm/pterm"
)

var DisableFooter bool

func Init(options *flags.Init) error {

	if options.List {
		return printTemplates()
	}

	if options.Quiet {
		pterm.DisableOutput()
	}

	if options.ProjectName == "" {
		return fmt.Errorf("please use the -n flag to specify a project name")
	}

	options.ProjectName = sanitizeFileName(options.ProjectName)

	if !templates.ValidTemplateName(options.TemplateName) {
		return fmt.Errorf("invalid template name: %s. Use -l flag to list valid templates", options.TemplateName)
	}

	err := templates.Install(options)
	if err != nil {
		return err
	}

	// Generate build assets
	buildAssetsOptions := &BuildAssetsOptions{
		Name:               options.ProjectName,
		Dir:                filepath.Join(options.ProjectDir, "build"),
		Silent:             true,
		ProductCompany:     options.ProductCompany,
		ProductName:        options.ProductName,
		ProductDescription: options.ProductDescription,
		ProductVersion:     options.ProductVersion,
		ProductIdentifier:  options.ProductIdentifier,
		ProductCopyright:   options.ProductCopyright,
		ProductComments:    options.ProductComments,
	}
	return GenerateBuildAssets(buildAssetsOptions)
}

func printTemplates() error {
	defaultTemplates := templates.GetDefaultTemplates()

	pterm.DefaultSection.Println("Available templates")

	table := pterm.TableData{{"Name", "Description"}}
	for _, template := range defaultTemplates {
		table = append(table, []string{template.Name, template.Description})
	}
	err := pterm.DefaultTable.WithHasHeader(true).WithBoxed(true).WithData(table).Render()
	pterm.Println()
	return err
}

func sanitizeFileName(fileName string) string {
	// Regular expression to match non-allowed characters in file names
	// You can adjust this based on the specific requirements of your file system
	reg := regexp.MustCompile(`[^a-zA-Z0-9_.-]`)

	// Replace matched characters with an underscore or any other safe character
	return reg.ReplaceAllString(fileName, "_")
}
