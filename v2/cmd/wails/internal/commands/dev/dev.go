package dev

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/leaanthony/clir"
	"github.com/leaanthony/slicer"
	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/internal/process"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
	"github.com/wailsapp/wails/v2/pkg/commands/build"
)

// AddSubcommand adds the `dev` command for the Wails application
func AddSubcommand(app *clir.Cli, w io.Writer) error {

	command := app.NewSubCommand("dev", "Development mode")

	outputType := "desktop"

	validTargetTypes := slicer.String([]string{"desktop", "hybrid", "server"})

	// Setup target type flag
	description := "Type of application to develop. Valid types: " + validTargetTypes.Join(",")
	command.StringFlag("t", description, &outputType)

	// Passthrough ldflags
	ldflags := ""
	command.StringFlag("ldflags", "optional ldflags", &ldflags)

	// compiler command
	compilerCommand := "go"
	command.StringFlag("compiler", "Use a different go compiler to build, eg go1.15beta1", &compilerCommand)

	// extensions to trigger rebuilds
	extensions := "go"
	command.StringFlag("m", "Extensions to trigger rebuilds (comma separated) eg go,js,css,html", &extensions)

	command.Action(func() error {

		// Validate inputs
		if !validTargetTypes.Contains(outputType) {
			return fmt.Errorf("output type '%s' is not valid", outputType)
		}

		// Create logger
		logger := clilogger.New(w)
		app.PrintBanner()

		// TODO: Check you are in a project directory

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return err
		}
		defer watcher.Close()

		var debugBinaryProcess *process.Process = nil
		var buildFrontend bool = true
		var extensionsThatTriggerARebuild = strings.Split(extensions, ",")

		// Setup signal handler
		quitChannel := make(chan os.Signal, 1)
		signal.Notify(quitChannel, os.Interrupt, os.Kill, syscall.SIGTERM)

		debounceQuit := make(chan bool, 1)

		// Do initial build
		logger.Println("Building application for development...")
		debugBinaryProcess = restartApp(logger, outputType, ldflags, compilerCommand, buildFrontend, debugBinaryProcess)

		go debounce(100*time.Millisecond, watcher.Events, debounceQuit, func(event fsnotify.Event) {
			// logger.Println("event: %+v", event)

			// Check for new directories
			if event.Op&fsnotify.Create == fsnotify.Create {
				// If this is a folder, add it to our watch list
				if fs.DirExists(event.Name) {
					if !strings.Contains(event.Name, "node_modules") {
						err := watcher.Add(event.Name)
						if err != nil {
							logger.Fatal("%s", err.Error())
						}
						logger.Println("Watching directory: %s", event.Name)
					}
				}
				return
			}

			// Check for file writes
			if event.Op&fsnotify.Write == fsnotify.Write {

				// logger.Println("modified file: %s", event.Name)
				var rebuild bool = false

				// Iterate all file patterns
				for _, pattern := range extensionsThatTriggerARebuild {
					rebuild = strings.HasSuffix(event.Name, pattern)
					if err != nil {
						logger.Fatal(err.Error())
					}
					if rebuild {
						// Only build frontend when the file isn't a Go file
						buildFrontend = !strings.HasSuffix(event.Name, "go")
						break
					}
				}

				if !rebuild {
					logger.Println("Filename change: %s did not match extension list %s", event.Name, extensions)
					return
				}

				if buildFrontend {
					logger.Println("Full rebuild triggered: %s updated", event.Name)
				} else {
					logger.Println("Partial build triggered: %s updated", event.Name)
				}

				// Do a rebuild

				// Try and build the app
				newBinaryProcess := restartApp(logger, outputType, ldflags, compilerCommand, buildFrontend, debugBinaryProcess)

				// If we have a new process, save it
				if newBinaryProcess != nil {
					debugBinaryProcess = newBinaryProcess
				}

			}
		})

		// Get project dir
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		// Get all subdirectories
		dirs, err := fs.GetSubdirectories(dir)
		if err != nil {
			return err
		}

		// Setup a watcher for non-node_modules directories
		dirs.Each(func(dir string) {
			if strings.Contains(dir, "node_modules") {
				return
			}
			logger.Println("Watching directory: %s", dir)
			err = watcher.Add(dir)
			if err != nil {
				logger.Fatal(err.Error())
			}
		})

		// Wait until we get a quit signal
		quit := false
		for quit == false {
			select {
			case <-quitChannel:
				println()
				// Notify debouncer to quit
				debounceQuit <- true
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

		logger.Println("Development mode exited")

		return nil
	})

	return nil
}

// Credit: https://drailing.net/2018/01/debounce-function-for-golang/
func debounce(interval time.Duration, input chan fsnotify.Event, quitChannel chan bool, cb func(arg fsnotify.Event)) {
	var item fsnotify.Event
	timer := time.NewTimer(interval)
exit:
	for {
		select {
		case item = <-input:
			timer.Reset(interval)
		case <-timer.C:
			if item.Name != "" {
				cb(item)
			}
		case <-quitChannel:
			break exit
		}
	}
}

func restartApp(logger *clilogger.CLILogger, outputType string, ldflags string, compilerCommand string, buildFrontend bool, debugBinaryProcess *process.Process) *process.Process {

	appBinary, err := buildApp(logger, outputType, ldflags, compilerCommand, buildFrontend)
	println()
	if err != nil {
		logger.Println("[ERROR] Build Failed: %s", err.Error())
		return nil
	}
	logger.Println("Build new binary: %s", appBinary)

	// Kill existing binary if need be
	if debugBinaryProcess != nil {
		killError := debugBinaryProcess.Kill()

		if killError != nil {
			logger.Fatal("Unable to kill debug binary (PID: %d)!", debugBinaryProcess.PID())
		}

		debugBinaryProcess = nil
	}

	// TODO: Generate `backend.js`

	// Start up new binary
	newProcess := process.NewProcess(logger, appBinary)
	err = newProcess.Start()
	if err != nil {
		// Remove binary
		deleteError := fs.DeleteFile(appBinary)
		if deleteError != nil {
			logger.Fatal("Unable to delete app binary: " + appBinary)
		}
		logger.Fatal("Unable to start application: %s", err.Error())
	}

	return newProcess
}

func buildApp(logger *clilogger.CLILogger, outputType string, ldflags string, compilerCommand string, buildFrontend bool) (string, error) {

	// Create random output file
	outputFile := fmt.Sprintf("debug-%d", time.Now().Unix())

	// Create BuildOptions
	buildOptions := &build.Options{
		Logger:         logger,
		OutputType:     outputType,
		Mode:           build.Debug,
		Pack:           false,
		Platform:       runtime.GOOS,
		LDFlags:        ldflags,
		Compiler:       compilerCommand,
		OutputFile:     outputFile,
		IgnoreFrontend: !buildFrontend,
	}

	return build.Build(buildOptions)

}
