package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NotNil(t, doc.Wake.Migration)
	assert.False(t, doc.Wake.Migration.Complete)
}

func TestAnalyseMigrationRecognizesShippedDefault(t *testing.T) {
	report, _, err := analyseMigration(filepath.Join("..", "..", "examples", "badge"))
	require.NoError(t, err)
	assert.Truef(t, report.Complete, "diagnostics: %#v", report.Diagnostics)
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
}

func TestAnalyseMigrationFindsLocalCustomTasks(t *testing.T) {
	root := t.TempDir()
	writeMigrationFixture(t, root)
	require.NoError(t, os.WriteFile(filepath.Join(root, "Taskfile.local.yml"), []byte("version: '3'\ntasks:\n  publish:\n    cmds:\n      - echo custom\n"), 0o644))
	report, doc, err := analyseMigration(root)
	require.NoError(t, err)
	assert.False(t, report.Complete)
	assert.False(t, doc.Wake.Migration.Complete)
	assert.Contains(t, report.Sources, "Taskfile.local.yml")
	require.NotEmpty(t, report.Diagnostics)
	assert.Equal(t, "unsupported-task", report.Diagnostics[0].Code)
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
