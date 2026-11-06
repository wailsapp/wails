package pipeline

import (
	"os"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestAppImageInspectionWaitsForRestoredExecutable(t *testing.T) {
	handler := &faultExecutorHandler{identityErr: os.ErrNotExist}
	executor := Executor{Handler: handler}
	node := Node{Kind: PackageArtifact, Spec: PackageSpec{Format: "appimage"}, Cache: CacheArtifact, Output: "bin/app.AppImage", Dependencies: []NodeKey{"compile"}}
	operations := map[NodeKey]OperationInspection{"compile": {Decision: "restore"}}
	inspection, _, err := executor.inspectNode(t.Context(), &faultExecutorCache{}, node, operations, nil, fixedActionKey)
	if err != nil || inspection.Decision != "run" {
		t.Fatal("inspection read an executable before restoring it", inspection, err)
	}
	operations["compile"] = OperationInspection{Decision: "cached"}
	if _, _, err := executor.inspectNode(t.Context(), &faultExecutorCache{}, node, operations, nil, fixedActionKey); err == nil {
		t.Fatal("inspection suppressed an identity error after dependencies were present")
	}
}

func TestAppImageGTKProfileUsesResolvedExecutableTags(t *testing.T) {
	for _, test := range []struct {
		name               string
		primary, companion CompileSpec
		configured, want   int
		conflict           bool
	}{
		{name: "headless", primary: CompileSpec{Console: true}, companion: CompileSpec{Console: true}},
		{name: "default GUI", companion: CompileSpec{Console: true}, want: 4},
		{name: "purego GTK3 companion", primary: CompileSpec{Console: true}, companion: CompileSpec{Tags: []string{"purego", "gtk3"}}, want: 3},
		{name: "mixed profiles", companion: CompileSpec{Tags: []string{"gtk3"}}, conflict: true},
		{name: "prebuilt override", configured: 3, want: 3},
	} {
		t.Run(test.name, func(t *testing.T) {
			plan := Plan{Nodes: map[NodeKey]Node{"primary": {Spec: test.primary}, "companion": {Spec: test.companion}}}
			got, err := appImageGTKVersion(plan, test.configured)
			if test.conflict {
				if err == nil || !strings.Contains(err.Error(), "same GTK version") {
					t.Fatal("mixed GTK profiles accepted", err)
				}
			} else if err != nil || got != test.want {
				t.Fatalf("GTK = %d, %v; want %d", got, err, test.want)
			}
		})
	}
}

func TestAppImagePlanningSupportsDesktopHostsWithoutDeploymentTools(t *testing.T) {
	for _, host := range []string{"linux", "darwin", "windows"} {
		config := testConfig(t)
		config.Frontend.Disabled = true
		config.Build.Go.Tags = []string{"purego", "gtk3"}
		config.Package.Linux.AppImage = manifest.AppImagePackage{RuntimeRoot: "target-root"}
		plan, err := PlanBuildForHost(config, Request{Verb: "package", TargetOS: "linux", TargetArch: "amd64", Formats: []string{"appimage"}}, testHost(host, "arm64", "go"))
		if err != nil {
			t.Fatalf("%s: %v", host, err)
		}
		spec := plan.Nodes["package:linux/amd64:appimage"].Spec.(PackageSpec)
		if spec.GTKVersion != 3 {
			t.Fatalf("%s lost GTK3 tags: %+v", host, spec)
		}
	}
}
