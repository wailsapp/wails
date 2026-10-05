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

	wailsBridge, err := buildAssets.ReadFile("build_assets/android/app/src/main/java/com/wails/app/WailsBridge.java")
	require.NoError(t, err)

	mainActivityJava := string(mainActivity)
	wailsBridgeJava := string(wailsBridge)

	// Then: a recreated Activity (configuration or theme overlay change) must
	// not shut down the Go app that the new Activity reattaches to.
	assert.Contains(t, mainActivityJava, "if (bridge != null && isFinishing() && !isChangingConfigurations()) {")

	// Then: a crashed WebView renderer recreates the Activity instead of
	// taking the whole app down.
	assert.Contains(t, mainActivityJava, "public boolean onRenderProcessGone(WebView view, RenderProcessGoneDetail detail)")
	assert.Contains(t, mainActivityJava, "recreate();")

	// Then: renderer-crash recovery is bounded so a page that crashes the
	// renderer on every load cannot loop forever.
	assert.Contains(t, mainActivityJava, "if (renderCrashCount > MAX_RENDER_CRASH_RECOVERIES) {")

	// Then: pending picker/camera requests outlive a recreated Activity.
	assert.Contains(t, mainActivityJava, "private static int pendingFilePickerCallbackID = -1;")
	assert.Contains(t, mainActivityJava, "private static File pendingCaptureFile;")

	// Then: window state requested by the still-running Go app (screen-capture
	// protection, keep-awake, brightness) is re-applied to the new window.
	assert.Contains(t, mainActivityJava, "bridge.restoreWindowState();")
	assert.Contains(t, wailsBridgeJava, "screenProtectWanted = enabled != 0;")
	assert.Contains(t, wailsBridgeJava, "activity.getWindow().addFlags(WindowManager.LayoutParams.FLAG_SECURE);")
}
