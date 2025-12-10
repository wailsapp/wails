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

//go:embed docker/Dockerfile.cross
var dockerfileContent string

// DependencyStatus represents the status of a dependency
type DependencyStatus struct {
	Name           string `json:"name"`
	Installed      bool   `json:"installed"`
	Version        string `json:"version,omitempty"`
	Status         string `json:"status"` // "installed", "not_installed", "needs_update"
	Required       bool   `json:"required"`
	Message        string `json:"message,omitempty"`
	InstallCommand string `json:"installCommand,omitempty"`
	HelpURL        string `json:"helpUrl,omitempty"`
	ImageBuilt     bool   `json:"imageBuilt"` // For Docker: whether wails-cross image exists
}

// DockerStatus represents Docker installation and image status
type DockerStatus struct {
	Installed    bool   `json:"installed"`
	Running      bool   `json:"running"`
	Version      string `json:"version,omitempty"`
	ImageBuilt   bool   `json:"imageBuilt"`
	ImageName    string `json:"imageName"`
	ImageSize    string `json:"imageSize,omitempty"`
	PullProgress int    `json:"pullProgress"`
	PullStatus   string `json:"pullStatus"` // "idle", "pulling", "complete", "error"
	PullError    string `json:"pullError,omitempty"`
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
	server       *http.Server
	state        WizardState
	stateMu      sync.RWMutex
	dockerStatus DockerStatus
	dockerMu     sync.RWMutex
	done         chan struct{}
	shutdown     chan struct{}
	buildWg      sync.WaitGroup // Tracks background Docker builds
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
	mux.HandleFunc("/api/dependencies/install", w.handleInstallDependency)
	mux.HandleFunc("/api/docker/status", w.handleDockerStatus)
	mux.HandleFunc("/api/docker/build", w.handleDockerBuild)
	mux.HandleFunc("/api/docker/start-background", w.handleDockerStartBackground)
	mux.HandleFunc("/api/wails-config", w.handleWailsConfig)
	mux.HandleFunc("/api/defaults", w.handleDefaults)
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

	// Check if Docker build is in progress
	w.dockerMu.RLock()
	dockerBuilding := w.dockerStatus.PullStatus == "pulling"
	w.dockerMu.RUnlock()

	response := map[string]interface{}{
		"status":         "closing",
		"dockerBuilding": dockerBuilding,
	}
	if dockerBuilding {
		response["message"] = "Docker image build will continue in the background"
	}
	json.NewEncoder(rw).Encode(response)

	// Wait for any running Docker builds to complete before shutting down
	go func() {
		w.buildWg.Wait()
		close(w.shutdown)
	}()
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

func (w *Wizard) handleDockerStatus(rw http.ResponseWriter, r *http.Request) {
	// Get fresh Docker info (installed, running, image status)
	freshStatus := w.checkDocker()

	w.dockerMu.Lock()
	// Preserve the build progress/status if currently building
	// but update the image built status and other base info
	if w.dockerStatus.PullStatus == "pulling" {
		// Keep the pulling status and progress, but update image built flag
		freshStatus.PullStatus = w.dockerStatus.PullStatus
		freshStatus.PullProgress = w.dockerStatus.PullProgress
		freshStatus.PullError = w.dockerStatus.PullError
	} else if w.dockerStatus.PullStatus == "complete" || w.dockerStatus.PullStatus == "error" {
		// Keep the final status until it's acknowledged
		freshStatus.PullStatus = w.dockerStatus.PullStatus
		freshStatus.PullProgress = w.dockerStatus.PullProgress
		freshStatus.PullError = w.dockerStatus.PullError
	}
	w.dockerStatus = freshStatus
	status := w.dockerStatus
	w.dockerMu.Unlock()

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(status)
}

func (w *Wizard) checkDocker() DockerStatus {
	status := DockerStatus{
		ImageName:  "wails-cross",
		PullStatus: "idle",
	}

	// Check if Docker is installed
	output, err := execCommand("docker", "--version")
	if err != nil {
		status.Installed = false
		return status
	}

	status.Installed = true
	// Parse version from "Docker version 24.0.7, build afdd53b"
	parts := strings.Split(output, ",")
	if len(parts) > 0 {
		status.Version = strings.TrimPrefix(strings.TrimSpace(parts[0]), "Docker version ")
	}

	// Check if Docker daemon is running
	if _, err := execCommand("docker", "info"); err != nil {
		status.Running = false
		return status
	}
	status.Running = true

	// Check if wails-cross image exists
	imageOutput, err := execCommand("docker", "image", "inspect", "wails-cross")
	status.ImageBuilt = err == nil && len(imageOutput) > 0

	// Get image size if it exists
	if status.ImageBuilt {
		sizeOutput, err := execCommand("docker", "images", "wails-cross", "--format", "{{.Size}}")
		if err == nil && len(sizeOutput) > 0 {
			status.ImageSize = strings.TrimSpace(sizeOutput)
		}
	}

	return status
}

func (w *Wizard) handleDockerBuild(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.dockerMu.Lock()
	w.dockerStatus.PullStatus = "pulling"
	w.dockerStatus.PullProgress = 0
	w.dockerMu.Unlock()

	// Build the Docker image in background
	w.buildWg.Add(1)
	go func() {
		defer w.buildWg.Done()

		// Write Dockerfile to temp directory
		tmpDir, err := os.MkdirTemp("", "wails-docker-build-*")
		if err != nil {
			w.dockerMu.Lock()
			w.dockerStatus.PullStatus = "error"
			w.dockerStatus.PullError = fmt.Sprintf("Failed to create temp dir: %v", err)
			w.dockerMu.Unlock()
			return
		}
		defer os.RemoveAll(tmpDir)

		dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
		if err := os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644); err != nil {
			w.dockerMu.Lock()
			w.dockerStatus.PullStatus = "error"
			w.dockerStatus.PullError = fmt.Sprintf("Failed to write Dockerfile: %v", err)
			w.dockerMu.Unlock()
			return
		}

		// Build the Docker image: docker build -t wails-cross -f Dockerfile .
		cmd := exec.Command("docker", "build", "-t", "wails-cross", "-f", dockerfilePath, tmpDir)
		output, err := cmd.CombinedOutput()

		w.dockerMu.Lock()
		if err != nil {
			w.dockerStatus.PullStatus = "error"
			w.dockerStatus.PullError = fmt.Sprintf("Build failed: %v\n%s", err, string(output))
		} else {
			w.dockerStatus.PullStatus = "complete"
			w.dockerStatus.ImageBuilt = true
			// Get image size
			if sizeOutput, sizeErr := execCommand("docker", "images", "wails-cross", "--format", "{{.Size}}"); sizeErr == nil && len(sizeOutput) > 0 {
				w.dockerStatus.ImageSize = strings.TrimSpace(sizeOutput)
			}
		}
		w.dockerStatus.PullProgress = 100
		w.dockerMu.Unlock()
	}()

	// Simulate progress updates while building (only increment, never decrease)
	go func() {
		for i := 5; i < 90; i += 5 {
			time.Sleep(2 * time.Second)
			w.dockerMu.Lock()
			if w.dockerStatus.PullStatus != "pulling" {
				w.dockerMu.Unlock()
				return
			}
			// Only update if new value is higher (prevents going backwards)
			if i > w.dockerStatus.PullProgress {
				w.dockerStatus.PullProgress = i
			}
			w.dockerMu.Unlock()
		}
	}()

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{"status": "started"})
}

