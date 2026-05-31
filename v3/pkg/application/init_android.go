//go:build android

package application

func init() {
	// On Android, we don't call runtime.LockOSThread()
	// The Android runtime handles thread management via JNI
	// and calling LockOSThread can interfere with the JNI environment
}
