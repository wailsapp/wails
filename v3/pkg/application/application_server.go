//go:build server

package application

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

// serverApp implements platformApp for server mode.
// It provides a minimal implementation that runs an HTTP server
// without any GUI components.
//
// Server mode is enabled by building with the "server" build tag:
//
//	go build -tags server
type serverApp struct {
	app         *App
	server      *http.Server
	listener    net.Listener
	broadcaster *WebSocketBroadcaster
}

// newPlatformApp creates a new server-mode platform app.
// This function is only compiled when building with the "server" tag.
func newPlatformApp(app *App) *serverApp {
	app.info("Server mode enabled (built with -tags server)")
	return &serverApp{
		app: app,
	}
}

// parsePort parses a port string into an integer.
func parsePort(s string) (int, error) {
	p, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if p < 1 || p > 65535 {
		return 0, errors.New("port out of range")
	}
	return p, nil
}

// run starts the HTTP server and blocks until shutdown.
func (h *serverApp) run() error {
	// Set up common events
	h.setupCommonEvents()

	// Create WebSocket broadcaster for events
	h.broadcaster = NewWebSocketBroadcaster(h.app)
	globalBroadcaster = h.broadcaster // Set global reference for browser ID lookups
	h.app.wailsEventListenerLock.Lock()
	h.app.wailsEventListeners = append(h.app.wailsEventListeners, h.broadcaster)
	h.app.wailsEventListenerLock.Unlock()

	opts := h.app.options.Server

	// Environment variables override config (useful for Docker/containers)
	host := os.Getenv("WAILS_SERVER_HOST")
	if host == "" {
		host = opts.Host
	}
	if host == "" {
		host = "localhost"
	}

	port := opts.Port
	if envPort := os.Getenv("WAILS_SERVER_PORT"); envPort != "" {
		if p, err := parsePort(envPort); err == nil {
			port = p
		}
	}
	if port == 0 {
		port = 8080
	}

	readTimeout := opts.ReadTimeout
	if readTimeout == 0 {
		readTimeout = 30 * time.Second
	}

	writeTimeout := opts.WriteTimeout
	if writeTimeout == 0 {
		writeTimeout = 30 * time.Second
	}

	idleTimeout := opts.IdleTimeout
	if idleTimeout == 0 {
		idleTimeout = 120 * time.Second
	}

	shutdownTimeout := opts.ShutdownTimeout
	if shutdownTimeout == 0 {
		shutdownTimeout = 30 * time.Second
	}

	addr := fmt.Sprintf("%s:%d", host, port)

	// Create HTTP handler from asset server
	handler := h.createHandler()

	h.server = &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	// Create listener
	var err error
	h.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	h.app.info("Server mode starting", "address", addr)

	// Start server in goroutine
	errCh := make(chan error, 1)
	go func() {
		if opts.TLS != nil {
			errCh <- h.server.ServeTLS(h.listener, opts.TLS.CertFile, opts.TLS.KeyFile)
		} else {
			errCh <- h.server.Serve(h.listener)
		}
	}()

	// Wait for shutdown signal or error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-quit:
		h.app.info("Shutdown signal received")
	case <-h.app.ctx.Done():
		h.app.info("Application context cancelled")
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := h.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	h.app.info("Server stopped")
	return nil
}

// customJS is the JavaScript that sets up WebSocket event connection for server mode.
// Events FROM frontend TO backend use the existing HTTP transport.
// This WebSocket is only for receiving broadcast events FROM backend TO all frontends.
const customJS = `(function() {
	var protocol = location.protocol === 'https:' ? 'wss:' : 'ws:';
	var clientId = window._wails && window._wails.clientId ? window._wails.clientId : '';
	var wsUrl = protocol + '//' + location.host + '/wails/events' + (clientId ? '?clientId=' + encodeURIComponent(clientId) : '');
	var ws;

	function connect() {
		ws = new WebSocket(wsUrl);
		ws.onopen = function() {
			console.log('[Wails] Event WebSocket connected');
		};
		ws.onmessage = function(e) {
			try {
				var event = JSON.parse(e.data);
				if (window._wails && window._wails.dispatchWailsEvent) {
					window._wails.dispatchWailsEvent(event);
				}
			} catch (err) {
				console.error('[Wails] Failed to parse event:', err);
			}
		};
		ws.onclose = function() {
			console.log('[Wails] Event WebSocket disconnected, reconnecting...');
			setTimeout(connect, 1000);
		};
		ws.onerror = function() {
			ws.close();
		};
	}

	connect();
})();`

