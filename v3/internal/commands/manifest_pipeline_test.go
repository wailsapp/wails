package commands

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/wake"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestRunManifestPipelineMarksRenderedFailure(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	require.NoError(t, manifest.WriteMinimal(root, manifest.Project{
		Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0",
	}))
	tools := t.TempDir()
	for _, name := range []string{"go", "npm", "cc"} {
		require.NoError(t, os.WriteFile(filepath.Join(tools, name), []byte("#!/bin/sh\nexit 7\n"), 0o755))
	}
	t.Setenv("PATH", tools)

	err := runManifestPipeline(manifestRunOptions{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	require.Error(t, err)
	assert.True(t, wake.IsReported(err), "execution failure was rendered by Pulse and must not be printed again by the CLI")
}

func TestRunManifestPipelineLeavesUnrenderedFailurePrintable(t *testing.T) {
	t.Chdir(t.TempDir())

	err := runManifestPipeline(manifestRunOptions{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	require.Error(t, err)
	assert.False(t, wake.IsReported(err), "errors raised before Pulse starts still need the CLI error printer")
}

func TestRebaseGeneratedOverlayRejectsReplacementOutsideStaging(t *testing.T) {
	root := t.TempDir()
	staged := filepath.Join(root, "stage")
	require.NoError(t, os.MkdirAll(staged, 0o755))
	overlay := filepath.Join(staged, "overlay.json")
	require.NoError(t, os.WriteFile(overlay, []byte(`{"Replace":{"virtual.go":"/outside/generated.go"}}`), 0o644))

	err := rebaseGeneratedOverlay(overlay, staged, filepath.Join(root, "published"))
	assert.ErrorContains(t, err, "outside staging root")
}

func TestResolveManifestPlanUsesTheDiscoveredManifestRoot(t *testing.T) {
	prependFakePlanTools(t, "npm")
	root := t.TempDir()
	require.NoError(t, manifest.WriteMinimal(root, manifest.Project{
		Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0",
	}))
	nested := filepath.Join(root, "frontend", "src")
	require.NoError(t, os.MkdirAll(nested, 0o755))
	loaded, err := manifest.Load(nested, "")
	require.NoError(t, err)
	t.Chdir(nested)

	resolved, _, _, err := resolveManifestPlan(manifestRunOptions{Verb: "build", Loaded: loaded})
	require.NoError(t, err)
	assert.Equal(t, root, resolved)
}

func TestResolveManifestPlanPropagatesAnonymousGarbleArgumentsOnce(t *testing.T) {
	root := t.TempDir()
	tools := t.TempDir()
	for _, name := range []string{"go", "garble", "npm", "cc"} {
		require.NoError(t, os.WriteFile(filepath.Join(tools, name), []byte("#!/bin/sh\nexit 0\n"), 0o755))
	}
	t.Setenv("PATH", tools)
	require.NoError(t, manifest.WriteMinimal(root, manifest.Project{
		Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0",
	}))
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	loaded.Config.Build.Obfuscation = true
	loaded.Config.Build.Go.GarbleArgs = []string{"-tiny"}

	_, _, plan, err := resolveManifestPlan(manifestRunOptions{
		Verb: "build", Loaded: loaded, TargetOS: "linux", TargetArch: "amd64", GarbleArgs: []string{"-literals"},
	})
	require.NoError(t, err)
	compile := plan.Nodes["target:linux/amd64:compile"].Spec.(pipeline.CompileSpec)
	assert.Equal(t, []string{"-tiny", "-literals"}, compile.GarbleArgs)
}

func TestResolveManifestPlanReturnsAdapterFailures(t *testing.T) {
	want := errors.New("injected failure")
	loaded := &manifest.Loaded{Config: manifest.Config{}}
	valid := manifestPlanOperations{
		getwd: func() (string, error) { return t.TempDir(), nil },
		load:  func(string, string) (*manifest.Loaded, error) { return loaded, nil },
		plan:  func(manifest.Config, pipeline.Request) (pipeline.Plan, error) { return pipeline.Plan{}, nil },
	}

	operations := valid
	operations.getwd = func() (string, error) { return "", want }
	_, _, _, err := resolveManifestPlanWithOperations(manifestRunOptions{}, operations)
	require.ErrorIs(t, err, want)

	operations = valid
	operations.load = func(string, string) (*manifest.Loaded, error) { return nil, want }
	_, _, _, err = resolveManifestPlanWithOperations(manifestRunOptions{}, operations)
	require.ErrorIs(t, err, want)

	operations = valid
	operations.plan = func(manifest.Config, pipeline.Request) (pipeline.Plan, error) { return pipeline.Plan{}, want }
	_, _, _, err = resolveManifestPlanWithOperations(manifestRunOptions{}, operations)
	require.ErrorIs(t, err, want)
}

func TestResolvedPlanCompilersSkipsMalformedCompileNodesAndSortsTargets(t *testing.T) {
	plan := pipeline.Plan{Nodes: map[pipeline.NodeKey]pipeline.Node{
		"bad": {
			Kind: pipeline.CompileApplication,
			Spec: pipeline.FrontendSpec{},
		},
		"windows": {
			Kind: pipeline.CompileApplication,
			Spec: pipeline.CompileSpec{TargetOS: "windows", TargetArch: "amd64", Tags: []string{"release"}},
		},
		"linux": {
			Kind: pipeline.CompileApplication,
			Spec: pipeline.CompileSpec{TargetOS: "linux", TargetArch: "arm64", Obfuscated: true, GarbleArgs: []string{"-tiny"}},
		},
		"frontend": {Kind: pipeline.BuildFrontend},
	}}

	assert.Equal(t, []manifestPlanCompiler{
		{Target: "linux/arm64", Obfuscated: true, GarbleArgs: []string{"-tiny"}},
		{Target: "windows/amd64", Tags: []string{"release"}},
	}, resolvedPlanCompilers(plan))
}

func TestPlanDocumentIsDerivedEntirelyFromTheResolvedImmutablePlan(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, manifest.WriteMinimal(root, manifest.Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.2.3"}))
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	config := loaded.Config
	config.Selected = manifest.Profile{Name: "release", Targets: []manifest.ProfileTarget{{Target: "linux/amd64", Formats: []string{"deb"}, Sign: true}}}
	config.Profile = "release"
	config.Signing.Linux.Enabled = true
	plan, err := pipeline.PlanBuild(config, pipeline.Request{Verb: "build"})
	require.NoError(t, err)

	first := planOutput(plan)
	second := planOutput(plan)
	assert.Equal(t, first, second)
	assert.Equal(t, "build", first.Request.Command)
	assert.Equal(t, "release", first.Request.Profile)
	assert.Equal(t, []string{"linux/amd64"}, first.Request.Targets)
	assert.Equal(t, []string{"deb"}, first.Request.Formats)
	require.Len(t, first.Artifacts, 1)
	assert.Equal(t, manifestPlanArtifact{Target: "linux/amd64", Format: "deb", Kind: "deb", Path: "bin/app_1.2.3_amd64.deb.signed", Signed: true}, first.Artifacts[0])
	for _, operation := range first.Operations {
		if operation.Cache == string(pipeline.CacheNever) {
			assert.Equal(t, "run", operation.Decision)
		} else {
			assert.Equal(t, "cache-check", operation.Decision)
		}
	}
	var collect manifestPlanOperation
	for _, operation := range first.Operations {
		if operation.ID == "collect:artifacts" {
			collect = operation
			break
		}
	}
	assert.Equal(t, "collect", collect.Stage)
}

func TestManifestHandlerCollectVerifiesEveryFinalArtifact(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "bin", "App.app"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, ".wails", "dev"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "bin", "app"), []byte("binary"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".wails", "dev", "app"), []byte("binary"), 0o755))
	handler := &manifestHandler{root: root}
	spec := pipeline.CollectSpec{Receipt: ".wails/artifacts/receipt.json", Artifacts: []pipeline.ArtifactReference{
		{Key: "linux", Path: "bin/app"},
		{Key: "darwin", Path: "bin/App.app"},
		{Key: "development", Path: ".wails/dev/app"},
	}}

	result, err := handler.collect(spec)
	require.NoError(t, err)
	assert.Equal(t, "collected 3 artifact(s)", result.Detail)

	spec.Artifacts = append(spec.Artifacts, pipeline.ArtifactReference{Key: "missing", Path: "bin/missing"})
	_, err = handler.collect(spec)
	require.ErrorContains(t, err, "collect artifact missing")

	outside := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(outside, "app"), []byte("outside"), 0o755))
	require.NoError(t, os.Symlink(outside, filepath.Join(root, ".wails", "outside")))
	_, err = handler.collect(pipeline.CollectSpec{Artifacts: []pipeline.ArtifactReference{{Key: "escape", Path: ".wails/outside/app"}}})
	require.ErrorContains(t, err, "resolves outside the project")
}

