package commands

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/rjeczalik/notify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/cache"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestLegacyDevRejectsManifestOnlyFlags(t *testing.T) {
	t.Chdir(t.TempDir())
	err := Dev(&DevOptions{Profile: "release"})
	assert.ErrorContains(t, err, "require an active wails.hcl")
}

func TestLegacyDevAppliesExplicitSessionOptionsBeforeWatcherStartup(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := listener.Addr().(*net.TCPAddr).Port
	require.NoError(t, listener.Close())
	t.Setenv("WAILS_MCP", "")
	t.Setenv("EXTRA_TAGS", "existing")
	err = Dev(&DevOptions{Config: filepath.Join(root, "missing.yml"), VitePort: port, Secure: true, Tags: "debug,existing"})
	require.Error(t, err)
	assert.Equal(t, fmt.Sprintf("https://localhost:%d", port), os.Getenv("FRONTEND_DEVSERVER_URL"))
	assert.Equal(t, fmt.Sprintf("%d", port), os.Getenv(wailsVitePort))
	assert.Equal(t, "existing,debug", os.Getenv("EXTRA_TAGS"))
}

func TestLegacyDevAcceptsEnvironmentPort(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := listener.Addr().(*net.TCPAddr).Port
	require.NoError(t, listener.Close())
	t.Setenv(wailsVitePort, fmt.Sprintf("%d", port))
	err = Dev(&DevOptions{Config: filepath.Join(root, "missing.yml")})
	require.Error(t, err)
	assert.Equal(t, fmt.Sprintf("http://localhost:%d", port), os.Getenv("FRONTEND_DEVSERVER_URL"))
}

func TestLegacyDevRejectsUnavailablePort(t *testing.T) {
	t.Chdir(t.TempDir())
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	err = Dev(&DevOptions{VitePort: port})
	assert.Error(t, err)
}

func TestLegacyDevUsesDefaultPortForInvalidEnvironmentValue(t *testing.T) {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", defaultVitePort))
	if err != nil {
		t.Skipf("default Dev port is occupied: %v", err)
	}
	require.NoError(t, listener.Close())
	root := t.TempDir()
	t.Chdir(root)
	t.Setenv(wailsVitePort, "not-a-port")
	err = Dev(&DevOptions{Config: filepath.Join(root, "missing.yml")})
	require.Error(t, err)
	assert.Equal(t, fmt.Sprintf("http://localhost:%d", defaultVitePort), os.Getenv("FRONTEND_DEVSERVER_URL"))
}

func TestManifestDevRejectsProductionProfiles(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	require.NoError(t, manifest.WriteMinimal(root, manifest.Project{Name: "dev", ProductName: "Dev", Identifier: "com.example.dev", Version: "1.0.0"}))

	err := Dev(&DevOptions{Profile: "release", Plan: true})
	assert.ErrorContains(t, err, "dev does not accept production profiles")
}