// createHandler creates the HTTP handler for server mode.
func (h *serverApp) createHandler() http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Serve custom.js for server mode (WebSocket event connection)
	mux.HandleFunc("/wails/custom.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(customJS))
	})

	// WebSocket endpoint for events
	mux.Handle("/wails/events", h.broadcaster)

	// Serve all other requests through the asset server
	mux.Handle("/", h.app.assets)

	return mux
}

// destroy stops the server and cleans up.
func (h *serverApp) destroy() {
	if h.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		h.server.Shutdown(ctx)
	}
	h.app.cleanup()
}

// setApplicationMenu is a no-op in server mode.
func (h *serverApp) setApplicationMenu(menu *Menu) {
	// No-op: server mode has no GUI
}

// name returns the application name.
func (h *serverApp) name() string {
	return h.app.options.Name
}

// getCurrentWindowID returns 0 in server mode (no windows).
func (h *serverApp) getCurrentWindowID() uint {
	return 0
}

// showAboutDialog is a no-op in server mode.
func (h *serverApp) showAboutDialog(name string, description string, icon []byte) {
	// No-op: server mode has no GUI
	h.app.warning("showAboutDialog called in server mode - operation ignored")
}

// setIcon is a no-op in server mode.
func (h *serverApp) setIcon(icon []byte) {
	// No-op: server mode has no GUI
}

// on is a no-op in server mode.
func (h *serverApp) on(id uint) {
	// No-op: server mode has no platform-specific event handling
}

// dispatchOnMainThread executes the function directly in server mode.
func (h *serverApp) dispatchOnMainThread(id uint) {
	// In server mode, there's no "main thread" concept from GUI frameworks
	// Execute the function directly
	mainThreadFunctionStoreLock.Lock()
	fn, ok := mainThreadFunctionStore[id]
	if ok {
		delete(mainThreadFunctionStore, id)
	}
	mainThreadFunctionStoreLock.Unlock()

	if ok && fn != nil {
		fn()
	}
}

// hide is a no-op in server mode.
func (h *serverApp) hide() {
	// No-op: server mode has no GUI
}

// show is a no-op in server mode.
func (h *serverApp) show() {
	// No-op: server mode has no GUI
}

// getPrimaryScreen returns nil in server mode.
func (h *serverApp) getPrimaryScreen() (*Screen, error) {
	return nil, errors.New("screen information not available in server mode")
}

// getScreens returns an error in server mode (screen info unavailable).
func (h *serverApp) getScreens() ([]*Screen, error) {
	return nil, errors.New("screen information not available in server mode")
}

// GetFlags returns the application flags for server mode.
func (h *serverApp) GetFlags(options Options) map[string]any {
	flags := make(map[string]any)
	flags["server"] = true
	if options.Flags != nil {
		for k, v := range options.Flags {
			flags[k] = v
		}
	}
	return flags
}

// isOnMainThread always returns true in server mode.
func (h *serverApp) isOnMainThread() bool {
	// In server mode, there's no main thread concept
	return true
}

// isDarkMode returns false in server mode.
func (h *serverApp) isDarkMode() bool {
	return false
}

// getAccentColor returns empty string in server mode.
func (h *serverApp) getAccentColor() string {
	return ""
}

// logPlatformInfo logs platform info for server mode.
func (a *App) logPlatformInfo() {
	a.info("Platform Info:", "mode", "server")
}

// platformEnvironment returns environment info for server mode.
func (a *App) platformEnvironment() map[string]any {
	return map[string]any{
		"mode": "server",
	}
}

// fatalHandler sets up fatal error handling for server mode.
func fatalHandler(errFunc func(error)) {
	// In server mode, fatal errors are handled via standard mechanisms
}

// newClipboardImpl creates a clipboard implementation for server mode.
func newClipboardImpl() clipboardImpl {
	return &serverClipboard{}
}

// serverClipboard is a no-op clipboard for server mode.
type serverClipboard struct{}

func (c *serverClipboard) setText(text string) bool {
	return false
}

func (c *serverClipboard) text() (string, bool) {
	return "", false
}

// newDialogImpl creates a dialog implementation for server mode.
func newDialogImpl(d *MessageDialog) messageDialogImpl {
	return &serverDialog{}
}

// serverDialog is a no-op dialog for server mode.
type serverDialog struct{}

func (d *serverDialog) show() {
	// No-op in server mode
}

// newOpenFileDialogImpl creates an open file dialog implementation for server mode.
func newOpenFileDialogImpl(d *OpenFileDialogStruct) openFileDialogImpl {
	return &serverOpenFileDialog{}
}

