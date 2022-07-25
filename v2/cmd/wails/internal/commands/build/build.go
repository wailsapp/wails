package build

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/wailsapp/wails/v2/internal/colour"
	"github.com/wailsapp/wails/v2/internal/project"
	"github.com/wailsapp/wails/v2/internal/system"

	"github.com/wailsapp/wails/v2/cmd/wails/internal"
	"github.com/wailsapp/wails/v2/internal/gomod"

	"github.com/leaanthony/clir"
	"github.com/leaanthony/slicer"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
	"github.com/wailsapp/wails/v2/pkg/commands/build"
)

// AddBuildSubcommand adds the `build` command for the Wails application
func AddBuildSubcommand(app *clir.Cli, w io.Writer) {
	command := app.NewSubCommand("build", "Builds the application")

	outputType := "desktop"
	validTargetTypes := slicer.String([]string{"desktop", "hybrid", "server"})

	command.StringFlag("outputType", fmt.Sprintf("Type of binary %s", validTargetTypes.AsSlice()), &outputType)

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

	defaultPlatform := os.Getenv("GOOS")
	if defaultPlatform == "" {
		defaultPlatform = runtime.GOOS
	}
	defaultArch := os.Getenv("GOARCH")
	if defaultArch == "" {
		if system.IsAppleSilicon {
			defaultArch = "arm64"
		} else {
			defaultArch = runtime.GOARCH
		}
	}
	platform := defaultPlatform + "/" + defaultArch

	command.StringFlag("platform", "Platform to target. Comma separate multiple platforms", &platform)

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

	nsis := false
	command.BoolFlag("nsis", "Generate NSIS installer for Windows", &nsis)

	trimpath := false
	command.BoolFlag("trimpath", "Remove all file system paths from the resulting executable", &trimpath)

	raceDetector := false
	command.BoolFlag("race", "Build with Go's race detector", &raceDetector)

	windowsConsole := false
	command.BoolFlag("windowsconsole", "Keep the console when building for Windows", &windowsConsole)

	dryRun := false
	command.BoolFlag("dryrun", "Dry run, prints the config for the command that would be executed", &dryRun)

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

		var targets slicer.StringSlicer
		targets.AddSlice(strings.Split(platform, ","))
		targets.Deduplicate()

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
			TrimPath:            trimpath,
			RaceDetector:        raceDetector,
			WindowsConsole:      windowsConsole,
		}

		// Start a new tabwriter
		if !quiet {
			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 8, 8, 0, '\t', 0)

			// Write out the system information
			_, _ = fmt.Fprintf(w, "App Type: \t%s\n", buildOptions.OutputType)
			_, _ = fmt.Fprintf(w, "Platforms: \t%s\n", platform)
			_, _ = fmt.Fprintf(w, "Compiler: \t%s\n", compilerPath)
			_, _ = fmt.Fprintf(w, "Build Mode: \t%s\n", modeString)
			_, _ = fmt.Fprintf(w, "Skip Frontend: \t%t\n", skipFrontend)
			_, _ = fmt.Fprintf(w, "Compress: \t%t\n", buildOptions.Compress)
			_, _ = fmt.Fprintf(w, "Package: \t%t\n", buildOptions.Pack)
			_, _ = fmt.Fprintf(w, "Clean Build Dir: \t%t\n", buildOptions.CleanBuildDirectory)
			_, _ = fmt.Fprintf(w, "LDFlags: \t\"%s\"\n", buildOptions.LDFlags)
			_, _ = fmt.Fprintf(w, "Tags: \t[%s]\n", strings.Join(buildOptions.UserTags, ","))
			_, _ = fmt.Fprintf(w, "Race Detector: \t%t\n", buildOptions.RaceDetector)
			if len(buildOptions.OutputFile) > 0 && targets.Length() == 1 {
				_, _ = fmt.Fprintf(w, "Output File: \t%s\n", buildOptions.OutputFile)
			}
			_, _ = fmt.Fprintf(w, "\n")
			err = w.Flush()
			if err != nil {
				return err
			}
		}
		err = checkGoModVersion(logger, updateGoMod)
		if err != nil {
			return err
		}

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		projectOptions, err := project.Load(cwd)
		if err != nil {
			return err
		}

		// Check platform
		validPlatformArch := slicer.String([]string{
			"darwin",
			"darwin/amd64",
			"darwin/arm64",
			"darwin/universal",
			"linux",
			"linux/amd64",
			"linux/arm64",
			"linux/arm",
			"windows",
			"windows/amd64",
			"windows/arm64",
			"windows/386",
		})

		outputBinaries := map[string]string{}

		// Allows cancelling the build after the first error. It would be nice if targets.Each would support funcs
		// returning an error.
		var targetErr error
		targets.Each(func(platform string) {
			if targetErr != nil {
				return
			}

			if !validPlatformArch.Contains(platform) {
				buildOptions.Logger.Println("platform '%s' is not supported - skipping. Supported platforms: %s", platform, validPlatformArch.Join(","))
				return
			}

			desiredFilename := projectOptions.OutputFilename
			if desiredFilename == "" {
				desiredFilename = projectOptions.Name
			}
			desiredFilename = strings.TrimSuffix(desiredFilename, ".exe")

			// Calculate platform and arch
			platformSplit := strings.Split(platform, "/")
			buildOptions.Platform = platformSplit[0]
			buildOptions.Arch = defaultArch
			if len(platformSplit) > 1 {
				buildOptions.Arch = platformSplit[1]
			}
			banner := "Building target: " + buildOptions.Platform + "/" + buildOptions.Arch
			logger.Println(banner)
			logger.Println(strings.Repeat("-", len(banner)))

			if compress && platform == "darwin/universal" {
				logger.Println("Warning: compress flag unsupported for universal binaries. Ignoring.")
				compress = false
			}

			switch buildOptions.Platform {
			case "linux":
				if runtime.GOOS != "linux" {
					logger.Println("Crosscompiling to Linux not currently supported.\n")
					return
				}
			case "darwin":
				if runtime.GOOS != "darwin" {
					logger.Println("Crosscompiling to Mac not currently supported.\n")
					return
				}
				macTargets := targets.Filter(func(platform string) bool {
					return strings.HasPrefix(platform, "darwin")
				})
				if macTargets.Length() == 2 {
					buildOptions.BundleName = fmt.Sprintf("%s-%s.app", desiredFilename, buildOptions.Arch)
				}
			}

			if targets.Length() > 1 {
				// target filename
				switch buildOptions.Platform {
				case "windows":
					desiredFilename = fmt.Sprintf("%s-%s", desiredFilename, buildOptions.Arch)
				case "linux", "darwin":
					desiredFilename = fmt.Sprintf("%s-%s-%s", desiredFilename, buildOptions.Platform, buildOptions.Arch)
				}
			}
			if buildOptions.Platform == "windows" {
				desiredFilename += ".exe"
			}
			buildOptions.OutputFile = desiredFilename

			if outputFilename != "" {
				buildOptions.OutputFile = outputFilename
			}

			if !dryRun {
				// Start Time
				start := time.Now()

				compiledBinary, err := build.Build(buildOptions)
				if err != nil {
					logger.Println("Error: %s", err.Error())
					targetErr = err
					return
				}

				buildOptions.IgnoreFrontend = true
				buildOptions.CleanBuildDirectory = false

				// Output stats
				buildOptions.Logger.Println(fmt.Sprintf("Built '%s' in %s.\n", compiledBinary, time.Since(start).Round(time.Millisecond).String()))

				outputBinaries[buildOptions.Platform+"/"+buildOptions.Arch] = compiledBinary
			} else {
				logger.Println("Dry run: skipped build.")
			}
		})

		if targetErr != nil {
			return targetErr
		}

		if dryRun {
			return nil
		}

		if nsis {
			amd64Binary := outputBinaries["windows/amd64"]
			arm64Binary := outputBinaries["windows/arm64"]
			if amd64Binary == "" && arm64Binary == "" {
				return fmt.Errorf("cannot build nsis installer - no windows targets")
			}

			if err := build.GenerateNSISInstaller(buildOptions, amd64Binary, arm64Binary); err != nil {
				return err
			}
		}
		return nil
	})
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
