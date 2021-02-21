package dev

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

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

// AddSubcommand adds the `dev` command for the Wails application
func AddSubcommand(app *clir.Cli, w io.Writer) error {

	command := app.NewSubCommand("dev", "Development mode")

	// Passthrough ldflags
	ldflags := ""
	command.StringFlag("ldflags", "optional ldflags", &ldflags)

	// compiler command
	compilerCommand := "go"
	command.StringFlag("compiler", "Use a different go compiler to build, eg go1.15beta1", &compilerCommand)

	// extensions to trigger rebuilds
	extensions := "go"
	command.StringFlag("e", "Extensions to trigger rebuilds (comma separated) eg go,js,css,html", &extensions)

	// extensions to trigger rebuilds
	showWarnings := false
	command.BoolFlag("w", "Show warnings", &showWarnings)

	loglevel := ""
	command.StringFlag("loglevel", "Loglevel to use - Trace, Debug, Info, Warning, Error", &loglevel)

	command.Action(func() error {

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
		var extensionsThatTriggerARebuild = strings.Split(extensions, ",")

		// Setup signal handler
		quitChannel := make(chan os.Signal, 1)
		signal.Notify(quitChannel, os.Interrupt, os.Kill, syscall.SIGTERM)

		debounceQuit := make(chan bool, 1)

		// Do initial build
		logger.Println("Building application for development...")
		newProcess, err := restartApp(logger, "dev", ldflags, compilerCommand, debugBinaryProcess, loglevel)
		if newProcess != nil {
			debugBinaryProcess = newProcess
		}
		if err != nil {
			return err
		}
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
						LogGreen("[New Directory] Watching new directory: %s", event.Name)
					}
				}
				return
			}

			// Check for file writes
			if event.Op&fsnotify.Write == fsnotify.Write {

				var rebuild bool = false

				// Iterate all file patterns
				for _, pattern := range extensionsThatTriggerARebuild {
					if strings.HasSuffix(event.Name, pattern) {
						rebuild = true
						break
					}
				}

				if !rebuild {
					if showWarnings {
						LogDarkYellow("[File change] %s did not match extension list (%s)", event.Name, extensions)
					}
					return
				}

				LogGreen("[Attempting rebuild] %s updated", event.Name)

				// Do a rebuild

				// Try and build the app
				newBinaryProcess, err := restartApp(logger, "dev", ldflags, compilerCommand, debugBinaryProcess, loglevel)
				if err != nil {
					fmt.Printf("Error during build: %s", err.Error())
					return
				}
				// If we have a new process, save it
				if newBinaryProcess != nil {
					debugBinaryProcess = newBinaryProcess
				}

			}
		})

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
				LogGreen("\nCaught quit")
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

		LogGreen("Development mode exited")

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

func restartApp(logger *clilogger.CLILogger, outputType string, ldflags string, compilerCommand string, debugBinaryProcess *process.Process, loglevel string) (*process.Process, error) {

	appBinary, err := buildApp(logger, outputType, ldflags, compilerCommand)
	println()
	if err != nil {
		LogRed("Build error - continuing to run current version")
		LogDarkYellow(err.Error())
		return nil, nil
	}

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
	newProcess := process.NewProcess(logger, appBinary, "-loglevel", loglevel)
	err = newProcess.Start()
	if err != nil {
		// Remove binary
		deleteError := fs.DeleteFile(appBinary)
		if deleteError != nil {
			logger.Fatal("Unable to delete app binary: " + appBinary)
		}
		logger.Fatal("Unable to start application: %s", err.Error())
	}

	return newProcess, nil
}

func buildApp(logger *clilogger.CLILogger, outputType string, ldflags string, compilerCommand string) (string, error) {

	// Create random output file
	outputFile := fmt.Sprintf("dev-%d", time.Now().Unix())

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
		IgnoreFrontend: true,
	}

	return build.Build(buildOptions)

}
