package setupwizard

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wailsapp/wails/v3/internal/browser"
	"github.com/wailsapp/wails/v3/internal/operatingsystem"
	"github.com/wailsapp/wails/v3/internal/version"
	"gopkg.in/yaml.v3"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

//go:embed assets/apple-sdk-license.pdf
var appleLicensePDF []byte

// DependencyStatus represents the status of a dependency
type DependencyStatus struct {
	Name           string `json:"name"`
	Installed      bool   `json:"installed"`
	Version        string `json:"version,omitempty"`
	Status         string `json:"status"` // "installed", "not_installed", "needs_update"
	Required       bool   `json:"required"`
	Message        string `json:"message,omitempty"`
	InstallCommand string `json:"installCommand,omitempty"`
	ConfigCommand  string `json:"configCommand,omitempty"` // shell lines to add (e.g. export JAVA_HOME=...), shown with a Copy button
	HelpURL        string `json:"helpUrl,omitempty"`
	HelpLabel      string `json:"helpLabel,omitempty"` // OS-specific link text, e.g. "Get Xcode from the App Store"
	ImageBuilt     bool   `json:"imageBuilt"`          // For Docker: whether wails-cross image exists
}

// DockerStatus represents Docker installation and image status
type PullProgress struct {
	Stage    string
	Progress int
}

type pullParser struct {
	layerSizes      map[string]float64
	layerDownloaded map[string]float64
	layerComplete   map[string]bool
	layersPending   map[string]bool
	stage           string
}

func newPullParser() *pullParser {
	return &pullParser{
		layerSizes:      make(map[string]float64),
		layerDownloaded: make(map[string]float64),
		layerComplete:   make(map[string]bool),
		layersPending:   make(map[string]bool),
		stage:           "Connecting",
	}
}

func parseSize(s string) float64 {
	s = strings.TrimSpace(s)
	var val float64
	var unit string
	fmt.Sscanf(s, "%f%s", &val, &unit)
	switch strings.ToUpper(unit) {
	case "KB", "KIB":
		return val * 1024
	case "MB", "MIB":
		return val * 1024 * 1024
	case "GB", "GIB":
		return val * 1024 * 1024 * 1024
	case "B":
		return val
	}
	return val
}

var sizeRegex = regexp.MustCompile(`(\d+(?:\.\d+)?[KMGT]?i?B)/(\d+(?:\.\d+)?[KMGT]?i?B)`)

func (p *pullParser) ParseLine(line string) PullProgress {
	if len(line) == 0 {
		return PullProgress{Stage: p.stage, Progress: p.calculateProgress()}
	}

	parts := strings.Fields(line)
	if len(parts) == 0 {
		return PullProgress{Stage: p.stage, Progress: p.calculateProgress()}
	}

	layerID := strings.TrimSuffix(parts[0], ":")

	switch {
	case strings.Contains(line, "Pulling from"):
		p.stage = "Connecting"
	case strings.Contains(line, "Pulling fs layer"):
		p.stage = "Downloading"
		p.layersPending[layerID] = true
	case strings.Contains(line, "Waiting"):
		p.stage = "Downloading"
		p.layersPending[layerID] = true
	case strings.Contains(line, "Downloading"):
		p.stage = "Downloading"
		p.layersPending[layerID] = true
		if matches := sizeRegex.FindStringSubmatch(line); len(matches) == 3 {
			p.layerDownloaded[layerID] = parseSize(matches[1])
			p.layerSizes[layerID] = parseSize(matches[2])
		}
	case strings.Contains(line, "Verifying"):
		p.stage = "Verifying"
	case strings.Contains(line, "Download complete"):
		p.stage = "Extracting"
		if size, ok := p.layerSizes[layerID]; ok {
			p.layerDownloaded[layerID] = size
		}
	case strings.Contains(line, "Extracting"):
		p.stage = "Extracting"
		if matches := sizeRegex.FindStringSubmatch(line); len(matches) == 3 {
			p.layerDownloaded[layerID] = parseSize(matches[1])
			p.layerSizes[layerID] = parseSize(matches[2])
		}
	case strings.Contains(line, "Pull complete"):
		p.stage = "Extracting"
		p.layerComplete[layerID] = true
		delete(p.layersPending, layerID)
		if size, ok := p.layerSizes[layerID]; ok {
			p.layerDownloaded[layerID] = size
		}
	}

	return PullProgress{Stage: p.stage, Progress: p.calculateProgress()}
}

func (p *pullParser) calculateProgress() int {
	totalLayers := len(p.layersPending) + len(p.layerComplete)
	if totalLayers == 0 {
		return 0
	}

	var totalSize, downloaded float64
	hasSizeInfo := len(p.layerSizes) > 0

	if hasSizeInfo {
		for id, size := range p.layerSizes {
			totalSize += size
			if p.layerComplete[id] {
				downloaded += size
			} else if dl, ok := p.layerDownloaded[id]; ok {
				downloaded += dl
			}
		}
		if totalSize > 0 {
			progress := int((downloaded / totalSize) * 100)
			if progress > 95 && len(p.layerComplete) < len(p.layerSizes) {
				progress = 95
			}
			return progress
		}
	}

	progress := (len(p.layerComplete) * 100) / totalLayers
	if progress > 95 && len(p.layersPending) > 0 {
		progress = 95
	}
	return progress
}

type DockerStatus struct {
	Installed     bool   `json:"installed"`
	Running       bool   `json:"running"`
	Version       string `json:"version,omitempty"`
	ImageBuilt    bool   `json:"imageBuilt"`
	ImageName     string `json:"imageName"`
	ImageSize     string `json:"imageSize,omitempty"`
	ImageVersion  string `json:"imageVersion,omitempty"`
	SDKVersion    string `json:"sdkVersion,omitempty"`
	UpdateAvail   bool   `json:"updateAvailable"`
	LatestVersion string `json:"latestVersion,omitempty"`
	PullProgress  int    `json:"pullProgress"`
	PullMessage   string `json:"pullMessage,omitempty"`
	PullStatus    string `json:"pullStatus"`
	PullError     string `json:"pullError,omitempty"`
	BytesTotal    int64  `json:"bytesTotal,omitempty"`
	BytesDone     int64  `json:"bytesDone,omitempty"`
	LayerCount    int    `json:"layerCount,omitempty"`
	LayersDone    int    `json:"layersDone,omitempty"`
}

const crossImageName = "ghcr.io/wailsapp/wails-cross"

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
	GitName      string `json:"gitName,omitempty"`
	GitEmail     string `json:"gitEmail,omitempty"`
}

