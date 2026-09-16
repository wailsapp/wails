package pipeline

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/cache"
)

func TestCompileInputsTrackNativeSources(t *testing.T) {
	for _, local := range []bool{false, true} {
		scope := "project"
		if local {
			scope = "local-module"
		}
		t.Run(scope, func(t *testing.T) {
			for _, extension := range []string{".c", ".cc", ".cpp", ".cxx", ".m", ".h", ".hh", ".hpp", ".hxx", ".f", ".F", ".for", ".f90", ".s", ".S", ".sx", ".swig", ".swigcxx", ".syso"} {
				t.Run(extension, func(t *testing.T) {
					config := testConfig(t)
					sourceRoot := config.Root
					if local {
						sourceRoot = t.TempDir()
						require.NoError(t, os.WriteFile(filepath.Join(sourceRoot, "go.mod"), []byte("module example/native\n"), 0600))
						require.NoError(t, os.WriteFile(filepath.Join(config.Root, "go.mod"), []byte("module example/app\nreplace example/native => "+strconv.Quote(filepath.ToSlash(sourceRoot))+"\n"), 0600))
					}
					path := filepath.Join(sourceRoot, "helper"+extension)
					require.NoError(t, os.WriteFile(path, []byte("before"), 0600))
					plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
					require.NoError(t, err)
					node, ok := plan.Nodes["target:linux/amd64:compile"]
					require.True(t, ok)
					store, err := cache.OpenCache(config.Root)
					require.NoError(t, err)
					before, err := snapshotNodeInputs(store, node)
					require.NoError(t, err)
					require.NoError(t, os.WriteFile(path, []byte("after!"), 0600))
					store.InvalidateObservations()
					after, err := snapshotNodeInputs(store, node)
					require.NoError(t, err)
					require.NotEqual(t, before, after, "native-source edit must invalidate compilation")
				})
			}
		})
	}
}

func TestExampleCompileTracksParentModuleSources(t *testing.T) {
	config := testConfig(t)
	module := config.Root
	config.Root = filepath.Join(module, "examples", "app")
	require.NoError(t, os.MkdirAll(config.Root, 0700))
	require.NoError(t, os.WriteFile(filepath.Join(module, "go.mod"), []byte("module example.com/parent\n"), 0600))
	shared := filepath.Join(module, "shared.go")
	require.NoError(t, os.WriteFile(shared, []byte("package parent\nconst Value = 1\n"), 0600))
	config.Frontend.Disabled = true
	plan, err := PlanBuild(config, Request{TargetOS: "linux", TargetArch: "amd64"})
	require.NoError(t, err)
	node := plan.Nodes["target:linux/amd64:compile"]
	store, err := cache.OpenCache(config.Root)
	require.NoError(t, err)
	before, err := snapshotNodeInputs(store, node)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(shared, []byte("package parent\nconst Value = 2\n"), 0600))
	store.InvalidateObservations()
	after, err := snapshotNodeInputs(store, node)
	require.NoError(t, err)
	require.NotEqual(t, before, after, "parent module source changes must invalidate an example's build")
}

func TestRunAndroidPlanUsesProductionTagsAndInstallableAPK(t *testing.T) {
	config := testConfig(t)
	config.Build.Go.Tags = []string{"devtools", "devtools"}
	plan, err := PlanBuild(config, Request{Verb: "run", TargetOS: "android", TargetArch: "arm64"})
	require.NoError(t, err)
	apk := false
	for _, node := range plan.Nodes {
		switch spec := node.Spec.(type) {
		case CompileSpec:
			require.True(t, spec.Production)
			require.Equal(t, []string{"devtools", "production", "android"}, spec.Tags)
		case PackageSpec:
			if spec.Format == "apk" {
				apk = true
				require.False(t, spec.Development)
			}
		}
	}
	require.True(t, apk, "run must produce an installable APK rather than an AAB")
	_, err = PlanBuild(config, Request{Verb: "build", TargetOS: "android", TargetArch: "arm64", Formats: []string{"apk"}})
	require.ErrorContains(t, err, "production APK")
}

