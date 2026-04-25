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
// Expected behavior:
// 1. Window opens in fullscreen mode
// 2. Pressing Esc should NOT exit fullscreen
// 3. Pressing Cmd+Ctrl+F should toggle fullscreen (exit)
// 4. When not in fullscreen, Esc should work normally (cancel operations)
//
// The test verifies that WailsWindow.cancelOperation: is overridden to
// prevent Esc from exiting fullscreen mode, allowing web content to handle
// the Esc key (e.g., closing modals).

func main() {
	fmt.Println("Test for issue #4649: Esc key should not exit fullscreen on macOS")
	fmt.Println("")
	fmt.Println("This test must be run on macOS.")
	fmt.Println("The fix overrides cancelOperation: in WailsWindow to check")
	fmt.Println("if the window is in fullscreen mode before allowing the default")
	fmt.Println("Esc behavior (which exits fullscreen).")
	fmt.Println("")
	fmt.Println("Manual verification steps:")
	fmt.Println("1. Create a wails v2 app with a fullscreen window")
	fmt.Println("2. Add an HTML modal dialog")
	fmt.Println("3. Press Esc - the modal should close but NOT exit fullscreen")
	fmt.Println("4. The window should remain in fullscreen mode")
	os.Exit(0)
}
