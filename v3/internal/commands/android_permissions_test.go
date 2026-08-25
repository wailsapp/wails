package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAndroidPermissionGatedActionsResumeAfterGrant(t *testing.T) {
	// Given
	mainActivity, err := buildAssets.ReadFile("build_assets/android/app/src/main/java/com/wails/app/MainActivity.java")
	require.NoError(t, err)
	wailsBridge, err := buildAssets.ReadFile("build_assets/android/app/src/main/java/com/wails/app/WailsBridge.java")
	require.NoError(t, err)

	mainActivityJava := string(mainActivity)
	wailsBridgeJava := string(wailsBridge)

	// Then
	assert.NotContains(t, mainActivityJava, "tap again once granted")
	assert.Contains(t, mainActivityJava, "public void onRequestPermissionsResult")
	assert.Contains(t, mainActivityJava, "if (requestCode == CAMERA_PERMISSION_REQUEST)")
	assert.Contains(t, mainActivityJava, "launchCameraCapture(pendingCaptureIsVideo)")
	assert.Contains(t, mainActivityJava, "camera permission denied")
	assert.Contains(t, mainActivityJava, "bridge.onRequestPermissionsResult(requestCode, grantResults)")

	assert.NotContains(t, wailsBridgeJava, "tap again once granted")
	assert.Contains(t, wailsBridgeJava, "private static final int LOCATION_PERMISSION_REQUEST = 1002")
	assert.Contains(t, wailsBridgeJava, "private boolean pendingLocationRequest")
	assert.Contains(t, wailsBridgeJava, "pendingLocationRequest = true")
	assert.Contains(t, wailsBridgeJava, "pendingLocationRequest = false")
	assert.Contains(t, wailsBridgeJava, "public void onRequestPermissionsResult")
	assert.Contains(t, wailsBridgeJava, "if (requestCode == LOCATION_PERMISSION_REQUEST)")
	assert.Contains(t, wailsBridgeJava, "getLocation()")
	assert.Contains(t, wailsBridgeJava, "location permission denied")
}
