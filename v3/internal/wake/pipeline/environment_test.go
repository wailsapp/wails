package pipeline

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestCurrentHostCapabilityProbeIsDeterministicAndCopiesFacts(t *testing.T) {
	root := t.TempDir()
	sdk := filepath.Join(root, "sdk")
	badNDK := filepath.Join(sdk, "ndk", "bad")
	goodNDK := filepath.Join(sdk, "ndk", "good")
	require.NoError(t, os.MkdirAll(goodNDK, 0o755))
	require.NoError(t, os.WriteFile(badNDK, []byte("not a directory"), 0o644))
	environment := map[string]string{"ANDROID_SDK_ROOT": sdk, "SIGNING_PASSWORD": "present"}
	operations := hostProbeOperations{
		hostOS: "linux", hostArch: "arm64",
		lookPath: func(name string) (string, error) {
			if name == "go" || name == "docker" || name == "xcrun" {
				return "/tools/" + name, nil
			}
			return "", fs.ErrNotExist
		},
		lookupEnv: func(name string) (string, bool) { value, ok := environment[name]; return value, ok },
		getenv:    func(name string) string { return environment[name] },
		stat:      os.Stat,
		glob:      func(string) ([]string, error) { return []string{badNDK, goodNDK}, nil },
		run: func(name string, arguments ...string) error {
			if name == "docker" || name == "xcrun" && !strings.Contains(strings.Join(arguments, " "), "iphonesimulator") {
				return nil
			}
			return errors.New("unavailable")
		},
	}
	host := currentHostCapabilitiesWithOperations([]string{"", "MISSING", "SIGNING_PASSWORD"}, operations)
	assert.Equal(t, "linux", host.os)
	assert.Equal(t, "arm64", host.arch)
	assert.True(t, host.hasTool("go"))
	assert.True(t, host.hasCredential("SIGNING_PASSWORD"))
	assert.True(t, host.androidSDK)
	assert.True(t, host.androidNDK)
	assert.True(t, host.hasDockerImage("wails-cross"))
	assert.True(t, host.hasAppleSDK("macosx"))
	assert.True(t, host.hasAppleSDK("iphoneos"))
	assert.False(t, host.hasAppleSDK("iphonesimulator"))
}

func TestCurrentHostCapabilitiesUsesTheRealProbeAdapterWithoutRequiringInstalledSDKs(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("PATH probing uses PATHEXT resolution on Windows")
	}
	tools := t.TempDir()
	for _, tool := range []string{"go", "docker", "xcrun"} {
		require.NoError(t, os.WriteFile(filepath.Join(tools, tool), []byte("#!/bin/sh\nexit 0\n"), 0o755))
	}
	t.Setenv("PATH", tools)
	t.Setenv("ANDROID_HOME", "")
	t.Setenv("ANDROID_SDK_ROOT", "")
	t.Setenv("ANDROID_NDK_HOME", "")
	t.Setenv("PIPELINE_TEST_CREDENTIAL", "present")
	host := CurrentHostCapabilities("PIPELINE_TEST_CREDENTIAL")
	assert.NotEmpty(t, host.os)
	assert.NotEmpty(t, host.arch)
	assert.True(t, host.hasTool("go"))
	assert.True(t, host.hasCredential("PIPELINE_TEST_CREDENTIAL"))
	assert.True(t, host.hasDockerImage("wails-cross"))
	assert.True(t, host.hasAppleSDK("iphonesimulator"))
}

func testHost(hostOS, hostArch string, extraTools ...string) HostCapabilities {
	tools := append([]string{"go", "npm", "cc"}, extraTools...)
	return NewHostCapabilitiesWithFacts(hostOS, hostArch, tools, nil, HostFacts{
		AndroidSDK: true, AndroidNDK: true,
		DockerImages: []string{"wails-cross"},
		AppleSDKs:    []string{"macosx", "iphoneos", "iphonesimulator"},
	})
}

