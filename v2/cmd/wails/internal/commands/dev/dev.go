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
	"syscall"
	"time"

	"github.com/google/shlex"
	"github.com/wailsapp/wails/v2/cmd/wails/internal"
	"github.com/wailsapp/wails/v2/internal/gomod"

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
	reloadDirs      string
	openBrowser     bool
	noReload        bool
	noGen           bool
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
	command.BoolFlag("nogen", "Disable generate module", &flags.noGen)
	command.StringFlag("wailsjsdir", "Directory to generate the Wails JS modules", &flags.wailsjsdir)
	command.StringFlag("tags", "tags to pass to Go compiler (quoted and space separated)", &flags.tags)
	command.IntFlag("v", "Verbosity level (0 - silent, 1 - standard, 2 - verbose)", &flags.verbosity)
	command.StringFlag("loglevel", "Loglevel to use - Trace, Debug, Info, Warning, Error", &flags.loglevel)
	command.BoolFlag("f", "Force build application", &flags.forceBuild)
	command.IntFlag("debounce", "The amount of time to wait to trigger a reload on change", &flags.debounceMS)
	command.StringFlag("devserver", "The address of the wails dev server", &flags.devServer)
	command.StringFlag("frontenddevserverurl", "The url of the external frontend dev server to use", &flags.frontendDevServerURL)
	command.StringFlag("appargs", "arguments to pass to the underlying app (quoted and space searated)", &flags.appargs)
	command.BoolFlag("save", "Save given flags as defaults", &flags.saveConfig)
	command.BoolFlag("race", "Build with Go's race detector", &flags.raceDetector)

	command.Action(func() error {
		// Create logger
		logger := clilogger.New(w)
		app.PrintBanner()

		userTags := []string{}
		for _, tag := range strings.Split(flags.tags, " ") {
			thisTag := strings.TrimSpace(tag)
			if thisTag != "" {
				userTags = append(userTags, thisTag)
			}
		}

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
		err = syncGoModVersion(cwd)
		if err != nil {
			return err
		}

		// Run go mod tidy to ensure we're up to date
		err = runCommand(cwd, false, "go", "mod", "tidy", "-compat=1.17")
		if err != nil {
			return err
		}

		if !flags.noGen {
			self := os.Args[0]
			if flags.tags != "" {
				err = runCommand(cwd, true, self, "generate", "module", "-tags", flags.tags)
			} else {
				err = runCommand(cwd, true, self, "generate", "module")
			}
			if err != nil {
				return err
			}
		}

		buildOptions := generateBuildOptions(flags)
		buildOptions.Logger = logger
		buildOptions.UserTags = internal.ParseUserTags(flags.tags)

		// Setup signal handler
		quitChannel := make(chan os.Signal, 1)
		signal.Notify(quitChannel, os.Interrupt, os.Kill, syscall.SIGTERM)
		exitCodeChannel := make(chan int, 1)

		// Do initial build
		logger.Println("Building application for development...")
		debugBinaryProcess, appBinary, err := restartApp(buildOptions, nil, flags, exitCodeChannel)
		if err != nil {
			return err
		}
		defer func() {
			if err := killProcessAndCleanupBinary(debugBinaryProcess, appBinary); err != nil {
				LogDarkYellow("Unable to kill process and cleanup binary: %s", err)
			}
		}()

		// frontend:dev:watcher command.
		if command := projectConfig.DevWatcherCommand; command != "" {
			closer, err := runFrontendDevWatcherCommand(cwd, command)
			if err != nil {
				return err
			}
			defer closer()
		}

		// open browser
		if flags.openBrowser {
			err = browser.OpenURL(devServerURL.String())
			if err != nil {
				return err
			}
		}

		// create the project files watcher
		watcher, err := initialiseWatcher(cwd, logger.Fatal)
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

		// Watch for changes and trigger restartApp()
		debugBinaryProcess = doWatcherLoop(buildOptions, debugBinaryProcess, flags, watcher, exitCodeChannel, quitChannel, devServerURL)

		// Kill the current program if running and remove dev binary
		if err := killProcessAndCleanupBinary(debugBinaryProcess, appBinary); err != nil {
			return err
		}

		// Reset the process and the binary so the defer knows about it and is a nop.
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
		IgnoreFrontend: false,
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
func runFrontendDevWatcherCommand(cwd string, devCommand string) (func(), error) {
	ctx, cancel := context.WithCancel(context.Background())
	dir := filepath.Join(cwd, "frontend")
	cmdSlice := strings.Split(devCommand, " ")
	cmd := exec.CommandContext(ctx, cmdSlice[0], cmdSlice[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = dir
	setParentGID(cmd)
	if err := cmd.Start(); err != nil {
		cancel()
		return nil, fmt.Errorf("Unable to start frontend DevWatcher: %w", err)
	}

	LogGreen("Running frontend DevWatcher command: '%s'", devCommand)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		if err := cmd.Wait(); err != nil {
			if err.Error() != "exit status 1" {
				LogRed("Error from DevWatcher '%s': %s", devCommand, err.Error())
			}
		}
		LogGreen("DevWatcher command exited!")
		wg.Done()
	}()

	return func() {
		killProc(cmd, devCommand)
		LogGreen("DevWatcher command killed!")
		cancel()
		wg.Wait()
	}, nil
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
		if strings.Contains(dir, ".git") {
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
		path, err := filepath.Abs(dir)
		if err != nil {
			LogRed("Unable to expand reloadDir '%s': %s", dir, err)
			continue
		}
		dirsThatTriggerAReload = append(dirsThatTriggerAReload, path)
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
		//reload := false
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
				// If we are using an external dev server all the reload of the frontend part can be skipped
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
					for path := range changedPaths {
						if strings.HasPrefix(path, assetDir) {
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
