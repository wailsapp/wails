package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestConfigCheckValidatesEveryProfileAndOneRequestedProfile(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, manifest.WriteMinimal(root, manifest.Project{Name: "check", ProductName: "Check", Identifier: "com.example.check", Version: "1.0.0"}))
	path := filepath.Join(root, manifest.Filename)
	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	raw = append(raw, []byte("\nprofile \"release\" {\n  target \"linux/amd64\" {}\n}\n")...)
	require.NoError(t, os.WriteFile(path, raw, 0o644))
	t.Chdir(root)

	require.NoError(t, ConfigCheck(&ConfigCheckOptions{}, nil))
	require.NoError(t, ConfigCheck(&ConfigCheckOptions{Profile: "release"}, nil))
	err = ConfigCheck(&ConfigCheckOptions{Profile: "release"}, []string{"other"})
	assert.ErrorContains(t, err, "not both")
}

func TestConfigCheckReportsInvalidProfilePlan(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, manifest.WriteMinimal(root, manifest.Project{Name: "check", ProductName: "Check", Identifier: "com.example.check", Version: "1.0.0"}))
	path := filepath.Join(root, manifest.Filename)
	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	raw = append(raw, []byte("\nprofile \"release\" {\n  target \"linux/amd64\" { formats = [\"aab\"] }\n}\n")...)
	require.NoError(t, os.WriteFile(path, raw, 0o644))
	t.Chdir(root)

	err = ConfigCheck(&ConfigCheckOptions{}, nil)
	assert.ErrorContains(t, err, `profile "release"`)
	assert.ErrorContains(t, err, `format "aab"`)
}
