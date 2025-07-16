package keygen

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Event names for the keygen service
const (
	EventUpdateAvailable  = "keygen:update-available"
	EventLicenseStatus    = "keygen:license-status"
	EventDownloadProgress = "keygen:download-progress"
	EventUpdateInstalled  = "keygen:update-installed"
)

// UpdateAvailableEvent is emitted when a new update is available
type UpdateAvailableEvent struct {
	Version     string `json:"version"`
	ReleaseDate string `json:"releaseDate"`
	Notes       string `json:"notes"`
	Mandatory   bool   `json:"mandatory"`
	DownloadURL string `json:"downloadUrl"`
	Size        int64  `json:"size"`
}

// LicenseStatusEvent is emitted when the license status changes
type LicenseStatusEvent struct {
	Valid       bool   `json:"valid"`
	Status      string `json:"status"`
	Key         string `json:"key"`
	Email       string `json:"email"`
	ExpiresAt   string `json:"expiresAt,omitempty"`
	EnteredAt   string `json:"enteredAt,omitempty"`
	LastChecked string `json:"lastChecked,omitempty"`
	Error       string `json:"error,omitempty"`
}

// DownloadProgressEvent is emitted during update downloads
type DownloadProgressEvent struct {
	BytesDownloaded int64   `json:"bytesDownloaded"`
	TotalBytes      int64   `json:"totalBytes"`
	Progress        float64 `json:"progress"`
	Speed           int64   `json:"speed"`         // bytes per second
	TimeRemaining   int     `json:"timeRemaining"` // seconds
	Status          string  `json:"status"`
	Error           string  `json:"error,omitempty"`
}

// EventEmitter provides methods to emit events from the keygen service
type EventEmitter struct {
	app *application.App
}

// NewEventEmitter creates a new event emitter
// Note: This requires the application instance to be available globally or passed through context
func NewEventEmitter(app *application.App) *EventEmitter {
	return &EventEmitter{
		app: app,
	}
}

// EmitUpdateAvailable emits an update available event
func (e *EventEmitter) EmitUpdateAvailable(event UpdateAvailableEvent) {
	if e.app != nil && e.app.Event != nil {
		e.app.Event.Emit(EventUpdateAvailable, event)
	}
}

// EmitLicenseStatus emits a license status event
func (e *EventEmitter) EmitLicenseStatus(event LicenseStatusEvent) {
	if e.app != nil && e.app.Event != nil {
		e.app.Event.Emit(EventLicenseStatus, event)
	}
}

// EmitDownloadProgress emits a download progress event
func (e *EventEmitter) EmitDownloadProgress(event DownloadProgressEvent) {
	if e.app != nil && e.app.Event != nil {
		e.app.Event.Emit(EventDownloadProgress, event)
	}
}

// EmitUpdateInstalled emits an update installed event
func (e *EventEmitter) EmitUpdateInstalled(event UpdateInstalledEvent) {
	if e.app != nil && e.app.Event != nil {
		e.app.Event.Emit(EventUpdateInstalled, event)
	}
}

// Service-level event emission helpers
// These can be used when the service has access to the application instance

// EmitUpdateAvailableFromService is a helper to emit update available events from the service
func EmitUpdateAvailableFromService(app *application.App, event UpdateAvailableEvent) {
	if app != nil && app.Event != nil {
		app.Event.Emit(EventUpdateAvailable, event)
	}
}

// EmitLicenseStatusFromService is a helper to emit license status events from the service
func EmitLicenseStatusFromService(app *application.App, event LicenseStatusEvent) {
	if app != nil && app.Event != nil {
		app.Event.Emit(EventLicenseStatus, event)
	}
}

// EmitDownloadProgressFromService is a helper to emit download progress events from the service
func EmitDownloadProgressFromService(app *application.App, event DownloadProgressEvent) {
	if app != nil && app.Event != nil {
		app.Event.Emit(EventDownloadProgress, event)
	}
}

// Alternative approach: Use window-based event emission if the service operates within a window context

// EmitUpdateAvailableFromWindow emits an update available event from a window
func EmitUpdateAvailableFromWindow(window application.Window, event UpdateAvailableEvent) {
	if window != nil {
		window.EmitEvent(EventUpdateAvailable, event)
	}
}

// EmitLicenseStatusFromWindow emits a license status event from a window
func EmitLicenseStatusFromWindow(window application.Window, event LicenseStatusEvent) {
	if window != nil {
		window.EmitEvent(EventLicenseStatus, event)
	}
}

// EmitDownloadProgressFromWindow emits a download progress event from a window
func EmitDownloadProgressFromWindow(window application.Window, event DownloadProgressEvent) {
	if window != nil {
		window.EmitEvent(EventDownloadProgress, event)
	}
}

// GetWindowFromContext retrieves the window from a context (used in service methods)
// Note: This uses the internal window context key "Window"
func GetWindowFromContext(ctx context.Context) application.Window {
	type contextKey string
	const windowKey contextKey = "Window"

	if window, ok := ctx.Value(windowKey).(application.Window); ok {
		return window
	}
	return nil
}
