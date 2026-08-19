//go:build !windows

package commands

import (
	"context"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestManifestDevSessionOwnsIncrementalLifecycle(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	require.NoError(t, os.MkdirAll(filepath.Join(root, "frontend"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "frontend", "package.json"), []byte("{}\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "frontend", "src.js"), []byte("export const value = 1;\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".gitignore"), []byte("ignored/\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/devsession\n\ngo 1.24\n"), 0o644))
	originalMain := `package main

import (
  "os"
  "os/signal"
  "strconv"
  "syscall"
)

func main() {
  _ = (backendState{}).generation()
  if log := os.Getenv("WAILS_TEST_BACKEND_LOG"); log != "" {
    file, _ := os.OpenFile(log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
    _, _ = file.WriteString(strconv.Itoa(os.Getpid()) + "\n")
    _ = file.Close()
  }
  stopped := make(chan os.Signal, 1)
  signal.Notify(stopped, os.Interrupt, syscall.SIGTERM)
  <-stopped
}

type backendState struct{}

func (backendState) generation() int { return 1 }
`
	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte(originalMain), 0o644))

	tools := filepath.Join(root, "tools")
	require.NoError(t, os.MkdirAll(tools, 0o755))
	npm := `#!/bin/sh
set -eu
printf '%s\n' "$*" >> "$WAILS_TEST_NPM_LOG"
case "${1:-}:${2:-}" in
  install:)
    mkdir -p node_modules
    ;;
  run:build)
    mkdir -p dist
    printf 'dev\n' > dist/index.html
    ;;
  run:dev|run:serve)
    exec "$WAILS_TEST_HELPER" -test.run='^TestManifestDevFrontendHelper$'
    ;;
  *)
    printf 'unsupported fake npm invocation: %s\n' "$*" >&2
    exit 2
    ;;
esac
`
	require.NoError(t, os.WriteFile(filepath.Join(tools, "npm"), []byte(npm), 0o755))
	t.Setenv("PATH", tools+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("WAILS_TEST_HELPER", os.Args[0])
	hcl := `version = 3

project {
  name = "dev-session"
  product_name = "Dev Session"
  identifier = "com.example.devsession"
  version = "1.0.0"
}

frontend {
  directory = "frontend"
  install = ["npm", "install"]
  build = ["npm", "run", "build"]
  dev = ["npm", "run", "dev"]
  output = "frontend/dist"
}

dev {
  debounce_ms = 25
  watch = ["**/*.go", "wails.hcl"]
  exclude = [".git", ".wails", "bin", "node_modules", "frontend"]
  use_git_ignore = true
  grace_period_ms = 100
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))

	backendLog := filepath.Join(root, "backend.log")
	frontendLog := filepath.Join(root, "frontend.log")
	npmLog := filepath.Join(root, "npm.log")
	t.Setenv("WAILS_TEST_BACKEND_LOG", backendLog)
	t.Setenv("WAILS_TEST_FRONTEND_LOG", frontendLog)
	t.Setenv("WAILS_TEST_NPM_LOG", npmLog)
	port := reserveManifestDevPort(t)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- runManifestDevContext(ctx, &DevOptions{VitePort: port})
		close(done)
	}()
	t.Cleanup(func() {
		cancel()
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Error("development session did not stop during cleanup")
		}
	})

	waitForManifestDevStart(t, backendLog, done, 30*time.Second)
	waitForManifestDevLines(t, frontendLog, 1, 5*time.Second)
	assert.FileExists(t, filepath.Join(root, ".wails", "dev", runtime.GOOS+"-"+runtime.GOARCH, "dev-session"))
	assert.NoFileExists(t, filepath.Join(root, "bin", "dev-session"))
	initialBackend := manifestDevLoggedPIDs(t, backendLog)[0]
	initialFrontend := manifestDevLoggedPIDs(t, frontendLog)[0]
	assert.True(t, manifestDevProcessAlive(initialBackend))
	assert.True(t, manifestDevProcessAlive(initialFrontend))

	// The persistent frontend process owns frontend-source changes through HMR.
	require.NoError(t, os.WriteFile(filepath.Join(root, "frontend", "src.js"), []byte("export const value = 2;\n"), 0o644))
	time.Sleep(250 * time.Millisecond)
	assert.Len(t, manifestDevLoggedPIDs(t, backendLog), 1, "frontend-only changes must not rebuild the backend")
	assert.Len(t, manifestDevLoggedPIDs(t, frontendLog), 1, "frontend-only changes must not restart the frontend server")

	// Ignored Go test changes are no-op events for the running session.
	require.NoError(t, os.WriteFile(filepath.Join(root, "ignored_test.go"), []byte("package main\n"), 0o644))
	time.Sleep(250 * time.Millisecond)
	assert.Len(t, manifestDevLoggedPIDs(t, backendLog), 1)
	require.NoError(t, os.MkdirAll(filepath.Join(root, "ignored", "deep"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "ignored", "deep", "service.go"), []byte("package ignored\n"), 0o644))
	time.Sleep(250 * time.Millisecond)
	assert.Len(t, manifestDevLoggedPIDs(t, backendLog), 1, "new gitignored directories must not refresh watches or rebuild")
	require.NoError(t, os.MkdirAll(filepath.Join(root, "new", "nested"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "new", "nested", "ignored_test.go"), []byte("package nested\n"), 0o644))
	time.Sleep(250 * time.Millisecond)
	assert.Len(t, manifestDevLoggedPIDs(t, backendLog), 1, "new directories containing only ignored inputs are no-ops")

	// Ignore policy reloads in place and does not schedule a build.
	require.NoError(t, os.WriteFile(filepath.Join(root, ".gitignore"), []byte("ignored/\n*.generated.go\n"), 0o644))
	time.Sleep(250 * time.Millisecond)
	assert.Len(t, manifestDevLoggedPIDs(t, backendLog), 1)
	assert.Len(t, manifestDevLoggedPIDs(t, frontendLog), 1)

	// Frontend command changes replace only the persistent frontend process.
	updatedHCL := strings.Replace(hcl, `dev = ["npm", "run", "dev"]`, `dev = ["npm", "run", "serve"]`, 1)
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(updatedHCL), 0o644))
	waitForManifestDevLinesOrExit(t, frontendLog, 2, done, 10*time.Second)
	replacementFrontend := manifestDevLoggedPIDs(t, frontendLog)[1]
	waitForManifestDevCondition(t, 5*time.Second, func() bool { return !manifestDevProcessAlive(initialFrontend) })
	assert.True(t, manifestDevProcessAlive(replacementFrontend))
	time.Sleep(400 * time.Millisecond)
	assert.Len(t, manifestDevLoggedPIDs(t, backendLog), 1, "frontend session changes must not restart the backend")

	// A failed same-port frontend transition restores the previous server.
	failedFrontendHCL := strings.Replace(updatedHCL, `dev = ["npm", "run", "serve"]`, `dev = ["npm", "run", "missing"]`, 1)
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(failedFrontendHCL), 0o644))
	waitForManifestDevLinesOrExit(t, frontendLog, 3, done, 10*time.Second)
	restoredFrontend := manifestDevLoggedPIDs(t, frontendLog)[2]
	waitForManifestDevCondition(t, 5*time.Second, func() bool { return !manifestDevProcessAlive(replacementFrontend) })
	assert.True(t, manifestDevProcessAlive(restoredFrontend))
	assert.Len(t, manifestDevLoggedPIDs(t, backendLog), 1)

	// Restoring the manifest to the active command is an all-cached no-op.
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(updatedHCL), 0o644))
	time.Sleep(400 * time.Millisecond)
	assert.Len(t, manifestDevLoggedPIDs(t, frontendLog), 3)
	assert.Len(t, manifestDevLoggedPIDs(t, backendLog), 1)
	activeFrontend := restoredFrontend

	// Watch-policy and frontend changes commit together, and the replacement
	// watcher must observe the next exact manifest edit.
	watchHCL := strings.Replace(updatedHCL, `watch = ["**/*.go", "wails.hcl"]`, `watch = ["**/*.go", "wails.hcl", "README.md"]`, 1)
	watchHCL = strings.Replace(watchHCL, `dev = ["npm", "run", "serve"]`, `dev = ["npm", "run", "dev"]`, 1)
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(watchHCL), 0o644))
	waitForManifestDevLinesOrExit(t, frontendLog, 4, done, 10*time.Second)
	waitForManifestDevCondition(t, 5*time.Second, func() bool { return !manifestDevProcessAlive(activeFrontend) })
	activeFrontend = manifestDevLoggedPIDs(t, frontendLog)[3]
	assert.True(t, manifestDevProcessAlive(activeFrontend))

	watchServeHCL := strings.Replace(watchHCL, `dev = ["npm", "run", "dev"]`, `dev = ["npm", "run", "serve"]`, 1)
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(watchServeHCL), 0o644))
	waitForManifestDevLinesOrExit(t, frontendLog, 5, done, 10*time.Second)
	waitForManifestDevCondition(t, 5*time.Second, func() bool { return !manifestDevProcessAlive(activeFrontend) })
	activeFrontend = manifestDevLoggedPIDs(t, frontendLog)[4]
	assert.True(t, manifestDevProcessAlive(activeFrontend))
	assert.Len(t, manifestDevLoggedPIDs(t, backendLog), 1)

	// A method-body edit compiles and swaps the backend without rebuilding the frontend.
	replacementMain := strings.Replace(originalMain, "return 1", "return 2", 1)
	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte(replacementMain), 0o644))
	waitForManifestDevLines(t, backendLog, 2, 30*time.Second)
	replacementBackend := manifestDevLoggedPIDs(t, backendLog)[1]
	assert.NotEqual(t, initialBackend, replacementBackend)
	waitForManifestDevCondition(t, 5*time.Second, func() bool { return !manifestDevProcessAlive(initialBackend) })
	assert.True(t, manifestDevProcessAlive(replacementBackend))
	assert.Equal(t, 1, strings.Count(readManifestDevFile(t, npmLog), "run build\n"), "method-body edits must not rebuild the frontend")

	// A failed rebuild preserves the healthy replacement backend.
	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\nfunc broken("), 0o644))
	time.Sleep(750 * time.Millisecond)
	assert.Len(t, manifestDevLoggedPIDs(t, backendLog), 2)
	assert.True(t, manifestDevProcessAlive(replacementBackend), "failed generation must preserve the current backend")

	// A compiled replacement that exits during readiness is rejected.
	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\nfunc main() {}\n"), 0o644))
	time.Sleep(750 * time.Millisecond)
	assert.Len(t, manifestDevLoggedPIDs(t, backendLog), 2)
	assert.True(t, manifestDevProcessAlive(replacementBackend), "an unhealthy replacement must not displace the current backend")

	// A subsequent valid edit recovers and replaces the healthy backend transactionally.
	recoveryMain := strings.Replace(originalMain, "return 1", "return 3", 1)
	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte(recoveryMain), 0o644))
	waitForManifestDevLines(t, backendLog, 3, 30*time.Second)
	recoveryBackend := manifestDevLoggedPIDs(t, backendLog)[2]
	waitForManifestDevCondition(t, 5*time.Second, func() bool { return !manifestDevProcessAlive(replacementBackend) })
	assert.True(t, manifestDevProcessAlive(recoveryBackend))

	cancel()
	require.NoError(t, <-done)
	waitForManifestDevCondition(t, 5*time.Second, func() bool {
		return !manifestDevProcessAlive(recoveryBackend) && !manifestDevProcessAlive(activeFrontend)
	})
	assert.False(t, manifestDevProcessAlive(recoveryBackend))
	assert.False(t, manifestDevProcessAlive(activeFrontend))
}

func TestManifestDevStartupFailuresAreTransactional(t *testing.T) {
	t.Run("initial build", func(t *testing.T) {
		fixture := newManifestDevStartupFixture(t, "package main\nfunc main() {}\n", "#!/bin/sh\nexit 7\n")
		err := runManifestDevContext(t.Context(), &DevOptions{VitePort: fixture.port})
		require.Error(t, err)
		assert.Empty(t, manifestDevLoggedPIDs(t, fixture.frontendLog))
		assert.Empty(t, manifestDevLoggedPIDs(t, fixture.backendLog))
	})

	t.Run("frontend launch", func(t *testing.T) {
		fixture := newManifestDevStartupFixture(t, "package main\nfunc main() {}\n", `#!/bin/sh
set -eu
case "${1:-}:${2:-}" in
  install:) mkdir -p node_modules ;;
  run:build) mkdir -p dist; printf 'dev\n' > dist/index.html; rm -- "$0" ;;
