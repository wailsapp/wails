package commands

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAndroidOpenURLReturnsLaunchFailures(t *testing.T) {
	generated, err := buildAssets.ReadFile("build_assets/android/app/src/main/java/com/wails/app/WailsBridge.java")
	require.NoError(t, err)
	example, err := os.ReadFile("../../examples/mobile/build/android/app/src/main/java/com/wails/app/WailsBridge.java")
	require.NoError(t, err)

	sources := map[string][]byte{
		"generated host": generated,
		"mobile example": example,
	}
	for name, source := range sources {
		t.Run(name, func(t *testing.T) {
			bridge := string(source)
			assert.Contains(t, bridge, "public String openURL(final String url)")
			assert.Contains(t, bridge, "new URI(url)")
			assert.Contains(t, bridge, "FutureTask<String>")
			assert.Contains(t, bridge, "launch.get(30, TimeUnit.SECONDS)")
			assert.Contains(t, bridge, "launch.cancel(false)")
			assert.NotContains(t, bridge, `Log.e(TAG, "openURL failed", e)`)
		})
	}
}
