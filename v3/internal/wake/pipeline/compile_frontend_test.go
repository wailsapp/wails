package pipeline

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/cache"
)

type frontendCompileHandler struct {
	root     string
	generate func() error
}

func (*frontendCompileHandler) Identity(context.Context, Node) (string, error) {
	return "compile-regression-v1", nil
}

func (h *frontendCompileHandler) Run(ctx context.Context, node Node) (RunResult, error) {
	if node.Kind == BuildFrontend {
		if h.generate != nil {
			if err := h.generate(); err != nil {
				return RunResult{}, err
			}
		}
		content, err := os.ReadFile(filepath.Join(h.root, "frontend/src/page.txt"))
		if err != nil {
			return RunResult{}, err
		}
		return RunResult{}, os.WriteFile(filepath.Join(h.root, "frontend/dist/page.txt"), content, 0644)
	}
	spec := node.Spec.(CompileSpec)
	command := exec.CommandContext(ctx, filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-buildvcs=false", "-o", spec.Output, spec.Package)
	command.Dir = h.root
	command.Env = append(os.Environ(), "CGO_ENABLED=0", "GOWORK=off")
	output, err := command.CombinedOutput()
	if err != nil {
		return RunResult{}, fmt.Errorf("go build: %w: %s", err, output)
	}
	return RunResult{}, nil
}

func writeCompileFixture(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func frontendCompilePlan(t *testing.T) (Plan, *frontendCompileHandler) {
	t.Helper()
	root := t.TempDir()
	for name, content := range map[string]string{
		"go.mod":                 "module example.test/app\ngo 1.26\n",
		"main.go":                "package main\nimport (\"fmt\"; _ \"embed\"; \"example.test/app/shared\")\n//go:embed frontend/dist/page.txt\nvar page string\nfunc main(){fmt.Print(page+shared.Value)}\n",
		"cmd/helper/main.go":     "package main\nimport (\"fmt\"; \"example.test/app/shared\")\nfunc main(){fmt.Print(shared.Value)}\n",
		"shared/value.go":        "package shared\nconst Value = \"shared-one\"\n",
		"frontend/src/page.txt":  "one/",
		"frontend/dist/page.txt": "old/",
	} {
		writeCompileFixture(t, root, name, content)
	}
	plan := Plan{Roots: []NodeKey{"app", "helper"}, Nodes: map[NodeKey]Node{
		"frontend": {Key: "frontend", Kind: BuildFrontend, Spec: FrontendSpec{}, Inputs: []InputSpec{{Label: "frontend", Root: "frontend/src", IncludeAll: true}}, Output: "frontend/dist", Cache: CacheArtifact},
	}}
	for _, name := range []string{"app", "helper"} {
		entry := "."
		if name == "helper" {
			entry = "./cmd/helper"
		}
		output := "bin/" + name
		if runtime.GOOS == "windows" {
			output += ".exe"
		}
		plan.Nodes[NodeKey(name)] = Node{Key: NodeKey(name), Kind: CompileApplication, Spec: CompileSpec{Package: entry, CGOEnabled: "0", Output: output}, Dependencies: []NodeKey{"frontend"}, Inputs: []InputSpec{{Label: "go-sources", Root: ".", IncludeGoEmbed: true, IncludeExtensions: goSourceExtensions(), ExcludeDirs: []string{".wails", "bin"}, ExcludeSuffixes: []string{"_test.go"}}}, Output: output, Cache: CacheArtifact}
	}
	refineCompileInputs(&plan, root)
	return plan, &frontendCompileHandler{root: root}
}

func TestFrontendEditsReuseIndependentCompanion(t *testing.T) {
	plan, handler := frontendCompilePlan(t)
	executor := Executor{Handler: handler}
	build := func(app, helper cache.LookupStatus) {
		t.Helper()
		results, err := executor.Execute(t.Context(), plan, ExecuteOptions{Root: handler.root, Workers: 2})
		if err != nil {
			t.Fatal(err)
		}
		if results["app"].Status != app || results["helper"].Status != helper {
			t.Fatalf("compile statuses app=%s helper=%s, want %s/%s", results["app"].Status, results["helper"].Status, app, helper)
		}
	}
	build(cache.LookupMiss, cache.LookupMiss)
	build(cache.LookupHit, cache.LookupHit)
	if err := os.Remove(filepath.Join(handler.root, "frontend/dist/page.txt")); err != nil {
		t.Fatal(err)
	}
	inspection, err := executor.Inspect(t.Context(), plan, handler.root)
	if err != nil {
		t.Fatal(err)
	}
	if inspection.Operations["app"].Decision != "run" {
		t.Fatal("inspection promised a cache hit before restoring embedded output")
	}
	build(cache.LookupHit, cache.LookupHit)
	writeCompileFixture(t, handler.root, "frontend/src/page.txt", "two/")
	inspection, err = executor.Inspect(t.Context(), plan, handler.root)
	if err != nil {
		t.Fatal(err)
	}
	if inspection.Operations["app"].Decision != "run" {
		t.Fatal("inspection promised a cache hit before new embedded output existed")
	}
	build(cache.LookupMiss, cache.LookupHit)
	output, err := exec.CommandContext(t.Context(), filepath.Join(handler.root, plan.Nodes["app"].Output)).CombinedOutput()
	if err != nil || string(output) != "two/shared-one" {
		t.Fatalf("app output %q: %v", output, err)
	}
	writeCompileFixture(t, handler.root, "shared/value.go", "package shared\nconst Value = \"shared-two\"\n")
	build(cache.LookupMiss, cache.LookupMiss)
	// Reusing a plan must discover both newly generated imports and later edits
	// to those imports after the frontend stage has finished.
	handler.generate = func() error {
		if err := os.MkdirAll(filepath.Join(handler.root, "generated"), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(handler.root, "cmd/helper/generated.go"), []byte("package main\nimport _ \"example.test/app/generated\"\n"), 0644); err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(handler.root, "generated/value.go"), []byte("package generated\nconst Value = 1\n"), 0644)
	}
	writeCompileFixture(t, handler.root, "frontend/src/page.txt", "three/")
	build(cache.LookupMiss, cache.LookupMiss)
	writeCompileFixture(t, handler.root, "generated/value.go", "package generated\nconst Value = 2\n")
	build(cache.LookupHit, cache.LookupMiss)
	handler.generate = func() error { return errors.New("frontend failed") }
	writeCompileFixture(t, handler.root, "frontend/src/page.txt", "four/")
	results, err := executor.Execute(t.Context(), plan, ExecuteOptions{Root: handler.root, Workers: 2})
	if err == nil {
		t.Fatal("failed frontend accepted")
	}
	if _, ok := results["helper"]; ok {
		t.Fatal("companion ran despite failed frontend dependency")
	}
}

func TestOrderOnlyFrontendRequiresCompletePureGoInputs(t *testing.T) {
	for _, mode := range []string{"0", "1"} {
		t.Run(mode, func(t *testing.T) {
			plan, handler := frontendCompilePlan(t)
			node := plan.Nodes["helper"]
			spec := node.Spec.(CompileSpec)
			spec.CGOEnabled = mode
			node.Spec, node.OrderOnlyDependencies = spec, nil
			plan.Nodes["helper"] = node
			refineCompileInputs(&plan, handler.root)
			if got := slices.Contains(plan.Nodes["helper"].OrderOnlyDependencies, "frontend"); got != (mode == "0") {
				t.Fatalf("CGO=%s order-only=%t", mode, got)
			}
			writeCompileFixture(t, handler.root, "cmd/helper/generated.go", "package main\nimport _ \"example.test/app/not-yet-generated\"\n")
			if got := resolveCompileInputs(plan.Nodes["helper"]); len(got.OrderOnlyDependencies) != 0 {
				t.Fatal("incomplete generated graph must retain frontend artifact dependency")
			}
		})
	}
}

func TestOrderOnlyDependenciesAreValidatedAndChangeActionIdentity(t *testing.T) {
	plan, _ := frontendCompilePlan(t)
	node := plan.Nodes["helper"]
	before, err := cache.ActionKey(string(node.Kind), nodeActionSpec(node, "tool"), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	node.OrderOnlyDependencies = nil
	after, err := cache.ActionKey(string(node.Kind), nodeActionSpec(node, "tool"), nil, nil)
	if err != nil || before == after {
		t.Fatalf("dependency role did not change identity: %v", err)
	}
	for _, invalid := range [][]NodeKey{{"missing"}, {"app"}, {"frontend", "frontend"}} {
		node.OrderOnlyDependencies = invalid
		plan.Nodes["helper"] = node
		if err := plan.Validate("."); err == nil {
			t.Fatalf("accepted invalid order-only dependencies %v", invalid)
		}
	}
}

func TestFrontendTracksIgnoredGeneratedGoSources(t *testing.T) {
	plan, handler := frontendCompilePlan(t)
	writeCompileFixture(t, handler.root, "cmd/helper/main.go", "package main\nimport \"fmt\"\nfunc main(){fmt.Print(generated)}\n")
	writeCompileFixture(t, handler.root, "cmd/helper/generated.go", "package main\nconst generated = \"before\"\n")
	writeCompileFixture(t, handler.root, "cmd/helper/.gitignore", "generated.go\n")
	node := plan.Nodes["helper"]
	node.Inputs[0].UseGitIgnore = true
	plan.Nodes["helper"] = node
	refineCompileInputs(&plan, handler.root)
	executor := Executor{Handler: handler}
	if _, err := executor.Execute(t.Context(), plan, ExecuteOptions{Root: handler.root}); err != nil {
		t.Fatal(err)
	}
	handler.generate = func() error {
		return os.WriteFile(filepath.Join(handler.root, "cmd/helper/generated.go"), []byte("package main\nconst generated = \"after\"\n"), 0644)
	}
	writeCompileFixture(t, handler.root, "frontend/src/page.txt", "changed/")
	if _, err := executor.Execute(t.Context(), plan, ExecuteOptions{Root: handler.root}); err != nil {
		t.Fatal(err)
	}
	output, err := exec.CommandContext(t.Context(), filepath.Join(handler.root, plan.Nodes["helper"].Output)).CombinedOutput()
	if err != nil || string(output) != "after" {
		t.Fatalf("ignored generated Go source stayed stale: %q, %v", output, err)
	}
}
