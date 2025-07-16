package keygen

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/keygen-sh/keygen-go/v3"
	"github.com/matryer/is"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// mockPlatformKeygen is a mock implementation of platformKeygen interface
type mockPlatformKeygen struct {
	fingerprint      string
	fingerprintErr   error
	installErr       error
	cacheDir         string
	installCallCount int
}

func (m *mockPlatformKeygen) GetMachineFingerprint() (string, error) {
	if m.fingerprintErr != nil {
		return "", m.fingerprintErr
	}
	if m.fingerprint == "" {
		return "test-fingerprint-12345", nil
	}
	return m.fingerprint, nil
}

func (m *mockPlatformKeygen) InstallUpdatePlatform(updatePath string) error {
	m.installCallCount++
	return m.installErr
}

func (m *mockPlatformKeygen) GetCacheDir() string {
	if m.cacheDir == "" {
		return os.TempDir()
	}
	return m.cacheDir
}

// mockEventEmitter for testing event emission
type mockEventEmitter struct {
	events []interface{}
	mu     sync.Mutex
}

func (m *mockEventEmitter) EmitLicenseStatus(event LicenseStatusEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, event)
}

func (m *mockEventEmitter) EmitUpdateAvailable(event UpdateAvailableEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, event)
}

func (m *mockEventEmitter) EmitDownloadProgress(event DownloadProgressEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, event)
}

func (m *mockEventEmitter) GetEvents() []interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]interface{}{}, m.events...)
}

// mockApp implements a minimal application interface for testing
type mockApp struct {
	Event *mockEventManager
}

type mockEventManager struct {
	events []interface{}
}

func (m *mockEventManager) Emit(name string, data ...interface{}) {
	m.events = append(m.events, map[string]interface{}{
		"name": name,
		"data": data,
	})
}

// Helper function to create a test service
func createTestService(t *testing.T) (*Service, *mockPlatformKeygen, *mockEventEmitter) {
	t.Helper()

	mockPlatform := &mockPlatformKeygen{
		cacheDir: filepath.Join(os.TempDir(), "keygen-test-"+t.Name()),
	}

	service := New(ServiceOptions{
		AccountID:      "test-account",
		ProductID:      "test-product",
		LicenseKey:     "test-license-key",
		PublicKey:      "test-public-key",
		CurrentVersion: "1.0.0",
		UpdateChannel:  "stable",
		AutoCheck:      false,
		CheckInterval:  24 * time.Hour,
		CacheDir:       mockPlatform.cacheDir,
	})

	// Replace platform implementation with mock
	service.impl = mockPlatform

	// Create mock event emitter
	mockEmitter := &mockEventEmitter{}
	service.eventEmitter = mockEmitter

	// Create cache directory
	os.MkdirAll(mockPlatform.cacheDir, 0755)

	// Cleanup function
	t.Cleanup(func() {
		os.RemoveAll(mockPlatform.cacheDir)
	})

	return service, mockPlatform, mockEmitter
}

// Test service initialization
func TestNew(t *testing.T) {
	i := is.New(t)

	options := ServiceOptions{
		AccountID:      "test-account",
		ProductID:      "test-product",
		LicenseKey:     "test-license",
		PublicKey:      "test-key",
		CurrentVersion: "1.0.0",
	}

	service := New(options)

	i.True(service != nil)
	i.Equal(service.config.Account, "test-account")
	i.Equal(service.config.Product, "test-product")
	i.Equal(service.config.LicenseKey, "test-license")
	i.Equal(service.config.PublicKey, "test-key")
	i.Equal(service.currentVersion, "1.0.0")

	// Check defaults
	i.Equal(service.config.Environment, "production")
	i.Equal(service.config.CheckInterval, 24*time.Hour)
	i.Equal(service.config.UpdateChannel, "stable")
}

func TestServiceName(t *testing.T) {
	i := is.New(t)
	service, _, _ := createTestService(t)

	i.Equal(service.ServiceName(), "github.com/wailsapp/wails/v3/pkg/services/keygen")
}

func TestServiceStartup(t *testing.T) {
	i := is.New(t)
	service, mockPlatform, _ := createTestService(t)

	mockApp := &mockApp{Event: &mockEventManager{}}

	err := service.ServiceStartup(context.Background(), application.ServiceOptions{
		App: mockApp,
	})

	i.NoErr(err)
	i.True(service.app != nil)
	i.Equal(service.cacheDir, mockPlatform.cacheDir)

	// Check cache directory was created
	_, err = os.Stat(service.cacheDir)
	i.NoErr(err)
}

