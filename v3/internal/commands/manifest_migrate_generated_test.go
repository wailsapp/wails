package commands

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/templates"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestGeneratedLegacyTemplatesMigrateWithoutBuilding(t *testing.T) {
	for _, template := range templates.GetDefaultTemplates() {
		t.Run(template.Name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			t.Setenv("XDG_CONFIG_HOME", t.TempDir())
			t.Setenv("WAILS_EXP_USE_WAKE", "restore")
			if err := os.Unsetenv("WAILS_EXP_USE_WAKE"); err != nil {
				t.Fatal(err)
			}
			options := &flags.Init{
				TemplateName: template.Name, ProjectName: "migration-demo", ProjectDir: ".",
				ModulePath: "example.com/migration-demo", ProductName: "Migration Demo",
				ProductIdentifier: "com.example.migrationdemo", ProductVersion: "1.2.3",
				SkipGoModTidy: true,
			}
			captureMigrationStdout(t, func() {
				if err := Init(options); err != nil {
					t.Fatal(err)
				}
			})
			root := options.ProjectDir
			legacyData, err := os.ReadFile(filepath.Join(root, "build", "config.yml"))
			if err != nil {
				t.Fatal(err)
			}
			var legacy legacyBuildConfiguration
			if err := yaml.Unmarshal(legacyData, &legacy); err != nil {
				t.Fatal(err)
			}
			before, err := snapshotMigrationInputs(root)
			if err != nil {
				t.Fatal(err)
			}
			for _, name := range []string{manifest.Filename, manifest.MigratedFilename, "frontend/dist", "frontend/node_modules"} {
				if _, err := os.Stat(filepath.Join(root, name)); !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("expected an unbuilt legacy project; %s: %v", name, err)
				}
			}
			t.Setenv("WAILS_EXP_USE_WAKE", "1")
			t.Chdir(root)
			report, _, err := analyseMigration(root)
			if err != nil {
				t.Fatal(err)
			}
			if !report.Complete {
				t.Fatalf("current generated template is not migratable: %+v", report.Diagnostics)
			}
			run := func(options MigrateOptions) {
				t.Helper()
				captureMigrationStdout(t, func() {
					if err := Migrate(&options); err != nil {
						t.Fatal(err)
					}
				})
			}
			run(MigrateOptions{DryRun: true})
			if _, err := os.Stat(manifest.MigratedFilename); !errors.Is(err, os.ErrNotExist) {
				t.Fatal("dry-run created a draft", err)
			}
			run(MigrateOptions{})
			draft, err := os.ReadFile(manifest.MigratedFilename)
			if err != nil {
				t.Fatal(err)
			}
			reviewed := append([]byte("# reviewed migration\n"), draft...)
			if err := os.WriteFile(manifest.MigratedFilename, reviewed, 0600); err != nil {
				t.Fatal(err)
			}
			run(MigrateOptions{})
			run(MigrateOptions{Activate: true, DryRun: true})
			if _, err := os.Stat(manifest.Filename); !errors.Is(err, os.ErrNotExist) {
				t.Fatal("activation dry-run created the active manifest", err)
			}
			run(MigrateOptions{Activate: true})
			active, err := os.ReadFile(manifest.Filename)
			if err != nil || !bytes.Equal(active, reviewed) {
				t.Fatal("activation did not preserve the reviewed draft", err)
			}
			if _, err := os.Stat(manifest.MigratedFilename); !errors.Is(err, os.ErrNotExist) {
				t.Fatal("activation retained the inactive draft", err)
			}
			after, err := snapshotMigrationInputs(root)
			if err != nil {
				t.Fatal(err)
			}
			delete(after, manifest.Filename)
			if !reflect.DeepEqual(before, after) {
				t.Fatal("migration changed legacy sources or generated build outputs")
			}
			loaded, err := manifest.Load(root, "")
			if err != nil {
				t.Fatal(err)
			}
			if loaded.Config.Project.ProductName != legacy.Info.ProductName || loaded.Config.Project.Version != legacy.Info.Version {
				t.Fatal("migration lost generated project metadata")
			}
			if _, err := pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"}); err != nil {
				t.Fatal("migrated project cannot produce a build plan", err)
			}
		})
	}
}