func TestHostResolutionChoosesTheFastestAvailableCompatibleToolchain(t *testing.T) {
	config := testConfig(t)
	request := Request{Verb: "build", TargetOS: "windows", TargetArch: "amd64"}

	plan, err := PlanBuildForHost(config, request, testHost("linux", "amd64", "zig", "docker"))
	require.NoError(t, err)
	assert.Equal(t, "zig", plan.Nodes["target:windows/amd64:compile"].Spec.(CompileSpec).Toolchain)

	plan, err = PlanBuildForHost(config, request, testHost("linux", "amd64", "docker"))
	require.NoError(t, err)
	assert.Equal(t, "docker", plan.Nodes["target:windows/amd64:compile"].Spec.(CompileSpec).Toolchain)

	_, err = PlanBuildForHost(config, request, testHost("linux", "amd64"))
	assert.ErrorContains(t, err, "requires zig or docker")
}

func TestHostResolutionUsesTheInjectedHostAsTheAnonymousDefaultTarget(t *testing.T) {
	plan, err := PlanBuildForHost(testConfig(t), Request{Verb: "build"}, testHost("windows", "arm64"))
	require.NoError(t, err)
	assert.Equal(t, "windows/arm64", plan.Target)
	assert.Equal(t, "native", plan.Nodes["target:windows/arm64:compile"].Spec.(CompileSpec).Toolchain)
}

func TestHostResolutionEnforcesExplicitToolchains(t *testing.T) {
	config := testConfig(t)
	config.Targets.Linux.ARM64.Toolchain = "native"
	_, err := PlanBuildForHost(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "arm64"}, testHost("darwin", "arm64"))
	assert.ErrorContains(t, err, `toolchain "native" cannot build linux/arm64 on darwin/arm64`)

	config.Targets.Linux.ARM64.Toolchain = "zig"
	_, err = PlanBuildForHost(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "arm64"}, testHost("darwin", "arm64"))
	assert.ErrorContains(t, err, `toolchain "zig" requires tool "zig"`)

	config.Targets.Linux.ARM64.Toolchain = "docker"
	_, err = PlanBuildForHost(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "arm64"}, testHost("windows", "amd64", "docker"))
	assert.ErrorContains(t, err, `toolchain "docker" cannot build linux/arm64 on windows/amd64`)

	android := configForAndroid(t)
	android.Targets.Android.ARM64.Toolchain = "docker"
	_, err = PlanBuildForHost(android, Request{Verb: "build", TargetOS: "android", TargetArch: "arm64"}, testHost("linux", "amd64", "docker"))
	assert.ErrorContains(t, err, `toolchain "docker" is not supported for target android/arm64`)
}

func TestHostResolutionRejectsUnavailablePackageOperations(t *testing.T) {
	config := testConfig(t)
	_, err := PlanBuildForHost(config, Request{Verb: "build", TargetOS: "darwin", TargetArch: "arm64", Formats: []string{"dmg"}}, testHost("linux", "amd64", "docker", "hdiutil"))
	assert.ErrorContains(t, err, "dmg packaging for darwin/arm64 requires a darwin host")

	_, err = PlanBuildForHost(config, Request{Verb: "build", TargetOS: "windows", TargetArch: "amd64", Formats: []string{"nsis"}}, testHost("windows", "amd64"))
	assert.ErrorContains(t, err, `nsis packaging for windows/amd64 requires tool "makensis"`)
}

func TestHostResolutionRejectsUnavailableSigningAndNotarization(t *testing.T) {
	config := testConfig(t)
	_, err := PlanBuildForHost(config, Request{Verb: "sign", TargetOS: "linux", TargetArch: "amd64", Formats: []string{"deb"}}, testHost("linux", "amd64"))
	assert.ErrorContains(t, err, "signing is not enabled for linux")

	config.Signing.Darwin.Enabled = true
	config.Signing.Darwin.Identity = "Developer ID Application: Example"
	config.Selected.Name = "release"
	config.Selected.Targets = []manifest.ProfileTarget{{Target: "darwin/arm64", Formats: []string{"dmg"}, Sign: true, Notarize: true}}
	_, err = PlanBuildForHost(config, Request{Verb: "build"}, testHost("darwin", "arm64", "hdiutil", "codesign", "xcrun"))
	assert.ErrorContains(t, err, "notarization requires signing.darwin.notarization credential")
}

