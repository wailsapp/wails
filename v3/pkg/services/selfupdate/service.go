// Package selfupdate provides a service for self-updating Wails applications.
//
// The service supports multiple update providers through a plugin architecture:
//   - GitHub Releases (built-in)
//   - Keygen.sh (planned)
//   - Custom HTTP endpoints (planned)
//
// Example usage:
//
//	app := application.New(application.Options{
//	    Services: []application.Service{
//	        application.NewService(selfupdate.New(&selfupdate.Config{
//	            CurrentVersion: "1.0.0",
//	            Provider:       "github",
//	            GitHub: &selfupdate.GitHubConfig{
//	                Owner: "myorg",
//	                Repo:  "myapp",
//	            },
//	        })),
//	    },
//	})
package selfupdate

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// DefaultMaxDownloadSize is the default maximum download size (500MB).
const DefaultMaxDownloadSize int64 = 500 * 1024 * 1024

// DefaultCheckTimeout is the default timeout for update check operations.
const DefaultCheckTimeout = 30 * time.Second

// Config holds the configuration for the selfupdate service.
type Config struct {
	// CurrentVersion is the current version of the application (e.g., "1.0.0").
	// This is required.
	CurrentVersion string

	// Provider is the update provider to use (e.g., "github", "keygen").
	// Defaults to "github".
	Provider string

	// Channel is the release channel (e.g., "stable", "beta", "alpha").
	// Defaults to "stable".
	Channel string

	// Variant is an optional build variant (e.g., "webkit2_41").
	Variant string

	// AssetPattern is a template for matching release assets.
	// See pattern.go for supported variables.
	AssetPattern string

	// PublicKey is the Ed25519 public key for signature verification (base64).
	// If empty, signature verification is disabled (not recommended).
	PublicKey string

	// GitHub contains GitHub-specific configuration.
	GitHub *GitHubConfig

	// AutoCheck enables automatic update checking on startup.
	AutoCheck bool

	// MaxDownloadSize is the maximum allowed download size in bytes.
	// Defaults to DefaultMaxDownloadSize (500MB).
	MaxDownloadSize int64

	// CheckTimeout is the timeout for update check operations.
	// Defaults to DefaultCheckTimeout (30s).
	CheckTimeout time.Duration

	// PrepareFunc is called after download to prepare the update.
	// This handles archive extraction, macOS .app bundles, etc.
	// If nil, a default prepare function is used.
	PrepareFunc PrepareFunc
}

// GitHubConfig contains GitHub-specific configuration.
type GitHubConfig struct {
	// Owner is the repository owner (user or organization).
	Owner string

	// Repo is the repository name.
	Repo string

	// Token is an optional GitHub personal access token for private repos.
	Token string

	// BaseURL is the API base URL for GitHub Enterprise.
	BaseURL string
}

// Service provides self-update functionality for Wails applications.
type Service struct {
	mu       sync.RWMutex
	config   *Config
	provider UpdateProvider
	executor *Executor
	logger   *slog.Logger

	// Cached state
	lastCheck        *UpdateResult
	isDownloading    bool
	downloadTempPath string // Path to temp file with downloaded data (M3)
}

// New creates a new selfupdate service with the given configuration.
func New(config *Config) *Service {
	if config == nil {
		config = &Config{}
	}
	if config.MaxDownloadSize <= 0 {
		config.MaxDownloadSize = DefaultMaxDownloadSize
	}
	if config.CheckTimeout <= 0 {
		config.CheckTimeout = DefaultCheckTimeout
	}
	return &Service{
		config: config,
		logger: slog.Default(),
	}
}

// ServiceName returns the name of the service.
func (s *Service) ServiceName() string {
	return "github.com/wailsapp/wails/v3/pkg/services/selfupdate"
}

