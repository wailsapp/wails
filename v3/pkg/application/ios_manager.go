package application

// IOSManager provides iOS-specific platform functionality.
// On non-iOS platforms, all methods are no-ops.
// Access via the package-level IOS variable: application.IOS.HapticsImpact("medium")
type IOSManager struct{}

// IOS is the package-level iOS manager instance.
// Use this to access iOS-specific features: application.IOS.HapticsImpact("medium")
var IOS = &IOSManager{}
