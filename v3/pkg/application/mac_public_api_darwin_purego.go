//go:build darwin && purego && !ios && !server && !private_mac_apis

package application

// Public-API replacements for every call in mac_private_api_darwin_purego.go,
// selected when private_mac_apis is absent. Where macOS offers no public
// equivalent the function is a documented no-op: the Go API is unchanged, only
// the visual result differs. This file mirrors mac_public_api_darwin.go.

// There is no public WKWebView transparency switch.
func wailsPrivateSetWebviewTransparent(webView id) {}

func wailsPrivateSetWebviewBackgroundColour(webView id, r, g, b, alpha int) {
	if webView.isNil() {
		return
	}
	colour := class("NSColor").send("colorWithRed:green:blue:alpha:",
		CGFloat(r)/255, CGFloat(g)/255, CGFloat(b)/255, CGFloat(alpha)/255)
	if macOSAtLeast(12, 0) && respondsTo(webView, "setUnderPageBackgroundColor:") {
		webView.send("setUnderPageBackgroundColor:", colour)
		return
	}
	webView.send("setWantsLayer:", true)
	layer := webView.send("layer")
	if !layer.isNil() {
		layer.send("setBackgroundColor:", colour.send("CGColor"))
	}
}

// NSGlassEffectView documents two styles: regular (0) and clear (1). Wails'
// light and dark styles have no native counterpart, so they are expressed with
// NSAppearance instead.
func wailsPrivateSetGlassStyle(glassView id, style MacLiquidGlassStyle) {
	if glassView.isNil() {
		return
	}
	if respondsTo(glassView, "setStyle:") {
		nativeStyle := 0
		if style == LiquidGlassStyleVibrant {
			nativeStyle = 1
		}
		glassView.send("setValue:forKey:", nsNumberInt(nativeStyle), nsString("style"))
	}
	appearance := id(0)
	switch style {
	case LiquidGlassStyleLight:
		appearance = class("NSAppearance").send("appearanceNamed:", nsString("NSAppearanceNameAqua"))
	case LiquidGlassStyleDark:
		appearance = class("NSAppearance").send("appearanceNamed:", nsString("NSAppearanceNameDarkAqua"))
	}
	glassView.send("setAppearance:", appearance)
}

// Cross-window glass grouping has no public AppKit equivalent. Each glass view
// stays independent.
func wailsPrivateSetGlassGrouping(glassView id, groupID string, groupSpacing float64) {}

// There is no public API that opens the Web Inspector programmatically.
func wailsPrivateOpenWebInspector(webView id) {}

// On macOS 13.3+ the public inspectable property enables the inspector.
func wailsPrivateEnableWebInspector(webView id) {
	if webView.isNil() || !macOSAtLeast(13, 3) || !respondsTo(webView, "setInspectable:") {
		return
	}
	webView.send("setInspectable:", true)
}
