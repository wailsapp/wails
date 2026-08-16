package pipeline

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"

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

func (*branchHandler) Identity(context.Context, Node) (string, error) { return "branch-v1", nil }
func (h *branchHandler) Run(_ context.Context, node Node) (RunResult, error) {
	h.mu.Lock()
	h.runs = append(h.runs, node.Key)
	h.mu.Unlock()
	if node.Key == "fail" {
		return RunResult{}, errors.New("expected failure")
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
	_, err = executor.Execute(context.Background(), plan, ExecuteOptions{Root: root, Reporter: report.Nop{}})
	require.NoError(t, err)
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

func TestPlannerRejectsInvalidTargetAndFormatCombinations(t *testing.T) {
	config := testConfig(t)
	_, err := PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "riscv64"})
	require.ErrorContains(t, err, "unsupported target architecture")
	_, err = PlanBuild(config, Request{Verb: "package", TargetOS: "linux", TargetArch: "amd64", Formats: []string{"nsis"}})
	require.ErrorContains(t, err, "not supported for linux")
	_, err = PlanBuild(config, Request{Verb: "build", TargetOS: "darwin", TargetArch: "universal"})
	require.NoError(t, err)
	config.Package.Darwin.DMG.Template = "packaging/custom-dmg"
	_, err = PlanBuild(config, Request{Verb: "package", TargetOS: "darwin", TargetArch: "arm64", Formats: []string{"dmg"}})
	require.ErrorContains(t, err, "custom templates are not supported")
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
		})
	}
}

func TestPackageAndSignHooksRemainBarriersForEveryFormat(t *testing.T) {
	config := testConfig(t)
	config.Hooks.AfterBuild.Script = "scripts/after-build.sh"
	config.Hooks.BeforePackage.Script = "scripts/before-package.sh"
	config.Hooks.AfterPackage.Script = "scripts/after-package.sh"
	config.Hooks.BeforeSign.Script = "scripts/before-sign.sh"
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
}

func TestExecutorContinuesIndependentBranchesAfterFailure(t *testing.T) {
	plan := Plan{Name: "branches", Roots: []NodeKey{"blocked", "independent-child"}, Nodes: map[NodeKey]Node{
		"fail":              {Key: "fail", Kind: RunHook, Cache: CacheNever},
		"blocked":           {Key: "blocked", Kind: RunHook, Dependencies: []NodeKey{"fail"}, Cache: CacheNever},
		"independent":       {Key: "independent", Kind: RunHook, Cache: CacheNever},
		"independent-child": {Key: "independent-child", Kind: RunHook, Dependencies: []NodeKey{"independent"}, Cache: CacheNever},
	}}
	handler := &branchHandler{}
	_, err := (Executor{Handler: handler}).Execute(context.Background(), plan, ExecuteOptions{Root: t.TempDir(), Workers: 2, Reporter: report.Nop{}})
	require.ErrorContains(t, err, "fail: expected failure")
	assert.Contains(t, handler.runs, NodeKey("independent"))
	assert.Contains(t, handler.runs, NodeKey("independent-child"))
	assert.NotContains(t, handler.runs, NodeKey("blocked"))
}

func testConfig(t *testing.T) manifest.Config {
	t.Helper()
	root := t.TempDir()
	require.NoError(t, manifest.WriteMinimal(root, manifest.Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"}))
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	return loaded.Config
}
