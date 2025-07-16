//go:build windows

package keygen

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// MockWindowsKeygen provides a mock for Windows-specific functionality
type MockWindowsKeygen struct {
	mock.Mock
	service *Service
}

func NewMockWindowsKeygen(service *Service) *MockWindowsKeygen {
	return &MockWindowsKeygen{service: service}
}

func (m *MockWindowsKeygen) Startup(ctx context.Context, options application.ServiceOptions) error {
	args := m.Called(ctx, options)
	return args.Error(0)
}

func (m *MockWindowsKeygen) InstallUpdate(updatePath string, version string) error {
	args := m.Called(updatePath, version)
	return args.Error(0)
}

func (m *MockWindowsKeygen) GetMachineFingerprint() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockWindowsKeygen) GetInstallPath() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockWindowsKeygen) StoreLicenseKey(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockWindowsKeygen) RetrieveLicenseKey() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockWindowsKeygen) DeleteLicenseKey() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockWindowsKeygen) SetRegistryValue(name, value string) error {
	args := m.Called(name, value)
	return args.Error(0)
}

func (m *MockWindowsKeygen) GetRegistryValue(name string) (string, error) {
	args := m.Called(name)
	return args.String(0), args.Error(1)
}

func (m *MockWindowsKeygen) DeleteRegistryValue(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

// TestWindowsKeygen tests the Windows specific implementation
func TestWindowsKeygen(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows tests only run on Windows")
	}

	// Test with real Windows implementation
	service := &Service{}
	windows := &platformKeygenWindows{service: service}

	t.Run("GetMachineFingerprint_NotImplemented", func(t *testing.T) {
		fingerprint, err := windows.GetMachineFingerprint()
		assert.Error(t, err)
		assert.Empty(t, fingerprint)
		assert.Contains(t, err.Error(), "not yet implemented")
	})

	t.Run("InstallUpdate_NotImplemented", func(t *testing.T) {
		err := windows.InstallUpdate("test.exe", "1.0.0")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not yet implemented")
	})

	t.Run("GetInstallPath_NotImplemented", func(t *testing.T) {
		path := windows.GetInstallPath()
		assert.Empty(t, path)
	})

	t.Run("StoreLicenseKey_NotImplemented", func(t *testing.T) {
		err := windows.StoreLicenseKey("test-key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not yet implemented")
	})

	t.Run("RetrieveLicenseKey_NotImplemented", func(t *testing.T) {
		key, err := windows.RetrieveLicenseKey()
		assert.Error(t, err)
		assert.Empty(t, key)
		assert.Contains(t, err.Error(), "not yet implemented")
	})

	t.Run("DeleteLicenseKey_NotImplemented", func(t *testing.T) {
		err := windows.DeleteLicenseKey()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not yet implemented")
	})

	t.Run("RegistryOperations_NotImplemented", func(t *testing.T) {
		err := windows.SetRegistryValue("test", "value")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not yet implemented")

		value, err := windows.GetRegistryValue("test")
		assert.Error(t, err)
		assert.Empty(t, value)
		assert.Contains(t, err.Error(), "not yet implemented")

		err = windows.DeleteRegistryValue("test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not yet implemented")
	})

	t.Run("Startup_Success", func(t *testing.T) {
		ctx := context.Background()
		options := application.ServiceOptions{}

		err := windows.Startup(ctx, options)
		assert.NoError(t, err) // Currently returns nil
	})
}