// handleDockerStartBackground checks if Docker is available and starts building in background
// This is called early in the wizard flow to get a head start on the image build
func (w *Wizard) handleDockerStartBackground(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	// Check Docker status first
	status := w.checkDocker()

	w.dockerMu.Lock()
	w.dockerStatus = status
	w.dockerMu.Unlock()

	// Only start build if Docker is installed, running, and image not built yet
	if !status.Installed || !status.Running || status.ImageBuilt {
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"started": false,
			"reason":  getDockerNotStartedReason(status),
			"status":  status,
		})
		return
	}

	// Check if already building
	w.dockerMu.RLock()
	alreadyBuilding := w.dockerStatus.PullStatus == "pulling"
	w.dockerMu.RUnlock()

	if alreadyBuilding {
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"started": false,
			"reason":  "already_building",
			"status":  status,
		})
		return
	}

	// Start building in background
	w.dockerMu.Lock()
	w.dockerStatus.PullStatus = "pulling"
	w.dockerStatus.PullProgress = 0
	w.dockerMu.Unlock()

	// Build the Docker image in background using embedded Dockerfile
	w.buildWg.Add(1)
	go func() {
		defer w.buildWg.Done()

		// Write Dockerfile to temp directory
		tmpDir, err := os.MkdirTemp("", "wails-docker-build-*")
		if err != nil {
			w.dockerMu.Lock()
			w.dockerStatus.PullStatus = "error"
			w.dockerStatus.PullError = fmt.Sprintf("Failed to create temp dir: %v", err)
			w.dockerMu.Unlock()
			return
		}
		defer os.RemoveAll(tmpDir)

		dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
		if err := os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644); err != nil {
			w.dockerMu.Lock()
			w.dockerStatus.PullStatus = "error"
			w.dockerStatus.PullError = fmt.Sprintf("Failed to write Dockerfile: %v", err)
			w.dockerMu.Unlock()
			return
		}

		// Build the Docker image: docker build -t wails-cross -f Dockerfile .
		cmd := exec.Command("docker", "build", "-t", "wails-cross", "-f", dockerfilePath, tmpDir)
		output, err := cmd.CombinedOutput()

		w.dockerMu.Lock()
		if err != nil {
			w.dockerStatus.PullStatus = "error"
			w.dockerStatus.PullError = fmt.Sprintf("Build failed: %v\n%s", err, string(output))
		} else {
			w.dockerStatus.PullStatus = "complete"
			w.dockerStatus.ImageBuilt = true
			// Get image size
			if sizeOutput, sizeErr := execCommand("docker", "images", "wails-cross", "--format", "{{.Size}}"); sizeErr == nil && len(sizeOutput) > 0 {
				w.dockerStatus.ImageSize = strings.TrimSpace(sizeOutput)
			}
		}
		w.dockerStatus.PullProgress = 100
		w.dockerMu.Unlock()
	}()

	// Simulate progress updates while building (only increment, never decrease)
	go func() {
		for i := 5; i < 90; i += 5 {
			time.Sleep(2 * time.Second)
			w.dockerMu.Lock()
			if w.dockerStatus.PullStatus != "pulling" {
				w.dockerMu.Unlock()
				return
			}
			// Only update if new value is higher (prevents going backwards)
			if i > w.dockerStatus.PullProgress {
				w.dockerStatus.PullProgress = i
			}
			w.dockerMu.Unlock()
		}
	}()

	json.NewEncoder(rw).Encode(map[string]interface{}{
		"started": true,
		"status":  status,
	})
}

