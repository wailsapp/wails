package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestNSISStagesRuntimeOnlyWhenRequired(t *testing.T) {
	for _, test := range []struct {
		name, kind, companionKind string
		custom, runtime           bool
	}{
		{name: "CLI", kind: "cli"},
		{name: "server", kind: "server"},
		{name: "GUI", kind: "gui", runtime: true},
		{name: "CLI companion", kind: "cli", companionKind: "cli"},
		{name: "GUI companion", kind: "cli", companionKind: "gui", runtime: true},
		{name: "unknown companion", kind: "cli", companionKind: "unknown", runtime: true},
		{name: "custom CLI script", kind: "cli", custom: true, runtime: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			assets := filepath.Join(root, "assets")
			if err := GenerateBuildAssets(&BuildAssetsOptions{Dir: assets, Name: "runtime-test", BinaryName: "runtime-test", ProductName: "Runtime Test", ProductVersion: "1.2.3", Silent: true}); err != nil {
				t.Fatal(err)
			}
			originalPath := filepath.Join(assets, "windows", "nsis", "project.nsi")
			original, err := os.ReadFile(originalPath)
			if err != nil {
				t.Fatal(err)
			}
			spec := pipeline.PackageSpec{TargetOS: "windows", TargetArch: "amd64", Format: "nsis", Assets: "assets", Profile: "default", Project: manifest.Project{Kind: test.kind}}
			if test.companionKind != "" {
				spec.Companions = []pipeline.NamedBinary{{Name: "helper.exe", Path: "bin/helper.exe", Kind: test.companionKind}}
			}
			wantScript := original
			var customPath string
			if test.custom {
				customPath = filepath.Join(root, "custom.nsi")
				wantScript = []byte("# user owned\n!insertmacro wails.webview2runtime\n")
				if err := os.WriteFile(customPath, wantScript, 0600); err != nil {
					t.Fatal(err)
				}
				spec.Config.Template = "custom.nsi"
			}
			directory, _, err := (&manifestHandler{root: root}).prepareNSISWorkspace(spec)
			if err != nil {
				t.Fatal(err)
			}
			_, runtimeErr := os.Stat(filepath.Join(directory, "MicrosoftEdgeWebview2Setup.exe"))
			if test.runtime && runtimeErr != nil || !test.runtime && !os.IsNotExist(runtimeErr) {
				t.Fatalf("runtime expected=%v stat=%v", test.runtime, runtimeErr)
			}
			staged, err := os.ReadFile(filepath.Join(directory, "project.nsi"))
			if err != nil || string(staged) != string(wantScript) {
				t.Fatalf("project script changed: %v", err)
			}
			after, err := os.ReadFile(originalPath)
			if err != nil || string(after) != string(original) {
				t.Fatalf("platform assets changed: %v", err)
			}
			if customPath != "" {
				after, err := os.ReadFile(customPath)
				if err != nil || string(after) != string(wantScript) {
					t.Fatalf("custom source changed: %v", err)
				}
			}
		})
	}
}

func TestNSISPackagePlanRetainsGUICompanionKind(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "cmd", "helper"), 0700); err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string]string{
		"go.mod":             "module example.org/runtimeplan\n\ngo 1.26.0\n",
		"main.go":            "package main\nfunc main(){}\n",
		"cmd/helper/main.go": "package main\nfunc main(){}\n",
		"wails.yaml": `version: 3
project:
  name: runtimeplan
  product_name: Runtime Plan
  identifier: org.example.runtimeplan
  kind: cli
  version: 1.0.0
build:
  environment:
    CGO_ENABLED: "0"
binaries:
  helper:
    package: cmd/helper
    kind: gui
`,
	} {
		if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(name)), []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
	}
	loaded, err := manifest.Load(root, "")
	if err != nil {
		t.Fatal(err)
	}
	plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build", TargetOS: "windows", TargetArch: "amd64", Formats: []string{"nsis"}, SkipFrontend: true, SkipBindings: true})
	if err != nil {
		t.Fatal(err)
	}
	node, ok := plan.Nodes[pipeline.NodeKey("package:windows/amd64:nsis")]
	if !ok {
		t.Fatal("missing NSIS package")
	}
	spec := node.Spec.(pipeline.PackageSpec)
	if len(spec.Companions) != 1 || spec.Companions[0].Kind != "gui" || !packageUsesWebView(spec) {
		t.Fatalf("lost GUI companion requirements: %+v", spec.Companions)
	}

	loaded.Config.Binaries[0].Kind = ""
	plan, err = pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build", TargetOS: "windows", TargetArch: "amd64", Formats: []string{"nsis"}, SkipFrontend: true, SkipBindings: true})
	if err != nil {
		t.Fatal(err)
	}
	spec = plan.Nodes[pipeline.NodeKey("package:windows/amd64:nsis")].Spec.(pipeline.PackageSpec)
	if len(spec.Companions) != 1 || spec.Companions[0].Kind != "cli" || packageUsesWebView(spec) {
		t.Fatalf("lost default CLI companion semantics: %+v", spec.Companions)
	}
}
