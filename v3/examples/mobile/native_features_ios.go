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
}
