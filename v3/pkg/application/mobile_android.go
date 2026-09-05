//go:build android

package application

// Mobile is the cross-platform mobile manager. On Android it dispatches to the
// Android manager; the compiler verifies androidManager satisfies MobileManager.
var Mobile MobileManager = Android
