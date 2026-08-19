// Package buildinfo owns the closed Wails v3 build capability registry.
// Manifest validation and pipeline planning both consult this package so a
// target or format cannot be accepted by one layer and rejected by another.
package buildinfo

import "path/filepath"

type Target struct {
	OS   string
	Arch string
}

type HostMask uint8

const (
	HostWindows HostMask = 1 << iota
	HostDarwin
	HostLinux
	HostAny = HostWindows | HostDarwin | HostLinux
)

type ToolchainMask uint8

const (
	ToolchainAuto ToolchainMask = 1 << iota
	ToolchainNative
	ToolchainZig
	ToolchainDocker
)

type RunnableKind uint8

const (
	RunnableNone RunnableKind = iota
	RunnableBinary
	RunnableApp
)

type outputStyle uint8

const (
	outputSuffix outputStyle = iota
	outputNSIS
	outputAppImage
	outputDeb
	outputRPM
	outputArchLinux
)

type TargetCapability struct {
	Target         Target
	Components     [2]Target
	ComponentCount uint8
	Formats        [4]string
	FormatCount    uint8
	Toolchains     ToolchainMask
	Runnable       RunnableKind
	RunnableSuffix string
}

func (c TargetCapability) Synthetic() bool { return c.ComponentCount != 0 }

func (c TargetCapability) Component(index int) Target {
	if index < 0 || index >= int(c.ComponentCount) {
		return Target{}
	}
	return c.Components[index]
}

func (c TargetCapability) SupportsFormat(format string, development bool) bool {
	capability, ok := LookupFormat(format)
	if !ok || capability.Platform != c.Target.OS || development && !capability.Development || !development && !capability.Production {
		return false
	}
	for index := range int(c.FormatCount) {
		if c.Formats[index] == format {
			return true
		}
	}
	return false
}

func (c TargetCapability) SupportsToolchain(name string) bool {
	return c.Toolchains&ToolchainMaskFor(name) != 0
}

type FormatCapability struct {
	Name                string
	Platform            string
	Production          bool
	Development         bool
	Hosts               HostMask
	RequiredTools       [3]string
	ToolCount           uint8
	RequiredDestination string
	style               outputStyle
	suffix              string
}

func (c FormatCapability) OutputPath(directory, binaryName, version, arch string) string {
	base := filepath.Join(directory, binaryName)
	switch c.style {
	case outputNSIS:
		return filepath.ToSlash(base + "-installer.exe")
	case outputAppImage:
		return filepath.ToSlash(base + "-" + arch + ".AppImage")
	case outputDeb:
		return filepath.ToSlash(base + "_" + version + "_" + arch + ".deb")
	case outputRPM:
		return filepath.ToSlash(base + "-" + version + "." + arch + ".rpm")
	case outputArchLinux:
		return filepath.ToSlash(base + "-" + version + "-" + arch + ".pkg.tar.zst")
	default:
		return filepath.ToSlash(base + c.suffix)
	}
}

func (c FormatCapability) SupportsHost(host string) bool {
	return c.Hosts&HostMaskFor(host) != 0
}

func (c FormatCapability) RequiredTool(index int) string {
	if index < 0 || index >= int(c.ToolCount) {
		return ""
	}
	return c.RequiredTools[index]
}

// Fixed arrays keep registry storage immutable. Lookup and enumeration return
// values, never registry-owned slices or maps.
var targetRegistry = [...]TargetCapability{
	{Target: Target{OS: "windows", Arch: "amd64"}, Formats: [4]string{"nsis", "msix"}, FormatCount: 2, Toolchains: ToolchainAuto | ToolchainNative | ToolchainZig | ToolchainDocker, Runnable: RunnableBinary, RunnableSuffix: ".exe"},
	{Target: Target{OS: "windows", Arch: "arm64"}, Formats: [4]string{"nsis", "msix"}, FormatCount: 2, Toolchains: ToolchainAuto | ToolchainNative | ToolchainZig | ToolchainDocker, Runnable: RunnableBinary, RunnableSuffix: ".exe"},
	{Target: Target{OS: "darwin", Arch: "amd64"}, Formats: [4]string{"dmg"}, FormatCount: 1, Toolchains: ToolchainAuto | ToolchainNative | ToolchainDocker, Runnable: RunnableApp, RunnableSuffix: ".app"},
	{Target: Target{OS: "darwin", Arch: "arm64"}, Formats: [4]string{"dmg"}, FormatCount: 1, Toolchains: ToolchainAuto | ToolchainNative | ToolchainDocker, Runnable: RunnableApp, RunnableSuffix: ".app"},
	{Target: Target{OS: "darwin", Arch: "universal"}, Components: [2]Target{{OS: "darwin", Arch: "amd64"}, {OS: "darwin", Arch: "arm64"}}, ComponentCount: 2, Formats: [4]string{"dmg"}, FormatCount: 1, Toolchains: ToolchainAuto | ToolchainNative | ToolchainDocker, Runnable: RunnableApp, RunnableSuffix: ".app"},
	{Target: Target{OS: "linux", Arch: "amd64"}, Formats: [4]string{"appimage", "deb", "rpm", "archlinux"}, FormatCount: 4, Toolchains: ToolchainAuto | ToolchainNative | ToolchainZig | ToolchainDocker, Runnable: RunnableBinary},
	{Target: Target{OS: "linux", Arch: "arm64"}, Formats: [4]string{"appimage", "deb", "rpm", "archlinux"}, FormatCount: 4, Toolchains: ToolchainAuto | ToolchainNative | ToolchainZig | ToolchainDocker, Runnable: RunnableBinary},
	{Target: Target{OS: "ios", Arch: "arm64"}, Formats: [4]string{"ipa"}, FormatCount: 1, Toolchains: ToolchainAuto | ToolchainNative, Runnable: RunnableApp, RunnableSuffix: ".app"},
	{Target: Target{OS: "android", Arch: "amd64"}, Formats: [4]string{"aab", "apk"}, FormatCount: 2, Toolchains: ToolchainAuto | ToolchainNative},
	{Target: Target{OS: "android", Arch: "arm64"}, Formats: [4]string{"aab", "apk"}, FormatCount: 2, Toolchains: ToolchainAuto | ToolchainNative},
	{Target: Target{OS: "android", Arch: "universal"}, Components: [2]Target{{OS: "android", Arch: "amd64"}, {OS: "android", Arch: "arm64"}}, ComponentCount: 2, Formats: [4]string{"aab", "apk"}, FormatCount: 2, Toolchains: ToolchainAuto | ToolchainNative},
}

