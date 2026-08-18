package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func TestRunManifestPipelineMarksRenderedFailure(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	require.NoError(t, manifest.WriteMinimal(root, manifest.Project{
		Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0",
	}))
	t.Setenv("PATH", t.TempDir())

	err := runManifestPipeline(manifestRunOptions{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	require.Error(t, err)
	assert.True(t, wake.IsReported(err), "execution failure was rendered by Pulse and must not be printed again by the CLI")
}

func TestRunManifestPipelineLeavesUnrenderedFailurePrintable(t *testing.T) {
	t.Chdir(t.TempDir())

	err := runManifestPipeline(manifestRunOptions{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
	require.Error(t, err)
	assert.False(t, wake.IsReported(err), "errors raised before Pulse starts still need the CLI error printer")
}

func TestResolveManifestPlanUsesTheDiscoveredManifestRoot(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, manifest.WriteMinimal(root, manifest.Project{
		Name: "app", ProductName: "App", Identifier: "com.example.app", Version: "1.0.0",
	}))
	nested := filepath.Join(root, "frontend", "src")
	require.NoError(t, os.MkdirAll(nested, 0o755))
	loaded, err := manifest.Load(nested, "")
	require.NoError(t, err)
	t.Chdir(nested)

	resolved, _, _, err := resolveManifestPlan(manifestRunOptions{Verb: "build", Loaded: loaded})
	require.NoError(t, err)
	assert.Equal(t, root, resolved)
}

func TestApplyGeneratedTargetSettings(t *testing.T) {
	root := t.TempDir()
	write := func(relative, content string) {
		path := filepath.Join(root, relative)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	}
	write("ios/xcode/main/Info.plist", "<key>CFBundleVersion</key><string>1.0.0</string><key>MinimumOSVersion</key><string>15.0</string>")
	write("android/app/build.gradle", "android {\n    versionCode 1\n}\n")
	require.NoError(t, applyGeneratedTargetSettings(root, pipeline.AssetsSpec{TargetOS: "ios", MinimumVersion: "17.0", Project: manifest.Project{BuildNumber: 42}}))
	plist, err := os.ReadFile(filepath.Join(root, "ios/xcode/main/Info.plist"))
	require.NoError(t, err)
	assert.Contains(t, string(plist), "<string>42</string>")
	assert.Contains(t, string(plist), "<string>17.0</string>")
	gradle, err := os.ReadFile(filepath.Join(root, "android/app/build.gradle"))
	require.NoError(t, err)
	assert.Contains(t, string(gradle), "versionCode 42")
}

func TestReplacePlistStringFillsSelfClosingValue(t *testing.T) {
	path := filepath.Join(t.TempDir(), "Info.plist")
	require.NoError(t, os.WriteFile(path, []byte("<key>CFBundleVersion</key>\n<string/>\n"), 0o644))
	require.NoError(t, replacePlistString(path, "CFBundleVersion", "13"))
	assert.Equal(t, "<key>CFBundleVersion</key>\n<string>13</string>\n", readTestFile(t, path))
}

func TestCopyManifestPathPreservesSymlink(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	destination := filepath.Join(root, "destination")
	require.NoError(t, os.MkdirAll(source, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(source, "actual"), []byte("payload"), 0o644))
	require.NoError(t, os.Symlink("actual", filepath.Join(source, "link")))
	require.NoError(t, copyManifestPath(source, destination))
	info, err := os.Lstat(filepath.Join(destination, "link"))
	require.NoError(t, err)
	assert.NotZero(t, info.Mode()&os.ModeSymlink)
	target, err := os.Readlink(filepath.Join(destination, "link"))
	require.NoError(t, err)
	assert.Equal(t, "actual", target)
}

func TestFindAndMovePackageIgnoresExistingDestination(t *testing.T) {
	root := t.TempDir()
	staging := filepath.Join(root, "staging")
	require.NoError(t, os.MkdirAll(staging, 0o755))
	expected := filepath.Join(root, "app.deb")
	require.NoError(t, os.WriteFile(expected, []byte("old"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(staging, "app_1.0.0_amd64.deb"), []byte("new"), 0o644))

	require.NoError(t, findAndMovePackage(staging, "app", "deb", expected))
	data, err := os.ReadFile(expected)
	require.NoError(t, err)
	assert.Equal(t, "new", string(data))
}

func TestManifestHandlerIdentityIncludesRelevantEnvironment(t *testing.T) {
	handler := &manifestHandler{root: t.TempDir()}
	node := pipeline.Node{Key: "assets", Kind: pipeline.GeneratePlatformAssets, Spec: pipeline.AssetsSpec{}}
	t.Setenv("PATH", "/first")
	first, err := handler.Identity(t.Context(), node)
	require.NoError(t, err)
	t.Setenv("PATH", "/second")
	second, err := handler.Identity(t.Context(), node)
	require.NoError(t, err)
	assert.NotEqual(t, first, second)
}

func TestDevelopmentFrontendIdentityIncludesSessionURLAndPort(t *testing.T) {
	node := pipeline.Node{Kind: pipeline.BuildFrontend, Spec: pipeline.FrontendSpec{Production: false}}
	t.Setenv("FRONTEND_DEVSERVER_URL", "http://localhost:9245")
	t.Setenv(wailsVitePort, "9245")
	first := relevantEnvironment(node, nil)
	t.Setenv("FRONTEND_DEVSERVER_URL", "http://localhost:9246")
	assert.NotEqual(t, first, relevantEnvironment(node, nil))
	t.Setenv("FRONTEND_DEVSERVER_URL", "http://localhost:9245")
	t.Setenv(wailsVitePort, "9246")
	assert.NotEqual(t, first, relevantEnvironment(node, nil))

	production := pipeline.Node{Kind: pipeline.BuildFrontend, Spec: pipeline.FrontendSpec{Production: true}}
	first = relevantEnvironment(production, nil)
	t.Setenv("FRONTEND_DEVSERVER_URL", "http://localhost:9247")
	t.Setenv(wailsVitePort, "9247")
	second := relevantEnvironment(production, nil)
	assert.Equal(t, first, second)
}

func TestDevelopmentEnvironmentOverridesAreIsolatedFromProcessState(t *testing.T) {
	node := pipeline.Node{Kind: pipeline.BuildFrontend, Spec: pipeline.FrontendSpec{Production: false}}
	t.Setenv("FRONTEND_DEVSERVER_URL", "http://global:9245")
	t.Setenv(wailsVitePort, "9245")
	first := relevantEnvironment(node, []string{"FRONTEND_DEVSERVER_URL=http://generation-one:9246", wailsVitePort + "=9246"})
	second := relevantEnvironment(node, []string{"FRONTEND_DEVSERVER_URL=http://generation-two:9247", wailsVitePort + "=9247"})
	assert.NotEqual(t, first, second)
	assert.Equal(t, "http://global:9245", os.Getenv("FRONTEND_DEVSERVER_URL"))
	assert.Equal(t, "9245", os.Getenv(wailsVitePort))
}
