package pipeline

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/buildinfo"
)

func TestCapabilityRegistryDefinesTheClosedV3TargetSet(t *testing.T) {
	assert.Equal(t, []string{
		"windows/amd64", "windows/arm64",
		"darwin/amd64", "darwin/arm64", "darwin/universal",
		"linux/amd64", "linux/arm64",
		"ios/arm64",
		"android/amd64", "android/arm64", "android/universal",
	}, supportedTargetNames())
	seen := map[string]bool{}
	for _, capability := range buildinfo.TargetCapabilities() {
		name := capability.Target.OS + "/" + capability.Target.Arch
		assert.False(t, seen[name], name)
		seen[name] = true
		assert.NotZero(t, capability.Toolchains, name)
		for index := range int(capability.ComponentCount) {
			component := capability.Component(index)
			assert.Equal(t, capability.Target.OS, component.OS)
			child, ok := lookupTarget(component.OS, component.Arch)
			require.True(t, ok, component)
			assert.False(t, child.Synthetic(), component)
		}
		assert.Empty(t, capability.Component(-1))
		assert.Empty(t, capability.Component(int(capability.ComponentCount)))
	}
	_, ok := lookupTarget("plan9", "amd64")
	assert.False(t, ok)
	assert.ErrorContains(t, unsupportedTargetError("plan9", "amd64"), "supported targets")
}

func TestCapabilityRegistryOwnsProductionAndDevelopmentFormats(t *testing.T) {
	seen := map[string]bool{}
	for _, format := range buildinfo.FormatCapabilities() {
		assert.False(t, seen[format.Name], format.Name)
		seen[format.Name] = true
		assert.NotZero(t, format.Hosts, format.Name)
		assert.NotEqual(t, format.Production, format.Development, format.Name)
		matched := false
		for _, target := range buildinfo.TargetCapabilities() {
			if target.Target.OS == format.Platform && target.SupportsFormat(format.Name, format.Development) {
				matched = true
			}
		}
		assert.True(t, matched, format.Name)
		for index := range int(format.ToolCount) {
			assert.NotEmpty(t, format.RequiredTool(index), format.Name)
		}
		assert.Empty(t, format.RequiredTool(-1))
		assert.Empty(t, format.RequiredTool(int(format.ToolCount)))
	}

	aab, ok := lookupFormat("aab")
	require.True(t, ok)
	assert.True(t, aab.Production)
	assert.False(t, aab.Development)
	apk, ok := lookupFormat("apk")
	require.True(t, ok)
	assert.False(t, apk.Production)
	assert.True(t, apk.Development)
	_, ok = lookupFormat("unknown")
	assert.False(t, ok)

	android, _ := lookupTarget("android", "arm64")
	assert.True(t, android.SupportsFormat("aab", false))
	assert.False(t, android.SupportsFormat("apk", false))
	assert.True(t, android.SupportsFormat("apk", true))
	assert.False(t, android.SupportsFormat("aab", true))
}

func TestCapabilityRegistryOwnsDeterministicOutputNames(t *testing.T) {
	tests := map[string]string{
		"nsis":      "dist/example-installer.exe",
		"msix":      "dist/example.msix",
		"dmg":       "dist/example.dmg",
		"appimage":  "dist/example-arm64.AppImage",
		"deb":       "dist/example_1.2.3_arm64.deb",
		"rpm":       "dist/example-1.2.3.arm64.rpm",
		"archlinux": "dist/example-1.2.3-arm64.pkg.tar.zst",
		"ipa":       "dist/example.ipa",
		"aab":       "dist/example.aab",
		"apk":       "dist/example.apk",
	}
	for name, expected := range tests {
		format, ok := lookupFormat(name)
		require.True(t, ok)
		assert.Equal(t, expected, format.OutputPath("dist", "example", "1.2.3", "arm64"), name)
	}
}

func TestCapabilityRegistryOwnsHostAndToolchainRules(t *testing.T) {
	for _, host := range []string{"windows", "darwin", "linux"} {
		assert.NotZero(t, hostBit(host))
	}
	assert.Zero(t, hostBit("plan9"))
	for _, toolchain := range []string{"", "auto", "native", "zig", "docker"} {
		assert.NotZero(t, toolchainBit(toolchain), toolchain)
	}
	assert.Zero(t, toolchainBit("magic"))

	msix, _ := lookupFormat("msix")
	assert.True(t, msix.SupportsHost("windows"))
	assert.False(t, msix.SupportsHost("linux"))
	dmg, _ := lookupFormat("dmg")
	assert.True(t, dmg.SupportsHost("darwin"))
	assert.False(t, dmg.SupportsHost("windows"))
	deb, _ := lookupFormat("deb")
	assert.True(t, deb.SupportsHost("windows"))
	assert.True(t, deb.SupportsHost("darwin"))
	assert.True(t, deb.SupportsHost("linux"))
	assert.False(t, deb.SupportsHost("plan9"))

	windows, _ := lookupTarget("windows", "amd64")
	assert.True(t, windows.SupportsToolchain("zig"))
	ios, _ := lookupTarget("ios", "arm64")
	assert.False(t, ios.SupportsToolchain("zig"))
}

func TestSyntheticTargetsExpandIntoConcreteCompileOperations(t *testing.T) {
	config := testConfig(t)

	darwin, err := PlanBuild(config, Request{Verb: "build", TargetOS: "darwin", TargetArch: "universal"})
	require.NoError(t, err)
	assert.NotContains(t, darwin.Nodes, NodeKey("target:darwin/universal:compile"))
	assert.Contains(t, darwin.Nodes, NodeKey("target:darwin/amd64:compile"))
	assert.Contains(t, darwin.Nodes, NodeKey("target:darwin/arm64:compile"))
	merge := darwin.Nodes[NodeKey("assemble:darwin/universal:binary")]
	assert.Contains(t, merge.Dependencies, NodeKey("target:darwin/amd64:compile"))
	assert.Contains(t, merge.Dependencies, NodeKey("target:darwin/arm64:compile"))
	assembly := darwin.Nodes[NodeKey("assemble:darwin/universal")]
	assert.Contains(t, assembly.Dependencies, NodeKey("assemble:darwin/universal:binary"))

	android, err := PlanBuild(config, Request{Verb: "build", TargetOS: "android", TargetArch: "universal"})
	require.NoError(t, err)
	assert.NotContains(t, android.Nodes, NodeKey("target:android/universal:compile"))
	assert.Contains(t, android.Nodes, NodeKey("target:android/amd64:compile"))
	assert.Contains(t, android.Nodes, NodeKey("target:android/arm64:compile"))
	packaging := android.Nodes[NodeKey("package:android/universal:aab")]
	assert.Contains(t, packaging.Dependencies, NodeKey("target:android/amd64:compile"))
	assert.Contains(t, packaging.Dependencies, NodeKey("target:android/arm64:compile"))
}
