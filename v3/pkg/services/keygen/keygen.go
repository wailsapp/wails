package keygen

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/keygen-sh/keygen-go/v3"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// platformKeygen defines the platform-specific interface for the Keygen service
type platformKeygen interface {
	// Platform-specific machine fingerprinting
	GetMachineFingerprint() (string, error)

	// Platform-specific update installation
	InstallUpdatePlatform(updatePath string) error

	// Platform-specific cache directory
	GetCacheDir() string
}

// Service represents the Keygen licensing and update service
type Service struct {
	app          *application.App
	eventEmitter *EventEmitter
	impl         platformKeygen
	client       *keygen.Client
	config       Config
	state        *LicenseState
	mu           sync.RWMutex

	// Update management
	currentVersion   string
	updateInProgress bool
	downloadCancel   context.CancelFunc

	// Cache management
	cacheDir string
}

// ServiceOptions contains options for the keygen service
type ServiceOptions struct {
	AccountID      string
	ProductID      string
	LicenseKey     string
	PublicKey      string
	Environment    string
	CacheDir       string
	AutoCheck      bool
	CheckInterval  time.Duration
	UpdateChannel  string
	CurrentVersion string
}

// New creates a new keygen service instance
func New(options ServiceOptions) *Service {
	// Set defaults
	if options.Environment == "" {
		options.Environment = "production"
	}
	if options.CheckInterval == 0 {
		options.CheckInterval = 24 * time.Hour
	}
	if options.UpdateChannel == "" {
		options.UpdateChannel = "stable"
	}

	// Create config
	config := Config{
		Account:       options.AccountID,
		Product:       options.ProductID,
		LicenseKey:    options.LicenseKey,
		PublicKey:     options.PublicKey,
		CacheDir:      options.CacheDir,
		AutoCheck:     options.AutoCheck,
		CheckInterval: options.CheckInterval,
		UpdateChannel: options.UpdateChannel,
		Environment:   options.Environment,
	}

	// Create Keygen client
	var keygenOptions []keygen.Option
	if options.Environment != "production" {
		keygenOptions = append(keygenOptions, keygen.WithEnvironment(options.Environment))
	}
	if options.PublicKey != "" {
		keygenOptions = append(keygenOptions, keygen.WithPublicKey(options.PublicKey))
	}

	client := keygen.NewClient(options.AccountID, keygenOptions...)

	return &Service{
		config:         config,
		client:         client,
		state:          &LicenseState{},
		currentVersion: options.CurrentVersion,
		impl:           NewPlatformKeygen(),
	}
}

// ServiceName returns the name of the service
func (s *Service) ServiceName() string {
	return "github.com/wailsapp/wails/v3/pkg/services/keygen"
}

// ServiceStartup is called when the application starts
func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	// Store the application reference
	if app, ok := options.App.(*application.App); ok {
		s.app = app
		s.eventEmitter = NewEventEmitter(s.app)
	}

	// Initialize cache directory
	if s.config.CacheDir == "" {
		s.cacheDir = s.impl.GetCacheDir()
	} else {
		s.cacheDir = s.config.CacheDir
	}

	// Ensure cache directory exists
	if err := os.MkdirAll(s.cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Load offline license if available
	if err := s.LoadOfflineLicense(); err != nil {
		// Log error but don't fail startup
		_ = err
	}

	// Start automatic update checking if enabled
	if s.config.AutoCheck {
		go s.startAutoUpdateCheck(ctx)
	}

	return nil
}

// ServiceShutdown is called when the application shuts down
func (s *Service) ServiceShutdown() error {
	// Cancel any ongoing downloads
	if s.downloadCancel != nil {
		s.downloadCancel()
	}

	// Save offline license
	if s.state != nil && s.state.Valid {
		_ = s.SaveOfflineLicense()
	}

	return nil
}

