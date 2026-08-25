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
	originalDirectory, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, os.Chdir(originalDirectory)) })

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