func TestHostResolutionValidatesFormatSpecificSigningRequirements(t *testing.T) {
	config := testConfig(t)
	config.Signing.Linux.Enabled = true
	request := Request{Verb: "sign", TargetOS: "linux", TargetArch: "amd64", Formats: []string{"deb"}}
	_, err := PlanBuildForHost(config, request, testHost("linux", "amd64"))
	assert.ErrorContains(t, err, "PGP key identifier")

	config.Signing.Linux.Certificate = "release@example.com"
	_, err = PlanBuildForHost(config, request, testHost("linux", "amd64"))
	assert.ErrorContains(t, err, `requires tool "dpkg-sig"`)
	_, err = PlanBuildForHost(config, request, testHost("linux", "amd64", "dpkg-sig"))
	require.NoError(t, err)

	config.Signing.Windows.Enabled = true
	_, err = PlanBuildForHost(config, Request{Verb: "sign", TargetOS: "windows", TargetArch: "amd64", Formats: []string{"msix"}}, testHost("windows", "amd64", "MakeAppx.exe", "signtool.exe"))
	assert.ErrorContains(t, err, "certificate or thumbprint")
}

func TestHostResolutionRejectsMissingBuildExecutables(t *testing.T) {
	config := testConfig(t)
	request := Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"}
	_, err := PlanBuildForHost(config, request, NewHostCapabilities("linux", "amd64", []string{"npm"}, nil))
	assert.ErrorContains(t, err, `requires tool "go"`)

	_, err = PlanBuildForHost(config, request, NewHostCapabilities("linux", "amd64", []string{"go"}, nil))
	assert.ErrorContains(t, err, `requires tool "npm"`)

	config.Build.Obfuscation = true
	_, err = PlanBuildForHost(config, request, NewHostCapabilities("linux", "amd64", []string{"go", "npm"}, nil))
	assert.ErrorContains(t, err, `requires tool "garble"`)
}

func TestHostResolutionRejectsMissingToolchainRuntimeFacts(t *testing.T) {
	config := testConfig(t)
	config.Targets.Windows.AMD64.Toolchain = "docker"
	_, err := PlanBuildForHost(config, Request{Verb: "build", TargetOS: "windows", TargetArch: "amd64"}, NewHostCapabilities("linux", "amd64", []string{"go", "npm", "docker"}, nil))
	assert.ErrorContains(t, err, `requires Docker image "wails-cross"`)

	host := NewHostCapabilitiesWithFacts("linux", "amd64", []string{"go", "npm", "java"}, nil, HostFacts{})
	_, err = PlanBuildForHost(configForAndroid(t), Request{Verb: "build", TargetOS: "android", TargetArch: "arm64"}, host)
	assert.ErrorContains(t, err, "Android SDK")

	host = NewHostCapabilitiesWithFacts("linux", "amd64", []string{"go", "npm", "java"}, nil, HostFacts{AndroidSDK: true})
	_, err = PlanBuildForHost(configForAndroid(t), Request{Verb: "build", TargetOS: "android", TargetArch: "arm64"}, host)
	assert.ErrorContains(t, err, "Android NDK")
}

func TestHostResolutionRejectsUnavailableAppleSDK(t *testing.T) {
	config := testConfig(t)
	config.Selected = manifest.Profile{Name: "release", Targets: []manifest.ProfileTarget{{Target: "ios/arm64", Destination: "simulator"}}}
	host := NewHostCapabilitiesWithFacts("darwin", "arm64", []string{"go", "npm", "xcrun", "codesign"}, nil, HostFacts{AppleSDKs: []string{"iphoneos"}})
	_, err := PlanBuildForHost(config, Request{Verb: "build"}, host)
	assert.ErrorContains(t, err, `requires Apple SDK "iphonesimulator"`)
}

