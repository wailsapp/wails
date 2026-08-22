package pipeline

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/report"
	"github.com/wailsapp/wails/v3/internal/wake/cache"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

type fakeHandler struct {
	mu   sync.Mutex
	root string
	runs []NodeKey
}

type branchHandler struct {
	mu   sync.Mutex
	runs []NodeKey
}

type cancelHandler struct{ started chan struct{} }

type parallelCancelHandler struct {
	want    int
	started chan struct{}
	once    sync.Once
	mu      sync.Mutex
	running int
}

func (*cancelHandler) Identity(context.Context, Node) (string, error) { return "cancel-v1", nil }
func (h *cancelHandler) Run(ctx context.Context, _ Node) (RunResult, error) {
	close(h.started)
	<-ctx.Done()
	return RunResult{}, ctx.Err()
}

func (*parallelCancelHandler) Identity(context.Context, Node) (string, error) {
	return "parallel-cancel-v1", nil
}

func (h *parallelCancelHandler) Run(ctx context.Context, _ Node) (RunResult, error) {
	h.mu.Lock()
	h.running++
	if h.running == h.want {
		h.once.Do(func() { close(h.started) })
	}
	h.mu.Unlock()
	<-ctx.Done()
	h.mu.Lock()
	h.running--
	h.mu.Unlock()
	return RunResult{}, ctx.Err()
}

type failureReporter struct {
	report.Nop
	failure report.Failure
}

func (r *failureReporter) StepFailed(_ report.StepID, failure report.Failure) {
	r.failure = failure
}

func (*branchHandler) Identity(context.Context, Node) (string, error) { return "branch-v1", nil }
func (h *branchHandler) Run(_ context.Context, node Node) (RunResult, error) {
	h.mu.Lock()
	h.runs = append(h.runs, node.Key)
	h.mu.Unlock()
	if node.Key == "fail" {
		return RunResult{Detail: "specific compiler output"}, errors.New("expected failure")
	}
	return RunResult{}, nil
}

func (f *fakeHandler) Identity(context.Context, Node) (string, error) { return "fake-v1", nil }
func (f *fakeHandler) Run(_ context.Context, n Node) (RunResult, error) {
	f.mu.Lock()
	f.runs = append(f.runs, n.Key)
	f.mu.Unlock()
	if n.Output != "" {
		path := filepath.Join(f.root, filepath.FromSlash(n.Output))
		if n.Kind == CompileApplication {
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return RunResult{}, err
			}
			if err := os.WriteFile(path, []byte(n.Key), 0o644); err != nil {
				return RunResult{}, err
			}
		} else {
			if err := os.MkdirAll(path, 0o755); err != nil {
				return RunResult{}, err
			}
			if err := os.WriteFile(filepath.Join(path, "result"), []byte(n.Key), 0o644); err != nil {
				return RunResult{}, err
			}
		}
	}
	return RunResult{}, nil
}