func TestManifestHandlerPublishesFilesAndDirectoriesWithoutMutatingStagedSources(t *testing.T) {
	root := t.TempDir()
	handler := &manifestHandler{root: root}
	stagedFile := filepath.Join(root, ".wails", "stage", "app")
	require.NoError(t, os.MkdirAll(filepath.Dir(stagedFile), 0o755))
	require.NoError(t, os.WriteFile(stagedFile, []byte("new"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "bin"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "bin", "app"), []byte("old"), 0o644))

	result, err := handler.publish(pipeline.PublishSpec{Source: ".wails/stage/app", Destination: "bin/app"})
	require.NoError(t, err)
	assert.Equal(t, "bin/app", result.Output)
	stagedData, err := os.ReadFile(stagedFile)
	require.NoError(t, err)
	publishedData, err := os.ReadFile(filepath.Join(root, "bin", "app"))
	require.NoError(t, err)
	assert.Equal(t, []byte("new"), stagedData)
	assert.Equal(t, stagedData, publishedData)

	stagedDirectory := filepath.Join(root, ".wails", "stage", "App.app")
	require.NoError(t, os.MkdirAll(stagedDirectory, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(stagedDirectory, "new"), []byte("bundle"), 0o644))
	destination := filepath.Join(root, "bin", "App.app")
	require.NoError(t, os.MkdirAll(destination, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(destination, "old"), []byte("old"), 0o644))
	_, err = handler.publish(pipeline.PublishSpec{Source: ".wails/stage/App.app", Destination: "bin/App.app"})
	require.NoError(t, err)
	assert.FileExists(t, filepath.Join(stagedDirectory, "new"))
	assert.FileExists(t, filepath.Join(destination, "new"))
	assert.NoFileExists(t, filepath.Join(destination, "old"))
}