func TestManifestDevValidatesBeforeStartingProcesses(t *testing.T) {
	prependFakePlanTools(t, "npm")
	t.Run("malformed manifest", func(t *testing.T) {
		root := t.TempDir()
		t.Chdir(root)
		require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte("version = ["), 0o644))
		err := runManifestDevContext(t.Context(), &DevOptions{Plan: true})
		assert.ErrorContains(t, err, "parse wails.hcl")
	})

	t.Run("malformed target", func(t *testing.T) {
		root := t.TempDir()
		t.Chdir(root)
		require.NoError(t, manifest.WriteMinimal(root, manifest.Project{Name: "dev", ProductName: "Dev", Identifier: "com.example.dev", Version: "1.0.0"}))
		err := runManifestDevContext(t.Context(), &DevOptions{Target: "linux"})
		assert.ErrorContains(t, err, "platform/architecture")
	})

	t.Run("foreign target", func(t *testing.T) {
		root := t.TempDir()
		t.Chdir(root)
		require.NoError(t, manifest.WriteMinimal(root, manifest.Project{Name: "dev", ProductName: "Dev", Identifier: "com.example.dev", Version: "1.0.0"}))
		foreign := "linux/amd64"
		if runtime.GOOS == "linux" && runtime.GOARCH == "amd64" {
			foreign = "windows/amd64"
		}
		err := runManifestDevContext(t.Context(), &DevOptions{Target: foreign})
		assert.ErrorContains(t, err, "dev target must be the host")
	})

	t.Run("plan", func(t *testing.T) {
		root := t.TempDir()
		t.Chdir(root)
		require.NoError(t, manifest.WriteMinimal(root, manifest.Project{Name: "dev", ProductName: "Dev", Identifier: "com.example.dev", Version: "1.0.0"}))
		require.NoError(t, runManifestDevContext(t.Context(), &DevOptions{Plan: true}))
		assert.NoDirExists(t, filepath.Join(root, ".wails"), "a Dev Plan must not mutate the project")
	})

	t.Run("occupied port", func(t *testing.T) {
		root := t.TempDir()
		t.Chdir(root)
		require.NoError(t, manifest.WriteMinimal(root, manifest.Project{Name: "dev", ProductName: "Dev", Identifier: "com.example.dev", Version: "1.0.0"}))
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		require.NoError(t, err)
		defer listener.Close()
		port := listener.Addr().(*net.TCPAddr).Port
		err = runManifestDevContext(t.Context(), &DevOptions{VitePort: port})
		assert.ErrorContains(t, err, "frontend port")
		assert.ErrorContains(t, err, "is unavailable")
	})

	t.Run("canceled startup", func(t *testing.T) {
		root := t.TempDir()
		t.Chdir(root)
		require.NoError(t, manifest.WriteMinimal(root, manifest.Project{Name: "dev", ProductName: "Dev", Identifier: "com.example.dev", Version: "1.0.0"}))
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		require.NoError(t, err)
		port := listener.Addr().(*net.TCPAddr).Port
		require.NoError(t, listener.Close())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		require.NoError(t, runManifestDevContext(ctx, &DevOptions{VitePort: port}))
	})
}

func TestDevPathExcludedAlwaysSkipsGeneratedAndDependencyTrees(t *testing.T) {
	config := manifest.Config{Build: manifest.Build{OutputDirectory: "release-output"}, Frontend: manifest.Frontend{Directory: "frontend"}, Dev: manifest.Dev{Exclude: []string{"tmp"}}}
	for _, path := range []string{".wails/cache", ".git/objects", "frontend/src/main.ts", "tools/node_modules/pkg", "tmp/generated.go"} {
		assert.True(t, devPathExcluded(config, path), path)
	}
	assert.False(t, devPathExcluded(config, "release-output/app"), "finite build output is not Dev watch policy")
	assert.False(t, devPathExcluded(config, "internal/service.go"))
}

func TestMatchesDevWatchDoubleStar(t *testing.T) {
	patterns := []string{"internal/**/generated/*.go", "**/*.go", "wails.hcl"}
	assert.True(t, matchesDevWatch(patterns, "internal/service.go"))
	assert.True(t, matchesDevWatch(patterns, "internal/deep/generated/service.go"))
	assert.True(t, matchesDevWatch(patterns, "main.go"))
	assert.True(t, matchesDevWatch(patterns, "wails.hcl"))
	assert.False(t, matchesDevWatch(patterns, "README.md"))
	assert.False(t, matchesDevWatch([]string{"internal/**/generated/*.go"}, "other/generated/service.go"))
	assert.False(t, matchesDevWatch([]string{"[invalid"}, "main.go"))
}

