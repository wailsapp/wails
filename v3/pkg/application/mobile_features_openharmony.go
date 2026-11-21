//go:build openharmony && cgo && !server

package application

// Genuinely-mobile native capabilities for OHOS, mirroring the Android surface
// in mobile_features_android.go. Each call is forwarded to a matching method on
// the ArkTS WailsBridge via the bridge helpers. Asynchronous results come back
// to the frontend as custom events through WailsOhNativeEmitEvent.
//
// The package-level Openharmony singleton is the entry point:
// application.Openharmony.Haptic(...), application.Openharmony.Share(...), etc.

type openharmonyManager struct{}

// Openharmony exposes OHOS-specific native capabilities.
var Openharmony openharmonyManager

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// --- one-way actions -------------------------------------------------------

func (openharmonyManager) Share(jsonPayload string) { ohBridgeVoid("share", jsonPayload) }

func (openharmonyManager) OpenURL(url string) {
	ohBridgeVoid("openURL", jsonStringField("url", url))
}

func (openharmonyManager) SetKeepAwake(on bool) { ohBridgeVoidInt("setKeepAwake", boolToInt(on)) }

func (openharmonyManager) SetTorch(on bool) { ohBridgeVoidInt("setTorch", boolToInt(on)) }

func (openharmonyManager) SetBrightness(val int) { ohBridgeVoidInt("setBrightness", val) }

func (openharmonyManager) SetOrientation(o string) {
	ohBridgeVoid("setOrientation", jsonStringField("orientation", o))
}

func (openharmonyManager) SetStatusBar(jsonPayload string) { ohBridgeVoid("setStatusBar", jsonPayload) }

func (openharmonyManager) Notify(jsonPayload string) { ohBridgeVoid("postNotification", jsonPayload) }

func (openharmonyManager) Haptic(jsonPayload string) { ohBridgeVoid("haptic", jsonPayload) }

func (openharmonyManager) BiometricAuthenticate(jsonPayload string) {
	ohBridgeVoid("authenticate", jsonPayload)
}

func (openharmonyManager) SetMotion(jsonPayload string) { ohBridgeVoid("setMotion", jsonPayload) }

func (openharmonyManager) SetProximity(jsonPayload string) { ohBridgeVoid("setProximity", jsonPayload) }

func (openharmonyManager) SetKeyboardWatch(on bool) {
	ohBridgeVoidInt("setKeyboardWatch", boolToInt(on))
}

func (openharmonyManager) SetScreenProtect(on bool) {
	ohBridgeVoidInt("setScreenProtect", boolToInt(on))
}

func (openharmonyManager) Speak(text string) { ohBridgeVoid("speak", jsonStringField("text", text)) }

func (openharmonyManager) StopSpeak() { ohBridgeVoid("stopSpeak", "{}") }

func (openharmonyManager) GetLocation() { ohBridgeVoid("getLocation", "{}") }

func (openharmonyManager) CapturePhoto(jsonPayload string) { ohBridgeVoid("capturePhoto", jsonPayload) }

func (openharmonyManager) CaptureVideo(jsonPayload string) { ohBridgeVoid("captureVideo", jsonPayload) }

func (openharmonyManager) StartForegroundService() { ohBridgeVoid("startForegroundService", "{}") }

func (openharmonyManager) StopForegroundService() { ohBridgeVoid("stopForegroundService", "{}") }

// --- secure storage --------------------------------------------------------

func (openharmonyManager) SecureSet(key, value string) {
	ohBridgeVoid("secureSet", string(mustJSON(map[string]string{"key": key, "value": value})))
}

func (openharmonyManager) SecureGet(key string) (string, bool) {
	return ohBridgeString("secureGet", jsonStringField("key", key))
}

func (openharmonyManager) SecureDelete(key string) {
	ohBridgeVoid("secureDelete", jsonStringField("key", key))
}

// --- info getters ----------------------------------------------------------

func (openharmonyManager) AppInfoJSON() string {
	s, _ := ohBridgeString("getAppInfoJson", "{}")
	return s
}
func (openharmonyManager) SafeAreaJSON() string {
	s, _ := ohBridgeString("getSafeAreaJson", "{}")
	return s
}
func (openharmonyManager) BrightnessJSON() string {
	s, _ := ohBridgeString("getBrightnessJson", "{}")
	return s
}
func (openharmonyManager) OrientationJSON() string {
	s, _ := ohBridgeString("getOrientationJson", "{}")
	return s
}
func (openharmonyManager) StorageJSON() string {
	s, _ := ohBridgeString("getStorageJson", "{}")
	return s
}
func (openharmonyManager) StoragePath() string {
	s, _ := ohBridgeString("getStoragePath", "{}")
	return s
}
func (openharmonyManager) PowerJSON() string { s, _ := ohBridgeString("getPowerJson", "{}"); return s }
func (openharmonyManager) NetworkJSON() string {
	s, _ := ohBridgeString("getNetworkJson", "{}")
	return s
}

func jsonStringField(key, value string) string {
	return string(mustJSON(map[string]string{key: value}))
}
