// \0vite/modulepreload-polyfill.js
(function polyfill() {
	const relList = document.createElement("link").relList;
	if (relList && relList.supports && relList.supports("modulepreload")) return;
	for (const link of document.querySelectorAll("link[rel=\"modulepreload\"]")) processPreload(link);
	new MutationObserver((mutations) => {
		for (const mutation of mutations) {
			if (mutation.type !== "childList") continue;
			for (const node of mutation.addedNodes) if (node.tagName === "LINK" && node.rel === "modulepreload") processPreload(node);
		}
	}).observe(document, {
		childList: true,
		subtree: true
	});
	function getFetchOpts(link) {
		const fetchOpts = {};
		if (link.integrity) fetchOpts.integrity = link.integrity;
		if (link.referrerPolicy) fetchOpts.referrerPolicy = link.referrerPolicy;
		if (link.crossOrigin === "use-credentials") fetchOpts.credentials = "include";
		else if (link.crossOrigin === "anonymous") fetchOpts.credentials = "omit";
		else fetchOpts.credentials = "same-origin";
		return fetchOpts;
	}
	function processPreload(link) {
		if (link.ep) return;
		link.ep = true;
		const fetchOpts = getFetchOpts(link);
		fetch(link.href, fetchOpts);
	}
})();
// ../../../internal/runtime/desktop/@wailsio/runtime/src/environment.ts
/**
* True when running inside a browser/webview with a DOM available.
* False under server-side rendering (e.g. `next build` prerendering),
* where application code may import the runtime module even though no
* Wails APIs can actually be used (#4679). Modules must not touch
* `window`/`document` at import time except behind this guard.
*/
var hasDOM = typeof window !== "undefined" && typeof document !== "undefined";
// ../../../internal/runtime/desktop/@wailsio/runtime/src/nanoid.ts
var urlAlphabet = "useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict";
function nanoid(size = 21) {
	const source = typeof crypto !== "undefined" && typeof crypto.getRandomValues === "function" ? crypto : null;
	if (source) {
		const bytes = new Uint8Array(size);
		source.getRandomValues(bytes);
		let id = "";
		for (let i = 0; i < size; i++) id += urlAlphabet[bytes[i] & 63];
		return id;
	}
	let id = "";
	let i = size | 0;
	while (i--) id += urlAlphabet[Math.random() * 64 | 0];
	return id;
}
// ../../../internal/runtime/desktop/@wailsio/runtime/src/runtime.ts
function runtimeURL() {
	return window.location.origin + "/wails/runtime";
}
var CHUNK_THRESHOLD = 524288;
/**
* Exception class that will be thrown in case the bound method returns an error.
* The value of the {@link RuntimeError#name} property is "RuntimeError".
*/
var RuntimeError = class extends Error {
	/**
	* Constructs a new RuntimeError instance.
	* @param message - The error message.
	* @param options - Options to be forwarded to the Error constructor.
	*/
	constructor(message, options) {
		super(message, options);
		this.name = "RuntimeError";
	}
};
var objectNames = Object.freeze({
	Call: 0,
	Clipboard: 1,
	Application: 2,
	Events: 3,
	ContextMenu: 4,
	Dialog: 5,
	Window: 6,
	Screens: 7,
	System: 8,
	Browser: 9,
	CancelCall: 10,
	IOS: 11,
	Android: 12
});
var clientId = nanoid();
/**
* Custom transport implementation (can be set by user)
*/
var customTransport = null;
/**
* Creates a new runtime caller with specified ID.
*
* @param object - The object to invoke the method on.
* @param windowName - The name of the window.
* @return The new runtime caller function.
*/
function newRuntimeCaller(object, windowName = "") {
	return function(method, args = null) {
		return runtimeCallWithID(object, method, windowName, args);
	};
}
async function runtimeCallWithID(objectID, method, windowName, args) {
	if (customTransport) return customTransport.call(objectID, method, windowName, args);
	let url = new URL(runtimeURL());
	let body = {
		object: objectID,
		method
	};
	if (args !== null && args !== void 0) body.args = args;
	let headers = {
		["x-wails-client-id"]: clientId,
		["Content-Type"]: "application/json"
	};
	if (windowName) headers["x-wails-window-name"] = windowName;
	const bodyStr = JSON.stringify(body);
	let response;
	if (bodyStr.length > CHUNK_THRESHOLD) response = await sendChunked(url, headers, bodyStr);
	else response = await fetch(url, {
		method: "POST",
		headers,
		body: bodyStr
	});
	if (!response.ok) {
		if (response.headers.get("Content-Type")?.includes("application/json")) {
			const json = await response.json();
			let err;
			switch (json.kind) {
				case "ReferenceError":
					err = new ReferenceError(json.message);
					break;
				case "TypeError":
					err = new TypeError(json.message);
					break;
				case "RuntimeError":
					err = new RuntimeError(json.message);
					break;
				default: err = new Error(json.message);
			}
			err.cause = json.cause;
			throw err;
		}
		throw new Error(await response.text());
	}
	if ((response.headers.get("Content-Type")?.indexOf("application/json") ?? -1) !== -1) return response.json();
	else return response.text();
}
async function sendChunked(url, headers, bodyStr) {
	const chunkId = nanoid();
	const bodyBytes = new TextEncoder().encode(bodyStr);
	const totalChunks = Math.ceil(bodyBytes.length / CHUNK_THRESHOLD);
	for (let i = 0; i < totalChunks - 1; i++) {
		const chunk = bodyBytes.subarray(i * CHUNK_THRESHOLD, (i + 1) * CHUNK_THRESHOLD);
		const resp = await fetch(url, {
			method: "POST",
			headers: {
				...headers,
				"x-wails-chunk-id": chunkId,
				"x-wails-chunk-index": String(i),
				"x-wails-chunk-total": String(totalChunks)
			},
			body: chunk
		});
		if (!resp.ok) throw new Error(await resp.text());
	}
	return fetch(url, {
		method: "POST",
		headers: {
			...headers,
			"x-wails-chunk-id": chunkId,
			"x-wails-chunk-index": String(totalChunks - 1),
			"x-wails-chunk-total": String(totalChunks)
		},
		body: bodyBytes.subarray((totalChunks - 1) * CHUNK_THRESHOLD)
	});
}
var androidBridge = hasDOM && typeof window.wails?.invokeAsync === "function" ? window.wails : null;
if (androidBridge) {
	const pending = /* @__PURE__ */ new Map();
	window._wailsAndroidCallback = (id, response, error) => {
		const promise = pending.get(id);
		if (!promise) return;
		pending.delete(id);
		if (error) {
			promise.reject(new Error(error));
			return;
		}
		try {
			const envelope = JSON.parse(response ?? "{}");
			if (!envelope.ok) {
				promise.reject(new Error(envelope.error ?? "unknown runtime call error"));
				return;
			}
			promise.resolve("text" in envelope ? envelope.text : envelope.data);
		} catch (e) {
			promise.reject(e);
		}
	};
	customTransport = { call(objectID, method, windowName, args) {
		return new Promise((resolve, reject) => {
			const id = nanoid();
			pending.set(id, {
				resolve,
				reject
			});
			try {
				androidBridge.invokeAsync(id, JSON.stringify({
					object: objectID,
					method,
					windowName,
					args: args ?? null,
					clientId
				}));
			} catch (e) {
				pending.delete(id);
				reject(e);
			}
		});
	} };
}
objectNames.System;
var _invoke = (function() {
	try {
		if (window.chrome?.webview?.postMessage) return window.chrome.webview.postMessage.bind(window.chrome.webview);
		else if (window.webkit?.messageHandlers?.["external"]?.postMessage) return window.webkit.messageHandlers["external"].postMessage.bind(window.webkit.messageHandlers["external"]);
		else if (window.wails?.invoke) {
			const bridge = window.wails;
			return (msg) => bridge.invoke(typeof msg === "string" ? msg : JSON.stringify(msg));
		}
	} catch (e) {}
	console.warn("\n%c⚠️ Browser Environment Detected %c\n\n%cOnly UI previews are available in the browser. For full functionality, please run the application in desktop mode.\nMore information at: https://v3.wails.io/learn/build/#using-a-browser-for-development\n", "background: #ffffff; color: #000000; font-weight: bold; padding: 4px 8px; border-radius: 4px; border: 2px solid #000000;", "background: transparent;", "color: #ffffff; font-style: italic; font-weight: bold;");
	return null;
})();
function invoke(msg) {
	_invoke?.(msg);
}
/**
* Checks if the current operating system is Windows.
*
* @return True if the operating system is Windows, otherwise false.
*/
function IsWindows() {
	return window._wails?.environment?.OS === "windows";
}
/**
* Checks if the current operating system is Linux.
*
* @returns Returns true if the current operating system is Linux, false otherwise.
*/
function IsLinux() {
	return window._wails?.environment?.OS === "linux";
}
/**
* Checks if the current environment is a macOS operating system.
*
* @returns True if the environment is macOS, false otherwise.
*/
function IsMac() {
	return window._wails?.environment?.OS === "darwin";
}
/**
* Reports whether the app is being run in debug mode.
*
* @returns True if the app is being run in debug mode.
*/
function IsDebug() {
	return Boolean(window._wails?.environment?.Debug);
}
// ../../../internal/runtime/desktop/@wailsio/runtime/src/utils.ts
/**
* Checks whether the webview supports the {@link MouseEvent#buttons} property.
* Looking at you macOS High Sierra!
*/
function canTrackButtons() {
	return new MouseEvent("mousedown").buttons === 0;
}
/**
* Resolves the closest HTMLElement ancestor of an event's target.
*/
function eventTarget(event) {
	if (event.target instanceof HTMLElement) return event.target;
	else if (!(event.target instanceof HTMLElement) && event.target instanceof Node) return event.target.parentElement ?? document.body;
	else return document.body;
}
/***
This technique for proper load detection is taken from HTMX:

BSD 2-Clause License

Copyright (c) 2020, Big Sky Software
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
this list of conditions and the following disclaimer in the documentation
and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

***/
var isReady = false;
if (hasDOM) document.addEventListener("DOMContentLoaded", () => {
	isReady = true;
});
function whenReady(callback) {
	if (isReady || document.readyState === "complete") callback();
	else document.addEventListener("DOMContentLoaded", callback);
}
// ../../../internal/runtime/desktop/@wailsio/runtime/src/contextmenu.ts
if (hasDOM) window.addEventListener("contextmenu", contextMenuHandler);
var call$2 = newRuntimeCaller(objectNames.ContextMenu);
var ContextMenuOpen = 0;
function openContextMenu(id, x, y, data) {
	call$2(ContextMenuOpen, {
		id,
		x,
		y,
		data
	});
}
function contextMenuHandler(event) {
	const target = eventTarget(event);
	const customContextMenu = window.getComputedStyle(target).getPropertyValue("--custom-contextmenu").trim();
	if (customContextMenu) {
		event.preventDefault();
		const data = window.getComputedStyle(target).getPropertyValue("--custom-contextmenu-data");
		openContextMenu(customContextMenu, event.clientX, event.clientY, data);
	} else processDefaultContextMenu(event, target);
}
function processDefaultContextMenu(event, target) {
	if (IsDebug()) return;
	switch (window.getComputedStyle(target).getPropertyValue("--default-contextmenu").trim()) {
		case "show": return;
		case "hide":
			event.preventDefault();
			return;
	}
	if (target.isContentEditable) return;
	const selection = window.getSelection();
	const hasSelection = selection && selection.toString().length > 0;
	if (hasSelection) for (let i = 0; i < selection.rangeCount; i++) {
		const rects = selection.getRangeAt(i).getClientRects();
		for (let j = 0; j < rects.length; j++) {
			const rect = rects[j];
			if (document.elementFromPoint(rect.left, rect.top) === target) return;
		}
	}
	if (target instanceof HTMLInputElement || target instanceof HTMLTextAreaElement) {
		if (hasSelection || !target.readOnly && !target.disabled) return;
	}
	event.preventDefault();
}
// ../../../internal/runtime/desktop/@wailsio/runtime/src/flags.ts
/**
* Retrieves the value associated with the specified key from the flag map.
*
* @param key - The key to retrieve the value for.
* @return The value associated with the specified key.
*/
function GetFlag(key) {
	try {
		return window._wails.flags[key];
	} catch (e) {
		throw new Error("Unable to retrieve flag '" + key + "': " + e, { cause: e });
	}
}
// ../../../internal/runtime/desktop/@wailsio/runtime/src/drag.ts
var canDrag = false;
var dragging = false;
var resizable = false;
var canResize = false;
var resizing = false;
var resizeEdge = "";
var defaultCursor = "auto";
var buttons = 0;
var buttonsTracked = false;
if (hasDOM) {
	buttonsTracked = canTrackButtons();
	window._wails = window._wails || {};
	window._wails.setResizable = (value) => {
		resizable = value;
		if (!resizable) {
			canResize = resizing = false;
			setResize();
		}
	};
}
var dragInitDone = false;
function isMobile() {
	const os = window._wails?.environment?.OS;
	if (os === "ios" || os === "android") return true;
	const ua = navigator.userAgent || navigator.vendor || window.opera || "";
	return /android|iphone|ipad|ipod|iemobile|wpdesktop/i.test(ua);
}
function tryInitDragHandlers() {
	if (dragInitDone) return;
	if (isMobile()) return;
	window.addEventListener("mousedown", update, { capture: true });
	window.addEventListener("mousemove", update, { capture: true });
	window.addEventListener("mouseup", update, { capture: true });
	for (const ev of [
		"click",
		"contextmenu",
		"dblclick"
	]) window.addEventListener(ev, suppressEvent, { capture: true });
	dragInitDone = true;
}
if (hasDOM) {
	tryInitDragHandlers();
	document.addEventListener("DOMContentLoaded", tryInitDragHandlers, { once: true });
	let dragEnvPolls = 0;
	const dragEnvPoll = window.setInterval(() => {
		if (dragInitDone) {
			window.clearInterval(dragEnvPoll);
			return;
		}
		tryInitDragHandlers();
		if (++dragEnvPolls > 100) window.clearInterval(dragEnvPoll);
	}, 50);
}
function suppressEvent(event) {
	if (event.type === "dblclick" && IsMac() && window === window.top && isDraggableEvent(event)) {
		event.stopImmediatePropagation();
		event.stopPropagation();
		event.preventDefault();
		invoke("wails:drag:doubleclick");
		return;
	}
	if (dragging || resizing) {
		event.stopImmediatePropagation();
		event.stopPropagation();
		event.preventDefault();
	}
}
function isDraggableEvent(event) {
	const target = eventTarget(event);
	return window.getComputedStyle(target).getPropertyValue("--wails-draggable").trim() === "drag" && event.offsetX >= 0 && event.offsetX < target.clientWidth && event.offsetY >= 0 && event.offsetY < target.clientHeight;
}
var MouseDown = 0;
var MouseUp = 1;
var MouseMove = 2;
function update(event) {
	let eventType, eventButtons = event.buttons;
	switch (event.type) {
		case "mousedown":
			eventType = MouseDown;
			if (!buttonsTracked) eventButtons = buttons | 1 << event.button;
			break;
		case "mouseup":
			eventType = MouseUp;
			if (!buttonsTracked) eventButtons = buttons & ~(1 << event.button);
			break;
		default:
			eventType = MouseMove;
			if (!buttonsTracked) eventButtons = buttons;
	}
	let released = buttons & ~eventButtons;
	let pressed = eventButtons & ~buttons;
	buttons = eventButtons;
	if (eventType === MouseDown && !(pressed & event.button)) {
		released |= 1 << event.button;
		pressed |= 1 << event.button;
	}
	if (eventType !== MouseMove && resizing || dragging && (eventType === MouseDown || event.button !== 0)) {
		event.stopImmediatePropagation();
		event.stopPropagation();
		event.preventDefault();
	}
	if (released & 1) primaryUp(event);
	if (pressed & 1) primaryDown(event);
	if (eventType === MouseMove) onMouseMove(event);
}
function primaryDown(event) {
	canDrag = false;
	canResize = false;
	if (!IsWindows()) {
		if (event.type === "mousedown" && event.button === 0 && event.detail !== 1) return;
	}
	if (resizeEdge) {
		if (event.type !== "mousedown") return;
		canResize = true;
		return;
	}
	canDrag = isDraggableEvent(event);
}
function primaryUp(event) {
	canDrag = false;
	dragging = false;
	canResize = false;
	resizing = false;
}
var cursorForEdge = Object.freeze({
	"se-resize": "nwse-resize",
	"sw-resize": "nesw-resize",
	"nw-resize": "nwse-resize",
	"ne-resize": "nesw-resize",
	"w-resize": "ew-resize",
	"n-resize": "ns-resize",
	"s-resize": "ns-resize",
	"e-resize": "ew-resize"
});
function setResize(edge) {
	if (edge) {
		if (!resizeEdge) defaultCursor = document.body.style.cursor;
		document.body.style.cursor = cursorForEdge[edge];
	} else if (!edge && resizeEdge) document.body.style.cursor = defaultCursor;
	resizeEdge = edge || "";
}
function onMouseMove(event) {
	if (canResize && resizeEdge) {
		resizing = true;
		invoke("wails:resize:" + resizeEdge);
	} else if (canDrag) {
		dragging = true;
		invoke("wails:drag");
	}
	if (dragging || resizing) {
		canDrag = canResize = false;
		return;
	}
	if (!resizable || !IsWindows() && !(IsLinux() && GetFlag("frameless"))) {
		if (resizeEdge) setResize();
		return;
	}
	const resizeHandleHeight = GetFlag("system.resizeHandleHeight") || 5;
	const resizeHandleWidth = GetFlag("system.resizeHandleWidth") || 5;
	const cornerExtra = GetFlag("resizeCornerExtra") || 10;
	const scrollbarWidth = Math.max(0, window.innerWidth - document.documentElement.clientWidth);
	const scrollbarHeight = Math.max(0, window.innerHeight - document.documentElement.clientHeight);
	const rightContentEdge = window.innerWidth - scrollbarWidth;
	const bottomContentEdge = window.innerHeight - scrollbarHeight;
	const rightBorder = event.clientX < rightContentEdge && rightContentEdge - event.clientX < resizeHandleWidth;
	const leftBorder = event.clientX < resizeHandleWidth;
	const topBorder = event.clientY < resizeHandleHeight;
	const bottomBorder = event.clientY < bottomContentEdge && bottomContentEdge - event.clientY < resizeHandleHeight;
	const rightCorner = event.clientX < rightContentEdge && rightContentEdge - event.clientX < resizeHandleWidth + cornerExtra;
	const leftCorner = event.clientX < resizeHandleWidth + cornerExtra;
	const topCorner = event.clientY < resizeHandleHeight + cornerExtra;
	const bottomCorner = event.clientY < bottomContentEdge && bottomContentEdge - event.clientY < resizeHandleHeight + cornerExtra;
	if (!leftCorner && !topCorner && !bottomCorner && !rightCorner) setResize();
	else if (rightCorner && bottomCorner) setResize("se-resize");
	else if (leftCorner && bottomCorner) setResize("sw-resize");
	else if (leftCorner && topCorner) setResize("nw-resize");
	else if (topCorner && rightCorner) setResize("ne-resize");
	else if (leftBorder) setResize("w-resize");
	else if (topBorder) setResize("n-resize");
	else if (bottomBorder) setResize("s-resize");
	else if (rightBorder) setResize("e-resize");
	else setResize();
}
// ../../../internal/runtime/desktop/@wailsio/runtime/src/appregion.ts
var regionProperty = "--wails-non-client-region";
var runtimeConfigReadyEvent = "wails:runtime-config-ready";
var validRegions = /* @__PURE__ */ new Set([
	"caption",
	"minimize",
	"maximize",
	"close"
]);
if (hasDOM) window._wails = window._wails || {};
var updatePending = false;
var lastPayload = "";
var observedElements = /* @__PURE__ */ new Set();
var resizeObserver;
var trackingStarted = false;
function normaliseRegionKind(value) {
	const region = value.trim().toLowerCase();
	if (validRegions.has(region)) return region;
}
function nonClientRegionForElement(element) {
	if (!(element instanceof HTMLElement)) return;
	const region = normaliseRegionKind(window.getComputedStyle(element).getPropertyValue(regionProperty));
	if (!region) return;
	const parent = element.parentElement;
	if (parent) {
		if (normaliseRegionKind(window.getComputedStyle(parent).getPropertyValue(regionProperty)) === region) return;
	}
	return region;
}
function isVisible(element) {
	const style = window.getComputedStyle(element);
	return style.display !== "none" && style.visibility !== "hidden" && style.contentVisibility !== "hidden";
}
function elementRegion(element) {
	if (!(element instanceof HTMLElement)) return;
	const kind = nonClientRegionForElement(element);
	if (!kind || !isVisible(element)) return;
	const rect = element.getBoundingClientRect();
	if (rect.width <= 0 || rect.height <= 0) return;
	const scale = window.devicePixelRatio || 1;
	const left = Math.floor(rect.left * scale);
	const top = Math.floor(rect.top * scale);
	const right = Math.ceil(rect.right * scale);
	const bottom = Math.ceil(rect.bottom * scale);
	if (right <= left || bottom <= top) return;
	return {
		kind,
		left,
		top,
		right,
		bottom
	};
}
function regionElements() {
	const elements = [];
	if (document.documentElement) elements.push(document.documentElement);
	if (document.body) {
		elements.push(document.body);
		for (const element of document.body.querySelectorAll("*")) elements.push(element);
	}
	return elements;
}
function observeRegionElements(elements) {
	if (typeof ResizeObserver === "undefined") return;
	resizeObserver ??= new ResizeObserver(scheduleUpdate);
	const nextElements = new Set(elements);
	for (const element of observedElements) if (!nextElements.has(element)) resizeObserver.unobserve(element);
	for (const element of nextElements) if (!observedElements.has(element)) resizeObserver.observe(element);
	observedElements = nextElements;
}
function updateNonClientRegions() {
	updatePending = false;
	const elements = regionElements();
	const regions = [];
	const activeElements = [];
	for (const element of elements) {
		const region = elementRegion(element);
		if (region) {
			regions.push(region);
			activeElements.push(element);
		}
	}
	observeRegionElements(activeElements);
	const payload = JSON.stringify({
		version: 1,
		regions
	});
	if (payload === lastPayload) return;
	lastPayload = payload;
	invoke("wails:non-client-region:" + payload);
}
function scheduleUpdate() {
	if (updatePending) return;
	updatePending = true;
	window.requestAnimationFrame(updateNonClientRegions);
}
function startNonClientRegionTracking() {
	if (trackingStarted) return;
	trackingStarted = true;
	scheduleUpdate();
	new MutationObserver(scheduleUpdate).observe(document.documentElement, {
		attributes: true,
		childList: true,
		subtree: true
	});
	window.addEventListener("resize", scheduleUpdate);
	window.addEventListener("scroll", scheduleUpdate, true);
	window.visualViewport?.addEventListener("resize", scheduleUpdate);
	window.visualViewport?.addEventListener("scroll", scheduleUpdate);
}
function tryStartNonClientRegionTracking() {
	const os = window._wails.environment?.OS;
	if (os === void 0) return false;
	const enabled = window._wails.flags?.nonClientRegionTracking;
	if (os === "windows") {
		if (enabled === true) whenReady(startNonClientRegionTracking);
		return true;
	}
	return true;
}
if (hasDOM && !tryStartNonClientRegionTracking()) window.addEventListener(runtimeConfigReadyEvent, tryStartNonClientRegionTracking, { once: true });
// ../../../internal/runtime/desktop/@wailsio/runtime/src/callable.ts
var fnToStr = Function.prototype.toString;
var reflectApply = typeof Reflect === "object" && Reflect !== null && Reflect.apply;
var badArrayLike;
var isCallableMarker;
if (typeof reflectApply === "function" && typeof Object.defineProperty === "function") try {
	badArrayLike = Object.defineProperty({}, "length", { get: function() {
		throw isCallableMarker;
	} });
	isCallableMarker = {};
	reflectApply(function() {
		throw 42;
	}, null, badArrayLike);
} catch (_) {
	if (_ !== isCallableMarker) reflectApply = null;
}
else reflectApply = null;
var constructorRegex = /^\s*class\b/;
var isES6ClassFn = function isES6ClassFunction(value) {
	try {
		var fnStr = fnToStr.call(value);
		return constructorRegex.test(fnStr);
	} catch (e) {
		return false;
	}
};
var tryFunctionObject = function tryFunctionToStr(value) {
	try {
		if (isES6ClassFn(value)) return false;
		fnToStr.call(value);
		return true;
	} catch (e) {
		return false;
	}
};
var toStr = Object.prototype.toString;
var objectClass = "[object Object]";
var fnClass = "[object Function]";
var genClass = "[object GeneratorFunction]";
var ddaClass = "[object HTMLAllCollection]";
var ddaClass2 = "[object HTML document.all class]";
var ddaClass3 = "[object HTMLCollection]";
var hasToStringTag = typeof Symbol === "function" && !!Symbol.toStringTag;
var isIE68 = !(0 in [,]);
var isDDA = function isDocumentDotAll() {
	return false;
};
if (typeof document === "object") {
	var all = document.all;
	if (toStr.call(all) === toStr.call(document.all)) isDDA = function isDocumentDotAll(value) {
		if ((isIE68 || !value) && (typeof value === "undefined" || typeof value === "object")) try {
			var str = toStr.call(value);
			return (str === ddaClass || str === ddaClass2 || str === ddaClass3 || str === objectClass) && value("") == null;
		} catch (e) {}
		return false;
	};
}
function isCallableRefApply(value) {
	if (isDDA(value)) return true;
	if (!value) return false;
	if (typeof value !== "function" && typeof value !== "object") return false;
	try {
		reflectApply(value, null, badArrayLike);
	} catch (e) {
		if (e !== isCallableMarker) return false;
	}
	return !isES6ClassFn(value) && tryFunctionObject(value);
}
function isCallableNoRefApply(value) {
	if (isDDA(value)) return true;
	if (!value) return false;
	if (typeof value !== "function" && typeof value !== "object") return false;
	if (hasToStringTag) return tryFunctionObject(value);
	if (isES6ClassFn(value)) return false;
	var strClass = toStr.call(value);
	if (strClass !== fnClass && strClass !== genClass && !/^\[object HTML/.test(strClass)) return false;
	return tryFunctionObject(value);
}
var callable_default = reflectApply ? isCallableRefApply : isCallableNoRefApply;
// ../../../internal/runtime/desktop/@wailsio/runtime/src/cancellable.ts
/**
* Exception class that will be used as rejection reason
* in case a {@link CancellablePromise} is cancelled successfully.
*
* The value of the {@link name} property is the string `"CancelError"`.
* The value of the {@link cause} property is the cause passed to the cancel method, if any.
*/
var CancelError = class extends Error {
	/**
	* Constructs a new `CancelError` instance.
	* @param message - The error message.
	* @param options - Options to be forwarded to the Error constructor.
	*/
	constructor(message, options) {
		super(message, options);
		this.name = "CancelError";
	}
};
/**
* Exception class that will be reported as an unhandled rejection
* in case a {@link CancellablePromise} rejects after being cancelled,
* or when the `oncancelled` callback throws or rejects.
*
* The value of the {@link name} property is the string `"CancelledRejectionError"`.
* The value of the {@link cause} property is the reason the promise rejected with.
*
* Because the original promise was cancelled,
* a wrapper promise will be passed to the unhandled rejection listener instead.
* The {@link promise} property holds a reference to the original promise.
*/
var CancelledRejectionError = class extends Error {
	/**
	* Constructs a new `CancelledRejectionError` instance.
	* @param promise - The promise that caused the error originally.
	* @param reason - The rejection reason.
	* @param info - An optional informative message specifying the circumstances in which the error was thrown.
	*               Defaults to the string `"Unhandled rejection in cancelled promise."`.
	*/
	constructor(promise, reason, info) {
		super((info ?? "Unhandled rejection in cancelled promise.") + " Reason: " + errorMessage(reason), { cause: reason });
		this.promise = promise;
		this.name = "CancelledRejectionError";
	}
};
var barrierSym = Symbol("barrier");
var cancelImplSym = Symbol("cancelImpl");
var species = Symbol.species ?? Symbol("speciesPolyfill");
/**
* A promise with an attached method for cancelling long-running operations (see {@link CancellablePromise#cancel}).
* Cancellation can optionally be bound to an {@link AbortSignal}
* for better composability (see {@link CancellablePromise#cancelOn}).
*
* Cancelling a pending promise will result in an immediate rejection
* with an instance of {@link CancelError} as reason,
* but whoever started the promise will be responsible
* for actually aborting the underlying operation.
* To this purpose, the constructor and all chaining methods
* accept optional cancellation callbacks.
*
* If a `CancellablePromise` still resolves after having been cancelled,
* the result will be discarded. If it rejects, the reason
* will be reported as an unhandled rejection,
* wrapped in a {@link CancelledRejectionError} instance.
* To facilitate the handling of cancellation requests,
* cancelled `CancellablePromise`s will _not_ report unhandled `CancelError`s
* whose `cause` field is the same as the one with which the current promise was cancelled.
*
* All usual promise methods are defined and return a `CancellablePromise`
* whose cancel method will cancel the parent operation as well, propagating the cancellation reason
* upwards through promise chains.
* Conversely, cancelling a promise will not automatically cancel dependent promises downstream:
* ```ts
* let root = new CancellablePromise((resolve, reject) => { ... });
* let child1 = root.then(() => { ... });
* let child2 = child1.then(() => { ... });
* let child3 = root.catch(() => { ... });
* child1.cancel(); // Cancels child1 and root, but not child2 or child3
* ```
* Cancelling a promise that has already settled is safe and has no consequence.
*
* The `cancel` method returns a promise that _always fulfills_
* after the whole chain has processed the cancel request
* and all attached callbacks up to that moment have run.
*
* All ES2024 promise methods (static and instance) are defined on CancellablePromise,
* but actual availability may vary with OS/webview version.
*
* In line with the proposal at https://github.com/tc39/proposal-rm-builtin-subclassing,
* `CancellablePromise` does not support transparent subclassing.
* Extenders should take care to provide their own method implementations.
* This might be reconsidered in case the proposal is retired.
*
* CancellablePromise is a wrapper around the DOM Promise object
* and is compliant with the [Promises/A+ specification](https://promisesaplus.com/)
* (it passes the [compliance suite](https://github.com/promises-aplus/promises-tests))
* if so is the underlying implementation.
*/
var CancellablePromise = class CancellablePromise extends Promise {
	/**
	* Creates a new `CancellablePromise`.
	*
	* @param executor - A callback used to initialize the promise. This callback is passed two arguments:
	*                   a `resolve` callback used to resolve the promise with a value
	*                   or the result of another promise (possibly cancellable),
	*                   and a `reject` callback used to reject the promise with a provided reason or error.
	*                   If the value provided to the `resolve` callback is a thenable _and_ cancellable object
	*                   (it has a `then` _and_ a `cancel` method),
	*                   cancellation requests will be forwarded to that object and the oncancelled will not be invoked anymore.
	*                   If any one of the two callbacks is called _after_ the promise has been cancelled,
	*                   the provided values will be cancelled and resolved as usual,
	*                   but their results will be discarded.
	*                   However, if the resolution process ultimately ends up in a rejection
	*                   that is not due to cancellation, the rejection reason
	*                   will be wrapped in a {@link CancelledRejectionError}
	*                   and bubbled up as an unhandled rejection.
	* @param oncancelled - It is the caller's responsibility to ensure that any operation
	*                      started by the executor is properly halted upon cancellation.
	*                      This optional callback can be used to that purpose.
	*                      It will be called _synchronously_ with a cancellation cause
	*                      when cancellation is requested, _after_ the promise has already rejected
	*                      with a {@link CancelError}, but _before_
	*                      any {@link then}/{@link catch}/{@link finally} callback runs.
	*                      If the callback returns a thenable, the promise returned from {@link cancel}
	*                      will only fulfill after the former has settled.
	*                      Unhandled exceptions or rejections from the callback will be wrapped
	*                      in a {@link CancelledRejectionError} and bubbled up as unhandled rejections.
	*                      If the `resolve` callback is called before cancellation with a cancellable promise,
	*                      cancellation requests on this promise will be diverted to that promise,
	*                      and the original `oncancelled` callback will be discarded.
	*/
	constructor(executor, oncancelled) {
		let resolve;
		let reject;
		super((res, rej) => {
			resolve = res;
			reject = rej;
		});
		if (this.constructor[species] !== Promise) throw new TypeError("CancellablePromise does not support transparent subclassing. Please refrain from overriding the [Symbol.species] static property.");
		let promise = {
			promise: this,
			resolve,
			reject,
			get oncancelled() {
				return oncancelled ?? null;
			},
			set oncancelled(cb) {
				oncancelled = cb ?? void 0;
			}
		};
		const state = {
			get root() {
				return state;
			},
			resolving: false,
			settled: false
		};
		Object.defineProperties(this, {
			[barrierSym]: {
				configurable: false,
				enumerable: false,
				writable: true,
				value: null
			},
			[cancelImplSym]: {
				configurable: false,
				enumerable: false,
				writable: false,
				value: cancellerFor(promise, state)
			}
		});
		const rejector = rejectorFor(promise, state);
		try {
			executor(resolverFor(promise, state), rejector);
		} catch (err) {
			if (state.resolving) console.log("Unhandled exception in CancellablePromise executor.", err);
			else rejector(err);
		}
	}
	/**
	* Cancels immediately the execution of the operation associated with this promise.
	* The promise rejects with a {@link CancelError} instance as reason,
	* with the {@link CancelError#cause} property set to the given argument, if any.
	*
	* Has no effect if called after the promise has already settled;
	* repeated calls in particular are safe, but only the first one
	* will set the cancellation cause.
	*
	* The `CancelError` exception _need not_ be handled explicitly _on the promises that are being cancelled:_
	* cancelling a promise with no attached rejection handler does not trigger an unhandled rejection event.
	* Therefore, the following idioms are all equally correct:
	* ```ts
	* new CancellablePromise((resolve, reject) => { ... }).cancel();
	* new CancellablePromise((resolve, reject) => { ... }).then(...).cancel();
	* new CancellablePromise((resolve, reject) => { ... }).then(...).catch(...).cancel();
	* ```
	* Whenever some cancelled promise in a chain rejects with a `CancelError`
	* with the same cancellation cause as itself, the error will be discarded silently.
	* However, the `CancelError` _will still be delivered_ to all attached rejection handlers
	* added by {@link then} and related methods:
	* ```ts
	* let cancellable = new CancellablePromise((resolve, reject) => { ... });
	* cancellable.then(() => { ... }).catch(console.log);
	* cancellable.cancel(); // A CancelError is printed to the console.
	* ```
	* If the `CancelError` is not handled downstream by the time it reaches
	* a _non-cancelled_ promise, it _will_ trigger an unhandled rejection event,
	* just like normal rejections would:
	* ```ts
	* let cancellable = new CancellablePromise((resolve, reject) => { ... });
	* let chained = cancellable.then(() => { ... }).then(() => { ... }); // No catch...
	* cancellable.cancel(); // Unhandled rejection event on chained!
	* ```
	* Therefore, it is important to either cancel whole promise chains from their tail,
	* as shown in the correct idioms above, or take care of handling errors everywhere.
	*
	* @returns A cancellable promise that _fulfills_ after the cancel callback (if any)
	* and all handlers attached up to the call to cancel have run.
	* If the cancel callback returns a thenable, the promise returned by `cancel`
	* will also wait for that thenable to settle.
	* This enables callers to wait for the cancelled operation to terminate
	* without being forced to handle potential errors at the call site.
	* ```ts
	* cancellable.cancel().then(() => {
	*     // Cleanup finished, it's safe to do something else.
	* }, (err) => {
	*     // Unreachable: the promise returned from cancel will never reject.
	* });
	* ```
	* Note that the returned promise will _not_ handle implicitly any rejection
	* that might have occurred already in the cancelled chain.
	* It will just track whether registered handlers have been executed or not.
	* Therefore, unhandled rejections will never be silently handled by calling cancel.
	*/
	cancel(cause) {
		return new CancellablePromise((resolve) => {
			Promise.all([this[cancelImplSym](new CancelError("Promise cancelled.", { cause })), currentBarrier(this)]).then(() => resolve(), () => resolve());
		});
	}
	/**
	* Binds promise cancellation to the abort event of the given {@link AbortSignal}.
	* If the signal has already aborted, the promise will be cancelled immediately.
	* When either condition is verified, the cancellation cause will be set
	* to the signal's abort reason (see {@link AbortSignal#reason}).
	*
	* Has no effect if called (or if the signal aborts) _after_ the promise has already settled.
	* Only the first signal to abort will set the cancellation cause.
	*
	* For more details about the cancellation process,
	* see {@link cancel} and the `CancellablePromise` constructor.
	*
	* This method enables `await`ing cancellable promises without having
	* to store them for future cancellation, e.g.:
	* ```ts
	* await longRunningOperation().cancelOn(signal);
	* ```
	* instead of:
	* ```ts
	* let promiseToBeCancelled = longRunningOperation();
	* await promiseToBeCancelled;
	* ```
	*
	* @returns This promise, for method chaining.
	*/
	cancelOn(signal) {
		if (signal.aborted) this.cancel(signal.reason);
		else signal.addEventListener("abort", () => void this.cancel(signal.reason), { capture: true });
		return this;
	}
	/**
	* Attaches callbacks for the resolution and/or rejection of the `CancellablePromise`.
	*
	* The optional `oncancelled` argument will be invoked when the returned promise is cancelled,
	* with the same semantics as the `oncancelled` argument of the constructor.
	* When the parent promise rejects or is cancelled, the `onrejected` callback will run,
	* _even after the returned promise has been cancelled:_
	* in that case, should it reject or throw, the reason will be wrapped
	* in a {@link CancelledRejectionError} and bubbled up as an unhandled rejection.
	*
	* @param onfulfilled The callback to execute when the Promise is resolved.
	* @param onrejected The callback to execute when the Promise is rejected.
	* @returns A `CancellablePromise` for the completion of whichever callback is executed.
	* The returned promise is hooked up to propagate cancellation requests up the chain, but not down:
	*
	*   - if the parent promise is cancelled, the `onrejected` handler will be invoked with a `CancelError`
	*     and the returned promise _will resolve regularly_ with its result;
	*   - conversely, if the returned promise is cancelled, _the parent promise is cancelled too;_
	*     the `onrejected` handler will still be invoked with the parent's `CancelError`,
	*     but its result will be discarded
	*     and the returned promise will reject with a `CancelError` as well.
	*
	* The promise returned from {@link cancel} will fulfill only after all attached handlers
	* up the entire promise chain have been run.
	*
	* If either callback returns a cancellable promise,
	* cancellation requests will be diverted to it,
	* and the specified `oncancelled` callback will be discarded.
	*/
	then(onfulfilled, onrejected, oncancelled) {
		if (!(this instanceof CancellablePromise)) throw new TypeError("CancellablePromise.prototype.then called on an invalid object.");
		if (!callable_default(onfulfilled)) onfulfilled = identity;
		if (!callable_default(onrejected)) onrejected = thrower;
		if (onfulfilled === identity && onrejected == thrower) return new CancellablePromise((resolve) => resolve(this));
		const barrier = {};
		this[barrierSym] = barrier;
		return new CancellablePromise((resolve, reject) => {
			super.then((value) => {
				if (this[barrierSym] === barrier) this[barrierSym] = null;
				barrier.resolve?.();
				try {
					resolve(onfulfilled(value));
				} catch (err) {
					reject(err);
				}
			}, (reason) => {
				if (this[barrierSym] === barrier) this[barrierSym] = null;
				barrier.resolve?.();
				try {
					resolve(onrejected(reason));
				} catch (err) {
					reject(err);
				}
			});
		}, async (cause) => {
			try {
				return oncancelled?.(cause);
			} finally {
				await this.cancel(cause);
			}
		});
	}
	/**
	* Attaches a callback for only the rejection of the Promise.
	*
	* The optional `oncancelled` argument will be invoked when the returned promise is cancelled,
	* with the same semantics as the `oncancelled` argument of the constructor.
	* When the parent promise rejects or is cancelled, the `onrejected` callback will run,
	* _even after the returned promise has been cancelled:_
	* in that case, should it reject or throw, the reason will be wrapped
	* in a {@link CancelledRejectionError} and bubbled up as an unhandled rejection.
	*
	* It is equivalent to
	* ```ts
	* cancellablePromise.then(undefined, onrejected, oncancelled);
	* ```
	* and the same caveats apply.
	*
	* @returns A Promise for the completion of the callback.
	* Cancellation requests on the returned promise
	* will propagate up the chain to the parent promise,
	* but not in the other direction.
	*
	* The promise returned from {@link cancel} will fulfill only after all attached handlers
	* up the entire promise chain have been run.
	*
	* If `onrejected` returns a cancellable promise,
	* cancellation requests will be diverted to it,
	* and the specified `oncancelled` callback will be discarded.
	* See {@link then} for more details.
	*/
	catch(onrejected, oncancelled) {
		return this.then(void 0, onrejected, oncancelled);
	}
	/**
	* Attaches a callback that is invoked when the CancellablePromise is settled (fulfilled or rejected). The
	* resolved value cannot be accessed or modified from the callback.
	* The returned promise will settle in the same state as the original one
	* after the provided callback has completed execution,
	* unless the callback throws or returns a rejecting promise,
	* in which case the returned promise will reject as well.
	*
	* The optional `oncancelled` argument will be invoked when the returned promise is cancelled,
	* with the same semantics as the `oncancelled` argument of the constructor.
	* Once the parent promise settles, the `onfinally` callback will run,
	* _even after the returned promise has been cancelled:_
	* in that case, should it reject or throw, the reason will be wrapped
	* in a {@link CancelledRejectionError} and bubbled up as an unhandled rejection.
	*
	* This method is implemented in terms of {@link then} and the same caveats apply.
	* It is polyfilled, hence available in every OS/webview version.
	*
	* @returns A Promise for the completion of the callback.
	* Cancellation requests on the returned promise
	* will propagate up the chain to the parent promise,
	* but not in the other direction.
	*
	* The promise returned from {@link cancel} will fulfill only after all attached handlers
	* up the entire promise chain have been run.
	*
	* If `onfinally` returns a cancellable promise,
	* cancellation requests will be diverted to it,
	* and the specified `oncancelled` callback will be discarded.
	* See {@link then} for more details.
	*/
	finally(onfinally, oncancelled) {
		if (!(this instanceof CancellablePromise)) throw new TypeError("CancellablePromise.prototype.finally called on an invalid object.");
		if (!callable_default(onfinally)) return this.then(onfinally, onfinally, oncancelled);
		return this.then((value) => CancellablePromise.resolve(onfinally()).then(() => value), (reason) => CancellablePromise.resolve(onfinally()).then(() => {
			throw reason;
		}), oncancelled);
	}
	/**
	* We use the `[Symbol.species]` static property, if available,
	* to disable the built-in automatic subclassing features from {@link Promise}.
	* It is critical for performance reasons that extenders do not override this.
	* Once the proposal at https://github.com/tc39/proposal-rm-builtin-subclassing
	* is either accepted or retired, this implementation will have to be revised accordingly.
	*
	* @ignore
	* @internal
	*/
	static get [species]() {
		return Promise;
	}
	static all(values) {
		let collected = Array.from(values);
		const promise = collected.length === 0 ? CancellablePromise.resolve(collected) : new CancellablePromise((resolve, reject) => {
			Promise.all(collected).then(resolve, reject);
		}, (cause) => cancelAll(promise, collected, cause));
		return promise;
	}
	static allSettled(values) {
		let collected = Array.from(values);
		const promise = collected.length === 0 ? CancellablePromise.resolve(collected) : new CancellablePromise((resolve, reject) => {
			Promise.allSettled(collected).then(resolve, reject);
		}, (cause) => cancelAll(promise, collected, cause));
		return promise;
	}
	static any(values) {
		let collected = Array.from(values);
		const promise = collected.length === 0 ? CancellablePromise.resolve(collected) : new CancellablePromise((resolve, reject) => {
			Promise.any(collected).then(resolve, reject);
		}, (cause) => cancelAll(promise, collected, cause));
		return promise;
	}
	static race(values) {
		let collected = Array.from(values);
		const promise = new CancellablePromise((resolve, reject) => {
			Promise.race(collected).then(resolve, reject);
		}, (cause) => cancelAll(promise, collected, cause));
		return promise;
	}
	/**
	* Creates a new cancelled CancellablePromise for the provided cause.
	*
	* @group Static Methods
	*/
	static cancel(cause) {
		const p = new CancellablePromise(() => {});
		p.cancel(cause);
		return p;
	}
	/**
	* Creates a new CancellablePromise that cancels
	* after the specified timeout, with the provided cause.
	*
	* If the {@link AbortSignal.timeout} factory method is available,
	* it is used to base the timeout on _active_ time rather than _elapsed_ time.
	* Otherwise, `timeout` falls back to {@link setTimeout}.
	*
	* @group Static Methods
	*/
	static timeout(milliseconds, cause) {
		const promise = new CancellablePromise(() => {});
		if (AbortSignal && typeof AbortSignal === "function" && AbortSignal.timeout && typeof AbortSignal.timeout === "function") AbortSignal.timeout(milliseconds).addEventListener("abort", () => void promise.cancel(cause));
		else setTimeout(() => void promise.cancel(cause), milliseconds);
		return promise;
	}
	static sleep(milliseconds, value) {
		return new CancellablePromise((resolve) => {
			setTimeout(() => resolve(value), milliseconds);
		});
	}
	/**
	* Creates a new rejected CancellablePromise for the provided reason.
	*
	* @group Static Methods
	*/
	static reject(reason) {
		return new CancellablePromise((_, reject) => reject(reason));
	}
	static resolve(value) {
		if (value instanceof CancellablePromise) return value;
		return new CancellablePromise((resolve) => resolve(value));
	}
	/**
	* Creates a new CancellablePromise and returns it in an object, along with its resolve and reject functions
	* and a getter/setter for the cancellation callback.
	*
	* This method is polyfilled, hence available in every OS/webview version.
	*
	* @group Static Methods
	*/
	static withResolvers() {
		let result = { oncancelled: null };
		result.promise = new CancellablePromise((resolve, reject) => {
			result.resolve = resolve;
			result.reject = reject;
		}, (cause) => {
			result.oncancelled?.(cause);
		});
		return result;
	}
};
/**
* Returns a callback that implements the cancellation algorithm for the given cancellable promise.
* The promise returned from the resulting function does not reject.
*/
function cancellerFor(promise, state) {
	let cancellationPromise = void 0;
	return (reason) => {
		if (!state.settled) {
			state.settled = true;
			state.reason = reason;
			promise.reject(reason);
			Promise.prototype.then.call(promise.promise, void 0, (err) => {
				if (err !== reason) throw err;
			});
		}
		if (!state.reason || !promise.oncancelled) return;
		cancellationPromise = new Promise((resolve) => {
			try {
				resolve(promise.oncancelled(state.reason.cause));
			} catch (err) {
				Promise.reject(new CancelledRejectionError(promise.promise, err, "Unhandled exception in oncancelled callback."));
			}
		}).catch((reason) => {
			Promise.reject(new CancelledRejectionError(promise.promise, reason, "Unhandled rejection in oncancelled callback."));
		});
		promise.oncancelled = null;
		return cancellationPromise;
	};
}
/**
* Returns a callback that implements the resolution algorithm for the given cancellable promise.
*/
function resolverFor(promise, state) {
	return (value) => {
		if (state.resolving) return;
		state.resolving = true;
		if (value === promise.promise) {
			if (state.settled) return;
			state.settled = true;
			promise.reject(/* @__PURE__ */ new TypeError("A promise cannot be resolved with itself."));
			return;
		}
		if (value != null && (typeof value === "object" || typeof value === "function")) {
			let then;
			try {
				then = value.then;
			} catch (err) {
				state.settled = true;
				promise.reject(err);
				return;
			}
			if (callable_default(then)) {
				try {
					let cancel = value.cancel;
					if (callable_default(cancel)) {
						const oncancelled = (cause) => {
							Reflect.apply(cancel, value, [cause]);
						};
						if (state.reason) cancellerFor({
							...promise,
							oncancelled
						}, state)(state.reason);
						else promise.oncancelled = oncancelled;
					}
				} catch {}
				const newState = {
					root: state.root,
					resolving: false,
					get settled() {
						return this.root.settled;
					},
					set settled(value) {
						this.root.settled = value;
					},
					get reason() {
						return this.root.reason;
					}
				};
				const rejector = rejectorFor(promise, newState);
				try {
					Reflect.apply(then, value, [resolverFor(promise, newState), rejector]);
				} catch (err) {
					rejector(err);
				}
				return;
			}
		}
		if (state.settled) return;
		state.settled = true;
		promise.resolve(value);
	};
}
/**
* Returns a callback that implements the rejection algorithm for the given cancellable promise.
*/
function rejectorFor(promise, state) {
	return (reason) => {
		if (state.resolving) return;
		state.resolving = true;
		if (state.settled) {
			try {
				if (reason instanceof CancelError && state.reason instanceof CancelError && Object.is(reason.cause, state.reason.cause)) return;
			} catch {}
			Promise.reject(new CancelledRejectionError(promise.promise, reason));
		} else {
			state.settled = true;
			promise.reject(reason);
		}
	};
}
/**
* Cancels all values in an array that look like cancellable thenables.
* Returns a promise that fulfills once all cancellation procedures for the given values have settled.
*/
function cancelAll(parent, values, cause) {
	const results = [];
	for (const value of values) {
		let cancel;
		try {
			if (!callable_default(value.then)) continue;
			cancel = value.cancel;
			if (!callable_default(cancel)) continue;
		} catch {
			continue;
		}
		let result;
		try {
			result = Reflect.apply(cancel, value, [cause]);
		} catch (err) {
			Promise.reject(new CancelledRejectionError(parent, err, "Unhandled exception in cancel method."));
			continue;
		}
		if (!result) continue;
		results.push((result instanceof Promise ? result : Promise.resolve(result)).catch((reason) => {
			Promise.reject(new CancelledRejectionError(parent, reason, "Unhandled rejection in cancel method."));
		}));
	}
	return Promise.all(results);
}
/**
* Returns its argument.
*/
function identity(x) {
	return x;
}
/**
* Throws its argument.
*/
function thrower(reason) {
	throw reason;
}
/**
* Attempts various strategies to convert an error to a string.
*/
function errorMessage(err) {
	try {
		if (err instanceof Error || typeof err !== "object" || err.toString !== Object.prototype.toString) return "" + err;
	} catch {}
	try {
		return JSON.stringify(err);
	} catch {}
	try {
		return Object.prototype.toString.call(err);
	} catch {}
	return "<could not convert error to string>";
}
/**
* Gets the current barrier promise for the given cancellable promise. If necessary, initialises the barrier.
*/
function currentBarrier(promise) {
	let pwr = promise[barrierSym] ?? {};
	if (!("promise" in pwr)) Object.assign(pwr, promiseWithResolvers());
	if (promise[barrierSym] == null) {
		pwr.resolve();
		promise[barrierSym] = pwr;
	}
	return pwr.promise;
}
var promiseWithResolvers = Promise.withResolvers;
if (promiseWithResolvers && typeof promiseWithResolvers === "function") promiseWithResolvers = promiseWithResolvers.bind(Promise);
else promiseWithResolvers = function() {
	let resolve;
	let reject;
	return {
		promise: new Promise((res, rej) => {
			resolve = res;
			reject = rej;
		}),
		resolve,
		reject
	};
};
// ../../../internal/runtime/desktop/@wailsio/runtime/src/calls.ts
if (hasDOM) window._wails = window._wails || {};
var call$1 = newRuntimeCaller(objectNames.Call);
var cancelCall = newRuntimeCaller(objectNames.CancelCall);
var callResponses = /* @__PURE__ */ new Map();
var CallBinding = 0;
var CancelMethod = 0;
/**
* Generates a unique ID using the nanoid library.
*
* @returns A unique ID that does not exist in the callResponses set.
*/
function generateID() {
	let result;
	do
		result = nanoid();
	while (callResponses.has(result));
	return result;
}
/**
* Call a bound method according to the given call options.
*
* In case of failure, the returned promise will reject with an exception
* among ReferenceError (unknown method), TypeError (wrong argument count or type),
* {@link RuntimeError} (method returned an error), or other (network or internal errors).
* The exception might have a "cause" field with the value returned
* by the application- or service-level error marshaling functions.
*
* @param options - A method call descriptor.
* @returns The result of the call.
*/
function Call(options) {
	const id = generateID();
	const result = CancellablePromise.withResolvers();
	callResponses.set(id, {
		resolve: result.resolve,
		reject: result.reject
	});
	const request = call$1(CallBinding, Object.assign({ "call-id": id }, options));
	let running = true;
	request.then((res) => {
		running = false;
		callResponses.delete(id);
		result.resolve(res);
	}, (err) => {
		running = false;
		callResponses.delete(id);
		result.reject(err);
	});
	const cancel = () => {
		callResponses.delete(id);
		return cancelCall(CancelMethod, { "call-id": id }).catch((err) => {
			console.error("Error while requesting binding call cancellation:", err);
		});
	};
	result.oncancelled = () => {
		if (running) return cancel();
		else return request.then(cancel);
	};
	return result.promise;
}
/**
* Calls a method by its numeric ID with the specified arguments.
* See {@link Call} for details.
*
* @param methodID - The ID of the method to call.
* @param args - The arguments to pass to the method.
* @return The result of the method call.
*/
function ByID(methodID, ...args) {
	return Call({
		methodID,
		args
	});
}
// ../../../internal/runtime/desktop/@wailsio/runtime/src/create.ts
/**
* Maps known event names to creation functions for their data types.
* Will be monkey-patched by the binding generator.
*/
var Events = {};
// ../../../internal/runtime/desktop/@wailsio/runtime/src/listener.ts
var eventListeners = /* @__PURE__ */ new Map();
var Listener = class {
	constructor(eventName, callback, maxCallbacks) {
		this.eventName = eventName;
		this.callback = callback;
		this.maxCallbacks = maxCallbacks || -1;
	}
	dispatch(data) {
		try {
			this.callback(data);
		} catch (err) {
			console.error(err);
		}
		if (this.maxCallbacks === -1) return false;
		this.maxCallbacks -= 1;
		return this.maxCallbacks === 0;
	}
};
function listenerOff(listener) {
	let listeners = eventListeners.get(listener.eventName);
	if (!listeners) return;
	listeners = listeners.filter((l) => l !== listener);
	if (listeners.length === 0) eventListeners.delete(listener.eventName);
	else eventListeners.set(listener.eventName, listeners);
}
// ../../../internal/runtime/desktop/@wailsio/runtime/src/events.ts
if (hasDOM) {
	window._wails = window._wails || {};
	window._wails.dispatchWailsEvent = dispatchWailsEvent;
}
objectNames.Events;
/**
* Represents a system event or a custom event emitted through wails-provided facilities.
*/
var WailsEvent = class {
	constructor(name, data) {
		this.name = name;
		this.data = data ?? null;
	}
};
function dispatchWailsEvent(event) {
	let listeners = eventListeners.get(event.name);
	if (!listeners) return;
	let wailsEvent = new WailsEvent(event.name, event.name in Events ? Events[event.name](event.data) : event.data);
	if ("sender" in event) wailsEvent.sender = event.sender;
	const expired = /* @__PURE__ */ new Set();
	for (const listener of listeners.slice()) if (listener.dispatch(wailsEvent)) expired.add(listener);
	if (expired.size > 0) {
		const live = eventListeners.get(event.name);
		if (live) {
			const remaining = live.filter((l) => !expired.has(l));
			if (remaining.length === 0) eventListeners.delete(event.name);
			else eventListeners.set(event.name, remaining);
		}
	}
}
/**
* Register a callback function to be called multiple times for a specific event.
*
* @param eventName - The name of the event to register the callback for.
* @param callback - The callback function to be called when the event is triggered.
* @param maxCallbacks - The maximum number of times the callback can be called for the event. Once the maximum number is reached, the callback will no longer be called.
* @returns A function that, when called, will unregister the callback from the event.
*/
function OnMultiple(eventName, callback, maxCallbacks) {
	let listeners = eventListeners.get(eventName) || [];
	const thisListener = new Listener(eventName, callback, maxCallbacks);
	listeners.push(thisListener);
	eventListeners.set(eventName, listeners);
	return () => listenerOff(thisListener);
}
/**
* Registers a callback function to be executed when the specified event occurs.
*
* @param eventName - The name of the event to register the callback for.
* @param callback - The callback function to be called when the event is triggered.
* @returns A function that, when called, will unregister the callback from the event.
*/
function On(eventName, callback) {
	return OnMultiple(eventName, callback, -1);
}
// ../../../internal/runtime/desktop/@wailsio/runtime/src/window.ts
var DROP_TARGET_ATTRIBUTE = "data-file-drop-target";
var DROP_TARGET_ACTIVE_CLASS = "file-drop-target-active";
var currentDropTarget = null;
var PositionMethod = 0;
var CenterMethod = 1;
var CloseMethod = 2;
var DisableSizeConstraintsMethod = 3;
var EnableSizeConstraintsMethod = 4;
var FocusMethod = 5;
var ForceReloadMethod = 6;
var FullscreenMethod = 7;
var GetScreenMethod = 8;
var GetZoomMethod = 9;
var HeightMethod = 10;
var HideMethod = 11;
var IsFocusedMethod = 12;
var IsFullscreenMethod = 13;
var IsMaximisedMethod = 14;
var IsMinimisedMethod = 15;
var MaximiseMethod = 16;
var MinimiseMethod = 17;
var NameMethod = 18;
var OpenDevToolsMethod = 19;
var RelativePositionMethod = 20;
var ReloadMethod = 21;
var ResizableMethod = 22;
var RestoreMethod = 23;
var SetPositionMethod = 24;
var SetAlwaysOnTopMethod = 25;
var SetBackgroundColourMethod = 26;
var SetFramelessMethod = 27;
var SetFullscreenButtonEnabledMethod = 28;
var SetMaxSizeMethod = 29;
var SetMinSizeMethod = 30;
var SetRelativePositionMethod = 31;
var SetResizableMethod = 32;
var SetSizeMethod = 33;
var SetTitleMethod = 34;
var SetZoomMethod = 35;
var ShowMethod = 36;
var SizeMethod = 37;
var ToggleFullscreenMethod = 38;
var ToggleMaximiseMethod = 39;
var ToggleFramelessMethod = 40;
var UnFullscreenMethod = 41;
var UnMaximiseMethod = 42;
var UnMinimiseMethod = 43;
var WidthMethod = 44;
var ZoomMethod = 45;
var ZoomInMethod = 46;
var ZoomOutMethod = 47;
var ZoomResetMethod = 48;
var SnapAssistMethod = 49;
var FilesDropped = 50;
var PrintMethod = 51;
var SetScreenMethod = 52;
/**
* Finds the nearest drop target element by walking up the DOM tree.
*/
function getDropTargetElement(element) {
	if (!element) return null;
	return element.closest(`[${DROP_TARGET_ATTRIBUTE}]`);
}
/**
* Check if we can use WebView2's postMessageWithAdditionalObjects (Windows)
* Also checks that EnableFileDrop is true for this window.
*/
function canResolveFilePaths() {
	if (window.chrome?.webview?.postMessageWithAdditionalObjects == null) return false;
	return window._wails?.flags?.enableFileDrop === true;
}
/**
* Send file drop to backend via WebView2 (Windows only)
*/
function resolveFilePaths(x, y, files) {
	if (window.chrome?.webview?.postMessageWithAdditionalObjects) window.chrome.webview.postMessageWithAdditionalObjects(`file:drop:${x}:${y}`, files);
}
var nativeDragActive = false;
/**
* Cleans up native drag state and hover effects.
* Called on drop or when drag leaves the window.
*/
function cleanupNativeDrag() {
	nativeDragActive = false;
	if (currentDropTarget) {
		currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
		currentDropTarget = null;
	}
}
/**
* Called from Go when a file drag enters the window on Linux/macOS.
*/
function handleDragEnter() {
	if (window._wails?.flags?.enableFileDrop === false) return;
	nativeDragActive = true;
}
/**
* Called from Go when a file drag leaves the window on Linux/macOS.
*/
function handleDragLeave() {
	cleanupNativeDrag();
}
/**
* Called from Go during file drag to update hover state on Linux/macOS.
* @param x - X coordinate in CSS pixels
* @param y - Y coordinate in CSS pixels
*/
function handleDragOver(x, y) {
	if (!nativeDragActive) return;
	if (window._wails?.flags?.enableFileDrop === false) return;
	const dropTarget = getDropTargetElement(document.elementFromPoint(x, y));
	if (currentDropTarget && currentDropTarget !== dropTarget) currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
	if (dropTarget) {
		dropTarget.classList.add(DROP_TARGET_ACTIVE_CLASS);
		currentDropTarget = dropTarget;
	} else currentDropTarget = null;
}
var callerSym = Symbol("caller");
/**
* The window within which the script is running.
*/
var thisWindow = new class Window {
	/**
	* Initialises a window object with the specified name.
	*
	* @private
	* @param name - The name of the target window.
	*/
	constructor(name = "") {
		this[callerSym] = newRuntimeCaller(objectNames.Window, name);
		for (const method of Object.getOwnPropertyNames(Window.prototype)) if (method !== "constructor" && typeof this[method] === "function") this[method] = this[method].bind(this);
	}
	/**
	* Gets the specified window.
	*
	* @param name - The name of the window to get.
	* @returns The corresponding window object.
	*/
	Get(name) {
		return new Window(name);
	}
	/**
	* Returns the absolute position of the window.
	*
	* @returns The current absolute position of the window.
	*/
	Position() {
		return this[callerSym](PositionMethod);
	}
	/**
	* Centers the window on the screen.
	*/
	Center() {
		return this[callerSym](CenterMethod);
	}
	/**
	* Closes the window.
	*/
	Close() {
		return this[callerSym](CloseMethod);
	}
	/**
	* Disables min/max size constraints.
	*/
	DisableSizeConstraints() {
		return this[callerSym](DisableSizeConstraintsMethod);
	}
	/**
	* Enables min/max size constraints.
	*/
	EnableSizeConstraints() {
		return this[callerSym](EnableSizeConstraintsMethod);
	}
	/**
	* Focuses the window.
	*/
	Focus() {
		return this[callerSym](FocusMethod);
	}
	/**
	* Forces the window to reload the page assets.
	*/
	ForceReload() {
		return this[callerSym](ForceReloadMethod);
	}
	/**
	* Switches the window to fullscreen mode.
	*/
	Fullscreen() {
		return this[callerSym](FullscreenMethod);
	}
	/**
	* Returns the screen that the window is on.
	*
	* @returns The screen the window is currently on.
	*/
	GetScreen() {
		return this[callerSym](GetScreenMethod);
	}
	/**
	* Returns the current zoom level of the window.
	*
	* @returns The current zoom level.
	*/
	GetZoom() {
		return this[callerSym](GetZoomMethod);
	}
	/**
	* Returns the height of the window.
	*
	* @returns The current height of the window.
	*/
	Height() {
		return this[callerSym](HeightMethod);
	}
	/**
	* Hides the window.
	*/
	Hide() {
		return this[callerSym](HideMethod);
	}
	/**
	* Returns true if the window is focused.
	*
	* @returns Whether the window is currently focused.
	*/
	IsFocused() {
		return this[callerSym](IsFocusedMethod);
	}
	/**
	* Returns true if the window is fullscreen.
	*
	* @returns Whether the window is currently fullscreen.
	*/
	IsFullscreen() {
		return this[callerSym](IsFullscreenMethod);
	}
	/**
	* Returns true if the window is maximised.
	*
	* @returns Whether the window is currently maximised.
	*/
	IsMaximised() {
		return this[callerSym](IsMaximisedMethod);
	}
	/**
	* Returns true if the window is minimised.
	*
	* @returns Whether the window is currently minimised.
	*/
	IsMinimised() {
		return this[callerSym](IsMinimisedMethod);
	}
	/**
	* Maximises the window.
	*/
	Maximise() {
		return this[callerSym](MaximiseMethod);
	}
	/**
	* Minimises the window.
	*/
	Minimise() {
		return this[callerSym](MinimiseMethod);
	}
	/**
	* Returns the name of the window.
	*
	* @returns The name of the window.
	*/
	Name() {
		return this[callerSym](NameMethod);
	}
	/**
	* Opens the development tools pane.
	*/
	OpenDevTools() {
		return this[callerSym](OpenDevToolsMethod);
	}
	/**
	* Returns the relative position of the window to the screen.
	*
	* @returns The current relative position of the window.
	*/
	RelativePosition() {
		return this[callerSym](RelativePositionMethod);
	}
	/**
	* Reloads the page assets.
	*/
	Reload() {
		return this[callerSym](ReloadMethod);
	}
	/**
	* Returns true if the window is resizable.
	*
	* @returns Whether the window is currently resizable.
	*/
	Resizable() {
		return this[callerSym](ResizableMethod);
	}
	/**
	* Restores the window to its previous state if it was previously minimised, maximised or fullscreen.
	*/
	Restore() {
		return this[callerSym](RestoreMethod);
	}
	/**
	* Sets the absolute position of the window.
	*
	* @param x - The desired horizontal absolute position of the window.
	* @param y - The desired vertical absolute position of the window.
	*/
	SetPosition(x, y) {
		return this[callerSym](SetPositionMethod, {
			x,
			y
		});
	}
	/**
	* Sets the window to be always on top.
	*
	* @param alwaysOnTop - Whether the window should stay on top.
	*/
	SetAlwaysOnTop(alwaysOnTop) {
		return this[callerSym](SetAlwaysOnTopMethod, { alwaysOnTop });
	}
	/**
	* Sets the background colour of the window.
	*
	* @param r - The desired red component of the window background.
	* @param g - The desired green component of the window background.
	* @param b - The desired blue component of the window background.
	* @param a - The desired alpha component of the window background.
	*/
	SetBackgroundColour(r, g, b, a) {
		return this[callerSym](SetBackgroundColourMethod, {
			r,
			g,
			b,
			a
		});
	}
	/**
	* Removes the window frame and title bar.
	*
	* @param frameless - Whether the window should be frameless.
	*/
	SetFrameless(frameless) {
		return this[callerSym](SetFramelessMethod, { frameless });
	}
	/**
	* Disables the system fullscreen button.
	*
	* @param enabled - Whether the fullscreen button should be enabled.
	*/
	SetFullscreenButtonEnabled(enabled) {
		return this[callerSym](SetFullscreenButtonEnabledMethod, { enabled });
	}
	/**
	* Sets the maximum size of the window.
	*
	* @param width - The desired maximum width of the window.
	* @param height - The desired maximum height of the window.
	*/
	SetMaxSize(width, height) {
		return this[callerSym](SetMaxSizeMethod, {
			width,
			height
		});
	}
	/**
	* Sets the minimum size of the window.
	*
	* @param width - The desired minimum width of the window.
	* @param height - The desired minimum height of the window.
	*/
	SetMinSize(width, height) {
		return this[callerSym](SetMinSizeMethod, {
			width,
			height
		});
	}
	/**
	* Sets the relative position of the window to the screen.
	*
	* @param x - The desired horizontal relative position of the window.
	* @param y - The desired vertical relative position of the window.
	*/
	SetRelativePosition(x, y) {
		return this[callerSym](SetRelativePositionMethod, {
			x,
			y
		});
	}
	/**
	* Sets whether the window is resizable.
	*
	* @param resizable - Whether the window should be resizable.
	*/
	SetResizable(resizable) {
		return this[callerSym](SetResizableMethod, { resizable });
	}
	/**
	* Sets the size of the window.
	*
	* @param width - The desired width of the window.
	* @param height - The desired height of the window.
	*/
	SetSize(width, height) {
		return this[callerSym](SetSizeMethod, {
			width,
			height
		});
	}
	/**
	* Sets the title of the window.
	*
	* @param title - The desired title of the window.
	*/
	SetTitle(title) {
		return this[callerSym](SetTitleMethod, { title });
	}
	/**
	* Sets the zoom level of the window.
	*
	* @param zoom - The desired zoom level.
	*/
	SetZoom(zoom) {
		return this[callerSym](SetZoomMethod, { zoom });
	}
	/**
	* Shows the window.
	*/
	Show() {
		return this[callerSym](ShowMethod);
	}
	/**
	* Returns the size of the window.
	*
	* @returns The current size of the window.
	*/
	Size() {
		return this[callerSym](SizeMethod);
	}
	/**
	* Toggles the window between fullscreen and normal.
	*/
	ToggleFullscreen() {
		return this[callerSym](ToggleFullscreenMethod);
	}
	/**
	* Toggles the window between maximised and normal.
	*/
	ToggleMaximise() {
		return this[callerSym](ToggleMaximiseMethod);
	}
	/**
	* Toggles the window between frameless and normal.
	*/
	ToggleFrameless() {
		return this[callerSym](ToggleFramelessMethod);
	}
	/**
	* Un-fullscreens the window.
	*/
	UnFullscreen() {
		return this[callerSym](UnFullscreenMethod);
	}
	/**
	* Un-maximises the window.
	*/
	UnMaximise() {
		return this[callerSym](UnMaximiseMethod);
	}
	/**
	* Un-minimises the window.
	*/
	UnMinimise() {
		return this[callerSym](UnMinimiseMethod);
	}
	/**
	* Returns the width of the window.
	*
	* @returns The current width of the window.
	*/
	Width() {
		return this[callerSym](WidthMethod);
	}
	/**
	* Zooms the window.
	*/
	Zoom() {
		return this[callerSym](ZoomMethod);
	}
	/**
	* Increases the zoom level of the webview content.
	*/
	ZoomIn() {
		return this[callerSym](ZoomInMethod);
	}
	/**
	* Decreases the zoom level of the webview content.
	*/
	ZoomOut() {
		return this[callerSym](ZoomOutMethod);
	}
	/**
	* Resets the zoom level of the webview content.
	*/
	ZoomReset() {
		return this[callerSym](ZoomResetMethod);
	}
	/**
	* Handles file drops originating from platform-specific code (e.g., macOS/Linux native drag-and-drop).
	* Gathers information about the drop target element and sends it back to the Go backend.
	*
	* @param filenames - An array of file paths (strings) that were dropped.
	* @param x - The x-coordinate of the drop event, in logical (CSS) pixels relative to the webview.
	* @param y - The y-coordinate of the drop event, in logical (CSS) pixels relative to the webview.
	*/
	HandlePlatformFileDrop(filenames, x, y) {
		if (window._wails?.flags?.enableFileDrop === false) return;
		const dropTarget = getDropTargetElement(document.elementFromPoint(x, y));
		if (!dropTarget) return;
		const elementDetails = {
			id: dropTarget.id,
			classList: Array.from(dropTarget.classList),
			attributes: {}
		};
		for (let i = 0; i < dropTarget.attributes.length; i++) {
			const attr = dropTarget.attributes[i];
			elementDetails.attributes[attr.name] = attr.value;
		}
		const payload = {
			filenames,
			x,
			y,
			elementDetails
		};
		this[callerSym](FilesDropped, payload);
		cleanupNativeDrag();
	}
	/**
	* Moves the window to the center of the specified screen's work area.
	*
	* @param screenID - The ID of the target screen.
	*/
	SetScreen(screenID) {
		return this[callerSym](SetScreenMethod, { screenID });
	}
	SnapAssist() {
		return this[callerSym](SnapAssistMethod);
	}
	/**
	* Opens the print dialog for the window.
	*/
	Print() {
		return this[callerSym](PrintMethod);
	}
}("");
/**
* Sets up global drag and drop event listeners for file drops.
* Handles visual feedback (hover state) and file drop processing.
*/
function setupDropTargetListeners() {
	const docElement = document.documentElement;
	let dragEnterCounter = 0;
	docElement.addEventListener("dragenter", (event) => {
		if (!event.dataTransfer?.types.includes("Files")) return;
		event.preventDefault();
		if (window._wails?.flags?.enableFileDrop === false) {
			event.dataTransfer.dropEffect = "none";
			return;
		}
		dragEnterCounter++;
		const dropTarget = getDropTargetElement(document.elementFromPoint(event.clientX, event.clientY));
		if (currentDropTarget && currentDropTarget !== dropTarget) currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
		if (dropTarget) {
			dropTarget.classList.add(DROP_TARGET_ACTIVE_CLASS);
			event.dataTransfer.dropEffect = "copy";
			currentDropTarget = dropTarget;
		} else {
			event.dataTransfer.dropEffect = "none";
			currentDropTarget = null;
		}
	}, false);
	docElement.addEventListener("dragover", (event) => {
		if (!event.dataTransfer?.types.includes("Files")) return;
		event.preventDefault();
		if (window._wails?.flags?.enableFileDrop === false) {
			event.dataTransfer.dropEffect = "none";
			return;
		}
		const dropTarget = getDropTargetElement(document.elementFromPoint(event.clientX, event.clientY));
		if (currentDropTarget && currentDropTarget !== dropTarget) currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
		if (dropTarget) {
			if (!dropTarget.classList.contains(DROP_TARGET_ACTIVE_CLASS)) dropTarget.classList.add(DROP_TARGET_ACTIVE_CLASS);
			event.dataTransfer.dropEffect = "copy";
			currentDropTarget = dropTarget;
		} else {
			event.dataTransfer.dropEffect = "none";
			currentDropTarget = null;
		}
	}, false);
	docElement.addEventListener("dragleave", (event) => {
		if (!event.dataTransfer?.types.includes("Files")) return;
		event.preventDefault();
		if (window._wails?.flags?.enableFileDrop === false) return;
		if (event.relatedTarget === null) return;
		dragEnterCounter--;
		if (dragEnterCounter === 0 || currentDropTarget && !currentDropTarget.contains(event.relatedTarget)) {
			if (currentDropTarget) {
				currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
				currentDropTarget = null;
			}
			dragEnterCounter = 0;
		}
	}, false);
	docElement.addEventListener("drop", (event) => {
		if (!event.dataTransfer?.types.includes("Files")) return;
		event.preventDefault();
		if (window._wails?.flags?.enableFileDrop === false) return;
		dragEnterCounter = 0;
		if (currentDropTarget) {
			currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
			currentDropTarget = null;
		}
		if (canResolveFilePaths()) {
			const files = [];
			if (event.dataTransfer.items) {
				for (const item of event.dataTransfer.items) if (item.kind === "file") {
					const file = item.getAsFile();
					if (file) files.push(file);
				}
			} else if (event.dataTransfer.files) for (const file of event.dataTransfer.files) files.push(file);
			if (files.length > 0) resolveFilePaths(event.clientX, event.clientY, files);
		}
	}, false);
}
if (typeof window !== "undefined" && typeof document !== "undefined") setupDropTargetListeners();
// ../../../internal/runtime/desktop/@wailsio/runtime/src/index.ts
if (hasDOM) window._wails = window._wails || {};
if (hasDOM) {
	window._wails.invoke = invoke;
	window._wails.clientId = clientId;
}
if (hasDOM) window._wails.handlePlatformFileDrop = thisWindow.HandlePlatformFileDrop.bind(thisWindow);
if (hasDOM) {
	window._wails.handleDragEnter = handleDragEnter;
	window._wails.handleDragLeave = handleDragLeave;
	window._wails.handleDragOver = handleDragOver;
}
if (hasDOM) invoke("wails:runtime:ready");
/**
* Loads a script from the given URL if it exists.
* Uses HEAD request to check existence, then injects a script tag.
* Silently ignores if the script doesn't exist.
*/
function loadOptionalScript(url) {
	return fetch(url, { method: "HEAD" }).then((response) => {
		if (response.ok) {
			if ((response.headers.get("content-type") || "").toLowerCase().includes("javascript")) {
				const script = document.createElement("script");
				script.src = url;
				document.head.appendChild(script);
			}
		}
	}).catch(() => {});
}
if (hasDOM) loadOptionalScript("/wails/custom.js");
// bindings/github.com/wailsapp/wails/v3/examples/file-association/greetservice.js
/**
* @param {string} name
* @returns {$CancellablePromise<string>}
*/
function Greet(name) {
	return ByID(1411160069, name);
}
// main.js
var resultElement = document.getElementById("result");
var timeElement = document.getElementById("time");
window.doGreet = () => {
	let name = document.getElementById("name").value;
	if (!name) name = "anonymous";
	Greet(name).then((result) => {
		resultElement.innerText = result;
	}).catch((err) => {
		console.log(err);
	});
};
On("time", (time) => {
	timeElement.innerText = time.data;
});
