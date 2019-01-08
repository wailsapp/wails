package main

import (
	"fmt"

	"github.com/wailsapp/wails/cmd"
)

func init() {

	projectHelper := cmd.NewProjectHelper()
	projectOptions := projectHelper.NewProjectOptions()
	commandDescription := `Generates a new Wails project using the given flags. 
Any flags that are required and not given will be prompted for.`

	initCommand := app.Command("init", "Initialises a new Wails project").
		LongDescription(commandDescription).
		BoolFlag("f", "Use defaults", &projectOptions.UseDefaults).
		StringFlag("dir", "Directory to create project in", &projectOptions.OutputDirectory).
		StringFlag("template", "Template name", &projectOptions.Template).
		StringFlag("name", "Project name", &projectOptions.Name).
		StringFlag("description", "Project description", &projectOptions.Description).
		StringFlag("output", "Output binary name", &projectOptions.BinaryName)

	initCommand.Action(func() error {

		logger.WhiteUnderline("Initialising project")
		fmt.Println()

		// Check if the system is initialised
		system := cmd.NewSystemHelper()
		err := system.CheckInitialised()
		if err != nil {
			return err
		}

		success, err := cmd.CheckDependenciesSilent(logger)
		if !success {
			return err
		}

		// Do we want to just force defaults?
		if projectOptions.UseDefaults {
			// Use defaults
			projectOptions.Defaults()
		} else {
			err = projectOptions.PromptForInputs()
			if err != nil {
				return err
			}
		}

		// Generate the project
		err = projectHelper.GenerateProject(projectOptions)
		if err != nil {
			logger.Error(err.Error())
		}
		return err
	})
}
