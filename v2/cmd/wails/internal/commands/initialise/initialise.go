package initialise

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/internal/templates"
)

// AddSubcommand adds the `init` command for the Wails application
func AddSubcommand(app *clir.Cli) error {

	// Load the template shortnames
	validShortNames, err := templates.TemplateShortNames()
	if err != nil {
		return err
	}

	command := app.NewSubCommand("init", "Initialise a new Wails project")

	// Setup template name flag
	templateName := "vanilla"
	description := "Name of template to use. Valid tempates: " + validShortNames.Join(" ")
	command.StringFlag("t", description, &templateName)

	// Setup project name
	projectName := ""
	command.StringFlag("n", "Name of project", &projectName)

	// Setup project directory
	projectDirectory := "."
	command.StringFlag("d", "Project directory", &projectDirectory)

	// Quiet Init
	quiet := false
	command.BoolFlag("q", "Supress output to console", &quiet)

	// List templates
	list := false
	command.BoolFlag("l", "List templates", &list)

	command.Action(func() error {

		// Create logger
		logger := logger.New()

		if !quiet {
			logger.AddOutput(os.Stdout)
		}

		// Are we listing templates?
		if list {
			app.PrintBanner()
			err := templates.OutputList(logger)
			logger.Writeln("")
			return err
		}

		// Validate output type
		if !validShortNames.Contains(templateName) {
			logger.Write(fmt.Sprintf("ERROR: Template '%s' is not valid", templateName))
			logger.Writeln("")
			command.PrintHelp()
			return nil
		}

		// Validate name
		if len(projectName) == 0 {
			logger.Writeln("ERROR: Project name required")
			logger.Writeln("")
			command.PrintHelp()
			return nil
		}

		if !quiet {
			app.PrintBanner()
		}

		task := fmt.Sprintf("Initialising Project %s", strings.Title(projectName))
		logger.Writeln(task)
		logger.Writeln(strings.Repeat("-", len(task)))

		// Create Template Options
		options := &templates.Options{
			ProjectName:  projectName,
			TargetDir:    projectDirectory,
			TemplateName: templateName,
			Logger:       logger,
		}

		return initProject(options)
	})

	return nil
}

// initProject is our main init command
func initProject(options *templates.Options) error {

	// Start Time
	start := time.Now()

	// Install the template
	err := templates.Install(options)
	if err != nil {
		return err
	}

	// Output stats
	elapsed := time.Since(start)
	options.Logger.Writeln("")
	options.Logger.Writeln(fmt.Sprintf("Initialised project '%s' in %s.", options.ProjectName, elapsed.Round(time.Millisecond).String()))
	options.Logger.Writeln("")

	return nil
}