// ValidateLicense validates the current license
func (s *Service) ValidateLicense() (*LicenseValidationResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.config.LicenseKey == "" {
		return &LicenseValidationResult{
			Valid:   false,
			Message: "No license key provided",
			Code:    ValidationCodeInvalid,
		}, NewConfigError("licenseKey", "license key is required")
	}

	// Create license
	license := s.client.License(s.config.LicenseKey)

	// Validate the license
	validation, err := license.Validate()
	if err != nil {
		s.state.Valid = false
		s.state.LastChecked = time.Now()

		// Emit license status event
		s.emitLicenseStatus(false, "License validation failed", err.Error())

		return &LicenseValidationResult{
			Valid:   false,
			Message: "License validation failed",
			Code:    ValidationCodeInvalid,
			State:   s.state,
		}, err
	}

	// Update state based on validation
	s.state.Valid = validation.Valid
	s.state.License = license
	s.state.LastChecked = time.Now()
	s.state.OfflineMode = false

	if validation.Valid {
		// Extract license metadata
		if meta := validation.License.Metadata; meta != nil {
			s.state.Entitlements = meta
		}

		// Check expiration
		if validation.License.Expiry != nil {
			expiry := *validation.License.Expiry
			s.state.ExpiresAt = &expiry
		}

		// Emit success event
		s.emitLicenseStatus(true, "License is valid", "")

		return &LicenseValidationResult{
			Valid:              true,
			Message:            "License is valid",
			Code:               ValidationCodeValid,
			State:              s.state,
			RequiresActivation: validation.License.MaxMachines != nil && *validation.License.MaxMachines > 0,
		}, nil
	}

	// Handle validation failure
	code := ValidationCodeInvalid
	message := "License is invalid"

	if validation.Code == "LICENSE_EXPIRED" {
		code = ValidationCodeExpired
		message = "License has expired"
	} else if validation.Code == "LICENSE_SUSPENDED" {
		code = ValidationCodeSuspended
		message = "License has been suspended"
	}

	s.emitLicenseStatus(false, message, validation.Code)

	return &LicenseValidationResult{
		Valid:   false,
		Message: message,
		Code:    code,
		State:   s.state,
	}, nil
}

// ActivateMachine activates the current machine
func (s *Service) ActivateMachine() (*MachineActivationResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state.License == nil {
		return nil, NewLicenseInvalidError("license must be validated before machine activation")
	}

	// Get machine fingerprint
	fingerprint, err := s.impl.GetMachineFingerprint()
	if err != nil {
		return nil, fmt.Errorf("failed to get machine fingerprint: %w", err)
	}

	// Get hostname
	hostname, _ := os.Hostname()

	// Create machine activation
	machine, err := s.state.License.ActivateMachine(fingerprint, keygen.ActivateMachineOptions{
		Name:     hostname,
		Platform: runtime.GOOS,
		Metadata: map[string]interface{}{
			"arch":    runtime.GOARCH,
			"version": s.currentVersion,
		},
	})

	if err != nil {
		return &MachineActivationResult{
			Success:     false,
			Message:     "Machine activation failed",
			Fingerprint: fingerprint,
		}, err
	}

	// Update state
	s.state.Machine = machine

	// Emit machine activated event
	if s.eventEmitter != nil {
		s.eventEmitter.EmitLicenseStatus(LicenseStatusEvent{
			Valid:       true,
			Status:      "active",
			Key:         s.config.LicenseKey,
			EnteredAt:   time.Now().Format(time.RFC3339),
			LastChecked: time.Now().Format(time.RFC3339),
		})
	}

	return &MachineActivationResult{
		Success:     true,
		Message:     "Machine activated successfully",
		Machine:     machine,
		Fingerprint: fingerprint,
	}, nil
}

// DeactivateMachine deactivates the current machine
func (s *Service) DeactivateMachine() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state.Machine == nil {
		return &KeygenError{Code: ErrMachineNotActivated, Message: "no machine is currently activated"}
	}

	// Deactivate the machine
	err := s.state.Machine.Deactivate()
	if err != nil {
		return fmt.Errorf("failed to deactivate machine: %w", err)
	}

	// Clear machine from state
	s.state.Machine = nil

	// Emit deactivation event
	if s.eventEmitter != nil {
		s.eventEmitter.EmitLicenseStatus(LicenseStatusEvent{
			Valid:  true,
			Status: "deactivated",
			Key:    s.config.LicenseKey,
			Error:  "Machine deactivated",
		})
	}

	return nil
}

