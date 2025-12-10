//go:build android

package application

import "fmt"

func init() {
	fmt.Println("🤖 [init_android.go] START init()")
	// On Android, we don't call runtime.LockOSThread()
	// The Android runtime handles thread management via JNI
	// and calling LockOSThread can interfere with the JNI environment
	fmt.Println("🤖 [init_android.go] END init() - no LockOSThread on Android")
}
