//go:build !windows

package commands

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestManifestHookRunsDirectlyWithStableEnvironment(t *testing.T) {
	root := t.TempDir()
	script := writeTestHook(t, root, "hook.sh", `#!/bin/sh
set -eu
mkdir -p generated
printf '%s\n%s\n%s\n%s\n%s\n%s\n' "$WAILS_PROJECT_DIR" "$WAILS_TARGET_OS" "$WAILS_TARGET_ARCH" "$WAILS_PROFILE" "$WAILS_OUTPUT" "$WAILS_PIPELINE_VERSION" > generated/environment
`)
	spec := pipeline.HookSpec{Phase: manifest.AfterBuild, Script: script, Directory: ".", Profile: "release", TargetOS: "linux", TargetArch: "amd64", ScopeOutput: "bin/app", DeclaredOutputs: []string{"generated/environment"}}
	handler := &manifestHandler{root: root, environment: []string{"WAILS_PROFILE=hostile", "WAILS_OUTPUT=hostile"}}
	result, err := handler.runHook(t.Context(), spec)
	require.NoError(t, err)
	assert.Empty(t, result.Detail)
	content, err := os.ReadFile(filepath.Join(root, "generated", "environment"))
	require.NoError(t, err)
	assert.Equal(t, strings.Join([]string{root, "linux", "amd64", "release", filepath.Join(root, "bin", "app"), "1", ""}, "\n"), string(content))
}

func TestManifestHookFailureRetainsOutputAndExitCode(t *testing.T) {
	root := t.TempDir()
	script := writeTestHook(t, root, "fail.sh", "#!/bin/sh\necho preflight failed\nexit 7\n")
	result, err := (&manifestHandler{root: root}).runHook(t.Context(), pipeline.HookSpec{Phase: manifest.BeforeBuild, Script: script})
	require.Error(t, err)
	assert.Equal(t, "preflight failed", result.Detail)
	var coder interface{ ExitCode() int }
	require.True(t, errors.As(err, &coder))
	assert.Equal(t, 7, coder.ExitCode())
}

func TestManifestHookCancellationTerminatesPromptly(t *testing.T) {
	root := t.TempDir()
	script := writeTestHook(t, root, "wait.sh", "#!/bin/sh\nsleep 30 &\necho $! > child.pid\nwait\n")
	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()
	started := time.Now()
	_, err := (&manifestHandler{root: root}).runHook(ctx, pipeline.HookSpec{Phase: manifest.BeforeBuild, Script: script})
	require.Error(t, err)
	assert.Less(t, time.Since(started), 3*time.Second)
	pidText, readErr := os.ReadFile(filepath.Join(root, "child.pid"))
	require.NoError(t, readErr)
	pid, parseErr := strconv.Atoi(strings.TrimSpace(string(pidText)))
	require.NoError(t, parseErr)
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if signalErr := syscall.Kill(pid, 0); errors.Is(signalErr, syscall.ESRCH) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("hook child process %d survived cancellation", pid)
}

