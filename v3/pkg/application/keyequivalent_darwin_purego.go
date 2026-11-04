//go:build darwin && purego && !ios && !server

package application

import "unicode/utf16"

// Key equivalents under purego.
//
// WKWebView swallows some key presses (Ctrl+Tab, Cmd+Shift+K) before the
// window's keyDown: is reached, so a key binding for them only fires if the
// window answers performKeyEquivalent: first. The cgo backend does this in
// webview_window_darwin.m, and this file is its purego port: the same naming
// of the press, the same parse-then-lookup, the same YES-or-super answer.

// keyEquivalentAccelerator names an AppKit key press the way a key binding is
// written. It is the pure half of cgo acceleratorStringFromKeyEvent: the
// characters reduce to their first UTF-16 unit, as the cgo side passes
// unichar, and the naming itself is macAccelerator's.
func keyEquivalentAccelerator(keyCode uint16, modifiers uint, characters string) string {
	var character rune
	hasCharacter := characters != ""
	if hasCharacter {
		character = rune(utf16.Encode([]rune(characters))[0])
	}
	return macAccelerator(keyCode, modifiers, character, hasCharacter)
}

// acceleratorFromNSEvent reads the three things AppKit reports a press with
// (key code, modifier flags, characters) and names the press. Empty when the
// press has no name, in which case nothing is dispatched.
func acceleratorFromNSEvent(event id) string {
	// keyCode is an unsigned short: read it as a word and truncate, as the
	// rest of this file reads small integer results.
	return keyEquivalentAccelerator(
		uint16(get[uint](event, "keyCode")),
		get[uint](event, "modifierFlags"),
		event.send("characters").string(),
	)
}

// keyEquivalentConsumed reports whether a key binding consumes the key
// equivalent. It mirrors cgo processWindowKeyEquivalent: a name that does not
// parse is not consumed, and one that does is looked up in the normalised form
// the bindings are stored under. lookup is the binding table; in production it
// is WebviewWindow.processKeyBinding, which runs the binding synchronously.
func keyEquivalentConsumed(accelerator string, lookup func(string) bool) bool {
	if accelerator == "" {
		return false
	}
	parsed, err := parseAccelerator(accelerator)
	if err != nil {
		return false
	}
	return lookup(parsed.String())
}

// processKeyEquivalent is the purego processWindowKeyEquivalent: it resolves
// the window by id and asks its bindings about the key equivalent. It runs on
// the main thread, from AppKit's performKeyEquivalent:, and so the binding runs
// before the responder chain sees the press.
func processKeyEquivalent(windowID uint, accelerator string) bool {
	if globalApplication == nil {
		return false
	}
	globalApplication.windowsLock.RLock()
	window, ok := globalApplication.windows[windowID]
	globalApplication.windowsLock.RUnlock()
	if !ok {
		return false
	}
	webviewWindow, ok := window.(*WebviewWindow)
	if !ok {
		return false
	}
	return keyEquivalentConsumed(accelerator, webviewWindow.processKeyBinding)
}
