//go:build android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func init() {
	// On Android the app is built as a c-shared library, so main() is not
	// called automatically. Register it to run when the Android Activity
	// initialises the native library.
	application.RegisterAndroidMain(main)
}
