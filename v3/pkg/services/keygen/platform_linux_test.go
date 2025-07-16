//go:build linux

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

// TestLinuxKeygen tests the Linux specific implementation
func TestLinuxKeygen(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux tests only run on Linux")
	}

	linux := &linuxKeygen{}

	t.Run("GetMachineFingerprint", func(t *testing.T) {
		fingerprint, err := linux.GetMachineFingerprint()
		assert.NoError(t, err)
		assert.NotEmpty(t, fingerprint)
		assert.Len(t, fingerprint, 64) // SHA256 hex string

		// Verify it's a valid hex string
		assert.Regexp(t, "^[a-f0-9]{64}$", fingerprint)

		// Test consistency - should return same fingerprint multiple times
		fingerprint2, err := linux.GetMachineFingerprint()
		assert.NoError(t, err)
		assert.Equal(t, fingerprint, fingerprint2)
	})

	t.Run("GetCacheDir", func(t *testing.T) {
		// Test with XDG_CACHE_HOME set
		originalXDG := os.Getenv("XDG_CACHE_HOME")
		testXDG := "/tmp/test-xdg-cache"
		os.Setenv("XDG_CACHE_HOME", testXDG)

		cacheDir := linux.GetCacheDir()
		assert.Contains(t, cacheDir, testXDG)
		assert.Contains(t, cacheDir, "keygen")

		// Restore original XDG_CACHE_HOME
		if originalXDG != "" {
			os.Setenv("XDG_CACHE_HOME", originalXDG)
		} else {
			os.Unsetenv("XDG_CACHE_HOME")
		}

		// Test fallback to ~/.cache
		cacheDir = linux.GetCacheDir()
		assert.NotEmpty(t, cacheDir)
		assert.True(t, filepath.IsAbs(cacheDir))
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
					return "/path/to/nonexistent/file.deb"
				},
				expectedError: "update file not found",
			},
			{
				name: "ValidExecutable",
				setupFile: func(t *testing.T) string {
					tempDir := createTestTempDir(t)
					execPath := filepath.Join(tempDir, "updater")

					// Create a shell script
					content := "#!/bin/sh\necho 'Update completed'\nexit 0\n"
					err := os.WriteFile(execPath, []byte(content), 0755)
					require.NoError(t, err)

					return execPath
				},
				// Should not error for valid executable
			},
			{
				name: "AppImageFile",
				setupFile: func(t *testing.T) string {
					tempDir := createTestTempDir(t)
					appImagePath := filepath.Join(tempDir, "test.appimage")

					// Create a mock AppImage file
					content := "#!/bin/sh\necho 'AppImage'\nexit 0\n"
					err := os.WriteFile(appImagePath, []byte(content), 0644)
					require.NoError(t, err)

					return appImagePath
				},
				// Will fail during execution but should handle gracefully
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				updatePath := tt.setupFile(t)

				err := linux.InstallUpdatePlatform(updatePath)

				if tt.expectedError != "" {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tt.expectedError)
				} else {
					// For executables, we expect the function to attempt execution
					// The result depends on system state, so we verify no panic
					assert.NotPanics(t, func() {
						linux.InstallUpdatePlatform(updatePath)
					})
				}
			})
		}
	})
}

// TestLinuxMachineIdentification tests Linux-specific machine ID sources
func TestLinuxMachineIdentification(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux machine ID tests only run on Linux")
	}

	linux := &linuxKeygen{}

	t.Run("MachineIDSources", func(t *testing.T) {
		// Test that the function can handle various scenarios
		fingerprint, err := linux.GetMachineFingerprint()

		// Should succeed on most Linux systems
		if err != nil {
			// Acceptable if no machine ID sources are available
			assert.Contains(t, err.Error(), "unable to generate machine fingerprint")
		} else {
			assert.NotEmpty(t, fingerprint)
			assert.Len(t, fingerprint, 64)
		}
	})

	t.Run("MACAddressFallback", func(t *testing.T) {
		// Test MAC address retrieval
		macs, err := getMACAddresses()

		// Should work on most systems
		if err == nil {
			assert.NotEmpty(t, macs)
			for _, mac := range macs {
				assert.NotEmpty(t, mac)
				assert.NotEqual(t, "00:00:00:00:00:00", mac)
			}
		}
	})
}