func TestPlanSharesFrontendAndCachesSecondRun(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "frontend"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\nfunc main(){}\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "frontend", "package.json"), []byte("{}"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "frontend", "package-lock.json"), []byte("{}"), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "frontend", "node_modules"), 0o755))
	require.NoError(t, manifest.WriteMinimal(root, manifest.Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"}))
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	plan, err := PlanBuild(loaded.Config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	require.NoError(t, err)
	assert.Len(t, plan.Nodes, 6)
	h := &fakeHandler{root: root}
	executor := Executor{Handler: h}
	results, err := executor.Execute(context.Background(), plan, ExecuteOptions{Root: root, Reporter: report.Nop{}})
	require.NoError(t, err)
	assert.Empty(t, results["frontend:install"].Artifact, "a receipt is not an artifact")
	assert.NotEmpty(t, results["publish:target:linux/amd64:compile"].Artifact, "publication must preserve the source artifact identity")
	first := len(h.runs)
	assert.Equal(t, 6, first)
	_, err = executor.Execute(context.Background(), plan, ExecuteOptions{Root: root, Reporter: report.Nop{}})
	require.NoError(t, err)
	assert.Equal(t, first, len(h.runs), "a warm build should run no handlers")
	publishKey := NodeKey("publish:target:linux/amd64:compile")
	assert.Equal(t, CacheArtifact, plan.Nodes[publishKey].Cache)
	publishedResult := filepath.Join(root, "bin", "app", "result")
	require.NoError(t, os.WriteFile(publishedResult, []byte("user-modified"), 0o644))
	_, err = executor.Execute(context.Background(), plan, ExecuteOptions{Root: root, Reporter: report.Nop{}})
	require.NoError(t, err)
	assert.Equal(t, 2, countNodeRuns(h.runs, publishKey), "modified published output must be regenerated")
	published, err := os.ReadFile(publishedResult)
	require.NoError(t, err)
	assert.Equal(t, []byte(publishKey), published)
}

// Regression test for https://github.com/wailsapp/wails/issues/1031.
// A backend implementation edit may require bindings to be reconsidered, but
// unchanged generated bindings must stop invalidation before frontend:build.
func TestDevelopmentBackendChangeDoesNotRebuildFrontend(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "frontend", "node_modules"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\nfunc generation() int { return 1 }\nfunc main() {}\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "frontend", "package.json"), []byte("{}"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "frontend", "package-lock.json"), []byte("{}"), 0o644))
	require.NoError(t, manifest.WriteMinimal(root, manifest.Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"}))
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	plan, err := PlanBuild(loaded.Config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64", Development: true})
	require.NoError(t, err)

	handler := &fakeHandler{root: root}
	executor := Executor{Handler: handler}
	_, err = executor.Execute(t.Context(), plan, ExecuteOptions{Root: root, Reporter: report.Nop{}})
	require.NoError(t, err)
	assert.Equal(t, 1, countNodeRuns(handler.runs, NodeKey("frontend:build")))
	assert.Equal(t, 1, countNodeRuns(handler.runs, NodeKey("target:linux/amd64:compile")))

	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\nfunc generation() int { return 2 }\nfunc main() {}\n"), 0o644))
	_, err = executor.Execute(t.Context(), plan, ExecuteOptions{Root: root, Reporter: report.Nop{}})
	require.NoError(t, err)

	assert.Equal(t, 2, countNodeRuns(handler.runs, NodeKey("frontend:bindings")), "backend changes must reconsider generated bindings")
	assert.Equal(t, 1, countNodeRuns(handler.runs, NodeKey("frontend:build")), "unchanged binding output must not rebuild the frontend")
	assert.Equal(t, 2, countNodeRuns(handler.runs, NodeKey("target:linux/amd64:compile")), "backend changes must rebuild the application")
}

func TestPlanInspectionReportsActualCacheDecisionsWithoutChangingFiles(t *testing.T) {
	config := performanceConfig(t)
	plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	require.NoError(t, err)
	handler := &fakeHandler{root: config.Root}
	executor := Executor{Handler: handler}

	cold, err := executor.Inspect(t.Context(), plan, config.Root)
	require.NoError(t, err)
	assert.NoDirExists(t, filepath.Join(config.Root, ".wails"))
	for key, operation := range cold.Operations {
		assert.Equal(t, "run", operation.Decision, key)
	}
	assert.NotEmpty(t, cold.Operations["target:linux/amd64:compile"].Inputs)

	_, err = executor.Execute(t.Context(), plan, ExecuteOptions{Root: config.Root, Reporter: report.Nop{}})
	require.NoError(t, err)
	warm, err := executor.Inspect(t.Context(), plan, config.Root)
	require.NoError(t, err)
	for key, operation := range warm.Operations {
		assert.Equal(t, "cached", operation.Decision, key)
	}

	output := filepath.Join(config.Root, filepath.FromSlash(plan.Nodes["target:linux/amd64:compile"].Output))
	require.NoError(t, os.Remove(output))
	restore, err := executor.Inspect(t.Context(), plan, config.Root)
	require.NoError(t, err)
	assert.Equal(t, "restore", restore.Operations["target:linux/amd64:compile"].Decision)
	assert.NoFileExists(t, output, "inspection must not restore a cached output")
}

func TestWindowsBuildGeneratesAssetsBeforeCompile(t *testing.T) {
	config := testConfig(t)
	plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: "windows", TargetArch: "amd64"})
	require.NoError(t, err)
	assets := NodeKey("target:windows/amd64:assets")
	compile := plan.Nodes["target:windows/amd64:compile"]
	require.Contains(t, plan.Nodes, assets)
	assert.Contains(t, compile.Dependencies, assets)
}

func TestDevelopmentPlanDoesNotInheritFiniteBuildPolicy(t *testing.T) {
	config := testConfig(t)
	config.Build.OutputDirectory = "release-output"
	config.Build.Go.Tags = []string{"release"}
	config.Build.Go.LinkerFlags = []string{"-X", "release=true"}
	config.Build.Go.CompilerFlags = []string{"all=-l"}
	config.Build.Go.GarbleArgs = []string{"-tiny"}
	config.Build.Obfuscation = true
	config.Build.TrimPath = true
	config.Build.Strip = true
	config.Dev.Tags = []string{"debug", "devtools"}
	config.Targets.Linux.AMD64.Tags = []string{"target-release"}

	request := Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64", Development: true, ExtraTags: []string{"cli-debug"}}
	plan, err := PlanBuild(config, request)
	require.NoError(t, err)
	compile := plan.Nodes["target:linux/amd64:compile"]
	assert.Equal(t, ".wails/dev/linux-amd64/app", compile.Output)
	spec := compile.Spec.(CompileSpec)
	assert.Equal(t, compile.Output, spec.Output)
	assert.Equal(t, []string{"debug", "devtools", "cli-debug"}, spec.Tags)
	assert.Empty(t, spec.LinkerFlags)
	assert.Empty(t, spec.CompilerFlags)
	assert.Empty(t, spec.GarbleArgs)
	assert.False(t, spec.Production)
	assert.False(t, spec.Obfuscated)
	assert.False(t, spec.TrimPath)
	assert.False(t, spec.Strip)
	bindings := plan.Nodes["frontend:bindings"].Spec.(BindingsSpec)
	assert.Equal(t, spec.Tags, bindings.Tags)
	assert.False(t, bindings.Obfuscated)
	frontend := plan.Nodes["frontend:build"]
	require.Len(t, frontend.Inputs, 1)
	assert.Equal(t, "frontend-dev-bootstrap", frontend.Inputs[0].Label)
	assert.False(t, frontend.Inputs[0].IncludeAll, "the persistent Dev server owns frontend source changes")
	assert.Contains(t, frontend.Inputs[0].Files, filepath.Join(config.Frontend.Directory, "package.json"))
	assert.True(t, plan.Nodes["frontend:bindings"].Inputs[0].UseGitIgnore)
	assert.True(t, compile.Inputs[0].UseGitIgnore)

	changedBuildPolicy := config
	changedBuildPolicy.Build.OutputDirectory = "somewhere-else"
	changedBuildPolicy.Build.Go.Tags = []string{"different-release"}
	changedBuildPolicy.Build.Go.LinkerFlags = []string{"different"}
	changedBuildPolicy.Build.Go.CompilerFlags = []string{"different"}
	changedBuildPolicy.Build.Go.GarbleArgs = []string{"different"}
	changedBuildPolicy.Build.Obfuscation = false
	changedBuildPolicy.Build.TrimPath = false
	changedBuildPolicy.Build.Strip = false
	changedPlan, err := PlanBuild(changedBuildPolicy, request)
	require.NoError(t, err)
	assert.Equal(t, plan, changedPlan, "finite build policy must not change the Dev Plan or its future action keys")
}

