package pipeline

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestAnonymousFormatsResolveToTheirCompatibleTargets(t *testing.T) {
	outcomes, err := resolveBuildOutcomes(manifest.Config{}, Request{
		Verb:    "build",
		Targets: []Target{{OS: "windows", Arch: "amd64"}, {OS: "linux", Arch: "arm64"}},
		Formats: []string{"nsis", "deb"},
	})
	require.NoError(t, err)
	require.Len(t, outcomes, 2)
	assert.Equal(t, buildOutcome{target: Target{OS: "linux", Arch: "arm64"}, formats: []string{"deb"}}, outcomes[0])
	assert.Equal(t, buildOutcome{target: Target{OS: "windows", Arch: "amd64"}, formats: []string{"nsis"}}, outcomes[1])
}

func TestAnonymousFormatsRejectUnmatchedFormatsAndTargets(t *testing.T) {
	_, err := resolveBuildOutcomes(manifest.Config{}, Request{
		Targets: []Target{{OS: "linux", Arch: "amd64"}},
		Formats: []string{"deb", "nsis"},
	})
	assert.ErrorContains(t, err, `format "nsis" is not supported for any selected target`)

	_, err = resolveBuildOutcomes(manifest.Config{}, Request{
		Targets: []Target{{OS: "linux", Arch: "amd64"}, {OS: "windows", Arch: "amd64"}},
		Formats: []string{"deb"},
	})
	assert.ErrorContains(t, err, "target windows/amd64 receives no compatible format")
}

func TestConfiguredPackageFormatsAreValidatedForAnonymousPackageRequests(t *testing.T) {
	config := manifest.Config{}
	config.Package.Linux.Formats = []string{"unknown"}
	_, err := resolveBuildOutcomes(config, Request{Verb: "package", TargetOS: "linux", TargetArch: "amd64"})
	assert.ErrorContains(t, err, `target linux/amd64: unknown package format "unknown"`)
}

func TestProfileAndAnonymousRequestsShareOutcomeSemantics(t *testing.T) {
	profile := manifest.Config{Selected: manifest.Profile{Name: "release", Targets: []manifest.ProfileTarget{
		{Target: "linux/arm64", Formats: []string{"rpm", "deb"}},
		{Target: "windows/amd64", Formats: []string{"nsis"}, Sign: true},
	}}}
	profileOutcomes, err := resolveBuildOutcomes(profile, Request{Verb: "build"})
	require.NoError(t, err)
	assert.Equal(t, []buildOutcome{
		{target: Target{OS: "linux", Arch: "arm64"}, formats: []string{"deb", "rpm"}},
		{target: Target{OS: "windows", Arch: "amd64"}, formats: []string{"nsis"}, sign: true},
	}, profileOutcomes)

	_, err = resolveBuildOutcomes(profile, Request{Verb: "build", Formats: []string{"deb"}})
	assert.ErrorContains(t, err, "complete build request")
}

func TestProductionAndroidDefaultsToAABAndRejectsAPK(t *testing.T) {
	outcomes, err := resolveBuildOutcomes(manifest.Config{}, Request{Verb: "build", TargetOS: "android", TargetArch: "universal"})
	require.NoError(t, err)
	assert.Equal(t, []buildOutcome{{target: Target{OS: "android", Arch: "universal"}, formats: []string{"aab"}}}, outcomes)

	_, err = resolveBuildOutcomes(manifest.Config{}, Request{Verb: "build", TargetOS: "android", TargetArch: "arm64", Formats: []string{"apk"}})
	assert.ErrorContains(t, err, "production APK is no longer supported")
}

func TestIOSProfileRequiresDestinationAndIPARequiresDevice(t *testing.T) {
	config := manifest.Config{Selected: manifest.Profile{Name: "store", Targets: []manifest.ProfileTarget{{Target: "ios/arm64", Formats: []string{"ipa"}}}}}
	_, err := resolveBuildOutcomes(config, Request{Verb: "build"})
	assert.ErrorContains(t, err, "destination")

	config.Selected.Targets[0].Destination = "simulator"
	_, err = resolveBuildOutcomes(config, Request{Verb: "build"})
	assert.ErrorContains(t, err, "requires destination = \"device\"")

	config.Selected.Targets[0].Destination = "device"
	_, err = resolveBuildOutcomes(config, Request{Verb: "build"})
	require.NoError(t, err)

	config.Selected.Targets[0] = manifest.ProfileTarget{Target: "darwin/arm64", Destination: "device"}
	_, err = resolveBuildOutcomes(config, Request{Verb: "build"})
	assert.ErrorContains(t, err, "destination is only valid for ios/arm64")
}

