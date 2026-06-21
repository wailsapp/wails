//go:build android

package runtime

// Android uses window.wails.invoke which is set up via addJavascriptInterface in WailsJSBridge
// We need to log the state to debug why it's not being detected
var invoke = `
console.log('[Wails Android Runtime] Injecting runtime, window.wails exists:', !!window.wails);
console.log('[Wails Android Runtime] window.wails.invoke exists:', !!(window.wails && window.wails.invoke));
window._wails.invoke=function(m){
    console.log('[Wails Android Runtime] _wails.invoke called:', m);
    return window.wails.invoke(typeof m==='string'?m:JSON.stringify(m));
};
console.log('[Wails Android Runtime] Runtime injection complete');
`
var flags = ""
