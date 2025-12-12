//go:build ios

package application

func init() {
	// On iOS, we don't call runtime.LockOSThread()
	// The iOS runtime handles thread management differently
	// and calling LockOSThread can interfere with signal handling
}
