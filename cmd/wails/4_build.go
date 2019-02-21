package main

import (
	"fmt"
	"os"

	"github.com/leaanthony/spinner"
	"github.com/wailsapp/wails/cmd"
)

func init() {

	var packageApp = false
	var forceRebuild = false
	var debugMode = false
	buildSpinner := spinner.NewSpinner()
	buildSpinner.SetSpinSpeed(50)

	commandDescription := `This command will check to ensure all pre-requistes are installed prior to building. If not, it will attempt to install them. Building comprises of a number of steps: install frontend dependencies, build frontend, pack frontend, compile main application.`
	initCmd := app.Command("build", "Builds your Wails project").
		LongDescription(commandDescription).
		BoolFlag("p", "Package application on successful build", &packageApp).
		BoolFlag("f", "Force rebuild of application components", &forceRebuild).
		BoolFlag("d", "Build in Debug mode", &debugMode)

	initCmd.Action(func() error {

		message := "Building Application"
		if forceRebuild {
			message += " (force rebuild)"
		}
		logger.PrintSmallBanner(message)
		fmt.Println()

		// Project options
		projectOptions := &cmd.ProjectOptions{}

		// Check we are in project directory
		// Check project.json loads correctly
		fs := cmd.NewFSHelper()
		err := projectOptions.LoadConfig(fs.Cwd())
		if err != nil {
			return err
		}

		// Validate config
		// Check if we have a frontend
		err = cmd.ValidateFrontendConfig(projectOptions)
		if err != nil {
			return err
		}

		// Program checker
		program := cmd.NewProgramHelper()

		if projectOptions.FrontEnd != nil {
			// npm
			if !program.IsInstalled("npm") {
				return fmt.Errorf("it appears npm is not installed. Please install and run again")
			}
		}

		// Save project directory
		projectDir := fs.Cwd()

		// Install deps
		if projectOptions.FrontEnd != nil {
			err = cmd.InstallFrontendDeps(projectDir, projectOptions, forceRebuild, "build")
			if err != nil {
				return err
			}
		}

		// Move to project directory
		err = os.Chdir(projectDir)
		if err != nil {
			return err
		}

		// Install dependencies
		err = cmd.InstallGoDependencies()
		if err != nil {
			return err
		}

		// Build application
		buildMode := cmd.BuildModeProd
		if debugMode {
			buildMode = cmd.BuildModeDebug
		}

		err = cmd.BuildApplication(projectOptions.BinaryName, forceRebuild, buildMode, packageApp, projectOptions)
		if err != nil {
			return err
		}

		logger.Yellow("Awesome! Project '%s' built!", projectOptions.Name)

		return nil

	})
}
