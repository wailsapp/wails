package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/leaanthony/spinner"
	"github.com/wailsapp/wails/cmd"
)

func init() {

	var bundle = false
	var forceRebuild = false
	buildSpinner := spinner.NewSpinner()
	buildSpinner.SetSpinSpeed(50)

	commandDescription := `This command will check to ensure all pre-requistes are installed prior to building. If not, it will attempt to install them. Building comprises of a number of steps: install frontend dependencies, build frontend, pack frontend, compile main application.`
	initCmd := app.Command("build", "Builds your Wails project").
		LongDescription(commandDescription).
		BoolFlag("b", "Bundle application on successful build", &bundle).
		BoolFlag("f", "Force rebuild of application components", &forceRebuild)

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
		if projectOptions.FrontEnd != nil {
			if projectOptions.FrontEnd.Dir == "" {
				return fmt.Errorf("Frontend directory not set in project.json")
			}
			if projectOptions.FrontEnd.Build == "" {
				return fmt.Errorf("Frontend build command not set in project.json")
			}
			if projectOptions.FrontEnd.Install == "" {
				return fmt.Errorf("Frontend install command not set in project.json")
			}
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

		// Install backend deps - needed?
		if projectOptions.FrontEnd != nil {
			// Install frontend deps
			err = os.Chdir(projectOptions.FrontEnd.Dir)
			if err != nil {
				return err
			}

			// Check if frontend deps have been updated
			buildSpinner.Start("Installing frontend dependencies...")

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
						buildSpinner.Success("Skipped frontend dependencies (-f to force rebuild)")
					}
				}
			}

			// Md5 sum package.json
			// Different? Build
			if requiresNPMInstall || forceRebuild {
				// Install dependencies
				err = program.RunCommand(projectOptions.FrontEnd.Install)
				if err != nil {
					buildSpinner.Error()
					return err
				}
				buildSpinner.Success()

				// Update md5sum file
				ioutil.WriteFile(md5sumFile, []byte(packageJSONMD5), 0644)
			}

			// Build frontend
			buildSpinner.Start("Building frontend...")
			err = program.RunCommand(projectOptions.FrontEnd.Build)
			if err != nil {
				buildSpinner.Error()
				return err
			}
			buildSpinner.Success()
		}

		// Run packr in project directory
		err = os.Chdir(projectDir)
		if err != nil {
			buildSpinner.Error()
			return err
		}

		// Support build tags
		buildTags := []string{}

		// Do we have any frameworks specified?
		if projectOptions.Framework != nil {
			buildSpinner.Start()
			buildSpinner.Success("Compiling support for " + projectOptions.Framework.Name)
			buildTags = append(buildTags, projectOptions.Framework.BuildTag)
		}

		// // Initialise Go Module - if go.mod doesn't exist
		// if !fs.FileExists("go.mod") {
		// 	buildSpinner.Start("Initialising Go module...")
		// 	err = program.RunCommand("go mod init " + projectOptions.BinaryName)
		// 	if err != nil {
		// 		buildSpinner.Error()
		// 		return err
		// 	}
		// 	buildSpinner.Success()
		// }

		buildSpinner.Start("Installing Dependencies...")
		installCommand := "go get"
		err = program.RunCommand(installCommand)
		if err != nil {
			buildSpinner.Error()
			return err
		}
		buildSpinner.Success()

		buildSpinner.Start("Packing + Compiling project...")

		buildCommand := "packr build"

		// Add build tags
		if len(buildTags) > 0 {
			buildCommand += fmt.Sprintf(" --tags '%s'", strings.Join(buildTags, " "))
		}

		if projectOptions.BinaryName != "" {
			buildCommand += " -o " + projectOptions.BinaryName
		}

		// If we are forcing a rebuild
		if forceRebuild {
			buildCommand += " -a"
		}

		err = program.RunCommand(buildCommand)
		if err != nil {
			buildSpinner.Error()
			return err
		}
		buildSpinner.Success()

		if bundle == false {
			logger.Yellow("Awesome! Project '%s' built!", projectOptions.Name)
			return nil
		}

		// Bundle app
		buildSpinner.Start("Bundling Application")
		bundler := cmd.NewBundleHelper()
		err = bundler.Bundle(projectOptions)
		if err != nil {
			buildSpinner.Error()
			return err
		}
		buildSpinner.Success()
		logger.Yellow("Awesome! Project '%s' built!", projectOptions.Name)
		return nil
	})
}
