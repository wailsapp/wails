package application

import (
	_ "embed"
	"strings"
)

// inlineEventShimJS is a small ES5 script that installs
// `window.wails.Events.On / Emit` and `window._wails.dispatchWailsEvent`
// for windows whose HTML is loaded directly via WebviewWindowOptions.HTML
// rather than served by the asset server. Those pages can't import
// `/wails/runtime.js` (their origin is the literal string "null"), so
// the framework injects this fallback at construction time when the
// owning code asked for it.
//
// Injection is gated on WebviewWindowOptions.AllowSimpleEventEmit so the
// shim only ships into windows that have opted into the simple postMessage
// emit path in the first place — same security boundary that gates the
// host-side handler. See the field's GoDoc for the threat model.
//
//go:embed inline_event_shim.js
var inlineEventShimJS string

// maybeInjectInlineEventShim prepends the inline-event shim to the
// supplied HTML when the window has opted in via AllowSimpleEventEmit.
// Returns the (possibly modified) HTML unchanged otherwise.
//
// The injected `<script>` is placed at the very top of the document so
// that any inline event handlers registered later in the page can rely
// on `window.wails.Events` being present.
func maybeInjectInlineEventShim(html string, allow bool) string {
	if !allow || html == "" {
		return html
	}
	const marker = "data-wails-inline-event-shim"
	if strings.Contains(html, marker) {
		// Already injected (e.g. caller wrapped manually). Don't duplicate.
		return html
	}
	script := `<script ` + marker + `>` + "\n" + inlineEventShimJS + "\n</script>\n"
	// If the body starts with <!doctype …>, keep the doctype on line 1 and
	// drop the script immediately after; otherwise just prepend.
	lower := strings.ToLower(html)
	if idx := strings.Index(lower, "<!doctype"); idx >= 0 {
		closing := strings.Index(html[idx:], ">")
		if closing >= 0 {
			cut := idx + closing + 1
			return html[:cut] + "\n" + script + html[cut:]
		}
	}
	return script + html
}
