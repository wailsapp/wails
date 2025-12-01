//go:build android && cgo

package application

// --- Haptic Feedback ---

// Vibrate triggers device vibration for the specified duration in milliseconds.
// Note: For cross-platform haptic feedback, use app.Mobile.Vibrate() instead.
func (m *AndroidManager) Vibrate(durationMs int) {
	AndroidVibrate(durationMs)
}

// --- Notifications ---

// ShowToast displays a native Android toast notification.
// Toasts are non-blocking and auto-dismiss after a short duration.
// Note: For cross-platform notifications, consider using app.Mobile.Toast() instead.
func (m *AndroidManager) ShowToast(message string) {
	AndroidShowToast(message)
}

// --- Device Information ---

// GetDeviceInfo returns device information as a JSON string.
// Fields include: platform, manufacturer, model, brand, sdkVersion, release
// Note: For cross-platform device info, use app.Mobile.DeviceInfo() instead.
func (m *AndroidManager) GetDeviceInfo() string {
	return AndroidGetDeviceInfo()
}

// GetScreenInfo returns screen information as a JSON string.
// Fields include: widthPixels, heightPixels, density, densityDpi, scaledDensity
// Note: For cross-platform screen info, use app.Mobile.ScreenInfo() instead.
func (m *AndroidManager) GetScreenInfo() string {
	return AndroidGetScreenInfo()
}

// --- Dark Mode ---

// IsDarkMode returns true if the system is in dark mode.
func (m *AndroidManager) IsDarkMode() bool {
	globalAppLock.RLock()
	app := globalApp
	globalAppLock.RUnlock()
	if app != nil {
		return app.isDarkMode()
	}
	return false
}

// --- Logging ---

// Log writes a log message with the specified level to Android logcat.
// Levels: "debug", "info", "warn", "error"
func (m *AndroidManager) Log(level, format string, args ...interface{}) {
	androidLogf(level, format, args...)
}

// Debug logs a debug message to Android logcat.
func (m *AndroidManager) Debug(format string, args ...interface{}) {
	androidLogf("debug", format, args...)
}

// Info logs an info message to Android logcat.
func (m *AndroidManager) Info(format string, args ...interface{}) {
	androidLogf("info", format, args...)
}

// Warn logs a warning message to Android logcat.
func (m *AndroidManager) Warn(format string, args ...interface{}) {
	androidLogf("warn", format, args...)
}

// Error logs an error message to Android logcat.
func (m *AndroidManager) Error(format string, args ...interface{}) {
	androidLogf("error", format, args...)
}

// --- Feature Detection ---

// IsAndroid returns true since this is the Android implementation.
func (m *AndroidManager) IsAndroid() bool {
	return true
}
