package pipeline

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/report"
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

func (*cancelHandler) Identity(context.Context, Node) (string, error) { return "cancel-v1", nil }
func (h *cancelHandler) Run(ctx context.Context, _ Node) (RunResult, error) {
	close(h.started)
	<-ctx.Done()
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
	assert.Len(t, plan.Nodes, 4)
	h := &fakeHandler{root: root}
	executor := Executor{Handler: h}
	results, err := executor.Execute(context.Background(), plan, ExecuteOptions{Root: root, Reporter: report.Nop{}})
	require.NoError(t, err)
	assert.Empty(t, results["frontend:install"].Artifact, "a receipt is not an artifact")
	first := len(h.runs)
	assert.Equal(t, 4, first)
	_, err = executor.Execute(context.Background(), plan, ExecuteOptions{Root: root, Reporter: report.Nop{}})
	require.NoError(t, err)
	assert.Equal(t, first, len(h.runs))
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

func TestPlannerBuildsOnePlanForMultipleTargets(t *testing.T) {
	config := testConfig(t)
	plan, err := PlanBuild(config, Request{Verb: "build", Targets: []Target{{OS: "linux", Arch: "amd64"}, {OS: "linux", Arch: "arm64"}}})
	require.NoError(t, err)
	assert.Equal(t, "linux/amd64,linux/arm64", plan.Target)
	assert.Contains(t, plan.Nodes, NodeKey("frontend:build"))
	assert.Contains(t, plan.Nodes, NodeKey("target:linux/amd64:compile"))
	assert.Contains(t, plan.Nodes, NodeKey("target:linux/arm64:compile"))
	assert.Len(t, plan.Roots, 2)
	amd64 := plan.Nodes["target:linux/amd64:compile"].Output
	arm64 := plan.Nodes["target:linux/arm64:compile"].Output
	assert.NotEqual(t, amd64, arm64)
}

func TestPlannerRejectsInvalidTargetAndFormatCombinations(t *testing.T) {
	config := testConfig(t)
	_, err := PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "riscv64"})
	require.ErrorContains(t, err, "unsupported target architecture")
	_, err = PlanBuild(config, Request{Verb: "package", TargetOS: "linux", TargetArch: "amd64", Formats: []string{"nsis"}})
	require.ErrorContains(t, err, "not supported for linux")
	_, err = PlanBuild(config, Request{Verb: "build", TargetOS: "darwin", TargetArch: "universal"})
	require.NoError(t, err)
}

func TestPlannerCarriesPackageCustomizationForEveryFormat(t *testing.T) {
	tests := []struct {
		platform string
		format   string
	}{
		{"windows", "nsis"}, {"windows", "msix"},
		{"darwin", "app"}, {"darwin", "dmg"},
		{"linux", "appimage"}, {"linux", "deb"}, {"linux", "rpm"}, {"linux", "archlinux"},
		{"ios", "app"}, {"ios", "ipa"},
		{"android", "apk"}, {"android", "aab"},
	}
	for _, test := range tests {
		t.Run(test.platform+"/"+test.format, func(t *testing.T) {
			config := testConfig(t)
			if test.platform == "ios" && test.format == "ipa" {
				config.Targets.IOS.ARM64.Variant = "device"
			}
			format := packageFormatPointer(&config.Package, test.platform, test.format)
			format.Template = "packaging/" + test.platform + "/" + test.format + ".tmpl"
			format.Options = map[string]any{"channel": "preview"}
			config.Associations = []manifest.Association{{Extensions: []string{"example"}}}
			config.Protocols = []manifest.Protocol{{Scheme: "example"}}
			plan, err := PlanBuild(config, Request{Verb: "package", TargetOS: test.platform, TargetArch: "arm64", Formats: []string{test.format}})
			require.NoError(t, err)
			node := plan.Nodes[NodeKey("package:"+test.platform+"/arm64:"+test.format)]
			spec := node.Spec.(PackageSpec)
			assert.Equal(t, format.Template, spec.Config.Template)
			assert.Equal(t, "preview", spec.Config.Options["channel"])
			assert.Equal(t, config.Root, spec.TemplateRoot)
			assert.Equal(t, config.Associations, spec.Associations)
			assert.Equal(t, config.Protocols, spec.Protocols)
			require.Len(t, node.Inputs, 1)
			assert.Equal(t, "package-template", node.Inputs[0].Label)
		})
	}
}