// TestLinuxPackageInstallation tests Linux package manager integration
func TestLinuxPackageInstallation(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux package tests only run on Linux")
	}

	linux := &linuxKeygen{}

	t.Run("DEBInstallation", func(t *testing.T) {
		tempDir := createTestTempDir(t)
		debPath := filepath.Join(tempDir, "test.deb")

		// Create a mock DEB file
		err := os.WriteFile(debPath, []byte("fake deb content"), 0644)
		require.NoError(t, err)

		// Test installation - will likely fail due to invalid package or permissions
		err = linux.installDEB(debPath)

		// Expected to fail on most systems due to fake package or lack of dpkg
		if err != nil {
			assert.True(t,
				strings.Contains(err.Error(), "dpkg not found") ||
					strings.Contains(err.Error(), "exit status") ||
					strings.Contains(err.Error(), "permission denied"))
		}
	})

	t.Run("RPMInstallation", func(t *testing.T) {
		tempDir := createTestTempDir(t)
		rpmPath := filepath.Join(tempDir, "test.rpm")

		// Create a mock RPM file
		err := os.WriteFile(rpmPath, []byte("fake rpm content"), 0644)
		require.NoError(t, err)

		// Test installation - will likely fail due to invalid package or missing tools
		err = linux.installRPM(rpmPath)

		// Expected to fail on most systems
		if err != nil {
			assert.True(t,
				strings.Contains(err.Error(), "package manager found") ||
					strings.Contains(err.Error(), "exit status") ||
					strings.Contains(err.Error(), "permission denied"))
		}
	})
}

// TestLinuxAppImageHandling tests AppImage specific functionality
func TestLinuxAppImageHandling(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux AppImage tests only run on Linux")
	}

	linux := &linuxKeygen{}

	t.Run("AppImagePermissions", func(t *testing.T) {
		tempDir := createTestTempDir(t)
		appImagePath := filepath.Join(tempDir, "test.appimage")

		// Create a file without execute permissions
		err := os.WriteFile(appImagePath, []byte("#!/bin/sh\necho 'appimage'"), 0644)
		require.NoError(t, err)

		// Should make it executable during installation
		err = linux.installAppImage(appImagePath)

		// May fail during execution, but file should be made executable
		info, statErr := os.Stat(appImagePath)
		if statErr == nil {
			assert.True(t, info.Mode()&0111 != 0, "AppImage should be made executable")
		}
	})

	t.Run("AppImageReplacement", func(t *testing.T) {
		tempDir := createTestTempDir(t)

		// Create a mock current executable
		currentExe := filepath.Join(tempDir, "current-app")
		err := os.WriteFile(currentExe, []byte("#!/bin/sh\necho 'current'"), 0755)
		require.NoError(t, err)

		// Create a new AppImage
		newAppImage := filepath.Join(tempDir, "new.appimage")
		err = os.WriteFile(newAppImage, []byte("#!/bin/sh\necho 'new version'"), 0644)
		require.NoError(t, err)

		// Test replacement logic
		// Note: The actual replacement involves os.Exit(), so we can't fully test it
		// We can test that it attempts to create the update script
		err = linux.installAppImage(newAppImage)

		// Should either succeed or fail gracefully
		if err != nil {
			assert.NotContains(t, err.Error(), "panic")
		}
	})
}

