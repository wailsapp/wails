//go:build ios

package application

import (
	"fmt"
	"net/http"
)

const (
	IOSHapticsImpact = 0
	IOSDeviceInfo    = 1
	IOSScrollSetEnabled                   = 2
	IOSScrollSetBounceEnabled             = 3
	IOSScrollSetIndicatorsEnabled         = 4
	IOSNavigationSetBackForwardGestures   = 5
	IOSLinksSetPreviewEnabled             = 6
	IOSDebugSetInspectableEnabled         = 7
	IOSUserAgentSet                       = 8
)

var iosMethodNames = map[int]string{
	IOSHapticsImpact: "Haptics.Impact",
	IOSDeviceInfo:    "Device.Info",
	IOSScrollSetEnabled:                 "Scroll.SetEnabled",
	IOSScrollSetBounceEnabled:           "Scroll.SetBounceEnabled",
	IOSScrollSetIndicatorsEnabled:       "Scroll.SetIndicatorsEnabled",
	IOSNavigationSetBackForwardGestures: "Navigation.SetBackForwardGesturesEnabled",
	IOSLinksSetPreviewEnabled:           "Links.SetPreviewEnabled",
	IOSDebugSetInspectableEnabled:       "Debug.SetInspectableEnabled",
	IOSUserAgentSet:                     "UserAgent.Set",
}

func (m *MessageProcessor) processIOSMethod(method int, rw http.ResponseWriter, r *http.Request, window Window, params QueryParams) {
	switch method {
	case IOSHapticsImpact:
		args, _ := params.Args()
		style := "medium"
		if s := args.String("style"); s != nil {
			style = *s
		}
		iosHapticsImpact(style)
		m.ok(rw)
	case IOSDeviceInfo:
		m.json(rw, iosDeviceInfo())
	case IOSScrollSetEnabled:
		args, _ := params.Args()
		enabled := true
		if b := args.Bool("enabled"); b != nil {
			enabled = *b
		}
		iosSetScrollEnabled(enabled)
		m.ok(rw)
	case IOSScrollSetBounceEnabled:
		args, _ := params.Args()
		enabled := true
		if b := args.Bool("enabled"); b != nil {
			enabled = *b
		}
		iosSetBounceEnabled(enabled)
		m.ok(rw)
	case IOSScrollSetIndicatorsEnabled:
		args, _ := params.Args()
		enabled := true
		if b := args.Bool("enabled"); b != nil {
			enabled = *b
		}
		iosSetScrollIndicatorsEnabled(enabled)
		m.ok(rw)
	case IOSNavigationSetBackForwardGestures:
		args, _ := params.Args()
		enabled := false
		if b := args.Bool("enabled"); b != nil {
			enabled = *b
		}
		iosSetBackForwardGesturesEnabled(enabled)
		m.ok(rw)
	case IOSLinksSetPreviewEnabled:
		args, _ := params.Args()
		enabled := true
		if b := args.Bool("enabled"); b != nil {
			enabled = *b
		}
		iosSetLinkPreviewEnabled(enabled)
		m.ok(rw)
	case IOSDebugSetInspectableEnabled:
		args, _ := params.Args()
		enabled := true
		if b := args.Bool("enabled"); b != nil {
			enabled = *b
		}
		iosSetInspectableEnabled(enabled)
		m.ok(rw)
	case IOSUserAgentSet:
		args, _ := params.Args()
		ua := ""
		if s := args.String("ua"); s != nil {
			ua = *s
		} else if s2 := args.String("userAgent"); s2 != nil {
			ua = *s2
		}
		iosSetCustomUserAgent(ua)
		m.ok(rw)
	default:
		m.httpError(rw, "Invalid iOS call:", errUnknownMethod(method))
		return
	}

	m.Info("Runtime call:", "method", "IOS."+iosMethodNames[method])
}

func errUnknownMethod(m int) error { return fmt.Errorf("unknown method: %d", m) }
