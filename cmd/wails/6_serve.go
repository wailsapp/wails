package main

import (
	"fmt"
	"runtime"

	"github.com/leaanthony/spinner"
	"github.com/wailsapp/wails/cmd"
)

func init() {

	var forceRebuild = false
	var verbose = false
	buildSpinner := spinner.NewSpinner()
	buildSpinner.SetSpinSpeed(50)

	commandDescription := `This command builds then serves your application in bridge mode. Useful for developing your app in a browser.`
	initCmd := app.Command("serve", "Run your Wails project in bridge mode").
		LongDescription(commandDescription).
		BoolFlag("verbose", "Verbose output", &verbose).
		BoolFlag("f", "Force rebuild of application components", &forceRebuild)

	initCmd.Action(func() error {

		message := "Serving Application"
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

		// Set project options
		projectOptions.Verbose = verbose
		projectOptions.Platform = runtime.GOOS

		// Save project directory
		projectDir := fs.Cwd()

		// Install the bridge library
		err = cmd.InstallBridge(projectDir, projectOptions)
		if err != nil {
			return err
		}

		// Install dependencies
		err = cmd.InstallGoDependencies(projectOptions.Verbose)
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
