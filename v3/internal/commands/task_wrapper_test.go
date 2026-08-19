package commands

import (
	"bytes"
	"errors"
	"os"
	"runtime"
	"testing"

	"github.com/pterm/pterm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/flags"
)

func TestBuildManifestRoutingContractAndFailures(t *testing.T) {
	want := errors.New("injected failure")
	operations := manifestCommandOperations{
		active: func() (bool, error) { return true, nil },
		plan:   func(manifestRunOptions, bool) error { return nil },
		run:    func(manifestRunOptions) error { return nil },
		task:   func(string, []string) error { return nil },
	}

	activeFailure := operations
	activeFailure.active = func() (bool, error) { return false, want }
	require.ErrorIs(t, buildWithOperations(&flags.Build{}, nil, activeFailure), want)

	require.Error(t, buildWithOperations(&flags.Build{Targets: "invalid"}, nil, operations))
	require.Error(t, buildWithOperations(&flags.Build{Profile: "release"}, []string{"debug"}, operations))
	require.Error(t, buildWithOperations(&flags.Build{Profile: "release", Targets: "linux/amd64"}, nil, operations))

	var planned manifestRunOptions
	planFailure := operations
	planFailure.plan = func(options manifestRunOptions, asJSON bool) error {
		planned = options
		assert.True(t, asJSON)
		return want
	}
	require.ErrorIs(t, buildWithOperations(&flags.Build{Plan: true, JSON: true, Formats: "deb", Targets: "linux/amd64"}, nil, planFailure), want)
	assert.Equal(t, "package", planned.Verb)

	var ran manifestRunOptions
	runFailure := operations
	runFailure.run = func(options manifestRunOptions) error {
		ran = options
		return want
	}
	require.ErrorIs(t, buildWithOperations(&flags.Build{Targets: "linux/amd64", Force: true}, nil, runFailure), want)
	assert.Equal(t, "build", ran.Verb)
	assert.True(t, ran.Force)

	var taskArgs []string
	taskFailure := operations
	taskFailure.active = func() (bool, error) { return false, nil }
	taskFailure.task = func(command string, args []string) error {
		assert.Equal(t, "build", command)
		taskArgs = append([]string(nil), args...)
		return want
	}
	require.ErrorIs(t, buildWithOperations(&flags.Build{Tags: "sqlite", Obfuscated: true, GarbleArgs: "-tiny"}, nil, taskFailure), want)
	assert.Equal(t, []string{"EXTRA_TAGS=sqlite", "OBFUSCATED=true", "GARBLE_ARGS=-tiny"}, taskArgs)
}

func TestManifestProfileRejectsTaskVariablesAndMultiplePositionals(t *testing.T) {
	_, err := manifestProfile("", []string{"CONFIG=release"})
	require.ErrorContains(t, err, "do not accept Task variables")
	_, err = manifestProfile("", []string{"release", "debug"})
	require.ErrorContains(t, err, "accept at most one profile")
}

func TestDeprecatedPackageAndSignManifestRoutingContracts(t *testing.T) {
	want := errors.New("injected failure")
	base := manifestCommandOperations{
		active: func() (bool, error) { return true, nil },
		run:    func(manifestRunOptions) error { return nil },
		task:   func(string, []string) error { return nil },
	}

	activeFailure := base
	activeFailure.active = func() (bool, error) { return false, want }
	require.ErrorIs(t, packageWithOperations(&flags.Package{}, nil, activeFailure), want)
	require.ErrorIs(t, signWithOperations(&flags.SignWrapper{}, nil, activeFailure), want)

	require.Error(t, packageWithOperations(&flags.Package{}, []string{"CONFIG=release"}, base))
	require.Error(t, signWithOperations(&flags.SignWrapper{}, []string{"CONFIG=release"}, base))
	require.Error(t, packageWithOperations(&flags.Package{Targets: "invalid"}, nil, base))
	require.Error(t, signWithOperations(&flags.SignWrapper{Targets: "invalid"}, nil, base))

	var requests []manifestRunOptions
	runFailure := base
	runFailure.run = func(options manifestRunOptions) error {
		requests = append(requests, options)
		return want
	}
	require.ErrorIs(t, packageWithOperations(&flags.Package{Profile: "release", Targets: "linux/amd64", Formats: "deb", Force: true}, nil, runFailure), want)
	require.ErrorIs(t, signWithOperations(&flags.SignWrapper{Profile: "release", Targets: "darwin/arm64", Formats: "app"}, nil, runFailure), want)
	require.Len(t, requests, 2)
	assert.Equal(t, "package", requests[0].Verb)
	assert.True(t, requests[0].Force)
	assert.Equal(t, "sign", requests[1].Verb)

	inactive := base
	inactive.active = func() (bool, error) { return false, nil }
	inactive.task = func(command string, args []string) error {
		assert.Empty(t, args)
		return want
	}
	require.ErrorIs(t, packageWithOperations(&flags.Package{}, nil, inactive), want)
	require.ErrorIs(t, signWithOperations(&flags.SignWrapper{}, nil, inactive), want)
}

