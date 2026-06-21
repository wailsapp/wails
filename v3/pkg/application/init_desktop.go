//go:build !ios

package application

import "runtime"

func init() {
	// Lock the main thread for desktop platforms
	// This ensures UI operations happen on the main thread
	runtime.LockOSThread()
}