esac
`)
		err := runManifestDevContext(t.Context(), &DevOptions{VitePort: fixture.port})
		require.Error(t, err)
		assert.Empty(t, manifestDevLoggedPIDs(t, fixture.frontendLog))
		assert.Empty(t, manifestDevLoggedPIDs(t, fixture.backendLog))
	})

	t.Run("frontend readiness", func(t *testing.T) {
		fixture := newManifestDevStartupFixture(t, "package main\nfunc main() {}\n", `#!/bin/sh
set -eu
case "${1:-}:${2:-}" in
  install:) mkdir -p node_modules ;;
  run:build) mkdir -p dist; printf 'dev\n' > dist/index.html ;;
  run:dev) exit 9 ;;
esac
`)
		err := runManifestDevContext(t.Context(), &DevOptions{VitePort: fixture.port})
		assert.ErrorContains(t, err, "frontend readiness")
		assert.Empty(t, manifestDevLoggedPIDs(t, fixture.backendLog))
	})

	t.Run("missing built binary", func(t *testing.T) {
		fixture := newManifestDevStartupFixture(t, manifestDevHoldingMain, manifestDevWorkingNPM)
		t.Setenv("WAILS_TEST_DELETE_BINARY", filepath.Join(fixture.root, ".wails", "dev", runtime.GOOS+"-"+runtime.GOARCH, "dev-session"))
		err := runManifestDevContext(t.Context(), &DevOptions{VitePort: fixture.port})
		require.Error(t, err)
		waitForManifestDevCondition(t, 5*time.Second, func() bool {
			pids := manifestDevLoggedPIDs(t, fixture.frontendLog)
			return len(pids) == 1 && !manifestDevProcessAlive(pids[0])
		})
		assert.Empty(t, manifestDevLoggedPIDs(t, fixture.backendLog))
	})

	t.Run("backend readiness", func(t *testing.T) {
		fixture := newManifestDevStartupFixture(t, "package main\nfunc main() {}\n", manifestDevWorkingNPM)
		err := runManifestDevContext(t.Context(), &DevOptions{VitePort: fixture.port})
		assert.ErrorContains(t, err, "backend readiness")
		waitForManifestDevCondition(t, 5*time.Second, func() bool {
			pids := manifestDevLoggedPIDs(t, fixture.frontendLog)
			return len(pids) == 1 && !manifestDevProcessAlive(pids[0])
		})
	})

	t.Run("watcher registration", func(t *testing.T) {
		if os.Geteuid() == 0 {
			t.Skip("permission failure requires an unprivileged user")
		}
		fixture := newManifestDevStartupFixture(t, manifestDevHoldingMain, `#!/bin/sh
