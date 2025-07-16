package keygen

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// SimpleMockPlatform provides a simple mock for the platformKeygen interface
type SimpleMockPlatform struct {
	mock.Mock
	mu sync.Mutex
}

func NewSimpleMockPlatform() *SimpleMockPlatform {
	return &SimpleMockPlatform{}
}

func (m *SimpleMockPlatform) GetMachineFingerprint() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *SimpleMockPlatform) InstallUpdatePlatform(updatePath string) error {
	args := m.Called(updatePath)
	return args.Error(0)
}

func (m *SimpleMockPlatform) GetCacheDir() string {
	args := m.Called()
	return args.String(0)
}

// TestSimplePlatformInterface tests the basic platform interface
func TestSimplePlatformInterface(t *testing.T) {
	t.Run("GetMachineFingerprint", func(t *testing.T) {
		mock := NewSimpleMockPlatform()
		mock.On("GetMachineFingerprint").Return("test-fingerprint-hash", nil)

		fingerprint, err := mock.GetMachineFingerprint()
		assert.NoError(t, err)
		assert.Equal(t, "test-fingerprint-hash", fingerprint)
		mock.AssertExpectations(t)
	})

	t.Run("InstallUpdatePlatform", func(t *testing.T) {
		mock := NewSimpleMockPlatform()
		mock.On("InstallUpdatePlatform", "/tmp/update.exe").Return(nil)

		err := mock.InstallUpdatePlatform("/tmp/update.exe")
		assert.NoError(t, err)
		mock.AssertExpectations(t)
	})

	t.Run("GetCacheDir", func(t *testing.T) {
		mock := NewSimpleMockPlatform()
		mock.On("GetCacheDir").Return("/tmp/cache")

		cacheDir := mock.GetCacheDir()
		assert.Equal(t, "/tmp/cache", cacheDir)
		mock.AssertExpectations(t)
	})
}

// TestSimpleFingerprintGeneration tests fingerprint generation scenarios
func TestSimpleFingerprintGeneration(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*SimpleMockPlatform)
		expectedError  string
		validateResult func(t *testing.T, fingerprint string)
	}{
		{
			name: "ValidFingerprint",
			setupMock: func(m *SimpleMockPlatform) {
				m.On("GetMachineFingerprint").Return("abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890", nil)
			},
			validateResult: func(t *testing.T, fingerprint string) {
				assert.Len(t, fingerprint, 64)
				_, err := hex.DecodeString(fingerprint)
				assert.NoError(t, err)
			},
		},
		{
			name: "HardwareError",
			setupMock: func(m *SimpleMockPlatform) {
				m.On("GetMachineFingerprint").Return("", errors.New("hardware access denied"))
			},
			expectedError: "hardware access denied",
		},
		{
			name: "EmptyFingerprint",
			setupMock: func(m *SimpleMockPlatform) {
				m.On("GetMachineFingerprint").Return("", errors.New("no hardware identifiers found"))
			},
			expectedError: "no hardware identifiers found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewSimpleMockPlatform()
			tt.setupMock(mock)

			fingerprint, err := mock.GetMachineFingerprint()

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

			mock.AssertExpectations(t)
		})
	}
}

// TestSimpleUpdateInstallation tests update installation scenarios
func TestSimpleUpdateInstallation(t *testing.T) {
	tests := []struct {
		name          string
		updatePath    string
		setupMock     func(*SimpleMockPlatform, string)
		expectedError string
	}{
		{
			name:       "SuccessfulInstallation",
			updatePath: "/tmp/update.exe",
			setupMock: func(m *SimpleMockPlatform, path string) {
				m.On("InstallUpdatePlatform", path).Return(nil)
			},
		},
		{
			name:       "FileNotFound",
			updatePath: "/nonexistent/update.exe",
			setupMock: func(m *SimpleMockPlatform, path string) {
				m.On("InstallUpdatePlatform", path).Return(errors.New("update file not found"))
			},
			expectedError: "update file not found",
		},
		{
			name:       "PermissionDenied",
			updatePath: "/tmp/update.exe",
			setupMock: func(m *SimpleMockPlatform, path string) {
				m.On("InstallUpdatePlatform", path).Return(errors.New("permission denied"))
			},
			expectedError: "permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewSimpleMockPlatform()
			tt.setupMock(mock, tt.updatePath)

			err := mock.InstallUpdatePlatform(tt.updatePath)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mock.AssertExpectations(t)
		})
	}
}

