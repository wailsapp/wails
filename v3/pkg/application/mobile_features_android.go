//go:build android

package application

// Genuinely-mobile native capabilities for Android, mirroring the iOS surface
// in mobile_features_ios.go. Each call is forwarded to a matching method on the
// Java WailsBridge via the reflective bridge helpers. Asynchronous results
// (biometric prompts, torch availability, …) come back to the frontend as
// custom events through the nativeEmitEvent JNI export.

// androidManager is the receiver for all Android-specific native feature methods.
// The package-level Android singleton is the entry point:
// application.Android.Haptic(...), application.Android.Share(...), and so on. It
// is only present on Android builds; all callers live in //go:build android files.
type androidManager struct{}

// Android exposes Android-specific native capabilities (share chooser, haptics,
// biometrics, secure storage, sensors, camera, toast, …).
var Android androidManager

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// --- Phase A: one-way actions -------------------------------------------------

// Share presents the Android share chooser (Intent.ACTION_SEND). The
// JSON payload may contain "text" and/or "url" keys.
func (androidManager) Share(jsonPayload string) { androidBridgeVoidString("share", jsonPayload) }

// OpenURL opens the given URL in the system browser (Intent.ACTION_VIEW).
func (androidManager) OpenURL(url string) { androidBridgeVoidString("openURL", url) }

// SetKeepAwake adds or clears FLAG_KEEP_SCREEN_ON on the activity window.
func (androidManager) SetKeepAwake(enabled bool) {
	androidBridgeVoidInt("setKeepAwake", boolToInt(enabled))
}

// SetTorch toggles the camera flash via CameraManager.setTorchMode.
func (androidManager) SetTorch(enabled bool) { androidBridgeVoidInt("setTorch", boolToInt(enabled)) }

// --- Phase B: state / query ---------------------------------------------------

// SafeAreaJSON returns the system-bar insets ({top,bottom,left,right}) in px.
func (androidManager) SafeAreaJSON() string { s, _ := androidBridgeString("getSafeAreaJson"); return s }

// SetBrightness sets the window brightness, 0-100 (negative restores the
// system default).
func (androidManager) SetBrightness(pct int) { androidBridgeVoidInt("setBrightness", pct) }

// BrightnessJSON returns the current brightness as {"value":0.0-1.0}.
func (androidManager) BrightnessJSON() string {
	s, _ := androidBridgeString("getBrightnessJson")
	return s
}

// AppInfoJSON returns {name,version,build,bundleId} for the app.
func (androidManager) AppInfoJSON() string { s, _ := androidBridgeString("getAppInfoJson"); return s }

// SetOrientation locks orientation to "portrait", "landscape" or "auto".
func (androidManager) SetOrientation(mode string) { androidBridgeVoidString("setOrientation", mode) }

// OrientationJSON returns the current orientation as {"orientation":"…"}.
func (androidManager) OrientationJSON() string {
	s, _ := androidBridgeString("getOrientationJson")
	return s
}

// SetStatusBar sets the status-bar appearance. JSON: {"style":"light|dark|
// default","hidden":bool}.
func (androidManager) SetStatusBar(jsonPayload string) {
	androidBridgeVoidString("setStatusBar", jsonPayload)
}

// --- Phase C: async results / permissions -------------------------------------

// BiometricAuthenticate shows the BiometricPrompt. The outcome is
// delivered to the frontend as the "common:biometric" event {ok, error}.
func (androidManager) BiometricAuthenticate(reason string) {
	androidBridgeVoidString("authenticate", reason)
}

// Notify posts a local notification. JSON: {"title","body","delay":seconds}.
func (androidManager) Notify(jsonPayload string) {
	androidBridgeVoidString("postNotification", jsonPayload)
}

// SecureSet stores a value in EncryptedSharedPreferences. JSON: {"key","value"}.
func (androidManager) SecureSet(jsonPayload string) {
	androidBridgeVoidString("secureSet", jsonPayload)
}

