//go:build openharmony && !server

package application

import (
	"encoding/json"

	"github.com/wailsapp/wails/v3/pkg/errs"
)

const (
	OpenharmonyDeviceInfo = 1
	OpenharmonyToast      = 2
)

var openharmonyMethodNames = map[int]string{
	OpenharmonyDeviceInfo: "Device.Info",
	OpenharmonyToast:      "Toast.Show",
}

// processOpenharmonyMethod handles OHOS-specific runtime calls, backed by the
// WailsBridge (same shape as processAndroidMethod).
func (m *MessageProcessor) processOpenharmonyMethod(req *RuntimeRequest, window Window) (any, error) {
	args := req.Args.AsMap()

	switch req.Method {
	case OpenharmonyDeviceInfo:
		return openharmonyDeviceInfo(), nil
	case OpenharmonyToast:
		message := ""
		if s := args.String("message"); s != nil {
			message = *s
		}
		openharmonyShowToast(message)
		return unit, nil
	default:
		return nil, errs.NewInvalidRuntimeCallErrorf("unknown OpenHarmony method: %d", req.Method)
	}
}

func openharmonyDeviceInfo() map[string]interface{} {
	info := map[string]interface{}{
		"platform": "openharmony",
	}
	if jsonStr, ok := ohBridgeString("getDeviceInfoJson", "{}"); ok && jsonStr != "" {
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &parsed); err == nil {
			for k, v := range parsed {
				info[k] = v
			}
		}
	}
	return info
}

func openharmonyShowToast(message string) {
	payload, _ := json.Marshal(map[string]string{"message": message})
	ohBridgeVoid("showToast", string(payload))
}