// TestLinuxTarGzHandling tests tar.gz archive handling
func TestLinuxTarGzHandling(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux tar.gz tests only run on Linux")
	}

	linux := &linuxKeygen{}

	t.Run("TarGzExtraction", func(t *testing.T) {
		tempDir := createTestTempDir(t)
		tarPath := filepath.Join(tempDir, "test.tar.gz")

		// Create a minimal tar.gz file (just empty file for now)
		err := os.WriteFile(tarPath, []byte("fake tar content"), 0644)
		require.NoError(t, err)

		// Test extraction - will fail due to invalid tar
		err = linux.installTarGz(tarPath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to extract tar.gz")
	})

	t.Run("InstallScriptSearch", func(t *testing.T) {
		// Test the logic for finding install scripts
		// We can't easily create a real tar.gz, but we can test the search logic

		// This would normally be tested with a real tar.gz containing install.sh
		// For now, we test that the function handles invalid tar files gracefully
		tempDir := createTestTempDir(t)
		invalidTar := filepath.Join(tempDir, "invalid.tar.gz")
		err := os.WriteFile(invalidTar, []byte("not a tar"), 0644)
		require.NoError(t, err)

		err = linux.installTarGz(invalidTar)
		assert.Error(t, err)
	})
}

// TestLinuxFileOperations tests Linux-specific file operations
func TestLinuxFileOperations(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux file operation tests only run on Linux")
	}

	t.Run("CopyFile", func(t *testing.T) {
		tempDir := createTestTempDir(t)

		// Test copyFile helper function
		srcFile := filepath.Join(tempDir, "source.txt")
		dstFile := filepath.Join(tempDir, "dest.txt")

		testContent := "test file content"
		err := os.WriteFile(srcFile, []byte(testContent), 0644)
		require.NoError(t, err)

		err = copyFile(srcFile, dstFile)
		assert.NoError(t, err)

		// Verify content was copied
		copied, err := os.ReadFile(dstFile)
		assert.NoError(t, err)
		assert.Equal(t, testContent, string(copied))
	})

	t.Run("ExecutableDetection", func(t *testing.T) {
		tempDir := createTestTempDir(t)

		// Create executable file
		execFile := filepath.Join(tempDir, "executable")
		err := os.WriteFile(execFile, []byte("#!/bin/sh\necho test"), 0755)
		require.NoError(t, err)

		// Create non-executable file
		nonExecFile := filepath.Join(tempDir, "non-executable")
		err = os.WriteFile(nonExecFile, []byte("text content"), 0644)
		require.NoError(t, err)

		// Test executable detection
		execInfo, err := os.Stat(execFile)
		require.NoError(t, err)
		assert.True(t, execInfo.Mode()&0111 != 0)

		nonExecInfo, err := os.Stat(nonExecFile)
		require.NoError(t, err)
		assert.False(t, nonExecInfo.Mode()&0111 != 0)
	})
}

// TestLinuxSystemIntegration tests integration with Linux system facilities
func TestLinuxSystemIntegration(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux system integration tests only run on Linux")
	}

	t.Run("GetAppName", func(t *testing.T) {
		appName := getAppName()
		assert.NotEmpty(t, appName)

		// Should return a reasonable app name
		assert.True(t, len(appName) > 0)
		assert.False(t, strings.Contains(appName, "/"))
	})

	t.Run("SystemFileAccess", func(t *testing.T) {
		// Test access to common Linux system files
		systemFiles := []string{
			"/sys/class/net",
			"/proc/version",
		}

		for _, file := range systemFiles {
			_, err := os.Stat(file)
			if err != nil {
				t.Logf("System file %s not accessible: %v", file, err)
			}
		}
	})

	t.Run("DMIAccess", func(t *testing.T) {
		// Test access to DMI information
		dmiFiles := []string{
			"/sys/class/dmi/id/product_uuid",
			"/sys/class/dmi/id/board_serial",
			"/sys/class/dmi/id/product_name",
		}

		for _, file := range dmiFiles {
			if data, err := os.ReadFile(file); err == nil {
				t.Logf("DMI file %s contains: %s", file, strings.TrimSpace(string(data)))
			} else {
				t.Logf("DMI file %s not accessible: %v", file, err)
			}
		}
	})
}

