package initialise

import (
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/buildassets"
	"io"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/cmd/wails/internal/commands/initialise/templates"

	"github.com/leaanthony/clir"
	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
	"github.com/wailsapp/wails/v2/pkg/git"
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
	command.BoolFlag("q", "Suppress output to console", &quiet)

	initGit := false
	gitInstalled := git.IsInstalled()
	if gitInstalled {
		// Git Init
		command.BoolFlag("g", "Initialise git repository", &initGit)
	}

	// VSCode project files
	vscode := false
	command.BoolFlag("vscode", "Generate VSCode project files", &vscode)

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
			ProjectName:    projectName,
			TargetDir:      projectDirectory,
			TemplateName:   templateName,
			Logger:         logger,
			GenerateVSCode: vscode,
			InitGit:        initGit,
		}

		// Try to discover author details from git config
		err := findAuthorDetails(options)
		if err != nil {
			return err
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

	// Install the default assets
	err = buildassets.Install(options.TargetDir, options.ProjectName)
	if err != nil {
		return err
	}

	if options.InitGit {
		err = initGit(options)
		if err != nil {
			return err
		}
	}

	// Output stats
	elapsed := time.Since(start)
	options.Logger.Println("")
	options.Logger.Println("Project Name:      " + options.ProjectName)
	options.Logger.Println("Project Directory: " + options.TargetDir)
	options.Logger.Println("Project Template:  " + options.TemplateName)
	if options.GenerateVSCode {
		options.Logger.Println("VSCode config files generated.")
	}
	if options.InitGit {
		options.Logger.Println("Git repository initialised.")
	}
	options.Logger.Println("")
	options.Logger.Println(fmt.Sprintf("Initialised project '%s' in %s.", options.ProjectName, elapsed.Round(time.Millisecond).String()))
	options.Logger.Println("")

	return nil
}

func initGit(options *templates.Options) error {
	err := git.InitRepo(options.TargetDir)
	if err != nil {
		return errors.Wrap(err, "Unable to initialise git repository:")
	}

	return nil
}

func findAuthorDetails(options *templates.Options) error {
	if git.IsInstalled() {
		name, err := git.Name()
		if err != nil {
			return err
		}
		options.AuthorName = strings.TrimSpace(name)

		email, err := git.Email()
		if err != nil {
			return err
		}
		options.AuthorEmail = strings.TrimSpace(email)
	}

	return nil
}
