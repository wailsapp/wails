package pipeline

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/report"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func BenchmarkPlanMultiTarget(b *testing.B) {
	config := performanceConfig(b)
	request := Request{Verb: "build", Targets: []Target{{OS: "linux", Arch: "amd64"}, {OS: "windows", Arch: "amd64"}, {OS: "darwin", Arch: "arm64"}}}
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if _, err := PlanBuild(config, request); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPlanSingleTarget(b *testing.B) {
	config := performanceConfig(b)
	request := Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"}
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if _, err := PlanBuild(config, request); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPlanValidation1000Nodes(b *testing.B) {
	plan := syntheticPlan(1000)
	require.NoError(b, plan.Validate("."))
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if err := plan.Validate("."); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkInspectColdPlan(b *testing.B) {
	config := performanceConfig(b)
	plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	require.NoError(b, err)
	executor := Executor{Handler: &fakeHandler{root: config.Root}}
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if _, err := executor.Inspect(context.Background(), plan, config.Root); err != nil {
			b.Fatal(err)
		}
	}
}

func syntheticPlan(count int) Plan {
	plan := Plan{Nodes: make(map[NodeKey]Node, count)}
	for index := range count {
		key := NodeKey(fmt.Sprintf("node-%04d", index))
		dependencies := []NodeKey(nil)
		if index > 0 {
			dependencies = []NodeKey{NodeKey(fmt.Sprintf("node-%04d", index-1))}
		}
		plan.Nodes[key] = Node{Key: key, Kind: CompileApplication, Spec: CompileSpec{}, Dependencies: dependencies, Cache: CacheArtifact}
	}
	plan.Roots = []NodeKey{NodeKey(fmt.Sprintf("node-%04d", count-1))}
	return plan
}

func TestConcurrentPlanConstructionIsDeterministic(t *testing.T) {
	config := performanceConfig(t)
	request := Request{Verb: "build", Targets: []Target{{OS: "linux", Arch: "amd64"}, {OS: "linux", Arch: "arm64"}}}
	want, err := PlanBuild(config, request)
	require.NoError(t, err)
	const workers = 64
	results := make(chan Plan, workers)
	errors := make(chan error, workers)
	var group sync.WaitGroup
	for range workers {
		group.Add(1)
		go func() {
			defer group.Done()
			plan, err := PlanBuild(config, request)
			if err != nil {
				errors <- err
				return
			}
			results <- plan
		}()
	}
	group.Wait()
	close(results)
	close(errors)
	for err := range errors {
		require.NoError(t, err)
	}
	for result := range results {
		assert.Equal(t, want, result)
	}
}

func BenchmarkWarmNoopExecutor(b *testing.B) {
	config := performanceConfig(b)
	plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	require.NoError(b, err)
	handler := &fakeHandler{root: config.Root}
	executor := Executor{Handler: handler}
	options := ExecuteOptions{Root: config.Root, Reporter: report.Nop{}}
	_, err = executor.Execute(context.Background(), plan, options)
	require.NoError(b, err)
	b.ResetTimer()
	for range b.N {
		if _, err := executor.Execute(context.Background(), plan, options); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWarmNoopMultiTargetExecutor(b *testing.B) {
	config := performanceConfig(b)
	plan, err := PlanBuild(config, Request{Verb: "build", Targets: []Target{{OS: "linux", Arch: "amd64"}, {OS: "windows", Arch: "amd64"}, {OS: "darwin", Arch: "arm64"}}})
	require.NoError(b, err)
	handler := &fakeHandler{root: config.Root}
	executor := Executor{Handler: handler}
	options := ExecuteOptions{Root: config.Root, Reporter: report.Nop{}}
	_, err = executor.Execute(context.Background(), plan, options)
	require.NoError(b, err)
	b.ResetTimer()
	for range b.N {
		if _, err := executor.Execute(context.Background(), plan, options); err != nil {
			b.Fatal(err)
		}
	}
}

func TestWarmNoopExecutorRunsNoHandlers(t *testing.T) {
	config := performanceConfig(t)
	plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	require.NoError(t, err)
	handler := &fakeHandler{root: config.Root}
	executor := Executor{Handler: handler}
	options := ExecuteOptions{Root: config.Root, Reporter: report.Nop{}}
	_, err = executor.Execute(context.Background(), plan, options)
	require.NoError(t, err)
	runs := len(handler.runs)
	results, err := executor.Execute(context.Background(), plan, options)
	require.NoError(t, err)
	assert.Equal(t, runs, len(handler.runs), "a warm build should run no handlers")
	for key, result := range results {
		assert.NotEqual(t, "miss", string(result.Status), key)
	}
}

func TestIndependentLinuxPackageBranchesRunConcurrently(t *testing.T) {
	config := performanceConfig(t)
	plan, err := PlanBuild(config, Request{Verb: "package", TargetOS: "linux", TargetArch: "amd64", Formats: []string{"deb", "rpm", "archlinux"}})
	require.NoError(t, err)
	for _, format := range []string{"deb", "rpm", "archlinux"} {
		assert.Empty(t, plan.Nodes[NodeKey("package:linux/amd64:"+format)].Claims.Exclusive)
	}
	handler := newOverlapHandler(&fakeHandler{root: config.Root}, PackageArtifact, 3)
	_, err = (Executor{Handler: handler}).Execute(context.Background(), plan, ExecuteOptions{Root: config.Root, Workers: 3, Force: true, Reporter: report.Nop{}})
	require.NoError(t, err)
	assert.Equal(t, 3, handler.maximum(), "all independent Linux package adapters should overlap")
}

func TestAppImageIsProcessIsolatedAndDarwinBundleAdaptersRemainExclusive(t *testing.T) {
	linux := performanceConfig(t)
	linuxPlan, err := PlanBuild(linux, Request{Verb: "package", TargetOS: "linux", TargetArch: "amd64", Formats: []string{"appimage", "deb"}})
	require.NoError(t, err)
	assert.Empty(t, linuxPlan.Nodes["package:linux/amd64:appimage"].Claims.Exclusive, "AppImage runs through an isolated CLI subprocess")
	assert.Empty(t, linuxPlan.Nodes["package:linux/amd64:deb"].Claims.Exclusive)

	darwin := performanceConfig(t)
	darwinPlan, err := PlanBuild(darwin, Request{Verb: "package", TargetOS: "darwin", TargetArch: "arm64", Formats: []string{"dmg"}})
	require.NoError(t, err)
	appClaim := darwinPlan.Nodes["assemble:darwin/arm64"].Claims.Exclusive
	assert.NotEmpty(t, appClaim)
	assert.Equal(t, appClaim, darwinPlan.Nodes["package:darwin/arm64:dmg"].Claims.Exclusive, "app assembly and DMG packaging share the same bundle workspace")
}

func TestIndependentMultiTargetCompilesRunConcurrently(t *testing.T) {
	config := performanceConfig(t)
	plan, err := PlanBuild(config, Request{Verb: "build", Targets: []Target{{OS: "linux", Arch: "amd64"}, {OS: "linux", Arch: "arm64"}}})
	require.NoError(t, err)
	amd64 := plan.Nodes["target:linux/amd64:compile"]
	arm64 := plan.Nodes["target:linux/arm64:compile"]
	assert.Empty(t, amd64.Claims.Exclusive)
	assert.Empty(t, arm64.Claims.Exclusive)
	workers := amd64.Claims.CPU + arm64.Claims.CPU
	handler := newOverlapHandler(&fakeHandler{root: config.Root}, CompileApplication, 2)
	_, err = (Executor{Handler: handler}).Execute(context.Background(), plan, ExecuteOptions{Root: config.Root, Workers: workers, MemoryLimitMB: amd64.Claims.MemoryMB + arm64.Claims.MemoryMB, Force: true, Reporter: report.Nop{}})
	require.NoError(t, err)
	assert.Equal(t, 2, handler.maximum(), "independent Target compiles should overlap when their CPU claims fit")
}

func TestMemoryClaimsBoundOtherwiseIndependentWork(t *testing.T) {
	root := t.TempDir()
	plan := Plan{Name: "memory", Roots: []NodeKey{"one", "two"}, Nodes: map[NodeKey]Node{
		"one": {Key: "one", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheNever, Claims: ResourceClaims{CPU: 1, MemoryMB: 1024}},
		"two": {Key: "two", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheNever, Claims: ResourceClaims{CPU: 1, MemoryMB: 1024}},
	}}
	handler := newOverlapHandler(&fakeHandler{root: root}, CompileApplication, 2)
	_, err := (Executor{Handler: handler}).Execute(context.Background(), plan, ExecuteOptions{Root: root, Workers: 2, MemoryLimitMB: 1024, Force: true, Reporter: report.Nop{}})
	require.NoError(t, err)
	assert.Equal(t, 1, handler.maximum(), "memory capacity should serialize Nodes whose claims do not fit together")
}

type overlapHandler struct {
	delegate Handler
	kind     NodeKind
	want     int
	ready    chan struct{}
	once     sync.Once

	mu      sync.Mutex
	running int
	max     int
}

func newOverlapHandler(delegate Handler, kind NodeKind, want int) *overlapHandler {
	return &overlapHandler{delegate: delegate, kind: kind, want: want, ready: make(chan struct{})}
}

func (h *overlapHandler) Identity(ctx context.Context, node Node) (string, error) {
	return h.delegate.Identity(ctx, node)
}

func (h *overlapHandler) Run(ctx context.Context, node Node) (RunResult, error) {
	if node.Kind != h.kind {
		return h.delegate.Run(ctx, node)
	}
	h.mu.Lock()
	h.running++
	if h.running > h.max {
		h.max = h.running
	}
	if h.running == h.want {
		h.once.Do(func() { close(h.ready) })
	}
	h.mu.Unlock()
	select {
	case <-h.ready:
	case <-ctx.Done():
		return RunResult{}, ctx.Err()
	case <-time.After(500 * time.Millisecond):
	}
	result, err := h.delegate.Run(ctx, node)
	h.mu.Lock()
	h.running--
	h.mu.Unlock()
	return result, err
}

func (h *overlapHandler) maximum() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.max
}

type testingTempDir interface {
	Helper()
	TempDir() string
	Fatalf(string, ...any)
}

func performanceConfig(t testingTempDir) manifest.Config {
	t.Helper()
	must := func(err error) {
		if err != nil {
			t.Fatalf("performance fixture: %v", err)
		}
	}
	root := t.TempDir()
	must(os.MkdirAll(filepath.Join(root, "frontend", "node_modules"), 0o755))
	must(os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example\n"), 0o644))
	must(os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\nfunc main(){}\n"), 0o644))
	must(os.WriteFile(filepath.Join(root, "frontend", "package.json"), []byte("{}"), 0o644))
	must(os.WriteFile(filepath.Join(root, "frontend", "package-lock.json"), []byte("{}"), 0o644))
	must(manifest.WriteMinimal(root, manifest.Project{Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0"}))
	loaded, err := manifest.Load(root, "")
	must(err)
	return loaded.Config
}
