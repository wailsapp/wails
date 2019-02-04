package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"

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
		log := cmd.NewLogger()
		message := "Building Application"
		if forceRebuild {
			message += " (force rebuild)"
		}
		log.WhiteUnderline(message)

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

		// Check pre-requisites are installed

		// Program checker
		program := cmd.NewProgramHelper()

		if projectOptions.FrontEnd != nil {
			// npm
			if !program.IsInstalled("npm") {
				return fmt.Errorf("it appears npm is not installed. Please install and run again")
			}
		}

		// packr
		if !program.IsInstalled("packr") {
			buildSpinner.Start("Installing packr...")
			err := program.InstallGoPackage("github.com/gobuffalo/packr/...")
			if err != nil {
				buildSpinner.Error()
				return err
			}
			buildSpinner.Success()
		}

		// Save project directory
		projectDir := fs.Cwd()

		// Install deps
		if projectOptions.FrontEnd != nil {

			// Install frontend deps
			err = os.Chdir(projectOptions.FrontEnd.Dir)
			if err != nil {
				return err
			}

			// Check if frontend deps have been updated
			feSpinner := spinner.New("Installing frontend dependencies (This may take a while)...")
			feSpinner.SetSpinSpeed(50)
			feSpinner.Start()

			requiresNPMInstall := true

			// Read in package.json MD5
			packageJSONMD5, err := fs.FileMD5("package.json")
			if err != nil {
				return err
			}

			const md5sumFile = "package.json.md5"

			// If we aren't forcing the install and the md5sum file exists
			if !forceRebuild && fs.FileExists(md5sumFile) {
				// Yes - read contents
				savedMD5sum, err := fs.LoadAsString(md5sumFile)
				// File exists
				if err == nil {
					// Compare md5
					if savedMD5sum == packageJSONMD5 {
						// Same - no need for reinstall
						requiresNPMInstall = false
						feSpinner.Success("Skipped frontend dependencies (-f to force rebuild)")
					}
				}
			}

			// Md5 sum package.json
			// Different? Build
			if requiresNPMInstall || forceRebuild {
				// Install dependencies
				err = program.RunCommand(projectOptions.FrontEnd.Install)
				if err != nil {
					feSpinner.Error()
					return err
				}
				feSpinner.Success()

				// Update md5sum file
				ioutil.WriteFile(md5sumFile, []byte(packageJSONMD5), 0644)
			}

			bridgeFile := "wailsbridge.prod.js"

			// Copy bridge to project
			_, filename, _, _ := runtime.Caller(1)
			bridgeFileSource := filepath.Join(path.Dir(filename), "..", "assets", "default", bridgeFile)
			bridgeFileTarget := filepath.Join(projectDir, projectOptions.FrontEnd.Dir, projectOptions.FrontEnd.Bridge, "wailsbridge.js")
			err = fs.CopyFile(bridgeFileSource, bridgeFileTarget)
			if err != nil {
				return err
			}

			// Build frontend
			err = cmd.BuildFrontend(projectOptions.FrontEnd.Build)
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
		buildMode := "prod"
		if debugMode {
			buildMode = "debug"
		}
		err = cmd.BuildApplication(projectOptions.BinaryName, forceRebuild, buildMode)
		if err != nil {
			return err
		}

		// Package application
		if packageApp {
			err = cmd.PackageApplication(projectOptions)
			if err != nil {
				return err
			}
		}

		logger.Yellow("Awesome! Project '%s' built!", projectOptions.Name)

		return nil

	})
}
