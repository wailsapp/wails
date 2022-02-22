package build

import (
	"fmt"
	"github.com/wailsapp/wails/v2/internal/colour"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/wailsapp/wails/v2/cmd/wails/internal"
	"github.com/wailsapp/wails/v2/internal/gomod"

	"github.com/wailsapp/wails/v2/internal/system"

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

	// Setup noPackage flag
	noPackage := false
	command.BoolFlag("noPackage", "Skips platform specific packaging", &noPackage)

	compilerCommand := "go"
	command.StringFlag("compiler", "Use a different go compiler to build, eg go1.15beta1", &compilerCommand)

	skipModTidy := false
	command.BoolFlag("m", "Skip mod tidy before compile", &skipModTidy)

	compress := false
	command.BoolFlag("upx", "Compress final binary with UPX (if installed)", &compress)

	compressFlags := ""
	command.StringFlag("upxflags", "Flags to pass to upx", &compressFlags)

	// Setup Platform flag
	platform := runtime.GOOS
	command.StringFlag("platform", "Platform to target", &platform)

	// Verbosity
	verbosity := 1
	command.IntFlag("v", "Verbosity level (0 - silent, 1 - default, 2 - verbose)", &verbosity)

	// ldflags to pass to `go`
	ldflags := ""
	command.StringFlag("ldflags", "optional ldflags", &ldflags)

	// tags to pass to `go`
	tags := ""
	command.StringFlag("tags", "tags to pass to Go compiler (quoted and space separated)", &tags)

	outputFilename := ""
	command.StringFlag("o", "Output filename", &outputFilename)

	// Clean build directory
	cleanBuildDirectory := false
	command.BoolFlag("clean", "Clean the build directory before building", &cleanBuildDirectory)

	webview2 := "download"
	command.StringFlag("webview2", "WebView2 installer strategy: download,embed,browser,error.", &webview2)

	skipFrontend := false
	command.BoolFlag("s", "Skips building the frontend", &skipFrontend)

	forceBuild := false
	command.BoolFlag("f", "Force build application", &forceBuild)

	updateGoMod := false
	command.BoolFlag("u", "Updates go.mod to use the same Wails version as the CLI", &updateGoMod)

	debug := false
	command.BoolFlag("debug", "Retains debug data in the compiled application", &debug)

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

		// Check platform
		validPlatformArch := slicer.String([]string{
			"darwin",
			"darwin/amd64",
			"darwin/arm64",
			"darwin/universal",
			"linux",
			//"linux/amd64",
			//"linux/arm-7",
			"windows",
			"windows/amd64",
			"windows/arm64",
		})
		if !validPlatformArch.Contains(platform) {
			return fmt.Errorf("platform %s is not supported. Platforms supported: %s", platform, validPlatformArch.Join(","))
		}

		if compress && platform == "darwin/universal" {
			logger.Println("Warning: compress flag unsupported for universal binaries. Ignoring.")
			compress = false
		}

		// Lookup compiler path
		compilerPath, err := exec.LookPath(compilerCommand)
		if err != nil {
			return fmt.Errorf("unable to find compiler: %s", compilerCommand)
		}

		// Tags
		userTags := []string{}
		for _, tag := range strings.Split(tags, " ") {
			thisTag := strings.TrimSpace(tag)
			if thisTag != "" {
				userTags = append(userTags, thisTag)
			}
		}

		// Webview2 installer strategy (download by default)
		wv2rtstrategy := ""
		webview2 = strings.ToLower(webview2)
		if webview2 != "" {
			validWV2Runtime := slicer.String([]string{"download", "embed", "browser", "error"})
			if !validWV2Runtime.Contains(webview2) {
				return fmt.Errorf("invalid option for flag 'webview2': %s", webview2)
			}
			// These are the build tags associated with the strategies
			switch webview2 {
			case "embed":
				wv2rtstrategy = "wv2runtime.embed"
			case "error":
				wv2rtstrategy = "wv2runtime.error"
			case "browser":
				wv2rtstrategy = "wv2runtime.browser"
			}
		}

		mode := build.Production
		modeString := "Production"
		if debug {
			mode = build.Debug
			modeString = "Debug"
		}

		// Create BuildOptions
		buildOptions := &build.Options{
			Logger:              logger,
			OutputType:          outputType,
			OutputFile:          outputFilename,
			CleanBuildDirectory: cleanBuildDirectory,
			Mode:                mode,
			Pack:                !noPackage,
			LDFlags:             ldflags,
			Compiler:            compilerCommand,
			SkipModTidy:         skipModTidy,
			Verbosity:           verbosity,
			ForceBuild:          forceBuild,
			IgnoreFrontend:      skipFrontend,
			Compress:            compress,
			CompressFlags:       compressFlags,
			UserTags:            userTags,
			WebView2Strategy:    wv2rtstrategy,
		}

		// Calculate platform and arch
		platformSplit := strings.Split(platform, "/")
		buildOptions.Platform = platformSplit[0]
		if system.IsAppleSilicon {
			buildOptions.Arch = "arm64"
		} else {
			buildOptions.Arch = runtime.GOARCH
		}
		if len(platformSplit) == 2 {
			buildOptions.Arch = platformSplit[1]
		}

		// Start a new tabwriter
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 8, 8, 0, '\t', 0)

		// Write out the system information
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "App Type: \t%s\n", buildOptions.OutputType)
		fmt.Fprintf(w, "Platform: \t%s\n", buildOptions.Platform)
		fmt.Fprintf(w, "Arch: \t%s\n", buildOptions.Arch)
		fmt.Fprintf(w, "Compiler: \t%s\n", compilerPath)
		fmt.Fprintf(w, "Build Mode: \t%s\n", modeString)
		fmt.Fprintf(w, "Skip Frontend: \t%t\n", skipFrontend)
		fmt.Fprintf(w, "Compress: \t%t\n", buildOptions.Compress)
		fmt.Fprintf(w, "Package: \t%t\n", buildOptions.Pack)
		fmt.Fprintf(w, "Clean Build Dir: \t%t\n", buildOptions.CleanBuildDirectory)
		fmt.Fprintf(w, "LDFlags: \t\"%s\"\n", buildOptions.LDFlags)
		fmt.Fprintf(w, "Tags: \t[%s]\n", strings.Join(buildOptions.UserTags, ","))
		if len(buildOptions.OutputFile) > 0 {
			fmt.Fprintf(w, "Output File: \t%s\n", buildOptions.OutputFile)
		}
		fmt.Fprintf(w, "\n")
		w.Flush()

		err = checkGoModVersion(logger, updateGoMod)
		if err != nil {
			return err
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
	buildOptions.Logger.Println("")
	buildOptions.Logger.Println(fmt.Sprintf("Built '%s' in %s.", outputFilename, elapsed.Round(time.Millisecond).String()))
	buildOptions.Logger.Println("")

	return nil
}

