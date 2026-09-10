package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	"github.com/wailsapp/wails/v3/internal/dev"
	"github.com/wailsapp/wails/v3/internal/report"
	"github.com/wailsapp/wails/v3/internal/wake/cache"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

type manifestProcess struct {
	tail *dev.Tail
	// Windows must stop the previous WebView2 owner before starting another
	// backend with the same browser profile. Keep its image for startup rollback.
	readiness    *dev.Readiness
	restartImage []byte
	cmd          *exec.Cmd
	done         chan struct{}
	mu           sync.Mutex
	err          error
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
	output          *dev.Output
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

func productionManifestDevOps(options ...*DevOptions) manifestDevOps {
	scheme := "http"
	if len(options) > 0 && options[0].Secure {
		scheme = "https"
	}
	level := report.Normal
	if len(options) > 0 {
		if options[0].Verbose {
			level = report.Verbose
		}
		if options[0].Quiet {
			level = report.Silent
		}
	}
	output := dev.NewOutput(os.Stdout, level)
	return manifestDevOps{
		output:    output,
		getwd:     os.Getwd,
		load:      manifest.Load,
		checkPort: checkManifestDevPort,
		build:     runManifestDevBuild,
		startFrontend: func(r string, c manifest.Config, h string, p int, u string) (*manifestProcess, error) {
			return startFrontendDev(r, c, h, p, u, output)
		},
		waitTCP: func(c context.Context, p *manifestProcess, a string, t time.Duration) error {
			return dev.WaitHTTP(c, p, scheme+"://"+a, t)
		},
		binaryPath: manifestDevBinaryPath,
		startApp: func(r, b, u string, p int) (*manifestProcess, error) {
			return startManifestReadyApp(r, b, u, p, output)
		},
		waitStable:     waitForManifestRuntime,
		startWatches:   startManifestWatches,
		restartWatches: restartManifestWatches,
		restoreFrontend: func(c context.Context, r string, cfg manifest.Config, p int, u, h string) (*manifestProcess, error) {
			return restoreManifestFrontend(c, r, cfg, p, u, h, output)
		},
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
	return runManifestDevContextWithOps(ctx, options, productionManifestDevOps(options))
}

func runManifestDevContextWithOps(ctx context.Context, options *DevOptions, ops manifestDevOps) error {
	output := ops.output
	if output == nil {
		output = dev.NewOutput(os.Stdout, report.Normal)
	}
	operations := dev.Operations{
		Getwd: ops.getwd, Load: ops.load, Target: splitTarget, CheckPort: ops.checkPort,
		Plan: func(l *manifest.Loaded, o, a string) error {
			return printManifestPlan(manifestRunOptions{Verb: "build", Loaded: l, TargetOS: o, TargetArch: a, Development: true, Tags: manifestDevTags(options)}, false)
		},
		Build: func(c context.Context, l *manifest.Loaded, o, a, u string, p int) (manifestPipelineRun, error) {
			return ops.build(dev.WithOutput(c, output), options, l, o, a, u, p)
		},
		StartFrontend: func(r string, c manifest.Config, h string, p int, u string) (dev.Process, error) {
			v, e := ops.startFrontend(r, c, h, p, u)
			if e != nil {
				return nil, e
			}
			return v, nil
		},
		WaitFrontendReady: func(c context.Context, p dev.Process, a string, t time.Duration) error {
			return ops.waitTCP(c, p.(*manifestProcess), a, t)
		},
		BinaryPath: ops.binaryPath,
		StartApp: func(r, b, u string, p int) (dev.Process, error) {
			v, e := ops.startApp(r, b, u, p)
			if e != nil {
				return nil, e
			}
			return v, nil
		},
		WaitReady: func(c context.Context, p dev.Process, t time.Duration) error {
			return ops.waitStable(c, p.(*manifestProcess), t)
		},
		StartWatches: func(r string, c manifest.Config) (dev.WatchSet, error) {
			v, e := ops.startWatches(r, c)
			if e != nil {
				return nil, e
			}
			return v, nil
		},
		RestartWatches: func(r string, c manifest.Config, w dev.WatchSet) (dev.WatchSet, error) {
			v, e := ops.restartWatches(r, c, w.(*manifestWatchSet))
			if v == nil {
				return nil, e
			}
			return v, e
		},
		RestoreFrontend: func(c context.Context, r string, cfg manifest.Config, p int, u, h string) (dev.Process, error) {
			v, e := ops.restoreFrontend(c, r, cfg, p, u, h)
			if e != nil {
				return nil, e
			}
			return v, nil
		},
		RestoreApp: func(c context.Context, r string, p dev.Process, u string, port int) (dev.Process, error) {
			v, e := restoreManifestWindowsApp(c, r, p.(*manifestProcess), u, port, output)
			if e != nil {
				return nil, e
			}
			return v, nil
		},
		BackendChanged: manifestBackendChanged,
	}
	if strings.HasPrefix(options.Target, "ios/") || strings.HasPrefix(options.Target, "android/") {
		cleanup, err := configureMobileDev(ctx, options, output, &operations)
		if err != nil {
			return err
		}
		defer cleanup()
	} else if options.Device != "" || options.Emulator != "" || options.Destination != "" {
		return fmt.Errorf("--device, --emulator and --destination require a mobile --target")
	}
	err := dev.Run(ctx, dev.Options{Profile: options.Profile, Target: options.Target, Plan: options.Plan, Secure: options.Secure, VitePort: options.VitePort, Host: options.Host, Output: output}, operations)
	if err != nil {
		output.Failure(0, "session", "failed", err, "")
		return report.MarkReported(err)
	}
	return nil
}
func (p *manifestProcess) Done() <-chan struct{} { return p.done }
func (p *manifestProcess) Stop(t time.Duration)  { p.stop(t) }
func (p *manifestProcess) Err() error {
	err := p.waitError()
	if err == nil || p.cmd == nil {
		return err
	}
	output := ""
	if p.tail != nil {
		output = p.tail.String()
	}
	code := -1
	if p.cmd.ProcessState != nil {
		code = p.cmd.ProcessState.ExitCode()
	}
	return &dev.Diagnostic{Component: filepath.Base(p.cmd.Path), Phase: "exit", Command: p.cmd.Args, Output: output, ExitCode: code, Err: err}
}
func (p *manifestProcess) NeedsStopBeforeReplacement() bool { return len(p.restartImage) > 0 }
func (w *manifestWatchSet) Events() <-chan notify.EventInfo { return w.events }
func (w *manifestWatchSet) Stop()                           { w.stop() }
func (w *manifestWatchSet) ContainsInput(c manifest.Config, p string) bool {
	return directoryContainsDevInputMatched(w.root, c, w.ignored, w.matcher, p)
}
func (w *manifestWatchSet) Classify(c manifest.Config, e notify.EventInfo) dev.EventKind {
	if isCreatedDevDirectory(e) {
		if shouldRefreshDevWatchesForDirectory(w.root, c, w.ignored, e.Path()) {
			return dev.Directory
		}
		return dev.Ignore
	}
	if isDevGitIgnoreEvent(w.root, c, e.Path()) {
		return dev.GitIgnore
	}
	if ignoreDevEventMatched(w.root, c, w.ignored, w.matcher, e.Path()) {
		return dev.Ignore
	}
	return dev.Source
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
	if err := registerManifestWatches(canonicalRoot, config, ignored, events); err != nil {
		stopManifestWatchChannel(events)
		return nil, fmt.Errorf("watch project: %w", err)
	}
	return &manifestWatchSet{events: events, ignored: ignored, matcher: newDevWatchMatcher(config.Dev.Watch), root: canonicalRoot}, nil
}

func (w *manifestWatchSet) stop() {
	if w != nil && w.events != nil {
		stopManifestWatchChannel(w.events)
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

func restoreManifestFrontend(ctx context.Context, root string, config manifest.Config, port int, frontendURL, host string, output ...*dev.Output) (*manifestProcess, error) {
	frontend, err := startFrontendDev(root, config, host, port, frontendURL, output...)
	if err != nil {
		return nil, fmt.Errorf("restore previous frontend: %w", err)
	}
	if err := dev.WaitHTTP(ctx, frontend, frontendURL, 30*time.Second); err != nil {
		frontend.stop(time.Duration(config.Dev.GracePeriodMS) * time.Millisecond)
		return nil, fmt.Errorf("restore previous frontend readiness: %w", err)
	}
	return frontend, nil
}

func startFrontendDev(root string, config manifest.Config, host string, port int, frontendURL string, output ...*dev.Output) (*manifestProcess, error) {
	if len(config.Frontend.Dev) > 0 {
		args := append([]string(nil), config.Frontend.Dev...)
		// Migration and explicit HCL defaults express package scripts as argv.
		// Recognised Vite scripts need the same endpoint configuration as
		// compiled defaults. Other user commands keep their exact argv/launcher.
		if manager, script, ok := frontendPackageScript(args); ok {
			config.Frontend.PackageManager = manager
			config.Frontend.DevCommand = script
			if frontendDevCommandUsesVite(root, config) {
				config.Frontend.Dev = nil
				return startFrontendDev(root, config, host, port, frontendURL, output...)
			}
		}
		env := declaredEnvironment(nil, config.Frontend.Environment)
		env = mergeEnvironment(env, []string{wailsVitePort + "=" + strconv.Itoa(port), "FRONTEND_DEVSERVER_URL=" + frontendURL})
		return startManifestProcessOutput(output, filepath.Join(root, config.Frontend.Directory), args[0], env, args[1:]...)
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
	env := declaredEnvironment(nil, config.Frontend.Environment)
	env = mergeEnvironment(env, []string{wailsVitePort + "=" + strconv.Itoa(port), "FRONTEND_DEVSERVER_URL=" + frontendURL})
	name, resolved, err := resolvePackageManagerProcess(manager, args, exec.LookPath, filepath.EvalSymlinks)
	if err != nil {
		return nil, err
	}
	return startManifestProcessOutput(output, filepath.Join(root, config.Frontend.Directory), name, env, resolved...)
}

func frontendPackageScript(args []string) (manager, script string, ok bool) {
	if len(args) < 2 || !containsString([]string{"npm", "pnpm", "yarn", "bun"}, args[0]) {
		return "", "", false
	}
	if len(args) == 3 && args[1] == "run" && args[2] != "" && !strings.HasPrefix(args[2], "-") {
		return args[0], args[2], true
	}
	if len(args) == 2 && args[0] == "yarn" && args[1] != "" && !strings.HasPrefix(args[1], "-") {
		return args[0], args[1], true
	}
	return "", "", false
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
			// npm's Windows shim lives beside node_modules/npm. Launch its
			// JavaScript entry point directly, avoiding cmd.exe's incompatible
			// quoting rules for installations under Program Files.
			node, err := lookPath("node")
			if err != nil {
				return "", nil, err
			}
			return node, append([]string{filepath.Join(filepath.Dir(npm), "node_modules", "npm", "bin", "npm-cli.js")}, args...), nil
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

func startManifestApp(root, binaryPath, frontendURL string, port int, output ...*dev.Output) (*manifestProcess, error) {
	return startManifestAppEnvironment(root, binaryPath, frontendURL, port, nil, output...)
}
func startManifestAppEnvironment(root, binaryPath, frontendURL string, port int, environment []string, output ...*dev.Output) (*manifestProcess, error) {
	env := append([]string{wailsVitePort + "=" + strconv.Itoa(port), "FRONTEND_DEVSERVER_URL=" + frontendURL}, environment...)
	var restartImage []byte
	if runtime.GOOS == "windows" {
		var err error
		restartImage, err = os.ReadFile(binaryPath)
		if err != nil {
			return nil, err
		}
	}
	process, err := startManifestProcessOutput(output, root, binaryPath, env)
	if err == nil {
		process.restartImage = restartImage
	}
	return process, err
}

func restoreManifestWindowsApp(ctx context.Context, root string, previous *manifestProcess, frontendURL string, port int, output ...*dev.Output) (*manifestProcess, error) {
	binaryPath := previous.cmd.Path
	if err := producePathTransactional(binaryPath, ".backend-rollback-*", func(staged string) error {
		return os.WriteFile(staged, previous.restartImage, 0o755)
	}); err != nil {
		return nil, err
	}
	var process *manifestProcess
	var err error
	if previous.readiness != nil {
		var out *dev.Output
		if len(output) > 0 {
			out = output[0]
		}
		process, err = startManifestReadyApp(root, binaryPath, frontendURL, port, out)
	} else {
		process, err = startManifestApp(root, binaryPath, frontendURL, port, output...)
	}
	if err != nil {
		return nil, err
	}
	wait := waitForProcessStable
	timeout := manifestBackendReadinessDelay
	if process.readiness != nil {
		wait = waitForManifestRuntime
		timeout = 30 * time.Second
	}
	if err := wait(ctx, process, timeout); err != nil {
		process.stop(time.Second)
		return nil, err
	}
	return process, nil
}
func startManifestProcess(dir, name string, env []string, args ...string) (*manifestProcess, error) {
	return startManifestProcessOutput(nil, dir, name, env, args...)
}
func startManifestProcessOutput(output []*dev.Output, dir, name string, env []string, args ...string) (*manifestProcess, error) {
	cmd := exec.CommandContext(context.Background(), name, args...)
	cmd.Dir = dir
	cmd.WaitDelay = 2 * time.Second
	cmd.Env = mergeEnvironment(os.Environ(), env)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	var closeOutput func()
	if len(output) > 0 && output[0] != nil {
		stdout, stderr := output[0].Writer(filepath.Base(name), "stdout"), output[0].Writer(filepath.Base(name), "stderr")
		cmd.Stdout, cmd.Stderr = stdout, stderr
		closeOutput = func() { stdout.Close(); stderr.Close() }
	}
	tail := &dev.Tail{}
	cmd.Stdout = io.MultiWriter(cmd.Stdout, tail)
	cmd.Stderr = io.MultiWriter(cmd.Stderr, tail)
	cleanup, err := startManifestOwnedProcess(cmd)
	if err != nil {
		if closeOutput != nil {
			closeOutput()
		}
		return nil, &dev.Diagnostic{Component: filepath.Base(name), Phase: "start", Command: cmd.Args, ExitCode: -1, Err: err}
	}
	result := &manifestProcess{cmd: cmd, done: make(chan struct{}), tail: tail}
	go func() {
		waitErr := cmd.Wait()
		if closeOutput != nil {
			closeOutput()
		}
		// A wrapper may exit before its descendants. Clean its group immediately,
		// before exposing completion, instead of signalling a stale PID in stop.
		cleanup()
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
	p.stopWithSignal(grace, os.Interrupt)
}

func (p *manifestProcess) stopWithSignal(grace time.Duration, termination os.Signal) {
	if p == nil || p.cmd == nil || p.cmd.Process == nil {
		return
	}
	select {
	case <-p.done:
		return
	default:
	}
	_ = signalManifestProcess(p.cmd.Process, termination)
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
	// The manifest controls the session itself, even when migrated or custom
	// source watch patterns contain only Go files.
	if rel == manifest.Filename {
		return false
	}
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
		return watchManifestDirectory(current, events)
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

func startManifestReadyApp(root, binaryPath, frontendURL string, port int, output *dev.Output) (*manifestProcess, error) {
	ready, err := dev.NewReadiness()
	if err != nil {
		return nil, err
	}
	p, err := startManifestAppEnvironment(root, binaryPath, frontendURL, port, ready.Environment(), output)
	if err != nil {
		ready.Close()
		return nil, err
	}
	p.readiness = ready
	go func() { <-p.done; ready.Close() }()
	return p, nil
}
func waitForManifestRuntime(ctx context.Context, p *manifestProcess, timeout time.Duration) error {
	if p.readiness == nil {
		return fmt.Errorf("backend launch is missing the Wails runtime readiness channel")
	}
	return p.readiness.Wait(ctx, p, timeout)
}
