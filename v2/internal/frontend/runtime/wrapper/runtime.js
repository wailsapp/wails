(()=>{var u=Object.defineProperty;var m=n=>u(n,"__esModule",{value:!0});var e=(n,i)=>{m(n);for(var o in i)u(n,o,{get:i[o],enumerable:!0})};var t={};e(t,{LogDebug:()=>W,LogError:()=>s,LogFatal:()=>x,LogInfo:()=>c,LogTrace:()=>p,LogWarning:()=>f});function p(n){window.runtime.LogTrace(n)}function W(n){window.runtime.LogDebug(n)}function c(n){window.runtime.LogInfo(n)}function f(n){window.runtime.LogWarning(n)}function s(n){window.runtime.LogError(n)}function x(n){window.runtime.LogFatal(n)}var w={};e(w,{EventsEmit:()=>g,EventsOn:()=>a,EventsOnMultiple:()=>l,EventsOnce:()=>S});function l(n,i,o){window.runtime.EventsOnMultiple(n,i,o)}function a(n,i){OnMultiple(n,i,-1)}function S(n,i){OnMultiple(n,i,1)}function g(n){let i=[n].slice.call(arguments);return window.runtime.EventsEmit.apply(null,i)}var r={};e(r,{WindowCenter:()=>M,WindowClose:()=>H,WindowFullscreen:()=>v,WindowGetPosition:()=>B,WindowGetSize:()=>O,WindowHide:()=>C,WindowMaximise:()=>T,WindowMinimise:()=>h,WindowReload:()=>L,WindowSetMaxSize:()=>F,WindowSetMinSize:()=>G,WindowSetPosition:()=>R,WindowSetRGBA:()=>D,WindowSetSize:()=>U,WindowSetTitle:()=>E,WindowShow:()=>P,WindowUnFullscreen:()=>z,WindowUnmaximise:()=>b,WindowUnminimise:()=>A});function L(){window.runtime.WindowReload()}function M(){window.runtime.WindowCenter()}function E(n){window.runtime.WindowSetTitle(n)}function v(){window.runtime.WindowFullscreen()}function z(){window.runtime.WindowUnFullscreen()}function O(){window.runtime.WindowGetSize()}function U(n,i){window.runtime.WindowSetSize(n,i)}function F(n,i){window.runtime.WindowSetMaxSize(n,i)}function G(n,i){window.runtime.WindowSetMinSize(n,i)}function R(n,i){window.runtime.WindowSetPosition(n,i)}function B(){window.runtime.WindowGetPosition()}function C(){window.runtime.WindowHide()}function P(){window.runtime.WindowShow()}function T(){window.runtime.WindowMaximise()}function b(){window.runtime.WindowUnmaximise()}function h(){window.runtime.WindowMinimise()}function A(){window.runtime.WindowUnminimise()}function D(n){window.runtime.WindowSetRGBA(n)}function H(){window.runtime.WindowClose()}var d={};e(d,{BrowserOpenURL:()=>I});function I(n){window.runtime.BrowserOpenURL(n)}var j={...t,...w,...r,...d};})();
