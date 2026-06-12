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

// --- Phase B: state / query ---------------------------------------------------

// AndroidSafeAreaJSON returns the system-bar insets ({top,bottom,left,right}) in px.
func AndroidSafeAreaJSON() string { s, _ := androidBridgeString("getSafeAreaJson"); return s }

// AndroidSetBrightness sets the window brightness, 0-100 (negative restores the
// system default).
func AndroidSetBrightness(pct int) { androidBridgeVoidInt("setBrightness", pct) }

// AndroidBrightnessJSON returns the current brightness as {"value":0.0-1.0}.
func AndroidBrightnessJSON() string { s, _ := androidBridgeString("getBrightnessJson"); return s }

// AndroidAppInfoJSON returns {name,version,build,bundleId} for the app.
func AndroidAppInfoJSON() string { s, _ := androidBridgeString("getAppInfoJson"); return s }

// AndroidSetOrientation locks orientation to "portrait", "landscape" or "auto".
func AndroidSetOrientation(mode string) { androidBridgeVoidString("setOrientation", mode) }

// AndroidOrientationJSON returns the current orientation as {"orientation":"…"}.
func AndroidOrientationJSON() string { s, _ := androidBridgeString("getOrientationJson"); return s }

// AndroidSetStatusBar sets the status-bar appearance. JSON: {"style":"light|dark|
// default","hidden":bool}.
func AndroidSetStatusBar(jsonPayload string) { androidBridgeVoidString("setStatusBar", jsonPayload) }

// --- Phase C: async results / permissions -------------------------------------

// AndroidBiometricAuthenticate shows the BiometricPrompt. The outcome is
// delivered to the frontend as the "native:biometric" event {ok, error}.
func AndroidBiometricAuthenticate(reason string) { androidBridgeVoidString("authenticate", reason) }

// AndroidNotify posts a local notification. JSON: {"title","body","delay":seconds}.
func AndroidNotify(jsonPayload string) { androidBridgeVoidString("postNotification", jsonPayload) }

// AndroidSecureSet stores a value in EncryptedSharedPreferences. JSON: {"key","value"}.
func AndroidSecureSet(jsonPayload string) { androidBridgeVoidString("secureSet", jsonPayload) }

// AndroidSecureGet reads a value from secure storage (empty if absent).
func AndroidSecureGet(key string) string {
	s, _ := androidBridgeStringString("secureGet", key)
	return s
}

// AndroidSecureDelete removes a value from secure storage.
func AndroidSecureDelete(key string) { androidBridgeVoidString("secureDelete", key) }