func TestManifestHandlerPublishEnforcesGeneratedSourceAndUserDestinationBoundaries(t *testing.T) {
	root := t.TempDir()
	handler := &manifestHandler{root: root}
	require.NoError(t, os.MkdirAll(filepath.Join(root, ".wails", "stage"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".wails", "stage", "app"), []byte("app"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "user-app"), []byte("user"), 0o644))

	_, err := handler.publish(pipeline.PublishSpec{Source: "user-app", Destination: "bin/app"})
	assert.ErrorContains(t, err, "not in the generated .wails workspace")
	_, err = handler.publish(pipeline.PublishSpec{Source: ".wails/stage/app", Destination: ".wails/final/app"})
	assert.ErrorContains(t, err, "inside the generated .wails workspace")
	_, err = handler.publish(pipeline.PublishSpec{Source: ".wails/stage/app", Destination: "../app"})
	assert.ErrorContains(t, err, "must be project-relative")

	outside := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(outside, "app"), []byte("outside"), 0o644))
	require.NoError(t, os.Symlink(filepath.Join(outside, "app"), filepath.Join(root, ".wails", "stage", "escape")))
	_, err = handler.publish(pipeline.PublishSpec{Source: ".wails/stage/escape", Destination: "bin/app"})
	assert.ErrorContains(t, err, "resolves outside the project")
}