set -eu
case "${1:-}:${2:-}" in
  install:) mkdir -p node_modules ;;
  run:build) mkdir -p dist; printf 'dev\n' > dist/index.html ;;
  run:dev) mkdir -p ../sealed; chmod 000 ../sealed; exec "$WAILS_TEST_HELPER" -test.run='^TestManifestDevFrontendHelper$' ;;
  *) exit 2 ;;
esac
`)
		sealed := filepath.Join(fixture.root, "sealed")
		t.Cleanup(func() { _ = os.Chmod(sealed, 0o755) })
		err := runManifestDevContext(t.Context(), &DevOptions{VitePort: fixture.port})
		assert.ErrorContains(t, err, "watch project")
		for _, path := range []string{fixture.frontendLog, fixture.backendLog} {
			for _, pid := range manifestDevLoggedPIDs(t, path) {
				waitForManifestDevCondition(t, 5*time.Second, func() bool { return !manifestDevProcessAlive(pid) })
			}
		}
	})
}

func TestManifestDevStartupCancellationIsClean(t *testing.T) {
	t.Run("frontend readiness", func(t *testing.T) {
		fixture := newManifestDevStartupFixture(t, manifestDevHoldingMain, manifestDevWorkingNPM)
		t.Setenv("WAILS_TEST_FRONTEND_DELAY", "500ms")
		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan error, 1)
		go func() { done <- runManifestDevContext(ctx, &DevOptions{VitePort: fixture.port}) }()
		waitForManifestDevLines(t, fixture.frontendLog, 1, 30*time.Second)
		cancel()
		select {
		case err := <-done:
			require.NoError(t, err)
		case <-time.After(5 * time.Second):
			t.Fatal("frontend-readiness cancellation did not stop the session")
		}
	})

	t.Run("backend readiness", func(t *testing.T) {
		fixture := newManifestDevStartupFixture(t, manifestDevHoldingMain, manifestDevWorkingNPM)
		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan error, 1)
		go func() { done <- runManifestDevContext(ctx, &DevOptions{VitePort: fixture.port}) }()
		waitForManifestDevStart(t, fixture.backendLog, done, 30*time.Second)
		cancel()
		select {
		case err := <-done:
			require.NoError(t, err)
		case <-time.After(5 * time.Second):
			t.Fatal("backend-readiness cancellation did not stop the session")
		}
		for _, pid := range manifestDevLoggedPIDs(t, fixture.frontendLog) {
			waitForManifestDevCondition(t, 5*time.Second, func() bool { return !manifestDevProcessAlive(pid) })
		}
		for _, pid := range manifestDevLoggedPIDs(t, fixture.backendLog) {
			waitForManifestDevCondition(t, 5*time.Second, func() bool { return !manifestDevProcessAlive(pid) })
		}
	})
}

func TestManifestDevUsesDefaultSecureSessionAddress(t *testing.T) {
	listener, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(defaultVitePort)))
	if err != nil {
		t.Skipf("default Dev port is already occupied: %v", err)
	}
	require.NoError(t, listener.Close())

	mainSource := strings.Replace(manifestDevHoldingMain,
		`func main() {`,
		`func main() {
  if output := os.Getenv("WAILS_TEST_BACKEND_ENV"); output != "" {
    _ = os.WriteFile(output, []byte(os.Getenv("FRONTEND_DEVSERVER_URL")+"\n"+os.Getenv("WAILS_VITE_PORT")), 0600)
  }`, 1)
	fixture := newManifestDevStartupFixture(t, mainSource, manifestDevWorkingNPM)
	environmentOutput := filepath.Join(fixture.root, "backend-env")
	t.Setenv("WAILS_TEST_BACKEND_ENV", environmentOutput)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- runManifestDevContext(ctx, &DevOptions{Secure: true}) }()
	waitForManifestDevStart(t, fixture.backendLog, done, 30*time.Second)
	require.Equal(t, "https://127.0.0.1:9245\n9245", strings.TrimSpace(readTestFile(t, environmentOutput)))
	cancel()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("secure development session did not stop")
	}
}

func TestManifestDevRetainsWatcherWhenRefreshFails(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("permission failure requires an unprivileged user")
	}
	fixture := newManifestDevStartupFixture(t, manifestDevHoldingMain, manifestDevWorkingNPM)
	manifestPath := filepath.Join(fixture.root, manifest.Filename)
	hcl := strings.Replace(readTestFile(t, manifestPath), "use_git_ignore = false", "use_git_ignore = true", 1)
	require.NoError(t, os.WriteFile(manifestPath, []byte(hcl), 0o644))
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- runManifestDevContext(ctx, &DevOptions{VitePort: fixture.port}) }()
	t.Cleanup(func() {
		cancel()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Error("development session did not stop during cleanup")
		}
	})
	waitForManifestDevStart(t, fixture.backendLog, done, 30*time.Second)
	// The backend writes its PID before the controller's stability check and
	// watcher registration complete.
	time.Sleep(300 * time.Millisecond)

	sealed := filepath.Join(fixture.root, "sealed")
	require.NoError(t, os.Mkdir(sealed, 0o755))
	require.NoError(t, os.Chmod(sealed, 0))
	t.Cleanup(func() { _ = os.Chmod(sealed, 0o755) })
	time.Sleep(300 * time.Millisecond)
	require.NoError(t, os.WriteFile(filepath.Join(fixture.root, ".gitignore"), []byte("generated/\n"), 0o644))
	time.Sleep(300 * time.Millisecond)
	select {
	case err := <-done:
		done <- err
		require.NoError(t, err, "failed watch refresh must not terminate the session")
		t.Fatal("development session stopped after failed watch refresh")
	default:
	}

	require.NoError(t, os.Chmod(sealed, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(fixture.root, "main.go"), []byte(strings.Replace(manifestDevHoldingMain, "func main() {", "var changed = true\n\nfunc main() {", 1)), 0o644))
	waitForManifestDevLinesOrExit(t, fixture.backendLog, 2, done, 15*time.Second)
}

func TestManifestDevReportsLostWorkingDirectory(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	require.NoError(t, os.Remove(root))
	err := runManifestDevContext(t.Context(), &DevOptions{})
	require.Error(t, err)
}

func TestManifestDevReportsUnexpectedOwnedProcessExit(t *testing.T) {
	for _, processName := range []string{"frontend", "backend"} {
		t.Run(processName, func(t *testing.T) {
			fixture := newManifestDevStartupFixture(t, manifestDevHoldingMain, manifestDevWorkingNPM)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			done := make(chan error, 1)
			go func() { done <- runManifestDevContext(ctx, &DevOptions{VitePort: fixture.port}) }()
			waitForManifestDevStart(t, fixture.backendLog, done, 30*time.Second)
			waitForManifestDevLines(t, fixture.frontendLog, 1, 5*time.Second)
			time.Sleep(200 * time.Millisecond)
			backendPID := manifestDevLoggedPIDs(t, fixture.backendLog)[0]
			frontendPID := manifestDevLoggedPIDs(t, fixture.frontendLog)[0]
			pid := frontendPID
			if processName == "backend" {
				pid = backendPID
			}
			require.NoError(t, syscall.Kill(pid, syscall.SIGTERM))

			select {
			case err := <-done:
				assert.ErrorContains(t, err, processName+" process exited unexpectedly")
			case <-time.After(5 * time.Second):
				t.Fatal("development session did not report the owned process exit")
			}
			waitForManifestDevCondition(t, 5*time.Second, func() bool {
				return !manifestDevProcessAlive(frontendPID) && !manifestDevProcessAlive(backendPID)
			})
		})
	}
}

func TestManifestDevCancelsStaleInFlightGeneration(t *testing.T) {
	mainSource := strings.Replace(manifestDevHoldingMain, "func main() {", `type backendState struct{}

func (backendState) generation() int { return 1 }

func main() {
  if log := os.Getenv("WAILS_TEST_GENERATION_LOG"); log != "" {
    file, _ := os.OpenFile(log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
    _, _ = file.WriteString(strconv.Itoa((backendState{}).generation()) + "\n")
    _ = file.Close()
  }`, 1)
	fixture := newManifestDevStartupFixture(t, mainSource, manifestDevWorkingNPM)
	realGo, err := exec.LookPath("go")
	require.NoError(t, err)
	goBuildLog := filepath.Join(fixture.root, "go-build.log")
	generationLog := filepath.Join(fixture.root, "generation.log")
	t.Setenv("WAILS_TEST_REAL_GO", realGo)
	t.Setenv("WAILS_TEST_GO_BUILD_LOG", goBuildLog)
	t.Setenv("WAILS_TEST_GENERATION_LOG", generationLog)
	goWrapper := `#!/bin/sh
set -eu
if [ "${1:-}" = "build" ]; then
  printf 'build\n' >> "$WAILS_TEST_GO_BUILD_LOG"
  sleep 0.6
fi
exec "$WAILS_TEST_REAL_GO" "$@"
`
	require.NoError(t, os.WriteFile(filepath.Join(fixture.root, "tools", "go"), []byte(goWrapper), 0o755))

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- runManifestDevContext(ctx, &DevOptions{VitePort: fixture.port}) }()
	waitForManifestDevStart(t, fixture.backendLog, done, 30*time.Second)
	time.Sleep(200 * time.Millisecond)

	generationTwo := strings.Replace(mainSource, "return 1", "return 2", 1)
	require.NoError(t, os.WriteFile(filepath.Join(fixture.root, "main.go"), []byte(generationTwo), 0o644))
	waitForManifestDevCondition(t, 10*time.Second, func() bool { return manifestDevFileLineCount(goBuildLog) >= 2 })
	generationThree := strings.Replace(mainSource, "return 1", "return 3", 1)
	require.NoError(t, os.WriteFile(filepath.Join(fixture.root, "main.go"), []byte(generationThree), 0o644))
	waitForManifestDevLines(t, fixture.backendLog, 2, 30*time.Second)
	time.Sleep(750 * time.Millisecond)
	assert.Equal(t, []string{"1", "3"}, strings.Fields(readManifestDevFile(t, generationLog)), "the canceled generation must never start")
	assert.Len(t, manifestDevLoggedPIDs(t, fixture.backendLog), 2)

	cancel()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("development session did not stop after stale-generation test")
	}
}

const manifestDevHoldingMain = `package main

import (
  "os"
  "os/signal"
  "strconv"
  "syscall"
)

func main() {
  if log := os.Getenv("WAILS_TEST_BACKEND_LOG"); log != "" {
    file, _ := os.OpenFile(log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
    _, _ = file.WriteString(strconv.Itoa(os.Getpid()) + "\n")
    _ = file.Close()
  }
  stopped := make(chan os.Signal, 1)
  signal.Notify(stopped, os.Interrupt, syscall.SIGTERM)
  <-stopped
}
`

const manifestDevWorkingNPM = `#!/bin/sh
set -eu
case "${1:-}:${2:-}" in
  install:) mkdir -p node_modules ;;
  run:build) mkdir -p dist; printf 'dev\n' > dist/index.html ;;
  run:dev) exec "$WAILS_TEST_HELPER" -test.run='^TestManifestDevFrontendHelper$' ;;
  *) exit 2 ;;
esac
`

type manifestDevStartupFixture struct {
	root, frontendLog, backendLog string
	port                          int
}

func newManifestDevStartupFixture(t *testing.T, mainSource, npmSource string) manifestDevStartupFixture {
	t.Helper()
	root := t.TempDir()
	t.Chdir(root)
	require.NoError(t, os.MkdirAll(filepath.Join(root, "frontend"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "frontend", "package.json"), []byte("{}\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/devstartup\n\ngo 1.24\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte(mainSource), 0o644))
	tools := filepath.Join(root, "tools")
	require.NoError(t, os.MkdirAll(tools, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(tools, "npm"), []byte(npmSource), 0o755))
	t.Setenv("PATH", tools+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("WAILS_TEST_HELPER", os.Args[0])
	frontendLog := filepath.Join(root, "frontend.log")
	backendLog := filepath.Join(root, "backend.log")
	t.Setenv("WAILS_TEST_FRONTEND_LOG", frontendLog)
	t.Setenv("WAILS_TEST_BACKEND_LOG", backendLog)
	hcl := `version = 3
project {
  name = "dev-session"
  product_name = "Dev Session"
  identifier = "com.example.devsession"
  version = "1.0.0"
}
frontend {
  directory = "frontend"
  install = ["npm", "install"]
  build = ["npm", "run", "build"]
  dev = ["npm", "run", "dev"]
  output = "frontend/dist"
}
dev {
  debounce_ms = 0
  watch = ["**/*.go", "wails.hcl"]
  exclude = [".git", ".wails", "bin", "node_modules", "frontend"]
  use_git_ignore = false
  grace_period_ms = 50
}
`
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte(hcl), 0o644))
	return manifestDevStartupFixture{root: root, frontendLog: frontendLog, backendLog: backendLog, port: reserveManifestDevPort(t)}
}

func TestManifestDevFrontendHelper(t *testing.T) {
	port := os.Getenv(wailsVitePort)
	if port == "" {
		return
	}
	if binary := os.Getenv("WAILS_TEST_DELETE_BINARY"); binary != "" {
		require.NoError(t, os.Remove(binary))
	}
	if log := os.Getenv("WAILS_TEST_FRONTEND_LOG"); log != "" {
		file, openErr := os.OpenFile(log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
		require.NoError(t, openErr)
		_, writeErr := file.WriteString(strconv.Itoa(os.Getpid()) + "\n")
		require.NoError(t, writeErr)
		require.NoError(t, file.Close())
	}
	if delay, err := time.ParseDuration(os.Getenv("WAILS_TEST_FRONTEND_DELAY")); err == nil && delay > 0 {
		time.Sleep(delay)
	}
	listener, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", port))
	require.NoError(t, err)
	defer listener.Close()
	for {
		connection, acceptErr := listener.Accept()
		if acceptErr != nil {
			return
		}
		_ = connection.Close()
	}
}

func reserveManifestDevPort(t *testing.T) int {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := listener.Addr().(*net.TCPAddr).Port
	require.NoError(t, listener.Close())
	return port
}

func waitForManifestDevLines(t *testing.T, path string, count int, timeout time.Duration) {
	t.Helper()
	waitForManifestDevCondition(t, timeout, func() bool { return len(manifestDevLoggedPIDs(t, path)) >= count })
}

func waitForManifestDevStart(t *testing.T, path string, done <-chan error, timeout time.Duration) {
	t.Helper()
	deadline := time.NewTimer(timeout)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer deadline.Stop()
	defer ticker.Stop()
	for {
		if len(manifestDevLoggedPIDs(t, path)) > 0 {
			return
		}
		select {
		case err := <-done:
			require.NoError(t, err, "development session exited before starting the backend")
			require.FailNow(t, "development session exited before starting the backend")
		case <-deadline.C:
			require.FailNow(t, "development session did not start", "timeout: %s", timeout)
		case <-ticker.C:
		}
	}
}

func waitForManifestDevLinesOrExit(t *testing.T, path string, count int, done <-chan error, timeout time.Duration) {
	t.Helper()
	deadline := time.NewTimer(timeout)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer deadline.Stop()
	defer ticker.Stop()
	for {
		if len(manifestDevLoggedPIDs(t, path)) >= count {
			return
		}
		select {
		case err := <-done:
			require.NoError(t, err, "development session exited before the expected transition")
			require.FailNow(t, "development session exited before the expected transition")
		case <-deadline.C:
			require.FailNow(t, "development session transition timed out", "timeout: %s", timeout)
		case <-ticker.C:
		}
	}
}

func waitForManifestDevCondition(t *testing.T, timeout time.Duration, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	require.FailNow(t, "condition was not satisfied before timeout", "timeout: %s", timeout)
}

func manifestDevLoggedPIDs(t *testing.T, path string) []int {
	t.Helper()
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	require.NoError(t, err)
	var result []int
	for _, line := range strings.Fields(string(data)) {
		pid, parseErr := strconv.Atoi(line)
		require.NoError(t, parseErr)
		result = append(result, pid)
	}
	return result
}

func readManifestDevFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return string(data)
}

func manifestDevFileLineCount(path string) int {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	return len(strings.Fields(string(data)))
}

func manifestDevProcessAlive(pid int) bool {
	return syscall.Kill(pid, 0) == nil
}
