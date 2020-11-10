package initialise

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/internal/templates"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
)

// AddSubcommand adds the `init` command for the Wails application
func AddSubcommand(app *clir.Cli, w io.Writer) error {

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
	projectDirectory := ""
	command.StringFlag("d", "Project directory", &projectDirectory)

	// Quiet Init
	quiet := false
	command.BoolFlag("q", "Supress output to console", &quiet)

	// List templates
	list := false
	command.BoolFlag("l", "List templates", &list)

	command.Action(func() error {

		// Create logger
		logger := clilogger.New(w)
		logger.Mute(quiet)

		// Are we listing templates?
		if list {
			app.PrintBanner()
			err := templates.OutputList(logger)
			logger.Println("")
			return err
		}

		// Validate output type
		if !validShortNames.Contains(templateName) {
			logger.Print(fmt.Sprintf("[ERROR] Template '%s' is not valid", templateName))
			logger.Println("")
			command.PrintHelp()
			return nil
		}

		// Validate name
		if len(projectName) == 0 {
			logger.Println("ERROR: Project name required")
			logger.Println("")
			command.PrintHelp()
			return nil
		}

		if !quiet {
			app.PrintBanner()
		}

		task := fmt.Sprintf("Initialising Project %s", strings.Title(projectName))
		logger.Println(task)
		logger.Println(strings.Repeat("-", len(task)))

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
	options.Logger.Println("")
	options.Logger.Println("Project Name:      " + options.ProjectName)
	options.Logger.Println("Project Directory: " + options.TargetDir)
	options.Logger.Println("Project Template:  " + options.TemplateName)
	options.Logger.Println("")
	options.Logger.Println(fmt.Sprintf("Initialised project '%s' in %s.", options.ProjectName, elapsed.Round(time.Millisecond).String()))
	options.Logger.Println("")

	return nil
}
