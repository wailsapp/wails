package commands

import (
	"context"
	"errors"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rjeczalik/notify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/cache"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

type manifestDevTestEvent struct {
	event notify.Event
	path  string
}

func (e manifestDevTestEvent) Event() notify.Event { return e.event }
func (e manifestDevTestEvent) Path() string        { return e.path }
func (e manifestDevTestEvent) Sys() any            { return nil }

func TestManifestDevAdapterCoversUnrepresentableStartupFailures(t *testing.T) {
	t.Run("Plan invariant", func(t *testing.T) {
		ops, _ := newManifestDevTestOps(t)
		ops.binaryPath = func(string, manifestPipelineRun, string, string) (string, error) {
			return "", errors.New("missing compile output invariant")
		}
		err := runManifestDevContextWithOps(t.Context(), &DevOptions{}, ops)
		assert.ErrorContains(t, err, "missing compile output invariant")
	})

	t.Run("watcher boundary", func(t *testing.T) {
		ops, _ := newManifestDevTestOps(t)
		ops.startWatches = func(string, manifest.Config) (*manifestWatchSet, error) {
			return nil, errors.New("watch adapter unavailable")
		}
		err := runManifestDevContextWithOps(t.Context(), &DevOptions{}, ops)
		assert.ErrorContains(t, err, "watch adapter unavailable")
	})
}

func TestManifestDevReportsClosedWatcherAdapter(t *testing.T) {
	ops, watches := newManifestDevTestOps(t)
	ready := make(chan struct{})
	ops.startWatches = func(string, manifest.Config) (*manifestWatchSet, error) {
		close(ready)
		return watches, nil
	}
	done := make(chan error, 1)
	go func() { done <- runManifestDevContextWithOps(context.Background(), &DevOptions{}, ops) }()
	<-ready
	close(watches.events)
	select {
	case err := <-done:
		assert.ErrorContains(t, err, "project watcher stopped unexpectedly")
	case <-time.After(2 * time.Second):
		t.Fatal("closed watcher was not reported")
	}
}

func TestManifestDevKeepsSessionWhenWatchReconfigurationFails(t *testing.T) {
	ops, watches := newManifestDevTestOps(t)
	initial, err := ops.load("", "")
	require.NoError(t, err)
	next := cloneManifestDevLoaded(initial)
	next.Config.Dev.Watch = []string{"**/*.go"}
	var loads atomic.Int32
	ops.load = func(string, string) (*manifest.Loaded, error) {
		if loads.Add(1) == 1 {
			return initial, nil
		}
		return next, nil
	}
	reconfigured := make(chan struct{})
	var watchStarts atomic.Int32
	ops.startWatches = func(string, manifest.Config) (*manifestWatchSet, error) {
		if watchStarts.Add(1) == 1 {
			return watches, nil
		}
		close(reconfigured)
		return nil, errors.New("synthetic watch reconfiguration failure")
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- runManifestDevContextWithOps(ctx, &DevOptions{}, ops) }()
	watches.events <- manifestDevTestEvent{event: notify.Write, path: filepath.Join(watches.root, "main.go")}
	select {
	case <-reconfigured:
	case <-time.After(2 * time.Second):
		t.Fatal("watch reconfiguration was not attempted")
	}
	cancel()
	assert.NoError(t, <-done)
}

func TestManifestDevFrontendRollbackStopsStagedWatcher(t *testing.T) {
	for _, restoreFails := range []bool{false, true} {
		t.Run(map[bool]string{false: "restored", true: "restore failure"}[restoreFails], func(t *testing.T) {
			ops, watches := newManifestDevTestOps(t)
			initial, err := ops.load("", "")
			require.NoError(t, err)
			next := cloneManifestDevLoaded(initial)
			next.Config.Dev.Watch = []string{"**/*.go"}
			next.Config.Frontend.Dev = []string{"replacement"}
			var loads atomic.Int32
			ops.load = func(string, string) (*manifest.Loaded, error) {
				if loads.Add(1) == 1 {
					return initial, nil
				}
				return next, nil
			}
			stagedWatches := &manifestWatchSet{root: watches.root, matcher: newDevWatchMatcher(next.Config.Dev.Watch)}
			var watchStarts atomic.Int32
			ops.startWatches = func(string, manifest.Config) (*manifestWatchSet, error) {
				if watchStarts.Add(1) == 1 {
					return watches, nil
				}
				return stagedWatches, nil
			}
			var frontendStarts atomic.Int32
			ops.startFrontend = func(string, manifest.Config, string, int, string) (*manifestProcess, error) {
				if frontendStarts.Add(1) == 1 {
					return openManifestDevTestProcess(), nil
				}
				return nil, errors.New("synthetic frontend launch failure")
			}
			restored := make(chan struct{})
			ops.restoreFrontend = func(context.Context, string, manifest.Config, int, string, string) (*manifestProcess, error) {
				close(restored)
				if restoreFails {
					return nil, errors.New("synthetic restore failure")
				}
				return openManifestDevTestProcess(), nil
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			done := make(chan error, 1)
			go func() { done <- runManifestDevContextWithOps(ctx, &DevOptions{}, ops) }()
			watches.events <- manifestDevTestEvent{event: notify.Write, path: filepath.Join(watches.root, "main.go")}
			select {
			case <-restored:
			case <-time.After(2 * time.Second):
				t.Fatal("previous frontend was not restored")
			}
			if restoreFails {
				assert.ErrorContains(t, <-done, "synthetic restore failure")
				return
			}
			cancel()
			assert.NoError(t, <-done)
		})
	}
}

func TestManifestDevBackendRollbackRestoresFrontendAndWatcher(t *testing.T) {
	for _, restoreFails := range []bool{false, true} {
		t.Run(map[bool]string{false: "restored", true: "restore failure"}[restoreFails], func(t *testing.T) {
			ops, watches, restored := manifestDevBackendFailureOps(t, restoreFails)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			done := make(chan error, 1)
			go func() { done <- runManifestDevContextWithOps(ctx, &DevOptions{}, ops) }()
			watches.events <- manifestDevTestEvent{event: notify.Write, path: filepath.Join(watches.root, "main.go")}
			select {
			case <-restored:
			case <-time.After(2 * time.Second):
				t.Fatal("backend rollback did not restore the previous frontend")
			}
			if restoreFails {
				assert.ErrorContains(t, <-done, "synthetic backend restore failure")
				return
			}
			cancel()
			assert.NoError(t, <-done)
		})
	}
}

func manifestDevBackendFailureOps(t *testing.T, restoreFails bool) (manifestDevOps, *manifestWatchSet, <-chan struct{}) {
	t.Helper()
	ops, watches := newManifestDevTestOps(t)
	initial, err := ops.load("", "")
	require.NoError(t, err)
	next := cloneManifestDevLoaded(initial)
	next.Config.Dev.Watch = []string{"**/*.go"}
	next.Config.Frontend.Dev = []string{"replacement"}
	var loads atomic.Int32
	ops.load = func(string, string) (*manifest.Loaded, error) {
		if loads.Add(1) == 1 {
			return initial, nil
		}
		return next, nil
	}
	compileKey := pipelineCompileKey(runtime.GOOS, runtime.GOARCH)
	var builds atomic.Int32
	ops.build = func(context.Context, *DevOptions, *manifest.Loaded, string, string, string, int) (manifestPipelineRun, error) {
		run := manifestDevTestRun()
		if builds.Add(1) > 1 {
			run.Results[compileKey] = pipeline.Result{Status: cache.LookupMiss}
		}
		return run, nil
	}
	stagedWatches := &manifestWatchSet{root: watches.root, matcher: newDevWatchMatcher(next.Config.Dev.Watch)}
	var watchStarts atomic.Int32
	ops.startWatches = func(string, manifest.Config) (*manifestWatchSet, error) {
		if watchStarts.Add(1) == 1 {
			return watches, nil
		}
		return stagedWatches, nil
	}
	var appStarts atomic.Int32
	ops.startApp = func(string, string, string, int) (*manifestProcess, error) {
		if appStarts.Add(1) == 1 {
			return openManifestDevTestProcess(), nil
		}
		return nil, errors.New("synthetic backend launch failure")
	}
	restored := make(chan struct{})
	ops.restoreFrontend = func(context.Context, string, manifest.Config, int, string, string) (*manifestProcess, error) {
		close(restored)
		if restoreFails {
			return nil, errors.New("synthetic backend restore failure")
		}
		return openManifestDevTestProcess(), nil
	}
	return ops, watches, restored
}

func TestManifestDevRebuildHandlesPlanInvariantFailure(t *testing.T) {
	ops, watches := newManifestDevTestOps(t)
	compileKey := pipelineCompileKey(runtime.GOOS, runtime.GOARCH)
	var builds atomic.Int32
	ops.build = func(context.Context, *DevOptions, *manifest.Loaded, string, string, string, int) (manifestPipelineRun, error) {
		run := manifestDevTestRun()
		if builds.Add(1) > 1 {
			run.Results[compileKey] = pipeline.Result{Status: cache.LookupMiss}
		}
		return run, nil
	}
	failed := make(chan struct{})
	var paths atomic.Int32
	ops.binaryPath = func(string, manifestPipelineRun, string, string) (string, error) {
		if paths.Add(1) == 1 {
			return "adapter", nil
		}
		close(failed)
		return "", errors.New("synthetic rebuild Plan invariant")
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- runManifestDevContextWithOps(ctx, &DevOptions{}, ops) }()
	watches.events <- manifestDevTestEvent{event: notify.Write, path: filepath.Join(watches.root, "main.go")}
	select {
	case <-failed:
	case <-time.After(2 * time.Second):
		t.Fatal("rebuild Plan invariant was not checked")
	}
	cancel()
	assert.NoError(t, <-done)
}

func TestManifestDevCancellationDuringReplacementIsClean(t *testing.T) {
	for _, phase := range []string{"frontend", "backend"} {
		t.Run(phase, func(t *testing.T) {
			ops, watches := newManifestDevTestOps(t)
			initial, err := ops.load("", "")
			require.NoError(t, err)
			next := cloneManifestDevLoaded(initial)
			if phase == "frontend" {
				next.Config.Frontend.Dev = []string{"replacement"}
			}
			var loads atomic.Int32
			ops.load = func(string, string) (*manifest.Loaded, error) {
				if loads.Add(1) == 1 {
					return initial, nil
				}
				return next, nil
			}
			compileKey := pipelineCompileKey(runtime.GOOS, runtime.GOARCH)
			var builds atomic.Int32
			ops.build = func(context.Context, *DevOptions, *manifest.Loaded, string, string, string, int) (manifestPipelineRun, error) {
				run := manifestDevTestRun()
				if phase == "backend" && builds.Add(1) > 1 {
					run.Results[compileKey] = pipeline.Result{Status: cache.LookupMiss}
				}
				return run, nil
			}
			ctx, cancel := context.WithCancel(context.Background())
			if phase == "frontend" {
				var waits atomic.Int32
				ops.waitTCP = func(context.Context, *manifestProcess, string, time.Duration) error {
					if waits.Add(1) == 1 {
						return nil
					}
					cancel()
					return context.Canceled
				}
			} else {
				var waits atomic.Int32
				ops.waitStable = func(context.Context, *manifestProcess, time.Duration) error {
					if waits.Add(1) == 1 {
						return nil
					}
					cancel()
					return context.Canceled
				}
			}
			done := make(chan error, 1)
			go func() { done <- runManifestDevContextWithOps(ctx, &DevOptions{}, ops) }()
			watches.events <- manifestDevTestEvent{event: notify.Write, path: filepath.Join(watches.root, "main.go")}
			select {
			case err := <-done:
				assert.NoError(t, err)
			case <-time.After(2 * time.Second):
				t.Fatal("replacement cancellation did not stop the session")
			}
		})
	}
}

func TestManifestDevCancelsQueuedAndActiveRebuildsOnShutdown(t *testing.T) {
	ops, watches := newManifestDevTestOps(t)
	firstRebuild := make(chan struct{})
	releaseFirst := make(chan struct{})
	thirdBuild := make(chan struct{})
	var builds atomic.Int32
	ops.build = func(ctx context.Context, _ *DevOptions, _ *manifest.Loaded, _, _, _ string, _ int) (manifestPipelineRun, error) {
		switch builds.Add(1) {
		case 1:
			return manifestDevTestRun(), nil
		case 2:
			close(firstRebuild)
			<-releaseFirst
			return manifestDevTestRun(), ctx.Err()
		default:
			close(thirdBuild)
			<-ctx.Done()
			return manifestDevTestRun(), ctx.Err()
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- runManifestDevContextWithOps(ctx, &DevOptions{}, ops) }()
	event := manifestDevTestEvent{event: notify.Write, path: filepath.Join(watches.root, "main.go")}
	watches.events <- event
	select {
	case <-firstRebuild:
	case <-time.After(2 * time.Second):
		t.Fatal("first incremental rebuild did not start")
	}
	watches.events <- event
	time.Sleep(20 * time.Millisecond)
	watches.events <- event
	time.Sleep(20 * time.Millisecond)
	close(releaseFirst)
	select {
	case <-thirdBuild:
	case <-time.After(2 * time.Second):
		t.Fatal("latest incremental rebuild did not supersede the queued generation")
	}
	// Keep the latest result in flight while shutdown observes and cancels the
	// active generation.
	cancel()
	assert.NoError(t, <-done)
}

func newManifestDevTestOps(t *testing.T) (manifestDevOps, *manifestWatchSet) {
	t.Helper()
	root := t.TempDir()
	loaded := &manifest.Loaded{Config: manifest.Config{
		Root:     root,
		Project:  manifest.Project{Name: "adapter", BinaryName: "adapter"},
		Frontend: manifest.Frontend{Directory: "frontend", Dev: []string{"frontend"}},
		Dev:      manifest.Dev{DebounceMS: 1, GracePeriodMS: 1},
	}}
	watches := &manifestWatchSet{events: make(chan notify.EventInfo, 16), root: root, matcher: newDevWatchMatcher(nil)}
	ops := manifestDevOps{
		getwd:     func() (string, error) { return root, nil },
		load:      func(string, string) (*manifest.Loaded, error) { return loaded, nil },
		checkPort: func(string, int) error { return nil },
		build: func(context.Context, *DevOptions, *manifest.Loaded, string, string, string, int) (manifestPipelineRun, error) {
			return manifestDevTestRun(), nil
		},
		startFrontend: func(string, manifest.Config, string, int, string) (*manifestProcess, error) {
			return openManifestDevTestProcess(), nil
		},
		waitTCP: func(context.Context, *manifestProcess, string, time.Duration) error { return nil },
		binaryPath: func(string, manifestPipelineRun, string, string) (string, error) {
			return filepath.Join(root, "adapter"), nil
		},
		startApp:     func(string, string, string, int) (*manifestProcess, error) { return openManifestDevTestProcess(), nil },
		waitStable:   func(context.Context, *manifestProcess, time.Duration) error { return nil },
		startWatches: func(string, manifest.Config) (*manifestWatchSet, error) { return watches, nil },
		restartWatches: func(_ string, _ manifest.Config, current *manifestWatchSet) (*manifestWatchSet, error) {
			return current, nil
		},
		restoreFrontend: func(context.Context, string, manifest.Config, int, string, string) (*manifestProcess, error) {
			return openManifestDevTestProcess(), nil
		},
	}
	return ops, watches
}

func manifestDevTestRun() manifestPipelineRun {
	key := pipelineCompileKey(runtime.GOOS, runtime.GOARCH)
	return manifestPipelineRun{
		Plan:    pipeline.Plan{Nodes: map[pipeline.NodeKey]pipeline.Node{key: {Key: key, Kind: pipeline.CompileApplication, Output: ".wails/dev/adapter"}}},
		Results: map[pipeline.NodeKey]pipeline.Result{key: {Status: cache.LookupHit}},
	}
}

func openManifestDevTestProcess() *manifestProcess {
	return &manifestProcess{done: make(chan struct{})}
}

func cloneManifestDevLoaded(source *manifest.Loaded) *manifest.Loaded {
	copy := *source
	copy.Config = source.Config
	copy.Config.Dev.Watch = append([]string(nil), source.Config.Dev.Watch...)
	copy.Config.Frontend.Dev = append([]string(nil), source.Config.Frontend.Dev...)
	return &copy
}
