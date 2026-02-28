//go:build android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

// registerAndroidRuntimeEventHandlers registers Go-side event listeners that mutate the Android WebView at runtime.
func registerAndroidRuntimeEventHandlers(app *application.App) {
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

	app.Event.On("android:setScrollEnabled", func(e *application.CustomEvent) {
		application.AndroidSetScrollEnabled(getBool(e.Data, "enabled", true))
	})
	app.Event.On("android:setBounceEnabled", func(e *application.CustomEvent) {
		application.AndroidSetBounceEnabled(getBool(e.Data, "enabled", true))
	})
	app.Event.On("android:setScrollIndicatorsEnabled", func(e *application.CustomEvent) {
		application.AndroidSetScrollIndicatorsEnabled(getBool(e.Data, "enabled", true))
	})
	app.Event.On("android:setBackForwardGesturesEnabled", func(e *application.CustomEvent) {
		application.AndroidSetBackForwardGesturesEnabled(getBool(e.Data, "enabled", false))
	})
	app.Event.On("android:setLinkPreviewEnabled", func(e *application.CustomEvent) {
		application.AndroidSetLinkPreviewEnabled(getBool(e.Data, "enabled", true))
	})
	app.Event.On("android:setCustomUserAgent", func(e *application.CustomEvent) {
		application.AndroidSetCustomUserAgent(getString(e.Data, "ua"))
	})
}
