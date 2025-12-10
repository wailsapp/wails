//go:build ios

package main

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// registerIOSRuntimeEventHandlers registers Go-side event listeners that mutate iOS WKWebView at runtime.
func registerIOSRuntimeEventHandlers(app *application.App) {
	// Helper to fetch boolean from event data. Accepts {"enabled":bool} or a bare bool.
	getBool := func(data any, key string, def bool) bool {
		switch v := data.(type) {
		case bool:
			return v
		case map[string]any:
			if raw, ok := v[key]; ok {
				if b, ok := raw.(bool); ok {
					return b
				}
			}
		}
		return def
	}
	// Helper to fetch string from event data. Accepts {"ua":string} or bare string.
	getString := func(data any, key string) string {
		switch v := data.(type) {
		case string:
			return v
		case map[string]any:
			if raw, ok := v[key]; ok {
				if s, ok := raw.(string); ok {
					return s
				}
			}
		}
		return ""
	}

	app.Event.On("ios:setScrollEnabled", func(e *application.CustomEvent) {
		application.IOSSetScrollEnabled(getBool(e.Data, "enabled", true))
	})
	app.Event.On("ios:setBounceEnabled", func(e *application.CustomEvent) {
		application.IOSSetBounceEnabled(getBool(e.Data, "enabled", true))
	})
	app.Event.On("ios:setScrollIndicatorsEnabled", func(e *application.CustomEvent) {
		application.IOSSetScrollIndicatorsEnabled(getBool(e.Data, "enabled", true))
	})
	app.Event.On("ios:setBackForwardGesturesEnabled", func(e *application.CustomEvent) {
		application.IOSSetBackForwardGesturesEnabled(getBool(e.Data, "enabled", false))
	})
	app.Event.On("ios:setLinkPreviewEnabled", func(e *application.CustomEvent) {
		application.IOSSetLinkPreviewEnabled(getBool(e.Data, "enabled", true))
	})
	app.Event.On("ios:setInspectableEnabled", func(e *application.CustomEvent) {
		application.IOSSetInspectableEnabled(getBool(e.Data, "enabled", true))
	})
	app.Event.On("ios:setCustomUserAgent", func(e *application.CustomEvent) {
		ua := getString(e.Data, "ua")
		application.IOSSetCustomUserAgent(ua)
	})
}
