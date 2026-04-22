//go:build android

package main

import "github.com/wailsapp/wails/v3/pkg/application"

func init() {
	// Register main function to be called when the Android app initializes
	// This is necessary because in c-shared build mode, main() is not automatically called
	application.RegisterAndroidMain(main)
}
