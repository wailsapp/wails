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
