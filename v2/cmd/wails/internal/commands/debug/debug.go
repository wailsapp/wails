package debug

import (
	"fmt"
	"github.com/wailsapp/wails/v2/internal/shell"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/leaanthony/clir"
	"github.com/leaanthony/slicer"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
	"github.com/wailsapp/wails/v2/pkg/commands/build"
)

// AddSubcommand adds the `debug` command for the Wails application
func AddSubcommand(app *clir.Cli, w io.Writer) error {

	outputType := "desktop"

	validTargetTypes := slicer.String([]string{"desktop", "hybrid", "server"})

	command := app.NewSubCommand("debug", "Builds the application then runs delve on the binary")

	// Setup target type flag
	description := "Type of application to build. Valid types: " + validTargetTypes.Join(",")
	command.StringFlag("t", description, &outputType)

	compilerCommand := "go"
	command.StringFlag("compiler", "Use a different go compiler to build, eg go1.15beta1", &compilerCommand)

	quiet := false
	command.BoolFlag("q", "Suppress output to console", &quiet)

	// ldflags to pass to `go`
	ldflags := ""
	command.StringFlag("ldflags", "optional ldflags", &ldflags)

	// Log to file
	logFile := ""
	command.StringFlag("l", "Log to file", &logFile)

	command.Action(func() error {

		// Create logger
		logger := clilogger.New(w)
		logger.Mute(quiet)

		// Validate output type
		if !validTargetTypes.Contains(outputType) {
			return fmt.Errorf("output type '%s' is not valid", outputType)
		}

		if !quiet {
			app.PrintBanner()
		}

		task := fmt.Sprintf("Building %s Application", strings.Title(outputType))
		logger.Println(task)
		logger.Println(strings.Repeat("-", len(task)))

		// Setup mode
		mode := build.Debug

		// Create BuildOptions
		buildOptions := &build.Options{
			Logger:     logger,
			OutputType: outputType,
			Mode:       mode,
			Pack:       false,
			Platform:   runtime.GOOS,
			LDFlags:    ldflags,
			Compiler:   compilerCommand,
			KeepAssets: false,
		}

		outputFilename, err := doDebugBuild(buildOptions)
		if err != nil {
			return err
		}

		// Check delve exists
		delveExists := shell.CommandExists("dlv")
		if !delveExists {
			return fmt.Errorf("cannot launch delve (Is it installed?)")
		}

		// Get cwd
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		// Launch delve
		println("Launching Delve on port 2345...")
		command := shell.CreateCommand(cwd, "dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", outputFilename)
		return command.Run()
	})

	return nil
}

// doDebugBuild is our main build command
func doDebugBuild(buildOptions *build.Options) (string, error) {

	// Start Time
	start := time.Now()

	outputFilename, err := build.Build(buildOptions)
	if err != nil {
		return "", err
	}

	// Output stats
	elapsed := time.Since(start)
	buildOptions.Logger.Println("")
	buildOptions.Logger.Println(fmt.Sprintf("Built '%s' in %s.", outputFilename, elapsed.Round(time.Millisecond).String()))
	buildOptions.Logger.Println("")

	return outputFilename, nil
}
