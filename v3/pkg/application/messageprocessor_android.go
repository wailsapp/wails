//go:build android

package application

import (
	"github.com/wailsapp/wails/v3/pkg/errs"
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

var iosMethodNames = map[int]string{}

func (m *MessageProcessor) processAndroidMethod(req *RuntimeRequest, window Window) (any, error) {
	args := req.Args.AsMap()

	switch req.Method {
	case AndroidHapticsVibrate:
		duration := 100 // default 100ms
		if d := args.Int("duration"); d != nil {
			duration = *d
		}
		androidHapticsVibrate(duration)
		return unit, nil
	case AndroidDeviceInfo:
		return androidDeviceInfo(), nil
	case AndroidToast:
		message := ""
		if s := args.String("message"); s != nil {
			message = *s
		}
		androidShowToast(message)
		return unit, nil
	default:
		return nil, errs.NewInvalidAndroidCallErrorf("unknown method: %d", req.Method)
	}
}

// processIOSMethod is a stub on Android
func (m *MessageProcessor) processIOSMethod(req *RuntimeRequest, window Window) (any, error) {
	return nil, errs.NewInvalidIOSCallErrorf("iOS methods not available on Android")
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
