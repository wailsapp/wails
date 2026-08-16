package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	wakeast "github.com/wailsapp/wails/v3/internal/wake/ast"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/migration"
)

func TestAnalyseMigrationTranslatesLegacyConfig(t *testing.T) {
	root := t.TempDir()
	writeMigrationFixture(t, root)
	report, doc, err := analyseMigration(root)
	require.NoError(t, err)
	assert.False(t, report.Complete)
	assert.Equal(t, "Example", doc.Project.ProductName)
	assert.Equal(t, "com.example.app", doc.Project.Identifier)
	assert.Equal(t, "pnpm", doc.Frontend.PackageManager)
	require.Len(t, doc.Associations, 1)
	assert.Equal(t, []string{"demo"}, doc.Associations[0].Extensions)
	assert.NotEmpty(t, report.CompletedBy)
	assert.NotEmpty(t, report.Sources)
	assert.Empty(t, doc.Wake, "migration state belongs in .wails/migration-report.json, not wails.toml")
}

func TestMigrationRecognizesHistoricalGeneratedDefaults(t *testing.T) {
	assert.True(t, knownStockTaskfiles["common"]["4ab65b0363866b550ef897db8f7aae12794789056dd4f2e82407a738a64f6819"])
	assert.True(t, legacyDevCommand("wails3 build DEV=true"))
	assert.False(t, legacyDevCommand("wails3 build DEV=true && publish"))
}

func TestAnalyseMigrationRecognizesShippedDefault(t *testing.T) {
	report, doc, err := analyseMigration(filepath.Join("..", "..", "examples", "badge"))
	require.NoError(t, err)
	assert.Truef(t, report.Complete, "diagnostics: %#v", report.Diagnostics)
	assert.Empty(t, doc.Wake, "completed migrations must look like native manifests")
	for _, taskfile := range report.Taskfiles {
		assert.Contains(t, []string{"current-default", "historical-default"}, taskfile.Classification, "%s", taskfile.File)
	}
}

func TestAnalyseMigrationFindsModifiedKnownRootTask(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.CopyFS(root, os.DirFS(filepath.Join("..", "..", "examples", "badge"))))
	path := filepath.Join(root, "Taskfile.yml")
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	original := string(data)
	modified := strings.Replace(original, `task: "{{OS}}:build"`, `echo custom build`, 1)
	require.NotEqual(t, original, modified, "fixture Taskfile.yml no longer contains the expected root dispatch command")
	data = []byte(modified)
	require.NoError(t, os.WriteFile(path, data, 0o644))
	report, _, err := analyseMigration(root)
	require.NoError(t, err)
	assert.False(t, report.Complete)
	assert.Contains(t, report.Diagnostics, MigrationDiagnostic{Severity: "warning", Code: "modified-task", File: "Taskfile.yml", Task: "build", Message: "root dispatch task was modified and requires a manifest field or user-owned hook script"})
	for _, taskfile := range report.Taskfiles {
		if taskfile.File == "Taskfile.yml" {
			assert.Equal(t, "customised", taskfile.Classification)
			assert.Contains(t, taskfile.ModifiedTasks, "build")
		}
	}
}

func TestCurrentCommonTaskfileDefaultsAreComparedStructurally(t *testing.T) {
	variants := stockTaskVariants("common")
	require.Len(t, variants, 4)
	assert.True(t, closestCanonical(variants[0], variants).exact)

	actual := make(map[string]*wakeast.Task, len(variants[0]))
	for name, task := range variants[0] {
		actual[name] = task.Clone()
	}
	require.NotEmpty(t, actual["build:frontend"].Cmds)
	actual["build:frontend"].Cmds[0].Cmd = "echo customised"
	diff := closestCanonical(actual, variants)
	assert.False(t, diff.exact)
	assert.Contains(t, diff.changed, "build:frontend")
}

func TestAnalyseMigrationFindsLocalCustomTasks(t *testing.T) {
	root := t.TempDir()
	writeMigrationFixture(t, root)
	require.NoError(t, os.WriteFile(filepath.Join(root, "Taskfile.local.yml"), []byte("version: '3'\ntasks:\n  publish:\n    cmds:\n      - echo custom\n"), 0o644))
	report, doc, err := analyseMigration(root)
	require.NoError(t, err)
	assert.False(t, report.Complete)
	assert.Empty(t, doc.Wake)
	assert.Contains(t, report.Sources, "Taskfile.local.yml")
	require.NotEmpty(t, report.Diagnostics)
	assert.Equal(t, "unsupported-task", report.Diagnostics[0].Code)
}

func TestBackupLegacySourcesPreservesTree(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "build", "Taskfile.yml")
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, []byte("version: '3'\n"), 0o644))
	digest, err := digestFile(path)
	require.NoError(t, err)
	backedUp, diagnostics := backupLegacySources(root, map[string]string{"build/Taskfile.yml": digest})
	assert.Empty(t, diagnostics)
	assert.Equal(t, []string{".wails/migration-backup/build/Taskfile.yml"}, backedUp)
	data, err := os.ReadFile(filepath.Join(root, ".wails", "migration-backup", "build", "Taskfile.yml"))
	require.NoError(t, err)
	assert.Equal(t, "version: '3'\n", string(data))
}

