package commands

import (
	"context"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/cache"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestLegacyDevRejectsManifestOnlyFlags(t *testing.T) {
	t.Chdir(t.TempDir())
	err := Dev(&DevOptions{Profile: "release"})
	assert.ErrorContains(t, err, "require an active wails.toml")
}

func TestDevPathExcludedAlwaysSkipsGeneratedAndDependencyTrees(t *testing.T) {
	config := manifest.Config{Build: manifest.Build{OutputDirectory: "out"}, Frontend: manifest.Frontend{Directory: "frontend"}}
	for _, path := range []string{".wails/cache", ".git/objects", "out/app", "frontend/src/main.ts", "tools/node_modules/pkg"} {
		assert.True(t, devPathExcluded(config, path), path)
	}
	assert.False(t, devPathExcluded(config, "internal/service.go"))
}

func TestMatchesDevWatchDoubleStar(t *testing.T) {
	patterns := []string{"internal/**/generated/*.go", "**/*.go", "wails.toml"}
	assert.True(t, matchesDevWatch(patterns, "internal/service.go"))
	assert.True(t, matchesDevWatch(patterns, "internal/deep/generated/service.go"))
	assert.True(t, matchesDevWatch(patterns, "main.go"))
	assert.True(t, matchesDevWatch(patterns, "wails.toml"))
	assert.False(t, matchesDevWatch(patterns, "README.md"))
	assert.False(t, matchesDevWatch([]string{"internal/**/generated/*.go"}, "other/generated/service.go"))
	assert.False(t, matchesDevWatch([]string{"[invalid"}, "main.go"))
}

func TestDevSessionHonoursGitIgnoreAndWatchesPolicyChanges(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "ignored"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".gitignore"), []byte("ignored/\n*.generated.go\n"), 0o644))
	ignored, err := loadDevGitIgnore(root, true)
	require.NoError(t, err)
	config := manifest.Config{Dev: manifest.Dev{UseGitIgnore: true, Watch: []string{"**/*.go"}}}
	assert.True(t, ignoreDevEvent(root, config, ignored, filepath.Join(root, "ignored", "service.go")))
	assert.True(t, ignoreDevEvent(root, config, ignored, filepath.Join(root, "service.generated.go")))
	assert.False(t, ignoreDevEvent(root, config, ignored, filepath.Join(root, ".gitignore")))

	disabled, err := loadDevGitIgnore(root, false)
	require.NoError(t, err)
	assert.Nil(t, disabled)
}

func TestNewDirectoryDetectionFindsOnlyWatchedInputs(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, ".gitignore"), []byte("ignored/\n"), 0o644))
	write := func(relative string) string {
		path := filepath.Join(root, relative)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte("package example\n"), 0o644))
		return filepath.Dir(path)
	}
	config := manifest.Config{Dev: manifest.Dev{UseGitIgnore: true, Watch: []string{"**/*.go"}}}
	ignored, err := loadDevGitIgnore(root, true)
	require.NoError(t, err)
	assert.True(t, shouldRefreshDevWatchesForDirectory(root, config, ignored, filepath.Join(root, "new")))
	assert.False(t, shouldRefreshDevWatchesForDirectory(root, config, ignored, filepath.Join(root, "ignored")))
	assert.False(t, shouldRefreshDevWatchesForDirectory(root, config, ignored, filepath.Join(root, ".wails", "cache")))
	assert.True(t, directoryContainsDevInput(root, config, ignored, write("new/deep/service.go")))
	assert.False(t, directoryContainsDevInput(root, config, ignored, write("docs/README.md")))
	assert.False(t, directoryContainsDevInput(root, config, ignored, write("ignored/service.go")))
	assert.True(t, isDevGitIgnoreEvent(root, config, filepath.Join(root, "nested", ".gitignore")))
}

func TestFrontendSessionChanged(t *testing.T) {
	base := manifest.Config{Frontend: manifest.Frontend{Directory: "frontend", PackageManager: "npm", DevCommand: "dev"}}
	assert.False(t, frontendSessionChanged(base, 9245, base, 9245))
	next := base
	next.Frontend.DevCommand = "serve"
	assert.True(t, frontendSessionChanged(base, 9245, next, 9245))
	assert.True(t, frontendSessionChanged(base, 9245, base, 5173))
}