func TestServiceShutdown(t *testing.T) {
	i := is.New(t)
	service, _, _ := createTestService(t)

	// Setup a valid license state
	service.state = &LicenseState{
		Valid: true,
		License: &keygen.License{
			ID: "test-license",
		},
	}

	err := service.ServiceShutdown()
	i.NoErr(err)

	// Check that offline license was saved
	licensePath := filepath.Join(service.cacheDir, "license.json")
	_, err = os.Stat(licensePath)
	i.NoErr(err)
}

func TestValidateLicense_NoKey(t *testing.T) {
	i := is.New(t)
	service, _, mockEmitter := createTestService(t)
	service.config.LicenseKey = ""

	result, err := service.ValidateLicense()

	i.True(err != nil)
	i.Equal(result.Valid, false)
	i.Equal(result.Code, ValidationCodeInvalid)
	i.Equal(result.Message, "No license key provided")

	// Check error type
	keygenErr, ok := err.(*KeygenError)
	i.True(ok)
	i.Equal(keygenErr.Code, ErrConfigInvalid)
}

func TestGetCurrentVersion(t *testing.T) {
	i := is.New(t)
	service, _, _ := createTestService(t)

	version := service.GetCurrentVersion()
	i.Equal(version, "1.0.0")
}

func TestSetUpdateChannel(t *testing.T) {
	i := is.New(t)
	service, _, _ := createTestService(t)

	// Test valid channels
	validChannels := []string{"stable", "beta", "alpha", "dev"}
	for _, channel := range validChannels {
		err := service.SetUpdateChannel(channel)
		i.NoErr(err)
		i.Equal(service.config.UpdateChannel, channel)
	}

	// Test invalid channel
	err := service.SetUpdateChannel("invalid")
	i.True(err != nil)
	keygenErr, ok := err.(*KeygenError)
	i.True(ok)
	i.Equal(keygenErr.Code, ErrUpdateChannelInvalid)
}

func TestCheckEntitlement(t *testing.T) {
	i := is.New(t)
	service, _, _ := createTestService(t)

	// Test with invalid license
	service.state.Valid = false
	enabled, err := service.CheckEntitlement("feature1")
	i.True(err != nil)
	i.Equal(enabled, false)

	// Test with valid license but no entitlements
	service.state.Valid = true
	service.state.Entitlements = nil
	enabled, err = service.CheckEntitlement("feature1")
	i.NoErr(err)
	i.Equal(enabled, false)

	// Test with entitlements
	service.state.Entitlements = map[string]interface{}{
		"feature1": true,
		"feature2": false,
		"feature3": "enabled",
		"feature4": map[string]interface{}{"enabled": true},
	}

	// Test enabled feature (boolean true)
	enabled, err = service.CheckEntitlement("feature1")
	i.NoErr(err)
	i.True(enabled)

	// Test disabled feature (boolean false)
	enabled, err = service.CheckEntitlement("feature2")
	i.NoErr(err)
	i.Equal(enabled, false)

	// Test enabled feature (non-false value)
	enabled, err = service.CheckEntitlement("feature3")
	i.NoErr(err)
	i.True(enabled)

	// Test enabled feature (complex value)
	enabled, err = service.CheckEntitlement("feature4")
	i.NoErr(err)
	i.True(enabled)

	// Test missing feature
	enabled, err = service.CheckEntitlement("missing")
	i.NoErr(err)
	i.Equal(enabled, false)
}