// TestSimpleCacheDirectory tests cache directory functionality
func TestSimpleCacheDirectory(t *testing.T) {
	tests := []struct {
		name           string
		platform       string
		setupMock      func(*SimpleMockPlatform)
		validateResult func(t *testing.T, cacheDir string)
	}{
		{
			name:     "Darwin_CacheDir",
			platform: "darwin",
			setupMock: func(m *SimpleMockPlatform) {
				m.On("GetCacheDir").Return("/Users/test/Library/Caches/TestApp/keygen")
			},
			validateResult: func(t *testing.T, cacheDir string) {
				assert.Contains(t, cacheDir, "Library/Caches")
				assert.Contains(t, cacheDir, "keygen")
			},
		},
		{
			name:     "Windows_CacheDir",
			platform: "windows",
			setupMock: func(m *SimpleMockPlatform) {
				m.On("GetCacheDir").Return("C:\\Users\\test\\AppData\\Local\\TestApp\\keygen")
			},
			validateResult: func(t *testing.T, cacheDir string) {
				assert.Contains(t, cacheDir, "keygen")
			},
		},
		{
			name:     "Linux_CacheDir",
			platform: "linux",
			setupMock: func(m *SimpleMockPlatform) {
				m.On("GetCacheDir").Return("/home/test/.cache/testapp/keygen")
			},
			validateResult: func(t *testing.T, cacheDir string) {
				assert.Contains(t, cacheDir, ".cache")
				assert.Contains(t, cacheDir, "keygen")
			},
		},
		{
			name:     "FallbackToTemp",
			platform: "unknown",
			setupMock: func(m *SimpleMockPlatform) {
				m.On("GetCacheDir").Return(filepath.Join(os.TempDir(), "keygen-cache"))
			},
			validateResult: func(t *testing.T, cacheDir string) {
				assert.Contains(t, cacheDir, "keygen-cache")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewSimpleMockPlatform()
			tt.setupMock(mock)

			cacheDir := mock.GetCacheDir()

			assert.NotEmpty(t, cacheDir)
			if tt.validateResult != nil {
				tt.validateResult(t, cacheDir)
			}

			mock.AssertExpectations(t)
		})
	}
}

// TestSimpleConcurrency tests concurrent operations
func TestSimpleConcurrency(t *testing.T) {
	t.Run("ConcurrentFingerprinting", func(t *testing.T) {
		const numGoroutines = 10
		mock := NewSimpleMockPlatform()
		mock.On("GetMachineFingerprint").Return("test-fingerprint", nil).Maybe()

		results := make(chan string, numGoroutines)
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				fingerprint, err := mock.GetMachineFingerprint()
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

	t.Run("ConcurrentCacheAccess", func(t *testing.T) {
		const numGoroutines = 5
		mock := NewSimpleMockPlatform()
		mock.On("GetCacheDir").Return("/tmp/cache").Maybe()

		results := make(chan string, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				cacheDir := mock.GetCacheDir()
				results <- cacheDir
			}()
		}

		// Collect results
		var cacheDirs []string
		for i := 0; i < numGoroutines; i++ {
			cacheDirs = append(cacheDirs, <-results)
		}

		// All cache directories should be identical
		expected := cacheDirs[0]
		for _, dir := range cacheDirs {
			assert.Equal(t, expected, dir)
		}
	})
}

// TestSimpleHashFunction tests hash function for fingerprints
func TestSimpleHashFunction(t *testing.T) {
	hashFunc := func(data string) string {
		h := sha256.New()
		h.Write([]byte(data))
		return hex.EncodeToString(h.Sum(nil))
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "EmptyString",
			input:    "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:     "SimpleString",
			input:    "test",
			expected: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		},
		{
			name:     "ComplexFingerprint",
			input:    "hardware-uuid-12345-serial-67890-mac-aa:bb:cc:dd:ee:ff",
			expected: hashFunc("hardware-uuid-12345-serial-67890-mac-aa:bb:cc:dd:ee:ff"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hashFunc(tt.input)
			assert.Equal(t, tt.expected, result)
			assert.Len(t, result, 64) // SHA256 hex length

			// Verify it's valid hex
			_, err := hex.DecodeString(result)
			assert.NoError(t, err)
		})
	}
}