func TestTransactionalReplacementRestoresPreviousWorkspaceWhenCommitRenameFails(t *testing.T) {
	root := t.TempDir()
	destination := filepath.Join(root, "workspace")
	staged := filepath.Join(root, "staged")
	require.NoError(t, os.MkdirAll(destination, 0o755))
	require.NoError(t, os.MkdirAll(staged, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(destination, "value"), []byte("previous"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(staged, "value"), []byte("next"), 0o644))
	want := errors.New("injected commit rename failure")
	renames := 0
	operations := replacePathOperations{mkdirTemp: os.MkdirTemp, lstat: os.Lstat, removeAll: os.RemoveAll, rename: func(old, next string) error {
		renames++
		if renames == 2 {
			return want
		}
		return os.Rename(old, next)
	}}

	err := replacePathTransactionalWithOperations(staged, destination, operations)
	require.ErrorIs(t, err, want)
	assert.Equal(t, "previous", string(readTestFile(t, filepath.Join(destination, "value"))))
	assert.Equal(t, "next", string(readTestFile(t, filepath.Join(staged, "value"))))
}

func TestTransactionalReplacementReportsRollbackFailure(t *testing.T) {
	root := t.TempDir()
	destination := filepath.Join(root, "workspace")
	staged := filepath.Join(root, "staged")
	require.NoError(t, os.MkdirAll(destination, 0o755))
	require.NoError(t, os.MkdirAll(staged, 0o755))
	commitErr := errors.New("commit failed")
	rollbackErr := errors.New("rollback failed")
	renames := 0
	operations := replacePathOperations{mkdirTemp: os.MkdirTemp, lstat: os.Lstat, removeAll: os.RemoveAll, rename: func(old, next string) error {
		renames++
		switch renames {
		case 1:
			return os.Rename(old, next)
		case 2:
			return commitErr
		default:
			return rollbackErr
		}
	}}

	err := replacePathTransactionalWithOperations(staged, destination, operations)
	require.ErrorIs(t, err, commitErr)
	require.ErrorIs(t, err, rollbackErr)
	assert.ErrorContains(t, err, "restore previous generated workspace")
}

func TestManifestSigningNeverMutatesItsInputArtifact(t *testing.T) {
	root := t.TempDir()
	input := filepath.Join(root, ".wails", "artifacts", "example.deb")
	require.NoError(t, os.MkdirAll(filepath.Dir(input), 0o755))
	require.NoError(t, os.WriteFile(input, []byte("unsigned"), 0o644))
	previousSign := manifestSign
	manifestSign = func(options *flags.Sign) error {
		return os.WriteFile(options.Input, []byte("signed"), 0o644)
	}
	t.Cleanup(func() { manifestSign = previousSign })

	_, err := (&manifestHandler{root: root}).sign(t.Context(), pipeline.SignSpec{
		TargetOS: "linux", Format: "deb", Input: ".wails/artifacts/example.deb",
		Config: manifest.SigningPlatform{Enabled: true, Certificate: "release@example.com"},
	})
	require.NoError(t, err)
	assert.Equal(t, "unsigned", string(readTestFile(t, input)))
	assert.Equal(t, "signed", string(readTestFile(t, input+".signed")))
}

func TestManifestCompileUsesResolvedZigToolchain(t *testing.T) {
	root := t.TempDir()
	tools := t.TempDir()
	record := filepath.Join(root, "zig-compile.txt")
	script := "#!/bin/sh\nprintf '%s\\n' \"$CC\" \"$CXX\" \"$@\" > \"$COMPILE_RECORD\"\noutput=\nwhile [ \"$#\" -gt 0 ]; do if [ \"$1\" = \"-o\" ]; then output=$2; shift; fi; shift; done\nmkdir -p \"$(dirname \"$output\")\"\nprintf binary > \"$output\"\n"
	require.NoError(t, os.WriteFile(filepath.Join(tools, "go"), []byte(script), 0o755))
	t.Setenv("PATH", tools)
	t.Setenv("COMPILE_RECORD", record)
	handler := &manifestHandler{root: root}
	_, err := handler.compile(t.Context(), pipeline.CompileSpec{TargetOS: "linux", TargetArch: "arm64", Output: ".wails/build/app", Toolchain: "zig", Tags: []string{"production"}})
	require.NoError(t, err)
	data, err := os.ReadFile(record)
	require.NoError(t, err)
	text := string(data)
	assert.Contains(t, text, "zig cc -target aarch64-linux-gnu")
	assert.Contains(t, text, "zig c++ -target aarch64-linux-gnu")
	assert.Contains(t, text, "-tags\nproduction")
	assert.Contains(t, text, "-o\n"+filepath.Join(root, ".wails", "build", ".compile-output-"))
	assert.FileExists(t, filepath.Join(root, ".wails", "build", "app"))
}

func TestManifestCompileFailurePreservesLastCompleteArtifact(t *testing.T) {
	root := t.TempDir()
	tools := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(tools, "go"), []byte("#!/bin/sh\noutput=\nwhile [ \"$#\" -gt 0 ]; do if [ \"$1\" = \"-o\" ]; then output=$2; shift; fi; shift; done\nmkdir -p \"$(dirname \"$output\")\"\nprintf partial > \"$output\"\nexit 9\n"), 0o755))
	t.Setenv("PATH", tools)
	output := filepath.Join(root, ".wails", "build", "app")
	require.NoError(t, os.MkdirAll(filepath.Dir(output), 0o755))
	require.NoError(t, os.WriteFile(output, []byte("last complete"), 0o755))

	_, err := (&manifestHandler{root: root}).compile(t.Context(), pipeline.CompileSpec{TargetOS: "linux", TargetArch: "amd64", Output: ".wails/build/app"})
	require.Error(t, err)
	assert.Equal(t, "last complete", string(readTestFile(t, output)))
	staging, globErr := filepath.Glob(filepath.Join(filepath.Dir(output), ".compile-output-*"))
	require.NoError(t, globErr)
	assert.Empty(t, staging)
}

