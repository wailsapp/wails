//go:build !windows

package commands

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/report"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestManifestHookReceivesStableEnvironmentAndWorkingDirectory(t *testing.T) {
	for _, key := range []string{"WAILS_PROJECT_DIR", "WAILS_TARGET_OS", "WAILS_TARGET_ARCH", "WAILS_PROFILE", "WAILS_OUTPUT", "WAILS_PIPELINE_VERSION"} {
		t.Setenv(key, "inherited-value-must-not-win")
	}
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "scripts"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "work"), 0o755))
	script := filepath.Join(root, "scripts", "record hook.sh")
	require.NoError(t, os.WriteFile(script, []byte(`#!/bin/sh
printf '%s\n%s\n%s\n%s\n%s\n%s\n%s\n' "$PWD" "$WAILS_PROJECT_DIR" "$WAILS_TARGET_OS" "$WAILS_TARGET_ARCH" "$WAILS_PROFILE" "$WAILS_OUTPUT" "$WAILS_PIPELINE_VERSION" > "$WAILS_PROJECT_DIR/hook.env"
`), 0o755))

	node := pipeline.Node{Key: "hook", Kind: pipeline.RunHook, Spec: pipeline.HookSpec{
		Phase: "after_build", TargetOS: "linux", TargetArch: "arm64", Profile: "release", Artifact: "bin/example",
		Hook: manifest.Hook{Script: "scripts/record hook.sh", Directory: "work"},
	}}
	_, err := (&manifestHandler{root: root}).Run(context.Background(), node)
	require.NoError(t, err)
	data, err := os.ReadFile(filepath.Join(root, "hook.env"))
	require.NoError(t, err)
	resolvedWork, err := filepath.EvalSymlinks(filepath.Join(root, "work"))
	require.NoError(t, err)
	assert.Equal(t, resolvedWork+"\n"+root+"\nlinux\narm64\nrelease\nbin/example\n1\n", string(data))
}

func TestCachedManifestHookRerunsOnlyWhenDeclaredInputsChange(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "scripts"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "version.txt"), []byte("one\n"), 0o644))
	script := filepath.Join(root, "scripts", "generate.sh")
	require.NoError(t, os.WriteFile(script, []byte(`#!/bin/sh
mkdir -p "$WAILS_PROJECT_DIR/generated"
count=0
if [ -f "$WAILS_PROJECT_DIR/generated/count" ]; then count=$(cat "$WAILS_PROJECT_DIR/generated/count"); fi
count=$((count + 1))
printf '%s' "$count" > "$WAILS_PROJECT_DIR/generated/count"
`), 0o755))

	node := pipeline.Node{Key: "hook:before_build:linux/amd64", Kind: pipeline.RunHook,
		Spec:   pipeline.HookSpec{Phase: "before_build", TargetOS: "linux", TargetArch: "amd64", Profile: "default", Artifact: "bin/app", Hook: manifest.Hook{Script: "scripts/generate.sh", Cache: true, Inputs: []string{"version.txt"}, Outputs: []string{"generated"}}},
		Inputs: []pipeline.InputSpec{{Label: "hook-script", Files: []string{"scripts/generate.sh", "version.txt"}}},
		Output: "generated", Cache: pipeline.CacheArtifact,
	}
	plan := pipeline.Plan{Name: "hook", Target: "linux/amd64", Nodes: map[pipeline.NodeKey]pipeline.Node{node.Key: node}, Roots: []pipeline.NodeKey{node.Key}}
	executor := pipeline.Executor{Handler: &manifestHandler{root: root}}
	_, err := executor.Execute(context.Background(), plan, pipeline.ExecuteOptions{Root: root, Reporter: report.Nop{}})
	require.NoError(t, err)
	_, err = executor.Execute(context.Background(), plan, pipeline.ExecuteOptions{Root: root, Reporter: report.Nop{}})
	require.NoError(t, err)
	data, err := os.ReadFile(filepath.Join(root, "generated", "count"))
	require.NoError(t, err)
	assert.Equal(t, "1", string(data))

	require.NoError(t, os.Rename(filepath.Join(root, "generated"), filepath.Join(root, "removed-generated")))
	_, err = executor.Execute(context.Background(), plan, pipeline.ExecuteOptions{Root: root, Reporter: report.Nop{}})
	require.NoError(t, err)
	data, err = os.ReadFile(filepath.Join(root, "generated", "count"))
	require.NoError(t, err)
	assert.Equal(t, "1", string(data), "missing hook outputs are restored without executing the hook")

	require.NoError(t, os.WriteFile(filepath.Join(root, "version.txt"), []byte("two\n"), 0o644))
	_, err = executor.Execute(context.Background(), plan, pipeline.ExecuteOptions{Root: root, Reporter: report.Nop{}})
	require.NoError(t, err)
	data, err = os.ReadFile(filepath.Join(root, "generated", "count"))
	require.NoError(t, err)
	assert.Equal(t, "2", string(data))

	file, err := os.OpenFile(script, os.O_APPEND|os.O_WRONLY, 0)
	require.NoError(t, err)
	_, err = file.WriteString("# implementation changed\n")
	require.NoError(t, err)
	require.NoError(t, file.Close())
	_, err = executor.Execute(context.Background(), plan, pipeline.ExecuteOptions{Root: root, Reporter: report.Nop{}})
	require.NoError(t, err)
	data, err = os.ReadFile(filepath.Join(root, "generated", "count"))
	require.NoError(t, err)
	assert.Equal(t, "3", string(data), "script content is an implicit cache input")
}

