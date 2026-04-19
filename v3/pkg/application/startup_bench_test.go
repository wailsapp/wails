//go:build server

package application

import (
	"context"
	"net/http"
	"testing"
	"time"
)

// benchSleepyService blocks ServiceStartup for a configured duration to
// model user work that currently sits on the startup critical path.
type benchSleepyService struct {
	delay time.Duration
}

func (s *benchSleepyService) ServiceStartup(ctx context.Context, _ ServiceOptions) error {
	select {
	case <-time.After(s.delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// BenchmarkRun_NoServices measures the framework-only cost of App.Run() from
// the moment it is called to the moment the platform event loop begins
// accepting work. It uses server-mode platform impl so no GUI is spun up.
func BenchmarkRun_NoServices(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		globalApplication = nil

		app := New(Options{
			Name: "bench-no-services",
			Server: ServerOptions{
				Host: "127.0.0.1",
				// Port 0 lets the OS pick an unused port.
				Port: 0,
			},
			Assets: AssetOptions{
				Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
		})

		start := time.Now()
		errCh := make(chan error, 1)
		go func() { errCh <- app.Run() }()

		// Wait until the app has reached its main loop.
		waitForRunning(b, app, 2*time.Second)
		b.ReportMetric(float64(time.Since(start).Microseconds()), "μs/run")

		app.Quit()
		select {
		case <-errCh:
		case <-time.After(2 * time.Second):
			b.Fatal("timeout waiting for app shutdown")
		}
	}
}

// BenchmarkRun_FiveSleepyServices measures Run() time when five services each
// sleep in ServiceStartup. With today's sequential loop the durations are
// additive; this benchmark is the control that P2 (non-blocking services)
// needs to beat.
func BenchmarkRun_FiveSleepyServices(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		globalApplication = nil

		services := make([]Service, 5)
		for j := range services {
			services[j] = NewService(&benchSleepyService{delay: 50 * time.Millisecond})
		}

		app := New(Options{
			Name:     "bench-sleepy-services",
			Services: services,
			Server: ServerOptions{
				Host: "127.0.0.1",
				Port: 0,
			},
			Assets: AssetOptions{
				Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
		})

		start := time.Now()
		errCh := make(chan error, 1)
		go func() { errCh <- app.Run() }()

		waitForRunning(b, app, 5*time.Second)
		b.ReportMetric(float64(time.Since(start).Milliseconds()), "ms/run")

		app.Quit()
		select {
		case <-errCh:
		case <-time.After(2 * time.Second):
			b.Fatal("timeout waiting for app shutdown")
		}
	}
}

func waitForRunning(tb testing.TB, app *App, timeout time.Duration) {
	tb.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		app.runLock.Lock()
		running := app.running
		app.runLock.Unlock()
		if running {
			return
		}
		time.Sleep(time.Millisecond)
	}
	tb.Fatalf("app did not reach running state within %s", timeout)
}
