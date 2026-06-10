//go:build mcp

package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const enabled = true

// portEnvVar overrides the configured port at runtime.
const portEnvVar = "WAILS_MCP_PORT"

// mcpServer is the real MCP server implementation, compiled in with the `mcp`
// build tag.
type mcpServer struct {
	app    *application.App
	config Config
	logger *slog.Logger

	// route is the asset server route the service is mounted on, used as the
	// same-origin callback channel for JavaScript results.
	route string

	httpServer *http.Server
	addr       string

	pendingMu sync.Mutex
	pending   map[string]chan evalResult

	tools []*tool
}

func (s *Service) startup(_ context.Context, options application.ServiceOptions) error {
	app := application.Get()
	if app == nil {
		return errors.New("mcp: no application instance")
	}

	server := &mcpServer{
		app:     app,
		config:  s.config,
		logger:  app.Logger,
		route:   options.Route,
		pending: make(map[string]chan evalResult),
	}
	server.registerTools()

	port := s.config.Port
	if env := os.Getenv(portEnvVar); env != "" {
		p, err := strconv.Atoi(env)
		if err != nil {
			return fmt.Errorf("mcp: invalid %s value %q: %w", portEnvVar, env, err)
		}
		port = p
	}
	if port < 0 {
		port = 0 // pick a free port
	}

	listener, err := net.Listen("tcp", net.JoinHostPort(s.config.Host, strconv.Itoa(port)))
	if err != nil {
		return fmt.Errorf("mcp: failed to listen: %w", err)
	}
	server.addr = listener.Addr().String()

	mux := http.NewServeMux()
	mux.HandleFunc("/mcp", server.handleMCP)
	mux.HandleFunc("/eval-result", server.handleEvalResult)
	mux.HandleFunc("/", server.handleStatus)

	server.httpServer = &http.Server{Handler: mux}
	go func() {
		if err := server.httpServer.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			server.logger.Error("mcp: server error", "error", err)
		}
	}()

	s.server = server
	server.logger.Info("MCP server started. Connect MCP clients using the streamable HTTP transport.",
		"url", fmt.Sprintf("http://%s/mcp", server.addr))
	return nil
}

func (s *Service) shutdown() error {
	if s.server == nil {
		return nil
	}
	server := s.server
	s.server = nil
	server.failPending(errors.New("mcp: service shutting down"))
	return server.httpServer.Close()
}

// serveHTTP handles requests on the asset server route (same-origin with the
// window content, so no CORS is involved).
func (s *Service) serveHTTP(w http.ResponseWriter, r *http.Request) {
	if s.server == nil {
		http.Error(w, "MCP service not started", http.StatusServiceUnavailable)
		return
	}
	switch r.URL.Path {
	case "/eval-result":
		s.server.handleEvalResult(w, r)
	default:
		s.server.handleStatus(w, r)
	}
}

func (m *mcpServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"service":  "wails-mcp",
		"endpoint": fmt.Sprintf("http://%s/mcp", m.addr),
		"tools":    len(m.tools),
	})
}

// failPending unblocks all in-flight evaluations, e.g. on shutdown.
func (m *mcpServer) failPending(err error) {
	m.pendingMu.Lock()
	defer m.pendingMu.Unlock()
	for id, ch := range m.pending {
		select {
		case ch <- evalResult{Ok: false, Error: err.Error()}:
		default:
		}
		delete(m.pending, id)
	}
}

// resolveWindow returns the named window, or a sensible default: the focused
// window if any, otherwise the first window.
func (m *mcpServer) resolveWindow(name string) (application.Window, error) {
	if name != "" {
		window, ok := m.app.Window.GetByName(name)
		if !ok {
			return nil, fmt.Errorf("no window named %q", name)
		}
		return window, nil
	}
	windows := m.app.Window.GetAll()
	if len(windows) == 0 {
		return nil, errors.New("the application has no windows")
	}
	for _, window := range windows {
		if window.IsFocused() {
			return window, nil
		}
	}
	return windows[0], nil
}

// evalTimeout returns the timeout to use for a tool call, honouring an
// optional `timeout_ms` argument.
func (m *mcpServer) evalTimeout(args map[string]any) time.Duration {
	if ms, ok := argFloat(args, "timeout_ms"); ok && ms > 0 {
		return time.Duration(ms) * time.Millisecond
	}
	return m.config.EvalTimeout
}