func TestSaveAndLoadOfflineLicense(t *testing.T) {
	i := is.New(t)
	service, _, _ := createTestService(t)

	// Create test license data
	testLicense := &keygen.License{
		ID:   "test-license-id",
		Key:  "test-license-key",
		Name: "Test License",
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	service.state = &LicenseState{
		Valid:       true,
		License:     testLicense,
		LastChecked: time.Now(),
		ExpiresAt:   &expiresAt,
		Entitlements: map[string]interface{}{
			"feature1": true,
			"feature2": "enabled",
		},
	}

	// Save offline license
	err := service.SaveOfflineLicense()
	i.NoErr(err)

	// Check file exists
	licensePath := filepath.Join(service.cacheDir, "license.json")
	_, err = os.Stat(licensePath)
	i.NoErr(err)

	// Create new service and load offline license
	service2, _, _ := createTestService(t)
	service2.cacheDir = service.cacheDir

	err = service2.LoadOfflineLicense()
	i.NoErr(err)

	// Check state was restored
	i.True(service2.state.OfflineMode)
	i.True(service2.state.Entitlements != nil)
	i.Equal(service2.state.Entitlements["feature1"], true)
	i.Equal(service2.state.Entitlements["feature2"], "enabled")
}

func TestClearLicenseCache(t *testing.T) {
	i := is.New(t)
	service, _, _ := createTestService(t)

	// Create a license file
	licensePath := filepath.Join(service.cacheDir, "license.json")
	err := os.WriteFile(licensePath, []byte("{}"), 0600)
	i.NoErr(err)

	// Set some state
	service.state = &LicenseState{
		Valid: true,
		Entitlements: map[string]interface{}{
			"feature1": true,
		},
	}

	// Clear cache
	err = service.ClearLicenseCache()
	i.NoErr(err)

	// Check state was cleared
	i.Equal(service.state.Valid, false)
	i.True(service.state.Entitlements == nil)

	// Check file was removed
	_, err = os.Stat(licensePath)
	i.True(os.IsNotExist(err))
}

func TestCompareVersions(t *testing.T) {
	i := is.New(t)

	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.1", "1.0.0", 1},
		{"1.0.0", "1.0.1", -1},
		{"2.0.0", "1.9.9", 1},
		{"1.9.9", "2.0.0", -1},
		{"v1.0.0", "1.0.0", 0},
		{"1.0.0", "v1.0.0", 0},
		{"1.0", "1.0.0", 0},
		{"1.0.0", "1.0", 0},
		{"1.2.3", "1.2", 1},
		{"1.2", "1.2.3", -1},
	}

	for _, tt := range tests {
		result := compareVersions(tt.v1, tt.v2)
		i.Equal(result, tt.expected)
	}
}

func TestProgressReader(t *testing.T) {
	i := is.New(t)

	// Create test data
	testData := []byte("Hello, World! This is test data for progress tracking.")
	reader := &progressReader{
		reader: httptest.NewRequest("GET", "/", nil).Body,
		progress: &DownloadProgress{
			TotalBytes: int64(len(testData)),
		},
		onProgress: func(p *DownloadProgress) {
			// Progress callback
		},
	}

	// Read data
	buf := make([]byte, 10)
	n, err := reader.Read(buf)

	// Since we're reading from an empty body, we expect 0 bytes and EOF
	i.Equal(n, 0)
	i.True(err != nil) // EOF
}

func TestDownloadProgress(t *testing.T) {
	i := is.New(t)
	service, _, mockEmitter := createTestService(t)

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1024")
		w.WriteHeader(http.StatusOK)
		// Write some test data
		data := make([]byte, 1024)
		w.Write(data)
	}))
	defer server.Close()

	// Create test artifact
	artifact := &keygen.Artifact{
		ID:       "test-artifact",
		Filename: "test-update.zip",
		Filesize: 1024,
		Checksum: "",
	}

	// Create download directory
	downloadDir := filepath.Join(service.cacheDir, "downloads")
	os.MkdirAll(downloadDir, 0755)

	progress := &DownloadProgress{
		ID:         artifact.ID,
		State:      DownloadStateDownloading,
		TotalBytes: artifact.Filesize,
		FilePath:   filepath.Join(downloadDir, artifact.Filename),
	}

	// Test progress update
	service.emitDownloadProgress(progress)

	events := mockEmitter.GetEvents()
	i.Equal(len(events), 1)

	event, ok := events[0].(DownloadProgressEvent)
	i.True(ok)
	i.Equal(event.Status, DownloadStateDownloading)
}

func TestMachineActivation(t *testing.T) {
	i := is.New(t)
	service, mockPlatform, _ := createTestService(t)

	// Test without validated license
	_, err := service.ActivateMachine()
	i.True(err != nil)
	keygenErr, ok := err.(*KeygenError)
	i.True(ok)
	i.Equal(keygenErr.Code, ErrLicenseInvalid)

	// Test with fingerprint error
	mockPlatform.fingerprintErr = fmt.Errorf("fingerprint error")
	service.state.License = &keygen.License{ID: "test"}

	_, err = service.ActivateMachine()
	i.True(err != nil)
}

func TestDeactivateMachine(t *testing.T) {
	i := is.New(t)
	service, _, _ := createTestService(t)

	// Test without activated machine
	err := service.DeactivateMachine()
	i.True(err != nil)
	keygenErr, ok := err.(*KeygenError)
	i.True(ok)
	i.Equal(keygenErr.Code, ErrMachineNotActivated)
}

