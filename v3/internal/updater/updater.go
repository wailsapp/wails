package updater

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Updater manages the complete update lifecycle
type Updater struct {
	config         Config
	currentVersion string
	checker        *Checker
	downloader     *Downloader

	// State management
	mu             sync.RWMutex
	state          UpdateState
	currentInfo    *UpdateInfo
	downloadedPath string
	lastError      error

	// Event handling
	eventHandlers []func(UpdateEvent)

	// Background checking
	stopChan chan struct{}
	stopOnce sync.Once
}

// New creates a new Updater instance
func New(config Config, currentVersion string) *Updater {
	return &Updater{
		config:         config,
		currentVersion: currentVersion,
		checker:        NewChecker(config),
		downloader:     NewDownloader(),
		state:          StateIdle,
		stopChan:       make(chan struct{}),
	}
}

// GetCurrentVersion returns the current application version
func (u *Updater) GetCurrentVersion() string {
	return u.currentVersion
}

// GetState returns the current update state
func (u *Updater) GetState() UpdateState {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.state
}

// GetUpdateInfo returns information about the available update (if any)
func (u *Updater) GetUpdateInfo() *UpdateInfo {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.currentInfo
}

// GetLastError returns the last error that occurred
func (u *Updater) GetLastError() error {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.lastError
}

// OnEvent registers an event handler
func (u *Updater) OnEvent(handler func(UpdateEvent)) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.eventHandlers = append(u.eventHandlers, handler)
}

// CheckForUpdate checks if a new version is available
func (u *Updater) CheckForUpdate(ctx context.Context) (*UpdateInfo, error) {
	u.setState(StateChecking, nil, nil, nil)

	info, err := u.checker.CheckForUpdate(ctx, u.currentVersion)
	if err != nil {
		u.setError(err)
		return nil, err
	}

	if info == nil {
		u.setState(StateIdle, nil, nil, nil)
		return nil, nil
	}

	u.mu.Lock()
	u.currentInfo = info
	u.mu.Unlock()

	u.setState(StateAvailable, info, nil, nil)
	return info, nil
}

// DownloadUpdate downloads the available update
func (u *Updater) DownloadUpdate(ctx context.Context, progress ProgressCallback) error {
	u.mu.RLock()
	info := u.currentInfo
	u.mu.RUnlock()

	if info == nil {
		return fmt.Errorf("no update available")
	}

	u.setState(StateDownloading, info, nil, nil)

	// Create a wrapper progress callback that emits events
	wrappedProgress := func(p DownloadProgress) {
		u.emitEvent(UpdateEvent{
			State:    StateDownloading,
			Info:     info,
			Progress: &p,
		})
		if progress != nil {
			progress(p)
		}
	}

	path, err := u.downloader.Download(ctx, info, wrappedProgress)
	if err != nil {
		u.setError(err)
		return err
	}

	u.mu.Lock()
	u.downloadedPath = path
	u.mu.Unlock()

	u.setState(StateReady, info, nil, nil)
	return nil
}

// ApplyUpdate applies the downloaded update
// This will typically restart the application
func (u *Updater) ApplyUpdate(ctx context.Context) error {
	u.mu.RLock()
	path := u.downloadedPath
	info := u.currentInfo
	u.mu.RUnlock()

	if path == "" {
		return fmt.Errorf("no update downloaded")
	}

	u.setState(StateInstalling, info, nil, nil)

	// Platform-specific installation
	if err := applyUpdate(ctx, path, info); err != nil {
		u.setError(err)
		return err
	}

	// If we get here, installation succeeded but app should restart
	return nil
}

// DownloadAndApply downloads and applies an update in one call
func (u *Updater) DownloadAndApply(ctx context.Context, progress ProgressCallback) error {
	if err := u.DownloadUpdate(ctx, progress); err != nil {
		return err
	}
	return u.ApplyUpdate(ctx)
}

// StartBackgroundChecks starts periodic background update checks
func (u *Updater) StartBackgroundChecks(ctx context.Context) {
	if u.config.CheckInterval <= 0 {
		return
	}

	go func() {
		ticker := time.NewTicker(u.config.CheckInterval)
		defer ticker.Stop()

		// Do an initial check
		u.CheckForUpdate(ctx)

		for {
			select {
			case <-ctx.Done():
				return
			case <-u.stopChan:
				return
			case <-ticker.C:
				u.CheckForUpdate(ctx)
			}
		}
	}()
}

// StopBackgroundChecks stops periodic background update checks
func (u *Updater) StopBackgroundChecks() {
	u.stopOnce.Do(func() {
		close(u.stopChan)
	})
}

// setState updates the current state and emits an event
func (u *Updater) setState(state UpdateState, info *UpdateInfo, progress *DownloadProgress, err error) {
	u.mu.Lock()
	u.state = state
	if err != nil {
		u.lastError = err
	}
	u.mu.Unlock()

	event := UpdateEvent{
		State:    state,
		Info:     info,
		Progress: progress,
	}
	if err != nil {
		event.Error = err.Error()
	}
	u.emitEvent(event)
}

// setError sets the error state
func (u *Updater) setError(err error) {
	u.setState(StateError, u.currentInfo, nil, err)
}

// emitEvent sends an event to all registered handlers
func (u *Updater) emitEvent(event UpdateEvent) {
	u.mu.RLock()
	handlers := make([]func(UpdateEvent), len(u.eventHandlers))
	copy(handlers, u.eventHandlers)
	u.mu.RUnlock()

	for _, handler := range handlers {
		handler(event)
	}
}

// Reset resets the updater state
func (u *Updater) Reset() {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.state = StateIdle
	u.currentInfo = nil
	u.downloadedPath = ""
	u.lastError = nil
}

// CleanupDownloads removes old downloaded update files
func (u *Updater) CleanupDownloads() error {
	return u.downloader.CleanupOldDownloads()
}
