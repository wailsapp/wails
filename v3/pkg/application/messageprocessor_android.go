//go:build android

package application

import (
	"fmt"
	"net/http"
)

const (
	AndroidHapticsVibrate = 0
	AndroidDeviceInfo     = 1
	AndroidToast          = 2
)

var androidMethodNames = map[int]string{
	AndroidHapticsVibrate: "Haptics.Vibrate",
	AndroidDeviceInfo:     "Device.Info",
	AndroidToast:          "Toast.Show",
}

func (m *MessageProcessor) processAndroidMethod(method int, rw http.ResponseWriter, r *http.Request, window Window, params QueryParams) {
	switch method {
	case AndroidHapticsVibrate:
		args, _ := params.Args()
		duration := 100 // default 100ms
		if d := args.Int("duration"); d != nil {
			duration = *d
		}
		androidHapticsVibrate(duration)
		m.ok(rw)
	case AndroidDeviceInfo:
		m.json(rw, androidDeviceInfo())
	case AndroidToast:
		args, _ := params.Args()
		message := ""
		if s := args.String("message"); s != nil {
			message = *s
		}
		androidShowToast(message)
		m.ok(rw)
	default:
		m.httpError(rw, "Invalid Android call:", fmt.Errorf("unknown method: %d", method))
		return
	}

	m.Info("Runtime call:", "method", "Android."+androidMethodNames[method])
}

// processIOSMethod is a stub on Android
func (m *MessageProcessor) processIOSMethod(method int, rw http.ResponseWriter, r *http.Request, window Window, params QueryParams) {
	m.httpError(rw, "iOS methods not available on Android:", fmt.Errorf("unknown method: %d", method))
}

// Android-specific runtime functions (stubs for now)

func androidHapticsVibrate(durationMs int) {
	// TODO: Implement via JNI to Android Vibrator service
	androidLogf("debug", "Haptics vibrate: %dms", durationMs)
}

func androidDeviceInfo() map[string]interface{} {
	// TODO: Implement via JNI to get actual device info
	return map[string]interface{}{
		"platform": "android",
		"model":    "Unknown",
		"version":  "Unknown",
	}
}

func androidShowToast(message string) {
	// TODO: Implement via JNI to Android Toast
	androidLogf("debug", "Toast: %s", message)
}
