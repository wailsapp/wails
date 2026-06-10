// Package mcp provides a Model Context Protocol (MCP) server for Wails
// applications. It lets MCP clients such as LLM-powered coding agents inspect,
// test and control a running application: listing and controlling windows,
// evaluating JavaScript, querying the DOM, simulating mouse and keyboard input
// with an animated on-screen cursor, calling bound service methods and working
// with application events.
//
// The service is compiled in only when the `mcp` build tag is set. Without the
// tag it is a no-op stub, so it can be registered unconditionally:
//
//	app := application.New(application.Options{
//	    Services: []application.Service{
//	        application.NewServiceWithOptions(mcp.New(), application.ServiceOptions{
//	            Route: "/wails-mcp",
//	        }),
//	    },
//	})
//
// Setting the WAILS_MCP environment variable to a truthy value when running
// `wails3 build` or `wails3 dev` adds the `mcp` build tag automatically.
//
// When enabled, the service listens on http://127.0.0.1:9099/mcp (configurable
// via Config or the WAILS_MCP_PORT environment variable) using the streamable
// HTTP transport, e.g. for Claude Code:
//
//	claude mcp add --transport http my-wails-app http://127.0.0.1:9099/mcp
//
// The optional Route mounts the service on the application's internal asset
// server. It is used as a same-origin callback channel for JavaScript results;
// without it the service falls back to a CORS-enabled localhost endpoint.
package mcp

import (
	"context"
	"net/http"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Config configures the MCP service.
type Config struct {
	// Host is the interface the MCP server listens on.
	// Defaults to 127.0.0.1. Binding to other interfaces exposes control of
	// the application to the network - do this only in trusted environments.
	Host string

	// Port is the port the MCP server listens on. Defaults to 9099.
	// Use -1 to pick a free port. The WAILS_MCP_PORT environment variable
	// overrides this at runtime.
	Port int

	// EvalTimeout is the default time to wait for JavaScript executed in a
	// window to report its result. Defaults to 30 seconds.
	EvalTimeout time.Duration

	// HideCursor disables the animated cursor overlay shown for simulated
	// mouse input.
	HideCursor bool
}

// Service is the MCP service. Use New or NewWithConfig to create one.
type Service struct {
	config Config
	server *mcpServer
}

// New creates a new MCP service with default configuration.
func New() *Service {
	return NewWithConfig(Config{})
}

// NewWithConfig creates a new MCP service with the given configuration.
func NewWithConfig(config Config) *Service {
	if config.Host == "" {
		config.Host = "127.0.0.1"
	}
	if config.Port == 0 {
		config.Port = 9099
	}
	if config.EvalTimeout <= 0 {
		config.EvalTimeout = 30 * time.Second
	}
	return &Service{config: config}
}

// ServiceName returns the name of the service.
func (s *Service) ServiceName() string {
	return "github.com/wailsapp/wails/v3/services/mcp"
}

// ServiceStartup is called when the application starts.
func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	return s.startup(ctx, options)
}

// ServiceShutdown is called when the application shuts down.
func (s *Service) ServiceShutdown() error {
	return s.shutdown()
}

// ServeHTTP handles requests on the service's asset server route. It is used
// as a same-origin callback channel for JavaScript evaluation results.
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveHTTP(w, r)
}

// Enabled reports whether the application was built with the `mcp` build tag.
func (s *Service) Enabled() bool {
	return enabled
}