func getDockerNotStartedReason(status DockerStatus) string {
	if !status.Installed {
		return "not_installed"
	}
	if !status.Running {
		return "not_running"
	}
	if status.ImageBuilt {
		return "already_built"
	}
	return "unknown"
}

// InstallRequest represents a request to install a dependency
type InstallRequest struct {
	Command string `json:"command"`
}

// InstallResponse represents the result of an install attempt
type InstallResponse struct {
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

func (w *Wizard) handleInstallDependency(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req InstallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	// Execute the install command
	// Split the command into parts
	parts := strings.Fields(req.Command)
	if len(parts) == 0 {
		json.NewEncoder(rw).Encode(InstallResponse{
			Success: false,
			Error:   "Empty command",
		})
		return
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		json.NewEncoder(rw).Encode(InstallResponse{
			Success: false,
			Output:  string(output),
			Error:   err.Error(),
		})
		return
	}

	json.NewEncoder(rw).Encode(InstallResponse{
		Success: true,
		Output:  string(output),
	})
}

func (w *Wizard) handleDefaults(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		defaults, err := LoadGlobalDefaults()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// Try to pre-populate author info from git config if empty
		if defaults.Author.Name == "" {
			if name, err := execCommand("git", "config", "--global", "user.name"); err == nil && name != "" {
				defaults.Author.Name = name
			}
		}

		json.NewEncoder(rw).Encode(defaults)

	case http.MethodPost:
		var defaults GlobalDefaults
		if err := json.NewDecoder(r.Body).Decode(&defaults); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		if err := SaveGlobalDefaults(defaults); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		path, _ := GetDefaultsPath()
		json.NewEncoder(rw).Encode(map[string]string{"status": "saved", "path": path})

	default:
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
