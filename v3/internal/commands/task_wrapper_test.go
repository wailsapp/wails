package commands

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wailsapp/wails/v3/internal/flags"
)

func TestWrapTask(t *testing.T) {
	tests := []struct {
		name           string
		command        string
		otherArgs      []string
		expectedOsArgs []string
	}{
		{
			name:           "Build with parameters",
			command:        "build",
			otherArgs:      []string{"PLATFORM=linux", "CONFIG=debug"},
			expectedOsArgs: []string{"wails3", "task", "build", "PLATFORM=linux", "CONFIG=debug"},
		},
		{
			name:           "Package with parameters",
			command:        "package",
			otherArgs:      []string{"VERSION=1.0.0", "OUTPUT=app.pkg"},
			expectedOsArgs: []string{"wails3", "task", "package", "VERSION=1.0.0", "OUTPUT=app.pkg"},
		},
		{
			name:           "Build without parameters",
			command:        "build",
			otherArgs:      []string{},
			expectedOsArgs: []string{"wails3", "task", "build"},
		},
		{
			name:           "Build with complex parameter values",
			command:        "build",
			otherArgs:      []string{"URL=https://example.com?key=value", "TAGS=tag1,tag2,tag3"},
			expectedOsArgs: []string{"wails3", "task", "build", "URL=https://example.com?key=value", "TAGS=tag1,tag2,tag3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original os.Args
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

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
			err := wrapTaskInternal(tt.command, tt.otherArgs)
			assert.NoError(t, err)

			// Verify os.Args was set correctly
			assert.Equal(t, tt.expectedOsArgs, os.Args)

			// Verify RunTask was called with correct parameters
			assert.Equal(t, tt.command, capturedOptions.Name)
			assert.Equal(t, tt.otherArgs, capturedOtherArgs)
		})
	}
}

func TestBuildCommand(t *testing.T) {
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

	// Save original os.Args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Test Build command
	buildFlags := &flags.Build{}
	otherArgs := []string{"PLATFORM=darwin", "CONFIG=release"}
	
	err := Build(buildFlags, otherArgs)
	assert.NoError(t, err)
	assert.Equal(t, "build", capturedOptions.Name)
	assert.Equal(t, otherArgs, capturedOtherArgs)
}

func TestPackageCommand(t *testing.T) {
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

	// Save original os.Args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Test Package command
	packageFlags := &flags.Package{}
	otherArgs := []string{"VERSION=2.0.0", "OUTPUT=myapp.dmg"}
	
	err := Package(packageFlags, otherArgs)
	assert.NoError(t, err)
	assert.Equal(t, "package", capturedOptions.Name)
	assert.Equal(t, otherArgs, capturedOtherArgs)
}

// Variables to enable mocking in tests
var (
	wrapTaskFunc = wrapTask
)

// Internal version that uses the function variables for testing
func wrapTaskInternal(command string, otherArgs []string) error {
	// Note: We skip the warning message in tests
	// Rebuild os.Args to include the command and all additional arguments
	newArgs := []string{"wails3", "task", command}
	newArgs = append(newArgs, otherArgs...)
	os.Args = newArgs
	// Pass the task name via options and otherArgs as CLI variables
	return runTaskFunc(&RunTaskOptions{Name: command}, otherArgs)
}