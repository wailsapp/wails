package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/go-git/go-billy/v5/osfs"
	gitignore "github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/rjeczalik/notify"
	"github.com/wailsapp/wails/v3/internal/wake"
	"github.com/wailsapp/wails/v3/internal/wake/cache"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

type manifestProcess struct {
	cmd  *exec.Cmd
	done chan struct{}
	mu   sync.Mutex
	err  error
}

type manifestRebuild struct {
	generation uint64
	loaded     *manifest.Loaded
	run        manifestPipelineRun
	err        error
}

type manifestWatchSet struct {
	events  chan notify.EventInfo
	ignored gitignore.Matcher
	matcher devWatchMatcher
	root    string
}

type devWatchMatcher struct {
	matchAll bool
	patterns [][]devWatchSegment
}

const manifestBackendReadinessDelay = 500 * time.Millisecond

type devWatchSegment struct {
	value      string
	glob       bool
	doubleStar bool
}

// manifestDevOps is the private adapter seam between the Dev session state
// machine and local operating-system effects. The production adapter below is
// used by every caller; tests substitute only effects the host cannot produce
// deterministically, such as a watcher failing during reconfiguration.
type manifestDevOps struct {
	getwd           func() (string, error)
	load            func(string, string) (*manifest.Loaded, error)
	checkPort       func(string, int) error
	build           func(context.Context, *DevOptions, *manifest.Loaded, string, string, string, int) (manifestPipelineRun, error)
	startFrontend   func(string, manifest.Config, string, int, string) (*manifestProcess, error)
	waitTCP         func(context.Context, *manifestProcess, string, time.Duration) error
	binaryPath      func(string, manifestPipelineRun, string, string) (string, error)
	startApp        func(string, string, string, int) (*manifestProcess, error)
	waitStable      func(context.Context, *manifestProcess, time.Duration) error
	startWatches    func(string, manifest.Config) (*manifestWatchSet, error)
	restartWatches  func(string, manifest.Config, *manifestWatchSet) (*manifestWatchSet, error)
	restoreFrontend func(context.Context, string, manifest.Config, int, string, string) (*manifestProcess, error)
}

func productionManifestDevOps() manifestDevOps {
	return manifestDevOps{
		getwd:           os.Getwd,
		load:            manifest.Load,
		checkPort:       checkManifestDevPort,
		build:           runManifestDevBuild,
		startFrontend:   startFrontendDev,
		waitTCP:         waitForProcessTCP,
		binaryPath:      manifestDevBinaryPath,
		startApp:        startManifestApp,
		waitStable:      waitForProcessStable,
		startWatches:    startManifestWatches,
		restartWatches:  restartManifestWatches,
		restoreFrontend: restoreManifestFrontend,
	}
}

func checkManifestDevPort(host string, port int) error {
	listener, err := net.Listen("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		return fmt.Errorf("frontend port %d is unavailable: %w", port, err)
	}
	return listener.Close()
}

func runManifestDev(options *DevOptions) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return runManifestDevContext(ctx, options)
}

func runManifestDevContext(ctx context.Context, options *DevOptions) error {
	return runManifestDevContextWithOps(ctx, options, productionManifestDevOps())
}

