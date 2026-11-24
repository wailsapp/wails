package pipeline

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenHarmonyPlanRegistersAssetsCompileAndHapPackage(t *testing.T) {
	config := testConfig(t)
	plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: "openharmony", TargetArch: "amd64", Formats: []string{"hap"}})
	require.NoError(t, err)

	assets, ok := plan.Nodes[NodeKey("target:openharmony/amd64:assets")]
	require.True(t, ok, "OpenHarmony builds must generate platform assets")
	assert.Equal(t, GeneratePlatformAssets, assets.Kind)
	assetsSpec, ok := assets.Spec.(AssetsSpec)
	require.True(t, ok)
	assert.Equal(t, "openharmony", assetsSpec.TargetOS)

	compile, ok := plan.Nodes[NodeKey("target:openharmony/amd64:compile")]
	require.True(t, ok, "OpenHarmony builds must compile the Go shared library")
	assert.Equal(t, CompileApplication, compile.Kind)
	compileSpec, ok := compile.Spec.(CompileSpec)
	require.True(t, ok)
	assert.Equal(t, "libwails_go.so", filepath.Base(filepath.ToSlash(compileSpec.Output)))
	assert.Equal(t, "openharmony", compileSpec.TargetOS)

	packaged, ok := plan.Nodes[NodeKey("package:openharmony/amd64:hap")]
	require.True(t, ok, "OpenHarmony builds must produce a HAP")
	assert.Equal(t, PackageArtifact, packaged.Kind)
	packageSpec, ok := packaged.Spec.(PackageSpec)
	require.True(t, ok)
	assert.Equal(t, "hap", packageSpec.Format)
	assert.Contains(t, compile.Dependencies, NodeKey("target:openharmony/amd64:assets"))
}
