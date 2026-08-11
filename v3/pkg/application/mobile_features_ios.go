//go:build ios

package application

/*
#cgo CFLAGS: -x objective-c -fmodules -fobjc-arc
#cgo LDFLAGS: -framework UIKit -framework LocalAuthentication -framework UserNotifications -framework AVFoundation -framework Security -framework CoreLocation -framework CoreMotion -framework SystemConfiguration
#include <stdlib.h>
#include "mobile_features_ios.h"
*/
import "C"

import (
	"encoding/json"
	"unsafe"
)

// This file exposes a set of "genuinely mobile" native capabilities to Wails
// applications running on iOS: share sheet, opening URLs, keeping the screen
// awake, the torch, safe-area insets, brightness, orientation, the status bar,
// biometric authentication, local notifications and Keychain-backed secure
// storage. Each maps to a small Objective-C bridge in mobile_features_ios.m.
//
// Results that arrive asynchronously (e.g. a biometric prompt) are delivered
// back to the frontend as custom events via iosEmitNativeEvent.

// iosManager is the receiver for all iOS-specific native feature methods. The
// package-level IOS singleton is the entry point: application.IOS.Haptic(...),
// application.IOS.Share(...), and so on. It is only present on iOS builds; all
// callers live in //go:build ios files.
type iosManager struct{}

// IOS exposes iOS-specific native capabilities (share sheet, haptics, biometrics,
// secure storage, sensors, camera, and WKWebView runtime options).
var IOS iosManager

func cString(s string) (*C.char, func()) {
	c := C.CString(s)
	return c, func() { C.free(unsafe.Pointer(c)) }
}

// --- Phase A: one-way actions -------------------------------------------------

// Share presents the native iOS share sheet (UIActivityViewController). The
// JSON payload may contain "text" and/or "url" keys.
func (iosManager) Share(jsonPayload string) {
	c, free := cString(jsonPayload)
	defer free()
	C.ios_share(c)
}

// OpenURL opens the given URL in the system browser (Safari).
func (iosManager) OpenURL(url string) {
	c, free := cString(url)
	defer free()
	C.ios_open_url(c)
}

// SetKeepAwake disables (true) or restores (false) the idle timer, keeping
// the screen on while the app is in the foreground.
func (iosManager) SetKeepAwake(enabled bool) { C.ios_set_keep_awake(C.bool(enabled)) }

// SetTorch turns the device torch (flashlight) on or off.
func (iosManager) SetTorch(enabled bool) { C.ios_set_torch(C.bool(enabled)) }

// --- Phase B: state / query --------------------------------------------------

func cStr(p *C.char) string {
	if p == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(p))
	return C.GoString(p)
}

// SafeAreaJSON returns the safe-area insets ({top,bottom,left,right}) in points.
func (iosManager) SafeAreaJSON() string { return cStr(C.ios_safe_area_json()) }

// SetBrightness sets the screen brightness (0.0 - 1.0).
func (iosManager) SetBrightness(value float64) { C.ios_set_brightness(C.double(value)) }

// GetBrightness returns the current screen brightness (0.0 - 1.0).
func (iosManager) GetBrightness() float64 { return float64(C.ios_get_brightness()) }

// AppInfoJSON returns {name,version,build,bundleId} from the app bundle.
func (iosManager) AppInfoJSON() string { return cStr(C.ios_app_info_json()) }

// SetOrientation locks orientation to "portrait", "landscape" or "auto".
func (iosManager) SetOrientation(mode string) {
	c, free := cString(mode)
	defer free()
	C.ios_set_orientation(c)
}

// GetOrientation returns the current interface orientation ("portrait" /
// "landscape" / "unknown").
func (iosManager) GetOrientation() string { return cStr(C.ios_get_orientation()) }

// SetStatusBar sets the status-bar style/visibility. JSON: {"style":"light|
// dark|default","hidden":bool}.
func (iosManager) SetStatusBar(jsonPayload string) {
	c, free := cString(jsonPayload)
	defer free()
	C.ios_set_status_bar(c)
}

// --- Phase C: async results / permissions ------------------------------------

// BiometricAuthenticate triggers Face ID / Touch ID. The outcome is delivered
// to the frontend as the "common:biometric" event {ok, error}.
func (iosManager) BiometricAuthenticate(reason string) {
	c, free := cString(reason)
	defer free()
	C.ios_biometric_authenticate(c)
}

// PostNotification schedules a local notification. JSON: {"title","body",
// "delay":seconds}.
func (iosManager) PostNotification(jsonPayload string) {
	c, free := cString(jsonPayload)
	defer free()
	C.ios_post_notification(c)
}