func TestGetLicenseInfo(t *testing.T) {
	i := is.New(t)
	service, _, _ := createTestService(t)

	// Test without license
	_, err := service.GetLicenseInfo()
	i.True(err != nil)

	// Test with license
	createdAt := time.Now()
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	maxMachines := 5

	service.state.License = &keygen.License{
		ID:          "test-license",
		CreatedAt:   createdAt,
		Expiry:      &expiresAt,
		MaxMachines: &maxMachines,
		Metadata: map[string]interface{}{
			"name":    "Test User",
			"email":   "test@example.com",
			"company": "Test Corp",
		},
	}
	service.state.LastChecked = time.Now()
	service.state.Entitlements = map[string]interface{}{
		"feature1": true,
	}

	info, err := service.GetLicenseInfo()
	i.NoErr(err)
	i.Equal(info.Key, "test-license-key")
	i.Equal(info.Status, "active")
	i.Equal(info.Name, "Test User")
	i.Equal(info.Email, "test@example.com")
	i.Equal(info.Company, "Test Corp")
	i.Equal(*info.MaxMachines, 5)
	i.True(info.ExpiresAt != nil)
	i.True(info.Entitlements != nil)
}

func TestAutoUpdateCheck(t *testing.T) {
	i := is.New(t)

	// Create service with auto-check enabled
	service := New(ServiceOptions{
		AccountID:      "test-account",
		ProductID:      "test-product",
		LicenseKey:     "test-license",
		CurrentVersion: "1.0.0",
		AutoCheck:      true,
		CheckInterval:  100 * time.Millisecond, // Short interval for testing
	})

	mockPlatform := &mockPlatformKeygen{
		cacheDir: filepath.Join(os.TempDir(), "keygen-test-auto"),
	}
	service.impl = mockPlatform

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// Start service
	err := service.ServiceStartup(ctx, application.ServiceOptions{})
	i.NoErr(err)

	// Wait for context to expire
	<-ctx.Done()

	// Cleanup
	os.RemoveAll(mockPlatform.cacheDir)
}

func TestEventEmission(t *testing.T) {
	i := is.New(t)
	service, _, mockEmitter := createTestService(t)

	// Test license status event
	service.emitLicenseStatus(true, "License is valid", "")

	events := mockEmitter.GetEvents()
	i.Equal(len(events), 1)

	event, ok := events[0].(LicenseStatusEvent)
	i.True(ok)
	i.True(event.Valid)
	i.Equal(event.Status, "active")
	i.Equal(event.Key, "test-license-key")
	i.Equal(event.Error, "")
}

// Test concurrent operations
func TestConcurrentOperations(t *testing.T) {
	i := is.New(t)
	service, _, _ := createTestService(t)

	// Set initial state
	service.state = &LicenseState{
		Valid: true,
		Entitlements: map[string]interface{}{
			"feature1": true,
			"feature2": false,
		},
	}

	var wg sync.WaitGroup
	errors := make([]error, 10)

	// Run concurrent entitlement checks
	for j := 0; j < 10; j++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			_, err := service.CheckEntitlement(fmt.Sprintf("feature%d", index%2+1))
			errors[index] = err
		}(j)
	}

	wg.Wait()

	// Check no errors occurred
	for _, err := range errors {
		i.NoErr(err)
	}
}

// Test error handling
func TestLoadOfflineLicenseErrors(t *testing.T) {
	i := is.New(t)
	service, _, _ := createTestService(t)

	// Test with corrupted data
	licensePath := filepath.Join(service.cacheDir, "license.json")
	err := os.WriteFile(licensePath, []byte("invalid json"), 0600)
	i.NoErr(err)

	err = service.LoadOfflineLicense()
	i.True(err != nil)
	keygenErr, ok := err.(*KeygenError)
	i.True(ok)
	i.Equal(keygenErr.Code, ErrCacheCorrupted)

	// Test with non-existent file (should not error)
	os.Remove(licensePath)
	err = service.LoadOfflineLicense()
	i.NoErr(err)
}

// Test download cancellation
func TestDownloadCancellation(t *testing.T) {
	i := is.New(t)
	service, _, _ := createTestService(t)

	// Set up download in progress
	ctx, cancel := context.WithCancel(context.Background())
	service.downloadCancel = cancel
	service.updateInProgress = true

	// Simulate shutdown
	err := service.ServiceShutdown()
	i.NoErr(err)

	// Check context was cancelled
	select {
	case <-ctx.Done():
		// Success - context was cancelled
	default:
		t.Fatal("Download context was not cancelled")
	}
}
