//go:build android && production

package application

// androidVerboseLogging is disabled in production builds: the framework's
// internal diagnostic logging compiles to a no-op.
const androidVerboseLogging = false