// WizardState represents the complete wizard state
type WizardState struct {
	Dependencies []DependencyStatus `json:"dependencies"`
	System       SystemInfo         `json:"system"`
	StartTime    time.Time          `json:"startTime"`
}

// Wizard is the setup wizard server
type Wizard struct {
	server          *http.Server
	state           WizardState
	stateMu         sync.RWMutex
	dockerStatus    DockerStatus
	dockerBuildLogs string
	dockerMu        sync.RWMutex
	done            chan struct{}
	doneOnce        sync.Once
	shutdown        chan struct{}
	shutdownOnce    sync.Once
	buildWg         sync.WaitGroup

	// Init-mode state. When initData is non-nil the wizard runs as the project
	// "init" wizard (wails3 init -ui) instead of the global setup wizard.
	initData   *InitData
	initResult *InitData
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
	mux.HandleFunc("/api/dependencies/mobile", w.handleCheckMobileDependencies)
	mux.HandleFunc("/api/dependencies/install", w.handleInstallDependency)
	mux.HandleFunc("/api/docker/status", w.handleDockerStatus)
	mux.HandleFunc("/api/docker/status/stream", w.handleDockerStatusStream)
	mux.HandleFunc("/api/docker/build", w.handleDockerBuild)
	mux.HandleFunc("/api/docker/logs", w.handleDockerLogs)
	mux.HandleFunc("/api/docker/start-background", w.handleDockerStartBackground)
	mux.HandleFunc("/api/wails-config", w.handleWailsConfig)
	mux.HandleFunc("/api/defaults", w.handleDefaults)
	mux.HandleFunc("/api/signing", w.handleSigning)
	mux.HandleFunc("/api/signing/status", w.handleSigningStatus)
	mux.HandleFunc("/api/signing/notarize/create", w.handleNotarizeCreate)
	mux.HandleFunc("/api/signing/notarize/validate", w.handleNotarizeValidate)
	mux.HandleFunc("/api/init", w.handleInit)
	mux.HandleFunc("/api/init/create", w.handleInitCreate)
	mux.HandleFunc("/api/signing/gpg/create", w.handleGPGCreate)
	mux.HandleFunc("/api/signing/gpg/export", w.handleGPGExport)
	mux.HandleFunc("/api/signing/windows/create-cert", w.handleWindowsCertCreate)
	mux.HandleFunc("/api/complete", w.handleComplete)
	mux.HandleFunc("/api/close", w.handleClose)
	mux.HandleFunc("/api/report-bug", w.handleReportBug)

	mux.HandleFunc("/assets/apple-sdk-license.pdf", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/pdf")
		rw.Write(appleLicensePDF)
	})

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
			path = "/index.html"
		}

		// Set the content type explicitly for static asset types. http.FileServer
		// otherwise relies on mime.TypeByExtension, which returns "" on hosts with
		// an incomplete MIME database (e.g. minimal containers, some Windows setups).
		// For .svg that fallback sniffs to "text/xml", which browsers refuse to
		// render inside an <img>, leaving the Wails logo blank. Pin the common types.
		if ct := staticContentType(path); ct != "" {
			rw.Header().Set("Content-Type", ct)
		}

		fileServer.ServeHTTP(rw, r)
	})
}

// staticContentType returns a stable Content-Type for the static asset
// extensions the setup wizard frontend ships, so rendering never depends on the
// host's MIME database. Returns "" for unknown extensions, leaving http.FileServer
// to handle them as before.
func staticContentType(path string) string {
	switch {
	case strings.HasSuffix(path, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(path, ".html"):
		return "text/html; charset=utf-8"
	case strings.HasSuffix(path, ".css"):
		return "text/css; charset=utf-8"
	case strings.HasSuffix(path, ".js"), strings.HasSuffix(path, ".mjs"):
		return "text/javascript; charset=utf-8"
	case strings.HasSuffix(path, ".json"):
		return "application/json"
	case strings.HasSuffix(path, ".png"):
		return "image/png"
	case strings.HasSuffix(path, ".webp"):
		return "image/webp"
	case strings.HasSuffix(path, ".ico"):
		return "image/x-icon"
	case strings.HasSuffix(path, ".woff2"):
		return "font/woff2"
	case strings.HasSuffix(path, ".woff"):
		return "font/woff"
	}
	return ""
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

	// Pre-fill author identity from git config (effective config: local repo
	// overrides global). Used to seed the GPG key and other signing forms.
	if name, err := execCommand("git", "config", "user.name"); err == nil {
		w.state.System.GitName = name
	}
	if email, err := execCommand("git", "config", "user.email"); err == nil {
		w.state.System.GitEmail = email
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

// handleCheckMobileDependencies reports the toolchain status for the requested
// mobile platforms, e.g. /api/dependencies/mobile?ios=true&android=true.
func (w *Wizard) handleCheckMobileDependencies(rw http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	deps := w.checkMobileDependencies(q.Get("ios") == "true", q.Get("android") == "true")
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

	w.doneOnce.Do(func() { close(w.done) })
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
		w.shutdownOnce.Do(func() { close(w.shutdown) })
	}()
}

func (w *Wizard) handleReportBug(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	// Get current step from query parameter
	currentStep := r.URL.Query().Get("step")
	if currentStep == "" {
		currentStep = "unknown"
	}

	// Gather system info
	w.stateMu.RLock()
	system := w.state.System
	w.stateMu.RUnlock()

	// Build a concise comment body - description first, then details table
	var sb strings.Builder
	sb.WriteString("**What went wrong?**\n\n\n\n")
	sb.WriteString("**What were you doing when the issue occurred?**\n\n\n\n")
	sb.WriteString("---\n\n")
	sb.WriteString("| | |\n")
	sb.WriteString("|--|--|\n")
	sb.WriteString(fmt.Sprintf("| Platform | %s |\n", system.OS))
	sb.WriteString(fmt.Sprintf("| Arch | %s |\n", system.Arch))
	sb.WriteString(fmt.Sprintf("| Wails | %s |\n", system.WailsVersion))
	sb.WriteString(fmt.Sprintf("| Go | %s |\n", system.GoVersion))
	sb.WriteString(fmt.Sprintf("| Step | %s |\n", currentStep))

	issueURL := "https://github.com/wailsapp/wails/issues/4904#issue-comment-box"
	commentBody := sb.String()

	// Return the body for the frontend to copy to clipboard
	// Frontend will handle opening the browser after showing the overlay
	json.NewEncoder(rw).Encode(map[string]interface{}{
		"status": "ready",
		"url":    issueURL,
		"body":   commentBody,
	})
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
	if w.dockerStatus.PullStatus == "pulling" || w.dockerStatus.PullStatus == "complete" || w.dockerStatus.PullStatus == "error" {
		freshStatus.PullStatus = w.dockerStatus.PullStatus
		freshStatus.PullProgress = w.dockerStatus.PullProgress
		freshStatus.PullMessage = w.dockerStatus.PullMessage
		freshStatus.PullError = w.dockerStatus.PullError
		freshStatus.BytesTotal = w.dockerStatus.BytesTotal
		freshStatus.BytesDone = w.dockerStatus.BytesDone
		freshStatus.LayerCount = w.dockerStatus.LayerCount
		freshStatus.LayersDone = w.dockerStatus.LayersDone
	}
	w.dockerStatus = freshStatus
	status := w.dockerStatus
	w.dockerMu.Unlock()

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(status)
}

func (w *Wizard) handleDockerStatusStream(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")

	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "SSE not supported", http.StatusInternalServerError)
		return
	}

	sendStatus := func() (done bool) {
		w.dockerMu.RLock()
		status := w.dockerStatus
		w.dockerMu.RUnlock()

		data, _ := json.Marshal(status)
		fmt.Fprintf(rw, "data: %s\n\n", data)
		flusher.Flush()

		return status.PullStatus == "complete" || status.PullStatus == "error"
	}

	if sendStatus() {
		return
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			if sendStatus() {
				return
			}
		}
	}
}