// GetLicenseInfo returns detailed license information
func (s *Service) GetLicenseInfo() (*LicenseInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.state.License == nil {
		return nil, NewLicenseInvalidError("no license information available")
	}

	license := s.state.License
	info := &LicenseInfo{
		Key:             s.config.LicenseKey,
		Status:          "active",
		CreatedAt:       license.CreatedAt,
		LastValidatedAt: &s.state.LastChecked,
		Entitlements:    s.state.Entitlements,
	}

	// Extract metadata
	if meta := license.Metadata; meta != nil {
		if name, ok := meta["name"].(string); ok {
			info.Name = name
		}
		if email, ok := meta["email"].(string); ok {
			info.Email = email
		}
		if company, ok := meta["company"].(string); ok {
			info.Company = company
		}
		info.Metadata = meta
	}

	// Set expiration
	if license.Expiry != nil {
		info.ExpiresAt = license.Expiry
	}

	// Set machine limits
	if license.MaxMachines != nil {
		info.MaxMachines = license.MaxMachines
	}

	return info, nil
}

// CheckEntitlement checks if a specific feature is enabled
func (s *Service) CheckEntitlement(feature string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.state.Valid {
		return false, NewLicenseInvalidError("license is not valid")
	}

	if s.state.Entitlements == nil {
		return false, nil
	}

	// Check if feature exists in entitlements
	if val, ok := s.state.Entitlements[feature]; ok {
		// Handle boolean values
		if enabled, ok := val.(bool); ok {
			return enabled, nil
		}
		// Any non-false value is considered enabled
		return true, nil
	}

	return false, nil
}

