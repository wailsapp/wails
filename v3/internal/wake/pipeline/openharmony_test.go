package pipeline

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
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

func TestHarmonyOSPlanRegistersAssetsCompileAndHosHapPackage(t *testing.T) {
	config := testConfig(t)
	plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: "harmonyos", TargetArch: "amd64", Formats: []string{"hos_hap"}})
	require.NoError(t, err)

	assets, ok := plan.Nodes[NodeKey("target:harmonyos/amd64:assets")]
	require.True(t, ok, "HarmonyOS builds must generate platform assets")
	assetsSpec, ok := assets.Spec.(AssetsSpec)
	require.True(t, ok)
	assert.Equal(t, "harmonyos", assetsSpec.TargetOS)

	compile, ok := plan.Nodes[NodeKey("target:harmonyos/amd64:compile")]
	require.True(t, ok, "HarmonyOS builds must compile the Go shared library")
	compileSpec, ok := compile.Spec.(CompileSpec)
	require.True(t, ok)
	assert.Equal(t, "libwails_go.so", filepath.Base(filepath.ToSlash(compileSpec.Output)))
	assert.Equal(t, "harmonyos", compileSpec.TargetOS)
	assert.Contains(t, compile.Dependencies, NodeKey("target:harmonyos/amd64:assets"))

	packaged, ok := plan.Nodes[NodeKey("package:harmonyos/amd64:hos_hap")]
	require.True(t, ok, "HarmonyOS builds must produce a commercial HAP")
	packageSpec, ok := packaged.Spec.(PackageSpec)
	require.True(t, ok)
	assert.Equal(t, "hos_hap", packageSpec.Format)
	assert.True(t, strings.HasSuffix(filepath.ToSlash(packageSpec.Output), ".hap"))
}

func TestHarmonyOSPlanUsesHarmonyOSPlatformSettings(t *testing.T) {
	config := testConfig(t)
	config.Targets.HarmonyOS.Identifier = "com.example.harmony"
	config.Signing.HarmonyOS = manifest.SigningPlatform{Enabled: true, Certificate: "h.p12", KeyAlias: "wails", Credential: "HOS_PASSWORD"}

	assert.Equal(t, "com.example.harmony", platformConfig(config.Targets, "harmonyos").Identifier)
	assert.Equal(t, config.Targets.HarmonyOS.AMD64, targetConfig(config.Targets, "harmonyos", "amd64"))
	assert.Equal(t, config.Signing.HarmonyOS, signingPlatform(config.Signing, "harmonyos"))
	assert.Equal(t, config.Signing.HarmonyOS.Certificate, signingPlatformForPlan(config, "harmonyos", "assets").Certificate)
}

func TestHarmonyOSHostChecksFollowOpenHarmonyFamily(t *testing.T) {
	config := testConfig(t)
	request := Request{Verb: "build", TargetOS: "harmonyos", TargetArch: "arm64", Formats: []string{"hos_hap"}, CGOEnabled: "1"}

	withoutSDK := NewHostCapabilitiesWithFacts("windows", "amd64", []string{"go", "npm"}, nil, HostFacts{})
	_, err := PlanBuildForHost(config, request, withoutSDK)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "OpenHarmony SDK")

	withSDK := NewHostCapabilitiesWithFacts("windows", "amd64", []string{"go", "npm"}, nil, HostFacts{OHSDK: true})
	plan, err := PlanBuildForHost(config, request, withSDK)
	require.NoError(t, err)
	compile := plan.Nodes[NodeKey("target:harmonyos/arm64:compile")].Spec.(CompileSpec)
	assert.Equal(t, "native", compile.Toolchain)
	assert.Equal(t, "harmonyos", compile.TargetOS)
}

func TestHarmonyOSSigningHostRequiresTheSharedOpenHarmonyFields(t *testing.T) {
	host := NewHostCapabilitiesWithFacts("windows", "amd64", nil, []string{"HOS_PASSWORD"}, HostFacts{OHSDK: true})
	valid := SignSpec{TargetOS: "harmonyos", TargetArch: "arm64", Format: "hos_hap", Config: manifest.SigningPlatform{Enabled: true, Certificate: "h.p12", KeyAlias: "wails", Credential: "HOS_PASSWORD"}}
	require.NoError(t, validateSigningHost(valid, host))

	missing := valid
	missing.Config.KeyAlias = ""
	err := validateSigningHost(missing, host)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "harmonyos signing requires certificate")

	noSDK := NewHostCapabilitiesWithFacts("windows", "amd64", nil, []string{"HOS_PASSWORD"}, HostFacts{})
	err = validateSigningHost(valid, noSDK)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "OpenHarmony SDK")
}

func TestHarmonyOSCompileModeRequiresCGO(t *testing.T) {
	_, err := compileMode(CompileSpec{TargetOS: "harmonyos", TargetArch: "arm64"}, "0")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mobile targets require CGO_ENABLED=1")
	mode, err := compileMode(CompileSpec{TargetOS: "harmonyos", TargetArch: "arm64"}, "1")
	require.NoError(t, err)
	assert.Equal(t, "1", mode)
}
