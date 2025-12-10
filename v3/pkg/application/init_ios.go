//go:build ios

package application

import "fmt"

func init() {
	fmt.Println("🔵 [init_ios.go] START init()")
	// On iOS, we don't call runtime.LockOSThread()
	// The iOS runtime handles thread management differently
	// and calling LockOSThread can interfere with signal handling
	fmt.Println("🔵 [init_ios.go] END init() - no LockOSThread on iOS")
}