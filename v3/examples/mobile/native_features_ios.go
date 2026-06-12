//go:build ios

package main

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// registerNativeFeatures wires the "native:*" events the frontend emits to the
// exported iOS mobile-feature APIs. Results that come back asynchronously
// (torch availability, biometric outcome, …) are emitted by the framework as
// "native:*" custom events the frontend listens for.
func registerNativeFeatures(app *application.App) {
	// Phase A — one-way actions
	app.Event.On("native:share", func(e *application.CustomEvent) {
		application.IOSShare(payloadJSON(e.Data))
	})
	app.Event.On("native:openURL", func(e *application.CustomEvent) {
		application.IOSOpenURL(eventString(e.Data, "url"))
	})
	app.Event.On("native:keepAwake", func(e *application.CustomEvent) {
		application.IOSSetKeepAwake(eventBool(e.Data, "enabled", false))
	})
	app.Event.On("native:torch", func(e *application.CustomEvent) {
		// Ignore the framework's own result echo (it carries "available").
		if m := firstMap(e.Data); m != nil {
			if _, isResult := m["available"]; isResult {
				return
			}
		}
		application.IOSSetTorch(eventBool(e.Data, "enabled", false))
	})

	// Phase B — state / query (request → response event)
	app.Event.On("native:getSafeArea", func(e *application.CustomEvent) {
		app.Event.Emit("native:safeArea", jsonToMap(application.IOSSafeAreaJSON()))
	})
	app.Event.On("native:setBrightness", func(e *application.CustomEvent) {
		application.IOSSetBrightness(eventFloat(e.Data, "value", 0.5))
	})
	app.Event.On("native:getBrightness", func(e *application.CustomEvent) {
		app.Event.Emit("native:brightness", map[string]any{"value": application.IOSGetBrightness()})
	})
	app.Event.On("native:getAppInfo", func(e *application.CustomEvent) {
		app.Event.Emit("native:appInfo", jsonToMap(application.IOSAppInfoJSON()))
	})
	app.Event.On("native:setOrientation", func(e *application.CustomEvent) {
		application.IOSSetOrientation(eventString(e.Data, "mode"))
	})
	app.Event.On("native:getOrientation", func(e *application.CustomEvent) {
		app.Event.Emit("native:orientation", jsonToMap(application.IOSGetOrientation()))
	})
	app.Event.On("native:setStatusBar", func(e *application.CustomEvent) {
		application.IOSSetStatusBar(payloadJSON(e.Data))
	})

	// Phase C — async results / permissions
	app.Event.On("native:authenticate", func(e *application.CustomEvent) {
		reason := eventString(e.Data, "reason")
		if reason == "" {
			reason = "Authenticate to continue"
		}
		application.IOSBiometricAuthenticate(reason)
	})
	app.Event.On("native:notify", func(e *application.CustomEvent) {
		application.IOSPostNotification(payloadJSON(e.Data))
	})
	app.Event.On("native:secureSet", func(e *application.CustomEvent) {
		application.IOSSecureSet(eventString(e.Data, "key"), eventString(e.Data, "value"))
	})
	app.Event.On("native:secureGet", func(e *application.CustomEvent) {
		key := eventString(e.Data, "key")
		app.Event.Emit("native:secureValue", map[string]any{
			"key": key, "value": application.IOSSecureGet(key),
		})
	})
	app.Event.On("native:secureDelete", func(e *application.CustomEvent) {
		application.IOSSecureDelete(eventString(e.Data, "key"))
	})
}