func TestHostAndToolchainValidationCoversEveryFailureContract(t *testing.T) {
	config := testConfig(t)
	_, err := PlanBuildForHost(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"}, NewHostCapabilities("plan9", "amd64", []string{"go", "npm", "cc"}, nil))
	assert.ErrorContains(t, err, "unsupported build host")
	_, err = PlanBuildForHost(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"}, NewHostCapabilities("linux", "", []string{"go", "npm", "cc"}, nil))
	assert.ErrorContains(t, err, "unsupported build host")

	assert.NoError(t, requireCommandTool(HostCapabilities{}, "custom command", "./tools/frontend"))
	assert.NoError(t, requireCommandTool(HostCapabilities{}, "empty command", ""))
	assert.Equal(t, "npm", firstCommand(nil, "npm"))
	assert.Equal(t, "pnpm", firstCommand([]string{"pnpm", "run", "build"}, "npm"))

	nativeLinux := CompileSpec{TargetOS: "linux", TargetArch: "amd64", Toolchain: "native"}
	err = validateCompileEnvironment(nativeLinux, NewHostCapabilities("linux", "amd64", []string{"go"}, nil))
	assert.ErrorContains(t, err, "requires a C compiler")
	assert.NoError(t, validateCompileEnvironment(nativeLinux, NewHostCapabilities("linux", "amd64", []string{"go", "gcc"}, nil)))

	docker := CompileSpec{TargetOS: "windows", TargetArch: "amd64", Toolchain: "docker"}
	err = validateCompileEnvironment(docker, NewHostCapabilities("linux", "amd64", []string{"go", "docker"}, nil))
	assert.ErrorContains(t, err, "Docker image")
	hostWithImage := NewHostCapabilitiesWithFacts("linux", "amd64", []string{"go", "docker"}, nil, HostFacts{DockerImages: []string{"wails-cross"}})
	assert.NoError(t, validateCompileEnvironment(docker, hostWithImage))

	_, err = resolveToolchain(CompileSpec{TargetOS: "windows", TargetArch: "amd64", Toolchain: "docker"}, NewHostCapabilities("linux", "amd64", []string{"go"}, nil))
	assert.ErrorContains(t, err, `requires tool "docker"`)
	_, err = resolveToolchain(CompileSpec{TargetOS: "plan9", TargetArch: "amd64", Toolchain: "docker"}, NewHostCapabilities("linux", "amd64", []string{"go", "docker"}, nil))
	assert.ErrorContains(t, err, "cannot build")
	assert.False(t, dockerToolchainSupports(CompileSpec{TargetOS: "plan9", TargetArch: "amd64"}, testHost("linux", "amd64", "docker")))
	assert.False(t, nativeToolchainSupports(CompileSpec{TargetOS: "android", TargetArch: "arm64"}, HostCapabilities{os: "plan9", arch: "amd64"}))
	androidCompile := CompileSpec{TargetOS: "android", TargetArch: "arm64", Toolchain: "native"}
	err = validateCompileEnvironment(androidCompile, NewHostCapabilitiesWithFacts("linux", "amd64", []string{"go"}, nil, HostFacts{}))
	assert.ErrorContains(t, err, "Android SDK")
}

func TestHostPlanningReportsInstallAndIOSAssemblyToolFailures(t *testing.T) {
	config := testConfig(t)
	config.Frontend.Install = []string{"pnpm", "install"}
	_, err := PlanBuildForHost(config, Request{Verb: "build", TargetOS: "linux", TargetArch: "amd64"}, testHost("linux", "amd64"))
	assert.ErrorContains(t, err, `frontend dependency installation requires tool "pnpm"`)

	config = testConfig(t)
	config.Selected = manifest.Profile{Name: "simulator", Targets: []manifest.ProfileTarget{{Target: "ios/arm64", Destination: "simulator"}}}
	host := NewHostCapabilitiesWithFacts("darwin", "arm64", []string{"go", "npm", "xcrun"}, nil, HostFacts{AppleSDKs: []string{"iphonesimulator"}})
	_, err = PlanBuildForHost(config, Request{Verb: "build"}, host)
	assert.ErrorContains(t, err, `iOS application assembly for ios/arm64 requires tool "codesign"`)
}

func TestSigningHostValidationCoversEveryPlatformContract(t *testing.T) {
	assert.ErrorContains(t, validateSigningHost(SignSpec{TargetOS: "darwin", Format: "dmg", Config: manifest.SigningPlatform{Enabled: true, Identity: "Developer ID"}}, testHost("linux", "amd64", "codesign")), "requires a darwin host")
	assert.ErrorContains(t, validateSigningHost(SignSpec{TargetOS: "darwin", Format: "dmg", Config: manifest.SigningPlatform{Enabled: true}}, testHost("darwin", "arm64", "codesign")), "identity")
	assert.ErrorContains(t, validateSigningHost(SignSpec{TargetOS: "darwin", Format: "dmg", Config: manifest.SigningPlatform{Enabled: true, Identity: "Developer ID"}}, testHost("darwin", "arm64")), "codesign")
	assert.ErrorContains(t, validateSigningHost(SignSpec{TargetOS: "darwin", Format: "dmg", Config: manifest.SigningPlatform{Enabled: true, Identity: "Developer ID", Notarize: true, NotarizationCredential: "NOTARY"}}, testHost("darwin", "arm64", "codesign")), "xcrun")
	assert.ErrorContains(t, validateSigningHost(SignSpec{TargetOS: "darwin", Format: "app", Config: manifest.SigningPlatform{Enabled: true, Identity: "Developer ID", Notarize: true, NotarizationCredential: "NOTARY"}}, testHost("darwin", "arm64", "codesign", "xcrun")), "ditto")

	assert.ErrorContains(t, validateSigningHost(SignSpec{TargetOS: "ios", Format: "ipa", Config: manifest.SigningPlatform{Enabled: true, Identity: "Apple Distribution"}}, testHost("linux", "amd64", "codesign")), "requires a darwin host")
	assert.ErrorContains(t, validateSigningHost(SignSpec{TargetOS: "ios", Format: "ipa", Config: manifest.SigningPlatform{Enabled: true}}, testHost("darwin", "arm64", "codesign")), "identity")
	assert.NoError(t, validateSigningHost(SignSpec{TargetOS: "ios", Format: "ipa", Config: manifest.SigningPlatform{Enabled: true, Identity: "Apple Distribution"}}, testHost("darwin", "arm64", "codesign")))

	assert.ErrorContains(t, validateSigningHost(SignSpec{TargetOS: "windows", Format: "msix", Config: manifest.SigningPlatform{Enabled: true, Certificate: "app.pfx"}}, testHost("linux", "amd64", "signtool.exe")), "requires a windows host")
	assert.NoError(t, validateSigningHost(SignSpec{TargetOS: "windows", Format: "msix", Config: manifest.SigningPlatform{Enabled: true, Thumbprint: "abc"}}, testHost("windows", "amd64", "signtool.exe")))

	android := manifest.SigningPlatform{Enabled: true, Certificate: "upload.jks", KeyAlias: "upload", Credential: "ANDROID_PASSWORD"}
	assert.ErrorContains(t, validateSigningHost(SignSpec{TargetOS: "android", Format: "aab", Config: manifest.SigningPlatform{Enabled: true}}, testHost("linux", "amd64", "jarsigner")), "certificate, key_alias, and credential")
	assert.ErrorContains(t, validateSigningHost(SignSpec{TargetOS: "android", Format: "aab", Config: android}, testHost("linux", "amd64", "jarsigner")), "credential")
	androidHost := NewHostCapabilities("linux", "amd64", []string{"jarsigner", "apksigner"}, []string{"ANDROID_PASSWORD"})
	assert.NoError(t, validateSigningHost(SignSpec{TargetOS: "android", Format: "aab", Config: android}, androidHost))
	assert.NoError(t, validateSigningHost(SignSpec{TargetOS: "android", Format: "apk", Config: android}, androidHost))

	linux := manifest.SigningPlatform{Enabled: true, Certificate: "release@example.com"}
	assert.NoError(t, validateSigningHost(SignSpec{TargetOS: "linux", Format: "rpm", Config: linux}, testHost("linux", "amd64", "rpmsign")))
	assert.ErrorContains(t, validateSigningHost(SignSpec{TargetOS: "linux", Format: "appimage", Config: linux}, testHost("linux", "amd64")), "not supported")
	assert.NoError(t, validateSigningHost(SignSpec{TargetOS: "unknown", Config: manifest.SigningPlatform{Enabled: true}}, HostCapabilities{}))
}

func TestHostCapabilityCollectionHelpersAreDeterministic(t *testing.T) {
	assert.Empty(t, uniqueSorted(nil))
	assert.Equal(t, []string{"a", "b"}, uniqueSorted([]string{"b", "a", "b"}))
	assert.False(t, containsSorted(nil, "a"))
	assert.Empty(t, firstNonemptyEnvironment(func(string) string { return "" }, "first", "second"))
}

func configForAndroid(t *testing.T) manifest.Config {
	t.Helper()
	return testConfig(t)
}
