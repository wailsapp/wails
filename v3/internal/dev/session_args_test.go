package dev

import (
	"context"
	"errors"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rjeczalik/notify"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

type argumentProcess struct {
	args      []string
	stopFirst bool
	done      chan struct{}
	once      sync.Once
}

func (p *argumentProcess) Done() <-chan struct{}            { return p.done }
func (p *argumentProcess) Stop(time.Duration)               { p.once.Do(func() { close(p.done) }) }
func (*argumentProcess) Err() error                         { return nil }
func (p *argumentProcess) NeedsStopBeforeReplacement() bool { return p.stopFirst }
func newArgumentProcess(args []string, stopFirst bool) *argumentProcess {
	return &argumentProcess{args: append([]string{}, args...), stopFirst: stopFirst, done: make(chan struct{})}
}

type argumentWatch struct{ events chan notify.EventInfo }

func (w *argumentWatch) Events() <-chan notify.EventInfo                    { return w.events }
func (*argumentWatch) Stop()                                                {}
func (*argumentWatch) Classify(manifest.Config, notify.EventInfo) EventKind { return Source }
func (*argumentWatch) ContainsInput(manifest.Config, string) bool           { return true }

type argumentEvent string

func (e argumentEvent) Path() string      { return string(e) }
func (argumentEvent) Event() notify.Event { return notify.Write }
func (argumentEvent) Sys() any            { return nil }

// Exercise both overlapping desktop replacement and the stop/restore lifecycle
// used by Windows and mobile, without needing either platform's host tools.
func TestDevelopmentArgumentsSurviveRebuildAndFailedReplacement(t *testing.T) {
	for _, stopFirst := range []bool{false, true} {
		t.Run(map[bool]string{false: "overlapping", true: "stop-and-restore"}[stopFirst], func(t *testing.T) {
			root := t.TempDir()
			var loaded atomic.Pointer[manifest.Loaded]
			setArgs := func(args ...string) {
				loaded.Store(&manifest.Loaded{Config: manifest.Config{Root: root, Dev: manifest.Dev{Args: args, ArgsSet: true, DebounceMS: 1, GracePeriodMS: 1}}})
			}
			setArgs("--config-path", "first profile.yaml")
			watches := &argumentWatch{events: make(chan notify.EventInfo, 8)}
			starts := make(chan *argumentProcess, 8)
			restores := make(chan *argumentProcess, 8)
			var changed atomic.Bool
			ops := Operations{
				Getwd: func() (string, error) { return root, nil }, Load: func(string, string) (*manifest.Loaded, error) { return loaded.Load(), nil },
				Target: func(string) (string, string, error) { return runtime.GOOS, runtime.GOARCH, nil }, CheckPort: func(string, int) error { return nil },
				Build: func(context.Context, *manifest.Loaded, string, string, string, int) (BuildResult, error) {
					name := "cached"
					if changed.Load() {
						name = "changed"
					}
					return BuildResult{Plan: pipeline.Plan{Name: name}}, nil
				},
				StartFrontend: func(string, manifest.Config, string, int, string) (Process, error) {
					return newArgumentProcess(nil, false), nil
				},
				WaitFrontendReady: func(context.Context, Process, string, time.Duration) error { return nil },
				BinaryPath:        func(string, BuildResult, string, string) (string, error) { return "app", nil },
				StartApp: func(_, _, _ string, _ int, args []string) (Process, error) {
					p := newArgumentProcess(args, stopFirst)
					starts <- p
					return p, nil
				},
				WaitReady: func(_ context.Context, p Process, _ time.Duration) error {
					if p.(*argumentProcess).args[0] == "bad" {
						return errors.New("rejected arguments")
					}
					return nil
				},
				StartWatches:   func(string, manifest.Config) (WatchSet, error) { return watches, nil },
				BackendChanged: func(r BuildResult, _, _ string) bool { return r.Plan.Name == "changed" },
				RestoreApp: func(_ context.Context, _ string, previous Process, _ string, _ int) (Process, error) {
					p := newArgumentProcess(previous.(*argumentProcess).args, stopFirst)
					restores <- p
					return p, nil
				},
			}
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			done := make(chan error, 1)
			go func() { done <- Run(ctx, Options{}, ops) }()
			receive := func(ch <-chan *argumentProcess) *argumentProcess {
				select {
				case p := <-ch:
					return p
				case err := <-done:
					t.Fatalf("session exited: %v", err)
				case <-time.After(5 * time.Second):
					t.Fatal("launch timed out")
				}
				return nil
			}
			trigger := func() { watches.events <- argumentEvent(filepath.Join(root, "wails.hcl")) }
			initial := receive(starts)
			require.Equal(t, []string{"--config-path", "first profile.yaml"}, initial.args)
			changed.Store(true)
			trigger()
			rebuilt := receive(starts)
			require.Equal(t, initial.args, rebuilt.args)
			changed.Store(false)
			setArgs("--config-path", "changed profile.yaml")
			trigger()
			active := receive(starts)
			require.Equal(t, []string{"--config-path", "changed profile.yaml"}, active.args)
			setArgs("bad")
			trigger()
			failed := receive(starts)
			select {
			case <-failed.Done():
			case <-time.After(5 * time.Second):
				t.Fatal("failed candidate was not stopped")
			}
			if stopFirst {
				active = receive(restores)
				require.Equal(t, []string{"--config-path", "changed profile.yaml"}, active.args)
			}
			select {
			case <-active.Done():
				t.Fatal("previous arguments were not retained after failure")
			default:
			}
			setArgs("--config-path", "recovered.yaml")
			trigger()
			recovered := receive(starts)
			require.Equal(t, []string{"--config-path", "recovered.yaml"}, recovered.args)
			cancel()
			require.NoError(t, <-done)
		})
	}
}

func TestAndroidArgumentsAreRejectedBeforeDeviceEffects(t *testing.T) {
	ops := Operations{
		Getwd: func() (string, error) { return t.TempDir(), nil },
		Load: func(string, string) (*manifest.Loaded, error) {
			return &manifest.Loaded{Config: manifest.Config{Dev: manifest.Dev{Args: []string{"--config-path", "testing.yaml"}}}}, nil
		},
		Target: func(string) (string, string, error) { return "android", "arm64", nil },
		ValidateTarget: func(string, string) error {
			t.Fatal("device discovery must not run for unsupported arguments")
			return nil
		},
	}
	require.ErrorContains(t, Run(t.Context(), Options{Plan: true}, ops), "Android development does not support")
}

func TestApplicationArgumentOverridesSurviveConfigurationChanges(t *testing.T) {
	for _, override := range [][]string{{"--config-path", "CLI profile.yaml"}, {}} {
		for _, configured := range [][]string{{"first"}, {"second"}, {}} {
			config := manifest.Config{Dev: manifest.Dev{Args: configured}}
			args, err := applicationArguments(config, "linux", "amd64", override)
			require.NoError(t, err)
			require.Equal(t, override, args)
			if len(args) > 0 {
				args[0] = "changed"
				require.Equal(t, "--config-path", override[0], "resolution must not mutate the CLI vector")
			}
		}
	}
}