func TestProductionPlanCarriesDeclaredEnvironmentAndCompilerPolicy(t *testing.T) {
	config := testConfig(t)
	config.Frontend.Environment = map[string]string{"PUBLIC_RELEASE": "true"}
	config.Build.Environment = map[string]string{"CGO_ENABLED": "0", "SHARED": "build"}
	config.Build.VCSInfo = true
	config.Targets.Linux.AMD64.Environment = map[string]string{"CC": "zig cc", "SHARED": "target"}
	config.Targets.Linux.AMD64.Toolchain = "zig"
	config.Targets.Linux.AMD64.LinkerFlags = []string{"-X target=value"}
	config.Targets.Linux.AMD64.CompilerFlags = []string{"all=-l"}
	config.Targets.Linux.AMD64.GarbleArgs = []string{"-literals"}
	config.Targets.Linux.AMD64.Obfuscated = true
	config.Targets.Linux.AMD64.ObfuscatedSet = true

	plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	require.NoError(t, err)
	assert.Equal(t, config.Frontend.Environment, plan.Nodes["frontend:install"].Spec.(InstallSpec).Environment)
	assert.Equal(t, config.Frontend.Environment, plan.Nodes["frontend:build"].Spec.(FrontendSpec).Environment)
	compile := plan.Nodes["target:linux/amd64:compile"].Spec.(CompileSpec)
	assert.Equal(t, map[string]string{"CC": "zig cc", "SHARED": "target"}, compile.Environment)
	assert.Equal(t, "zig", compile.Toolchain)
	assert.Equal(t, []string{"-X target=value"}, compile.LinkerFlags)
	assert.Equal(t, []string{"all=-l"}, compile.CompilerFlags)
	assert.Equal(t, []string{"-literals"}, compile.GarbleArgs)
	assert.True(t, compile.Obfuscated)
	assert.True(t, compile.VCSInfo)
}

func TestExplicitTargetFalseClearsInheritedObfuscation(t *testing.T) {
	config := testConfig(t)
	config.Build.Obfuscation = true
	config.Targets.Linux.AMD64.ObfuscatedSet = true
	config.Targets.Linux.AMD64.Obfuscated = false

	plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	require.NoError(t, err)
	compile := plan.Nodes["target:linux/amd64:compile"].Spec.(CompileSpec)
	assert.False(t, compile.Obfuscated)
	assert.False(t, plan.Nodes["frontend:bindings"].Spec.(BindingsSpec).Obfuscated)
}

func TestAnonymousGarbleArgumentsAreAppliedOnceAndAffectActionIdentity(t *testing.T) {
	config := testConfig(t)
	config.Build.Obfuscation = true
	config.Build.Go.GarbleArgs = []string{"-tiny"}

	base, err := PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	require.NoError(t, err)
	overridden, err := PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64", GarbleArgs: []string{"-literals"}})
	require.NoError(t, err)

	baseCompile := base.Nodes["target:linux/amd64:compile"]
	overriddenCompile := overridden.Nodes["target:linux/amd64:compile"]
	assert.Equal(t, []string{"-tiny", "-literals"}, overriddenCompile.Spec.(CompileSpec).GarbleArgs)
	baseKey, err := cache.ActionKey(string(baseCompile.Kind), baseCompile.Spec, nil, nil)
	require.NoError(t, err)
	overriddenKey, err := cache.ActionKey(string(overriddenCompile.Kind), overriddenCompile.Spec, nil, nil)
	require.NoError(t, err)
	assert.NotEqual(t, baseKey, overriddenKey)
}