// CheckForUpdates checks for available updates
func (s *Service) CheckForUpdates() (*UpdateInfo, error) {
	s.mu.RLock()
	product := s.config.Product
	channel := s.config.UpdateChannel
	currentVersion := s.currentVersion
	s.mu.RUnlock()

	// List releases for the product
	releases, err := s.client.Releases(product, keygen.ReleaseListOptions{
		Channel: channel,
		Limit:   1,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}

	if len(releases) == 0 {
		return &UpdateInfo{
			Available:      false,
			CurrentVersion: currentVersion,
			LatestVersion:  currentVersion,
		}, nil
	}

	latestRelease := releases[0]

	// Check if update is available
	isNewer := compareVersions(latestRelease.Version, currentVersion) > 0

	if !isNewer {
		return &UpdateInfo{
			Available:      false,
			CurrentVersion: currentVersion,
			LatestVersion:  latestRelease.Version,
		}, nil
	}

	// Build update info
	updateInfo := &UpdateInfo{
		Available:      true,
		CurrentVersion: currentVersion,
		LatestVersion:  latestRelease.Version,
		ReleaseID:      latestRelease.ID,
		ReleaseName:    latestRelease.Name,
		ReleaseNotes:   latestRelease.Description,
		PublishedAt:    &latestRelease.PublishedAt,
		Channel:        latestRelease.Channel,
		Artifacts:      []ReleaseArtifact{},
	}

	// Check if update is critical
	if meta := latestRelease.Metadata; meta != nil {
		if critical, ok := meta["critical"].(bool); ok {
			updateInfo.Critical = critical
		}
	}

	// Get artifacts for current platform
	artifacts, err := latestRelease.Artifacts(keygen.ArtifactListOptions{})
	if err == nil {
		for _, artifact := range artifacts {
			// Filter by platform
			if artifact.Platform == runtime.GOOS && artifact.Arch == runtime.GOARCH {
				updateInfo.Artifacts = append(updateInfo.Artifacts, ReleaseArtifact{
					ID:       artifact.ID,
					Platform: artifact.Platform,
					Arch:     artifact.Arch,
					Filename: artifact.Filename,
					Size:     artifact.Filesize,
					Checksum: artifact.Checksum,
				})

				if updateInfo.Size == 0 {
					updateInfo.Size = artifact.Filesize
				}
			}
		}
	}

	// Emit update available event
	if s.eventEmitter != nil {
		s.eventEmitter.EmitUpdateAvailable(UpdateAvailableEvent{
			Version:     latestRelease.Version,
			ReleaseDate: latestRelease.PublishedAt.Format(time.RFC3339),
			Notes:       latestRelease.Description,
			Mandatory:   updateInfo.Critical,
			DownloadURL: "", // Set by download process
			Size:        updateInfo.Size,
		})
	}

	return updateInfo, nil
}

// DownloadUpdate downloads an available update
func (s *Service) DownloadUpdate(releaseID string) (*DownloadProgress, error) {
	s.mu.Lock()
	if s.updateInProgress {
		s.mu.Unlock()
		return nil, &KeygenError{Code: ErrUpdateInProgress, Message: "Update download already in progress"}
	}
	s.updateInProgress = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.updateInProgress = false
		s.downloadCancel = nil
		s.mu.Unlock()
	}()

	// Create download context
	ctx, cancel := context.WithCancel(context.Background())
	s.downloadCancel = cancel

	// Get release
	release, err := s.client.Release(s.config.Product, releaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get release: %w", err)
	}

	// Find artifact for current platform
	artifacts, err := release.Artifacts(keygen.ArtifactListOptions{
		Platform: runtime.GOOS,
		Arch:     runtime.GOARCH,
	})

	if err != nil || len(artifacts) == 0 {
		return nil, fmt.Errorf("no artifact available for platform %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	artifact := artifacts[0]

	// Create download directory
	downloadDir := filepath.Join(s.cacheDir, "downloads")
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create download directory: %w", err)
	}

	// Determine download path
	downloadPath := filepath.Join(downloadDir, artifact.Filename)

	// Create progress tracker
	progress := &DownloadProgress{
		ID:         artifact.ID,
		State:      DownloadStateDownloading,
		TotalBytes: artifact.Filesize,
		FilePath:   downloadPath,
	}

	// Start download in background
	go s.performDownload(ctx, artifact, downloadPath, progress)

	return progress, nil
}

// InstallUpdate installs a downloaded update
func (s *Service) InstallUpdate() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Find the latest downloaded update
	downloadDir := filepath.Join(s.cacheDir, "downloads")
	entries, err := os.ReadDir(downloadDir)
	if err != nil {
		return fmt.Errorf("failed to read download directory: %w", err)
	}

	if len(entries) == 0 {
		return &KeygenError{Code: ErrUpdateNotAvailable, Message: "No update available to install"}
	}

	// Get the most recent download
	var latestFile string
	var latestTime time.Time

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().After(latestTime) {
			latestTime = info.ModTime()
			latestFile = filepath.Join(downloadDir, entry.Name())
		}
	}

	if latestFile == "" {
		return &KeygenError{Code: ErrUpdateNotAvailable, Message: "No update file found"}
	}

	// Delegate to platform-specific installation
	if err := s.impl.InstallUpdatePlatform(latestFile); err != nil {
		return &KeygenError{
			Code:    ErrUpdateInstallFailed,
			Message: "Failed to install update",
			Err:     err,
		}
	}

	// Clean up download
	_ = os.Remove(latestFile)

	// Emit update installed event
	if s.eventEmitter != nil {
		// Note: We should emit an UpdateInstalledEvent type, but events.go doesn't define it
		// For now, we'll emit a generic event
		if s.app != nil && s.app.Event != nil {
			s.app.Event.Emit("keygen:update-installed", map[string]interface{}{
				"version":     latestFile,
				"installedAt": time.Now().Format(time.RFC3339),
			})
		}
	}

	return nil
}

// GetCurrentVersion returns the current application version
func (s *Service) GetCurrentVersion() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentVersion
}

// SetUpdateChannel sets the update channel
func (s *Service) SetUpdateChannel(channel string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	validChannels := []string{"stable", "beta", "alpha", "dev"}
	valid := false
	for _, c := range validChannels {
		if c == channel {
			valid = true
			break
		}
	}

	if !valid {
		return &KeygenError{
			Code:    ErrUpdateChannelInvalid,
			Message: fmt.Sprintf("Invalid update channel: %s", channel),
		}
	}

	s.config.UpdateChannel = channel
	return nil
}

// SaveOfflineLicense saves the current license for offline use
func (s *Service) SaveOfflineLicense() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.state.License == nil {
		return NewLicenseInvalidError("no license to save")
	}

	// Create offline data
	offlineData := map[string]interface{}{
		"license":      s.state.License,
		"lastChecked":  s.state.LastChecked,
		"expiresAt":    s.state.ExpiresAt,
		"entitlements": s.state.Entitlements,
	}

	// Marshal to JSON
	data, err := json.Marshal(offlineData)
	if err != nil {
		return fmt.Errorf("failed to marshal offline license: %w", err)
	}

	// Save to cache
	licensePath := filepath.Join(s.cacheDir, "license.json")
	if err := os.WriteFile(licensePath, data, 0600); err != nil {
		return &KeygenError{
			Code:    ErrCacheWriteFailed,
			Message: "Failed to save offline license",
			Err:     err,
		}
	}

	return nil
}