// Exercise the planner's inputs through the executor's actual warm cache with
// a real CGo binary, rather than only comparing snapshot hashes.
func TestNativeSourceEditRebuildsWarmExecutable(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux CGo acceptance fixture")
	}
	if _, err := exec.LookPath("g++"); err != nil {
		t.Skip("requires g++")
	}
	for _, local := range []bool{false, true} {
		t.Run(fmt.Sprintf("local=%t", local), func(t *testing.T) {
			config := testConfig(t)
			sourceRoot := config.Root
			main := "package main\n/* int answer(); */\nimport \"C\"\nimport \"fmt\"\nfunc main() { fmt.Print(C.answer()) }\n"
			module := "module example/app\ngo 1.24\n"
			if local {
				sourceRoot = t.TempDir()
				require.NoError(t, os.WriteFile(filepath.Join(sourceRoot, "go.mod"), []byte("module example/native\ngo 1.24\n"), 0600))
				require.NoError(t, os.WriteFile(filepath.Join(sourceRoot, "native.go"), []byte("package native\n/* int answer(); */\nimport \"C\"\nfunc Answer() int { return int(C.answer()) }\n"), 0600))
				module += "require example/native v0.0.0\nreplace example/native => " + filepath.ToSlash(sourceRoot) + "\n"
				main = "package main\nimport (\"fmt\"; \"example/native\")\nfunc main() { fmt.Print(native.Answer()) }\n"
			}
			require.NoError(t, os.WriteFile(filepath.Join(config.Root, "go.mod"), []byte(module), 0600))
			require.NoError(t, os.WriteFile(filepath.Join(config.Root, "main.go"), []byte(main), 0600))
			source := filepath.Join(sourceRoot, "helper.cxx")
			require.NoError(t, os.WriteFile(source, []byte("extern \"C\" int answer() { return 1; }"), 0600))
			full, err := PlanBuild(config, Request{Verb: "build", TargetOS: runtime.GOOS, TargetArch: runtime.GOARCH})
			require.NoError(t, err)
			key := NodeKey("target:" + runtime.GOOS + "/" + runtime.GOARCH + ":compile")
			node := full.Nodes[key]
			node.Dependencies = nil
			plan := Plan{Roots: []NodeKey{key}, Artifacts: []NodeKey{key}, Nodes: map[NodeKey]Node{key: node}}
			executor := Executor{Handler: nativeCompileHandler{root: config.Root}}
			build := func(want cache.LookupStatus, answer string) {
				results, err := executor.Execute(context.Background(), plan, ExecuteOptions{Root: config.Root})
				require.NoError(t, err)
				require.Equal(t, want, results[key].Status)
				output, err := exec.Command(filepath.Join(config.Root, node.Output)).CombinedOutput()
				require.NoError(t, err, string(output))
				require.Equal(t, answer, string(output))
			}
			build(cache.LookupMiss, "1")
			build(cache.LookupHit, "1")
			require.NoError(t, os.WriteFile(source, []byte("extern \"C\" int answer() { return 2; }"), 0600))
			build(cache.LookupMiss, "2")
			build(cache.LookupHit, "2")
		})
	}
}

type nativeCompileHandler struct{ root string }

func (nativeCompileHandler) Identity(context.Context, Node) (string, error) {
	return "native-cxx-test-v1", nil
}
func (h nativeCompileHandler) Run(ctx context.Context, node Node) (RunResult, error) {
	command := exec.CommandContext(ctx, "go", "build", "-o", node.Output, ".")
	command.Dir = h.root
	command.Env = append(os.Environ(), "CGO_ENABLED=1", "GOWORK=off")
	output, err := command.CombinedOutput()
	return RunResult{Detail: string(output)}, err
}

func TestDisabledFrontendStillTracksGoSourcesInFrontendDirectory(t *testing.T) {
	config := testConfig(t)
	config.Frontend.Disabled = true
	require.NoError(t, os.WriteFile(filepath.Join(config.Root, "go.mod"), []byte("module example.com/disabled\n"), 0600))
	directory := filepath.Join(config.Root, config.Frontend.Directory)
	require.NoError(t, os.MkdirAll(directory, 0700))
	source := filepath.Join(directory, "helper.go")
	require.NoError(t, os.WriteFile(source, []byte("package frontend\nconst Value = 1\n"), 0600))
	plan, err := PlanBuild(config, Request{Verb: "run", TargetOS: "linux", TargetArch: "amd64"})
	require.NoError(t, err)
	node := plan.Nodes["target:linux/amd64:compile"]
	store, err := cache.OpenCache(config.Root)
	require.NoError(t, err)
	before, err := snapshotNodeInputs(store, node)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(source, []byte("package frontend\nconst Value = 2\n"), 0600))
	store.InvalidateObservations()
	after, err := snapshotNodeInputs(store, node)
	require.NoError(t, err)
	require.NotEqual(t, before, after, "disabled frontend directory may contain imported Go packages")
}
