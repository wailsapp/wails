//go:build android && !production

package application

// androidVerboseLogging enables the framework's internal diagnostic logging
// in debug builds. See android_logging_production.go for the release value.
const androidVerboseLogging = true
