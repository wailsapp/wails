package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

func TestMacOSOutputAncestorMetadataRevokesMarkerAndRelocates(t *testing.T) {
	for _, attribute := range []string{
		"com.apple.fileprovider.fpfs#P",
		"com.apple.fileprovider.detached#B",
	} {
		t.Run(attribute, func(t *testing.T) {
			providerRoot := t.TempDir()
			project := filepath.Join(providerRoot, "fresh", "project")
			bin := filepath.Join(project, "bin")
			require.NoError(t, os.MkdirAll(bin, 0o755))
			artifact := filepath.Join(bin, "artifact")
			require.NoError(t, os.WriteFile(artifact, []byte("original"), 0o644))
			require.NoError(t, markMacOSOutputLocal(bin))
			require.True(t, macOSOutputIsLocal(bin))

			// The provider tags an ancestor after the first probe was cached.
			// Domain IDs are provider-owned; the other markers can be set in a test.
			require.NoError(t, unix.Setxattr(providerRoot, attribute, []byte("test"), 0))
			require.False(t, macOSOutputIsLocal(bin), "cached trust survived provider metadata")
			require.Error(t, ensureMacOSOutputIsLocal(bin), "provider output root was accepted")
			t.Setenv("WAILS_MACOS_OUTPUT_ROOT", t.TempDir())
			t.Chdir(project)
			require.NoError(t, ToolPrepareMacOSOutput(&PrepareMacOSOutputOptions{Dir: "bin"}))
			linked, err := isOutputLink(bin)
			require.NoError(t, err)
			require.True(t, linked)
			require.True(t, macOSOutputIsLocal(bin))
			data, err := os.ReadFile(artifact)
			require.NoError(t, err)
			require.Equal(t, "original", string(data), "relocation damaged output")
		})
	}
}
