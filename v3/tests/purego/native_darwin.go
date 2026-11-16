//go:build darwin

package main

import (
	"fmt"

	"github.com/ebitengine/purego/objc"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// checkPlatform reads back native state the Go API does not expose. It uses
// purego/objc, so it works with either backend.
func checkPlatform(app *application.App, window *application.WebviewWindow) error {
	hasShadow := func(w *application.WebviewWindow) func() bool {
		return func() bool {
			var shadow bool
			application.InvokeSync(func() {
				shadow = objc.Send[bool](objc.ID(w.NativeWindow()), objc.RegisterName("hasShadow"))
			})
			return shadow
		}
	}
	if !eventually(hasShadow(window)) {
		return fmt.Errorf("shown window has no shadow")
	}
	plain := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "No shadow", Width: 200, Height: 120, Hidden: true,
		Mac: application.MacWindow{DisableShadow: true},
	})
	defer plain.Close()
	plain.Show()
	if !eventually(plain.IsVisible) {
		return fmt.Errorf("DisableShadow window is not visible")
	}
	if eventually(hasShadow(plain)) {
		return fmt.Errorf("DisableShadow window has a shadow after Show")
	}
	return nil
}
