package commands

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/templates"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestInitBuiltInTemplatesCreateOnlyWailsHCL(t *testing.T) {
	originalDirectory, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, os.Chdir(originalDirectory)) })
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	templateNames := make([]string, 0, len(templates.GetDefaultTemplates()))
	for _, template := range templates.GetDefaultTemplates() {
		templateNames = append(templateNames, template.Name)
	}
	sort.Strings(templateNames)
	require.NotEmpty(t, templateNames)

	for _, templateName := range templateNames {
		t.Run(templateName, func(t *testing.T) {
			require.NoError(t, os.Chdir(originalDirectory))
			options := &flags.Init{
				TemplateName:      templateName,
				ProjectName:       "HCL Only",
				ProjectDir:        filepath.Join(t.TempDir(), "project"),
				ModulePath:        "example.com/hcl-only",
				ProductName:       "HCL Only",
				ProductIdentifier: "com.example.hclonly",
				ProductVersion:    "0.1.0",
				SkipGoModTidy:     true,
			}
			require.NoError(t, Init(options))
			root := options.ProjectDir

			assert.FileExists(t, filepath.Join(root, manifest.Filename))
			assert.NoFileExists(t, filepath.Join(root, "wails.toml"))
			assert.NoFileExists(t, filepath.Join(root, "Taskfile.yml"))
			assert.NoFileExists(t, filepath.Join(root, "Taskfile.yaml"))

			loaded, err := manifest.Load(root, "")
			require.NoError(t, err)
			assert.Equal(t, "HCL_Only", loaded.Config.Project.Name)
			assert.Equal(t, "bin", loaded.Config.Build.OutputDirectory)
			assert.Equal(t, "frontend/dist", loaded.Config.Frontend.OutputDirectory)
			_, err = pipeline.PlanBuild(loaded.Config, pipeline.Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
			require.NoError(t, err)

			nested := filepath.Join(root, "frontend", "src")
			require.NoError(t, os.MkdirAll(nested, 0o755))
			nestedLoad, err := manifest.Load(nested, "")
			require.NoError(t, err)
			assert.Equal(t, root, nestedLoad.Config.Root)
		})
	}
}

func TestInitInsideAnotherManifestProjectCreatesItsOwnManifest(t *testing.T) {
	originalDirectory, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, os.Chdir(originalDirectory)) })
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	parent := t.TempDir()
	require.NoError(t, manifest.WriteMinimal(parent, manifest.Project{
		Name: "parent", ProductName: "Parent", Identifier: "com.example.parent", Version: "1.0.0",
	}))
	options := &flags.Init{
		TemplateName:      "vanilla",
		ProjectName:       "Child",
		ProjectDir:        filepath.Join(parent, "projects"),
		ModulePath:        "example.com/child",
		ProductName:       "Child",
		ProductIdentifier: "com.example.child",
		ProductVersion:    "0.1.0",
		SkipGoModTidy:     true,
	}
	require.NoError(t, Init(options))
	childManifest := filepath.Join(options.ProjectDir, manifest.Filename)
	assert.FileExists(t, childManifest)
	loaded, err := manifest.Load(options.ProjectDir, "")
	require.NoError(t, err)
	assert.Equal(t, "Child", loaded.Config.Project.Name)
}

func TestCommunityTemplateWithFullyMigratableTaskfileActivatesHCLAndPreservesLegacySource(t *testing.T) {
	root := t.TempDir()
	legacy := []byte("version: '3'\ntasks: {}\n")
	require.NoError(t, os.WriteFile(filepath.Join(root, "Taskfile.yml"), legacy, 0o600))
	options := &flags.Init{ProjectDir: root, ProjectName: "Community", ProductName: "Community App", ProductIdentifier: "com.example.community", ProductVersion: "2.0.0", Quiet: true}

	require.NoError(t, initialiseTemplateBuildManifest(options))
	assert.FileExists(t, filepath.Join(root, manifest.Filename))
	assert.NoFileExists(t, filepath.Join(root, manifest.MigratedFilename))
	actual, err := os.ReadFile(filepath.Join(root, "Taskfile.yml"))
	require.NoError(t, err)
	assert.Equal(t, legacy, actual)
	info, err := os.Stat(filepath.Join(root, "Taskfile.yml"))
	require.NoError(t, err)
	if runtime.GOOS != "windows" {
		assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
	}
	assert.NoFileExists(t, filepath.Join(root, ".wails", "migration-report.json"))
}

func TestCommunityTemplateWithCustomTaskWritesOnlyInactiveDraft(t *testing.T) {
	root := t.TempDir()
	legacy := []byte("version: '3'\ntasks:\n  build:\n    cmds:\n      - task: bespoke\n  bespoke:\n    cmds: ['echo custom']\n")
	require.NoError(t, os.WriteFile(filepath.Join(root, "Taskfile.yml"), legacy, 0o640))
	options := &flags.Init{ProjectDir: root, ProjectName: "Community", ProductName: "Community App", ProductIdentifier: "com.example.community", ProductVersion: "2.0.0", Quiet: true}

	require.NoError(t, initialiseTemplateBuildManifest(options))
	assert.NoFileExists(t, filepath.Join(root, manifest.Filename))
	assert.FileExists(t, filepath.Join(root, manifest.MigratedFilename))
	draft, err := os.ReadFile(filepath.Join(root, manifest.MigratedFilename))
	require.NoError(t, err)
	assert.Contains(t, string(draft), "# BLOCKED: Taskfile.yml [bespoke]")
	actual, err := os.ReadFile(filepath.Join(root, "Taskfile.yml"))
	require.NoError(t, err)
	assert.Equal(t, legacy, actual)
	info, err := os.Stat(filepath.Join(root, "Taskfile.yml"))
	require.NoError(t, err)
	if runtime.GOOS != "windows" {
		assert.Equal(t, os.FileMode(0o640), info.Mode().Perm())
	}
	assert.NoFileExists(t, filepath.Join(root, ".wails", "migration-report.json"))
}
