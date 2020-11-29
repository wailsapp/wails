package main

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/leaanthony/spinner"
	"github.com/wailsapp/wails/cmd"
)

func init() {

	fs := cmd.NewFSHelper()

	var forceRebuild = false
	var verbose = false
	var configPath = fs.Cwd()
	var outPath = filepath.Join(fs.Cwd(), "build")
	buildSpinner := spinner.NewSpinner()
	buildSpinner.SetSpinSpeed(50)

	commandDescription := `This command builds then serves your application in bridge mode. Useful for developing your app in a browser.`
	initCmd := app.Command("serve", "Run your Wails project in bridge mode").
		LongDescription(commandDescription).
		BoolFlag("verbose", "Verbose output", &verbose).
		BoolFlag("f", "Force rebuild of application components", &forceRebuild).
		StringFlag("c", "Specify location of project.json", &configPath).
		StringFlag("o", "Specify where the built executable should be placed", &outPath)

	initCmd.Action(func() error {

		message := "Serving Application"
		logger.PrintSmallBanner(message)
		fmt.Println()

		// Check Mewn is installed
		err := cmd.CheckMewn(verbose)
		if err != nil {
			return err
		}

		// Project options
		projectOptions := &cmd.ProjectOptions{}

		// Check if configPath exists
		if ok := fs.DirExists(configPath); !ok {
			return fmt.Errorf("Unable to find 'project.json'. Please make sure the specified path is valid")
		}

		err = projectOptions.LoadConfig(configPath)
		if err != nil {
			return err
		}

		// Set project options
		projectOptions.Verbose = verbose
		projectOptions.Platform = runtime.GOOS

		// Install the bridge library
		err = cmd.InstallBridge(projectOptions)
		if err != nil {
			return err
		}

		// Install dependencies
		err = cmd.InstallGoDependencies(projectOptions.Verbose, projectOptions.MainPackage)
		if err != nil {
			return err
		}

		buildMode := cmd.BuildModeBridge

		err = cmd.BuildApplication(projectOptions.BinaryName, outPath, forceRebuild, buildMode, false, projectOptions)
		if err != nil {
			return err
		}

		logger.Yellow("Awesome! Project '%s' built!", projectOptions.Name)

		return cmd.ServeProject(outPath, projectOptions, logger)
	})
}
