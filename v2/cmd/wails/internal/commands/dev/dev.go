package dev

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/wailsapp/wails/v2/cmd/wails/internal"
	"github.com/wailsapp/wails/v2/internal/gomod"

	"github.com/leaanthony/slicer"
	"github.com/wailsapp/wails/v2/internal/project"

	"github.com/pkg/browser"
	"github.com/wailsapp/wails/v2/internal/colour"

	"github.com/fsnotify/fsnotify"
	"github.com/leaanthony/clir"
	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/internal/process"
	"github.com/wailsapp/wails/v2/pkg/clilogger"
	"github.com/wailsapp/wails/v2/pkg/commands/build"
)

const defaultDevServerURL = "http://localhost:34115"

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

type devFlags struct {
	ldflags         string
	compilerCommand string
	assetDir        string
	extensions      string
	openBrowser     bool
	noReload        bool
	wailsjsdir      string
	tags            string
	verbosity       int
	loglevel        string
	forceBuild      bool
	debounceMS      int
	devServerURL    string
}

// AddSubcommand adds the `dev` command for the Wails application
func AddSubcommand(app *clir.Cli, w io.Writer) error {

	command := app.NewSubCommand("dev", "Development mode")

	flags := defaultDevFlags()
	command.StringFlag("ldflags", "optional ldflags", &flags.ldflags)
	command.StringFlag("compiler", "Use a different go compiler to build, eg go1.15beta1", &flags.compilerCommand)
	command.StringFlag("assetdir", "Serve assets from the given directory", &flags.assetDir)
	command.StringFlag("e", "Extensions to trigger rebuilds (comma separated) eg go", &flags.extensions)
	command.BoolFlag("browser", "Open application in browser", &flags.openBrowser)
	command.BoolFlag("noreload", "Disable reload on asset change", &flags.noReload)
	command.StringFlag("wailsjsdir", "Directory to generate the Wails JS modules", &flags.wailsjsdir)
	command.StringFlag("tags", "tags to pass to Go compiler (quoted and space separated)", &flags.tags)
	command.IntFlag("v", "Verbosity level (0 - silent, 1 - standard, 2 - verbose)", &flags.verbosity)
	command.StringFlag("loglevel", "Loglevel to use - Trace, Debug, Info, Warning, Error", &flags.loglevel)
	command.BoolFlag("f", "Force build application", &flags.forceBuild)
	command.IntFlag("debounce", "The amount of time to wait to trigger a reload on change", &flags.debounceMS)
	command.StringFlag("devserverurl", "The url of the dev server to use", &flags.devServerURL)

	command.Action(func() error {
		// Create logger
		logger := clilogger.New(w)
		app.PrintBanner()

		experimental := false
		userTags := []string{}
		for _, tag := range strings.Split(flags.tags, " ") {
			thisTag := strings.TrimSpace(tag)
			if thisTag != "" {
				userTags = append(userTags, thisTag)
			}
			if thisTag == "exp" {
				experimental = true
			}
		}

		if runtime.GOOS == "darwin" && !experimental {
			return fmt.Errorf("MacOS version coming soon!")
		}

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		projectConfig, err := loadAndMergeProjectConfig(cwd, &flags)
		if err != nil {
			return err
		}

		// Update go.mod to use current wails version
		err = syncGoModVersion(cwd)
		if err != nil {
			return err
		}

		// Run go mod tidy to ensure we're up to date
		err = runCommand(cwd, false, "go", "mod", "tidy")
		if err != nil {
			return err
		}

		if flags.tags != "" {
			err = runCommand(cwd, true, "wails", "generate", "module", "-tags", flags.tags)
		} else {
			err = runCommand(cwd, true, "wails", "generate", "module")
		}
		if err != nil {
			return err
		}

		// frontend:dev server command
		if projectConfig.DevCommand != "" {
			var devCommandWaitGroup sync.WaitGroup
			closer := runFrontendDevCommand(cwd, projectConfig.DevCommand, &devCommandWaitGroup)
			defer closer(&devCommandWaitGroup)
		}

		buildOptions := generateBuildOptions(flags)
		buildOptions.Logger = logger
		buildOptions.UserTags = internal.ParseUserTags(flags.tags)

		var debugBinaryProcess *process.Process = nil

		// Setup signal handler
		quitChannel := make(chan os.Signal, 1)
		signal.Notify(quitChannel, os.Interrupt, os.Kill, syscall.SIGTERM)
		exitCodeChannel := make(chan int, 1)

		// Do initial build
		logger.Println("Building application for development...")
		newProcess, appBinary, err := restartApp(buildOptions, debugBinaryProcess, flags, exitCodeChannel)
		if err != nil {
			return err
		}
		if newProcess != nil {
			debugBinaryProcess = newProcess
		}

		// open browser
		if flags.openBrowser {
			url := defaultDevServerURL
			if flags.devServerURL != "" {
				url = flags.devServerURL
			}
			err = browser.OpenURL(url)
			if err != nil {
				return err
			}
		}

		// create the project files watcher
		watcher, err := initialiseWatcher(cwd, logger.Fatal)
		defer func(watcher *fsnotify.Watcher) {
			err := watcher.Close()
			if err != nil {
				logger.Fatal(err.Error())
			}
		}(watcher)

		LogGreen("Watching (sub)/directory: %s", cwd)
		LogGreen("Using Dev Server URL: %s", flags.devServerURL)
		LogGreen("Using reload debounce setting of %d milliseconds", flags.debounceMS)

		// Watch for changes and trigger restartApp()
		doWatcherLoop(buildOptions, debugBinaryProcess, flags, watcher, exitCodeChannel, quitChannel)

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

func runCommand(dir string, exitOnError bool, command string, args ...string) error {
	LogGreen("Executing: " + command + " " + strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		println(string(output))
		if exitOnError {
			os.Exit(1)
		}
		return err
	}
	return nil
}

// defaultDevFlags generates devFlags with default options
func defaultDevFlags() devFlags {
	return devFlags{
		devServerURL:    defaultDevServerURL,
		compilerCommand: "go",
		verbosity:       1,
		extensions:      "go",
		debounceMS:      100,
	}
}

// generateBuildOptions creates a build.Options using the flags
func generateBuildOptions(flags devFlags) *build.Options {
	result := &build.Options{
		OutputType:     "dev",
		Mode:           build.Dev,
		Arch:           runtime.GOARCH,
		Pack:           false,
		Platform:       runtime.GOOS,
		LDFlags:        flags.ldflags,
		Compiler:       flags.compilerCommand,
		ForceBuild:     flags.forceBuild,
		IgnoreFrontend: false,
		Verbosity:      flags.verbosity,
		WailsJSDir:     flags.wailsjsdir,
	}
	switch runtime.GOOS {
	case "darwin":
		result.Pack = false
	}
	return result
}

// loadAndMergeProjectConfig reconciles flags passed to the CLI with project config settings and updates
// the project config if necessary
func loadAndMergeProjectConfig(cwd string, flags *devFlags) (*project.Project, error) {
	projectConfig, err := project.Load(cwd)
	if err != nil {
		return nil, err
	}

	var shouldSaveConfig bool

	if projectConfig.AssetDirectory == "" && flags.assetDir == "" {
		return nil, fmt.Errorf("No asset directory provided. Please use -assetdir to indicate which directory contains your built assets.")
	}

	if flags.assetDir == "" && projectConfig.AssetDirectory != "" {
		flags.assetDir = projectConfig.AssetDirectory
	}

	if flags.assetDir != projectConfig.AssetDirectory {
		projectConfig.AssetDirectory = filepath.ToSlash(flags.assetDir)
	}

	flags.assetDir, err = filepath.Abs(flags.assetDir)
	if err != nil {
		return nil, err
	}

	if flags.devServerURL == defaultDevServerURL && projectConfig.DevServerURL != defaultDevServerURL && projectConfig.DevServerURL != "" {
		flags.devServerURL = projectConfig.DevServerURL
	}

	if flags.devServerURL != projectConfig.DevServerURL {
		projectConfig.DevServerURL = flags.devServerURL
		shouldSaveConfig = true
	}

	if flags.wailsjsdir == "" && projectConfig.WailsJSDir != "" {
		flags.wailsjsdir = projectConfig.WailsJSDir
	}

	if flags.wailsjsdir == "" {
		flags.wailsjsdir = "./frontend"
	}

	if flags.wailsjsdir != projectConfig.WailsJSDir {
		projectConfig.WailsJSDir = filepath.ToSlash(flags.wailsjsdir)
		shouldSaveConfig = true
	}

	if flags.debounceMS == 100 && projectConfig.DebounceMS != 100 {
		if projectConfig.DebounceMS == 0 {
			projectConfig.DebounceMS = 100
		}
		flags.debounceMS = projectConfig.DebounceMS
	}

	if flags.debounceMS != projectConfig.DebounceMS {
		projectConfig.DebounceMS = flags.debounceMS
		shouldSaveConfig = true
	}

	if shouldSaveConfig {
		err = projectConfig.Save()
		if err != nil {
			return nil, err
		}
	}

	return projectConfig, nil
}

// runFrontendDevCommand will run the `frontend:dev` command if it was given, ex- `npm run dev`
func runFrontendDevCommand(cwd string, devCommand string, wg *sync.WaitGroup) func(group *sync.WaitGroup) {
	LogGreen("Running frontend dev command: '%s'", devCommand)
	ctx, cancel := context.WithCancel(context.Background())
	dir := filepath.Join(cwd, "frontend")
	cmdSlice := strings.Split(devCommand, " ")
	wg.Add(1)
	cmd := exec.CommandContext(ctx, cmdSlice[0], cmdSlice[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = dir
	go func(ctx context.Context, devCommand string, cwd string, wg *sync.WaitGroup) {
		err := cmd.Run()
		if err != nil {
			if err.Error() != "exit status 1" {
				LogRed("Error from '%s': %s", devCommand, err.Error())
			}
		}
		LogGreen("Dev command exited!")
		wg.Done()
	}(ctx, devCommand, cwd, wg)

	return func(wg *sync.WaitGroup) {
		if runtime.GOOS == "windows" {
			// Credit: https://stackoverflow.com/a/44551450
			// For whatever reason, killing an npm script on windows just doesn't exit properly with cancel
			if cmd != nil && cmd.Process != nil {
				kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid))
				kill.Stderr = os.Stderr
				kill.Stdout = os.Stdout
				err := kill.Run()
				if err != nil {
					if err.Error() != "exit status 1" {
						LogRed("Error from '%s': %s", devCommand, err.Error())
					}
				}
			}
		} else {
			cancel()
		}
		wg.Wait()
	}
}

// initialiseWatcher creates the project directory watcher that will trigger recompile
func initialiseWatcher(cwd string, logFatal func(string, ...interface{})) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// Get all subdirectories
	dirs, err := fs.GetSubdirectories(cwd)
	if err != nil {
		return nil, err
	}

	// Setup a watcher for non-node_modules directories
	dirs.Each(func(dir string) {
		if strings.Contains(dir, "node_modules") {
			return
		}
		// Ignore build directory
		if strings.HasPrefix(dir, filepath.Join(cwd, "build")) {
			return
		}
		// Ignore dot directories
		if strings.HasPrefix(dir, ".") {
			return
		}
		err = watcher.Add(dir)
		if err != nil {
			logFatal(err.Error())
		}
	})
	return watcher, nil
}

// restartApp does the actual rebuilding of the application when files change
func restartApp(buildOptions *build.Options, debugBinaryProcess *process.Process, flags devFlags, exitCodeChannel chan int) (*process.Process, string, error) {

	appBinary, err := build.Build(buildOptions)
	println()
	if err != nil {
		LogRed("Build error - continuing to run current version")
		LogDarkYellow(err.Error())
		return nil, "", nil
	}

	// Kill existing binary if need be
	if debugBinaryProcess != nil {
		killError := debugBinaryProcess.Kill()

		if killError != nil {
			buildOptions.Logger.Fatal("Unable to kill debug binary (PID: %d)!", debugBinaryProcess.PID())
		}

		debugBinaryProcess = nil
	}
	args := slicer.StringSlicer{}

	// Set environment variables accordingly
	os.Setenv("loglevel", flags.loglevel)
	os.Setenv("assetdir", flags.assetDir)
	os.Setenv("devserverurl", flags.devServerURL)

	// Start up new binary with correct args
	newProcess := process.NewProcess(appBinary, args.AsSlice()...)
	err = newProcess.Start(exitCodeChannel)
	if err != nil {
		// Remove binary
		if fs.FileExists(appBinary) {
			deleteError := fs.DeleteFile(appBinary)
			if deleteError != nil {
				buildOptions.Logger.Fatal("Unable to delete app binary: " + appBinary)
			}
		}
		buildOptions.Logger.Fatal("Unable to start application: %s", err.Error())
	}

	return newProcess, appBinary, nil
}

// doWatcherLoop is the main watch loop that runs while dev is active
func doWatcherLoop(buildOptions *build.Options, debugBinaryProcess *process.Process, flags devFlags, watcher *fsnotify.Watcher, exitCodeChannel chan int, quitChannel chan os.Signal) {
	// Main Loop
	var (
		err              error
		newBinaryProcess *process.Process
	)
	var extensionsThatTriggerARebuild = sliceToMap(strings.Split(flags.extensions, ","))
	quit := false
	interval := time.Duration(flags.debounceMS) * time.Millisecond
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
						timer.Reset(interval)
						continue
					}
				}

				if strings.HasPrefix(item.Name, flags.assetDir) {
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
							buildOptions.Logger.Fatal("%s", err.Error())
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
				newBinaryProcess, _, err = restartApp(buildOptions, debugBinaryProcess, flags, exitCodeChannel)
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
}