func TestWrapTask(t *testing.T) {
	// Get current platform info for expected values
	currentOS := runtime.GOOS
	currentArch := runtime.GOARCH

	// foreignOS is a GOOS guaranteed to differ from the host, so cross-compile
	// cases are exercised regardless of which platform the test runs on.
	foreignOS := "windows"
	if currentOS == "windows" {
		foreignOS = "linux"
	}

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
			// Host-OS build runs the root Taskfile's `build` task (not the
			// platform-prefixed one) so root customisations are honoured. The
			// root task dispatches via the GOOS variable we pass through (#5615).
			name:             "Host build runs root build task",
			command:          "build",
			otherArgs:        []string{"CONFIG=debug"},
			expectedTaskName: "build",
			expectedArgs:     []string{"CONFIG=debug", "GOOS=" + currentOS, "ARCH=" + currentArch},
			expectedOsArgs:   []string{"wails3", "task", "build", "CONFIG=debug", "GOOS=" + currentOS, "ARCH=" + currentArch},
		},
		{
			name:             "Host package runs root package task",
			command:          "package",
			otherArgs:        []string{"VERSION=1.0.0", "OUTPUT=app.pkg"},
			expectedTaskName: "package",
			expectedArgs:     []string{"VERSION=1.0.0", "OUTPUT=app.pkg", "GOOS=" + currentOS, "ARCH=" + currentArch},
			expectedOsArgs:   []string{"wails3", "task", "package", "VERSION=1.0.0", "OUTPUT=app.pkg", "GOOS=" + currentOS, "ARCH=" + currentArch},
		},
		{
			name:             "Host build without parameters runs root build task",
			command:          "build",
			otherArgs:        []string{},
			expectedTaskName: "build",
			expectedArgs:     []string{"GOOS=" + currentOS, "ARCH=" + currentArch},
			expectedOsArgs:   []string{"wails3", "task", "build", "GOOS=" + currentOS, "ARCH=" + currentArch},
		},
		{
			// sign has no root dispatch task, so it always targets the platform.
			name:             "Sign always targets platform task",
			command:          "sign",
			otherArgs:        []string{"IDENTITY=Developer ID"},
			expectedTaskName: currentOS + ":sign",
			expectedArgs:     []string{"IDENTITY=Developer ID", "GOOS=" + currentOS, "ARCH=" + currentArch},
			expectedOsArgs:   []string{"wails3", "task", currentOS + ":sign", "IDENTITY=Developer ID", "GOOS=" + currentOS, "ARCH=" + currentArch},
		},
		{
			// Cross-OS build still runs the root `build` task; the GOOS variable
			// carries the target so the root Taskfile dispatches to it.
			name:             "Cross-OS GOOS override runs root build task with GOOS var",
			command:          "build",
			otherArgs:        []string{"GOOS=" + foreignOS, "CONFIG=release"},
			expectedTaskName: "build",
			expectedArgs:     []string{"CONFIG=release", "GOOS=" + foreignOS, "ARCH=" + currentArch},
			expectedOsArgs:   []string{"wails3", "task", "build", "CONFIG=release", "GOOS=" + foreignOS, "ARCH=" + currentArch},
		},
		{
			// GOARCH alone (cross-arch, same OS) runs the root task; the root
			// dispatch passes GOOS/ARCH through to the platform build.
			name:             "GOARCH-only override runs root build task",
			command:          "build",
			otherArgs:        []string{"GOARCH=arm64"},
			expectedTaskName: "build",
			expectedArgs:     []string{"GOOS=" + currentOS, "ARCH=arm64"},
			expectedOsArgs:   []string{"wails3", "task", "build", "GOOS=" + currentOS, "ARCH=arm64"},
		},
		{
			name:             "Cross-OS GOOS and GOARCH override",
			command:          "package",
			otherArgs:        []string{"GOOS=" + foreignOS, "GOARCH=386", "VERSION=2.0"},
			expectedTaskName: "package",
			expectedArgs:     []string{"VERSION=2.0", "GOOS=" + foreignOS, "ARCH=386"},
			expectedOsArgs:   []string{"wails3", "task", "package", "VERSION=2.0", "GOOS=" + foreignOS, "ARCH=386"},
		},
		{
			name:             "Environment GOOS (cross) is used when no arg override",
			command:          "build",
			otherArgs:        []string{"CONFIG=debug"},
			envGOOS:          foreignOS,
			expectedTaskName: "build",
			expectedArgs:     []string{"CONFIG=debug", "GOOS=" + foreignOS, "ARCH=" + currentArch},
			expectedOsArgs:   []string{"wails3", "task", "build", "CONFIG=debug", "GOOS=" + foreignOS, "ARCH=" + currentArch},
		},
		{
			// Arg GOOS takes precedence over the GOOS environment variable; the
			// passed-through GOOS var reflects the arg value.
			name:             "Arg GOOS overrides environment GOOS",
			command:          "build",
			otherArgs:        []string{"GOOS=" + foreignOS},
			envGOOS:          "darwin",
			expectedTaskName: "build",
			expectedArgs:     []string{"GOOS=" + foreignOS, "ARCH=" + currentArch},
			expectedOsArgs:   []string{"wails3", "task", "build", "GOOS=" + foreignOS, "ARCH=" + currentArch},
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
	assert.Equal(t, "build", capturedOptions.Name)
	assert.Equal(t, []string{"CONFIG=release", "GOOS=" + currentOS, "ARCH=" + currentArch}, capturedOtherArgs)
}

