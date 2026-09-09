//go:build linux && !android

package runtime

// On webkit2gtk, `messageHandlers.external.postMessage` only works when
// `this` is bound to the handler object. Assigning the bare function
// reference (as we did historically) silently swallows messages when
// called as `window._wails.invoke(msg)` — the page's invoke loses the
// receiver. Wrap it like darwin does so callers can invoke without
// thinking about receiver binding.
var invoke = "window._wails.invoke=function(msg){window.webkit.messageHandlers.external.postMessage(msg);};"