func (w *Wizard) checkDocker() DockerStatus {
	status := DockerStatus{
		ImageName:  crossImageName,
		PullStatus: "idle",
	}

	output, err := execCommand("docker", "--version")
	if err != nil {
		status.Installed = false
		return status
	}

	status.Installed = true
	parts := strings.Split(output, ",")
	if len(parts) > 0 {
		status.Version = strings.TrimPrefix(strings.TrimSpace(parts[0]), "Docker version ")
	}

	if _, err := execCommand("docker", "info"); err != nil {
		status.Running = false
		return status
	}
	status.Running = true

	imageOutput, err := execCommand("docker", "image", "inspect", crossImageName)
	status.ImageBuilt = err == nil && len(imageOutput) > 0

	if status.ImageBuilt {
		sizeOutput, err := execCommand("docker", "images", crossImageName, "--format", "{{.Size}}")
		if err == nil && len(sizeOutput) > 0 {
			status.ImageSize = strings.TrimSpace(sizeOutput)
		}
		versionOutput, err := execCommand("docker", "inspect", crossImageName, "--format", "{{index .Config.Labels \"org.opencontainers.image.version\"}}")
		if err == nil && len(versionOutput) > 0 {
			status.ImageVersion = strings.TrimSpace(versionOutput)
		}
		sdkOutput, err := execCommand("docker", "inspect", crossImageName, "--format", "{{index .Config.Labels \"io.wails.sdk.version\"}}")
		if err == nil && len(sdkOutput) > 0 {
			status.SDKVersion = strings.TrimSpace(sdkOutput)
		}
	}

	return status
}

type dockerPullEvent struct {
	Status         string `json:"status"`
	ID             string `json:"id"`
	Progress       string `json:"progress"`
	ProgressDetail struct {
		Current int64 `json:"current"`
		Total   int64 `json:"total"`
	} `json:"progressDetail"`
	Error string `json:"error"`
}

type layerProgress struct {
	dlTotal   int64
	dlCurrent int64
	dlDone    bool
	exTotal   int64
	exCurrent int64
	exDone    bool
}

func (w *Wizard) startDockerPull() {
	w.dockerMu.Lock()
	if w.dockerStatus.PullStatus == "pulling" {
		w.dockerMu.Unlock()
		return
	}
	w.dockerStatus.PullStatus = "pulling"
	w.dockerStatus.PullProgress = 0
	w.dockerStatus.PullMessage = "Connecting"
	// Reset stale state from previous attempts
	w.dockerStatus.PullError = ""
	w.dockerStatus.BytesTotal = 0
	w.dockerStatus.BytesDone = 0
	w.dockerStatus.LayerCount = 0
	w.dockerStatus.LayersDone = 0
	w.dockerBuildLogs = ""
	w.dockerMu.Unlock()

	w.buildWg.Add(1)
	go func() {
		defer w.buildWg.Done()

		if err := w.pullViaDockerAPI(); err != nil {
			// Reset status so the SSE stream stays open and handleClose detects
			// the active build while the CLI fallback runs.
			w.dockerMu.Lock()
			w.dockerStatus.PullStatus = "pulling"
			w.dockerStatus.PullMessage = "Retrying via Docker CLI"
			w.dockerStatus.PullError = ""
			w.dockerMu.Unlock()
			w.pullViaDockerCLI()
		}
	}()
}