// serverOpenFileDialog is a no-op open file dialog for server mode.
type serverOpenFileDialog struct{}

func (d *serverOpenFileDialog) show() (chan string, error) {
	ch := make(chan string, 1)
	close(ch)
	return ch, errors.New("file dialogs not available in server mode")
}

// newSaveFileDialogImpl creates a save file dialog implementation for server mode.
func newSaveFileDialogImpl(d *SaveFileDialogStruct) saveFileDialogImpl {
	return &serverSaveFileDialog{}
}

// serverSaveFileDialog is a no-op save file dialog for server mode.
type serverSaveFileDialog struct{}

func (d *serverSaveFileDialog) show() (chan string, error) {
	ch := make(chan string, 1)
	close(ch)
	return ch, errors.New("file dialogs not available in server mode")
}

// newMenuImpl creates a menu implementation for server mode.
func newMenuImpl(menu *Menu) menuImpl {
	return &serverMenu{}
}

// serverMenu is a no-op menu for server mode.
type serverMenu struct{}

func (m *serverMenu) update() {
	// No-op in server mode
}

// newPlatformLock creates a platform-specific single instance lock for server mode.
func newPlatformLock(manager *singleInstanceManager) (platformLock, error) {
	return &serverLock{}, nil
}

// serverLock is a basic lock for server mode.
type serverLock struct{}

func (l *serverLock) acquire(uniqueID string) error {
	return nil
}

func (l *serverLock) release() {
	// No-op in server mode
}

func (l *serverLock) notify(data string) error {
	return errors.New("single instance not supported in server mode")
}

// newSystemTrayImpl creates a system tray implementation for server mode.
func newSystemTrayImpl(s *SystemTray) systemTrayImpl {
	return &serverSystemTray{parent: s}
}

// serverSystemTray is a no-op system tray for server mode.
type serverSystemTray struct {
	parent *SystemTray
}

func (t *serverSystemTray) setLabel(label string)             {}
func (t *serverSystemTray) setTooltip(tooltip string)         {}
func (t *serverSystemTray) run()                              {}
func (t *serverSystemTray) setIcon(icon []byte)               {}
func (t *serverSystemTray) setMenu(menu *Menu)                {}
func (t *serverSystemTray) setIconPosition(pos IconPosition)  {}
func (t *serverSystemTray) setTemplateIcon(icon []byte)       {}
func (t *serverSystemTray) destroy()                          {}
func (t *serverSystemTray) setDarkModeIcon(icon []byte)       {}
func (t *serverSystemTray) bounds() (*Rect, error)            { return nil, errors.New("system tray not available in server mode") }
func (t *serverSystemTray) getScreen() (*Screen, error)       { return nil, errors.New("system tray not available in server mode") }
func (t *serverSystemTray) positionWindow(w Window, o int) error { return errors.New("system tray not available in server mode") }
func (t *serverSystemTray) openMenu()                         {}
func (t *serverSystemTray) Show()                             {}
func (t *serverSystemTray) Hide()                             {}

// newWindowImpl creates a webview window implementation for server mode.
func newWindowImpl(parent *WebviewWindow) *serverWebviewWindow {
	return &serverWebviewWindow{parent: parent}
}

// serverWebviewWindow is a no-op webview window for server mode.
type serverWebviewWindow struct {
	parent *WebviewWindow
}