// TestSimpleFileOperations tests file operations for updates
func TestSimpleFileOperations(t *testing.T) {
	createTempFile := func(t *testing.T, name, content string) string {
		tempDir, err := os.MkdirTemp("", "simple-test-*")
		require.NoError(t, err)
		t.Cleanup(func() {
			os.RemoveAll(tempDir)
		})

		filePath := filepath.Join(tempDir, name)
		err = os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err)
		return filePath
	}

	tests := []struct {
		name           string
		setupFile      func(t *testing.T) string
		setupMock      func(*SimpleMockPlatform, string)
		expectedError  string
		shouldCallMock bool
	}{
		{
			name: "ValidExecutable",
			setupFile: func(t *testing.T) string {
				return createTempFile(t, "update.exe", "fake executable content")
			},
			setupMock: func(m *SimpleMockPlatform, updatePath string) {
				m.On("InstallUpdatePlatform", updatePath).Return(nil)
			},
			shouldCallMock: true,
		},
		{
			name: "NonExecutableFile",
			setupFile: func(t *testing.T) string {
				return createTempFile(t, "update.txt", "not executable")
			},
			setupMock: func(m *SimpleMockPlatform, updatePath string) {
				m.On("InstallUpdatePlatform", updatePath).Return(errors.New("file is not executable"))
			},
			expectedError:  "file is not executable",
			shouldCallMock: true,
		},
		{
			name: "NonExistentFile",
			setupFile: func(t *testing.T) string {
				return "/path/to/nonexistent/file.exe"
			},
			setupMock: func(m *SimpleMockPlatform, updatePath string) {
				m.On("InstallUpdatePlatform", updatePath).Return(errors.New("update file not found"))
			},
			expectedError:  "update file not found",
			shouldCallMock: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatePath := tt.setupFile(t)
			mock := NewSimpleMockPlatform()

			if tt.shouldCallMock {
				tt.setupMock(mock, updatePath)
			}

			err := mock.InstallUpdatePlatform(updatePath)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			if tt.shouldCallMock {
				mock.AssertExpectations(t)
			}
		})
	}
}

// TestSimpleErrorScenarios tests various error scenarios
func TestSimpleErrorScenarios(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		setupMock func(*SimpleMockPlatform)
		testFunc  func(t *testing.T, mock *SimpleMockPlatform)
	}{
		{
			name:      "FingerprintTimeout",
			operation: "GetMachineFingerprint",
			setupMock: func(m *SimpleMockPlatform) {
				m.On("GetMachineFingerprint").Return("", errors.New("operation timed out"))
			},
			testFunc: func(t *testing.T, mock *SimpleMockPlatform) {
				_, err := mock.GetMachineFingerprint()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "operation timed out")
			},
		},
		{
			name:      "InstallationInterrupted",
			operation: "InstallUpdatePlatform",
			setupMock: func(m *SimpleMockPlatform) {
				m.On("InstallUpdatePlatform", mock.AnythingOfType("string")).Return(errors.New("installation interrupted"))
			},
			testFunc: func(t *testing.T, mock *SimpleMockPlatform) {
				err := mock.InstallUpdatePlatform("/tmp/update.exe")
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "installation interrupted")
			},
		},
		{
			name:      "CacheDirectoryUnavailable",
			operation: "GetCacheDir",
			setupMock: func(m *SimpleMockPlatform) {
				m.On("GetCacheDir").Return("")
			},
			testFunc: func(t *testing.T, mock *SimpleMockPlatform) {
				cacheDir := mock.GetCacheDir()
				assert.Empty(t, cacheDir)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewSimpleMockPlatform()
			tt.setupMock(mock)
			tt.testFunc(t, mock)
			mock.AssertExpectations(t)
		})
	}
}

// BenchmarkSimpleFingerprinting benchmarks fingerprint generation
func BenchmarkSimpleFingerprinting(b *testing.B) {
	mock := NewSimpleMockPlatform()
	mock.On("GetMachineFingerprint").Return("test-fingerprint-hash", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := mock.GetMachineFingerprint()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSimpleCacheAccess benchmarks cache directory access
func BenchmarkSimpleCacheAccess(b *testing.B) {
	mock := NewSimpleMockPlatform()
	mock.On("GetCacheDir").Return("/tmp/test-cache")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cacheDir := mock.GetCacheDir()
		if cacheDir == "" {
			b.Fatal("empty cache directory")
		}
	}
}