// pullViaDockerAPI attempts to pull the image using Docker's HTTP API directly.
// This provides detailed progress tracking with layer-by-layer download status.
// Uses API v1.44 (Docker 25.0+). If this fails for any reason (older Docker version,
// permission issues, etc.), the caller falls back to pullViaDockerCLI which works
// with any Docker version but provides less detailed progress.
func (w *Wizard) pullViaDockerAPI() error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("windows named pipes not supported, falling back to CLI")
	}
	socketPath := "/var/run/docker.sock"

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
		Timeout: 30 * time.Minute,
	}

	imageParts := strings.SplitN(crossImageName, ":", 2)
	imageName := imageParts[0]
	tag := "latest"
	if len(imageParts) > 1 {
		tag = imageParts[1]
	}

	url := fmt.Sprintf("http://localhost/v1.44/images/create?fromImage=%s&tag=%s", imageName, tag)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to Docker API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Docker API returned status %d", resp.StatusCode)
	}

	layers := make(map[string]*layerProgress)
	var logs strings.Builder
	var maxTotal int64
	decoder := json.NewDecoder(resp.Body)

	for {
		var event dockerPullEvent
		if err := decoder.Decode(&event); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("docker API stream error: %w", err)
		}

		logs.WriteString(fmt.Sprintf("%s %s %s\n", event.ID, event.Status, event.Progress))

		if event.Error != "" {
			w.dockerMu.Lock()
			w.dockerStatus.PullStatus = "error"
			w.dockerStatus.PullError = event.Error
			w.dockerStatus.PullMessage = "Failed"
			w.dockerBuildLogs = logs.String()
			w.dockerMu.Unlock()
			return fmt.Errorf("docker pull error: %s", event.Error)
		}

		if event.ID != "" {
			if layers[event.ID] == nil {
				layers[event.ID] = &layerProgress{}
			}
			lp := layers[event.ID]

			switch event.Status {
			case "Downloading":
				lp.dlCurrent = event.ProgressDetail.Current
				if event.ProgressDetail.Total > 0 {
					lp.dlTotal = event.ProgressDetail.Total
				}
			case "Download complete":
				lp.dlDone = true
				lp.dlCurrent = lp.dlTotal
			case "Extracting":
				lp.exCurrent = event.ProgressDetail.Current
				if event.ProgressDetail.Total > 0 {
					lp.exTotal = event.ProgressDetail.Total
				}
			case "Pull complete":
				lp.dlDone = true
				lp.exDone = true
				lp.dlCurrent = lp.dlTotal
				lp.exCurrent = lp.exTotal
			case "Already exists":
				lp.dlDone = true
				lp.exDone = true
			}
		}

		var dlTotal, dlDone, exTotal, exDone int64
		var layerCount, dlComplete, exComplete int
		for _, lp := range layers {
			layerCount++
			if lp.dlTotal > 0 {
				dlTotal += lp.dlTotal
				dlDone += lp.dlCurrent
			}
			if lp.exTotal > 0 {
				exTotal += lp.exTotal
				exDone += lp.exCurrent
			}
			if lp.dlDone {
				dlComplete++
			}
			if lp.exDone {
				exComplete++
			}
		}

		if dlTotal > maxTotal {
			maxTotal = dlTotal
		}

		var progress int
		var message string
		downloadDone := maxTotal > 0 && dlDone >= maxTotal
		allExtracted := layerCount > 0 && exComplete == layerCount

		if allExtracted {
			progress = 100
			message = "Finalizing"
		} else if downloadDone {
			if exTotal > 0 {
				progress = 90 + int(exDone*10/exTotal)
			} else if layerCount > 0 {
				progress = 90 + (exComplete * 10 / layerCount)
			} else {
				progress = 95
			}
			message = "Extracting"
		} else if maxTotal > 0 {
			progress = int(dlDone * 90 / maxTotal)
			message = fmt.Sprintf("%s/%s", formatBytesMB(dlDone), formatBytesMB(maxTotal))
		} else if layerCount > 0 {
			message = fmt.Sprintf("Preparing %d layers", layerCount)
		} else {
			message = "Connecting"
		}

		w.dockerMu.Lock()
		w.dockerStatus.PullProgress = progress
		w.dockerStatus.PullMessage = message
		w.dockerStatus.BytesTotal = maxTotal
		w.dockerStatus.BytesDone = dlDone
		w.dockerStatus.LayerCount = layerCount
		w.dockerStatus.LayersDone = exComplete
		w.dockerMu.Unlock()
	}

	w.dockerMu.Lock()
	w.dockerBuildLogs = logs.String()
	w.dockerStatus.PullStatus = "complete"
	w.dockerStatus.ImageBuilt = true
	w.dockerStatus.PullProgress = 100
	w.dockerStatus.PullMessage = "Complete"
	if sizeOutput, sizeErr := execCommand("docker", "images", crossImageName, "--format", "{{.Size}}"); sizeErr == nil && len(sizeOutput) > 0 {
		w.dockerStatus.ImageSize = strings.TrimSpace(sizeOutput)
	}
	w.dockerMu.Unlock()
	return nil
}

func (w *Wizard) pullViaDockerCLI() {
	cmd := exec.Command("docker", "pull", crossImageName)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		w.dockerMu.Lock()
		w.dockerStatus.PullStatus = "error"
		w.dockerStatus.PullError = fmt.Sprintf("Failed to create pipe: %v", err)
		w.dockerMu.Unlock()
		return
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		w.dockerMu.Lock()
		w.dockerStatus.PullStatus = "error"
		w.dockerStatus.PullError = fmt.Sprintf("Failed to start: %v", err)
		w.dockerMu.Unlock()
		return
	}

	done := make(chan struct{})
	var downloadDetected atomic.Bool

	go func() {
		ticker := time.NewTicker(300 * time.Millisecond)
		defer ticker.Stop()
		progress := 0

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if downloadDetected.Load() && progress < 80 {
					progress++
					w.dockerMu.Lock()
					w.dockerStatus.PullProgress = progress
					w.dockerMu.Unlock()
				}
			}
		}
	}()

	var lastOutput strings.Builder
	buf := make([]byte, 4096)

	for {
		n, readErr := stdout.Read(buf)
		if n > 0 {
			chunk := string(buf[:n])
			lastOutput.WriteString(chunk)

			if !downloadDetected.Load() && (strings.Contains(chunk, "Pulling") || strings.Contains(chunk, "Downloading")) {
				downloadDetected.Store(true)
				w.dockerMu.Lock()
				w.dockerStatus.PullMessage = "Downloading"
				w.dockerMu.Unlock()
			}
		}
		if readErr != nil {
			break
		}
	}

	close(done)

	err = cmd.Wait()
	w.dockerMu.Lock()
	w.dockerBuildLogs = lastOutput.String()
	if err != nil {
		w.dockerStatus.PullStatus = "error"
		w.dockerStatus.PullError = fmt.Sprintf("Pull failed: %v", err)
		w.dockerStatus.PullMessage = "Failed"
	} else {
		w.dockerStatus.PullStatus = "complete"
		w.dockerStatus.ImageBuilt = true
		w.dockerStatus.PullProgress = 100
		w.dockerStatus.PullMessage = "Complete"
		if sizeOutput, sizeErr := execCommand("docker", "images", crossImageName, "--format", "{{.Size}}"); sizeErr == nil && len(sizeOutput) > 0 {
			w.dockerStatus.ImageSize = strings.TrimSpace(sizeOutput)
		}
	}
	w.dockerMu.Unlock()
}

func formatBytesMB(b int64) string {
	mb := float64(b) / (1024 * 1024)
	if mb < 1 {
		return fmt.Sprintf("%.1f MB", mb)
	}
	return fmt.Sprintf("%.0f MB", mb)
}

func (w *Wizard) handleDockerBuild(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.startDockerPull()

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{"status": "started"})
}

func (w *Wizard) handleDockerLogs(rw http.ResponseWriter, r *http.Request) {
	w.dockerMu.RLock()
	logs := w.dockerBuildLogs
	w.dockerMu.RUnlock()

	rw.Header().Set("Content-Type", "text/plain")
	rw.Write([]byte(logs))
}