func TestDevWatchPatternsHavePathGlobSemantics(t *testing.T) {
	tests := []struct {
		name, pattern, value string
		want                 bool
	}{
		{name: "double star matches zero directories", pattern: "internal/**/planner.go", value: "internal/planner.go", want: true},
		{name: "double star matches many directories", pattern: "internal/**/planner.go", value: "internal/wake/deep/planner.go", want: true},
		{name: "multiple double stars", pattern: "**/generated/**/models.ts", value: "frontend/generated/bindings/deep/models.ts", want: true},
		{name: "question mark stays within segment", pattern: "internal/test?.go", value: "internal/test1.go", want: true},
		{name: "character class", pattern: "internal/test[0-9].go", value: "internal/test7.go", want: true},
		{name: "segment wildcard does not cross slash", pattern: "internal/*.go", value: "internal/deep/file.go", want: false},
		{name: "suffix mismatch", pattern: "**/*.go", value: "frontend/main.ts", want: false},
		{name: "too few segments", pattern: "internal/*/file.go", value: "internal/file.go", want: false},
		{name: "too many segments", pattern: "internal/file.go", value: "internal/deep/file.go", want: false},
		{name: "invalid segment glob", pattern: "internal/[.go", value: "internal/x.go", want: false},
		{name: "leading relative markers are normalized", pattern: "./**/*.go/", value: "./internal/main.go/", want: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, matchDevPathPattern(test.pattern, test.value))
		})
	}
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
	assert.False(t, isDevGitIgnoreEvent(root, manifest.Config{}, filepath.Join(root, ".gitignore")))
	assert.False(t, isDevGitIgnoreEvent(root, config, filepath.Join(root, "main.go")))
	assert.True(t, isDevGitIgnoreEvent(root, config, filepath.Join(root, ".gitignore")))

	disabled, err := loadDevGitIgnore(root, false)
	require.NoError(t, err)
	assert.Nil(t, disabled)
}

func TestManifestWatchSetupReportsMissingRoots(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing")
	_, err := startManifestWatches(missing, manifest.Config{})
	assert.ErrorContains(t, err, "resolve project root")
	events := make(chan notify.EventInfo, 1)
	err = registerManifestWatches(missing, manifest.Config{}, nil, events)
	assert.Error(t, err)

	root := t.TempDir()
	current, err := startManifestWatches(root, manifest.Config{})
	require.NoError(t, err)
	t.Cleanup(current.stop)
	restarted, err := restartManifestWatches(missing, manifest.Config{}, current)
	assert.ErrorContains(t, err, "resolve project root")
	assert.Same(t, current, restarted, "a failed replacement must retain the healthy watcher")
	changed := filepath.Join(root, "main.go")
	require.NoError(t, os.WriteFile(changed, []byte("package main\n"), 0o644))
	select {
	case event := <-current.events:
		canonicalChanged, resolveErr := filepath.EvalSymlinks(changed)
		require.NoError(t, resolveErr)
		assert.Equal(t, canonicalChanged, event.Path())
	case <-time.After(2 * time.Second):
		t.Fatal("healthy watcher stopped after failed replacement")
	}
}

func TestRestartManifestWatchesContinuesDeliveringEvents(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, ".gitignore"), []byte("ignored/\n"), 0o644))
	require.NoError(t, manifest.WriteMinimal(root, manifest.Project{Name: "watch", ProductName: "Watch", Identifier: "com.example.watch", Version: "1.0.0"}))
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	current, err := startManifestWatches(root, loaded.Config)
	require.NoError(t, err)
	t.Cleanup(current.stop)

	restarted, err := restartManifestWatches(root, loaded.Config, current)
	require.NoError(t, err)
	t.Cleanup(restarted.stop)
	require.NoError(t, os.WriteFile(filepath.Join(root, manifest.Filename), []byte("version = 3\n"), 0o644))

	select {
	case event := <-restarted.events:
		expected, statErr := os.Stat(filepath.Join(root, manifest.Filename))
		require.NoError(t, statErr)
		actual, statErr := os.Stat(event.Path())
		require.NoError(t, statErr)
		assert.True(t, os.SameFile(expected, actual))
		assert.False(t, ignoreDevEvent(restarted.root, loaded.Config, restarted.ignored, event.Path()))
	case <-time.After(2 * time.Second):
		t.Fatal("replacement watcher did not deliver the manifest write")
	}
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
	assert.False(t, directoryContainsDevInput(root, config, ignored, filepath.Join(root, "missing")))
	gitignoreOnly := filepath.Join(root, "gitignore-only")
	require.NoError(t, os.MkdirAll(gitignoreOnly, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(gitignoreOnly, ".gitignore"), []byte("*.tmp\n"), 0o644))
	assert.False(t, directoryContainsDevInput(root, config, ignored, gitignoreOnly))
	assert.True(t, shouldRefreshDevWatchesForDirectory(root, config, ignored, root))
	assert.True(t, ignoreDevEvent(root, config, ignored, root))
}

