//go:build android

package runtime

// Android uses window.wails.invoke which is set up via addJavascriptInterface in WailsJSBridge.
//
// The critical fix: Android's shouldInterceptRequest does NOT expose POST bodies, so the default
// fetch-based transport loses the request body for /wails/runtime calls (bound methods, clipboard, etc).
// We patch window.fetch to route these calls through the synchronous window.wails.invoke bridge instead.
var invoke = `
window._wails.invoke=function(m){
    return window.wails.invoke(typeof m==='string'?m:JSON.stringify(m));
};
(function(){
    var _origFetch=window.fetch;
    window.fetch=function(input,init){
        var url=typeof input==='string'?input:(input instanceof URL?input.toString():input.url);
        if(url&&url.indexOf('/wails/runtime')!==-1&&init&&init.method==='POST'&&init.body){
            return new Promise(function(resolve,reject){
                try{
                    var resp=window.wails.invoke(typeof init.body==='string'?init.body:JSON.stringify(init.body));
                    resolve(new Response(resp||'null',{status:200,headers:{'Content-Type':'application/json'}}));
                }catch(e){reject(e);}
            });
        }
        return _origFetch.apply(window,arguments);
    };
})();
`
var flags = ""
