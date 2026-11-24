package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
)

func TestOpenHarmonyManifestPlatformsAndABI(t *testing.T) {
	targets := manifest.Targets{}
	platform := manifestPlatform(targets, "openharmony")
	require.NotNil(t, platform)
	assert.Equal(t, &targets.OpenHarmony, platform)
	assert.Nil(t, manifestPlatform(targets, "plan9"))

	signing := manifest.Signing{}
	signingPlatform := manifestSigningPlatform(signing, "openharmony")
	require.NotNil(t, signingPlatform)
	assert.Equal(t, &signing.OpenHarmony, signingPlatform)
	assert.Nil(t, manifestSigningPlatform(signing, "plan9"))

	for _, test := range []struct {
		arch string
		want string
	}{
		{arch: "amd64", want: "x86_64"},
		{arch: "arm64", want: "arm64-v8a"},
	} {
		t.Run(test.arch, func(t *testing.T) {
			// Pure mapping: no DEVECO_SDK_HOME or DevEco toolchain is required.
			got, err := openHarmonyABI(test.arch)
			require.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}

	_, err := openHarmonyABI("riscv64")
	assert.ErrorContains(t, err, "unsupported OpenHarmony architecture")
}

func TestHarmonyOSManifestPlatforms(t *testing.T) {
	targets := manifest.Targets{}
	platform := manifestPlatform(targets, "harmonyos")
	require.NotNil(t, platform)
	assert.Equal(t, &targets.HarmonyOS, platform)
	assert.Nil(t, manifestPlatform(targets, "plan9"))

	signing := manifest.Signing{}
	signingPlatform := manifestSigningPlatform(signing, "harmonyos")
	require.NotNil(t, signingPlatform)
	assert.Equal(t, &signing.HarmonyOS, signingPlatform)
	assert.Nil(t, manifestSigningPlatform(signing, "plan9"))
}
