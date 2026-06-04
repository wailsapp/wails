package wake

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestLocalOverrideEndToEnd drives the public Execute entrypoint (the same one
// the CLI's runWakeTask calls) to prove a Taskfile.local.yml override takes
// effect through the full discover -> merge -> execute path.
func TestLocalOverrideEndToEnd(t *testing.T) {
	t.Setenv("WAILS_USE_WAKE", "true")

	dir := t.TempDir()
	// Mirror the CLI: runWakeTask runs with the project dir as the process cwd,
	// and tasks without an explicit `dir:` execute there.
	t.Chdir(dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "Taskfile.yml"), []byte(`
version: "3"
tasks:
  greet:
    cmds:
      - echo base > greet.out
  base-only:
    cmds:
      - echo base-only > base.out
`), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "Taskfile.local.yml"), []byte(`
version: "3"
tasks:
  greet:
    cmds:
      - echo overridden > greet.out
  local-only:
    cmds:
      - echo local-only > local.out
`), 0644))

	// Overridden task: local cmds replace base cmds.
	require.NoError(t, Execute("greet", ExecuteOptions{Dir: dir, Silent: true}))
	got, err := os.ReadFile(filepath.Join(dir, "greet.out"))
	require.NoError(t, err)
	require.Equal(t, "overridden\n", string(got))

	// Local-only task is reachable through the merged taskfile.
	require.NoError(t, Execute("local-only", ExecuteOptions{Dir: dir, Silent: true}))
	got, err = os.ReadFile(filepath.Join(dir, "local.out"))
	require.NoError(t, err)
	require.Equal(t, "local-only\n", string(got))

	// Untouched base-only task still runs.
	require.NoError(t, Execute("base-only", ExecuteOptions{Dir: dir, Silent: true}))
	got, err = os.ReadFile(filepath.Join(dir, "base.out"))
	require.NoError(t, err)
	require.Equal(t, "base-only\n", string(got))
}

// TestOverridesDisabledOptOut verifies WAILS_NO_OVERRIDES=true skips local
// override loading entirely, so the committed base Taskfile runs unmodified.
func TestOverridesDisabledOptOut(t *testing.T) {
	t.Setenv("WAILS_USE_WAKE", "true")
	t.Setenv("WAILS_NO_OVERRIDES", "true")

	dir := t.TempDir()
	t.Chdir(dir)

	require.NoError(t, os.WriteFile(filepath.Join(dir, "Taskfile.yml"), []byte(`
version: "3"
tasks:
  greet:
    cmds:
      - echo base > greet.out
`), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "Taskfile.local.yml"), []byte(`
version: "3"
tasks:
  greet:
    cmds:
      - echo overridden > greet.out
  local-only:
    cmds:
      - echo local-only > local.out
`), 0644))

	// Base command runs; the override is ignored.
	require.NoError(t, Execute("greet", ExecuteOptions{Dir: dir, Silent: true}))
	got, err := os.ReadFile(filepath.Join(dir, "greet.out"))
	require.NoError(t, err)
	require.Equal(t, "base\n", string(got))

	// Local-only task is not present, so executing it fails.
	require.Error(t, Execute("local-only", ExecuteOptions{Dir: dir, Silent: true}))
}
