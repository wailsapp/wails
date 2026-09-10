package commands

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestRunApplicationBuildsTaggedGoOnlyProjectAndUsesLaunchDefaults(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/runfixture\n\ngo 1.25\n"), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte(`package main
import ("os"; "encoding/json")
func main() {
 wd,_:=os.Getwd()
 data,_:=json.Marshal([]any{os.Args[1:],os.Getenv("RUN_VALUE"),wd,selected})
 if err:=os.WriteFile("result.json",data,0600);err!=nil { panic(err) }
 if os.Getenv("EXIT_TEST") == "yes" { os.Exit(42) }
}`), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(root, "selected.go"), []byte("//go:build run_fixture && production\n\npackage main\nconst selected = true\n"), 0600))
	source := `version = 3
project {
 name = "runfixture"
 product_name = "Run fixture"
 identifier = "com.example.runfixture"
 version = "1.0.0"
}
frontend { disabled = true }
run {
 tags = ["run_fixture"]
 args = ["default", "space value"]
 environment = { RUN_VALUE = "configured" }
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, "wails.hcl"), []byte(source), 0600))
	t.Setenv("WAILS_EXP_USE_WAKE", "1")
	t.Setenv("GOWORK", "off")
	t.Setenv("RUN_VALUE", "inherited")
	t.Chdir(root)
	require.NoError(t, RunApplication(&RunOptions{}, nil))
	data, err := os.ReadFile("result.json")
	require.NoError(t, err)
	var result []any
	require.NoError(t, json.Unmarshal(data, &result))
	require.Equal(t, []any{[]any{"default", "space value"}, "configured", root, true}, result)
	binary := filepath.Join(root, "bin", "runfixture")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	if runtime.GOOS == "darwin" {
		binary = filepath.Join(root, "bin", "runfixture.app", "Contents", "MacOS", "runfixture")
	}
	before, err := os.Stat(binary)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(root, "wails.hcl"), []byte(strings.ReplaceAll(source, "configured", "changed")), 0600))
	require.NoError(t, RunApplication(&RunOptions{}, []string{}))
	after, err := os.Stat(binary)
	require.NoError(t, err)
	require.Equal(t, before.ModTime(), after.ModTime(), "runtime-only change must preserve cached executable")
	data, err = os.ReadFile("result.json")
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(data, &result))
	require.Equal(t, []any{[]any{}, "changed", root, true}, result)
	t.Setenv("EXIT_TEST", "yes")
	err = RunApplication(&RunOptions{}, nil)
	code, ok := ApplicationExitCode(err)
	require.True(t, ok)
	require.Equal(t, 42, code)
}

func TestRunApplicationFallbackAndInvalidManifest(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	t.Setenv("GOWORK", "off")
	require.NoError(t, os.WriteFile("go.mod", []byte("module example.com/fallback\n\ngo 1.25\n"), 0600))
	require.NoError(t, os.WriteFile("main.go", []byte(`package main
import "os"
func main(){ if len(os.Args)!=2 || os.Args[1]!="--app-flag" { panic(os.Args) }; if err:=os.WriteFile("ran",[]byte("yes"),0600);err!=nil{panic(err)} }
`), 0600))
	require.NoError(t, RunApplication(&RunOptions{}, []string{"--app-flag"}))
	require.FileExists(t, "ran")
	require.NoError(t, os.Remove("ran"))
	require.NoError(t, os.WriteFile("wails.hcl", []byte("invalid"), 0600))
	t.Setenv("WAILS_EXP_USE_WAKE", "1")
	require.Error(t, RunApplication(&RunOptions{}, []string{"--app-flag"}))
	require.NoFileExists(t, "ran")
}

func TestRunProcessCancellation(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix process cancellation fixture")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	start := time.Now()
	err := runApplicationProcess(ctx, t.TempDir(), "sh", []string{"-c", "sleep 30 & wait"}, nil)
	require.True(t, errors.Is(err, context.DeadlineExceeded), "%v", err)
	require.Less(t, time.Since(start), 5*time.Second)
}

func TestRunMobileRejectsUnsupportedLaunchSettingsBeforeDeployment(t *testing.T) {
	for _, settings := range []manifest.Run{{Args: []string{"--flag"}}, {Environment: map[string]string{"VALUE": "x"}}} {
		err := configureRunTarget(&RunOptions{}, "android", "arm64", settings, &manifestRunOptions{})
		require.ErrorContains(t, err, "does not support")
	}
	require.Error(t, configureRunTarget(&RunOptions{Device: "device", Emulator: "emulator"}, "android", "arm64", manifest.Run{}, &manifestRunOptions{}))
	require.Error(t, configureRunTarget(&RunOptions{Destination: "invalid", Plan: true}, "ios", "arm64", manifest.Run{}, &manifestRunOptions{}))
	require.Error(t, configureRunTarget(&RunOptions{Device: "device"}, "linux", "amd64", manifest.Run{}, &manifestRunOptions{}))
}

func TestRunBundleRegistrationPreservesBundlePathAndReportsFailure(t *testing.T) {
	root := t.TempDir()
	bundle := filepath.Join(root, "bin", "My App.app")
	executable := filepath.Join(bundle, "Contents", "MacOS", "my-app")
	failure := errors.New("registration failed")
	err := registerRunBundle(context.Background(), root, executable, func(_ context.Context, directory, tool string, args ...string) (string, error) {
		require.Equal(t, root, directory)
		require.Equal(t, "lsregister", filepath.Base(tool))
		require.Equal(t, []string{"-f", bundle}, args)
		return "", failure
	})
	require.ErrorIs(t, err, failure)
}

func TestIOSDeviceRunEnvironmentUsesJSONAndPreservesArguments(t *testing.T) {
	settings := manifest.Run{Environment: map[string]string{"VALUE": "space \"quote\" $literal", "EMPTY": ""}, Args: []string{"--help", "space value"}}
	args := iosDeviceRunCommand("device-id", "com.example.app", settings)
	index := -1
	for i, arg := range args {
		if arg == "--environment-variables" {
			index = i
			break
		}
	}
	require.NotEqual(t, -1, index)
	var decoded map[string]string
	require.NoError(t, json.Unmarshal([]byte(args[index+1]), &decoded))
	require.Equal(t, settings.Environment, decoded)
	require.Equal(t, []string{"com.example.app", "--help", "space value"}, args[index+2:])
	require.NotContains(t, iosDeviceRunCommand("device-id", "com.example.app", manifest.Run{}), "--environment-variables")
}
