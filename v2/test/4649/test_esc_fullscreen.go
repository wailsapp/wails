//go:build darwin
// +build darwin

package main

import (
	"fmt"
	"os"
)

// This is a manual test program for issue #4649.
// Run on macOS: go run v2/test/4649/test_esc_fullscreen.go
//
// Expected behavior when mac.Options.DisableEscapeExitsFullscreen is true:
// 1. Window opens in fullscreen mode
// 2. Pressing Esc should NOT exit fullscreen
// 3. Pressing Cmd+Ctrl+F should toggle fullscreen (exit)
// 4. When DisableEscapeExitsFullscreen is false (default), Esc exits fullscreen normally
//
// The fix adds a DisableEscapeExitsFullscreen option to mac.Options.
// When true, WailsWindow.cancelOperation: swallows the Esc event while
// fullscreen, allowing web content to handle it (e.g., closing modals).
// When false (the default), system behavior is preserved unchanged.

func main() {
	fmt.Println("Test for issue #4649: Esc key behaviour in fullscreen on macOS")
	fmt.Println("")
	fmt.Println("This test must be run on macOS.")
	fmt.Println("Set mac.Options{DisableEscapeExitsFullscreen: true} to opt in.")
	fmt.Println("")
	fmt.Println("Manual verification steps (with DisableEscapeExitsFullscreen: true):")
	fmt.Println("1. Create a wails v2 app with a fullscreen window and the option set")
	fmt.Println("2. Add an HTML modal dialog")
	fmt.Println("3. Press Esc - the modal should close but NOT exit fullscreen")
	fmt.Println("4. The window should remain in fullscreen mode")
	fmt.Println("")
	fmt.Println("Manual verification steps (with DisableEscapeExitsFullscreen: false / default):")
	fmt.Println("1. Same app without the option")
	fmt.Println("2. Press Esc in fullscreen - the window should EXIT fullscreen (default macOS behavior)")
	os.Exit(0)
}
