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
	port       int
	url        string
	run        manifestPipelineRun
	err        error
}

type manifestWatchSet struct {
	events  chan notify.EventInfo
	ignored gitignore.Matcher
}

func runManifestDev(options *DevOptions) error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	loaded, err := manifest.Load(root, options.Profile)
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
	sessionCtx, stopSession := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSession()
	port := options.VitePort
	if port == 0 {
		port = loaded.Config.Dev.Port
	}
	if port == 0 {
		port = defaultVitePort
	}
	host := "127.0.0.1"
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return fmt.Errorf("frontend port %d is unavailable: %w", port, err)
	}
	_ = listener.Close()
	scheme := "http"
	if options.Secure {
		scheme = "https"
	}
	frontendURL := fmt.Sprintf("%s://%s:%d", scheme, host, port)

	if _, err := runManifestDevBuild(sessionCtx, options, loaded, goos, goarch, frontendURL, port); err != nil {
		if errors.Is(err, context.Canceled) && sessionCtx.Err() != nil {
			return nil
		}
		return err
	}
	frontend, err := startFrontendDev(root, loaded.Config, host, port, frontendURL)
	if err != nil {
		return err
	}
	if err := waitForProcessTCP(sessionCtx, frontend, net.JoinHostPort(host, strconv.Itoa(port)), 30*time.Second); err != nil {
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
	app, err := startManifestApp(root, loaded.Config, goos, frontendURL, port)
	if err != nil {
		return err
	}
	if err := waitForProcessStable(sessionCtx, app, 150*time.Millisecond); err != nil {
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
	fmt.Println("Backend built and started")

	watches, err := startManifestWatches(root, loaded.Config)
	if err != nil {
		return err
	}
	defer func() { watches.stop() }()
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
		currentPort := port
		currentURL := frontendURL
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
			nextLoaded, loadErr := manifest.Load(root, options.Profile)
			nextPort := currentPort
			nextURL := currentURL
			var run manifestPipelineRun
			if loadErr == nil {
				nextPort = resolvedDevPort(options, nextLoaded.Config)
				nextURL = fmt.Sprintf("%s://%s:%d", scheme, host, nextPort)
				run, loadErr = runManifestDevBuild(buildCtx, options, nextLoaded, goos, goarch, nextURL, nextPort)
			}
			result := manifestRebuild{generation: current, loaded: nextLoaded, port: nextPort, url: nextURL, run: run, err: loadErr}
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
			if sessionCtx.Err() != nil {
				return nil
			}
			return manifestProcessExitError("frontend", frontend)
		case <-app.done:
			if sessionCtx.Err() != nil {
				return nil
			}
			return manifestProcessExitError("backend", app)
		case event, ok := <-watches.events:
			if !ok {
				return fmt.Errorf("project watcher stopped unexpectedly")
			}
			createdDirectory := isCreatedDevDirectory(event)
			if createdDirectory && !shouldRefreshDevWatchesForDirectory(root, loaded.Config, watches.ignored, event.Path()) {
				continue
			}
			if createdDirectory {
				nextWatches, watchErr := startManifestWatches(root, loaded.Config)
				if watchErr != nil {
					fmt.Fprintf(os.Stderr, "new directory watch failed; retaining the previous watch policy: %v\n", watchErr)
					continue
				}
				oldWatches := watches
				watches = nextWatches
				oldWatches.stop()
				if !directoryContainsDevInput(root, loaded.Config, watches.ignored, event.Path()) {
					continue
				}
			} else if isDevGitIgnoreEvent(root, loaded.Config, event.Path()) {
				nextWatches, watchErr := startManifestWatches(root, loaded.Config)
				if watchErr != nil {
					fmt.Fprintf(os.Stderr, "gitignore reload failed; retaining the previous watch policy: %v\n", watchErr)
					continue
				}
				oldWatches := watches
				watches = nextWatches
				oldWatches.stop()
				fmt.Println("Watch policy reloaded")
				continue
			} else if ignoreDevEvent(root, loaded.Config, watches.ignored, event.Path()) {
				continue
			}
			if timer == nil {
				timer = time.NewTimer(debounce)
			} else {
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
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
				nextWatches, err = startManifestWatches(root, nextLoaded.Config)
				if err != nil {
					fmt.Fprintf(os.Stderr, "watch reconfiguration failed; keeping the current session: %v\n", err)
					continue
				}
			}

			frontendChanged := frontendSessionChanged(loaded.Config, port, nextLoaded.Config, result.port)
			backendChanged := manifestBackendChanged(result.run, goos, goarch, frontendURL, result.url)
			oldFrontend := frontend
			var nextFrontend *manifestProcess
			oldFrontendStopped := false
			if frontendChanged {
				if result.port == port {
					oldFrontend.stop(time.Duration(loaded.Config.Dev.GracePeriodMS) * time.Millisecond)
					oldFrontendStopped = true
				}
				nextFrontend, err = startFrontendDev(root, nextLoaded.Config, host, result.port, result.url)
				if err == nil {
					err = waitForProcessTCP(sessionCtx, nextFrontend, net.JoinHostPort(host, strconv.Itoa(result.port)), 30*time.Second)
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
					if oldFrontendStopped {
						frontend, err = restoreManifestFrontend(sessionCtx, root, loaded.Config, port, frontendURL, host)
						if err != nil {
							return err
						}
					}
					fmt.Fprintf(os.Stderr, "frontend restart failed; keeping the current session: %v\n", transitionErr)
					continue
				}
			}

			var nextApp *manifestProcess
			if backendChanged {
				nextApp, err = startManifestApp(root, nextLoaded.Config, goos, result.url, result.port)
				if err == nil {
					err = waitForProcessStable(sessionCtx, nextApp, 150*time.Millisecond)
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
						if oldFrontendStopped {
							frontend, err = restoreManifestFrontend(sessionCtx, root, loaded.Config, port, frontendURL, host)
							if err != nil {
								return err
							}
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
				if !oldFrontendStopped {
					oldFrontend.stop(time.Duration(loaded.Config.Dev.GracePeriodMS) * time.Millisecond)
				}
				frontend = nextFrontend
				fmt.Printf("Frontend process restarted at %s\n", result.url)
			}
			if nextWatches != nil {
				oldWatches := watches
				watches = nextWatches
				oldWatches.stop()
			}
			loaded = nextLoaded
			port = result.port
			frontendURL = result.url
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
	return runManifestPipelineResult(manifestRunOptions{Context: ctx, Verb: "build", Profile: options.Profile, Loaded: loaded, TargetOS: goos, TargetArch: goarch, Environment: environment, Development: true, Tags: envTags()})
}

func manifestBackendChanged(run manifestPipelineRun, goos, goarch, currentURL, nextURL string) bool {
	if currentURL != nextURL {
		return true
	}
	result, exists := run.Results[pipelineCompileKey(goos, goarch)]
	if !exists || result.Status == cache.LookupMiss {
		return true
	}
	for key, node := range run.Plan.Nodes {
		if node.Kind != pipeline.RunHook {
			continue
		}
		spec, ok := node.Spec.(pipeline.HookSpec)
		if !ok || spec.Phase != "after_build" {
			continue
		}
		hookResult, ok := run.Results[key]
		if !ok || hookResult.Status == cache.LookupMiss {
			return true
		}
	}
	return false
}

func pipelineCompileKey(goos, goarch string) pipeline.NodeKey {
	return pipeline.NodeKey("target:" + goos + "/" + goarch + ":compile")
}

func devWatchSessionChanged(current, next manifest.Config) bool {
	return current.Dev.UseGitIgnore != next.Dev.UseGitIgnore ||
		current.Build.OutputDirectory != next.Build.OutputDirectory ||
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
	found := false
	_ = filepath.WalkDir(directory, func(current string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil || found {
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
		if !ignoreDevEvent(root, config, ignored, current) {
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	return found
}

func startManifestWatches(root string, config manifest.Config) (*manifestWatchSet, error) {
	ignored, err := loadDevGitIgnore(root, config.Dev.UseGitIgnore)
	if err != nil {
		return nil, fmt.Errorf("load .gitignore: %w", err)
	}
	events := make(chan notify.EventInfo, 1024)
	if err := registerManifestWatches(root, config, ignored, events); err != nil {
		notify.Stop(events)
		return nil, fmt.Errorf("watch project: %w", err)
	}
	return &manifestWatchSet{events: events, ignored: ignored}, nil
}

func (w *manifestWatchSet) stop() {
	if w != nil && w.events != nil {
		notify.Stop(w.events)
	}
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

func reportManifestRebuildFailure(err error) {
	if wake.IsReported(err) {
		fmt.Fprintln(os.Stderr, "Current app is still running.")
		return
	}
	fmt.Fprintf(os.Stderr, "Rebuild failed; current app is still running: %v\n", err)
}

func resolvedDevPort(options *DevOptions, config manifest.Config) int {
	if options.VitePort != 0 {
		return options.VitePort
	}
	if config.Dev.Port != 0 {
		return config.Dev.Port
	}
	return defaultVitePort
}

func frontendSessionChanged(current manifest.Config, currentPort int, next manifest.Config, nextPort int) bool {
	return currentPort != nextPort || current.Frontend.Directory != next.Frontend.Directory || current.Frontend.PackageManager != next.Frontend.PackageManager || current.Frontend.DevCommand != next.Frontend.DevCommand
}

func startFrontendDev(root string, config manifest.Config, host string, port int, frontendURL string) (*manifestProcess, error) {
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
	name := manager
	if manager == "npm" {
		npm, err := exec.LookPath("npm")
		if err != nil {
			return nil, err
		}
		script, err := filepath.EvalSymlinks(npm)
		if err != nil {
			script = npm
		}
		node, err := exec.LookPath("node")
		if err != nil {
			return nil, err
		}
		name = node
		args = append([]string{script}, args...)
	}
	return startManifestProcess(dir, name, env, args...)
}
func startManifestApp(root string, config manifest.Config, goos, frontendURL string, port int) (*manifestProcess, error) {
	name := config.Project.BinaryName
	if goos == "windows" {
		name += ".exe"
	}
	env := []string{wailsVitePort + "=" + strconv.Itoa(port), "FRONTEND_DEVSERVER_URL=" + frontendURL}
	return startManifestProcess(root, filepath.Join(root, config.Build.OutputDirectory, name), env)
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
	return strings.HasSuffix(rel, "_test.go") || !matchesDevWatch(config.Dev.Watch, rel)
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
	paths := []string{config.Build.OutputDirectory, config.Frontend.Directory}
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
	if len(patterns) == 0 {
		return true
	}
	rel = filepath.ToSlash(rel)
	for _, pattern := range patterns {
		if matchDevPathPattern(filepath.ToSlash(pattern), rel) {
			return true
		}
	}
	return false
}

func matchDevPathPattern(pattern, value string) bool {
	pattern = strings.Trim(strings.TrimPrefix(pattern, "./"), "/")
	value = strings.Trim(strings.TrimPrefix(value, "./"), "/")
	patternParts := strings.Split(pattern, "/")
	valueParts := strings.Split(value, "/")
	type state struct{ pattern, value int }
	memo := map[state]bool{}
	visited := map[state]bool{}
	var match func(int, int) bool
	match = func(patternIndex, valueIndex int) bool {
		key := state{patternIndex, valueIndex}
		if visited[key] {
			return memo[key]
		}
		visited[key] = true
		if patternIndex == len(patternParts) {
			memo[key] = valueIndex == len(valueParts)
			return memo[key]
		}
		if patternParts[patternIndex] == "**" {
			memo[key] = match(patternIndex+1, valueIndex) || (valueIndex < len(valueParts) && match(patternIndex, valueIndex+1))
			return memo[key]
		}
		if valueIndex == len(valueParts) {
			return false
		}
		segmentMatch, err := path.Match(patternParts[patternIndex], valueParts[valueIndex])
		if err != nil || !segmentMatch {
			return false
		}
		memo[key] = match(patternIndex+1, valueIndex+1)
		return memo[key]
	}
	return match(0, 0)
}