func TestMigrateCompleteProjectRetiresLegacySourcesAndKeepsStatePrivate(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.CopyFS(root, os.DirFS(filepath.Join("..", "..", "examples", "badge"))))
	require.NoError(t, os.Remove(filepath.Join(root, manifest.Filename)))
	t.Chdir(root)

	require.NoError(t, Migrate(&MigrateOptions{Backup: true}))
	_, err := os.Stat(filepath.Join(root, "Taskfile.yml"))
	assert.ErrorIs(t, err, os.ErrNotExist)
	_, err = os.Stat(filepath.Join(root, ".wails", "migration-backup", "Taskfile.yml"))
	require.NoError(t, err)
	report, exists, err := migration.Read(root)
	require.NoError(t, err)
	assert.True(t, exists)
	assert.True(t, report.Complete)
	assert.Contains(t, report.Removed, "Taskfile.yml")

	manifestData, err := os.ReadFile(filepath.Join(root, manifest.Filename))
	require.NoError(t, err)
	assert.NotContains(t, string(manifestData), "migration")
}

func TestCompleteReviewedMigrationUsesOriginalSourceDigests(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "Taskfile.yml"), []byte("version: '3'\ntasks:\n  publish:\n    cmds: [echo custom]\n"), 0o644))
	t.Chdir(root)
	require.NoError(t, Migrate(&MigrateOptions{}))

	require.NoError(t, Migrate(&MigrateOptions{Complete: true, Backup: true}))
	_, err := os.Stat(filepath.Join(root, "Taskfile.yml"))
	assert.ErrorIs(t, err, os.ErrNotExist)
	report, exists, err := migration.Read(root)
	require.NoError(t, err)
	assert.True(t, exists)
	assert.True(t, report.Complete)
	assert.Contains(t, report.BackedUp, ".wails/migration-backup/Taskfile.yml")
}

func TestCompleteReviewedMigrationRefusesChangedSource(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "Taskfile.yml")
	require.NoError(t, os.WriteFile(path, []byte("version: '3'\ntasks:\n  publish:\n    cmds: [echo custom]\n"), 0o644))
	t.Chdir(root)
	require.NoError(t, Migrate(&MigrateOptions{}))
	require.NoError(t, os.WriteFile(path, []byte("version: '3'\ntasks:\n  publish:\n    cmds: [echo changed]\n"), 0o644))

	err := Migrate(&MigrateOptions{Complete: true})
	require.ErrorContains(t, err, "legacy sources changed after analysis")
	_, statErr := os.Stat(path)
	require.NoError(t, statErr)
	report, exists, readErr := migration.Read(root)
	require.NoError(t, readErr)
	assert.True(t, exists)
	assert.False(t, report.Complete)
}

func TestAnalyseMigrationTranslatesScriptFileHook(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "scripts"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "scripts", "before.sh"), []byte("#!/bin/sh\n"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "Taskfile.yml"), []byte(`version: '3'
tasks:
  before-build:
    cmds:
      - ./scripts/before.sh
`), 0o644))

	report, doc, err := analyseMigration(root)
	require.NoError(t, err)
	assert.True(t, report.Complete)
	assert.Equal(t, "scripts/before.sh", doc.Hooks.BeforeBuild.Script)
	assert.Contains(t, report.Diagnostics, MigrationDiagnostic{Severity: "info", Code: "script-hook", File: "Taskfile.yml", Task: "before-build", Message: "translated script-file task to hooks.before_build"})
}

func TestRemoveLegacySourcesRequiresMatchingDigest(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "Taskfile.yml")
	require.NoError(t, os.WriteFile(path, []byte("original"), 0o644))
	digest, err := digestFile(path)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, []byte("changed"), 0o644))
	removed, diagnostics := removeLegacySources(root, map[string]string{"Taskfile.yml": digest}, true)
	assert.Empty(t, removed)
	require.Len(t, diagnostics, 1)
	assert.Equal(t, "modified-source", diagnostics[0].Code)
	_, err = os.Stat(path)
	assert.NoError(t, err)
}

func TestRemoveLegacySourcesRejectsEscapedPath(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "project")
	require.NoError(t, os.MkdirAll(root, 0o755))
	outside := filepath.Join(parent, "outside-taskfile.yml")
	require.NoError(t, os.WriteFile(outside, []byte("version: '3'\n"), 0o644))
	digest, err := digestFile(outside)
	require.NoError(t, err)
	removed, diagnostics := removeLegacySources(root, map[string]string{"../outside-taskfile.yml": digest}, true)
	assert.Empty(t, removed)
	require.Len(t, diagnostics, 1)
	assert.Equal(t, "remove-outside-project", diagnostics[0].Code)
	_, err = os.Stat(outside)
	assert.NoError(t, err)
}

func writeMigrationFixture(t *testing.T, root string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "build"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "Taskfile.yml"), []byte(`version: '3'
vars:
  APP_NAME: example
  PACKAGE_MANAGER: pnpm
includes:
  common: ./build/Taskfile.yml
tasks:
  build:
    cmds:
      - task: common:build:frontend
`), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "build", "Taskfile.yml"), []byte(`version: '3'
tasks:
  build:frontend:
    cmds:
      - pnpm run build
`), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "build", "config.yml"), []byte(`info:
  productName: Example
  productIdentifier: com.example.app
  version: 1.2.3
dev_mode:
  executes:
    - cmd: wails3 task build
      type: blocking
fileAssociations:
  - ext: demo
    name: Demo
`), 0o644))
}
