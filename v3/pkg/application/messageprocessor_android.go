//go:build android

package application

import (
	"encoding/json"

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

// iosMethodNames is referenced by the shared messageprocessor debug logging;
// empty on Android.
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

// Android-specific runtime functions, backed by the WailsBridge

func androidHapticsVibrate(durationMs int) {
	androidBridgeVoidInt("vibrate", durationMs)
}

func androidDeviceInfo() map[string]interface{} {
	info := map[string]interface{}{
		"platform": "android",
	}
	if jsonStr, ok := androidBridgeString("getDeviceInfoJson"); ok && jsonStr != "" {
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &parsed); err == nil {
			for k, v := range parsed {
				info[k] = v
			}
		}
	}
	return info
}

func androidShowToast(message string) {
	androidBridgeVoidString("showToast", message)
}
