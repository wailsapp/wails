// Package dev owns the Wails HCL development session. It is an internal,
// Wails-specific controller, not a general-purpose process supervisor.
package dev

import (
	"context"
	"errors"
	"fmt"
	"github.com/rjeczalik/notify"
	"github.com/wailsapp/wails/v3/internal/report"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
	"io"
	"maps"
	"net"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Options struct {
	Host            string
	Profile, Target string
	Plan, Secure    bool
	VitePort        int
	Output          *Output
}
type BuildResult struct {
	Plan    pipeline.Plan
	Results map[pipeline.NodeKey]pipeline.Result
}
type rebuild struct {
	generation uint64
	loaded     *manifest.Loaded
	run        BuildResult
	err        error
}
type Process interface {
	Done() <-chan struct{}
	Stop(time.Duration)
	Err() error
	NeedsStopBeforeReplacement() bool
}
type EventKind int

const (
	Source EventKind = iota
	Ignore
	Directory
	GitIgnore
)

type WatchSet interface {
	Events() <-chan notify.EventInfo
	Stop()
	Classify(manifest.Config, notify.EventInfo) EventKind
	ContainsInput(manifest.Config, string) bool
}

// Operations adapts Wails build and host tools; Run owns their lifetime.
type Operations struct {
	ValidateTarget    func(string, string) error
	Getwd             func() (string, error)
	Load              func(string, string) (*manifest.Loaded, error)
	Target            func(string) (string, string, error)
	Plan              func(*manifest.Loaded, string, string) error
	CheckPort         func(string, int) error
	Build             func(context.Context, *manifest.Loaded, string, string, string, int) (BuildResult, error)
	StartFrontend     func(string, manifest.Config, string, int, string) (Process, error)
	WaitFrontendReady func(context.Context, Process, string, time.Duration) error
	BinaryPath        func(string, BuildResult, string, string) (string, error)
	StartApp          func(string, string, string, int) (Process, error)
	WaitReady         func(context.Context, Process, time.Duration) error
	StartWatches      func(string, manifest.Config) (WatchSet, error)
	RestartWatches    func(string, manifest.Config, WatchSet) (WatchSet, error)
	RestoreFrontend   func(context.Context, string, manifest.Config, int, string, string) (Process, error)
	RestoreApp        func(context.Context, string, Process, string, int) (Process, error)
	BackendChanged    func(BuildResult, string, string) bool
}

func sessionProcessExit(ctx context.Context, name string, process Process) error {
	if ctx.Err() != nil {
		return nil
	}
	return processFailure(process, name, "exit", name+" process exited unexpectedly")
}
func Run(ctx context.Context, options Options, ops Operations) (resultErr error) {
	component, phase := "session", "configuration"
	var generation uint64
	defer func() {
		var diagnostic *Diagnostic
		if resultErr != nil && !errors.As(resultErr, &diagnostic) {
			resultErr = &Diagnostic{Component: component, Phase: phase, Generation: generation, ExitCode: exitCode(resultErr), Err: resultErr}
		}
	}()
	if options.Profile != "" {
		return fmt.Errorf("dev does not accept production profiles; configure development policy in the dev block")
	}
	output := options.Output
	if output == nil {
		output = NewOutput(io.Discard, report.Normal)
	}
	sessionCtx, stopSession := context.WithCancel(ctx)
	defer stopSession()
	root, err := ops.Getwd()
	if err != nil {
		return err
	}
	loaded, err := ops.Load(root, "")
	if err != nil {
		return err
	}
	root = loaded.Config.Root
	goos, goarch, err := ops.Target(options.Target)
	if err != nil {
		return err
	}
	if goos == "" {
		goos, goarch = runtime.GOOS, runtime.GOARCH
	}
	if ops.ValidateTarget != nil {
		if err := ops.ValidateTarget(goos, goarch); err != nil {
			return err
		}
	} else if goos != runtime.GOOS || goarch != runtime.GOARCH {
		return fmt.Errorf("dev target must be the host %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	if options.Plan {
		return ops.Plan(loaded, goos, goarch)
	}
	port := options.VitePort
	if port == 0 {
		port = 9245
	}
	host := options.Host
	if host == "" {
		host = "127.0.0.1"
	}
	if err := ops.CheckPort(host, port); err != nil {
		return err
	}
	scheme := "http"
	if options.Secure {
		scheme = "https"
	}
	frontendURL := scheme + "://" + net.JoinHostPort(host, strconv.Itoa(port))

	component, phase = "build", "initial"
	initialRun, err := ops.Build(sessionCtx, loaded, goos, goarch, frontendURL, port)
	if err != nil {
		if errors.Is(err, context.Canceled) && sessionCtx.Err() != nil {
			return nil
		}
		return err
	}
	component, phase = "frontend", "start"
	frontend, err := ops.StartFrontend(root, loaded.Config, host, port, frontendURL)
	if err != nil {
		return err
	}
	if err := ops.WaitFrontendReady(sessionCtx, frontend, net.JoinHostPort(host, strconv.Itoa(port)), 30*time.Second); err != nil {
		frontend.Stop(time.Duration(loaded.Config.Dev.GracePeriodMS) * time.Millisecond)
		if errors.Is(err, context.Canceled) && sessionCtx.Err() != nil {
			return nil
		}
		return fmt.Errorf("frontend readiness: %w", err)
	}
	defer func() {
		if frontend != nil {
			frontend.Stop(time.Duration(loaded.Config.Dev.GracePeriodMS) * time.Millisecond)
		}
	}()
	output.Status(0, "frontend", "ready", "Frontend process started at "+frontendURL)
	binaryPath, err := ops.BinaryPath(root, initialRun, goos, goarch)
	if err != nil {
		return err
	}
	component, phase = "backend", "start"
	app, err := ops.StartApp(root, binaryPath, frontendURL, port)
	if err != nil {
		return err
	}
	if err := ops.WaitReady(sessionCtx, app, 30*time.Second); err != nil {
		app.Stop(time.Duration(loaded.Config.Dev.GracePeriodMS) * time.Millisecond)
		if errors.Is(err, context.Canceled) && sessionCtx.Err() != nil {
			return nil
		}
		return err
	}
	defer func() {
		if app != nil {
			app.Stop(time.Duration(loaded.Config.Dev.GracePeriodMS) * time.Millisecond)
			discard(app)
		}
	}()
	component, phase = "watcher", "start"
	watches, err := ops.StartWatches(root, loaded.Config)
	if err != nil {
		return err
	}
	defer func() { watches.Stop() }()
	component, phase = "session", "running"
	output.Status(0, "backend", "ready", "Backend built and started")
	debounce := time.Duration(loaded.Config.Dev.DebounceMS) * time.Millisecond
	if debounce <= 0 {
		debounce = 250 * time.Millisecond
	}
	var timer *time.Timer
	var timerC <-chan time.Time
	results := make(chan rebuild, 1)
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
		output.Status(generation, "build", "started", "Rebuilding")
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
				case results <- rebuild{generation: current, err: buildErr}:
				case <-sessionCtx.Done():
				}
				return
			}
			nextLoaded, loadErr := ops.Load(root, "")
			var run BuildResult
			if loadErr == nil {
				run, loadErr = ops.Build(WithGeneration(buildCtx, current), nextLoaded, goos, goarch, frontendURL, port)
			}
			result := rebuild{generation: current, loaded: nextLoaded, run: run, err: loadErr}
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
		case <-frontend.Done():
			return sessionProcessExit(sessionCtx, "frontend", frontend)
		case <-app.Done():
			return sessionProcessExit(sessionCtx, "backend", app)
		case event, ok := <-watches.Events():
			if !ok {
				return fmt.Errorf("project watcher stopped unexpectedly")
			}
			kind := watches.Classify(loaded.Config, event)
			createdDirectory := kind == Directory
			if kind == Ignore {
				continue
			}
			if createdDirectory {
				nextWatches, watchErr := ops.RestartWatches(root, loaded.Config, watches)
				if nextWatches != nil {
					watches = nextWatches
				}
				if watchErr != nil {
					output.Failure(generation, "watcher", "reload", watchErr, "new directory watch failed; restored the previous watch policy")
					continue
				}
				if !watches.ContainsInput(loaded.Config, event.Path()) {
					continue
				}
			} else if kind == GitIgnore {
				nextWatches, watchErr := ops.RestartWatches(root, loaded.Config, watches)
				if nextWatches != nil {
					watches = nextWatches
				}
				if watchErr != nil {
					output.Failure(generation, "watcher", "reload", watchErr, "gitignore reload failed; restored the previous watch policy")
					continue
				}
				output.Status(generation, "session", "updated", "Watch policy reloaded")
				continue
			} else if kind == Ignore {
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
					output.Failure(generation, "build", "rebuild", result.err, "Current app is still running.")
				}
				continue
			}
			nextLoaded := result.loaded
			var nextWatches WatchSet
			if WatchSessionChanged(loaded.Config, nextLoaded.Config) {
				nextWatches, err = ops.StartWatches(root, nextLoaded.Config)
				if err != nil {
					output.Failure(generation, "watcher", "reconfigure", err, "watch reconfiguration failed; keeping the current session")
					continue
				}
			}

			frontendChanged := FrontendSessionChanged(loaded.Config, nextLoaded.Config)
			backendChanged := ops.BackendChanged(result.run, goos, goarch)
			oldFrontend := frontend
			var nextFrontend Process
			if frontendChanged {
				oldFrontend.Stop(time.Duration(loaded.Config.Dev.GracePeriodMS) * time.Millisecond)
				nextFrontend, err = ops.StartFrontend(root, nextLoaded.Config, host, port, frontendURL)
				if err == nil {
					err = ops.WaitFrontendReady(sessionCtx, nextFrontend, net.JoinHostPort(host, strconv.Itoa(port)), 30*time.Second)
				}
				if err != nil {
					transitionErr := err
					if nextFrontend != nil {
						nextFrontend.Stop(time.Duration(nextLoaded.Config.Dev.GracePeriodMS) * time.Millisecond)
					}
					if nextWatches != nil {
						nextWatches.Stop()
					}
					if errors.Is(transitionErr, context.Canceled) && sessionCtx.Err() != nil {
						return nil
					}
					frontend, err = ops.RestoreFrontend(sessionCtx, root, loaded.Config, port, frontendURL, host)
					if err != nil {
						return err
					}
					output.Failure(generation, "frontend", "restart", transitionErr, "frontend restart failed; keeping the current session")
					continue
				}
			}

			var nextApp Process
			var stoppedApp Process
			if backendChanged {
				binaryPath, pathErr := ops.BinaryPath(root, result.run, goos, goarch)
				if pathErr != nil {
					err = pathErr
				} else {
					if app.NeedsStopBeforeReplacement() {
						stoppedApp = app
						stoppedApp.Stop(time.Duration(loaded.Config.Dev.GracePeriodMS) * time.Millisecond)
					}
					nextApp, err = ops.StartApp(root, binaryPath, frontendURL, port)
				}
				if err == nil {
					err = ops.WaitReady(sessionCtx, nextApp, 30*time.Second)
				}
				if err != nil {
					transitionErr := err
					if nextApp != nil {
						nextApp.Stop(time.Duration(nextLoaded.Config.Dev.GracePeriodMS) * time.Millisecond)
						discard(nextApp)
					}
					if nextWatches != nil {
						nextWatches.Stop()
					}
					if nextFrontend != nil {
						nextFrontend.Stop(time.Duration(nextLoaded.Config.Dev.GracePeriodMS) * time.Millisecond)
						frontend, err = ops.RestoreFrontend(sessionCtx, root, loaded.Config, port, frontendURL, host)
						if err != nil {
							return err
						}
					}
					if errors.Is(transitionErr, context.Canceled) && sessionCtx.Err() != nil {
						return nil
					}
					if stoppedApp != nil {
						app, err = ops.RestoreApp(sessionCtx, root, stoppedApp, frontendURL, port)
						if err != nil {
							return fmt.Errorf("replacement failed (%v) and restoring the previous backend failed: %w", transitionErr, err)
						}
						discard(stoppedApp)
					}
					output.Failure(generation, "backend", "restart", transitionErr, "restart failed; keeping the current app running")
					continue
				}
			}

			if backendChanged {
				oldApp := app
				app = nextApp
				oldApp.Stop(time.Duration(loaded.Config.Dev.GracePeriodMS) * time.Millisecond)
				discard(oldApp)
			}
			if frontendChanged {
				frontend = nextFrontend
				output.Status(generation, "frontend", "ready", "Frontend process restarted at "+frontendURL)
			}
			if nextWatches != nil {
				oldWatches := watches
				watches = nextWatches
				oldWatches.Stop()
			}
			loaded = nextLoaded
			debounce = time.Duration(loaded.Config.Dev.DebounceMS) * time.Millisecond
			if debounce <= 0 {
				debounce = 250 * time.Millisecond
			}
			switch {
			case backendChanged:
				output.Status(generation, "session", "updated", "Backend rebuilt and restarted")
			case frontendChanged:
				output.Status(generation, "session", "updated", "Frontend session updated; backend unchanged")
			default:
				output.Status(generation, "session", "updated", "Build is current; backend unchanged")
			}
		}
	}
}

func WatchSessionChanged(current, next manifest.Config) bool {
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

func FrontendSessionChanged(current, next manifest.Config) bool {
	return current.Frontend.Directory != next.Frontend.Directory || current.Frontend.PackageManager != next.Frontend.PackageManager || current.Frontend.DevCommand != next.Frontend.DevCommand || !equalStrings(current.Frontend.Dev, next.Frontend.Dev) || !maps.Equal(current.Frontend.Environment, next.Frontend.Environment)
}

func discard(p Process) {
	if d, ok := p.(interface{ Discard() }); ok {
		d.Discard()
	}
}
