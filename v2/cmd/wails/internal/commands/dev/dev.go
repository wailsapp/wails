package dev

import (
	"fmt"
	"github.com/leaanthony/slicer"
	"github.com/wailsapp/wails/v2/internal/project"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/browser"
	"github.com/wailsapp/wails/v2/internal/colour"

	"github.com/fsnotify/fsnotify"
	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/internal/process"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
	"github.com/wailsapp/wails/v2/pkg/commands/build"
)

func LogGreen(message string, args ...interface{}) {
	text := fmt.Sprintf(message, args...)
	println(colour.Green(text))
}

func LogRed(message string, args ...interface{}) {
	text := fmt.Sprintf(message, args...)
	println(colour.Red(text))
}

func LogDarkYellow(message string, args ...interface{}) {
	text := fmt.Sprintf(message, args...)
	println(colour.DarkYellow(text))
}

func sliceToMap(input []string) map[string]struct{} {
	result := map[string]struct{}{}
	for _, value := range input {
		result[value] = struct{}{}
	}
	return result
}

// AddSubcommand adds the `dev` command for the Wails application
func AddSubcommand(app *clir.Cli, w io.Writer) error {

	command := app.NewSubCommand("dev", "Development mode")

	// Passthrough ldflags
	ldflags := ""
	command.StringFlag("ldflags", "optional ldflags", &ldflags)

	// compiler command
	compilerCommand := "go"
	command.StringFlag("compiler", "Use a different go compiler to build, eg go1.15beta1", &compilerCommand)

	assetDir := ""
	command.StringFlag("assetdir", "Serve assets from the given directory", &assetDir)

	// extensions to trigger rebuilds of application
	extensions := "go"
	command.StringFlag("e", "Extensions to trigger rebuilds (comma separated) eg go", &extensions)

	openBrowser := false
	command.BoolFlag("browser", "Open application in browser", &openBrowser)

	noreload := false
	command.BoolFlag("noreload", "Disable reload on asset change", &noreload)

	wailsjsdir := ""
	command.StringFlag("wailsjsdir", "Directory to generate the Wails JS modules", &wailsjsdir)

	// tags to pass to `go`
	tags := ""
	command.StringFlag("tags", "tags to pass to Go compiler (quoted and space separated)", &tags)

	// Verbosity
	verbosity := 1
	command.IntFlag("v", "Verbosity level (0 - silent, 1 - standard, 2 - verbose)", &verbosity)

	loglevel := ""
	command.StringFlag("loglevel", "Loglevel to use - Trace, Dev, Info, Warning, Error", &loglevel)

	forceBuild := false
	command.BoolFlag("f", "Force build application", &forceBuild)

	command.Action(func() error {

		// Create logger
		logger := clilogger.New(w)
		app.PrintBanner()

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		projectConfig, err := project.Load(cwd)
		if err != nil {
			return err
		}

		if projectConfig.AssetDirectory == "" && assetDir == "" {
			return fmt.Errorf("No asset directory provided. Please use -assetdir to indicate which directory contains your built assets.")
		}

		if assetDir == "" && projectConfig.AssetDirectory != "" {
			assetDir = projectConfig.AssetDirectory
		}

		if assetDir != projectConfig.AssetDirectory {
			projectConfig.AssetDirectory = filepath.ToSlash(assetDir)
			err := projectConfig.Save()
			if err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}

		if wailsjsdir == "" && projectConfig.WailsJSDir != "" {
			wailsjsdir = projectConfig.WailsJSDir
		}

		if wailsjsdir == "" {
			wailsjsdir = "./frontend"
		}

		if wailsjsdir != projectConfig.WailsJSDir {
			projectConfig.WailsJSDir = filepath.ToSlash(wailsjsdir)
			err := projectConfig.Save()
			if err != nil {
				return err
			}
		}

		buildOptions := &build.Options{
			Logger:         logger,
			OutputType:     "dev",
			Mode:           build.Dev,
			Arch:           runtime.GOARCH,
			Pack:           true,
			Platform:       runtime.GOOS,
			LDFlags:        ldflags,
			Compiler:       compilerCommand,
			ForceBuild:     forceBuild,
			IgnoreFrontend: false,
			Verbosity:      verbosity,
			WailsJSDir:     wailsjsdir,
		}

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return err
		}
		defer func(watcher *fsnotify.Watcher) {
			err := watcher.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(watcher)

		var debugBinaryProcess *process.Process = nil
		var extensionsThatTriggerARebuild = sliceToMap(strings.Split(extensions, ","))

		// Setup signal handler
		quitChannel := make(chan os.Signal, 1)
		signal.Notify(quitChannel, os.Interrupt, os.Kill, syscall.SIGTERM)
		exitCodeChannel := make(chan int, 1)

		var passthruArgs []string
		//if len(os.Args) > 2 {
		//	passthruArgs = os.Args[2:]
		//}

		// Do initial build
		logger.Println("Building application for development...")
		newProcess, appBinary, err := restartApp(logger, buildOptions, debugBinaryProcess, loglevel, passthruArgs, assetDir, false, exitCodeChannel)
		if err != nil {
			return err
		}
		if newProcess != nil {
			debugBinaryProcess = newProcess
		}

		// open browser
		if openBrowser {
			err = browser.OpenURL("http://localhost:34115")
			if err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
		var newBinaryProcess *process.Process

		// Get project dir
		projectDir, err := os.Getwd()
		if err != nil {
			return err
		}

		// Get all subdirectories
		dirs, err := fs.GetSubdirectories(projectDir)
		if err != nil {
			return err
		}

		LogGreen("Watching (sub)/directory: %s", projectDir)

		// Setup a watcher for non-node_modules directories
		dirs.Each(func(dir string) {
			if strings.Contains(dir, "node_modules") {
				return
			}
			// Ignore build directory
			if strings.HasPrefix(dir, filepath.Join(projectDir, "build")) {
				return
			}
			//println("Watching", dir)
			err = watcher.Add(dir)
			if err != nil {
				logger.Fatal(err.Error())
			}
		})

		// Main Loop
		quit := false
		// Use 100ms debounce
		interval := 100 * time.Millisecond
		timer := time.NewTimer(interval)
		rebuild := false
		reload := false
		for quit == false {
			//reload := false
			select {
			case exitCode := <-exitCodeChannel:
				if exitCode == 0 {
					quit = true
				}
			case item := <-watcher.Events:
				// Check for file writes
				if item.Op&fsnotify.Write == fsnotify.Write {
					// Ignore directories
					if fs.DirExists(item.Name) {
						continue
					}

					// Iterate all file patterns
					ext := filepath.Ext(item.Name)
					if ext != "" {
						ext = ext[1:]
						if _, exists := extensionsThatTriggerARebuild[ext]; exists {
							rebuild = true
							continue
						}
					}

					if strings.HasPrefix(item.Name, assetDir) {
						reload = true
					}
					timer.Reset(interval)
				}
				// Check for new directories
				if item.Op&fsnotify.Create == fsnotify.Create {
					// If this is a folder, add it to our watch list
					if fs.DirExists(item.Name) {
						//node_modules is BANNED!
						if !strings.Contains(item.Name, "node_modules") {
							err := watcher.Add(item.Name)
							if err != nil {
								logger.Fatal("%s", err.Error())
							}
							LogGreen("Added new directory to watcher: %s", item.Name)
						}
					}
				}
			case <-timer.C:
				if rebuild {
					rebuild = false
					LogGreen("[Rebuild triggered] files updated")
					// Try and build the app
					newBinaryProcess, _, err = restartApp(logger, buildOptions, debugBinaryProcess, loglevel, passthruArgs, assetDir, false, exitCodeChannel)
					if err != nil {
						LogRed("Error during build: %s", err.Error())
						continue
					}
					// If we have a new process, save it
					if newBinaryProcess != nil {
						debugBinaryProcess = newBinaryProcess
					}
				}
				if reload {
					reload = false
					_, err = http.Get("http://localhost:34115/wails/reload")
					if err != nil {
						LogRed("Error during refresh: %s", err.Error())
					}
				}
			case <-quitChannel:
				LogGreen("\nCaught quit")
				quit = true
			}
		}

		// Kill the current program if running
		if debugBinaryProcess != nil {
			err := debugBinaryProcess.Kill()
			if err != nil {
				return err
			}
		}

		// Remove dev binary
		err = os.Remove(appBinary)
		if err != nil {
			return err
		}

		LogGreen("Development mode exited")

		return nil
	})

	return nil
}

func restartApp(logger *clilogger.CLILogger, buildOptions *build.Options, debugBinaryProcess *process.Process, loglevel string, passthruArgs []string, assetDir string, firstRun bool, exitCodeChannel chan int) (*process.Process, string, error) {

	appBinary, err := build.Build(buildOptions)
	println()
	if err != nil {
		if firstRun {
			return nil, "", err
		}
		LogRed("Build error - continuing to run current version")
		LogDarkYellow(err.Error())
		return nil, "", nil
	}

	// Kill existing binary if need be
	if debugBinaryProcess != nil {
		killError := debugBinaryProcess.Kill()

		if killError != nil {
			logger.Fatal("Unable to kill debug binary (PID: %d)!", debugBinaryProcess.PID())
		}

		debugBinaryProcess = nil
	}

	// Start up new binary with correct args
	args := slicer.StringSlicer{}
	args.Add("-loglevel", loglevel)
	if assetDir != "" {
		args.Add("-assetdir", assetDir)
	}

	if len(passthruArgs) > 0 {
		args.AddSlice(passthruArgs)
	}
	newProcess := process.NewProcess(appBinary, args.AsSlice()...)
	err = newProcess.Start(exitCodeChannel)
	if err != nil {
		// Remove binary
		deleteError := fs.DeleteFile(appBinary)
		if deleteError != nil {
			logger.Fatal("Unable to delete app binary: " + appBinary)
		}
		logger.Fatal("Unable to start application: %s", err.Error())
	}

	return newProcess, appBinary, nil
}
