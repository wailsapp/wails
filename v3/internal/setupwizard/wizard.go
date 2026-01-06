package setupwizard

import (
	"bufio"
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
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
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
	HelpURL        string `json:"helpUrl,omitempty"`
	ImageBuilt     bool   `json:"imageBuilt"` // For Docker: whether wails-cross image exists
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
	shutdown        chan struct{}
	buildWg         sync.WaitGroup
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
	mux.HandleFunc("/api/docker/status/stream", w.handleDockerStatusStream)
	mux.HandleFunc("/api/docker/build", w.handleDockerBuild)
	mux.HandleFunc("/api/docker/build-with-sdk", w.handleDockerBuildWithSDK)
	mux.HandleFunc("/api/docker/logs", w.handleDockerLogs)
	mux.HandleFunc("/api/docker/start-background", w.handleDockerStartBackground)
	mux.HandleFunc("/api/wails-config", w.handleWailsConfig)
	mux.HandleFunc("/api/defaults", w.handleDefaults)
	mux.HandleFunc("/api/complete", w.handleComplete)
	mux.HandleFunc("/api/close", w.handleClose)

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
	rw.Header().Set("Access-Control-Allow-Origin", "*")

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
	w.dockerMu.Unlock()

	w.buildWg.Add(1)
	go func() {
		defer w.buildWg.Done()

		if err := w.pullViaDockerAPI(); err != nil {
			w.pullViaDockerCLI()
		}
	}()
}

