//go:build !android || !cgo

package application

// --- Haptic Feedback ---

// Vibrate is a no-op on non-Android platforms.
func (m *AndroidManager) Vibrate(durationMs int) {}

// --- Notifications ---

// ShowToast is a no-op on non-Android platforms.
func (m *AndroidManager) ShowToast(message string) {}

// --- Device Information ---

// GetDeviceInfo returns an empty JSON object on non-Android platforms.
func (m *AndroidManager) GetDeviceInfo() string {
	return `{"platform":"unknown"}`
}

// GetScreenInfo returns an empty JSON object on non-Android platforms.
func (m *AndroidManager) GetScreenInfo() string {
	return `{}`
}

// --- Dark Mode ---

// IsDarkMode returns false on non-Android platforms.
func (m *AndroidManager) IsDarkMode() bool {
	return false
}

// --- Logging ---

// Log is a no-op on non-Android platforms.
func (m *AndroidManager) Log(level, format string, args ...interface{}) {}

// Debug is a no-op on non-Android platforms.
func (m *AndroidManager) Debug(format string, args ...interface{}) {}

// Info is a no-op on non-Android platforms.
func (m *AndroidManager) Info(format string, args ...interface{}) {}

// Warn is a no-op on non-Android platforms.
func (m *AndroidManager) Warn(format string, args ...interface{}) {}

// Error is a no-op on non-Android platforms.
func (m *AndroidManager) Error(format string, args ...interface{}) {}

// --- Feature Detection ---

// IsAndroid returns false on non-Android platforms.
func (m *AndroidManager) IsAndroid() bool {
	return false
}