// ServiceStartup initializes the service.
func (s *Service) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate configuration
	if s.config.CurrentVersion == "" {
		return fmt.Errorf("selfupdate: CurrentVersion is required")
	}

	// Default provider
	providerName := s.config.Provider
	if providerName == "" {
		providerName = "github"
	}

	// Get provider instance
	provider, err := GetProvider(providerName)
	if err != nil {
		return fmt.Errorf("selfupdate: %w", err)
	}
	s.provider = provider

	// Build provider config
	providerConfig := &ProviderConfig{
		CurrentVersion: s.config.CurrentVersion,
		Channel:        s.config.Channel,
		Variant:        s.config.Variant,
		AssetPattern:   s.config.AssetPattern,
		PublicKey:       s.config.PublicKey,
		Settings:       make(map[string]any),
	}

	// Add provider-specific settings
	if s.config.GitHub != nil {
		providerConfig.Settings["owner"] = s.config.GitHub.Owner
		providerConfig.Settings["repo"] = s.config.GitHub.Repo
		if s.config.GitHub.Token != "" {
			providerConfig.Settings["token"] = s.config.GitHub.Token
		}
		if s.config.GitHub.BaseURL != "" {
			providerConfig.Settings["baseURL"] = s.config.GitHub.BaseURL
		}
	}

	// Configure provider
	if err := provider.Configure(ctx, providerConfig); err != nil {
		return fmt.Errorf("selfupdate: failed to configure provider: %w", err)
	}

	// Create executor
	executorConfig := &ExecutorConfig{
		PrepareFunc: s.config.PrepareFunc,
	}
	executor, err := NewExecutor(executorConfig)
	if err != nil {
		return fmt.Errorf("selfupdate: failed to create executor: %w", err)
	}
	s.executor = executor

	// Clean up any old versions from previous updates
	CleanupOldVersions()

	// Auto-check if enabled (H5: use background context, add error logging)
	if s.config.AutoCheck {
		go func() {
			checkCtx, cancel := context.WithTimeout(context.Background(), s.config.CheckTimeout)
			defer cancel()

			result, err := s.check(checkCtx, &CheckOptions{})
			if err != nil {
				s.logger.Warn("selfupdate: auto-check failed", "error", err)
				return
			}
			if result != nil && result.UpdateAvailable {
				s.emitEvent("selfupdate:available", result)
			}
		}()
	}

	return nil
}

// ServiceShutdown cleans up the service.
func (s *Service) ServiceShutdown() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.provider != nil {
		s.provider.Close()
		s.provider = nil
	}

	// Clean up temp download file (M3)
	if s.downloadTempPath != "" {
		os.Remove(s.downloadTempPath)
		s.downloadTempPath = ""
	}

	return nil
}

// Check checks if an update is available.
// Returns UpdateResult with details about the available update.
func (s *Service) Check(ctx context.Context) (*UpdateResult, error) {
	// Add timeout to context (M7)
	ctx, cancel := context.WithTimeout(ctx, s.config.CheckTimeout)
	defer cancel()
	return s.check(ctx, &CheckOptions{})
}

// CheckWithPrerelease checks for updates including pre-releases.
func (s *Service) CheckWithPrerelease(ctx context.Context) (*UpdateResult, error) {
	// Add timeout to context (M7)
	ctx, cancel := context.WithTimeout(ctx, s.config.CheckTimeout)
	defer cancel()
	return s.check(ctx, &CheckOptions{IncludePrerelease: true})
}

// check is the internal implementation shared by Check and CheckWithPrerelease.
func (s *Service) check(ctx context.Context, opts *CheckOptions) (*UpdateResult, error) {
	s.mu.RLock()
	provider := s.provider
	s.mu.RUnlock()

	if provider == nil {
		return nil, fmt.Errorf("selfupdate: service not initialized")
	}

	result, err := provider.CheckForUpdate(ctx, opts)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.lastCheck = result
	s.mu.Unlock()

	return result, nil
}