// SecureGet reads a value from secure storage (empty if absent).
func (androidManager) SecureGet(key string) string {
	s, _ := androidBridgeStringString("secureGet", key)
	return s
}

// SecureDelete removes a value from secure storage.
func (androidManager) SecureDelete(key string) { androidBridgeVoidString("secureDelete", key) }

// --- Phase D: sensors & hardware ---------------------------------------------

// Haptic plays a haptic pattern (impact-light/medium/heavy, success,
// warning, error, selection) via the Vibrator.
func (androidManager) Haptic(hapticType string) { androidBridgeVoidString("haptic", hapticType) }

// GetLocation requests a one-shot location fix; the result arrives as the
// "common:location" event {lat,lng,accuracy,error}.
func (androidManager) GetLocation() { androidBridgeVoid("getLocation") }

// SetMotion starts/stops accelerometer updates, streamed as
// "common:motion" {x,y,z} events.
func (androidManager) SetMotion(enabled bool) { androidBridgeVoidInt("setMotion", boolToInt(enabled)) }

// SetProximity enables/disables the proximity sensor; changes arrive as
// "common:proximity" {near} events.
func (androidManager) SetProximity(enabled bool) {
	androidBridgeVoidInt("setProximity", boolToInt(enabled))
}

// Speak speaks the given text via Android TextToSpeech.
func (androidManager) Speak(text string) { androidBridgeVoidString("speak", text) }

// StopSpeak stops any in-progress speech.
func (androidManager) StopSpeak() { androidBridgeVoid("stopSpeak") }

// StorageJSON returns disk space as {"free":bytes,"total":bytes}.
func (androidManager) StorageJSON() string { s, _ := androidBridgeString("getStorageJson"); return s }

// StoragePath returns the absolute path to the app's private internal files
// directory (activity.getFilesDir()), suitable for databases and other
// persistent files. The directory always exists.
func (androidManager) StoragePath() string { s, _ := androidBridgeString("getStoragePath"); return s }

// PowerJSON returns {"level":0-1,"charging":bool,"lowPower":bool}.
func (androidManager) PowerJSON() string { s, _ := androidBridgeString("getPowerJson"); return s }

// NetworkJSON returns {"connected":bool,"type":"wifi|cellular|ethernet|none"}.
func (androidManager) NetworkJSON() string { s, _ := androidBridgeString("getNetworkJson"); return s }

// SetKeyboardWatch starts/stops emitting "common:keyboard"
// {visible,height} events as the soft keyboard shows and hides.
func (androidManager) SetKeyboardWatch(enabled bool) {
	androidBridgeVoidInt("setKeyboardWatch", boolToInt(enabled))
}

// SetScreenProtect toggles FLAG_SECURE (blocks screenshots & screen
// recording) and reports state as a "common:screenCapture" event.
func (androidManager) SetScreenProtect(enabled bool) {
	androidBridgeVoidInt("setScreenProtect", boolToInt(enabled))
}

// --- Phase E: camera & background -------------------------------------------

// CapturePhoto launches the system camera to take a photo; the result
// arrives as the "common:capture" event {type:"photo",path,size,thumb}.
func (androidManager) CapturePhoto() { androidBridgeVoidString("capturePhoto", "{}") }

// CaptureVideo launches the system camera to record a video; the result
// arrives as the "common:capture" event {type:"video",path,size}.
func (androidManager) CaptureVideo() { androidBridgeVoidString("captureVideo", "{}") }

// StartForegroundService starts a foreground service (with an ongoing
// notification) so the process keeps running for long-running background work.
// JSON: {"title","text"}.
func (androidManager) StartForegroundService(jsonPayload string) {
	androidBridgeVoidString("startForegroundService", jsonPayload)
}

// StopForegroundService stops the foreground service.
func (androidManager) StopForegroundService() { androidBridgeVoid("stopForegroundService") }
