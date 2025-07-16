package keygen

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/keygen-sh/keygen-go/v3"
	"github.com/matryer/is"
)

func TestConfigMarshaling(t *testing.T) {
	i := is.New(t)

	config := Config{
		Account:       "test-account",
		Product:       "test-product",
		LicenseKey:    "test-key",
		PublicKey:     "test-public-key",
		CacheDir:      "/tmp/cache",
		AutoCheck:     true,
		CheckInterval: 24 * time.Hour,
		UpdateChannel: "stable",
		Environment:   "production",
	}

	// Marshal to JSON
	data, err := json.Marshal(config)
	i.NoErr(err)

	// Unmarshal back
	var config2 Config
	err = json.Unmarshal(data, &config2)
	i.NoErr(err)

	// Compare
	i.Equal(config.Account, config2.Account)
	i.Equal(config.Product, config2.Product)
	i.Equal(config.LicenseKey, config2.LicenseKey)
	i.Equal(config.PublicKey, config2.PublicKey)
	i.Equal(config.CacheDir, config2.CacheDir)
	i.Equal(config.AutoCheck, config2.AutoCheck)
	i.Equal(config.CheckInterval, config2.CheckInterval)
	i.Equal(config.UpdateChannel, config2.UpdateChannel)
	i.Equal(config.Environment, config2.Environment)
}