func (w *Wizard) handleDockerStartBackground(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	status := w.checkDocker()

	w.dockerMu.Lock()
	w.dockerStatus = status
	w.dockerMu.Unlock()

	if !status.Installed || !status.Running || status.ImageBuilt {
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"started": false,
			"reason":  getDockerNotStartedReason(status),
			"status":  status,
		})
		return
	}

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

	w.startDockerPull()

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

// allowedInstallers is the set of executables permitted by handleInstallDependency.
// The wizard only ever sends package-manager install commands, so we restrict to
// those to prevent an adversarial local process from invoking arbitrary binaries.
var allowedInstallers = map[string]bool{
	// Linux package managers (invoked directly or via sudo)
	"sudo": true, "apt": true, "apt-get": true, "dnf": true, "yum": true,
	"pacman": true, "zypper": true, "emerge": true, "xbps-install": true,
	"eopkg": true, "nix-env": true,
	// macOS
	"brew": true,
	// Mobile toolchains: install SDK packages / simulator runtimes
	"sdkmanager": true, "xcodebuild": true,
	// Windows
	"winget": true, "choco": true, "scoop": true,
}

// allowedSudoSubcmds lists the second word (the real executable) when the
// command starts with "sudo", so callers cannot do e.g. "sudo rm -rf /".
var allowedSudoSubcmds = map[string]bool{
	"apt": true, "apt-get": true, "dnf": true, "yum": true,
	"pacman": true, "zypper": true, "emerge": true, "xbps-install": true,
	"eopkg": true, "nix-env": true,
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

	parts := strings.Fields(req.Command)
	if len(parts) == 0 {
		json.NewEncoder(rw).Encode(InstallResponse{
			Success: false,
			Error:   "Empty command",
		})
		return
	}

	if !allowedInstallers[parts[0]] {
		json.NewEncoder(rw).Encode(InstallResponse{
			Success: false,
			Error:   fmt.Sprintf("command not allowed: %s", parts[0]),
		})
		return
	}
	if parts[0] == "sudo" && (len(parts) < 2 || !allowedSudoSubcmds[parts[1]]) {
		json.NewEncoder(rw).Encode(InstallResponse{
			Success: false,
			Error:   "sudo subcommand not allowed",
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
		// Load existing defaults first so that fields absent from the POST body
		// retain their current values (merge semantics, not replace semantics).
		existing, err := LoadGlobalDefaults()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&existing); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		if err := SaveGlobalDefaults(existing); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		path, _ := GetDefaultsPath()
		json.NewEncoder(rw).Encode(map[string]string{"status": "saved", "path": path})

	default:
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (w *Wizard) handleSigningStatus(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	status := checkSigningStatus()
	json.NewEncoder(rw).Encode(status)
}

func (w *Wizard) handleSigning(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		defaults, err := LoadGlobalDefaults()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(rw).Encode(defaults.Signing)

	case http.MethodPost:
		var signing SigningDefaults
		if err := json.NewDecoder(r.Body).Decode(&signing); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		defaults, err := LoadGlobalDefaults()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		defaults.Signing = signing

		if err := SaveGlobalDefaults(defaults); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(rw).Encode(map[string]string{"status": "saved"})

	default:
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (w *Wizard) handleNotarizeValidate(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	if runtime.GOOS != "darwin" {
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"valid": false,
			"error": "Notarization profiles are only available on macOS",
		})
		return
	}

	profile := r.URL.Query().Get("profile")
	if profile == "" {
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"valid": false,
			"error": "Profile name is required",
		})
		return
	}

	cmd := exec.Command("xcrun", "notarytool", "history", "--keychain-profile", profile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := string(output)
		if strings.Contains(errMsg, "No Keychain password item found") {
			json.NewEncoder(rw).Encode(map[string]interface{}{
				"valid": false,
				"error": "Profile not found in keychain",
			})
		} else {
			json.NewEncoder(rw).Encode(map[string]interface{}{
				"valid": false,
				"error": strings.TrimSpace(errMsg),
			})
		}
		return
	}

	json.NewEncoder(rw).Encode(map[string]interface{}{
		"valid": true,
	})
}

type notarizeCreateRequest struct {
	ProfileName string `json:"profileName"`
	AppleID     string `json:"appleID"`
	TeamID      string `json:"teamID"`
	Password    string `json:"password"`
}

func (w *Wizard) handleNotarizeCreate(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if runtime.GOOS != "darwin" {
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"success": false,
			"error":   "Notarization profiles can only be created on macOS",
		})
		return
	}

	var req notarizeCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	if req.ProfileName == "" || req.AppleID == "" || req.TeamID == "" || req.Password == "" {
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"success": false,
			"error":   "All fields are required",
		})
		return
	}

	cmd := exec.Command("xcrun", "notarytool", "store-credentials",
		req.ProfileName,
		"--apple-id", req.AppleID,
		"--team-id", req.TeamID,
		"--password-stdin",
		"--validate",
	)
	cmd.Stdin = strings.NewReader(req.Password)
	output, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := strings.TrimSpace(string(output))
		// Clean up the error message
		if strings.Contains(errMsg, "Error:") {
			lines := strings.Split(errMsg, "\n")
			for _, line := range lines {
				if strings.Contains(line, "Error:") {
					errMsg = strings.TrimSpace(strings.TrimPrefix(line, "Error:"))
					break
				}
			}
		}
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"success": false,
			"error":   errMsg,
		})
		return
	}

	// Save the profile name to defaults
	defs, err := LoadGlobalDefaults()
	if err != nil {
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"success": false,
			"error":   "Profile created but failed to load defaults: " + err.Error(),
		})
		return
	}
	defs.Signing.Darwin.KeychainProfile = req.ProfileName
	defs.Signing.Darwin.TeamID = req.TeamID
	if err := SaveGlobalDefaults(defs); err != nil {
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"success": false,
			"error":   "Profile created but failed to save defaults: " + err.Error(),
		})
		return
	}

	json.NewEncoder(rw).Encode(map[string]interface{}{
		"success": true,
		"output":  strings.TrimSpace(string(output)),
	})
}

type gpgCreateRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Passphrase string `json:"passphrase"`
}

// gpgFieldClean rejects values that could break out of a gpg batch param file
// (newlines) or be interpreted as a batch directive (a leading '%'). The values
// are passed through a parameter file, not a shell, so the only injection vector
// is the file format itself.
func gpgFieldClean(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}
	if strings.ContainsAny(s, "\r\n") || strings.HasPrefix(s, "%") {
		return "", false
	}
	return s, true
}