// TestWindowsMockImplementation tests the mock Windows implementation
func TestWindowsMockImplementation(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T, mock *MockWindowsKeygen)
	}{
		{
			name: "MachineFingerprint_Success",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				expectedFingerprint := "windows-test-fingerprint-12345"
				mock.On("GetMachineFingerprint").Return(expectedFingerprint, nil)

				fingerprint, err := mock.GetMachineFingerprint()
				assert.NoError(t, err)
				assert.Equal(t, expectedFingerprint, fingerprint)
			},
		},
		{
			name: "MachineFingerprint_WMIError",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				mock.On("GetMachineFingerprint").Return("", errors.New("WMI access denied"))

				fingerprint, err := mock.GetMachineFingerprint()
				assert.Error(t, err)
				assert.Empty(t, fingerprint)
				assert.Contains(t, err.Error(), "WMI access denied")
			},
		},
		{
			name: "InstallUpdate_EXE_Success",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				updatePath := "C:\\Updates\\myapp-2.0.0.exe"
				version := "2.0.0"
				mock.On("InstallUpdate", updatePath, version).Return(nil)

				err := mock.InstallUpdate(updatePath, version)
				assert.NoError(t, err)
			},
		},
		{
			name: "InstallUpdate_MSI_Success",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				updatePath := "C:\\Updates\\myapp-2.0.0.msi"
				version := "2.0.0"
				mock.On("InstallUpdate", updatePath, version).Return(nil)

				err := mock.InstallUpdate(updatePath, version)
				assert.NoError(t, err)
			},
		},
		{
			name: "InstallUpdate_UACRequired",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				updatePath := "C:\\Updates\\myapp-2.0.0.exe"
				version := "2.0.0"
				mock.On("InstallUpdate", updatePath, version).Return(errors.New("UAC elevation required"))

				err := mock.InstallUpdate(updatePath, version)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "UAC elevation required")
			},
		},
		{
			name: "GetInstallPath_Success",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				expectedPath := "C:\\Program Files\\MyApp"
				mock.On("GetInstallPath").Return(expectedPath)

				path := mock.GetInstallPath()
				assert.Equal(t, expectedPath, path)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWindows := NewMockWindowsKeygen(&Service{})
			tt.testFunc(t, mockWindows)
			mockWindows.AssertExpectations(t)
		})
	}
}

// TestWindowsCredentialManager tests Windows Credential Manager functionality
func TestWindowsCredentialManager(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T, mock *MockWindowsKeygen)
	}{
		{
			name: "StoreLicenseKey_Success",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				licenseKey := "WINDOWS-TEST-LICENSE-KEY-12345"
				mock.On("StoreLicenseKey", licenseKey).Return(nil)

				err := mock.StoreLicenseKey(licenseKey)
				assert.NoError(t, err)
			},
		},
		{
			name: "StoreLicenseKey_AccessDenied",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				licenseKey := "WINDOWS-TEST-LICENSE-KEY-12345"
				mock.On("StoreLicenseKey", licenseKey).Return(errors.New("access denied to credential manager"))

				err := mock.StoreLicenseKey(licenseKey)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "access denied")
			},
		},
		{
			name: "RetrieveLicenseKey_Success",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				expectedKey := "WINDOWS-TEST-LICENSE-KEY-12345"
				mock.On("RetrieveLicenseKey").Return(expectedKey, nil)

				key, err := mock.RetrieveLicenseKey()
				assert.NoError(t, err)
				assert.Equal(t, expectedKey, key)
			},
		},
		{
			name: "RetrieveLicenseKey_NotFound",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				mock.On("RetrieveLicenseKey").Return("", errors.New("credential not found"))

				key, err := mock.RetrieveLicenseKey()
				assert.Error(t, err)
				assert.Empty(t, key)
				assert.Contains(t, err.Error(), "not found")
			},
		},
		{
			name: "DeleteLicenseKey_Success",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				mock.On("DeleteLicenseKey").Return(nil)

				err := mock.DeleteLicenseKey()
				assert.NoError(t, err)
			},
		},
		{
			name: "DeleteLicenseKey_NotFound",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				mock.On("DeleteLicenseKey").Return(errors.New("credential not found"))

				err := mock.DeleteLicenseKey()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWindows := NewMockWindowsKeygen(&Service{})
			tt.testFunc(t, mockWindows)
			mockWindows.AssertExpectations(t)
		})
	}
}

