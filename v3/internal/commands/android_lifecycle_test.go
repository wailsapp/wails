package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAndroidActivityRecreationKeepsGoAppRunning(t *testing.T) {
	// Given
	mainActivity, err := buildAssets.ReadFile("build_assets/android/app/src/main/java/com/wails/app/MainActivity.java")
	require.NoError(t, err)

	mainActivityJava := string(mainActivity)

	// Then: a recreated Activity (configuration or theme overlay change) must
	// not shut down the Go app that the new Activity reattaches to.
	assert.Contains(t, mainActivityJava, "if (bridge != null && isFinishing() && !isChangingConfigurations()) {")
}
