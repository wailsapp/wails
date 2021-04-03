package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/leaanthony/spinner"
	"github.com/wailsapp/wails/cmd"
)

// getSupportedPlatforms returns a slice of platform/architecture
// targets that are buildable using the cross-platform 'x' option.
func getSupportedPlatforms() []string {
	return []string{
		"darwin/amd64",
		"linux/amd64",
		"linux/arm-7",
		"windows/amd64",
	}
}

func init() {

	var packageApp = false
	var forceRebuild = false
	var debugMode = false
	var usefirebug = false
	var gopath = ""
	var typescriptFilename = ""
	var verbose = false
	var platform = ""
	var ldflags = ""
	var tags = ""

	buildSpinner := spinner.NewSpinner()
	buildSpinner.SetSpinSpeed(50)

	commandDescription := `This command will check to ensure all pre-requistes are installed prior to building. If not, it will attempt to install them. Building comprises of a number of steps: install frontend dependencies, build frontend, pack frontend, compile main application.`
	initCmd := app.Command("build", "Builds your Wails project").
		LongDescription(commandDescription).
		BoolFlag("p", "Package application on successful build", &packageApp).
		BoolFlag("f", "Force rebuild of application components", &forceRebuild).
		BoolFlag("d", "Build in Debug mode", &debugMode).
		BoolFlag("firebug", "Enable firebug console for debug builds", &usefirebug).
		BoolFlag("verbose", "Verbose output", &verbose).
		StringFlag("t", "Generate Typescript definitions to given file (at runtime)", &typescriptFilename).
		StringFlag("ldflags", "Extra options for -ldflags", &ldflags).
		StringFlag("gopath", "Specify your GOPATH location. Mounted to /go during cross-compilation.", &gopath).
		StringFlag("tags", "Build tags to pass to the go compiler (quoted and space separated)", &tags)

	var b strings.Builder
	for _, plat := range getSupportedPlatforms() {
		_, err := fmt.Fprintf(&b, " - %s\n", plat)
		if err != nil {
			log.Fatal(err)
		}
	}
	initCmd.StringFlag("x",
		fmt.Sprintf("Cross-compile application to specified platform via xgo\n%s", b.String()),
		&platform)

	initCmd.Action(func() error {

		message := "Building Application"
		if packageApp {
			message = "Packaging Application"
		}
		if forceRebuild {
			message += " (force rebuild)"
		}
		logger.PrintSmallBanner(message)
		fmt.Println()

		// Project options
		projectOptions := &cmd.ProjectOptions{}
		projectOptions.Verbose = verbose
		projectOptions.UseFirebug = usefirebug

		// Check we are in project directory
		// Check project.json loads correctly
		fs := cmd.NewFSHelper()
		err := projectOptions.LoadConfig(fs.Cwd())
		if err != nil {
			return fmt.Errorf("unable to find 'project.json'. Please check you are in a Wails project directory")
		}

		// Set firebug flag
		projectOptions.UseFirebug = usefirebug

		// Check that this platform is supported
		if !projectOptions.PlatformSupported() {
			logger.Yellow("WARNING: This project is unsupported on %s - it probably won't work!\n         Valid platforms: %s\n", runtime.GOOS, strings.Join(projectOptions.Platforms, ", "))
		}

		// Set cross-compile
		projectOptions.Platform = runtime.GOOS
		if len(platform) > 0 {
			supported := false
			for _, plat := range getSupportedPlatforms() {
				if plat == platform {
					supported = true
				}
			}
			if !supported {
				return fmt.Errorf("unsupported platform '%s' specified.\nPlease run `wails build -h` to see the supported platform/architecture options", platform)
			}

			projectOptions.CrossCompile = true
			plat := strings.Split(platform, "/")
			projectOptions.Platform = plat[0]
			projectOptions.Architecture = plat[1]
		}

		// Add ldflags
		projectOptions.LdFlags = ldflags
		projectOptions.GoPath = gopath

		// Add tags
		projectOptions.Tags = tags

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
		err = cmd.InstallGoDependencies(projectOptions.Verbose)
		if err != nil {
			return err
		}

		// Build application
		buildMode := cmd.BuildModeProd
		if debugMode {
			buildMode = cmd.BuildModeDebug
		}

		// Save if we wish to dump typescript or not
		if typescriptFilename != "" {
			projectOptions.SetTypescriptDefsFilename(typescriptFilename)
		}

		// Update go.mod if it is out of sync with current version
		outofsync, err := cmd.GoModOutOfSync()
		if err != nil {
			return err
		}
		gomodVersion, err := cmd.GetWailsVersion()
		if err != nil {
			return err
		}
		if outofsync {
			syncMessage := fmt.Sprintf("Updating go.mod (Wails version %s => %s)", gomodVersion, cmd.Version)
			buildSpinner := spinner.NewSpinner(syncMessage)
			buildSpinner.Start()
			err := cmd.UpdateGoModVersion()
			if err != nil {
				buildSpinner.Error(err.Error())
				return err
			}
			buildSpinner.Success()
		}

		err = cmd.BuildApplication(projectOptions.BinaryName, forceRebuild, buildMode, packageApp, projectOptions)
		if err != nil {
			return err
		}

		if projectOptions.Platform == "windows" {
			logger.Yellow("*** Please note: Windows builds use mshtml which is only compatible with IE11. For more information, please read https://wails.app/guides/windows/ ***")
		}

		logger.Yellow("Awesome! Project '%s' built!", projectOptions.Name)

		return nil

	})
}
