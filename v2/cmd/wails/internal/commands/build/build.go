package build

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/leaanthony/clir"
	"github.com/leaanthony/slicer"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
	"github.com/wailsapp/wails/v2/pkg/commands/build"
)

// AddBuildSubcommand adds the `build` command for the Wails application
func AddBuildSubcommand(app *clir.Cli, w io.Writer) {

	outputType := "desktop"

	validTargetTypes := slicer.String([]string{"desktop", "hybrid", "server"})

	command := app.NewSubCommand("build", "Builds the application")

	// Setup target type flag
	//description := "Type of application to build. Valid types: " + validTargetTypes.Join(",")
	//command.StringFlag("t", description, &outputType)

	// Setup production flag
	production := false
	command.BoolFlag("production", "Build in production mode", &production)

	// Setup pack flag
	pack := false
	command.BoolFlag("package", "Create a platform specific package", &pack)

	compilerCommand := "go"
	command.StringFlag("compiler", "Use a different go compiler to build, eg go1.15beta1", &compilerCommand)

	// Setup Platform flag
	platform := runtime.GOOS
	command.StringFlag("platform", "Platform to target", &platform)

	// Verbosity
	verbosity := 1
	command.IntFlag("v", "Verbosity level (0 - silent, 1 - default, 2 - verbose)", &verbosity)

	// ldflags to pass to `go`
	ldflags := ""
	command.StringFlag("ldflags", "optional ldflags", &ldflags)

	// Log to file
	//logFile := ""
	//command.StringFlag("l", "Log to file", &logFile)

	// Retain assets
	keepAssets := false
	command.BoolFlag("k", "Keep generated assets", &keepAssets)

	// Retain assets
	outputFilename := ""
	command.StringFlag("o", "Output filename", &outputFilename)

	appleIdentity := ""
	if runtime.GOOS == "darwin" {
		command.StringFlag("sign", "Signs your app with the given identity.", &appleIdentity)
	}

	command.Action(func() error {

		quiet := verbosity == 0

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

		// Ensure package is used with apple identity
		if appleIdentity != "" && pack == false {
			return fmt.Errorf("must use `-package` flag when using `-sign`")
		}

		// Setup mode
		mode := build.Debug
		if production {
			mode = build.Production
		}

		// Check platform
		validPlatformArch := slicer.String([]string{
			"darwin",
			"darwin/amd64",
			"darwin/arm64",
			"darwin/universal",
			//"linux/amd64",
			//"linux/arm-7",
			//"windows/amd64",
		})
		if !validPlatformArch.Contains(platform) {
			return fmt.Errorf("platform %s is not supported", platform)
		}

		// Create BuildOptions
		buildOptions := &build.Options{
			Logger:        logger,
			OutputType:    outputType,
			OutputFile:    outputFilename,
			Mode:          mode,
			Pack:          pack,
			LDFlags:       ldflags,
			Compiler:      compilerCommand,
			KeepAssets:    keepAssets,
			AppleIdentity: appleIdentity,
			Verbosity:     verbosity,
		}

		// Calculate platform and arch
		platformSplit := strings.Split(platform, "/")
		buildOptions.Platform = platformSplit[0]
		buildOptions.Arch = runtime.GOARCH
		if len(platformSplit) == 2 {
			buildOptions.Arch = platformSplit[1]
		}

		// Start a new tabwriter
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 8, 8, 0, '\t', 0)

		buildModeText := "debug"
		if production {
			buildModeText = "production"
		}

		// Write out the system information
		fmt.Fprintf(w, "App Type: \t%s\n", buildOptions.OutputType)
		fmt.Fprintf(w, "Platform: \t%s\n", buildOptions.Platform)
		fmt.Fprintf(w, "Arch: \t%s\n", buildOptions.Arch)
		fmt.Fprintf(w, "Compiler: \t%s\n", buildOptions.Compiler)
		fmt.Fprintf(w, "Build Mode: \t%s\n", buildModeText)
		fmt.Fprintf(w, "Package: \t%t\n", buildOptions.Pack)
		fmt.Fprintf(w, "KeepAssets: \t%t\n", buildOptions.KeepAssets)
		fmt.Fprintf(w, "LDFlags: \t\"%s\"\n", buildOptions.LDFlags)
		if len(buildOptions.OutputFile) > 0 {
			fmt.Fprintf(w, "Output File: \t%s\n", buildOptions.OutputFile)
		}
		fmt.Fprintf(w, "\n")
		w.Flush()

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
	buildOptions.Logger.Println("")
	buildOptions.Logger.Println(fmt.Sprintf("Built '%s' in %s.", outputFilename, elapsed.Round(time.Millisecond).String()))
	buildOptions.Logger.Println("")

	return nil
}