func TestBindingAndCompileSnapshotsShareTreeSelectors(t *testing.T) {
	config := performanceConfig(t)
	local := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(local, "go.mod"), []byte("module example.com/local\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(config.Root, "go.mod"), []byte("module example.com/app\n\nreplace example.com/local => "+filepath.ToSlash(local)+"\n"), 0o644))
	plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	require.NoError(t, err)
	binding := plan.Nodes["frontend:bindings"]
	compile := plan.Nodes["target:linux/amd64:compile"]
	compileByRoot := map[string]InputSpec{}
	for _, input := range compile.Inputs {
		compileByRoot[input.Root] = input
	}
	for _, input := range binding.Inputs {
		if !input.SemanticGo {
			continue
		}
		matching, ok := compileByRoot[input.Root]
		require.True(t, ok, input.Root)
		assert.Equal(t, matching.IncludeNames, input.IncludeNames)
		assert.Equal(t, matching.IncludeExtensions, input.IncludeExtensions)
		assert.Equal(t, matching.ExcludeDirs, input.ExcludeDirs)
		assert.Equal(t, matching.ExcludeSuffixes, input.ExcludeSuffixes)
		assert.Equal(t, matching.UseGitIgnore, input.UseGitIgnore)
	}
}

func TestPlannerBuildsOnePlanForMultipleTargets(t *testing.T) {
	config := testConfig(t)
	plan, err := PlanBuild(config, Request{Verb: "build", Targets: []Target{{OS: "linux", Arch: "amd64"}, {OS: "linux", Arch: "arm64"}}})
	require.NoError(t, err)
	assert.Equal(t, "linux/amd64,linux/arm64", plan.Target)
	assert.Contains(t, plan.Nodes, NodeKey("frontend:build"))
	assert.Contains(t, plan.Nodes, NodeKey("target:linux/amd64:compile"))
	assert.Contains(t, plan.Nodes, NodeKey("target:linux/arm64:compile"))
	require.Equal(t, []NodeKey{"collect:artifacts"}, plan.Roots)
	assert.ElementsMatch(t, []NodeKey{"publish:target:linux/amd64:compile", "publish:target:linux/arm64:compile"}, plan.Artifacts)
	collect := plan.Nodes["collect:artifacts"]
	assert.Equal(t, CollectArtifacts, collect.Kind)
	assert.Equal(t, CacheArtifact, collect.Cache)
	assert.ElementsMatch(t, plan.Artifacts, collect.Dependencies)
	amd64 := plan.Nodes["target:linux/amd64:compile"].Output
	arm64 := plan.Nodes["target:linux/arm64:compile"].Output
	assert.NotEqual(t, amd64, arm64)
}

func TestPlannerRejectsInvalidTargetAndFormatCombinations(t *testing.T) {
	config := testConfig(t)
	_, err := PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "riscv64"})
	require.ErrorContains(t, err, `unsupported target "linux/riscv64"`)
	_, err = PlanBuild(config, Request{Verb: "package", TargetOS: "linux", TargetArch: "amd64", Formats: []string{"nsis"}})
	require.ErrorContains(t, err, `format "nsis" is not supported for any selected target`)
	_, err = PlanBuild(config, Request{Verb: "build", TargetOS: "darwin", TargetArch: "universal"})
	require.NoError(t, err)
}

func TestPlannerRejectsDuplicateCommaSeparatedTargetsAndFormats(t *testing.T) {
	config := testConfig(t)
	_, err := PlanBuild(config, Request{Verb: "build", Targets: []Target{{OS: "linux", Arch: "amd64"}, {OS: "linux", Arch: "amd64"}}})
	require.ErrorContains(t, err, "duplicate target linux/amd64")
	_, err = PlanBuild(config, Request{Verb: "package", TargetOS: "linux", TargetArch: "amd64", Formats: []string{"deb", "deb"}})
	require.ErrorContains(t, err, "duplicate package format deb")
}

func TestPlannerCarriesPackageCustomizationForEveryFormat(t *testing.T) {
	tests := []struct {
		platform string
		format   string
		config   manifest.PackageFormat
	}{
		{"windows", "nsis", manifest.PackageFormat{Template: "packaging/windows/installer.nsi"}},
		{"windows", "msix", manifest.PackageFormat{Publisher: "CN=Example"}},
		{"darwin", "dmg", manifest.PackageFormat{Background: "packaging/darwin/background.png"}},
		{"linux", "appimage", manifest.PackageFormat{Categories: []string{"Development", "IDE"}}},
		{"linux", "deb", manifest.PackageFormat{Maintainer: "Example"}},
		{"linux", "rpm", manifest.PackageFormat{Maintainer: "Example"}},
		{"linux", "archlinux", manifest.PackageFormat{Maintainer: "Example"}},
	}
	for _, test := range tests {
		t.Run(test.platform+"/"+test.format, func(t *testing.T) {
			config := testConfig(t)
			setTestPackageFormat(&config.Package, test.platform, test.format, test.config)
			config.Associations = []manifest.Association{{Extensions: []string{"example"}}}
			config.Protocols = []manifest.Protocol{{Scheme: "example"}}
			plan, err := PlanBuild(config, Request{Verb: "package", TargetOS: test.platform, TargetArch: "arm64", Formats: []string{test.format}})
			require.NoError(t, err)
			node := plan.Nodes[NodeKey("package:"+test.platform+"/arm64:"+test.format)]
			spec := node.Spec.(PackageSpec)
			assert.Equal(t, test.config, spec.Config)
			assert.Equal(t, config.Associations, spec.Associations)
			assert.Equal(t, config.Protocols, spec.Protocols)
			if test.config.Template != "" {
				require.Len(t, node.Inputs, 1)
				assert.Equal(t, "package-template", node.Inputs[0].Label)
			}
		})
	}
}

