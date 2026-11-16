package commands

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/setupwizard"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestInitWizardChoicesConfigureTheNewProject(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		t.Run(map[bool]string{false: "legacy", true: "wake"}[enabled], func(t *testing.T) {
			t.Chdir(t.TempDir())
			t.Setenv("XDG_CONFIG_HOME", t.TempDir())
			t.Setenv("WAILS_EXP_USE_WAKE", "")
			if !enabled {
				if err := os.Unsetenv("WAILS_EXP_USE_WAKE"); err != nil {
					t.Fatal(err)
				}
			}
			options := &flags.Init{UI: true, ProjectName: "created", ProjectDir: ".", TemplateName: "vanilla-js", ModulePath: "example.com/created", ProductName: "Created App", ProductIdentifier: "org.example.created", ProductVersion: "1.2.3", SkipGoModTidy: true}
			wizard := func(*flags.Init) (*setupwizard.WakeSetupConfig, error) {
				if !enabled {
					return nil, nil
				}
				return &setupwizard.WakeSetupConfig{ApplicationType: "gui", Backend: "purego", ExperimentalPureGo: true, Targets: []setupwizard.WakeSetupTarget{{OS: "windows", Arch: "arm64", Formats: []string{"msix"}}}}, nil
			}
			if err := initWithWizard(options, wizard); err != nil {
				t.Fatal(err)
			}
			for _, path := range []string{"main.go", "frontend/package.json", "frontend/src/main.js"} {
				if _, err := os.Stat(filepath.Join(options.ProjectDir, path)); err != nil {
					t.Fatal("missing scaffold", path, err)
				}
			}
			if !enabled {
				if _, err := os.Stat(filepath.Join(options.ProjectDir, "Taskfile.yml")); err != nil {
					t.Fatal("legacy entry missing", err)
				}
				if _, err := os.Stat(filepath.Join(options.ProjectDir, "wails.yaml")); !os.IsNotExist(err) {
					t.Fatal("legacy init wrote Wake manifest")
				}
				return
			}
			loaded, err := manifest.Load(options.ProjectDir, "setup")
			if err != nil {
				t.Fatal(err)
			}
			if loaded.Config.Project.ProductName != "Created App" || loaded.Config.Frontend.Bindings.TypeScript || loaded.Config.Selected.Targets[0].Target != "windows/arm64" || loaded.Config.Selected.Targets[0].Formats[0] != "msix" {
				t.Fatalf("init choices or template metadata lost: %+v", loaded.Config)
			}
			plan, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build"})
			if err != nil {
				t.Fatal(err)
			}
			compiled := false
			for _, node := range plan.Nodes {
				if spec, ok := node.Spec.(pipeline.CompileSpec); ok {
					compiled = true
					if spec.TargetOS != "windows" || spec.TargetArch != "arm64" || spec.CGOEnabled != "0" || !slices.Contains(spec.Tags, "purego") {
						t.Fatalf("choices did not control compilation: %+v", spec)
					}
				}
			}
			if !compiled {
				t.Fatal("no compiled executable")
			}
		})
	}
}
