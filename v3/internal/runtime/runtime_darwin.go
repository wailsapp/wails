//go:build darwin

package runtime

var invoke = "window._wails.invoke=function(msg){window.webkit.messageHandlers.external.postMessage(msg);};"
var flags = ""
