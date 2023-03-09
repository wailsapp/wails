package commands

import (
	"fmt"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/templates"

	"github.com/pterm/pterm"
)

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

	if !templates.ValidTemplateName(options.TemplateName) {
		return fmt.Errorf("invalid template name: %s. Use -l flag to view available templates", options.TemplateName)
	}

	return templates.Install(options)
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
