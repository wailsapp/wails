package commands

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wailsapp/wails/v3/internal/flags"
)

func TestWrapTask(t *testing.T) {
	// Get current platform info for expected values
	currentOS := runtime.GOOS
	currentArch := runtime.GOARCH

	tests := []struct {
		name             string
		command          string
		otherArgs        []string
		envGOOS          string
		envGOARCH        string
		expectedTaskName string
		expectedArgs     []string
		expectedOsArgs   []string
	}{
		{
			name:             "Build with parameters uses current platform",
			command:          "build",
			otherArgs:        []string{"CONFIG=debug"},
			expectedTaskName: currentOS + ":build",
			expectedArgs:     []string{"CONFIG=debug", "ARCH=" + currentArch},
			expectedOsArgs:   []string{"wails3", "task", currentOS + ":build", "CONFIG=debug", "ARCH=" + currentArch},
		},
		{
			name:             "Package with parameters uses current platform",
			command:          "package",
			otherArgs:        []string{"VERSION=1.0.0", "OUTPUT=app.pkg"},
			expectedTaskName: currentOS + ":package",
			expectedArgs:     []string{"VERSION=1.0.0", "OUTPUT=app.pkg", "ARCH=" + currentArch},
			expectedOsArgs:   []string{"wails3", "task", currentOS + ":package", "VERSION=1.0.0", "OUTPUT=app.pkg", "ARCH=" + currentArch},
		},
		{
			name:             "Build without parameters",
			command:          "build",
			otherArgs:        []string{},
			expectedTaskName: currentOS + ":build",
			expectedArgs:     []string{"ARCH=" + currentArch},
			expectedOsArgs:   []string{"wails3", "task", currentOS + ":build", "ARCH=" + currentArch},
		},
		{
			name:             "GOOS override changes task prefix",
			command:          "build",
			otherArgs:        []string{"GOOS=darwin", "CONFIG=release"},
			expectedTaskName: "darwin:build",
			expectedArgs:     []string{"CONFIG=release", "ARCH=" + currentArch},
			expectedOsArgs:   []string{"wails3", "task", "darwin:build", "CONFIG=release", "ARCH=" + currentArch},
		},
		{
			name:             "GOARCH override changes ARCH arg",
			command:          "build",
			otherArgs:        []string{"GOARCH=arm64"},
			expectedTaskName: currentOS + ":build",
			expectedArgs:     []string{"ARCH=arm64"},
			expectedOsArgs:   []string{"wails3", "task", currentOS + ":build", "ARCH=arm64"},
		},
		{
			name:             "Both GOOS and GOARCH override",
			command:          "package",
			otherArgs:        []string{"GOOS=windows", "GOARCH=386", "VERSION=2.0"},
			expectedTaskName: "windows:package",
			expectedArgs:     []string{"VERSION=2.0", "ARCH=386"},
			expectedOsArgs:   []string{"wails3", "task", "windows:package", "VERSION=2.0", "ARCH=386"},
		},
		{
			name:             "Environment GOOS is used when no arg override",
			command:          "build",
			otherArgs:        []string{"CONFIG=debug"},
			envGOOS:          "darwin",
			expectedTaskName: "darwin:build",
			expectedArgs:     []string{"CONFIG=debug", "ARCH=" + currentArch},
			expectedOsArgs:   []string{"wails3", "task", "darwin:build", "CONFIG=debug", "ARCH=" + currentArch},
		},
		{
			name:             "Arg GOOS overrides environment GOOS",
			command:          "build",
			otherArgs:        []string{"GOOS=linux"},
			envGOOS:          "darwin",
			expectedTaskName: "linux:build",
			expectedArgs:     []string{"ARCH=" + currentArch},
			expectedOsArgs:   []string{"wails3", "task", "linux:build", "ARCH=" + currentArch},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save and restore os.Args
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			// Save and restore environment variables
			originalGOOS := os.Getenv("GOOS")
			originalGOARCH := os.Getenv("GOARCH")
			defer func() {
				if originalGOOS == "" {
					os.Unsetenv("GOOS")
				} else {
					os.Setenv("GOOS", originalGOOS)
				}
				if originalGOARCH == "" {
					os.Unsetenv("GOARCH")
				} else {
					os.Setenv("GOARCH", originalGOARCH)
				}
			}()

			// Set test environment
			if tt.envGOOS != "" {
				os.Setenv("GOOS", tt.envGOOS)
			} else {
				os.Unsetenv("GOOS")
			}
			if tt.envGOARCH != "" {
				os.Setenv("GOARCH", tt.envGOARCH)
			} else {
				os.Unsetenv("GOARCH")
			}

			// Mock RunTask to capture the arguments
			originalRunTask := runTaskFunc
			var capturedOptions *RunTaskOptions
			var capturedOtherArgs []string
			runTaskFunc = func(options *RunTaskOptions, otherArgs []string) error {
				capturedOptions = options
				capturedOtherArgs = otherArgs
				return nil
			}
			defer func() { runTaskFunc = originalRunTask }()

			// Execute wrapTask
			err := wrapTask(tt.command, tt.otherArgs)
			assert.NoError(t, err)

			// Verify os.Args was set correctly
			assert.Equal(t, tt.expectedOsArgs, os.Args)

			// Verify RunTask was called with correct parameters
			assert.Equal(t, tt.expectedTaskName, capturedOptions.Name)
			assert.Equal(t, tt.expectedArgs, capturedOtherArgs)
		})
	}
}