func TestPlannerSnapshotsBuiltInDMGResources(t *testing.T) {
	config := testConfig(t)
	config.Package.Darwin.DMG.Background = "packaging/background.png"
	config.Package.Darwin.DMG.Files = map[string]string{"Read Me": "README.md", "License": "LICENSE"}
	plan, err := PlanBuild(config, Request{Verb: "package", TargetOS: "darwin", TargetArch: "arm64", Formats: []string{"dmg"}})
	require.NoError(t, err)
	node := plan.Nodes[NodeKey("package:darwin/arm64:dmg")]
	require.Len(t, node.Inputs, 1)
	assert.Equal(t, "package-resources", node.Inputs[0].Label)
	assert.ElementsMatch(t, []string{"packaging/background.png", "README.md", "LICENSE"}, node.Inputs[0].Files)
}

func setTestPackageFormat(packages *manifest.Packages, platform, format string, value manifest.PackageFormat) {
	switch platform + "/" + format {
	case "windows/nsis":
		packages.Windows.NSIS = manifest.NSISPackage{Template: value.Template, InstallScope: value.InstallScope}
	case "windows/msix":
		packages.Windows.MSIX = manifest.MSIXPackage{Publisher: value.Publisher, Manifest: value.Manifest}
	case "darwin/dmg":
		packages.Darwin.DMG = manifest.DMGPackage{Template: value.Template, Background: value.Background, VolumeIcon: value.VolumeIcon, FileIcon: value.FileIcon, Files: value.Files, WindowWidth: value.WindowWidth, WindowHeight: value.WindowHeight}
	case "linux/appimage":
		packages.Linux.AppImage = manifest.AppImagePackage{Icon: value.Icon, DesktopEntry: value.DesktopEntry, Categories: value.Categories}
	case "linux/deb":
		packages.Linux.Deb = manifest.LinuxPackage{Template: value.Template, Maintainer: value.Maintainer}
	case "linux/rpm":
		packages.Linux.RPM = manifest.LinuxPackage{Template: value.Template, Maintainer: value.Maintainer}
	case "linux/archlinux":
		packages.Linux.ArchLinux = manifest.LinuxPackage{Template: value.Template, Maintainer: value.Maintainer}
	}
}

func TestPlannerDoesNotMutateRequestedOrConfiguredFormats(t *testing.T) {
	config := testConfig(t)
	config.Package.Linux.Formats = []string{"rpm", "deb"}
	requested := []string{"rpm", "deb"}
	_, err := PlanBuild(config, Request{Verb: "package", TargetOS: "linux", TargetArch: "amd64", Formats: requested})
	require.NoError(t, err)
	assert.Equal(t, []string{"rpm", "deb"}, requested)
	assert.Equal(t, []string{"rpm", "deb"}, config.Package.Linux.Formats)
}

func TestPlannerFindsParentGoModuleMetadata(t *testing.T) {
	module := t.TempDir()
	project := filepath.Join(module, "examples", "app")
	require.NoError(t, os.MkdirAll(project, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(module, "go.mod"), []byte("module example\n"), 0o644))
	files := goMetadataFiles(project)
	assert.Contains(t, files, filepath.Join(module, "go.mod"))
	assert.Contains(t, files, filepath.Join(module, "go.sum"))
}

func TestTargetCommonAndArchitectureOverridesReachSpecs(t *testing.T) {
	config := testConfig(t)
	config.Targets.Darwin.ProductName = "Mac App"
	config.Targets.Darwin.Identifier = "com.example.mac"
	config.Targets.Darwin.MinimumVersion = "12.0"
	config.Targets.Darwin.BuildNumber = 4
	config.Targets.Darwin.ARM64.MinimumVersion = "13.0"
	config.Targets.Darwin.ARM64.BuildNumber = 8
	config.Targets.Darwin.ARM64.Tags = []string{"metal"}
	plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: "darwin", TargetArch: "arm64"})
	require.NoError(t, err)
	compile := plan.Nodes["target:darwin/arm64:compile"].Spec.(CompileSpec)
	assert.Equal(t, "13.0", compile.MinimumVersion)
	assert.Contains(t, compile.Tags, "metal")
	pkg := plan.Nodes["assemble:darwin/arm64"].Spec.(PackageSpec)
	assert.Equal(t, "Mac App", pkg.Project.ProductName)
	assert.Equal(t, "com.example.mac", pkg.Project.Identifier)
	assert.Equal(t, 8, pkg.Project.BuildNumber)
}