func (w *Wizard) pullViaDockerAPI() error {
	socketPath := "/var/run/docker.sock"
	if runtime.GOOS == "windows" {
		socketPath = "//./pipe/docker_engine"
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				if runtime.GOOS == "windows" {
					return nil, fmt.Errorf("windows named pipes not supported, falling back to CLI")
				}
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
			if err.Error() == "EOF" {
				break
			}
			break
		}

		logs.WriteString(fmt.Sprintf("%s %s %s\n", event.ID, event.Status, event.Progress))

		if event.Error != "" {
			w.dockerMu.Lock()
			w.dockerStatus.PullStatus = "error"
			w.dockerStatus.PullError = event.Error
			w.dockerStatus.PullMessage = "Failed"
			w.dockerBuildLogs = logs.String()
			w.dockerMu.Unlock()
			return nil
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

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func formatBytesMB(b int64) string {
	mb := float64(b) / (1024 * 1024)
	if mb < 1 {
		return fmt.Sprintf("%.1f MB", mb)
	}
	return fmt.Sprintf("%.0f MB", mb)
}

func (w *Wizard) startDockerBuildLocal(localSDKPath string) {
	w.dockerMu.Lock()
	if w.dockerStatus.PullStatus == "pulling" {
		w.dockerMu.Unlock()
		return
	}
	w.dockerStatus.PullStatus = "pulling"
	w.dockerStatus.PullProgress = 0
	w.dockerMu.Unlock()

	w.buildWg.Add(1)
	go func() {
		defer w.buildWg.Done()

		tmpDir, err := os.MkdirTemp("", "wails-docker-build-*")
		if err != nil {
			w.dockerMu.Lock()
			w.dockerStatus.PullStatus = "error"
			w.dockerStatus.PullError = fmt.Sprintf("Failed to create temp dir: %v", err)
			w.dockerMu.Unlock()
			return
		}
		defer os.RemoveAll(tmpDir)

		dockerfile := dockerfileContent
		sdkFileName := filepath.Base(localSDKPath)
		destPath := filepath.Join(tmpDir, sdkFileName)

		srcFile, err := os.Open(localSDKPath)
		if err != nil {
			w.dockerMu.Lock()
			w.dockerStatus.PullStatus = "error"
			w.dockerStatus.PullError = fmt.Sprintf("Failed to open SDK file: %v", err)
			w.dockerMu.Unlock()
			return
		}
		defer srcFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			w.dockerMu.Lock()
			w.dockerStatus.PullStatus = "error"
			w.dockerStatus.PullError = fmt.Sprintf("Failed to create SDK copy: %v", err)
			w.dockerMu.Unlock()
			return
		}

		if _, err := destFile.ReadFrom(srcFile); err != nil {
			destFile.Close()
			w.dockerMu.Lock()
			w.dockerStatus.PullStatus = "error"
			w.dockerStatus.PullError = fmt.Sprintf("Failed to copy SDK file: %v", err)
			w.dockerMu.Unlock()
			return
		}
		destFile.Close()

		sdkDirName := strings.TrimSuffix(strings.TrimSuffix(sdkFileName, ".xz"), ".tar")
		localSDKDockerfile := strings.Replace(dockerfile,
			`RUN curl -L "https://github.com/wailsapp/macosx-sdks/releases/download/${MACOS_SDK_VERSION}/MacOSX${MACOS_SDK_VERSION}.sdk.tar.xz" \
    | tar -xJ -C /opt \
    && mv /opt/MacOSX${MACOS_SDK_VERSION}.sdk /opt/macos-sdk`,
			fmt.Sprintf(`COPY %s /tmp/sdk.tar.xz
RUN tar -xJf /tmp/sdk.tar.xz -C /opt \
    && mv /opt/%s /opt/macos-sdk \
    && rm /tmp/sdk.tar.xz`, sdkFileName, sdkDirName),
			1)

		dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
		if err := os.WriteFile(dockerfilePath, []byte(localSDKDockerfile), 0644); err != nil {
			w.dockerMu.Lock()
			w.dockerStatus.PullStatus = "error"
			w.dockerStatus.PullError = fmt.Sprintf("Failed to write Dockerfile: %v", err)
			w.dockerMu.Unlock()
			return
		}

		cmd := exec.Command("docker", "build", "--progress=plain", "-t", crossImageName, "-f", dockerfilePath, tmpDir)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			w.dockerMu.Lock()
			w.dockerStatus.PullStatus = "error"
			w.dockerStatus.PullError = fmt.Sprintf("Failed to create stdout pipe: %v", err)
			w.dockerMu.Unlock()
			return
		}
		cmd.Stderr = cmd.Stdout

		if err := cmd.Start(); err != nil {
			w.dockerMu.Lock()
			w.dockerStatus.PullStatus = "error"
			w.dockerStatus.PullError = fmt.Sprintf("Failed to start build: %v", err)
			w.dockerMu.Unlock()
			return
		}

		scanner := bufio.NewScanner(stdout)
		stepProgressRegex := regexp.MustCompile(`\[\s*(\d+)/(\d+)\]`)
		var lastOutput strings.Builder

		for scanner.Scan() {
			line := scanner.Text()
			lastOutput.WriteString(line + "\n")

			if matches := stepProgressRegex.FindStringSubmatch(line); len(matches) == 3 {
				current, err1 := strconv.Atoi(matches[1])
				total, err2 := strconv.Atoi(matches[2])
				if err1 == nil && err2 == nil && total > 0 {
					progress := int(float64(current) / float64(total) * 90)
					w.dockerMu.Lock()
					if progress > w.dockerStatus.PullProgress {
						w.dockerStatus.PullProgress = progress
					}
					w.dockerMu.Unlock()
				}
			}
		}

		err = cmd.Wait()
		w.dockerMu.Lock()
		w.dockerBuildLogs = lastOutput.String()
		if err != nil {
			w.dockerStatus.PullStatus = "error"
			w.dockerStatus.PullError = fmt.Sprintf("Build failed: %v", err)
		} else {
			w.dockerStatus.PullStatus = "complete"
			w.dockerStatus.ImageBuilt = true
			w.dockerStatus.PullProgress = 100
			if sizeOutput, sizeErr := execCommand("docker", "images", crossImageName, "--format", "{{.Size}}"); sizeErr == nil && len(sizeOutput) > 0 {
				w.dockerStatus.ImageSize = strings.TrimSpace(sizeOutput)
			}
		}
		w.dockerMu.Unlock()
	}()
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

func (w *Wizard) handleDockerBuildWithSDK(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(200 << 20); err != nil {
		http.Error(rw, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("sdk")
	if err != nil {
		http.Error(rw, "No SDK file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	tmpFile, err := os.CreateTemp("", "macos-sdk-*.tar.xz")
	if err != nil {
		http.Error(rw, "Failed to create temp file", http.StatusInternalServerError)
		return
	}
	tmpPath := tmpFile.Name()

	if _, err := tmpFile.ReadFrom(file); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		http.Error(rw, "Failed to save SDK file", http.StatusInternalServerError)
		return
	}
	tmpFile.Close()

	_ = header
	w.startDockerBuildLocal(tmpPath)

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
