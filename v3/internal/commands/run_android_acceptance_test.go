package commands

import (
	"archive/zip"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

// This opt-in check needs an installed Android SDK/NDK and Java. It exercises
// the real example and run pipeline, independently of emulator availability.
func TestRunAndroidAPKAcceptance(t *testing.T) {
	if os.Getenv("WAILS_TEST_ANDROID_RUN") != "1" {
		t.Skip("set WAILS_TEST_ANDROID_RUN=1 with Android SDK/NDK and Java installed")
	}
	root, err := filepath.Abs(filepath.Join("..", "..", "examples", "android"))
	require.NoError(t, err)
	loaded, err := manifest.Load(root, "")
	require.NoError(t, err)
	settings, err := loaded.Config.RunForTarget("android", "amd64")
	require.NoError(t, err)
	options := manifestRunOptions{Context: context.Background(), Verb: "run", Loaded: loaded, TargetOS: "android", TargetArch: "amd64", Tags: settings.Tags}
	require.NoError(t, configureRunTarget(&RunOptions{}, "android", "amd64", settings, &options))
	result, err := runManifestPipelineResult(options)
	require.NoError(t, err)
	apk, err := runArtifactPath(root, result.Plan, "android", loaded.Config.Project.BinaryName)
	require.NoError(t, err)
	archive, err := zip.OpenReader(apk)
	require.NoError(t, err)
	defer archive.Close()
	files := map[string]bool{}
	for _, file := range archive.File {
		files[file.Name] = true
	}
	for _, name := range []string{"AndroidManifest.xml", "classes.dex", "lib/x86_64/libwails.so"} {
		require.True(t, files[name], "APK must contain %s", name)
	}
}