func TestMobilePackagePlansUseResolvedDestinations(t *testing.T) {
	config := testConfig(t)
	for _, tc := range []struct {
		os, format, variant, key string
	}{{"ios", "", "simulator", "assemble:ios/arm64"}, {"android", "aab", "", "package:android/arm64:aab"}} {
		t.Run(tc.os, func(t *testing.T) {
			plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: tc.os, TargetArch: "arm64", Formats: nonemptyStrings(tc.format)})
			require.NoError(t, err)
			node, ok := plan.Nodes[NodeKey(tc.key)]
			require.True(t, ok)
			spec := node.Spec.(PackageSpec)
			assert.Equal(t, tc.variant, spec.Destination)
			assert.Contains(t, node.Dependencies, NodeKey("target:"+tc.os+"/arm64:assets"))
			if tc.os == "ios" {
				assert.Equal(t, CacheNever, node.Cache, "iOS assembly signs and must not be reusable")
			}
		})
	}
}

func TestPlannerAcceptanceMatrixCoversEverySupportedTargetAndFormat(t *testing.T) {
	tests := []struct {
		platform string
		arch     string
		formats  []string
	}{
		{"windows", "amd64", []string{"nsis", "msix"}},
		{"windows", "arm64", []string{"nsis", "msix"}},
		{"darwin", "amd64", []string{"dmg"}},
		{"darwin", "arm64", []string{"dmg"}},
		{"darwin", "universal", []string{"dmg"}},
		{"linux", "amd64", []string{"appimage", "deb", "rpm", "archlinux"}},
		{"linux", "arm64", []string{"appimage", "deb", "rpm", "archlinux"}},
		{"ios", "arm64", []string{"ipa"}},
		{"android", "arm64", []string{"aab"}},
		{"android", "amd64", []string{"aab"}},
		{"android", "universal", []string{"aab"}},
	}
	for _, test := range tests {
		for _, format := range test.formats {
			t.Run(test.platform+"/"+test.arch+"/"+format, func(t *testing.T) {
				config := testConfig(t)
				request := Request{Verb: "package", TargetOS: test.platform, TargetArch: test.arch, Formats: []string{format}}
				if test.platform == "ios" && format == "ipa" {
					config.Selected = manifest.Profile{Name: "release", Targets: []manifest.ProfileTarget{{Target: "ios/arm64", Destination: "device", Formats: []string{"ipa"}}}}
					request = Request{Verb: "build"}
				}
				plan, err := PlanBuild(config, request)
				require.NoError(t, err)
				key := NodeKey("package:" + test.platform + "/" + test.arch + ":" + format)
				assert.Contains(t, plan.Nodes, key)
				publishKey := NodeKey("publish:" + test.platform + "/" + test.arch + ":" + format)
				assert.Contains(t, plan.Artifacts, publishKey)
				assert.Equal(t, []NodeKey{key}, plan.Nodes[publishKey].Dependencies)
			})
		}
	}
}

func nonemptyStrings(value string) []string {
	if value == "" {
		return nil
	}
	return []string{value}
}

func TestPlannerFiltersRegistrationsForTargetAndCarriesCapabilities(t *testing.T) {
	config := testConfig(t)
	config.Targets.Windows.Capabilities = []string{"internetClient"}
	config.Associations = []manifest.Association{
		{Extensions: []string{"shared"}},
		{Extensions: []string{"windows"}, Platforms: []string{"windows"}},
		{Extensions: []string{"darwin"}, Platforms: []string{"darwin"}},
	}
	config.Protocols = []manifest.Protocol{
		{Scheme: "shared"},
		{Scheme: "windows", Platforms: []string{"windows"}},
		{Scheme: "darwin", Platforms: []string{"darwin"}},
	}

	plan, err := PlanBuild(config, Request{Verb: "package", TargetOS: "windows", TargetArch: "amd64", Formats: []string{"msix"}})
	require.NoError(t, err)
	assets := plan.Nodes["target:windows/amd64:assets"].Spec.(AssetsSpec)
	packaging := plan.Nodes["package:windows/amd64:msix"].Spec.(PackageSpec)
	assert.Equal(t, []string{"internetClient"}, assets.Capabilities)
	assert.Equal(t, []string{"internetClient"}, packaging.Capabilities)
	assert.Equal(t, []manifest.Association{config.Associations[0], config.Associations[1]}, assets.Associations)
	assert.Equal(t, assets.Associations, packaging.Associations)
	assert.Equal(t, []manifest.Protocol{config.Protocols[0], config.Protocols[1]}, assets.Protocols)
	assert.Equal(t, assets.Protocols, packaging.Protocols)
}

func TestPlannerRejectsIPAPackagingForSimulator(t *testing.T) {
	config := testConfig(t)
	_, err := PlanBuild(config, Request{Verb: "package", TargetOS: "ios", TargetArch: "arm64", Formats: []string{"ipa"}})
	require.ErrorContains(t, err, `profile destination = "device"`)
}

