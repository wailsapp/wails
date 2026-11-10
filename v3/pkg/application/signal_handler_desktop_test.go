//go:build !windows && !ios && !android

package application

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

// signalQuitApp observes Quit without starting a native event loop.
type signalQuitApp struct {
	platformApp
	quit chan struct{}
}

func (a *signalQuitApp) isOnMainThread() bool { return true }
func (a *signalQuitApp) destroy()             { close(a.quit) }

func TestDefaultSignalHandler(t *testing.T) {
	for _, sig := range []syscall.Signal{syscall.SIGINT, syscall.SIGTERM} {
		for _, disabled := range []string{"false", "true"} {
			t.Run(sig.String()+"/disabled="+disabled, func(t *testing.T) {
				// Bound the child by the test's own deadline rather than a fixed
				// timeout, which a slow machine can use up starting the child.
				// The child times out first and prints its goroutines; the
				// context stops one that cannot.
				args := []string{"-test.run=^TestDefaultSignalHandlerProcess$"}
				ctx := context.Background()
				if deadline, ok := t.Deadline(); ok {
					left := time.Until(deadline)
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, left*9/10)
					defer cancel()
					args = append(args, "-test.timeout="+(left*8/10).String())
				}
				cmd := exec.CommandContext(ctx, os.Args[0], args...)
				cmd.Env = append(os.Environ(), "WAILS_TEST_SIGNAL="+sig.String(), "WAILS_TEST_SIGNAL_DISABLED="+disabled)
				output, err := cmd.CombinedOutput()
				if disabled == "false" {
					if err != nil {
						t.Fatalf("signal did not reach Quit: %v\n%s", err, output)
					}
					return
				}
				var exitErr *exec.ExitError
				if !errors.As(err, &exitErr) {
					t.Fatalf("expected default signal termination, got %v\n%s", err, output)
				}
				status, ok := exitErr.Sys().(syscall.WaitStatus)
				if !ok || !status.Signaled() || status.Signal() != sig {
					t.Fatalf("expected termination by %v, got %v\n%s", sig, err, output)
				}
			})
		}
	}
}

func TestDefaultSignalHandlerProcess(t *testing.T) {
	name := os.Getenv("WAILS_TEST_SIGNAL")
	if name == "" {
		return
	}
	var sig syscall.Signal
	switch name {
	case syscall.SIGINT.String():
		sig = syscall.SIGINT
	case syscall.SIGTERM.String():
		sig = syscall.SIGTERM
	default:
		t.Fatalf("unexpected signal %q", name)
	}

	impl := &signalQuitApp{quit: make(chan struct{})}
	app := &App{impl: impl, Logger: slog.Default()}
	globalApplication = app
	app.setupSignalHandler(Options{
		DisableDefaultSignalHandler: os.Getenv("WAILS_TEST_SIGNAL_DISABLED") == "true",
	})
	if err := syscall.Kill(os.Getpid(), sig); err != nil {
		t.Fatal(err)
	}
	// Bounded by the -test.timeout TestDefaultSignalHandler passes.
	<-impl.quit
}