// LoadOfflineLicense loads a previously saved offline license
func (s *Service) LoadOfflineLicense() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	licensePath := filepath.Join(s.cacheDir, "license.json")
	data, err := os.ReadFile(licensePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Not an error if no offline license exists
		}
		return &KeygenError{
			Code:    ErrCacheReadFailed,
			Message: "Failed to load offline license",
			Err:     err,
		}
	}

	// Unmarshal offline data
	var offlineData map[string]interface{}
	if err := json.Unmarshal(data, &offlineData); err != nil {
		return &KeygenError{
			Code:    ErrCacheCorrupted,
			Message: "Offline license data is corrupted",
			Err:     err,
		}
	}

	// Restore state
	s.state.OfflineMode = true

	if lastChecked, ok := offlineData["lastChecked"].(string); ok {
		if t, err := time.Parse(time.RFC3339, lastChecked); err == nil {
			s.state.LastChecked = t
		}
	}

	if entitlements, ok := offlineData["entitlements"].(map[string]interface{}); ok {
		s.state.Entitlements = entitlements
	}

	return nil
}

// ClearLicenseCache clears all cached license data
func (s *Service) ClearLicenseCache() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Clear in-memory state
	s.state = &LicenseState{}

	// Remove cache files
	licensePath := filepath.Join(s.cacheDir, "license.json")
	if err := os.Remove(licensePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove license cache: %w", err)
	}

	return nil
}

// SetApplication sets the application reference for event emission
// This would typically be called by the application during service initialization
func (s *Service) SetApplication(app *application.App) {
	s.app = app
	s.eventEmitter = NewEventEmitter(app)
}

// Helper methods

// startAutoUpdateCheck starts the automatic update checking routine
func (s *Service) startAutoUpdateCheck(ctx context.Context) {
	ticker := time.NewTicker(s.config.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, _ = s.CheckForUpdates()
		}
	}
}

// performDownload performs the actual download
func (s *Service) performDownload(ctx context.Context, artifact *keygen.Artifact, downloadPath string, progress *DownloadProgress) {
	// Create temp file
	tempPath := downloadPath + ".tmp"
	tempFile, err := os.Create(tempPath)
	if err != nil {
		progress.State = DownloadStateFailed
		progress.Error = err.Error()
		s.emitDownloadProgress(progress)
		return
	}
	defer tempFile.Close()
	defer os.Remove(tempPath)

	// Get download URL
	downloadURL, err := artifact.URL()
	if err != nil {
		progress.State = DownloadStateFailed
		progress.Error = err.Error()
		s.emitDownloadProgress(progress)
		return
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", downloadURL, nil)
	if err != nil {
		progress.State = DownloadStateFailed
		progress.Error = err.Error()
		s.emitDownloadProgress(progress)
		return
	}

	// Perform request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		progress.State = DownloadStateFailed
		progress.Error = err.Error()
		s.emitDownloadProgress(progress)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		progress.State = DownloadStateFailed
		progress.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
		s.emitDownloadProgress(progress)
		return
	}

	// Create progress reader
	reader := &progressReader{
		reader:   resp.Body,
		progress: progress,
		onProgress: func(p *DownloadProgress) {
			s.emitDownloadProgress(p)
		},
	}

	// Download with progress
	startTime := time.Now()
	hasher := sha256.New()
	writer := io.MultiWriter(tempFile, hasher)

	_, err = io.Copy(writer, reader)
	if err != nil {
		if ctx.Err() != nil {
			progress.State = DownloadStateCancelled
		} else {
			progress.State = DownloadStateFailed
			progress.Error = err.Error()
		}
		s.emitDownloadProgress(progress)
		return
	}

	// Verify checksum
	checksum := hex.EncodeToString(hasher.Sum(nil))
	if artifact.Checksum != "" && checksum != artifact.Checksum {
		progress.State = DownloadStateFailed
		progress.Error = "Checksum mismatch"
		s.emitDownloadProgress(progress)
		return
	}

	// Move temp file to final location
	if err := os.Rename(tempPath, downloadPath); err != nil {
		progress.State = DownloadStateFailed
		progress.Error = err.Error()
		s.emitDownloadProgress(progress)
		return
	}

	// Update final progress
	progress.State = DownloadStateCompleted
	progress.BytesDownloaded = progress.TotalBytes
	progress.Percentage = 100.0
	progress.Speed = 0
	progress.ETA = 0

	// Calculate average speed
	duration := time.Since(startTime).Seconds()
	if duration > 0 {
		avgSpeed := float64(progress.TotalBytes) / duration
		progress.Speed = int64(avgSpeed)
	}

	s.emitDownloadProgress(progress)
}

