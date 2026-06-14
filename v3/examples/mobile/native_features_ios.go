//go:build ios

package main

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// registerNativeFeatures wires the "common:*" events the frontend emits to the
// exported iOS mobile-feature APIs. Results that come back asynchronously
// (torch availability, biometric outcome, …) are emitted by the framework as
// "common:*" custom events the frontend listens for.
func registerNativeFeatures(app *application.App) {
	// Phase A — one-way actions
	app.Event.On("common:share", func(e *application.CustomEvent) {
		application.IOS.Share(payloadJSON(e.Data))
	})
	app.Event.On("common:openURL", func(e *application.CustomEvent) {
		application.IOS.OpenURL(eventString(e.Data, "url"))
	})
	app.Event.On("common:keepAwake", func(e *application.CustomEvent) {
		application.IOS.SetKeepAwake(eventBool(e.Data, "enabled", false))
	})
	app.Event.On("common:torch", func(e *application.CustomEvent) {
		// Ignore the framework's own result echo (it carries "available").
		if m := firstMap(e.Data); m != nil {
			if _, isResult := m["available"]; isResult {
				return
			}
		}
		application.IOS.SetTorch(eventBool(e.Data, "enabled", false))
	})

	// Phase B — state / query (request → response event)
	app.Event.On("common:getSafeArea", func(e *application.CustomEvent) {
		app.Event.Emit("common:safeArea", jsonToMap(application.IOS.SafeAreaJSON()))
	})
	app.Event.On("common:setBrightness", func(e *application.CustomEvent) {
		application.IOS.SetBrightness(eventFloat(e.Data, "value", 0.5))
	})
	app.Event.On("common:getBrightness", func(e *application.CustomEvent) {
		app.Event.Emit("common:brightness", map[string]any{"value": application.IOS.GetBrightness()})
	})
	app.Event.On("common:getAppInfo", func(e *application.CustomEvent) {
		app.Event.Emit("common:appInfo", jsonToMap(application.IOS.AppInfoJSON()))
	})
	app.Event.On("common:setOrientation", func(e *application.CustomEvent) {
		application.IOS.SetOrientation(eventString(e.Data, "mode"))
	})
	app.Event.On("common:getOrientation", func(e *application.CustomEvent) {
		app.Event.Emit("common:orientation", jsonToMap(application.IOS.GetOrientation()))
	})
	app.Event.On("common:setStatusBar", func(e *application.CustomEvent) {
		application.IOS.SetStatusBar(payloadJSON(e.Data))
	})

	// Phase C — async results / permissions
	app.Event.On("common:authenticate", func(e *application.CustomEvent) {
		reason := eventString(e.Data, "reason")
		if reason == "" {
			reason = "Authenticate to continue"
		}
		application.IOS.BiometricAuthenticate(reason)
	})
	app.Event.On("common:notify", func(e *application.CustomEvent) {
		application.IOS.PostNotification(payloadJSON(e.Data))
	})
	app.Event.On("common:secureSet", func(e *application.CustomEvent) {
		application.IOS.SecureSet(eventString(e.Data, "key"), eventString(e.Data, "value"))
	})
	app.Event.On("common:secureGet", func(e *application.CustomEvent) {
		key := eventString(e.Data, "key")
		app.Event.Emit("common:secureValue", map[string]any{
			"key": key, "value": application.IOS.SecureGet(key),
		})
	})
	app.Event.On("common:secureDelete", func(e *application.CustomEvent) {
		application.IOS.SecureDelete(eventString(e.Data, "key"))
	})

	// Phase D — sensors & hardware
	app.Event.On("common:haptic", func(e *application.CustomEvent) {
		application.IOS.Haptic(eventString(e.Data, "type"))
	})
	app.Event.On("common:getLocation", func(e *application.CustomEvent) {
		application.IOS.GetLocation()
	})
	app.Event.On("common:watchMotion", func(e *application.CustomEvent) {
		application.IOS.SetMotion(eventBool(e.Data, "enabled", false))
	})
	app.Event.On("common:watchProximity", func(e *application.CustomEvent) {
		application.IOS.SetProximity(eventBool(e.Data, "enabled", false))
	})
	app.Event.On("common:speak", func(e *application.CustomEvent) {
		application.IOS.Speak(eventString(e.Data, "text"))
	})
	app.Event.On("common:stopSpeak", func(e *application.CustomEvent) {
		application.IOS.StopSpeak()
	})
	app.Event.On("common:getStorage", func(e *application.CustomEvent) {
		app.Event.Emit("common:storage", jsonToMap(application.IOS.StorageJSON()))
	})
	app.Event.On("common:getPower", func(e *application.CustomEvent) {
		app.Event.Emit("common:power", jsonToMap(application.IOS.PowerJSON()))
	})
	app.Event.On("common:getNetwork", func(e *application.CustomEvent) {
		app.Event.Emit("common:network", jsonToMap(application.IOS.NetworkJSON()))
	})
	app.Event.On("common:watchKeyboard", func(e *application.CustomEvent) {
		application.IOS.SetKeyboardWatch(eventBool(e.Data, "enabled", false))
	})
	app.Event.On("common:setScreenProtect", func(e *application.CustomEvent) {
		application.IOS.SetScreenProtect(eventBool(e.Data, "enabled", false))
	})

	// Phase E — camera & background
	app.Event.On("common:capturePhoto", func(e *application.CustomEvent) {
		application.IOS.CapturePhoto()
	})
	app.Event.On("common:captureVideo", func(e *application.CustomEvent) {
		application.IOS.CaptureVideo()
	})
	app.Event.On("ios:beginBackgroundTask", func(e *application.CustomEvent) {
		application.IOS.BeginBackgroundTask(int(eventFloat(e.Data, "seconds", 20)))
	})
	// iOS has no foreground service; a background-task window is the closest analogue.
	app.Event.On("common:startForegroundService", func(e *application.CustomEvent) {
		application.IOS.BeginBackgroundTask(30)
	})
	app.Event.On("common:stopForegroundService", func(e *application.CustomEvent) {
		application.IOS.EndBackgroundTask()
	})
}