// SecureSet stores a value in the Keychain under key.
func (iosManager) SecureSet(key, value string) {
	ck, freeK := cString(key)
	defer freeK()
	cv, freeV := cString(value)
	defer freeV()
	C.ios_secure_set(ck, cv)
}

// SecureGet reads a Keychain value (empty if absent).
func (iosManager) SecureGet(key string) string {
	c, free := cString(key)
	defer free()
	return cStr(C.ios_secure_get(c))
}

// SecureDelete removes a Keychain value.
func (iosManager) SecureDelete(key string) {
	c, free := cString(key)
	defer free()
	C.ios_secure_delete(c)
}

// --- Phase D: sensors & hardware ---------------------------------------------

// Haptic plays a haptic feedback pattern. type is one of impact-light,
// impact-medium, impact-heavy, success, warning, error, selection.
func (iosManager) Haptic(hapticType string) {
	c, free := cString(hapticType)
	defer free()
	C.ios_haptic(c)
}

// GetLocation requests a one-shot location fix. The result is delivered to
// the frontend as the "common:location" event {lat,lng,accuracy,error}.
func (iosManager) GetLocation() { C.ios_get_location() }

// SetMotion starts (true) or stops (false) accelerometer updates, streamed
// to the frontend as "common:motion" {x,y,z} events.
func (iosManager) SetMotion(enabled bool) { C.ios_set_motion(C.bool(enabled)) }

// SetProximity enables/disables the proximity sensor; changes arrive as
// "common:proximity" {near} events.
func (iosManager) SetProximity(enabled bool) { C.ios_set_proximity(C.bool(enabled)) }

// Speak speaks the given text via AVSpeechSynthesizer.
func (iosManager) Speak(text string) {
	c, free := cString(text)
	defer free()
	C.ios_speak(c)
}

// StopSpeak stops any in-progress speech.
func (iosManager) StopSpeak() { C.ios_stop_speak() }

// StorageJSON returns disk space as {"free":bytes,"total":bytes}.
func (iosManager) StorageJSON() string { return cStr(C.ios_storage_json()) }

// StoragePath returns the absolute path to the app's Application Support
// directory, suitable for databases and other persistent files (the iOS analog
// of Android's getFilesDir()). The directory is created if it does not yet exist.
func (iosManager) StoragePath() string { return cStr(C.ios_storage_path()) }

// PowerJSON returns {"level":0-1,"charging":bool,"lowPower":bool}.
func (iosManager) PowerJSON() string { return cStr(C.ios_power_json()) }

// NetworkJSON returns {"connected":bool,"type":"wifi|cellular|none"}.
func (iosManager) NetworkJSON() string { return cStr(C.ios_network_json()) }

// SetKeyboardWatch starts/stops emitting "common:keyboard" {visible,height}
// events as the software keyboard shows and hides.
func (iosManager) SetKeyboardWatch(enabled bool) { C.ios_set_keyboard_watch(C.bool(enabled)) }

// SetScreenProtect starts/stops screenshot & screen-recording detection,
// reported as "common:screenCapture" events. (iOS cannot block screenshots, so
// this is detection-only; on Android the same control sets FLAG_SECURE.)
func (iosManager) SetScreenProtect(enabled bool) { C.ios_set_screen_protect(C.bool(enabled)) }

// --- Phase E: camera & background --------------------------------------------

// CapturePhoto presents the camera to take a photo; the result is delivered
// as the "common:capture" event {type:"photo",path,size,thumb}.
func (iosManager) CapturePhoto() { C.ios_capture_photo() }

// CaptureVideo presents the camera to record a video; the result is delivered
// as the "common:capture" event {type:"video",path,size}.
func (iosManager) CaptureVideo() { C.ios_capture_video() }

// BeginBackgroundTask opens a UIApplication background-task window so short
// work can finish after the app is backgrounded. iOS grants a limited amount of
// time; the granted/remaining seconds are reported as the "ios:backgroundTask"
// event. Sustained background execution requires a declared UIBackgroundMode.
func (iosManager) BeginBackgroundTask(seconds int) { C.ios_begin_background_task(C.int(seconds)) }

// EndBackgroundTask ends the background-task window opened by
// BeginBackgroundTask.
func (iosManager) EndBackgroundTask() { C.ios_end_background_task() }

// iosEmitNativeEvent is called from the Objective-C bridge to deliver an
// asynchronous result to the frontend as a custom event.
//
//export iosEmitNativeEvent
func iosEmitNativeEvent(cname *C.char, cjson *C.char) {
	name := C.GoString(cname)
	var data map[string]any
	if cjson != nil {
		if s := C.GoString(cjson); s != "" {
			_ = json.Unmarshal([]byte(s), &data)
		}
	}
	app := globalApplication
	if app == nil {
		return
	}
	app.Event.Emit(name, data)
}