func TestPlannerSnapshotsBuiltInDMGResources(t *testing.T) {
	config := testConfig(t)
	config.Package.Darwin.DMG.Options = map[string]any{
		"background": "packaging/background.png",
		"files":      "Read Me=README.md,License=LICENSE",
	}
	plan, err := PlanBuild(config, Request{Verb: "package", TargetOS: "darwin", TargetArch: "arm64", Formats: []string{"dmg"}})
	require.NoError(t, err)
	node := plan.Nodes[NodeKey("package:darwin/arm64:dmg")]
	require.Len(t, node.Inputs, 1)
	assert.Equal(t, "package-resources", node.Inputs[0].Label)
	assert.ElementsMatch(t, []string{"packaging/background.png", "README.md", "LICENSE"}, node.Inputs[0].Files)
}

func packageFormatPointer(packages *manifest.Packages, platform, format string) *manifest.PackageFormat {
	var selected *manifest.PackagePlatform
	switch platform {
	case "windows":
		selected = &packages.Windows
	case "darwin":
		selected = &packages.Darwin
	case "linux":
		selected = &packages.Linux
	case "ios":
		selected = &packages.IOS
	case "android":
		selected = &packages.Android
	}
	switch format {
	case "nsis":
		return &selected.NSIS
	case "msix":
		return &selected.MSIX
	case "app":
		return &selected.App
	case "dmg":
		return &selected.DMG
	case "appimage":
		return &selected.AppImage
	case "deb":
		return &selected.Deb
	case "rpm":
		return &selected.RPM
	case "archlinux":
		return &selected.ArchLinux
	case "ipa":
		return &selected.IPA
	case "apk":
		return &selected.APK
	default:
		return &selected.AAB
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
	plan, err := PlanBuild(config, Request{Verb: "package", TargetOS: "darwin", TargetArch: "arm64", Formats: []string{"app"}})
	require.NoError(t, err)
	compile := plan.Nodes["target:darwin/arm64:compile"].Spec.(CompileSpec)
	assert.Equal(t, "13.0", compile.MinimumVersion)
	assert.Contains(t, compile.Tags, "metal")
	pkg := plan.Nodes["package:darwin/arm64:app"].Spec.(PackageSpec)
	assert.Equal(t, "Mac App", pkg.Project.ProductName)
	assert.Equal(t, "com.example.mac", pkg.Project.Identifier)
	assert.Equal(t, 8, pkg.Project.BuildNumber)
}

func TestMobilePackagePlansUseTargetVariants(t *testing.T) {
	config := testConfig(t)
	for _, tc := range []struct {
		os, format, variant string
	}{{"ios", "app", "simulator"}, {"android", "apk", ""}} {
		t.Run(tc.os, func(t *testing.T) {
			plan, err := PlanBuild(config, Request{Verb: "package", TargetOS: tc.os, TargetArch: "arm64", Formats: []string{tc.format}})
			require.NoError(t, err)
			key := NodeKey("package:" + tc.os + "/arm64:" + tc.format)
			node, ok := plan.Nodes[key]
			require.True(t, ok)
			spec := node.Spec.(PackageSpec)
			assert.Equal(t, tc.variant, spec.Variant)
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
		{"darwin", "amd64", []string{"app", "dmg"}},
		{"darwin", "arm64", []string{"app", "dmg"}},
		{"darwin", "universal", []string{"app", "dmg"}},
		{"linux", "amd64", []string{"appimage", "deb", "rpm", "archlinux"}},
		{"linux", "arm64", []string{"appimage", "deb", "rpm", "archlinux"}},
		{"ios", "arm64", []string{"app", "ipa"}},
		{"android", "arm64", []string{"apk", "aab"}},
		{"android", "amd64", []string{"apk", "aab"}},
	}
	for _, test := range tests {
		for _, format := range test.formats {
			t.Run(test.platform+"/"+test.arch+"/"+format, func(t *testing.T) {
				config := testConfig(t)
				if test.platform == "ios" && format == "ipa" {
					config.Targets.IOS.ARM64.Variant = "device"
				}
				plan, err := PlanBuild(config, Request{Verb: "package", TargetOS: test.platform, TargetArch: test.arch, Formats: []string{format}})
				require.NoError(t, err)
				key := NodeKey("package:" + test.platform + "/" + test.arch + ":" + format)
				assert.Contains(t, plan.Nodes, key)
				assert.Contains(t, plan.Roots, key)
			})
		}
	}
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
	require.ErrorContains(t, err, `variant = "device"`)
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

func TestPackageAndSignHooksRemainBarriersForEveryFormat(t *testing.T) {
	config := testConfig(t)
	config.Hooks.AfterBuild.Script = "scripts/after-build.sh"
	config.Hooks.BeforePackage.Script = "scripts/before-package.sh"
	config.Hooks.AfterPackage.Script = "scripts/after-package.sh"
	config.Hooks.BeforeSign.Script = "scripts/before-sign.sh"
	config.Hooks.AfterSign.Script = "scripts/after-sign.sh"
	plan, err := PlanBuild(config, Request{Verb: "sign", TargetOS: "linux", TargetArch: "amd64", Formats: []string{"deb", "rpm"}})
	require.NoError(t, err)
	for _, format := range []string{"deb", "rpm"} {
		packageKey := NodeKey("package:linux/amd64:" + format)
		packageNode := plan.Nodes[packageKey]
		assert.Contains(t, packageNode.Dependencies, NodeKey("hook:before_package:linux/amd64"))
		signNode, ok := plan.Nodes[NodeKey(string(packageKey)+":sign")]
		require.True(t, ok)
		assert.Equal(t, plan.Nodes[packageKey].Output, signNode.Spec.(SignSpec).Input)
		assert.Contains(t, signNode.Dependencies, packageKey)
		assert.Contains(t, signNode.Dependencies, NodeKey("hook:before_sign:linux/amd64"))
	}
	packageOutputs := []string{plan.Nodes["package:linux/amd64:deb"].Output, plan.Nodes["package:linux/amd64:rpm"].Output}
	assert.Equal(t, commonOutput(packageOutputs), plan.Nodes["hook:before_package:linux/amd64"].Spec.(HookSpec).Artifact)
	assert.Equal(t, TargetScope, plan.Nodes["hook:before_package:linux/amd64"].Scope)
	assert.Equal(t, commonOutput(packageOutputs), plan.Nodes["hook:after_package:linux/amd64"].Spec.(HookSpec).Artifact)
	assert.Equal(t, PackageScope, plan.Nodes["hook:after_package:linux/amd64"].Scope)
	assert.Equal(t, commonOutput(packageOutputs), plan.Nodes["hook:before_sign:linux/amd64"].Spec.(HookSpec).Artifact)
	signedOutputs := []string{packageOutputs[0] + ".signed", packageOutputs[1] + ".signed"}
	assert.Equal(t, commonOutput(signedOutputs), plan.Nodes["hook:after_sign:linux/amd64"].Spec.(HookSpec).Artifact)
	assert.Equal(t, []NodeKey{"hook:after_sign:linux/amd64"}, plan.Roots)
}

func TestHookSpecsExposeThePhaseArtifactSeparatelyFromCacheOutputs(t *testing.T) {
	config := testConfig(t)
	config.Hooks.BeforeBuild = manifest.Hook{Script: "scripts/before-build.sh", Cache: true, Inputs: []string{"version.txt"}, Outputs: []string{"generated/version.go"}}
	config.Hooks.AfterBuild.Script = "scripts/after-build.sh"
	config.Hooks.BeforePackage.Script = "scripts/before-package.sh"
	config.Hooks.AfterPackage.Script = "scripts/after-package.sh"
	plan, err := PlanBuild(config, Request{Verb: "package", TargetOS: "linux", TargetArch: "amd64", Formats: []string{"deb"}})
	require.NoError(t, err)

	binary := plan.Nodes["target:linux/amd64:compile"].Output
	artifact := plan.Nodes["package:linux/amd64:deb"].Output
	beforeBuild := plan.Nodes["hook:before_build:project"]
	assert.Equal(t, "generated/version.go", beforeBuild.Output, "Node output belongs to hook caching")
	beforeBuildSpec := beforeBuild.Spec.(HookSpec)
	assert.Empty(t, beforeBuildSpec.Artifact, "project preparation has no phase Artifact yet")
	assert.Empty(t, beforeBuildSpec.TargetOS)
	assert.Empty(t, beforeBuildSpec.TargetArch)
	assert.Equal(t, "default", beforeBuildSpec.Profile)
	assert.Equal(t, binary, plan.Nodes["hook:after_build:linux/amd64"].Spec.(HookSpec).Artifact)
	assert.Equal(t, artifact, plan.Nodes["hook:before_package:linux/amd64"].Spec.(HookSpec).Artifact)
	assert.Equal(t, artifact, plan.Nodes["hook:after_package:linux/amd64"].Spec.(HookSpec).Artifact)
}

func TestPlannerRejectsInvalidCachedHookOutputsFromDirectConfig(t *testing.T) {
	config := testConfig(t)
	config.Hooks.BeforeBuild = manifest.Hook{Script: "scripts/build.sh", Cache: true, Inputs: []string{"main.go"}, Outputs: []string{"one.txt", "two.txt"}}
	_, err := PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	require.ErrorContains(t, err, "hook before_build outputs")
}

func TestProjectBeforeBuildHookIsSharedAcrossMultiTargetPlan(t *testing.T) {
	config := testConfig(t)
	config.Hooks.BeforeBuild.Script = "scripts/before-build.sh"
	plan, err := PlanBuild(config, Request{Verb: "build", Targets: []Target{{OS: "linux", Arch: "amd64"}, {OS: "windows", Arch: "amd64"}, {OS: "darwin", Arch: "arm64"}}})
	require.NoError(t, err)

	hook, ok := plan.Nodes["hook:before_build:project"]
	require.True(t, ok)
	assert.Equal(t, ProjectScope, hook.Scope)
	spec := hook.Spec.(HookSpec)
	assert.Empty(t, spec.TargetOS)
	assert.Empty(t, spec.TargetArch)
	assert.Empty(t, spec.Artifact)
	for key, node := range plan.Nodes {
		if strings.HasPrefix(string(key), "hook:before_build:") {
			assert.Equal(t, NodeKey("hook:before_build:project"), key)
		}
		if key == "frontend:bindings" {
			assert.Contains(t, node.Dependencies, NodeKey("hook:before_build:project"))
		}
	}
}

func TestExecutorContinuesIndependentBranchesAfterFailure(t *testing.T) {
	plan := Plan{Name: "branches", Roots: []NodeKey{"blocked", "independent-child"}, Nodes: map[NodeKey]Node{
		"fail":              {Key: "fail", Kind: RunHook, Cache: CacheNever},
		"blocked":           {Key: "blocked", Kind: RunHook, Dependencies: []NodeKey{"fail"}, Cache: CacheNever},
		"independent":       {Key: "independent", Kind: RunHook, Cache: CacheNever},
		"independent-child": {Key: "independent-child", Kind: RunHook, Dependencies: []NodeKey{"independent"}, Cache: CacheNever},
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
		"work": {Key: "work", Kind: RunHook, Cache: CacheNever},
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
		"work": {Key: "work", Kind: RunHook, Cache: CacheNever},
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

func testConfig(t *testing.T) manifest.Config {
	t.Helper()
	root := t.TempDir()
	require.NoError(t, manifest.WriteMinimal(root, manifest.Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"}))
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	return loaded.Config
}