func (w *Wizard) handleGPGCreate(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gpg := gpgBinary()
	if gpg == "" {
		json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": "gpg is not installed"})
		return
	}

	var req gpgCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": "Invalid request: " + err.Error()})
		return
	}

	name, okName := gpgFieldClean(req.Name)
	email, okEmail := gpgFieldClean(req.Email)
	if !okName || !okEmail || !strings.Contains(email, "@") {
		json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": "A valid name and email are required"})
		return
	}

	// Build a batch key-parameter file in a private temp dir. The passphrase, if
	// any, lives only inside this 0600 file and is removed immediately after.
	tmpDir, err := os.MkdirTemp("", "wails-gpg-*")
	if err != nil {
		json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": "Failed to create temp dir: " + err.Error()})
		return
	}
	defer os.RemoveAll(tmpDir)

	var params strings.Builder
	params.WriteString("Key-Type: RSA\nKey-Length: 4096\n")
	params.WriteString("Subkey-Type: RSA\nSubkey-Length: 4096\n")
	params.WriteString("Name-Real: " + name + "\n")
	params.WriteString("Name-Email: " + email + "\n")
	params.WriteString("Expire-Date: 0\n")
	if req.Passphrase == "" {
		params.WriteString("%no-protection\n")
	} else {
		params.WriteString("Passphrase: " + strings.ReplaceAll(strings.Trim(req.Passphrase, "\r\n"), "\n", "") + "\n")
	}
	params.WriteString("%commit\n")

	paramFile := filepath.Join(tmpDir, "keyparams")
	if err := os.WriteFile(paramFile, []byte(params.String()), 0o600); err != nil {
		json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": "Failed to write key params: " + err.Error()})
		return
	}

	cmd := exec.Command(gpg, "--batch", "--pinentry-mode", "loopback", "--gen-key", paramFile)
	if out, err := cmd.CombinedOutput(); err != nil {
		json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": cleanToolError(string(out), err)})
		return
	}

	// Find the freshly created key by matching the email in the UID list.
	keyID := ""
	for _, k := range listGPGSecretKeys() {
		if strings.Contains(k.UID, email) {
			keyID = k.KeyID // last match wins → newest
		}
	}

	// Export the key to an armoured file. The build pipeline signs DEB/RPM
	// packages with a key *file path* (the Taskfile's PGP_KEY var → `wails3 tool
	// sign --pgp-key`), not a keyring ID — so generating the key alone wouldn't
	// be usable by a build. Export it and record the path alongside the ID.
	keyPath := ""
	if keyID != "" {
		if p, err := exportGPGSecretKey(gpg, keyID, req.Passphrase); err == nil {
			keyPath = p
		}
		if defs, err := LoadGlobalDefaults(); err == nil {
			defs.Signing.Linux.GPGKeyID = keyID
			if keyPath != "" {
				defs.Signing.Linux.GPGKeyPath = keyPath
			}
			if err := SaveGlobalDefaults(defs); err != nil {
				json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": "Key created but failed to save config: " + err.Error()})
				return
			}
		}
	}

	json.NewEncoder(rw).Encode(map[string]any{"success": true, "keyID": keyID, "keyPath": keyPath})
}

// gpgKeyIDPattern matches a hex GPG key id or fingerprint (16–40 hex chars,
// optional 0x prefix). keyID is used to build an export file path, so it is
// validated against this allowlist to prevent path traversal.
var gpgKeyIDPattern = regexp.MustCompile(`^(?:0x)?[A-Fa-f0-9]{16,40}$`)

// exportGPGSecretKey writes the ASCII-armoured private (and public) key for keyID
// into the signing output dir and returns the private key's path — the artifact
// the build consumes via the Taskfile PGP_KEY variable. The passphrase, if any,
// is supplied over stdin so it never appears in the process list, and the private
// file is written 0600.
func exportGPGSecretKey(gpg, keyID, passphrase string) (string, error) {
	// keyID flows into the export file path; reject anything that isn't a plain
	// hex key id/fingerprint so it can't escape the signing directory.
	if !gpgKeyIDPattern.MatchString(keyID) {
		return "", fmt.Errorf("invalid GPG key id %q", keyID)
	}
	dir, err := signingOutputDir()
	if err != nil {
		return "", err
	}
	privPath := filepath.Join(dir, keyID+".asc")
	pubPath := filepath.Join(dir, keyID+".pub.asc")

	args := []string{"--batch", "--yes", "--pinentry-mode", "loopback"}
	if passphrase != "" {
		args = append(args, "--passphrase-fd", "0")
	}
	args = append(args, "--armor", "--export-secret-keys", keyID)
	priv := exec.Command(gpg, args...)
	if passphrase != "" {
		priv.Stdin = strings.NewReader(passphrase)
	}
	out, err := priv.Output()
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(privPath, out, 0o600); err != nil {
		return "", err
	}

	// Public key is non-sensitive — export best-effort for distribution.
	if pub, err := exec.Command(gpg, "--armor", "--export", keyID).Output(); err == nil {
		_ = os.WriteFile(pubPath, pub, 0o644)
	}

	return privPath, nil
}

type gpgExportRequest struct {
	KeyID      string `json:"keyID"`
	Passphrase string `json:"passphrase"`
}

// handleGPGExport exports an already-existing keyring key to a file and records
// its path in the global config. This fixes keys configured by ID only (e.g.
// generated before the wizard exported automatically): the build signs with a
// key file, so an ID with no file path is not usable on its own.
func (w *Wizard) handleGPGExport(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gpg := gpgBinary()
	if gpg == "" {
		json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": "gpg is not installed"})
		return
	}

	var req gpgExportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": "Invalid request: " + err.Error()})
		return
	}

	keyID := strings.TrimSpace(req.KeyID)
	if keyID == "" {
		if defs, err := LoadGlobalDefaults(); err == nil {
			keyID = defs.Signing.Linux.GPGKeyID
		}
	}
	if keyID == "" {
		json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": "No GPG key id configured to export"})
		return
	}

	keyPath, err := exportGPGSecretKey(gpg, keyID, req.Passphrase)
	if err != nil {
		// The most common cause of a non-interactive export failure is a
		// passphrase-protected key with no passphrase supplied.
		resp := map[string]any{"success": false, "error": cleanToolError("", err)}
		if req.Passphrase == "" {
			resp["needsPassphrase"] = true
		}
		json.NewEncoder(rw).Encode(resp)
		return
	}

	if defs, err := LoadGlobalDefaults(); err == nil {
		defs.Signing.Linux.GPGKeyID = keyID
		defs.Signing.Linux.GPGKeyPath = keyPath
		if err := SaveGlobalDefaults(defs); err != nil {
			json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": "Key exported but failed to save config: " + err.Error()})
			return
		}
	}

	json.NewEncoder(rw).Encode(map[string]any{"success": true, "keyID": keyID, "keyPath": keyPath})
}

