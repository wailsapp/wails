(()=>{var G=Object.defineProperty;var c=(e,n)=>{for(var o in n)G(e,o,{get:n[o],enumerable:!0})};var x={};c(x,{LogDebug:()=>j,LogError:()=>Y,LogFatal:()=>$,LogInfo:()=>V,LogLevel:()=>Q,LogPrint:()=>U,LogTrace:()=>H,LogWarning:()=>X,SetLogLevel:()=>q});function f(e,n){window.WailsInvoke("L"+e+n)}function H(e){f("T",e)}function U(e){f("P",e)}function j(e){f("D",e)}function V(e){f("I",e)}function X(e){f("W",e)}function Y(e){f("E",e)}function $(e){f("F",e)}function q(e){f("S",e)}var Q={TRACE:1,DEBUG:2,INFO:3,WARNING:4,ERROR:5};var v=class{constructor(n,o,i){this.eventName=n,this.maxCallbacks=i||-1,this.Callback=t=>(o.apply(null,t),this.maxCallbacks===-1?!1:(this.maxCallbacks-=1,this.maxCallbacks===0))}},a={};function p(e,n,o){a[e]=a[e]||[];let i=new v(e,n,o);return a[e].push(i),()=>Z(i)}function y(e,n){return p(e,n,-1)}function E(e,n){return p(e,n,1)}function O(e){let n=e.name;if(a[n]){let o=a[n].slice();for(let i=a[n].length-1;i>=0;i-=1){let t=a[n][i],r=e.data;t.Callback(r)&&o.splice(i,1)}o.length===0?g(n):a[n]=o}}function C(e){let n;try{n=JSON.parse(e)}catch{let i="Invalid JSON passed to Notify: "+e;throw new Error(i)}O(n)}function T(e){let n={name:e,data:[].slice.apply(arguments).slice(1)};O(n),window.WailsInvoke("EE"+JSON.stringify(n))}function g(e){delete a[e],window.WailsInvoke("EX"+e)}function M(e,...n){g(e),n.length>0&&n.forEach(o=>{g(o)})}function Z(e){let n=e.eventName;a[n]=a[n].filter(o=>o!==e),a[n].length===0&&g(n)}var u={};function K(){var e=new Uint32Array(1);return window.crypto.getRandomValues(e)[0]}function _(){return Math.random()*9007199254740991}var W;window.crypto?W=K:W=_;function s(e,n,o){return o==null&&(o=0),new Promise(function(i,t){var r;do r=e+"-"+W();while(u[r]);var l;o>0&&(l=setTimeout(function(){t(Error("Call to "+e+" timed out. Request ID: "+r))},o)),u[r]={timeoutHandle:l,reject:t,resolve:i};try{let d={name:e,args:n,callbackID:r};window.WailsInvoke("C"+JSON.stringify(d))}catch(d){console.error(d)}})}window.ObfuscatedCall=(e,n,o)=>(o==null&&(o=0),new Promise(function(i,t){var r;do r=e+"-"+W();while(u[r]);var l;o>0&&(l=setTimeout(function(){t(Error("Call to method "+e+" timed out. Request ID: "+r))},o)),u[r]={timeoutHandle:l,reject:t,resolve:i};try{let d={id:e,args:n,callbackID:r};window.WailsInvoke("c"+JSON.stringify(d))}catch(d){console.error(d)}}));function L(e){let n;try{n=JSON.parse(e)}catch(t){let r=`Invalid JSON passed to callback: ${t.message}. Message: ${e}`;throw runtime.LogDebug(r),new Error(r)}let o=n.callbackid,i=u[o];if(!i){let t=`Callback '${o}' not registered!!!`;throw console.error(t),new Error(t)}clearTimeout(i.timeoutHandle),delete u[o],n.error?i.reject(n.error):i.resolve(n.result)}window.go={};function z(e){try{e=JSON.parse(e)}catch(n){console.error(n)}window.go=window.go||{},Object.keys(e).forEach(n=>{window.go[n]=window.go[n]||{},Object.keys(e[n]).forEach(o=>{window.go[n][o]=window.go[n][o]||{},Object.keys(e[n][o]).forEach(i=>{window.go[n][o][i]=function(){let t=0;function r(){let l=[].slice.call(arguments);return s([n,o,i].join("."),l,t)}return r.setTimeout=function(l){t=l},r.getTimeout=function(){return t},r}()})})})}var D={};c(D,{WindowCenter:()=>re,WindowFullscreen:()=>le,WindowGetPosition:()=>We,WindowGetSize:()=>fe,WindowHide:()=>me,WindowIsFullscreen:()=>we,WindowIsMaximised:()=>he,WindowIsMinimised:()=>be,WindowIsNormal:()=>ye,WindowMaximise:()=>ve,WindowMinimise:()=>Se,WindowReload:()=>ee,WindowReloadApp:()=>ne,WindowSetAlwaysOnTop:()=>ge,WindowSetBackgroundColour:()=>Ee,WindowSetDarkTheme:()=>te,WindowSetLightTheme:()=>ie,WindowSetMaxSize:()=>ue,WindowSetMinSize:()=>ce,WindowSetPosition:()=>pe,WindowSetSize:()=>de,WindowSetSystemDefaultTheme:()=>oe,WindowSetTitle:()=>se,WindowShow:()=>xe,WindowToggleMaximise:()=>De,WindowUnfullscreen:()=>ae,WindowUnmaximise:()=>ke,WindowUnminimise:()=>Ie});function ee(){window.location.reload()}function ne(){window.WailsInvoke("WR")}function oe(){window.WailsInvoke("WASDT")}function ie(){window.WailsInvoke("WALT")}function te(){window.WailsInvoke("WADT")}function re(){window.WailsInvoke("Wc")}function se(e){window.WailsInvoke("WT"+e)}function le(){window.WailsInvoke("WF")}function ae(){window.WailsInvoke("Wf")}function we(){return s(":wails:WindowIsFullscreen")}function de(e,n){window.WailsInvoke("Ws:"+e+":"+n)}function fe(){return s(":wails:WindowGetSize")}function ue(e,n){window.WailsInvoke("WZ:"+e+":"+n)}function ce(e,n){window.WailsInvoke("Wz:"+e+":"+n)}function ge(e){window.WailsInvoke("WATP:"+(e?"1":"0"))}function pe(e,n){window.WailsInvoke("Wp:"+e+":"+n)}function We(){return s(":wails:WindowGetPos")}function me(){window.WailsInvoke("WH")}function xe(){window.WailsInvoke("WS")}function ve(){window.WailsInvoke("WM")}function De(){window.WailsInvoke("Wt")}function ke(){window.WailsInvoke("WU")}function he(){return s(":wails:WindowIsMaximised")}function Se(){window.WailsInvoke("Wm")}function Ie(){window.WailsInvoke("Wu")}function be(){return s(":wails:WindowIsMinimised")}function ye(){return s(":wails:WindowIsNormal")}function Ee(e,n,o,i){let t=JSON.stringify({r:e||0,g:n||0,b:o||0,a:i||255});window.WailsInvoke("Wr:"+t)}var k={};c(k,{ScreenGetAll:()=>Oe});function Oe(){return s(":wails:ScreenGetAll")}var h={};c(h,{BrowserOpenURL:()=>Ce});function Ce(e){window.WailsInvoke("BO:"+e)}var S={};c(S,{ClipboardGetText:()=>Me,ClipboardSetText:()=>Te});function Te(e){return s(":wails:ClipboardSetText",[e])}function Me(){return s(":wails:ClipboardGetText")}function R(e){let n=e.target;switch(window.getComputedStyle(n).getPropertyValue("--default-contextmenu").trim()){case"show":return;case"hide":e.preventDefault();return;default:if(n.isContentEditable)return;let t=window.getSelection(),r=t.toString().length>0;if(r)for(let l=0;l<t.rangeCount;l++){let I=t.getRangeAt(l).getClientRects();for(let m=0;m<I.length;m++){let b=I[m];if(document.elementFromPoint(b.left,b.top)===n)return}}if((n.tagName==="INPUT"||n.tagName==="TEXTAREA")&&(r||!n.readOnly&&!n.disabled))return;e.preventDefault()}}function F(e){return window.WailsInvoke("DOD:"+JSON.stringify(e))}function P(e){return window.WailsInvoke("DOMD:"+JSON.stringify(e))}function A(e){return window.WailsInvoke("DOF:"+JSON.stringify(e))}function B(e){return window.WailsInvoke("DOMF:"+JSON.stringify(e))}function J(e){return window.WailsInvoke("DSF:"+JSON.stringify(e))}function N(e){return window.WailsInvoke("DM:"+JSON.stringify(e))}function ze(){window.WailsInvoke("Q")}function Re(){window.WailsInvoke("S")}function Fe(){window.WailsInvoke("H")}function Pe(){return s(":wails:Environment")}window.runtime={...x,...D,...h,...k,...S,EventsOn:y,EventsOnce:E,EventsOnMultiple:p,EventsEmit:T,EventsOff:M,Environment:Pe,Show:Re,Hide:Fe,Quit:ze,OpenDirectoryDialog:F,OpenMultipleDirectoriesDialog:P,OpenFileDialog:A,OpenMultipleFilesDialog:B,SaveFileDialog:J,MessageDialog:N};window.wails={Callback:L,EventsNotify:C,SetBindings:z,eventListeners:a,callbacks:u,flags:{disableScrollbarDrag:!1,disableDefaultContextMenu:!1,enableResize:!1,defaultCursor:null,borderThickness:6,shouldDrag:!1,deferDragToMouseMove:!0,cssDragProperty:"--wails-draggable",cssDragValue:"drag"}};window.wailsbindings&&(window.wails.SetBindings(window.wailsbindings),delete window.wails.SetBindings);delete window.wailsbindings;var Ae=function(e){var n=window.getComputedStyle(e.target).getPropertyValue(window.wails.flags.cssDragProperty);return n&&(n=n.trim()),!(n!==window.wails.flags.cssDragValue||e.buttons!==1||e.detail!==1)};window.wails.setCSSDragProperties=function(e,n){window.wails.flags.cssDragProperty=e,window.wails.flags.cssDragValue=n};window.addEventListener("mousedown",e=>{if(window.wails.flags.resizeEdge){window.WailsInvoke("resize:"+window.wails.flags.resizeEdge),e.preventDefault();return}if(Ae(e)){if(window.wails.flags.disableScrollbarDrag&&(e.offsetX>e.target.clientWidth||e.offsetY>e.target.clientHeight))return;window.wails.flags.deferDragToMouseMove?window.wails.flags.shouldDrag=!0:(e.preventDefault(),window.WailsInvoke("drag"));return}else window.wails.flags.shouldDrag=!1});window.addEventListener("mouseup",()=>{window.wails.flags.shouldDrag=!1});function w(e){document.documentElement.style.cursor=e||window.wails.flags.defaultCursor,window.wails.flags.resizeEdge=e}window.addEventListener("mousemove",function(e){if(window.wails.flags.shouldDrag&&(window.wails.flags.shouldDrag=!1,(e.buttons!==void 0?e.buttons:e.which)>0)){window.WailsInvoke("drag");return}if(!window.wails.flags.enableResize)return;window.wails.flags.defaultCursor==null&&(window.wails.flags.defaultCursor=document.documentElement.style.cursor),window.outerWidth-e.clientX<window.wails.flags.borderThickness&&window.outerHeight-e.clientY<window.wails.flags.borderThickness&&(document.documentElement.style.cursor="se-resize");let n=window.outerWidth-e.clientX<window.wails.flags.borderThickness,o=e.clientX<window.wails.flags.borderThickness,i=e.clientY<window.wails.flags.borderThickness,t=window.outerHeight-e.clientY<window.wails.flags.borderThickness;!o&&!n&&!i&&!t&&window.wails.flags.resizeEdge!==void 0?w():n&&t?w("se-resize"):o&&t?w("sw-resize"):o&&i?w("nw-resize"):i&&n?w("ne-resize"):o?w("w-resize"):i?w("n-resize"):t?w("s-resize"):n&&w("e-resize")});window.addEventListener("contextmenu",function(e){window.wails.flags.disableDefaultContextMenu?e.preventDefault():R(e)});window.WailsInvoke("runtime:ready");})();