func TestManifestCompileUsesResolvedDockerToolchainAndReadOnlyLocalMounts(t *testing.T) {
	root := t.TempDir()
	localRoot := t.TempDir()
	tools := t.TempDir()
	record := filepath.Join(root, "docker-compile.txt")
	script := "#!/bin/sh\nprintf '%s\\n' \"$@\" > \"$DOCKER_RECORD\"\noutput=\nwhile [ \"$#\" -gt 0 ]; do if [ \"$1\" = \"-o\" ]; then output=$2; shift; fi; shift; done\nmkdir -p \"$(dirname \"$output\")\"\nprintf binary > \"$output\"\n"
	require.NoError(t, os.WriteFile(filepath.Join(tools, "docker"), []byte(script), 0o755))
	t.Setenv("PATH", tools)
	t.Setenv("DOCKER_RECORD", record)
	previousHost := manifestHostOS
	manifestHostOS = "linux"
	t.Cleanup(func() { manifestHostOS = previousHost })
	handler := &manifestHandler{root: root}
	_, err := handler.compile(t.Context(), pipeline.CompileSpec{TargetOS: "darwin", TargetArch: "arm64", Output: ".wails/build/app", Toolchain: "docker", LocalRoots: []string{localRoot}, MinimumVersion: "12.0", Tags: []string{"production"}})
	require.NoError(t, err)
	data, err := os.ReadFile(record)
	require.NoError(t, err)
	arguments := strings.Split(strings.TrimSpace(string(data)), "\n")
	assert.Contains(t, arguments, "run")
	assert.Contains(t, arguments, "--rm")
	assert.Contains(t, arguments, root+":"+root)
	assert.Contains(t, arguments, localRoot+":"+localRoot+":ro")
	assert.Contains(t, arguments, "GOOS=darwin")
	assert.Contains(t, arguments, "GOARCH=arm64")
	assert.Contains(t, arguments, "CC=zcc-darwin-arm64")
	assert.Contains(t, arguments, "MACOSX_DEPLOYMENT_TARGET=12.0")
	assert.Contains(t, arguments, "--entrypoint")
	assert.Contains(t, arguments, "go")
	assert.Contains(t, arguments, "wails-cross")
	assert.Condition(t, func() bool {
		for _, argument := range arguments {
			if strings.Contains(argument, filepath.Join(root, ".wails", "build", ".docker-compile-output-")) {
				return true
			}
		}
		return false
	})
	assert.FileExists(t, filepath.Join(root, ".wails", "build", "app"))
}

func TestManifestDockerLinuxCrossCompileUsesNativeContainerCompiler(t *testing.T) {
	root := t.TempDir()
	tools := t.TempDir()
	record := filepath.Join(root, "docker-linux.txt")
	require.NoError(t, os.WriteFile(filepath.Join(tools, "docker"), []byte("#!/bin/sh\nprintf '%s\\n' \"$@\" > \"$DOCKER_RECORD\"\noutput=\nwhile [ \"$#\" -gt 0 ]; do if [ \"$1\" = \"-o\" ]; then output=$2; shift; fi; shift; done\nmkdir -p \"$(dirname \"$output\")\"\nprintf binary > \"$output\"\n"), 0o755))
	t.Setenv("PATH", tools)
	t.Setenv("DOCKER_RECORD", record)
	previousHost := manifestHostOS
	manifestHostOS = "linux"
	t.Cleanup(func() { manifestHostOS = previousHost })

	_, err := (&manifestHandler{root: root}).compile(t.Context(), pipeline.CompileSpec{TargetOS: "linux", TargetArch: "arm64", Output: ".wails/build/app", Toolchain: "docker"})
	require.NoError(t, err)
	arguments := strings.Split(strings.TrimSpace(readTestFile(t, record)), "\n")
	assert.Contains(t, arguments, "CC=gcc")
	assert.Contains(t, arguments, "CXX=g++")
}