type windowsCertCreateRequest struct {
	CommonName string `json:"commonName"`
	Password   string `json:"password"`
}

func (w *Wizard) handleWindowsCertCreate(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !toolAvailable("openssl") {
		json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": "openssl is not installed"})
		return
	}

	var req windowsCertCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": "Invalid request: " + err.Error()})
		return
	}

	// CN is passed to openssl -subj. Strip the characters that would let it add
	// extra subject fields or break the arg. Args go via exec (no shell).
	cn := strings.TrimSpace(req.CommonName)
	cn = strings.Map(func(r rune) rune {
		if r == '/' || r == '\n' || r == '\r' || r == '=' {
			return -1
		}
		return r
	}, cn)
	if cn == "" {
		json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": "A certificate name is required"})
		return
	}
	if req.Password == "" {
		json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": "A password is required to protect the .pfx"})
		return
	}

	outDir, err := signingOutputDir()
	if err != nil {
		json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": err.Error()})
		return
	}

	tmpDir, err := os.MkdirTemp("", "wails-wincert-*")
	if err != nil {
		json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": "Failed to create temp dir: " + err.Error()})
		return
	}
	defer os.RemoveAll(tmpDir)

	keyPEM := filepath.Join(tmpDir, "key.pem")
	certPEM := filepath.Join(tmpDir, "cert.pem")
	pfxPath := filepath.Join(outDir, "wails-codesign-selfsigned.pfx")

	// 1) Self-signed cert + key with the codeSigning EKU.
	reqCmd := exec.Command("openssl", "req", "-x509", "-newkey", "rsa:4096",
		"-keyout", keyPEM, "-out", certPEM, "-days", "3650", "-nodes",
		"-subj", "/CN="+cn, "-addext", "extendedKeyUsage=codeSigning")
	if out, err := reqCmd.CombinedOutput(); err != nil {
		json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": cleanToolError(string(out), err)})
		return
	}

	// 2) Bundle into a password-protected .pfx. Password is passed via env, not
	// argv, so it does not appear in the process list.
	pfxCmd := exec.Command("openssl", "pkcs12", "-export", "-out", pfxPath,
		"-inkey", keyPEM, "-in", certPEM, "-name", cn, "-passout", "env:WAILS_PFX_PASS")
	pfxCmd.Env = append(os.Environ(), "WAILS_PFX_PASS="+req.Password)
	if out, err := pfxCmd.CombinedOutput(); err != nil {
		json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": cleanToolError(string(out), err)})
		return
	}

	if defs, err := LoadGlobalDefaults(); err == nil {
		defs.Signing.Windows.CertificatePath = pfxPath
		if err := SaveGlobalDefaults(defs); err != nil {
			json.NewEncoder(rw).Encode(map[string]any{"success": false, "error": "Certificate created but failed to save config: " + err.Error()})
			return
		}
	}

	json.NewEncoder(rw).Encode(map[string]any{
		"success":  true,
		"path":     pfxPath,
		"selfSign": true,
	})
}

// signingOutputDir returns a stable directory for generated signing artifacts.
func signingOutputDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot resolve home directory: %w", err)
	}
	dir := filepath.Join(home, ".wails", "signing")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("cannot create signing directory: %w", err)
	}
	return dir, nil
}

// cleanToolError condenses a command's combined output into a single-line error.
func cleanToolError(output string, err error) string {
	out := strings.TrimSpace(output)
	if out == "" {
		return err.Error()
	}
	lines := strings.Split(out, "\n")
	// Prefer a line that looks like an error message.
	for _, line := range lines {
		l := strings.TrimSpace(line)
		if strings.Contains(strings.ToLower(l), "error") || strings.Contains(l, "failed") {
			return l
		}
	}
	return strings.TrimSpace(lines[len(lines)-1])
}

type signingStatusResponse struct {
	Host        string               `json:"host"`
	Darwin      darwinSigningStatus  `json:"darwin"`
	Windows     windowsSigningStatus `json:"windows"`
	Linux       linuxSigningStatus   `json:"linux"`
	ConfigError string               `json:"configError,omitempty"`
}

type darwinSigningStatus struct {
	HasIdentity     bool     `json:"hasIdentity"`
	Identity        string   `json:"identity,omitempty"`
	Identities      []string `json:"identities,omitempty"`
	HasNotarization bool     `json:"hasNotarization"`
	TeamID          string   `json:"teamID,omitempty"`
	ConfigSource    string   `json:"configSource,omitempty"`
	// RcodesignAvailable reports whether the `rcodesign` tool is installed —
	// the way to sign macOS apps from a non-macOS host.
	RcodesignAvailable bool `json:"rcodesignAvailable"`
}

type windowsSigningStatus struct {
	HasCertificate  bool   `json:"hasCertificate"`
	CertificateType string `json:"certificateType,omitempty"`
	HasSignTool     bool   `json:"hasSignTool"`
	TimestampServer string `json:"timestampServer,omitempty"`
	ConfigSource    string `json:"configSource,omitempty"`
	// OsslsigncodeAvailable reports whether `osslsigncode` is installed — the
	// cross-platform way to Authenticode-sign Windows binaries from macOS/Linux.
	OsslsigncodeAvailable bool `json:"osslsigncodeAvailable"`
	// OpensslAvailable gates the "generate self-signed certificate" option.
	OpensslAvailable bool `json:"opensslAvailable"`
}

type linuxSigningStatus struct {
	HasGPGKey    bool         `json:"hasGpgKey"`
	GPGKeyID     string       `json:"gpgKeyID,omitempty"`
	GPGAvailable bool         `json:"gpgAvailable"`
	GPGKeys      []gpgKeyInfo `json:"gpgKeys,omitempty"`
	ConfigSource string       `json:"configSource,omitempty"`
}

type gpgKeyInfo struct {
	KeyID string `json:"keyID"`
	UID   string `json:"uid"`
}

func checkSigningStatus() signingStatusResponse {
	globalDefaults, err := LoadGlobalDefaults()

	resp := signingStatusResponse{
		Host:    runtime.GOOS,
		Darwin:  checkDarwinSigningStatus(globalDefaults),
		Windows: checkWindowsSigningStatus(globalDefaults),
		Linux:   checkLinuxSigningStatus(globalDefaults),
	}

	if err != nil {
		resp.ConfigError = err.Error()
	}

	return resp
}

