package webserver

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/wailsapp/wails/v2/internal/assetdb"
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/messagedispatcher"
	"github.com/wailsapp/wails/v2/internal/subsystem"
)

// WebServer serves the application over http
type WebServer struct {
	port       int
	ip         string
	assets     *assetdb.AssetDB
	logger     *logger.Logger
	dispatcher *messagedispatcher.Dispatcher
	event      *subsystem.Event
	server     *http.Server
	bindings   *binding.Bindings

	lock        sync.Mutex
	connections map[string]*WebClient
}

// NewWebServer creates a new WebServer
func NewWebServer(logger *logger.Logger) *WebServer {

	// Return a WebServer with default values
	return &WebServer{
		assets:      db,
		connections: make(map[string]*WebClient),
		logger:      logger,
	}
}

// URL returns the URL that the server is serving from
func (w *WebServer) URL() string {
	return fmt.Sprintf("http://%s:%d", w.ip, w.port)
}

// SetPort sets the server port to listen on
func (w *WebServer) SetPort(port int) {
	w.port = port
}

// SetIP sets the server ip to listen on
func (w *WebServer) SetIP(ip string) {
	w.ip = ip
}

// SetBindings provides the webserver with the mapping of bindings to provide to a client
func (w *WebServer) SetBindings(bindings *binding.Bindings) {
	w.bindings = bindings
}

// Start the webserver
func (w *WebServer) Start(dispatcher *messagedispatcher.Dispatcher, event *subsystem.Event) error {
	var err error

	// Create the server
	w.server = &http.Server{Addr: fmt.Sprintf("%s:%d", w.ip, w.port)}
	w.event = event

	// Set up the Web Server's routes
	err = w.setUpRoutes()
	if err != nil {
		return err
	}

	// Save the dispatcher
	w.dispatcher = dispatcher

	// Start the WebServer
	err = w.server.ListenAndServe()

	// Return any error except http.ErrServerClosed
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