func TestFrontendSessionChanged(t *testing.T) {
	base := manifest.Config{Frontend: manifest.Frontend{Directory: "frontend", PackageManager: "npm", DevCommand: "dev", Dev: []string{"npm", "run", "dev"}}}
	assert.False(t, frontendSessionChanged(base, base))
	next := base
	next.Frontend.DevCommand = "serve"
	assert.True(t, frontendSessionChanged(base, next))
	next = base
	next.Frontend.Dev = []string{"pnpm", "run", "serve"}
	assert.True(t, frontendSessionChanged(base, next))
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
	for _, manager := range []string{"pnpm", "bun"} {
		args, managerErr := frontendDevArgs(manager, "dev", serverArgs)
		require.NoError(t, managerErr)
		assert.Equal(t, append([]string{"run", "dev"}, serverArgs...), args)
	}
	_, err = frontendDevArgs("deno", "dev", nil)
	assert.ErrorContains(t, err, "unsupported frontend.package_manager")
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
	require.NoError(t, os.WriteFile(filepath.Join(frontend, "package.json"), []byte(`{"scripts":{"dev":"node ./node_modules/.bin/vite.cmd"}}`), 0o644))
	assert.True(t, frontendDevCommandUsesVite(root, config))
	require.NoError(t, os.WriteFile(filepath.Join(frontend, "package.json"), []byte(`{`), 0o644))
	assert.False(t, frontendDevCommandUsesVite(root, config))
	require.NoError(t, os.Remove(filepath.Join(frontend, "package.json")))
	assert.False(t, frontendDevCommandUsesVite(root, config))
}

func TestStartFrontendDevValidatesResolvedCommands(t *testing.T) {
	root := t.TempDir()
	config := manifest.Config{Frontend: manifest.Frontend{Directory: "frontend"}}
	_, err := startFrontendDev(root, config, "127.0.0.1", 9245, "http://127.0.0.1:9245")
	assert.ErrorContains(t, err, "frontend.dev_command is not set")
	config.Frontend.DevCommand = "dev"
	config.Frontend.PackageManager = "deno"
	_, err = startFrontendDev(root, config, "127.0.0.1", 9245, "http://127.0.0.1:9245")
	assert.ErrorContains(t, err, "unsupported frontend.package_manager")
}

func TestStartFrontendDevRunsResolvedNPMThroughNode(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the fixture uses POSIX shell scripts")
	}
	root := t.TempDir()
	frontend := filepath.Join(root, "frontend")
	tools := filepath.Join(root, "tools")
	require.NoError(t, os.MkdirAll(frontend, 0o755))
	require.NoError(t, os.MkdirAll(tools, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(frontend, "package.json"), []byte(`{"scripts":{"dev":"vite"}}`), 0o644))
	invocation := filepath.Join(root, "npm-invocation")
	npm := "#!/bin/sh\nprintf '%s\\n%s\\n%s\\n' \"$*\" \"$FRONTEND_DEVSERVER_URL\" \"$WAILS_VITE_PORT\" > \"$WAILS_TEST_NPM_INVOCATION\"\ntrap 'exit 0' INT TERM\nwhile :; do sleep 1; done\n"
	node := "#!/bin/sh\nexec \"$@\"\n"
	require.NoError(t, os.WriteFile(filepath.Join(tools, "npm"), []byte(npm), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(tools, "node"), []byte(node), 0o755))
	t.Setenv("PATH", tools)
	t.Setenv("WAILS_TEST_NPM_INVOCATION", invocation)
	process, err := startFrontendDev(root, manifest.Config{Frontend: manifest.Frontend{Directory: "frontend", PackageManager: "npm", DevCommand: "dev"}}, "127.0.0.1", 9754, "http://127.0.0.1:9754")
	require.NoError(t, err)
	defer process.stop(100 * time.Millisecond)
	deadline := time.Now().Add(2 * time.Second)
	for {
		if _, statErr := os.Stat(invocation); statErr == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("npm invocation was not recorded")
		}
		time.Sleep(10 * time.Millisecond)
	}
	assert.Equal(t, "run dev -- --host 127.0.0.1 --port 9754 --strictPort\nhttp://127.0.0.1:9754\n9754\n", readTestFile(t, invocation))
}

