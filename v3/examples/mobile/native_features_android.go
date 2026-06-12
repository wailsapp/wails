//go:build android

package main

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// registerNativeFeatures wires the "native:*" events the frontend emits to the
// exported Android mobile-feature APIs. Asynchronous results come back as
// "native:*" custom events via the framework's nativeEmitEvent bridge.
func registerNativeFeatures(app *application.App) {
	// Phase A — one-way actions
	app.Event.On("native:share", func(e *application.CustomEvent) {
		application.AndroidShare(payloadJSON(e.Data))
	})
	app.Event.On("native:openURL", func(e *application.CustomEvent) {
		application.AndroidOpenURL(eventString(e.Data, "url"))
	})
	app.Event.On("native:keepAwake", func(e *application.CustomEvent) {
		application.AndroidSetKeepAwake(eventBool(e.Data, "enabled", false))
	})
	app.Event.On("native:torch", func(e *application.CustomEvent) {
		if m := firstMap(e.Data); m != nil {
			if _, isResult := m["available"]; isResult {
				return
			}
		}
		application.AndroidSetTorch(eventBool(e.Data, "enabled", false))
	})

	// Phase B — state / query (request → response event)
	app.Event.On("native:getSafeArea", func(e *application.CustomEvent) {
		app.Event.Emit("native:safeArea", jsonToMap(application.AndroidSafeAreaJSON()))
	})
	app.Event.On("native:setBrightness", func(e *application.CustomEvent) {
		application.AndroidSetBrightness(int(eventFloat(e.Data, "value", 0.5) * 100))
	})
	app.Event.On("native:getBrightness", func(e *application.CustomEvent) {
		app.Event.Emit("native:brightness", jsonToMap(application.AndroidBrightnessJSON()))
	})
	app.Event.On("native:getAppInfo", func(e *application.CustomEvent) {
		app.Event.Emit("native:appInfo", jsonToMap(application.AndroidAppInfoJSON()))
	})
	app.Event.On("native:setOrientation", func(e *application.CustomEvent) {
		application.AndroidSetOrientation(eventString(e.Data, "mode"))
	})
	app.Event.On("native:getOrientation", func(e *application.CustomEvent) {
		app.Event.Emit("native:orientation", jsonToMap(application.AndroidOrientationJSON()))
	})
	app.Event.On("native:setStatusBar", func(e *application.CustomEvent) {
		application.AndroidSetStatusBar(payloadJSON(e.Data))
	})
}
