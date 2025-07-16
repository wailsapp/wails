package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/keygen"
)

// App struct
type App struct {
	window application.Window
	keygen *keygen.Service
}

// SetLicenseKey sets the license key for validation
func (a *App) SetLicenseKey(key string) error {
	a.keygen = keygen.New(keygen.ServiceOptions{
		AccountID:      "demo",
		ProductID:      "prod_demo",
		LicenseKey:     key,
		PublicKey:      "",
		CurrentVersion: "1.0.0",
		AutoCheck:      true,
		UpdateChannel:  "stable",
	})

	// Initialize the service
	ctx := context.Background()
	err := a.keygen.ServiceStartup(ctx, application.ServiceOptions{
		App: application.Get(),
	})
	if err != nil {
		return fmt.Errorf("failed to initialize keygen service: %w", err)
	}

	return nil
}

// ValidateLicense validates the current license
func (a *App) ValidateLicense() (*keygen.LicenseValidationResult, error) {
	if a.keygen == nil {
		return nil, fmt.Errorf("keygen service not initialized")
	}

	result, err := a.keygen.ValidateLicense()
	if err != nil {
		slog.Error("License validation failed", "error", err)
		return result, err
	}

	slog.Info("License validation completed", "valid", result.Valid, "message", result.Message)
	return result, nil
}

// ActivateMachine activates the current machine for the license
func (a *App) ActivateMachine() (*keygen.MachineActivationResult, error) {
	if a.keygen == nil {
		return nil, fmt.Errorf("keygen service not initialized")
	}

	result, err := a.keygen.ActivateMachine()
	if err != nil {
		slog.Error("Machine activation failed", "error", err)
		return result, err
	}

	slog.Info("Machine activation completed", "success", result.Success, "fingerprint", result.Fingerprint)
	return result, nil
}

// DeactivateMachine deactivates the current machine
func (a *App) DeactivateMachine() error {
	if a.keygen == nil {
		return fmt.Errorf("keygen service not initialized")
	}

	err := a.keygen.DeactivateMachine()
	if err != nil {
		slog.Error("Machine deactivation failed", "error", err)
		return err
	}

	slog.Info("Machine deactivated successfully")
	return nil
}

// GetLicenseInfo returns detailed license information
func (a *App) GetLicenseInfo() (*keygen.LicenseInfo, error) {
	if a.keygen == nil {
		return nil, fmt.Errorf("keygen service not initialized")
	}

	info, err := a.keygen.GetLicenseInfo()
	if err != nil {
		slog.Error("Failed to get license info", "error", err)
		return nil, err
	}

	return info, nil
}

// CheckEntitlement checks if a specific feature is enabled
func (a *App) CheckEntitlement(feature string) (bool, error) {
	if a.keygen == nil {
		return false, fmt.Errorf("keygen service not initialized")
	}

	enabled, err := a.keygen.CheckEntitlement(feature)
	if err != nil {
		slog.Error("Failed to check entitlement", "feature", feature, "error", err)
		return false, err
	}

	slog.Info("Entitlement check", "feature", feature, "enabled", enabled)
	return enabled, nil
}

// CheckForUpdates checks for available updates
func (a *App) CheckForUpdates() (*keygen.UpdateInfo, error) {
	if a.keygen == nil {
		return nil, fmt.Errorf("keygen service not initialized")
	}

	updateInfo, err := a.keygen.CheckForUpdates()
	if err != nil {
		slog.Error("Failed to check for updates", "error", err)
		return nil, err
	}

	slog.Info("Update check completed",
		"available", updateInfo.Available,
		"current", updateInfo.CurrentVersion,
		"latest", updateInfo.LatestVersion)

	return updateInfo, nil
}

// DownloadUpdate downloads an available update
func (a *App) DownloadUpdate(releaseID string) (*keygen.DownloadProgress, error) {
	if a.keygen == nil {
		return nil, fmt.Errorf("keygen service not initialized")
	}

	progress, err := a.keygen.DownloadUpdate(releaseID)
	if err != nil {
		slog.Error("Failed to start update download", "error", err)
		return nil, err
	}

	slog.Info("Update download started", "releaseID", releaseID)
	return progress, nil
}

// InstallUpdate installs a downloaded update
func (a *App) InstallUpdate() error {
	if a.keygen == nil {
		return fmt.Errorf("keygen service not initialized")
	}

	err := a.keygen.InstallUpdate()
	if err != nil {
		slog.Error("Failed to install update", "error", err)
		return err
	}

	slog.Info("Update installed successfully")
	return nil
}

// SetUpdateChannel sets the update channel (stable, beta, alpha, dev)
func (a *App) SetUpdateChannel(channel string) error {
	if a.keygen == nil {
		return fmt.Errorf("keygen service not initialized")
	}

	err := a.keygen.SetUpdateChannel(channel)
	if err != nil {
		slog.Error("Failed to set update channel", "channel", channel, "error", err)
		return err
	}

	slog.Info("Update channel changed", "channel", channel)
	return nil
}

// GetCurrentVersion returns the current application version
func (a *App) GetCurrentVersion() string {
	if a.keygen == nil {
		return "1.0.0"
	}

	return a.keygen.GetCurrentVersion()
}

// SaveOfflineLicense saves the license for offline use
func (a *App) SaveOfflineLicense() error {
	if a.keygen == nil {
		return fmt.Errorf("keygen service not initialized")
	}

	err := a.keygen.SaveOfflineLicense()
	if err != nil {
		slog.Error("Failed to save offline license", "error", err)
		return err
	}

	slog.Info("Offline license saved successfully")
	return nil
}

// LoadOfflineLicense loads a previously saved offline license
func (a *App) LoadOfflineLicense() error {
	if a.keygen == nil {
		return fmt.Errorf("keygen service not initialized")
	}

	err := a.keygen.LoadOfflineLicense()
	if err != nil {
		slog.Error("Failed to load offline license", "error", err)
		return err
	}

	slog.Info("Offline license loaded successfully")
	return nil
}

// ClearLicenseCache clears all cached license data
func (a *App) ClearLicenseCache() error {
	if a.keygen == nil {
		return fmt.Errorf("keygen service not initialized")
	}

	err := a.keygen.ClearLicenseCache()
	if err != nil {
		slog.Error("Failed to clear license cache", "error", err)
		return err
	}

	slog.Info("License cache cleared successfully")
	return nil
}
