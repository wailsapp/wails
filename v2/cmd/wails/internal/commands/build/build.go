package build

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/leaanthony/clir"
	"github.com/leaanthony/slicer"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/commands/build"
)

// AddBuildSubcommand adds the `build` command for the Wails application
func AddBuildSubcommand(app *clir.Cli) {

	outputType := "desktop"

	validTargetTypes := slicer.String([]string{"desktop", "hybrid", "server"})

	command := app.NewSubCommand("build", "Builds the application")

	// Setup target type flag
	description := "Type of application to build. Valid types: " + validTargetTypes.Join(",")
	command.StringFlag("t", description, &outputType)

	// Setup production flag
	production := false
	command.BoolFlag("production", "Build in production mode", &production)

	// Setup pack flag
	pack := false
	command.BoolFlag("pack", "Create a platform specific package", &pack)

	compilerCommand := "go"
	command.StringFlag("compiler", "Use a different go compiler to build, eg go1.15beta1", &compilerCommand)

	// Setup Platform flag
	platform := runtime.GOOS
	command.StringFlag("platform", "Platform to target", &platform)

	// Quiet Build
	quiet := false
	command.BoolFlag("q", "Supress output to console", &quiet)

	// ldflags to pass to `go`
	ldflags := ""
	command.StringFlag("ldflags", "optional ldflags", &ldflags)

	// Log to file
	logFile := ""
	command.StringFlag("l", "Log to file", &logFile)

	command.Action(func() error {

		// Create logger
		logger := logger.New()

		if !quiet {
			logger.AddOutput(os.Stdout)
		}

		// Validate output type
		if !validTargetTypes.Contains(outputType) {
			return fmt.Errorf("output type '%s' is not valid", outputType)
		}

		if !quiet {
			app.PrintBanner()
		}

		task := fmt.Sprintf("Building %s Application", strings.Title(outputType))
		logger.Writeln(task)
		logger.Writeln(strings.Repeat("-", len(task)))

		// Setup mode
		mode := build.Debug
		if production {
			mode = build.Production
		}

		// Create BuildOptions
		buildOptions := &build.Options{
			Logger:     logger,
			OutputType: outputType,
			Mode:       mode,
			Pack:       pack,
			Platform:   platform,
			LDFlags:    ldflags,
			Compiler:   compilerCommand,
		}

		return doBuild(buildOptions)
	})
}

// doBuild is our main build command
func doBuild(buildOptions *build.Options) error {

	// Start Time
	start := time.Now()

	outputFilename, err := build.Build(buildOptions)
	if err != nil {
		return err
	}
	// Output stats
	elapsed := time.Since(start)
	buildOptions.Logger.Writeln("")
	buildOptions.Logger.Writeln(fmt.Sprintf("Built '%s' in %s.", outputFilename, elapsed.Round(time.Millisecond).String()))
	buildOptions.Logger.Writeln("")

	return nil
}
