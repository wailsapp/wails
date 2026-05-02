//go:build windows

package runtime

var invoke = `window._wails.invoke=window.chrome.webview.postMessage;`