func TestManifestHookCacheHitRestoreAndInvalidation(t *testing.T) {
	root := t.TempDir()
	script := writeTestHook(t, root, "generate.sh", `#!/bin/sh
set -eu
mkdir -p generated
cat version.txt > generated/version.txt
printf 'run\n' >> runs.txt
`)
	require.NoError(t, os.WriteFile(filepath.Join(root, "version.txt"), []byte("one\n"), 0o644))
	node := pipeline.Node{Key: "hook:before_build", Kind: pipeline.RunHook, Label: "Run before build hook", Scope: pipeline.ProjectScope,
		Spec:   pipeline.HookSpec{Phase: manifest.BeforeBuild, Script: script, Profile: "default", DeclaredOutputs: []string{"generated/version.txt"}},
		Inputs: []pipeline.InputSpec{{Label: "hook-script", Files: []string{script}}, {Label: "hook-inputs", Files: []string{"version.txt"}}}, Output: "generated", Cache: pipeline.CacheArtifact}
	plan := pipeline.Plan{Name: "hook", Roots: []pipeline.NodeKey{node.Key}, Nodes: map[pipeline.NodeKey]pipeline.Node{node.Key: node}}
	executor := pipeline.Executor{Handler: &manifestHandler{root: root}}
	run := func() pipeline.Result {
		results, err := executor.Execute(t.Context(), plan, pipeline.ExecuteOptions{Root: root, Workers: 1})
		require.NoError(t, err)
		return results[node.Key]
	}
	assert.Equal(t, "miss", string(run().Status))
	assert.Equal(t, "hit", string(run().Status))
	require.NoError(t, os.RemoveAll(filepath.Join(root, "generated")))
	assert.Equal(t, "restored", string(run().Status))
	restored, err := os.ReadFile(filepath.Join(root, "generated", "version.txt"))
	require.NoError(t, err)
	assert.Equal(t, "one\n", string(restored))
	require.NoError(t, os.WriteFile(filepath.Join(root, "version.txt"), []byte("two\n"), 0o644))
	assert.Equal(t, "miss", string(run().Status))
	require.NoError(t, os.Chmod(filepath.Join(root, filepath.FromSlash(script)), 0o744))
	assert.Equal(t, "miss", string(run().Status))
	executor = pipeline.Executor{Handler: &manifestHandler{root: root, environment: []string{"HOOK_VARIANT=two"}}}
	assert.Equal(t, "miss", string(run().Status))
	runs, err := os.ReadFile(filepath.Join(root, "runs.txt"))
	require.NoError(t, err)
	assert.Equal(t, "run\nrun\nrun\nrun\n", string(runs))
}

func TestManifestHookRejectsSymlinksInCachedOutputs(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "source.txt"), []byte("source\n"), 0o644))
	script := writeTestHook(t, root, "link.sh", `#!/bin/sh
set -eu
mkdir -p generated
ln -s ../source.txt generated/link.txt
`)
	_, err := (&manifestHandler{root: root}).runHook(t.Context(), pipeline.HookSpec{
		Phase:           manifest.BeforeBuild,
		Script:          script,
		DeclaredOutputs: []string{"generated/link.txt"},
	})
	require.ErrorContains(t, err, "output contains unsupported symlink")
}

func TestManifestHookRejectsSpecialFilesInCachedOutputs(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "generated"), 0o755))
	require.NoError(t, syscall.Mkfifo(filepath.Join(root, "generated", "pipe"), 0o600))
	script := writeTestHook(t, root, "noop.sh", "#!/bin/sh\n:\n")
	_, err := (&manifestHandler{root: root}).runHook(t.Context(), pipeline.HookSpec{
		Phase:           manifest.BeforeBuild,
		Script:          script,
		DeclaredOutputs: []string{"generated"},
	})
	require.ErrorContains(t, err, "output contains unsupported file type")
}

func TestManifestHookPlanUsesLifecyclePhaseAndScopeOutput(t *testing.T) {
	node := pipeline.Node{Key: "hook:after_build:linux-amd64", Kind: pipeline.RunHook, Scope: pipeline.TargetScope, Cache: pipeline.CacheNever,
		Spec: pipeline.HookSpec{Phase: manifest.AfterBuild, ScopeOutput: ".wails/dev/linux-amd64/app"}}
	plan := pipeline.Plan{Intent: pipeline.BuildIntent{Command: "build", Targets: []pipeline.TargetIntent{{Target: pipeline.Target{OS: "linux", Arch: "amd64"}}}}, Nodes: map[pipeline.NodeKey]pipeline.Node{node.Key: node}}
	output := planOutput(plan)
	require.Len(t, output.Operations, 1)
	assert.Equal(t, "after_build", output.Operations[0].Stage)
	assert.Equal(t, ".wails/dev/linux-amd64/app", output.Operations[0].Output)
}

func writeTestHook(t *testing.T, root, name, content string) string {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "scripts"), 0o755))
	path := filepath.Join(root, "scripts", name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o755))
	relative, err := filepath.Rel(root, path)
	require.NoError(t, err)
	return filepath.ToSlash(relative)
}