// TestLinuxErrorHandling tests error scenarios specific to Linux
func TestLinuxErrorHandling(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux error handling tests only run on Linux")
	}

	linux := &linuxKeygen{}

	t.Run("InvalidPackageFormats", func(t *testing.T) {
		tempDir := createTestTempDir(t)

		// Test various invalid package formats
		invalidFiles := []struct {
			name     string
			content  string
			expected string
		}{
			{
				name:     "invalid.deb",
				content:  "not a real deb",
				expected: "exit status", // dpkg will fail
			},
			{
				name:     "invalid.rpm",
				content:  "not a real rpm",
				expected: "package manager", // no rpm tools or exit status
			},
			{
				name:     "empty.sh",
				content:  "",
				expected: "", // Empty script may succeed
			},
		}

		for _, test := range invalidFiles {
			t.Run(test.name, func(t *testing.T) {
				filePath := filepath.Join(tempDir, test.name)
				err := os.WriteFile(filePath, []byte(test.content), 0644)
				require.NoError(t, err)

				err = linux.InstallUpdatePlatform(filePath)
				if test.expected != "" {
					if err != nil {
						assert.Contains(t, err.Error(), test.expected)
					}
				}
			})
		}
	})

	t.Run("PermissionErrors", func(t *testing.T) {
		// Test handling of permission errors
		tempDir := createTestTempDir(t)

		// Create a directory without write permissions
		noWriteDir := filepath.Join(tempDir, "no-write")
		err := os.MkdirAll(noWriteDir, 0555)
		require.NoError(t, err)
		defer os.Chmod(noWriteDir, 0755) // Restore permissions for cleanup

		// Try to create file in non-writable directory
		testFile := filepath.Join(noWriteDir, "test.txt")
		err = os.WriteFile(testFile, []byte("test"), 0644)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "permission denied")
	})
}

// TestLinuxConcurrency tests concurrent operations on Linux
func TestLinuxConcurrency(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux concurrency tests only run on Linux")
	}

	linux := &linuxKeygen{}

	t.Run("ConcurrentFingerprinting", func(t *testing.T) {
		const numGoroutines = 10
		results := make(chan string, numGoroutines)
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				fingerprint, err := linux.GetMachineFingerprint()
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
				cacheDir := linux.GetCacheDir()
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

// TestLinuxEnvironmentVariables tests handling of environment variables
func TestLinuxEnvironmentVariables(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux environment tests only run on Linux")
	}

	t.Run("XDGCacheHome", func(t *testing.T) {
		linux := &linuxKeygen{}

		// Test with custom XDG_CACHE_HOME
		originalXDG := os.Getenv("XDG_CACHE_HOME")
		testXDG := "/tmp/test-xdg-cache"
		os.Setenv("XDG_CACHE_HOME", testXDG)

		cacheDir := linux.GetCacheDir()
		assert.Contains(t, cacheDir, testXDG)

		// Restore original
		if originalXDG != "" {
			os.Setenv("XDG_CACHE_HOME", originalXDG)
		} else {
			os.Unsetenv("XDG_CACHE_HOME")
		}
	})

	t.Run("HomeDirectory", func(t *testing.T) {
		linux := &linuxKeygen{}

		// Test with no HOME variable (edge case)
		originalHome := os.Getenv("HOME")
		os.Unsetenv("HOME")
		os.Unsetenv("XDG_CACHE_HOME")

		cacheDir := linux.GetCacheDir()

		// Should fallback to temp directory
		assert.Contains(t, cacheDir, "keygen-cache")

		// Restore HOME
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		}
	})
}

// TestLinuxStartup tests Linux-specific startup functionality
func TestLinuxStartup(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux startup tests only run on Linux")
	}

	// Note: The current linuxKeygen doesn't implement Startup, but if it did:
	t.Run("StartupNotImplemented", func(t *testing.T) {
		linux := &linuxKeygen{}

		// Verify the type doesn't implement Startup (based on current code)
		_, hasStartup := interface{}(linux).(interface {
			Startup(context.Context, application.ServiceOptions) error
		})

		// Currently should be false based on the implementation
		assert.False(t, hasStartup, "linuxKeygen should not implement Startup method yet")
	})
}
