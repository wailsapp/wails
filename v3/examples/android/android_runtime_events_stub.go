//go:build !android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

// registerAndroidRuntimeEventHandlers is a no-op on non-Android platforms.
func registerAndroidRuntimeEventHandlers(app *application.App) {}