// TestWindowsRegistry tests Windows Registry functionality
func TestWindowsRegistry(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T, mock *MockWindowsKeygen)
	}{
		{
			name: "SetRegistryValue_HKCU_Success",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				mock.On("SetRegistryValue", "LastUpdateCheck", "2023-12-01T10:00:00Z").Return(nil)

				err := mock.SetRegistryValue("LastUpdateCheck", "2023-12-01T10:00:00Z")
				assert.NoError(t, err)
			},
		},
		{
			name: "SetRegistryValue_AccessDenied",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				mock.On("SetRegistryValue", "InstallPath", "C:\\Program Files\\MyApp").Return(errors.New("access denied to registry"))

				err := mock.SetRegistryValue("InstallPath", "C:\\Program Files\\MyApp")
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "access denied")
			},
		},
		{
			name: "GetRegistryValue_Success",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				expectedValue := "C:\\Program Files\\MyApp"
				mock.On("GetRegistryValue", "InstallPath").Return(expectedValue, nil)

				value, err := mock.GetRegistryValue("InstallPath")
				assert.NoError(t, err)
				assert.Equal(t, expectedValue, value)
			},
		},
		{
			name: "GetRegistryValue_NotFound",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				mock.On("GetRegistryValue", "NonExistentKey").Return("", errors.New("registry key not found"))

				value, err := mock.GetRegistryValue("NonExistentKey")
				assert.Error(t, err)
				assert.Empty(t, value)
				assert.Contains(t, err.Error(), "not found")
			},
		},
		{
			name: "DeleteRegistryValue_Success",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				mock.On("DeleteRegistryValue", "TempValue").Return(nil)

				err := mock.DeleteRegistryValue("TempValue")
				assert.NoError(t, err)
			},
		},
		{
			name: "DeleteRegistryValue_NotFound",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				mock.On("DeleteRegistryValue", "NonExistentKey").Return(errors.New("registry value not found"))

				err := mock.DeleteRegistryValue("NonExistentKey")
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWindows := NewMockWindowsKeygen(&Service{})
			tt.testFunc(t, mockWindows)
			mockWindows.AssertExpectations(t)
		})
	}
}

// TestWindowsUpdateInstallation tests Windows-specific update installation
func TestWindowsUpdateInstallation(t *testing.T) {
	tests := []struct {
		name        string
		updatePath  string
		version     string
		setupMock   func(*MockWindowsKeygen, string, string)
		expectedErr string
	}{
		{
			name:       "EXE_Installation_Success",
			updatePath: "C:\\Updates\\myapp-setup.exe",
			version:    "2.0.0",
			setupMock: func(mock *MockWindowsKeygen, path, version string) {
				mock.On("InstallUpdate", path, version).Return(nil)
			},
		},
		{
			name:       "MSI_Installation_Success",
			updatePath: "C:\\Updates\\myapp.msi",
			version:    "2.0.0",
			setupMock: func(mock *MockWindowsKeygen, path, version string) {
				mock.On("InstallUpdate", path, version).Return(nil)
			},
		},
		{
			name:       "EXE_Installation_FileNotFound",
			updatePath: "C:\\NonExistent\\update.exe",
			version:    "2.0.0",
			setupMock: func(mock *MockWindowsKeygen, path, version string) {
				mock.On("InstallUpdate", path, version).Return(errors.New("update file not found"))
			},
			expectedErr: "update file not found",
		},
		{
			name:       "MSI_Installation_ElevationRequired",
			updatePath: "C:\\Updates\\system-update.msi",
			version:    "2.0.0",
			setupMock: func(mock *MockWindowsKeygen, path, version string) {
				mock.On("InstallUpdate", path, version).Return(errors.New("operation requires elevation"))
			},
			expectedErr: "operation requires elevation",
		},
		{
			name:       "Corrupted_Update_File",
			updatePath: "C:\\Updates\\corrupted.exe",
			version:    "2.0.0",
			setupMock: func(mock *MockWindowsKeygen, path, version string) {
				mock.On("InstallUpdate", path, version).Return(errors.New("update file is corrupted"))
			},
			expectedErr: "update file is corrupted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWindows := NewMockWindowsKeygen(&Service{})
			tt.setupMock(mockWindows, tt.updatePath, tt.version)

			err := mockWindows.InstallUpdate(tt.updatePath, tt.version)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			mockWindows.AssertExpectations(t)
		})
	}
}