func runManifestDevContextWithOps(ctx context.Context, options *DevOptions, ops manifestDevOps) error {
	if options.Profile != "" {
		return fmt.Errorf("dev does not accept production profiles; configure development policy in the dev block")
	}
	sessionCtx, stopSession := context.WithCancel(ctx)
	defer stopSession()
	root, err := ops.getwd()
	if err != nil {
		return err
	}
	loaded, err := ops.load(root, "")
	if err != nil {
		return err
	}
	goos, goarch, err := splitTarget(options.Target)
	if err != nil {
		return err
	}
	if goos == "" {
		goos, goarch = runtime.GOOS, runtime.GOARCH
	}
	if goos != runtime.GOOS || goarch != runtime.GOARCH {
		return fmt.Errorf("dev target must be the host %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	if options.Plan {
		return printManifestPlan(manifestRunOptions{Verb: "build", Loaded: loaded, TargetOS: goos, TargetArch: goarch, Development: true, Tags: manifestDevTags(options)}, false)
	}
	port := options.VitePort
	if port == 0 {
		port = defaultVitePort
	}
	host := "127.0.0.1"
	if err := ops.checkPort(host, port); err != nil {
		return err
	}
	scheme := "http"
	if options.Secure {
		scheme = "https"
	}
	frontendURL := fmt.Sprintf("%s://%s:%d", scheme, host, port)

	initialRun, err := ops.build(sessionCtx, options, loaded, goos, goarch, frontendURL, port)
	if err != nil {
		if errors.Is(err, context.Canceled) && sessionCtx.Err() != nil {
			return nil
		}
		return err
	}
	frontend, err := ops.startFrontend(root, loaded.Config, host, port, frontendURL)
	if err != nil {
		return err
	}
	if err := ops.waitTCP(sessionCtx, frontend, net.JoinHostPort(host, strconv.Itoa(port)), 30*time.Second); err != nil {
		frontend.stop(time.Duration(loaded.Config.Dev.GracePeriodMS) * time.Millisecond)
		if errors.Is(err, context.Canceled) && sessionCtx.Err() != nil {
			return nil
		}
		return fmt.Errorf("frontend readiness: %w", err)
	}
	defer func() {
		if frontend != nil {
			frontend.stop(time.Duration(loaded.Config.Dev.GracePeriodMS) * time.Millisecond)
		}
	}()
	fmt.Printf("Frontend process started at %s\n", frontendURL)
	binaryPath, err := ops.binaryPath(root, initialRun, goos, goarch)
	if err != nil {
		return err
	}
	app, err := ops.startApp(root, binaryPath, frontendURL, port)
	if err != nil {
		return err
	}
	if err := ops.waitStable(sessionCtx, app, manifestBackendReadinessDelay); err != nil {
		app.stop(time.Duration(loaded.Config.Dev.GracePeriodMS) * time.Millisecond)
		if errors.Is(err, context.Canceled) && sessionCtx.Err() != nil {
			return nil
		}
		return fmt.Errorf("backend readiness: %w", err)
	}
	defer func() {
		if app != nil {
			app.stop(time.Duration(loaded.Config.Dev.GracePeriodMS) * time.Millisecond)
		}
	}()
	watches, err := ops.startWatches(root, loaded.Config)
	if err != nil {
		return err
	}
	defer func() { watches.stop() }()
	fmt.Println("Backend built and started")
	debounce := time.Duration(loaded.Config.Dev.DebounceMS) * time.Millisecond
	if debounce <= 0 {
		debounce = 250 * time.Millisecond
	}
	var timer *time.Timer
	var timerC <-chan time.Time
	results := make(chan manifestRebuild, 1)
	var generation uint64
	var cancelBuild context.CancelFunc
	var buildMu sync.Mutex
	var buildWG sync.WaitGroup
	defer func() {
		stopSession()
		if cancelBuild != nil {
			cancelBuild()
		}
		buildWG.Wait()
	}()
	startRebuild := func() {
		if cancelBuild != nil {
			cancelBuild()
		}
		generation++
		current := generation
		buildCtx, cancel := context.WithCancel(sessionCtx)
		cancelBuild = cancel
		buildWG.Add(1)
		go func() {
			defer buildWG.Done()
			buildMu.Lock()
			defer buildMu.Unlock()
			if buildErr := buildCtx.Err(); buildErr != nil {
				select {
				case results <- manifestRebuild{generation: current, err: buildErr}:
				case <-sessionCtx.Done():
				}
				return
			}
			nextLoaded, loadErr := ops.load(root, "")
			var run manifestPipelineRun
			if loadErr == nil {
				run, loadErr = ops.build(buildCtx, options, nextLoaded, goos, goarch, frontendURL, port)
			}
			result := manifestRebuild{generation: current, loaded: nextLoaded, run: run, err: loadErr}
			select {
			case results <- result:
			case <-sessionCtx.Done():
			}
		}()
	}
	for {
		select {
		case <-sessionCtx.Done():
			if timer != nil {
				timer.Stop()
			}
			return nil
		case <-frontend.done:
			return manifestSessionProcessExit(sessionCtx, "frontend", frontend)
		case <-app.done:
			return manifestSessionProcessExit(sessionCtx, "backend", app)
		case event, ok := <-watches.events:
			if !ok {
				return fmt.Errorf("project watcher stopped unexpectedly")
			}
			createdDirectory := isCreatedDevDirectory(event)
			if createdDirectory && !shouldRefreshDevWatchesForDirectory(watches.root, loaded.Config, watches.ignored, event.Path()) {
				continue
			}
			if createdDirectory {
				nextWatches, watchErr := ops.restartWatches(root, loaded.Config, watches)
				if nextWatches != nil {
					watches = nextWatches
				}
				if watchErr != nil {
					fmt.Fprintf(os.Stderr, "new directory watch failed; restored the previous watch policy: %v\n", watchErr)
					continue
				}
				if !directoryContainsDevInputMatched(watches.root, loaded.Config, watches.ignored, watches.matcher, event.Path()) {
					continue
				}
			} else if isDevGitIgnoreEvent(watches.root, loaded.Config, event.Path()) {
				nextWatches, watchErr := ops.restartWatches(root, loaded.Config, watches)
				if nextWatches != nil {
					watches = nextWatches
				}
				if watchErr != nil {
					fmt.Fprintf(os.Stderr, "gitignore reload failed; restored the previous watch policy: %v\n", watchErr)
					continue
				}
				fmt.Println("Watch policy reloaded")
				continue
			} else if ignoreDevEventMatched(watches.root, loaded.Config, watches.ignored, watches.matcher, event.Path()) {
				continue
			}
			if timer == nil {
				timer = time.NewTimer(debounce)
			} else {
				timer.Reset(debounce)
			}
			timerC = timer.C
		case <-timerC:
			timerC = nil
			startRebuild()
		case result := <-results:
			if result.generation != generation {
				continue
			}
			cancelBuild()
			cancelBuild = nil
			if result.err != nil {
				if !errors.Is(result.err, context.Canceled) {
					reportManifestRebuildFailure(result.err)
				}
				continue
			}
			nextLoaded := result.loaded
			var nextWatches *manifestWatchSet
			if devWatchSessionChanged(loaded.Config, nextLoaded.Config) {
				nextWatches, err = ops.startWatches(root, nextLoaded.Config)
				if err != nil {
					fmt.Fprintf(os.Stderr, "watch reconfiguration failed; keeping the current session: %v\n", err)
					continue
				}
			}

			frontendChanged := frontendSessionChanged(loaded.Config, nextLoaded.Config)
			backendChanged := manifestBackendChanged(result.run, goos, goarch)
			oldFrontend := frontend
			var nextFrontend *manifestProcess
			if frontendChanged {
				oldFrontend.stop(time.Duration(loaded.Config.Dev.GracePeriodMS) * time.Millisecond)
				nextFrontend, err = ops.startFrontend(root, nextLoaded.Config, host, port, frontendURL)
				if err == nil {
					err = ops.waitTCP(sessionCtx, nextFrontend, net.JoinHostPort(host, strconv.Itoa(port)), 30*time.Second)
				}
				if err != nil {
					transitionErr := err
					if nextFrontend != nil {
						nextFrontend.stop(time.Duration(nextLoaded.Config.Dev.GracePeriodMS) * time.Millisecond)
					}
					if nextWatches != nil {
						nextWatches.stop()
					}
					if errors.Is(transitionErr, context.Canceled) && sessionCtx.Err() != nil {
						return nil
					}
					frontend, err = ops.restoreFrontend(sessionCtx, root, loaded.Config, port, frontendURL, host)
					if err != nil {
						return err
					}
					fmt.Fprintf(os.Stderr, "frontend restart failed; keeping the current session: %v\n", transitionErr)
					continue
				}
			}

			var nextApp *manifestProcess
			if backendChanged {
				binaryPath, pathErr := ops.binaryPath(root, result.run, goos, goarch)
				if pathErr != nil {
					err = pathErr
				} else {
					nextApp, err = ops.startApp(root, binaryPath, frontendURL, port)
				}
				if err == nil {
					err = ops.waitStable(sessionCtx, nextApp, manifestBackendReadinessDelay)
				}
				if err != nil {
					transitionErr := err
					if nextApp != nil {
						nextApp.stop(time.Duration(nextLoaded.Config.Dev.GracePeriodMS) * time.Millisecond)
					}
					if nextWatches != nil {
						nextWatches.stop()
					}
					if nextFrontend != nil {
						nextFrontend.stop(time.Duration(nextLoaded.Config.Dev.GracePeriodMS) * time.Millisecond)
						frontend, err = ops.restoreFrontend(sessionCtx, root, loaded.Config, port, frontendURL, host)
						if err != nil {
							return err
						}
					}
					if errors.Is(transitionErr, context.Canceled) && sessionCtx.Err() != nil {
						return nil
					}
					fmt.Fprintf(os.Stderr, "restart failed; keeping the current app running: %v\n", transitionErr)
					continue
				}
			}

			if backendChanged {
				oldApp := app
				app = nextApp
				oldApp.stop(time.Duration(loaded.Config.Dev.GracePeriodMS) * time.Millisecond)
			}
			if frontendChanged {
				frontend = nextFrontend
				fmt.Printf("Frontend process restarted at %s\n", frontendURL)
			}
			if nextWatches != nil {
				oldWatches := watches
				watches = nextWatches
				oldWatches.stop()
			}
			loaded = nextLoaded
			debounce = time.Duration(loaded.Config.Dev.DebounceMS) * time.Millisecond
			if debounce <= 0 {
				debounce = 250 * time.Millisecond
			}
			switch {
			case backendChanged:
				fmt.Println("Backend rebuilt and restarted")
			case frontendChanged:
				fmt.Println("Frontend session updated; backend unchanged")
			default:
				fmt.Println("Build is current; backend unchanged")
			}
		}
	}
}

func runManifestDevBuild(ctx context.Context, options *DevOptions, loaded *manifest.Loaded, goos, goarch, frontendURL string, port int) (manifestPipelineRun, error) {
	environment := []string{
		"FRONTEND_DEVSERVER_URL=" + frontendURL,
		wailsVitePort + "=" + strconv.Itoa(port),
	}
	return runManifestPipelineResult(manifestRunOptions{Context: ctx, Verb: "build", Loaded: loaded, TargetOS: goos, TargetArch: goarch, Environment: environment, Development: true, Tags: manifestDevTags(options)})
}

func manifestDevTags(options *DevOptions) []string {
	return appendUniqueStrings(splitComma(options.Tags), envTags()...)
}

func manifestBackendChanged(run manifestPipelineRun, goos, goarch string) bool {
	result, exists := run.Results[pipelineCompileKey(goos, goarch)]
	if !exists || result.Status == cache.LookupMiss {
		return true
	}
	hook, exists := run.Results[pipelineAfterBuildHookKey(goos, goarch)]
	return exists && hook.Status == cache.LookupMiss
}

func pipelineCompileKey(goos, goarch string) pipeline.NodeKey {
	return pipeline.NodeKey("target:" + goos + "/" + goarch + ":compile")
}

func pipelineAfterBuildHookKey(goos, goarch string) pipeline.NodeKey {
	return pipeline.NodeKey("hook:after_build:" + goos + "-" + goarch)
}

func devWatchSessionChanged(current, next manifest.Config) bool {
	return current.Dev.UseGitIgnore != next.Dev.UseGitIgnore ||
		current.Frontend.Directory != next.Frontend.Directory ||
		!equalStrings(current.Dev.Watch, next.Dev.Watch) ||
		!equalStrings(current.Dev.Exclude, next.Dev.Exclude)
}

func equalStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func isDevGitIgnoreEvent(root string, config manifest.Config, eventPath string) bool {
	if !config.Dev.UseGitIgnore {
		return false
	}
	rel, err := filepath.Rel(root, eventPath)
	if err != nil {
		return false
	}
	rel = filepath.ToSlash(rel)
	return rel == ".gitignore" || strings.HasSuffix(rel, "/.gitignore")
}

func isCreatedDevDirectory(event notify.EventInfo) bool {
	if event == nil || event.Event()&notify.Create == 0 {
		return false
	}
	info, err := os.Stat(event.Path())
	return err == nil && info.IsDir()
}

func directoryContainsDevInput(root string, config manifest.Config, ignored gitignore.Matcher, directory string) bool {
	return directoryContainsDevInputMatched(root, config, ignored, newDevWatchMatcher(config.Dev.Watch), directory)
}

func directoryContainsDevInputMatched(root string, config manifest.Config, ignored gitignore.Matcher, matcher devWatchMatcher, directory string) bool {
	found := false
	_ = filepath.WalkDir(directory, func(current string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return filepath.SkipAll
		}
		if entry.IsDir() {
			rel, err := filepath.Rel(root, current)
			if err != nil {
				return filepath.SkipDir
			}
			rel = filepath.ToSlash(rel)
			if rel != "." && (devPathExcluded(config, rel) || ignored != nil && ignored.Match(strings.Split(rel, "/"), true)) {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Base(current) == ".gitignore" {
			return nil
		}
		if !ignoreDevEventMatched(root, config, ignored, matcher, current) {
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	return found
}

func startManifestWatches(root string, config manifest.Config) (*manifestWatchSet, error) {
	canonicalRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return nil, fmt.Errorf("resolve project root: %w", err)
	}
	ignored, err := loadDevGitIgnore(root, config.Dev.UseGitIgnore)
	if err != nil {
		return nil, fmt.Errorf("load .gitignore: %w", err)
	}
	events := make(chan notify.EventInfo, 1024)
	if err := registerManifestWatches(root, config, ignored, events); err != nil {
		notify.Stop(events)
		return nil, fmt.Errorf("watch project: %w", err)
	}
	return &manifestWatchSet{events: events, ignored: ignored, matcher: newDevWatchMatcher(config.Dev.Watch), root: canonicalRoot}, nil
}

func (w *manifestWatchSet) stop() {
	if w != nil && w.events != nil {
		notify.Stop(w.events)
	}
}

func restartManifestWatches(root string, config manifest.Config, current *manifestWatchSet) (*manifestWatchSet, error) {
	next, err := startManifestWatches(root, config)
	if err != nil {
		return current, err
	}
	current.stop()
	return next, nil
}

func restoreManifestFrontend(ctx context.Context, root string, config manifest.Config, port int, frontendURL, host string) (*manifestProcess, error) {
	frontend, err := startFrontendDev(root, config, host, port, frontendURL)
	if err != nil {
		return nil, fmt.Errorf("restore previous frontend: %w", err)
	}
	if err := waitForProcessTCP(ctx, frontend, net.JoinHostPort(host, strconv.Itoa(port)), 30*time.Second); err != nil {
		frontend.stop(time.Duration(config.Dev.GracePeriodMS) * time.Millisecond)
		return nil, fmt.Errorf("restore previous frontend readiness: %w", err)
	}
	return frontend, nil
}

func manifestProcessExitError(name string, process *manifestProcess) error {
	if err := process.waitError(); err != nil {
		return fmt.Errorf("%s process exited unexpectedly: %w", name, err)
	}
	return fmt.Errorf("%s process exited unexpectedly", name)
}

func manifestSessionProcessExit(ctx context.Context, name string, process *manifestProcess) error {
	if ctx.Err() != nil {
		return nil
	}
	return manifestProcessExitError(name, process)
}

func reportManifestRebuildFailure(err error) {
	if wake.IsReported(err) {
		fmt.Fprintln(os.Stderr, "Current app is still running.")
		return
	}
	fmt.Fprintf(os.Stderr, "Rebuild failed; current app is still running: %v\n", err)
}

func frontendSessionChanged(current, next manifest.Config) bool {
	return current.Frontend.Directory != next.Frontend.Directory || current.Frontend.PackageManager != next.Frontend.PackageManager || current.Frontend.DevCommand != next.Frontend.DevCommand || !equalStrings(current.Frontend.Dev, next.Frontend.Dev)
}

func startFrontendDev(root string, config manifest.Config, host string, port int, frontendURL string) (*manifestProcess, error) {
	if len(config.Frontend.Dev) > 0 {
		args := append([]string(nil), config.Frontend.Dev...)
		env := []string{wailsVitePort + "=" + strconv.Itoa(port), "FRONTEND_DEVSERVER_URL=" + frontendURL}
		return startManifestProcess(filepath.Join(root, config.Frontend.Directory), args[0], env, args[1:]...)
	}
	manager := config.Frontend.PackageManager
	if config.Frontend.DevCommand == "" {
		return nil, fmt.Errorf("frontend.dev_command is not set in %s", manifest.Filename)
	}
	var serverArgs []string
	if frontendDevCommandUsesVite(root, config) {
		serverArgs = []string{"--host", host, "--port", strconv.Itoa(port), "--strictPort"}
	}
	args, err := frontendDevArgs(manager, config.Frontend.DevCommand, serverArgs)
	if err != nil {
		return nil, err
	}
	env := []string{wailsVitePort + "=" + strconv.Itoa(port), "FRONTEND_DEVSERVER_URL=" + frontendURL}
	return startPackageManagerProcess(filepath.Join(root, config.Frontend.Directory), manager, env, args...)
}

func frontendDevArgs(manager, command string, serverArgs []string) ([]string, error) {
	var args []string
	switch manager {
	case "npm":
		args = []string{"run", command}
		if len(serverArgs) > 0 {
			args = append(args, "--")
		}
		args = append(args, serverArgs...)
	case "pnpm", "bun":
		args = append([]string{"run", command}, serverArgs...)
	case "yarn":
		args = append([]string{command}, serverArgs...)
	default:
		return nil, fmt.Errorf("unsupported frontend.package_manager %q", manager)
	}
	return args, nil
}

func frontendDevCommandUsesVite(root string, config manifest.Config) bool {
	data, err := os.ReadFile(filepath.Join(root, config.Frontend.Directory, "package.json"))
	if err != nil {
		return false
	}
	var packageJSON struct {
		Scripts map[string]string `json:"scripts"`
	}
	if json.Unmarshal(data, &packageJSON) != nil {
		return false
	}
	command := strings.NewReplacer("&&", " ", "||", " ", ";", " ", "|", " ").Replace(packageJSON.Scripts[config.Frontend.DevCommand])
	for _, field := range strings.Fields(command) {
		field = strings.Trim(strings.ReplaceAll(field, `\`, "/"), `"'`)
		base := path.Base(field)
		if base == "vite" || base == "vite.cmd" {
			return true
		}
	}
	return false
}
func startPackageManagerProcess(dir, manager string, env []string, args ...string) (*manifestProcess, error) {
	name, resolvedArgs, err := resolvePackageManagerProcess(manager, args, exec.LookPath, filepath.EvalSymlinks)
	if err != nil {
		return nil, err
	}
	return startManifestProcess(dir, name, env, resolvedArgs...)
}

func resolvePackageManagerProcess(manager string, args []string, lookPath func(string) (string, error), evalSymlinks func(string) (string, error)) (string, []string, error) {
	name := manager
	if manager == "npm" {
		npm, err := lookPath("npm")
		if err != nil {
			return "", nil, err
		}
		switch strings.ToLower(filepath.Ext(npm)) {
		case ".cmd", ".bat":
			commandInterpreter, err := lookPath("cmd")
			if err != nil {
				return "", nil, err
			}
			return commandInterpreter, append([]string{"/d", "/s", "/c", npm}, args...), nil
		}
		script, err := evalSymlinks(npm)
		if err != nil {
			script = npm
		}
		node, err := lookPath("node")
		if err != nil {
			return "", nil, err
		}
		name = node
		args = append([]string{script}, args...)
	}
	return name, args, nil
}
func manifestDevBinaryPath(root string, run manifestPipelineRun, goos, goarch string) (string, error) {
	node, exists := run.Plan.Nodes[pipelineCompileKey(goos, goarch)]
	if !exists || node.Output == "" {
		return "", fmt.Errorf("development Plan for %s/%s has no compile output", goos, goarch)
	}
	return filepath.Join(root, filepath.FromSlash(node.Output)), nil
}

func startManifestApp(root, binaryPath, frontendURL string, port int) (*manifestProcess, error) {
	env := []string{wailsVitePort + "=" + strconv.Itoa(port), "FRONTEND_DEVSERVER_URL=" + frontendURL}
	return startManifestProcess(root, binaryPath, env)
}
func startManifestProcess(dir, name string, env []string, args ...string) (*manifestProcess, error) {
	cmd := exec.CommandContext(context.Background(), name, args...)
	cmd.Dir = dir
	cmd.Env = mergeEnvironment(os.Environ(), env)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	configureManifestProcess(cmd)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	result := &manifestProcess{cmd: cmd, done: make(chan struct{})}
	go func() {
		waitErr := cmd.Wait()
		result.mu.Lock()
		result.err = waitErr
		result.mu.Unlock()
		close(result.done)
	}()
	return result, nil
}

func (p *manifestProcess) waitError() error {
	if p == nil {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.err
}

func (p *manifestProcess) stop(grace time.Duration) {
	if p == nil || p.cmd == nil || p.cmd.Process == nil {
		return
	}
	select {
	case <-p.done:
		return
	default:
	}
	_ = signalManifestProcess(p.cmd.Process, os.Interrupt)
	if grace <= 0 {
		grace = 1500 * time.Millisecond
	}
	select {
	case <-p.done:
		return
	case <-time.After(grace):
		_ = killManifestProcess(p.cmd.Process)
		<-p.done
	}
}

func waitForProcessTCP(ctx context.Context, process *manifestProcess, address string, timeout time.Duration) error {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	deadline := time.NewTimer(timeout)
	defer deadline.Stop()
	ticker := time.NewTicker(25 * time.Millisecond)
	defer ticker.Stop()
	for {
		connection, err := net.DialTimeout("tcp", address, 100*time.Millisecond)
		if err == nil {
			_ = connection.Close()
			return nil
		}
		select {
		case <-process.done:
			processErr := process.waitError()
			if processErr == nil {
				return fmt.Errorf("process exited before becoming ready")
			}
			return fmt.Errorf("process exited before becoming ready: %w", processErr)
		case <-ctx.Done():
			return ctx.Err()
		case <-deadline.C:
			return fmt.Errorf("timed out waiting for %s", address)
		case <-ticker.C:
		}
	}
}

func waitForProcessStable(ctx context.Context, process *manifestProcess, duration time.Duration) error {
	if duration <= 0 {
		duration = 150 * time.Millisecond
	}
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-process.done:
		processErr := process.waitError()
		if processErr == nil {
			return fmt.Errorf("process exited during startup")
		}
		return fmt.Errorf("process exited during startup: %w", processErr)
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
func ignoreDevEvent(root string, config manifest.Config, ignored gitignore.Matcher, path string) bool {
	return ignoreDevEventMatched(root, config, ignored, newDevWatchMatcher(config.Dev.Watch), path)
}

func ignoreDevEventMatched(root string, config manifest.Config, ignored gitignore.Matcher, matcher devWatchMatcher, path string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return true
	}
	rel = filepath.ToSlash(rel)
	if rel == "." || devPathExcluded(config, rel) {
		return true
	}
	if config.Dev.UseGitIgnore && (rel == ".gitignore" || strings.HasSuffix(rel, "/.gitignore")) {
		return false
	}
	if ignored != nil && ignored.Match(strings.Split(rel, "/"), false) {
		return true
	}
	return strings.HasSuffix(rel, "_test.go") || !matcher.Match(rel)
}

func registerManifestWatches(root string, config manifest.Config, ignored gitignore.Matcher, events chan<- notify.EventInfo) error {
	return filepath.WalkDir(root, func(current string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, current)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel != "." && devPathExcluded(config, rel) {
			return filepath.SkipDir
		}
		if rel != "." && ignored != nil && ignored.Match(strings.Split(rel, "/"), true) {
			return filepath.SkipDir
		}
		return notify.Watch(current, events, notify.All)
	})
}

func loadDevGitIgnore(root string, enabled bool) (gitignore.Matcher, error) {
	if !enabled {
		return nil, nil
	}
	patterns, err := gitignore.ReadPatterns(osfs.New(root), nil)
	if err != nil {
		return nil, err
	}
	return gitignore.NewMatcher(patterns), nil
}

func devPathExcluded(config manifest.Config, rel string) bool {
	rel = strings.Trim(filepath.ToSlash(rel), "/")
	for _, segment := range strings.Split(rel, "/") {
		if segment == ".git" || segment == ".wails" || segment == "node_modules" {
			return true
		}
	}
	paths := []string{config.Frontend.Directory}
	paths = append(paths, config.Dev.Exclude...)
	for _, excluded := range paths {
		clean := strings.Trim(filepath.ToSlash(excluded), "/")
		if clean != "" && (rel == clean || strings.HasPrefix(rel, clean+"/")) {
			return true
		}
	}
	return false
}

func shouldRefreshDevWatchesForDirectory(root string, config manifest.Config, ignored gitignore.Matcher, directory string) bool {
	relative, err := filepath.Rel(root, directory)
	if err != nil {
		return false
	}
	relative = strings.Trim(filepath.ToSlash(relative), "/")
	if relative == "" || relative == "." {
		return true
	}
	if devPathExcluded(config, relative) {
		return false
	}
	return ignored == nil || !ignored.Match(strings.Split(relative, "/"), true)
}

func matchesDevWatch(patterns []string, rel string) bool {
	return newDevWatchMatcher(patterns).Match(rel)
}

func matchDevPathPattern(pattern, value string) bool {
	return newDevWatchMatcher([]string{pattern}).Match(value)
}

func newDevWatchMatcher(patterns []string) devWatchMatcher {
	matcher := devWatchMatcher{matchAll: len(patterns) == 0, patterns: make([][]devWatchSegment, 0, len(patterns))}
	for _, pattern := range patterns {
		pattern = strings.Trim(strings.TrimPrefix(filepath.ToSlash(pattern), "./"), "/")
		parts := strings.Split(pattern, "/")
		segments := make([]devWatchSegment, len(parts))
		for index, part := range parts {
			segments[index] = devWatchSegment{value: part, glob: strings.ContainsAny(part, `*?[\`), doubleStar: part == "**"}
		}
		matcher.patterns = append(matcher.patterns, segments)
	}
	return matcher
}

func (m devWatchMatcher) Match(value string) bool {
	if m.matchAll {
		return true
	}
	value = strings.Trim(strings.TrimPrefix(filepath.ToSlash(value), "./"), "/")
	valueParts := strings.Split(value, "/")
	first := make([]bool, len(valueParts)+1)
	second := make([]bool, len(valueParts)+1)
	for _, patternParts := range m.patterns {
		next, current := first, second
		clear(next)
		next[len(valueParts)] = true
		for patternIndex := len(patternParts) - 1; patternIndex >= 0; patternIndex-- {
			clear(current)
			part := patternParts[patternIndex]
			if part.doubleStar {
				current[len(valueParts)] = next[len(valueParts)]
				for valueIndex := len(valueParts) - 1; valueIndex >= 0; valueIndex-- {
					current[valueIndex] = next[valueIndex] || current[valueIndex+1]
				}
			} else {
				for valueIndex := 0; valueIndex < len(valueParts); valueIndex++ {
					if !next[valueIndex+1] {
						continue
					}
					if !part.glob {
						current[valueIndex] = part.value == valueParts[valueIndex]
						continue
					}
					segmentMatch, err := path.Match(part.value, valueParts[valueIndex])
					current[valueIndex] = err == nil && segmentMatch
				}
			}
			next, current = current, next
		}
		if next[0] {
			return true
		}
	}
	return false
}