func TestManifestPodmanCrossCompileUsesSELinuxSafeMounts(t *testing.T) {
	root := t.TempDir()
	localRoot := t.TempDir()
	tools := t.TempDir()
	record := filepath.Join(root, "podman-linux.txt")
	script := "#!/bin/sh\nprintf '%s\\n' \"$@\" > \"$PODMAN_RECORD\"\noutput=\nwhile [ \"$#\" -gt 0 ]; do if [ \"$1\" = \"-o\" ]; then output=$2; shift; fi; shift; done\nmkdir -p \"$(dirname \"$output\")\"\nprintf binary > \"$output\"\n"
	require.NoError(t, os.WriteFile(filepath.Join(tools, "podman"), []byte(script), 0o755))
	t.Setenv("PATH", tools)
	t.Setenv("PODMAN_RECORD", record)
	previousHost := manifestHostOS
	manifestHostOS = "linux"
	t.Cleanup(func() { manifestHostOS = previousHost })

	_, err := (&manifestHandler{root: root}).compile(t.Context(), pipeline.CompileSpec{
		TargetOS: "linux", TargetArch: "arm64", Output: ".wails/build/app",
		Toolchain: "docker", ContainerRuntime: "podman", ContainerImage: "ghcr.io/wailsapp/wails-cross:latest", LocalRoots: []string{localRoot},
	})
	require.NoError(t, err)
	arguments := strings.Split(strings.TrimSpace(readTestFile(t, record)), "\n")
	assert.Contains(t, arguments, root+":"+root+":Z")
	assert.Contains(t, arguments, localRoot+":"+localRoot+":ro,Z")
	assert.Contains(t, arguments, "--platform")
	assert.Contains(t, arguments, "linux/arm64")
	assert.Contains(t, arguments, "ghcr.io/wailsapp/wails-cross:latest")
}

func TestZigTargetRejectsTargetsOutsideTheClosedZigRegistry(t *testing.T) {
	for target, want := range map[string]string{
		"windows/amd64": "x86_64-windows-gnu",
		"windows/arm64": "aarch64-windows-gnu",
		"linux/amd64":   "x86_64-linux-gnu",
		"linux/arm64":   "aarch64-linux-gnu",
	} {
		platform, arch, _ := strings.Cut(target, "/")
		actual, err := zigTarget(platform, arch)
		require.NoError(t, err)
		assert.Equal(t, want, actual)
	}
	_, err := zigTarget("darwin", "arm64")
	assert.ErrorContains(t, err, "does not support")
}

func TestManifestAssetGenerationPreservesTheLastCompleteWorkspaceOnFailure(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, ".wails", "build", "default", "linux-amd64", "assets")
	require.NoError(t, os.MkdirAll(output, 0o755))
	sentinel := filepath.Join(output, "last-complete")
	require.NoError(t, os.WriteFile(sentinel, []byte("preserve"), 0o644))
	handler := &manifestHandler{root: root, config: manifest.Config{}}
	_, err := handler.assets(pipeline.AssetsSpec{
		TargetOS: "linux", TargetArch: "amd64", Directory: ".wails/build/default/linux-amd64/assets",
		Project: manifest.Project{Name: "app", BinaryName: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0", Icon: "missing.png"},
	})
	assert.ErrorContains(t, err, "project.icon")
	assert.Equal(t, "preserve", string(readTestFile(t, sentinel)))
	entries, err := filepath.Glob(filepath.Join(filepath.Dir(output), ".assets-*"))
	require.NoError(t, err)
	assert.Empty(t, entries, "failed staging directories must be disposable")
}

func TestPrintManifestPlanReturnsResolutionErrors(t *testing.T) {
	t.Chdir(t.TempDir())
	require.Error(t, printManifestPlan(manifestRunOptions{Verb: "build"}, false))
}

func TestExecutionPathGuardRejectsCrossHostEscapesAndWailsState(t *testing.T) {
	root := t.TempDir()
	for _, value := range []string{"../secret", `..\\secret`, "/tmp/secret", `C:\\temp\\secret`, `\\\\server\\share\\secret`, ".wails/cache/secret", `safe\\.WAILS\\secret`} {
		_, err := pathInsideProject(root, value)
		require.Error(t, err, value)
	}
	path, err := pathInsideProject(root, "assets/icon.png")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(root, "assets", "icon.png"), path)

	outside := t.TempDir()
	require.NoError(t, os.Symlink(outside, filepath.Join(root, "linked")))
	_, err = pathInsideProject(root, "linked/not-created-yet")
	require.ErrorContains(t, err, "resolves outside the project")
	_, err = (&manifestHandler{root: root}).install(t.Context(), pipeline.InstallSpec{Directory: "linked"})
	require.ErrorContains(t, err, "resolves outside the project")
}

