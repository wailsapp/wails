(()=>{var H=Object.defineProperty;var g=(e,t)=>{for(var o in t)H(e,o,{get:t[o],enumerable:!0})};var E={};g(E,{LogDebug:()=>J,LogError:()=>$,LogFatal:()=>Y,LogInfo:()=>N,LogLevel:()=>Q,LogPrint:()=>j,LogTrace:()=>U,LogWarning:()=>X,SetLogLevel:()=>q});function u(e,t){window.WailsInvoke("L"+e+t)}function U(e){u("T",e)}function j(e){u("P",e)}function J(e){u("D",e)}function N(e){u("I",e)}function X(e){u("W",e)}function $(e){u("E",e)}function Y(e){u("F",e)}function q(e){u("S",e)}var Q={TRACE:1,DEBUG:2,INFO:3,WARNING:4,ERROR:5};var h=class{constructor(t,o,n){this.eventName=t,this.maxCallbacks=n||-1,this.Callback=i=>(o.apply(null,i),this.maxCallbacks===-1?!1:(this.maxCallbacks-=1,this.maxCallbacks===0))}},w={};function v(e,t,o){w[e]=w[e]||[];let n=new h(e,t,o);return w[e].push(n),()=>Z(n)}function W(e,t){return v(e,t,-1)}function T(e,t){return v(e,t,1)}function L(e){let t=e.name;if(w[t]){let o=w[t].slice();for(let n=w[t].length-1;n>=0;n-=1){let i=w[t][n],r=e.data;i.Callback(r)&&o.splice(n,1)}o.length===0?m(t):w[t]=o}}function P(e){let t;try{t=JSON.parse(e)}catch{let n="Invalid JSON passed to Notify: "+e;throw new Error(n)}L(t)}function O(e){let t={name:e,data:[].slice.apply(arguments).slice(1)};L(t),window.WailsInvoke("EE"+JSON.stringify(t))}function m(e){delete w[e],window.WailsInvoke("EX"+e)}function A(e,...t){m(e),t.length>0&&t.forEach(o=>{m(o)})}function Z(e){let t=e.eventName;w[t]=w[t].filter(o=>o!==e),w[t].length===0&&m(t)}var c={};function K(){var e=new Uint32Array(1);return window.crypto.getRandomValues(e)[0]}function _(){return Math.random()*9007199254740991}var x;window.crypto?x=K:x=_;function a(e,t,o){return o==null&&(o=0),new Promise(function(n,i){var r;do r=e+"-"+x();while(c[r]);var l;o>0&&(l=setTimeout(function(){i(Error("Call to "+e+" timed out. Request ID: "+r))},o)),c[r]={timeoutHandle:l,reject:i,resolve:n};try{let d={name:e,args:t,callbackID:r};window.WailsInvoke("C"+JSON.stringify(d))}catch(d){console.error(d)}})}window.ObfuscatedCall=(e,t,o)=>(o==null&&(o=0),new Promise(function(n,i){var r;do r=e+"-"+x();while(c[r]);var l;o>0&&(l=setTimeout(function(){i(Error("Call to method "+e+" timed out. Request ID: "+r))},o)),c[r]={timeoutHandle:l,reject:i,resolve:n};try{let d={id:e,args:t,callbackID:r};window.WailsInvoke("c"+JSON.stringify(d))}catch(d){console.error(d)}}));function M(e){let t;try{t=JSON.parse(e)}catch(i){let r=`Invalid JSON passed to callback: ${i.message}. Message: ${e}`;throw runtime.LogDebug(r),new Error(r)}let o=t.callbackid,n=c[o];if(!n){let i=`Callback '${o}' not registered!!!`;throw console.error(i),new Error(i)}clearTimeout(n.timeoutHandle),delete c[o],t.error?n.reject(t.error):n.resolve(t.result)}window.go={};function z(e){try{e=JSON.parse(e)}catch(t){console.error(t)}window.go=window.go||{},Object.keys(e).forEach(t=>{window.go[t]=window.go[t]||{},Object.keys(e[t]).forEach(o=>{window.go[t][o]=window.go[t][o]||{},Object.keys(e[t][o]).forEach(n=>{window.go[t][o][n]=function(){let i=0;function r(){let l=[].slice.call(arguments);return a([t,o,n].join("."),l,i)}return r.setTimeout=function(l){i=l},r.getTimeout=function(){return i},r}()})})})}var y={};g(y,{WindowCenter:()=>re,WindowFullscreen:()=>le,WindowGetPosition:()=>me,WindowGetSize:()=>fe,WindowHide:()=>ve,WindowIsFullscreen:()=>we,WindowIsMaximised:()=>he,WindowIsMinimised:()=>Se,WindowIsNormal:()=>ke,WindowMaximise:()=>xe,WindowMinimise:()=>ye,WindowReload:()=>ee,WindowReloadApp:()=>te,WindowSetAlwaysOnTop:()=>pe,WindowSetBackgroundColour:()=>Ie,WindowSetDarkTheme:()=>ie,WindowSetLightTheme:()=>ne,WindowSetMaxSize:()=>ue,WindowSetMinSize:()=>ce,WindowSetPosition:()=>ge,WindowSetSize:()=>de,WindowSetSystemDefaultTheme:()=>oe,WindowSetTitle:()=>se,WindowShow:()=>We,WindowToggleMaximise:()=>De,WindowUnfullscreen:()=>ae,WindowUnmaximise:()=>Ee,WindowUnminimise:()=>be});function ee(){window.location.reload()}function te(){window.WailsInvoke("WR")}function oe(){window.WailsInvoke("WASDT")}function ne(){window.WailsInvoke("WALT")}function ie(){window.WailsInvoke("WADT")}function re(){window.WailsInvoke("Wc")}function se(e){window.WailsInvoke("WT"+e)}function le(){window.WailsInvoke("WF")}function ae(){window.WailsInvoke("Wf")}function we(){return a(":wails:WindowIsFullscreen")}function de(e,t){window.WailsInvoke("Ws:"+e+":"+t)}function fe(){return a(":wails:WindowGetSize")}function ue(e,t){window.WailsInvoke("WZ:"+e+":"+t)}function ce(e,t){window.WailsInvoke("Wz:"+e+":"+t)}function pe(e){window.WailsInvoke("WATP:"+(e?"1":"0"))}function ge(e,t){window.WailsInvoke("Wp:"+e+":"+t)}function me(){return a(":wails:WindowGetPos")}function ve(){window.WailsInvoke("WH")}function We(){window.WailsInvoke("WS")}function xe(){window.WailsInvoke("WM")}function De(){window.WailsInvoke("Wt")}function Ee(){window.WailsInvoke("WU")}function he(){return a(":wails:WindowIsMaximised")}function ye(){window.WailsInvoke("Wm")}function be(){window.WailsInvoke("Wu")}function Se(){return a(":wails:WindowIsMinimised")}function ke(){return a(":wails:WindowIsNormal")}function Ie(e,t,o,n){let i=JSON.stringify({r:e||0,g:t||0,b:o||0,a:n||255});window.WailsInvoke("Wr:"+i)}var b={};g(b,{ScreenGetAll:()=>Ce});function Ce(){return a(":wails:ScreenGetAll")}var S={};g(S,{BrowserOpenURL:()=>Te});function Te(e){window.WailsInvoke("BO:"+e)}var k={};g(k,{ClipboardGetText:()=>Pe,ClipboardSetText:()=>Le});function Le(e){return a(":wails:ClipboardSetText",[e])}function Pe(){return a(":wails:ClipboardGetText")}var I={};g(I,{CanResolveFilePaths:()=>B,DragAndDropOff:()=>Me,DragAndDropOn:()=>Ae,ResolveFilePaths:()=>Oe});var s={registered:!1,defaultUseDropTarget:!0,useDropTarget:!0,prevElement:null};function R(e){if(e.preventDefault(),e.stopPropagation(),!s.useDropTarget)return;let t=document.elementFromPoint(e.x,e.y);if(t===s.prevElement)return;let o=t.style,n=null;Object.keys(o).findIndex(i=>o[i]===window.wails.flags.cssDropProperty)<0&&(t=t.closest(`[style*='${window.wails.flags.cssDropProperty}']`)),t!==null&&(n=window.getComputedStyle(t).getPropertyValue(window.wails.flags.cssDropProperty),n&&(n=n.trim()),n===window.wails.flags.cssDropValue?t.classList.add("wails-drop-target-active"):s.prevElement&&(t.classList.remove("wails-drop-target-active"),s.prevElement.classList.remove("wails-drop-target-active")),s.prevElement=t)}function V(e){if(e.preventDefault(),e.stopPropagation(),!s.useDropTarget)return;let t=document.elementFromPoint(e.x,e.y),o=window.getComputedStyle(t).getPropertyValue(window.wails.flags.cssDropProperty);o&&(o=o.trim()),o!==window.wails.flags.cssDropValue&&s.prevElement&&(t.classList.remove("wails-drop-target-active"),s.prevElement.classList.remove("wails-drop-target-active"))}function F(e){if(e.preventDefault(),e.stopPropagation(),!s.useDropTarget)return;let t=document.elementFromPoint(e.x,e.y),o=window.getComputedStyle(t).getPropertyValue(window.wails.flags.cssDropProperty);if(o&&(o=o.trim()),o!==window.wails.flags.cssDropValue){s.prevElement&&(t.classList.remove("wails-drop-target-active"),s.prevElement.classList.remove("wails-drop-target-active"));return}if(B()){let n=[];e.dataTransfer.items?n=[...e.dataTransfer.items].map((i,r)=>{if(i.kind==="file")return i.getAsFile()}):n=[...e.dataTransfer.files],window.runtime.ResolveFilePaths(e.x,e.y,n)}s.prevElement&&s.prevElement.classList.remove("wails-drop-target-active")}function B(){return window.chrome?.webview?.postMessageWithAdditionalObjects!=null}function Oe(e,t,o){if(window.chrome?.webview?.postMessageWithAdditionalObjects){chrome.webview.postMessageWithAdditionalObjects(`file:drop:${e}:${t}`,o);return}console.warn("unsupported platform")}function Ae(e,t){if(!window.wails.flags.enableWailsDragAndDrop||typeof e!="function"||s.registered)return;s.registered=!0;let o=typeof t;s.useDropTarget=o==="undefined"||o!=="boolean"?s.defaultUseDropTarget:t,window.addEventListener("dragover",R),window.addEventListener("dragleave",V),window.addEventListener("drop",F);let n=e;s.useDropTarget&&(n=function(i,r,l){let d=document.elementFromPoint(i,r);if(!d)return;let p=window.getComputedStyle(d).getPropertyValue(window.wails.flags.cssDropProperty);p&&(p=p.trim()),p===window.wails.flags.cssDropValue&&e(i,r,l)}),W("wails:dnd:drop",n)}function Me(){window.removeEventListener("dragover",R),window.removeEventListener("dragleave",V),window.removeEventListener("drop",F),EventsOff("wails:dnd:drop"),s.registered=!1}function G(e){let t=e.target;switch(window.getComputedStyle(t).getPropertyValue("--default-contextmenu").trim()){case"show":return;case"hide":e.preventDefault();return;default:if(t.isContentEditable)return;let i=window.getSelection(),r=i.toString().length>0;if(r)for(let l=0;l<i.rangeCount;l++){let p=i.getRangeAt(l).getClientRects();for(let D=0;D<p.length;D++){let C=p[D];if(document.elementFromPoint(C.left,C.top)===t)return}}if((t.tagName==="INPUT"||t.tagName==="TEXTAREA")&&(r||!t.readOnly&&!t.disabled))return;e.preventDefault()}}function Re(){window.WailsInvoke("Q")}function Ve(){window.WailsInvoke("S")}function Fe(){window.WailsInvoke("H")}function Be(){return a(":wails:Environment")}window.runtime={...E,...y,...S,...b,...k,...I,EventsOn:W,EventsOnce:T,EventsOnMultiple:v,EventsEmit:O,EventsOff:A,Environment:Be,Show:Ve,Hide:Fe,Quit:Re};window.wails={Callback:M,EventsNotify:P,SetBindings:z,eventListeners:w,callbacks:c,flags:{disableScrollbarDrag:!1,disableDefaultContextMenu:!1,enableResize:!1,defaultCursor:null,borderThickness:6,shouldDrag:!1,deferDragToMouseMove:!0,cssDragProperty:"--wails-draggable",cssDragValue:"drag",cssDropProperty:"--wails-drop-target",cssDropValue:"drop",enableWailsDragAndDrop:!1,wailsDropPreviousElement:null}};window.wailsbindings&&(window.wails.SetBindings(window.wailsbindings),delete window.wails.SetBindings);delete window.wailsbindings;var Ge=function(e){var t=window.getComputedStyle(e.target).getPropertyValue(window.wails.flags.cssDragProperty);return t&&(t=t.trim()),!(t!==window.wails.flags.cssDragValue||e.buttons!==1||e.detail!==1)};window.wails.setCSSDragProperties=function(e,t){window.wails.flags.cssDragProperty=e,window.wails.flags.cssDragValue=t};window.wails.setCSSDropProperties=function(e,t){window.wails.flags.cssDropProperty=e,window.wails.flags.cssDropValue=t};window.addEventListener("mousedown",e=>{if(window.wails.flags.resizeEdge){window.WailsInvoke("resize:"+window.wails.flags.resizeEdge),e.preventDefault();return}if(Ge(e)){if(window.wails.flags.disableScrollbarDrag&&(e.offsetX>e.target.clientWidth||e.offsetY>e.target.clientHeight))return;window.wails.flags.deferDragToMouseMove?window.wails.flags.shouldDrag=!0:(e.preventDefault(),window.WailsInvoke("drag"));return}else window.wails.flags.shouldDrag=!1});window.addEventListener("mouseup",()=>{window.wails.flags.shouldDrag=!1});function f(e){document.documentElement.style.cursor=e||window.wails.flags.defaultCursor,window.wails.flags.resizeEdge=e}window.addEventListener("mousemove",function(e){if(window.wails.flags.shouldDrag&&(window.wails.flags.shouldDrag=!1,(e.buttons!==void 0?e.buttons:e.which)>0)){window.WailsInvoke("drag");return}if(!window.wails.flags.enableResize)return;window.wails.flags.defaultCursor==null&&(window.wails.flags.defaultCursor=document.documentElement.style.cursor),window.outerWidth-e.clientX<window.wails.flags.borderThickness&&window.outerHeight-e.clientY<window.wails.flags.borderThickness&&(document.documentElement.style.cursor="se-resize");let t=window.outerWidth-e.clientX<window.wails.flags.borderThickness,o=e.clientX<window.wails.flags.borderThickness,n=e.clientY<window.wails.flags.borderThickness,i=window.outerHeight-e.clientY<window.wails.flags.borderThickness;!o&&!t&&!n&&!i&&window.wails.flags.resizeEdge!==void 0?f():t&&i?f("se-resize"):o&&i?f("sw-resize"):o&&n?f("nw-resize"):n&&t?f("ne-resize"):o?f("w-resize"):n?f("n-resize"):i?f("s-resize"):t&&f("e-resize")});window.addEventListener("contextmenu",function(e){window.wails.flags.disableDefaultContextMenu?e.preventDefault():G(e)});window.WailsInvoke("runtime:ready");})();
