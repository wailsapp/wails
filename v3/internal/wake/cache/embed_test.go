package cache

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestEmbeddedResourcesBypassSourceFiltersAndTrackMembership(t *testing.T) {
	for _, pattern := range []string{"build/data.json", "build/*.json", "build", "all:build", "\"build/with space.json\""} {
		t.Run(pattern, func(t *testing.T) {
			root := t.TempDir()
			require.NoError(t, os.MkdirAll(filepath.Join(root, "build"), 0755))
			resource := filepath.Join(root, "build", "data.json")
			if pattern == "\"build/with space.json\"" {
				resource = filepath.Join(root, "build", "with space.json")
			}
			require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\nimport _ \"embed\"\n//go:embed "+pattern+"\nvar data string\n"), 0644))
			require.NoError(t, os.WriteFile(resource, []byte("first"), 0644))
			store, err := OpenCache(root)
			require.NoError(t, err)
			options := SnapshotOptions{Label: "compile", Root: root, IncludeGoEmbed: true, IncludeExtensions: []string{".go"}, ExcludeDirs: []string{"build"}}
			before, err := store.Snapshot(options)
			require.NoError(t, err)
			require.NoError(t, store.Save())
			store, err = OpenCache(root)
			require.NoError(t, err)
			warm, err := store.Snapshot(options)
			require.NoError(t, err)
			require.Equal(t, before, warm)
			require.NoError(t, os.WriteFile(resource, []byte("second"), 0644))
			after, err := store.Snapshot(options)
			require.NoError(t, err)
			require.NotEqual(t, before, after)
			require.NoError(t, os.Remove(resource))
			deleted, err := store.Snapshot(options)
			require.NoError(t, err)
			require.NotEqual(t, after, deleted)
			require.NoError(t, os.WriteFile(resource, []byte("first"), 0644))
			restored, err := store.Snapshot(options)
			require.NoError(t, err)
			require.Equal(t, before, restored)
		})
	}
}
func TestParseEmbedPatterns(t *testing.T) {
	patterns, err := parseEmbedPatterns("file.json \"with space.json\" `raw space` all:assets/*")
	require.NoError(t, err)
	require.Equal(t, []string{"file.json", "with space.json", "raw space", "all:assets/*"}, patterns)
	_, err = parseEmbedPatterns("\"unterminated")
	require.Error(t, err)
}
