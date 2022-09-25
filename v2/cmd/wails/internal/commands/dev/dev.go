package dev

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/wailsapp/wails/v2/pkg/commands/bindings"
	"github.com/wailsapp/wails/v2/pkg/commands/buildtags"

	"github.com/google/shlex"
	buildcmd "github.com/wailsapp/wails/v2/cmd/wails/internal/commands/build"

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

func LogGreen(message string, args ...interface{}) {
	if len(message) == 0 {
		return
	}
	text := fmt.Sprintf(message, args...)
	println(colour.Green(text))
}

func LogRed(message string, args ...interface{}) {
	if len(message) == 0 {
		return
	}
	text := fmt.Sprintf(message, args...)
	println(colour.Red(text))
}

func LogDarkYellow(message string, args ...interface{}) {
	if len(message) == 0 {
		return
	}
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
	reloadDirs      string
	openBrowser     bool
	noReload        bool
	skipBindings    bool
	wailsjsdir      string
	tags            string
	verbosity       int
	loglevel        string
	forceBuild      bool
	debounceMS      int
	devServer       string
	appargs         string
	saveConfig      bool
	raceDetector    bool

	frontendDevServerURL string
	skipFrontend         bool
}

// AddSubcommand adds the `dev` command for the Wails application
func AddSubcommand(app *clir.Cli, w io.Writer) error {

	command := app.NewSubCommand("dev", "Development mode")

	flags := defaultDevFlags()
	command.StringFlag("ldflags", "optional ldflags", &flags.ldflags)
	command.StringFlag("compiler", "Use a different go compiler to build, eg go1.15beta1", &flags.compilerCommand)
	command.StringFlag("assetdir", "Serve assets from the given directory instead of using the provided asset FS", &flags.assetDir)
	command.StringFlag("e", "Extensions to trigger rebuilds (comma separated) eg go", &flags.extensions)
	command.StringFlag("reloaddirs", "Additional directories to trigger reloads (comma separated)", &flags.reloadDirs)
	command.BoolFlag("browser", "Open application in browser", &flags.openBrowser)
	command.BoolFlag("noreload", "Disable reload on asset change", &flags.noReload)
	command.BoolFlag("skipbindings", "Skip bindings generation", &flags.skipBindings)
	command.StringFlag("wailsjsdir", "Directory to generate the Wails JS modules", &flags.wailsjsdir)
	command.StringFlag("tags", "Build tags to pass to Go compiler. Must be quoted. Space or comma (but not both) separated", &flags.tags)
	command.IntFlag("v", "Verbosity level (0 - silent, 1 - standard, 2 - verbose)", &flags.verbosity)
	command.StringFlag("loglevel", "Loglevel to use - Trace, Debug, Info, Warning, Error", &flags.loglevel)
	command.BoolFlag("f", "Force build application", &flags.forceBuild)
	command.IntFlag("debounce", "The amount of time to wait to trigger a reload on change", &flags.debounceMS)
	command.StringFlag("devserver", "The address of the wails dev server", &flags.devServer)
	command.StringFlag("frontenddevserverurl", "The url of the external frontend dev server to use", &flags.frontendDevServerURL)
	command.StringFlag("appargs", "arguments to pass to the underlying app (quoted and space separated)", &flags.appargs)
	command.BoolFlag("save", "Save given flags as defaults", &flags.saveConfig)
	command.BoolFlag("race", "Build with Go's race detector", &flags.raceDetector)
	command.BoolFlag("s", "Skips building the frontend", &flags.skipFrontend)

	command.Action(func() error {
		// Create logger
		logger := clilogger.New(w)
		app.PrintBanner()

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		projectConfig, err := loadAndMergeProjectConfig(cwd, &flags)
		if err != nil {
			return err
		}

		devServer := flags.devServer
		if _, _, err := net.SplitHostPort(devServer); err != nil {
			return fmt.Errorf("DevServer is not of the form 'host:port', please check your wails.json")
		}

		devServerURL, err := url.Parse("http://" + devServer)
		if err != nil {
			return err
		}

		// Update go.mod to use current wails version
		err = buildcmd.SyncGoMod(logger, true)
		if err != nil {
			return err
		}

		// Run go mod tidy to ensure we're up-to-date
		err = runCommand(cwd, false, "go", "mod", "tidy")
		if err != nil {
			return err
		}

		buildOptions := generateBuildOptions(flags)
		buildOptions.SkipBindings = flags.skipBindings
		buildOptions.Logger = logger

		userTags, err := buildtags.Parse(flags.tags)
		if err != nil {
			return err
		}

		buildOptions.UserTags = userTags

		if !buildOptions.SkipBindings {
			if flags.verbosity == build.VERBOSE {
				LogGreen("Generating Bindings...")
			}
			stdout, err := bindings.GenerateBindings(bindings.Options{
				Tags: buildOptions.UserTags,
			})
			if err != nil {
				return err
			}
			if flags.verbosity == build.VERBOSE {
				LogGreen(stdout)
			}
		}

		// Setup signal handler
		quitChannel := make(chan os.Signal, 1)
		signal.Notify(quitChannel, os.Interrupt, os.Kill, syscall.SIGTERM)
		exitCodeChannel := make(chan int, 1)

		// Build the frontend if requested, but ignore building the application itself.
		ignoreFrontend := buildOptions.IgnoreFrontend
		if !ignoreFrontend {
			logger.Println("Building frontend for development...")
			buildOptions.IgnoreApplication = true
			if _, err := build.Build(buildOptions); err != nil {
				return err
			}
			buildOptions.IgnoreApplication = false
		}

		// frontend:dev:watcher command.
		frontendDevAutoDiscovery := projectConfig.IsFrontendDevServerURLAutoDiscovery()
		if command := projectConfig.DevWatcherCommand; command != "" {
			closer, devServerURL, err := runFrontendDevWatcherCommand(cwd, command, frontendDevAutoDiscovery)
			if err != nil {
				return err
			}
			if devServerURL != "" {
				projectConfig.FrontendDevServerURL = devServerURL
				flags.frontendDevServerURL = devServerURL
			}
			defer closer()
		} else if frontendDevAutoDiscovery {
			return fmt.Errorf("Unable to auto discover frontend:dev:serverUrl without a frontend:dev:watcher command, please either set frontend:dev:watcher or remove the auto discovery from frontend:dev:serverUrl")
		}

		// Do initial build but only for the application.
		logger.Println("Building application for development...")
		buildOptions.IgnoreFrontend = true
		debugBinaryProcess, appBinary, err := restartApp(buildOptions, nil, flags, exitCodeChannel)
		buildOptions.IgnoreFrontend = ignoreFrontend || flags.frontendDevServerURL != ""
		if err != nil {
			return err
		}
		defer func() {
			if err := killProcessAndCleanupBinary(debugBinaryProcess, appBinary); err != nil {
				LogDarkYellow("Unable to kill process and cleanup binary: %s", err)
			}
		}()

		// open browser
		if flags.openBrowser {
			err = browser.OpenURL(devServerURL.String())
			if err != nil {
				return err
			}
		}

		// create the project files watcher
		watcher, err := initialiseWatcher(cwd)
		if err != nil {
			return err
		}

		defer func(watcher *fsnotify.Watcher) {
			err := watcher.Close()
			if err != nil {
				logger.Fatal(err.Error())
			}
		}(watcher)

		LogGreen("Watching (sub)/directory: %s", cwd)
		LogGreen("Using DevServer URL: %s", devServerURL)
		if flags.frontendDevServerURL != "" {
			LogGreen("Using Frontend DevServer URL: %s", flags.frontendDevServerURL)
		}
		LogGreen("Using reload debounce setting of %d milliseconds", flags.debounceMS)

		// Show dev server URL in terminal after 3 seconds
		go func() {
			time.Sleep(3 * time.Second)
			LogGreen("\n\nTo develop in the browser and call your bound Go methods from Javascript, navigate to: %s", devServerURL)
		}()

		// Watch for changes and trigger restartApp()
		debugBinaryProcess = doWatcherLoop(buildOptions, debugBinaryProcess, flags, watcher, exitCodeChannel, quitChannel, devServerURL)

		// Kill the current program if running and remove dev binary
		if err := killProcessAndCleanupBinary(debugBinaryProcess, appBinary); err != nil {
			return err
		}

		// Reset the process and the binary so defer knows about it and is a nop.
		debugBinaryProcess = nil
		appBinary = ""

		LogGreen("Development mode exited")

		return nil
	})
	return nil
}