func TestApplyGeneratedTargetSettings(t *testing.T) {
	root := t.TempDir()
	write := func(relative, content string) {
		path := filepath.Join(root, relative)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	}
	write("ios/xcode/main/Info.plist", "<key>CFBundleVersion</key><string>1.0.0</string><key>MinimumOSVersion</key><string>15.0</string>")
	write("android/app/build.gradle", "android {\n    defaultConfig {\n        applicationId \"com.wails.app\"\n        minSdk 21\n        versionCode 1\n        versionName \"1.0\"\n    }\n}\n")
	require.NoError(t, applyGeneratedTargetSettings(root, pipeline.AssetsSpec{TargetOS: "ios", MinimumVersion: "17.0", Project: manifest.Project{BuildNumber: 42}}))
	plist, err := os.ReadFile(filepath.Join(root, "ios/xcode/main/Info.plist"))
	require.NoError(t, err)
	assert.Contains(t, string(plist), "<string>42</string>")
	assert.Contains(t, string(plist), "<string>17.0</string>")
	gradle, err := os.ReadFile(filepath.Join(root, "android/app/build.gradle"))
	require.NoError(t, err)
	assert.Contains(t, string(gradle), "versionCode 42")
	require.NoError(t, applyGeneratedTargetSettings(root, pipeline.AssetsSpec{TargetOS: "android", MinimumVersion: "26", Project: manifest.Project{Identifier: "com.example.badge", Version: "2.4.1", BuildNumber: 42}}))
	gradle, err = os.ReadFile(filepath.Join(root, "android/app/build.gradle"))
	require.NoError(t, err)
	assert.Contains(t, string(gradle), `applicationId "com.example.badge"`)
	assert.Contains(t, string(gradle), "minSdk 26")
	assert.Contains(t, string(gradle), `versionName "2.4.1"`)
}

func TestApplyGeneratedTargetSettingsSupportsGradleAssignmentSyntax(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "android", "app", "build.gradle")
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, []byte("android {\n    defaultConfig {\n        applicationId = \"com.wails.app\"\n        minSdk = 21\n        versionCode = 1\n        versionName = \"1.0\"\n    }\n}\n"), 0o644))

	require.NoError(t, applyGeneratedTargetSettings(root, pipeline.AssetsSpec{
		TargetOS:       "android",
		MinimumVersion: "26",
		Project: manifest.Project{
			Identifier:  "com.example.badge",
			Version:     "2.4.1",
			BuildNumber: 42,
		},
	}))
	gradle := readTestFile(t, path)
	assert.Contains(t, gradle, `applicationId = "com.example.badge"`)
	assert.Contains(t, gradle, "minSdk = 26")
	assert.Contains(t, gradle, "versionCode = 42")
	assert.Contains(t, gradle, `versionName = "2.4.1"`)
}

func TestReplacePlistStringFillsSelfClosingValue(t *testing.T) {
	path := filepath.Join(t.TempDir(), "Info.plist")
	require.NoError(t, os.WriteFile(path, []byte("<key>CFBundleVersion</key>\n<string/>\n"), 0o644))
	require.NoError(t, replacePlistString(path, "CFBundleVersion", "13"))
	assert.Equal(t, "<key>CFBundleVersion</key>\n<string>13</string>\n", readTestFile(t, path))
}

func TestCopyManifestPathPreservesSymlink(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	destination := filepath.Join(root, "destination")
	require.NoError(t, os.MkdirAll(source, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(source, "actual"), []byte("payload"), 0o644))
	require.NoError(t, os.Symlink("actual", filepath.Join(source, "link")))
	require.NoError(t, copyManifestPath(source, destination))
	info, err := os.Lstat(filepath.Join(destination, "link"))
	require.NoError(t, err)
	assert.NotZero(t, info.Mode()&os.ModeSymlink)
	target, err := os.Readlink(filepath.Join(destination, "link"))
	require.NoError(t, err)
	assert.Equal(t, "actual", target)
}