// All webviewWindowImpl methods as no-ops for server mode
func (w *serverWebviewWindow) setTitle(title string)                      {}
func (w *serverWebviewWindow) setSize(width, height int)                  {}
func (w *serverWebviewWindow) setAlwaysOnTop(alwaysOnTop bool)            {}
func (w *serverWebviewWindow) setURL(url string)                          {}
func (w *serverWebviewWindow) setResizable(resizable bool)                {}
func (w *serverWebviewWindow) setMinSize(width, height int)               {}
func (w *serverWebviewWindow) setMaxSize(width, height int)               {}
func (w *serverWebviewWindow) execJS(js string)                           {}
func (w *serverWebviewWindow) setBackgroundColour(color RGBA)             {}
func (w *serverWebviewWindow) run()                                       {}
func (w *serverWebviewWindow) center()                                    {}
func (w *serverWebviewWindow) size() (int, int)                           { return 0, 0 }
func (w *serverWebviewWindow) width() int                                 { return 0 }
func (w *serverWebviewWindow) height() int                                { return 0 }
func (w *serverWebviewWindow) destroy()                                   {}
func (w *serverWebviewWindow) reload()                                    {}
func (w *serverWebviewWindow) forceReload()                               {}
func (w *serverWebviewWindow) openDevTools()                              {}
func (w *serverWebviewWindow) zoomReset()                                 {}
func (w *serverWebviewWindow) zoomIn()                                    {}
func (w *serverWebviewWindow) zoomOut()                                   {}
func (w *serverWebviewWindow) getZoom() float64                           { return 1.0 }
func (w *serverWebviewWindow) setZoom(zoom float64)                       {}
func (w *serverWebviewWindow) close()                                     {}
func (w *serverWebviewWindow) zoom()                                      {}
func (w *serverWebviewWindow) setHTML(html string)                        {}
func (w *serverWebviewWindow) on(eventID uint)                            {}
func (w *serverWebviewWindow) minimise()                                  {}
func (w *serverWebviewWindow) unminimise()                                {}
func (w *serverWebviewWindow) maximise()                                  {}
func (w *serverWebviewWindow) unmaximise()                                {}
func (w *serverWebviewWindow) fullscreen()                                {}
func (w *serverWebviewWindow) unfullscreen()                              {}
func (w *serverWebviewWindow) isMinimised() bool                          { return false }
func (w *serverWebviewWindow) isMaximised() bool                          { return false }
func (w *serverWebviewWindow) isFullscreen() bool                         { return false }
func (w *serverWebviewWindow) isNormal() bool                             { return true }
func (w *serverWebviewWindow) isVisible() bool                            { return false }
func (w *serverWebviewWindow) isFocused() bool                            { return false }
func (w *serverWebviewWindow) focus()                                     {}
func (w *serverWebviewWindow) show()                                      {}
func (w *serverWebviewWindow) hide()                                      {}
func (w *serverWebviewWindow) getScreen() (*Screen, error)                { return nil, errors.New("screens not available in server mode") }
func (w *serverWebviewWindow) setFrameless(frameless bool)                {}
func (w *serverWebviewWindow) openContextMenu(menu *Menu, data *ContextMenuData) {}
func (w *serverWebviewWindow) nativeWindow() unsafe.Pointer               { return nil }
func (w *serverWebviewWindow) startDrag() error                           { return errors.New("drag not available in server mode") }
func (w *serverWebviewWindow) startResize(border string) error            { return errors.New("resize not available in server mode") }
func (w *serverWebviewWindow) print() error                               { return errors.New("print not available in server mode") }
func (w *serverWebviewWindow) setEnabled(enabled bool)                    {}
func (w *serverWebviewWindow) physicalBounds() Rect                       { return Rect{} }
func (w *serverWebviewWindow) setPhysicalBounds(bounds Rect)              {}
func (w *serverWebviewWindow) bounds() Rect                               { return Rect{} }
func (w *serverWebviewWindow) setBounds(bounds Rect)                      {}
func (w *serverWebviewWindow) position() (int, int)                       { return 0, 0 }
func (w *serverWebviewWindow) setPosition(x int, y int)                   {}
func (w *serverWebviewWindow) relativePosition() (int, int)               { return 0, 0 }
func (w *serverWebviewWindow) setRelativePosition(x int, y int)           {}
func (w *serverWebviewWindow) flash(enabled bool)                         {}
func (w *serverWebviewWindow) handleKeyEvent(acceleratorString string)    {}
func (w *serverWebviewWindow) getBorderSizes() *LRTB                      { return &LRTB{} }
func (w *serverWebviewWindow) setMinimiseButtonState(state ButtonState)   {}
func (w *serverWebviewWindow) setMaximiseButtonState(state ButtonState)   {}
func (w *serverWebviewWindow) setCloseButtonState(state ButtonState)      {}
func (w *serverWebviewWindow) isIgnoreMouseEvents() bool                  { return false }
func (w *serverWebviewWindow) setIgnoreMouseEvents(ignore bool)           {}
func (w *serverWebviewWindow) cut()                                       {}
func (w *serverWebviewWindow) copy()                                      {}
func (w *serverWebviewWindow) paste()                                     {}
func (w *serverWebviewWindow) undo()                                      {}
func (w *serverWebviewWindow) delete()                                    {}
func (w *serverWebviewWindow) selectAll()                                 {}
func (w *serverWebviewWindow) redo()                                      {}
func (w *serverWebviewWindow) showMenuBar()                               {}
func (w *serverWebviewWindow) hideMenuBar()                               {}
func (w *serverWebviewWindow) toggleMenuBar()                             {}
func (w *serverWebviewWindow) setMenu(menu *Menu)                         {}
func (w *serverWebviewWindow) snapAssist()                                {}
func (w *serverWebviewWindow) setContentProtection(enabled bool)          {}