func TestLicenseStateMarshaling(t *testing.T) {
	i := is.New(t)

	now := time.Now()
	expiresAt := now.Add(30 * 24 * time.Hour)

	state := LicenseState{
		Valid: true,
		License: &keygen.License{
			ID:   "license-123",
			Key:  "test-key",
			Name: "Test License",
		},
		LastChecked: now,
		OfflineMode: false,
		Entitlements: map[string]interface{}{
			"feature1": true,
			"feature2": "enabled",
			"feature3": map[string]interface{}{
				"limit": 100,
			},
		},
		ExpiresAt: &expiresAt,
		Machine: &keygen.Machine{
			ID:          "machine-123",
			Fingerprint: "fingerprint-123",
			Name:        "Test Machine",
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(state)
	i.NoErr(err)

	// Unmarshal back
	var state2 LicenseState
	err = json.Unmarshal(data, &state2)
	i.NoErr(err)

	// Compare basic fields
	i.Equal(state.Valid, state2.Valid)
	i.Equal(state.OfflineMode, state2.OfflineMode)
	i.Equal(state.LastChecked.Unix(), state2.LastChecked.Unix())
	i.True(state2.License != nil)
	i.True(state2.Machine != nil)
	i.True(state2.ExpiresAt != nil)
	i.Equal(state.ExpiresAt.Unix(), state2.ExpiresAt.Unix())

	// Check entitlements
	i.Equal(state2.Entitlements["feature1"], true)
	i.Equal(state2.Entitlements["feature2"], "enabled")
	feature3, ok := state2.Entitlements["feature3"].(map[string]interface{})
	i.True(ok)
	i.Equal(feature3["limit"], float64(100)) // JSON numbers unmarshal as float64
}

func TestLicenseValidationResultMarshaling(t *testing.T) {
	i := is.New(t)

	result := LicenseValidationResult{
		Valid:   true,
		Message: "License is valid",
		Code:    ValidationCodeValid,
		State: &LicenseState{
			Valid:       true,
			LastChecked: time.Now(),
		},
		RequiresActivation: true,
	}

	// Marshal to JSON
	data, err := json.Marshal(result)
	i.NoErr(err)

	// Unmarshal back
	var result2 LicenseValidationResult
	err = json.Unmarshal(data, &result2)
	i.NoErr(err)

	// Compare
	i.Equal(result.Valid, result2.Valid)
	i.Equal(result.Message, result2.Message)
	i.Equal(result.Code, result2.Code)
	i.Equal(result.RequiresActivation, result2.RequiresActivation)
	i.True(result2.State != nil)
	i.Equal(result.State.Valid, result2.State.Valid)
}

func TestMachineActivationResultMarshaling(t *testing.T) {
	i := is.New(t)

	result := MachineActivationResult{
		Success:     true,
		Message:     "Machine activated successfully",
		Machine:     &keygen.Machine{ID: "machine-123"},
		Fingerprint: "fingerprint-123",
	}

	// Marshal to JSON
	data, err := json.Marshal(result)
	i.NoErr(err)

	// Unmarshal back
	var result2 MachineActivationResult
	err = json.Unmarshal(data, &result2)
	i.NoErr(err)

	// Compare
	i.Equal(result.Success, result2.Success)
	i.Equal(result.Message, result2.Message)
	i.Equal(result.Fingerprint, result2.Fingerprint)
	i.True(result2.Machine != nil)
}

func TestLicenseInfoMarshaling(t *testing.T) {
	i := is.New(t)

	now := time.Now()
	expiresAt := now.Add(30 * 24 * time.Hour)
	maxMachines := 5

	info := LicenseInfo{
		Key:             "license-key",
		Name:            "John Doe",
		Email:           "john@example.com",
		Company:         "ACME Corp",
		Status:          "active",
		ExpiresAt:       &expiresAt,
		CreatedAt:       now,
		LastValidatedAt: &now,
		Entitlements: map[string]interface{}{
			"feature1": true,
		},
		Metadata: map[string]interface{}{
			"custom": "data",
		},
		MaxMachines:  &maxMachines,
		MachineCount: 3,
	}

	// Marshal to JSON
	data, err := json.Marshal(info)
	i.NoErr(err)

	// Unmarshal back
	var info2 LicenseInfo
	err = json.Unmarshal(data, &info2)
	i.NoErr(err)

	// Compare
	i.Equal(info.Key, info2.Key)
	i.Equal(info.Name, info2.Name)
	i.Equal(info.Email, info2.Email)
	i.Equal(info.Company, info2.Company)
	i.Equal(info.Status, info2.Status)
	i.True(info2.ExpiresAt != nil)
	i.True(info2.LastValidatedAt != nil)
	i.Equal(*info.MaxMachines, *info2.MaxMachines)
	i.Equal(info.MachineCount, info2.MachineCount)
}

func TestUpdateInfoMarshaling(t *testing.T) {
	i := is.New(t)

	publishedAt := time.Now()

	info := UpdateInfo{
		Available:      true,
		CurrentVersion: "1.0.0",
		LatestVersion:  "1.1.0",
		ReleaseID:      "release-123",
		ReleaseName:    "Version 1.1.0",
		ReleaseNotes:   "Bug fixes and improvements",
		PublishedAt:    &publishedAt,
		Critical:       true,
		Size:           1024 * 1024 * 50, // 50MB
		Channel:        "stable",
		Artifacts: []ReleaseArtifact{
			{
				ID:           "artifact-123",
				Platform:     "darwin",
				Arch:         "amd64",
				Filename:     "app-darwin-amd64.zip",
				Size:         1024 * 1024 * 50,
				Checksum:     "sha256:abcdef123456",
				SignatureURL: "https://example.com/sig",
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(info)
	i.NoErr(err)

	// Unmarshal back
	var info2 UpdateInfo
	err = json.Unmarshal(data, &info2)
	i.NoErr(err)

	// Compare
	i.Equal(info.Available, info2.Available)
	i.Equal(info.CurrentVersion, info2.CurrentVersion)
	i.Equal(info.LatestVersion, info2.LatestVersion)
	i.Equal(info.ReleaseID, info2.ReleaseID)
	i.Equal(info.ReleaseName, info2.ReleaseName)
	i.Equal(info.ReleaseNotes, info2.ReleaseNotes)
	i.Equal(info.Critical, info2.Critical)
	i.Equal(info.Size, info2.Size)
	i.Equal(info.Channel, info2.Channel)
	i.Equal(len(info.Artifacts), len(info2.Artifacts))

	if len(info2.Artifacts) > 0 {
		i.Equal(info.Artifacts[0].ID, info2.Artifacts[0].ID)
		i.Equal(info.Artifacts[0].Platform, info2.Artifacts[0].Platform)
		i.Equal(info.Artifacts[0].Arch, info2.Artifacts[0].Arch)
		i.Equal(info.Artifacts[0].Filename, info2.Artifacts[0].Filename)
		i.Equal(info.Artifacts[0].Size, info2.Artifacts[0].Size)
		i.Equal(info.Artifacts[0].Checksum, info2.Artifacts[0].Checksum)
		i.Equal(info.Artifacts[0].SignatureURL, info2.Artifacts[0].SignatureURL)
	}
}

func TestDownloadProgressMarshaling(t *testing.T) {
	i := is.New(t)

	progress := DownloadProgress{
		ID:              "download-123",
		State:           DownloadStateDownloading,
		BytesDownloaded: 1024 * 512,  // 512KB
		TotalBytes:      1024 * 1024, // 1MB
		Percentage:      50.0,
		Speed:           1024 * 100, // 100KB/s
		ETA:             5,          // 5 seconds
		Error:           "",
		FilePath:        "/tmp/download.zip",
	}

	// Marshal to JSON
	data, err := json.Marshal(progress)
	i.NoErr(err)

	// Unmarshal back
	var progress2 DownloadProgress
	err = json.Unmarshal(data, &progress2)
	i.NoErr(err)

	// Compare
	i.Equal(progress.ID, progress2.ID)
	i.Equal(progress.State, progress2.State)
	i.Equal(progress.BytesDownloaded, progress2.BytesDownloaded)
	i.Equal(progress.TotalBytes, progress2.TotalBytes)
	i.Equal(progress.Percentage, progress2.Percentage)
	i.Equal(progress.Speed, progress2.Speed)
	i.Equal(progress.ETA, progress2.ETA)
	i.Equal(progress.Error, progress2.Error)
	i.Equal(progress.FilePath, progress2.FilePath)
}

func TestReleaseInfoMarshaling(t *testing.T) {
	i := is.New(t)

	publishedAt := time.Now()

	info := ReleaseInfo{
		ID:          "release-123",
		Version:     "1.1.0",
		Name:        "Version 1.1.0",
		Description: "Bug fixes and improvements",
		PublishedAt: publishedAt,
		Channel:     "stable",
		Metadata: map[string]interface{}{
			"critical": true,
			"size":     1024 * 1024 * 50,
		},
		Artifacts: []ReleaseArtifact{
			{
				ID:       "artifact-123",
				Platform: "darwin",
				Arch:     "amd64",
				Filename: "app-darwin-amd64.zip",
				Size:     1024 * 1024 * 50,
				Checksum: "sha256:abcdef123456",
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(info)
	i.NoErr(err)

	// Unmarshal back
	var info2 ReleaseInfo
	err = json.Unmarshal(data, &info2)
	i.NoErr(err)

	// Compare
	i.Equal(info.ID, info2.ID)
	i.Equal(info.Version, info2.Version)
	i.Equal(info.Name, info2.Name)
	i.Equal(info.Description, info2.Description)
	i.Equal(info.PublishedAt.Unix(), info2.PublishedAt.Unix())
	i.Equal(info.Channel, info2.Channel)
	i.Equal(len(info.Artifacts), len(info2.Artifacts))

	// Check metadata
	i.True(info2.Metadata != nil)
	i.Equal(info2.Metadata["critical"], true)
	i.Equal(info2.Metadata["size"], float64(1024*1024*50)) // JSON numbers unmarshal as float64
}

func TestUpdateInstalledEventMarshaling(t *testing.T) {
	i := is.New(t)

	event := UpdateInstalledEvent{
		Version:     "1.1.0",
		InstalledAt: time.Now(),
	}

	// Marshal to JSON
	data, err := json.Marshal(event)
	i.NoErr(err)

	// Unmarshal back
	var event2 UpdateInstalledEvent
	err = json.Unmarshal(data, &event2)
	i.NoErr(err)

	// Compare
	i.Equal(event.Version, event2.Version)
	i.Equal(event.InstalledAt.Unix(), event2.InstalledAt.Unix())
}

func TestDownloadStateConstants(t *testing.T) {
	i := is.New(t)

	// Test that constants have expected values
	i.Equal(DownloadStatePending, "pending")
	i.Equal(DownloadStateDownloading, "downloading")
	i.Equal(DownloadStateCompleted, "completed")
	i.Equal(DownloadStateFailed, "failed")
	i.Equal(DownloadStateCancelled, "cancelled")
}

func TestValidationCodeConstants(t *testing.T) {
	i := is.New(t)

	// Test that constants have expected values
	i.Equal(ValidationCodeValid, "VALID")
	i.Equal(ValidationCodeInvalid, "INVALID")
	i.Equal(ValidationCodeExpired, "EXPIRED")
	i.Equal(ValidationCodeSuspended, "SUSPENDED")
	i.Equal(ValidationCodeOverdue, "OVERDUE")
	i.Equal(ValidationCodeNoMachines, "NO_MACHINES")
	i.Equal(ValidationCodeMachineLimitReached, "MACHINE_LIMIT_REACHED")
	i.Equal(ValidationCodeFingerprintMismatch, "FINGERPRINT_MISMATCH")
}

// Test edge cases
func TestTypesEdgeCases(t *testing.T) {
	i := is.New(t)

	// Test Config with empty/nil values
	config := Config{}
	data, err := json.Marshal(config)
	i.NoErr(err)

	var config2 Config
	err = json.Unmarshal(data, &config2)
	i.NoErr(err)
	i.Equal(config.CheckInterval, time.Duration(0))

	// Test LicenseState with nil pointers
	state := LicenseState{
		Valid:        false,
		License:      nil,
		ExpiresAt:    nil,
		Machine:      nil,
		Entitlements: nil,
	}
	data, err = json.Marshal(state)
	i.NoErr(err)

	var state2 LicenseState
	err = json.Unmarshal(data, &state2)
	i.NoErr(err)
	i.True(state2.License == nil)
	i.True(state2.ExpiresAt == nil)
	i.True(state2.Machine == nil)
	i.True(state2.Entitlements == nil)

	// Test UpdateInfo with empty artifacts
	info := UpdateInfo{
		Available: false,
		Artifacts: []ReleaseArtifact{},
	}
	data, err = json.Marshal(info)
	i.NoErr(err)

	var info2 UpdateInfo
	err = json.Unmarshal(data, &info2)
	i.NoErr(err)
	i.Equal(len(info2.Artifacts), 0)
}

// Test nested structures
func TestNestedStructures(t *testing.T) {
	i := is.New(t)

	// Create deeply nested entitlements
	entitlements := map[string]interface{}{
		"features": map[string]interface{}{
			"api": map[string]interface{}{
				"rate_limit": 1000,
				"endpoints":  []string{"v1", "v2"},
			},
			"storage": map[string]interface{}{
				"quota_gb": 100,
				"types":    []string{"images", "documents"},
			},
		},
		"limits": map[string]interface{}{
			"users":    10,
			"projects": 5,
		},
	}

	state := LicenseState{
		Valid:        true,
		Entitlements: entitlements,
	}

	// Marshal and unmarshal
	data, err := json.Marshal(state)
	i.NoErr(err)

	var state2 LicenseState
	err = json.Unmarshal(data, &state2)
	i.NoErr(err)

	// Verify nested structure
	features, ok := state2.Entitlements["features"].(map[string]interface{})
	i.True(ok)

	api, ok := features["api"].(map[string]interface{})
	i.True(ok)
	i.Equal(api["rate_limit"], float64(1000))

	endpoints, ok := api["endpoints"].([]interface{})
	i.True(ok)
	i.Equal(len(endpoints), 2)
	i.Equal(endpoints[0], "v1")
	i.Equal(endpoints[1], "v2")

	limits, ok := state2.Entitlements["limits"].(map[string]interface{})
	i.True(ok)
	i.Equal(limits["users"], float64(10))
	i.Equal(limits["projects"], float64(5))
}
