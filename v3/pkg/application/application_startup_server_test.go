//go:build server

package application

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/pkg/events"
)

type blockingStartupService struct {
	started chan struct{}
	release chan struct{}
}

func (s *blockingStartupService) ServiceStartup(ctx context.Context, _ ServiceOptions) error {
	close(s.started)
	select {
	case <-s.release:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func TestApplicationLifecycleEventsSurroundServiceStartup(t *testing.T) {
	resetGlobalApp()
	t.Cleanup(resetGlobalApp)

	service := &blockingStartupService{
		started: make(chan struct{}),
		release: make(chan struct{}),
	}
	app := New(Options{
		Name:     "Lifecycle startup test",
		Services: []Service{NewService(service)},
		Server: ServerOptions{
			Host: "127.0.0.1",
			Port: 18084,
		},
		Assets: AssetOptions{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
		},
	})

	starting := make(chan struct{})
	initialized := make(chan struct{})
	app.Event.OnApplicationEvent(events.Common.ApplicationStarting, func(*ApplicationEvent) {
		close(starting)
	})
	app.Event.OnApplicationEvent(events.Common.ApplicationInitialized, func(*ApplicationEvent) {
		close(initialized)
	})

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run()
	}()

	waitForSignal(t, starting, "ApplicationStarting")
	waitForSignal(t, service.started, "service startup")

	select {
	case <-initialized:
		t.Fatal("ApplicationInitialized was emitted before service startup completed")
	default:
	}

	resp, err := http.Get("http://127.0.0.1:18084/health")
	if err != nil {
		t.Fatalf("application was unresponsive during service startup: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("health status during service startup = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	close(service.release)
	waitForSignal(t, initialized, "ApplicationInitialized")
	app.Quit()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("app.Run() returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for app shutdown")
	}
}

func TestQuitCancelsServiceStartupBeforeWaitingForShutdown(t *testing.T) {
	resetGlobalApp()
	t.Cleanup(resetGlobalApp)

	service := &blockingStartupService{
		started: make(chan struct{}),
		release: make(chan struct{}),
	}
	app := New(Options{
		Name:     "Lifecycle cancellation test",
		Services: []Service{NewService(service)},
		Server: ServerOptions{
			Host: "127.0.0.1",
			Port: 18085,
		},
		Assets: AssetOptions{
			Handler: http.NotFoundHandler(),
		},
	})

	initialized := make(chan struct{})
	app.Event.OnApplicationEvent(events.Common.ApplicationInitialized, func(*ApplicationEvent) {
		close(initialized)
	})

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run()
	}()

	waitForSignal(t, service.started, "service startup")
	app.Quit()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("app.Run() returned nil after service startup was cancelled")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("shutdown deadlocked waiting for service startup")
	}

	select {
	case <-initialized:
		t.Fatal("ApplicationInitialized was emitted after service startup was cancelled")
	default:
	}
}

func waitForSignal(t *testing.T, signal <-chan struct{}, name string) {
	t.Helper()
	select {
	case <-signal:
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for %s", name)
	}
}