func checkGoModVersion(logger *clilogger.CLILogger, updateGoMod bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	gomodFilename := filepath.Join(cwd, "go.mod")
	gomodData, err := os.ReadFile(gomodFilename)
	if err != nil {
		return err
	}
	outOfSync, err := gomod.GoModOutOfSync(gomodData, internal.Version)
	if err != nil {
		return err
	}
	if !outOfSync {
		return nil
	}
	gomodversion, err := gomod.GetWailsVersionFromModFile(gomodData)
	if err != nil {
		return err
	}

	if updateGoMod {
		return syncGoModVersion(cwd)
	}

	logger.Println("Warning: go.mod is using Wails '%s' but the CLI is '%s'. Consider updating your project's `go.mod` file.\n", gomodversion.String(), internal.Version)
	return nil
}

func LogGreen(message string, args ...interface{}) {
	text := fmt.Sprintf(message, args...)
	println(colour.Green(text))
}

func syncGoModVersion(cwd string) error {
	gomodFilename := filepath.Join(cwd, "go.mod")
	gomodData, err := os.ReadFile(gomodFilename)
	if err != nil {
		return err
	}
	outOfSync, err := gomod.GoModOutOfSync(gomodData, internal.Version)
	if err != nil {
		return err
	}
	if !outOfSync {
		return nil
	}
	LogGreen("Updating go.mod to use Wails '%s'", internal.Version)
	newGoData, err := gomod.UpdateGoModVersion(gomodData, internal.Version)
	if err != nil {
		return err
	}
	return os.WriteFile(gomodFilename, newGoData, 0755)
}
