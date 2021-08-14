package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/leaanthony/spinner"
	"github.com/wailsapp/wails/cmd"
)

func init() {

	projectHelper := cmd.NewProjectHelper()
	projectOptions := projectHelper.NewProjectOptions()
	commandDescription := `Generates a new Wails project using the given flags. 
Any flags that are required and not given will be prompted for.`
	build := false

	initCommand := app.Command("init", "Initialises a new Wails project").
		LongDescription(commandDescription).
		BoolFlag("f", "Use defaults", &projectOptions.UseDefaults).
		StringFlag("dir", "Directory to create project in", &projectOptions.OutputDirectory).
		StringFlag("template", "Template name", &projectOptions.Template).
		StringFlag("name", "Project name", &projectOptions.Name).
		StringFlag("description", "Project description", &projectOptions.Description).
		StringFlag("output", "Output binary name", &projectOptions.BinaryName).
		BoolFlag("build", "Build project after generating", &build)

	initCommand.Action(func() error {

		logger.PrintSmallBanner("Initialising project")
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

		genSpinner := spinner.NewSpinner()
		genSpinner.SetSpinSpeed(50)
		genSpinner.Start("Generating project...")

		// Generate the project
		err = projectHelper.GenerateProject(projectOptions)
		if err != nil {
			genSpinner.Error()
			return err
		}
		genSpinner.Success()
		if !build {
			logger.Yellow("Project '%s' initialised. Run `wails build` to build it.", projectOptions.Name)
			return nil
		}

		// Build the project
		cwd, _ := os.Getwd()
		projectDir := filepath.Join(cwd, projectOptions.OutputDirectory)
		program := cmd.NewProgramHelper()
		buildSpinner := spinner.NewSpinner()
		buildSpinner.SetSpinSpeed(50)
		buildSpinner.Start("Building project (this may take a while)...")
		err = program.RunCommandArray([]string{"wails", "build"}, projectDir)
		if err != nil {
			buildSpinner.Error(err.Error())
			return err
		}
		buildSpinner.Success()
		logger.Yellow("Project '%s' built in directory '%s'!", projectOptions.Name, projectOptions.OutputDirectory)

		return err
	})
}