var formatRegistry = [...]FormatCapability{
	{Name: "nsis", Platform: "windows", Production: true, Hosts: HostAny, RequiredTools: [3]string{"makensis"}, ToolCount: 1, style: outputNSIS},
	{Name: "msix", Platform: "windows", Production: true, Hosts: HostWindows, RequiredTools: [3]string{"MakeAppx.exe"}, ToolCount: 1, suffix: ".msix"},
	{Name: "dmg", Platform: "darwin", Production: true, Hosts: HostDarwin, RequiredTools: [3]string{"hdiutil"}, ToolCount: 1, suffix: ".dmg"},
	{Name: "appimage", Platform: "linux", Production: true, Hosts: HostLinux, style: outputAppImage},
	{Name: "deb", Platform: "linux", Production: true, Hosts: HostAny, style: outputDeb},
	{Name: "rpm", Platform: "linux", Production: true, Hosts: HostAny, style: outputRPM},
	{Name: "archlinux", Platform: "linux", Production: true, Hosts: HostAny, style: outputArchLinux},
	{Name: "ipa", Platform: "ios", Production: true, Hosts: HostDarwin, RequiredTools: [3]string{"xcrun", "codesign", "zip"}, ToolCount: 3, RequiredDestination: "device", suffix: ".ipa"},
	{Name: "aab", Platform: "android", Production: true, Hosts: HostAny, RequiredTools: [3]string{"java"}, ToolCount: 1, suffix: ".aab"},
	{Name: "apk", Platform: "android", Development: true, Hosts: HostAny, RequiredTools: [3]string{"java"}, ToolCount: 1, suffix: ".apk"},
}

func LookupTarget(platform, arch string) (TargetCapability, bool) {
	for _, capability := range targetRegistry {
		if capability.Target.OS == platform && capability.Target.Arch == arch {
			return capability, true
		}
	}
	return TargetCapability{}, false
}

func LookupFormat(name string) (FormatCapability, bool) {
	for _, capability := range formatRegistry {
		if capability.Name == name {
			return capability, true
		}
	}
	return FormatCapability{}, false
}

func TargetCapabilities() []TargetCapability {
	result := make([]TargetCapability, len(targetRegistry))
	copy(result, targetRegistry[:])
	return result
}

func FormatCapabilities() []FormatCapability {
	result := make([]FormatCapability, len(formatRegistry))
	copy(result, formatRegistry[:])
	return result
}

func SupportedTargetNames() []string {
	result := make([]string, len(targetRegistry))
	for index, capability := range targetRegistry {
		result[index] = capability.Target.OS + "/" + capability.Target.Arch
	}
	return result
}

func HostMaskFor(host string) HostMask {
	switch host {
	case "windows":
		return HostWindows
	case "darwin":
		return HostDarwin
	case "linux":
		return HostLinux
	default:
		return 0
	}
}

func ToolchainMaskFor(name string) ToolchainMask {
	switch name {
	case "", "auto":
		return ToolchainAuto
	case "native":
		return ToolchainNative
	case "zig":
		return ToolchainZig
	case "docker":
		return ToolchainDocker
	default:
		return 0
	}
}

func RequiredHostName(mask HostMask) string {
	switch mask {
	case HostWindows:
		return "windows"
	case HostDarwin:
		return "darwin"
	case HostLinux:
		return "linux"
	default:
		return "supported"
	}
}
