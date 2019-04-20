package main

import (
	"fmt"

	"github.com/leaanthony/spinner"
	"github.com/wailsapp/wails/cmd"
)

func init() {

	var forceRebuild = false
	buildSpinner := spinner.NewSpinner()
	buildSpinner.SetSpinSpeed(50)

	commandDescription := `This command builds then serves your application in bridge mode. Useful for developing your app in a browser.`
	initCmd := app.Command("serve", "Run your Wails project in bridge mode.").
		LongDescription(commandDescription).
		BoolFlag("f", "Force rebuild of application components", &forceRebuild)

	initCmd.Action(func() error {

		message := "Serving Application"
		logger.PrintSmallBanner(message)
		fmt.Println()

		// Check Mewn is installed
		err := cmd.CheckMewn()
		if err != nil {
			return err
		}

		// Project options
		projectOptions := &cmd.ProjectOptions{}

		// Check we are in project directory
		// Check project.json loads correctly
		fs := cmd.NewFSHelper()
		err = projectOptions.LoadConfig(fs.Cwd())
		if err != nil {
			return err
		}

		// Install dependencies
		err = cmd.InstallGoDependencies()
		if err != nil {
			return err
		}

		buildMode := cmd.BuildModeBridge
		err = cmd.BuildApplication(projectOptions.BinaryName, forceRebuild, buildMode, false, projectOptions)
		if err != nil {
			return err
		}

		logger.Yellow("Awesome! Project '%s' built!", projectOptions.Name)
		return cmd.ServeProject(projectOptions, logger)
	})
}