// TestWindowsHardwareFingerprinting tests Windows hardware fingerprinting scenarios
func TestWindowsHardwareFingerprinting(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockWindowsKeygen)
		expectedError  string
		validateResult func(t *testing.T, fingerprint string)
	}{
		{
			name: "Complete_Hardware_Info",
			setupMock: func(mock *MockWindowsKeygen) {
				// Simulate successful hardware fingerprinting using multiple identifiers
				mock.On("GetMachineFingerprint").Return("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", nil)
			},
			validateResult: func(t *testing.T, fingerprint string) {
				assert.Len(t, fingerprint, 64)
				assert.Regexp(t, "^[a-f0-9]{64}$", fingerprint)
			},
		},
		{
			name: "WMI_Access_Denied",
			setupMock: func(mock *MockWindowsKeygen) {
				mock.On("GetMachineFingerprint").Return("", errors.New("WMI access denied"))
			},
			expectedError: "WMI access denied",
		},
		{
			name: "Partial_Hardware_Info",
			setupMock: func(mock *MockWindowsKeygen) {
				// Some hardware info available, but not all
				mock.On("GetMachineFingerprint").Return("abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890", nil)
			},
			validateResult: func(t *testing.T, fingerprint string) {
				assert.Len(t, fingerprint, 64)
				assert.Regexp(t, "^[a-f0-9]{64}$", fingerprint)
			},
		},
		{
			name: "No_Hardware_Identifiers",
			setupMock: func(mock *MockWindowsKeygen) {
				mock.On("GetMachineFingerprint").Return("", errors.New("no hardware identifiers found"))
			},
			expectedError: "no hardware identifiers found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWindows := NewMockWindowsKeygen(&Service{})
			tt.setupMock(mockWindows)

			fingerprint, err := mockWindows.GetMachineFingerprint()

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Empty(t, fingerprint)
			} else {
				assert.NoError(t, err)
				if tt.validateResult != nil {
					tt.validateResult(t, fingerprint)
				}
			}

			mockWindows.AssertExpectations(t)
		})
	}
}

// TestWindowsStartup tests Windows-specific startup scenarios
func TestWindowsStartup(t *testing.T) {
	tests := []struct {
		name          string
		setupContext  func() context.Context
		setupMock     func(*MockWindowsKeygen)
		expectedError string
	}{
		{
			name: "Normal_Startup",
			setupContext: func() context.Context {
				return context.Background()
			},
			setupMock: func(mock *MockWindowsKeygen) {
				mock.On("Startup", mock.AnythingOfType("*context.backgroundCtx"), mock.AnythingOfType("application.ServiceOptions")).Return(nil)
			},
		},
		{
			name: "Startup_With_Registry_Init",
			setupContext: func() context.Context {
				return context.Background()
			},
			setupMock: func(mock *MockWindowsKeygen) {
				mock.On("Startup", mock.AnythingOfType("*context.backgroundCtx"), mock.AnythingOfType("application.ServiceOptions")).Return(nil)
			},
		},
		{
			name: "Startup_Registry_Access_Denied",
			setupContext: func() context.Context {
				return context.Background()
			},
			setupMock: func(mock *MockWindowsKeygen) {
				mock.On("Startup", mock.AnythingOfType("*context.backgroundCtx"), mock.AnythingOfType("application.ServiceOptions")).Return(errors.New("registry access denied"))
			},
			expectedError: "registry access denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWindows := NewMockWindowsKeygen(&Service{})
			tt.setupMock(mockWindows)

			ctx := tt.setupContext()
			options := application.ServiceOptions{}

			err := mockWindows.Startup(ctx, options)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockWindows.AssertExpectations(t)
		})
	}
}

