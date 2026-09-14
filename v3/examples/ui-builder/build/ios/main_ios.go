//go:build ios

package main

import (
	"C"
)

// For iOS builds we export a function the WailsAppDelegate can call.
// This wrapper keeps the user's main.go unmodified.

// WailsIOSMain runs the user's main() (application.New / app.Run). It is invoked
// by the WailsAppDelegate from didFinishLaunchingWithOptions — i.e. only AFTER
// UIKit has launched — on a BACKGROUND thread, so the Go runtime never starts
// concurrently with UIApplicationMain (that race intermittently corrupts the
// FrontBoard launch handshake on a physical device → blank cold launch /
// scene-create watchdog 0x8BADF00D). Keeping all app setup off the OS main
// thread also leaves UIApplicationMain unobstructed on the main thread.
//
//export WailsIOSMain
func WailsIOSMain() {
	main()
}
