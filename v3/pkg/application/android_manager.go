package application

// AndroidManager provides Android-specific platform functionality.
// On non-Android platforms, all methods are no-ops or return sensible defaults.
// Access via the package-level Android variable: application.Android.Vibrate(50)
type AndroidManager struct{}

// Android is the package-level Android manager instance.
// Use this to access Android-specific features: application.Android.Vibrate(50)
var Android = &AndroidManager{}