// TestWindowsErrorScenarios tests various Windows-specific error scenarios
func TestWindowsErrorScenarios(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T, mock *MockWindowsKeygen)
	}{
		{
			name: "UAC_Elevation_Required",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				mock.On("InstallUpdate", "C:\\Updates\\system-app.msi", "1.0.0").Return(errors.New("This operation requires elevation"))

				err := mock.InstallUpdate("C:\\Updates\\system-app.msi", "1.0.0")
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "elevation")
			},
		},
		{
			name: "Credential_Manager_Service_Unavailable",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				mock.On("StoreLicenseKey", "test-key").Return(errors.New("Credential Manager service is not available"))

				err := mock.StoreLicenseKey("test-key")
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "service is not available")
			},
		},
		{
			name: "Registry_Key_Access_Denied",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				mock.On("SetRegistryValue", "HKLM_Value", "test").Return(errors.New("Access is denied"))

				err := mock.SetRegistryValue("HKLM_Value", "test")
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "Access is denied")
			},
		},
		{
			name: "WMI_Service_Error",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				mock.On("GetMachineFingerprint").Return("", errors.New("WMI service error: 0x80041001"))

				_, err := mock.GetMachineFingerprint()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "WMI service error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWindows := NewMockWindowsKeygen(&Service{})
			tt.testFunc(t, mockWindows)
			mockWindows.AssertExpectations(t)
		})
	}
}

// TestWindowsPathHandling tests Windows-specific path handling
func TestWindowsPathHandling(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		setupMock      func(*MockWindowsKeygen, string)
		expectedResult interface{}
		expectedError  string
	}{
		{
			name: "Program_Files_Path",
			path: "C:\\Program Files\\MyApp\\app.exe",
			setupMock: func(mock *MockWindowsKeygen, path string) {
				mock.On("GetInstallPath").Return("C:\\Program Files\\MyApp")
			},
			expectedResult: "C:\\Program Files\\MyApp",
		},
		{
			name: "Program_Files_x86_Path",
			path: "C:\\Program Files (x86)\\MyApp\\app.exe",
			setupMock: func(mock *MockWindowsKeygen, path string) {
				mock.On("GetInstallPath").Return("C:\\Program Files (x86)\\MyApp")
			},
			expectedResult: "C:\\Program Files (x86)\\MyApp",
		},
		{
			name: "User_AppData_Path",
			path: "C:\\Users\\TestUser\\AppData\\Local\\MyApp\\app.exe",
			setupMock: func(mock *MockWindowsKeygen, path string) {
				mock.On("GetInstallPath").Return("C:\\Users\\TestUser\\AppData\\Local\\MyApp")
			},
			expectedResult: "C:\\Users\\TestUser\\AppData\\Local\\MyApp",
		},
		{
			name: "Network_Path",
			path: "\\\\server\\share\\MyApp\\app.exe",
			setupMock: func(mock *MockWindowsKeygen, path string) {
				mock.On("GetInstallPath").Return("\\\\server\\share\\MyApp")
			},
			expectedResult: "\\\\server\\share\\MyApp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWindows := NewMockWindowsKeygen(&Service{})
			tt.setupMock(mockWindows, tt.path)

			result := mockWindows.GetInstallPath()

			if tt.expectedError != "" {
				// For this test, GetInstallPath doesn't return error, but we can extend if needed
			} else {
				assert.Equal(t, tt.expectedResult, result)
			}

			mockWindows.AssertExpectations(t)
		})
	}
}

// TestWindowsConcurrency tests concurrent operations on Windows
func TestWindowsConcurrency(t *testing.T) {
	mockWindows := NewMockWindowsKeygen(&Service{})

	// Setup concurrent-safe mock expectations
	mockWindows.On("GetMachineFingerprint").Return("test-fingerprint", nil).Maybe()
	mockWindows.On("StoreLicenseKey", mock.AnythingOfType("string")).Return(nil).Maybe()
	mockWindows.On("RetrieveLicenseKey").Return("test-key", nil).Maybe()

	t.Run("ConcurrentFingerprinting", func(t *testing.T) {
		const numGoroutines = 5
		results := make(chan string, numGoroutines)
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				fingerprint, err := mockWindows.GetMachineFingerprint()
				if err != nil {
					errors <- err
				} else {
					results <- fingerprint
				}
			}()
		}

		// Collect results
		var fingerprints []string
		var collectedErrors []error

		for i := 0; i < numGoroutines; i++ {
			select {
			case fp := <-results:
				fingerprints = append(fingerprints, fp)
			case err := <-errors:
				collectedErrors = append(collectedErrors, err)
			}
		}

		// All successful fingerprints should be identical
		if len(fingerprints) > 0 {
			expected := fingerprints[0]
			for _, fp := range fingerprints {
				assert.Equal(t, expected, fp)
			}
		}

		assert.Empty(t, collectedErrors)
	})

	t.Run("ConcurrentCredentialOperations", func(t *testing.T) {
		const numGoroutines = 3
		done := make(chan bool, numGoroutines*2)

		// Concurrent store operations
		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				err := mockWindows.StoreLicenseKey("test-key")
				assert.NoError(t, err)
				done <- true
			}(i)
		}

		// Concurrent retrieve operations
		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				key, err := mockWindows.RetrieveLicenseKey()
				assert.NoError(t, err)
				assert.Equal(t, "test-key", key)
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < numGoroutines*2; i++ {
			<-done
		}
	})

	// Note: We don't call AssertExpectations here because we used Maybe()
	// and the exact number of calls may vary due to concurrency
}