func TestBuildOutcomeSelectionRejectsEveryInvalidRequestShape(t *testing.T) {
	tests := []struct {
		name    string
		config  manifest.Config
		request Request
		message string
	}{
		{name: "empty profile", config: manifest.Config{Selected: manifest.Profile{Name: "release"}}, message: "requires at least one target"},
		{name: "invalid profile target", config: manifest.Config{Selected: manifest.Profile{Name: "release", Targets: []manifest.ProfileTarget{{Target: "linux"}}}}, message: "invalid target"},
		{name: "duplicate profile target", config: manifest.Config{Selected: manifest.Profile{Name: "release", Targets: []manifest.ProfileTarget{{Target: "linux/amd64"}, {Target: "linux/amd64"}}}}, message: "duplicate target"},
		{name: "unsupported profile target", config: manifest.Config{Selected: manifest.Profile{Name: "release", Targets: []manifest.ProfileTarget{{Target: "plan9/amd64"}}}}, message: "supported targets"},
		{name: "android profile without AAB", config: manifest.Config{Selected: manifest.Profile{Name: "release", Targets: []manifest.ProfileTarget{{Target: "android/arm64"}}}}, message: "must select the aab"},
		{name: "notarize non-darwin", config: manifest.Config{Selected: manifest.Profile{Name: "release", Targets: []manifest.ProfileTarget{{Target: "linux/amd64", Notarize: true, Sign: true}}}}, message: "cannot be notarized"},
		{name: "notarize unsigned", config: manifest.Config{Selected: manifest.Profile{Name: "release", Targets: []manifest.ProfileTarget{{Target: "darwin/arm64", Notarize: true}}}}, message: "must be signed"},
		{name: "duplicate anonymous target", request: Request{Targets: []Target{{OS: "linux", Arch: "amd64"}, {OS: "linux", Arch: "amd64"}}}, message: "duplicate target"},
		{name: "unsupported anonymous target", request: Request{TargetOS: "plan9", TargetArch: "amd64"}, message: "supported targets"},
		{name: "unsupported verb", request: Request{Verb: "deploy", TargetOS: "linux", TargetArch: "amd64"}, message: "unsupported pipeline verb"},
		{name: "empty format", request: Request{TargetOS: "linux", TargetArch: "amd64", Formats: []string{" "}}, message: "format cannot be empty"},
		{name: "duplicate format", request: Request{TargetOS: "linux", TargetArch: "amd64", Formats: []string{"deb", "deb"}}, message: "duplicate package format"},
		{name: "unknown format", request: Request{TargetOS: "linux", TargetArch: "amd64", Formats: []string{"unknown"}}, message: "unknown package format"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := resolveBuildOutcomes(test.config, test.request)
			require.ErrorContains(t, err, test.message)
		})
	}
}

func TestFormatAndTargetParsersCoverProductionAndDevelopmentFailures(t *testing.T) {
	linux, ok := lookupTarget("linux", "amd64")
	require.True(t, ok)
	_, err := resolveFormatsForTarget(linux, []string{"unknown"}, false)
	assert.ErrorContains(t, err, "unknown package format")
	_, err = resolveFormatsForTarget(linux, []string{"apk"}, false)
	assert.ErrorContains(t, err, "production APK")
	_, err = resolveFormatsForTarget(linux, []string{"apk"}, true)
	assert.ErrorContains(t, err, "in development")
	_, err = resolveFormatsForTarget(linux, []string{"deb", "deb"}, false)
	assert.ErrorContains(t, err, "duplicate package format")

	for _, value := range []string{"", "linux", "/amd64", "linux/", "linux/amd64/extra"} {
		_, err := parseTargetName(value)
		assert.Error(t, err, value)
	}
}

func TestAnonymousSelectionDefaultsEachPartialTargetAndUsesConfiguredFormats(t *testing.T) {
	outcomes, err := resolveAnonymousOutcomes(manifest.Config{}, Request{Targets: []Target{{}}})
	require.NoError(t, err)
	require.Len(t, outcomes, 1)
	assert.NotEmpty(t, outcomes[0].target.OS)
	assert.NotEmpty(t, outcomes[0].target.Arch)

	config := manifest.Config{}
	config.Package.Linux.Formats = []string{"deb"}
	outcomes, err = resolveAnonymousOutcomes(config, Request{Verb: "package", TargetOS: "linux", TargetArch: "amd64"})
	require.NoError(t, err)
	assert.Equal(t, []string{"deb"}, outcomes[0].formats)
}
