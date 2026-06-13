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

func cString(s string) (*C.char, func()) {
	c := C.CString(s)
	return c, func() { C.free(unsafe.Pointer(c)) }
}

// --- Phase A: one-way actions -------------------------------------------------

// IOSShare presents the native iOS share sheet (UIActivityViewController). The
// JSON payload may contain "text" and/or "url" keys.
func IOSShare(jsonPayload string) {
	c, free := cString(jsonPayload)
	defer free()
	C.ios_share(c)
}

// IOSOpenURL opens the given URL in the system browser (Safari).
func IOSOpenURL(url string) {
	c, free := cString(url)
	defer free()
	C.ios_open_url(c)
}

// IOSSetKeepAwake disables (true) or restores (false) the idle timer, keeping
// the screen on while the app is in the foreground.
func IOSSetKeepAwake(enabled bool) { C.ios_set_keep_awake(C.bool(enabled)) }

// IOSSetTorch turns the device torch (flashlight) on or off.
func IOSSetTorch(enabled bool) { C.ios_set_torch(C.bool(enabled)) }

// --- Phase B: state / query --------------------------------------------------

func cStr(p *C.char) string {
	if p == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(p))
	return C.GoString(p)
}

// IOSSafeAreaJSON returns the safe-area insets ({top,bottom,left,right}) in points.
func IOSSafeAreaJSON() string { return cStr(C.ios_safe_area_json()) }

// IOSSetBrightness sets the screen brightness (0.0 - 1.0).
func IOSSetBrightness(value float64) { C.ios_set_brightness(C.double(value)) }

// IOSGetBrightness returns the current screen brightness (0.0 - 1.0).
func IOSGetBrightness() float64 { return float64(C.ios_get_brightness()) }

// IOSAppInfoJSON returns {name,version,build,bundleId} from the app bundle.
func IOSAppInfoJSON() string { return cStr(C.ios_app_info_json()) }

// IOSSetOrientation locks orientation to "portrait", "landscape" or "auto".
func IOSSetOrientation(mode string) {
	c, free := cString(mode)
	defer free()
	C.ios_set_orientation(c)
}

// IOSGetOrientation returns the current interface orientation ("portrait" /
// "landscape" / "unknown").
func IOSGetOrientation() string { return cStr(C.ios_get_orientation()) }

// IOSSetStatusBar sets the status-bar style/visibility. JSON: {"style":"light|
// dark|default","hidden":bool}.
func IOSSetStatusBar(jsonPayload string) {
	c, free := cString(jsonPayload)
	defer free()
	C.ios_set_status_bar(c)
}

// --- Phase C: async results / permissions ------------------------------------

// IOSBiometricAuthenticate triggers Face ID / Touch ID. The outcome is delivered
// to the frontend as the "native:biometric" event {ok, error}.
func IOSBiometricAuthenticate(reason string) {
	c, free := cString(reason)
	defer free()
	C.ios_biometric_authenticate(c)
}

// IOSPostNotification schedules a local notification. JSON: {"title","body",
// "delay":seconds}.
func IOSPostNotification(jsonPayload string) {
	c, free := cString(jsonPayload)
	defer free()
	C.ios_post_notification(c)
}

// IOSSecureSet stores a value in the Keychain under key.
func IOSSecureSet(key, value string) {
	ck, freeK := cString(key)
	defer freeK()
	cv, freeV := cString(value)
	defer freeV()
	C.ios_secure_set(ck, cv)
}

// IOSSecureGet reads a Keychain value (empty if absent).
func IOSSecureGet(key string) string {
	c, free := cString(key)
	defer free()
	return cStr(C.ios_secure_get(c))
}

// IOSSecureDelete removes a Keychain value.
func IOSSecureDelete(key string) {
	c, free := cString(key)
	defer free()
	C.ios_secure_delete(c)
}

// --- Phase D: sensors & hardware ---------------------------------------------

// IOSHaptic plays a haptic feedback pattern. type is one of impact-light,
// impact-medium, impact-heavy, success, warning, error, selection.
func IOSHaptic(hapticType string) {
	c, free := cString(hapticType)
	defer free()
	C.ios_haptic(c)
}

// IOSGetLocation requests a one-shot location fix. The result is delivered to
// the frontend as the "native:location" event {lat,lng,accuracy,error}.
func IOSGetLocation() { C.ios_get_location() }

// IOSSetMotion starts (true) or stops (false) accelerometer updates, streamed
// to the frontend as "native:motion" {x,y,z} events.
func IOSSetMotion(enabled bool) { C.ios_set_motion(C.bool(enabled)) }

// IOSSetProximity enables/disables the proximity sensor; changes arrive as
// "native:proximity" {near} events.
func IOSSetProximity(enabled bool) { C.ios_set_proximity(C.bool(enabled)) }

// IOSSpeak speaks the given text via AVSpeechSynthesizer.
func IOSSpeak(text string) {
	c, free := cString(text)
	defer free()
	C.ios_speak(c)
}

// IOSStopSpeak stops any in-progress speech.
func IOSStopSpeak() { C.ios_stop_speak() }

// IOSStorageJSON returns disk space as {"free":bytes,"total":bytes}.
func IOSStorageJSON() string { return cStr(C.ios_storage_json()) }

// IOSPowerJSON returns {"level":0-1,"charging":bool,"lowPower":bool}.
func IOSPowerJSON() string { return cStr(C.ios_power_json()) }

// IOSNetworkJSON returns {"connected":bool,"type":"wifi|cellular|none"}.
func IOSNetworkJSON() string { return cStr(C.ios_network_json()) }

// IOSSetKeyboardWatch starts/stops emitting "native:keyboard" {visible,height}
// events as the software keyboard shows and hides.
func IOSSetKeyboardWatch(enabled bool) { C.ios_set_keyboard_watch(C.bool(enabled)) }

// IOSSetScreenProtect starts/stops screenshot & screen-recording detection,
// reported as "native:screenCapture" events. (iOS cannot block screenshots, so
// this is detection-only; on Android the same control sets FLAG_SECURE.)
func IOSSetScreenProtect(enabled bool) { C.ios_set_screen_protect(C.bool(enabled)) }

// --- Phase E: camera & background --------------------------------------------

// IOSCapturePhoto presents the camera to take a photo; the result is delivered
// as the "native:capture" event {type:"photo",path,size,thumb}.
func IOSCapturePhoto() { C.ios_capture_photo() }

// IOSCaptureVideo presents the camera to record a video; the result is delivered
// as the "native:capture" event {type:"video",path,size}.
func IOSCaptureVideo() { C.ios_capture_video() }

// IOSBeginBackgroundTask opens a UIApplication background-task window so short
// work can finish after the app is backgrounded. iOS grants a limited amount of
// time; the granted/remaining seconds are reported as the "native:backgroundTask"
// event. Sustained background execution requires a declared UIBackgroundMode.
func IOSBeginBackgroundTask(seconds int) { C.ios_begin_background_task(C.int(seconds)) }

// IOSEndBackgroundTask ends the background-task window opened by
// IOSBeginBackgroundTask.
func IOSEndBackgroundTask() { C.ios_end_background_task() }

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