func TestStartPackageManagerProcessReportsMissingLaunchers(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	_, err := startPackageManagerProcess(t.TempDir(), "npm", nil, "run", "dev")
	assert.ErrorContains(t, err, "executable file not found")

	tools := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(tools, "npm"), []byte("#!/bin/sh\n"), 0o755))
	t.Setenv("PATH", tools)
	_, err = startPackageManagerProcess(t.TempDir(), "npm", nil, "run", "dev")
	assert.ErrorContains(t, err, "executable file not found")

	_, err = startPackageManagerProcess(t.TempDir(), "pnpm", nil, "run", "dev")
	assert.ErrorContains(t, err, "executable file not found")
}

func TestResolveNPMProcessFallsBackToUnresolvedLauncher(t *testing.T) {
	lookup := func(name string) (string, error) { return "/tools/" + name, nil }
	name, args, err := resolvePackageManagerProcess("npm", []string{"run", "dev"}, lookup, func(string) (string, error) {
		return "", errors.New("synthetic symlink resolution failure")
	})
	require.NoError(t, err)
	assert.Equal(t, "/tools/node", name)
	assert.Equal(t, []string{"/tools/npm", "run", "dev"}, args)
}

func TestResolveNPMProcessUsesCommandInterpreterForWindowsLauncher(t *testing.T) {
	lookup := func(name string) (string, error) {
		switch name {
		case "npm":
			return `C:\tools\npm.cmd`, nil
		case "cmd":
			return `C:\Windows\System32\cmd.exe`, nil
		default:
			return "", os.ErrNotExist
		}
	}
	name, args, err := resolvePackageManagerProcess("npm", []string{"run", "dev"}, lookup, func(string) (string, error) {
		t.Fatal("batch launchers must not be passed to node")
		return "", nil
	})
	require.NoError(t, err)
	assert.Equal(t, `C:\Windows\System32\cmd.exe`, name)
	assert.Equal(t, []string{"/d", "/s", "/c", `C:\tools\npm.cmd`, "run", "dev"}, args)
}

func TestRestoreManifestFrontendReportsLaunchAndReadinessFailures(t *testing.T) {
	root := t.TempDir()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := listener.Addr().(*net.TCPAddr).Port
	require.NoError(t, listener.Close())
	url := fmt.Sprintf("http://127.0.0.1:%d", port)
	config := manifest.Config{Frontend: manifest.Frontend{Directory: ".", Dev: []string{"definitely-missing-wails-frontend"}}, Dev: manifest.Dev{GracePeriodMS: 25}}
	_, err = restoreManifestFrontend(t.Context(), root, config, port, url, "127.0.0.1")
	assert.ErrorContains(t, err, "restore previous frontend")

	config.Frontend.Dev = []string{os.Args[0], "-test.run=TestManifestReadinessHelper"}
	t.Setenv("WAILS_READINESS_HELPER", "exit")
	_, err = restoreManifestFrontend(t.Context(), root, config, port, url, "127.0.0.1")
	assert.ErrorContains(t, err, "restore previous frontend readiness")

	t.Setenv("WAILS_READINESS_HELPER", "listen")
	frontend, err := restoreManifestFrontend(t.Context(), root, config, port, url, "127.0.0.1")
	require.NoError(t, err)
	frontend.stop(25 * time.Millisecond)
}

func TestDevWatchSessionChanged(t *testing.T) {
	base := manifest.Config{Build: manifest.Build{OutputDirectory: "bin"}, Frontend: manifest.Frontend{Directory: "frontend"}, Dev: manifest.Dev{Watch: []string{"**/*.go"}, Exclude: []string{"tmp"}, UseGitIgnore: true}}
	assert.False(t, devWatchSessionChanged(base, base))
	next := base
	next.Dev.Watch = []string{"**/*.go", "wails.hcl"}
	assert.True(t, devWatchSessionChanged(base, next))
	next = base
	next.Frontend.Directory = "web"
	assert.True(t, devWatchSessionChanged(base, next))
	next = base
	next.Build.OutputDirectory = "release-output"
	assert.False(t, devWatchSessionChanged(base, next), "finite build output is not Dev session policy")
}

