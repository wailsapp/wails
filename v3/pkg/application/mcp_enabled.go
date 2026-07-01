//go:build mcp && !ios && !android

package application

// MCP (Model Context Protocol) server — compiled in only with -tags mcp.
//
// The server starts automatically inside App.Run() when this tag is present;
// no user code changes are required. Configure with environment variables:
//
//	WAILS_MCP_HOST       bind address (default 127.0.0.1)
//	WAILS_MCP_PORT       port number  (default 9099; 0 = free port)
//	WAILS_MCP_TIMEOUT    JS eval timeout in milliseconds (default 30000)
//	WAILS_MCP_HIDE_CURSOR  set to 1/true to disable the animated cursor overlay

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	mcpHostEnvVar       = "WAILS_MCP_HOST"
	mcpPortEnvVar       = "WAILS_MCP_PORT"
	mcpTimeoutEnvVar    = "WAILS_MCP_TIMEOUT"
	mcpHideCursorEnvVar = "WAILS_MCP_HIDE_CURSOR"

	mcpDefaultHost    = "127.0.0.1"
	mcpDefaultPort    = 9099
	mcpDefaultTimeout = 30 * time.Second
)

type mcpServer struct {
	app    *App
	logger *slog.Logger

	hideCursor  bool
	evalTimeout time.Duration
	addr        string

	httpServer *http.Server

	pendingMu sync.Mutex
	pending   map[string]chan mcpEvalResult

	tools []*mcpTool
}

// startMCPServer reads configuration from environment variables, starts the
// MCP HTTP listener and registers a shutdown hook. Called from App.Run().
func startMCPServer(a *App) error {
	host := os.Getenv(mcpHostEnvVar)
	if host == "" {
		host = mcpDefaultHost
	}

	port := mcpDefaultPort
	if env := os.Getenv(mcpPortEnvVar); env != "" {
		p, err := strconv.Atoi(env)
		if err != nil {
			return fmt.Errorf("invalid %s value %q: %w", mcpPortEnvVar, env, err)
		}
		port = p
	}
	if port < 0 {
		return fmt.Errorf("invalid %s value: must be 0 (free port) or a positive port number, got %d", mcpPortEnvVar, port)
	}
	// port == 0 means ask the OS for a free port

	evalTimeout := mcpDefaultTimeout
	if env := os.Getenv(mcpTimeoutEnvVar); env != "" {
		ms, err := strconv.ParseInt(env, 10, 64)
		if err != nil || ms <= 0 {
			return fmt.Errorf("invalid %s value %q: must be a positive integer (milliseconds)", mcpTimeoutEnvVar, env)
		}
		evalTimeout = time.Duration(ms) * time.Millisecond
	}

	hideCursor := false
	switch strings.ToLower(strings.TrimSpace(os.Getenv(mcpHideCursorEnvVar))) {
	case "1", "true", "on", "yes":
		hideCursor = true
	}

	listener, err := net.Listen("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	server := &mcpServer{
		app:         a,
		logger:      a.Logger,
		hideCursor:  hideCursor,
		evalTimeout: evalTimeout,
		addr:        listener.Addr().String(),
		pending:     make(map[string]chan mcpEvalResult),
	}
	server.registerTools()

	mux := http.NewServeMux()
	mux.HandleFunc("/mcp", server.handleMCP)
	mux.HandleFunc("/eval-result", server.handleEvalResult)
	mux.HandleFunc("/", server.handleStatus)

	server.httpServer = &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      120 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	go func() {
		if err := server.httpServer.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			server.logger.Error("mcp: server error", "error", err)
		}
	}()

	a.OnShutdown(func() {
		server.failPending(errors.New("mcp: application shutting down"))
		_ = server.httpServer.Close()
	})

	server.logger.Info("MCP server started. Connect MCP clients using the streamable HTTP transport.",
		"url", fmt.Sprintf("http://%s/mcp", server.addr))
	return nil
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

// failPending unblocks all in-flight JS evaluations, e.g. on shutdown.
func (m *mcpServer) failPending(err error) {
	m.pendingMu.Lock()
	defer m.pendingMu.Unlock()
	for id, ch := range m.pending {
		select {
		case ch <- mcpEvalResult{Ok: false, Error: err.Error()}:
		default:
		}
		delete(m.pending, id)
	}
}

// resolveWindow returns the named window, or a sensible default: the focused
// window if any, otherwise the first window.
func (m *mcpServer) resolveWindow(name string) (Window, error) {
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

// mcpEvalTimeout returns the timeout to use for a tool call, honouring an
// optional timeout_ms argument.
func (m *mcpServer) mcpEvalTimeout(args map[string]any) time.Duration {
	if ms, ok := mcpArgFloat(args, "timeout_ms"); ok && ms > 0 {
		return time.Duration(ms) * time.Millisecond
	}
	return m.evalTimeout
}
