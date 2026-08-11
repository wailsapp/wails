//go:build android

package main

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// registerNativeFeatures wires the "common:*" events the frontend emits to the
// exported Android mobile-feature APIs. Asynchronous results come back as
// "common:*" custom events via the framework's nativeEmitEvent bridge.
func registerNativeFeatures(app *application.App) {
	// Phase A — one-way actions
	app.Event.On("common:share", func(e *application.CustomEvent) {
		application.Android.Share(payloadJSON(e.Data))
	})
	app.Event.On("common:openURL", func(e *application.CustomEvent) {
		application.Android.OpenURL(eventString(e.Data, "url"))
	})
	app.Event.On("common:keepAwake", func(e *application.CustomEvent) {
		application.Android.SetKeepAwake(eventBool(e.Data, "enabled", false))
	})
	app.Event.On("common:torch", func(e *application.CustomEvent) {
		if m := firstMap(e.Data); m != nil {
			if _, isResult := m["available"]; isResult {
				return
			}
		}
		application.Android.SetTorch(eventBool(e.Data, "enabled", false))
	})

	// Phase B — state / query (request → response event)
	app.Event.On("common:getSafeArea", func(e *application.CustomEvent) {
		app.Event.Emit("common:safeArea", jsonToMap(application.Android.SafeAreaJSON()))
	})
	app.Event.On("common:setBrightness", func(e *application.CustomEvent) {
		application.Android.SetBrightness(int(eventFloat(e.Data, "value", 0.5) * 100))
	})
	app.Event.On("common:getBrightness", func(e *application.CustomEvent) {
		app.Event.Emit("common:brightness", jsonToMap(application.Android.BrightnessJSON()))
	})
	app.Event.On("common:getAppInfo", func(e *application.CustomEvent) {
		app.Event.Emit("common:appInfo", jsonToMap(application.Android.AppInfoJSON()))
	})
	app.Event.On("common:setOrientation", func(e *application.CustomEvent) {
		application.Android.SetOrientation(eventString(e.Data, "mode"))
	})
	app.Event.On("common:getOrientation", func(e *application.CustomEvent) {
		app.Event.Emit("common:orientation", jsonToMap(application.Android.OrientationJSON()))
	})
	app.Event.On("common:setStatusBar", func(e *application.CustomEvent) {
		application.Android.SetStatusBar(payloadJSON(e.Data))
	})

	// Phase C — async results / permissions
	app.Event.On("common:authenticate", func(e *application.CustomEvent) {
		reason := eventString(e.Data, "reason")
		if reason == "" {
			reason = "Authenticate to continue"
		}
		application.Android.BiometricAuthenticate(reason)
	})
	app.Event.On("common:notify", func(e *application.CustomEvent) {
		application.Android.Notify(payloadJSON(e.Data))
	})
	app.Event.On("common:secureSet", func(e *application.CustomEvent) {
		application.Android.SecureSet(payloadJSON(e.Data))
	})
	app.Event.On("common:secureGet", func(e *application.CustomEvent) {
		key := eventString(e.Data, "key")
		app.Event.Emit("common:secureValue", map[string]any{
			"key": key, "value": application.Android.SecureGet(key),
		})
	})
	app.Event.On("common:secureDelete", func(e *application.CustomEvent) {
		application.Android.SecureDelete(eventString(e.Data, "key"))
	})

	// Phase D — sensors & hardware
	app.Event.On("common:haptic", func(e *application.CustomEvent) {
		application.Android.Haptic(eventString(e.Data, "type"))
	})
	app.Event.On("common:getLocation", func(e *application.CustomEvent) {
		application.Android.GetLocation()
	})
	app.Event.On("common:watchMotion", func(e *application.CustomEvent) {
		application.Android.SetMotion(eventBool(e.Data, "enabled", false))
	})
	app.Event.On("common:watchProximity", func(e *application.CustomEvent) {
		application.Android.SetProximity(eventBool(e.Data, "enabled", false))
	})
	app.Event.On("common:speak", func(e *application.CustomEvent) {
		application.Android.Speak(eventString(e.Data, "text"))
	})
	app.Event.On("common:stopSpeak", func(e *application.CustomEvent) {
		application.Android.StopSpeak()
	})
	app.Event.On("common:getStorage", func(e *application.CustomEvent) {
		app.Event.Emit("common:storage", jsonToMap(application.Android.StorageJSON()))
	})
	app.Event.On("common:getPower", func(e *application.CustomEvent) {
		app.Event.Emit("common:power", jsonToMap(application.Android.PowerJSON()))
	})
	app.Event.On("common:getNetwork", func(e *application.CustomEvent) {
		app.Event.Emit("common:network", jsonToMap(application.Android.NetworkJSON()))
	})
	app.Event.On("common:watchKeyboard", func(e *application.CustomEvent) {
		application.Android.SetKeyboardWatch(eventBool(e.Data, "enabled", false))
	})
	app.Event.On("common:setScreenProtect", func(e *application.CustomEvent) {
		application.Android.SetScreenProtect(eventBool(e.Data, "enabled", false))
	})

	// Phase E — camera & background
	app.Event.On("common:capturePhoto", func(e *application.CustomEvent) {
		application.Android.CapturePhoto()
	})
	app.Event.On("common:captureVideo", func(e *application.CustomEvent) {
		application.Android.CaptureVideo()
	})
	app.Event.On("common:startForegroundService", func(e *application.CustomEvent) {
		application.Android.StartForegroundService(payloadJSON(e.Data))
	})
	app.Event.On("common:stopForegroundService", func(e *application.CustomEvent) {
		application.Android.StopForegroundService()
	})
}
