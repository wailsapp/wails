package setupwizard

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/pkg/browser"
	"github.com/wailsapp/wails/v3/internal/operatingsystem"
	"github.com/wailsapp/wails/v3/internal/version"
	"gopkg.in/yaml.v3"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

// DependencyStatus represents the status of a dependency
type DependencyStatus struct {
	Name      string `json:"name"`
	Installed bool   `json:"installed"`
	Version   string `json:"version,omitempty"`
	Status    string `json:"status"` // "installed", "not_installed", "needs_update"
	Required  bool   `json:"required"`
	Message   string `json:"message,omitempty"`
}

// WailsConfigInfo represents the info section of wails.yaml
type WailsConfigInfo struct {
	CompanyName       string `json:"companyName" yaml:"companyName"`
	ProductName       string `json:"productName" yaml:"productName"`
	ProductIdentifier string `json:"productIdentifier" yaml:"productIdentifier"`
	Description       string `json:"description" yaml:"description"`
	Copyright         string `json:"copyright" yaml:"copyright"`
	Comments          string `json:"comments,omitempty" yaml:"comments,omitempty"`
	Version           string `json:"version" yaml:"version"`
}

// WailsConfig represents the wails.yaml configuration
type WailsConfig struct {
	Info WailsConfigInfo `json:"info" yaml:"info"`
}

// SystemInfo contains detected system information
type SystemInfo struct {
	OS           string `json:"os"`
	Arch         string `json:"arch"`
	WailsVersion string `json:"wailsVersion"`
	GoVersion    string `json:"goVersion"`
	HomeDir      string `json:"homeDir"`
	OSName       string `json:"osName,omitempty"`
	OSVersion    string `json:"osVersion,omitempty"`
}

// WizardState represents the complete wizard state
type WizardState struct {
	Dependencies []DependencyStatus `json:"dependencies"`
	System       SystemInfo         `json:"system"`
	StartTime    time.Time          `json:"startTime"`
}

// Wizard is the setup wizard server
type Wizard struct {
	server   *http.Server
	state    WizardState
	stateMu  sync.RWMutex
	done     chan struct{}
	shutdown chan struct{}
}

// New creates a new setup wizard
func New() *Wizard {
	return &Wizard{
		done:     make(chan struct{}),
		shutdown: make(chan struct{}),
		state: WizardState{
			StartTime: time.Now(),
		},
	}
}

// Run starts the wizard and opens it in the browser
func (w *Wizard) Run() error {
	// Initialize system info
	w.initSystemInfo()

	// Find an available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("failed to find available port: %w", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://127.0.0.1:%d", port)

	// Set up HTTP routes
	mux := http.NewServeMux()
	w.setupRoutes(mux)

	w.server = &http.Server{
		Handler: mux,
	}

	// Start server in goroutine
	go func() {
		if err := w.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		}
	}()

	fmt.Printf("Setup wizard running at %s\n", url)

	// Open browser
	if err := browser.OpenURL(url); err != nil {
		fmt.Printf("Please open %s in your browser\n", url)
	}

	// Wait for completion or shutdown
	select {
	case <-w.done:
		fmt.Println("\nSetup completed successfully!")
	case <-w.shutdown:
		fmt.Println("\nSetup wizard closed.")
	}

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return w.server.Shutdown(ctx)
}

func (w *Wizard) setupRoutes(mux *http.ServeMux) {
	// API routes
	mux.HandleFunc("/api/state", w.handleState)
	mux.HandleFunc("/api/dependencies/check", w.handleCheckDependencies)
	mux.HandleFunc("/api/wails-config", w.handleWailsConfig)
	mux.HandleFunc("/api/complete", w.handleComplete)
	mux.HandleFunc("/api/close", w.handleClose)

	// Serve frontend
	frontendDist, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(frontendDist))

	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		// Try to serve the file
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		// Check if file exists
		if _, err := fs.Stat(frontendDist, strings.TrimPrefix(path, "/")); err != nil {
			// Serve index.html for SPA routing
			r.URL.Path = "/"
		}

		fileServer.ServeHTTP(rw, r)
	})
}

func (w *Wizard) initSystemInfo() {
	w.stateMu.Lock()
	defer w.stateMu.Unlock()

	homeDir, _ := os.UserHomeDir()

	w.state.System = SystemInfo{
		OS:           runtime.GOOS,
		Arch:         runtime.GOARCH,
		WailsVersion: version.String(),
		GoVersion:    runtime.Version(),
		HomeDir:      homeDir,
	}

	// Get OS details
	if info, err := operatingsystem.Info(); err == nil {
		w.state.System.OSName = info.Name
		w.state.System.OSVersion = info.Version
	}
}

func (w *Wizard) handleState(rw http.ResponseWriter, r *http.Request) {
	w.stateMu.RLock()
	defer w.stateMu.RUnlock()

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(w.state)
}

func (w *Wizard) handleCheckDependencies(rw http.ResponseWriter, r *http.Request) {
	deps := w.checkAllDependencies()

	w.stateMu.Lock()
	w.state.Dependencies = deps
	w.stateMu.Unlock()

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(deps)
}

func (w *Wizard) handleWailsConfig(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	// Find wails.yaml in current directory or parent directories
	configPath := findWailsConfig()

	switch r.Method {
	case http.MethodGet:
		if configPath == "" {
			json.NewEncoder(rw).Encode(nil)
			return
		}

		data, err := os.ReadFile(configPath)
		if err != nil {
			json.NewEncoder(rw).Encode(nil)
			return
		}

		var config WailsConfig
		if err := yaml.Unmarshal(data, &config); err != nil {
			json.NewEncoder(rw).Encode(nil)
			return
		}

		json.NewEncoder(rw).Encode(config)

	case http.MethodPost:
		var config WailsConfig
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		if configPath == "" {
			configPath = "wails.yaml"
		}

		data, err := yaml.Marshal(&config)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := os.WriteFile(configPath, data, 0644); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(rw).Encode(map[string]string{"status": "saved", "path": configPath})

	default:
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func findWailsConfig() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		configPath := filepath.Join(dir, "wails.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}

func (w *Wizard) handleComplete(rw http.ResponseWriter, r *http.Request) {
	w.stateMu.RLock()
	state := w.state
	w.stateMu.RUnlock()

	duration := time.Since(state.StartTime)

	response := map[string]interface{}{
		"status":   "complete",
		"duration": duration.String(),
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(response)

	close(w.done)
}

func (w *Wizard) handleClose(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{"status": "closing"})

	close(w.shutdown)
}

// execCommand runs a command and returns its output
func execCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	return strings.TrimSpace(string(output)), err
}

// commandExists checks if a command exists in PATH
func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
