//go:build android

package application

// Genuinely-mobile native capabilities for Android, mirroring the iOS surface
// in mobile_features_ios.go. Each call is forwarded to a matching method on the
// Java WailsBridge via the reflective bridge helpers. Asynchronous results
// (biometric prompts, torch availability, …) come back to the frontend as
// custom events through the nativeEmitEvent JNI export.

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// --- Phase A: one-way actions -------------------------------------------------

// AndroidShare presents the Android share chooser (Intent.ACTION_SEND). The
// JSON payload may contain "text" and/or "url" keys.
func AndroidShare(jsonPayload string) { androidBridgeVoidString("share", jsonPayload) }

// AndroidOpenURL opens the given URL in the system browser (Intent.ACTION_VIEW).
func AndroidOpenURL(url string) { androidBridgeVoidString("openURL", url) }

// AndroidSetKeepAwake adds or clears FLAG_KEEP_SCREEN_ON on the activity window.
func AndroidSetKeepAwake(enabled bool) { androidBridgeVoidInt("setKeepAwake", boolToInt(enabled)) }

// AndroidSetTorch toggles the camera flash via CameraManager.setTorchMode.
func AndroidSetTorch(enabled bool) { androidBridgeVoidInt("setTorch", boolToInt(enabled)) }
