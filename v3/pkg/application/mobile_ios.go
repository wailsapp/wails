//go:build ios

package application

// Mobile is the cross-platform mobile manager. On iOS it dispatches to the IOS
// manager; the compiler verifies iosManager satisfies MobileManager.
var Mobile MobileManager = IOS