func TestPlannerSnapshotsLocalModuleReplacements(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "app")
	dependency := filepath.Join(parent, "dependency")
	require.NoError(t, os.MkdirAll(root, 0o755))
	require.NoError(t, os.MkdirAll(dependency, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example/app\n\ngo 1.24\n\nreplace example/dependency => ../dependency\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dependency, "go.mod"), []byte("module example/dependency\n"), 0o644))

	inputs, err := goLocalSourceInputs(root)
	require.NoError(t, err)
	require.Len(t, inputs, 1)
	assert.Equal(t, dependency, inputs[0].Root)
}

func TestPlannerSnapshotsLocalWorkspaceReplacements(t *testing.T) {
	workspace := t.TempDir()
	root := filepath.Join(workspace, "app")
	dependency := filepath.Join(workspace, "dependency")
	require.NoError(t, os.MkdirAll(root, 0o755))
	require.NoError(t, os.MkdirAll(dependency, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example/app\n\ngo 1.24\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(workspace, "go.work"), []byte("go 1.24\n\nuse ./app\nreplace example/dependency => ./dependency\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dependency, "go.mod"), []byte("module example/dependency\n"), 0o644))

	inputs, err := goLocalSourceInputs(root)
	require.NoError(t, err)
	require.Len(t, inputs, 1)
	assert.Equal(t, dependency, inputs[0].Root)
}

func TestPlannerPrivateHelpersCoverDefensiveAndFormatSpecificPaths(t *testing.T) {
	plan := Plan{Nodes: map[NodeKey]Node{}}
	node := Node{Key: "one", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheNever}
	assert.Equal(t, NodeKey("one"), addNode(&plan, node, []OriginReference{{Field: "build"}}))
	assert.PanicsWithValue(t, "duplicate planner key: one", func() { addNode(&plan, node, nil) })

	root := t.TempDir()
	templateDirectory := filepath.Join(root, "template")
	require.NoError(t, os.MkdirAll(templateDirectory, 0o755))
	tests := []struct {
		platform string
		format   string
		config   manifest.PackageFormat
		want     []InputSpec
	}{
		{"windows", "msix", manifest.PackageFormat{Manifest: "Package.appxmanifest"}, []InputSpec{{Label: "package-resources", Files: []string{"Package.appxmanifest"}}}},
		{"darwin", "dmg", manifest.PackageFormat{Background: "background.png", VolumeIcon: "volume.icns", FileIcon: "file.icns", Files: map[string]string{"z": "z.txt", "a": "a.txt"}}, []InputSpec{{Label: "package-resources", Files: []string{"background.png", "volume.icns", "file.icns", "a.txt", "z.txt"}}}},
		{"linux", "appimage", manifest.PackageFormat{Icon: "icon.png", DesktopEntry: "app.desktop"}, []InputSpec{{Label: "package-resources", Files: []string{"icon.png", "app.desktop"}}}},
		{"linux", "deb", manifest.PackageFormat{PreInstall: "pre", PostInstall: "post", PreRemove: "prerm", PostRemove: "postrm"}, []InputSpec{{Label: "package-resources", Files: []string{"pre", "post", "prerm", "postrm"}}}},
		{"linux", "rpm", manifest.PackageFormat{Template: "template"}, []InputSpec{{Label: "package-template", Root: "template", IncludeAll: true}}},
		{"linux", "archlinux", manifest.PackageFormat{Template: "template.file"}, []InputSpec{{Label: "package-template", Files: []string{"template.file"}}}},
	}
	for _, test := range tests {
		t.Run(test.platform+"/"+test.format, func(t *testing.T) {
			assert.Equal(t, test.want, packageInputs(root, test.platform, test.format, test.config))
		})
	}
	assetConfig := manifest.Config{
		Project:      manifest.Project{Icon: "project.png"},
		Associations: []manifest.Association{{Icon: "association.png"}},
		Signing: manifest.Signing{
			Windows: manifest.SigningPlatform{Certificate: "private.pfx", Entitlements: "app.entitlements", ProvisioningProfile: "profile.mobileprovision"},
		},
	}
	assert.Equal(t, []InputSpec{{Label: "platform-assets", Files: []string{manifest.Filename, "project.png", "app.entitlements", "profile.mobileprovision", "association.png"}}}, assetInputs(assetConfig))

	config := manifest.Config{Build: manifest.Build{OutputDirectory: "bin"}, Project: manifest.Project{BinaryName: "app"}}
	assert.Equal(t, "bin/app.unknown", packageOutput(config, "linux", "amd64", "unknown", false))
}

func TestNonReproducibleArtifactsAndTheirPublishOperationsAreNeverCached(t *testing.T) {
	config := testConfig(t)
	config.Signing.Linux.Enabled = true
	config.Signing.Linux.Certificate = "release@example.com"
	plan, err := PlanBuild(config, Request{Verb: "sign", TargetOS: "linux", TargetArch: "amd64", Formats: []string{"deb"}})
	require.NoError(t, err)
	sign := plan.Nodes["package:linux/amd64:deb:sign"]
	publish := plan.Nodes["publish:linux/amd64:deb:sign"]
	assert.Equal(t, CacheNever, sign.Cache)
	assert.Equal(t, CacheNever, publish.Cache)
	assert.Empty(t, publish.Marker)
	assert.True(t, publish.Artifact.Signed)
	assert.Equal(t, CacheNever, plan.Nodes["collect:artifacts"].Cache, "a receipt over non-reproducible outputs must be regenerated")
}

func TestSigningInputsAreRewrittenToFrozenStagedCopies(t *testing.T) {
	signing := manifest.Signing{Darwin: manifest.SigningPlatform{
		Entitlements:        "user-assets/release.entitlements",
		ProvisioningProfile: "user-assets/distribution.mobileprovision",
	}}
	resolved := signingPlatformForPlan(signing, "darwin", ".wails/generated/darwin/assets")
	assert.Equal(t, ".wails/generated/darwin/assets/signing/entitlements.entitlements", resolved.Entitlements)
	assert.Equal(t, ".wails/generated/darwin/assets/signing/provisioning-profile.mobileprovision", resolved.ProvisioningProfile)
	assert.Equal(t, "user-assets/release.entitlements", signing.Darwin.Entitlements, "planning must not mutate user configuration")
}

func TestExecutorContinuesIndependentBranchesAfterFailure(t *testing.T) {
	plan := Plan{Name: "branches", Roots: []NodeKey{"blocked", "independent-child"}, Nodes: map[NodeKey]Node{
		"fail":              {Key: "fail", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheNever},
		"blocked":           {Key: "blocked", Kind: CompileApplication, Spec: CompileSpec{}, Dependencies: []NodeKey{"fail"}, Cache: CacheNever},
		"independent":       {Key: "independent", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheNever},
		"independent-child": {Key: "independent-child", Kind: CompileApplication, Spec: CompileSpec{}, Dependencies: []NodeKey{"independent"}, Cache: CacheNever},
	}}
	handler := &branchHandler{}
	reporter := &failureReporter{}
	_, err := (Executor{Handler: handler}).Execute(context.Background(), plan, ExecuteOptions{Root: t.TempDir(), Workers: 2, Reporter: reporter})
	require.ErrorContains(t, err, "fail: expected failure")
	assert.Contains(t, handler.runs, NodeKey("independent"))
	assert.Contains(t, handler.runs, NodeKey("independent-child"))
	assert.NotContains(t, handler.runs, NodeKey("blocked"))
	assert.Equal(t, "specific compiler output", reporter.failure.Output)
}

func TestExecutorPreservesCancellationWithoutReportingFailure(t *testing.T) {
	plan := Plan{Name: "cancel", Roots: []NodeKey{"work"}, Nodes: map[NodeKey]Node{
		"work": {Key: "work", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheNever},
	}}
	handler := &cancelHandler{started: make(chan struct{})}
	reporter := &failureReporter{}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := (Executor{Handler: handler}).Execute(ctx, plan, ExecuteOptions{Root: t.TempDir(), Reporter: reporter})
		done <- err
	}()
	<-handler.started
	cancel()
	require.ErrorIs(t, <-done, context.Canceled)
	assert.NoError(t, reporter.failure.Err)
}

func TestExecutorPreservesDeadlineWithoutReportingFailure(t *testing.T) {
	plan := Plan{Name: "deadline", Roots: []NodeKey{"work"}, Nodes: map[NodeKey]Node{
		"work": {Key: "work", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheNever},
	}}
	handler := &cancelHandler{started: make(chan struct{})}
	reporter := &failureReporter{}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	done := make(chan error, 1)
	go func() {
		_, err := (Executor{Handler: handler}).Execute(ctx, plan, ExecuteOptions{Root: t.TempDir(), Reporter: reporter})
		done <- err
	}()
	<-handler.started
	require.ErrorIs(t, <-done, context.DeadlineExceeded)
	assert.NoError(t, reporter.failure.Err)
}

func TestExecutorCancelsAtMaximumParallelismWithoutLeakingWorkers(t *testing.T) {
	const workers = 32
	plan := Plan{Name: "parallel-cancel", Roots: make([]NodeKey, workers), Nodes: make(map[NodeKey]Node, workers)}
	for index := range workers {
		key := NodeKey(fmt.Sprintf("work-%02d", index))
		plan.Roots[index] = key
		plan.Nodes[key] = Node{Key: key, Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheNever}
	}
	handler := &parallelCancelHandler{want: workers, started: make(chan struct{})}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := (Executor{Handler: handler}).Execute(ctx, plan, ExecuteOptions{Root: t.TempDir(), Workers: workers, MemoryLimitMB: workers})
		done <- err
	}()
	select {
	case <-handler.started:
	case <-time.After(5 * time.Second):
		t.Fatal("workers did not reach maximum parallelism")
	}
	cancel()
	select {
	case err := <-done:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(5 * time.Second):
		t.Fatal("executor did not release all workers after cancellation")
	}
	handler.mu.Lock()
	defer handler.mu.Unlock()
	assert.Zero(t, handler.running)
}

func testConfig(t *testing.T) manifest.Config {
	t.Helper()
	root := t.TempDir()
	require.NoError(t, manifest.WriteMinimal(root, manifest.Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"}))
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	return loaded.Config
}
