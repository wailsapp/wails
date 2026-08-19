package pipeline

import (
	"fmt"
	"testing"

	"github.com/wailsapp/wails/v3/internal/wake/buildinfo"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestExhaustiveProductionCapabilityMatrix(t *testing.T) {
	base := testConfig(t)
	formats := []string{""}
	for _, format := range buildinfo.FormatCapabilities() {
		formats = append(formats, format.Name)
	}
	hosts := []Target{{OS: "windows", Arch: "amd64"}, {OS: "darwin", Arch: "arm64"}, {OS: "linux", Arch: "amd64"}}
	toolchains := []string{"auto", "native", "zig", "docker"}
	allTools := []string{"go", "npm", "cc", "zig", "docker", "makensis", "MakeAppx.exe", "signtool.exe", "hdiutil", "xcrun", "codesign", "ditto", "zip", "java", "jarsigner", "apksigner", "dpkg-sig", "rpmsign"}
	facts := HostFacts{AndroidSDK: true, AndroidNDK: true, DockerImages: []string{"wails-cross"}, AppleSDKs: []string{"macosx", "iphoneos", "iphonesimulator"}}
	cases := 0
	for _, target := range buildinfo.TargetCapabilities() {
		for _, format := range formats {
			for _, host := range hosts {
				for _, toolchain := range toolchains {
					for _, sign := range []bool{false, true} {
						for _, notarize := range []bool{false, true} {
							cases++
							config := base
							setMatrixToolchain(&config.Targets, target.Target, toolchain)
							configureMatrixSigning(&config.Signing)
							selected := manifest.ProfileTarget{Target: target.Target.OS + "/" + target.Target.Arch, Sign: sign, Notarize: notarize}
							if format != "" {
								selected.Formats = []string{format}
							}
							if target.Target.OS == "ios" {
								selected.Destination = "simulator"
								if format == "ipa" {
									selected.Destination = "device"
								}
							}
							config.Selected = manifest.Profile{Name: "matrix", Targets: []manifest.ProfileTarget{selected}}
							hostCapabilities := NewHostCapabilitiesWithFacts(host.OS, host.Arch, allTools, []string{"ANDROID_PASSWORD"}, facts)
							_, err := PlanBuildForHost(config, Request{Verb: "build"}, hostCapabilities)
							wantSuccess := matrixCombinationSupported(target, format, host, toolchain, sign, notarize)
							if (err == nil) != wantSuccess {
								t.Fatalf("matrix mismatch target=%s/%s format=%q host=%s/%s toolchain=%s sign=%t notarize=%t want_success=%t: %v", target.Target.OS, target.Target.Arch, format, host.OS, host.Arch, toolchain, sign, notarize, wantSuccess, err)
							}
						}
					}
				}
			}
		}
	}
	if cases != len(buildinfo.TargetCapabilities())*len(formats)*len(hosts)*len(toolchains)*4 {
		t.Fatalf("matrix generated %d cases", cases)
	}
}

func matrixCombinationSupported(target buildinfo.TargetCapability, format string, host Target, toolchain string, sign, notarize bool) bool {
	if notarize && (!sign || target.Target.OS != "darwin") {
		return false
	}
	if format != "" {
		formatCapability, ok := buildinfo.LookupFormat(format)
		if !ok || !formatCapability.Production || !target.SupportsFormat(format, false) || !formatCapability.SupportsHost(host.OS) {
			return false
		}
	} else if target.Runnable == buildinfo.RunnableNone {
		return false
	}
	if !target.SupportsToolchain(toolchain) || !matrixToolchainAvailable(target.Target, host, toolchain) {
		return false
	}
	if target.Target.OS == "ios" && host.OS != "darwin" {
		return false
	}
	if !sign {
		return true
	}
	switch target.Target.OS {
	case "windows":
		return host.OS == "windows"
	case "darwin", "ios":
		return host.OS == "darwin"
	case "android":
		return true
	case "linux":
		return format == "deb" || format == "rpm"
	default:
		return false
	}
}

func matrixToolchainAvailable(target, host Target, toolchain string) bool {
	native := target.OS == "android" || target.OS == "ios" && host.OS == "darwin" || target.OS == host.OS && (target.Arch == host.Arch || target.OS == "darwin")
	switch toolchain {
	case "native":
		return native
	case "zig":
		return target.OS == "windows" || target.OS == "linux"
	case "docker":
		return host.OS != "windows" && target.OS != "ios" && target.OS != "android"
	case "auto":
		return native || target.OS == "windows" || target.OS == "linux" || host.OS != "windows" && target.OS == "darwin"
	default:
		panic(fmt.Sprintf("unexpected matrix toolchain %q", toolchain))
	}
}

func setMatrixToolchain(targets *manifest.Targets, target Target, toolchain string) {
	platform := map[string]*manifest.Platform{
		"windows": &targets.Windows, "darwin": &targets.Darwin, "linux": &targets.Linux, "ios": &targets.IOS, "android": &targets.Android,
	}[target.OS]
	settings := map[string]*manifest.Target{
		"amd64": &platform.AMD64, "arm64": &platform.ARM64, "universal": &platform.Universal,
	}[target.Arch]
	settings.Toolchain = toolchain
}

func configureMatrixSigning(signing *manifest.Signing) {
	signing.Windows = manifest.SigningPlatform{Enabled: true, Certificate: "windows.pfx"}
	signing.Darwin = manifest.SigningPlatform{Enabled: true, Identity: "Developer ID", NotarizationCredential: "notary"}
	signing.Linux = manifest.SigningPlatform{Enabled: true, Certificate: "release@example.com"}
	signing.IOS = manifest.SigningPlatform{Enabled: true, Identity: "Apple Distribution"}
	signing.Android = manifest.SigningPlatform{Enabled: true, Certificate: "android.keystore", KeyAlias: "upload", Credential: "ANDROID_PASSWORD"}
}

func TestDevelopmentAPKMatrixIsSeparateFromProductionFormats(t *testing.T) {
	for _, target := range buildinfo.TargetCapabilities() {
		if target.Target.OS != "android" {
			continue
		}
		config := testConfig(t)
		plan, err := PlanBuild(config, Request{Verb: "build", TargetOS: "android", TargetArch: target.Target.Arch, Formats: []string{"apk"}, Development: true})
		if err != nil {
			t.Fatalf("development APK for %s/%s: %v", target.Target.OS, target.Target.Arch, err)
		}
		if _, ok := plan.Nodes[NodeKey("package:android/"+target.Target.Arch+":apk")]; !ok {
			t.Fatalf("development APK node missing for %s/%s", target.Target.OS, target.Target.Arch)
		}
	}
}