func TestFrontendDevArgsPinTheLoopbackHost(t *testing.T) {
	serverArgs := []string{"--host", "127.0.0.1", "--port", "9245", "--strictPort"}
	npm, err := frontendDevArgs("npm", "dev", serverArgs)
	require.NoError(t, err)
	assert.Equal(t, []string{"run", "dev", "--", "--host", "127.0.0.1", "--port", "9245", "--strictPort"}, npm)
	yarn, err := frontendDevArgs("yarn", "dev", serverArgs)
	require.NoError(t, err)
	assert.Equal(t, []string{"dev", "--host", "127.0.0.1", "--port", "9245", "--strictPort"}, yarn)
	custom, err := frontendDevArgs("npm", "serve", nil)
	require.NoError(t, err)
	assert.Equal(t, []string{"run", "serve"}, custom)
}

func TestFrontendDevFlagsAreOnlyAddedForViteScripts(t *testing.T) {
	root := t.TempDir()
	frontend := filepath.Join(root, "frontend")
	require.NoError(t, os.MkdirAll(frontend, 0o755))
	config := manifest.Config{Frontend: manifest.Frontend{Directory: "frontend", DevCommand: "dev"}}
	require.NoError(t, os.WriteFile(filepath.Join(frontend, "package.json"), []byte(`{"scripts":{"dev":"cross-env MODE=dev vite"}}`), 0o644))
	assert.True(t, frontendDevCommandUsesVite(root, config))
	require.NoError(t, os.WriteFile(filepath.Join(frontend, "package.json"), []byte(`{"scripts":{"dev":"ng serve"}}`), 0o644))
	assert.False(t, frontendDevCommandUsesVite(root, config))
}

func TestDevWatchSessionChanged(t *testing.T) {
	base := manifest.Config{Build: manifest.Build{OutputDirectory: "bin"}, Frontend: manifest.Frontend{Directory: "frontend"}, Dev: manifest.Dev{Watch: []string{"**/*.go"}, Exclude: []string{"tmp"}, UseGitIgnore: true}}
	assert.False(t, devWatchSessionChanged(base, base))
	next := base
	next.Dev.Watch = []string{"**/*.go", "wails.toml"}
	assert.True(t, devWatchSessionChanged(base, next))
	next = base
	next.Frontend.Directory = "web"
	assert.True(t, devWatchSessionChanged(base, next))
}

func TestManifestBackendChangedOnlyForExecutedCompileOrURLChange(t *testing.T) {
	key := pipelineCompileKey("linux", "amd64")
	for _, status := range []cache.LookupStatus{cache.LookupHit, cache.LookupRestored} {
		run := manifestPipelineRun{Results: map[pipeline.NodeKey]pipeline.Result{key: {Status: status}}}
		assert.False(t, manifestBackendChanged(run, "linux", "amd64", "http://localhost:9245", "http://localhost:9245"), status)
	}
	run := manifestPipelineRun{Results: map[pipeline.NodeKey]pipeline.Result{key: {Status: cache.LookupMiss}}}
	assert.True(t, manifestBackendChanged(run, "linux", "amd64", "http://localhost:9245", "http://localhost:9245"))
	assert.True(t, manifestBackendChanged(manifestPipelineRun{}, "linux", "amd64", "http://localhost:9245", "http://localhost:9245"))
	assert.True(t, manifestBackendChanged(run, "linux", "amd64", "http://localhost:9245", "http://localhost:9246"))
	hookKey := pipeline.NodeKey("hook:after_build:linux/amd64")
	run = manifestPipelineRun{
		Plan: pipeline.Plan{Nodes: map[pipeline.NodeKey]pipeline.Node{hookKey: {Kind: pipeline.RunHook, Spec: pipeline.HookSpec{Phase: "after_build"}}}},
		Results: map[pipeline.NodeKey]pipeline.Result{
			key:     {Status: cache.LookupHit},
			hookKey: {Status: cache.LookupMiss},
		},
	}
	assert.True(t, manifestBackendChanged(run, "linux", "amd64", "http://localhost:9245", "http://localhost:9245"))
}

