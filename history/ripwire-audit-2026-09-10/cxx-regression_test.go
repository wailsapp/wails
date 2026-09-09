package pipeline

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/cache"
)

func TestReviewCompileInputsTrackCXX(t *testing.T) {
	c := testConfig(t)
	if err := os.WriteFile(filepath.Join(c.Root, "main.go"), []byte("package main\nfunc main() {}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(c.Root, "helper.cxx")
	if err := os.WriteFile(p, []byte("extern \"C\" int answer() { return 1; }"), 0600); err != nil {
		t.Fatal(err)
	}
	plan, err := PlanBuild(c, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	if err != nil {
		t.Fatal(err)
	}
	node, ok := plan.Nodes["target:linux/amd64:compile"]
	if !ok {
		t.Fatal("plan is missing node target:linux/amd64:compile")
	}
	store, err := cache.OpenCache(c.Root)
	if err != nil {
		t.Fatal(err)
	}
	before, err := snapshotNodeInputs(store, node)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("extern \"C\" int answer() { return 2; }"), 0600); err != nil {
		t.Fatal(err)
	}
	store.InvalidateObservations()
	after, err := snapshotNodeInputs(store, node)
	if err != nil {
		t.Fatal(err)
	}
	changed := len(before) != len(after)
	for i := range before {
		if changed {
			break
		}
		if before[i] != after[i] {
			changed = true
		}
	}
	if !changed {
		t.Fatal("CXX source edit must invalidate compile inputs")
	}
}
