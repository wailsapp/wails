package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/leaanthony/slicer"
	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v2/cmd/wails/flags"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/gomod"
	"github.com/wailsapp/wails/v2/internal/colour"
	"github.com/wailsapp/wails/v2/internal/project"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
	"github.com/wailsapp/wails/v2/pkg/commands/build"
)

func buildApplication(f *flags.Build) error {
	if f.NoColour {
		pterm.DisableColor()
		colour.ColourEnabled = false
	}

	quiet := f.Verbosity == flags.Quiet

	// Create logger
	logger := clilogger.New(os.Stdout)
	logger.Mute(quiet)

	if quiet {
		pterm.DisableOutput()
	} else {
		app.PrintBanner()
	}

	err := f.Process()
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

	// Set obfuscation from project file
	if projectOptions.Obfuscated {
		f.Obfuscated = projectOptions.Obfuscated
	}

	// Set garble args from project file
	if projectOptions.GarbleArgs != "" {
		f.GarbleArgs = projectOptions.GarbleArgs
	}

	// Create BuildOptions
	buildOptions := &build.Options{
		Logger:            logger,
		OutputType:        "desktop",
		OutputFile:        f.OutputFilename,
		CleanBinDirectory: f.Clean,
		Mode:              f.GetBuildMode(),
		Devtools:          f.Debug || f.Devtools,
		Pack:              !f.NoPackage,
		LDFlags:           f.LdFlags,
		Compiler:          f.Compiler,
		SkipModTidy:       f.SkipModTidy,
		Verbosity:         f.Verbosity,
		ForceBuild:        f.ForceBuild,
		IgnoreFrontend:    f.SkipFrontend,
		Compress:          f.Upx,
		CompressFlags:     f.UpxFlags,
		UserTags:          f.GetTags(),
		WebView2Strategy:  f.GetWebView2Strategy(),
		TrimPath:          f.TrimPath,
		RaceDetector:      f.RaceDetector,
		WindowsConsole:    f.WindowsConsole,
		Obfuscated:        f.Obfuscated,
		GarbleArgs:        f.GarbleArgs,
		SkipBindings:      f.SkipBindings,
		ProjectData:       projectOptions,
	}

	tableData := pterm.TableData{
		{"Platform(s)", f.Platform},
		{"Compiler", f.GetCompilerPath()},
		{"Skip Bindings", bool2Str(f.SkipBindings)},
		{"Build Mode", f.GetBuildModeAsString()},
		{"Devtools", bool2Str(buildOptions.Devtools)},
		{"Frontend Directory", projectOptions.GetFrontendDir()},
		{"Obfuscated", bool2Str(f.Obfuscated)},
	}
	if f.Obfuscated {
		tableData = append(tableData, []string{"Garble Args", f.GarbleArgs})
	}
	tableData = append(tableData, pterm.TableData{
		{"Skip Frontend", bool2Str(f.SkipFrontend)},
		{"Compress", bool2Str(f.Upx)},
		{"Package", bool2Str(!f.NoPackage)},
		{"Clean Bin Dir", bool2Str(f.Clean)},
		{"LDFlags", f.LdFlags},
		{"Tags", "[" + strings.Join(f.GetTags(), ",") + "]"},
		{"Race Detector", bool2Str(f.RaceDetector)},
	}...)
	if len(buildOptions.OutputFile) > 0 && f.GetTargets().Length() == 1 {
		tableData = append(tableData, []string{"Output File", f.OutputFilename})
	}
	pterm.DefaultSection.Println("Build Options")

	err = pterm.DefaultTable.WithData(tableData).Render()
	if err != nil {
		return err
	}

	if !f.NoSyncGoMod {
		err = gomod.SyncGoMod(logger, f.UpdateWailsVersionGoMod)
		if err != nil {
			return err
		}
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
	targets := f.GetTargets()
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
		buildOptions.Arch = f.GetDefaultArch()
		if len(platformSplit) > 1 {
			buildOptions.Arch = platformSplit[1]
		}
		banner := "Building target: " + buildOptions.Platform + "/" + buildOptions.Arch
		pterm.DefaultSection.Println(banner)

		if f.Upx && platform == "darwin/universal" {
			pterm.Warning.Println("Warning: compress flag unsupported for universal binaries. Ignoring.")
			f.Upx = false
		}

		switch buildOptions.Platform {
		case "linux":
			if runtime.GOOS != "linux" {
				pterm.Warning.Println("Crosscompiling to Linux not currently supported.")
				return
			}
		case "darwin":
			if runtime.GOOS != "darwin" {
				pterm.Warning.Println("Crosscompiling to Mac not currently supported.")
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

		if f.OutputFilename != "" {
			buildOptions.OutputFile = f.OutputFilename
		}

		if f.Obfuscated && f.SkipBindings {
			pterm.Warning.Println("obfuscated flag overrides skipbindings flag.")
			buildOptions.SkipBindings = false
		}

		if !f.DryRun {
			// Start Time
			start := time.Now()

			compiledBinary, err := build.Build(buildOptions)
			if err != nil {
				pterm.Error.Println(err.Error())
				targetErr = err
				return
			}

			buildOptions.IgnoreFrontend = true
			buildOptions.CleanBinDirectory = false

			// Output stats
			buildOptions.Logger.Println(fmt.Sprintf("Built '%s' in %s.\n", compiledBinary, time.Since(start).Round(time.Millisecond).String()))

			outputBinaries[buildOptions.Platform+"/"+buildOptions.Arch] = compiledBinary
		} else {
			pterm.Info.Println("Dry run: skipped build.")
		}
	})

	if targetErr != nil {
		return targetErr
	}

	if f.DryRun {
		return nil
	}

	if f.NSIS {
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
}