func killProcessAndCleanupBinary(process *process.Process, binary string) error {
	if process != nil && process.Running {
		if err := process.Kill(); err != nil {
			return err
		}
	}

	if binary != "" {
		err := os.Remove(binary)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}

func runCommand(dir string, exitOnError bool, command string, args ...string) error {
	LogGreen("Executing: " + command + " " + strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		println(string(output))
		println(err.Error())
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
		Pack:           true,
		Platform:       runtime.GOOS,
		LDFlags:        flags.ldflags,
		Compiler:       flags.compilerCommand,
		ForceBuild:     flags.forceBuild,
		IgnoreFrontend: flags.skipFrontend,
		Verbosity:      flags.verbosity,
		WailsJSDir:     flags.wailsjsdir,
		RaceDetector:   flags.raceDetector,
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

	if flags.assetDir == "" && projectConfig.AssetDirectory != "" {
		flags.assetDir = projectConfig.AssetDirectory
	}

	if flags.assetDir != projectConfig.AssetDirectory {
		projectConfig.AssetDirectory = filepath.ToSlash(flags.assetDir)
	}

	if flags.assetDir != "" {
		flags.assetDir, err = filepath.Abs(flags.assetDir)
		if err != nil {
			return nil, err
		}
	}

	if flags.reloadDirs == "" && projectConfig.ReloadDirectories != "" {
		flags.reloadDirs = projectConfig.ReloadDirectories
	}

	if flags.reloadDirs != projectConfig.ReloadDirectories {
		projectConfig.ReloadDirectories = filepath.ToSlash(flags.reloadDirs)
	}

	if flags.devServer == "" && projectConfig.DevServer != "" {
		flags.devServer = projectConfig.DevServer
	}

	if flags.frontendDevServerURL == "" && projectConfig.FrontendDevServerURL != "" {
		flags.frontendDevServerURL = projectConfig.FrontendDevServerURL
	}

	if flags.wailsjsdir == "" && projectConfig.WailsJSDir != "" {
		flags.wailsjsdir = projectConfig.WailsJSDir
	}

	if flags.wailsjsdir == "" {
		flags.wailsjsdir = "./frontend"
	}

	if flags.wailsjsdir != projectConfig.WailsJSDir {
		projectConfig.WailsJSDir = filepath.ToSlash(flags.wailsjsdir)
	}

	if flags.debounceMS == 100 && projectConfig.DebounceMS != 100 {
		if projectConfig.DebounceMS == 0 {
			projectConfig.DebounceMS = 100
		}
		flags.debounceMS = projectConfig.DebounceMS
	}

	if flags.debounceMS != projectConfig.DebounceMS {
		projectConfig.DebounceMS = flags.debounceMS
	}

	if flags.appargs == "" && projectConfig.AppArgs != "" {
		flags.appargs = projectConfig.AppArgs
	}

	if flags.saveConfig {
		err = projectConfig.Save()
		if err != nil {
			return nil, err
		}
	}

	return projectConfig, nil
}

// runFrontendDevWatcherCommand will run the `frontend:dev:watcher` command if it was given, ex- `npm run dev`
func runFrontendDevWatcherCommand(cwd string, devCommand string, discoverViteServerURL bool) (func(), string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	scanner := NewStdoutScanner()
	dir := filepath.Join(cwd, "frontend")
	cmdSlice := strings.Split(devCommand, " ")
	cmd := exec.CommandContext(ctx, cmdSlice[0], cmdSlice[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = scanner
	cmd.Dir = dir
	setParentGID(cmd)

	if err := cmd.Start(); err != nil {
		cancel()
		return nil, "", fmt.Errorf("unable to start frontend DevWatcher: %w", err)
	}

	var viteServerURL string
	if discoverViteServerURL {
		select {
		case serverURL := <-scanner.ViteServerURLChan:
			viteServerURL = serverURL
		case <-time.After(time.Second * 10):
			cancel()
			return nil, "", errors.New("failed to find Vite server URL")
		}
	}

	LogGreen("Running frontend DevWatcher command: '%s'", devCommand)
	var wg sync.WaitGroup
	wg.Add(1)

	const (
		stateRunning   int32 = 0
		stateCanceling       = 1
		stateStopped         = 2
	)
	state := stateRunning
	go func() {
		if err := cmd.Wait(); err != nil {
			wasRunning := atomic.CompareAndSwapInt32(&state, stateRunning, stateStopped)
			if err.Error() != "exit status 1" && wasRunning {
				LogRed("Error from DevWatcher '%s': %s", devCommand, err.Error())
			}
		}
		atomic.StoreInt32(&state, stateStopped)
		wg.Done()
	}()

	return func() {
		if atomic.CompareAndSwapInt32(&state, stateRunning, stateCanceling) {
			killProc(cmd, devCommand)
		}
		cancel()
		wg.Wait()
	}, viteServerURL, nil
}

// restartApp does the actual rebuilding of the application when files change
func restartApp(buildOptions *build.Options, debugBinaryProcess *process.Process, flags devFlags, exitCodeChannel chan int) (*process.Process, string, error) {

	appBinary, err := build.Build(buildOptions)
	println()
	if err != nil {
		LogRed("Build error - " + err.Error())

		msg := "Continuing to run current version"
		if debugBinaryProcess == nil {
			msg = "No version running, build will be retriggered as soon as changes have been detected"
		}
		LogDarkYellow(msg)
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

	// parse appargs if any
	args, err := shlex.Split(flags.appargs)

	if err != nil {
		buildOptions.Logger.Fatal("Unable to parse appargs: %s", err.Error())
	}

	// Set environment variables accordingly
	os.Setenv("loglevel", flags.loglevel)
	os.Setenv("assetdir", flags.assetDir)
	os.Setenv("devserver", flags.devServer)
	os.Setenv("frontenddevserverurl", flags.frontendDevServerURL)

	// Start up new binary with correct args
	newProcess := process.NewProcess(appBinary, args...)
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
func doWatcherLoop(buildOptions *build.Options, debugBinaryProcess *process.Process, flags devFlags, watcher *fsnotify.Watcher, exitCodeChannel chan int, quitChannel chan os.Signal, devServerURL *url.URL) *process.Process {
	// Main Loop
	var extensionsThatTriggerARebuild = sliceToMap(strings.Split(flags.extensions, ","))
	var dirsThatTriggerAReload []string
	for _, dir := range strings.Split(flags.reloadDirs, ",") {
		if dir == "" {
			continue
		}
		thePath, err := filepath.Abs(dir)
		if err != nil {
			LogRed("Unable to expand reloadDir '%s': %s", dir, err)
			continue
		}
		dirsThatTriggerAReload = append(dirsThatTriggerAReload, thePath)
	}

	quit := false
	interval := time.Duration(flags.debounceMS) * time.Millisecond
	timer := time.NewTimer(interval)
	rebuild := false
	reload := false
	assetDir := ""
	changedPaths := map[string]struct{}{}

	assetDirURL := joinPath(devServerURL, "/wails/assetdir")
	reloadURL := joinPath(devServerURL, "/wails/reload")
	for quit == false {
		// reload := false
		select {
		case exitCode := <-exitCodeChannel:
			if exitCode == 0 {
				quit = true
			}
		case err := <-watcher.Errors:
			LogDarkYellow(err.Error())
		case item := <-watcher.Events:
			// Check for file writes
			if item.Op&fsnotify.Write == fsnotify.Write {
				// Ignore directories
				itemName := item.Name
				if fs.DirExists(itemName) {
					continue
				}

				// Iterate all file patterns
				ext := filepath.Ext(itemName)
				if ext != "" {
					ext = ext[1:]
					if _, exists := extensionsThatTriggerARebuild[ext]; exists {
						rebuild = true
						timer.Reset(interval)
						continue
					}
				}

				for _, reloadDir := range dirsThatTriggerAReload {
					if strings.HasPrefix(itemName, reloadDir) {
						reload = true
						break
					}
				}

				if !reload {
					changedPaths[filepath.Dir(itemName)] = struct{}{}
				}

				timer.Reset(interval)
			}
			// Check for new directories
			if item.Op&fsnotify.Create == fsnotify.Create {
				// If this is a folder, add it to our watch list
				if fs.DirExists(item.Name) {
					// node_modules is BANNED!
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
				newBinaryProcess, _, err := restartApp(buildOptions, debugBinaryProcess, flags, exitCodeChannel)
				if err != nil {
					LogRed("Error during build: %s", err.Error())
					continue
				}
				// If we have a new process, saveConfig it
				if newBinaryProcess != nil {
					debugBinaryProcess = newBinaryProcess
				}
			}

			if flags.frontendDevServerURL != "" {
				// If we are using an external dev server, the reloading of the frontend part can be skipped
				continue
			}
			if len(changedPaths) != 0 {
				if assetDir == "" {
					resp, err := http.Get(assetDirURL)
					if err != nil {
						LogRed("Error during retrieving assetdir: %s", err.Error())
					} else {
						content, err := io.ReadAll(resp.Body)
						if err != nil {
							LogRed("Error reading assetdir from devserver: %s", err.Error())
						} else {
							assetDir = string(content)
						}
						resp.Body.Close()
					}
				}

				if assetDir != "" {
					for thePath := range changedPaths {
						if strings.HasPrefix(thePath, assetDir) {
							reload = true
							break
						}
					}
				} else if len(dirsThatTriggerAReload) == 0 {
					LogRed("Reloading couldn't be triggered: Please specify -assetdir or -reloaddirs")
				}

				changedPaths = map[string]struct{}{}
			}
			if reload {
				reload = false
				_, err := http.Get(reloadURL)
				if err != nil {
					LogRed("Error during refresh: %s", err.Error())
				}
			}
		case <-quitChannel:
			LogGreen("\nCaught quit")
			quit = true
		}
	}
	return debugBinaryProcess
}

func joinPath(url *url.URL, subPath string) string {
	u := *url
	u.Path = path.Join(u.Path, subPath)
	return u.String()
}
