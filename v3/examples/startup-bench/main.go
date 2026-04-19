// The startup-bench example is a controlled workload for measuring Wails v3
// startup performance. It registers five services that each sleep 200ms in
// ServiceStartup (simulating DB/auth/network work) and renders a trivial
// index.html.
//
// Build and run with the wails_trace_startup tag to emit a Chrome-trace JSON
// file:
//
//	WAILS_TRACE_STARTUP_OUTPUT=trace.json \
//	    go run -tags wails_trace_startup ./v3/examples/startup-bench
//
// Load trace.json in chrome://tracing or https://ui.perfetto.dev to inspect
// startup timing.
package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/*
var assets embed.FS

type slowService struct {
	name  string
	delay time.Duration
}

func (s *slowService) ServiceName() string { return s.name }

func (s *slowService) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	select {
	case <-time.After(s.delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Ping lets the frontend call into the service to prove it's bound.
func (s *slowService) Ping() string { return "pong from " + s.name }

func main() {
	svcs := make([]application.Service, 5)
	for i := range svcs {
		svcs[i] = application.NewService(&slowService{
			name:  fmt.Sprintf("slow-%d", i),
			delay: 200 * time.Millisecond,
		})
	}

	app := application.New(application.Options{
		Name:     "startup-bench",
		Services: svcs,
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Wails v3 startup benchmark",
		URL:   "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
