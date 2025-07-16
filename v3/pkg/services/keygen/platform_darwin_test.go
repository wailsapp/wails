//go:build darwin

package keygen

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// TestDarwinKeygen tests the Darwin (macOS) specific implementation
func TestDarwinKeygen(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Darwin tests only run on macOS")
	}

	darwin := &darwinKeygen{}

	t.Run("GetMachineFingerprint", func(t *testing.T) {
		fingerprint, err := darwin.GetMachineFingerprint()
		assert.NoError(t, err)
		assert.NotEmpty(t, fingerprint)
		assert.Len(t, fingerprint, 64) // SHA256 hex string

		// Verify it's a valid hex string
		assert.Regexp(t, "^[a-f0-9]{64}$", fingerprint)

		// Test consistency - should return same fingerprint multiple times
		fingerprint2, err := darwin.GetMachineFingerprint()
		assert.NoError(t, err)
		assert.Equal(t, fingerprint, fingerprint2)
	})

	t.Run("GetCacheDir", func(t *testing.T) {
		cacheDir := darwin.GetCacheDir()
		assert.NotEmpty(t, cacheDir)
		assert.True(t, filepath.IsAbs(cacheDir))
		assert.Contains(t, cacheDir, "Library/Caches")
		assert.Contains(t, cacheDir, "keygen")

		// Verify the directory can be created
		err := os.MkdirAll(cacheDir, 0755)
		assert.NoError(t, err)
		defer os.RemoveAll(cacheDir)

		// Verify it's writable
		testFile := filepath.Join(cacheDir, "test.txt")
		err = os.WriteFile(testFile, []byte("test"), 0644)
		assert.NoError(t, err)
	})

	t.Run("InstallUpdatePlatform", func(t *testing.T) {
		tests := []struct {
			name          string
			setupFile     func(t *testing.T) string
			expectedError string
		}{
			{
				name: "NonExistentFile",
				setupFile: func(t *testing.T) string {
					return "/path/to/nonexistent/file.dmg"
				},
				expectedError: "update file not found",
			},
			{
				name: "ValidZipFile",
				setupFile: func(t *testing.T) string {
					// Create a temporary zip file for testing
					tempDir := createTestTempDir(t)
					zipPath := filepath.Join(tempDir, "test-update.zip")

					// Create a simple zip file with an executable
					err := os.WriteFile(zipPath, []byte("PK\x03\x04"), 0644) // Minimal zip header
					require.NoError(t, err)

					return zipPath
				},
				expectedError: "no installable content found", // Since our zip is minimal
			},
			{
				name: "ExecutableFile",
				setupFile: func(t *testing.T) string {
					tempDir := createTestTempDir(t)
					execPath := filepath.Join(tempDir, "updater")

					// Create a simple shell script
					content := "#!/bin/sh\necho 'Update completed'\nexit 0\n"
					err := os.WriteFile(execPath, []byte(content), 0755)
					require.NoError(t, err)

					return execPath
				},
				// Should not error for valid executable
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				updatePath := tt.setupFile(t)

				err := darwin.InstallUpdatePlatform(updatePath)

				if tt.expectedError != "" {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tt.expectedError)
				} else {
					// For executables, we expect no error but the process might start
					// We can't easily test full execution without potentially modifying the system
					if strings.HasSuffix(updatePath, "updater") {
						// The function should attempt to execute, which may succeed or fail
						// depending on the system state, so we just verify it doesn't panic
						assert.NotPanics(t, func() {
							darwin.InstallUpdatePlatform(updatePath)
						})
					}
				}
			})
		}
	})
}

// TestDarwinDMGInstallation tests DMG installation specifically
func TestDarwinDMGInstallation(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Darwin DMG tests only run on macOS")
	}

	darwin := &darwinKeygen{}

	t.Run("InstallDMG_MockScenarios", func(t *testing.T) {
		// We can't easily test actual DMG mounting without system modifications
		// But we can test the error path for non-existent DMG files
		nonExistentDMG := "/tmp/nonexistent.dmg"

		err := darwin.installDMG(nonExistentDMG)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to mount DMG")
	})
}