func TestManifestDevBinaryPathUsesTheCompilePlanOutput(t *testing.T) {
	root := t.TempDir()
	key := pipelineCompileKey("linux", "amd64")
	run := manifestPipelineRun{Plan: pipeline.Plan{Nodes: map[pipeline.NodeKey]pipeline.Node{
		key: {Key: key, Kind: pipeline.CompileApplication, Output: ".wails/dev/linux-amd64/app"},
	}}}
	path, err := manifestDevBinaryPath(root, run, "linux", "amd64")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(root, ".wails", "dev", "linux-amd64", "app"), path)

	_, err = manifestDevBinaryPath(root, manifestPipelineRun{}, "linux", "amd64")
	assert.ErrorContains(t, err, "has no compile output")
}

func TestManifestBackendChangedOnlyForExecutedCompileOrURLChange(t *testing.T) {
	key := pipelineCompileKey("linux", "amd64")
	for _, status := range []cache.LookupStatus{cache.LookupHit, cache.LookupRestored} {
		run := manifestPipelineRun{Results: map[pipeline.NodeKey]pipeline.Result{key: {Status: status}}}
		assert.False(t, manifestBackendChanged(run, "linux", "amd64"), status)
	}
	run := manifestPipelineRun{Results: map[pipeline.NodeKey]pipeline.Result{key: {Status: cache.LookupMiss}}}
	assert.True(t, manifestBackendChanged(run, "linux", "amd64"))
	assert.True(t, manifestBackendChanged(manifestPipelineRun{}, "linux", "amd64"))
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

func TestManifestSessionSuppressesOwnedExitAfterCancellation(t *testing.T) {
	process := &manifestProcess{done: make(chan struct{})}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	assert.NoError(t, manifestSessionProcessExit(ctx, "frontend", process))
	assert.EqualError(t, manifestSessionProcessExit(context.Background(), "frontend", process), "frontend process exited unexpectedly")
}

func TestManifestRebuildFailureReportsOrdinaryErrors(t *testing.T) {
	reportManifestRebuildFailure(errors.New("ordinary rebuild failure"))
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
	cleanExit := &manifestProcess{done: make(chan struct{})}
	close(cleanExit.done)
	err = waitForProcessTCP(context.Background(), cleanExit, "127.0.0.1:1", time.Second)
	assert.EqualError(t, err, "process exited before becoming ready")
}

func TestManifestProcessReadinessTimesOut(t *testing.T) {
	process := startReadinessHelper(t, "hold")
	defer process.stop(100 * time.Millisecond)
	err := waitForProcessTCP(context.Background(), process, "127.0.0.1:1", 20*time.Millisecond)
	assert.ErrorContains(t, err, "timed out waiting")
}

func TestManifestProcessReadinessRequiresBackendToStayAlive(t *testing.T) {
	stable := startReadinessHelper(t, "hold")
	defer stable.stop(100 * time.Millisecond)
	require.NoError(t, waitForProcessStable(context.Background(), stable, 50*time.Millisecond))

	exited := &manifestProcess{done: make(chan struct{})}
	go func() {
		time.Sleep(20 * time.Millisecond)
		exited.mu.Lock()
		exited.err = errors.New("synthetic startup failure")
		exited.mu.Unlock()
		close(exited.done)
	}()
	err := waitForProcessStable(context.Background(), exited, 200*time.Millisecond)
	require.ErrorContains(t, err, "process exited during startup")
	cleanExit := &manifestProcess{done: make(chan struct{})}
	close(cleanExit.done)
	err = waitForProcessStable(context.Background(), cleanExit, time.Second)
	assert.EqualError(t, err, "process exited during startup")
}

func TestManifestProcessReadinessHonoursCancellation(t *testing.T) {
	process := startReadinessHelper(t, "hold")
	defer process.stop(100 * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	started := time.Now()
	err := waitForProcessTCP(ctx, process, "127.0.0.1:1", 0)
	require.ErrorIs(t, err, context.Canceled)
	assert.Less(t, time.Since(started), 250*time.Millisecond)

	err = waitForProcessStable(ctx, process, 0)
	require.ErrorIs(t, err, context.Canceled)
}

func TestManifestProcessStopIsIdempotent(t *testing.T) {
	ready := filepath.Join(t.TempDir(), "ready")
	t.Setenv("WAILS_READINESS_READY", ready)
	process := startReadinessHelper(t, "graceful-interrupt")
	waitForManifestDevTestFile(t, ready)
	process.stop(0)
	process.stop(100 * time.Millisecond)
	select {
	case <-process.done:
		assert.NoError(t, process.waitError())
	default:
		t.Fatal("process was not reaped")
	}
}

func TestManifestProcessStopKillsAnInterruptIgnoringProcess(t *testing.T) {
	ready := filepath.Join(t.TempDir(), "ready")
	t.Setenv("WAILS_READINESS_READY", ready)
	process := startReadinessHelper(t, "ignore-interrupt")
	waitForManifestDevTestFile(t, ready)
	process.stop(10 * time.Millisecond)
	select {
	case <-process.done:
		assert.Error(t, process.waitError())
	default:
		t.Fatal("interrupt-ignoring process was not killed and reaped")
	}
	assert.NoError(t, (*manifestProcess)(nil).waitError())
	(*manifestProcess)(nil).stop(0)
	(&manifestProcess{}).stop(0)
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
	case "ignore-interrupt":
		signal.Ignore(os.Interrupt)
		_ = os.WriteFile(os.Getenv("WAILS_READINESS_READY"), nil, 0o600)
		time.Sleep(5 * time.Minute)
	case "graceful-interrupt":
		stopped := make(chan os.Signal, 1)
		signal.Notify(stopped, os.Interrupt)
		_ = os.WriteFile(os.Getenv("WAILS_READINESS_READY"), nil, 0o600)
		<-stopped
	case "exit":
		os.Exit(3)
	case "capture-env":
		data := []byte(os.Getenv("FRONTEND_DEVSERVER_URL") + "\n" + os.Getenv(wailsVitePort))
		if err := os.WriteFile(os.Getenv("WAILS_ENV_OUTPUT"), data, 0o600); err != nil {
			os.Exit(4)
		}
	case "listen":
		listener, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", os.Getenv(wailsVitePort)))
		if err != nil {
			os.Exit(5)
		}
		defer listener.Close()
		time.Sleep(5 * time.Minute)
	}
}

func waitForManifestDevTestFile(t *testing.T, path string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("process readiness marker was not written")
}

func BenchmarkMatchDevPathPattern(b *testing.B) {
	cases := []struct {
		pattern string
		path    string
	}{
		{pattern: "**/*.go", path: "internal/wake/pipeline/planner.go"},
		{pattern: "internal/**/test?.go", path: "internal/wake/cache/test1.go"},
		{pattern: "wails.hcl", path: "wails.hcl"},
		{pattern: "frontend/**/generated/**/*.ts", path: "frontend/src/deep/generated/bindings/models.ts"},
	}
	b.ReportAllocs()
	for b.Loop() {
		for _, item := range cases {
			if !matchDevPathPattern(item.pattern, item.path) {
				b.Fatal("benchmark fixture did not match", item)
			}
		}
	}
}

func BenchmarkCompiledDevWatchMatcher(b *testing.B) {
	patterns := []string{"**/*.go", "internal/**/test?.go", "wails.hcl", "frontend/**/generated/**/*.ts"}
	matcher := newDevWatchMatcher(patterns)
	paths := []string{
		"internal/wake/pipeline/planner.go",
		"internal/wake/cache/test1.go",
		"wails.hcl",
		"frontend/src/deep/generated/bindings/models.ts",
	}
	b.ReportAllocs()
	for b.Loop() {
		for _, value := range paths {
			if !matcher.Match(value) {
				b.Fatal("benchmark fixture did not match", value)
			}
		}
	}
}

func BenchmarkUncompiledDevWatchMatcher(b *testing.B) {
	patterns := []string{"**/*.go", "internal/**/test?.go", "wails.hcl", "frontend/**/generated/**/*.ts"}
	paths := []string{
		"internal/wake/pipeline/planner.go",
		"internal/wake/cache/test1.go",
		"wails.hcl",
		"frontend/src/deep/generated/bindings/models.ts",
	}
	b.ReportAllocs()
	for b.Loop() {
		for _, value := range paths {
			if !matchesDevWatch(patterns, value) {
				b.Fatal("benchmark fixture did not match", value)
			}
		}
	}
}
