package pipeline

import (
	"github.com/wailsapp/wails/v3/internal/wake/cache"
	"os"
	"path/filepath"
	"testing"
)

func TestCompileInputsTrackEmbeddedResources(t *testing.T) {
	c := testConfig(t)
	if err := os.WriteFile(filepath.Join(c.Root, "main.go"), []byte("package main\nimport _ \"embed\"\n//go:embed data.json\nvar data string\nfunc main() {}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(c.Root, "data.json")
	if err := os.WriteFile(p, []byte("before"), 0600); err != nil {
		t.Fatal(err)
	}
	plan, err := PlanBuild(c, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	if err != nil {
		t.Fatal(err)
	}
	node := plan.Nodes["target:linux/amd64:compile"]
	store, err := cache.OpenCache(c.Root)
	if err != nil {
		t.Fatal(err)
	}
	before, err := snapshotNodeInputs(store, node)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("after"), 0600); err != nil {
		t.Fatal(err)
	}
	store.InvalidateObservations()
	after, err := snapshotNodeInputs(store, node)
	if err != nil {
		t.Fatal(err)
	}
	changed := false
	for i := range before {
		if before[i] != after[i] {
			changed = true
		}
	}
	if !changed {
		t.Fatal("embedded resource edit must invalidate compile inputs")
	}
}
