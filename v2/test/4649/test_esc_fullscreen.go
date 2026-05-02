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
// The fix is opt-in via mac.Options.DisableEscapeExitsFullscreen.
// Default behaviour (flag = false) preserves standard macOS fullscreen:
// Esc exits fullscreen.
//
// Expected behaviour with the flag enabled (DisableEscapeExitsFullscreen: true):
// 1. Window opens in fullscreen mode
// 2. Pressing Esc should NOT exit fullscreen — the keypress is swallowed
//    so web content (modals, custom Esc handlers) can consume it
// 3. Pressing Cmd+Ctrl+F should still toggle fullscreen (exit)
// 4. When NOT in fullscreen, Esc should still work normally
//
// Implementation: WailsWindow.cancelOperation: checks the
// disableEscapeExitsFullscreen property (set from
// mac.Options.DisableEscapeExitsFullscreen at window creation) AND
// the fullscreen styleMask before swallowing the keypress.

func main() {
	fmt.Println("Test for issue #4649: opt-in flag to keep fullscreen on Esc (macOS)")
	fmt.Println("")
	fmt.Println("Set mac.Options.DisableEscapeExitsFullscreen = true to opt in.")
	fmt.Println("")
	fmt.Println("Manual verification (opt-in):")
	fmt.Println("1. Configure &options.App{ Mac: &mac.Options{ DisableEscapeExitsFullscreen: true } }")
	fmt.Println("2. Open a fullscreen window with an HTML modal dialog")
	fmt.Println("3. Press Esc - the modal should close but NOT exit fullscreen")
	fmt.Println("")
	fmt.Println("Manual verification (default, flag unset):")
	fmt.Println("1. Same app but with the flag omitted/false")
	fmt.Println("2. Press Esc - the window should exit fullscreen (standard macOS behaviour)")
	os.Exit(0)
}
