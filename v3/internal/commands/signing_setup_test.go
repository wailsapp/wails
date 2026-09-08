package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestSigningSetupUpdatesHCLWithoutTouchingTaskfiles(t *testing.T) {
	t.Setenv("WAILS_EXP_USE_WAKE", "1")
	originalDirectory, err := os.Getwd()
	require.NoError(t, err)
	defer func() { require.NoError(t, os.Chdir(originalDirectory)) }()

	root := t.TempDir()
	require.NoError(t, manifest.WriteMinimal(root, manifest.Project{
		Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0",
	}))
	taskfile := filepath.Join(root, "build", "linux", "Taskfile.yml")
	require.NoError(t, os.MkdirAll(filepath.Dir(taskfile), 0o755))
	require.NoError(t, os.WriteFile(taskfile, []byte("user-owned taskfile\n"), 0o640))
	nested := filepath.Join(root, "frontend", "src")
	require.NoError(t, os.MkdirAll(nested, 0o755))
	require.NoError(t, os.Chdir(nested))

	previous := runLinuxSigningSetup
	runLinuxSigningSetup = func(save signingSetupSave) error {
		return save("linux", manifest.SigningPlatform{
			Enabled:     true,
			Certificate: "signing-key.asc",
			Identity:    "origin",
		}, map[string]string{"PGP_KEY": "must-not-be-written"})
	}
	t.Cleanup(func() { runLinuxSigningSetup = previous })

	require.NoError(t, SigningSetup(&flags.SigningSetup{Platforms: []string{"linux"}}))
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	assert.Equal(t, "signing-key.asc", loaded.Config.Signing.Linux.Certificate)
	assert.Equal(t, "origin", loaded.Config.Signing.Linux.Identity)
	contents, err := os.ReadFile(taskfile)
	require.NoError(t, err)
	assert.Equal(t, "user-owned taskfile\n", string(contents))
}

func TestSigningSetupAcceptsCommaSeparatedPlatforms(t *testing.T) {
	assert.Equal(t, []string{"darwin", "windows", "linux"}, normaliseSigningPlatforms([]string{"darwin, windows", "linux"}))
}

func TestSigningSetupWithoutExperimentUsesTaskfiles(t *testing.T) {
	t.Setenv("WAILS_EXP_USE_WAKE", "restore")
	require.NoError(t, os.Unsetenv("WAILS_EXP_USE_WAKE"))
	t.Chdir(t.TempDir())
	require.NoError(t, os.MkdirAll("build/linux", 0755))
	require.NoError(t, os.WriteFile("build/linux/Taskfile.yml", []byte("version: '3'\nvars:\n  PGP_KEY: old.asc\ntasks:\n  custom:\n    cmds: ['echo keep']\n"), 0644))
	require.NoError(t, os.WriteFile("wails.hcl", []byte("do not read or modify"), 0644))
	previous := runLinuxSigningSetup
	t.Cleanup(func() { runLinuxSigningSetup = previous })
	runLinuxSigningSetup = func(save signingSetupSave) error {
		return save("linux", manifest.SigningPlatform{Certificate: "new.asc"}, map[string]string{"PGP_KEY": "new.asc"})
	}
	require.NoError(t, SigningSetup(&flags.SigningSetup{Platforms: []string{"linux"}}))
	data, err := os.ReadFile("build/linux/Taskfile.yml")
	require.NoError(t, err)
	require.Contains(t, string(data), "new.asc")
	require.Contains(t, string(data), "echo keep")
	data, err = os.ReadFile("wails.hcl")
	require.NoError(t, err)
	require.Equal(t, "do not read or modify", string(data))
}
