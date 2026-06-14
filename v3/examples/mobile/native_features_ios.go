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

	// Phase D — sensors & hardware
	app.Event.On("native:haptic", func(e *application.CustomEvent) {
		application.IOSHaptic(eventString(e.Data, "type"))
	})
	app.Event.On("native:getLocation", func(e *application.CustomEvent) {
		application.IOSGetLocation()
	})
	app.Event.On("native:watchMotion", func(e *application.CustomEvent) {
		application.IOSSetMotion(eventBool(e.Data, "enabled", false))
	})
	app.Event.On("native:watchProximity", func(e *application.CustomEvent) {
		application.IOSSetProximity(eventBool(e.Data, "enabled", false))
	})
	app.Event.On("native:speak", func(e *application.CustomEvent) {
		application.IOSSpeak(eventString(e.Data, "text"))
	})
	app.Event.On("native:stopSpeak", func(e *application.CustomEvent) {
		application.IOSStopSpeak()
	})
	app.Event.On("native:getStorage", func(e *application.CustomEvent) {
		app.Event.Emit("native:storage", jsonToMap(application.IOSStorageJSON()))
	})
	app.Event.On("native:getPower", func(e *application.CustomEvent) {
		app.Event.Emit("native:power", jsonToMap(application.IOSPowerJSON()))
	})
	app.Event.On("native:getNetwork", func(e *application.CustomEvent) {
		app.Event.Emit("native:network", jsonToMap(application.IOSNetworkJSON()))
	})
	app.Event.On("native:watchKeyboard", func(e *application.CustomEvent) {
		application.IOSSetKeyboardWatch(eventBool(e.Data, "enabled", false))
	})
	app.Event.On("native:setScreenProtect", func(e *application.CustomEvent) {
		application.IOSSetScreenProtect(eventBool(e.Data, "enabled", false))
	})

	// Phase E — camera & background
	app.Event.On("native:capturePhoto", func(e *application.CustomEvent) {
		application.IOSCapturePhoto()
	})
	app.Event.On("native:captureVideo", func(e *application.CustomEvent) {
		application.IOSCaptureVideo()
	})
	app.Event.On("native:beginBackgroundTask", func(e *application.CustomEvent) {
		application.IOSBeginBackgroundTask(int(eventFloat(e.Data, "seconds", 20)))
	})
	// iOS has no foreground service; a background-task window is the closest analogue.
	app.Event.On("native:startForegroundService", func(e *application.CustomEvent) {
		application.IOSBeginBackgroundTask(30)
	})
	app.Event.On("native:stopForegroundService", func(e *application.CustomEvent) {
		application.IOSEndBackgroundTask()
	})
}