func TestDefaultManifestHookAlwaysRuns(t *testing.T) {
	root := t.TempDir()
	script := filepath.Join(root, "run.sh")
	require.NoError(t, os.WriteFile(script, []byte(`#!/bin/sh
count=0
if test -f count; then count=$(cat count); fi
printf '%s' "$((count + 1))" > count
`), 0o755))
	node := pipeline.Node{Key: "hook:after_build:linux/amd64", Kind: pipeline.RunHook, Spec: pipeline.HookSpec{Phase: "after_build", TargetOS: "linux", TargetArch: "amd64", Profile: "default", Artifact: "bin/app", Hook: manifest.Hook{Script: "run.sh"}}, Cache: pipeline.CacheNever}
	plan := pipeline.Plan{Name: "hook", Target: "linux/amd64", Nodes: map[pipeline.NodeKey]pipeline.Node{node.Key: node}, Roots: []pipeline.NodeKey{node.Key}}
	executor := pipeline.Executor{Handler: &manifestHandler{root: root}}
	for range 2 {
		_, err := executor.Execute(context.Background(), plan, pipeline.ExecuteOptions{Root: root, Reporter: report.Nop{}})
		require.NoError(t, err)
	}
	data, err := os.ReadFile(filepath.Join(root, "count"))
	require.NoError(t, err)
	assert.Equal(t, "2", string(data))
}

func TestManifestHookRejectsSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.sh")
	require.NoError(t, os.WriteFile(outside, []byte("#!/bin/sh\n"), 0o755))
	require.NoError(t, os.Symlink(outside, filepath.Join(root, "escaped.sh")))

	_, err := (&manifestHandler{root: root}).hook(context.Background(), pipeline.HookSpec{Hook: manifest.Hook{Script: "escaped.sh"}})
	require.ErrorContains(t, err, "path resolves outside the project")
}

func TestManifestHookCancellationTerminatesChildProcessGroup(t *testing.T) {
	root := t.TempDir()
	script := filepath.Join(root, "wait.sh")
	require.NoError(t, os.WriteFile(script, []byte(`#!/bin/sh
sleep 30 &
printf '%s' "$!" > child.pid
wait
`), 0o755))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		_, err := (&manifestHandler{root: root}).hook(ctx, pipeline.HookSpec{Hook: manifest.Hook{Script: "wait.sh"}})
		done <- err
	}()

	var pid int
	require.Eventually(t, func() bool {
		data, err := os.ReadFile(filepath.Join(root, "child.pid"))
		if err != nil {
			return false
		}
		pid, err = strconv.Atoi(string(data))
		return err == nil && pid > 0
	}, 2*time.Second, 10*time.Millisecond)
	cancel()
	select {
	case err := <-done:
		require.Error(t, err)
	case <-time.After(3 * time.Second):
		t.Fatal("cancelled hook did not return")
	}
	require.Eventually(t, func() bool {
		return syscall.Kill(pid, 0) == syscall.ESRCH
	}, 2*time.Second, 10*time.Millisecond, "hook child process %d was left behind", pid)
}

func TestManifestHookFailureKeepsCapturedOutput(t *testing.T) {
	root := t.TempDir()
	script := filepath.Join(root, "fail.sh")
	require.NoError(t, os.WriteFile(script, []byte("#!/bin/sh\nprintf 'specific hook failure\\n'\nexit 7\n"), 0o755))

	result, err := (&manifestHandler{root: root}).hook(context.Background(), pipeline.HookSpec{Hook: manifest.Hook{Script: "fail.sh"}})
	require.Error(t, err)
	assert.Equal(t, "specific hook failure", result.Detail)

	reporter := &hookFailureReporter{}
	node := pipeline.Node{Key: "hook:after_build:linux/amd64", Kind: pipeline.RunHook, Spec: pipeline.HookSpec{Phase: "after_build", Hook: manifest.Hook{Script: "fail.sh"}}, Cache: pipeline.CacheNever}
	plan := pipeline.Plan{Name: "hook", Nodes: map[pipeline.NodeKey]pipeline.Node{node.Key: node}, Roots: []pipeline.NodeKey{node.Key}}
	_, err = (pipeline.Executor{Handler: &manifestHandler{root: root}}).Execute(context.Background(), plan, pipeline.ExecuteOptions{Root: root, Reporter: reporter})
	require.Error(t, err)
	assert.Equal(t, 7, reporter.failure.ExitCode)
	assert.Equal(t, "specific hook failure", reporter.failure.Output)
}

type hookFailureReporter struct {
	report.Nop
	failure report.Failure
}

func (r *hookFailureReporter) StepFailed(_ report.StepID, failure report.Failure) {
	r.failure = failure
}