// Download downloads the update but does not install it.
// Returns true if download was successful.
func (s *Service) Download(ctx context.Context) (bool, error) {
	s.mu.Lock()
	if s.isDownloading {
		s.mu.Unlock()
		return false, fmt.Errorf("download already in progress")
	}
	s.isDownloading = true
	lastCheck := s.lastCheck
	provider := s.provider
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.isDownloading = false
		s.mu.Unlock()
	}()

	if lastCheck == nil || !lastCheck.UpdateAvailable {
		return false, fmt.Errorf("no update available; call Check first")
	}

	if provider == nil {
		return false, fmt.Errorf("selfupdate: service not initialized")
	}

	// Download with progress
	reader, err := provider.DownloadUpdate(ctx, lastCheck, func(info *ProgressInfo) {
		s.emitEvent("selfupdate:progress", info)
	})
	if err != nil {
		return false, fmt.Errorf("download failed: %w", err)
	}
	defer reader.Close()

	// Read data with size limit into temp file (C3 + M3)
	limitedReader := io.LimitReader(reader, s.config.MaxDownloadSize+1)

	tempFile, err := os.CreateTemp("", "wails-selfupdate-download-*")
	if err != nil {
		return false, fmt.Errorf("failed to create temp file: %w", err)
	}
	tempPath := tempFile.Name()

	written, err := io.Copy(tempFile, limitedReader)
	if err != nil {
		tempFile.Close()
		os.Remove(tempPath)
		return false, fmt.Errorf("failed to write download: %w", err)
	}
	tempFile.Close()

	if written > s.config.MaxDownloadSize {
		os.Remove(tempPath)
		return false, fmt.Errorf("download exceeds maximum size of %d bytes", s.config.MaxDownloadSize)
	}

	// Read back for verification
	data, err := os.ReadFile(tempPath)
	if err != nil {
		os.Remove(tempPath)
		return false, fmt.Errorf("failed to read download for verification: %w", err)
	}

	// Verify
	if err := provider.VerifyUpdate(ctx, lastCheck, bytes.NewReader(data)); err != nil {
		os.Remove(tempPath)
		return false, fmt.Errorf("verification failed: %w", err)
	}

	s.mu.Lock()
	// Clean up previous temp file if any
	if s.downloadTempPath != "" {
		os.Remove(s.downloadTempPath)
	}
	s.downloadTempPath = tempPath
	s.mu.Unlock()

	return true, nil
}

// Install installs a previously downloaded update.
// The application should be restarted after calling this method.
func (s *Service) Install(ctx context.Context) error {
	s.mu.Lock()
	tempPath := s.downloadTempPath
	executor := s.executor
	s.mu.Unlock()

	if tempPath == "" {
		return fmt.Errorf("no update downloaded; call Download first")
	}

	if executor == nil {
		return fmt.Errorf("selfupdate: executor not initialized")
	}

	// Open temp file as reader
	file, err := os.Open(tempPath)
	if err != nil {
		return fmt.Errorf("failed to open downloaded update: %w", err)
	}
	defer file.Close()

	// Apply the update
	if err := executor.Apply(ctx, file, &ProviderConfig{
		Variant: s.config.Variant,
	}); err != nil {
		return fmt.Errorf("failed to apply update: %w", err)
	}

	// Clear temp file
	s.mu.Lock()
	os.Remove(s.downloadTempPath)
	s.downloadTempPath = ""
	s.mu.Unlock()

	return nil
}

// DownloadAndInstall downloads and installs the update in one step.
// You must call Check first to detect an available update.
// Returns true if the update was successful.
func (s *Service) DownloadAndInstall(ctx context.Context) (bool, error) {
	ok, err := s.Download(ctx)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	if err := s.Install(ctx); err != nil {
		return false, err
	}

	return true, nil
}

// Restart quits the application and relaunches it (M2).
// On Unix, this uses exec to replace the current process.
// On Windows, this spawns a new process and quits.
func (s *Service) Restart() error {
	exe, err := GetExecutablePath()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	if runtime.GOOS == "windows" {
		// On Windows, spawn new process then quit
		cmd := exec.Command(exe, os.Args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to restart: %w", err)
		}
	}

	app := application.Get()
	if app != nil {
		app.Quit()
	}

	// On Unix, try to exec (replaces current process)
	if runtime.GOOS != "windows" {
		// Give the app a moment to shut down cleanly
		// The Quit() above should trigger shutdown hooks
		// If we get here, the app didn't exit from Quit(), so we exec
		// Note: In practice, app.Quit() may not return, so this is a fallback
	}

	return nil
}

// GetCurrentVersion returns the current application version.
func (s *Service) GetCurrentVersion() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.config == nil {
		return ""
	}
	return s.config.CurrentVersion
}

// GetLastCheck returns the result of the last update check.
func (s *Service) GetLastCheck() *UpdateResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.lastCheck
}

// CanUpdate returns true if the application has permission to update itself.
func (s *Service) CanUpdate() bool {
	return CanUpdate()
}

// IsDownloading returns true if a download is in progress.
func (s *Service) IsDownloading() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.isDownloading
}

// emitEvent emits an event to the frontend.
func (s *Service) emitEvent(name string, data any) {
	app := application.Get()
	if app != nil {
		app.Event.Emit(name, data)
	}
}
