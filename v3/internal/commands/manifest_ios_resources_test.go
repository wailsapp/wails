package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"howett.net/plist"
)

func TestIOSAssetMetadataPreservesExplicitPlistSettings(t *testing.T) {
	root := t.TempDir()
	infoPath := filepath.Join(root, "Info.plist")
	partialPath := filepath.Join(root, "asset-info.plist")
	info := map[string]any{"CFBundleIdentifier": "com.example.custom", "UIBackgroundModes": []any{"fetch"}, "CFBundleIcons": map[string]any{"CFBundlePrimaryIcon": map[string]any{"CFBundleIconName": "CustomIcon"}}}
	partial := map[string]any{"CFBundleIcons": map[string]any{"CFBundlePrimaryIcon": map[string]any{"CFBundleIconName": "AppIcon", "CFBundleIconFiles": []any{"AppIcon60x60"}}}, "CFBundleIcons~ipad": map[string]any{"CFBundlePrimaryIcon": map[string]any{"CFBundleIconName": "AppIcon"}}}
	data, err := plist.Marshal(info, plist.BinaryFormat)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(infoPath, data, 0644))
	data, err = plist.Marshal(partial, plist.XMLFormat)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(partialPath, data, 0644))
	require.NoError(t, mergeIOSAssetInfo(infoPath, partialPath))
	data, err = os.ReadFile(infoPath)
	require.NoError(t, err)
	var got map[string]any
	_, err = plist.Unmarshal(data, &got)
	require.NoError(t, err)
	require.Equal(t, "com.example.custom", got["CFBundleIdentifier"])
	require.Equal(t, []any{"fetch"}, got["UIBackgroundModes"])
	icon := got["CFBundleIcons"].(map[string]any)["CFBundlePrimaryIcon"].(map[string]any)
	require.Equal(t, "CustomIcon", icon["CFBundleIconName"])
	require.Equal(t, []any{"AppIcon60x60"}, icon["CFBundleIconFiles"])
	require.Contains(t, got, "CFBundleIcons~ipad")
}

func TestIOSAssetMetadataRejectsInvalidPlistsWithoutChangingInfo(t *testing.T) {
	for _, invalid := range []string{"Info.plist", "asset-info.plist"} {
		t.Run(invalid, func(t *testing.T) {
			root := t.TempDir()
			info := filepath.Join(root, "Info.plist")
			partial := filepath.Join(root, "asset-info.plist")
			for _, path := range []string{info, partial} {
				require.NoError(t, os.WriteFile(path, []byte(`<plist version="1.0"><dict/></plist>`), 0644))
			}
			require.NoError(t, os.WriteFile(filepath.Join(root, invalid), []byte("invalid"), 0644))
			before, err := os.ReadFile(info)
			require.NoError(t, err)
			require.ErrorContains(t, mergeIOSAssetInfo(info, partial), "read iOS plist")
			after, err := os.ReadFile(info)
			require.NoError(t, err)
			require.Equal(t, before, after)
		})
	}
}
