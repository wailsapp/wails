package buildinfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistryDefinesClosedTargetsAndDefensiveEnumeration(t *testing.T) {
	want := []string{
		"windows/amd64", "windows/arm64", "darwin/amd64", "darwin/arm64", "darwin/universal",
		"linux/amd64", "linux/arm64", "ios/arm64", "android/amd64", "android/arm64", "android/universal",
	}
	assert.Equal(t, want, SupportedTargetNames())
	capabilities := TargetCapabilities()
	require.Len(t, capabilities, len(want))
	capabilities[0].Target.OS = "mutated"
	first, ok := LookupTarget("windows", "amd64")
	require.True(t, ok)
	assert.Equal(t, "windows", first.Target.OS)

	for _, capability := range TargetCapabilities() {
		assert.NotZero(t, capability.Toolchains)
		for index := range int(capability.ComponentCount) {
			component := capability.Component(index)
			assert.Equal(t, capability.Target.OS, component.OS)
			child, found := LookupTarget(component.OS, component.Arch)
			require.True(t, found)
			assert.False(t, child.Synthetic())
		}
		assert.Empty(t, capability.Component(-1))
		assert.Empty(t, capability.Component(int(capability.ComponentCount)))
	}
	_, ok = LookupTarget("plan9", "amd64")
	assert.False(t, ok)
}

func TestRegistryDefinesFormatsModesAndDefensiveEnumeration(t *testing.T) {
	formats := FormatCapabilities()
	require.NotEmpty(t, formats)
	formats[0].Name = "mutated"
	_, ok := LookupFormat("nsis")
	require.True(t, ok)

	seen := map[string]bool{}
	for _, format := range FormatCapabilities() {
		assert.False(t, seen[format.Name])
		seen[format.Name] = true
		assert.NotZero(t, format.Hosts)
		assert.NotEqual(t, format.Production, format.Development)
		matched := false
		for _, target := range TargetCapabilities() {
			matched = matched || target.SupportsFormat(format.Name, format.Development)
		}
		assert.True(t, matched, format.Name)
		for index := range int(format.ToolCount) {
			assert.NotEmpty(t, format.RequiredTool(index))
		}
		assert.Empty(t, format.RequiredTool(-1))
		assert.Empty(t, format.RequiredTool(int(format.ToolCount)))
	}
	_, ok = LookupFormat("unknown")
	assert.False(t, ok)
	android, _ := LookupTarget("android", "arm64")
	assert.True(t, android.SupportsFormat("aab", false))
	assert.False(t, android.SupportsFormat("apk", false))
	assert.True(t, android.SupportsFormat("apk", true))
	assert.False(t, android.SupportsFormat("aab", true))
	assert.False(t, android.SupportsFormat("unknown", false))
	windows, _ := LookupTarget("windows", "amd64")
	assert.False(t, windows.SupportsFormat("dmg", false))
	assert.False(t, (TargetCapability{Target: Target{OS: "android", Arch: "arm64"}}).SupportsFormat("aab", false))
}

func TestRegistryOwnsOutputHostAndToolchainRules(t *testing.T) {
	wantPaths := map[string]string{
		"nsis": "dist/example-installer.exe", "msix": "dist/example.msix", "dmg": "dist/example.dmg",
		"appimage": "dist/example-arm64.AppImage", "deb": "dist/example_1.2.3_arm64.deb",
		"rpm": "dist/example-1.2.3.arm64.rpm", "archlinux": "dist/example-1.2.3-arm64.pkg.tar.zst",
		"ipa": "dist/example.ipa", "aab": "dist/example.aab", "apk": "dist/example.apk",
	}
	for name, want := range wantPaths {
		format, ok := LookupFormat(name)
		require.True(t, ok)
		assert.Equal(t, want, format.OutputPath("dist", "example", "1.2.3", "arm64"))
	}
	for _, host := range []string{"windows", "darwin", "linux"} {
		assert.NotZero(t, HostMaskFor(host))
	}
	assert.Zero(t, HostMaskFor("plan9"))
	for _, toolchain := range []string{"", "auto", "native", "zig", "docker"} {
		assert.NotZero(t, ToolchainMaskFor(toolchain))
	}
	assert.Zero(t, ToolchainMaskFor("magic"))

	assert.Equal(t, "windows", RequiredHostName(HostWindows))
	assert.Equal(t, "darwin", RequiredHostName(HostDarwin))
	assert.Equal(t, "linux", RequiredHostName(HostLinux))
	assert.Equal(t, "supported", RequiredHostName(HostAny))

	msix, _ := LookupFormat("msix")
	assert.True(t, msix.SupportsHost("windows"))
	assert.False(t, msix.SupportsHost("linux"))
	deb, _ := LookupFormat("deb")
	assert.True(t, deb.SupportsHost("windows"))
	assert.True(t, deb.SupportsHost("darwin"))
	assert.True(t, deb.SupportsHost("linux"))
	assert.False(t, deb.SupportsHost("plan9"))
	windows, _ := LookupTarget("windows", "amd64")
	assert.True(t, windows.SupportsToolchain("zig"))
	ios, _ := LookupTarget("ios", "arm64")
	assert.False(t, ios.SupportsToolchain("zig"))
}