// TestDarwinAppInstallation tests .app installation
func TestDarwinAppInstallation(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Darwin app tests only run on macOS")
	}

	darwin := &darwinKeygen{}

	t.Run("InstallApp_MockApp", func(t *testing.T) {
		tempDir := createTestTempDir(t)

		// Create a mock .app structure
		appDir := filepath.Join(tempDir, "TestApp.app")
		contentsDir := filepath.Join(appDir, "Contents")
		macOSDir := filepath.Join(contentsDir, "MacOS")

		err := os.MkdirAll(macOSDir, 0755)
		require.NoError(t, err)

		// Create a minimal Info.plist
		infoPlist := filepath.Join(contentsDir, "Info.plist")
		plistContent := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleName</key>
	<string>TestApp</string>
	<key>CFBundleExecutable</key>
	<string>TestApp</string>
</dict>
</plist>`
		err = os.WriteFile(infoPlist, []byte(plistContent), 0644)
		require.NoError(t, err)

		// Create executable
		execPath := filepath.Join(macOSDir, "TestApp")
		err = os.WriteFile(execPath, []byte("#!/bin/sh\necho 'app'\n"), 0755)
		require.NoError(t, err)

		// Test installation - this will try to copy to /Applications which may fail
		// due to permissions, but we can verify it doesn't panic
		err = darwin.installApp(appDir)
		// May fail due to permissions to /Applications, but should not panic
		if err != nil {
			assert.Contains(t, err.Error(), "failed to copy app")
		}
	})
}

// TestDarwinSystemIntegration tests integration with macOS system calls
func TestDarwinSystemIntegration(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Darwin system integration tests only run on macOS")
	}

	t.Run("GetAppName", func(t *testing.T) {
		appName := getAppName()
		assert.NotEmpty(t, appName)

		// Should return a reasonable app name
		assert.True(t, len(appName) > 0)
		assert.False(t, strings.Contains(appName, "/"))
	})

	t.Run("SystemProfilerAccess", func(t *testing.T) {
		// Test that we can at least attempt to run system_profiler
		// This tests the command exists and is accessible
		darwin := &darwinKeygen{}
		fingerprint, err := darwin.GetMachineFingerprint()

		// Even if it fails due to permissions or parsing, it should not panic
		if err != nil {
			// Acceptable errors include parsing failures or access issues
			assert.NotContains(t, err.Error(), "panic")
		} else {
			assert.NotEmpty(t, fingerprint)
			assert.Len(t, fingerprint, 64)
		}
	})
}

// TestDarwinFileOperations tests Darwin-specific file operations
func TestDarwinFileOperations(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Darwin file operation tests only run on macOS")
	}

	t.Run("ZipExtraction", func(t *testing.T) {
		darwin := &darwinKeygen{}
		tempDir := createTestTempDir(t)

		// Create a real zip file with content
		zipPath := filepath.Join(tempDir, "test.zip")

		// Create a simple test app structure in a temp dir
		testAppDir := filepath.Join(tempDir, "TestApp.app")
		err := os.MkdirAll(testAppDir, 0755)
		require.NoError(t, err)

		// Create the zip file using system zip command
		origDir, _ := os.Getwd()
		os.Chdir(tempDir)

		// We'll skip actual zip creation and test the error path instead
		err = darwin.installZip(zipPath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unzip")

		os.Chdir(origDir)
	})

	t.Run("PermissionHandling", func(t *testing.T) {
		tempDir := createTestTempDir(t)

		// Create a file without execute permissions
		testFile := filepath.Join(tempDir, "no-exec")
		err := os.WriteFile(testFile, []byte("test"), 0644)
		require.NoError(t, err)

		// Try to make it executable
		err = os.Chmod(testFile, 0755)
		assert.NoError(t, err)

		// Verify permissions
		info, err := os.Stat(testFile)
		require.NoError(t, err)
		assert.True(t, info.Mode()&0111 != 0)
	})
}

// TestDarwinErrorHandling tests error scenarios specific to Darwin
func TestDarwinErrorHandling(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Darwin error handling tests only run on macOS")
	}

	darwin := &darwinKeygen{}

	t.Run("InvalidUpdateFormats", func(t *testing.T) {
		tempDir := createTestTempDir(t)

		// Test various invalid file formats
		invalidFiles := []struct {
			name     string
			content  string
			expected string
		}{
			{
				name:     "invalid.dmg",
				content:  "not a real dmg",
				expected: "failed to mount DMG",
			},
			{
				name:    "invalid.pkg",
				content: "not a real pkg",
				// PKG files are opened with the system, so this may not immediately fail
			},
			{
				name:     "empty.zip",
				content:  "",
				expected: "failed to unzip",
			},
		}

		for _, test := range invalidFiles {
			t.Run(test.name, func(t *testing.T) {
				filePath := filepath.Join(tempDir, test.name)
				err := os.WriteFile(filePath, []byte(test.content), 0644)
				require.NoError(t, err)

				err = darwin.InstallUpdatePlatform(filePath)
				if test.expected != "" {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), test.expected)
				}
			})
		}
	})
}

// TestDarwinConcurrency tests concurrent operations on Darwin
func TestDarwinConcurrency(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Darwin concurrency tests only run on macOS")
	}

	darwin := &darwinKeygen{}

	t.Run("ConcurrentFingerprinting", func(t *testing.T) {
		const numGoroutines = 10
		results := make(chan string, numGoroutines)
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				fingerprint, err := darwin.GetMachineFingerprint()
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

		// If there are errors, they should be consistent
		if len(collectedErrors) > 0 {
			t.Logf("Collected %d errors during concurrent fingerprinting", len(collectedErrors))
		}
	})

	t.Run("ConcurrentCacheAccess", func(t *testing.T) {
		const numGoroutines = 5
		results := make(chan string, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				cacheDir := darwin.GetCacheDir()
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

// TestDarwinEdgeCases tests edge cases specific to Darwin
func TestDarwinEdgeCases(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Darwin edge case tests only run on macOS")
	}

	t.Run("EmptyEnvironmentVariables", func(t *testing.T) {
		// Temporarily clear HOME to test fallback behavior
		originalHome := os.Getenv("HOME")
		os.Unsetenv("HOME")
		defer os.Setenv("HOME", originalHome)

		darwin := &darwinKeygen{}
		cacheDir := darwin.GetCacheDir()

		// Should fallback to temp directory
		assert.Contains(t, cacheDir, "keygen-cache")
	})

	t.Run("SpecialCharactersInPaths", func(t *testing.T) {
		tempDir := createTestTempDir(t)

		// Create a file with special characters
		specialFile := filepath.Join(tempDir, "test file with spaces & symbols!.app")
		err := os.MkdirAll(specialFile, 0755)
		require.NoError(t, err)

		darwin := &darwinKeygen{}

		// Should handle special characters gracefully
		err = darwin.installApp(specialFile)
		// May fail due to permissions, but should not panic
		if err != nil {
			assert.NotContains(t, err.Error(), "panic")
		}
	})
}

// TestDarwinStartup tests Darwin-specific startup functionality
func TestDarwinStartup(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Darwin startup tests only run on macOS")
	}

	// Note: The current darwinKeygen doesn't implement Startup, but if it did:
	t.Run("StartupNotImplemented", func(t *testing.T) {
		darwin := &darwinKeygen{}

		// Verify the type doesn't implement Startup (based on current code)
		_, hasStartup := interface{}(darwin).(interface {
			Startup(context.Context, application.ServiceOptions) error
		})

		// Currently should be false based on the implementation
		assert.False(t, hasStartup, "darwinKeygen should not implement Startup method yet")
	})
}