func checkDarwinSigningStatus(cfg GlobalDefaults) darwinSigningStatus {
	status := darwinSigningStatus{}

	if cfg.Signing.Darwin.Identity != "" {
		status.HasIdentity = true
		status.Identity = cfg.Signing.Darwin.Identity
		status.TeamID = cfg.Signing.Darwin.TeamID
		status.ConfigSource = "defaults.yaml"
	}

	if cfg.Signing.Darwin.KeychainProfile != "" || cfg.Signing.Darwin.APIKeyID != "" {
		status.HasNotarization = true
	}

	if runtime.GOOS == "darwin" {
		identities := getMacOSSigningIdentities()
		status.Identities = identities
		if len(identities) > 0 && !status.HasIdentity {
			status.HasIdentity = true
			status.Identity = identities[0]
			status.ConfigSource = "keychain"
		}
	}

	status.RcodesignAvailable = toolAvailable("rcodesign")

	return status
}

func checkWindowsSigningStatus(cfg GlobalDefaults) windowsSigningStatus {
	status := windowsSigningStatus{}

	if cfg.Signing.Windows.CertificatePath != "" {
		status.HasCertificate = true
		status.CertificateType = "file"
		status.ConfigSource = "defaults.yaml"
	} else if cfg.Signing.Windows.Thumbprint != "" {
		status.HasCertificate = true
		status.CertificateType = "store"
		status.ConfigSource = "defaults.yaml"
	} else if cfg.Signing.Windows.CloudProvider != "" {
		status.HasCertificate = true
		status.CertificateType = "cloud:" + cfg.Signing.Windows.CloudProvider
		status.ConfigSource = "defaults.yaml"
	}

	status.TimestampServer = cfg.Signing.Windows.TimestampServer
	if status.TimestampServer == "" {
		status.TimestampServer = "http://timestamp.digicert.com"
	}

	if runtime.GOOS == "windows" {
		_, err := exec.LookPath("signtool.exe")
		status.HasSignTool = err == nil
	}

	status.OsslsigncodeAvailable = toolAvailable("osslsigncode")
	status.OpensslAvailable = toolAvailable("openssl")

	return status
}

func checkLinuxSigningStatus(cfg GlobalDefaults) linuxSigningStatus {
	status := linuxSigningStatus{}

	status.GPGAvailable = gpgBinary() != ""
	status.GPGKeys = listGPGSecretKeys()

	if cfg.Signing.Linux.GPGKeyPath != "" || cfg.Signing.Linux.GPGKeyID != "" {
		status.HasGPGKey = true
		status.GPGKeyID = cfg.Signing.Linux.GPGKeyID
		status.ConfigSource = "defaults.yaml"
	}

	if !status.HasGPGKey && len(status.GPGKeys) > 0 {
		status.HasGPGKey = true
		status.GPGKeyID = status.GPGKeys[0].KeyID
		status.ConfigSource = "gpg"
	}

	return status
}

// toolAvailable reports whether an executable is on PATH.
func toolAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// gpgBinary returns the path to a usable gpg binary ("gpg" preferred, then
// "gpg2"), or "" if neither is installed. gpg exists on macOS and Windows too
// (Homebrew, Gpg4win), so GPG signing setup is not Linux-only.
func gpgBinary() string {
	for _, name := range []string{"gpg", "gpg2"} {
		if p, err := exec.LookPath(name); err == nil {
			return p
		}
	}
	return ""
}

// listGPGSecretKeys returns the secret keys available to gpg, parsed from the
// stable --with-colons machine format (works identically on macOS/Linux/Windows).
func listGPGSecretKeys() []gpgKeyInfo {
	gpg := gpgBinary()
	if gpg == "" {
		return nil
	}

	cmd := exec.Command(gpg, "--list-secret-keys", "--with-colons", "--keyid-format", "long")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}
	return parseGPGSecretKeys(output)
}

// parseGPGSecretKeys parses `gpg --list-secret-keys --with-colons` output. Split
// out as a pure function so it can be tested without gpg installed.
func parseGPGSecretKeys(output []byte) []gpgKeyInfo {
	var keys []gpgKeyInfo
	var current *gpgKeyInfo
	for _, line := range strings.Split(string(output), "\n") {
		fields := strings.Split(line, ":")
		if len(fields) == 0 {
			continue
		}
		switch fields[0] {
		case "sec":
			// Field 5 (index 4) is the long key ID. Start a new key; its UID
			// arrives on the following uid record.
			if len(fields) > 4 && fields[4] != "" {
				keys = append(keys, gpgKeyInfo{KeyID: fields[4]})
				current = &keys[len(keys)-1]
			}
		case "uid":
			// Field 10 (index 9) is the user-ID string; attach the first one.
			if current != nil && current.UID == "" && len(fields) > 9 {
				current.UID = unescapeGPGColon(fields[9])
			}
		}
	}
	return keys
}

// unescapeGPGColon decodes the \xNN escapes gpg uses in --with-colons output
// (e.g. \x3a for ':').
func unescapeGPGColon(s string) string {
	if !strings.Contains(s, "\\x") {
		return s
	}
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if i+3 < len(s) && s[i] == '\\' && s[i+1] == 'x' {
			var v int
			if _, err := fmt.Sscanf(s[i+2:i+4], "%02x", &v); err == nil {
				b.WriteByte(byte(v))
				i += 3
				continue
			}
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func getMacOSSigningIdentities() []string {
	if runtime.GOOS != "darwin" {
		return nil
	}

	cmd := exec.Command("security", "find-identity", "-v", "-p", "codesigning")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	seen := make(map[string]bool)
	var identities []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// Identity lines look like:  1) <40-hex-hash> "Apple Development: Name (TEAM)"
		// Accept ALL code-signing identity types — Apple Development, Apple
		// Distribution, Developer ID, Mac Developer, etc. (the `-p codesigning`
		// policy already restricts the list to signing-capable certs). The
		// trailing "N valid identities found" line has no quotes, so requiring a
		// quoted name is enough to skip it.
		start := strings.Index(line, "\"")
		end := strings.LastIndex(line, "\"")
		if start == -1 || end <= start {
			continue
		}
		identity := line[start+1 : end]
		if identity == "" || seen[identity] {
			continue
		}
		seen[identity] = true
		identities = append(identities, identity)
	}

	return identities
}

