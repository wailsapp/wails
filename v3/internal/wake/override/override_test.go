package override

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadLocalReturnsEmptyWhenAbsent(t *testing.T) {
	dir := t.TempDir()
	locals, err := LoadLocal(dir)
	require.NoError(t, err)
	require.Empty(t, locals)
}

func TestLoadLocalParsesOverride(t *testing.T) {
	dir := t.TempDir()
	ovPath := filepath.Join(dir, "Taskfile.local.yml")
	err := os.WriteFile(ovPath, []byte(`
version: "3"
tasks:
  build:
    cmds:
      - echo local-build
`), 0644)
	require.NoError(t, err)

	locals, err := LoadLocal(dir)
	require.NoError(t, err)
	require.Len(t, locals, 1)
	require.Contains(t, locals[0].Tasks, "build")
	require.Equal(t, "echo local-build", locals[0].Tasks["build"].Cmds[0].Cmd)
}

func TestLoadLocalStacksOverrideThenLocal(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "Taskfile.override.yml"), []byte(`
version: "3"
tasks:
  build:
    cmds:
      - echo from-override
`), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "Taskfile.local.yml"), []byte(`
version: "3"
tasks:
  build:
    cmds:
      - echo from-local
`), 0644))

	locals, err := LoadLocal(dir)
	require.NoError(t, err)
	// Lowest-to-highest precedence: override first, local last.
	require.Len(t, locals, 2)
	require.Equal(t, "echo from-override", locals[0].Tasks["build"].Cmds[0].Cmd)
	require.Equal(t, "echo from-local", locals[1].Tasks["build"].Cmds[0].Cmd)
}

func TestLoadLocalOneFilePerLayer(t *testing.T) {
	dir := t.TempDir()
	// Both extensions of the same layer present: only the .yml wins.
	require.NoError(t, os.WriteFile(filepath.Join(dir, "Taskfile.local.yml"), []byte(`
version: "3"
tasks:
  build:
    cmds:
      - echo yml
`), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "Taskfile.local.yaml"), []byte(`
version: "3"
tasks:
  build:
    cmds:
      - echo yaml
`), 0644))

	locals, err := LoadLocal(dir)
	require.NoError(t, err)
	require.Len(t, locals, 1)
	require.Equal(t, "echo yml", locals[0].Tasks["build"].Cmds[0].Cmd)
}

func TestLoadLocalReturnsParseError(t *testing.T) {
	dir := t.TempDir()
	// Invalid YAML in the override file must surface as an error, not be
	// silently ignored (fail-closed).
	require.NoError(t, os.WriteFile(filepath.Join(dir, "Taskfile.local.yml"),
		[]byte("version: \"3\"\ntasks:\n  build:\n    cmds: [unterminated"), 0644))

	_, err := LoadLocal(dir)
	require.Error(t, err)
}

func TestLoadLocalReturnsIncludeError(t *testing.T) {
	dir := t.TempDir()
	// A non-optional include pointing at a missing file must error.
	require.NoError(t, os.WriteFile(filepath.Join(dir, "Taskfile.local.yml"), []byte(`
version: "3"
includes:
  missing:
    taskfile: ./does-not-exist/Taskfile.yml
tasks:
  build:
    cmds:
      - echo hi
`), 0644))

	_, err := LoadLocal(dir)
	require.Error(t, err)
}

// TestLoadLocalReadsThroughSymlink documents the security contract: a symlinked
// override file is read like any other (os.ReadFile follows symlinks). This is
// no broader a capability than the base Taskfile, which is also just a file the
// runner reads and executes.
func TestLoadLocalReadsThroughSymlink(t *testing.T) {
	dir := t.TempDir()
	external := filepath.Join(t.TempDir(), "external.yml")
	require.NoError(t, os.WriteFile(external, []byte(`
version: "3"
tasks:
  build:
    cmds:
      - echo via-symlink
`), 0644))
	require.NoError(t, os.Symlink(external, filepath.Join(dir, "Taskfile.local.yml")))

	locals, err := LoadLocal(dir)
	require.NoError(t, err)
	require.Len(t, locals, 1)
	require.Equal(t, "echo via-symlink", locals[0].Tasks["build"].Cmds[0].Cmd)
}

// TestLoadLocalIncludeResolvesRelativeToFile documents that includes in an
// override file resolve relative to that file's directory (including `../`),
// matching base-Taskfile include semantics — not a sandbox.
func TestLoadLocalIncludeResolvesRelativeToFile(t *testing.T) {
	root := t.TempDir()
	project := filepath.Join(root, "project")
	require.NoError(t, os.MkdirAll(project, 0755))

	// Sibling taskfile one level up from the project dir.
	require.NoError(t, os.WriteFile(filepath.Join(root, "shared.yml"), []byte(`
version: "3"
tasks:
  shared-task:
    cmds:
      - echo shared
`), 0644))

	require.NoError(t, os.WriteFile(filepath.Join(project, "Taskfile.local.yml"), []byte(`
version: "3"
includes:
  shared:
    taskfile: ../shared.yml
`), 0644))

	locals, err := LoadLocal(project)
	require.NoError(t, err)
	require.Len(t, locals, 1)
	require.Contains(t, locals[0].Tasks, "shared:shared-task")
}

func TestLoadLocalResolvesIncludes(t *testing.T) {
	dir := t.TempDir()

	sub := filepath.Join(dir, "extra")
	require.NoError(t, os.MkdirAll(sub, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(sub, "Taskfile.yml"), []byte(`
version: "3"
tasks:
  deploy:
    cmds:
      - echo deploy
`), 0644))

	require.NoError(t, os.WriteFile(filepath.Join(dir, "Taskfile.local.yml"), []byte(`
version: "3"
includes:
  extra:
    taskfile: ./extra/Taskfile.yml
    dir: ./extra
tasks:
  build:
    cmds:
      - echo local-build
`), 0644))

	locals, err := LoadLocal(dir)
	require.NoError(t, err)
	require.Len(t, locals, 1)
	require.Contains(t, locals[0].Tasks, "build")
	require.Contains(t, locals[0].Tasks, "extra:deploy")
}
