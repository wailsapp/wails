//go:build android && !server

package application

import (
	"github.com/wailsapp/wails/v3/pkg/errs"
)

const (
	AndroidHapticsVibrate                   = 0
	AndroidDeviceInfo                       = 1
	AndroidToast                            = 2
	AndroidScrollSetEnabled                 = 3
	AndroidScrollSetBounceEnabled           = 4
	AndroidScrollSetIndicatorsEnabled       = 5
	AndroidNavigationSetBackForwardGestures = 6
	AndroidLinksSetPreviewEnabled           = 7
	AndroidUserAgentSet                     = 8
)

var androidMethodNames = map[int]string{
	AndroidHapticsVibrate:                   "Haptics.Vibrate",
	AndroidDeviceInfo:                       "Device.Info",
	AndroidToast:                            "Toast.Show",
	AndroidScrollSetEnabled:                 "Scroll.SetEnabled",
	AndroidScrollSetBounceEnabled:           "Scroll.SetBounceEnabled",
	AndroidScrollSetIndicatorsEnabled:       "Scroll.SetIndicatorsEnabled",
	AndroidNavigationSetBackForwardGestures: "Navigation.SetBackForwardGesturesEnabled",
	AndroidLinksSetPreviewEnabled:           "Links.SetPreviewEnabled",
	AndroidUserAgentSet:                     "UserAgent.Set",
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
	case AndroidScrollSetEnabled:
		enabled := true
		if b := args.Bool("enabled"); b != nil {
			enabled = *b
		}
		androidSetScrollEnabled(enabled)
		return unit, nil
	case AndroidScrollSetBounceEnabled:
		enabled := true
		if b := args.Bool("enabled"); b != nil {
			enabled = *b
		}
		androidSetBounceEnabled(enabled)
		return unit, nil
	case AndroidScrollSetIndicatorsEnabled:
		enabled := true
		if b := args.Bool("enabled"); b != nil {
			enabled = *b
		}
		androidSetScrollIndicatorsEnabled(enabled)
		return unit, nil
	case AndroidNavigationSetBackForwardGestures:
		enabled := false
		if b := args.Bool("enabled"); b != nil {
			enabled = *b
		}
		androidSetBackForwardGesturesEnabled(enabled)
		return unit, nil
	case AndroidLinksSetPreviewEnabled:
		enabled := true
		if b := args.Bool("enabled"); b != nil {
			enabled = *b
		}
		androidSetLinkPreviewEnabled(enabled)
		return unit, nil
	case AndroidUserAgentSet:
		ua := ""
		if s := args.String("ua"); s != nil {
			ua = *s
		} else if s2 := args.String("userAgent"); s2 != nil {
			ua = *s2
		}
		androidSetCustomUserAgent(ua)
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
	return androidDeviceInfoViaJNI()
}

func androidShowToast(message string) {
	// TODO: Implement via JNI to Android Toast
	androidLogf("debug", "Toast: %s", message)
}
