(()=>{var z=Object.defineProperty;var c=(e,n)=>{for(var o in n)z(e,o,{get:n[o],enumerable:!0})};var W={};c(W,{LogDebug:()=>A,LogError:()=>U,LogFatal:()=>G,LogInfo:()=>B,LogLevel:()=>P,LogPrint:()=>R,LogTrace:()=>C,LogWarning:()=>H,SetLogLevel:()=>J});function a(e,n){window.WailsInvoke("L"+e+n)}function C(e){a("T",e)}function R(e){a("P",e)}function A(e){a("D",e)}function B(e){a("I",e)}function H(e){a("W",e)}function U(e){a("E",e)}function G(e){a("F",e)}function J(e){a("S",e)}var P={TRACE:1,DEBUG:2,INFO:3,WARNING:4,ERROR:5};var p=class{constructor(n,o){o=o||-1,this.Callback=i=>(n.apply(null,i),o===-1?!1:(o-=1,o===0))}},s={};function u(e,n,o){s[e]=s[e]||[];let i=new p(n,o);s[e].push(i)}function E(e,n){u(e,n,-1)}function h(e,n){u(e,n,1)}function I(e){let n=e.name;if(s[n]){let o=s[n].slice();for(let i=0;i<s[n].length;i+=1){let t=s[n][i],r=e.data;t.Callback(r)&&o.splice(i,1)}s[n]=o}}function b(e){let n;try{n=JSON.parse(e)}catch{let i="Invalid JSON passed to Notify: "+e;throw new Error(i)}I(n)}function y(e){let n={name:e,data:[].slice.apply(arguments).slice(1)};I(n),window.WailsInvoke("EE"+JSON.stringify(n))}function T(e){delete s[e],window.WailsInvoke("EX"+e)}var f={};function M(){var e=new Uint32Array(1);return window.crypto.getRandomValues(e)[0]}function F(){return Math.random()*9007199254740991}var v;window.crypto?v=M:v=F;function l(e,n,o){return o==null&&(o=0),new Promise(function(i,t){var r;do r=e+"-"+v();while(f[r]);var d;o>0&&(d=setTimeout(function(){t(Error("Call to "+e+" timed out. Request ID: "+r))},o)),f[r]={timeoutHandle:d,reject:t,resolve:i};try{let g={name:e,args:n,callbackID:r};window.WailsInvoke("C"+JSON.stringify(g))}catch(g){console.error(g)}})}function L(e){let n;try{n=JSON.parse(e)}catch(t){let r=`Invalid JSON passed to callback: ${t.message}. Message: ${e}`;throw runtime.LogDebug(r),new Error(r)}let o=n.callbackid,i=f[o];if(!i){let t=`Callback '${o}' not registered!!!`;throw console.error(t),new Error(t)}clearTimeout(i.timeoutHandle),delete f[o],n.error?i.reject(n.error):i.resolve(n.result)}window.go={};function D(e){try{e=JSON.parse(e)}catch(n){console.error(n)}window.go=window.go||{},Object.keys(e).forEach(n=>{window.go[n]=window.go[n]||{},Object.keys(e[n]).forEach(o=>{window.go[n][o]=window.go[n][o]||{},Object.keys(e[n][o]).forEach(i=>{window.go[n][o][i]=function(){let t=0;function r(){let d=[].slice.call(arguments);return l([n,o,i].join("."),d,t)}return r.setTimeout=function(d){t=d},r.getTimeout=function(){return t},r}()})})})}var x={};c(x,{WindowCenter:()=>Q,WindowFullscreen:()=>N,WindowGetPosition:()=>te,WindowGetSize:()=>_,WindowHide:()=>re,WindowMaximise:()=>le,WindowMinimise:()=>de,WindowReload:()=>j,WindowReloadApp:()=>X,WindowSetAlwaysOnTop:()=>oe,WindowSetBackgroundColour:()=>ce,WindowSetDarkTheme:()=>$,WindowSetLightTheme:()=>V,WindowSetMaxSize:()=>ee,WindowSetMinSize:()=>ne,WindowSetPosition:()=>ie,WindowSetSize:()=>K,WindowSetSystemDefaultTheme:()=>Y,WindowSetTitle:()=>q,WindowShow:()=>se,WindowToggleMaximise:()=>we,WindowUnfullscreen:()=>Z,WindowUnmaximise:()=>ae,WindowUnminimise:()=>fe});function j(){window.location.reload()}function X(){window.WailsInvoke("WR")}function Y(){window.WailsInvoke("WASDT")}function V(){window.WailsInvoke("WALT")}function $(){window.WailsInvoke("WADT")}function Q(){window.WailsInvoke("Wc")}function q(e){window.WailsInvoke("WT"+e)}function N(){window.WailsInvoke("WF")}function Z(){window.WailsInvoke("Wf")}function K(e,n){window.WailsInvoke("Ws:"+e+":"+n)}function _(){return l(":wails:WindowGetSize")}function ee(e,n){window.WailsInvoke("WZ:"+e+":"+n)}function ne(e,n){window.WailsInvoke("Wz:"+e+":"+n)}function oe(e){window.WailsInvoke("WATP:"+(e?"1":"0"))}function ie(e,n){window.WailsInvoke("Wp:"+e+":"+n)}function te(){return l(":wails:WindowGetPos")}function re(){window.WailsInvoke("WH")}function se(){window.WailsInvoke("WS")}function le(){window.WailsInvoke("WM")}function we(){window.WailsInvoke("Wt")}function ae(){window.WailsInvoke("WU")}function de(){window.WailsInvoke("Wm")}function fe(){window.WailsInvoke("Wu")}function ce(e,n,o,i){let t=JSON.stringify({r:e||0,g:n||0,b:o||0,a:i||255});window.WailsInvoke("Wr:"+t)}var m={};c(m,{ScreenGetAll:()=>ue});function ue(){return l(":wails:ScreenGetAll")}var k={};c(k,{BrowserOpenURL:()=>ge});function ge(e){window.WailsInvoke("BO:"+e)}function We(){window.WailsInvoke("Q")}function pe(){window.WailsInvoke("S")}function ve(){window.WailsInvoke("H")}function xe(){return l(":wails:Environment")}window.runtime={...W,...x,...k,...m,EventsOn:E,EventsOnce:h,EventsOnMultiple:u,EventsEmit:y,EventsOff:T,Environment:xe,Show:pe,Hide:ve,Quit:We};window.wails={Callback:L,EventsNotify:b,SetBindings:D,eventListeners:s,callbacks:f,flags:{disableScrollbarDrag:!1,disableWailsDefaultContextMenu:!1,enableResize:!1,defaultCursor:null,borderThickness:6,shouldDrag:!1}};window.wails.SetBindings(window.wailsbindings);delete window.wails.SetBindings;window.addEventListener("mouseup",()=>{window.wails.flags.shouldDrag=!1});var me=function(e){return window.getComputedStyle(e.target).getPropertyValue("--wails-draggable")==="drag"},O=function(e){let n=e.target;for(;n!=null&&!n.hasAttribute("data-wails-no-drag");){if(n.hasAttribute("data-wails-drag"))return!0;n=n.parentElement}return!1},S=O;window.wails.useCSSDrag=function(e){e===!1?(console.log("Using original drag detection"),S=O):(console.log("Using CSS drag detection"),S=me)};window.addEventListener("mousedown",e=>{if(window.wails.flags.resizeEdge){window.WailsInvoke("resize:"+window.wails.flags.resizeEdge),e.preventDefault();return}if(S(e)){if(window.wails.flags.disableScrollbarDrag&&(e.offsetX>e.target.clientWidth||e.offsetY>e.target.clientHeight))return;window.wails.flags.shouldDrag=!0}});function w(e){document.body.style.cursor=e||window.wails.flags.defaultCursor,window.wails.flags.resizeEdge=e}window.addEventListener("mousemove",function(e){if(window.wails.flags.shouldDrag){window.WailsInvoke("drag");return}if(!window.wails.flags.enableResize)return;window.wails.flags.defaultCursor==null&&(window.wails.flags.defaultCursor=document.body.style.cursor),window.outerWidth-e.clientX<window.wails.flags.borderThickness&&window.outerHeight-e.clientY<window.wails.flags.borderThickness&&(document.body.style.cursor="se-resize");let n=window.outerWidth-e.clientX<window.wails.flags.borderThickness,o=e.clientX<window.wails.flags.borderThickness,i=e.clientY<window.wails.flags.borderThickness,t=window.outerHeight-e.clientY<window.wails.flags.borderThickness;!o&&!n&&!i&&!t&&window.wails.flags.resizeEdge!==void 0?w():n&&t?w("se-resize"):o&&t?w("sw-resize"):o&&i?w("nw-resize"):i&&n?w("ne-resize"):o?w("w-resize"):i?w("n-resize"):t?w("s-resize"):n&&w("e-resize")});window.addEventListener("contextmenu",function(e){window.wails.flags.disableWailsDefaultContextMenu&&e.preventDefault()});window.WailsInvoke("runtime:ready");})();