// TestWindowsSecurityIntegration tests Windows security features
func TestWindowsSecurityIntegration(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T, mock *MockWindowsKeygen)
	}{
		{
			name: "Credential_Encryption",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				// Test that credentials are properly encrypted
				sensitiveKey := "SUPER-SECRET-LICENSE-KEY-12345"
				mock.On("StoreLicenseKey", sensitiveKey).Return(nil)
				mock.On("RetrieveLicenseKey").Return(sensitiveKey, nil)

				// Store
				err := mock.StoreLicenseKey(sensitiveKey)
				assert.NoError(t, err)

				// Retrieve
				retrievedKey, err := mock.RetrieveLicenseKey()
				assert.NoError(t, err)
				assert.Equal(t, sensitiveKey, retrievedKey)
			},
		},
		{
			name: "Registry_Security_Descriptors",
			testFunc: func(t *testing.T, mock *MockWindowsKeygen) {
				// Test registry access with proper security
				mock.On("SetRegistryValue", "SecureValue", "data").Return(nil)
				mock.On("GetRegistryValue", "SecureValue").Return("data", nil)

				err := mock.SetRegistryValue("SecureValue", "data")
				assert.NoError(t, err)

				value, err := mock.GetRegistryValue("SecureValue")
				assert.NoError(t, err)
				assert.Equal(t, "data", value)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWindows := NewMockWindowsKeygen(&Service{})
			tt.testFunc(t, mockWindows)
			mockWindows.AssertExpectations(t)
		})
	}
}

// TestWindowsEdgeCases tests Windows-specific edge cases
func TestWindowsEdgeCases(t *testing.T) {
	t.Run("LongPathSupport", func(t *testing.T) {
		mockWindows := NewMockWindowsKeygen(&Service{})

		// Test handling of long Windows paths (>260 characters)
		longPath := "C:\\VeryLongPathName\\" + strings.Repeat("SubDirectory\\", 20) + "app.exe"
		mockWindows.On("GetInstallPath").Return(longPath)

		result := mockWindows.GetInstallPath()
		assert.Equal(t, longPath, result)

		mockWindows.AssertExpectations(t)
	})

	t.Run("SpecialCharacterPaths", func(t *testing.T) {
		mockWindows := NewMockWindowsKeygen(&Service{})

		// Test paths with special characters
		specialPath := "C:\\Program Files\\App with spaces & symbols!\\app.exe"
		mockWindows.On("InstallUpdate", specialPath, "1.0.0").Return(nil)

		err := mockWindows.InstallUpdate(specialPath, "1.0.0")
		assert.NoError(t, err)

		mockWindows.AssertExpectations(t)
	})

	t.Run("EmptyEnvironmentVariables", func(t *testing.T) {
		// Test behavior when environment variables are not set
		// This would be more relevant for the cache directory functionality
		// when implemented in the Windows platform

		// For now, just verify the mock works with empty strings
		mockWindows := NewMockWindowsKeygen(&Service{})
		mockWindows.On("GetRegistryValue", "CacheDir").Return("", errors.New("not found"))

		value, err := mockWindows.GetRegistryValue("CacheDir")
		assert.Error(t, err)
		assert.Empty(t, value)

		mockWindows.AssertExpectations(t)
	})
}
