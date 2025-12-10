//go:build ios

package application

import (
	"github.com/wailsapp/wails/v3/pkg/errs"
)

const (
	IOSHapticsImpact                    = 0
	IOSDeviceInfo                       = 1
	IOSScrollSetEnabled                 = 2
	IOSScrollSetBounceEnabled           = 3
	IOSScrollSetIndicatorsEnabled       = 4
	IOSNavigationSetBackForwardGestures = 5
	IOSLinksSetPreviewEnabled           = 6
	IOSDebugSetInspectableEnabled       = 7
	IOSUserAgentSet                     = 8
)

var iosMethodNames = map[int]string{
	IOSHapticsImpact:                    "Haptics.Impact",
	IOSDeviceInfo:                       "Device.Info",
	IOSScrollSetEnabled:                 "Scroll.SetEnabled",
	IOSScrollSetBounceEnabled:           "Scroll.SetBounceEnabled",
	IOSScrollSetIndicatorsEnabled:       "Scroll.SetIndicatorsEnabled",
	IOSNavigationSetBackForwardGestures: "Navigation.SetBackForwardGesturesEnabled",
	IOSLinksSetPreviewEnabled:           "Links.SetPreviewEnabled",
	IOSDebugSetInspectableEnabled:       "Debug.SetInspectableEnabled",
	IOSUserAgentSet:                     "UserAgent.Set",
}

func (m *MessageProcessor) processIOSMethod(req *RuntimeRequest, window Window) (any, error) {
	args := req.Args.AsMap()

	switch req.Method {
	case IOSHapticsImpact:
		style := "medium"
		if s := args.String("style"); s != nil {
			style = *s
		}
		iosHapticsImpact(style)
		return unit, nil
	case IOSDeviceInfo:
		return iosDeviceInfo(), nil
	case IOSScrollSetEnabled:
		enabled := true
		if b := args.Bool("enabled"); b != nil {
			enabled = *b
		}
		iosSetScrollEnabled(enabled)
		return unit, nil
	case IOSScrollSetBounceEnabled:
		enabled := true
		if b := args.Bool("enabled"); b != nil {
			enabled = *b
		}
		iosSetBounceEnabled(enabled)
		return unit, nil
	case IOSScrollSetIndicatorsEnabled:
		enabled := true
		if b := args.Bool("enabled"); b != nil {
			enabled = *b
		}
		iosSetScrollIndicatorsEnabled(enabled)
		return unit, nil
	case IOSNavigationSetBackForwardGestures:
		enabled := false
		if b := args.Bool("enabled"); b != nil {
			enabled = *b
		}
		iosSetBackForwardGesturesEnabled(enabled)
		return unit, nil
	case IOSLinksSetPreviewEnabled:
		enabled := true
		if b := args.Bool("enabled"); b != nil {
			enabled = *b
		}
		iosSetLinkPreviewEnabled(enabled)
		return unit, nil
	case IOSDebugSetInspectableEnabled:
		enabled := true
		if b := args.Bool("enabled"); b != nil {
			enabled = *b
		}
		iosSetInspectableEnabled(enabled)
		return unit, nil
	case IOSUserAgentSet:
		ua := ""
		if s := args.String("ua"); s != nil {
			ua = *s
		} else if s2 := args.String("userAgent"); s2 != nil {
			ua = *s2
		}
		iosSetCustomUserAgent(ua)
		return unit, nil
	default:
		return nil, errs.NewInvalidIOSCallErrorf("unknown method: %d", req.Method)
	}
}

// processAndroidMethod is a stub on iOS
func (m *MessageProcessor) processAndroidMethod(req *RuntimeRequest, window Window) (any, error) {
	return nil, errs.NewInvalidAndroidCallErrorf("Android methods not available on iOS")
}
