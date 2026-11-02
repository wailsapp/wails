//go:build darwin && purego && !ios && !server && private_mac_apis

package application

import "strings"

// This file is the CGO-free counterpart of mac_private_api_darwin.go. It is the
// only purego file that may reference undocumented macOS selectors and KVC
// keys, so App Store builds exclude it with the default !private_mac_apis
// selection of mac_public_api_darwin_purego.go.
//
// Every function here must be called on the main thread, exactly like the
// cgo bridge it mirrors.

// wailsPrivateSetWebviewTransparent draws the webview without a background.
// WKWebView has no public API for this; drawsBackground is a private property
// that WebKit has honoured since WebKit1.
func wailsPrivateSetWebviewTransparent(webView id) {
	wailsSetKVCValue(webView, "drawsBackground", nsNumberBool(false))
}

// wailsPrivateSetWebviewBackgroundColour sets the webview background colour
// through the private backgroundColor key.
func wailsPrivateSetWebviewBackgroundColour(webView id, r, g, b, alpha int) {
	if webView.isNil() {
		return
	}
	colour := class("NSColor").send("colorWithRed:green:blue:alpha:",
		CGFloat(r)/255, CGFloat(g)/255, CGFloat(b)/255, CGFloat(alpha)/255)
	wailsSetKVCValue(webView, "backgroundColor", colour)
}

// wailsPrivateSetGlassStyle applies Wails' historic glass style mapping:
// vibrant borrows the light style, and light/dark are passed through as raw
// style values even though only regular (0) and clear (1) are documented.
func wailsPrivateSetGlassStyle(glassView id, style MacLiquidGlassStyle) {
	if !respondsTo(glassView, "setStyle:") {
		return
	}
	nativeStyle := style
	if style == LiquidGlassStyleVibrant {
		nativeStyle = LiquidGlassStyleLight
	}
	wailsSetKVCValue(glassView, "style", nsNumberInt(int(nativeStyle)))
}

// wailsPrivateSetGlassGrouping groups glass views across windows through the
// private setGroupIdentifier:/setGroupName: and groupSpacing selectors.
func wailsPrivateSetGlassGrouping(glassView id, groupID string, groupSpacing float64) {
	if glassView.isNil() {
		return
	}
	if groupID != "" {
		identifier := nsString(groupID)
		switch {
		case respondsTo(glassView, "setGroupIdentifier:"):
			glassView.send("setGroupIdentifier:", identifier)
		case respondsTo(glassView, "setGroupName:"):
			glassView.send("setGroupName:", identifier)
		}
	}
	if groupSpacing > 0 {
		wailsSetKVCValue(glassView, "groupSpacing", nsNumberDouble(groupSpacing))
	}
}

// wailsPrivateOpenWebInspector shows the Web Inspector through the private
// -[WKWebView _inspector] selector. A missing selector raises an uncatchable
// NSException in Go, so every send is guarded.
func wailsPrivateOpenWebInspector(webView id) {
	if webView.isNil() || !respondsTo(webView, "_inspector") {
		return
	}
	inspector := webView.send("_inspector")
	if !inspector.isNil() && respondsTo(inspector, "show") {
		inspector.send("show")
	}
}

// wailsPrivateEnableWebInspector enables developer extras through the private
// developerExtrasEnabled preference.
func wailsPrivateEnableWebInspector(webView id) {
	if webView.isNil() {
		return
	}
	prefs := webView.send("configuration").send("preferences")
	wailsSetKVCValue(prefs, "developerExtrasEnabled", nsNumberBool(true))
}

// kvcSetterSelectors returns the two setter selectors that KVC looks for when
// writing key: set<Key>: and _set<Key>:, with the first letter of the key
// capitalised. Setting the key itself still goes through setValue:forKey:
// with the original key.
func kvcSetterSelectors(key string) (string, string) {
	setter := "set" + strings.ToUpper(key[:1]) + key[1:] + ":"
	return setter, "_" + setter
}

// wailsSetKVCValue performs the KVC write that the cgo bridge wraps in
// @try/@catch. It skips the write unless the receiver implements a setter for
// the key, because an unsupported key raises an uncatchable NSException.
func wailsSetKVCValue(o id, key string, value id) {
	if o.isNil() {
		return
	}
	setter, privateSetter := kvcSetterSelectors(key)
	if !respondsTo(o, setter) && !respondsTo(o, privateSetter) {
		return
	}
	o.send("setValue:forKey:", value, nsString(key))
}
