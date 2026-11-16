package pipeline

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/internal/report"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestDevelopmentBuildProducesCompanionsWithoutHooks(t *testing.T) {
	for _, platform := range []string{"linux", "windows", "darwin"} {
		t.Run(platform, func(t *testing.T) {
			config := performanceConfig(t)
			config.Frontend.Disabled = true
			config.Binaries = []manifest.Binary{{Name: "helper", Package: "cmd/helper", Kind: "cli"}}
			if err := os.MkdirAll(filepath.Join(config.Root, "cmd", "helper"), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(config.Root, "cmd", "helper", "main.go"), []byte("package main\nfunc main() {}\n"), 0644); err != nil {
				t.Fatal(err)
			}
			plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: platform, TargetArch: "amd64", Development: true})
			if err != nil {
				t.Fatal(err)
			}
			results, err := (Executor{Handler: &fakeHandler{root: config.Root}}).Execute(context.Background(), plan, ExecuteOptions{Root: config.Root, Reporter: report.Nop{}})
			if err != nil {
				t.Fatal(err)
			}
			found := false
			for key, node := range plan.Nodes {
				if spec, ok := node.Spec.(CompileSpec); ok && spec.Package == "cmd/helper" {
					found = true
					if _, ok := results[key]; !ok {
						t.Fatal("development build omitted companion compilation")
					}
					if _, err := os.Stat(filepath.Join(config.Root, spec.Output)); err != nil {
						t.Fatalf("companion output: %v", err)
					}
				}
			}
			if !found {
				t.Fatal("development build has no companion program")
			}
		})
	}
}
