package pipeline

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestPlannerInternalBoundariesPreserveErrors(t *testing.T) {
	config := testConfig(t)
	plan, err := PlanBuild(config, Request{TargetOS: "linux", TargetArch: "amd64"})
	require.NoError(t, err)
	assert.Equal(t, "build", plan.Name)

	combined := Plan{Nodes: map[NodeKey]Node{
		"shared": {Key: "shared", Kind: CompileApplication, Spec: CompileSpec{TargetOS: "linux"}, Cache: CacheArtifact},
	}}
	child := Plan{Target: "windows/amd64", Nodes: map[NodeKey]Node{
		"shared": {Key: "shared", Kind: CompileApplication, Spec: CompileSpec{TargetOS: "windows"}, Cache: CacheArtifact},
	}}
	require.ErrorContains(t, addPlanNodes(&combined, child), "resolves differently")
	config.Targets.Darwin.Universal.Tags = []string{"universal"}
	config.Targets.Darwin.ARM64.Tags = []string{"arm64"}
	_, err = PlanBuild(config, Request{Verb: "build", Targets: []Target{{OS: "darwin", Arch: "universal"}, {OS: "darwin", Arch: "arm64"}}})
	assert.ErrorContains(t, err, "resolves differently across targets")

	defaulted, err := planTarget(config, Request{}, false)
	require.NoError(t, err)
	assert.Equal(t, runtime.GOOS+"/"+runtime.GOARCH, defaulted.Target)
	_, err = planTarget(config, Request{TargetOS: "plan9", TargetArch: "amd64"}, false)
	assert.ErrorContains(t, err, "supported targets")

	_, err = planTarget(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64", Formats: []string{"unknown"}}, false)
	assert.ErrorContains(t, err, "unknown package format")
	_, err = planTarget(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64", Formats: []string{"nsis"}, Development: true}, false)
	assert.ErrorContains(t, err, "not supported for linux/amd64 in development")
	_, err = planTarget(config, Request{Verb: "build", TargetOS: "ios", TargetArch: "arm64", Formats: []string{"ipa"}, destination: "simulator"}, false)
	assert.ErrorContains(t, err, `requires profile destination = "device"`)
}

func TestPlannerCoversDevelopmentAppAndExplicitObfuscation(t *testing.T) {
	config := testConfig(t)
	plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: "darwin", TargetArch: "arm64", Development: true})
	require.NoError(t, err)
	assembly := plan.Nodes["assemble:darwin/arm64"]
	assert.Equal(t, "bin/app.app", assembly.Output)

	plan, err = PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64", Obfuscated: true})
	require.NoError(t, err)
	compile := plan.Nodes["target:linux/amd64:compile"].Spec.(CompileSpec)
	assert.True(t, compile.Obfuscated)
	assert.Contains(t, compile.Tags, "wails_obfuscated")
}

func TestPlannerGoMetadataFailuresAreObservable(t *testing.T) {
	injected := errors.New("injected absolute path failure")
	_, err := goLocalSourceInputsWithAbs("project", func(string) (string, error) { return "", injected })
	require.ErrorIs(t, err, injected)
	assert.Equal(t, []string{"go.mod", "go.sum", "go.work", "go.work.sum"}, goMetadataFilesWithAbs("project", func(string) (string, error) { return "", injected }))

	t.Run("go mod read", func(t *testing.T) {
		root := t.TempDir()
		require.NoError(t, os.Mkdir(filepath.Join(root, "go.mod"), 0o755))
		_, err := goLocalSourceInputs(root)
		require.Error(t, err)
	})
	t.Run("go mod parse", func(t *testing.T) {
		root := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("not a module file\n"), 0o644))
		_, err := goLocalSourceInputs(root)
		require.Error(t, err)
	})
	t.Run("go work parse", func(t *testing.T) {
		root := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(root, "go.work"), []byte("not a workspace file\n"), 0o644))
		_, err := goLocalSourceInputs(root)
		require.Error(t, err)
	})
	t.Run("planner propagation", func(t *testing.T) {
		config := testConfig(t)
		require.NoError(t, os.Mkdir(filepath.Join(config.Root, "go.mod"), 0o755))
		_, err := PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"})
		require.Error(t, err)
	})
}

func TestDevelopmentPlanningMarksLocalModuleInputsForGitIgnore(t *testing.T) {
	config := testConfig(t)
	local := filepath.Join(filepath.Dir(config.Root), "local-module")
	require.NoError(t, os.MkdirAll(local, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(local, "go.mod"), []byte("module example/local\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(config.Root, "go.mod"), []byte("module example/app\ngo 1.23\nreplace example/local => ../local-module\n"), 0o644))
	config.Dev.UseGitIgnore = true
	plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64", Development: true})
	require.NoError(t, err)
	compile := plan.Nodes["target:linux/amd64:compile"]
	found := false
	for _, input := range compile.Inputs {
		if input.Label == "go-local-source" {
			found = true
			assert.True(t, input.UseGitIgnore)
		}
	}
	assert.True(t, found)
}

func TestPlannerClosedPlatformHelpersCoverEverySlot(t *testing.T) {
	packages := manifest.Packages{}
	assert.Equal(t, packages.Windows, packagePlatform(packages, "windows"))
	assert.Equal(t, packages.Darwin, packagePlatform(packages, "darwin"))
	assert.Equal(t, packages.Linux, packagePlatform(packages, "linux"))
	assert.Equal(t, packages.IOS, packagePlatform(packages, "ios"))
	assert.Equal(t, packages.Android, packagePlatform(packages, "android"))

	targets := manifest.Targets{}
	assert.Equal(t, targets.Windows.ARM, targetConfig(targets, "windows", "arm"))
	assert.Equal(t, targets.Windows.X86, targetConfig(targets, "windows", "386"))
	assert.Equal(t, targets.Darwin.Universal, targetConfig(targets, "darwin", "universal"))
	assert.Equal(t, manifest.Target{}, targetConfig(targets, "unknown", "unknown"))
	assert.Equal(t, 2, max(2, 1))

	config := manifest.Config{}
	config.Targets.Windows.Manifest = "windows.manifest"
	config.Targets.Darwin.AssetsCar = "Assets.car"
	config.Targets.Linux.DesktopEntry = "application.desktop"
	config.Targets.IOS.InfoPlist = "Info.plist"
	config.Targets.Android.Icon = "android.png"
	files := assetInputs(config)[0].Files
	for _, path := range []string{"windows.manifest", "Assets.car", "application.desktop", "Info.plist", "android.png"} {
		assert.Contains(t, files, path)
	}

	assert.Nil(t, originsForNode(manifest.Config{}, Node{Kind: CollectArtifacts}, "linux/amd64"))
	assert.Panics(t, func() { registeredPackageFormat(manifest.Packages{}, "plan9", "pkg") })
}