func TestBuildCommandWithTags(t *testing.T) {
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

	// Test Build command with tags
	buildFlags := &flags.Build{}
	buildFlags.Tags = "gtk4"
	otherArgs := []string{"CONFIG=release"}

	err := Build(buildFlags, otherArgs)
	assert.NoError(t, err)
	assert.Equal(t, "build", capturedOptions.Name)
	assert.Equal(t, []string{"CONFIG=release", "EXTRA_TAGS=gtk4", "GOOS=" + currentOS, "ARCH=" + currentArch}, capturedOtherArgs)
}

func TestBuildCommandWithMultipleTags(t *testing.T) {
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

	// Test Build command with multiple comma-separated tags
	buildFlags := &flags.Build{}
	buildFlags.Tags = "gtk4,server"

	err := Build(buildFlags, nil)
	assert.NoError(t, err)
	assert.Equal(t, "build", capturedOptions.Name)
	assert.Equal(t, []string{"EXTRA_TAGS=gtk4,server", "GOOS=" + currentOS, "ARCH=" + currentArch}, capturedOtherArgs)
}

func TestBuildCommandWithObfuscation(t *testing.T) {
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

	buildFlags := &flags.Build{
		Tags:       "gtk4",
		Obfuscated: true,
		GarbleArgs: "-literals -tiny",
	}

	err := Build(buildFlags, nil)
	assert.NoError(t, err)
	assert.Equal(t, "build", capturedOptions.Name)
	assert.Equal(t, []string{"EXTRA_TAGS=gtk4", "OBFUSCATED=true", "GARBLE_ARGS=-literals -tiny", "GOOS=" + currentOS, "ARCH=" + currentArch}, capturedOtherArgs)
}

func TestBuildCommandWithoutTags(t *testing.T) {
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

	// Test Build command without tags - no EXTRA_TAGS should be present
	buildFlags := &flags.Build{}

	err := Build(buildFlags, nil)
	assert.NoError(t, err)
	assert.Equal(t, "build", capturedOptions.Name)
	assert.Equal(t, []string{"GOOS=" + currentOS, "ARCH=" + currentArch}, capturedOtherArgs)
}

func TestManifestProfileRejectsDuplicateFlagAndPositionalValue(t *testing.T) {
	profile, err := manifestProfile("release", []string{"debug"})
	assert.Empty(t, profile)
	assert.ErrorContains(t, err, "either positionally or with --profile")
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

	var output bytes.Buffer
	previousWarningOutput := pterm.Warning.Writer
	pterm.Warning.Writer = &output
	pterm.SetDefaultOutput(&output)
	t.Cleanup(func() {
		pterm.Warning.Writer = previousWarningOutput
		pterm.SetDefaultOutput(os.Stdout)
	})
	err := Package(packageFlags, otherArgs)
	assert.NoError(t, err)
	assert.Contains(t, output.String(), "wails3 package is deprecated")
	assert.Equal(t, "package", capturedOptions.Name)
	assert.Equal(t, []string{"VERSION=2.0.0", "OUTPUT=myapp.dmg", "GOOS=" + currentOS, "ARCH=" + currentArch}, capturedOtherArgs)
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

	var output bytes.Buffer
	previousWarningOutput := pterm.Warning.Writer
	pterm.Warning.Writer = &output
	pterm.SetDefaultOutput(&output)
	t.Cleanup(func() {
		pterm.Warning.Writer = previousWarningOutput
		pterm.SetDefaultOutput(os.Stdout)
	})
	err := SignWrapper(signFlags, otherArgs)
	assert.NoError(t, err)
	assert.Contains(t, output.String(), "wails3 sign is deprecated")
	assert.Equal(t, currentOS+":sign", capturedOptions.Name)
	assert.Equal(t, []string{"IDENTITY=Developer ID", "GOOS=" + currentOS, "ARCH=" + currentArch}, capturedOtherArgs)
}
