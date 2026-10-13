//go:build windows

package tui

// runSpinner on Windows falls back to plain progress lines: the glyph
// terminal library does not support Windows.
func runSpinner(message string, fn func() error) error {
	return runPlain(message, fn)
}
