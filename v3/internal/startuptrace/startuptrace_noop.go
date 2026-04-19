//go:build !wails_trace_startup

// Package startuptrace records named timestamps during application startup
// and emits them as Chrome trace JSON when the wails_trace_startup build tag
// is set. Without the tag, all functions compile to empty bodies and the
// inliner removes them entirely.
package startuptrace

// Mark records a named event at the current time.
func Mark(name string) {}

// MarkWindow records a named event scoped to a specific window.
func MarkWindow(windowID uint, name string) {}

// Flush writes the trace to the file named by WAILS_TRACE_STARTUP_OUTPUT.
func Flush() {}

// Enabled reports whether tracing is compiled in.
func Enabled() bool { return false }
