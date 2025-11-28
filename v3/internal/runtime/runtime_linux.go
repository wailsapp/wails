//go:build linux && !android

package runtime

var invoke = "window._wails.invoke=window.webkit.messageHandlers.external.postMessage;"
var flags = ""