func TestManifestProcessesReceiveExplicitDevEnvironment(t *testing.T) {
	output := filepath.Join(t.TempDir(), "environment")
	process, err := startManifestProcess(t.TempDir(), os.Args[0], []string{
		"WAILS_READINESS_HELPER=capture-env",
		"WAILS_ENV_OUTPUT=" + output,
		"FRONTEND_DEVSERVER_URL=http://localhost:7777",
		wailsVitePort + "=7777",
	}, "-test.run=TestManifestReadinessHelper")
	require.NoError(t, err)
	select {
	case <-process.done:
	case <-time.After(2 * time.Second):
		process.stop(100 * time.Millisecond)
		t.Fatal("environment helper did not exit")
	}
	require.NoError(t, process.waitError())
	data, err := os.ReadFile(output)
	require.NoError(t, err)
	assert.Equal(t, "http://localhost:7777\n7777", string(data))
}

func TestManifestProcessReadinessWaitsForListeningFrontend(t *testing.T) {
	process := startReadinessHelper(t, "hold")
	defer process.stop(100 * time.Millisecond)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	address := listener.Addr().String()
	require.NoError(t, listener.Close())
	listenErrors := make(chan error, 1)
	go func() {
		time.Sleep(50 * time.Millisecond)
		ready, listenErr := net.Listen("tcp", address)
		listenErrors <- listenErr
		if listenErr == nil {
			defer ready.Close()
			time.Sleep(300 * time.Millisecond)
		}
	}()

	err = waitForProcessTCP(context.Background(), process, address, time.Second)
	require.NoError(t, <-listenErrors)
	require.NoError(t, err)
}

func TestManifestProcessReadinessReportsEarlyExit(t *testing.T) {
	process := startReadinessHelper(t, "exit")
	err := waitForProcessTCP(context.Background(), process, "127.0.0.1:1", time.Second)
	require.ErrorContains(t, err, "process exited before becoming ready")
}

func TestManifestProcessReadinessRequiresBackendToStayAlive(t *testing.T) {
	stable := startReadinessHelper(t, "hold")
	defer stable.stop(100 * time.Millisecond)
	require.NoError(t, waitForProcessStable(context.Background(), stable, 50*time.Millisecond))

	exited := startReadinessHelper(t, "exit")
	err := waitForProcessStable(context.Background(), exited, 200*time.Millisecond)
	require.ErrorContains(t, err, "process exited during startup")
}

func TestManifestProcessReadinessHonoursCancellation(t *testing.T) {
	process := startReadinessHelper(t, "hold")
	defer process.stop(100 * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	started := time.Now()
	err := waitForProcessTCP(ctx, process, "127.0.0.1:1", time.Second)
	require.ErrorIs(t, err, context.Canceled)
	assert.Less(t, time.Since(started), 250*time.Millisecond)

	err = waitForProcessStable(ctx, process, time.Second)
	require.ErrorIs(t, err, context.Canceled)
}

func TestManifestProcessStopIsIdempotent(t *testing.T) {
	process := startReadinessHelper(t, "hold")
	process.stop(100 * time.Millisecond)
	process.stop(100 * time.Millisecond)
	select {
	case <-process.done:
		assert.Error(t, process.waitError())
	default:
		t.Fatal("process was not reaped")
	}
}

func startReadinessHelper(t *testing.T, mode string) *manifestProcess {
	t.Helper()
	command := exec.Command(os.Args[0], "-test.run=TestManifestReadinessHelper")
	command.Env = append(os.Environ(), "WAILS_READINESS_HELPER="+mode)
	configureManifestProcess(command)
	require.NoError(t, command.Start())
	process := &manifestProcess{cmd: command, done: make(chan struct{})}
	go func() {
		process.mu.Lock()
		process.err = command.Wait()
		process.mu.Unlock()
		close(process.done)
	}()
	return process
}

func TestManifestReadinessHelper(t *testing.T) {
	switch os.Getenv("WAILS_READINESS_HELPER") {
	case "hold":
		time.Sleep(5 * time.Minute)
	case "exit":
		os.Exit(3)
	case "capture-env":
		data := []byte(os.Getenv("FRONTEND_DEVSERVER_URL") + "\n" + os.Getenv(wailsVitePort))
		if err := os.WriteFile(os.Getenv("WAILS_ENV_OUTPUT"), data, 0o600); err != nil {
			os.Exit(4)
		}
	}
}