func TestFindAndMovePackageIgnoresExistingDestination(t *testing.T) {
	root := t.TempDir()
	staging := filepath.Join(root, "staging")
	require.NoError(t, os.MkdirAll(staging, 0o755))
	expected := filepath.Join(root, "app.deb")
	require.NoError(t, os.WriteFile(expected, []byte("old"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(staging, "app_1.0.0_amd64.deb"), []byte("new"), 0o644))

	require.NoError(t, findAndMovePackage(staging, "app", "deb", expected))
	data, err := os.ReadFile(expected)
	require.NoError(t, err)
	assert.Equal(t, "new", string(data))
}

func TestManifestHandlerIdentityIncludesRelevantEnvironment(t *testing.T) {
	handler := &manifestHandler{root: t.TempDir()}
	node := pipeline.Node{Key: "assets", Kind: pipeline.GeneratePlatformAssets, Spec: pipeline.AssetsSpec{}}
	t.Setenv("PATH", "/first")
	first, err := handler.Identity(t.Context(), node)
	require.NoError(t, err)
	t.Setenv("PATH", "/second")
	second, err := handler.Identity(t.Context(), node)
	require.NoError(t, err)
	assert.NotEqual(t, first, second)
}

func TestManifestDockerIdentityIncludesTheExactCrossImage(t *testing.T) {
	tools := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(tools, "docker"), []byte("#!/bin/sh\nexit 0\n"), 0o755))
	t.Setenv("PATH", tools)
	previous := manifestDockerImageIdentity
	image := "sha256:first"
	manifestDockerImageIdentity = func(context.Context, string, string) (string, error) { return image, nil }
	t.Cleanup(func() { manifestDockerImageIdentity = previous })
	handler := &manifestHandler{root: t.TempDir()}
	node := pipeline.Node{Key: "compile", Kind: pipeline.CompileApplication, Spec: pipeline.CompileSpec{Toolchain: "docker"}}
	first, err := handler.Identity(t.Context(), node)
	require.NoError(t, err)
	assert.Contains(t, first, "sha256:first")
	image = "sha256:second"
	second, err := handler.Identity(t.Context(), node)
	require.NoError(t, err)
	assert.NotEqual(t, first, second)
}

func TestDevelopmentFrontendIdentityIncludesSessionURLAndPort(t *testing.T) {
	node := pipeline.Node{Kind: pipeline.BuildFrontend, Spec: pipeline.FrontendSpec{Production: false}}
	t.Setenv("FRONTEND_DEVSERVER_URL", "http://localhost:9245")
	t.Setenv(wailsVitePort, "9245")
	first := relevantEnvironment(node, nil)
	t.Setenv("FRONTEND_DEVSERVER_URL", "http://localhost:9246")
	assert.NotEqual(t, first, relevantEnvironment(node, nil))
	t.Setenv("FRONTEND_DEVSERVER_URL", "http://localhost:9245")
	t.Setenv(wailsVitePort, "9246")
	assert.NotEqual(t, first, relevantEnvironment(node, nil))

	production := pipeline.Node{Kind: pipeline.BuildFrontend, Spec: pipeline.FrontendSpec{Production: true}}
	first = relevantEnvironment(production, nil)
	t.Setenv("FRONTEND_DEVSERVER_URL", "http://localhost:9247")
	t.Setenv(wailsVitePort, "9247")
	second := relevantEnvironment(production, nil)
	assert.Equal(t, first, second)
}

func TestDevelopmentEnvironmentOverridesAreIsolatedFromProcessState(t *testing.T) {
	node := pipeline.Node{Kind: pipeline.BuildFrontend, Spec: pipeline.FrontendSpec{Production: false}}
	t.Setenv("FRONTEND_DEVSERVER_URL", "http://global:9245")
	t.Setenv(wailsVitePort, "9245")
	first := relevantEnvironment(node, []string{"FRONTEND_DEVSERVER_URL=http://generation-one:9246", wailsVitePort + "=9246"})
	second := relevantEnvironment(node, []string{"FRONTEND_DEVSERVER_URL=http://generation-two:9247", wailsVitePort + "=9247"})
	assert.NotEqual(t, first, second)
	assert.Equal(t, "http://global:9245", os.Getenv("FRONTEND_DEVSERVER_URL"))
	assert.Equal(t, "9245", os.Getenv(wailsVitePort))
}