func TestBuildCommand(t *testing.T) {
	currentOS := runtime.GOOS
	currentArch := runtime.GOARCH

	// Save original RunTask
	originalRunTask := runTaskFunc
	defer func() { runTaskFunc = originalRunTask }()

	// Mock RunTask to capture the arguments
	var capturedOptions *RunTaskOptions
	var capturedOtherArgs []string
	runTaskFunc = func(options *RunTaskOptions, otherArgs []string) error {
		capturedOptions = options
		capturedOtherArgs = otherArgs
		return nil
	}

	// Save original os.Args and environment
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	originalGOOS := os.Getenv("GOOS")
	originalGOARCH := os.Getenv("GOARCH")
	defer func() {
		if originalGOOS == "" {
			os.Unsetenv("GOOS")
		} else {
			os.Setenv("GOOS", originalGOOS)
		}
		if originalGOARCH == "" {
			os.Unsetenv("GOARCH")
		} else {
			os.Setenv("GOARCH", originalGOARCH)
		}
	}()
	os.Unsetenv("GOOS")
	os.Unsetenv("GOARCH")

	// Test Build command
	buildFlags := &flags.Build{}
	otherArgs := []string{"CONFIG=release"}

	err := Build(buildFlags, otherArgs)
	assert.NoError(t, err)
	assert.Equal(t, currentOS+":build", capturedOptions.Name)
	assert.Equal(t, []string{"CONFIG=release", "ARCH=" + currentArch}, capturedOtherArgs)
}

func TestPackageCommand(t *testing.T) {
	currentOS := runtime.GOOS
	currentArch := runtime.GOARCH

	// Save original RunTask
	originalRunTask := runTaskFunc
	defer func() { runTaskFunc = originalRunTask }()

	// Mock RunTask to capture the arguments
	var capturedOptions *RunTaskOptions
	var capturedOtherArgs []string
	runTaskFunc = func(options *RunTaskOptions, otherArgs []string) error {
		capturedOptions = options
		capturedOtherArgs = otherArgs
		return nil
	}

	// Save original os.Args and environment
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	originalGOOS := os.Getenv("GOOS")
	originalGOARCH := os.Getenv("GOARCH")
	defer func() {
		if originalGOOS == "" {
			os.Unsetenv("GOOS")
		} else {
			os.Setenv("GOOS", originalGOOS)
		}
		if originalGOARCH == "" {
			os.Unsetenv("GOARCH")
		} else {
			os.Setenv("GOARCH", originalGOARCH)
		}
	}()
	os.Unsetenv("GOOS")
	os.Unsetenv("GOARCH")

	// Test Package command
	packageFlags := &flags.Package{}
	otherArgs := []string{"VERSION=2.0.0", "OUTPUT=myapp.dmg"}

	err := Package(packageFlags, otherArgs)
	assert.NoError(t, err)
	assert.Equal(t, currentOS+":package", capturedOptions.Name)
	assert.Equal(t, []string{"VERSION=2.0.0", "OUTPUT=myapp.dmg", "ARCH=" + currentArch}, capturedOtherArgs)
}

func TestSignWrapperCommand(t *testing.T) {
	currentOS := runtime.GOOS
	currentArch := runtime.GOARCH

	// Save original RunTask
	originalRunTask := runTaskFunc
	defer func() { runTaskFunc = originalRunTask }()

	// Mock RunTask to capture the arguments
	var capturedOptions *RunTaskOptions
	var capturedOtherArgs []string
	runTaskFunc = func(options *RunTaskOptions, otherArgs []string) error {
		capturedOptions = options
		capturedOtherArgs = otherArgs
		return nil
	}

	// Save original os.Args and environment
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	originalGOOS := os.Getenv("GOOS")
	originalGOARCH := os.Getenv("GOARCH")
	defer func() {
		if originalGOOS == "" {
			os.Unsetenv("GOOS")
		} else {
			os.Setenv("GOOS", originalGOOS)
		}
		if originalGOARCH == "" {
			os.Unsetenv("GOARCH")
		} else {
			os.Setenv("GOARCH", originalGOARCH)
		}
	}()
	os.Unsetenv("GOOS")
	os.Unsetenv("GOARCH")

	// Test SignWrapper command
	signFlags := &flags.SignWrapper{}
	otherArgs := []string{"IDENTITY=Developer ID"}

	err := SignWrapper(signFlags, otherArgs)
	assert.NoError(t, err)
	assert.Equal(t, currentOS+":sign", capturedOptions.Name)
	assert.Equal(t, []string{"IDENTITY=Developer ID", "ARCH=" + currentArch}, capturedOtherArgs)
}
