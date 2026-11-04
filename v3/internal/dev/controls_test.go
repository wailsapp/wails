package dev

import (
	"context"
	"errors"
	"io"
	"runtime"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/internal/report"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestSessionControlsPreserveAppOnFailureAndQuitCleansUp(t *testing.T) {
	root := t.TempDir()
	controls := make(chan string, 8)
	starts := make(chan *argumentProcess, 8)
	frontends := make(chan *argumentProcess, 8)
	failures := make(chan Event, 8)
	forced := make(chan bool, 8)
	output := NewOutput(io.Discard, report.Normal)
	output.Observe(func(event Event) {
		if event.Diagnostic != nil {
			failures <- event
		}
	})
	config := manifest.Config{Root: root, Dev: manifest.Dev{GracePeriodMS: 1}}
	ops := Operations{
		Getwd: func() (string, error) { return root, nil }, Load: func(string, string) (*manifest.Loaded, error) { return &manifest.Loaded{Config: config}, nil },
		Target: func(string) (string, string, error) { return runtime.GOOS, runtime.GOARCH, nil }, CheckPort: func(string, int) error { return nil },
		Build: func(ctx context.Context, _ *manifest.Loaded, _, _, _ string, _ int) (BuildResult, error) {
			forced <- IsForcedRebuild(ctx)
			return BuildResult{}, nil
		},
		GenerateBindings: func(context.Context, *manifest.Loaded) error { return errors.New("bindings fixture failed") },
		StartFrontend: func(string, manifest.Config, string, int, string) (Process, error) {
			p := newArgumentProcess(nil, false)
			frontends <- p
			return p, nil
		},
		WaitFrontendReady: func(context.Context, Process, string, time.Duration) error { return nil },
		BinaryPath:        func(string, BuildResult, string, string) (string, error) { return "app", nil },
		StartApp: func(string, string, string, int, []string) (Process, error) {
			p := newArgumentProcess(nil, false)
			starts <- p
			return p, nil
		},
		WaitReady:      func(context.Context, Process, time.Duration) error { return nil },
		BackendChanged: func(BuildResult, string, string) bool { return false },
	}
	ops.StartWatches = func(string, manifest.Config) (WatchSet, error) { return &argumentWatch{}, nil }
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- Run(ctx, Options{Controls: controls, Output: output}, ops) }()
	receive := func(channel <-chan *argumentProcess) *argumentProcess {
		t.Helper()
		select {
		case p := <-channel:
			return p
		case err := <-done:
			t.Fatalf("session exited: %v", err)
		case <-time.After(5 * time.Second):
			t.Fatal("session operation timed out")
		}
		return nil
	}
	initial := receive(starts)
	front := receive(frontends)
	if <-forced {
		t.Fatal("initial build was forced")
	}
	controls <- "b"
	select {
	case event := <-failures:
		if event.Diagnostic.Err.Error() != "bindings fixture failed" {
			t.Fatal(event.Diagnostic.Err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("bindings failure was not reported")
	}
	select {
	case <-initial.Done():
		t.Fatal("failed bindings stopped the current app")
	default:
	}
	controls <- "r"
	replacement := receive(starts)
	if !<-forced {
		t.Fatal("r did not bypass the build cache")
	}
	select {
	case <-initial.Done():
	case <-time.After(5 * time.Second):
		t.Fatal("unchanged backend was not restarted")
	}
	controls <- "f"
	nextFront := receive(frontends)
	select {
	case <-front.Done():
	case <-time.After(5 * time.Second):
		t.Fatal("frontend was not replaced")
	}
	controls <- "q"
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("quit did not finish")
	}
	for _, p := range []*argumentProcess{replacement, nextFront} {
		select {
		case <-p.Done():
		default:
			t.Fatal("quit left a process running")
		}
	}
}

func TestFrontendControlRestoresLastSessionAfterFailedReadiness(t *testing.T) {
	for _, restoreFails := range []bool{false, true} {
		t.Run(map[bool]string{false: "restored", true: "restore-failed"}[restoreFails], func(t *testing.T) {
			previous := newArgumentProcess(nil, false)
			candidate := newArgumentProcess(nil, false)
			restored := newArgumentProcess(nil, false)
			ops := Operations{
				StartFrontend:     func(string, manifest.Config, string, int, string) (Process, error) { return candidate, nil },
				WaitFrontendReady: func(context.Context, Process, string, time.Duration) error { return errors.New("frontend not ready") },
				RestoreFrontend: func(context.Context, string, manifest.Config, int, string, string) (Process, error) {
					if restoreFails {
						return nil, errors.New("restore failed")
					}
					return restored, nil
				},
			}
			next, err := restartSessionFrontend(t.Context(), previous, manifest.Config{}, frontendEndpoint{host: "127.0.0.1", port: 9245}, ops)
			if err == nil {
				t.Fatal("failed frontend must be reported")
			}
			for _, process := range []*argumentProcess{previous, candidate} {
				select {
				case <-process.Done():
				default:
					t.Fatal("frontend process was not stopped")
				}
			}
			if restoreFails {
				if next != nil {
					t.Fatal("unrecoverable frontend must end the session")
				}
			} else {
				if next != restored {
					t.Fatal("previous frontend policy was not restored")
				}
				select {
				case <-restored.Done():
					t.Fatal("restored frontend stopped")
				default:
				}
			}
		})
	}
}