// emitLicenseStatus emits a license status event
func (s *Service) emitLicenseStatus(valid bool, message, error string) {
	if s.eventEmitter == nil {
		return
	}

	event := LicenseStatusEvent{
		Valid:       valid,
		Status:      "invalid",
		Key:         s.config.LicenseKey,
		LastChecked: time.Now().Format(time.RFC3339),
		Error:       error,
	}

	if valid {
		event.Status = "active"
	}

	if s.state.ExpiresAt != nil {
		event.ExpiresAt = s.state.ExpiresAt.Format(time.RFC3339)
	}

	if s.state.License != nil && s.state.License.Metadata != nil {
		if email, ok := s.state.License.Metadata["email"].(string); ok {
			event.Email = email
		}
	}

	s.eventEmitter.EmitLicenseStatus(event)
}

// emitDownloadProgress emits a download progress event
func (s *Service) emitDownloadProgress(progress *DownloadProgress) {
	if s.eventEmitter == nil {
		return
	}

	event := DownloadProgressEvent{
		BytesDownloaded: progress.BytesDownloaded,
		TotalBytes:      progress.TotalBytes,
		Progress:        progress.Percentage,
		Speed:           progress.Speed,
		TimeRemaining:   int(progress.ETA),
		Status:          progress.State,
		Error:           progress.Error,
	}

	s.eventEmitter.EmitDownloadProgress(event)
}

// progressReader wraps an io.Reader to track download progress
type progressReader struct {
	reader            io.Reader
	progress          *DownloadProgress
	onProgress        func(*DownloadProgress)
	lastUpdate        time.Time
	bytesAtLastUpdate int64
}

func (r *progressReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)

	if n > 0 {
		r.progress.BytesDownloaded += int64(n)
		r.progress.Percentage = float64(r.progress.BytesDownloaded) / float64(r.progress.TotalBytes) * 100

		// Calculate speed and ETA
		now := time.Now()
		if now.Sub(r.lastUpdate) >= time.Second {
			duration := now.Sub(r.lastUpdate).Seconds()
			bytesThisPeriod := r.progress.BytesDownloaded - r.bytesAtLastUpdate
			r.progress.Speed = int64(float64(bytesThisPeriod) / duration)

			if r.progress.Speed > 0 {
				remaining := r.progress.TotalBytes - r.progress.BytesDownloaded
				r.progress.ETA = remaining / r.progress.Speed
			}

			r.lastUpdate = now
			r.bytesAtLastUpdate = r.progress.BytesDownloaded
		}

		// Emit progress update
		if r.onProgress != nil {
			r.onProgress(r.progress)
		}
	}

	return n, err
}

// compareVersions compares two semantic versions
// Returns: 1 if v1 > v2, -1 if v1 < v2, 0 if equal
func compareVersions(v1, v2 string) int {
	// Simple semantic version comparison
	// Strip "v" prefix if present
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	// Pad shorter version with zeros
	for len(parts1) < len(parts2) {
		parts1 = append(parts1, "0")
	}
	for len(parts2) < len(parts1) {
		parts2 = append(parts2, "0")
	}

	// Compare each part
	for i := 0; i < len(parts1); i++ {
		n1, _ := strconv.Atoi(parts1[i])
		n2, _ := strconv.Atoi(parts2[i])

		if n1 > n2 {
			return 1
		} else if n1 < n2 {
			return -1
		}
	}

	return 0
}
